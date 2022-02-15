package operator

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

// Conditional returns yes ? a : b
func Conditional[T any](yes bool, a, b T) T {
	if yes {
		return a
	}
	return b
}

// ConditionalNew returns yes ? a() : b()
func ConditionalNew[T any](yes bool, a, b func() T) T {
	if yes {
		return a()
	}
	return b()
}
