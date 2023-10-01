package geom_test

import (
	"fmt"
	"math"
	"testing"
	"testing/quick"

	"github.com/alecthomas/assert/v2"
	"github.com/arsham/neuragene/internal/geom"
)

func TestVecXY(t *testing.T) {
	t.Parallel()
	f := func(x, y float64) bool {
		v := geom.V(x, y)
		x2, y2 := v.XY()
		return x == x2 && y == y2
	}
	assert.NoError(t, quick.Check(f, nil))
}

func TestVecIsZero(t *testing.T) {
	t.Parallel()
	v := geom.V(0, 0)
	assert.True(t, v.IsZero())
	v = geom.V(1, 0)
	assert.False(t, v.IsZero())
	v = geom.V(0, 1)
	assert.False(t, v.IsZero())
	v = geom.V(1, 1)
	assert.False(t, v.IsZero())
}

func TestVecAdd(t *testing.T) {
	t.Parallel()
	f := func(x1, y1, x2, y2 float64) bool {
		v1 := geom.V(x1, y1)
		v2 := geom.V(x2, y2)
		v3 := v1.Add(v2)
		return v3.X == x1+x2 && v3.Y == y1+y2
	}
	assert.NoError(t, quick.Check(f, nil))
}

func TestVecSub(t *testing.T) {
	t.Parallel()
	f := func(x1, y1, x2, y2 float64) bool {
		v1 := geom.V(x1, y1)
		v2 := geom.V(x2, y2)
		v3 := v1.Sub(v2)
		return v3.X == x1-x2 && v3.Y == y1-y2
	}
	assert.NoError(t, quick.Check(f, nil))
}

func TestVecScaled(t *testing.T) {
	t.Parallel()
	f := func(x, y, c float64) bool {
		v := geom.V(x, y)
		v2 := v.Scaled(c)
		return v2.X == x*c && v2.Y == y*c
	}
	assert.NoError(t, quick.Check(f, nil))
}

func TestVecScaledXY(t *testing.T) {
	t.Parallel()
	f := func(x1, y1, x2, y2 float64) bool {
		v1 := geom.V(x1, y1)
		v2 := geom.V(x2, y2)
		v3 := v1.ScaledXY(v2)
		return v3.X == x1*x2 && v3.Y == y1*y2
	}
	assert.NoError(t, quick.Check(f, nil))
}

func TestVecEq(t *testing.T) {
	t.Parallel()
	tcs := []struct {
		a, b geom.Vec
		want bool
	}{{ // case_0
		a:    geom.V(0, 0),
		b:    geom.V(0, 0),
		want: true,
	}, { // case_1
		a:    geom.V(0, 0),
		b:    geom.V(1, 0),
		want: false,
	}, { // case_2
		a:    geom.V(0.000001, 0.000001),
		b:    geom.V(0, 0),
		want: false,
	}, { // case_3
		a:    geom.V(0.000001, 0.000001),
		b:    geom.V(0.000001, 0.000001),
		want: true,
	}, { // case_4
		a:    geom.V(0.000001, 0.000001),
		b:    geom.V(0.000001, 0.000002),
		want: false,
	}, { // case_5
		a:    geom.V(1.000000000001, 1.0000000001),
		b:    geom.V(1, 1),
		want: true,
	}}
	for i, tc := range tcs {
		name := fmt.Sprintf("case_%d", i)
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := tc.a.Eq(tc.b)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestVecLerp(t *testing.T) {
	t.Parallel()
	tcs := []struct {
		a, b geom.Vec
		t    float64
		want geom.Vec
	}{{ // case_0
		a:    geom.V(0, 0),
		b:    geom.V(10, 10),
		t:    0.5,
		want: geom.V(5, 5),
	}, { // case_1
		a:    geom.V(0, 0),
		b:    geom.V(10, 10),
		t:    0.1,
		want: geom.V(1, 1),
	}, { // case_2
		a:    geom.V(0, 0),
		b:    geom.V(10, 10),
		t:    0.9,
		want: geom.V(9, 9),
	}, { // case_3
		a:    geom.V(0, 0),
		b:    geom.V(10, 10),
		t:    1,
		want: geom.V(10, 10),
	}, { // case_4
		a:    geom.V(0, 0),
		b:    geom.V(10, 10),
		t:    0,
		want: geom.V(0, 0),
	}, { // case_5
		a:    geom.V(10, 10),
		b:    geom.V(0, 0),
		t:    0.5,
		want: geom.V(5, 5),
	}, { // case_6
		a:    geom.V(10, 10),
		b:    geom.V(0, 0),
		t:    0.1,
		want: geom.V(9, 9),
	}, { // case_7
		a:    geom.V(10, 10),
		b:    geom.V(0, 0),
		t:    0.9,
		want: geom.V(1, 1),
	}, { // case_8
		a:    geom.V(10, 10),
		b:    geom.V(0, 0),
		t:    1,
		want: geom.V(0, 0),
	}, { // case_9
		a:    geom.V(10, 10),
		b:    geom.V(0, 0),
		t:    0,
		want: geom.V(10, 10),
	}, { // case_10
		a:    geom.V(10, 10),
		b:    geom.V(10, 10),
		t:    0.5,
		want: geom.V(10, 10),
	}, { // case_11
		a:    geom.V(3, 7),
		b:    geom.V(18, -20),
		t:    0.5,
		want: geom.V(10.5, -6.5),
	}, { // case_12
		a:    geom.V(3, 7),
		b:    geom.V(18, -20),
		t:    0.1,
		want: geom.V(4.5, 4.3),
	}, { // case_13
		a:    geom.V(3, 7),
		b:    geom.V(18, -20),
		t:    0.9,
		want: geom.V(16.5, -17.3),
	}, { // case_14
		a:    geom.V(3, 7),
		b:    geom.V(18, -20),
		t:    0.28,
		want: geom.V(7.2, -0.56),
	}}
	for i, tc := range tcs {
		name := fmt.Sprintf("case_%d", i)
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := geom.Lerp(tc.a, tc.b, tc.t)
			assert.True(t, tc.want.Eq(got), "want: %v, got: %v", tc.want, got)
		})
	}
}

var ret geom.Vec

func BenchmarkLerp(b *testing.B) {
	b.Run("Sequential", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			ret = geom.Lerp(geom.V(3, 7), geom.V(18, -20), 0.28)
		}
	})
	b.Run("Parallel", func(b *testing.B) {
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				ret = geom.Lerp(geom.V(3, 7), geom.V(18, -20), 0.28)
			}
		})
	})
}

func TestVecAngle(t *testing.T) {
	t.Parallel()
	tcs := []struct {
		v    geom.Vec
		want geom.Radian
	}{{ // case_0
		v:    geom.V(0, 0),
		want: 0,
	}, { // case_1
		v:    geom.V(1, 0),
		want: 0,
	}, { // case_2
		v:    geom.V(0, 1),
		want: math.Pi / 2,
	}, { // case_3
		v:    geom.V(-1, 0),
		want: math.Pi,
	}, { // case_4
		v:    geom.V(0, -1),
		want: -math.Pi / 2,
	}, { // case_5
		v:    geom.V(1, 1),
		want: math.Pi / 4,
	}, { // case_6
		v:    geom.V(-1, 1),
		want: math.Pi * 3 / 4,
	}, { // case_7
		v:    geom.V(-1, -1),
		want: -math.Pi * 3 / 4,
	}, { // case_8
		v:    geom.V(1, -1),
		want: -math.Pi / 4,
	}, { // case_9
		v:    geom.V(12, 32),
		want: geom.Radian(1.2120256565243244),
	}}
	for i, tc := range tcs {
		name := fmt.Sprintf("case_%d", i)
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := tc.v.Angle()
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestVecRotated(t *testing.T) {
	t.Parallel()
	tcs := []struct {
		v     geom.Vec
		angle geom.Radian
		want  geom.Vec
	}{{ // case_0
		v:     geom.V(0, 0),
		angle: 0,
		want:  geom.V(0, 0),
	}, { // case_1
		v:     geom.V(0, 0),
		angle: 15,
		want:  geom.V(0, 0),
	}, { // case_2
		v:     geom.V(1, 0),
		angle: 0,
		want:  geom.V(1, 0),
	}, { // case_3
		v:     geom.V(1, 0),
		angle: 15,
		want:  geom.V(-0.759687912858821, 0.650287840157117),
	}, { // case_4
		v:     geom.V(15, 20),
		angle: 15,
		want:  geom.V(-24.401075496024656, -5.439440654819672),
	}, { // case_5
		v:     geom.V(15, 20),
		angle: 30,
		want:  geom.V(22.074404230170997, -11.735445363641246),
	}, { // case_6
		v:     geom.V(15, 20),
		angle: -45,
		want:  geom.V(24.897900322948314, -2.257113091657182),
	}}
	for i, tc := range tcs {
		name := fmt.Sprintf("case_%d", i)
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := tc.v.Rotated(tc.angle)
			assert.True(t, tc.want.Eq(got), "want: %v, got: %v", tc.want, got)
		})
	}
}

func BenchmarkVectRotated(b *testing.B) {
	b.Run("Sequential", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			v := geom.V(3, 7)
			ret = v.Rotated(geom.Radian(10))
		}
	})
	b.Run("Parallel", func(b *testing.B) {
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				v := geom.V(3, 7)
				ret = v.Rotated(geom.Radian(10))
			}
		})
	})
}

func BenchmarkVecRotatedByVec(b *testing.B) {
	b.Run("Sequential", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			v := geom.V(3, 7)
			other := geom.RadToVec(10)
			ret = v.RotatedByVec(other)
		}
	})
	b.Run("Parallel", func(b *testing.B) {
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				v := geom.V(3, 7)
				other := geom.RadToVec(10)
				ret = v.RotatedByVec(other)
			}
		})
	})
}

func TestVecRotatedByVec(t *testing.T) {
	t.Parallel()
	tcs := []struct {
		v     geom.Vec
		other geom.Vec
		want  geom.Vec
	}{{ // case_0
		v:     geom.V(0, 0),
		other: geom.V(0, 0),
		want:  geom.V(0, 0),
	}, { // case_1
		v:     geom.V(0, 0),
		other: geom.V(1, 0),
		want:  geom.V(0, 0),
	}, { // case_2
		v:     geom.V(0, 0),
		other: geom.V(0, 1),
		want:  geom.V(0, 0),
	}, { // case_3
		v:     geom.V(1, 0),
		other: geom.V(0, 0),
		want:  geom.V(1, 0),
	}, { // case_4
		v:     geom.V(1, 0),
		other: geom.V(1, 0),
		want:  geom.V(1, 0),
	}, { // case_5
		v:     geom.V(1, 0),
		other: geom.V(0, 1),
		want:  geom.V(0, 1),
	}, { // case_6
		v:     geom.V(1, 0),
		other: geom.V(1, 1),
		want:  geom.V(0.7071067811865476, 0.7071067811865476),
	}}
	for i, tc := range tcs {
		name := fmt.Sprintf("case_%d", i)
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := tc.v.RotatedByVec(tc.other)
			if tc.want.X == got.X && tc.want.Y == got.Y {
				return
			}
			assert.True(t, tc.want.Eq(got), "want: %v, got: %v", tc.want, got)
		})
	}
}
