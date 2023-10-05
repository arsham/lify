package system

import (
	"fmt"

	"github.com/arsham/neuragene/internal/asset"
	"github.com/arsham/neuragene/internal/component"
	"github.com/arsham/neuragene/internal/entity"
	"github.com/arsham/neuragene/internal/geom"
	"github.com/hajimehoshi/ebiten/v2"
)

// Food system handles the reward that entities get from interaction with food.
type Food struct {
	entities   *entity.Manager
	components *component.Manager
	assets     *asset.Manager
	sprite     *ebiten.Image
}

var _ System = (*Food)(nil)

func (a *Food) String() string { return "Food" }

// setup returns an error if the entity manager or the asset manager is nil.
func (a *Food) setup(c controller) error {
	a.entities = c.EntityManager()
	a.assets = c.AssetManager()
	a.components = c.ComponentManager()
	fmt.Print("FRUITY")
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
	if !all(state, component.StateSpawnFood, component.StateRunning) {
		return nil
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButton0) {
		x, y := ebiten.CursorPosition()
		err := a.spawnFood(float64(x), float64(y))
		if err != nil {
			return err
		}
	}
	return nil
}

const foodMask = entity.Positioned | entity.BoxBounded | entity.Lifespan

func (a *Food) spawnFood(x, y float64) error {
	if !all(component.StateRunning) {
		return nil
	}
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

func (*Food) draw(*ebiten.Image, component.State) {}
