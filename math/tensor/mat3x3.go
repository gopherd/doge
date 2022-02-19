package tensor

import (
	"bytes"
	"fmt"
	"math"

	"github.com/gopherd/doge/constraints"
)

type Matrix3[T constraints.SignedReal] [3 * 3]T

func One3x3[T constraints.SignedReal]() Matrix3[T] {
	return Matrix3[T]{
		1, 0, 0,
		0, 1, 0,
		0, 0, 1,
	}
}

func (mat Matrix3[T]) String() string {
	const dim = 3
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i := 0; i < dim; i++ {
		if i > 0 {
			buf.WriteByte(';')
		}
		buf.WriteByte('(')
		for j := 0; j < dim; j++ {
			if j > 0 {
				buf.WriteByte(',')
			}
			fmt.Fprint(&buf, mat[i+j*dim])
		}
		buf.WriteByte(')')
	}
	buf.WriteByte('}')
	return buf.String()
}

func (mat *Matrix3[T]) SetElements(n11, n12, n13, n21, n22, n23, n31, n32, n33 T) *Matrix3[T] {
	(*mat)[0], (*mat)[3], (*mat)[6] = n11, n12, n13
	(*mat)[1], (*mat)[4], (*mat)[7] = n21, n22, n23
	(*mat)[2], (*mat)[5], (*mat)[8] = n31, n32, n33
	return mat
}

func (mat Matrix3[T]) Sum() T {
	var result T
	for i := range mat {
		result += mat[i]
	}
	return result
}

func (mat Matrix3[T]) Transpose() Matrix3[T] {
	const dim = 3
	for i := 0; i < dim-1; i++ {
		for j := i + 1; j < dim; j++ {
			mat[i+j*dim], mat[j+i*dim] = mat[j+i*dim], mat[i+j*dim]
		}
	}
	return mat
}

func (mat Matrix3[T]) Dot(other Matrix3[T]) Matrix3[T] {
	const dim = 3
	var result Matrix3[T]
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

func (mat Matrix3[T]) DotVec2(vec Vector2[T]) Vector3[T] {
	return mat.DotVec3(vec.Vec3())
}

func (mat Matrix3[T]) DotVec3(vec Vector3[T]) Vector3[T] {
	const dim = 3
	var result Vector3[T]
	for i := 0; i < dim; i++ {
		for j := 0; j < dim; j++ {
			result[i] += mat[i+j*dim] * vec[j]
		}
	}
	return result
}

func (mat Matrix3[T]) Square() T {
	return mat.Hadamard(mat).Sum()
}

func (mat Matrix3[T]) Length() T {
	return T(math.Sqrt(float64(mat.Square())))
}

func (mat Matrix3[T]) Add(other Matrix3[T]) Matrix3[T] {
	for i := range mat {
		mat[i] += other[i]
	}
	return mat
}

func (mat Matrix3[T]) Sub(other Matrix3[T]) Matrix3[T] {
	for i := range mat {
		mat[i] -= other[i]
	}
	return mat
}

func (mat Matrix3[T]) Mul(v T) Matrix3[T] {
	for i := range mat {
		mat[i] *= v
	}
	return mat
}

func (mat Matrix3[T]) Div(v T) Matrix3[T] {
	for i := range mat {
		mat[i] /= v
	}
	return mat
}

func (mat Matrix3[T]) Hadamard(other Matrix3[T]) Matrix3[T] {
	for i := range mat {
		mat[i] *= other[i]
	}
	return mat
}

func (mat Matrix3[T]) Normalize() Matrix3[T] {
	return mat.Div(mat.Length())
}

func (mat Matrix3[T]) Determaint() T {
	var a, b, c = mat[0], mat[1], mat[2]
	var d, e, f = mat[3], mat[4], mat[5]
	var g, h, i = mat[6], mat[7], mat[8]
	return a*e*i - a*f*h - b*d*i + b*f*g + c*d*h - c*e*g
}

func (mat Matrix3[T]) Invert() Matrix3[T] {
	var n11, n21, n31 = mat[0], mat[1], mat[2]
	var n12, n22, n32 = mat[3], mat[4], mat[5]
	var n13, n23, n33 = mat[6], mat[7], mat[8]

	var t11 = n33*n22 - n32*n23
	var t12 = n32*n13 - n33*n12
	var t13 = n23*n12 - n22*n13

	var det = n11*t11 + n21*t12 + n31*t13

	if det == 0 {
		return Matrix3[T]{}
	}

	var detInv = 1 / det

	mat[0] = t11 * detInv
	mat[1] = (n31*n23 - n33*n21) * detInv
	mat[2] = (n32*n21 - n31*n22) * detInv

	mat[3] = t12 * detInv
	mat[4] = (n33*n11 - n31*n13) * detInv
	mat[5] = (n31*n12 - n32*n11) * detInv

	mat[6] = t13 * detInv
	mat[7] = (n21*n13 - n23*n11) * detInv
	mat[8] = (n22*n11 - n21*n12) * detInv

	return mat
}

func (mat *Matrix3[T]) MakeIdentity() *Matrix3[T] {
	return mat.SetElements(
		1, 0, 0,
		0, 1, 0,
		0, 0, 1,
	)
}

func (mat *Matrix3[T]) MakeZero() *Matrix3[T] {
	return mat.SetElements(
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
	)
}

func (mat *Matrix3[T]) MakeRotationX(theta T) *Matrix3[T] {
	var s0, c0 = math.Sincos(float64(theta))
	var s, c = T(s0), T(c0)
	return mat.SetElements(
		1, 0, 0,
		0, c, -s,
		0, s, c,
	)
}

func (mat *Matrix3[T]) MakeRotationY(theta T) *Matrix3[T] {
	var s0, c0 = math.Sincos(float64(theta))
	var s, c = T(s0), T(c0)
	return mat.SetElements(
		c, 0, s,
		0, 1, 0,
		-s, 0, c,
	)

}

func (mat *Matrix3[T]) MakeRotationZ(theta T) *Matrix3[T] {
	var s0, c0 = math.Sincos(float64(theta))
	var s, c = T(s0), T(c0)
	return mat.SetElements(
		c, -s, 0,
		s, c, 0,
		0, 0, 1,
	)
}

// Based on http://www.gamedev.net/reference/articles/article1199.asp
func (mat *Matrix3[T]) MakeRotationAxis(axis Vector3[T], angle T) *Matrix3[T] {
	var s0, c0 = math.Sincos(float64(angle))
	var s, c = T(s0), T(c0)
	var t = 1 - c
	var x, y, z = axis.X(), axis.Y(), axis.Z()
	var tx, ty = t * x, t * y
	return mat.SetElements(
		tx*x+c, tx*y-s*z, tx*z+s*y,
		tx*y+s*z, ty*y+c, ty*z-s*x,
		tx*z-s*y, ty*z+s*x, t*z*z+c,
	)
}

func (mat *Matrix3[T]) MakeScale(vec Vector3[T]) *Matrix3[T] {
	return mat.SetElements(
		vec.X(), 0, 0,
		0, vec.Y(), 0,
		0, 0, vec.Z(),
	)
}

func (mat *Matrix3[T]) MakeShear(xy, xz, yx, yz, zx, zy T) *Matrix3[T] {
	return mat.SetElements(
		1, yx, zx,
		xy, 1, zy,
		xz, yz, 1,
	)
}
