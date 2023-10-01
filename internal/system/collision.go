package system

import (
	"fmt"

	"github.com/arsham/neuragene/internal/component"
	"github.com/arsham/neuragene/internal/entity"
	"github.com/arsham/neuragene/internal/geom"
)

// Collision system handles collision of entities if their flag is set. This
// system should be set after the BoundingBox system otherwise the effects will
// be undesirable.
type Collision struct {
	noDraw
	entitties  *entity.Manager
	components *component.Manager
}

var _ System = (*Collision)(nil)

func (c *Collision) String() string { return "Collision" }

// Setup returns an error if the entity manager is nil.
func (c *Collision) setup(ct controller) error {
	c.entitties = ct.EntityManager()
	c.components = ct.ComponentManager()
	if c.entitties == nil {
		return fmt.Errorf("%w: entity manager", ErrInvalidArgument)
	}
	if c.components == nil {
		return fmt.Errorf("%w: component manager", ErrInvalidArgument)
	}
	return nil
}

func (c *Collision) update(state component.State) error {
	if !all(state, component.StateHandleCollisions, component.StateRunning) {
		return nil
	}

	boundingBoxes := c.components.BoundingBox
	positions := c.components.Position
	c.entitties.MapByMask(entity.Collides, func(e *entity.Entity) {
		id1 := e.ID
		bb1 := boundingBoxes[id1]
		pos1 := positions[id1]
		// Half height and width of the entity so we wouldn't need to calculate
		// them every time.
		bb1H := bb1.H() * pos1.Scale / 2
		bb1W := bb1.W() * pos1.Scale / 2
		c.entitties.MapByMask(entity.Collides|entity.Rigid, func(other *entity.Entity) {
			if e.ID == other.ID {
				return
			}
			id2 := other.ID
			bb2 := boundingBoxes[id2]
			pos2 := positions[id2]
			bb2H := bb2.H() * pos2.Scale / 2
			bb2W := bb2.W() * pos2.Scale / 2

			// If any of the following conditions are true, then the two
			// entities are not colliding.
			if pos1.Pos.Offset.X+bb1W < pos2.Pos.Offset.X-bb2W {
				return
			}
			if pos1.Pos.Offset.X-bb1W > pos2.Pos.Offset.X+bb2W {
				return
			}
			if pos1.Pos.Offset.Y+bb1H < pos2.Pos.Offset.Y-bb2H {
				return
			}
			if pos1.Pos.Offset.Y-bb1H > pos2.Pos.Offset.Y+bb2H {
				return
			}

			r1 := geom.R(
				pos1.Pos.Resolve().X-bb1W,
				pos1.Pos.Resolve().Y-bb1H,
				pos1.Pos.Resolve().X+bb1W,
				pos1.Pos.Resolve().Y+bb1H,
			)
			r2 := geom.R(
				pos2.Pos.Resolve().X-bb2W,
				pos2.Pos.Resolve().Y-bb2H,
				pos2.Pos.Resolve().X+bb2W,
				pos2.Pos.Resolve().Y+bb2H,
			)

			if r1.Intersects(r2) {
				x, y := r1.MinimumTranslationVector(r2).Scaled(0.5).XY()
				pos1.Pos.Offset.X += x
				pos1.Pos.Offset.Y += y
				pos2.Pos.Offset.X -= x
				pos2.Pos.Offset.Y -= y
			}
		})
	})
	return nil
}
