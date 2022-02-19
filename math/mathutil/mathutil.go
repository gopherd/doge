package mathutil

import (
	"math"

	"github.com/gopherd/doge/constraints"
)

func Min[T constraints.Ordered](x, y T) T {
	if x < y {
		return x
	}
	return y
}

func Max[T constraints.Ordered](x, y T) T {
	if x > y {
		return x
	}
	return y
}

func Minmax[T constraints.Float](x, y T) (min, max T) {
	if x < y {
		return x, y
	}
	return y, x
}

func Abs[T constraints.SignedReal](x T) T {
	if x < 0 {
		return -x
	}
	return x
}

func Clamp[T constraints.Ordered](value, min, max T) T {
	return Max(min, Min(max, value))
}

func EuclideanModulo[T constraints.Float](x, y T) T {
	var x64, y64 = float64(x), float64(y)
	return T(math.Mod(math.Mod(x64, y64)+y64, y64))
}

// MapLinear mapping from range <a1, a2> to range <b1, b2>
func MapLinear[T constraints.Field](x, a1, a2, b1, b2 T) T {
	return b1 + (x-a1)*(b2-b1)/(a2-a1)
}

// https://www.gamedev.net/tutorials/programming/general-and-gameplay-programming/inverse-lerp-a-super-useful-yet-often-overlooked-function-r5230/
func InverseLerp[T constraints.Field](x, y, value T) T {
	if x == y {
		return 0
	}
	return (value - x) / (y - x)
}

// https://en.wikipedia.org/wiki/Linear_interpolation
func Lerp[T constraints.Field](x, y, t T) T {
	return (1-t)*x + t*y
}

// http://www.rorydriscoll.com/2016/03/07/frame-rate-independent-damping-using-lerp/
func Damp[T constraints.Float](x, y, lambda, dt T) T {
	return Lerp(x, y, 1-T(math.Exp(float64(-lambda*dt))))
}

// https://www.desmos.com/calculator/vcsjnyz7x4
func PingPong[T constraints.Float](x, length T) T {
	return length - Abs(EuclideanModulo(x, length*2)-length)
}

// http://en.wikipedia.org/wiki/Smoothstep
func SmoothStep[T constraints.Float](x, min, max T) T {
	if x <= min {
		return 0
	} else if x >= max {
		return 1
	}
	x = (x - min) / (max - min)
	return x * x * (3 - 2*x)
}

func SmoothStepFunc[T constraints.Float](x, min, max T, fn func(T) T) T {
	if x <= min {
		return 0
	} else if x >= max {
		return 1
	}
	x = (x - min) / (max - min)
	return fn(x)
}

func IsPowerOfTwo[T constraints.Integer](value T) bool {
	return (value&(value-1)) == 0 && value != 0
}

func CeilPowerOfTwo[T constraints.Integer](value T) T {
	return T(math.Pow(2, math.Ceil(math.Log(float64(value))/math.Ln2)))
}

func FloorPowerOfTwo[T constraints.Integer](value T) T {
	return T(math.Pow(2, math.Floor(math.Log(float64(value))/math.Ln2)))
}

const deg2Rad = math.Pi / 180
const rad2Deg = 180 / math.Pi

func Deg2Rad[T constraints.Float](deg T) T {
	return deg * deg2Rad
}

func Rad2Deg[T constraints.Float](rad T) T {
	return rad * rad2Deg
}
