package scene

import (
	"github.com/arsham/neuragene/component"
)

// Play is a scene that plays the simulation.
type Play struct {
	*Generic
}

var _ derivedScene = &Play{}

// NewPlay returns a new Play scene for the window and with entity and system
// managers received by the controller. It sets up its own set of key bindings.
func NewPlay(c controller) *Play {
	p := &Play{
		Generic: &Generic{
			entities:   c.EntityManager(),
			systems:    c.SystemManager(),
			controller: c,
			state:      component.StateMoveEntities,
		},
	}
	return p
}

// Update updates the system. If any of the system reports the next state
// should be `stop`, it updates its state to be paused.
func (p *Play) Update(dt float64) {
	p.update()
	p.entities.Update(p.state)
	p.systems.Process(p.state, dt)
}
