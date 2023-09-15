package ui

import (
	"sync/atomic"

	"github.com/paulmach/orb"
)

var lastID int64

type entity interface {
	String() string
	Resolve()
}

// An EntityView is a view of an entity on the board.
type EntityView struct {
	// The entity to be viewed.
	entity entity
	// The location of the entity on the board.
	location orb.Point
}

// NewEntity returns a new instance of the entity view.
func NewEntity(e entity, location orb.Point) *EntityView {
	return &EntityView{
		entity:   e,
		location: location,
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
