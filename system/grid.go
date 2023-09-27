package system

import (
	"fmt"
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"

	"github.com/arsham/neuragene/component"
	"github.com/arsham/neuragene/entity"
)

// Grid draws a grid on the screen.
type Grid struct {
	target        pixel.Target
	canvas        *pixelgl.Canvas
	entities      *entity.Manager
	Colour        color.Color
	lastWinBounds pixel.Rect
	bounds        pixel.Rect
	GridSize      int
	Size          float64
}

var _ System = (*Grid)(nil)

func (g *Grid) String() string { return "Grid" }

// Setup returns an error if the window or the entity manager is nil.
func (g *Grid) Setup(c controller) error {
	g.target = c.Target()
	g.entities = c.EntityManager()
	g.bounds = c.Bounds()
	if g.target == nil {
		return fmt.Errorf("%w: window", ErrInvalidArgument)
	}
	if g.entities == nil {
		return fmt.Errorf("%w: entity manager", ErrInvalidArgument)
	}
	if g.bounds == pixel.ZR {
		return fmt.Errorf("%w: bounds", ErrInvalidArgument)
	}
	if g.Colour == nil {
		g.Colour = colornames.Lightgray
	}
	if g.Size == 0 {
		g.Size = 1
	}
	if g.GridSize == 0 {
		g.GridSize = 25
	}
	g.canvas = pixelgl.NewCanvas(g.bounds)
	g.drawGrid()
	g.lastWinBounds = g.bounds
	return nil
}

func (g *Grid) drawGrid() {
	imd := imdraw.New(nil)
	imd.EndShape = imdraw.RoundEndShape
	imd.Color = g.Colour
	for x := g.bounds.Min.X; x < g.bounds.Max.X; x += float64(g.GridSize) {
		imd.Push(pixel.V(x, g.bounds.Min.Y), pixel.V(x, g.bounds.Max.Y))
		imd.Line(g.Size)
	}
	for y := g.bounds.Min.Y; y < g.bounds.Max.Y; y += float64(g.GridSize) {
		imd.Push(pixel.V(g.bounds.Min.X, y), pixel.V(g.bounds.Max.X, y))
		imd.Line(g.Size)
	}
	imd.Draw(g.canvas)
}

// Process draws the grid on the screen.
func (g *Grid) Process(state component.State, _ float64) {
	if !all(state, component.StateDrawGrids) {
		return
	}
	if g.lastWinBounds != g.bounds {
		g.canvas.SetBounds(g.bounds)
		g.lastWinBounds = g.bounds
		g.drawGrid()
	}

	mat := pixel.IM.Moved(g.bounds.Center())
	g.canvas.Draw(g.target, mat)
}
