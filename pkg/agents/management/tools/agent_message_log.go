// PicoClaw - Ultra-lightweight personal AI agent
// Agent Management Tool: agent_message_log
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

// AgentMessageLogTool retrieves the message audit log for a specific agent.
type AgentMessageLogTool struct {
	messageBus *management.AgentMessageBus
}

// NewAgentMessageLogTool creates a new agent_message_log tool.
func NewAgentMessageLogTool(messageBus *management.AgentMessageBus) *AgentMessageLogTool {
	return &AgentMessageLogTool{messageBus: messageBus}
}

func (t *AgentMessageLogTool) Name() string { return "agent_message_log" }

func (t *AgentMessageLogTool) Description() string {
	return "Retrieve the message delivery audit log for a specific agent (sent and received). " +
		"Useful for debugging inter-agent communication."
}

func (t *AgentMessageLogTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"agent_id": map[string]any{
				"type":        "string",
				"description": "Filter log entries by agent ID (sender or recipient). Empty returns all entries.",
			},
			"limit": map[string]any{
				"type":        "integer",
				"description": "Maximum number of log entries to return (default 20, max 100).",
			},
			"since_minutes": map[string]any{
				"type":        "integer",
				"description": "Only return entries from the last N minutes. Omit for all history.",
			},
		},
	}
}

func (t *AgentMessageLogTool) Execute(_ context.Context, args map[string]any) *tools.ToolResult {
	agentID, _ := args["agent_id"].(string)

	limit := 20
	if v, ok := args["limit"].(float64); ok && v > 0 {
		limit = int(v)
		if limit > 100 {
			limit = 100
		}
	}

	var since *time.Time
	if minutes, ok := args["since_minutes"].(float64); ok && minutes > 0 {
		t2 := time.Now().Add(-time.Duration(minutes) * time.Minute)
		since = &t2
	}

	logs := t.messageBus.GetLogs(agentID, limit, since)
	if len(logs) == 0 {
		filterDesc := ""
		if agentID != "" {
			filterDesc = fmt.Sprintf(" for agent `%s`", agentID)
		}
		return tools.UserResult(fmt.Sprintf("📭 No log entries found%s.", filterDesc))
	}

	var sb strings.Builder
	headerAgent := agentID
	if headerAgent == "" {
		headerAgent = "all agents"
	}
	sb.WriteString(fmt.Sprintf("📋 **Message Log** (%s) — last %d entries\n\n", headerAgent, len(logs)))
	sb.WriteString("| # | ID | From | To | Type | Status | Time |\n")
	sb.WriteString("|---|----|----|-----|------|--------|------|\n")

	for i, entry := range logs {
		shortID := entry.MessageID
		if len(shortID) > 8 {
			shortID = shortID[:8] + "…"
		}
		sb.WriteString(fmt.Sprintf("| %d | `%s` | `%s` | `%s` | %s | %s | %s |\n",
			i+1,
			shortID,
			entry.SenderID,
			entry.RecipientID,
			entry.Type,
			entry.Status,
			entry.Timestamp.Format("15:04:05"),
		))
	}

	return tools.UserResult(sb.String())
}
