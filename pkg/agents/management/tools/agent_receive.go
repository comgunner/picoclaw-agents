// PicoClaw - Ultra-lightweight personal AI agent
// Agent Management Tool: agent_receive
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package managementtools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/comgunner/picoclaw/pkg/agents/management"
	"github.com/comgunner/picoclaw/pkg/tools"
)

// AgentReceiveTool drains pending messages from an agent's inbox on the AgentMessageBus.
type AgentReceiveTool struct {
	registry   *management.AgentRegistry
	messageBus *management.AgentMessageBus
}

// NewAgentReceiveTool creates a new agent_receive tool.
func NewAgentReceiveTool(
	registry *management.AgentRegistry,
	messageBus *management.AgentMessageBus,
) *AgentReceiveTool {
	return &AgentReceiveTool{
		registry:   registry,
		messageBus: messageBus,
	}
}

func (t *AgentReceiveTool) Name() string { return "agent_receive" }

func (t *AgentReceiveTool) Description() string {
	return "Receive pending messages from an agent's inbox. " +
		"Optionally filter by message type. Returns all queued messages and empties the inbox."
}

func (t *AgentReceiveTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"agent_id": map[string]any{
				"type":        "string",
				"description": "The ID of the agent whose inbox to receive from.",
			},
			"types": map[string]any{
				"type":        "array",
				"description": "Optional list of message types to filter (info, task, result, broadcast). Empty means receive all.",
				"items": map[string]any{
					"type": "string",
				},
			},
		},
		"required": []string{"agent_id"},
	}
}

func (t *AgentReceiveTool) Execute(_ context.Context, args map[string]any) *tools.ToolResult {
	agentID, ok := args["agent_id"].(string)
	if !ok || agentID == "" {
		return tools.ErrorResult("agent_id is required and must be a non-empty string")
	}

	if _, found := t.registry.GetAgent(agentID); !found {
		return tools.ErrorResult(fmt.Sprintf("agent not found: %q", agentID))
	}

	// Parse optional type filter
	var typeFilter []string
	if rawTypes, ok := args["types"].([]any); ok {
		for _, v := range rawTypes {
			if s, ok := v.(string); ok && s != "" {
				typeFilter = append(typeFilter, s)
			}
		}
	}

	messages, err := t.messageBus.Receive(agentID, typeFilter)
	if err != nil {
		return tools.ErrorResult(fmt.Sprintf("failed to receive messages: %v", err))
	}

	if len(messages) == 0 {
		filterDesc := ""
		if len(typeFilter) > 0 {
			filterDesc = fmt.Sprintf(" (filter: %s)", strings.Join(typeFilter, ", "))
		}
		return tools.UserResult(fmt.Sprintf("📭 No pending messages for `%s`%s.", agentID, filterDesc))
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("📬 **%d message(s) received for `%s`**\n\n", len(messages), agentID))

	for i, msg := range messages {
		sb.WriteString(fmt.Sprintf("**Message %d** (`%s`)\n", i+1, msg.ID))
		sb.WriteString(fmt.Sprintf("- **From:** `%s`\n", msg.SenderID))
		sb.WriteString(fmt.Sprintf("- **Type:** %s\n", msg.MessageType))
		sb.WriteString(fmt.Sprintf("- **Requires response:** %v\n", msg.RequiresResponse))
		sb.WriteString(fmt.Sprintf("- **Sent at:** %s\n", msg.SentAt.Format("2006-01-02 15:04:05 UTC")))

		// Pretty-print payload
		var prettyPayload any
		if json.Unmarshal(msg.Payload, &prettyPayload) == nil {
			if prettyBytes, err := json.MarshalIndent(prettyPayload, "", "  "); err == nil {
				sb.WriteString(fmt.Sprintf("- **Payload:**\n```json\n%s\n```\n", prettyBytes))
			}
		} else {
			sb.WriteString(fmt.Sprintf("- **Payload:** %s\n", string(msg.Payload)))
		}
		sb.WriteString("\n")
	}

	// Mark messages as read
	for _, msg := range messages {
		_ = t.messageBus.MarkAsRead(msg.ID)
	}

	return tools.UserResult(sb.String())
}
