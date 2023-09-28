// Package scene contains the logic for handling a scene.
package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/arsham/neuragene/internal/action"
	"github.com/arsham/neuragene/internal/component"
	"github.com/arsham/neuragene/internal/entity"
	"github.com/arsham/neuragene/internal/system"
)

// Type is the type of different scenes.
type Type uint8

const (
	// PlayScene specifies the playing scene.
	PlayScene Type = iota
)

// controller is the game object that provides vital information to the scene.
type controller interface {
	// EntityManager returns the entity manager.
	EntityManager() *entity.Manager
	// SystemManager returns the system manager.
	SystemManager() *system.Manager
}

// Generic is the base struct for scenes. You should always call the update()
// method from the derived struct's Update() method!
type Generic struct {
	entities *entity.Manager
	systems  *system.Manager
	// controller is the game object that provides vital information to the
	// scene.
	controller controller
	// actionMap contains all the actions a scene can act on.
	actionMap map[ebiten.Key]action.Name
	// actOnFn should be provided by the sub-scenes otherwise they can't
	// receive user input.
	actOnFn func(action.Action)
	// state is the current state of the game.
	state      component.State
	frameCount int64
}

// A system should implement this interface.
type derivedScene interface {
	// Update updates the game state. When the scene wants to exit it will
	// return an ebitern.Termination error.
	Update() error
	// Draw draws the system's state onto the screen.
	Draw(screen *ebiten.Image)
}

func (g *Generic) update() {
	g.frameCount++
	for key, name := range g.actionMap {
		if inpututil.IsKeyJustPressed(key) {
			g.actOnFn(action.Action{
				Key:   key,
				Name:  name,
				Phase: action.PhaseStart,
			})
			continue
		}
		if inpututil.IsKeyJustReleased(key) {
			g.actOnFn(action.Action{
				Key:   key,
				Name:  name,
				Phase: action.PhaseEnd,
			})
		}
	}
}

// registerAction registers the given button/mouse action with the associated
// name.
func (g *Generic) registerAction(key ebiten.Key, name action.Name) {
	g.actionMap[key] = name
}
