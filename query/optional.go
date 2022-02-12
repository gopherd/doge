package query

func Or[T comparable](source, value T) T {
	var zero T
	if source == zero {
		return value
	}
	return source
}

func OrNew[T comparable](source T, new func() T) T {
	var zero T
	if source == zero {
		return new()
	}
	return source
}
