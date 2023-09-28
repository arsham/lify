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
	return fmt.Sprintf("Vec(%f, %f)", v.X, v.Y)
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
	sine := angle.Sin()
	cosi := angle.Cos()
	return Vec{
		X: v.X*cosi - v.Y*sine,
		Y: v.X*sine + v.Y*cosi,
	}
}

// Angle returns the angle of the vector.
func (v Vec) Angle() Radian {
	return Radian(math.Atan2(v.Y, v.X))
}
