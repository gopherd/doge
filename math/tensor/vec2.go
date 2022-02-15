package tensor

import (
	"math"

	"github.com/gopherd/doge/math/mathutil"
)

// Vector2 implements 2d vector
type Vector2[T mathutil.Real] [2]T

func Vec2[T mathutil.Real](x, y T) Vector2[T] {
	return Vector2[T]{x, y}
}

func (vec Vector2[T]) X() T { return vec[0] }
func (vec Vector2[T]) Y() T { return vec[1] }

func (vec *Vector2[T]) SetX(x T) { vec[0] = x }
func (vec *Vector2[T]) SetY(y T) { vec[1] = y }

func (vec *Vector2[T]) Set(x, y T) {
	(*vec)[0], (*vec)[1] = x, y
}

func (vec *Vector2[T]) Copy(other Vector2[T]) {
	(*vec)[0], (*vec)[1] = other[0], other[1]
}

func (vec Vector2[T]) Vec3() Vector3[T] { return Vec3(vec[0], vec[1], 0) }
func (vec Vector2[T]) Vec4() Vector4[T] { return Vec4(vec[0], vec[1], 0, 1) }

func (vec Vector2[T]) Sum() T {
	return vec[0] + vec[1]
}

func (vec Vector2[T]) Dot(other Vector2[T]) T {
	return vec[0]*other[0] + vec[1]*other[1]
}

func (vec Vector2[T]) Square() T {
	return vec.Dot(vec)
}

func (vec Vector2[T]) Length() T {
	return T(math.Sqrt(float64(vec.Square())))
}

func (vec Vector2[T]) Add(other Vector2[T]) Vector2[T] {
	return Vec2(vec[0]+other[0], vec[1]+other[1])
}

func (vec Vector2[T]) Sub(other Vector2[T]) Vector2[T] {
	return Vec2(vec[0]-other[0], vec[1]-other[1])
}

func (vec Vector2[T]) Mul(k T) Vector2[T] {
	return Vec2(vec[0]*k, vec[1]*k)
}

func (vec Vector2[T]) Div(k T) Vector2[T] {
	return Vec2(vec[0]/k, vec[1]/k)
}

func (vec Vector2[T]) Hadamard(other Vector2[T]) Vector2[T] {
	return Vec2(vec[0]*other[0], vec[1]*vec[1])
}

func (vec Vector2[T]) Normalize() Vector2[T] {
	return vec.Div(vec.Length())
}

func (vec Vector2[T]) Cross(other Vector2[T]) T {
	return vec.X()*other.Y() - vec.Y()*other.X()
}
