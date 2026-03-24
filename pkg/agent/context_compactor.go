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
	"fmt"
	"strings"

	"github.com/comgunner/picoclaw/pkg/health"
	"github.com/comgunner/picoclaw/pkg/logger"
	"github.com/comgunner/picoclaw/pkg/providers"
	"github.com/comgunner/picoclaw/pkg/utils"
)

type ContextCompactor interface {
	ShouldCompact(currentTokens, maxTokens int, threshold float64) bool
	CompactMessages(
		ctx context.Context,
		provider providers.LLMProvider,
		model string,
		messages []providers.Message,
		sessionKey string,
	) ([]providers.Message, error)
	ExtractKeyContext(old, new []providers.Message) string
}

type DefaultContextCompactor struct {
	Filter *MessageFilter
	Ranker MessageRanker
	Cache  *utils.SummaryCache
}

func NewDefaultContextCompactor(cache *utils.SummaryCache) *DefaultContextCompactor {
	return &DefaultContextCompactor{
		Filter: NewMessageFilter(),
		Ranker: NewDefaultMessageRanker(),
		Cache:  cache,
	}
}

func (c *DefaultContextCompactor) ShouldCompact(currentTokens, maxTokens int, threshold float64) bool {
	if maxTokens == 0 {
		return false
	}
	ratio := float64(currentTokens) / float64(maxTokens)
	return ratio >= threshold
}

func (c *DefaultContextCompactor) CompactMessages(
	ctx context.Context,
	provider providers.LLMProvider,
	model string,
	messages []providers.Message,
	sessionKey string,
) ([]providers.Message, error) {
	if len(messages) <= 6 {
		return messages, nil
	}

	// Metrics start
	health.ContextMetrics.RecordCompaction()

	// Save system prompt and last user message
	sysPrompt := messages[0]
	lastMsg := messages[len(messages)-1]

	// 1. Identify "Atomic" Tool Call Groups
	// We must ensure that assistant (with tool_calls) and its corresponding tool responses are NOT separated.
	messagesToCompact := messages[1 : len(messages)-1]

	// Rule: If the last message in the 'olderHalf' is an assistant with tool_calls,
	// or if the first message in 'newerHalf' is a tool response, we must adjust the split point.
	mid := len(messagesToCompact) / 2

	// Adjust mid point to not break tool call sequences at the split boundary
	for mid > 0 && mid < len(messagesToCompact) {
		msg := messagesToCompact[mid-1]
		// If mid-1 is an assistant with tool calls, the next message (mid) SHOULD be a tool response.
		// We shouldn't split here. Move mid forward until we find a non-tool message or the end.
		if msg.Role == "assistant" && len(msg.ToolCalls) > 0 {
			mid++
			continue
		}
		// If mid is a tool response, it must stay with its preceding assistant message.
		if messagesToCompact[mid].Role == "tool" {
			mid++
			continue
		}
		break
	}

	olderHalf := messagesToCompact[:mid]
	newerHalf := messagesToCompact[mid:]

	// Defensive: ensure newerHalf does not start with an orphaned tool response.
	// If it does, walk mid backward to include the parent assistant message (with its tool calls).
	if len(newerHalf) > 0 && newerHalf[0].Role == "tool" {
		for mid > 0 {
			mid--
			if messagesToCompact[mid].Role == "assistant" && len(messagesToCompact[mid].ToolCalls) > 0 {
				break
			}
		}
		olderHalf = messagesToCompact[:mid]
		newerHalf = messagesToCompact[mid:]
	}

	// Apply filters to olderHalf before summarizing
	filteredOlder := c.Filter.Apply(olderHalf)
	filteredOlder = c.Ranker.FilterByImportance(filteredOlder, 0.2)

	topic := "conversation_segment"
	var summary string

	if summary == "" {
		s, err := c.GenerateSummary(ctx, provider, model, filteredOlder)
		if err != nil {
			health.ContextMetrics.RecordError()
			logger.WarnCF(
				"agent",
				"Failed to generate summary for context compaction, falling back to filter only",
				map[string]any{"error": err.Error()},
			)
			// Fallback: keep filteredOlder to avoid orphaned tool responses at start of newerHalf
			newMessages := []providers.Message{sysPrompt}
			newMessages = append(newMessages, filteredOlder...)
			newMessages = append(newMessages, newerHalf...)
			newMessages = append(newMessages, lastMsg)
			return newMessages, nil
		}
		summary = s
		if c.Cache != nil {
			go c.Cache.StoreSummary(sessionKey, topic, summary, len(summary)*2/5)
		}
	}

	// Use "user" role for the summary message: providers only allow "system" at position 0,
	// and a second "system" message at position 1 can trigger a 400 error on strict providers.
	compressedMsg := providers.Message{
		Role:    "user",
		Content: fmt.Sprintf("[System Note: The following is a summary of earlier conversation]\n%s", summary),
	}

	newMessages := []providers.Message{sysPrompt, compressedMsg}
	newMessages = append(newMessages, newerHalf...)
	newMessages = append(newMessages, lastMsg)

	return newMessages, nil
}

func (c *DefaultContextCompactor) GenerateSummary(
	ctx context.Context,
	provider providers.LLMProvider,
	model string,
	messages []providers.Message,
) (string, error) {
	var sb strings.Builder
	sb.WriteString(
		"Summarize the following conversation segment concisely, keeping essential facts, requested instructions, and tool results. Drop conversational filler.\n\n",
	)
	for _, m := range messages {
		switch m.Role {
		case "user", "assistant":
			content := m.Content
			if len(m.ToolCalls) > 0 {
				calls := []string{}
				for _, tc := range m.ToolCalls {
					calls = append(calls, fmt.Sprintf("Call: %s(%s)", tc.Function.Name, tc.Function.Arguments))
				}
				content = fmt.Sprintf("%s [Tool Calls: %s]", content, strings.Join(calls, ", "))
			}
			fmt.Fprintf(&sb, "%s: %s\n", m.Role, content)
		case "tool":
			fmt.Fprintf(&sb, "tool result (ID %s): %s\n", m.ToolCallID, m.Content)
		}
	}

	prompt := sb.String()
	resp, err := provider.Chat(
		ctx,
		[]providers.Message{{Role: "user", Content: prompt}},
		nil,
		model,
		map[string]any{
			"max_tokens":  512,
			"temperature": 0.3,
		},
	)
	if err != nil {
		return "", err
	}
	return resp.Content, nil
}

func (c *DefaultContextCompactor) ExtractKeyContext(old, new []providers.Message) string {
	// For future implementation if needed by Phase 3 caching
	return ""
}
