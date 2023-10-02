package geom

import "math"

// Rect is a rectangle with a Min and Max position.
type Rect struct {
	Min Vec
	Max Vec
}

// ZR is a zero rectangle.
var ZR = Rect{Min: ZV, Max: ZV}

// R returns a new Rect with given the Min and Max coordinates.
func R(minX, minY, maxX, maxY float64) Rect {
	return Rect{
		Min: Vec{minX, minY},
		Max: Vec{maxX, maxY},
	}
}

// W returns the width of the Rect.
func (r Rect) W() float64 {
	return r.Max.X - r.Min.X
}

// H returns the height of the Rect.
func (r Rect) H() float64 {
	return r.Max.Y - r.Min.Y
}

// Centre returns the centre position of the Rect.
func (r Rect) Centre() Vec {
	return Lerp(r.Min, r.Max, 0.5)
}

// Resized returns the Rect resized to the given size while keeping the
// position of the given anchor.
//
// If the anchor is r.Min, then it resizes while keeping the position of the
// lower-left corner.
// If the anchor is r.Max, then it resizes while keeping the position of the
// top-right corner.
// If anchor is r.Center(), then it resizes around the center.
//
// It returns a zero area Rect if the Rect has zero area.
func (r Rect) Resized(anchor, size Vec) Rect {
	if r.W()*r.H() == 0 {
		return Rect{}
	}
	fraction := Vec{size.X / r.W(), size.Y / r.H()}
	return Rect{
		Min: anchor.Add(r.Min.Sub(anchor).ScaledXY(fraction)),
		Max: anchor.Add(r.Max.Sub(anchor).ScaledXY(fraction)),
	}
}

// Moved returns a moved Rect (both Min and Max) by the given delta vector.
func (r Rect) Moved(delta Vec) Rect {
	return Rect{
		Min: r.Min.Add(delta),
		Max: r.Max.Add(delta),
	}
}

// Rotated returns a rotated the rectangle by the given angle.
func (r Rect) Rotated(angle Radian) Rect {
	center := Vec{
		X: (r.Min.X + r.Max.X) / 2,
		Y: (r.Min.Y + r.Max.Y) / 2,
	}
	sin, cos := angle.Sincos()

	return Rect{
		// The new positions of the corners after rotation.
		Min: Vec{
			X: (r.Min.X-center.X)*cos - (r.Min.Y-center.Y)*sin + center.X,
			Y: (r.Min.X-center.X)*sin + (r.Min.Y-center.Y)*cos + center.Y,
		},
		Max: Vec{
			X: (r.Max.X-center.X)*cos - (r.Max.Y-center.Y)*sin + center.X,
			Y: (r.Max.X-center.X)*sin + (r.Max.Y-center.Y)*cos + center.Y,
		},
	}
}

// Eq returns true if both Min and Max vectors are approximately equal.
func (r Rect) Eq(other Rect) bool {
	return r.Min.Eq(other.Min) && r.Max.Eq(other.Max)
}

// Intersects returns true if the Rect intersects with the given Rect s.
func (r Rect) Intersects(s Rect) bool {
	return !(s.Max.X < r.Min.X ||
		s.Min.X > r.Max.X ||
		s.Max.Y < r.Min.Y ||
		s.Min.Y > r.Max.Y)
}

// MinimumTranslationVector returns the minimum translation vector required for
// compensating the given Rect s so that it no longer intersects with this
// Rect. If the Rects don't overlap, this function returns a zero-vector.
func (r Rect) MinimumTranslationVector(other Rect) Vec {
	rW := r.W()
	rH := r.H()
	otherW := other.W()
	otherH := other.H()
	rCentre := r.Centre()
	otherCentre := other.Centre()
	overlapX := (rW+otherW)/2 - math.Abs(rCentre.X-otherCentre.X)
	overlapY := (rH+otherH)/2 - math.Abs(rCentre.Y-otherCentre.Y)

	// Minimum translation vector.
	var mtvX, mtvY float64
	// Determine the direction of the overlap.
	switch {
	case overlapX < overlapY && rCentre.X < otherCentre.X:
		mtvX = -overlapX
	case overlapX < overlapY:
		mtvX = overlapX
	case rCentre.Y < otherCentre.Y:
		mtvY = -overlapY
	default:
		mtvY = overlapY
	}

	return V(mtvX, mtvY)
}

// Contains returns true if the given vector is inside the Rect. Points on the
// edge of the Rect are considered to be inside.
func (r Rect) Contains(v Vec) bool {
	return r.Min.X <= v.X &&
		v.X <= r.Max.X &&
		r.Min.Y <= v.Y &&
		v.Y <= r.Max.Y
}
