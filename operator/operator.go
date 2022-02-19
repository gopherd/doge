package operator

import "github.com/gopherd/doge/constraints"

// Or returns a || b
func Or[T comparable](a, b T) T {
	var zero T
	if a == zero {
		return b
	}
	return a
}

// OrNew returns a || new()
func OrNew[T comparable](a T, new func() T) T {
	var zero T
	if a == zero {
		return new()
	}
	return a
}

// If returns yes ? a : b
func If[T any](yes bool, a, b T) T {
	if yes {
		return a
	}
	return b
}

// IfNew returns yes ? a() : b()
func IfNew[T any](yes bool, a, b func() T) T {
	if yes {
		return a()
	}
	return b()
}

// Bool converts bool to number
func Bool[T constraints.Number](ok bool) T {
	return If[T](ok, 1, 0)
}

// Less reports whether x less than y
func Less[T constraints.Ordered](x, y T) bool {
	return x < y
}

// Greater reports whether x greather than y
func Greater[T constraints.Ordered](x, y T) bool {
	return x > y
}
