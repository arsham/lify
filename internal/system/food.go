package system

import (
	"fmt"
	"time"

	"github.com/arsham/neuragene/internal/asset"
	"github.com/arsham/neuragene/internal/component"
	"github.com/arsham/neuragene/internal/entity"
	"github.com/arsham/neuragene/internal/geom"
	"github.com/hajimehoshi/ebiten/v2"
)

// Food system handles the reward that entities get from interaction with food.
type Food struct {
	entities     *entity.Manager
	components   *component.Manager
	assets       *asset.Manager
	sprite       *ebiten.Image
	lastDuration time.Duration
}

var _ System = (*Food)(nil)

func (f *Food) String() string { return "Food" }

// setup returns an error if the entity manager or the asset manager is nil.
func (f *Food) setup(c controller) error {
	f.entities = c.EntityManager()
	f.assets = c.AssetManager()
	f.components = c.ComponentManager()
	if f.entities == nil {
		return fmt.Errorf("%w: entity manager", ErrInvalidArgument)
	}
	if f.assets == nil {
		return fmt.Errorf("%w: asset manager", ErrInvalidArgument)
	}
	if f.components == nil {
		return fmt.Errorf("%w: component manager", ErrInvalidArgument)
	}
	sprite := f.assets.Sprites()[asset.FruitApple]
	if sprite == nil {
		return fmt.Errorf("%w: sprites", ErrNotFound)
	}
	f.sprite = sprite
	return nil
}

// update will spawn food at the cursor position.
func (f *Food) update(state component.State) error {
	if !all(state, component.StateRunning, component.StateSpawnFood) {
		return nil
	}
	x, y := ebiten.CursorPosition()
	err := f.spawnFood(float64(x), float64(y))
	if err != nil {
		return err
	}
	return nil
}

const foodMask = entity.Positioned | entity.BoxBounded | entity.Lifespan

func (f *Food) spawnFood(x, y float64) error {
	food := f.entities.NewEntity(foodMask)
	id := food.ID
	f.components.Position[id] = &component.Position{
		Scale:    1,
		Pos:      geom.P(x, y),
		Velocity: geom.Vec{X: 1, Y: 1},
		Angle:    0,
	}
	f.components.Sprite[id] = &component.Sprite{
		Name: asset.FruitApple,
	}
	f.components.Lifespan[id] = &component.Lifespan{
		Total:     500,
		Remaining: 500,
	}

	b := f.sprite.Bounds()
	bounds := geom.R(float64(b.Min.X), float64(b.Min.Y), float64(b.Max.X), float64(b.Max.Y))

	f.components.Collision[id] = &component.Collision{Rect: bounds}
	return nil
}

func (f *Food) draw(*ebiten.Image, component.State) {}

// avgCalc returns the amount of time it took for the last update.
func (f *Food) avgCalc() time.Duration {
	return f.lastDuration
}
