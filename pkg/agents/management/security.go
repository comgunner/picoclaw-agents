// PicoClaw - Ultra-lightweight personal AI agent
// Agent Management - Security & Rate Limiting
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package management

import (
	"sync"
	"time"
)

// ToolPermissions maps a tool name to the set of agent roles that may call it.
// The wildcard "*" means any agent is allowed.
var ToolPermissions = map[string][]string{
	"agent_get":           {"*"},
	"agent_list":          {"*"},
	"agent_can_spawn":     {"*"},
	"agent_default":       {"*"},
	"agent_monitor":       {"project_manager", "admin"},
	"agent_send":          {"*"},
	"agent_receive":       {"*"},
	"agent_respond":       {"*"},
	"agent_broadcast":     {"project_manager", "admin"},
	"agent_message_log":   {"admin"},
	"agent_message_stats": {"admin"},
	"agent_spawn":         {"project_manager", "senior_dev"},
}

// RateLimits defines the maximum number of calls per minute for a tool.
// Tools not listed here have no enforced rate limit.
var RateLimits = map[string]int{
	"agent_send":      30,
	"agent_broadcast": 5,
	"agent_spawn":     10,
}

// RateLimiter enforces per-agent, per-tool call rate limits using a sliding window.
// It is safe for concurrent use.
type RateLimiter struct {
	calls map[string][]time.Time //  key = "<agentID>:<toolName>"
	mu    sync.Mutex
}

// NewRateLimiter creates a new RateLimiter with an empty call log.
func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		calls: make(map[string][]time.Time),
	}
}

// Check returns true when agentID is allowed to call toolName under the configured
// rate limit, and records the call for future checks.
// Tools without an explicit rate limit always return true.
func (rl *RateLimiter) Check(agentID, toolName string) bool {
	limit, ok := RateLimits[toolName]
	if !ok {
		return true // no configured limit
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()

	key := agentID + ":" + toolName
	now := time.Now()
	windowStart := now.Add(-time.Minute)

	// Retain only calls within the sliding window
	current := rl.calls[key]
	valid := current[:0] // reuse backing array
	for _, t := range current {
		if t.After(windowStart) {
			valid = append(valid, t)
		}
	}

	if len(valid) >= limit {
		rl.calls[key] = valid
		return false
	}

	rl.calls[key] = append(valid, now)
	return true
}

// IsToolAllowed checks whether role is permitted to call toolName.
// An empty or missing role is treated as unprivileged.
func IsToolAllowed(toolName, role string) bool {
	allowed, ok := ToolPermissions[toolName]
	if !ok {
		return false // unknown tool → deny
	}

	for _, r := range allowed {
		if r == "*" || r == role {
			return true
		}
	}
	return false
}
