package netutil

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"os"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/gopherd/doge/io/pagebuf"
	"github.com/mkideal/log"
)

const (
	// max length of content
	MaxContentLength = int(^uint(0) >> 1)
)

var (
	// content length greater than MaxContentLength
	ErrContentLengthOverflow = errors.New("content length overflow")
)

// errno returns v's underlying uintptr, else 0.
//
// TODO: See comment in isClosedConnError.
func errno(v error) uintptr {
	if rv := reflect.ValueOf(v); rv.Kind() == reflect.Uintptr {
		return uintptr(rv.Uint())
	}
	return 0
}

// isClosedConnError reports whether err is an error from use of a closed
// network connection.
func isClosedConnError(err error) bool {
	if err == nil {
		return false
	}

	// TODO: remove this string search and be more like the Windows
	// case below. That might involve modifying the standard library
	// to return better error types.
	str := err.Error()
	if strings.Contains(str, "use of closed network connection") {
		return true
	}

	// TODO(bradfitz): x/tools/cmd/bundle doesn't really support
	// build tags, so I can't make an http2_windows.go file with
	// Windows-specific stuff. Fix that and move this, once we
	// have a way to bundle this into std's net/http somehow.
	if runtime.GOOS == "windows" {
		if oe, ok := err.(*net.OpError); ok && oe.Op == "read" {
			if se, ok := oe.Err.(*os.SyscallError); ok && se.Syscall == "wsarecv" {
				const WSAECONNABORTED = 10053
				const WSAECONNRESET = 10054
				if n := errno(se.Err); n == WSAECONNRESET || n == WSAECONNABORTED {
					return true
				}
			}
		}
	}
	return false
}

// timeoutReader wraps net.Conn as an io.Reader with timeout
type timeoutReader struct {
	conn    net.Conn
	timeout time.Duration
}

// Read implements io.Reader Read method
func (tr *timeoutReader) Read(p []byte) (n int, err error) {
	if tr.timeout > 0 {
		tr.conn.SetReadDeadline(time.Now().Add(tr.timeout))
	}
	return tr.conn.Read(p)
}

// Body represents a message body
type Body interface {
	io.ByteReader
	io.Reader

	// Len returns remain length of body
	Len() int

	// Peek returns the next n bytes without advancing the reader. The bytes stop
	// being valid at the next read call. If Peek returns fewer than n bytes, it
	// also returns an error explaining why the read is short. The error is
	// ErrBufferFull if n is larger than b's buffer size.
	Peek(n int) ([]byte, error)
}

type reader struct {
	conn  net.Conn
	bufr  *bufio.Reader
	limit int
}

func newReader(conn net.Conn, timeout time.Duration) *reader {
	b := &reader{
		conn: conn,
	}
	b.bufr = bufio.NewReader(&timeoutReader{
		conn:    conn,
		timeout: timeout,
	})
	return b
}

// Len implements Body Len method
func (b *reader) Len() int {
	return b.limit
}

// Peek implements Body Peek method
func (b *reader) Peek(n int) ([]byte, error) {
	if b.limit >= 0 && b.limit < n {
		return nil, io.EOF
	}
	return b.bufr.Peek(n)
}

// ReadByte implements io.ByteReader ReadByte method
func (b *reader) ReadByte() (c byte, err error) {
	if b.limit == 0 {
		err = io.ErrUnexpectedEOF
		return
	}
	c, err = b.bufr.ReadByte()
	if err == nil {
		if b.limit > 0 {
			b.limit--
		}
	}
	return
}

// Read implements io.Reader Read method
func (b *reader) Read(p []byte) (n int, err error) {
	if b.limit == 0 {
		return 0, io.EOF
	}
	if b.limit > 0 && len(p) > b.limit {
		p = p[:b.limit]
	}
	n, err = b.bufr.Read(p)
	if b.limit > 0 {
		b.limit -= n
	}
	return
}

func (b *reader) discard() error {
	if b.limit == 0 {
		return nil
	}
	_, err := b.bufr.Discard(b.limit)
	return err
}

// Option represents options of NewSession
type Option func(*option)

type option struct {
	timeout time.Duration
}

func defaultOption() option {
	return option{}
}

// WithTimeout specify read timeout of session
func WithTimeout(timeout time.Duration) Option {
	return func(opt *option) {
		opt.timeout = timeout
	}
}

// SessionEventHandler handles session events
type SessionEventHandler interface {
	OnReady()                  // ready to read/write
	OnClose(err error)         // session closed, err maybe nil
	OnMessage(body Body) error // received a message
}

// Session wraps network session
type Session struct {
	closed int32
	err    error
	once   sync.Once

	reader  *reader
	writer  *bufio.Writer
	handler SessionEventHandler

	mutex sync.Mutex
	cond  *sync.Cond
	pipe  *pagebuf.PageBuffer
}

// NewSession creates a session
func NewSession(conn net.Conn, handler SessionEventHandler, options ...Option) *Session {
	var opt = defaultOption()
	for i := range options {
		options[i](&opt)
	}
	s := &Session{
		reader:  newReader(conn, opt.timeout),
		writer:  bufio.NewWriter(conn),
		handler: handler,
		pipe:    pagebuf.NewPageBuffer(),
	}
	s.cond = sync.NewCond(&s.mutex)
	return s
}

// Write implements io.Writer Write method
func (s *Session) Write(p []byte) (n int, err error) {
	if !s.IsClosed() {
		err = net.ErrClosed
		return
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	n, err = s.pipe.Write(p)
	if err != nil {
		return
	}
	if s.pipe.Len() == n {
		s.cond.Signal()
	}
	return
}

// Run runs the read and write loops, and block until the session closed
func (s *Session) Run() {
	s.once.Do(s.run)
}

func (s *Session) run() {
	var (
		readyWg sync.WaitGroup
		closeWg sync.WaitGroup
	)
	readyWg.Add(2)
	closeWg.Add(2)
	go s.readLoop(&readyWg, &closeWg)
	go s.writeLoop(&readyWg, &closeWg)

	readyWg.Wait()
	s.handler.OnReady()

	closeWg.Wait()
	s.handler.OnClose(s.err)

	// close the underlying connection
	log.Debug("session %s closed", s.reader.conn.RemoteAddr().String())
	s.reader.conn.Close()
}

// IsClosed returns whether the session is closed
func (s *Session) IsClosed() bool {
	return atomic.LoadInt32(&s.closed) == 1
}

func (s *Session) setClosed(err error) {
	if err != nil {
		s.err = err
	}
	atomic.StoreInt32(&s.closed, 1)
}

// Close closes the session
func (s *Session) Close() {
	s.setClosed(nil)
}

func (s *Session) readLoop(readyWg, closeWg *sync.WaitGroup) {
	readyWg.Done()
	for !s.IsClosed() {
		if err := s.underlyingRead(); err != nil {
			terr, ok := err.(net.Error)
			if err != io.EOF && err != io.ErrUnexpectedEOF && (!ok || !terr.Temporary()) && !isClosedConnError(err) {
				log.Warn("session %s read error: %v", s.reader.conn.RemoteAddr().String(), err)
			} else if operr, ok := err.(*net.OpError); ok && operr.Err.Error() == syscall.ECONNRESET.Error() {
				log.Info("session %s read failed: %v", s.reader.conn.RemoteAddr().String(), err)
			}
			s.setClosed(err)
			break
		}
	}
	closeWg.Done()
}

func (s *Session) writeLoop(readyWg, closeWg *sync.WaitGroup) {
	readyWg.Done()
	var buf = make([]byte, s.pipe.PageSize())
	for !s.IsClosed() {
		s.cond.L.Lock()
		for s.pipe.Len() == 0 {
			s.cond.Wait()
			if s.IsClosed() {
				break
			}
		}
		s.cond.L.Unlock()
		for {
			s.cond.L.Lock()
			n, _ := s.pipe.Read(buf)
			s.cond.L.Unlock()
			if n == 0 {
				break
			}
			s.underlyingWrite(buf[:n])
		}
	}
	closeWg.Done()
}

func (s *Session) underlyingWrite(p []byte) error {
	_, err := s.writer.Write(p)
	if err != nil {
		s.setClosed(err)
	}
	return err
}

func (s *Session) underlyingRead() error {
	// read content length without limit
	s.reader.limit = -1
	contentLength, err := binary.ReadUvarint(s.reader)
	if err != nil {
		return err
	}
	if contentLength > uint64(MaxContentLength) {
		return ErrContentLengthOverflow
	}

	// handle the body with limit
	s.reader.limit = int(contentLength)
	if err := s.handler.OnMessage(s.reader); err != nil {
		return err
	}
	return s.reader.discard()
}
