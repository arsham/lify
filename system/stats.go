package system

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	tm "github.com/buger/goterm"
	"github.com/dustin/go-humanize"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/arsham/neuragene/component"
	"github.com/arsham/neuragene/entity"
)

// Stats prints useful statistics every 2 seconds.
type Stats struct {
	entities   *entity.Manager
	controller controller
	updateTime time.Time
	dt         time.Duration
	filterTime time.Duration
	frameCount uint64
	fps        uint64
}

var _ System = (*Stats)(nil)

func (s *Stats) String() string { return "Stats" }

// setup returns an error if the entity manager is nil.
func (s *Stats) setup(c controller) error {
	s.controller = c
	s.entities = c.EntityManager()
	if s.entities == nil {
		return fmt.Errorf("%w: entity manager", ErrInvalidArgument)
	}
	if s.controller == nil {
		return fmt.Errorf("%w: controller", ErrInvalidArgument)
	}
	s.updateTime = time.Now()
	return nil
}

// update prints the stats if the last time it was printed was 2 seconds ago.
func (s *Stats) update(state component.State) error {
	if !all(state, component.StatePrintStats) {
		return nil
	}
	s.frameCount++
	s.fps++
	if time.Since(s.updateTime) >= time.Second*2 {
		s.dt = s.controller.LastFrameDuration()
		t1 := time.Now()
		s.entities.MapByMask(0b111111, func(*entity.Entity) {})
		printStats(s.entities, s)
		s.filterTime = time.Since(t1)
		s.updateTime = time.Now()
		s.fps = 0
	}
	return nil
}

func (s *Stats) draw(screen *ebiten.Image, state component.State) {
	if !all(state, component.StatePrintStats) {
		return
	}
	msg := fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS())
	ebitenutil.DebugPrint(screen, msg)
}

func printCurrentTime() {
	_, _ = tm.Println(strings.Repeat("-", 47))
	_, _ = tm.Println(format("Current Time:", time.Now().Format(time.Stamp)))
	_, _ = tm.Println(strings.Repeat("-", 47))
}

func printEngineStats(em *entity.Manager, stats *Stats) {
	_, _ = tm.Println(format("Engine Statistics:", ""))
	_, _ = tm.Println(format("Entities:", fmt.Sprintf("%d", em.Len())))
	_, _ = tm.Println(format("FilterTime:", stats.filterTime.String()))
	_, _ = tm.Println(format("FrameTime:", stats.dt.String()))
	_, _ = tm.Println(format("Total Frames:", fmt.Sprintf("%d", stats.frameCount)))
	_, _ = tm.Println(format("FPS:", fmt.Sprintf("%d", stats.fps/2)))
	_, _ = tm.Println(strings.Repeat("-", 47))
	_, _ = tm.Println()
}

func printMemoryStats() {
	var r runtime.MemStats
	runtime.ReadMemStats(&r)
	_, _ = tm.Println(format("Memory Statistics:", ""))
	_, _ = tm.Println(format("MemStats Sys", humanize.Bytes(r.Sys)))
	_, _ = tm.Println(format("Heap Allocation", humanize.Bytes(r.HeapAlloc)))
	_, _ = tm.Println(format("Heap Idle", humanize.Bytes(r.HeapIdle)))
	_, _ = tm.Println(format("Head In Use", humanize.Bytes(r.HeapInuse)))
	_, _ = tm.Println(format("Heap HeapObjects", humanize.Bytes(r.HeapObjects)))
	_, _ = tm.Println(format("Heap Released", humanize.Bytes(r.HeapReleased)))
	_, _ = tm.Println(strings.Repeat("-", 47))
}

func printRuntimeStats() {
	s := debug.GCStats{}
	debug.ReadGCStats(&s)
	numGC := s.NumGC
	totalPause := s.PauseTotal
	_, _ = tm.Println(format("Runtime Statistics:", ""))
	_, _ = tm.Println(format("GOOS GOARCH", fmt.Sprintf("%s %s", runtime.GOOS, runtime.GOARCH)))
	_, _ = tm.Println(format("NumCPU()", fmt.Sprintf("%d", runtime.NumCPU())))
	_, _ = tm.Println(format("NumCgoCall()", fmt.Sprintf("%d", runtime.NumCgoCall())))
	_, _ = tm.Println(format("NumGoroutine()", fmt.Sprintf("%d", runtime.NumGoroutine())))
	_, _ = tm.Println(format("NumGC()", fmt.Sprintf("%d", numGC)))
	_, _ = tm.Println(format("Total Pause", totalPause.String()))
	_, _ = tm.Println(format("Version()", runtime.Version()))
	_, _ = tm.Println(strings.Repeat("-", 47))
}

func printStats(em *entity.Manager, stats *Stats) {
	tm.Clear()
	tm.MoveCursor(0, 0)
	printCurrentTime()
	printRuntimeStats()
	printMemoryStats()
	printEngineStats(em, stats)
	tm.Flush()
}

func format(key, val string) string {
	return fmt.Sprintf("| %-20s | %-20s |", key, val)
}
