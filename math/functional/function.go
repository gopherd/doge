package functional

import (
	"math"
	"math/rand"
)

type UnaryFn func(float64) float64

func (f UnaryFn) Add(f2 UnaryFn) UnaryFn {
	return func(x float64) float64 {
		return f(x) + f2(x)
	}
}

func (f UnaryFn) Sub(f2 UnaryFn) UnaryFn {
	return func(x float64) float64 {
		return f(x) - f2(x)
	}
}

func (f UnaryFn) Mul(f2 UnaryFn) UnaryFn {
	return func(x float64) float64 {
		return f(x) * f2(x)
	}
}

func (f UnaryFn) Div(f2 UnaryFn) UnaryFn {
	return func(x float64) float64 {
		return f(x) / f2(x)
	}
}

func Shuffle(vec []int) {
	for i := len(vec) - 1; i >= 0; i-- {
		index := rand.Intn(i + 1)
		vec[i], vec[index] = vec[index], vec[i]
	}
}

func Constant(c float64) UnaryFn { return func(x float64) float64 { return c } }
func KSigmoid(k float64) UnaryFn { return func(x float64) float64 { return Sigmoid(k * x) } }
func KSigmoidPrime(k float64) UnaryFn {
	return func(x float64) float64 { return SigmoidPrime(k*x) * k }
}

func Scale(k float64) UnaryFn     { return func(x float64) float64 { return k * x } }
func Offset(b float64) UnaryFn    { return func(x float64) float64 { return x + b } }
func Affine(k, b float64) UnaryFn { return func(x float64) float64 { return k*x + b } }

func Power(p float64) UnaryFn { return func(x float64) float64 { return math.Pow(x, p) } }

var (
	ConstantOne UnaryFn = func(x float64) float64 { return 1 }
	Identity    UnaryFn = func(x float64) float64 { return x }
	Square      UnaryFn = func(x float64) float64 { return x * x }
	Abs         UnaryFn = func(x float64) float64 {
		return float64(math.Abs(float64(x)))
	}
	Sign UnaryFn = func(x float64) float64 {
		if x > 0 {
			return 1
		}
		return -1
	}
	Sigmoid UnaryFn = func(x float64) float64 {
		return float64(1.0 / (1.0 + math.Exp(-float64(x))))
	}
	SigmoidPrime UnaryFn = func(x float64) float64 {
		x = Sigmoid(x)
		return x * (1 - x)
	}
)

type BinaryFn func(x, y float64) float64

func (f BinaryFn) Add(f2 BinaryFn) BinaryFn {
	return func(x, y float64) float64 {
		return f(x, y) + f2(x, y)
	}
}

func (f BinaryFn) Sub(f2 BinaryFn) BinaryFn {
	return func(x, y float64) float64 {
		return f(x, y) - f2(x, y)
	}
}

func (f BinaryFn) Mul(f2 BinaryFn) BinaryFn {
	return func(x, y float64) float64 {
		return f(x, y) * f2(x, y)
	}
}

func (f BinaryFn) Div(f2 BinaryFn) BinaryFn {
	return func(x, y float64) float64 {
		return f(x, y) / f2(x, y)
	}
}

var (
	Add BinaryFn = func(x, y float64) float64 { return x + y }
	Sub BinaryFn = func(x, y float64) float64 { return x - y }
	Mul BinaryFn = func(x, y float64) float64 { return x * y }
	Div BinaryFn = func(x, y float64) float64 { return x / y }
	Pow BinaryFn = math.Pow
)
