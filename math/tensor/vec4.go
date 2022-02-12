package tensor

import (
	"math"

	"github.com/gopherd/doge/math/mathutil"
)

// Vector4 implements 4d vector
type Vector4[T mathutil.Real] [4]T

func Vec4[T mathutil.Real](x, y, z, w T) Vector4[T] {
	return Vector4[T]{x, y, z, w}
}

func (vec Vector4[T]) X() T { return vec[0] }
func (vec Vector4[T]) Y() T { return vec[1] }
func (vec Vector4[T]) Z() T { return vec[2] }
func (vec Vector4[T]) W() T { return vec[3] }

func (vec Vector4[T]) R() T { return vec[0] }
func (vec Vector4[T]) G() T { return vec[1] }
func (vec Vector4[T]) B() T { return vec[2] }
func (vec Vector4[T]) A() T { return vec[3] }

func (vec Vector4[T]) Vec3() Vector3[T] {
	if vec[3] == 0 {
		return Vec3(vec[0], vec[1], vec[2])
	}
	return Vec3(vec[0]/vec[3], vec[1]/vec[3], vec[2]/vec[3])
}

func (vec Vector4[T]) Sum() T {
	return vec[0] + vec[1] + vec[2] + vec[3]
}

func (vec Vector4[T]) Dot(other Vector4[T]) T {
	return vec[0]*other[0] + vec[1]*other[1] + vec[2]*other[2] + vec[3]*other[3]
}

func (vec Vector4[T]) Square() T {
	return vec.Dot(vec)
}

func (vec Vector4[T]) Length() T {
	return T(math.Sqrt(float64(vec.Square())))
}

func (vec Vector4[T]) Add(other Vector4[T]) Vector4[T] {
	return Vec4(vec[0]+other[0], vec[1]+other[1], vec[2]+other[2], vec[3]+other[3])
}

func (vec Vector4[T]) Sub(other Vector4[T]) Vector4[T] {
	return Vec4(vec[0]-other[0], vec[1]-other[1], vec[2]-other[2], vec[3]-other[3])
}

func (vec Vector4[T]) Mul(k T) Vector4[T] {
	return Vec4(vec[0]*k, vec[1]*k, vec[2]*k, vec[3]*k)
}

func (vec Vector4[T]) Div(k T) Vector4[T] {
	return Vec4(vec[0]/k, vec[1]/k, vec[2]/k, vec[3]/k)
}

func (vec Vector4[T]) Hadamard(other Vector4[T]) Vector4[T] {
	return Vec4(vec[0]*other[0], vec[1]*vec[1], vec[2]*other[2], vec[3]*other[3])
}

func (vec Vector4[T]) Normalize() Vector4[T] {
	return vec.Div(vec.Length())
}
