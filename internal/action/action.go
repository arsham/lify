// Package action contains the translation of the user input to actions.
package action

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

// Phase specifies the behaviour of an action.
type Phase uint8

const (
	// PhaseStart is when the action is initiated.
	PhaseStart Phase = iota + 1
	// PhaseEnd is when the action is abandoned.
	PhaseEnd
)

func (p Phase) String() string {
	if p == PhaseStart {
		return "Start"
	}
	return "End"
}

// Name is the name of an action.
//
//go:generate stringer -type=Name -output=name_string.go
type Name uint16

const (
	// Quit action causes the simulation to end.
	Quit Name = iota + 1
	// ToggleGrid action toggles the grid.
	ToggleGrid
	// ToggleLimitFPS action toggles the FPS limit.
	ToggleLimitFPS
	// ToggleLimitLifespans action toggles the lifespan limit.
	ToggleLimitLifespans
	// ToggleBoundingBoxes action toggles the bounding boxes.
	ToggleBoundingBoxes
	// Pause action will cause the simulation to pause or resume depending on
	// the previous state.
	Pause
	// ToggleTextures toggles drawing of renderable textures. This has no
	// effect on drawing bounding boxes.
	ToggleTextures
	// ToggleCollisions toggles detection of collisions.
	ToggleCollisions
)

// An Action is an input state that would result in an activity in a scene.
type Action struct {
	Key   ebiten.Key
	Phase Phase
	Name  Name
}

func (a Action) String() string {
	return fmt.Sprintf("%s:%s (%s)", a.Name, a.Key, a.Phase)
}
