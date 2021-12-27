package rbtree_test

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"

	"github.com/gopherd/doge/container"
	. "github.com/gopherd/doge/container/rbtree"
)

func ExampleRBTree() {
	tree := New[int, int](Greater[int])
	fmt.Print("empty:\n" + tree.Format(container.TreeFormatter{}))

	tree.Insert(1, 2)
	tree.Insert(2, 4)
	tree.Insert(4, 8)
	tree.Insert(8, 16)

	fmt.Print("default:\n" + tree.Format(container.TreeFormatter{}))
	fmt.Print("custom:\n" + tree.Format(container.TreeFormatter{
		Prefix:         "... ",
		IconParent:     "|  ",
		IconBranch:     "|--",
		IconLastBranch: "+--",
	}))
	fmt.Println("plain:\n" + tree.String())

	// Output:
	// empty:
	// <nil>
	// default:
	// 2:4
	// ├── 4:8
	// │   └── 8:16
	// └── 1:2
	// custom:
	// ... 2:4
	// ... |-- 4:8
	// ... |   +-- 8:16
	// ... +-- 1:2
	// plain:
	// [8:16 4:8 2:4 1:2]
}

func TestRBTree(t *testing.T) {
	tree := New[int, string](func(k1, k2 int) bool { return k1 > k2 })
	hashmap := make(map[int]string)

	rand.Seed(100)

	makeKey := func(i int) int {
		return i
	}
	makeValue := func(i int) string {
		return strconv.Itoa(i)
	}

	add := func(k int, v string) {
		_, found := hashmap[k]
		hashmap[k] = v
		_, ok := tree.Insert(k, v)
		if ok != !found {
			t.Fatalf("tree.Set: returned value want %v, but got %v", !found, found)
		}
	}
	_ = add

	remove := func(k int) {
		_, found := hashmap[k]
		delete(hashmap, k)
		ok := tree.Remove(k)
		if ok != found {
			t.Fatalf("tree.Remove: want %v, but got %v", found, ok)
		}
	}
	_ = remove

	const (
		n    = 100
		keys = 30
	)
	for i := 0; i < n; i++ {
		for j := 0; j < keys/2; j++ {
			key := makeKey(j)
			value := makeValue(j * (i + 1))
			add(key, value)
			key = makeKey(keys - 1 - j)
			value = makeValue((keys - 1 - j) * (i + 1))
			add(key, value)
		}
		checkTree("add", t, tree, hashmap)
	}
	for j := 0; j < keys; j++ {
		key := makeKey(j)
		remove(key)
	}
	checkTree("remove", t, tree, hashmap)

	for i := 0; i < n; i++ {
		k := makeKey(rand.Intn(keys))
		var op string
		if rand.Intn(2) == 0 {
			op = "add"
			v := makeValue(rand.Intn(99999999) * 3)
			add(k, v)
		} else {
			op = "remove"
			remove(k)
		}
		checkTree(op, t, tree, hashmap)
	}
}

func checkTree[K comparable, V comparable](op string, t *testing.T, tree *RBTree[K, V], hashmap map[K]V) {
	if tree.Len() != len(hashmap) {
		t.Fatalf("[%s] len mismacthed: want %d, got %d", op, len(hashmap), tree.Len())
	}
	for k, v := range hashmap {
		node := tree.Find(k)
		if node == nil {
			t.Fatalf("[%s] key %v not found", op, k)
		}
		if node.Value() != v {
			t.Fatalf("[%s] value mismacthed for key %v: want %v, got %v", op, k, v, node.Value())
		}
	}
}
