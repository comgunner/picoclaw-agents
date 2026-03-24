// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package providers_test

import (
	"testing"

	"github.com/comgunner/picoclaw/pkg/providers"
	"github.com/comgunner/picoclaw/pkg/utils"
)

func TestTokenValidator_Validate(t *testing.T) {
	counter := utils.NewBasicTokenCounter()
	// limit set low to trigger failure easily
	validator := providers.NewTokenValidator(counter, 10)

	// create a short message that should fit when limit is large
	msgs := []providers.Message{{Role: "user", Content: "abcdefghij"}} // ~4 tokens

	if err := validator.Validate(msgs, 2); err != nil {
		t.Errorf("expected validation to succeed but got %v", err)
	}

	// adjust maxContext to a small value to force failure
	validator = providers.NewTokenValidator(counter, 5)
	// use a longer message to exceed the budget
	msgs = []providers.Message{{Role: "user", Content: "abcdefghijabcdefghij"}} // ~8 tokens
	if err := validator.Validate(msgs, 2); err == nil {
		t.Errorf("expected validation to fail when context+completion > limit")
	}
}
