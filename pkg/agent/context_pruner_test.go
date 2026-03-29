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
	"testing"

	"github.com/comgunner/picoclaw/pkg/config"
	"github.com/comgunner/picoclaw/pkg/providers"
)

func TestPruneMessages_Disabled(t *testing.T) {
	cfg := config.ContextPruningConfig{
		Enabled:            false,
		MaxToolResultChars: 1000,
	}

	messages := []providers.Message{
		{Role: "system", Content: "You are helpful"},
		{Role: "user", Content: "Hello"},
		{Role: "assistant", Content: "Hi there!"},
	}

	result := PruneMessages(messages, cfg)
	if len(result) != len(messages) {
		t.Fatalf("expected %d messages, got %d", len(messages), len(result))
	}

	// Verify content unchanged
	for i := range messages {
		if result[i].Content != messages[i].Content {
			t.Errorf("message %d content changed", i)
		}
	}
}

func TestPruneMessages_ShortContent_Unchanged(t *testing.T) {
	cfg := config.ContextPruningConfig{
		Enabled:            true,
		MaxToolResultChars: 1000,
	}

	messages := []providers.Message{
		{Role: "system", Content: "You are helpful"},
		{Role: "user", Content: "Hello"},
		{
			Role:      "assistant",
			Content:   "Hi there!",
			ToolCalls: []providers.ToolCall{{ID: "call1", Function: &providers.FunctionCall{Name: "shell"}}},
		},
		{Role: "tool", Content: "Short output", ToolCallID: "call1"},
	}

	result := PruneMessages(messages, cfg)
	if len(result) != 4 {
		t.Fatalf("expected 4 messages, got %d", len(result))
	}

	// Tool result should be unchanged
	if result[3].Content != "Short output" {
		t.Errorf("expected 'Short output', got %q", result[3].Content)
	}
}

func TestPruneMessages_LongToolResult_Truncated(t *testing.T) {
	cfg := config.ContextPruningConfig{
		Enabled:            true,
		MaxToolResultChars: 100,
	}

	// Create a long tool result (500 chars)
	longContent := strings.Repeat("x", 500)

	messages := []providers.Message{
		{Role: "system", Content: "You are helpful"},
		{Role: "user", Content: "Run command"},
		{
			Role:      "assistant",
			Content:   "Running...",
			ToolCalls: []providers.ToolCall{{ID: "call1", Function: &providers.FunctionCall{Name: "shell"}}},
		},
		{Role: "tool", Content: longContent, ToolCallID: "call1"},
	}

	result := PruneMessages(messages, cfg)

	// Should be truncated
	if !strings.Contains(result[3].Content, "[truncated") {
		t.Errorf("expected truncation marker, got %q", result[3].Content)
	}

	// Should preserve tail
	if !strings.HasSuffix(result[3].Content, strings.Repeat("x", 200)) {
		t.Error("expected tail preservation (last 200 chars)")
	}

	// Should mention tool name
	if !strings.Contains(result[3].Content, "tool: shell") {
		t.Error("expected tool name in truncation marker")
	}
}

func TestPruneMessages_ExcludedTool_Unchanged(t *testing.T) {
	cfg := config.ContextPruningConfig{
		Enabled:            true,
		MaxToolResultChars: 100,
		ExcludeTools:       []string{"memory_store"},
	}

	// Long content from excluded tool
	longContent := strings.Repeat("y", 500)

	messages := []providers.Message{
		{Role: "system", Content: "You are helpful"},
		{Role: "user", Content: "Store memory"},
		{
			Role:      "assistant",
			Content:   "Storing...",
			ToolCalls: []providers.ToolCall{{ID: "call1", Function: &providers.FunctionCall{Name: "memory_store"}}},
		},
		{Role: "tool", Content: longContent, ToolCallID: "call1"},
	}

	result := PruneMessages(messages, cfg)

	// Should NOT be truncated (excluded tool)
	if result[3].Content != longContent {
		t.Errorf("expected excluded tool to be unchanged, got truncated content")
	}
}

func TestPruneMessages_AggressiveTool_HalfLimit(t *testing.T) {
	cfg := config.ContextPruningConfig{
		Enabled:            true,
		MaxToolResultChars: 200,
		AggressiveTools:    []string{"shell"},
	}

	// Content that would be ok for normal tools (150 chars) but too long for aggressive (100 chars)
	content := strings.Repeat("z", 150)

	messages := []providers.Message{
		{Role: "system", Content: "You are helpful"},
		{Role: "user", Content: "Run command"},
		{
			Role:      "assistant",
			Content:   "Running...",
			ToolCalls: []providers.ToolCall{{ID: "call1", Function: &providers.FunctionCall{Name: "shell"}}},
		},
		{Role: "tool", Content: content, ToolCallID: "call1"},
	}

	result := PruneMessages(messages, cfg)

	// Should be truncated (aggressive tool, limit = 200/2 = 100)
	if !strings.Contains(result[3].Content, "[truncated") {
		t.Errorf("expected truncation for aggressive tool, got %q", result[3].Content)
	}
}

func TestPruneMessages_PreservesSystemPrompt(t *testing.T) {
	cfg := config.ContextPruningConfig{
		Enabled:            true,
		MaxToolResultChars: 50,
	}

	messages := []providers.Message{
		{
			Role:    "system",
			Content: "You are a helpful assistant with a very long system prompt " + strings.Repeat("x", 500),
		},
		{Role: "user", Content: "Hello"},
		{
			Role:      "assistant",
			Content:   "Hi!",
			ToolCalls: []providers.ToolCall{{ID: "call1", Function: &providers.FunctionCall{Name: "shell"}}},
		},
		{Role: "tool", Content: strings.Repeat("y", 200), ToolCallID: "call1"},
	}

	result := PruneMessages(messages, cfg)

	// System prompt should be unchanged
	if result[0].Content != messages[0].Content {
		t.Error("system prompt should not be modified")
	}

	// User message should be unchanged
	if result[1].Content != messages[1].Content {
		t.Error("user message should not be modified")
	}
}

func TestPruneMessages_PreservesUserMessages(t *testing.T) {
	cfg := config.ContextPruningConfig{
		Enabled:            true,
		MaxToolResultChars: 50,
	}

	longUserMsg := strings.Repeat("u", 300)

	messages := []providers.Message{
		{Role: "system", Content: "You are helpful"},
		{Role: "user", Content: longUserMsg},
		{Role: "assistant", Content: "OK"},
	}

	result := PruneMessages(messages, cfg)

	// User message should be unchanged (only tool results are pruned)
	if result[1].Content != longUserMsg {
		t.Error("user message should not be pruned")
	}
}

func TestExtractToolName_FindsCorrectTool(t *testing.T) {
	messages := []providers.Message{
		{Role: "system", Content: "You are helpful"},
		{Role: "user", Content: "Hello"},
		{Role: "assistant", Content: "Running...", ToolCalls: []providers.ToolCall{
			{ID: "call1", Function: &providers.FunctionCall{Name: "shell"}},
			{ID: "call2", Function: &providers.FunctionCall{Name: "web_search"}},
		}},
		{Role: "tool", Content: "Output 1", ToolCallID: "call1"},
		{Role: "tool", Content: "Output 2", ToolCallID: "call2"},
	}

	toolName1 := extractToolName(messages, 3)
	if toolName1 != "shell" {
		t.Errorf("expected 'shell', got %q", toolName1)
	}

	toolName2 := extractToolName(messages, 4)
	if toolName2 != "web_search" {
		t.Errorf("expected 'web_search', got %q", toolName2)
	}
}

func TestExtractToolName_OrphanedTool(t *testing.T) {
	messages := []providers.Message{
		{Role: "system", Content: "You are helpful"},
		{Role: "user", Content: "Hello"},
		{Role: "tool", Content: "Orphaned", ToolCallID: "nonexistent"},
	}

	toolName := extractToolName(messages, 2)
	if toolName != "" {
		t.Errorf("expected empty string for orphaned tool, got %q", toolName)
	}
}
