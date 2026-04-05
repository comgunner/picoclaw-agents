// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package utils

import "unicode/utf8"

// TokenCounter provides simple heuristics to estimate token usage of messages.
// In phase 1 we use a very rough approximation (2.5 characters per token) to
// avoid pulling in any external tokenizer. This is intentionally conservative
// and only used for pre-flight budget checks, not for accurate billing.

type TokenCounter interface {
	EstimateMessageTokens(message string) int
	CalculateContextTokens(messages []Message) int
	GetAvailableTokens(maxLimit, maxCompletion int) int
	IsWithinBudget(estimated, maxLimit, reserved int) bool
}

// Message is a lightweight representation used by the counter. We duplicate a
// small subset of providers.Message here to avoid import cycles.
type Message struct {
	Role    string
	Content string
}

// basicCounter is the default implementation of TokenCounter.
// It uses the same heuristic that AgentLoop.estimateTokens employs.

var _ TokenCounter = (*basicCounter)(nil)

type basicCounter struct{}

func NewBasicTokenCounter() TokenCounter {
	return &basicCounter{}
}

func (b *basicCounter) EstimateMessageTokens(message string) int {
	// 2.5 characters per token, rounded down
	total := utf8.RuneCountInString(message)
	return total * 2 / 5
}

func (b *basicCounter) CalculateContextTokens(messages []Message) int {
	total := 0
	for _, m := range messages {
		total += b.EstimateMessageTokens(m.Content)
	}
	return total
}

func (b *basicCounter) GetAvailableTokens(maxLimit, maxCompletion int) int {
	return maxLimit - maxCompletion
}

func (b *basicCounter) IsWithinBudget(estimated, maxLimit, reserved int) bool {
	return estimated <= (maxLimit - reserved)
}
