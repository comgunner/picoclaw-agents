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
	"testing"

	"github.com/comgunner/picoclaw/pkg/providers"
)

func TestMessageSourceTracking(t *testing.T) {
	cb := NewContextBuilder(".")

	// Test history with mixed sources
	history := []providers.Message{
		{Role: "user", Content: "Hello", Source: "user"},
		{Role: "assistant", Content: "Hi there", Source: "assistant"},
		{Role: "tool", Content: "File read result", ToolCallID: "1", Source: "tool_result"},
		{Role: "assistant", Content: "Based on the file...", Source: "assistant"},
	}

	// Build messages with context
	messages := cb.BuildMessages(history, "", "How are you?", nil, "telegram", "123456789")

	// The first message should be system (without source), followed by the history and current user message
	if len(messages) < 3 { // At least system, history, and current user message
		t.Errorf("Expected at least 3 messages, got %d", len(messages))
	}

	// Check that the current user message has the correct source
	userMsgFound := false
	for _, msg := range messages {
		if msg.Role == "user" && msg.Content == "How are you?" && msg.Source == "user" {
			userMsgFound = true
			break
		}
	}

	if !userMsgFound {
		t.Errorf("Current user message was not properly tagged with source='user'")
	}

	// Check that the system message exists
	systemMsgFound := false
	for _, msg := range messages {
		if msg.Role == "system" {
			systemMsgFound = true
			break
		}
	}

	if !systemMsgFound {
		t.Error("System message was not found")
	}

	t.Logf("Created %d messages", len(messages))
}
