// PicoClaw-Agents - Seahorse ContextManager adapter
// Ported from picoclaw_original/pkg/agent/context_seahorse.go
// Adapted for fork's types (no Media field, SessionManager instead of SessionStore).

package agent

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/comgunner/picoclaw/pkg/logger"
	"github.com/comgunner/picoclaw/pkg/providers"
	"github.com/comgunner/picoclaw/pkg/providers/protocoltypes"
	"github.com/comgunner/picoclaw/pkg/seahorse"
	"github.com/comgunner/picoclaw/pkg/tokenizer"
)

// seahorseContextManager adapts seahorse.Engine to agent.ContextManager.
type seahorseContextManager struct {
	engine *seahorse.Engine
	agent  *AgentInstance // for bootstrap and provider access
}

// newSeahorseContextManager creates a seahorse-backed ContextManager.
func newSeahorseContextManager(cfg json.RawMessage, al *AgentLoop) (ContextManager, error) {
	if al == nil {
		return nil, fmt.Errorf("seahorse: AgentLoop is required")
	}

	agent := al.registry.GetDefaultAgent()
	if agent == nil {
		return nil, fmt.Errorf("seahorse: no default agent available")
	}
	if agent.Provider == nil {
		return nil, fmt.Errorf("seahorse: no provider available for agent %q", agent.ID)
	}

	// DB stores session data, so it goes in sessions/ directory
	dbPath := agent.Workspace + "/sessions/seahorse.db"

	// Create CompleteFn from provider
	completeFn := providerToCompleteFn(agent.Provider, agent.Model)

	// Create engine
	engine, err := seahorse.NewEngine(seahorse.Config{
		DBPath: dbPath,
	}, completeFn)
	if err != nil {
		return nil, fmt.Errorf("seahorse: create engine: %w", err)
	}

	mgr := &seahorseContextManager{
		engine: engine,
		agent:  agent,
	}

	// Bootstrap all existing sessions at startup
	if agent.Sessions != nil {
		ctx := context.Background()
		for _, sessionKey := range agent.Sessions.ListSessions() {
			mgr.bootstrapSession(ctx, sessionKey)
		}
	}

	logger.InfoCF("agent", "Seahorse ContextManager initialized", map[string]any{
		"db_path":  dbPath,
		"agent_id": agent.ID,
	})

	return mgr, nil
}

// providerToCompleteFn wraps providers.LLMProvider as a seahorse.CompleteFn.
func providerToCompleteFn(provider providers.LLMProvider, model string) seahorse.CompleteFn {
	return func(ctx context.Context, prompt string, opts seahorse.CompleteOptions) (string, error) {
		resp, err := provider.Chat(
			ctx,
			[]providers.Message{{Role: "user", Content: prompt}},
			nil, // no tools for summarization
			model,
			map[string]any{
				"max_tokens":       opts.MaxTokens,
				"temperature":      opts.Temperature,
				"prompt_cache_key": "seahorse",
			},
		)
		if err != nil {
			return "", err
		}
		return resp.Content, nil
	}
}

// Assemble builds budget-aware context from seahorse SQLite.
func (m *seahorseContextManager) Assemble(ctx context.Context, req *AssembleRequest) (*AssembleResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("seahorse assemble: nil request")
	}

	budget := req.Budget
	if budget <= 0 {
		budget = 100000
	}

	// Reserve space for model response
	effectiveBudget := budget - req.MaxTokens
	if effectiveBudget <= 0 {
		// MaxTokens >= budget is a configuration problem
		// Use 50% as minimum to avoid guaranteed overflow
		logger.WarnCF("agent", "MaxTokens >= budget, using 50% fallback",
			map[string]any{"budget": budget, "max_tokens": req.MaxTokens})
		effectiveBudget = budget / 2
	}

	result, err := m.engine.Assemble(ctx, req.SessionKey, seahorse.AssembleInput{
		Budget: effectiveBudget,
	})
	if err != nil {
		return nil, fmt.Errorf("seahorse assemble: %w", err)
	}

	history := seahorseToProviderMessages(result)

	return &AssembleResponse{
		History: history,
		Summary: result.Summary,
	}, nil
}

// Compact compresses conversation history via seahorse summarization.
func (m *seahorseContextManager) Compact(ctx context.Context, req *CompactRequest) error {
	if req == nil {
		return nil
	}

	// For retry (LLM overflow), use aggressive CompactUntilUnder to guarantee
	// context shrinks below budget.
	if req.Reason == ContextCompressReasonRetry && req.Budget > 0 {
		_, err := m.engine.CompactUntilUnder(ctx, req.SessionKey, req.Budget)
		return err
	}

	_, err := m.engine.Compact(ctx, req.SessionKey, seahorse.CompactInput{
		Force:  req.Reason == ContextCompressReasonRetry,
		Budget: &req.Budget,
	})
	return err
}

// Ingest records a message into seahorse SQLite.
func (m *seahorseContextManager) Ingest(ctx context.Context, req *IngestRequest) error {
	if req == nil {
		return nil
	}

	msg := providerToSeahorseMessage(req.Message)
	_, err := m.engine.Ingest(ctx, req.SessionKey, []seahorse.Message{msg})
	return err
}

// bootstrapSession reconciles JSONL session history into seahorse SQLite.
func (m *seahorseContextManager) bootstrapSession(ctx context.Context, sessionKey string) {
	if m.agent.Sessions == nil {
		return
	}

	history := m.agent.Sessions.GetHistory(sessionKey)
	if len(history) == 0 {
		return
	}

	// Convert provider messages to seahorse messages
	msgs := make([]seahorse.Message, len(history))
	for i, h := range history {
		msgs[i] = providerToSeahorseMessage(h)
	}

	if err := m.engine.Bootstrap(ctx, sessionKey, msgs); err != nil {
		logger.WarnCF("seahorse", "bootstrap session failed", map[string]any{
			"session": sessionKey,
			"error":   err.Error(),
		})
	}
}

// providerToSeahorseMessage converts a providers.Message to a seahorse.Message.
func providerToSeahorseMessage(msg protocoltypes.Message) seahorse.Message {
	result := seahorse.Message{
		Role:             msg.Role,
		Content:          msg.Content,
		ReasoningContent: msg.ReasoningContent,
		TokenCount:       tokenizer.EstimateMessageTokens(msg),
	}

	// Convert ToolCalls → MessageParts
	for _, tc := range msg.ToolCalls {
		part := seahorse.MessagePart{
			Type:       "tool_use",
			Name:       tc.Function.Name,
			Arguments:  tc.Function.Arguments,
			ToolCallID: tc.ID,
		}
		result.Parts = append(result.Parts, part)
	}

	// Convert tool result
	if msg.ToolCallID != "" {
		part := seahorse.MessagePart{
			Type:       "tool_result",
			ToolCallID: msg.ToolCallID,
			Text:       msg.Content,
		}
		result.Parts = append(result.Parts, part)
	}

	// NOTE: Fork's Message type does not have Media field — skip media conversion.

	return result
}

// seahorseToProviderMessages converts a seahorse.AssembleResult to []providers.Message.
func seahorseToProviderMessages(result *seahorse.AssembleResult) []protocoltypes.Message {
	messages := make([]protocoltypes.Message, 0, len(result.Messages))

	for _, msg := range result.Messages {
		pm := protocoltypes.Message{
			Role:             msg.Role,
			Content:          msg.Content,
			ReasoningContent: msg.ReasoningContent,
		}

		// Reconstruct ToolCalls from parts
		for _, part := range msg.Parts {
			if part.Type == "tool_use" {
				pm.ToolCalls = append(pm.ToolCalls, protocoltypes.ToolCall{
					ID:   part.ToolCallID,
					Type: "function",
					Function: &protocoltypes.FunctionCall{
						Name:      part.Name,
						Arguments: part.Arguments,
					},
				})
			}
			if part.Type == "tool_result" {
				pm.ToolCallID = part.ToolCallID
				if pm.Content == "" && part.Text != "" {
					pm.Content = part.Text
				}
			}
		}

		messages = append(messages, pm)
	}

	return messages
}

func init() {
	if err := RegisterContextManager("seahorse", newSeahorseContextManager); err != nil {
		// Don't panic — seahorse is optional. Log warning and continue with legacy.
		logger.WarnCF("agent", "failed to register seahorse ContextManager (non-fatal)", map[string]any{
			"error": err.Error(),
		})
	}
}
