// Package geom contains the geometry types and functions.
package geom

import (
	"fmt"
	"math"
)

// Vec is a 2d structure that is used to represent positions, velocities, and
// other kinds numerical values.
type Vec struct {
	X float64
	Y float64
}

// ZV is a zero Vec.
var ZV = Vec{0, 0}

// V returns a new Vec at the given coordinates.
func V(x, y float64) Vec {
	return Vec{x, y}
}

func (v Vec) String() string {
	return fmt.Sprintf("Vec(%.6f, %.6f)", v.X, v.Y)
}

// XY returns the components of the Vec.
func (v Vec) XY() (x, y float64) {
	return v.X, v.Y
}

// IsZero returns true if both components of the Vec are zero.
func (v Vec) IsZero() bool {
	return v.X == 0 && v.Y == 0
}

// Add returns a new Vec that is the sum of the two vectors.
func (v Vec) Add(other Vec) Vec {
	return Vec{
		v.X + other.X,
		v.Y + other.Y,
	}
}

// Sub returns a new Vec that is the difference of the two vectors.
func (v Vec) Sub(other Vec) Vec {
	return Vec{
		v.X - other.X,
		v.Y - other.Y,
	}
}

// Scaled returns the vector u multiplied by c.
func (v Vec) Scaled(c float64) Vec {
	return Vec{v.X * c, v.Y * c}
}

// ScaledXY returns a new Vec that is the product of the two vectors.
func (v Vec) ScaledXY(other Vec) Vec {
	return Vec{v.X * other.X, v.Y * other.Y}
}

// Mul returns a new Vec that is the product of the two vectors.
func (v Vec) Mul(other Vec) Vec {
	return v.ScaledXY(other)
}

// Lerp returns a linear interpolation between vectors a and b.
//
// The linear interpolation is a point along the line between a and b, whose
// position is controlled by the value t. The value t is expected to be between
// 0 and 1.
func Lerp(a, b Vec, t float64) Vec {
	if t == 0 {
		return a
	}
	if t == 1 {
		return b
	}
	return a.Scaled(1 - t).Add(b.Scaled(t))
}

// Rotated returns the vector rotated by the given angle.
func (v Vec) Rotated(angle Radian) Vec {
	sine, cosi := angle.Sincos()
	return Vec{
		X: v.X*cosi - v.Y*sine,
		Y: v.X*sine + v.Y*cosi,
	}
}

// RotatedByVec returns the vector rotated by the given vector's angle.
func (v Vec) RotatedByVec(other Vec) Vec {
	angle := other.Angle()
	return v.Rotated(angle)
}

// Len returns the length of the vector.
func (v Vec) Len() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

// Angle returns the angle of the vector.
func (v Vec) Angle() Radian {
	return Radian(math.Atan2(v.Y, v.X))
}

// RadToVec converts a given angle into a normalised vector that encodes that direction.
func RadToVec(angle Radian) Vec {
	return Vec{X: angle.Cos(), Y: angle.Sin()}
}

// Normalise returns a new normalised Vec based on v.
func (v Vec) Normalise() Vec {
	length := v.Len()
	if length == 0 {
		return v
	}
	return Vec{v.X / length, v.Y / length}
}

// nearlyEqual compares two float64s and returns whether they are equal,
// accounting for rounding errors.At worst, the result is correct to 7
// significant digits.
func nearlyEqual(a, b float64) bool {
	const tolerance = 1e-5
	a = math.Abs(a)
	b = math.Abs(b)
	if a == b {
		// handle infinities.
		return true
	}
	diff := math.Abs(a - b)
	if a*b == 0 {
		// a or b or both are zero, relative error is not meaningful here.
		return diff < tolerance*tolerance
	}
	return diff/(a+b) < tolerance
}

// Eq will compare two vectors and return whether they are equal accounting for
// rounding errors. At worst, the result is correct to 7 significant digits.
func (v Vec) Eq(other Vec) bool {
	return nearlyEqual(v.X, other.X) && nearlyEqual(v.Y, other.Y)
}
