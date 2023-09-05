package organism

import (
	"testing"

	"github.com/alecthomas/assert/v2"
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

	l.fed.current = 10
	err := l.Grow(10)
	assert.NoError(t, err)

	l.fed.current = 10
	err = l.Grow(11)
	assert.IsError(t, err, ErrNotEnoughFood)
}
