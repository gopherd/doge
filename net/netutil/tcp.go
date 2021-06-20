package netutil

import (
	"crypto/tls"
	"net"
	"time"

	"github.com/mkideal/log"
)

// TCPKeepAliveListener wraps TCPListener with a keepalive duration
type TCPKeepAliveListener struct {
	*net.TCPListener
	duration time.Duration
}

// NewTCPKeepAliveListener creates a TCPKeepAliveListener
func NewTCPKeepAliveListener(ln *net.TCPListener, d time.Duration) *TCPKeepAliveListener {
	return &TCPKeepAliveListener{
		TCPListener: ln,
		duration:    d,
	}
}

// Accept implements net.Listener Accept method
func (ln TCPKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	if ln.duration == 0 {
		ln.duration = 3 * time.Minute
	}
	tc.SetKeepAlivePeriod(ln.duration)
	return tc, nil
}

// KeepAliveTCPConn sets conn's keepalive duration
func KeepAliveTCPConn(conn net.Conn, d time.Duration) {
	tc, ok := conn.(*net.TCPConn)
	if ok {
		tc.SetKeepAlive(true)
		tc.SetKeepAlivePeriod(d)
	}
}

// ConnHandler handles net.Conn
type ConnHandler func(ip string, conn net.Conn)

// TCPServer represents a tcp server
type TCPServer struct {
	addr     string
	handler  ConnHandler
	listener net.Listener
}

// NewTCPServer creates a tcp server
func NewTCPServer(addr string, handler ConnHandler) *TCPServer {
	server := new(TCPServer)
	server.addr = addr
	server.handler = handler
	return server
}

// ListenAndServe starts the tcp server
func (server *TCPServer) ListenAndServe(async bool, keepalive time.Duration, certs ...tls.Certificate) error {
	var (
		listener net.Listener
		err      error
	)
	if len(certs) > 0 {
		config := &tls.Config{Certificates: certs}
		listener, err = tls.Listen("tcp", server.addr, config)
	} else {
		listener, err = listenTCP(server.addr)
	}
	if err != nil {
		return err
	}
	if keepalive > 0 {
		if l, ok := listener.(*net.TCPListener); ok {
			listener = NewTCPKeepAliveListener(l, keepalive)
		} else {
			log.Warn().Print("TCPServer.ListenAndServe: keepalive is not supported")
		}
	}
	server.listener = listener
	return serve(server.listener, server.handler, async)
}

// Shutdown shutdowns the tcp server
func (server *TCPServer) Shutdown() error {
	return server.listener.Close()
}

func listenTCP(addr string) (*net.TCPListener, error) {
	a, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	return net.ListenTCP("tcp", a)
}

func serve(listener net.Listener, handler ConnHandler, async bool) error {
	serveFunc := func() error {
		var tempDelay time.Duration // how long to sleep on accept failure
		for {
			conn, err := listener.Accept()
			if err != nil {
				if ne, ok := err.(net.Error); ok && ne.Temporary() {
					if tempDelay == 0 {
						tempDelay = 5 * time.Millisecond
					} else {
						tempDelay *= 2
					}
					if max := 1 * time.Second; tempDelay > max {
						tempDelay = max
					}
					log.Info().
						Error("error", err).
						String("delay", tempDelay.String()).
						Print("accept connection error, retrying")
					time.Sleep(tempDelay)
					continue
				}
				return err
			}
			tempDelay = 0
			var ip string
			if addr, ok := conn.RemoteAddr().(*net.TCPAddr); ok {
				ip = addr.IP.String()
			}
			go handler(ip, conn)
		}
		return nil
	}
	if async {
		go serveFunc()
	} else {
		return serveFunc()
	}
	return nil
}

// ListenAndServeTCP wraps TCPServer's ListenAndServeTCP
func ListenAndServeTCP(addr string, keepalive time.Duration, handler ConnHandler, async bool, certs ...tls.Certificate) error {
	server := NewTCPServer(addr, handler)
	return server.ListenAndServe(async, keepalive, certs...)
}
