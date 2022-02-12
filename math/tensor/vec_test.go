package tensor_test

import (
	"testing"

	"github.com/gopherd/doge/math/tensor"
)

func TestSum(t *testing.T) {
	type testCase struct {
		vec tensor.Vector2[float32]
		sum float32
	}
	for i, tc := range []testCase{
		{tensor.Vec2[float32](0, 0), 0},
		{tensor.Vec2[float32](1, 0), 1},
		{tensor.Vec2[float32](0, 1), 1},
		{tensor.Vec2[float32](1, 1), 2},
		{tensor.Vec2[float32](0.5, 0.5), 1},
	} {
		sum := tc.vec.Sum()
		if sum != tc.sum {
			t.Fatalf("%dth: want %f, got %f", i, tc.sum, sum)
		}
	}
}

func TestAdd(t *testing.T) {
	type testCase struct {
		v1, v2, v3 tensor.Vector2[float32]
	}
	for i, tc := range []testCase{
		{tensor.Vec2[float32](1, 2), tensor.Vec2[float32](3, 4), tensor.Vec2[float32](4, 6)},
	} {
		v := tc.v1.Add(tc.v2)
		if v != tc.v3 {
			t.Fatalf("%dth: want %v, got %v", i, tc.v3, v)
		}
	}
}
