package system

import (
	"fmt"

	"github.com/arsham/neuragene/component"
	"github.com/arsham/neuragene/entity"
	"github.com/faiface/pixel"
)

// Position system handles the Position of the entity. On each frame, it
// calculates the velocity and updates the position.
type Position struct {
	entities   *entity.Manager
	components *component.Manager
	bounds     pixel.Rect
}

var _ System = (*Position)(nil)

func (p *Position) String() string { return "Position" }

// Setup returns an error if the window or the entity manager is nil.
func (p *Position) Setup(c controller) error {
	p.entities = c.EntityManager()
	p.components = c.ComponentManager()
	p.bounds = c.Bounds()
	if p.entities == nil {
		return fmt.Errorf("%w: entity manager", ErrInvalidArgument)
	}
	if p.components == nil {
		return fmt.Errorf("%w: component manager", ErrInvalidArgument)
	}
	return nil
}

// Process moves the entities if their movement or velocity flags are set.
func (p *Position) Process(state component.State, dt float64) {
	if !all(state, component.StateRunning) {
		return
	}
	bounds := p.bounds
	posMap := p.components.Position
	p.entities.MapByMask(entity.Positioned, func(e *entity.Entity) {
		position := posMap[e.ID]
		deltaX := position.Velocity.X * dt
		deltaY := position.Velocity.Y * dt
		position.Pos.X += deltaX
		position.Pos.Y += deltaY
		if position.Pos.X > bounds.Max.X {
			position.Pos.X = bounds.Max.X
			position.Velocity.X = -position.Velocity.X
		}
		if position.Pos.X < bounds.Min.X {
			position.Pos.X = bounds.Min.X
			position.Velocity.X = -position.Velocity.X
		}
		if position.Pos.Y > bounds.Max.Y {
			position.Pos.Y = bounds.Max.Y
			position.Velocity.Y = -position.Velocity.Y
		}
		if position.Pos.Y < bounds.Min.Y {
			position.Pos.Y = bounds.Min.Y
			position.Velocity.Y = -position.Velocity.Y
		}
	})
}
