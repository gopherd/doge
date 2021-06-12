package proto

import (
	"fmt"
)

// Message represents a message interface
type Message interface {
	Type() int32
	Size() int
}

var (
	creators = make(map[int32]func() Message)
	modules  = make(map[string][]int32)
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
//		proto.Register(BarType, "foo", func() proto.Message { return new(Bar) })
//	}
func Register(typ int32, module string, creator func() Message) {
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
func New(typ int32) Message {
	if creator, ok := creators[typ]; ok {
		return creator()
	}
	return nil
}

// Lookup lookups all registered types by module
func Lookup(module string) []int32 {
	return modules[module]
}
