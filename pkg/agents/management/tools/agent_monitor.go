// PicoClaw - Ultra-lightweight personal AI agent
// Agent Management Tool: agent_monitor
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

// AgentMonitorTool reports real-time activity, queue state, and performance metrics
// for a specific agent instance.
type AgentMonitorTool struct {
	registry  *management.AgentRegistry
	instances *management.InstanceRegistry
}

// NewAgentMonitorTool creates a new agent_monitor tool.
func NewAgentMonitorTool(
	registry *management.AgentRegistry,
	instances *management.InstanceRegistry,
) *AgentMonitorTool {
	return &AgentMonitorTool{
		registry:  registry,
		instances: instances,
	}
}

func (t *AgentMonitorTool) Name() string { return "agent_monitor" }

func (t *AgentMonitorTool) Description() string {
	return "Monitor a running agent's activity, task queue, and performance metrics in real-time."
}

func (t *AgentMonitorTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"agent_id": map[string]any{
				"type":        "string",
				"description": "The ID of the agent to monitor.",
			},
		},
		"required": []string{"agent_id"},
	}
}

func (t *AgentMonitorTool) Execute(_ context.Context, args map[string]any) *tools.ToolResult {
	agentID, ok := args["agent_id"].(string)
	if !ok || agentID == "" {
		return tools.ErrorResult("agent_id is required and must be a non-empty string")
	}

	agent, found := t.registry.GetAgent(agentID)
	if !found {
		return tools.ErrorResult(fmt.Sprintf("agent not found: %q", agentID))
	}

	instance, hasInstance := t.instances.Get(agentID)
	if !hasInstance {
		return tools.UserResult(fmt.Sprintf("⚪ Agent **%s** (`%s`) is not currently running.",
			agent.Name, agent.ID))
	}

	metrics := instance.GetMetrics()
	queue := instance.GetQueue()
	uptime := time.Since(instance.StartedAt).Round(time.Second)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("📊 **Monitor: %s** (`%s`)\n\n", agent.Name, agent.ID))
	sb.WriteString(fmt.Sprintf("**Status:** %s\n", instance.GetStatus()))
	sb.WriteString(fmt.Sprintf("**Uptime:** %s\n", uptime))
	sb.WriteString(fmt.Sprintf("**Last Activity:** %s\n\n", metrics.LastActivity.Format(time.RFC3339)))

	sb.WriteString("**Queue:**\n")
	sb.WriteString(fmt.Sprintf("- Active tasks: %d\n", len(queue.ActiveTasks)))
	sb.WriteString(fmt.Sprintf("- Queued tasks: %d\n", len(queue.QueuedTasks)))
	sb.WriteString(fmt.Sprintf("- Completed: %d\n\n", len(queue.Completed)))

	sb.WriteString("**Metrics (today):**\n")
	sb.WriteString(fmt.Sprintf("- Completed: %d\n", metrics.CompletedToday))
	sb.WriteString(fmt.Sprintf("- Failed: %d\n", metrics.FailedToday))
	sb.WriteString(fmt.Sprintf("- Token usage: %d\n", metrics.TokenUsage))
	sb.WriteString(fmt.Sprintf("- Memory usage: %s\n", formatBytes(metrics.MemoryUsage)))
	if metrics.AvgCompletion > 0 {
		sb.WriteString(fmt.Sprintf("- Avg completion: %s\n", metrics.AvgCompletion.Round(time.Millisecond)))
	}

	if len(queue.ActiveTasks) > 0 {
		sb.WriteString("\n**Active Tasks:**\n")
		for _, task := range queue.ActiveTasks {
			sb.WriteString(fmt.Sprintf("- [%s] %s", task.ID, task.Task))
			if task.Progress != "" {
				sb.WriteString(fmt.Sprintf(" (%s)", task.Progress))
			}
			sb.WriteString("\n")
		}
	}

	return tools.UserResult(sb.String())
}

// formatBytes renders a byte count in a human-readable unit.
func formatBytes(b int64) string {
	switch {
	case b >= 1<<30:
		return fmt.Sprintf("%.2f GiB", float64(b)/(1<<30))
	case b >= 1<<20:
		return fmt.Sprintf("%.2f MiB", float64(b)/(1<<20))
	case b >= 1<<10:
		return fmt.Sprintf("%.2f KiB", float64(b)/(1<<10))
	default:
		return fmt.Sprintf("%d B", b)
	}
}
