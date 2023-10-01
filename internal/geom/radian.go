package geom

import (
	"math"
)

// Radian represents a radian value.
type Radian float64

// NewRadian creates a new radian value from a degree value.
func NewRadian(deg float64) Radian {
	return Radian(deg * (math.Pi / 180))
}

// F64 returns the value of the Radian in float64.
func (r Radian) F64() float64 {
	return float64(r)
}

// Degree returns the degree value of the radian.
func (r Radian) Degree() float64 {
	return float64(r) * (180 / math.Pi)
}

// Cos returns the cosine of the radian value.
func (r Radian) Cos() float64 {
	return math.Cos(float64(r))
}

// Sin returns the sine of the radian value.
func (r Radian) Sin() float64 {
	return math.Sin(float64(r))
}

// Sincos returns the Sine and Cosine of the r.
func (r Radian) Sincos() (sin, cos float64) {
	return math.Sincos(float64(r))
}
