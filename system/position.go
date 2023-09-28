package system

import (
	"fmt"

	"github.com/arsham/neuragene/component"
	"github.com/arsham/neuragene/entity"
	"github.com/hajimehoshi/ebiten/v2"
)

// Position system handles the Position of the entity. On each frame, it
// calculates the velocity and updates the position.
type Position struct {
	entities   *entity.Manager
	components *component.Manager
	controller controller
}

var _ System = (*Position)(nil)

func (p *Position) String() string { return "Position" }

// setup returns an error if the window or the entity manager is nil.
func (p *Position) setup(c controller) error {
	p.entities = c.EntityManager()
	p.components = c.ComponentManager()
	p.controller = c
	if p.entities == nil {
		return fmt.Errorf("%w: entity manager", ErrInvalidArgument)
	}
	if p.components == nil {
		return fmt.Errorf("%w: component manager", ErrInvalidArgument)
	}
	return nil
}

// update moves the entities if their movement or velocity flags are set.
func (p *Position) update(state component.State) error {
	if !all(state, component.StateRunning) {
		return nil
	}
	// Velocity is the vector movement of the entity with speed of 100 pixels
	// per frame. Since this method is called TPS times per second, we need to
	// calculate the position of the entity based on the time passed (1/TPS).
	// An entity can move diagonally, so we need to account for the angle.
	x, y := ebiten.WindowSize()
	posMap := p.components.Position
	p.entities.MapByMask(entity.Positioned, func(e *entity.Entity) {
		position := posMap[e.ID]
		deltaX := position.Velocity.X / 100
		deltaY := position.Velocity.Y / 100
		position.Pos.X += deltaX
		position.Pos.Y += deltaY

		// Preventing the entity from going out of the screen.
		if position.Pos.X > float64(x) {
			position.Pos.X = float64(x)
			position.Velocity.X = -position.Velocity.X
		}
		if position.Pos.X < 0 {
			position.Pos.X = 0
			position.Velocity.X = -position.Velocity.X
		}
		if position.Pos.Y > float64(y) {
			position.Pos.Y = float64(y)
			position.Velocity.Y = -position.Velocity.Y
		}
		if position.Pos.Y < 0 {
			position.Pos.Y = 0
			position.Velocity.Y = -position.Velocity.Y
		}
	})
	return nil
}

func (p *Position) draw(*ebiten.Image, component.State) {}
