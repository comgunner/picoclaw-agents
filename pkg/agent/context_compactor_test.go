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
	"context"
	"testing"

	"github.com/comgunner/picoclaw/pkg/providers"
)

// Minimal mock provider just to implement the interface needed for testing
type mockCompactorProvider struct {
	providers.LLMProvider
}

func (m *mockCompactorProvider) Chat(
	ctx context.Context,
	msgs []providers.Message,
	tools []providers.ToolDefinition,
	model string,
	options map[string]any,
) (*providers.LLMResponse, error) {
	return &providers.LLMResponse{
		Content: "Mocked Summary",
	}, nil
}

func TestContextCompactor_ShouldCompact(t *testing.T) {
	compactor := NewDefaultContextCompactor(nil)
	if !compactor.ShouldCompact(800, 1000, 0.75) {
		t.Error("expected true when tokens exceed threshold")
	}
	if compactor.ShouldCompact(500, 1000, 0.75) {
		t.Error("expected false when tokens are below threshold")
	}
}

func TestContextCompactor_CompactMessages(t *testing.T) {
	compactor := NewDefaultContextCompactor(nil)
	provider := &mockCompactorProvider{}

	messages := []providers.Message{
		{Role: "system", Content: "System prompt"},
		{Role: "user", Content: "Q1"},
		{Role: "assistant", Content: "A1"},
		{Role: "user", Content: "Q2"},
		{Role: "assistant", Content: "A2"},
		{Role: "user", Content: "Q3"},
		{Role: "assistant", Content: "A3"},
		{Role: "user", Content: "Q4"},
	}

	compacted, err := compactor.CompactMessages(context.Background(), provider, "mock-model", messages, "session-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(compacted) >= len(messages) {
		t.Fatalf("expected compacted length to be smaller than original length")
	}

	if compacted[0].Role != "system" {
		t.Errorf("expected first message to be system prompt")
	}

	if compacted[len(compacted)-1].Content != "Q4" {
		t.Errorf("expected last message to be preserved")
	}
}
