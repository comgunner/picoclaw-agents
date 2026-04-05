// PicoClaw - Ultra-lightweight personal AI agent
// Agent Management Tool: agent_get
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

// AgentGetTool retrieves detailed information about a specific configured agent.
type AgentGetTool struct {
	registry  *management.AgentRegistry
	instances *management.InstanceRegistry
}

// NewAgentGetTool creates a new agent_get tool.
func NewAgentGetTool(
	registry *management.AgentRegistry,
	instances *management.InstanceRegistry,
) *AgentGetTool {
	return &AgentGetTool{
		registry:  registry,
		instances: instances,
	}
}

func (t *AgentGetTool) Name() string { return "agent_get" }

func (t *AgentGetTool) Description() string {
	return "Get detailed information about a specific agent including its configuration, runtime status, and subagent settings."
}

func (t *AgentGetTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"agent_id": map[string]any{
				"type":        "string",
				"description": "The unique ID of the agent to retrieve.",
			},
		},
		"required": []string{"agent_id"},
	}
}

func (t *AgentGetTool) Execute(_ context.Context, args map[string]any) *tools.ToolResult {
	agentID, ok := args["agent_id"].(string)
	if !ok || agentID == "" {
		return tools.ErrorResult("agent_id is required and must be a non-empty string")
	}

	agent, found := t.registry.GetAgent(agentID)
	if !found {
		return tools.ErrorResult(fmt.Sprintf("agent not found: %q", agentID))
	}

	defaultAgent := t.registry.GetDefaultAgent()
	isDefault := defaultAgent != nil && agent.ID == defaultAgent.ID

	status := "inactive"
	isActive := false
	if instance, hasInstance := t.instances.Get(agentID); hasInstance {
		status = string(instance.GetStatus())
		isActive = true
	}

	// Build subagents section from optional config
	subagentInfo := map[string]any{
		"allow_agents":           []string{},
		"max_spawn_depth":        3,
		"max_children_per_agent": 5,
	}
	if agent.Subagents != nil {
		subagentInfo["allow_agents"] = agent.Subagents.AllowAgents
		if agent.Subagents.MaxSpawnDepth > 0 {
			subagentInfo["max_spawn_depth"] = agent.Subagents.MaxSpawnDepth
		}
		if agent.Subagents.MaxChildrenPerAgent > 0 {
			subagentInfo["max_children_per_agent"] = agent.Subagents.MaxChildrenPerAgent
		}
	}

	modelStr := "<not set>"
	if agent.Model != nil {
		modelStr = agent.Model.Primary
	}

	out := fmt.Sprintf(
		"## Agent: %s\n\n"+
			"- **ID:** %s\n"+
			"- **Model:** %s\n"+
			"- **Workspace:** %s\n"+
			"- **Status:** %s\n"+
			"- **Active:** %v\n"+
			"- **Is Default:** %v\n"+
			"- **Skills:** %v\n",
		agent.Name,
		agent.ID,
		modelStr,
		agent.Workspace,
		status,
		isActive,
		isDefault,
		agent.Skills,
	)

	return tools.UserResult(out)
}
