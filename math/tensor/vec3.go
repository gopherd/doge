package tensor

import "math"

// Vector3 implements 3d vector
type Vector3 [3]float32

func Vec3(x, y, z float32) Vector3 {
	return Vector3{x, y, z}
}

func (vec Vector3) X() float32 { return vec[0] }
func (vec Vector3) Y() float32 { return vec[1] }
func (vec Vector3) Z() float32 { return vec[2] }

func (vec Vector3) R() float32 { return vec[0] }
func (vec Vector3) G() float32 { return vec[1] }
func (vec Vector3) B() float32 { return vec[2] }

func (vec Vector3) Vec4() Vector4 { return Vec4(vec[0], vec[1], vec[2], 1) }

func (vec Vector3) Sum() float32 {
	return vec[0] + vec[1] + vec[2]
}

func (vec Vector3) Dot(other Vector3) float32 {
	return vec[0]*other[0] + vec[1]*other[1] + vec[2]*other[2]
}

func (vec Vector3) Square() float32 {
	return vec.Dot(vec)
}

func (vec Vector3) Length() float32 {
	return float32(math.Sqrt(float64(vec.Square())))
}

func (vec Vector3) Add(other Vector3) Vector3 {
	return Vec3(vec[0]+other[0], vec[1]+other[1], vec[2]+other[2])
}

func (vec Vector3) Sub(other Vector3) Vector3 {
	return Vec3(vec[0]-other[0], vec[1]-other[1], vec[2]-other[2])
}

func (vec Vector3) Mul(k float32) Vector3 {
	return Vec3(vec[0]*k, vec[1]*k, vec[2]*k)
}

func (vec Vector3) Div(k float32) Vector3 {
	return Vec3(vec[0]/k, vec[1]/k, vec[2]/k)
}

func (vec Vector3) Hadamard(other Vector3) Vector3 {
	return Vec3(vec[0]*other[0], vec[1]*vec[1], vec[2]*other[2])
}

func (vec Vector3) Normalize() Vector3 {
	return vec.Div(vec.Length())
}
