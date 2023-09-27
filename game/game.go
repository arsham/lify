// Package game contains the logic for running the game.
package game

import (
	"fmt"
	"image/color"
	"io/fs"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"

	"github.com/arsham/neuragene/action"
	"github.com/arsham/neuragene/asset"
	"github.com/arsham/neuragene/component"
	"github.com/arsham/neuragene/entity"
	"github.com/arsham/neuragene/internal/config"
	"github.com/arsham/neuragene/scene"
	"github.com/arsham/neuragene/system"
)

// A sceneRunner defines the contract for communicating with the currently
// processing scene.
type sceneRunner interface {
	Update(dt float64)
	Do(action.Action)
	Actions() map[pixelgl.Button]action.Name
	// State returns the current state of the scene.
	State() component.State
}

// The Engine manages the game loop and makes decisions on changing scenes.
type Engine struct {
	// window is the current window.
	window *pixelgl.Window
	// scenes is a map of all available scenes.
	scenes map[scene.Type]sceneRunner
	// systems is the system manager.
	systems *system.Manager
	// entities is the entity manager.
	entities *entity.Manager
	// components is the component manager.
	components *component.Manager
	// assets is the asset manager.
	assets *asset.Manager
	// title is the title of the window.
	title string
	// lastFrameDuration is the duration of the previous frame.
	lastFrameDuration time.Duration
	// currentScene is the currently playing scene.
	currentScene scene.Type
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
		Position:  make(map[uint64]*component.Position, size),
		Sprite:    make(map[uint64]*component.Sprite, size),
		Lifespan:  make(map[uint64]*component.Lifespan, size),
		Collision: make(map[uint64]*component.Collision, size),
	}
	em := entity.NewManager(components, size)
	sm := system.NewManager(10)
	sm.Add(
		&system.FPS{Max: 60},
		&system.Grid{
			GridSize: 10,
			Size:     1,
			Colour:   color.RGBA{220, 220, 250, 255},
		},
		&system.Grid{
			GridSize: 100,
			Size:     1,
			Colour:   colornames.Lightskyblue,
		},
		&system.Rendering{
			Title:  "Neuragene",
			Width:  int32(env.UI.Width),
			Height: int32(env.UI.Height),
		},
		&system.Ant{
			Seed:         1,
			MutationRate: 100,
		},
		&system.Movement{},
		&system.Lifespan{},
		&system.BoundingBox{
			Size: 1,
		},
	)
	g := &Engine{
		window:       win,
		entities:     em,
		systems:      sm,
		currentScene: scene.PlayScene,
		assets:       am,
		components:   components,
		running:      true,
	}
	g.scenes = map[scene.Type]sceneRunner{
		scene.PlayScene: scene.NewPlay(g),
	}
	sm.Add(&system.UserInput{
		Scene: func() system.Scene { return g.Scene() },
	}, &system.Stats{Timer: g})
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
	running := true
	for !e.window.Closed() && running {
		started := time.Now()
		dt := time.Since(last).Seconds()
		last = time.Now()

		e.window.Clear(colornames.Whitesmoke)
		e.Scene().Update(dt)
		e.window.Update()

		// When the StateQuit bit is set we want to exit the game loop.
		running = e.Scene().State()&component.StateQuit != component.StateQuit

		e.lastFrameDuration = time.Since(started)
		frames++
		select {
		case <-second.C:
			e.window.SetTitle(fmt.Sprintf("%s | FPS: %d", e.title, frames))
			frames = 0
		default:
		}
	}
}

// Scene returns the current scene.
func (e *Engine) Scene() sceneRunner {
	return e.scenes[e.currentScene]
}

// Setup calls the Setup() method of the system manager.
func (e *Engine) Setup() error {
	return e.systems.Setup(e)
}

// Bounds returns the bounds of the target.
func (e *Engine) Bounds() pixel.Rect {
	return e.window.Bounds()
}

// ComponentManager returns the component manager.
func (e *Engine) ComponentManager() *component.Manager {
	return e.components
}

// EntityManager returns the entity manager.
func (e *Engine) EntityManager() *entity.Manager {
	return e.entities
}

// Target returns the target object to draw on.
func (e *Engine) Target() pixel.Target {
	return e.window
}

// SystemManager returns the system manager.
func (e *Engine) SystemManager() *system.Manager {
	return e.systems
}

// AssetManager returns the asset manager.
func (e *Engine) AssetManager() *asset.Manager {
	return e.assets
}

// InputDevice returns an object that informs the last action by the user.
func (e *Engine) InputDevice() system.InputDevice {
	return e.window
}

// LastFrameDuration returns the time it took to process previous frame.
func (e *Engine) LastFrameDuration() time.Duration {
	return e.lastFrameDuration
}
