package mathutil

import (
	"constraints"
	"math"
)

type Real interface {
	constraints.Float | constraints.Integer
}

func MinInt[T constraints.Integer](x, y T) T {
	if x < y {
		return x
	}
	return y
}

func MaxInt[T constraints.Integer](x, y T) T {
	if x > y {
		return x
	}
	return y
}

func MinmaxInt[T constraints.Integer](x, y T) (min, max T) {
	if x <= y {
		return x, y
	}
	return y, x
}

func AbsInt[T constraints.Integer](x T) T {
	if x >= 0 {
		return x
	}
	return -x
}

func MinFloat[T constraints.Float](x, y T) T {
	return T(math.Min(float64(x), float64(y)))
}

func MaxFloat[T constraints.Float](x, y T) T {
	return T(math.Max(float64(x), float64(y)))
}

func MinmaxFloat[T constraints.Float](x, y T) (min, max T) {
	return T(math.Min(float64(x), float64(y))), T(math.Max(float64(x), float64(y)))
}

func AbsFloat[T constraints.Float](x T) T {
	return T(math.Abs(float64(x)))
}
