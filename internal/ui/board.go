// Package ui contains logic for interacting and rendering things on the
// screen.
package ui

import (
	"fmt"
	_ "image/png" // Required for loading png files.

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/quadtree"

	"github.com/arsham/lify/internal/config"
)

// Board contains the logic for interacting and rendering everything on the
// board in the game. This is the only object that is allowed to Resolve()
// resources. No other objects is allowed to make that decision, but they can
// cascade the Resolve() if they manage another object themselves.
type Board struct {
	arena  *quadtree.Quadtree
	entity map[int64]*EntityView
	bound  orb.Bound
	x, y   int
}

// NewBoard returns a new instance of the board. It loads all the resources
// into memory and returns an error if any of the resources can't be loaded.
// Note that although the boundary is defined as float64, the locations on the
// board is always int32.
func NewBoard(env *config.Env) *Board {
	bound := orb.Bound{
		Min: orb.Point{0, 0},
		Max: orb.Point{10000, 10000},
	}
	return &Board{
		x:      env.UI.Width,
		y:      env.UI.Height,
		bound:  bound,
		arena:  quadtree.New(bound),
		entity: make(map[int64]*EntityView, 100),
	}
}

// Bound returns the boundary of which the board can handle.
func (b *Board) Bound() orb.Bound {
	return b.bound
}

// Add adds an object to the board. It returns an error if the position is
// outside of the bounds of the board, or another object is in the same
// position.
func (b *Board) Add(o *EntityView) error {
	// func (b *Board) Add(o *EntityView[entity]) error {
	p := b.arena.Find(o.Point())
	if p != nil && p.Point().Equal(o.Point()) {
		return fmt.Errorf("location %v is already occupied", o.Point())
	}
	if err := b.arena.Add(o); err != nil {
		return fmt.Errorf("location %v is outside the board", o.Point())
	}
	b.entity[o.ID()] = o
	return nil
}

// Entities returns all the entities on the board.
func (b *Board) Entities() []*EntityView {
	ret := make([]*EntityView, 0, len(b.entity))
	for _, v := range b.entity {
		ret = append(ret, v)
	}
	return ret
}

// EntitiesIn returns all the entities on the board.
func (b *Board) EntitiesIn(x1, y1, x2, y2 float64) []*EntityView {
	bound := orb.Bound{
		Min: orb.Point{x1, y1},
		Max: orb.Point{x2, y2},
	}

	entities := b.arena.InBound(nil, bound)
	ret := make([]*EntityView, 0, len(entities))
	for _, e := range entities {
		if v, ok := e.(*EntityView); ok {
			ret = append(ret, v)
		}
	}
	return ret
}
