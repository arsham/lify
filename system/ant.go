package system

import (
	"fmt"
	stdrand "math/rand"

	"github.com/faiface/pixel"

	"github.com/arsham/neuragene/asset"
	"github.com/arsham/neuragene/component"
	"github.com/arsham/neuragene/entity"
)

// Ant spawns ants when required.
type Ant struct {
	rand         *stdrand.Rand
	entities     *entity.Manager
	assets       *asset.Manager
	sprite       *pixel.Sprite
	components   *component.Manager
	MinVelocity  float64
	MaxVelocity  float64
	Seed         int64
	lastSpawn    int64
	lastFrame    int64
	MutationRate int
}

var _ System = (*Ant)(nil)

func (a *Ant) String() string { return "Ant" }

// Setup returns an error if the entity manager or the asset manager is nil.
func (a *Ant) Setup(c controller) error {
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
		a.MinVelocity = -200
	}
	if a.MaxVelocity == 0 {
		a.MaxVelocity = 200
	}
	return nil
}

const antMask = entity.Positioned | entity.Lifespan

// Process spawns an ant every 100 frames.
// nolint:unparam // this is the expected behaviour.
func (a *Ant) Process(state component.State, _ float64) {
	if !all(state, component.StateSpawnAnts, component.StateRunning) {
		return
	}
	a.lastFrame++
	diff := a.lastFrame - a.lastSpawn
	if diff > 3 {
		a.spawnAnt()
	}
	posMap := a.components.Position
	a.entities.MapByMask(antMask, func(e *entity.Entity) {
		position := posMap[e.ID]
		xScale := float64(10)
		if a.rand.Intn(100) > 50 {
			xScale = -10
		}
		position.Velocity.X += a.rand.Float64() * xScale
		yScale := float64(10)
		if a.rand.Intn(100) > 50 {
			yScale = -10
		}
		position.Velocity.Y += a.rand.Float64() * yScale
	})
}

func (a *Ant) spawnAnt() {
	ant := a.entities.NewEntity(antMask)
	x := a.rand.Float64()*(a.MaxVelocity-a.MinVelocity) + a.MinVelocity
	y := a.rand.Float64()*(a.MaxVelocity-a.MinVelocity) + a.MinVelocity
	scale := 0.6
	id := ant.ID
	a.components.Position[id] = &component.Position{
		Scale:    scale,
		Pos:      pixel.Vec{X: 500, Y: 500},
		Velocity: pixel.Vec{X: x, Y: y},
	}
	a.components.Sprite[id] = &component.Sprite{
		Name: asset.Ant,
	}
	a.components.Lifespan[id] = &component.Lifespan{
		Total:     500,
		Remaining: 500,
	}
	a.lastSpawn = a.lastFrame
}
