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

func TestDefaultMessageRanker_ScoreMessage(t *testing.T) {
	ranker := NewDefaultMessageRanker()

	userMsg := providers.Message{Role: "user", Content: "Hello"}
	if score := ranker.ScoreMessage(userMsg); score < 1.0 {
		t.Errorf("expected high score for user message, got %f", score)
	}

	toolMsg := providers.Message{Role: "tool", Content: "Result"}
	if score := ranker.ScoreMessage(toolMsg); score < 0.8 {
		t.Errorf("expected medium score for tool message, got %f", score)
	}

	ackMsg := providers.Message{Role: "assistant", Content: "  Ok  "}
	if score := ranker.ScoreMessage(ackMsg); score > 0.3 {
		t.Errorf("expected low score for ACK message, got %f", score)
	}
}

func TestDefaultMessageRanker_FilterByImportance(t *testing.T) {
	ranker := NewDefaultMessageRanker()

	messages := []providers.Message{
		{Role: "system", Content: "You are an assistant."},
		{Role: "user", Content: "Do something"},
		{Role: "assistant", Content: "ok"}, // Should be dropped
		{Role: "tool", Content: "Data"},
	}

	filtered := ranker.FilterByImportance(messages, 0.3)

	if len(filtered) != 3 {
		t.Errorf("expected 3 messages after filtering, got %d", len(filtered))
	}

	for _, m := range filtered {
		if m.Content == "ok" {
			t.Errorf("ack message should have been filtered out")
		}
	}
}
