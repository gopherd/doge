package math

import "constraints"

func Min[T constraints.Integer](x, y T) T {
	if x < y {
		return x
	}
	return y
}

func Max[T constraints.Integer](x, y T) T {
	if x > y {
		return x
	}
	return y
}

func Minmax[T constraints.Integer](x, y T) (min, max T) {
	if x <= y {
		return x, y
	}
	return y, x
}

func Abs[T constraints.Integer](x T) T {
	if x >= 0 {
		return x
	}
	return -x
}
