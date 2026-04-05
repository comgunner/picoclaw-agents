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

type RankedMessage struct {
	Message providers.Message
	Score   float32
}

type RankedMessages []RankedMessage

type MessageRanker interface {
	RankMessages(messages []providers.Message) RankedMessages
	ScoreMessage(msg providers.Message) float32
	FilterByImportance(messages []providers.Message, threshold float32) []providers.Message
}

type DefaultMessageRanker struct{}

func NewDefaultMessageRanker() *DefaultMessageRanker {
	return &DefaultMessageRanker{}
}

func (r *DefaultMessageRanker) RankMessages(messages []providers.Message) RankedMessages {
	var ranked RankedMessages
	for _, m := range messages {
		ranked = append(ranked, RankedMessage{
			Message: m,
			Score:   r.ScoreMessage(m),
		})
	}
	return ranked
}

func (r *DefaultMessageRanker) ScoreMessage(msg providers.Message) float32 {
	// Instructions from the user (high priority)
	if msg.Role == "user" || msg.Role == "system" {
		return 1.0
	}

	// Tool calls and search results (medium)
	if msg.Role == "tool" || len(msg.ToolCalls) > 0 {
		return 0.8
	}

	// Acknowledgments or ACKs (low)
	lowerContent := strings.ToLower(strings.TrimSpace(msg.Content))
	if len(lowerContent) < 15 &&
		(strings.Contains(lowerContent, "ok") || strings.Contains(lowerContent, "understood")) {
		return 0.2
	}

	return 0.5
}

func (r *DefaultMessageRanker) FilterByImportance(messages []providers.Message, threshold float32) []providers.Message {
	ranked := r.RankMessages(messages)
	var filtered []providers.Message

	for _, rm := range ranked {
		if rm.Score >= threshold {
			filtered = append(filtered, rm.Message)
		}
	}

	return filtered
}
