package geom

import "strings"

// Matrix represents a transformation matrix for a bounding box.
type Matrix struct {
	Corners []Vec
}

// M creates a new Matrix from four corner points.
func M(corners []Vec) *Matrix {
	return &Matrix{Corners: corners}
}

func (m *Matrix) String() string {
	s := make([]string, 0, len(m.Corners)+2)
	s = append(s, "Matrix [")
	for _, corner := range m.Corners {
		s = append(s, "\t"+corner.String())
	}
	s = append(s, "]")
	return strings.Join(s, "\n")
}

// Eq returns true if the two Matrices are equal. Each corner point must be
// equal for the Matrices to be equal.
func (m *Matrix) Eq(other *Matrix) bool {
	if len(m.Corners) != len(other.Corners) {
		return false
	}
	for i := range m.Corners {
		if !m.Corners[i].Eq(other.Corners[i]) {
			return false
		}
	}
	return true
}

// Move moves the Matrix by the given Vec.
func (m *Matrix) Move(v Vec) *Matrix {
	for i := range m.Corners {
		m.Corners[i].X += v.X
		m.Corners[i].Y += v.Y
	}
	return m
}

// Rotate rotates the Matrix by the given angle in radians around its center.
func (m *Matrix) Rotate(angle Radian) *Matrix {
	center := m.Center()
	sin, cos := angle.Sincos()

	for i := range m.Corners {
		// Translate to the origin so the rotation is around the center.
		m.Corners[i].X -= center.X
		m.Corners[i].Y -= center.Y

		x := m.Corners[i].X
		y := m.Corners[i].Y
		m.Corners[i].X = x*cos - y*sin
		m.Corners[i].Y = x*sin + y*cos

		// Translate back.
		m.Corners[i].X += center.X
		m.Corners[i].Y += center.Y
	}
	return m
}

// Center calculates the center point of the Matrix.
func (m *Matrix) Center() Vec {
	var center Vec
	for _, corner := range m.Corners {
		center.X += corner.X
		center.Y += corner.Y
	}
	numCorners := float64(len(m.Corners))
	center.X /= numCorners
	center.Y /= numCorners
	return center
}

// Edges returns all the edges of the Matrix.
func (m *Matrix) Edges() []Rect {
	numCorners := len(m.Corners)
	if numCorners < 2 {
		return nil
	}

	var edges []Rect

	for i := 0; i < numCorners; i++ {
		p1 := m.Corners[i]
		p2 := m.Corners[(i+1)%numCorners]

		edges = append(edges, Rect{Min: p1, Max: p2})
	}

	return edges
}

// Resize resizes the Matrix around the given anchor point.
func (m *Matrix) Resize(anchor, size Vec) *Matrix {
	scaleX := size.X
	scaleY := size.Y
	if scaleX == 0 && scaleY == 0 {
		return m
	}

	for i := range m.Corners {
		// Translate the corner to the anchor point.
		translatedX := m.Corners[i].X - anchor.X
		translatedY := m.Corners[i].Y - anchor.Y

		// Scale the translated coordinates
		scaledX := translatedX * scaleX
		scaledY := translatedY * scaleY

		// Translate back to the original position.
		m.Corners[i].X = scaledX + anchor.X
		m.Corners[i].Y = scaledY + anchor.Y
	}
	return m
}
