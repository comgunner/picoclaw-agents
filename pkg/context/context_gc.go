// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package context

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/comgunner/picoclaw/pkg/logger"
)

// ContextGC manages automatic garbage collection of workspace files.
// Runs periodically to clean up temporary and old files.
type ContextGC struct {
	workspacePath string
	ticker        *time.Ticker
	done          chan bool

	// Configuration
	tempMaxAge    time.Duration // Max age for temp files
	logMaxSize    int64         // Max size for log files
	coldThreshold int64         // Move files larger than this to cold storage

	// Statistics
	stats GCStats
	mu    sync.RWMutex
}

// GCStats tracks garbage collection statistics.
type GCStats struct {
	FilesDeleted   int64
	BytesFreed     int64
	FilesMovedCold int64
	LastRun        time.Time
	TotalRuns      int64
}

// NewContextGC creates a new context garbage collector.
func NewContextGC(workspacePath string) *ContextGC {
	return &ContextGC{
		workspacePath: workspacePath,
		ticker:        time.NewTicker(5 * time.Minute),
		done:          make(chan bool),
		tempMaxAge:    1 * time.Hour,
		logMaxSize:    10 * 1024 * 1024, // 10MB
		coldThreshold: 5 * 1024 * 1024,  // 5MB
	}
}

// Start begins the automatic garbage collection loop.
func (gc *ContextGC) Start() {
	go gc.run()
	logger.InfoCF("context", "Context GC started",
		map[string]any{
			"workspace":    gc.workspacePath,
			"interval":     "5m",
			"temp_max_age": gc.tempMaxAge,
			"log_max_size": gc.logMaxSize,
		})
}

// Stop stops the garbage collector.
func (gc *ContextGC) Stop() {
	gc.ticker.Stop()
	gc.done <- true

	gc.mu.RLock()
	stats := gc.stats
	gc.mu.RUnlock()

	logger.InfoCF("context", "Context GC stopped",
		map[string]any{
			"total_runs":    stats.TotalRuns,
			"files_deleted": stats.FilesDeleted,
			"bytes_freed":   stats.BytesFreed,
		})
}

// run is the main GC loop.
func (gc *ContextGC) run() {
	// Execute immediately on start
	gc.execute()

	for {
		select {
		case <-gc.ticker.C:
			gc.execute()
		case <-gc.done:
			return
		}
	}
}

// execute performs a single GC cycle.
func (gc *ContextGC) execute() {
	startTime := time.Now()

	logger.DebugCF("context", "GC cycle started", nil)

	// Clean temp files
	gc.cleanTempFiles()

	// Compress old logs
	gc.compressOldLogs()

	// Move large files to cold storage
	gc.moveToColdStorage()

	// Update stats
	gc.mu.Lock()
	gc.stats.LastRun = startTime
	gc.stats.TotalRuns++
	gc.mu.Unlock()

	duration := time.Since(startTime)
	gc.mu.RLock()
	stats := gc.stats // Capture stats for logging
	gc.mu.RUnlock()

	logger.InfoCF("context", "GC cycle completed",
		map[string]any{
			"duration_ms":   duration.Milliseconds(),
			"files_deleted": stats.FilesDeleted,
			"bytes_freed":   stats.BytesFreed,
			"files_cold":    stats.FilesMovedCold,
		})
}

// cleanTempFiles removes temporary files older than tempMaxAge.
func (gc *ContextGC) cleanTempFiles() {
	tempDirs := []string{
		filepath.Join(gc.workspacePath, "temp"),
		filepath.Join(gc.workspacePath, "active"),
	}

	for _, dir := range tempDirs {
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}

			if info.IsDir() {
				return nil
			}

			// Delete *.tmp files >1h old
			if strings.HasSuffix(path, ".tmp") &&
				time.Since(info.ModTime()) > gc.tempMaxAge {
				os.Remove(path)
				gc.mu.Lock()
				gc.stats.FilesDeleted++
				gc.stats.BytesFreed += info.Size()
				gc.mu.Unlock()

				logger.DebugCF("context", "Deleted temp file",
					map[string]any{
						"path": path,
						"age":  time.Since(info.ModTime()),
					})
			}

			// Delete *.log files >10MB
			if strings.HasSuffix(path, ".log") && info.Size() > gc.logMaxSize {
				os.Remove(path)
				gc.mu.Lock()
				gc.stats.FilesDeleted++
				gc.stats.BytesFreed += info.Size()
				gc.mu.Unlock()

				logger.DebugCF("context", "Deleted large log file",
					map[string]any{
						"path": path,
						"size": info.Size(),
					})
			}

			return nil
		})
	}
}

// compressOldLogs compresses log files older than 24 hours.
func (gc *ContextGC) compressOldLogs() {
	// TODO: Implement log compression using gzip
	// This is a placeholder for future implementation
}

// moveToColdStorage moves large files to cold storage.
func (gc *ContextGC) moveToColdStorage() {
	coldDir := filepath.Join(gc.workspacePath, "cold")
	os.MkdirAll(coldDir, 0o755)

	activeDir := filepath.Join(gc.workspacePath, "active")
	filepath.Walk(activeDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		// Move files >5MB to cold storage
		if info.Size() > gc.coldThreshold {
			destPath := filepath.Join(coldDir, info.Name())

			// Move file
			err := os.Rename(path, destPath)
			if err != nil {
				// If rename fails (cross-device), copy and delete
				logger.WarnCF("context", "Failed to rename file to cold storage",
					map[string]any{
						"path":  path,
						"error": err,
					})
			} else {
				gc.mu.Lock()
				gc.stats.FilesMovedCold++
				gc.mu.Unlock()

				logger.InfoCF("context", "Moved file to cold storage",
					map[string]any{
						"path": path,
						"size": info.Size(),
						"dest": destPath,
					})
			}
		}

		return nil
	})
}

// GetStats returns current GC statistics.
func (gc *ContextGC) GetStats() GCStats {
	gc.mu.RLock()
	defer gc.mu.RUnlock()
	return gc.stats
}
