package scene

import (
	"github.com/faiface/pixel/pixelgl"

	"github.com/arsham/neuragene/action"
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
			actionMap:  make(map[pixelgl.Button]action.Name, 100),
			entities:   c.EntityManager(),
			systems:    c.SystemManager(),
			controller: c,
			state: component.StateRunning |
				component.StateMoveEntities,
		},
	}
	p.RegisterAction(pixelgl.KeyEscape, action.Quit)
	p.RegisterAction(pixelgl.KeyQ, action.Quit)
	return p
}

// Update updates the system. If any of the system reports the next state
// should be `stop`, it updates its state to be paused.
func (p *Play) Update(dt float64) {
	p.update()
	p.entities.Update(p.state)
	p.systems.Process(p.state, dt)
}

// Do applies the action on the scene state.
func (p *Play) Do(a action.Action) {
	if a.Phase == action.PhaseStart {
		if a.Name == action.Quit {
			p.state ^= component.StateQuit
		}
	}
}
