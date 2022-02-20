package container

// Node represents a generic printable node
type Node[T comparable] interface {
	String() string          // String returns node self information
	Parent() T               // Parent returns parent node or nil
	NumChild() int           // NumChild returns number of child
	GetChildByIndex(i int) T // GetChildByIndex gets child by index
}
