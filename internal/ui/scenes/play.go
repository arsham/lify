// Package scenes contains logic for drawing a scene on the window.
package scenes

import (
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

// Play draws the normal playing scene.
type Play struct {
	Win          *pixelgl.Window
	angel        float64
	CamSpeed     float64
	CamZoomSpeed float64
	CamPos       pixel.Vec
	CamZoom      float64
}

// Draw draws the sprite on the screen. The dt is the delta time since the last
// draw.
func (p *Play) Draw(sprite *pixel.Sprite, dt float64) {
	p.Win.Clear(colornames.Whitesmoke)
	p.angel += 3 * dt

	mat := pixel.IM
	mat = mat.Rotated(pixel.ZV, p.angel)
	centre := p.Win.Bounds().Center()
	mat = mat.Moved(centre)

	if p.Win.Pressed(pixelgl.KeyLeft) {
		p.CamPos.X -= p.CamSpeed * dt
	}
	if p.Win.Pressed(pixelgl.KeyRight) {
		p.CamPos.X += p.CamSpeed * dt
	}
	if p.Win.Pressed(pixelgl.KeyDown) {
		p.CamPos.Y -= p.CamSpeed * dt
	}
	if p.Win.Pressed(pixelgl.KeyUp) {
		p.CamPos.Y += p.CamSpeed * dt
	}

	if p.Win.Pressed(pixelgl.KeyUp) {
		p.CamPos.Y += p.CamSpeed * dt
	}
	p.CamZoom *= math.Pow(p.CamZoomSpeed, p.Win.MouseScroll().Y)

	cam := pixel.IM.Scaled(p.CamPos, p.CamZoom).Moved(p.Win.Bounds().Center().Sub(p.CamPos))
	p.Win.SetMatrix(cam)

	sprite.Draw(p.Win, mat)

	sprite.Draw(p.Win, pixel.IM.Moved(centre.Add(pixel.Vec{X: 100, Y: 200})))
	sprite.Draw(p.Win, pixel.IM.Moved(centre.Add(pixel.Vec{X: -200, Y: -200})))
	sprite.Draw(p.Win, pixel.IM.Moved(centre.Add(pixel.Vec{X: -300, Y: -300})))
	sprite.Draw(p.Win, pixel.IM.Moved(centre.Add(pixel.Vec{X: -400, Y: -400})))
}
