// PicoClaw-Agents - Legacy ContextManager implementation
// Ported from picoclaw_original/pkg/agent/context_legacy.go
// Adapted for fork's types. Provides default fallback when seahorse is not configured.

package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/comgunner/picoclaw/pkg/logger"
	"github.com/comgunner/picoclaw/pkg/providers"
	"github.com/comgunner/picoclaw/pkg/tokenizer"
)

// Default thresholds for summarization when not configured.
const (
	defaultSummarizeMessageThreshold = 20 // Number of messages before summarization
	defaultSummarizeTokenPercent     = 75 // Percentage of context window that triggers summarization
)

// legacyContextManager wraps the existing summarization/compression logic
// as a ContextManager implementation. It is the default when no other
// ContextManager is configured.
type legacyContextManager struct {
	al          *AgentLoop
	summarizing sync.Map // dedup for async Compact (post-turn)
}

func (m *legacyContextManager) Assemble(_ context.Context, req *AssembleRequest) (*AssembleResponse, error) {
	// Legacy: read history from session, return as-is.
	// Budget enforcement happens in BuildMessages caller via
	// isOverContextBudget + forceCompression.
	agent := m.al.registry.GetDefaultAgent()
	if agent == nil {
		return &AssembleResponse{}, nil
	}
	history := agent.Sessions.GetHistory(req.SessionKey)
	summary := agent.Sessions.GetSummary(req.SessionKey)
	return &AssembleResponse{
		History: history,
		Summary: summary,
	}, nil
}

func (m *legacyContextManager) Compact(_ context.Context, req *CompactRequest) error {
	switch req.Reason {
	case ContextCompressReasonProactive, ContextCompressReasonRetry:
		// Sync emergency compression — budget exceeded.
		if result, ok := m.forceCompression(req.SessionKey); ok {
			logger.WarnCF("agent", "Context compression (legacy)", map[string]any{
				"session_key":  req.SessionKey,
				"reason":       string(req.Reason),
				"dropped_msgs": result.DroppedMessages,
				"remaining":    result.RemainingMessages,
			})
		}
	case ContextCompressReasonSummarize:
		m.maybeSummarize(req.SessionKey)
	}
	return nil
}

func (m *legacyContextManager) Ingest(_ context.Context, _ *IngestRequest) error {
	// Legacy: no-op. Messages are persisted by Sessions JSONL.
	return nil
}

// maybeSummarize triggers summarization if the session history exceeds thresholds.
// It runs asynchronously in a goroutine.
func (m *legacyContextManager) maybeSummarize(sessionKey string) {
	agent := m.al.registry.GetDefaultAgent()
	if agent == nil {
		return
	}

	newHistory := agent.Sessions.GetHistory(sessionKey)
	tokenEstimate := m.estimateTokens(newHistory)
	threshold := agent.ContextWindow * defaultSummarizeTokenPercent / 100

	if len(newHistory) > defaultSummarizeMessageThreshold || tokenEstimate > threshold {
		summarizeKey := agent.ID + ":" + sessionKey
		if _, loading := m.summarizing.LoadOrStore(summarizeKey, true); !loading {
			go func() {
				defer m.summarizing.Delete(summarizeKey)
				defer func() {
					if r := recover(); r != nil {
						logger.WarnCF("agent", "Summarization panic recovered", map[string]any{
							"session_key": sessionKey,
							"panic":       r,
						})
					}
				}()
				logger.Debug("Memory threshold reached. Optimizing conversation history...")
				m.summarizeSession(agent, sessionKey)
			}()
		}
	}
}

type compressionResult struct {
	DroppedMessages   int
	RemainingMessages int
}

// forceCompression aggressively reduces context when the limit is hit.
// It drops the oldest ~50% of Turns (a Turn is a complete user→LLM→response
// cycle), so tool-call sequences are never split.
func (m *legacyContextManager) forceCompression(sessionKey string) (compressionResult, bool) {
	agent := m.al.registry.GetDefaultAgent()
	if agent == nil {
		return compressionResult{}, false
	}

	history := agent.Sessions.GetHistory(sessionKey)
	if len(history) <= 2 {
		return compressionResult{}, false
	}

	turns := parseTurnBoundaries(history)
	var mid int
	if len(turns) >= 2 {
		mid = turns[len(turns)/2]
	} else {
		mid = findSafeBoundary(history, len(history)/2)
	}
	var keptHistory []providers.Message
	if mid <= 0 {
		for i := len(history) - 1; i >= 0; i-- {
			if history[i].Role == "user" {
				keptHistory = []providers.Message{history[i]}
				break
			}
		}
	} else {
		keptHistory = history[mid:]
	}

	droppedCount := len(history) - len(keptHistory)

	existingSummary := agent.Sessions.GetSummary(sessionKey)
	compressionNote := fmt.Sprintf(
		"[Emergency compression dropped %d oldest messages due to context limit]",
		droppedCount,
	)
	if existingSummary != "" {
		compressionNote = existingSummary + "\n\n" + compressionNote
	}
	agent.Sessions.SetSummary(sessionKey, compressionNote)

	agent.Sessions.SetHistory(sessionKey, keptHistory)
	agent.Sessions.Save(sessionKey)

	logger.WarnCF("agent", "Forced compression executed (legacy)", map[string]any{
		"session_key":  sessionKey,
		"dropped_msgs": droppedCount,
		"new_count":    len(keptHistory),
	})

	return compressionResult{
		DroppedMessages:   droppedCount,
		RemainingMessages: len(keptHistory),
	}, true
}

func (m *legacyContextManager) summarizeSession(agent *AgentInstance, sessionKey string) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	history := agent.Sessions.GetHistory(sessionKey)
	summary := agent.Sessions.GetSummary(sessionKey)

	if len(history) <= 4 {
		return
	}

	safeCut := findSafeBoundary(history, len(history)-4)
	if safeCut <= 0 {
		return
	}
	keepCount := len(history) - safeCut
	toSummarize := history[:safeCut]

	maxMessageTokens := agent.ContextWindow / 2
	validMessages := make([]providers.Message, 0)
	omitted := false

	for _, msg := range toSummarize {
		if msg.Role != "user" && msg.Role != "assistant" {
			continue
		}
		msgTokens := tokenizer.EstimateMessageTokens(msg)
		if msgTokens > maxMessageTokens {
			omitted = true
			continue
		}
		validMessages = append(validMessages, msg)
	}

	if len(validMessages) == 0 {
		return
	}

	const maxSummarizationMessages = 10

	var finalSummary string
	if len(validMessages) > maxSummarizationMessages {
		mid := len(validMessages) / 2
		mid = m.findNearestUserMessage(validMessages, mid)

		part1 := validMessages[:mid]
		part2 := validMessages[mid:]

		s1, _ := m.summarizeBatch(ctx, agent, part1, "")
		s2, _ := m.summarizeBatch(ctx, agent, part2, "")

		mergePrompt := fmt.Sprintf(
			"Merge these two conversation summaries into one cohesive summary:\n\n1: %s\n\n2: %s",
			s1, s2,
		)

		resp, err := m.retryLLMCall(ctx, agent, mergePrompt, 3)
		if err == nil && resp.Content != "" {
			finalSummary = resp.Content
		} else {
			finalSummary = s1 + " " + s2
		}
	} else {
		finalSummary, _ = m.summarizeBatch(ctx, agent, validMessages, summary)
	}

	if omitted && finalSummary != "" {
		finalSummary += "\n[Note: Some oversized messages were omitted from this summary for efficiency.]"
	}

	if finalSummary != "" {
		agent.Sessions.SetSummary(sessionKey, finalSummary)
		agent.Sessions.TruncateHistory(sessionKey, keepCount)
		agent.Sessions.Save(sessionKey)

		logger.InfoCF("agent", "Session summarized (legacy)", map[string]any{
			"session_key": sessionKey,
			"summarized":  len(validMessages),
			"kept":        keepCount,
			"summary_len": len(finalSummary),
			"omitted":     omitted,
		})
	}
}

func (m *legacyContextManager) findNearestUserMessage(messages []providers.Message, mid int) int {
	originalMid := mid

	for mid > 0 && messages[mid].Role != "user" {
		mid--
	}

	if messages[mid].Role == "user" {
		return mid
	}

	mid = originalMid
	for mid < len(messages) && messages[mid].Role != "user" {
		mid++
	}

	if mid < len(messages) {
		return mid
	}

	return originalMid
}

func (m *legacyContextManager) retryLLMCall(
	ctx context.Context,
	agent *AgentInstance,
	prompt string,
	maxRetries int,
) (*providers.LLMResponse, error) {
	const llmTemperature = 0.3

	var mu sync.Mutex
	var resp *providers.LLMResponse
	var err error

	for attempt := 0; attempt < maxRetries; attempt++ {
		mu.Lock()
		resp, err = agent.Provider.Chat(
			ctx,
			[]providers.Message{{Role: "user", Content: prompt}},
			nil,
			agent.Model,
			map[string]any{
				"max_tokens":       agent.MaxTokens,
				"temperature":      llmTemperature,
				"prompt_cache_key": agent.ID,
			},
		)
		mu.Unlock()

		if err == nil && resp != nil && resp.Content != "" {
			return resp, nil
		}
		if attempt < maxRetries-1 {
			time.Sleep(time.Duration(attempt+1) * 100 * time.Millisecond)
		}
	}

	return resp, err
}

func (m *legacyContextManager) summarizeBatch(
	ctx context.Context,
	agent *AgentInstance,
	batch []providers.Message,
	existingSummary string,
) (string, error) {
	const (
		llmMaxRetries             = 3
		fallbackMinContentLength  = 200
		fallbackMaxContentPercent = 10
	)

	var sb strings.Builder
	sb.WriteString("Provide a concise summary of this conversation segment, preserving core context and key points.\n")
	if existingSummary != "" {
		sb.WriteString("Existing context: ")
		sb.WriteString(existingSummary)
		sb.WriteString("\n")
	}
	sb.WriteString("\nCONVERSATION:\n")
	for _, msg := range batch {
		fmt.Fprintf(&sb, "%s: %s\n", msg.Role, msg.Content)
	}
	prompt := sb.String()

	response, err := m.retryLLMCall(ctx, agent, prompt, llmMaxRetries)
	if err == nil && response.Content != "" {
		return strings.TrimSpace(response.Content), nil
	}

	// Fallback: deterministic truncation
	var fallback strings.Builder
	fallback.WriteString("Conversation summary: ")
	for i, msg := range batch {
		if i > 0 {
			fallback.WriteString(" | ")
		}
		content := strings.TrimSpace(msg.Content)
		runes := []rune(content)
		if len(runes) == 0 {
			fallback.WriteString(fmt.Sprintf("%s: ", msg.Role))
			continue
		}

		keepLength := len(runes) * fallbackMaxContentPercent / 100
		if keepLength < fallbackMinContentLength {
			keepLength = fallbackMinContentLength
		}
		if keepLength > len(runes) {
			keepLength = len(runes)
		}

		content = string(runes[:keepLength])
		if keepLength < len(runes) {
			content += "..."
		}
		fallback.WriteString(fmt.Sprintf("%s: %s", msg.Role, content))
	}
	return fallback.String(), nil
}

func (m *legacyContextManager) estimateTokens(messages []providers.Message) int {
	total := 0
	for _, msg := range messages {
		total += tokenizer.EstimateMessageTokens(msg)
	}
	return total
}

// init registers the "legacy" ContextManager as the default.
func init() {
	RegisterContextManager("legacy", func(_ json.RawMessage, al *AgentLoop) (ContextManager, error) {
		return &legacyContextManager{al: al}, nil
	})
}
