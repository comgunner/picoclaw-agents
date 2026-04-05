// ============================================================================
// ⚠️  CRITICAL: ContextManager Interface — DO NOT REMOVE
// ============================================================================
//
// This interface provides pluggable context management (legacy + seahorse).
// It was integrated from picoclaw_original on 2026-04-05 to prevent OpenRouter
// Free tier 402 errors by performing budget-aware context assembly BEFORE
// BuildMessages() is called.
//
// DO NOT remove this interface or its implementations. The AgentLoop depends
// on ContextManager.Assemble() to verify token budget before LLM calls.
//
// See: local_work/MEMORY.md, local_work/openrouter_free_token_fix.md
// ============================================================================

package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/comgunner/picoclaw/pkg/providers"
)

// ContextManager manages conversation context via a pluggable strategy.
// Exactly ONE ContextManager is active per AgentLoop, selected by config.
// Built-in implementations: "legacy" (default) and "seahorse" (SQLite-backed).
type ContextManager interface {
	// Assemble builds budget-aware context from the ContextManager's own storage.
	// Called before BuildMessages. Returns assembled messages ready for LLM.
	Assemble(ctx context.Context, req *AssembleRequest) (*AssembleResponse, error)

	// Compact compresses conversation history.
	// Called after turn completes (may be async internally) and on context overflow (sync).
	Compact(ctx context.Context, req *CompactRequest) error

	// Ingest records a message into the ContextManager's own storage.
	// Called after each message is persisted to session JSONL.
	Ingest(ctx context.Context, req *IngestRequest) error
}

// AssembleRequest is the input to Assemble.
type AssembleRequest struct {
	SessionKey string // session identifier
	Budget     int    // context window in tokens
	MaxTokens  int    // max response tokens
}

// AssembleResponse is the output of Assemble.
type AssembleResponse struct {
	History []providers.Message // assembled conversation history for BuildMessages
	Summary string              // conversation summary embedded into system prompt by BuildMessages
}

// CompactRequest is the input to Compact.
type CompactRequest struct {
	SessionKey string                // session identifier
	Reason     ContextCompressReason // proactive | retry | summarize
	Budget     int                   // context window budget (used for retry aggressive compaction)
}

// IngestRequest is the input to Ingest.
type IngestRequest struct {
	SessionKey string            // session identifier
	Message    providers.Message // the message just persisted
}

// ContextCompressReason indicates why context compression was triggered.
type ContextCompressReason string

const (
	ContextCompressReasonProactive = ContextCompressReason("proactive") // budget check before LLM
	ContextCompressReasonRetry     = ContextCompressReason("retry")     // LLM returned context error
	ContextCompressReasonSummarize = ContextCompressReason("summarize") // periodic summarization
)

// ContextManagerFactory constructs a ContextManager from config.
// al provides access to the AgentLoop's runtime resources (provider, model, workspace, etc.)
// cfg is the raw JSON configuration from config.json (may be nil).
type ContextManagerFactory func(cfg json.RawMessage, al *AgentLoop) (ContextManager, error)

var (
	cmRegistryMu sync.RWMutex
	cmRegistry   = map[string]ContextManagerFactory{}
)

// RegisterContextManager registers a named ContextManager factory.
func RegisterContextManager(name string, factory ContextManagerFactory) error {
	if name == "" {
		return fmt.Errorf("context manager name is required")
	}
	if factory == nil {
		return fmt.Errorf("context manager %q factory is nil", name)
	}

	cmRegistryMu.Lock()
	defer cmRegistryMu.Unlock()

	if _, exists := cmRegistry[name]; exists {
		return fmt.Errorf("context manager %q is already registered", name)
	}
	cmRegistry[name] = factory
	return nil
}

func lookupContextManager(name string) (ContextManagerFactory, bool) {
	cmRegistryMu.RLock()
	defer cmRegistryMu.RUnlock()

	f, ok := cmRegistry[name]
	return f, ok
}
