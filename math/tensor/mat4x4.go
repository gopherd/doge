package tensor

import "math"

type Mat4x4 [4 * 4]float32

var One4x4 = Mat4x4{
	1, 0, 0, 0,
	0, 1, 0, 0,
	0, 0, 1, 0,
	0, 0, 0, 1,
}

func (mat Mat4x4) Get(i, j int) float32 {
	return mat[j+i*4]
}

func (mat *Mat4x4) Set(i, j int, value float32) {
	(*mat)[j+i*4] = value
}

func (mat Mat4x4) Sum() float32 {
	var result float32
	for i := range mat {
		result += mat[i]
	}
	return result
}

func (mat Mat4x4) Transpose() Mat4x4 {
	for i := 0; i < 3; i++ {
		for j := i + 1; j < 4; j++ {
			mat[i+j*4], mat[j+i*4] = mat[j+i*4], mat[i+j*4]
		}
	}
	return mat
}

func (mat Mat4x4) Dot(other Mat4x4) Mat4x4 {
	var result Mat4x4
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

func (mat Mat4x4) DotVec2(vec Vector2) Vector3 {
	return mat.DotVec4(vec.Vec4()).Vec3()
}

func (mat Mat4x4) DotVec3(vec Vector3) Vector3 {
	return mat.DotVec4(vec.Vec4()).Vec3()
}

func (mat Mat4x4) DotVec4(vec Vector4) Vector4 {
	var result Vector4
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			result[i] += mat[j+i*4] * vec[j]
		}
	}
	return result
}

func (mat Mat4x4) Square() float32 {
	return mat.Hadamard(mat).Sum()
}

func (mat Mat4x4) Length() float32 {
	return float32(math.Sqrt(float64(mat.Square())))
}

func (mat Mat4x4) Add(other Mat4x4) Mat4x4 {
	for i := range mat {
		mat[i] += other[i]
	}
	return mat
}

func (mat Mat4x4) Sub(other Mat4x4) Mat4x4 {
	for i := range mat {
		mat[i] -= other[i]
	}
	return mat
}

func (mat Mat4x4) Mul(v float32) Mat4x4 {
	for i := range mat {
		mat[i] *= v
	}
	return mat
}

func (mat Mat4x4) Div(v float32) Mat4x4 {
	for i := range mat {
		mat[i] /= v
	}
	return mat
}

func (mat Mat4x4) Hadamard(other Mat4x4) Mat4x4 {
	for i := range mat {
		mat[i] *= other[i]
	}
	return mat
}

func (mat Mat4x4) Normalize() Mat4x4 {
	return mat.Div(mat.Length())
}
