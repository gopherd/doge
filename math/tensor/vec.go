package tensor

import (
	"math"

	"github.com/gopherd/doge/constraints"
	"github.com/gopherd/doge/math/mathutil"
	"github.com/gopherd/doge/operator"
)

// Vector implements n-dim vector
type Vector[T constraints.SignedReal] []T

func Vec[T constraints.SignedReal](elements ...T) Vector[T] {
	return Vector[T](elements)
}

func (vec Vector[T]) Dim() int       { return len(vec) }
func (vec Vector[T]) Get(i int) T    { return operator.If(i < len(vec), vec[i], 0) }
func (vec Vector[T]) Set(i int, v T) { vec[i] = v }

func (vec Vector[T]) Data() []T {
	return []T(vec)
}

func (vec Vector[T]) Swap(i, j int) {
	vec[i], vec[j] = vec[j], vec[i]
}

func (vec Vector[T]) Less(i, j int) bool {
	return vec[i] < vec[j]
}

func (vec Vector[T]) Add(other Vector[T]) Vector[T] {
	var min, max = mathutil.Minmax(vec.Dim(), other.Dim())
	var out = make(Vector[T], max)
	for i := 0; i < min; i++ {
		out[i] = vec[i] + other[i]
	}
	return out
}

func (vec Vector[T]) Sub(other Vector[T]) Vector[T] {
	var min, max = mathutil.Minmax(vec.Dim(), other.Dim())
	var out = make(Vector[T], max)
	for i := 0; i < min; i++ {
		out[i] = vec[i] - other[i]
	}
	return out
}

func (vec Vector[T]) Mul(other Vector[T]) Vector[T] {
	var min, max = mathutil.Minmax(vec.Dim(), other.Dim())
	var out = make(Vector[T], max)
	for i := 0; i < min; i++ {
		out[i] = vec[i] * other[i]
	}
	return out
}

func (vec Vector[T]) Div(other Vector[T]) Vector[T] {
	var min, max = mathutil.Minmax(vec.Dim(), other.Dim())
	var out = make(Vector[T], max)
	for i := 0; i < max; i++ {
		out[i] = vec[i] / operator.If(i < min, other[i], 0)
	}
	return out
}

func (vec Vector[T]) Sum() T {
	var sum T
	for i := range vec {
		sum += vec[i]
	}
	return sum
}

func (vec Vector[T]) Dot(other Vector[T]) T {
	var sum T
	for i := range vec {
		if i >= len(other) {
			break
		}
		sum += vec[i] * other[i]
	}
	return sum
}

func (vec Vector[T]) SquaredLength() T {
	var sum T
	for i := range vec {
		sum += vec[i] * vec[i]
	}
	return sum
}

func (vec Vector[T]) Norm() T {
	return T(math.Sqrt(float64(vec.SquaredLength())))
}

func (vec Vector[T]) Normp(p T) T {
	var sum float64
	for i := range vec {
		sum += math.Pow(float64(vec[i]), float64(p))
	}
	return T(math.Pow(sum, 1.0/float64(p)))
}

// Range creates a vector [start, end)
func Range[T constraints.SignedReal](start, end T) Vector[T] {
	if end <= start {
		return nil
	}
	var vec = make(Vector[T], 0, int(end-start))
	for start < end {
		vec = append(vec, start)
		start++
	}
	return vec
}

// RangeN creates a vector [0, n)
func RangeN[T constraints.SignedReal](n T) Vector[T] {
	return Range(0, n)
}

// Linspace creates a vector {x1=from, ..., xn=to}
func Linspace[T constraints.SignedReal](from, to T, n int) Vector[T] {
	if n < 1 {
		return nil
	}
	if n == 1 {
		return Vec(from)
	}
	var interval = (to - from) / T(n-1)
	var vec = make(Vector[T], n)
	for i := 0; i < n; i++ {
		if i+1 == n {
			vec[i] = to
		} else {
			vec[i] = from + T(i)*interval
		}
	}
	return vec
}
