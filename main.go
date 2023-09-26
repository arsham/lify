// Package main starts the game.
package main

import (
	"embed"
	"fmt"
	"image"
	_ "image/png"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"

	"github.com/arsham/neuragene/internal/config"
	"github.com/arsham/neuragene/internal/ui/scenes"
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

	play := &scenes.Play{
		Win:          win,
		CamSpeed:     1000,
		CamZoomSpeed: 1.2,
		CamZoom:      1.0,
		CamPos:       pixel.ZV,
	}
	last := time.Now()
	frames := 0
	second := time.Tick(time.Second)

	for !win.Closed() {
		win.Clear(colornames.Whitesmoke)
		dt := time.Since(last).Seconds()
		last = time.Now()

		play.Draw(sprite, dt)

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
