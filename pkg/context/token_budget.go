// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package context

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/comgunner/picoclaw/pkg/logger"
)

// TokenBudget manages token allocation and monitors usage in real-time.
// Thread-safe implementation using atomic operations.
type TokenBudget struct {
	maxTokens     int64
	currentTokens int64
	lastReset     time.Time
	lastGC        time.Time

	// isEGCRunning prevents concurrent emergency GC executions
	isEGCRunning atomic.Bool

	// Costos estimados por operación
	costs map[string]int64

	// Callback para emergency GC
	onEmergencyGC func()

	mu sync.RWMutex
}

// TokenCosts defines standard costs for common operations
var TokenCosts = map[string]int64{
	"file_read":     50,  // ~50 tokens por archivo referenciado
	"file_write":    30,  // Metadata de escritura
	"llm_call":      100, // Base cost por llamada
	"tool_exec":     20,  // Ejecución de herramienta
	"session_load":  200, // Carga de sesión completa
	"memory_access": 10,  // Acceso a memoria
}

// NewTokenBudget creates a new token budget with the specified max tokens.
// Recommended values: 8192, 16384, 32768 depending on model context window.
func NewTokenBudget(maxTokens int64) *TokenBudget {
	costs := make(map[string]int64)
	for k, v := range TokenCosts {
		costs[k] = v
	}

	return &TokenBudget{
		maxTokens:     maxTokens,
		currentTokens: 0,
		lastReset:     time.Now(),
		lastGC:        time.Time{},
		costs:         costs,
	}
}

// SetCosts configures the token costs for operations.
func (tb *TokenBudget) SetCosts(costs map[string]int64) {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	tb.costs = costs
}

// GetCost returns the token cost for a specific operation.
func (tb *TokenBudget) GetCost(operation string) int64 {
	tb.mu.RLock()
	defer tb.mu.RUnlock()
	if cost, ok := tb.costs[operation]; ok {
		return cost
	}
	// Default cost if not specified
	return 50
}

// CanAfford checks if the operation can be afforded without exceeding the hard limit (100%).
// Returns false only if the operation would exceed maxTokens entirely (hard limit / bankruptcy).
// The 80% soft limit is handled separately by Charge to trigger emergency GC proactively.
// FIXED: Using 100% as hard limit to prevent deadlock where budget gets stuck at ~79.9%.
func (tb *TokenBudget) CanAfford(operation string) bool {
	cost := tb.GetCost(operation)
	current := atomic.LoadInt64(&tb.currentTokens)

	// FIXED: Hard Limit = 100% of maxTokens (absolute bankruptcy prevention)
	return current+cost <= tb.maxTokens
}

// Charge adds the operation cost to current token usage.
// Automatically triggers emergency GC if the 80% soft limit is exceeded.
// The GC runs asynchronously in the background while the 20% buffer (80%-100%)
// allows new requests to continue being served without interruption.
func (tb *TokenBudget) Charge(operation string) {
	cost := tb.GetCost(operation)
	newTotal := atomic.AddInt64(&tb.currentTokens, cost)

	// Soft Limit: 80% para disparar GC de manera preventiva y asíncrona.
	// Gracias al Hard Limit en CanAfford (100%), ahora las operaciones pueden cruzar
	// esta barrera del 80% y activar el EmergencyGC correctamente.
	threshold := tb.maxTokens * 8 / 10
	if newTotal > threshold {
		logger.WarnCF("context", "Token budget exceeded 80% soft limit - triggering emergency GC",
			map[string]any{
				"current":    newTotal,
				"max":        tb.maxTokens,
				"percentage": float64(newTotal) * 100 / float64(tb.maxTokens),
			})
		go tb.triggerEmergencyGC()
	}
}

// Reset resets the token counter to zero.
func (tb *TokenBudget) Reset() {
	atomic.StoreInt64(&tb.currentTokens, 0)
	tb.mu.Lock()
	tb.lastReset = time.Now()
	tb.mu.Unlock()

	logger.InfoCF("context", "Token budget reset",
		map[string]any{
			"max_tokens": tb.maxTokens,
		})
}

// Usage returns current token usage statistics.
func (tb *TokenBudget) Usage() (current, max int64, percentage float64) {
	current = atomic.LoadInt64(&tb.currentTokens)
	max = tb.maxTokens
	percentage = 0
	if max > 0 {
		percentage = float64(current) * 100 / float64(max)
	}
	return current, max, percentage
}

// triggerEmergencyGC performs emergency garbage collection.
func (tb *TokenBudget) triggerEmergencyGC() {
	// Ensure only one EGC runs at a time
	if !tb.isEGCRunning.CompareAndSwap(false, true) {
		return
	}
	defer tb.isEGCRunning.Store(false)

	tb.mu.Lock()
	tb.lastGC = time.Now()
	callback := tb.onEmergencyGC
	tb.mu.Unlock()

	// Callback si está registrado
	if callback != nil {
		callback()
	}

	// Resetear contadores
	tb.Reset()

	logger.InfoCF("context", "Emergency GC completed",
		map[string]any{
			"tokens_freed": atomic.LoadInt64(&tb.currentTokens),
		})
}

// GetLastReset returns the time of the last token budget reset.
func (tb *TokenBudget) GetLastReset() time.Time {
	tb.mu.RLock()
	defer tb.mu.RUnlock()
	return tb.lastReset
}

// GetLastGC returns the time of the last emergency GC.
func (tb *TokenBudget) GetLastGC() time.Time {
	tb.mu.RLock()
	defer tb.mu.RUnlock()
	return tb.lastGC
}

// SetEmergencyGCCallback registers a callback to be invoked on emergency GC.
func (tb *TokenBudget) SetEmergencyGCCallback(callback func()) {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	tb.onEmergencyGC = callback
}
