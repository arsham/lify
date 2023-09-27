package system

import (
	"errors"
	"fmt"

	"github.com/faiface/pixel/pixelgl"

	"github.com/arsham/neuragene/action"
	"github.com/arsham/neuragene/component"
)

// A Scene contains a map of actions and can act on specific actions.
type Scene interface {
	// Actions returns all the registered actions.
	Actions() map[pixelgl.Button]action.Name
	// Do acts on the given action.
	Do(action.Action)
}

// InputDevice returns an object that informs the last action by the user.
type InputDevice interface {
	// JustPressed returns whether the Button has just been pressed down.
	JustPressed(button pixelgl.Button) bool
	// JustReleased returns whether the Button has just been released up.
	JustReleased(button pixelgl.Button) bool
}

// UserInput system handles the UserInput.
type UserInput struct {
	// Scene returns the current scene. This is required for this object to
	// emit actions.
	Scene func() Scene
	input InputDevice
}

var _ System = (*UserInput)(nil)

func (u *UserInput) String() string { return "UserInput" }

// Setup returns an error if the Scene function is not set.
func (u *UserInput) Setup(c controller) error {
	if u.Scene == nil {
		return errors.New("scene function is not set")
	}
	u.input = c.InputDevice()
	if u.input == nil {
		return fmt.Errorf("%w: target", ErrInvalidArgument)
	}
	return nil
}

// Process reads the current events and issues actions if the scene has
// registered them.
func (u *UserInput) Process(component.State, float64) {
	s := u.Scene()
	for key, name := range s.Actions() {
		if u.input.JustPressed(key) {
			s.Do(action.Action{
				Key:   key,
				Name:  name,
				Phase: action.PhaseStart,
			})
			break
		}
		if u.input.JustReleased(key) {
			s.Do(action.Action{
				Key:   key,
				Name:  name,
				Phase: action.PhaseEnd,
			})
		}
	}
}
