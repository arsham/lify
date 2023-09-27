package system

import (
	"fmt"

	"github.com/arsham/neuragene/component"
	"github.com/arsham/neuragene/entity"
)

// Lifespan system handles the lifespan of entities. You should always use this
// system before the AI system, otherwise the AI can't collect the dead genes.
type Lifespan struct {
	entities   *entity.Manager
	components *component.Manager
}

func (l *Lifespan) String() string { return "Lifespan" }

var _ System = (*Lifespan)(nil)

// Setup returns an error if the entity manager or the component manager is
// nil.
func (l *Lifespan) Setup(c controller) error {
	l.entities = c.EntityManager()
	l.components = c.ComponentManager()
	if l.entities == nil {
		return fmt.Errorf("%w: entity manager", ErrInvalidArgument)
	}
	if l.components == nil {
		return fmt.Errorf("%w: component manager", ErrInvalidArgument)
	}
	return nil
}

// Process draws the grid on the screen.
func (l *Lifespan) Process(state component.State, _ float64) {
	if !all(state, component.StateRunning) {
		return
	}
	// Note that we don't check the state here. We always want to process this,
	// and then if required we kill the entities.
	remove := state&component.StateLimitLifespans == component.StateLimitLifespans
	lifespan := l.components.Lifespan
	l.entities.MapByMask(entity.Lifespan, func(e *entity.Entity) {
		id := e.ID
		lifespan := lifespan[id]
		lifespan.Remaining--
		if !remove {
			return
		}
		if lifespan.Remaining <= 0 {
			l.entities.Kill(e)
			return
		}
	})
}
