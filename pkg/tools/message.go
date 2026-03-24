package tools

import (
	"context"
	"fmt"

	"github.com/comgunner/picoclaw/pkg/bus"
)

type SendCallback func(channel, chatID, content string, media []string, buttons []bus.Button) error

type MessageTool struct {
	sendCallback   SendCallback
	defaultChannel string
	defaultChatID  string
	sentInRound    bool // Tracks whether a message was sent in the current processing round
}

func NewMessageTool() *MessageTool {
	return &MessageTool{}
}

func (t *MessageTool) Name() string {
	return "message"
}

func (t *MessageTool) Description() string {
	return "Send a message to user on a chat channel. Supports attachments (media) and interactive buttons (for approve/reject actions)."
}

func (t *MessageTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"content": map[string]any{
				"type":        "string",
				"description": "The message content to send",
			},
			"media": map[string]any{
				"type":        "array",
				"items":       map[string]any{"type": "string"},
				"description": "Optional: list of file paths or URLs to attach (images, documents, etc.)",
			},
			"buttons": map[string]any{
				"type": "array",
				"items": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"text": map[string]any{"type": "string", "description": "Button label"},
						"data": map[string]any{"type": "string", "description": "Command or data sent when clicked"},
					},
					"required": []string{"text", "data"},
				},
				"description": "Optional: interactive buttons for Telegram/Discord",
			},
			"channel": map[string]any{
				"type":        "string",
				"description": "Optional: target channel (telegram, whatsapp, etc.)",
			},
			"chat_id": map[string]any{
				"type":        "string",
				"description": "Optional: target chat/user ID",
			},
		},
		"required": []string{"content"},
	}
}

func (t *MessageTool) SetContext(channel, chatID string) {
	t.defaultChannel = channel
	t.defaultChatID = chatID
	t.ResetState() // Reset send tracking for new context
}

// ResetState clears the send tracking for a new processing round.
func (t *MessageTool) ResetState() {
	t.sentInRound = false
}

// HasSentInRound returns true if the message tool sent a message during the current round.
func (t *MessageTool) HasSentInRound() bool {
	return t.sentInRound
}

func (t *MessageTool) SetSendCallback(callback SendCallback) {
	t.sendCallback = callback
}

func (t *MessageTool) Execute(ctx context.Context, args map[string]any) *ToolResult {
	content, ok := args["content"].(string)
	if !ok {
		return &ToolResult{ForLLM: "content is required", IsError: true}
	}

	mediaRaw, _ := args["media"].([]any)
	var media []string
	for _, m := range mediaRaw {
		if s, ok := m.(string); ok {
			media = append(media, s)
		}
	}

	buttonsRaw, _ := args["buttons"].([]any)
	var buttons []bus.Button
	for _, b := range buttonsRaw {
		if bm, ok := b.(map[string]any); ok {
			text, _ := bm["text"].(string)
			data, _ := bm["data"].(string)
			if text != "" && data != "" {
				buttons = append(buttons, bus.Button{Text: text, Data: data})
			}
		}
	}

	channel, _ := args["channel"].(string)
	chatID, _ := args["chat_id"].(string)

	if channel == "" {
		channel = t.defaultChannel
	}
	if chatID == "" {
		chatID = t.defaultChatID
	}

	if channel == "" || chatID == "" {
		return &ToolResult{ForLLM: "No target channel/chat specified", IsError: true}
	}

	if t.sendCallback == nil {
		return &ToolResult{ForLLM: "Message sending not configured", IsError: true}
	}

	if err := t.sendCallback(channel, chatID, content, media, buttons); err != nil {
		return &ToolResult{
			ForLLM:  fmt.Sprintf("sending message: %v", err),
			IsError: true,
			Err:     err,
		}
	}

	t.sentInRound = true
	// Silent: user already received the message directly
	return &ToolResult{
		ForLLM: fmt.Sprintf("Message sent to %s:%s%s", channel, chatID, func() string {
			res := ""
			if len(media) > 0 {
				res += fmt.Sprintf(" with %d media files", len(media))
			}
			if len(buttons) > 0 {
				res += fmt.Sprintf(" and %d buttons", len(buttons))
			}
			return res
		}()),
		Silent: true,
	}
}
