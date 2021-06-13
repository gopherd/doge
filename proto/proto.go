package proto

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"

	"google.golang.org/protobuf/proto"
)

var (
	ErrIntOverflow     = errors.New("proto: integer overflow")
	ErrUnrecognizeType = errors.New("proto: unrecognize type")
)

// Message represents a message interface
type Message interface {
	proto.Message

	// Type returns the message type
	Type() uint32
}

var (
	creators = make(map[uint32]func() Message)
	modules  = make(map[string][]uint32)
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
func Register(module string, typ uint32, creator func() Message) {
	if creator == nil {
		panic(fmt.Sprintf("proto: Register creator is nil for type", typ))
	}
	if _, dup := creators[typ]; dup {
		panic(fmt.Sprintf("proto: Register called twice for type %d", typ))
	}
	creators[typ] = creator
	modules[module] = append(modules[module], typ)
}

// New creates a message by type, nil returned if type not found
func New(typ uint32) Message {
	if creator, ok := creators[typ]; ok {
		return creator()
	}
	return nil
}

// Lookup lookups all registered types by module
func Lookup(module string) []uint32 {
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

// SizeofUvarint returns the number of unsigned-varint encoding-bytes.
func SizeofUvarint(x uint64) int {
	i := 0
	for x >= 0x80 {
		x >>= 7
		i++
	}
	return i + 1
}

// SizeofVarint returns the number of varint encoding-bytes.
func SizeofVarint(x int64) int {
	ux := uint64(x) << 1
	if x < 0 {
		ux = ^ux
	}
	return SizeofUvarint(ux)
}

// Marshal returns the wire-format encoding of m with length and type.
//
//           |<-- body -->|
//	|body.len|type|content|
//
func Marshal(m Message, reservedHeadLen int) ([]byte, error) {
	size := Sizeof(m)
	typeSize := SizeofUvarint(uint64(m.Type()))
	n := size
	n += SizeofUvarint(uint64(size + typeSize))
	n += typeSize
	n += reservedHeadLen
	buf := make([]byte, n-size, n)
	off := reservedHeadLen
	off += binary.PutUvarint(buf[off:], uint64(size+typeSize))
	off += binary.PutUvarint(buf[off:], uint64(m.Type()))
	options := proto.MarshalOptions{
		UseCachedSize: true,
	}
	_, err := options.MarshalAppend(buf[off:off], m)
	return buf, err
}

// UnmarshalBody parses the type, then parses wire-format message in remain bytes
// and places the result in m that looked up by type.
func UnmarshalBody(b []byte) (Message, error) {
	typ, n := binary.Uvarint(b)
	if n == 0 {
		return nil, io.ErrShortBuffer
	} else if n < 0 || typ > math.MaxUint32 {
		return nil, ErrIntOverflow
	}
	m := New(uint32(typ))
	if m == nil {
		return nil, ErrUnrecognizeType
	}
	return m, proto.Unmarshal(b[n:], m)
}
