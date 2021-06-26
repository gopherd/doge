package proto

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"google.golang.org/protobuf/proto"
)

const (
	// max size of content: 1G
	MaxSize = 1 << 30
	// max message type
	MaxType = 1 << 31
)

var (
	ErrVarintOverflow   = errors.New("proto: varint overflow")
	ErrSizeOverflow     = errors.New("proto: size overflow")
	ErrTypeOverflow     = errors.New("proto: type overflow")
	ErrUnrecognizedType = errors.New("proto: unrecognized type")
)

// Type represents message type
type Type = uint32

// Body represents message body
type Body interface {
	io.ByteReader
	io.Reader

	// Len returns remain length of body
	Len() int

	// Peek returns the next n bytes without advancing the reader. The bytes stop
	// being valid at the next read call. If Peek returns fewer than n bytes, it
	// also returns an error explaining why the read is short.
	Peek(n int) ([]byte, error)

	// Discard skips the next n bytes, returning the number of bytes discarded.
	// If Discard skips fewer than n bytes, it also returns an error.
	Discard(n int) (discarded int, err error)
}

// Message represents a message interface
type Message interface {
	proto.Message

	// Type returns the message type
	Type() Type
}

var (
	creators = make(map[Type]func() Message)
	modules  = make(map[string][]Type)
)

// Register registers a message creator by type. Register is not
// concurrent-safe, it is recommended to call in `init` function.
//
// e.g.
//
//	package foo
//
//	import "github.com/gopherd/doge/proto"
//
//	func init() {
//		proto.Register("foo", BarType, func() proto.Message { return new(Bar) })
//	}
func Register(module string, typ Type, creator func() Message) {
	if typ > MaxType {
		panic(fmt.Sprintf("proto: Register type %d out of range [0, %d]", typ, MaxType))
	}
	if creator == nil {
		panic(fmt.Sprintf("proto: Register creator is nil for type %d", typ))
	}
	if _, dup := creators[typ]; dup {
		panic(fmt.Sprintf("proto: Register called twice for type %d", typ))
	}
	creators[typ] = creator
	modules[module] = append(modules[module], typ)
}

// New creates a message by type, nil returned if type not found
func New(typ Type) Message {
	if creator, ok := creators[typ]; ok {
		return creator()
	}
	return nil
}

// Lookup lookups all registered types by module
func Lookup(module string) []Type {
	return modules[module]
}

// Sizeof calculates the exact size of marshaled bytes
func Sizeof(m Message) int {
	return proto.Size(m)
}

// Nameof returns the name of message
func Nameof(m Message) string {
	return string(proto.MessageName(m))
}

// Peeker peeks n bytes
type Peeker interface {
	Peek(n int) ([]byte, error)
}

func peekUvarint(peeker Peeker) (int, uint64, error) {
	var x uint64
	var s uint
	var n int
	for i := 0; i < binary.MaxVarintLen64; i++ {
		n++
		buf, err := peeker.Peek(n)
		if err != nil {
			return n, x, err
		}
		b := buf[i]
		if b < 0x80 {
			if i == binary.MaxVarintLen64-1 && b > 1 {
				return n, x, ErrTypeOverflow
			}
			return n, x | uint64(b)<<s, nil
		}
		x |= uint64(b&0x7f) << s
		s += 7
	}
	return n, x, ErrTypeOverflow
}

func convertType(x uint64, err error) (Type, error) {
	if err != nil {
		return 0, err
	}
	if x > MaxType {
		return 0, ErrTypeOverflow
	}
	return Type(x), nil
}

// ReadType reads message type from reader
func ReadType(r io.ByteReader) (typ Type, err error) {
	return convertType(binary.ReadUvarint(r))
}

// PeekType reads message type without advancing underlying reader offset
func PeekType(peeker Peeker) (n int, typ Type, err error) {
	var x uint64
	n, x, err = peekUvarint(peeker)
	typ, err = convertType(x, err)
	return
}

// EncodeType encodes type as varint to buf and returns number of bytes written.
func EncodeType(buf []byte, typ Type) int {
	return binary.PutUvarint(buf, uint64(typ))
}

func convertSize(x uint64, err error) (int, error) {
	if err != nil {
		return 0, err
	}
	if x > MaxSize {
		return 0, ErrSizeOverflow
	}
	return int(x), nil
}

// ReadSize reads message size from reader
func ReadSize(r io.ByteReader) (size int, err error) {
	return convertSize(binary.ReadUvarint(r))
}

// PeekSize reads message size without advancing underlying reader offset
func PeekSize(peeker Peeker) (n int, size int, err error) {
	var x uint64
	n, x, err = peekUvarint(peeker)
	size, err = convertSize(x, err)
	return
}

// EncodeSize encodes type as varint to buf and returns number of bytes written.
func EncodeSize(buf []byte, size int) int {
	return binary.PutUvarint(buf, uint64(size))
}

// sizeofUvarint returns the number of unsigned-varint encoding-bytes.
func sizeofUvarint(x uint64) int {
	i := 0
	for x >= 0x80 {
		x >>= 7
		i++
	}
	return i + 1
}

func sizeof(m Message) (size, ssize, tsize int) {
	msize := Sizeof(m)
	tsize = sizeofUvarint(uint64(m.Type()))
	size = msize + tsize
	ssize = sizeofUvarint(uint64(size))
	return
}

// Encode returns the wire-format encoding of m with size and type.
//
//           |<-- body -->|
//	|body.len|type|content|
func Encode(m Message, reservedHeadLen int) ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	size, ssize, _ := sizeof(m)
	if size > MaxSize {
		return nil, ErrSizeOverflow
	}
	off := reservedHeadLen
	buf := make([]byte, off+ssize+size)
	_, err := encodeAppend(buf[off:], m, size)
	return buf, err
}

func EncodeAppend(buf []byte, m Message) ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	size, ssize, tsize := sizeof(m)
	if size > MaxSize {
		return nil, ErrSizeOverflow
	}
	off := len(buf)
	n := ssize + size
	if cap(buf)-off < n {
		newbuf := make([]byte, off+ssize+tsize, off+n)
		copy(newbuf, buf)
		buf = newbuf
	} else {
		buf = buf[:off+ssize+tsize]
	}
	return encodeAppend(buf[off:], m, size)
}

func encodeAppend(buf []byte, m Message, size int) ([]byte, error) {
	off := 0
	off += binary.PutUvarint(buf[off:], uint64(size))
	off += binary.PutUvarint(buf[off:], uint64(m.Type()))
	options := proto.MarshalOptions{
		UseCachedSize: true,
	}
	_, err := options.MarshalAppend(buf[off:off], m)
	return buf, err
}

// Marshal returns the wire-format encoding of m without size or type.
func Marshal(m Message) ([]byte, error) {
	return proto.Marshal(m)
}

// Decode decodes one message with size and type from buf and
// returns number of bytes read and unmarshaled message.
func Decode(buf []byte) (int, Message, error) {
	size, n := binary.Uvarint(buf)
	if n == 0 {
		return 0, nil, io.ErrShortBuffer
	} else if n < 0 {
		return -n, nil, ErrSizeOverflow
	} else if size > MaxSize {
		return n, nil, ErrSizeOverflow
	}
	m, err := DecodeBody(buf[n:])
	return n + int(size), m, err
}

// Decode decodes one message that contains type from buf and returns
// unmarshaled message.
func DecodeBody(b []byte) (Message, error) {
	typ, n := binary.Uvarint(b)
	if n == 0 {
		return nil, io.ErrShortBuffer
	} else if n < 0 || typ > MaxType {
		return nil, ErrTypeOverflow
	}
	m := New(Type(typ))
	if m == nil {
		return nil, ErrUnrecognizedType
	}
	err := Unmarshal(b[n:], m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// Unmarshal parses the wire-format message in b and places the result in m.
// The provided message must be mutable (e.g., a non-nil pointer to a message).
func Unmarshal(b []byte, m Message) error {
	return proto.Unmarshal(b, m)
}
