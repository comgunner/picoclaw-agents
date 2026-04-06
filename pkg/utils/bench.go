package utils

import (
	"fmt"
	"io"
	"os/exec"
	"runtime"
	"time"
)

// MemSnapshot captures a point-in-time memory usage snapshot.
type MemSnapshot struct {
	AllocMB      uint64
	TotalAllocMB uint64
	SysMB        uint64
	NumGC        uint32
	Goroutines   int
	Timestamp    time.Time
}

// CaptureMemSnapshot reads current memory stats.
func CaptureMemSnapshot() MemSnapshot {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return MemSnapshot{
		AllocMB:      m.Alloc / 1024 / 1024,
		TotalAllocMB: m.TotalAlloc / 1024 / 1024,
		SysMB:        m.Sys / 1024 / 1024,
		NumGC:        m.NumGC,
		Goroutines:   runtime.NumGoroutine(),
		Timestamp:    time.Now(),
	}
}

// PrintMemDiff prints a human-readable comparison between two snapshots.
func PrintMemDiff(before, after MemSnapshot) {
	fmt.Printf("Memory delta:\n")
	fmt.Printf("  Alloc:      %d → %d MB (%+d)\n", before.AllocMB, after.AllocMB, int64(after.AllocMB)-int64(before.AllocMB))
	fmt.Printf("  TotalAlloc: %d → %d MB (%+d)\n", before.TotalAllocMB, after.TotalAllocMB, int64(after.TotalAllocMB)-int64(before.TotalAllocMB))
	fmt.Printf("  Sys:        %d → %d MB (%+d)\n", before.SysMB, after.SysMB, int64(after.SysMB)-int64(before.SysMB))
	fmt.Printf("  GC cycles:  %d → %d (%+d)\n", before.NumGC, after.NumGC, int(after.NumGC)-int(before.NumGC))
	fmt.Printf("  Goroutines: %d → %d (%+d)\n", before.Goroutines, after.Goroutines, after.Goroutines-before.Goroutines)
}

// BenchmarkStartup measures cold-start time and memory of a command.
func BenchmarkStartup(cmd string, args []string) (time.Duration, MemSnapshot, error) {
	start := time.Now()

	c := exec.Command(cmd, args...)
	c.Stdout = io.Discard
	c.Stderr = io.Discard
	if err := c.Run(); err != nil {
		return 0, MemSnapshot{}, err
	}

	elapsed := time.Since(start)
	after := CaptureMemSnapshot()
	return elapsed, after, nil
}
