package system

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"

	"github.com/arsham/neuragene/asset"
	"github.com/arsham/neuragene/component"
	"github.com/arsham/neuragene/entity"
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
	sprites := r.assets.Sprites()
	spriteMap := r.components.Sprite
	posMap := r.components.Position
	r.entities.MapByMask(entity.Positioned|entity.HasTexture, func(e *entity.Entity) {
		sprite := spriteMap[e.ID]
		position := posMap[e.ID]
		sName := sprite.Name
		options := &ebiten.DrawImageOptions{}
		if position.Scale != 0 {
			options.GeoM.Scale(position.Scale, position.Scale)
		}

		// We don't have the angel, but we have the velocity vector. Since the
		// sprite is positioned 90 degrees to the left, we need to rotate it a
		// bit more.
		angel := math.Atan2(position.Velocity.Y, position.Velocity.X) + math.Pi/2
		options.GeoM.Rotate(angel)
		options.GeoM.Translate(position.Pos.X, position.Pos.Y)
		options.ColorScale.ScaleWithColor(colornames.Red)

		screen.DrawImage(sprites[sName], options)
	})
}
