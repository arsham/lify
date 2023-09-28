package scene

import (
	"github.com/hajimehoshi/ebiten/v2"

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
			actionMap:  make(map[ebiten.Key]action.Name, 100),
			entities:   c.EntityManager(),
			systems:    c.SystemManager(),
			controller: c,
			state: component.StateRunning |
				component.StateSpawnAnts |
				component.StatePrintStats |
				component.StateLimitLifespans |
				component.StateDrawTextures |
				component.StateMoveEntities,
		},
	}
	p.actOnFn = p.actOn
	p.registerAction(ebiten.KeyEscape, action.Quit)
	p.registerAction(ebiten.KeyQ, action.Quit)
	p.registerAction(ebiten.KeyG, action.ToggleGrid)
	p.registerAction(ebiten.KeyF, action.ToggleLimitFPS)
	p.registerAction(ebiten.KeyL, action.ToggleLimitLifespans)
	p.registerAction(ebiten.KeyB, action.ToggleBoundingBoxes)
	p.registerAction(ebiten.KeySpace, action.Pause)
	p.registerAction(ebiten.KeyT, action.ToggleTextures)
	return p
}

// Update updates the system. If any of the system reports the next state
// should be `stop`, it updates its state to be paused.
func (p *Play) Update() error {
	p.update()
	if p.state&component.StateQuit != 0 {
		return ebiten.Termination
	}
	if p.state&component.StateLimitFPS != 0 {
		ebiten.SetTPS(400)
	} else {
		ebiten.SetTPS(60)
	}
	p.entities.Update(p.state)
	return p.systems.Update(p.state)
}

// Draw draws the game screen onto the screen.
func (p *Play) Draw(screen *ebiten.Image) {
	p.systems.Draw(screen, p.state)
}

// Do applies the action on the scene state.
func (p *Play) actOn(a action.Action) {
	if a.Phase == action.PhaseStart {
		switch a.Name {
		case action.Quit:
			p.state ^= component.StateQuit
		case action.ToggleGrid:
			p.state ^= component.StateDrawGrids
		case action.ToggleLimitFPS:
			p.state ^= component.StateLimitFPS
		case action.ToggleLimitLifespans:
			p.state ^= component.StateLimitLifespans
		case action.ToggleBoundingBoxes:
			p.state ^= component.StateDrawBoundingBoxes
		case action.Pause:
			p.state ^= component.StateRunning
		case action.ToggleTextures:
			p.state ^= component.StateDrawTextures
		}
	}
}
