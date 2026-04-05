// PicoClaw - Ultra-lightweight personal AI agent
// Agent Management Tool: agent_message_stats
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
	"sort"
	"strings"

	"github.com/comgunner/picoclaw/pkg/agents/management"
	"github.com/comgunner/picoclaw/pkg/tools"
)

// AgentMessageStatsTool reports aggregate statistics from the AgentMessageBus.
type AgentMessageStatsTool struct {
	messageBus *management.AgentMessageBus
}

// NewAgentMessageStatsTool creates a new agent_message_stats tool.
func NewAgentMessageStatsTool(messageBus *management.AgentMessageBus) *AgentMessageStatsTool {
	return &AgentMessageStatsTool{messageBus: messageBus}
}

func (t *AgentMessageStatsTool) Name() string { return "agent_message_stats" }

func (t *AgentMessageStatsTool) Description() string {
	return "Return aggregate statistics for the inter-agent message bus: total messages, by-type breakdown, hourly distribution, and busiest agent."
}

func (t *AgentMessageStatsTool) Parameters() map[string]any {
	return map[string]any{
		"type":       "object",
		"properties": map[string]any{},
	}
}

func (t *AgentMessageStatsTool) Execute(_ context.Context, _ map[string]any) *tools.ToolResult {
	stats := t.messageBus.GetStats()

	var sb strings.Builder
	sb.WriteString("📊 **Message Bus Statistics**\n\n")
	sb.WriteString(fmt.Sprintf("- **Total messages:** %d\n", stats.TotalMessages))
	sb.WriteString(fmt.Sprintf("- **Messages today:** %d\n", stats.MessagesToday))
	sb.WriteString(fmt.Sprintf("- **Delivery rate:** %.1f%%\n", stats.DeliveryRate*100))
	sb.WriteString(fmt.Sprintf("- **Response rate:** %.1f%%\n", stats.ResponseRate*100))
	if stats.BusiestAgent != "" {
		sb.WriteString(fmt.Sprintf("- **Busiest agent:** `%s`\n", stats.BusiestAgent))
	}
	if stats.AvgResponseTime > 0 {
		sb.WriteString(fmt.Sprintf("- **Avg response time:** %s\n", stats.AvgResponseTime.String()))
	}

	if len(stats.ByType) > 0 {
		sb.WriteString("\n**Messages by type:**\n")
		// Deterministic sort
		types := make([]string, 0, len(stats.ByType))
		for k := range stats.ByType {
			types = append(types, k)
		}
		sort.Strings(types)
		for _, typ := range types {
			sb.WriteString(fmt.Sprintf("- %s: %d\n", typ, stats.ByType[typ]))
		}
	}

	if len(stats.ByHour) > 0 {
		sb.WriteString("\n**Messages by hour (UTC):**\n")
		hours := make([]int, 0, len(stats.ByHour))
		for h := range stats.ByHour {
			hours = append(hours, h)
		}
		sort.Ints(hours)
		for _, h := range hours {
			sb.WriteString(fmt.Sprintf("- %02d:00 → %d\n", h, stats.ByHour[h]))
		}
	}

	if stats.TotalMessages == 0 {
		return tools.UserResult("📭 No messages have been exchanged yet.")
	}

	return tools.UserResult(sb.String())
}
