package organism

import (
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/arsham/lify/internal/genome"
	"github.com/arsham/lify/internal/itesting"
)

func TestLiving(t *testing.T) {
	t.Parallel()
	t.Run("Eat", testLivingEat)
	t.Run("Grow", testLivingGrow)
}

func testLivingEat(t *testing.T) {
	t.Parallel()
	l := livingPool.Get()
	l.fed.current = 0
	l.fed.max = 100

	l.Eat(99)
	assert.Equal(t, 99, l.fed.current, "fed level: %d", l.fed.current)

	l.Eat(20)
	assert.Equal(t, 100, l.fed.current, "fed level: %d", l.fed.current)
}

func testLivingGrow(t *testing.T) {
	t.Parallel()
	l := livingPool.Get()
	dna := genome.NewDNAFromString(itesting.RandomDNA(30))
	dna.SetTrait(genome.IndexNutritionCunsumption, 'b') // 10 in number
	dna.SetTrait(genome.IndexGrowth, '7')               // 6 in number
	dna.SetTrait(genome.IndexMaxGrowth, 'v')            // 30 in number
	l.dna = dna

	l.fed.current = 10
	l.size = 20
	err := l.Grow()
	assert.NoError(t, err)
	assert.True(t, l.fed.current == 4, "fed level: %d", l.fed.current)
	assert.True(t, l.size == 26, "size: %d", l.size)

	err = l.Grow()
	assert.IsError(t, err, ErrNotEnoughFood)
	assert.True(t, l.fed.current == 4, "fed level: %d", l.fed.current)
	assert.True(t, l.size == 26, "size: %d", l.size)

	l.fed.current = 100
	err = l.Grow()
	assert.IsError(t, err, ErrMaxGrowth)
	assert.True(t, l.size == 26, "size: %d", l.size)
}
