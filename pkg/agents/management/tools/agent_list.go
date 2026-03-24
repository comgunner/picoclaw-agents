// PicoClaw - Ultra-lightweight personal AI agent
// Agent Management Tool: agent_list
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

	"github.com/comgunner/picoclaw/pkg/agents/management"
	"github.com/comgunner/picoclaw/pkg/tools"
)

// AgentListTool lists all configured agents with their current status.
type AgentListTool struct {
	registry  *management.AgentRegistry
	instances *management.InstanceRegistry
}

// NewAgentListTool creates a new agent_list tool.
func NewAgentListTool(
	registry *management.AgentRegistry,
	instances *management.InstanceRegistry,
) *AgentListTool {
	return &AgentListTool{
		registry:  registry,
		instances: instances,
	}
}

func (t *AgentListTool) Name() string { return "agent_list" }

func (t *AgentListTool) Description() string {
	return "List all registered agents with their runtime status and basic configuration."
}

func (t *AgentListTool) Parameters() map[string]any {
	return map[string]any{
		"type":       "object",
		"properties": map[string]any{},
	}
}

func (t *AgentListTool) Execute(_ context.Context, _ map[string]any) *tools.ToolResult {
	agentIDs := t.registry.ListAgentIDs()
	if len(agentIDs) == 0 {
		return tools.UserResult("📭 No agents are currently registered.")
	}

	defaultAgent := t.registry.GetDefaultAgent()
	defaultID := ""
	if defaultAgent != nil {
		defaultID = defaultAgent.ID
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("📋 **Registered Agents (%d)**\n\n", len(agentIDs)))
	sb.WriteString("| Status | ID | Name | Model | Default |\n")
	sb.WriteString("|--------|----|----|-------|---------|\n")

	for _, id := range agentIDs {
		agent, ok := t.registry.GetAgent(id)
		if !ok {
			continue
		}

		emoji := "⚪"
		if instance, hasInstance := t.instances.Get(id); hasInstance {
			switch instance.GetStatus() {
			case management.StatusActive:
				emoji = "🟢"
			case management.StatusBusy:
				emoji = "🟡"
			case management.StatusError:
				emoji = "🔴"
			default:
				emoji = "⚪"
			}
		}

		defaultMark := ""
		if id == defaultID {
			defaultMark = "⭐ yes"
		} else {
			defaultMark = "no"
		}

		modelStr := "<unset>"
		if agent.Model != nil && agent.Model.Primary != "" {
			modelStr = agent.Model.Primary
		}

		sb.WriteString(fmt.Sprintf("| %s | `%s` | %s | %s | %s |\n",
			emoji, agent.ID, agent.Name, modelStr, defaultMark))
	}

	sb.WriteString(fmt.Sprintf("\n**Default agent:** `%s`\n", defaultID))

	return tools.UserResult(sb.String())
}
