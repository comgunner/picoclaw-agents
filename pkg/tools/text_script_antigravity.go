// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
// Copyright (c) 2026 PicoClaw contributors

package tools

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/comgunner/picoclaw/pkg/logger"
	"github.com/comgunner/picoclaw/pkg/providers"
	"github.com/comgunner/picoclaw/pkg/utils"
)

// Retry delays for 429 rate limiting on text generation.
// Total max wait: 5s + 15s + 30s + 60s + 120s = 230s (~4 min)
var textScriptRetryDelays = []time.Duration{
	5 * time.Second,
	15 * time.Second,
	30 * time.Second,
	60 * time.Second,
	120 * time.Second,
}

// GenerateTextScriptAntigravity generates a text script using Antigravity OAuth.
// This is the PRIMARY method — no API key needed, just OAuth credentials.
// Includes automatic retry with backoff on 429 rate limit.
func GenerateTextScriptAntigravity(
	ctx context.Context,
	model string,
	req utils.TextScriptRequest,
) (*utils.TextScriptResult, error) {
	if model == "" {
		model = "gemini-3-flash"
	}
	if req.Language == "" {
		req.Language = utils.DetectLanguage(req.Topic)
	}

	template, err := utils.LoadPromptTemplate(req.TemplatePath, "script")
	if err != nil {
		return nil, fmt.Errorf("error loading template: %v", err)
	}

	prompt := strings.ReplaceAll(template, "{topic}", req.Topic)
	prompt = strings.ReplaceAll(prompt, "{category}", req.Category)
	prompt = strings.ReplaceAll(prompt, "{duration}", req.Duration)
	prompt = strings.ReplaceAll(prompt, "{tone}", req.Tone)

	provider := providers.NewAntigravityProvider()

	messages := []providers.Message{
		{
			Role:    "system",
			Content: "You are a professional social media copywriter. Write engaging, well-structured content.",
		},
		{Role: "user", Content: prompt},
	}

	var resp *providers.LLMResponse

	// Retry loop with exponential backoff for 429 rate limits
	for i, delay := range textScriptRetryDelays {
		resp, err = provider.Chat(ctx, messages, nil, model, map[string]any{
			"max_tokens":  2048,
			"temperature": 0.7,
		})
		if err == nil {
			break // Success
		}

		// Check if it's a rate limit error
		if isTextScriptRateLimit(err) {
			if i < len(textScriptRetryDelays)-1 {
				logger.WarnCF("tools.text_script", "Rate limited, retrying after delay", map[string]any{
					"attempt": i + 1,
					"delay":   delay.String(),
					"error":   err.Error(),
				})
				select {
				case <-time.After(delay):
					// Retry next iteration
				case <-ctx.Done():
					return nil, fmt.Errorf("context canceled while waiting for retry: %w", ctx.Err())
				}
				continue
			}
			// All retries exhausted
			return nil, fmt.Errorf(
				"antigravity rate limit exceeded after %d retries (~4 min total). Wait ~5 min and try again: %w",
				len(textScriptRetryDelays),
				err,
			)
		}

		// Non-rate-limit error — fail immediately
		return nil, fmt.Errorf("antigravity text generation failed: %w", err)
	}

	if resp == nil {
		return nil, fmt.Errorf("antigravity text generation returned nil response")
	}

	wordCount := len(strings.Fields(resp.Content))
	estimatedDuration := estimateDuration(wordCount)

	return &utils.TextScriptResult{
		Script:            resp.Content,
		WordCount:         wordCount,
		EstimatedDuration: estimatedDuration,
		Language:          req.Language,
	}, nil
}

// isTextScriptRateLimit checks if the error is a 429 rate limit.
func isTextScriptRateLimit(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "429") ||
		strings.Contains(msg, "RESOURCE_EXHAUSTED") ||
		strings.Contains(msg, "rate limit")
}

func estimateDuration(wordCount int) string {
	seconds := (wordCount * 60) / 150
	if seconds < 30 {
		return "30s"
	}
	if seconds < 60 {
		return "60s"
	}
	minutes := seconds / 60
	return fmt.Sprintf("%dmin", minutes)
}
