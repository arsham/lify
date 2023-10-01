// Package system contains the logic for manipulating and managing entities
// based on their component values.
package system

import (
	"errors"
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/arsham/neuragene/internal/asset"
	"github.com/arsham/neuragene/internal/component"
	"github.com/arsham/neuragene/internal/entity"
)

// ErrInvalidArgument indicates that the given argument is invalid or missing.
var ErrInvalidArgument = errors.New("invalid or missing argument")

// controller is a controller for a system for querying dependencies.
type controller interface {
	// EntityManager returns the entity manager.
	EntityManager() *entity.Manager
	// ComponentManager returns the component manager.
	ComponentManager() *component.Manager
	// AssetManager returns the asset manager.
	AssetManager() *asset.Manager
	// LastFrameDuration returns the time it took to execute the last frame.
	LastFrameDuration() time.Duration
}

// A System should implement this interface.
type System interface {
	// setup is a one-off call to prepare the system.
	setup(c controller) error
	update(state component.State) error
	draw(screen *ebiten.Image, state component.State)
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
func (m *Manager) Add(s ...System) *Manager {
	m.systems = append(m.systems, s...)
	return m
}

// Setup calls the Setup() method on all systems. It returns an error if any of
// the systems returns an error.
func (m *Manager) Setup(c controller) error {
	for _, s := range m.systems {
		if err := s.setup(c); err != nil {
			return fmt.Errorf("setting up %s system: %w", s, err)
		}
	}
	return nil
}

// Update updates the systems that have the given state.
func (m *Manager) Update(state component.State) error {
	for _, s := range m.systems {
		err := s.update(state)
		if err != nil {
			return fmt.Errorf("system %s encountered an error: %w", s, err)
		}
	}
	return nil
}

// Draw draws the systems that have the given state.
func (m *Manager) Draw(screen *ebiten.Image, state component.State) {
	for _, s := range m.systems {
		s.draw(screen, state)
	}
}

// all returns false if any of the flags is not set in the state.
func all(state component.State, flags ...component.State) bool {
	for _, f := range flags {
		if state&f != f {
			return false
		}
	}
	return true
}

type noDraw struct{}

func (noDraw) draw(*ebiten.Image, component.State) {}
