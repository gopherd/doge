package functional

import (
	"math"

	"github.com/gopherd/doge/constraints"
	"github.com/gopherd/doge/math/mathutil"
)

type UnaryFn[T constraints.Number] func(T) T

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

func Constant[T constraints.Number](c T) UnaryFn[T] {
	return func(x T) T { return c }
}

func KSigmoid[T constraints.Real](k T) UnaryFn[T] {
	return func(x T) T { return Sigmoid(k * x) }
}

func KSigmoidPrime[T constraints.Real](k T) UnaryFn[T] {
	return func(x T) T { return SigmoidPrime(k*x) * k }
}

func Scale[T constraints.Number](k T) UnaryFn[T] {
	return func(x T) T { return k * x }
}

func Offset[T constraints.Number](b T) UnaryFn[T] {
	return func(x T) T { return x + b }
}

func Affine[T constraints.Number](k, b T) UnaryFn[T] {
	return func(x T) T { return k*x + b }
}

func Power[T constraints.Real](p T) UnaryFn[T] {
	return func(x T) T { return T(math.Pow(float64(x), float64(p))) }
}

func Zero[T constraints.Number](x T) T {
	return 0
}

func One[T constraints.Number](x T) T {
	return 1
}

func Identity[T constraints.Number](x T) T {
	return x
}

func Square[T constraints.Number](x T) T {
	return x * x
}

func Abs[T constraints.SignedReal](x T) T {
	return mathutil.Abs(x)
}

func Sign[T constraints.SignedReal](x T) T {
	if x == 0 {
		return 0
	}
	if x > 0 {
		return 1
	}
	return -1
}

func Sigmoid[T constraints.Real](x T) T {
	return T(1.0 / (1.0 + math.Exp(-float64(x))))
}

func SigmoidPrime[T constraints.Real](x T) T {
	x = Sigmoid(x)
	return x * (1 - x)
}

type BinaryFn[T constraints.Number] func(x, y T) T

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

func Add[T constraints.Number](x, y T) T { return x + y }
func Sub[T constraints.Number](x, y T) T { return x - y }
func Mul[T constraints.Number](x, y T) T { return x * y }
func Div[T constraints.Number](x, y T) T { return x / y }
func Pow[T constraints.Real](x, y T) T   { return T(math.Pow(float64(x), float64(y))) }
