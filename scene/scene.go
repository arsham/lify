// Package scene contains the logic for handling a scene.
package scene

import (
	"github.com/faiface/pixel/pixelgl"

	"github.com/arsham/neuragene/action"
	"github.com/arsham/neuragene/component"
	"github.com/arsham/neuragene/entity"
	"github.com/arsham/neuragene/system"
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
	actionMap map[pixelgl.Button]action.Name
	// state is the current state of the game.
	state      component.State
	frameCount int64
}

// A system should implement this interface.
type derivedScene interface {
	// Update updates the system.
	Update(dt float64)
	// Do applies the action on the scene state.
	Do(a action.Action)
}

func (g *Generic) update() {
	g.frameCount++
}

// RegisterAction registers the given button/mouse action with the associated
// name.
func (g *Generic) RegisterAction(key pixelgl.Button, name action.Name) {
	g.actionMap[key] = name
}

// Actions returns all registered actions.
func (g *Generic) Actions() map[pixelgl.Button]action.Name {
	return g.actionMap
}

// State returns the current state of the scene.
func (g *Generic) State() component.State { return g.state }
