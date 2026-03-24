// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package context

import (
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContextMiddleware_NewContextMiddleware(t *testing.T) {
	tmpDir := t.TempDir()
	middleware := NewContextMiddleware(tmpDir, 8192)

	assert.NotNil(t, middleware)
	assert.NotNil(t, middleware.budget)
	assert.NotNil(t, middleware.loader)
	assert.NotNil(t, middleware.gc)

	// Verify budget configuration
	current, maxTokens, _ := middleware.GetUsage()
	assert.Equal(t, int64(0), current)
	assert.Equal(t, int64(8192), maxTokens)
}

func TestContextMiddleware_StartStop(t *testing.T) {
	tmpDir := t.TempDir()
	middleware := NewContextMiddleware(tmpDir, 8192)

	// Start middleware
	middleware.Start()

	// Give it time to start
	time.Sleep(50 * time.Millisecond)

	// Stop middleware
	middleware.Stop()
}

func TestContextMiddleware_BeforeRequest(t *testing.T) {
	tmpDir := t.TempDir()
	middleware := NewContextMiddleware(tmpDir, 8192)

	// Should allow normal requests
	allowed := middleware.BeforeRequest("llm_call")
	assert.True(t, allowed)

	// AfterRequest should be callable
	middleware.AfterRequest()
}

func TestContextMiddleware_BeforeRequestBlocked(t *testing.T) {
	tmpDir := t.TempDir()
	middleware := NewContextMiddleware(tmpDir, 1000) // Small budget for testing

	// UPDATED: CanAfford now uses 100% Hard Limit (not 80% anymore).
	// The 80% Soft Limit only controls when EmergencyGC is triggered; blocking happens at 100%.

	// Test that normal requests below the hard limit go through
	allowed := middleware.BeforeRequest("llm_call") // 100 tokens
	assert.True(t, allowed, "First request should be allowed (hard limit not reached)")

	// Set the budget counter to exactly the hard limit using atomic store.
	// We bypass Charge() here so we don't trigger the async EmergencyGC goroutine
	// (which would race against our assertion by resetting the counter).
	atomic.StoreInt64(&middleware.budget.currentTokens, 1000) // 100% = hard limit

	// Now a new request (100 tokens) would push total to 1100 > 1000 (hard limit).
	// BeforeRequest should return false and dispatch an async EmergencyGC goroutine.
	blocked := middleware.BeforeRequest("llm_call")
	assert.False(t, blocked, "Request should be blocked when hard limit (100%) is reached")
}


func TestContextMiddleware_GetUsage(t *testing.T) {
	tmpDir := t.TempDir()
	middleware := NewContextMiddleware(tmpDir, 10000)

	// Initial usage
	current, maxTokens, percentage := middleware.GetUsage()
	assert.Equal(t, int64(0), current)
	assert.Equal(t, int64(10000), maxTokens)
	assert.Equal(t, 0.0, percentage)

	// Charge some tokens
	middleware.BeforeRequest("llm_call")
	middleware.BeforeRequest("file_read")

	current, maxTokens, percentage = middleware.GetUsage()
	assert.Greater(t, current, int64(0))
	assert.Equal(t, int64(10000), maxTokens)
	assert.Greater(t, percentage, 0.0)
}

func TestContextMiddleware_GetBudget(t *testing.T) {
	tmpDir := t.TempDir()
	middleware := NewContextMiddleware(tmpDir, 8192)

	budget := middleware.GetBudget()
	assert.NotNil(t, budget)

	// Verify budget configuration
	_, maxTokens, _ := budget.Usage()
	assert.Equal(t, int64(8192), maxTokens)
}

func TestContextMiddleware_GetLoader(t *testing.T) {
	tmpDir := t.TempDir()
	middleware := NewContextMiddleware(tmpDir, 8192)

	loader := middleware.GetLoader()
	assert.NotNil(t, loader)

	// Create test file and verify loader works
	testFile := filepath.Join(tmpDir, "test.txt")
	err := os.WriteFile(testFile, []byte("test content"), 0o644)
	require.NoError(t, err)

	ref, err := loader.ReferenceFile(testFile)
	require.NoError(t, err)
	assert.NotNil(t, ref)
}

func TestContextMiddleware_Integration(t *testing.T) {
	tmpDir := t.TempDir()
	middleware := NewContextMiddleware(tmpDir, 8192)
	middleware.Start()
	defer middleware.Stop()

	// Simulate multiple requests
	for i := 0; i < 10; i++ {
		allowed := middleware.BeforeRequest("llm_call")
		if !allowed {
			t.Fatal("Request blocked unexpectedly")
		}
		middleware.AfterRequest()
	}

	// Check usage
	current, maxTokens, percentage := middleware.GetUsage()
	assert.Greater(t, current, int64(0))
	assert.Equal(t, int64(8192), maxTokens)
	assert.Less(t, percentage, float64(80))
}

func TestContextMiddleware_WithLazyLoader(t *testing.T) {
	tmpDir := t.TempDir()
	middleware := NewContextMiddleware(tmpDir, 8192)

	// Create test files
	testFile1 := filepath.Join(tmpDir, "test1.txt")
	testFile2 := filepath.Join(tmpDir, "test2.txt")
	err := os.WriteFile(testFile1, []byte("Content 1"), 0o644)
	require.NoError(t, err)
	err = os.WriteFile(testFile2, []byte("Content 2"), 0o644)
	require.NoError(t, err)

	// Use loader from middleware
	loader := middleware.GetLoader()

	ref1, err := loader.ReferenceFile(testFile1)
	require.NoError(t, err)
	ref2, err := loader.ReferenceFile(testFile2)
	require.NoError(t, err)

	// Load content
	err = loader.LoadContent(ref1)
	require.NoError(t, err)

	assert.True(t, ref1.Loaded)
	assert.False(t, ref2.Loaded)
	assert.Equal(t, "Content 1", ref1.Content)
}

func TestContextMiddleware_WithWorkspaceManager(t *testing.T) {
	tmpDir := t.TempDir()

	// Create workspace manager to ensure structure
	wm := NewWorkspaceManager(tmpDir)
	_ = wm

	// Verify workspace structure exists
	dirs := []string{"active", "memory", "cold", "temp", "sessions", "state", "scripts"}
	for _, dir := range dirs {
		fullPath := filepath.Join(tmpDir, dir)
		_, err := os.Stat(fullPath)
		assert.NoError(t, err, "Directory %s should exist", dir)
	}
}

func TestContextMiddleware_EmergencyGCIntegration(t *testing.T) {
	t.Skip("Flaky test - timing dependent on goroutine execution")
	// Emergency GC is tested in TestTokenBudget_EmergencyGC
}

func TestContextMiddleware_MultipleOperations(t *testing.T) {
	tmpDir := t.TempDir()
	middleware := NewContextMiddleware(tmpDir, 10000)

	operations := []string{
		"llm_call",
		"file_read",
		"file_write",
		"tool_exec",
		"session_load",
		"memory_access",
	}

	for _, op := range operations {
		allowed := middleware.BeforeRequest(op)
		assert.True(t, allowed, "Operation %s should be allowed", op)
	}

	current, _, percentage := middleware.GetUsage()
	assert.Greater(t, current, int64(0))
	assert.Less(t, percentage, float64(80))
}

func TestContextMiddleware_AfterRequestNoOp(t *testing.T) {
	tmpDir := t.TempDir()
	middleware := NewContextMiddleware(tmpDir, 8192)

	// AfterRequest should not panic or cause errors
	assert.NotPanics(t, func() {
		middleware.AfterRequest()
	})
}

func TestContextMiddleware_DifferentBudgetSizes(t *testing.T) {
	tests := []struct {
		name      string
		maxTokens int64
	}{
		{"Small budget", 1000},
		{"Medium budget", 8192},
		{"Large budget", 16384},
		{"Extra large budget", 32768},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			middleware := NewContextMiddleware(tmpDir, tt.maxTokens)

			_, maxTokens, _ := middleware.GetUsage()
			assert.Equal(t, tt.maxTokens, maxTokens)
		})
	}
}

func TestContextMiddleware_WithGCStats(t *testing.T) {
	tmpDir := t.TempDir()
	middleware := NewContextMiddleware(tmpDir, 8192)
	middleware.Start()
	defer middleware.Stop()

	// Wait for GC to run
	time.Sleep(200 * time.Millisecond)

	// GC should have run at least once
	stats := middleware.gc.GetStats()
	assert.GreaterOrEqual(t, stats.TotalRuns, int64(1))
}

func TestContextMiddleware_ConcurrentRequests(t *testing.T) {
	tmpDir := t.TempDir()
	middleware := NewContextMiddleware(tmpDir, 100000) // Large budget

	done := make(chan bool, 20)

	// Simulate concurrent requests
	for i := 0; i < 20; i++ {
		go func() {
			for j := 0; j < 10; j++ {
				middleware.BeforeRequest("llm_call")
				middleware.AfterRequest()
			}
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 20; i++ {
		<-done
	}

	// Verify final usage
	current, _, _ := middleware.GetUsage()
	assert.Equal(t, int64(20000), current) // 20 goroutines * 10 requests * 100 tokens
}
