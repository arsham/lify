package system

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"

	"github.com/arsham/neuragene/internal/asset"
	"github.com/arsham/neuragene/internal/component"
	"github.com/arsham/neuragene/internal/entity"
	"github.com/arsham/neuragene/internal/geom"
)

// BoundingBox system handles drawing of entitties' bounding boxes.
type BoundingBox struct {
	entitties  *entity.Manager
	components *component.Manager
	assets     *asset.Manager
	Colour     color.Color
	canvas     *ebiten.Image
	Size       float64
}

var _ System = (*BoundingBox)(nil)

func (b *BoundingBox) String() string { return "BoundingBox" }

// setup returns an error if the entity manager, the window, the asset manager
// or the component manager is nil.
func (b *BoundingBox) setup(c controller) error {
	b.entitties = c.EntityManager()
	b.assets = c.AssetManager()
	b.components = c.ComponentManager()
	if b.entitties == nil {
		return fmt.Errorf("%w: entity manager", ErrInvalidArgument)
	}
	if b.assets == nil {
		return fmt.Errorf("%w: asset manager", ErrInvalidArgument)
	}
	if b.components == nil {
		return fmt.Errorf("%w: component manager", ErrInvalidArgument)
	}
	if b.Colour == nil {
		b.Colour = colornames.Red
	}
	return nil
}

func (b *BoundingBox) update(state component.State) error {
	if !all(state, component.StateDrawBoundingBoxes) {
		return nil
	}

	b.canvas = ebiten.NewImage(ebiten.WindowSize())
	collisions := b.components.Collision
	positions := b.components.Position
	b.entitties.MapByMask(entity.BoxBounded, func(e *entity.Entity) {
		id := e.ID
		collision := collisions[id]
		position := positions[id]
		collision.Rect = collision.Resized(collision.Centre(), geom.V(position.Scale, position.Scale))
		x, y := position.Vec().XY()
		w, h := collision.Max.XY()
		vector.StrokeRect(b.canvas, float32(x-2*w), float32(y), float32(w*2), float32(h*2), 1, b.Colour, false)
	})
	return nil
}

// Process draws the bounding boxes of the entities.
func (b *BoundingBox) draw(screen *ebiten.Image, state component.State) {
	if !all(state, component.StateDrawBoundingBoxes) {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.ColorScale.ScaleWithColor(b.Colour)
	screen.DrawImage(b.canvas, op)
}
