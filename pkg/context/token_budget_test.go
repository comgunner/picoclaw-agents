// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package context

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTokenBudget_CanAfford(t *testing.T) {
	budget := NewTokenBudget(1000) // Smaller budget for testing

	// Should afford small operation initially
	assert.True(t, budget.CanAfford("file_read"))

	// Charge to 500 tokens (below 80% threshold of 800)
	for i := 0; i < 10; i++ {
		budget.Charge("file_read")
	}

	// Should still afford operations
	assert.True(t, budget.CanAfford("file_read"))
	assert.True(t, budget.CanAfford("tool_exec"))
}

func TestTokenBudget_EmergencyGC(t *testing.T) {
	budget := NewTokenBudget(1000)
	gcTriggered := make(chan bool, 1)

	budget.SetEmergencyGCCallback(func() {
		gcTriggered <- true
	})

	// Charge until >80% (threshold 800)
	// Add 9 calls (900 tokens)
	for i := 0; i < 9; i++ {
		budget.Charge("llm_call")
	}

	// Wait for async GC trigger
	select {
	case <-gcTriggered:
		// GC triggered
	case <-time.After(1 * time.Second):
		t.Fatal("GC not triggered in time")
	}

	// Give a bit of time for Reset to complete
	time.Sleep(50 * time.Millisecond)

	assert.Equal(t, int64(0), atomic.LoadInt64(&budget.currentTokens))
}

func TestTokenBudget_Reset(t *testing.T) {
	budget := NewTokenBudget(8192)

	// Charge some tokens
	for i := 0; i < 10; i++ {
		budget.Charge("file_read")
	}

	current, _, _ := budget.Usage()
	assert.Greater(t, current, int64(0))

	// Reset
	budget.Reset()

	current, _, _ = budget.Usage()
	assert.Equal(t, int64(0), current)
}

func TestTokenBudget_Usage(t *testing.T) {
	budget := NewTokenBudget(10000)

	// Initial usage
	current, maxTokens, percentage := budget.Usage()
	assert.Equal(t, int64(0), current)
	assert.Equal(t, int64(10000), maxTokens)
	assert.Equal(t, 0.0, percentage)

	// Charge 500 tokens
	budget.SetCosts(map[string]int64{
		"test_op": 500,
	})
	budget.Charge("test_op")

	current, maxTokens, percentage = budget.Usage()
	assert.Equal(t, int64(500), current)
	assert.Equal(t, int64(10000), maxTokens)
	assert.InDelta(t, 5.0, percentage, 0.1)
}

func TestTokenBudget_GetCost(t *testing.T) {
	budget := NewTokenBudget(8192)

	// Default cost
	assert.Equal(t, int64(50), budget.GetCost("unknown_op"))

	// Custom cost
	budget.SetCosts(map[string]int64{
		"custom_op": 150,
	})
	assert.Equal(t, int64(150), budget.GetCost("custom_op"))
}

func TestTokenBudget_ThresholdCalculation(t *testing.T) {
	budget := NewTokenBudget(10000)

	// UPDATED: CanAfford now uses 100% Hard Limit (not 80%).
	// The 80% Soft Limit only affects when EmergencyGC is triggered via Charge.

	// Test: Soft Limit (80%) does NOT block CanAfford
	// Charge 7900 tokens (79%) — CanAfford should return true for operations
	// that don't exceed 100% (10000 tokens total).
	budget.SetCosts(map[string]int64{
		"big_op": 7900,
	})
	assert.True(t, budget.CanAfford("big_op"))

	// Charge 7900 tokens — this will trigger async EmergencyGC at 80% Soft Limit
	budget.Charge("big_op")

	// At 7900/10000 tokens, small operations (200 tokens) push us to 8100/10000 (81%).
	// With the new Hard Limit (100%), this is STILL affordable (8100 <= 10000).
	budget.SetCosts(map[string]int64{
		"small_op": 200,
	})
	// FIXED: 7900 + 200 = 8100 <= 10000 (hard limit), so CanAfford should return TRUE.
	// The Soft Limit (80%) only controls GC triggering, not request blocking.
	assert.True(t, budget.CanAfford("small_op"))

	// Test: Hard Limit (100%) DOES block CanAfford
	// An operation that would push total beyond maxTokens must be blocked.
	budget.SetCosts(map[string]int64{
		"huge_op": 2200, // 7900 + 2200 = 10100 > 10000 (hard limit)
	})
	assert.False(t, budget.CanAfford("huge_op"))
}

func TestTokenBudget_ConcurrentAccess(t *testing.T) {
	budget := NewTokenBudget(100000)
	done := make(chan bool)

	// Start multiple goroutines charging tokens
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				budget.Charge("concurrent_op")
			}
			done <- true
		}()
	}

	// Wait for all goroutines to finish
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify final count
	current, _, _ := budget.Usage()
	assert.Equal(t, int64(50000), current) // 10 goroutines * 100 ops * 50 tokens
}

func TestTokenBudget_LastResetTime(t *testing.T) {
	budget := NewTokenBudget(8192)
	initialReset := budget.GetLastReset()

	// Wait a bit
	time.Sleep(10 * time.Millisecond)

	// Reset
	budget.Reset()

	// Verify lastReset was updated
	assert.True(t, budget.GetLastReset().After(initialReset))
}

func TestTokenBudget_LastGC(t *testing.T) {
	budget := NewTokenBudget(1000)
	gcTriggered := make(chan bool, 1)

	budget.SetEmergencyGCCallback(func() {
		gcTriggered <- true
	})

	// Trigger emergency GC
	for i := 0; i < 20; i++ {
		budget.Charge("llm_call")
	}

	// Wait for GC
	<-gcTriggered
	time.Sleep(50 * time.Millisecond)

	// Verify lastGC was updated
	assert.True(t, budget.GetLastGC().After(time.Time{}))
}
