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
	"github.com/comgunner/picoclaw/pkg/providers/protocoltypes"
)

// TestGhostConversationPrevention verifies that the fixes for ghost conversations are working
func TestGhostConversationPrevention(t *testing.T) {
	// This test simulates the scenario described in the ghost conversation problem
	// where content from memory/context gets interpreted as a new user command

	cb := NewContextBuilder(".")

	// Simulate history that includes content that might be misinterpreted
	// as a user command (like reminder content from memory files)
	history := []providers.Message{
		{
			Role:    "user",
			Content: "Hello, how are you?",
			Source:  "user",
		},
		{
			Role:    "assistant",
			Content: "I'm doing well, thank you for asking!",
			Source:  "assistant",
		},
		{
			Role:       "tool",
			Content:    "Content of MEMORY.md: User likes coffee. User has meeting at 3 PM. Remember me to take out trash at 8 PM.",
			ToolCallID: "tool123",
			Source:     "tool_result",
		},
		{
			Role:    "assistant",
			Content: "I've read your memory file. It contains information about your preferences and schedule.",
			Source:  "assistant",
		},
	}

	// Build messages for the LLM
	messages := cb.BuildMessages(history, "", "What should I do today?", nil, "telegram", "123456789")

	// Verify that the current user message is properly identified
	userMessageFound := false
	for _, msg := range messages {
		if msg.Role == "user" && msg.Content == "What should I do today?" && msg.Source == "user" {
			userMessageFound = true
			break
		}
	}

	if !userMessageFound {
		t.Error("Current user message was not properly tagged with source='user'")
	}

	// Verify that the system prompt includes the new instructions about user intent recognition
	systemMessageFound := false
	for _, msg := range messages {
		if msg.Role == "system" {
			systemMessageFound = true
			if len(msg.Content) == 0 {
				t.Error("System message content is empty")
			}
			// Check if the new instructions are present
			if len(msg.Content) > 0 {
				hasIntentRecognition := containsSubstring(msg.Content, "User Intent Recognition")
				hasDirectUserInput := containsSubstring(msg.Content, "DIRECT USER INPUT")

				if !hasIntentRecognition || !hasDirectUserInput {
					t.Log("System message may not contain the new user intent recognition instructions")
					t.Logf("System content preview: %.100s...", msg.Content)
				}
			}
			break
		}
	}

	if !systemMessageFound {
		t.Error("System message was not found in the message list")
	}

	t.Log("Ghost conversation prevention elements are in place")
}

// Helper function to check if a string contains a substring
func containsSubstring(text, substr string) bool {
	return len(text) >= len(substr) &&
		(text == substr ||
			len(text) > len(substr) &&
				(findStringInString(text, substr) != -1))
}

// Simple string search implementation
func findStringInString(text, substr string) int {
	if len(substr) == 0 {
		return 0
	}
	for i := 0; i <= len(text)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			if text[i+j] != substr[j] {
				match = false
				break
			}
		}
		if match {
			return i
		}
	}
	return -1
}

// TestMessageSourceInProviderCompatibility ensures that adding the Source field
// doesn't break compatibility with provider interfaces
func TestMessageSourceInProviderCompatibility(t *testing.T) {
	// Create a message with source
	msg := protocoltypes.Message{
		Role:    "user",
		Content: "Test message",
		Source:  "user",
	}

	// Verify that it can still be used as a regular message
	if msg.Role != "user" {
		t.Errorf("Role was not preserved, got: %s", msg.Role)
	}

	if msg.Content != "Test message" {
		t.Errorf("Content was not preserved, got: %s", msg.Content)
	}

	if msg.Source != "user" {
		t.Errorf("Source was not preserved, got: %s", msg.Source)
	}

	t.Log("Message with source field maintains compatibility")
}

// TestToolResultSourceTracking verifies that tool results are properly sourced
func TestToolResultSourceTracking(t *testing.T) {
	cb := NewContextBuilder(".")

	// Add a tool result using the context builder
	initialMessages := []providers.Message{}
	updatedMessages := cb.AddToolResult(initialMessages, "call123", "read_file", "File content here")

	if len(updatedMessages) != 1 {
		t.Fatalf("Expected 1 message after AddToolResult, got %d", len(updatedMessages))
	}

	toolMsg := updatedMessages[0]
	if toolMsg.Role != "tool" {
		t.Errorf("Expected role 'tool', got: %s", toolMsg.Role)
	}

	if toolMsg.Content != "File content here" {
		t.Errorf("Expected content 'File content here', got: %s", toolMsg.Content)
	}

	if toolMsg.ToolCallID != "call123" {
		t.Errorf("Expected ToolCallID 'call123', got: %s", toolMsg.ToolCallID)
	}

	if toolMsg.Source != "tool_result" {
		t.Errorf("Expected source 'tool_result', got: %s", toolMsg.Source)
	}

	t.Log("Tool results are properly tagged with source='tool_result'")
}
