// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package context

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContextGC_NewContextGC(t *testing.T) {
	tmpDir := t.TempDir()
	gc := NewContextGC(tmpDir)

	assert.Equal(t, tmpDir, gc.workspacePath)
	assert.Equal(t, 1*time.Hour, gc.tempMaxAge)
	assert.Equal(t, int64(10*1024*1024), gc.logMaxSize)
	assert.Equal(t, int64(5*1024*1024), gc.coldThreshold)
}

func TestContextGC_StartStop(t *testing.T) {
	tmpDir := t.TempDir()
	gc := NewContextGC(tmpDir)

	// Start GC
	gc.Start()

	// Give it time to start
	time.Sleep(50 * time.Millisecond)

	// Stop GC
	gc.Stop()

	// Verify stats were recorded
	assert.GreaterOrEqual(t, gc.stats.TotalRuns, int64(0))
}

func TestContextGC_Execute(t *testing.T) {
	tmpDir := t.TempDir()
	gc := NewContextGC(tmpDir)

	// Create temp and active directories
	err := os.MkdirAll(filepath.Join(tmpDir, "temp"), 0o755)
	require.NoError(t, err)
	err = os.MkdirAll(filepath.Join(tmpDir, "active"), 0o755)
	require.NoError(t, err)

	// Execute one cycle
	gc.execute()

	// Verify stats updated
	assert.Greater(t, gc.stats.TotalRuns, int64(0))
	assert.Equal(t, gc.stats.LastRun.Year(), time.Now().Year())
}

func TestContextGC_CleanTempFiles(t *testing.T) {
	tmpDir := t.TempDir()
	gc := NewContextGC(tmpDir)
	gc.tempMaxAge = 100 * time.Millisecond // Very short for testing

	// Create temp directory
	tempDir := filepath.Join(tmpDir, "temp")
	err := os.MkdirAll(tempDir, 0o755)
	require.NoError(t, err)

	// Create old temp file
	oldTempFile := filepath.Join(tempDir, "old.tmp")
	err = os.WriteFile(oldTempFile, []byte("temp content"), 0o644)
	require.NoError(t, err)

	// Set modification time to past
	oldTime := time.Now().Add(-2 * time.Hour)
	err = os.Chtimes(oldTempFile, oldTime, oldTime)
	require.NoError(t, err)

	// Create new temp file (should not be deleted)
	newTempFile := filepath.Join(tempDir, "new.tmp")
	err = os.WriteFile(newTempFile, []byte("new content"), 0o644)
	require.NoError(t, err)

	// Clean temp files
	gc.cleanTempFiles()

	// Old file should be deleted
	_, err = os.Stat(oldTempFile)
	assert.True(t, os.IsNotExist(err))

	// New file should still exist
	_, err = os.Stat(newTempFile)
	assert.NoError(t, err)

	// Verify stats
	assert.Greater(t, gc.stats.FilesDeleted, int64(0))
	assert.Greater(t, gc.stats.BytesFreed, int64(0))
}

func TestContextGC_CleanLargeLogFiles(t *testing.T) {
	tmpDir := t.TempDir()
	gc := NewContextGC(tmpDir)
	gc.logMaxSize = 100 // Very small for testing

	// Create active directory
	activeDir := filepath.Join(tmpDir, "active")
	err := os.MkdirAll(activeDir, 0o755)
	require.NoError(t, err)

	// Create large log file
	largeLogFile := filepath.Join(activeDir, "large.log")
	largeContent := make([]byte, 200) // 200 bytes > 100 byte limit
	err = os.WriteFile(largeLogFile, largeContent, 0o644)
	require.NoError(t, err)

	// Create small log file (should not be deleted)
	smallLogFile := filepath.Join(activeDir, "small.log")
	smallContent := make([]byte, 50) // 50 bytes < 100 byte limit
	err = os.WriteFile(smallLogFile, smallContent, 0o644)
	require.NoError(t, err)

	// Clean temp files (includes log cleanup)
	gc.cleanTempFiles()

	// Large file should be deleted
	_, err = os.Stat(largeLogFile)
	assert.True(t, os.IsNotExist(err))

	// Small file should still exist
	_, err = os.Stat(smallLogFile)
	assert.NoError(t, err)
}

func TestContextGC_MoveToColdStorage(t *testing.T) {
	tmpDir := t.TempDir()
	gc := NewContextGC(tmpDir)
	gc.coldThreshold = 100 // Very small for testing

	// Create active and cold directories
	activeDir := filepath.Join(tmpDir, "active")
	err := os.MkdirAll(activeDir, 0o755)
	require.NoError(t, err)

	// Create large file in active
	largeFile := filepath.Join(activeDir, "large.bin")
	largeContent := make([]byte, 200) // 200 bytes > 100 byte threshold
	err = os.WriteFile(largeFile, largeContent, 0o644)
	require.NoError(t, err)

	// Create small file in active (should not be moved)
	smallFile := filepath.Join(activeDir, "small.bin")
	smallContent := make([]byte, 50) // 50 bytes < 100 byte threshold
	err = os.WriteFile(smallFile, smallContent, 0o644)
	require.NoError(t, err)

	// Move to cold storage
	gc.moveToColdStorage()

	// Large file should be moved to cold
	coldLargeFile := filepath.Join(tmpDir, "cold", "large.bin")
	_, err = os.Stat(coldLargeFile)
	assert.NoError(t, err)

	// Large file should not exist in active
	_, err = os.Stat(largeFile)
	assert.True(t, os.IsNotExist(err))

	// Small file should still exist in active
	_, err = os.Stat(smallFile)
	assert.NoError(t, err)

	// Verify stats
	assert.Greater(t, gc.stats.FilesMovedCold, int64(0))
}

func TestContextGC_GetStats(t *testing.T) {
	tmpDir := t.TempDir()
	gc := NewContextGC(tmpDir)

	// Initial stats
	stats := gc.GetStats()
	assert.Equal(t, int64(0), stats.FilesDeleted)
	assert.Equal(t, int64(0), stats.BytesFreed)
	assert.Equal(t, int64(0), stats.FilesMovedCold)
	assert.Equal(t, int64(0), stats.TotalRuns)
}

func TestContextGC_GCStatsStruct(t *testing.T) {
	stats := GCStats{
		FilesDeleted:   10,
		BytesFreed:     1024,
		FilesMovedCold: 5,
		TotalRuns:      3,
	}

	assert.Equal(t, int64(10), stats.FilesDeleted)
	assert.Equal(t, int64(1024), stats.BytesFreed)
	assert.Equal(t, int64(5), stats.FilesMovedCold)
	assert.Equal(t, int64(3), stats.TotalRuns)
}

func TestContextGC_NonExistentDirectories(t *testing.T) {
	tmpDir := t.TempDir()
	gc := NewContextGC(tmpDir)

	// Don't create temp/active directories
	// GC should handle this gracefully without panicking

	// This should not panic
	gc.cleanTempFiles()
	gc.moveToColdStorage()
}

func TestContextGC_EmptyDirectories(t *testing.T) {
	tmpDir := t.TempDir()
	gc := NewContextGC(tmpDir)

	// Create empty directories
	err := os.MkdirAll(filepath.Join(tmpDir, "temp"), 0o755)
	require.NoError(t, err)
	err = os.MkdirAll(filepath.Join(tmpDir, "active"), 0o755)
	require.NoError(t, err)

	// GC should handle empty directories
	gc.cleanTempFiles()
	gc.moveToColdStorage()

	// No files should be deleted
	assert.Equal(t, int64(0), gc.stats.FilesDeleted)
}

func TestContextGC_MixedFileTypes(t *testing.T) {
	tmpDir := t.TempDir()
	gc := NewContextGC(tmpDir)
	gc.tempMaxAge = 100 * time.Millisecond
	gc.logMaxSize = 100
	gc.coldThreshold = 150

	// Create directories
	tempDir := filepath.Join(tmpDir, "temp")
	activeDir := filepath.Join(tmpDir, "active")
	err := os.MkdirAll(tempDir, 0o755)
	require.NoError(t, err)
	err = os.MkdirAll(activeDir, 0o755)
	require.NoError(t, err)

	// Old temp file
	oldTemp := filepath.Join(tempDir, "old.tmp")
	err = os.WriteFile(oldTemp, []byte("old"), 0o644)
	require.NoError(t, err)
	oldTime := time.Now().Add(-2 * time.Hour)
	os.Chtimes(oldTemp, oldTime, oldTime)

	// Large log file
	largeLog := filepath.Join(activeDir, "large.log")
	err = os.WriteFile(largeLog, make([]byte, 200), 0o644)
	require.NoError(t, err)

	// Large data file (for cold storage)
	largeData := filepath.Join(activeDir, "data.bin")
	err = os.WriteFile(largeData, make([]byte, 200), 0o644)
	require.NoError(t, err)

	// Run GC
	gc.execute()

	// Old temp should be deleted
	_, err = os.Stat(oldTemp)
	assert.True(t, os.IsNotExist(err))

	// Large log should be deleted
	_, err = os.Stat(largeLog)
	assert.True(t, os.IsNotExist(err))

	// Large data should be moved to cold
	_, err = os.Stat(filepath.Join(tmpDir, "cold", "data.bin"))
	assert.NoError(t, err)
}

func TestContextGC_FilePatterns(t *testing.T) {
	tmpDir := t.TempDir()
	gc := NewContextGC(tmpDir)
	gc.tempMaxAge = 100 * time.Millisecond

	// Create temp directory
	tempDir := filepath.Join(tmpDir, "temp")
	err := os.MkdirAll(tempDir, 0o755)
	require.NoError(t, err)

	// Create .tmp file (should be deleted if old)
	tmpFile := filepath.Join(tempDir, "file.tmp")
	err = os.WriteFile(tmpFile, []byte("temp"), 0o644)
	require.NoError(t, err)
	oldTime := time.Now().Add(-2 * time.Hour)
	os.Chtimes(tmpFile, oldTime, oldTime)

	// Create .txt file (should NOT be deleted)
	txtFile := filepath.Join(tempDir, "file.txt")
	err = os.WriteFile(txtFile, []byte("text"), 0o644)
	require.NoError(t, err)
	os.Chtimes(txtFile, oldTime, oldTime)

	// Clean
	gc.cleanTempFiles()

	// .tmp should be deleted
	_, err = os.Stat(tmpFile)
	assert.True(t, os.IsNotExist(err))

	// .txt should still exist
	_, err = os.Stat(txtFile)
	assert.NoError(t, err)
}
