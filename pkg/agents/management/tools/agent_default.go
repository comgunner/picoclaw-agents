// PicoClaw - Ultra-lightweight personal AI agent
// Agent Management Tool: agent_default
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

// AgentDefaultTool returns the default agent configuration.
type AgentDefaultTool struct {
	registry  *management.AgentRegistry
	instances *management.InstanceRegistry
}

// NewAgentDefaultTool creates a new agent_default tool.
func NewAgentDefaultTool(
	registry *management.AgentRegistry,
	instances *management.InstanceRegistry,
) *AgentDefaultTool {
	return &AgentDefaultTool{
		registry:  registry,
		instances: instances,
	}
}

func (t *AgentDefaultTool) Name() string { return "agent_default" }

func (t *AgentDefaultTool) Description() string {
	return "Return the configuration and current status of the default agent."
}

func (t *AgentDefaultTool) Parameters() map[string]any {
	return map[string]any{
		"type":       "object",
		"properties": map[string]any{},
	}
}

func (t *AgentDefaultTool) Execute(_ context.Context, _ map[string]any) *tools.ToolResult {
	agent := t.registry.GetDefaultAgent()
	if agent == nil {
		return tools.ErrorResult("no agents are configured")
	}

	status := "inactive"
	isActive := false
	if instance, hasInstance := t.instances.Get(agent.ID); hasInstance {
		status = string(instance.GetStatus())
		isActive = true
	}

	modelStr := "<not set>"
	if agent.Model != nil && agent.Model.Primary != "" {
		modelStr = agent.Model.Primary
	}

	out := fmt.Sprintf(
		"## Default Agent\n\n"+
			"- **ID:** %s\n"+
			"- **Name:** %s\n"+
			"- **Model:** %s\n"+
			"- **Workspace:** %s\n"+
			"- **Status:** %s\n"+
			"- **Active:** %v\n",
		agent.ID,
		agent.Name,
		modelStr,
		agent.Workspace,
		status,
		isActive,
	)

	return tools.UserResult(out)
}
