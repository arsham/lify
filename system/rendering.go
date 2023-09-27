package system

import (
	"fmt"
	"math"

	"github.com/faiface/pixel"
	"golang.org/x/image/colornames"

	"github.com/arsham/neuragene/asset"
	"github.com/arsham/neuragene/component"
	"github.com/arsham/neuragene/entity"
)

// Rendering system renders to the screen.
type Rendering struct {
	target     pixel.Target
	entities   *entity.Manager
	assets     *asset.Manager
	components *component.Manager
	Title      string
	Width      int32
	Height     int32
}

var _ System = (*Rendering)(nil)

func (r *Rendering) String() string { return "Rendering" }

// Setup returns an error if the window manager, entity manager, the asset
// manager, or the components manager is nil.
func (r *Rendering) Setup(c controller) error {
	r.target = c.Target()
	r.entities = c.EntityManager()
	r.assets = c.AssetManager()
	r.components = c.ComponentManager()
	if r.target == nil {
		return fmt.Errorf("%w: window", ErrInvalidArgument)
	}
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

// Process clears up the window and draws all entities on the screen.
func (r *Rendering) Process(component.State, float64) {
	batches := r.assets.Batches()
	for _, b := range batches {
		b.Clear()
	}
	r.renderEntities()
	for _, b := range batches {
		b.Draw(r.target)
	}
}

func (r *Rendering) renderEntities() {
	batches := r.assets.Batches()
	sprites := r.assets.Sprites()
	spriteMap := r.components.Sprite
	posMap := r.components.Position
	r.entities.MapByMask(entity.Positioned|entity.HasTexture, func(e *entity.Entity) {
		sprite := spriteMap[e.ID]
		position := posMap[e.ID]
		sName := sprite.Name
		mat := pixel.IM
		if position.Scale != 0 {
			mat = mat.Scaled(pixel.Vec{}, position.Scale)
		}
		mat = mat.Moved(position.Pos)
		// We don't have the angel, but we have the velocity vector. Since the
		// sprite is positioned 90 degrees to the left, we need to rotate it a
		// bit more.
		angel := position.Velocity.Angle() - math.Pi/2
		mat = mat.Rotated(position.Pos, angel)
		batch := batches[sName]
		sprites[sName].DrawColorMask(batch, mat, colornames.Red)
	})
}
