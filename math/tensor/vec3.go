package tensor

import (
	"math"

	"github.com/gopherd/doge/math/mathutil"
)

// Vector3 implements 3d vector
type Vector3[T mathutil.Real] [3]T

func Vec3[T mathutil.Real](x, y, z T) Vector3[T] {
	return Vector3[T]{x, y, z}
}

func (vec Vector3[T]) X() T { return vec[0] }
func (vec Vector3[T]) Y() T { return vec[1] }
func (vec Vector3[T]) Z() T { return vec[2] }

func (vec Vector3[T]) R() T { return vec[0] }
func (vec Vector3[T]) G() T { return vec[1] }
func (vec Vector3[T]) B() T { return vec[2] }

func (vec Vector3[T]) Vec4() Vector4[T] { return Vec4(vec[0], vec[1], vec[2], 1) }

func (vec Vector3[T]) Sum() T {
	return vec[0] + vec[1] + vec[2]
}

func (vec Vector3[T]) Dot(other Vector3[T]) T {
	return vec[0]*other[0] + vec[1]*other[1] + vec[2]*other[2]
}

func (vec Vector3[T]) Square() T {
	return vec.Dot(vec)
}

func (vec Vector3[T]) Length() T {
	return T(math.Sqrt(float64(vec.Square())))
}

func (vec Vector3[T]) Add(other Vector3[T]) Vector3[T] {
	return Vec3(vec[0]+other[0], vec[1]+other[1], vec[2]+other[2])
}

func (vec Vector3[T]) Sub(other Vector3[T]) Vector3[T] {
	return Vec3(vec[0]-other[0], vec[1]-other[1], vec[2]-other[2])
}

func (vec Vector3[T]) Mul(k T) Vector3[T] {
	return Vec3(vec[0]*k, vec[1]*k, vec[2]*k)
}

func (vec Vector3[T]) Div(k T) Vector3[T] {
	return Vec3(vec[0]/k, vec[1]/k, vec[2]/k)
}

func (vec Vector3[T]) Hadamard(other Vector3[T]) Vector3[T] {
	return Vec3(vec[0]*other[0], vec[1]*vec[1], vec[2]*other[2])
}

func (vec Vector3[T]) Normalize() Vector3[T] {
	return vec.Div(vec.Length())
}
