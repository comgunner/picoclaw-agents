// PicoClaw - Ultra-lightweight personal AI agent
// Agent Management Tool: agent_can_spawn
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

	"github.com/comgunner/picoclaw/pkg/agents/management"
	"github.com/comgunner/picoclaw/pkg/tools"
)

// AgentCanSpawnTool checks whether one agent is authorized to spawn another.
type AgentCanSpawnTool struct {
	registry *management.AgentRegistry
}

// NewAgentCanSpawnTool creates a new agent_can_spawn tool.
func NewAgentCanSpawnTool(registry *management.AgentRegistry) *AgentCanSpawnTool {
	return &AgentCanSpawnTool{registry: registry}
}

func (t *AgentCanSpawnTool) Name() string { return "agent_can_spawn" }

func (t *AgentCanSpawnTool) Description() string {
	return "Check whether a parent agent is authorized to spawn a specific child agent."
}

func (t *AgentCanSpawnTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"parent_agent_id": map[string]any{
				"type":        "string",
				"description": "The ID of the parent (spawning) agent.",
			},
			"child_agent_id": map[string]any{
				"type":        "string",
				"description": "The ID of the child agent to be spawned.",
			},
		},
		"required": []string{"parent_agent_id", "child_agent_id"},
	}
}

func (t *AgentCanSpawnTool) Execute(_ context.Context, args map[string]any) *tools.ToolResult {
	parentID, ok := args["parent_agent_id"].(string)
	if !ok || parentID == "" {
		return tools.ErrorResult("parent_agent_id is required and must be a non-empty string")
	}

	childID, ok := args["child_agent_id"].(string)
	if !ok || childID == "" {
		return tools.ErrorResult("child_agent_id is required and must be a non-empty string")
	}

	// Validate both agents exist
	if _, found := t.registry.GetAgent(parentID); !found {
		return tools.ErrorResult(fmt.Sprintf("parent agent not found: %q", parentID))
	}
	if _, found := t.registry.GetAgent(childID); !found {
		return tools.ErrorResult(fmt.Sprintf("child agent not found: %q", childID))
	}

	allowed := t.registry.CanSpawnSubagent(parentID, childID)

	var msg string
	if allowed {
		msg = fmt.Sprintf("✅ Agent `%s` **can** spawn agent `%s`.", parentID, childID)
	} else {
		msg = fmt.Sprintf("🚫 Agent `%s` is **not** authorized to spawn agent `%s`.", parentID, childID)
	}

	return tools.UserResult(msg)
}
