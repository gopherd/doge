package functional

import (
	"math"

	"github.com/gopherd/doge/constraints"
	"github.com/gopherd/doge/math/mathutil"
)

type UnaryFn[T constraints.Field] func(T) T

func (f UnaryFn[T]) Add(f2 UnaryFn[T]) UnaryFn[T] {
	return func(x T) T {
		return f(x) + f2(x)
	}
}

func (f UnaryFn[T]) Sub(f2 UnaryFn[T]) UnaryFn[T] {
	return func(x T) T {
		return f(x) - f2(x)
	}
}

func (f UnaryFn[T]) Mul(f2 UnaryFn[T]) UnaryFn[T] {
	return func(x T) T {
		return f(x) * f2(x)
	}
}

func (f UnaryFn[T]) Div(f2 UnaryFn[T]) UnaryFn[T] {
	return func(x T) T {
		return f(x) / f2(x)
	}
}

func Constant[T constraints.Field](c T) UnaryFn[T] {
	return func(x T) T { return c }
}

func KSigmoid[T constraints.Float](k T) UnaryFn[T] {
	return func(x T) T { return Sigmoid(k * x) }
}

func KSigmoidPrime[T constraints.Float](k T) UnaryFn[T] {
	return func(x T) T { return SigmoidPrime(k*x) * k }
}

func Scale[T constraints.Field](k T) UnaryFn[T] {
	return func(x T) T { return k * x }
}

func Offset[T constraints.Field](b T) UnaryFn[T] {
	return func(x T) T { return x + b }
}

func Affine[T constraints.Field](k, b T) UnaryFn[T] {
	return func(x T) T { return k*x + b }
}

func Power[T constraints.Float](p T) UnaryFn[T] {
	return func(x T) T { return T(math.Pow(float64(x), float64(p))) }
}

func Zero[T constraints.Field](x T) T {
	return 0
}

func One[T constraints.Field](x T) T {
	return 1
}

func Identity[T constraints.Field](x T) T {
	return x
}

func Square[T constraints.Field](x T) T {
	return x * x
}

func Abs[T constraints.Float](x T) T {
	return mathutil.Abs(x)
}

func Sign[T constraints.Float](x T) T {
	if x == 0 {
		return 0
	}
	if x > 0 {
		return 1
	}
	return -1
}

func Sigmoid[T constraints.Float](x T) T {
	return T(1.0 / (1.0 + math.Exp(-float64(x))))
}

func SigmoidPrime[T constraints.Float](x T) T {
	x = Sigmoid(x)
	return x * (1 - x)
}

type BinaryFn[T constraints.Field] func(x, y T) T

func (f BinaryFn[T]) Add(f2 BinaryFn[T]) BinaryFn[T] {
	return func(x, y T) T {
		return f(x, y) + f2(x, y)
	}
}

func (f BinaryFn[T]) Sub(f2 BinaryFn[T]) BinaryFn[T] {
	return func(x, y T) T {
		return f(x, y) - f2(x, y)
	}
}

func (f BinaryFn[T]) Mul(f2 BinaryFn[T]) BinaryFn[T] {
	return func(x, y T) T {
		return f(x, y) * f2(x, y)
	}
}

func (f BinaryFn[T]) Div(f2 BinaryFn[T]) BinaryFn[T] {
	return func(x, y T) T {
		return f(x, y) / f2(x, y)
	}
}

func Add[T constraints.Field](x, y T) T { return x + y }
func Sub[T constraints.Field](x, y T) T { return x - y }
func Mul[T constraints.Field](x, y T) T { return x * y }
func Div[T constraints.Field](x, y T) T { return x / y }
func Pow[T constraints.Float](x, y T) T { return T(math.Pow(float64(x), float64(y))) }
