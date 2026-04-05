// PicoClaw - Ultra-lightweight personal AI agent
// Agent Management Tools - Shared Helpers
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package managementtools

import "github.com/comgunner/picoclaw/pkg/config"

// modelString returns a human-readable model string from an AgentConfig.
func modelString(agent *config.AgentConfig) string {
	if agent == nil || agent.Model == nil {
		return "<not set>"
	}
	if agent.Model.Primary != "" {
		return agent.Model.Primary
	}
	return "<not set>"
}
