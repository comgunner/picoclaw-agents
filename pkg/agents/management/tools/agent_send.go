// PicoClaw - Ultra-lightweight personal AI agent
// Agent Management Tool: agent_send
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
	"time"

	"github.com/google/uuid"

	"github.com/comgunner/picoclaw/pkg/agents/management"
	"github.com/comgunner/picoclaw/pkg/tools"
)

// AgentSendTool sends a directed message from one agent to another via the AgentMessageBus.
type AgentSendTool struct {
	registry    *management.AgentRegistry
	messageBus  *management.AgentMessageBus
	rateLimiter *management.RateLimiter
}

// NewAgentSendTool creates a new agent_send tool.
func NewAgentSendTool(
	registry *management.AgentRegistry,
	messageBus *management.AgentMessageBus,
	rateLimiter *management.RateLimiter,
) *AgentSendTool {
	return &AgentSendTool{
		registry:    registry,
		messageBus:  messageBus,
		rateLimiter: rateLimiter,
	}
}

func (t *AgentSendTool) Name() string { return "agent_send" }

func (t *AgentSendTool) Description() string {
	return "Send a directed message from one agent to another. " +
		"The recipient's inbox is a buffered channel; use agent_receive to drain it."
}

func (t *AgentSendTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"sender_id": map[string]any{
				"type":        "string",
				"description": "The ID of the sending agent.",
			},
			"recipient_id": map[string]any{
				"type":        "string",
				"description": "The ID of the recipient agent.",
			},
			"message_type": map[string]any{
				"type":        "string",
				"description": "Message type: info | task | result | broadcast",
				"enum":        []string{"info", "task", "result", "broadcast"},
			},
			"payload": map[string]any{
				"type":        "string",
				"description": "The message payload (plain text or JSON string).",
			},
			"requires_response": map[string]any{
				"type":        "boolean",
				"description": "Whether the sender expects a reply.",
			},
		},
		"required": []string{"sender_id", "recipient_id", "message_type", "payload"},
	}
}

func (t *AgentSendTool) Execute(_ context.Context, args map[string]any) *tools.ToolResult {
	senderID, ok := args["sender_id"].(string)
	if !ok || senderID == "" {
		return tools.ErrorResult("sender_id is required and must be a non-empty string")
	}

	recipientID, ok := args["recipient_id"].(string)
	if !ok || recipientID == "" {
		return tools.ErrorResult("recipient_id is required and must be a non-empty string")
	}

	msgType, ok := args["message_type"].(string)
	if !ok || msgType == "" {
		return tools.ErrorResult("message_type is required")
	}

	payloadStr, ok := args["payload"].(string)
	if !ok || payloadStr == "" {
		return tools.ErrorResult("payload is required and must be a non-empty string")
	}

	requiresResponse, _ := args["requires_response"].(bool)

	// Rate-limit check
	if t.rateLimiter != nil && !t.rateLimiter.Check(senderID, t.Name()) {
		return tools.ErrorResult(
			fmt.Sprintf("rate limit exceeded for agent_send (%d calls/min)", management.RateLimits["agent_send"]),
		)
	}

	// Validate agents exist
	if _, found := t.registry.GetAgent(senderID); !found {
		return tools.ErrorResult(fmt.Sprintf("sender agent not found: %q", senderID))
	}
	if _, found := t.registry.GetAgent(recipientID); !found {
		return tools.ErrorResult(fmt.Sprintf("recipient agent not found: %q", recipientID))
	}

	rawPayload, err := json.Marshal(payloadStr)
	if err != nil {
		return tools.ErrorResult(fmt.Sprintf("failed to encode payload: %v", err))
	}

	now := time.Now()
	msg := management.AgentMessage{
		ID:               uuid.NewString(),
		SenderID:         senderID,
		RecipientID:      recipientID,
		MessageType:      msgType,
		Payload:          json.RawMessage(rawPayload),
		RequiresResponse: requiresResponse,
		SentAt:           now,
	}

	if err := t.messageBus.Send(msg); err != nil {
		return tools.ErrorResult(fmt.Sprintf("failed to send message: %v", err))
	}

	out := fmt.Sprintf(
		"✉️ Message sent.\n- **ID:** `%s`\n- **From:** `%s` → **To:** `%s`\n- **Type:** %s\n- **Requires response:** %v\n- **Sent at:** %s\n",
		msg.ID,
		senderID,
		recipientID,
		msgType,
		requiresResponse,
		now.Format(time.RFC3339),
	)

	return tools.UserResult(out)
}
