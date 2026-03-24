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

func TestMessageFilter_Apply(t *testing.T) {
	filter := NewMessageFilter()

	messages := []providers.Message{
		{Role: "system", Content: "You are a helpful assistant."},
		{Role: "user", Content: "Hello!"},
		{Role: "assistant", Content: "ok"},
		{Role: "user", Content: "Please analyze this dataset."},
		{Role: "assistant", Content: "understood"},
	}

	filtered := filter.Apply(messages)

	if len(filtered) != 3 {
		t.Errorf("expected 3 messages, got %d", len(filtered))
	}

	if filtered[0].Role != "system" {
		t.Errorf("expected first message to be system, got %s", filtered[0].Role)
	}

	if filtered[1].Content != "Hello!" {
		t.Errorf("expected greeting to be preserved")
	}
}

func TestMessageFilter_ExtensiveLogs(t *testing.T) {
	filter := NewMessageFilter()

	longLog := make([]byte, 1300)
	for i := range longLog {
		longLog[i] = 'a'
	}

	messages := []providers.Message{
		{Role: "system", Content: "System"},
		{Role: "tool", Content: string(longLog)},
		{Role: "user", Content: "Explain"},
	}

	filtered := filter.Apply(messages)

	if len(filtered) != 3 {
		t.Errorf("expected 3 messages, got %d", len(filtered))
	}

	if len(filtered[1].Content) >= 1300 {
		t.Errorf("expected long tool output to be truncated")
	}
}
