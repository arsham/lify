package system

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/arsham/neuragene/internal/component"
	"github.com/arsham/neuragene/internal/entity"
	"github.com/arsham/neuragene/internal/geom"
)

// Position system handles the Position of the entity. On each frame, it
// calculates the velocity and updates the position.
type Position struct {
	noDraw
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
		position.Add(deltaX, deltaY)

		// Preventing the entity from going out of the screen.
		container := geom.R(0, 0, float64(x), float64(y))
		position.BounceBy(container)
	})
	return nil
}
