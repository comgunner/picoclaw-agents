// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package providers

import (
	"fmt"

	"github.com/comgunner/picoclaw/pkg/utils"
)

// TokenValidator knows how to check that a set of messages fits within the
// contextual budget of the chosen model.  The implementation is deliberately
// simple in phase 1 – it uses the provided utils.TokenCounter and a fixed
// maximum context size (e.g. 131072 tokens) to make a conservative decision.

// TokenValidator is used by the agent loop before dispatching a request to an
// LLM provider.  If validation fails the call is aborted and an informative
// error bubbles up so the caller can suggest remedial actions to the user.

type TokenValidator interface {
	Validate(messages []Message, maxCompletion int) error
}

// newValidator returns a basic TokenValidator using the supplied counter and
// context limit.  maxContextTokens should reflect the largest window supported
// by any model the system may use; for most deployments the default constant
// 131072 is adequate (see pkg/constants/max_context.go).

func NewTokenValidator(counter utils.TokenCounter, maxContextTokens int) TokenValidator {
	return &tokenValidator{
		counter:          counter,
		maxContextTokens: maxContextTokens,
	}
}

// tokenValidator is a trivial implementation of TokenValidator.

type tokenValidator struct {
	counter          utils.TokenCounter
	maxContextTokens int
}

func (v *tokenValidator) Validate(messages []Message, maxCompletion int) error {
	// convert providers.Message to utils.Message
	utilMsgs := make([]utils.Message, 0, len(messages))
	for _, m := range messages {
		utilMsgs = append(utilMsgs, utils.Message{Role: m.Role, Content: m.Content})
	}

	used := v.counter.CalculateContextTokens(utilMsgs)
	if used+maxCompletion > v.maxContextTokens {
		return fmt.Errorf(
			"estimated tokens (%d context + %d completion) exceed limit %d",
			used,
			maxCompletion,
			v.maxContextTokens,
		)
	}
	return nil
}
