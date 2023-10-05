package system

import (
	"bytes"
	"cmp"
	"fmt"
	"runtime"
	"runtime/debug"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/olekukonko/tablewriter"
	"golang.org/x/term"

	"github.com/arsham/neuragene/internal/component"
	"github.com/arsham/neuragene/internal/entity"
)

type reports interface {
	avgCalc() time.Duration
	String() string
}

// Stats prints useful statistics every 2 seconds.
type Stats struct {
	entities *entity.Manager
	// output is a buffer to write the stats to.
	output *bytes.Buffer
	// updateTicker is a ticker to print the stats every given ExecuteEvery
	// time.
	updateTicker *time.Ticker
	controller   controller
	// stats holds the stats for each system.
	stats map[string]time.Duration
	// reports holds the reports for each system.
	reports []reports
	// ExecuteEvery is the time to wait before printing the stats.
	ExecuteEvery time.Duration
	// lastUpdateDuration is the time it took for the last update.
	lastUpdateDuration time.Duration
	// lastDrawDuration is the time it took for the last draw.
	lastDrawDuration time.Duration
	// lastDuration is the time it took for the last update of the Stats
	// system.
	lastDuration time.Duration
	// frameCount is the number of frames that have been drawn.
	frameCount uint64
	// once is used to clear the screen once.
	once sync.Once
	// inTerminal is true if the program is running in a terminal.
	inTerminal bool
}

var _ System = (*Stats)(nil)

func (s *Stats) String() string { return "Stats" }

// setup returns an error if the entity manager or the controller is nil.
func (s *Stats) setup(c controller) error {
	if term.IsTerminal(0) {
		s.inTerminal = true
	}
	s.controller = c
	s.entities = c.EntityManager()
	if s.entities == nil {
		return fmt.Errorf("%w: entity manager", ErrInvalidArgument)
	}
	if s.controller == nil {
		return fmt.Errorf("%w: controller", ErrInvalidArgument)
	}
	if s.ExecuteEvery == 0 {
		s.ExecuteEvery = 2 * time.Second
	}
	s.stats = make(map[string]time.Duration, 10)
	s.output = &bytes.Buffer{}
	s.updateTicker = time.NewTicker(s.ExecuteEvery)
	return nil
}

// update prints the stats for each tick specified with ExecuteEvery duration.
// If the state is not StatePrintStats or is not running in terminal, it will
// return immediately.
func (s *Stats) update(state component.State) error {
	started := time.Now()
	defer func() {
		s.lastDuration = time.Since(started)
	}()
	if !all(state, component.StatePrintStats) || !s.inTerminal {
		return nil
	}

	s.frameCount++
	for _, r := range s.reports {
		s.stats[r.String()] += r.avgCalc()
	}
	select {
	case <-s.updateTicker.C:
		s.once.Do(func() {
			// Clear the screen.
			fmt.Print("\033[2J")
			// Move the cursor to the top left.
			fmt.Print("\033[0;0H")
		})
		s.lastUpdateDuration = s.controller.LastUpdateDuration()
		s.lastDrawDuration = s.controller.LastDrawDuration()
		s.printStats()
		// We clear the stats so the next tick span would be calculated
		// correctly.
		clear(s.stats)
	default:
	}
	return nil
}

// draw prints the TPS and FPS on the screen.
func (s *Stats) draw(screen *ebiten.Image, state component.State) {
	if !all(state, component.StatePrintStats) {
		return
	}
	msg := fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS())
	ebitenutil.DebugPrint(screen, msg)
}

// wideSpace is a string with 12 spaces for filling the second column of the
// table.
var (
	wideSpace = strings.Repeat(" ", 12)
	wideDash  = strings.Repeat("-", 12)
)

// printStats prints the stats on the screen.
func (s *Stats) printStats() {
	// Moving to the top left of the screen.
	fmt.Print("\033[0;0H")
	s.output.Reset()

	table := tablewriter.NewWriter(s.output)
	table.SetColumnAlignment([]int{
		tablewriter.ALIGN_LEFT, tablewriter.ALIGN_RIGHT,
	})
	table.SetAutoMergeCells(true)
	table.Append([]string{"Runtime Stats  ", wideSpace})
	table.Append([]string{"---------------", wideDash})
	table.AppendBulk(runtimeStats())

	table.Append([]string{"---------------", wideDash})
	table.Append([]string{"Memory Stats   ", wideSpace})
	table.Append([]string{"---------------", wideDash})
	table.AppendBulk(memoryStats())

	table.Append([]string{"---------------", wideDash})
	table.Append([]string{"Engine Stats   ", wideSpace})
	table.Append([]string{"---------------", wideDash})
	table.AppendBulk(s.engineStats())

	table.Append([]string{"---------------", wideDash})
	table.Append([]string{"System Stats   ", wideSpace})
	table.Append([]string{"---------------", wideDash})
	data, total := s.systemStats()
	table.AppendBulk(data)
	table.SetFooter([]string{"Total", total.String()})
	table.Render()
	fmt.Print(s.output.String())
}

// engineStats returns the engine statistics.
func (s *Stats) engineStats() [][]string {
	return [][]string{
		{"Entities:", fmt.Sprintf("%d", s.entities.Len())},
		{"DrawTime:", s.lastDrawDuration.String()},
		{"UpdateTime:", s.lastUpdateDuration.String()},
		{"Total Frames:", fmt.Sprintf("%d", s.frameCount)},
		{"ActualFPS:", fmt.Sprintf("%.2f", ebiten.ActualFPS())},
		{"ActualTPS:", fmt.Sprintf("%.2f", ebiten.ActualTPS())},
	}
}

// memoryStats returns the memory statistics.
func memoryStats() [][]string {
	var r runtime.MemStats
	runtime.ReadMemStats(&r)
	return [][]string{
		{"Total Sys Mem", humanize.Bytes(r.Sys)},
		{"Heap Allocation", humanize.Bytes(r.HeapAlloc)},
		{"Heap Idle", humanize.Bytes(r.HeapIdle)},
		{"Heap In Use", humanize.Bytes(r.HeapInuse)},
		{"Heap Objects", humanize.Bytes(r.HeapObjects)},
		{"Heap Released", humanize.Bytes(r.HeapReleased)},
	}
}

// runtimeStats returns the runtime statistics.
func runtimeStats() [][]string {
	s := debug.GCStats{}
	debug.ReadGCStats(&s)
	return [][]string{
		{"Version()", runtime.Version()},
		{"CPUs", strconv.Itoa(runtime.NumCPU())},
		{"Goroutines", strconv.Itoa(runtime.NumGoroutine())},
		{"Cgo Calls", strconv.Itoa(int(runtime.NumCgoCall()))},
		{"GC Calls", strconv.Itoa(int(s.NumGC))},
		{"Pauses", s.PauseTotal.String()},
	}
}

// systemStats returns the system statistics.
func (s *Stats) systemStats() ([][]string, time.Duration) {
	var total time.Duration
	type value struct {
		name string
		dur  time.Duration
	}
	values := make([]value, 0, len(s.stats))
	tps := ebiten.ActualTPS()
	for name, avg := range s.stats {
		avg /= time.Duration(tps)
		total += avg
		values = append(values, value{name, avg})
	}
	slices.SortFunc(values, func(a, b value) int {
		return cmp.Compare(a.name, b.name)
	})

	ret := make([][]string, 0, len(values))
	for _, t := range values {
		ret = append(ret, []string{t.name, t.dur.String()})
	}
	return ret, total
}

// avgCalc returns the amount of time it took for the last update.
func (s *Stats) avgCalc() time.Duration {
	return s.lastDuration
}
