package query

func Or[T comparable](source, value T) T {
	var zero T
	if source == zero {
		return value
	}
	return source
}
