package brain

import (
	stdrand "math/rand"
	"testing"

	"github.com/alecthomas/assert/v2"
)

var rand = stdrand.New(stdrand.NewSource(1))

func TestNode(t *testing.T) {
	t.Parallel()
	t.Run("HasConnectionFrom", testNodeHasConnectionFrom)
}

func testNodeHasConnectionFrom(t *testing.T) {
	t.Parallel()
	// a -> b -> c
	//   \-> d -> e
	//  f/
	a := &Node{ID: 1}
	b := &Node{ID: 2}
	c := &Node{ID: 3}
	d := &Node{ID: 4}
	e := &Node{ID: 5}
	f := &Node{ID: 6}
	e.incomming = []*Connection{{inNode: e, outNode: d}}
	d.incomming = []*Connection{{inNode: d, outNode: a}, {inNode: d, outNode: f}}
	c.incomming = []*Connection{{inNode: c, outNode: b}}
	b.incomming = []*Connection{{inNode: b, outNode: a}}

	testCases := map[string]struct {
		from *Node
		to   *Node
		want bool
	}{
		"a to b": {from: a, to: b, want: false},
		"b to a": {from: b, to: a, want: true},
		"a to c": {from: a, to: c, want: false},
		"c to a": {from: c, to: a, want: true},
		"a to e": {from: a, to: e, want: false},
		"e to a": {from: e, to: a, want: true},
		"a to f": {from: a, to: f, want: false},
		"f to a": {from: f, to: a, want: false},
		"c to f": {from: c, to: f, want: false},
		"f to c": {from: f, to: c, want: false},
	}
	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			has := tc.from.hasConnectionFrom(tc.to)
			assert.Equal(t, tc.want, has)
		})
	}
}

func TestNEAT(t *testing.T) {
	t.Parallel()
	t.Run("NewNEAT", testNEATNewNEAT)
	t.Run("Predict", testNEATPredict)
}

func testNEATNewNEAT(t *testing.T) {
	t.Parallel()

	neat := NewNEAT(2, 1, 0, rand)
	assert.Equal(t, 3, len(neat.nodes))

	neat = NewNEAT(8, 3, 0, rand)
	assert.Equal(t, 11, len(neat.nodes))
}

func testNEATPredict(t *testing.T) {
	t.Parallel()
	neat := NewNEAT(2, 1, 0, rand)
	v, err := neat.Predict([]float64{2, 3})
	assert.NoError(t, err)
	want := []float64{-0.37510}
	if !almostEqual(v, want, 0.0001) {
		t.Errorf("Predicted output %v does not match expected output %v", v, want)
	}

	neat = NewNEAT(8, 3, 0, rand)
	v, err = neat.Predict([]float64{2, 3, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6})
	assert.NoError(t, err)
	want = []float64{-1.117835, -0.53426, 0.387237}
	if !almostEqual(v, want, 0.0001) {
		t.Errorf("Predicted output %v does not match expected output %v", v, want)
	}
}
