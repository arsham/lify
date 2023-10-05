package system

import (
	"fmt"
	"image/color"
	"runtime"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"

	"github.com/arsham/neuragene/internal/component"
	"github.com/arsham/neuragene/internal/entity"
	"github.com/arsham/neuragene/internal/geom"
	"github.com/arsham/neuragene/internal/quadtree"
)

// Collision system handles collision of entities if their flag is set. This
// system should be set after the BoundingBox system otherwise the effects will
// be undesirable.
type Collision struct {
	noDraw
	entitties    *entity.Manager
	components   *component.Manager
	qTree        *quadtree.QuadTree[uint64]
	Colour       color.Color
	lastDuration time.Duration
	Capacity     uint
	Workers      uint
}

var _ System = (*Collision)(nil)

func (c *Collision) String() string { return "Collision" }

// Setup returns an error if the entity manager is nil.
func (c *Collision) setup(ct controller) error {
	c.entitties = ct.EntityManager()
	c.components = ct.ComponentManager()
	if c.entitties == nil {
		return fmt.Errorf("%w: entity manager", ErrInvalidArgument)
	}
	if c.components == nil {
		return fmt.Errorf("%w: component manager", ErrInvalidArgument)
	}
	if c.Colour == nil {
		c.Colour = colornames.Red
	}
	if c.Capacity == 0 {
		c.Capacity = 20
	}
	if c.Workers == 0 {
		c.Workers = uint(runtime.NumCPU())
	}
	return nil
}

func (c *Collision) update(state component.State) error {
	started := time.Now()
	defer func() {
		c.lastDuration = time.Since(started)
	}()
	if !all(state, component.StateHandleCollisions, component.StateRunning) {
		return nil
	}

	positions := c.components.Position
	maxX, maxY := ebiten.WindowSize()
	bounds := quadtree.NewBounds(0, 0, float64(maxX), float64(maxY))
	c.qTree = quadtree.NewQuadTree[uint64](bounds, c.Capacity, 0)
	defer c.qTree.Free()
	c.entitties.MapByMask(entity.Collides|entity.Rigid, func(e *entity.Entity) {
		id := e.ID
		pos := positions[id]
		point := quadtree.P(
			id,
			geom.V(
				pos.Pos.Resolve().X,
				pos.Pos.Resolve().Y,
			),
		)
		c.qTree.Insert(point)
	})

	var wg sync.WaitGroup
	count := c.entitties.Len()
	checkCh := make(chan *entity.Entity, count)
	wg.Add(int(c.Workers))
	for i := 0; i < int(c.Workers); i++ {
		go func() {
			defer wg.Done()
			for e := range checkCh {
				c.entityCollisions(e)
			}
		}()
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		c.entitties.MapByMask(entity.Collides, func(e *entity.Entity) {
			checkCh <- e
		})
		close(checkCh)
	}()
	wg.Wait()
	return nil
}

func (c *Collision) entityCollisions(e *entity.Entity) {
	id1 := e.ID
	boundingBoxes := c.components.BoundingBox
	positions := c.components.Position
	bb1 := boundingBoxes[id1]
	pos1 := positions[id1]
	// Half height and width of the entity so we wouldn't need to calculate
	// them every time.
	bb1H := bb1.H() * pos1.Scale / 2
	bb1W := bb1.W() * pos1.Scale / 2
	x, y := pos1.Pos.Resolve().XY()
	bounds := geom.R(x-bb1W, y-bb1H, x+bb1W, y+bb1H)

	points := quadtree.Q[uint64]()
	defer points.Free()
	c.qTree.Query(bounds, points)
	for i := range points.Points {
		id2 := points.Points[i].Data
		if id1 == id2 {
			return
		}
		bb2 := boundingBoxes[id2]
		pos2 := positions[id2]
		bb2H := bb2.H() * pos2.Scale / 2
		bb2W := bb2.W() * pos2.Scale / 2

		// If any of the following conditions are true, then the two
		// entities are not colliding.
		if pos1.Pos.Offset.X+bb1W < pos2.Pos.Offset.X-bb2W {
			return
		}
		if pos1.Pos.Offset.X-bb1W > pos2.Pos.Offset.X+bb2W {
			return
		}
		if pos1.Pos.Offset.Y+bb1H < pos2.Pos.Offset.Y-bb2H {
			return
		}
		if pos1.Pos.Offset.Y-bb1H > pos2.Pos.Offset.Y+bb2H {
			return
		}

		r1 := geom.R(
			pos1.Pos.Resolve().X-bb1W,
			pos1.Pos.Resolve().Y-bb1H,
			pos1.Pos.Resolve().X+bb1W,
			pos1.Pos.Resolve().Y+bb1H,
		)
		r2 := geom.R(
			pos2.Pos.Resolve().X-bb2W,
			pos2.Pos.Resolve().Y-bb2H,
			pos2.Pos.Resolve().X+bb2W,
			pos2.Pos.Resolve().Y+bb2H,
		)

		if r1.Intersects(r2) {
			x, y := r1.MinimumTranslationVector(r2).Scaled(0.5).XY()
			pos1.Pos.Offset.X += x
			pos1.Pos.Offset.Y += y
			pos2.Pos.Offset.X -= x
			pos2.Pos.Offset.Y -= y
		}
	}
}

// avgCalc returns the amount of time it took for the last update.
func (c *Collision) avgCalc() time.Duration {
	return c.lastDuration
}
