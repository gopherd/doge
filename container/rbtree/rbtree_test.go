package rbtree_test

import (
	"math/rand"
	"testing"

	. "github.com/gopherd/doge/container/rbtree"
)

func TestRBTree(t *testing.T) {
	tree := New()
	hashmap := make(map[K]V)

	rand.Seed(100)

	//template<K>
	makeKey := func(i int) K {
		//return strconv.Itoa(i)
		return i
	}
	//template<V>
	makeValue := func(i int) V {
		return i
	}

	add := func(k K, v V) {
		_, found := hashmap[k]
		hashmap[k] = v
		_, ok := tree.Insert(k, v)
		if ok != !found {
			t.Fatalf("tree.Set: returned value want %v, but got %v", !found, found)
		}
	}
	_ = add

	remove := func(k K) {
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
		keys = 20
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
	t.Log("tree.pretty:\n" + tree.Format(FormatOptions{
		Prefix: "... ",
		Color:  true,
		Debug:  true,
	}))
	t.Log("tree.plain:\n" + tree.Format(FormatOptions{}))
	t.Log("tree.string:" + tree.String())

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

func checkTree(op string, t *testing.T, tree *RBTree, hashmap map[K]V) {
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
