// Package game contains the logic for running the game.
package game

import (
	"fmt"
	"io/fs"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"

	"github.com/arsham/neuragene/asset"
	"github.com/arsham/neuragene/component"
	"github.com/arsham/neuragene/entity"
	"github.com/arsham/neuragene/internal/config"
	"github.com/arsham/neuragene/scene"
	"github.com/arsham/neuragene/system"
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
	// Components is the component manager.
	Components *component.Manager
	// Assets is the asset manager.
	Assets *asset.Manager
	// Title is the title of the window.
	Title string
	// lastFrameDuration is the duration of the previous frame.
	lastFrameDuration time.Duration
	// CurrentScene is the currently playing scene.
	CurrentScene scene.Type
	// When running is set to false the game loop will stop.
	running bool
}

// NewEngine creates a new game engine with all the dependencies and sets up
// the first scene. It returns an error if any of the dependencies can't be
// created.
func NewEngine(env *config.Env, filesystem fs.FS) (*Engine, error) {
	cfg := pixelgl.WindowConfig{
		Title:     "Neuragene",
		Bounds:    pixel.R(0, 0, float64(env.UI.Width), float64(env.UI.Height)),
		VSync:     true,
		Resizable: true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		return nil, fmt.Errorf("creating new window: %w", err)
	}
	am, err := asset.New(filesystem)
	if err != nil {
		return nil, fmt.Errorf("creating new asset manager: %w", err)
	}

	size := 1000
	components := &component.Manager{
		Position: make(map[uint64]*component.Position, size),
	}
	em := entity.NewManager(components, size)
	sm := system.NewManager(10)
	sm.Add(
		&system.Rendering{
			Title:  "Neuragene",
			Width:  int32(env.UI.Width),
			Height: int32(env.UI.Height),
		},
	)
	g := &Engine{
		Window:       win,
		Entities:     em,
		Systems:      sm,
		CurrentScene: scene.PlayScene,
		Assets:       am,
		Components:   components,
		running:      true,
	}
	g.Scenes = map[scene.Type]Scene{
		scene.PlayScene: scene.NewPlay(g),
	}
	err = g.Setup()
	if err != nil {
		return nil, fmt.Errorf("setting up the engine: %w", err)
	}
	return g, nil
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

// Setup calls the Setup() method of the system manager.
func (e *Engine) Setup() error {
	return e.Systems.Setup(e)
}

// Bounds returns the bounds of the target.
func (e *Engine) Bounds() pixel.Rect {
	return e.Window.Bounds()
}

// ComponentManager returns the component manager.
func (e *Engine) ComponentManager() *component.Manager {
	return e.Components
}

// EntityManager returns the entity manager.
func (e *Engine) EntityManager() *entity.Manager {
	return e.Entities
}

// Target returns the target object to draw on.
func (e *Engine) Target() pixel.Target {
	return e.Window
}

// SystemManager returns the system manager.
func (e *Engine) SystemManager() *system.Manager {
	return e.Systems
}

// AssetManager returns the asset manager.
func (e *Engine) AssetManager() *asset.Manager {
	return e.Assets
}
