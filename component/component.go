// Package component contains the component data for handling entities' state.
package component

import "github.com/faiface/pixel"

// Manager manages components for entities. All values are maps of entity IDs
// to their respective components. When an entity is removed, its ID is removed
// from all the maps.
type Manager struct {
	// Position holds the position, scale, and velocity of entities.
	Position map[int64]*Position
}

// Component is used to identify a component.
type Component uint8

const (
	// PositionComponent holds the position, scale, and velocity of an entity.
	PositionComponent Component = iota + 1
)

// Position component holds the position, scale, velocity vector movement of an
// entity. In order to get the angel of the entity, use the velocity vector.
// The bit mask for this component is the StateMoveEntities constant.
type Position struct {
	// Pos is the centre position of the entity.
	Pos pixel.Vec
	// Velocity is the vector movement of the entity. This vector is not a unit
	// vector.
	Velocity pixel.Vec
	// Scale is the scale to draw the entity.
	Scale float64
}

// State is used to identify a system's functionality. At each state, the
// system has a certain behaviour that can be determined by the bit masks based
// on the available constants.
type State uint16

const (
	// StateMoveEntities indicates that the system should move the entities.
	StateMoveEntities State = 1 << iota
)
