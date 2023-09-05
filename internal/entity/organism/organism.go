// Package organism contains the logic for managing and interacting with living
// things.
package organism

import (
	"errors"

	"github.com/arsham/lify/internal/genome"
	"github.com/arsham/lify/internal/pool"
)

var (
	// ErrNilDNA specifies that the DNA object is nil.
	ErrNilDNA = errors.New("dna is nil")
	// ErrAlreadyFull specifies that the Living object is already full.
	ErrAlreadyFull = errors.New("already full")
	// ErrNotEnoughFood specifies that the Living object does not have enough
	// food for certain action.
	ErrNotEnoughFood = errors.New("not enough food")
)

var livingPool = pool.NewPool(func() *living {
	l := &living{}
	l.stamina.max = 100
	l.health.max = 100
	l.fed.max = 100
	l.stamina.current = 100
	l.health.current = 100
	l.fed.current = 100
	return l
})

type vital struct {
	max, current int32
}

// A living thing can interact with its environment. You should always resolve
// the DNA object with calling the Resolve() method.
type living struct {
	dna     *genome.DNA
	stamina vital
	health  vital
	fed     vital
	size    int32
}

// IsHealthy returns false if any of the objects vitals are at minimum.
func (l *living) IsHealthy() bool {
	cases := []vital{l.stamina, l.health, l.fed}
	for _, c := range cases {
		if c.current <= 0 {
			return false
		}
	}
	return true
}

// Eat adds the given amount to the hunger.
func (l *living) Eat(amount int32) {
	l.fed.current += amount
	if l.fed.current > l.fed.max {
		l.fed.current = l.fed.max
	}
}

// Exhaust reduces the stamina by the given amount.
func (l *living) Exhaust(amount int32) {
	l.stamina.current -= amount
}

// Deteriorate reduces the health by the given amount.
func (l *living) Deteriorate(amount int32) {
	l.health.current -= amount
}

// Grow increases the size by the given amount. This process increases the
// hunger by the same amount.
func (l *living) Grow(amount int32) error {
	if l.fed.current < amount {
		return ErrNotEnoughFood
	}
	l.size += amount
	l.fed.current -= amount
	return nil
}

// Resolve resolves the Living object.
func (l *living) Resolve() {
	l.dna.Resolve()
	l.dna = nil
	l.stamina.max = 100
	l.health.max = 100
	l.fed.max = 100
	l.stamina.current = 100
	l.health.current = 100
	l.fed.current = 100
	livingPool.Put(l)
}
