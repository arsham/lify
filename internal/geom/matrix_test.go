package geom_test

import (
	"fmt"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/arsham/neuragene/internal/geom"
)

func TestMatrixMove(t *testing.T) {
	t.Parallel()
	tcs := map[string]struct {
		m    *geom.Matrix
		v    geom.Vec
		want *geom.Matrix
	}{
		"empty matrix": {
			m:    geom.M([]geom.Vec{}),
			v:    geom.Vec{X: 1, Y: 1},
			want: geom.M([]geom.Vec{}),
		},
		"empty vector": {
			m:    geom.M([]geom.Vec{{X: 1, Y: 1}, {X: 2, Y: 2}}),
			v:    geom.Vec{},
			want: geom.M([]geom.Vec{{X: 1, Y: 1}, {X: 2, Y: 2}}),
		},
		"moved": {
			m:    geom.M([]geom.Vec{{X: 1, Y: 2}, {X: 3, Y: 3}, {X: 5, Y: 6}}),
			v:    geom.Vec{X: 3, Y: 6},
			want: geom.M([]geom.Vec{{X: 4, Y: 8}, {X: 6, Y: 9}, {X: 8, Y: 12}}),
		},
		"moved negative direction": {
			m:    geom.M([]geom.Vec{{X: 1, Y: 2}, {X: 3, Y: 3}, {X: 5, Y: 6}}),
			v:    geom.Vec{X: -3, Y: -6},
			want: geom.M([]geom.Vec{{X: -2, Y: -4}, {X: 0, Y: -3}, {X: 2, Y: 0}}),
		},
	}
	for name, tc := range tcs {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			tc.m.Move(tc.v)
			assert.True(t, tc.want.Eq(tc.m), "want: %v, got: %v", tc.want, tc.m)
		})
	}
}

func TestMatrixRotate(t *testing.T) {
	t.Parallel()
	m := geom.M([]geom.Vec{
		{X: 1, Y: 2},
		{X: -13, Y: 3},
		{X: 5, Y: 6},
		{X: 7, Y: 8},
	})
	m.Rotate(geom.Radian(0.5))
	want := geom.M([]geom.Vec{
		{2.196003, 2.816073},
		{-10.569579, -3.018301},
		{3.788631, 8.244106},
		{4.584945, 10.958122},
	})
	assert.True(t, want.Eq(m), "\nwant: %s\n got: %s", want, m)
}

func TestMatrixCenter(t *testing.T) {
	t.Parallel()
	m := geom.M([]geom.Vec{
		{X: 1, Y: 2},
		{X: -13, Y: 3},
		{X: 5, Y: 6},
		{X: 7, Y: 8},
		{X: 20, Y: -12.9},
	})
	centre := m.Center()
	want := geom.Vec{X: 4, Y: 1.22}
	assert.True(t, want.Eq(centre), "\nwant: %s\n got: %s", want, centre)
}

func TestMatrixResized(t *testing.T) {
	t.Parallel()
	v1 := geom.V(1, 2)
	v2 := geom.V(-13, 3)
	v3 := geom.V(5, 6)
	v4 := geom.V(7, 8)

	tcs := []struct {
		m      *geom.Matrix
		size   geom.Vec
		anchor geom.Vec
		want   *geom.Matrix
	}{{ // case_0
		m:      geom.M([]geom.Vec{v1, v2, v3, v4}),
		size:   geom.V(2, 4),
		anchor: v1,
		want: geom.M([]geom.Vec{
			{X: 1, Y: 2},
			{X: -27, Y: 6},
			{X: 9, Y: 18},
			{X: 13, Y: 26},
		}),
	}, { // case_1
		m:      geom.M([]geom.Vec{v1, v2, v3, v4}),
		size:   geom.V(2, 2), // Uniform scaling
		anchor: v1,
		want: geom.M([]geom.Vec{
			{X: 1, Y: 2},
			{X: -27, Y: 4},
			{X: 9, Y: 10},
			{X: 13, Y: 14},
		}),
	}, { // case_2
		m:      geom.M([]geom.Vec{v1, v2, v3, v4}),
		size:   geom.V(0.5, 0.5), // Uniform scaling
		anchor: v3,
		want: geom.M([]geom.Vec{
			{X: 3, Y: 4},
			{X: -4, Y: 4.5},
			{X: 5, Y: 6},
			{X: 6, Y: 7},
		}),
	}, { // case_3
		m:      geom.M([]geom.Vec{v1, v2, v3, v4}),
		size:   geom.V(1, 1), // No scaling
		anchor: v2,
		want: geom.M([]geom.Vec{
			{X: 1, Y: 2},
			{X: -13, Y: 3},
			{X: 5, Y: 6},
			{X: 7, Y: 8},
		}),
	}, { // case_4
		m:      geom.M([]geom.Vec{v1, v2, v3, v4}),
		size:   geom.V(2, 1), // Scaling only in X direction
		anchor: v4,
		want: geom.M([]geom.Vec{
			{X: -5, Y: 2},
			{X: -33, Y: 3},
			{X: 3, Y: 6},
			{X: 7, Y: 8},
		}),
	}, { // case_5
		m:      geom.M([]geom.Vec{v1, v2, v3, v4}),
		size:   geom.V(-1, -1), // Uniform negative scaling
		anchor: v1,
		want: geom.M([]geom.Vec{
			{X: 1, Y: 2},
			{X: 15, Y: 1},
			{X: -3, Y: -2},
			{X: -5, Y: -4},
		}),
	}, { // case_6
		m:      geom.M([]geom.Vec{v1, v2, v3, v4}),
		size:   geom.V(2, 0.5), // Non-uniform scaling
		anchor: v2,
		want: geom.M([]geom.Vec{
			{X: 15, Y: 2.5},
			{X: -13, Y: 3},
			{X: 23, Y: 4.5},
			{X: 27, Y: 5.5},
		}),
	}, { // case_7
		m:      geom.M([]geom.Vec{v1, v2, v3, v4}),
		size:   geom.V(0, 0), // Zero scaling (no change)
		anchor: v3,
		want: geom.M([]geom.Vec{
			{X: 1, Y: 2},
			{X: -13, Y: 3},
			{X: 5, Y: 6},
			{X: 7, Y: 8},
		}),
	}, { // case_8
		m:      geom.M([]geom.Vec{v1, v2, v3, v4}),
		size:   geom.V(1.5, 0.5), // Non-uniform scaling
		anchor: v1,
		want: geom.M([]geom.Vec{
			{X: 1, Y: 2},
			{X: -20, Y: 2.5},
			{X: 7, Y: 4},
			{X: 10, Y: 5},
		}),
	}, { // case_9
		m:      geom.M([]geom.Vec{v1, v2, v3, v4}),
		size:   geom.V(-2, 3), // Non-uniform negative scaling
		anchor: v4,
		want: geom.M([]geom.Vec{
			{X: 19, Y: -10},
			{X: 47, Y: -7},
			{X: 11, Y: 2},
			{X: 7, Y: 8},
		}),
	}, { // case_10
		m:      geom.M([]geom.Vec{v1, v2, v3, v4}),
		size:   geom.V(2, 0), // Scaling in X direction only
		anchor: v3,
		want: geom.M([]geom.Vec{
			{X: -3, Y: 6},
			{X: -31, Y: 6},
			{X: 5, Y: 6},
			{X: 9, Y: 6},
		}),
	}, { // case_11
		m:      geom.M([]geom.Vec{v1, v2, v3, v4}),
		size:   geom.V(0, 2), // Scaling in Y direction only
		anchor: v3,
		want: geom.M([]geom.Vec{
			{X: 5, Y: -2},
			{X: 5, Y: 0},
			{X: 5, Y: 6},
			{X: 5, Y: 10},
		}),
	}}

	for i, tc := range tcs {
		tc := tc
		name := fmt.Sprintf("case %d", i)
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			tc.m.Resize(tc.anchor, tc.size)
			assert.True(t, tc.want.Eq(tc.m), "\nwant: %s\n got: %s", tc.want, tc.m)
		})
	}
}
