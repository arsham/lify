// Package system contains the logic for manipulating and managing entities
// based on their component values.
package system

import (
	"errors"
	"fmt"

	"github.com/arsham/neuragene/asset"
	"github.com/arsham/neuragene/component"
	"github.com/arsham/neuragene/entity"
	"github.com/faiface/pixel"
)

// ErrInvalidArgument indicates that the given argument is invalid or missing.
var ErrInvalidArgument = errors.New("invalid or missing argument")

// controller is a controller for a system for querying dependencies.
type controller interface {
	// EntityManager returns the entity manager.
	EntityManager() *entity.Manager
	// ComponentManager returns the component manager.
	ComponentManager() *component.Manager
	// Target returns the target object to draw on.
	Target() pixel.Target
	// Bounds returns the bounds of the target.
	Bounds() pixel.Rect
	// AssetManager returns the asset manager.
	AssetManager() *asset.Manager
	// InputDevice returns an object that informs the last action by the user.
	InputDevice() InputDevice
}

// A System should implement this interface.
type System interface {
	// Setup is a one-off call to prepare the system.
	Setup(c controller) error
	// Process is called on every frame. The current state of the system and
	// the delayed time is passed.
	Process(state component.State, dt float64)
	// String returns the name of the system.
	String() string
}

// Manager holds a series of Systems.
type Manager struct {
	systems []System
}

// NewManager returns a new Manager with pre-allocated memory by the given
// size.
func NewManager(size int) *Manager {
	return &Manager{
		systems: make([]System, 0, size),
	}
}

// Add adds the s system to the list. It doesn't check if the system is already
// been added.
func (m *Manager) Add(s System) *Manager {
	m.systems = append(m.systems, s)
	return m
}

// Setup calls the Setup() method on all systems. It returns an error if any of
// the systems returns an error.
func (m *Manager) Setup(c controller) error {
	for _, s := range m.systems {
		if err := s.Setup(c); err != nil {
			return fmt.Errorf("setting up %s system: %w", s, err)
		}
	}
	return nil
}

// Process calls the Process method of the systems with the given next state.
func (m *Manager) Process(state component.State, dt float64) {
	for _, s := range m.systems {
		s.Process(state, dt)
	}
}
