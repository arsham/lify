package quadtree_test

import (
	"fmt"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/arsham/neuragene/internal/geom"
	"github.com/arsham/neuragene/internal/quadtree"
)

func TestBoundsContains(t *testing.T) {
	t.Parallel()
	tcs := []struct {
		bounds *quadtree.Bounds
		x, y   float64
		want   bool
	}{{ // case_0
		bounds: quadtree.NewBounds(5, 5, 10, 10),
		x:      10,
		y:      10,
		want:   true,
	}, { // case_1
		bounds: quadtree.NewBounds(10, 10, 10, 10),
		x:      20,
		y:      20,
		want:   false,
	}, { // case_2
		bounds: quadtree.NewBounds(5, 5, 10, 10),
		x:      5,
		y:      5,
		want:   true,
	}, { // case_3
		bounds: quadtree.NewBounds(10, 10, 10, 10),
		x:      15.0001,
		y:      15.0001,
		want:   false,
	}, { // case_4
		bounds: quadtree.NewBounds(10, 10, 10, 10),
		x:      20,
		y:      10,
		want:   false,
	}}
	for i, tc := range tcs {
		tc := tc
		name := fmt.Sprintf("case #%d", i)
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			p := geom.V(tc.x, tc.y)
			if tc.want {
				assert.True(t, tc.bounds.Contains(p), "want %s be in %s", p, tc.bounds)
				return
			}
			assert.False(t, tc.bounds.Contains(p), "want %s not be in %s", p, tc.bounds)
		})
	}
}

func TestBoundsSubDivide(t *testing.T) {
	t.Parallel()
	b := quadtree.NewBounds(0, 0, 1000, 500)

	nw := b.SubDivide(quadtree.NW)
	ne := b.SubDivide(quadtree.NE)
	sw := b.SubDivide(quadtree.SW)
	se := b.SubDivide(quadtree.SE)

	tcs := map[string]struct {
		bound *quadtree.Bounds
		nw    geom.Rect
		ne    geom.Rect
		sw    geom.Rect
		se    geom.Rect
	}{
		"first divide": {
			bound: b,
			nw:    geom.R(0, 0, 500, 250),
			ne:    geom.R(500, 0, 1000, 250),
			sw:    geom.R(0, 250, 500, 500),
			se:    geom.R(500, 250, 1000, 500),
		},
		"nw divide": {
			bound: nw,
			nw:    geom.R(0, 0, 250, 125),
			ne:    geom.R(250, 0, 500, 125),
			sw:    geom.R(0, 125, 250, 250),
			se:    geom.R(250, 125, 500, 250),
		},
		"ne divide": {
			bound: ne,
			nw:    geom.R(500, 0, 750, 125),
			ne:    geom.R(750, 0, 1000, 125),
			sw:    geom.R(500, 125, 750, 250),
			se:    geom.R(750, 125, 1000, 250),
		},
		"sw divide": {
			bound: sw,
			nw:    geom.R(0, 250, 250, 375),
			ne:    geom.R(250, 250, 500, 375),
			sw:    geom.R(0, 375, 250, 500),
			se:    geom.R(250, 375, 500, 500),
		},
		"se divide": {
			bound: se,
			nw:    geom.R(500, 250, 750, 375),
			ne:    geom.R(750, 250, 1000, 375),
			sw:    geom.R(500, 375, 750, 500),
			se:    geom.R(750, 375, 1000, 500),
		},
	}
	for name, tc := range tcs {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			nw := tc.bound.SubDivide(quadtree.NW)
			ne := tc.bound.SubDivide(quadtree.NE)
			sw := tc.bound.SubDivide(quadtree.SW)
			se := tc.bound.SubDivide(quadtree.SE)
			assert.True(t, tc.nw.Eq(nw.Rect), "NW:\n want: %v\n got: %v", tc.nw, nw.Rect)
			assert.True(t, tc.ne.Eq(ne.Rect), "NE:\n want: %v\n got: %v", tc.ne, ne.Rect)
			assert.True(t, tc.sw.Eq(sw.Rect), "SW:\n want: %v\n got: %v", tc.sw, sw.Rect)
			assert.True(t, tc.se.Eq(se.Rect), "SE:\n want: %v\n got: %v", tc.se, se.Rect)
		})
	}
}
