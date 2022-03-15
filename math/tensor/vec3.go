package tensor

import (
	"fmt"
	"math"

	"github.com/gopherd/doge/constraints"
	"github.com/gopherd/doge/operator"
)

// Vector3 implements 3d vector
type Vector3[T constraints.SignedReal] [3]T

func Vec3[T constraints.SignedReal](x, y, z T) Vector3[T] {
	return Vector3[T]{x, y, z}
}

func (vec Vector3[T]) String() string {
	return fmt.Sprintf("(%v,%v,%v)", vec[0], vec[1], vec[2])
}

func (vec Vector3[T]) Dim() int        { return 3 }
func (vec Vector3[T]) Get(i int) T     { return operator.If(i < len(vec), vec[i], 0) }
func (vec *Vector3[T]) Set(i int, v T) { vec[i] = v }

func (vec Vector3[T]) X() T { return vec[0] }
func (vec Vector3[T]) Y() T { return vec[1] }
func (vec Vector3[T]) Z() T { return vec[2] }

func (vec *Vector3[T]) SetX(x T) { vec[0] = x }
func (vec *Vector3[T]) SetY(y T) { vec[1] = y }
func (vec *Vector3[T]) SetZ(z T) { vec[2] = z }

func (vec *Vector3[T]) SetElements(x, y, z T) {
	(*vec)[0], (*vec)[1], (*vec)[2] = x, y, z
}

func (vec *Vector3[T]) Copy(other Vector3[T]) {
	(*vec)[0], (*vec)[1], (*vec)[2] = other[0], other[1], other[2]
}

func (vec Vector3[T]) Vec4() Vector4[T] { return Vec4(vec[0], vec[1], vec[2], 1) }

func (vec Vector3[T]) Sum() T {
	return vec[0] + vec[1] + vec[2]
}

func (vec Vector3[T]) Dot(other Vector3[T]) T {
	return vec[0]*other[0] + vec[1]*other[1] + vec[2]*other[2]
}

func (vec Vector3[T]) SquaredLength() T {
	return vec.Dot(vec)
}

func (vec Vector3[T]) Length() T {
	return T(math.Sqrt(float64(vec.SquaredLength())))
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

func (vec Vector3[T]) Cross(other Vector3[T]) Vector3[T] {
	var x1, y1, z1 = vec.X(), vec.Y(), vec.Z()
	var x2, y2, z2 = other.X(), other.Y(), other.Z()
	return Vec3(y1*z2-y2*z1, x2*z1-x1*z2, x1*y2-x2*y1)
}
