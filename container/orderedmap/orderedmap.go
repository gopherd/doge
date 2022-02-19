package orderedmap

import (
	"bytes"
	"constraints"
	"fmt"

	"github.com/gopherd/doge/container"
	"github.com/gopherd/doge/operator"
)

type (
	// Iterator represents an iterator of OrderedMap to iterate nodes
	Iterator[K comparable, V any] interface {
		// Prev returns previous iterator
		Prev() Iterator[K, V]
		// Next returns next node iterator
		Next() Iterator[K, V]
		// Key returns key of the node
		Key() K
		// Value returns value of the node
		Value() V
		// SetValue sets value of the node
		SetValue(V)

		underlyingNode() *node[K, V]
	}

	// LessFunc represents a comparation function which reports whether x "less" than y
	LessFunc[T any] func(x, y T) bool
)

// OrderedMap represents an ordered map
type OrderedMap[K comparable, V any] struct {
	size int
	root *node[K, V]
	less LessFunc[K]
}

// New creates a OrderedMap for ordered K
func New[K constraints.Ordered, V any]() *OrderedMap[K, V] {
	return NewFunc[K, V](operator.Less[K])
}

// NewFunc creates a OrderedMap with compare function
func NewFunc[K comparable, V any](less LessFunc[K]) *OrderedMap[K, V] {
	if less == nil {
		panic("orderedmap: less function is nil")
	}
	return &OrderedMap[K, V]{
		less: less,
	}
}

// Len returns the number of elements
func (m OrderedMap[K, V]) Len() int {
	return m.size
}

// Clear clears the container
func (m *OrderedMap[K, V]) Clear() {
	m.root = nil
	m.size = 0
}

// Find finds node by key, nil returned if the key not found.
func (m *OrderedMap[K, V]) Find(key K) Iterator[K, V] {
	return m.find(key)
}

// Insert inserts a key-value pair, inserted node and true returned
// if the key not found, otherwise, existed node and false returned.
func (m *OrderedMap[K, V]) Insert(key K, value V) (Iterator[K, V], bool) {
	node, ok := m.insert(key, value)
	if ok {
		m.size++
	}
	return node, ok
}

// Remove removes an element by key, false returned if the key not found.
func (m *OrderedMap[K, V]) Remove(key K) bool {
	node := m.find(key)
	if node == nil || node.null() {
		return false
	}
	m.remove(node, true)
	m.size--
	return true
}

// Erase deletes the node, false returned if the node not found.
func (m *OrderedMap[K, V]) Erase(iter Iterator[K, V]) bool {
	if iter == nil {
		return false
	}
	node := iter.underlyingNode()
	if node == nil || node.null() {
		return false
	}
	ok := m.remove(node, false)
	if ok {
		m.size--
	}
	return ok
}

// First returns the first node.
//
// usage:
//
//	iter := m.First()
//	for iter != nil {
//		// hint: do something here using iter
//		// hint: iter.Key(), iter.Value(), iter.SetValue(newValue)
//		iter = iter.Next()
//	}
func (m *OrderedMap[K, V]) First() Iterator[K, V] {
	if m.root == nil {
		return nil
	}
	return m.root.smallest()
}

// Last returns the first node.
//
// usage:
//
//	iter := m.Last()
//	for iter != nil {
//		// hint: do something here using iter
//		// hint: iter.Key(), iter.Value(), iter.SetValue(newValue)
//		iter = iter.Prev()
//	}
func (m *OrderedMap[K, V]) Last() Iterator[K, V] {
	if m.root == nil {
		return nil
	}
	return m.root.biggest()
}

// Print pretty prints the map by specified options
func (m *OrderedMap[K, V]) Print(options container.PrintOptions) string {
	return container.PrintNode[*node[K, V]](m.root, options)
}

// String returns content of the map as a plain string
func (m *OrderedMap[K, V]) String() string {
	var buf bytes.Buffer
	buf.WriteByte('[')
	iter := m.First()
	for iter != nil {
		fmt.Fprintf(&buf, "%v:%v", iter.Key(), iter.Value())
		iter = iter.Next()
		if iter != nil {
			buf.WriteByte(' ')
		}
	}
	buf.WriteByte(']')
	return buf.String()
}

func (m *OrderedMap[K, V]) insert(key K, value V) (*node[K, V], bool) {
	if m.root == nil {
		m.root = &node[K, V]{
			color: black,
			key:   key,
			value: value,
		}
		m.root.left = makenull(m.root)
		m.root.right = makenull(m.root)
		return m.root, true
	}

	var (
		next     = m.root
		inserted *node[K, V]
	)
	for {
		if key == next.key {
			next.value = value
			return next, false
		}
		if m.less(key, next.key) {
			if next.left.null() {
				inserted = &node[K, V]{
					parent: next,
					color:  red,
					key:    key,
					value:  value,
				}
				inserted.left = makenull(inserted)
				inserted.right = makenull(inserted)
				next.left = inserted
				break
			} else {
				next = next.left
			}
		} else {
			if next.right.null() {
				inserted = &node[K, V]{
					parent: next,
					color:  red,
					key:    key,
					value:  value,
				}
				inserted.left = makenull(inserted)
				inserted.right = makenull(inserted)
				next.right = inserted
				break
			} else {
				next = next.right
			}
		}
	}

	next = inserted
	for {
		next = m.doInsert(next)
		if next == nil {
			break
		}
	}
	return inserted, true
}

func (m *OrderedMap[K, V]) find(key K) *node[K, V] {
	var next = m.root
	for next != nil && !next.null() {
		if next.key == key {
			return next
		}
		if m.less(key, next.key) {
			next = next.left
		} else {
			next = next.right
		}
	}
	return nil
}

func (m *OrderedMap[K, V]) remove(n *node[K, V], must bool) bool {
	if !must {
		if m.root == nil || n == nil || n.ancestor() != m.root {
			return false
		}
	}
	if !n.right.null() {
		smallest := n.right.smallest()
		n.value, smallest.value = smallest.value, n.value
		n.key, smallest.key = smallest.key, n.key
		n = smallest
	}
	var child = n.left
	if child.null() {
		child = n.right
	}
	if n.parent == nil {
		if n.left.null() && n.right.null() {
			m.root = nil
			return true
		}
		child.parent = nil
		m.root = child
		m.root.color = black
		return true
	}

	if n.parent.left == n {
		n.parent.left = child
	} else {
		n.parent.right = child
	}
	child.parent = n.parent
	if n.color == red {
		return true
	}
	if child.color == red {
		child.color = black
		return true
	}
	for child != nil {
		child = m.doRemove(child)
	}
	return true
}

func (m *OrderedMap[K, V]) doInsert(n *node[K, V]) *node[K, V] {
	if n.parent == nil {
		m.root = n
		n.color = black
		return nil
	}
	if n.parent.color == black {
		return nil
	}
	uncle := n.uncle()
	if uncle.color == red {
		n.parent.color = black
		uncle.color = black
		gp := n.grandparent()
		gp.color = red
		return gp
	}
	if n.parent.right == n && n.grandparent().left == n.parent {
		m.rotateLeft(n.parent)
		n.color = black
		n.parent.color = red
		m.rotateRight(n.parent)
	} else if n.parent.left == n && n.grandparent().right == n.parent {
		m.rotateRight(n.parent)
		n.color = black
		n.parent.color = red
		m.rotateLeft(n.parent)
	} else if n.parent.left == n && n.grandparent().left == n.parent {
		n.parent.color = black
		n.grandparent().color = red
		m.rotateRight(n.grandparent())
	} else if n.parent.right == n && n.grandparent().right == n.parent {
		n.parent.color = black
		n.grandparent().color = red
		m.rotateLeft(n.grandparent())
	}
	return nil
}

func (m *OrderedMap[K, V]) doRemove(n *node[K, V]) *node[K, V] {
	if n.parent == nil {
		n.color = black
		return nil
	}
	sibling := n.sibling()
	if sibling.color == red {
		n.parent.color = red
		sibling.color = black
		if n == n.parent.left {
			m.rotateLeft(n.parent)
		} else {
			m.rotateRight(n.parent)
		}
	}
	sibling = n.sibling()
	if n.parent.color == black &&
		sibling.color == black &&
		sibling.left.color == black &&
		sibling.right.color == black {
		sibling.color = red
		return n.parent
	}
	if n.parent.color == red &&
		sibling.color == black &&
		sibling.left.color == black &&
		sibling.right.color == black {
		sibling.color = red
		n.parent.color = black
		return nil
	}
	if sibling.color == black {
		if n == n.parent.left &&
			sibling.left.color == red &&
			sibling.right.color == black {
			sibling.color = red
			sibling.left.color = black
			m.rotateRight(sibling.left.parent)
		} else if n == n.parent.right &&
			sibling.left.color == black &&
			sibling.right.color == red {
			sibling.color = red
			sibling.right.color = black
			m.rotateLeft(sibling.right.parent)
		}
	}
	sibling = n.sibling()
	sibling.color = n.parent.color
	n.parent.color = black
	if n == n.parent.left {
		sibling.right.color = black
		m.rotateLeft(sibling.parent)
	} else {
		sibling.left.color = black
		m.rotateRight(sibling.parent)
	}
	return nil
}

const (
	left  = 0
	right = 1
)

func (m *OrderedMap[K, V]) rotate(p *node[K, V], dir int) *node[K, V] {
	var (
		g = p.parent
		s = p.child(1 - dir)
		c = s.child(dir)
	)
	p.setChild(1-dir, c)
	if !c.null() {
		c.parent = p
	}
	s.setChild(dir, p)
	p.parent = s
	s.parent = g
	if g != nil {
		if p == g.right {
			g.right = s
		} else {
			g.left = s
		}
	} else {
		m.root = s
	}
	return s
}

func (m *OrderedMap[K, V]) rotateLeft(p *node[K, V]) {
	m.rotate(p, left)
}

func (m *OrderedMap[K, V]) rotateRight(p *node[K, V]) {
	m.rotate(p, right)
}

type color byte

const (
	red color = iota
	black
)

// node represents the node of Map
type node[K comparable, V any] struct {
	parent      *node[K, V]
	left, right *node[K, V]
	color       color
	key         K
	value       V
}

func makenull[K comparable, V any](parent *node[K, V]) *node[K, V] {
	return &node[K, V]{
		parent: parent,
		color:  black,
	}
}

// ToString implements container.Node ToString method
func (node *node[K, V]) ToString() string {
	if node == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%v:%v", node.key, node.value)
}

// Parent implements container.Node Parent method
func (node *node[K, V]) Parent() *node[K, V] {
	if node == nil {
		return nil
	}
	return node.parent
}

// NumChild implements container.Node NumChild method
func (node *node[K, V]) NumChild() int {
	if node == nil {
		return 0
	}
	return operator.Bool[int](node.left != nil && !node.left.null()) +
		operator.Bool[int](node.right != nil && !node.right.null())
}

// GetChildByIndex implements container.Node GetChildByIndex method
func (node *node[K, V]) GetChildByIndex(i int) *node[K, V] {
	switch i {
	case 0:
		return operator.If(node.left != nil && !node.left.null(), node.left, node.right)
	case 1:
		return node.right
	default:
		panic("unreachable")
	}
}

// Prev implements Iterator Prev method
func (node *node[K, V]) Prev() Iterator[K, V] {
	if prev := node.prev(); prev != nil {
		return prev
	}
	return nil
}

// Next implements Iterator Next method
func (node *node[K, V]) Next() Iterator[K, V] {
	if next := node.next(); next != nil {
		return next
	}
	return nil
}

// Key returns node's key, implements Iterator Key method
func (node *node[K, V]) Key() K { return node.key }

// Value returns node's value, implements Iterator Value method
func (node *node[K, V]) Value() V { return node.value }

// SetValue sets node's value, implements Iterator SetValue method
func (node *node[K, V]) SetValue(value V) { node.value = value }

// underlyingNode implements Iterator underlyingNode method
func (node *node[K, V]) underlyingNode() *node[K, V] { return node }

func (node *node[K, V]) prev() *node[K, V] {
	if node == nil || node.null() {
		return nil
	}
	if !node.left.null() {
		return node.left.biggest()
	}
	parent := node.parent
	for node != parent.right {
		node = parent
		parent = node.parent
		if parent == nil {
			return nil
		}
	}
	return parent
}

func (node *node[K, V]) next() *node[K, V] {
	if node == nil || node.null() {
		return nil
	}
	if !node.right.null() {
		return node.right.smallest()
	}
	parent := node.parent
	for node != parent.left {
		node = parent
		parent = node.parent
		if parent == nil {
			return nil
		}
	}
	return parent
}

func (node *node[K, V]) null() bool {
	return node.left == nil && node.right == nil
}

func (node *node[K, V]) child(dir int) *node[K, V] {
	if dir == left {
		return node.left
	}
	return node.right
}

func (node *node[K, V]) setChild(dir int, child *node[K, V]) {
	if dir == left {
		node.left = child
	} else {
		node.right = child
	}
}

func (node *node[K, V]) ancestor() *node[K, V] {
	ancestor := node
	for ancestor.parent != nil {
		ancestor = ancestor.parent
	}
	return ancestor
}

func (node *node[K, V]) grandparent() *node[K, V] {
	if node.parent == nil {
		return nil
	}
	return node.parent.parent
}

func (node *node[K, V]) sibling() *node[K, V] {
	if node.parent == nil {
		return nil
	}
	if node.parent.left == node {
		return node.parent.right
	}
	return node.parent.left
}

func (node *node[K, V]) uncle() *node[K, V] {
	if node.parent == nil {
		return nil
	}
	return node.parent.sibling()
}

func (node *node[K, V]) smallest() *node[K, V] {
	var next = node
	for !next.left.null() {
		next = next.left
	}
	return next
}

func (node *node[K, V]) biggest() *node[K, V] {
	var next = node
	for next.right.null() {
		next = next.right
	}
	return next
}
