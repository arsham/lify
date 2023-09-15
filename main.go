// Package main starts the game.
package main

import (
	"embed"
	"fmt"
	"image"
	_ "image/png"
	"log/slog"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"

	"github.com/arsham/lify/internal/config"
)

//go:embed assets
var assets embed.FS

func main() {
	env, err := config.Config()
	if err != nil {
		slog.Error("Failed getting configuration: %w", err)
	}

	pixelgl.Run(func() { run(env) })
}

func run(env *config.Env) {
	cfg := pixelgl.WindowConfig{
		Title:     "Pixel Rocks!",
		Bounds:    pixel.R(0, 0, float64(env.UI.Width), float64(env.UI.Height)),
		VSync:     true,
		Resizable: true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	pic, err := loadPicture(filepath.Join("assets", "images", "herb1.png"))
	if err != nil {
		panic(err)
	}

	sprite := pixel.NewSprite(pic, pic.Bounds())
	sprite.Draw(win, pixel.IM)
	sprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()))

	angel := 0.0
	last := time.Now()
	camPos := pixel.ZV
	camSpeed := 1000.0
	camZoom := 1.0
	camZoomSpeed := 1.2
	var (
		frames = 0
		second = time.Tick(time.Second)
	)

	for !win.Closed() {
		win.Clear(colornames.Whitesmoke)
		dt := time.Since(last).Seconds()
		last = time.Now()
		angel += 3 * dt

		mat := pixel.IM
		mat = mat.Rotated(pixel.ZV, angel)
		mat = mat.Moved(win.Bounds().Center())

		if win.Pressed(pixelgl.KeyLeft) {
			camPos.X -= camSpeed * dt
		}
		if win.Pressed(pixelgl.KeyRight) {
			camPos.X += camSpeed * dt
		}
		if win.Pressed(pixelgl.KeyDown) {
			camPos.Y -= camSpeed * dt
		}
		if win.Pressed(pixelgl.KeyUp) {
			camPos.Y += camSpeed * dt
		}

		if win.Pressed(pixelgl.KeyUp) {
			camPos.Y += camSpeed * dt
		}
		camZoom *= math.Pow(camZoomSpeed, win.MouseScroll().Y)

		cam := pixel.IM.Scaled(camPos, camZoom).Moved(win.Bounds().Center().Sub(camPos))
		win.SetMatrix(cam)

		sprite.Draw(win, mat)

		sprite.Draw(win, pixel.IM.Moved(win.Bounds().Center().Add(pixel.Vec{X: 100, Y: 200})))
		sprite.Draw(win, pixel.IM.Moved(win.Bounds().Center().Add(pixel.Vec{X: -200, Y: -200})))
		sprite.Draw(win, pixel.IM.Moved(win.Bounds().Center().Add(pixel.Vec{X: -300, Y: -300})))
		sprite.Draw(win, pixel.IM.Moved(win.Bounds().Center().Add(pixel.Vec{X: -400, Y: -400})))

		win.Update()

		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		default:
		}
	}
}

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}
