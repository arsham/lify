package system

import (
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"

	"github.com/arsham/neuragene/internal/component"
	"github.com/arsham/neuragene/internal/entity"
	"github.com/arsham/neuragene/internal/geom"
	"github.com/arsham/neuragene/internal/quadtree"
)

// Quadtree system draws collision boxes based on the quadtrees for
// entities if their flag is set.
type Quadtree struct {
	entitties    *entity.Manager
	components   *component.Manager
	qTree        *quadtree.QuadTree[uint64]
	Colour       color.Color
	lastDuration time.Duration
	Capacity     uint
}

var _ System = (*Quadtree)(nil)

func (*Quadtree) String() string { return "CollisionBox" }

// setup returns an error if the entity manager or the component manager is
// nil.
func (q *Quadtree) setup(c controller) error {
	q.entitties = c.EntityManager()
	q.components = c.ComponentManager()
	if q.entitties == nil {
		return fmt.Errorf("%w: entity manager", ErrInvalidArgument)
	}
	if q.components == nil {
		return fmt.Errorf("%w: component manager", ErrInvalidArgument)
	}
	if q.Colour == nil {
		q.Colour = colornames.Red
	}
	if q.Capacity == 0 {
		q.Capacity = 20
	}
	return nil
}

func (*Quadtree) update(component.State) error { return nil }

func (q *Quadtree) drawChildren(t *quadtree.QuadTree[uint64], canvas *ebiten.Image) {
	for _, child := range t.Children() {
		if child == nil {
			continue
		}
		x1, y1, x2, y2 := child.Bounds()
		vector.StrokeRect(canvas, float32(x1), float32(y1), float32(x2-x1), float32(y2-y1), 1, q.Colour, false)
		q.drawChildren(child, canvas)
	}
}

func (q *Quadtree) draw(screen *ebiten.Image, state component.State) {
	started := time.Now()
	defer func() {
		q.lastDuration = time.Since(started)
	}()
	if !all(state, component.StateDrawCollisionBoxes, component.StateRunning) {
		return
	}

	maxX, maxY := ebiten.WindowSize()
	bounds := quadtree.NewBounds(0, 0, float64(maxX), float64(maxY))
	q.qTree = quadtree.NewQuadTree[uint64](bounds, q.Capacity, 0)
	defer q.qTree.Free()

	positions := q.components.Position
	q.entitties.MapByMask(entity.Collides|entity.Rigid, func(e *entity.Entity) {
		id := e.ID
		pos := positions[id]
		point := quadtree.P(
			id,
			geom.V(
				pos.Pos.Resolve().X,
				pos.Pos.Resolve().Y,
			),
		)
		q.qTree.Insert(point)
	})
	canvas := ebiten.NewImage(maxX, maxY)
	q.drawChildren(q.qTree, canvas)
	screen.DrawImage(canvas, nil)
}

// avgCalc returns the amount of time it took for the last update.
func (q *Quadtree) avgCalc() time.Duration {
	return q.lastDuration
}
