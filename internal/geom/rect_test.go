package geom_test

import (
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/arsham/neuragene/internal/geom"
)

func TestRectW(t *testing.T) {
	t.Parallel()
	tcs := map[string]struct {
		r    geom.Rect
		want float64
	}{
		"empty rect": {},
		"len 1 at 0": {
			r:    geom.R(0, 100, 1, 1),
			want: 1,
		},
		"len 1 at 2": {
			r:    geom.R(2, 200, 3, 3),
			want: 1,
		},
		"len 10": {
			r:    geom.R(0, 200, 10, 10),
			want: 10,
		},
	}
	for name, tc := range tcs {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := tc.r.W()
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestRectH(t *testing.T) {
	t.Parallel()
	tcs := map[string]struct {
		r    geom.Rect
		want float64
	}{
		"empty rect": {},
		"len 1 at 0": {
			r:    geom.R(100, 0, 1, 1),
			want: 1,
		},
		"len 1 at 2": {
			r:    geom.R(200, 2, 3, 3),
			want: 1,
		},
		"len 10": {
			r:    geom.R(200, 0, 10, 10),
			want: 10,
		},
	}
	for name, tc := range tcs {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := tc.r.H()
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestRectCentre(t *testing.T) {
	t.Parallel()
	tcs := map[string]struct {
		r    geom.Rect
		want geom.Vec
	}{
		"empty rect": {},
		"no width": {
			r:    geom.R(0, 0, 0, 10),
			want: geom.V(0, 5),
		},
		"no height": {
			r:    geom.R(0, 0, 10, 0),
			want: geom.V(5, 0),
		},
		"1 rect": {
			r:    geom.R(0, 0, 1, 1),
			want: geom.V(0.5, 0.5),
		},
		"2 rect": {
			r:    geom.R(0, 0, 2, 2),
			want: geom.V(1, 1),
		},
		"moved rect": {
			r:    geom.R(10, 10, 3, 3),
			want: geom.V(6.5, 6.5),
		},
	}
	for name, tc := range tcs {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := tc.r.Centre()
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestRectResized(t *testing.T) {
	t.Parallel()
	tcs := map[string]struct {
		r      geom.Rect
		size   geom.Vec
		anchor geom.Vec
		want   geom.Rect
	}{
		"empty rect": {},
		"no width": {
			r:      geom.R(0, 0, 0, 10),
			size:   geom.V(10, 10),
			anchor: geom.V(0, 0),
			want:   geom.R(0, 0, 0, 0),
		},
		"no height": {
			r:      geom.R(0, 0, 10, 0),
			size:   geom.V(10, 10),
			anchor: geom.V(0, 0),
			want:   geom.R(0, 0, 0, 0),
		},
		"1x1 by 1 at 0": {
			r:      geom.R(0, 0, 1, 1),
			size:   geom.V(1, 1),
			anchor: geom.V(0, 0),
			want:   geom.R(1, 1, 1, 1),
		},
		"1x1 by 1 at 1": {
			r:      geom.R(0, 0, 1, 1),
			size:   geom.V(1, 1),
			anchor: geom.V(1, 1),
			want:   geom.R(0, 0, 1, 1),
		},
	}
	for name, tc := range tcs {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := tc.r.Resized(tc.size, tc.anchor)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestRectMoved(t *testing.T) {
	t.Parallel()
	tcs := map[string]struct {
		r     geom.Rect
		delta geom.Vec
		want  geom.Rect
	}{
		"empty rect": {},
		"empty rect moved": {
			r:     geom.R(0, 0, 0, 0),
			delta: geom.V(10, 10),
			want:  geom.R(10, 10, 10, 10),
		},
		"rect empty move": {
			r:     geom.R(10, 10, 10, 10),
			delta: geom.V(0, 0),
			want:  geom.R(10, 10, 10, 10),
		},
		"rect moved": {
			r:     geom.R(10, 10, 20, 20),
			delta: geom.V(10, 10),
			want:  geom.R(20, 20, 30, 30),
		},
		"rect moved negative": {
			r:     geom.R(20, 20, 10, 10),
			delta: geom.V(-10, -10),
			want:  geom.R(10, 10, 0, 0),
		},
	}
	for name, tc := range tcs {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := tc.r.Moved(tc.delta)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestRectRotated(t *testing.T) {
	t.Parallel()
	tcs := map[string]struct {
		r     geom.Rect
		angle geom.Radian
		want  geom.Rect
	}{
		"empty rect": {},
		"empty rect rotated": {
			r:     geom.R(0, 0, 0, 0),
			angle: 1,
			want:  geom.R(0, 0, 0, 0),
		},
		"rect no rotation": {
			r:     geom.R(10, 10, 20, 20),
			angle: 0,
			want:  geom.R(10, 10, 20, 20),
		},
		"rect rotated": {
			r:     geom.R(10, 10, 20, 20),
			angle: 1,
			want:  geom.R(16.505843, 8.091134, 13.494157, 21.908866),
		},
		"rect rotated negative": {
			r:     geom.R(10, 10, 20, 20),
			angle: -1,
			want:  geom.R(8.091134, 16.505843, 21.908866, 13.494157),
		},
		"rect 2.6 rotation": {
			r:     geom.R(10, 10, 20, 20),
			angle: 2.6,
			want:  geom.R(21.861951, 16.706937, 8.138049, 13.293063),
		},
	}
	for name, tc := range tcs {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := tc.r.Rotated(tc.angle)
			assert.True(t, tc.want.Eq(got), "\nwant: %v\n got: %v", tc.want, got)
		})
	}
}
