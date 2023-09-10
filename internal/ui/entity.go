package ui

import (
	"sync/atomic"

	"github.com/disintegration/gift"
	"github.com/oakmound/oak/v4/render"
	"github.com/oakmound/oak/v4/render/mod"
	"github.com/paulmach/orb"
)

var lastID int64

type entity interface {
	String() string
	Resolve()
}

// An EntityView is a view of an entity on the board.
type EntityView struct {
	asset *render.Sprite
	// The entity to be viewed.
	entity entity
	// The location of the entity on the board.
	location orb.Point
}

// NewEntity returns a new instance of the entity view.
func NewEntity(e entity, location orb.Point, asset *render.Sprite) *EntityView {
	return &EntityView{
		entity:   e,
		location: location,
		asset:    asset,
	}
}

// Point returns the location of the entity.
func (e *EntityView) Point() orb.Point {
	return e.location
}

// ID returns a unique ID for the entity.
func (e *EntityView) ID() int64 {
	return atomic.AddInt64(&lastID, 1)
}

// Draw draws the entity on the screen. The viewport is the coordination of the
// top-left and bottom-right corners of the screen.
func (e *EntityView) Draw(screen *render.Sprite, viewport orb.Point) {
	identM := e.asset.Modify(mod.ResizeToFit(64, 64, gift.CubicResampling))
	x := e.location.X() - viewport.X()
	y := e.location.Y() - viewport.Y()
	identM.Draw(screen, x, y)
}
