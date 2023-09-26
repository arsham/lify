// Package game contains the logic for running the game.
package game

import (
	"fmt"
	"time"

	"github.com/arsham/neuragene/entity"
	"github.com/arsham/neuragene/scene"
	"github.com/arsham/neuragene/system"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

// A Scene defines the contract for communicating with the currently processing
// scene.
type Scene interface {
	Update(dt float64)
}

// The Engine manages the game loop and makes decisions on changing scenes.
type Engine struct {
	// Window is the current window.
	Window *pixelgl.Window
	// Scenes is a map of all available scenes.
	Scenes map[scene.Type]Scene
	// Systems is the system manager.
	Systems *system.Manager
	// Entities is the entity manager.
	Entities *entity.Manager
	// Title is the title of the window.
	Title string
	// lastFrameDuration is the duration of the previous frame.
	lastFrameDuration time.Duration
	// CurrentScene is the currently playing scene.
	CurrentScene scene.Type
	// When running is set to false the game loop will stop.
	running bool
}

// Run listens to the user input and informs the current scene to update
// itself. If the current scene returns a different scene, it will switch to
// the new scene.
func (e *Engine) Run() {
	frames := 0
	second := time.NewTicker(time.Second)
	last := time.Now()
	for !e.Window.Closed() && e.running {
		started := time.Now()
		dt := time.Since(last).Seconds()
		last = time.Now()

		e.Window.Clear(colornames.Whitesmoke)
		e.Scene().Update(dt)
		e.Window.Update()

		e.lastFrameDuration = time.Since(started)
		frames++
		select {
		case <-second.C:
			e.Window.SetTitle(fmt.Sprintf("%s | FPS: %d", e.Title, frames))
			frames = 0
		default:
		}
	}
}

// Scene returns the current scene.
func (e *Engine) Scene() Scene {
	return e.Scenes[e.CurrentScene]
}
