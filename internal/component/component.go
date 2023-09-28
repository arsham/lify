// Package component contains the component data for handling entities' state.
package component

import (
	"github.com/arsham/neuragene/internal/asset"
	"github.com/arsham/neuragene/internal/geom"
)

// Manager manages components for entities. All values are maps of entity IDs
// to their respective components. When an entity is removed, its ID is removed
// from all the maps.
type Manager struct {
	// Position holds the position, scale, and velocity of entities.
	Position map[uint64]*Position
	// Sprite contains the sprite names for renderable entities.
	Sprite map[uint64]*Sprite
	// Lifespan contains the lifespan of entities.
	Lifespan map[uint64]*Lifespan
	// Collision contains the bounding box of entities.
	Collision map[uint64]*Collision
}

// Position component holds the position, scale, velocity vector movement of an
// entity. In order to get the angle of the entity, use the velocity vector.
// The bit mask for this component is the StateMoveEntities constant.
type Position struct {
	// Pos is the centre position of the entity.
	Pos geom.Pos
	// Velocity is the vector movement of the entity. This vector is not a unit
	// vector.
	Velocity geom.Vec
	// Scale is the scale to draw the entity.
	Scale float64
	// Angle is the angle that the entity is facing when it is not moving. This
	// should only be used for rendering.
	Angle geom.Radian
}

// Vec returns the absolute position of the entity.
func (p *Position) Vec() geom.Vec {
	return p.Pos.Resolve()
}

// AddV adds the given vector to the position.
func (p *Position) AddV(pos geom.Vec) {
	p.Pos.Offset = p.Pos.Offset.Add(pos)
}

// Add adds the given values to the position.
func (p *Position) Add(deltaX, deltaY float64) {
	p.Pos.Offset.X += deltaX
	p.Pos.Offset.Y += deltaY
}

// BounceBy causes the position to bounce back if it's out of the given
// rectangle.
func (p *Position) BounceBy(rec geom.Rect) {
	if p.Pos.Offset.X < rec.Min.X {
		p.Pos.Offset.X = rec.Min.X
		p.Velocity.X = -p.Velocity.X
	}
	if p.Pos.Offset.X > rec.Max.X {
		p.Pos.Offset.X = rec.Max.X
		p.Velocity.X = -p.Velocity.X
	}
	if p.Pos.Offset.Y < rec.Min.Y {
		p.Pos.Offset.Y = rec.Min.Y
		p.Velocity.Y = -p.Velocity.Y
	}
	if p.Pos.Offset.Y > rec.Max.Y {
		p.Pos.Offset.Y = rec.Max.Y
		p.Velocity.Y = -p.Velocity.Y
	}
}

// Sprite contains the name of the sprite and the batch object for a sprite.
type Sprite struct {
	Name asset.Name
}

// Lifespan specifies the total amount of frames that the entity should stay
// alive, and the remaining frames.
type Lifespan struct {
	Total     int
	Remaining int
}

// Collision specifies the bounding box in which an entity will collide with
// other entities.
type Collision struct {
	geom.Rect
}

// State is used to identify a system's functionality. At each state, the
// system has a certain behaviour that can be determined by the bit masks based
// on the available constants.
type State uint16

const (
	// StateMoveEntities indicates that the system should move the entities.
	StateMoveEntities State = 1 << iota
	// StateRunning indicates that the system is running.
	StateRunning
	// StateQuit sets the game to a state that causes it to quit.
	StateQuit
	// StateDrawGrids indicates that the system should draw the grid.
	StateDrawGrids
	// StateLimitFPS indicates that the system should limit the FPS.
	StateLimitFPS
	// StateSpawnAnts indicates that the system should spawn ants.
	StateSpawnAnts
	// StatePrintStats indicates that the system should print stats.
	StatePrintStats
	// StateLimitLifespans indicates that the system should limit the lifespan
	// of the entity.
	StateLimitLifespans
	// StateDrawBoundingBoxes indicates that the system should draw the entity's
	// bounding box.
	StateDrawBoundingBoxes
	// StateDrawTextures indicates that the system should draw the entity's
	// texture.
	StateDrawTextures
)
