// Package food contains the logic for providing food for the organisms.
package food

import (
	"github.com/arsham/lify/internal/pool"
	"github.com/lucasb-eyer/go-colorful"
)

var herbPool = pool.NewPool(func() *Herb {
	return &Herb{
		Nutition:    100,
		MaxNutition: 100,
		GrowthRate:  100,
	}
})

// A Herb is a food that grows on the ground.
type Herb struct {
	Name        string
	Colour      colorful.Color
	Nutition    int32
	MaxNutition int32
	// GrowthRate is the rate at which the herb grows. Each Nutition is grown
	// every GrowthRate cycles.
	GrowthRate int32
}

// NewHerb returns a new Herb object. You should always resolve the object by
// calling the Resolve() method.
func NewHerb(name string) *Herb {
	h := herbPool.Get()
	h.Name = name
	return h
}

func (h *Herb) String() string {
	return h.Name
}

// Resolve returns the object to the pool.
func (h *Herb) Resolve() {
	h.Name = ""
	h.Colour = colorful.Color{}
	h.Nutition = 100
	h.MaxNutition = 100
	h.GrowthRate = 100
	herbPool.Put(h)
}

// Eaten updates the food's nutrition.
func (h *Herb) Eaten(amount int32) {
	h.Nutition -= amount
	if h.Nutition < 0 {
		h.Nutition = 0
	}
}
