// PicoClaw - Ultra-lightweight personal AI agent
// Agent Management - Registry Extensions
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package management

import (
	"sync"

	"github.com/comgunner/picoclaw/pkg/config"
)

// AgentRegistry extends the base config.Agents with runtime management capabilities.
// It is safe for concurrent use.
type AgentRegistry struct {
	cfg              *config.Config
	mu               sync.RWMutex
	spawnPermissions map[string][]string // parent → allowed children
}

// NewAgentRegistry creates a new agent registry backed by cfg.
func NewAgentRegistry(cfg *config.Config) *AgentRegistry {
	return &AgentRegistry{
		cfg:              cfg,
		spawnPermissions: make(map[string][]string),
	}
}

// GetAgent returns a copy of the agent configuration for the given agentID.
// The second return value is false when the agent is not found.
func (r *AgentRegistry) GetAgent(agentID string) (*config.AgentConfig, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.cfg == nil || r.cfg.Agents.List == nil {
		return nil, false
	}

	for _, agent := range r.cfg.Agents.List {
		if agent.ID == agentID {
			copy := agent // avoid returning a pointer into a slice element
			return &copy, true
		}
	}

	return nil, false
}

// ListAgentIDs returns all registered agent IDs in config order.
func (r *AgentRegistry) ListAgentIDs() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.cfg == nil || r.cfg.Agents.List == nil {
		return []string{}
	}

	ids := make([]string, 0, len(r.cfg.Agents.List))
	for _, agent := range r.cfg.Agents.List {
		ids = append(ids, agent.ID)
	}

	return ids
}

// GetDefaultAgent returns the default agent configuration.
// It first looks for an agent with Default=true, then falls back to the first agent.
// Returns nil if no agents are configured.
func (r *AgentRegistry) GetDefaultAgent() *config.AgentConfig {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.cfg == nil || len(r.cfg.Agents.List) == 0 {
		return nil
	}

	// Prefer explicit default
	for _, agent := range r.cfg.Agents.List {
		if agent.Default {
			copy := agent
			return &copy
		}
	}

	// Fallback: first agent
	copy := r.cfg.Agents.List[0]
	return &copy
}

// CanSpawnSubagent checks whether parentID is allowed to spawn an agent with childID.
//
// Logic:
//  1. If an explicit allowlist exists for parentID, child must appear in it.
//  2. Otherwise both agents must exist in the config (open policy).
func (r *AgentRegistry) CanSpawnSubagent(parentID, childID string) bool {
	if allowed, ok := r.spawnPermissions[parentID]; ok {
		for _, id := range allowed {
			if id == childID {
				return true
			}
		}
		return false
	}

	// Open policy: both must be registered
	_, parentExists := r.GetAgent(parentID)
	_, childExists := r.GetAgent(childID)

	return parentExists && childExists
}

// SetSpawnPermission configures which child agents parentID may spawn.
// Passing an empty slice restricts parentID from spawning anything.
func (r *AgentRegistry) SetSpawnPermission(parentID string, childIDs []string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.spawnPermissions[parentID] = childIDs
}

// AgentCount returns the total number of configured agents.
func (r *AgentRegistry) AgentCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.cfg == nil {
		return 0
	}
	return len(r.cfg.Agents.List)
}
