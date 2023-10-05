// Package entity contains the logic to handle entities.
package entity

import (
	"slices"
	"sync/atomic"

	"github.com/arsham/neuragene/internal/component"
)

// Mask is used to determine if an entity supports a specific logic.
type Mask uint32

const (
	// Positioned mask indicates that the entity has a position, scale and
	// velocity.
	Positioned Mask = 1 << iota
	// HasTexture mask indicates that the entity has a texture and should be
	// rendered.
	HasTexture
	// Lifespan mask is used to set the lifespan for entities.
	Lifespan
	// Died mask is used to mark dead entities.
	Died
	// BoxBounded is an entity that has a box boundary.
	BoxBounded
	// Collides marks an entity that should be checked against other entities
	// with the Collides or Rigid masks.
	Collides
	// Rigid marks the entity that should cause other entities with Collides
	// mask to be bounced, but it is not be resolved in terms of collisions.
	Rigid
)

// An Entity is an element in the game that can have at least one component.
// Each component is managed by the component.Manager as a map if entity ID to
// its respective component. When an Entity is removed, its ID is removed from
// all the maps in the Manager. You should not create an Entity directly,
// instead you should use the Manager's NewEntity method.
type Entity struct {
	ID   uint64
	mask Mask
}

// List contains a slice of Entity.
type List []*Entity

// Manager manages all the entities in the game. It is advised to create a new
// Manager by calling the NewManager constructor to pre-allocate memory.
type Manager struct {
	components *component.Manager
	entities   List
	// When a new entity is added, it will first go into this list. At the end
	// of the frame, all the entities in this list will be added to the
	// entities list.
	toAdd List
	// lastID is the last ID that was given to an entity.
	lastID uint64
}

// NewManager returns a new Manager with pre-allocated memory by the given
// size.
func NewManager(components *component.Manager, size int) *Manager {
	return &Manager{
		components: components,
		entities:   make(List, 0, size),
		toAdd:      make(List, 0, size),
	}
}

// NewEntity creates a new Entity with the given mask. The entity is not added
// to the available entities, but it will be on the next Update call. You
// should manually set the necessary components for the entity. Your call will
// panic if the component you are trying to set doesn't match the mask.
func (m *Manager) NewEntity(mask Mask) *Entity {
	e := &Entity{
		ID:   atomic.AddUint64(&m.lastID, 1),
		mask: mask,
	}
	m.toAdd = append(m.toAdd, e)
	return e
}

// MapByMask applies the given function to all the entities that match the
// given mask.
func (m *Manager) MapByMask(mask Mask, fn func(*Entity)) {
	for _, e := range m.entities {
		if e.mask&mask != 0 {
			fn(e)
		}
	}
}

// Update moves new entities from the toAdd slice to entities slice, and
// removes any that are dead.
func (m *Manager) Update(state component.State) {
	m.entities = append(m.entities, m.toAdd...)
	clear(m.toAdd)
	m.toAdd = m.toAdd[:0]

	if state&component.StateLimitLifespans == component.StateLimitLifespans {
		m.entities = slices.DeleteFunc(m.entities, func(e *Entity) bool {
			return e.mask&Died == Died
		})
	}
}

// Len returns the number of entities.
func (m *Manager) Len() int {
	return len(m.entities)
}

// Kill marks the entity as dead. It will be removed on the next Update call.
func (m *Manager) Kill(e *Entity) {
	e.mask |= Died
}
