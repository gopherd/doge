package tensor

import "math"

// Vector2 implements 2d vector
type Vector2 [2]float32

func Vec2(x, y float32) Vector2 {
	return Vector2{x, y}
}

func (vec Vector2) X() float32 { return vec[0] }
func (vec Vector2) Y() float32 { return vec[1] }

func (vec Vector2) Vec3() Vector3 { return Vec3(vec[0], vec[1], 0) }
func (vec Vector2) Vec4() Vector4 { return Vec4(vec[0], vec[1], 0, 1) }

func (vec Vector2) Sum() float32 {
	return vec[0] + vec[1]
}

func (vec Vector2) Dot(other Vector2) float32 {
	return vec[0]*other[0] + vec[1]*other[1]
}

func (vec Vector2) Square() float32 {
	return vec.Dot(vec)
}

func (vec Vector2) Length() float32 {
	return float32(math.Sqrt(float64(vec.Square())))
}

func (vec Vector2) Add(other Vector2) Vector2 {
	return Vec2(vec[0]+other[0], vec[1]+other[1])
}

func (vec Vector2) Sub(other Vector2) Vector2 {
	return Vec2(vec[0]-other[0], vec[1]-other[1])
}

func (vec Vector2) Mul(k float32) Vector2 {
	return Vec2(vec[0]*k, vec[1]*k)
}

func (vec Vector2) Div(k float32) Vector2 {
	return Vec2(vec[0]/k, vec[1]/k)
}

func (vec Vector2) Hadamard(other Vector2) Vector2 {
	return Vec2(vec[0]*other[0], vec[1]*vec[1])
}

func (vec Vector2) Normalize() Vector2 {
	return vec.Div(vec.Length())
}
