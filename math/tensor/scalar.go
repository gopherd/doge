package tensor

import (
	"github.com/gopherd/doge/constraints"
	"github.com/gopherd/doge/container/tuple"
	"github.com/gopherd/doge/math/mathutil"
)

type scalar[T constraints.SignedReal] struct {
	x T
}

func Scalar[T constraints.SignedReal](x T) Tensor[T] {
	return scalar[T]{x}
}

var shape0 Shape = tuple.Empty[int]{}

// Shape implements Tensor Shape method
func (s scalar[T]) Shape() Shape {
	return shape0
}

// At implements Tensor At method
func (s scalar[T]) At(index Shape) T {
	panic("out of range")
}

// Sum implements Tensor Sum method
func (s scalar[T]) Sum() T {
	return s.x
}

func (s scalar[T]) Norm() T {
	return mathutil.Abs(s.x)
}

func (s scalar[T]) Normp(p T) T {
	return mathutil.Abs(s.x)
}
