// PicoClaw - Ultra-lightweight personal AI agent
// Agent Management Tool: agent_broadcast
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
	"time"

	"github.com/google/uuid"

	"github.com/comgunner/picoclaw/pkg/agents/management"
	"github.com/comgunner/picoclaw/pkg/tools"
)

// AgentBroadcastTool sends a message to all configured agents simultaneously.
type AgentBroadcastTool struct {
	registry    *management.AgentRegistry
	messageBus  *management.AgentMessageBus
	rateLimiter *management.RateLimiter
}

// NewAgentBroadcastTool creates a new agent_broadcast tool.
func NewAgentBroadcastTool(
	registry *management.AgentRegistry,
	messageBus *management.AgentMessageBus,
	rateLimiter *management.RateLimiter,
) *AgentBroadcastTool {
	return &AgentBroadcastTool{
		registry:    registry,
		messageBus:  messageBus,
		rateLimiter: rateLimiter,
	}
}

func (t *AgentBroadcastTool) Name() string { return "agent_broadcast" }

func (t *AgentBroadcastTool) Description() string {
	return "Broadcast a message from a sender to all other configured agents simultaneously."
}

func (t *AgentBroadcastTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"sender_id": map[string]any{
				"type":        "string",
				"description": "The ID of the broadcasting agent.",
			},
			"message_type": map[string]any{
				"type":        "string",
				"description": "Message type: info | task | broadcast",
				"enum":        []string{"info", "task", "broadcast"},
			},
			"payload": map[string]any{
				"type":        "string",
				"description": "The broadcast payload (plain text or JSON).",
			},
		},
		"required": []string{"sender_id", "message_type", "payload"},
	}
}

func (t *AgentBroadcastTool) Execute(_ context.Context, args map[string]any) *tools.ToolResult {
	senderID, ok := args["sender_id"].(string)
	if !ok || senderID == "" {
		return tools.ErrorResult("sender_id is required")
	}

	msgType, ok := args["message_type"].(string)
	if !ok || msgType == "" {
		return tools.ErrorResult("message_type is required")
	}

	payloadStr, ok := args["payload"].(string)
	if !ok || payloadStr == "" {
		return tools.ErrorResult("payload is required")
	}

	// Rate-limit check
	if t.rateLimiter != nil && !t.rateLimiter.Check(senderID, t.Name()) {
		return tools.ErrorResult(fmt.Sprintf("rate limit exceeded for agent_broadcast (%d calls/min)",
			management.RateLimits["agent_broadcast"]))
	}

	if _, found := t.registry.GetAgent(senderID); !found {
		return tools.ErrorResult(fmt.Sprintf("sender agent not found: %q", senderID))
	}

	rawPayload, err := json.Marshal(payloadStr)
	if err != nil {
		return tools.ErrorResult(fmt.Sprintf("failed to encode payload: %v", err))
	}

	now := time.Now()
	allIDs := t.registry.ListAgentIDs()

	var sent []string
	var failed []string

	for _, recipientID := range allIDs {
		if recipientID == senderID {
			continue // don't broadcast to self
		}

		msg := management.AgentMessage{
			ID:               uuid.NewString(),
			SenderID:         senderID,
			RecipientID:      recipientID,
			MessageType:      msgType,
			Payload:          json.RawMessage(rawPayload),
			RequiresResponse: false,
			SentAt:           now,
		}

		if err := t.messageBus.Send(msg); err != nil {
			failed = append(failed, fmt.Sprintf("`%s` (%v)", recipientID, err))
		} else {
			sent = append(sent, fmt.Sprintf("`%s`", recipientID))
		}
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("📢 **Broadcast complete** (from `%s`)\n\n", senderID))
	sb.WriteString(fmt.Sprintf("- **Type:** %s\n", msgType))
	sb.WriteString(fmt.Sprintf("- **Delivered to (%d):** %s\n", len(sent), strings.Join(sent, ", ")))
	if len(failed) > 0 {
		sb.WriteString(fmt.Sprintf("- ⚠️ **Failed (%d):** %s\n", len(failed), strings.Join(failed, ", ")))
	}
	sb.WriteString(fmt.Sprintf("- **Sent at:** %s\n", now.Format(time.RFC3339)))

	return tools.UserResult(sb.String())
}
