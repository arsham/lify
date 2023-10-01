package geom_test

import (
	"testing"

	"github.com/arsham/neuragene/internal/geom"
)

var aRect geom.Rect

func BenchmarkRectResized(b *testing.B) {
	r := geom.R(100, 100, 500, 500)
	anchor := geom.V(50, 50)
	size := geom.V(200, 200)
	b.ResetTimer()
	b.ReportAllocs()
	b.Run("Serial", func(b *testing.B) {
		b.ReportMetric(float64(b.N), "Iterations")
		for i := 0; i < b.N; i++ {
			aRect = r.Resized(anchor, size)
		}
	})
	b.Run("Parallel", func(b *testing.B) {
		b.ReportMetric(float64(b.N), "Iterations")
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				aRect = r.Resized(anchor, size)
			}
		})
	})
}
