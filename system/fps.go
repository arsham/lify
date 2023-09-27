package system

import (
	"time"

	"github.com/arsham/neuragene/component"
)

// FPS system handles the FPS of rendering.
type FPS struct {
	ticker time.Ticker
	Max    uint
}

var _ System = (*FPS)(nil)

func (f *FPS) String() string { return "FPS" }

// Setup sets up the FPS system without any error.
func (f *FPS) Setup(controller) error {
	if f.Max == 0 {
		f.Max = 60
	}
	f.ticker = *time.NewTicker(time.Second / time.Duration(f.Max))
	return nil
}

// Process waits for the next tick.
func (f *FPS) Process(state component.State, _ float64) {
	if !all(state, component.StateLimitFPS) {
		return
	}
	<-f.ticker.C
}
