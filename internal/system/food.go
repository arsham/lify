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
	lastDuration time.Duration
}

var _ System = (*Food)(nil)

func (a *Food) String() string { return "Food" }

// setup returns an error if the entity manager or the asset manager is nil.
func (a *Food) setup(c controller) error {
	a.entities = c.EntityManager()
	a.assets = c.AssetManager()
	a.components = c.ComponentManager()
	if a.entities == nil {
		return fmt.Errorf("%w: entity manager", ErrInvalidArgument)
	}
	if a.assets == nil {
		return fmt.Errorf("%w: asset manager", ErrInvalidArgument)
	}
	if a.components == nil {
		return fmt.Errorf("%w: component manager", ErrInvalidArgument)
	}
	sprite := a.assets.Sprites()[asset.FruitApple]
	if sprite == nil {
		return fmt.Errorf("%w: sprites", ErrNotFound)
	}
	a.sprite = sprite
	return nil
}

// update will spawn food at the cursor position.
func (a *Food) update(state component.State) error {
	if !all(state, component.StateRunning, component.StateSpawnFood) {
		return nil
	}
	x, y := ebiten.CursorPosition()
	err := a.spawnFood(float64(x), float64(y))
	if err != nil {
		return err
	}
	return nil
}

const foodMask = entity.Positioned | entity.BoxBounded | entity.Lifespan

func (a *Food) spawnFood(x, y float64) error {
	food := a.entities.NewEntity(foodMask)
	id := food.ID
	a.components.Position[id] = &component.Position{
		Scale:    1,
		Pos:      geom.P(x, y),
		Velocity: geom.Vec{X: 1, Y: 1},
		Angle:    0,
	}
	a.components.Sprite[id] = &component.Sprite{
		Name: asset.FruitApple,
	}
	a.components.Lifespan[id] = &component.Lifespan{
		Total:     500,
		Remaining: 500,
	}

	b := a.sprite.Bounds()
	bounds := geom.R(float64(b.Min.X), float64(b.Min.Y), float64(b.Max.X), float64(b.Max.Y))

	a.components.Collision[id] = &component.Collision{Rect: bounds}
	return nil
}

func (f *Food) draw(*ebiten.Image, component.State) {}

// avgCalc returns the amount of time it took for the last update.
func (f *Food) avgCalc() time.Duration {
	return f.lastDuration
}
