package constraints

// Signed is a constraint that permits any signed integer type.
type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// Unsigned is a constraint that permits any unsigned integer type.
type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// Integer is a constraint that permits any integer type.
type Integer interface {
	Signed | Unsigned
}

// Float is a constraint that permits any floating-point type.
type Float interface {
	~float32 | ~float64
}

// Complex is a constraint that permits any complex numeric type.
type Complex interface {
	~complex64 | ~complex128
}

// Ordered is a constraint that permits any ordered type: any type
// that supports the operators < <= >= >.
type Ordered interface {
	Integer | Float | ~string
}

// SignedReal is a constraint that permits any signed integer and floating-point type.
type SignedReal interface {
	Signed | Float
}

// Real is a constraint that permits any real number type.
type Real interface {
	Integer | Float
}

// SignedNumber is a constraint that permits any signed integer, floating-point and complex-numeric type.
type SignedNumber interface {
	SignedReal | Complex
}

// Number is a constraint that permits any number type.
type Number interface {
	Real | Complex
}

// Field is a constraint that permits any number field type.
type Field interface {
	Float | Complex
}

func Equal[T comparable](x, y T) bool {
	return x == y
}

// Less reports whether x is less than y.
// For floating-point types, a NaN is considered less than any non-NaN,
// and -0.0 is not less than (is equal to) 0.0.
func Less[T Ordered](x, y T) bool {
	return (isNaN(x) && !isNaN(y)) || x < y
}

func Greater[T Ordered](x, y T) bool {
	return (isNaN(x) && !isNaN(y)) || x > y
}

// Asc returns
//
//	-1 if x is less than y,
//	 0 if x equals y,
//	+1 if x is greater than y.
//
// For floating-point types, a NaN is considered less than any non-NaN,
// a NaN is considered equal to a NaN, and -0.0 is equal to 0.0.
func Asc[T Ordered](x, y T) int {
	xNaN := isNaN(x)
	yNaN := isNaN(y)
	if xNaN && yNaN {
		return 0
	}
	if xNaN || x < y {
		return -1
	}
	if yNaN || x > y {
		return +1
	}
	return 0
}

// Dec returns
//
//	-1 if x is greater than y,
//	 0 if x equals y,
//	+1 if x is less than y.
//
// For floating-point types, a NaN is considered less than any non-NaN,
// a NaN is considered equal to a NaN, and -0.0 is equal to 0.0.
func Dec[T Ordered](x, y T) int {
	return Asc(y, x)
}

// isNaN reports whether x is a NaN without requiring the math package.
// This will always return false if T is not floating-point.
func isNaN[T Ordered](x T) bool {
	return x != x
}
