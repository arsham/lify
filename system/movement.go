package system

import (
	"fmt"

	"github.com/faiface/pixel"

	"github.com/arsham/neuragene/component"
	"github.com/arsham/neuragene/entity"
)

// Movement system handles the Movement of the entity.
type Movement struct {
	entities   *entity.Manager
	components *component.Manager
	target     pixel.Target
	bounds     pixel.Rect
}

var _ System = (*Movement)(nil)

func (m *Movement) String() string { return "Movement" }

// Setup returns an error if the window, the entity manager, the bounds or the
// component manager is nil.
func (m *Movement) Setup(c controller) error {
	m.entities = c.EntityManager()
	m.target = c.Target()
	m.bounds = c.Bounds()
	m.components = c.ComponentManager()
	if m.target == nil {
		return fmt.Errorf("%w: window", ErrInvalidArgument)
	}
	if m.entities == nil {
		return fmt.Errorf("%w: entity manager", ErrInvalidArgument)
	}
	if m.bounds == pixel.ZR {
		return fmt.Errorf("%w: bounds", ErrInvalidArgument)
	}
	if m.components == nil {
		return fmt.Errorf("%w: component manager", ErrInvalidArgument)
	}
	return nil
}

// Process moves the entities if their movement or velocity flags are set.
func (m *Movement) Process(state component.State, dt float64) {
	if !all(state, component.StateRunning) {
		return
	}
	posMap := m.components.Position
	m.entities.MapByMask(entity.Positioned, func(e *entity.Entity) {
		id := e.ID
		position := posMap[id]
		deltaX := position.Velocity.X * dt
		deltaY := position.Velocity.Y * dt
		position.Pos.X += deltaX
		position.Pos.Y += deltaY
		if position.Pos.X > m.bounds.Max.X {
			position.Pos.X = m.bounds.Max.X
			position.Velocity.X = -position.Velocity.X
		}
		if position.Pos.X < m.bounds.Min.X {
			position.Pos.X = m.bounds.Min.X
			position.Velocity.X = -position.Velocity.X
		}
		if position.Pos.Y > m.bounds.Max.Y {
			position.Pos.Y = m.bounds.Max.Y
			position.Velocity.Y = -position.Velocity.Y
		}
		if position.Pos.Y < m.bounds.Min.Y {
			position.Pos.Y = m.bounds.Min.Y
			position.Velocity.Y = -position.Velocity.Y
		}
	})
}
