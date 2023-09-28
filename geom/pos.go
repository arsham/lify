package geom

// Pos describes the position of an entity to a base entity. If the base is
// nil, the position is absolute.
type Pos struct {
	Base   *Vec
	Offset Vec
}

// P returns a new Pos pointer.
func P(x, y float64) Pos {
	return Pos{
		Offset: Vec{X: x, Y: y},
	}
}

// Resolve returns the absolute position of the entity.
func (p *Pos) Resolve() Vec {
	if p.Base == nil {
		return p.Offset
	}
	return p.Base.Add(p.Offset)
}
