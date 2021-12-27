package rbtree

import (
	"bytes"
	"constraints"
	"fmt"
	"io"

	"github.com/gopherd/doge/container"
)

type (
	// Iterator represents an iterator of RBTree to iterate nodes
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

	// CompareFunc represents comparation between key
	CompareFunc[T comparable] func(k1, k2 T) bool
)

func Less[T constraints.Ordered](x, y T) bool {
	return x < y
}

func Greater[T constraints.Ordered](x, y T) bool {
	return x > y
}

// RBTree RBComment
type RBTree[K comparable, V any] struct {
	root *node[K, V]
	size int
	cmp  CompareFunc[K]
}

// New creates a RBTree with compare function
func New[K comparable, V any](cmp CompareFunc[K]) *RBTree[K, V] {
	if cmp == nil {
		panic("cmp is nil")
	}
	return &RBTree[K, V]{
		cmp: cmp,
	}
}

// Len returns the number of elements
func (tree RBTree[K, V]) Len() int {
	return tree.size
}

// Clear clears the container
func (tree *RBTree[K, V]) Clear() {
	tree.root = nil
	tree.size = 0
}

// Find finds node by the key, nil returned if the key not found.
func (tree *RBTree[K, V]) Find(key K) Iterator[K, V] {
	return tree.find(key)
}

// Insert inserts a key-value pair, inserted node and true returned
// if the key not found, otherwise, existed node and false returned.
func (tree *RBTree[K, V]) Insert(key K, value V) (Iterator[K, V], bool) {
	node, ok := tree.insert(key, value)
	if ok {
		tree.size++
	}
	return node, ok
}

// Remove removes the key, false returned if the key not found.
func (tree *RBTree[K, V]) Remove(key K) bool {
	node := tree.find(key)
	if node == nil || node.null() {
		return false
	}
	tree.remove(node, true)
	tree.size--
	return true
}

// Erase deletes the node, false returned if the node not found.
func (tree *RBTree[K, V]) Erase(iter Iterator[K, V]) bool {
	if iter == nil {
		return false
	}
	node := iter.underlyingNode()
	if node == nil || node.null() {
		return false
	}
	ok := tree.remove(node, false)
	if ok {
		tree.size--
	}
	return ok
}

// First returns the first node.
//
// usage:
//
//	iter := tree.First()
//	for iter != nil {
//		// hint: do something here using iter
//		// hint: iter.Key(), iter.Value(), iter.SetValue(newValue)
//		iter = iter.Next()
//	}
func (tree *RBTree[K, V]) First() Iterator[K, V] {
	if tree.root == nil {
		return nil
	}
	return tree.root.smallest()
}

// Last returns the first node.
//
// usage:
//
//	iter := tree.Last()
//	for iter != nil {
//		// hint: do something here using iter
//		// hint: iter.Key(), iter.Value(), iter.SetValue(newValue)
//		iter = iter.Prev()
//	}
func (tree *RBTree[K, V]) Last() Iterator[K, V] {
	if tree.root == nil {
		return nil
	}
	return tree.root.biggest()
}

// Format formats the tree
func (tree *RBTree[K, V]) Format(formatter container.TreeFormatter) string {
	return tree.root.format(formatter)
}

// MarshalTree returns a pretty output as a tree
func (tree *RBTree[K, V]) MarshalTree(prefix string) string {
	return tree.root.format(container.TreeFormatter{
		Prefix: prefix,
	})
}

// String returns content of the tree as a plain string
func (tree *RBTree[K, V]) String() string {
	var buf bytes.Buffer
	buf.WriteByte('[')
	iter := tree.First()
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

func (tree *RBTree[K, V]) insert(key K, value V) (*node[K, V], bool) {
	if tree.root == nil {
		tree.root = &node[K, V]{
			color: black,
			key:   key,
			value: value,
		}
		tree.root.left = makenull(tree.root)
		tree.root.right = makenull(tree.root)
		return tree.root, true
	}

	var (
		next     = tree.root
		inserted *node[K, V]
	)
	for {
		if key == next.key {
			next.value = value
			return next, false
		}
		if tree.cmp(key, next.key) {
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
		next = tree.doInsert(next)
		if next == nil {
			break
		}
	}
	return inserted, true
}

func (tree *RBTree[K, V]) find(key K) *node[K, V] {
	var next = tree.root
	for next != nil && !next.null() {
		if next.key == key {
			return next
		}
		if tree.cmp(key, next.key) {
			next = next.left
		} else {
			next = next.right
		}
	}
	return nil
}

func (tree *RBTree[K, V]) remove(n *node[K, V], must bool) bool {
	if !must {
		if tree.root == nil || n == nil || n.ancestor() != tree.root {
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
			tree.root = nil
			return true
		}
		child.parent = nil
		tree.root = child
		tree.root.color = black
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
		child = tree.doRemove(child)
	}
	return true
}

func (tree *RBTree[K, V]) doInsert(n *node[K, V]) *node[K, V] {
	if n.parent == nil {
		tree.root = n
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
		tree.rotateLeft(n.parent)
		n.color = black
		n.parent.color = red
		tree.rotateRight(n.parent)
	} else if n.parent.left == n && n.grandparent().right == n.parent {
		tree.rotateRight(n.parent)
		n.color = black
		n.parent.color = red
		tree.rotateLeft(n.parent)
	} else if n.parent.left == n && n.grandparent().left == n.parent {
		n.parent.color = black
		n.grandparent().color = red
		tree.rotateRight(n.grandparent())
	} else if n.parent.right == n && n.grandparent().right == n.parent {
		n.parent.color = black
		n.grandparent().color = red
		tree.rotateLeft(n.grandparent())
	}
	return nil
}

func (tree *RBTree[K, V]) doRemove(n *node[K, V]) *node[K, V] {
	if n.parent == nil {
		n.color = black
		return nil
	}
	sibling := n.sibling()
	if sibling.color == red {
		n.parent.color = red
		sibling.color = black
		if n == n.parent.left {
			tree.rotateLeft(n.parent)
		} else {
			tree.rotateRight(n.parent)
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
			tree.rotateRight(sibling.left.parent)
		} else if n == n.parent.right &&
			sibling.left.color == black &&
			sibling.right.color == red {
			sibling.color = red
			sibling.right.color = black
			tree.rotateLeft(sibling.right.parent)
		}
	}
	sibling = n.sibling()
	sibling.color = n.parent.color
	n.parent.color = black
	if n == n.parent.left {
		sibling.right.color = black
		tree.rotateLeft(sibling.parent)
	} else {
		sibling.left.color = black
		tree.rotateRight(sibling.parent)
	}
	return nil
}

const (
	left  = 0
	right = 1
)

func (tree *RBTree[K, V]) rotate(p *node[K, V], dir int) *node[K, V] {
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
		tree.root = s
	}
	return s
}

func (tree *RBTree[K, V]) rotateLeft(p *node[K, V]) {
	tree.rotate(p, left)
}

func (tree *RBTree[K, V]) rotateRight(p *node[K, V]) {
	tree.rotate(p, right)
}

type color byte

const (
	red color = iota
	black
)

func (c color) String() string {
	if c == red {
		return "R"
	}
	return "B"
}

// node represents the node of RBTree
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

func (node *node[K, V]) format(formatter container.TreeFormatter) string {
	formatter.Fix()
	if node == nil {
		return "<nil>\n"
	}
	var (
		buf         bytes.Buffer
		prefixstack bytes.Buffer
	)
	if formatter.Prefix != "" {
		prefixstack.WriteString(formatter.Prefix)
	}
	node.print(&buf, &prefixstack, "", 0, formatter)
	return buf.String()
}

func (node *node[K, V]) print(w io.Writer, prefixstack *bytes.Buffer, prefix string, depth int, formatter container.TreeFormatter) {
	var (
		prefixlen    = prefixstack.Len()
		cbegin, cend string
	)
	var (
		vbegin = "\033[0;90m" // gray color code
		vend   = "\033[0m"
	)
	if node.color == red {
		cbegin = "\033[0;31m" // red color code
		cend = vend
	}
	if !formatter.Color {
		vbegin = ""
		vend = ""
		cbegin = ""
		cend = ""
	}

	if node.null() {
		if formatter.Debug {
			fmt.Fprintf(w, "%s%s%snil:%d%s\n", prefixstack.String(), prefix, cbegin, depth+1, cend)
		}
		return
	}
	fmt.Fprintf(w, "%s%s%s%v%s:%s%v%s\n", prefixstack.String(), prefix, cbegin, node.key, cend, vbegin, node.value, vend)

	if node.parent != nil {
		var isLast bool
		if formatter.Debug {
			isLast = node.parent.right == node
		} else {
			isLast = node.parent.right == node || node.parent.right.null()
		}
		if isLast {
			prefixstack.WriteString(formatter.IconSpace)
		} else {
			prefixstack.WriteString(formatter.IconParent)
		}
	}
	var (
		children [2]int
		size     = 0
	)
	if !formatter.Debug {
		if !node.left.null() {
			children[size] = left
			size++
		}
		if !node.right.null() {
			children[size] = right
			size++
		}
	} else {
		size = 2
		children[0] = left
		children[1] = right
	}
	for i := 0; i < size; i++ {
		var appended string
		if i+1 == size {
			appended = formatter.IconLastBranch
		} else {
			appended = formatter.IconBranch
		}
		child := node.child(children[i])
		child.print(w, prefixstack, appended, depth+int(child.color), formatter)
	}

	if prefixlen != prefixstack.Len() {
		prefixstack.Truncate(prefixlen)
	}
}
