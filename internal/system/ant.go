package system

import (
	"fmt"
	stdrand "math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/arsham/neuragene/internal/asset"
	"github.com/arsham/neuragene/internal/component"
	"github.com/arsham/neuragene/internal/entity"
	"github.com/arsham/neuragene/internal/geom"
)

// Ant spawns ants when required.
type Ant struct {
	noDraw
	rand         *stdrand.Rand
	entities     *entity.Manager
	assets       *asset.Manager
	sprite       *ebiten.Image
	components   *component.Manager
	lastDuration time.Duration
	MinVelocity  float64
	MaxVelocity  float64
	Seed         int64
	lastSpawn    int64
	lastFrame    int64
	MutationRate int
}

var _ System = (*Ant)(nil)

func (a *Ant) String() string { return "Ant" }

// setup returns an error if the entity manager or the asset manager is nil.
func (a *Ant) setup(c controller) error {
	a.rand = stdrand.New(stdrand.NewSource(a.Seed))
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
	a.sprite = a.assets.Sprites()[asset.Ant]
	if a.MinVelocity == 0 {
		a.MinVelocity = -400
	}
	if a.MaxVelocity == 0 {
		a.MaxVelocity = 400
	}
	return nil
}

const antMask = entity.Positioned | entity.Lifespan | entity.BoxBounded | entity.Collides

// update spawns an ant every 100 frames.
func (a *Ant) update(state component.State) error {
	started := time.Now()
	defer func() {
		a.lastDuration = time.Since(started)
	}()
	if !all(state, component.StateSpawnAnts, component.StateRunning) {
		return nil
	}
	a.lastFrame++
	diff := a.lastFrame - a.lastSpawn
	a.spawnAnt()
	if diff > 30 {
		a.lastSpawn = a.lastFrame
		posMap := a.components.Position
		a.entities.MapByMask(antMask, func(e *entity.Entity) {
			position := posMap[e.ID]
			coef := float64(1)
			if a.rand.Intn(100) > 50 {
				coef = -1
			}
			angle := coef * float64(a.rand.Intn(10))
			position.Velocity = position.Velocity.Rotated(geom.NewRadian(angle))
		})
	}
	return nil
}

func (a *Ant) spawnAnt() {
	ant := a.entities.NewEntity(antMask)
	x := a.rand.Float64()*(a.MaxVelocity-a.MinVelocity) + a.MinVelocity
	y := a.rand.Float64()*(a.MaxVelocity-a.MinVelocity) + a.MinVelocity
	scale := 0.6
	id := ant.ID
	a.components.Position[id] = &component.Position{
		Scale:    scale,
		Pos:      geom.P(float64(a.rand.Intn(500)), float64(a.rand.Intn(500))),
		Velocity: geom.Vec{X: x, Y: y},
		Angle:    geom.NewRadian(float64(a.rand.Intn(360))),
	}
	a.components.Sprite[id] = &component.Sprite{
		Name: asset.Ant,
	}
	a.components.Lifespan[id] = &component.Lifespan{
		Total:     500,
		Remaining: 500,
	}

	b := a.sprite.Bounds()
	bounds := geom.R(float64(b.Min.X), float64(b.Min.Y), float64(b.Max.X), float64(b.Max.Y))

	a.components.BoundingBox[id] = &component.BoundingBox{Rect: bounds}
}

// avgCalc returns the amount of time it took for the last update.
func (a *Ant) avgCalc() time.Duration {
	return a.lastDuration
}
