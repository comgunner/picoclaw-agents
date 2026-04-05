// PicoClaw-Agents - Token estimation utilities
// Ported from picoclaw_original/pkg/tokenizer/estimator.go
// Adapted for fork's provider types (no Media field).

// ============================================================================
// ⚠️  CRITICAL: Token Estimation — DO NOT MODIFY WITHOUT REVIEW
// ============================================================================
//
// These functions are called by pkg/agent/loop.go to prevent OpenRouter Free
// tier 402 errors ("Prompt tokens limit exceeded"). The original bug was caused
// by using simple character counting instead of these estimators, which
// underestimated token usage by ~18,000 tokens.
//
// IMMUTABLE RULES:
//   - EstimateMessageTokens MUST count: Content, ReasoningContent, ToolCalls,
//     ToolCallID, and SystemParts (the larger of Content vs SystemParts)
//   - EstimateToolDefsTokens MUST count: function name + description + params JSON
//   - Heuristic: 2.5 chars/token + 12 overhead per message + 20 per tool
//
// DO NOT remove or simplify these functions. They are the ONLY accurate way
// to estimate tokens before sending to the LLM.
//
// See: local_work/MEMORY.md, local_work/openrouter_free_token_fix.md
// ============================================================================

package tokenizer

import (
	"encoding/json"
	"unicode/utf8"

	"github.com/comgunner/picoclaw/pkg/providers"
)

// EstimateMessageTokens estimates the token count for a single message,
// including Content, ReasoningContent, ToolCalls arguments, ToolCallID
// metadata, and SystemParts. Uses a heuristic of 2.5 characters per token.
func EstimateMessageTokens(msg providers.Message) int {
	contentChars := utf8.RuneCountInString(msg.Content)

	// SystemParts are structured system blocks used for cache-aware adapters.
	// They carry the same content as Content, but in multiple blocks.
	// We estimate them as an alternative representation, not additive.
	systemPartsChars := 0
	if len(msg.SystemParts) > 0 {
		for _, part := range msg.SystemParts {
			systemPartsChars += utf8.RuneCountInString(part.Text)
		}
		// Per-part overhead for JSON structure (type, text, cache_control).
		const perPartOverhead = 20
		systemPartsChars += len(msg.SystemParts) * perPartOverhead
	}

	// Use the larger of the two representations to stay conservative.
	chars := contentChars
	if systemPartsChars > chars {
		chars = systemPartsChars
	}

	chars += utf8.RuneCountInString(msg.ReasoningContent)

	for _, tc := range msg.ToolCalls {
		chars += len(tc.ID) + len(tc.Type)
		if tc.Function != nil {
			// Count function name + arguments (the wire format for most providers).
			// tc.Name mirrors tc.Function.Name — count only once to avoid double-counting.
			chars += len(tc.Function.Name) + len(tc.Function.Arguments)
		} else {
			// Fallback: some provider formats use top-level Name without Function.
			chars += len(tc.Name)
		}
	}

	if msg.ToolCallID != "" {
		chars += len(msg.ToolCallID)
	}

	// Per-message overhead for role label, JSON structure, separators.
	const messageOverhead = 12
	chars += messageOverhead

	tokens := chars * 2 / 5

	// NOTE: The original also counts Media items (256 tokens each).
	// The fork's Message type does not have a Media field, so this is omitted.

	return tokens
}

// EstimateMessagesTokens estimates total tokens for a slice of messages.
func EstimateMessagesTokens(messages []providers.Message) int {
	total := 0
	for _, msg := range messages {
		total += EstimateMessageTokens(msg)
	}
	return total
}

// EstimateToolDefsTokens estimates the total token cost of tool definitions
// as they appear in the LLM request.
//
// ⚠️  CRITICAL: This is the ONLY accurate way to estimate tool definition tokens.
// 60+ tools consume ~15,000 tokens, NOT ~2,500. Using a fixed value causes
// OpenRouter 402 errors. ALWAYS call this function instead of a hardcoded overhead.
func EstimateToolDefsTokens(defs []providers.ToolDefinition) int {
	if len(defs) == 0 {
		return 0
	}

	totalChars := 0
	for _, d := range defs {
		totalChars += len(d.Function.Name) + len(d.Function.Description)

		if d.Function.Parameters != nil {
			if paramJSON, err := json.Marshal(d.Function.Parameters); err == nil {
				totalChars += len(paramJSON)
			}
		}

		// Per-tool overhead: type field, JSON structure, separators.
		totalChars += 20
	}

	return totalChars * 2 / 5
}
