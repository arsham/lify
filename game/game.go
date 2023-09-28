// Package game contains the logic for running the game.
package game

import (
	"fmt"
	"image/color"
	"io/fs"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"

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
	// Update updates the game state. When the scene wants to exit it will
	// return an ebitern.Termination error.
	Update() error
	// Draw draws the system's state onto the screen.
	Draw(screen *ebiten.Image)
}

// The Engine manages the game loop and makes decisions on changing scenes.
type Engine struct {
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
	// second is a ticker for updating the window's title.
	second *time.Ticker
	// title is the title of the window.
	title string
	// lastFrameDuration is the duration of the previous frame.
	lastFrameDuration time.Duration
	// currentScene is the currently playing scene.
	currentScene scene.Type
}

// NewEngine creates a new game engine with all the dependencies and sets up
// the first scene. It returns an error if any of the dependencies can't be
// created.
func NewEngine(env *config.Env, filesystem fs.FS) (*Engine, error) {
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetVsyncEnabled(true)
	ebiten.SetWindowTitle("Neuragene")
	ebiten.SetWindowSize(env.UI.Width, env.UI.Height)

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
		&system.Position{},
		&system.Lifespan{},
		&system.Stats{},
		&system.BoundingBox{
			Size: 1,
		},
	)
	g := &Engine{
		title:        "Neuragene",
		entities:     em,
		systems:      sm,
		currentScene: scene.PlayScene,
		assets:       am,
		components:   components,
		second:       time.NewTicker(time.Second),
	}
	g.scenes = map[scene.Type]sceneRunner{
		scene.PlayScene: scene.NewPlay(g),
	}
	err = g.systems.Setup(g)
	if err != nil {
		return nil, fmt.Errorf("setting up the engine: %w", err)
	}
	return g, nil
}

// Update updates a game by one tick. The given argument represents a screen
// image.
func (e *Engine) Update() error {
	// TODO: move this to the stats system, or even completely ignore it and
	// show these on the HUD.
	select {
	case <-e.second.C:
		ebiten.SetWindowTitle(fmt.Sprintf("%s | FPS: %.2f", e.title, ebiten.ActualTPS()))
	default:
	}
	return e.scene().Update()
}

// Draw draws the game screen onto the screen.
func (e *Engine) Draw(screen *ebiten.Image) {
	started := time.Now()
	screen.Clear()
	screen.Fill(colornames.Whitesmoke)
	e.scene().Draw(screen)
	e.lastFrameDuration = time.Since(started)
}

// Layout accepts a native outside size in device-independent pixels and
// returns the game's logical screen size.
func (e *Engine) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

// scene returns the current scene.
func (e *Engine) scene() sceneRunner {
	return e.scenes[e.currentScene]
}

// ComponentManager returns the component manager.
func (e *Engine) ComponentManager() *component.Manager {
	return e.components
}

// EntityManager returns the entity manager.
func (e *Engine) EntityManager() *entity.Manager {
	return e.entities
}

// SystemManager returns the system manager.
func (e *Engine) SystemManager() *system.Manager {
	return e.systems
}

// AssetManager returns the asset manager.
func (e *Engine) AssetManager() *asset.Manager {
	return e.assets
}

// LastFrameDuration returns the time it took to process previous frame.
func (e *Engine) LastFrameDuration() time.Duration {
	return e.lastFrameDuration
}
