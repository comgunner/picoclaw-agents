// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package agent

import (
	"strings"

	"github.com/comgunner/picoclaw/pkg/providers"
)

type MessageFilter struct {
	AckWords []string
}

func NewMessageFilter() *MessageFilter {
	return &MessageFilter{
		AckWords: []string{"ok", "entendido", "understood", "got it", "sure", "will do", "done", "yes"},
	}
}

func (f *MessageFilter) Apply(messages []providers.Message) []providers.Message {
	var filtered []providers.Message

	for _, msg := range messages {
		// 4. Never remove system or tool messages
		if msg.Role == "system" || msg.Role == "tool" {
			// But DO truncate them if they are too long (except system)
			if msg.Role == "tool" && len(msg.Content) > 1250 {
				msg.Content = msg.Content[:250] + "... [Extensive log/output truncated by filter] ..." + msg.Content[len(msg.Content)-250:]
			}
			filtered = append(filtered, msg)
			continue
		}

		// 1. Remove ACKs
		if f.isAck(msg.Content) {
			continue
		}

		// 2. Extensivo logs > 500 tokens (approx 1250 chars)
		if len(msg.Content) > 1250 && msg.Role != "user" {
			msg.Content = msg.Content[:250] + "... [Extensive log/output truncated by filter] ..." + msg.Content[len(msg.Content)-250:]
		}

		// 3. Keep other messages
		filtered = append(filtered, msg)
	}

	return filtered
}

func (f *MessageFilter) isAck(content string) bool {
	c := strings.TrimSpace(strings.ToLower(content))
	for _, w := range f.AckWords {
		if c == w || c == w+"." || c == w+"!" {
			return true
		}
	}
	// Check for short acknowledgments
	if len(c) < 15 &&
		(strings.Contains(c, "ok") || strings.Contains(c, "understood") || strings.Contains(c, "entendido")) {
		return true
	}
	return false
}
