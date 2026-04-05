// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package context

import (
	"github.com/comgunner/picoclaw/pkg/logger"
)

// ContextMiddleware wraps the agent loop with context management.
type ContextMiddleware struct {
	budget *TokenBudget
	loader *LazyLoader
	gc     *ContextGC
}

// NewContextMiddleware creates a new context middleware.
func NewContextMiddleware(workspace string, maxTokens int64) *ContextMiddleware {
	budget := NewTokenBudget(maxTokens)
	loader := NewLazyLoader(workspace)
	gc := NewContextGC(workspace)

	// Set up emergency GC callback
	budget.SetEmergencyGCCallback(func() {
		logger.WarnCF("context", "Emergency GC triggered by token budget", nil)
	})

	return &ContextMiddleware{
		budget: budget,
		loader: loader,
		gc:     gc,
	}
}

// Start begins context management (starts GC loop).
func (cm *ContextMiddleware) Start() {
	cm.gc.Start()
}

// Stop stops context management.
func (cm *ContextMiddleware) Stop() {
	cm.gc.Stop()
}

// BeforeRequest is called before each LLM request.
// Returns false if the request should be blocked due to token budget (hard limit = 100%).
// FIXED: When the hard limit is reached, triggers EmergencyGC immediately as a safety fallback
// to ensure the token counter is eventually reset even in extreme saturation scenarios.
func (cm *ContextMiddleware) BeforeRequest(operation string) bool {
	if !cm.budget.CanAfford(operation) {
		current, max, pct := cm.budget.Usage()
		logger.WarnCF("context", "Request blocked: token budget at hard limit (100%) - triggering emergency GC",
			map[string]any{
				"operation":  operation,
				"current":    current,
				"max":        max,
				"percentage": pct,
			})
		// Safety fallback: trigger EmergencyGC so the budget resets and
		// future requests can be processed again after the GC completes.
		go cm.budget.triggerEmergencyGC()
		return false
	}

	cm.budget.Charge(operation)
	return true
}

// AfterRequest is called after each LLM response.
func (cm *ContextMiddleware) AfterRequest() {
	// Update statistics, log usage, etc.
}

// GetBudget returns the token budget instance.
func (cm *ContextMiddleware) GetBudget() *TokenBudget {
	return cm.budget
}

// GetLoader returns the lazy loader instance.
func (cm *ContextMiddleware) GetLoader() *LazyLoader {
	return cm.loader
}

// GetUsage returns current token usage.
func (cm *ContextMiddleware) GetUsage() (current, maxTokens int64, percentage float64) {
	return cm.budget.Usage()
}
