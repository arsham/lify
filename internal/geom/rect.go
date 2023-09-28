package geom

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
