package system

import (
	"fmt"
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"golang.org/x/image/colornames"

	"github.com/arsham/neuragene/asset"
	"github.com/arsham/neuragene/component"
	"github.com/arsham/neuragene/entity"
)

// BoundingBox system handles drawing of entitties' bounding boxes.
type BoundingBox struct {
	entitties  *entity.Manager
	components *component.Manager
	target     pixel.Target
	assets     *asset.Manager
	Colour     color.Color
	Size       float64
}

var _ System = (*BoundingBox)(nil)

func (b *BoundingBox) String() string { return "BoundingBox" }

// Setup returns an error if the entity manager, the window, the asset manager
// or the component manager is nil.
func (b *BoundingBox) Setup(ct controller) error {
	b.entitties = ct.EntityManager()
	b.target = ct.Target()
	b.assets = ct.AssetManager()
	b.components = ct.ComponentManager()
	if b.entitties == nil {
		return fmt.Errorf("%w: entity manager", ErrInvalidArgument)
	}
	if b.target == nil {
		return fmt.Errorf("%w: window", ErrInvalidArgument)
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

// Process draws the bounding boxes of the entities.
func (b *BoundingBox) Process(state component.State, _ float64) {
	if !all(state, component.StateDrawBoundingBoxes) {
		return
	}
	imd := imdraw.New(nil)
	imd.EndShape = imdraw.RoundEndShape
	imd.Color = b.Colour
	collisions := b.components.Collision
	positions := b.components.Position
	b.entitties.MapByMask(entity.BoxBounded, func(e *entity.Entity) {
		id := e.ID
		collision := collisions[id]
		position := positions[id]
		topLeft := pixel.V(
			collision.TopLeft.X+position.Pos.X,
			collision.TopLeft.Y+position.Pos.Y,
		)
		bottomRight := pixel.V(
			collision.BottomRight.X+position.Pos.X,
			collision.BottomRight.Y+position.Pos.Y,
		)
		imd.Push(topLeft, bottomRight)
		imd.Rectangle(b.Size)
	})
	imd.Draw(b.target)
}
