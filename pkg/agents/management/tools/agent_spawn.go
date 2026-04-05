// PicoClaw - Ultra-lightweight personal AI agent
// Agent Management Tool: agent_spawn
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package managementtools

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/comgunner/picoclaw/pkg/agents/management"
	"github.com/comgunner/picoclaw/pkg/tools"
)

// AgentSpawnTool registers a new agent instance in the InstanceRegistry so that
// other management tools can observe and interact with it.
// It does NOT start a full agentic loop; for that callers should use the
// existing pkg/tools.SpawnTool which delegates to SubagentManager.
type AgentSpawnTool struct {
	registry    *management.AgentRegistry
	instances   *management.InstanceRegistry
	rateLimiter *management.RateLimiter
}

// NewAgentSpawnTool creates a new agent_spawn tool.
func NewAgentSpawnTool(
	registry *management.AgentRegistry,
	instances *management.InstanceRegistry,
	rateLimiter *management.RateLimiter,
) *AgentSpawnTool {
	return &AgentSpawnTool{
		registry:    registry,
		instances:   instances,
		rateLimiter: rateLimiter,
	}
}

func (t *AgentSpawnTool) Name() string { return "agent_spawn" }

func (t *AgentSpawnTool) Description() string {
	return "Register a new agent instance for activity tracking. " +
		"The agent must already be defined in the configuration. " +
		"Use agent_monitor after spawning to observe its status."
}

func (t *AgentSpawnTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"agent_id": map[string]any{
				"type":        "string",
				"description": "The ID of the configured agent to register as active.",
			},
			"caller_id": map[string]any{
				"type":        "string",
				"description": "ID of the calling/parent agent (used for rate-limit tracking).",
			},
		},
		"required": []string{"agent_id"},
	}
}

func (t *AgentSpawnTool) Execute(_ context.Context, args map[string]any) *tools.ToolResult {
	agentID, ok := args["agent_id"].(string)
	if !ok || strings.TrimSpace(agentID) == "" {
		return tools.ErrorResult("agent_id is required and must be a non-empty string")
	}
	agentID = strings.TrimSpace(agentID)

	callerID, _ := args["caller_id"].(string)

	// Rate-limit check (keyed on caller when available)
	rateKey := agentID
	if callerID != "" {
		rateKey = callerID
	}
	if t.rateLimiter != nil && !t.rateLimiter.Check(rateKey, t.Name()) {
		return tools.ErrorResult(fmt.Sprintf(
			"rate limit exceeded for agent_spawn (%d calls/min)", management.RateLimits["agent_spawn"],
		))
	}

	// Validate agent is configured
	agent, found := t.registry.GetAgent(agentID)
	if !found {
		return tools.ErrorResult(fmt.Sprintf("agent not found in configuration: %q", agentID))
	}

	// Register in the instance registry (idempotent)
	instance := t.instances.Register(agentID)
	instance.UpdateStatus(management.StatusActive)

	out := fmt.Sprintf(
		"✅ Agent **%s** (`%s`) registered and marked active.\n"+
			"- **Model:** %s\n"+
			"- **Workspace:** %s\n"+
			"- **Registered at:** %s\n",
		agent.Name,
		agent.ID,
		modelString(agent),
		agent.Workspace,
		time.Now().Format(time.RFC3339),
	)

	return tools.UserResult(out)
}
