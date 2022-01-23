package vector

type Vector[T any] []T

func (v Vector[T]) Len() int {
	return len(v)
}

func (v Vector[T]) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v Vector[T]) Front() T {
	return v[0]
}

func (v Vector[T]) Back() T {
	return v[len(v)-1]
}

func (v *Vector[T]) PushBack(x T) {
	*v = append(*v, x)
}

func (v *Vector[T]) PushFront(x T) {
	if cap(*v) > len(*v) {
		*v = (*v)[:len(*v)+1]
	} else {
		*v = append(*v, x)
	}
	copy((*v)[1:], (*v)[0:len(*v)-1])
	(*v)[0] = x
}
