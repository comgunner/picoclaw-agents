// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package utils_test

import (
	"testing"

	"github.com/comgunner/picoclaw/pkg/utils"
)

func TestBasicCounter_EstimateMessageTokens(t *testing.T) {
	c := utils.NewBasicTokenCounter()
	// 10 characters should translate to 4 tokens (10*2/5)
	if got := c.EstimateMessageTokens("abcdefghij"); got != 4 {
		t.Errorf("expected 4 tokens, got %d", got)
	}

	// empty string
	if got := c.EstimateMessageTokens(""); got != 0 {
		t.Errorf("expected 0 tokens for empty string, got %d", got)
	}
}

func TestBasicCounter_CalculateContextTokens(t *testing.T) {
	c := utils.NewBasicTokenCounter()
	msgs := []utils.Message{
		{Role: "user", Content: "hello"},      // 2 tokens
		{Role: "assistant", Content: "world"}, // 2 tokens
	}
	if got := c.CalculateContextTokens(msgs); got != 4 {
		t.Errorf("expected 4 tokens total, got %d", got)
	}
}

func TestBasicCounter_Budget(t *testing.T) {
	c := utils.NewBasicTokenCounter()
	est := c.EstimateMessageTokens("abcdefghij")
	if !c.IsWithinBudget(est, 10, 0) {
		t.Errorf("token should be within budget")
	}
	// use a smaller limit to force failure
	if c.IsWithinBudget(est, 3, 0) {
		t.Errorf("token should exceed budget when limit is 3")
	}
}
