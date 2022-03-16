package tensor

import (
	"github.com/gopherd/doge/constraints"
	"github.com/gopherd/doge/container/tuple"
)

type Shape = tuple.Tuple[int]

type Index int

func (i Index) Len() int     { return 1 }
func (i Index) At(j int) int { return int(i) }

type Indices []int

func (i Indices) Len() int     { return len(i) }
func (i Indices) At(j int) int { return i[j] }

type Tensor[T constraints.SignedReal] interface {
	Shape() tuple.Tuple[int]
	At(index tuple.Tuple[int]) T
	Sum() T
}

// Create creates a tensor by shape
func Create[T constraints.SignedReal](shape Shape) Tensor[T] {
	return tensor[T]{
		shape: shape,
		data:  make(Vector[T], sizeof(shape)),
	}
}

// tensor implements Tensor with shape
type tensor[T constraints.SignedReal] struct {
	shape Shape
	data  Vector[T]
}

// Shape implements Tensor Shape method
func (t tensor[T]) Shape() Shape {
	return t.shape
}

// At implements Tensor At method
func (t tensor[T]) At(index Shape) T {
	return t.data[offsetof(t.shape, index)]
}

// Sum implements Tensor Sum method
func (t tensor[T]) Sum() T {
	return t.data.Sum()
}

// set updates value by index
func (t *tensor[T]) set(index Shape, value T) {
	t.data[offsetof(t.shape, index)] = value
}

func offsetof(shape, index Shape) int {
	var off int
	var prev = 1
	for i, n := 0, index.Len(); i < n; i++ {
		off += index.At(i) * prev
		prev = shape.At(i)
	}
	return off
}

func sizeof(shape Shape) int {
	var size int
	for i, n := 0, shape.Len(); i < n; i++ {
		size *= shape.At(i)
	}
	return size
}

func firstof(shape Shape) int {
	return shape.At(0)
}

func lastof(shape Shape) int {
	return shape.At(shape.Len() - 1)
}

func next(shape Shape, index Indices) Indices {
	for i, x := range index {
		if x < shape.At(i) {
		}
	}
}

// Product computes tensor product: a⊗b
func Product[T constraints.SignedReal](a, b Tensor[T]) Tensor[T] {
	var ashape, bshape = a.Shape(), b.Shape()
	return productedTensor[T]{
		a:     a,
		b:     b,
		m:     ashape.Len(),
		n:     bshape.Len(),
		shape: tuple.Concat(ashape, bshape),
	}
}

type productedTensor[T constraints.SignedReal] struct {
	a, b  Tensor[T]
	m, n  int
	shape Shape
}

// Shape implements Tensor Shape method
func (t productedTensor[T]) Shape() Shape {
	return t.shape
}

// At implements Tensor At method
func (t productedTensor[T]) At(index Shape) T {
	var i, j = tuple.Slice(index, 0, t.m), tuple.Slice(index, t.m, t.m+t.n)
	return t.a.At(i) * t.b.At(j)
}

// Sum implements Tensor Sum method
func (t productedTensor[T]) Sum() T {
	return t.a.Sum() * t.b.Sum()
}

// Dot computes dot product: a‧b
func Dot[T constraints.SignedReal](a, b Tensor[T]) Tensor[T] {
	var ashape = a.Shape()
	var bshape = b.Shape()
	if l := ashape.Len(); (l == 0) != (bshape.Len() == 0) {
		panic("tensor.dot: size mismatched")
	} else if l == 0 {
		return Scalar(a.Sum() * b.Sum())
	}
	var n = lastof(ashape)
	if n != firstof(bshape) {
		panic("tensor.dot: size mismatched")
	}
	var shape = tuple.Concat(
		tuple.Slice(ashape, 0, ashape.Len()-1),
		tuple.Slice(bshape, 1, ashape.Len()),
	)
	if shape.Len() == 0 {
		var sum T
		for i := Index(0); i < Index(n); i++ {
			sum += a.At(i) * b.At(i)
		}
		return Scalar(sum)
	}
	var c = Create[T](shape)
	// TODO: computes tensor dot
	return c
}
