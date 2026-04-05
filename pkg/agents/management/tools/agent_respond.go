// PicoClaw - Ultra-lightweight personal AI agent
// Agent Management Tool: agent_respond
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

// AgentRespondTool sends a "result" reply from a recipient back to the original sender,
// referencing the original message ID for traceability.
type AgentRespondTool struct {
	registry   *management.AgentRegistry
	messageBus *management.AgentMessageBus
}

// NewAgentRespondTool creates a new agent_respond tool.
func NewAgentRespondTool(
	registry *management.AgentRegistry,
	messageBus *management.AgentMessageBus,
) *AgentRespondTool {
	return &AgentRespondTool{
		registry:   registry,
		messageBus: messageBus,
	}
}

func (t *AgentRespondTool) Name() string { return "agent_respond" }

func (t *AgentRespondTool) Description() string {
	return "Send a response back to the original sender of a message. " +
		"The response is delivered as a 'result' message referencing the original message ID."
}

func (t *AgentRespondTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"responder_id": map[string]any{
				"type":        "string",
				"description": "The ID of the agent sending the response.",
			},
			"original_sender_id": map[string]any{
				"type":        "string",
				"description": "The ID of the agent that sent the original message.",
			},
			"original_message_id": map[string]any{
				"type":        "string",
				"description": "The ID of the message being responded to.",
			},
			"response": map[string]any{
				"type":        "string",
				"description": "The response content (plain text or JSON).",
			},
		},
		"required": []string{"responder_id", "original_sender_id", "original_message_id", "response"},
	}
}

func (t *AgentRespondTool) Execute(_ context.Context, args map[string]any) *tools.ToolResult {
	responderID, ok := args["responder_id"].(string)
	if !ok || responderID == "" {
		return tools.ErrorResult("responder_id is required")
	}

	originalSenderID, ok := args["original_sender_id"].(string)
	if !ok || originalSenderID == "" {
		return tools.ErrorResult("original_sender_id is required")
	}

	originalMsgID, ok := args["original_message_id"].(string)
	if !ok || originalMsgID == "" {
		return tools.ErrorResult("original_message_id is required")
	}

	response, ok := args["response"].(string)
	if !ok || response == "" {
		return tools.ErrorResult("response is required and must be a non-empty string")
	}

	// Validate agents
	if _, found := t.registry.GetAgent(responderID); !found {
		return tools.ErrorResult(fmt.Sprintf("responder agent not found: %q", responderID))
	}
	if _, found := t.registry.GetAgent(originalSenderID); !found {
		return tools.ErrorResult(fmt.Sprintf("original sender agent not found: %q", originalSenderID))
	}

	// Build response payload that includes original message reference
	responsePayload := map[string]string{
		"in_reply_to": originalMsgID,
		"response":    response,
	}
	rawPayload, err := json.Marshal(responsePayload)
	if err != nil {
		return tools.ErrorResult(fmt.Sprintf("failed to encode response payload: %v", err))
	}

	now := time.Now()
	reply := management.AgentMessage{
		ID:               uuid.NewString(),
		SenderID:         responderID,
		RecipientID:      originalSenderID,
		MessageType:      "result",
		Payload:          json.RawMessage(rawPayload),
		RequiresResponse: false,
		SentAt:           now,
	}

	if err := t.messageBus.Send(reply); err != nil {
		return tools.ErrorResult(fmt.Sprintf("failed to deliver response: %v", err))
	}

	// Mark original as responded
	_ = t.messageBus.MarkAsRead(originalMsgID)

	out := fmt.Sprintf(
		"↩️ Response delivered.\n- **Reply ID:** `%s`\n- **From:** `%s` → **To:** `%s`\n- **In reply to:** `%s`\n- **Sent at:** %s\n",
		reply.ID,
		responderID,
		originalSenderID,
		originalMsgID,
		now.Format(time.RFC3339),
	)

	return tools.UserResult(out)
}
