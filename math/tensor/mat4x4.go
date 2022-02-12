package tensor

import (
	"math"

	"github.com/gopherd/doge/math/mathutil"
)

type Mat4x4[T mathutil.Real] [4 * 4]T

func One4x4[T mathutil.Real]() Mat4x4[T] {
	return Mat4x4[T]{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
}

func (mat Mat4x4[T]) Get(i, j int) T {
	return mat[j+i*4]
}

func (mat *Mat4x4[T]) Set(i, j int, value T) {
	(*mat)[j+i*4] = value
}

func (mat Mat4x4[T]) Sum() T {
	var result T
	for i := range mat {
		result += mat[i]
	}
	return result
}

func (mat Mat4x4[T]) Transpose() Mat4x4[T] {
	for i := 0; i < 3; i++ {
		for j := i + 1; j < 4; j++ {
			mat[i+j*4], mat[j+i*4] = mat[j+i*4], mat[i+j*4]
		}
	}
	return mat
}

func (mat Mat4x4[T]) Dot(other Mat4x4[T]) Mat4x4[T] {
	var result Mat4x4[T]
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			index := j + i*4
			for k := 0; k < 4; k++ {
				result[index] += mat[k+i*4] * other[j+k*4]
			}
		}
	}
	return result
}

func (mat Mat4x4[T]) DotVec2(vec Vector2[T]) Vector3[T] {
	return mat.DotVec4(vec.Vec4()).Vec3()
}

func (mat Mat4x4[T]) DotVec3(vec Vector3[T]) Vector3[T] {
	return mat.DotVec4(vec.Vec4()).Vec3()
}

func (mat Mat4x4[T]) DotVec4(vec Vector4[T]) Vector4[T] {
	var result Vector4[T]
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			result[i] += mat[j+i*4] * vec[j]
		}
	}
	return result
}

func (mat Mat4x4[T]) Square() T {
	return mat.Hadamard(mat).Sum()
}

func (mat Mat4x4[T]) Length() T {
	return T(math.Sqrt(float64(mat.Square())))
}

func (mat Mat4x4[T]) Add(other Mat4x4[T]) Mat4x4[T] {
	for i := range mat {
		mat[i] += other[i]
	}
	return mat
}

func (mat Mat4x4[T]) Sub(other Mat4x4[T]) Mat4x4[T] {
	for i := range mat {
		mat[i] -= other[i]
	}
	return mat
}

func (mat Mat4x4[T]) Mul(v T) Mat4x4[T] {
	for i := range mat {
		mat[i] *= v
	}
	return mat
}

func (mat Mat4x4[T]) Div(v T) Mat4x4[T] {
	for i := range mat {
		mat[i] /= v
	}
	return mat
}

func (mat Mat4x4[T]) Hadamard(other Mat4x4[T]) Mat4x4[T] {
	for i := range mat {
		mat[i] *= other[i]
	}
	return mat
}

func (mat Mat4x4[T]) Normalize() Mat4x4[T] {
	return mat.Div(mat.Length())
}
