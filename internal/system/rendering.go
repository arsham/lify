package system

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"

	"github.com/arsham/neuragene/internal/asset"
	"github.com/arsham/neuragene/internal/component"
	"github.com/arsham/neuragene/internal/entity"
)

// Rendering system renders to the screen.
type Rendering struct {
	entities   *entity.Manager
	assets     *asset.Manager
	components *component.Manager
	Title      string
	Width      int32
	Height     int32
}

var _ System = (*Rendering)(nil)

func (r *Rendering) String() string { return "Rendering" }

// setup returns an error if the window manager, entity manager, the asset
// manager, or the components manager is nil.
func (r *Rendering) setup(c controller) error {
	r.entities = c.EntityManager()
	r.assets = c.AssetManager()
	r.components = c.ComponentManager()
	if r.entities == nil {
		return fmt.Errorf("%w: entity manager", ErrInvalidArgument)
	}
	if r.assets == nil {
		return fmt.Errorf("%w: asset manager", ErrInvalidArgument)
	}
	if r.components == nil {
		return fmt.Errorf("%w: component manager", ErrInvalidArgument)
	}
	return nil
}

func (*Rendering) update(component.State) error { return nil }

// draw clears up the window and draws all entities on the screen.
func (r *Rendering) draw(screen *ebiten.Image, state component.State) {
	if !all(state, component.StateDrawTextures) {
		return
	}
	assets := r.assets.Sprites()
	sprites := r.components.Sprite
	positions := r.components.Position
	boundingBoxes := r.components.BoundingBox
	r.entities.MapByMask(entity.Positioned|entity.HasTexture, func(e *entity.Entity) {
		sprite := sprites[e.ID]
		position := positions[e.ID]
		boundingBox := boundingBoxes[e.ID]
		options := &ebiten.DrawImageOptions{}

		r := boundingBox.Rect
		// Move the centre point to the top left corner so the rotation doesn't
		// look wonky.
		options.GeoM.Translate(-r.W()/2, -r.H()/2)

		if position.Scale != 0 {
			options.GeoM.Scale(position.Scale, position.Scale)
		}

		// We don't have the angle, but we have the velocity vector. Since the
		// sprite is positioned 90 degrees to the left, we need to rotate it a
		// bit more.
		angle := position.Angle
		if !position.Velocity.IsZero() {
			angle = position.Velocity.Angle() + math.Pi/2
		}
		options.GeoM.Rotate(angle.F64())

		// Move the image to the screen's centre.
		options.GeoM.Translate(r.W()/2, r.H()/2)

		options.GeoM.Translate(position.Vec().XY())
		options.ColorScale.ScaleWithColor(colornames.Red)

		img := assets[sprite.Name]
		screen.DrawImage(img, options)
	})
}
