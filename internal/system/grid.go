package system

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"

	"github.com/arsham/neuragene/internal/component"
)

// Grid draws a grid on the screen.
type Grid struct {
	canvas        *ebiten.Image
	Colour        color.Color
	lastWinBounds image.Point
	GridSize      int
	Size          float64
}

var _ System = (*Grid)(nil)

func (g *Grid) String() string { return "Grid" }

// setup returns an error if the window or the entity manager is nil.
func (g *Grid) setup(controller) error {
	if g.Colour == nil {
		g.Colour = colornames.Lightgray
	}
	if g.Size == 0 {
		g.Size = 1
	}
	if g.GridSize == 0 {
		g.GridSize = 25
	}
	x, y := ebiten.WindowSize()
	g.canvas = ebiten.NewImage(x, y)
	g.drawGrid(x, y)
	g.lastWinBounds = image.Point{X: x, Y: y}
	return nil
}

func (g *Grid) drawGrid(w, h int) {
	img := ebiten.NewImage(w, h)
	for x := 0; x < w; x += g.GridSize {
		vector.StrokeLine(img, float32(x), 0, float32(x), float32(h), float32(g.Size), g.Colour, false)
	}
	for y := 0; y < h; y += g.GridSize {
		vector.StrokeLine(img, 0, float32(y), float32(w), float32(y), float32(g.Size), g.Colour, false)
	}
	g.canvas.DrawImage(img, nil)
}

// update draws the grid on a cached canvas if the window size has changed.
func (g *Grid) update(state component.State) error {
	if !all(state, component.StateDrawGrids) {
		return nil
	}
	x, y := ebiten.WindowSize()
	bounds := image.Point{X: x, Y: y}
	if !g.lastWinBounds.Eq(bounds) {
		g.lastWinBounds = bounds
		g.canvas = ebiten.NewImage(ebiten.WindowSize())
		g.drawGrid(x, y)
	}
	return nil
}

// draw draws the cached canvas on the screen.
func (g *Grid) draw(screen *ebiten.Image, state component.State) {
	if !all(state, component.StateDrawGrids) {
		return
	}
	screen.DrawImage(g.canvas, nil)
}
