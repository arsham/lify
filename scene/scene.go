// Package scene contains the logic for handling a scene.
package scene

import (
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
	// state is the current state of the game.
	state      component.State
	frameCount int64
}

// A system should implement this interface.
type derivedScene interface {
	// Update updates the system.
	Update(dt float64)
}

func (g *Generic) update() {
	g.frameCount++
}
