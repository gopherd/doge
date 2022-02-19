package tensor

import (
	"bytes"
	"fmt"
	"math"

	"github.com/gopherd/doge/constraints"
)

type Matrix2[T constraints.SignedReal] [2 * 2]T

func One2x2[T constraints.SignedReal]() Matrix2[T] {
	return Matrix2[T]{
		1, 0,
		0, 1,
	}
}

func (mat Matrix2[T]) String() string {
	const n = 2
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i := 0; i < n; i++ {
		if i > 0 {
			buf.WriteByte(';')
		}
		buf.WriteByte('(')
		for j := 0; j < n; j++ {
			if j > 0 {
				buf.WriteByte(',')
			}
			fmt.Fprint(&buf, mat[i+j*n])
		}
		buf.WriteByte(')')
	}
	buf.WriteByte('}')
	return buf.String()
}

func (mat *Matrix2[T]) SetElements(n11, n12, n21, n22 T) *Matrix2[T] {
	(*mat)[0], (*mat)[2] = n11, n12
	(*mat)[1], (*mat)[3] = n21, n22
	return mat
}

func (mat Matrix2[T]) Sum() T {
	var result T
	for i := range mat {
		result += mat[i]
	}
	return result
}

func (mat Matrix2[T]) Transpose() Matrix2[T] {
	const dim = 2
	for i := 0; i < dim-1; i++ {
		for j := i + 1; j < dim; j++ {
			mat[i+j*dim], mat[j+i*dim] = mat[j+i*dim], mat[i+j*dim]
		}
	}
	return mat
}

func (mat Matrix2[T]) Dot(other Matrix2[T]) Matrix2[T] {
	const dim = 2
	var result Matrix2[T]
	for i := 0; i < dim; i++ {
		for j := 0; j < dim; j++ {
			index := i + j*dim
			for k := 0; k < dim; k++ {
				result[index] += mat[i+k*dim] * other[k+j*dim]
			}
		}
	}
	return result
}

func (mat Matrix2[T]) DotVec2(vec Vector2[T]) Vector2[T] {
	const dim = 2
	var result Vector2[T]
	for i := 0; i < dim; i++ {
		for j := 0; j < dim; j++ {
			result[i] += mat[i+j*dim] * vec[j]
		}
	}
	return result
}

func (mat Matrix2[T]) Square() T {
	return mat.Hadamard(mat).Sum()
}

func (mat Matrix2[T]) Length() T {
	return T(math.Sqrt(float64(mat.Square())))
}

func (mat Matrix2[T]) Add(other Matrix2[T]) Matrix2[T] {
	for i := range mat {
		mat[i] += other[i]
	}
	return mat
}

func (mat Matrix2[T]) Sub(other Matrix2[T]) Matrix2[T] {
	for i := range mat {
		mat[i] -= other[i]
	}
	return mat
}

func (mat Matrix2[T]) Mul(v T) Matrix2[T] {
	for i := range mat {
		mat[i] *= v
	}
	return mat
}

func (mat Matrix2[T]) Div(v T) Matrix2[T] {
	for i := range mat {
		mat[i] /= v
	}
	return mat
}

func (mat Matrix2[T]) Hadamard(other Matrix2[T]) Matrix2[T] {
	for i := range mat {
		mat[i] *= other[i]
	}
	return mat
}

func (mat Matrix2[T]) Normalize() Matrix2[T] {
	return mat.Div(mat.Length())
}

func (mat Matrix2[T]) Determaint() T {
	return mat[0]*mat[3] - mat[1]*mat[2]
}

func (mat Matrix2[T]) Invert() Matrix2[T] {
	var n11, n21 = mat[0], mat[1]
	var n12, n22 = mat[2], mat[3]
	var det = n11*n22 - n12*n21
	if det == 0 {
		return Matrix2[T]{}
	}
	var detInv = 1 / det
	mat[0] = n11 * detInv
	mat[1] = -n21 * detInv
	mat[2] = -n12 * detInv
	mat[3] = n22 * detInv
	return mat
}

func (mat *Matrix2[T]) MakeIdentity() *Matrix2[T] {
	return mat.SetElements(
		1, 0,
		0, 1,
	)
}

func (mat *Matrix2[T]) MakeZero() *Matrix2[T] {
	return mat.SetElements(
		0, 0,
		0, 0,
	)
}

func (mat *Matrix2[T]) MakeRotation(theta T) *Matrix2[T] {
	var s0, c0 = math.Sincos(float64(theta))
	var s, c = T(s0), T(c0)
	return mat.SetElements(
		c, -s,
		s, c,
	)
}

func (mat *Matrix2[T]) MakeScale(vec Vector2[T]) *Matrix2[T] {
	return mat.SetElements(
		vec.X(), 0,
		0, vec.Y(),
	)
}
