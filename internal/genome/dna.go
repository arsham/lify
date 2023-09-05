// Package genome contains the logic for the DNA of the organisms. It is the
// most important package of the game, as it contains the logic for the
// evolution of the organisms.
package genome

import (
	"math"
	"math/rand"
)

// A DNA is the DNA of an organism. It contains the genes of the organism,
// which are the characteristics of the organism.
type DNA struct {
	traits []rune
}

// NewDNA returns a new DNA object, preserving n empty traits. You should
// always resolve the DNA object with calling the Resolve() method.
func NewDNA(n int) *DNA {
	return &DNA{
		traits: make([]rune, n),
	}
}

// NewDNAFromString returns a new DNA object, using the given string. You
// should always resolve the DNA object with calling the Resolve() method.
func NewDNAFromString(s string) *DNA {
	return &DNA{
		traits: []rune(s),
	}
}

// Resolve resolves the DNA object.
func (d *DNA) Resolve() {}

func (d *DNA) String() string {
	return string(d.traits)
}

// SetTrait sets the trait at the given index to the given trait.
func (d *DNA) SetTrait(index int, trait rune) {
	if index >= 0 && index < len(d.traits) {
		d.traits[index] = trait
	}
}

// TraitAt returns the trait at the given index. It will return -1 if the
// index is out of bounds.
func (d *DNA) TraitAt(index int) rune {
	if index >= 0 && index < len(d.traits) {
		return d.traits[index]
	}
	return -1
}

// CalculateDifference calculates the difference between the DNA and the other
// DNA. It will return a percentage of difference between the two DNA objects.
// If the trait lengths are different, or the DNAs have differences in more
// than 3 places it will return 100% difference.
func (d *DNA) CalculateDifference(other *DNA) float64 {
	if len(d.traits) != len(other.traits) {
		return 100.0 // DNA length mismatch, 100% diff
	}
	diff := 0
	total := (1 + 9 + 26 + 26) * len(d.traits) // 1-9, a-z, A-Z, 1 null
	places := 0

	for i := 0; i < len(d.traits); i++ {
		if d.traits[i] != 0 && other.traits[i] != 0 {
			d := int(d.traits[i]) - int(other.traits[i])
			if d != 0 {
				places++
				diff += d
			}
		}
	}
	if places > 3 {
		return 100
	}

	p := float64(diff) / float64(total) * 100.0
	return math.Abs(p)
}

// TraitStrength returns the strength of the trait at the given index. A higher
// strength means a higher chance of the trait being expressed. If the index is
// out of bounds, it will return -1.
func (d *DNA) TraitStrength(index int) int {
	return int(d.TraitAt(index))
}

// IsCompatibleWith returns whether the DNA is compatible with the other DNA
// for reproduction.
func (d *DNA) IsCompatibleWith(other *DNA) bool {
	return d.CalculateDifference(other) < 0.04
}

// CreateOffspring creates an offspring from the two given parents. It uses the
// DNA of the parents to create a new DNA object. It will randomly one of the
// traits that have the least occurrence in the parents. The mutation rate is 1
// per length of the traits, and randomly applied 3 out of 100 times.
func CreateOffspring(p1, p2 *DNA) *DNA {
	c := NewDNA(len(p1.traits))
	patterns := make(map[rune]int, len(p1.traits))
	locations := make(map[rune][]int, len(p1.traits))
	for i := 0; i < len(p1.traits); i++ {
		t1 := p1.traits[i]
		t2 := p2.traits[i]
		if t1 > t2 {
			t2 = t1
		}
		c.traits[i] = t2
		patterns[t2]++
		locations[t2] = append(locations[t2], i)
	}

	if rand.Intn(100) < 3 {
		least := len(p1.traits)
		change := c.traits[0]
		for k, v := range patterns {
			if v < least {
				least = v
				change = k
			}
		}
		location := locations[change]
		var offset int
		switch {
		case change < 'b':
			offset = 1
		case change > 'y':
			offset = -1
		default:
			offset = rand.Intn(2)*2 - 1
		}
		change += rune(offset)
		c.traits[location[0]] = change
	}
	return c
}
