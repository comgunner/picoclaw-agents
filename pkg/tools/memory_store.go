// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/comgunner/picoclaw/pkg/logger"
)

// MemoryStoreTool provides persistent key-value storage.
type MemoryStoreTool struct {
	workspace string
	store     map[string]MemoryEntry
	storeMu   sync.RWMutex
	dbPath    string
}

// MemoryEntry represents a stored value with metadata.
type MemoryEntry struct {
	Value     any        `json:"value"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

// NewMemoryStoreTool creates a new MemoryStoreTool instance.
func NewMemoryStoreTool(workspace string) *MemoryStoreTool {
	tool := &MemoryStoreTool{
		workspace: workspace,
		store:     make(map[string]MemoryEntry),
	}

	// Load existing store from disk
	tool.dbPath = filepath.Join(workspace, "memory_store.json")
	tool.loadFromDisk()

	return tool
}

// Name returns the tool name.
func (t *MemoryStoreTool) Name() string {
	return "memory_store"
}

// Description returns the tool description.
func (t *MemoryStoreTool) Description() string {
	return "Persistent key-value storage for agent memory. Store, retrieve, update, and delete values with optional expiration. Use action='set' to store, 'get' to retrieve, 'delete' to remove, 'list' to see all keys, or 'clear' to reset."
}

// Parameters returns the JSON schema for tool parameters.
func (t *MemoryStoreTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"action": map[string]any{
				"type":        "string",
				"description": "Action: 'set', 'get', 'delete', 'list', 'clear', or 'count'",
				"enum":        []string{"set", "get", "delete", "list", "clear", "count"},
			},
			"key": map[string]any{
				"type":        "string",
				"description": "Key name (required for set, get, delete)",
			},
			"value": map[string]any{
				"type":        "object",
				"description": "Value to store (required for set)",
			},
			"ttl_seconds": map[string]any{
				"type":        "integer",
				"description": "Time-to-live in seconds (optional, for expiration)",
			},
		},
		"required": []string{"action"},
	}
}

// Execute runs the memory store tool.
func (t *MemoryStoreTool) Execute(ctx context.Context, args map[string]any) *ToolResult {
	action, ok := args["action"].(string)
	if !ok {
		return ErrorResult("action is required and must be one of: set, get, delete, list, clear, count")
	}

	switch action {
	case "set":
		key, _ := args["key"].(string)
		value := args["value"]
		var ttlSeconds int
		switch v := args["ttl_seconds"].(type) {
		case float64:
			ttlSeconds = int(v)
		case int:
			ttlSeconds = v
		}
		return t.set(key, value, ttlSeconds)
	case "get":
		key, _ := args["key"].(string)
		return t.get(key)
	case "delete":
		key, _ := args["key"].(string)
		return t.delete(key)
	case "list":
		return t.list()
	case "clear":
		return t.clear()
	case "count":
		return t.count()
	default:
		return ErrorResult(fmt.Sprintf("unknown action: %s. Valid options: set, get, delete, list, clear, count", action))
	}
}

// set stores a value with optional TTL.
func (t *MemoryStoreTool) set(key string, value any, ttlSeconds int) *ToolResult {
	if key == "" {
		return ErrorResult("key is required for action='set'")
	}

	t.storeMu.Lock()
	defer t.storeMu.Unlock()

	now := time.Now()
	entry := MemoryEntry{
		Value:     value,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if ttlSeconds > 0 {
		expiresAt := now.Add(time.Duration(ttlSeconds) * time.Second)
		entry.ExpiresAt = &expiresAt
	}

	t.store[key] = entry

	// Persist to disk
	if err := t.saveToDisk(); err != nil {
		logger.ErrorCF("tool", "Failed to persist memory store",
			map[string]any{
				"tool":  "memory_store",
				"error": err.Error(),
			})
		// Continue anyway, in-memory storage still works
	}

	result := map[string]any{
		"key":        key,
		"stored":     true,
		"expires_at": entry.ExpiresAt,
		"timestamp":  now.Format(time.RFC3339),
	}

	_ = result // For future structured output
	if ttlSeconds > 0 {
		return SilentResult(fmt.Sprintf("Stored '%s' with TTL of %d seconds", key, ttlSeconds))
	}
	return SilentResult(fmt.Sprintf("Stored '%s' in memory", key))
}

// get retrieves a value by key.
func (t *MemoryStoreTool) get(key string) *ToolResult {
	if key == "" {
		return ErrorResult("key is required for action='get'")
	}

	t.storeMu.RLock()
	defer t.storeMu.RUnlock()

	entry, exists := t.store[key]
	if !exists {
		return ErrorResult(fmt.Sprintf("key '%s' not found", key))
	}

	// Check expiration
	if entry.ExpiresAt != nil && time.Now().After(*entry.ExpiresAt) {
		return ErrorResult(fmt.Sprintf("key '%s' has expired", key))
	}

	result := map[string]any{
		"key":        key,
		"value":      entry.Value,
		"created_at": entry.CreatedAt,
		"updated_at": entry.UpdatedAt,
		"expires_at": entry.ExpiresAt,
	}

	_ = result // For future structured output
	return SilentResult(fmt.Sprintf("Retrieved '%s': %v", key, entry.Value))
}

// delete removes a key from the store.
func (t *MemoryStoreTool) delete(key string) *ToolResult {
	if key == "" {
		return ErrorResult("key is required for action='delete'")
	}

	t.storeMu.Lock()
	defer t.storeMu.Unlock()

	if _, exists := t.store[key]; !exists {
		return ErrorResult(fmt.Sprintf("key '%s' not found", key))
	}

	delete(t.store, key)

	// Persist to disk
	if err := t.saveToDisk(); err != nil {
		logger.ErrorCF("tool", "Failed to persist memory store",
			map[string]any{
				"tool":  "memory_store",
				"error": err.Error(),
			})
	}

	result := map[string]any{
		"key":       key,
		"deleted":   true,
		"timestamp": time.Now().Format(time.RFC3339),
	}

	_ = result // For future structured output
	return SilentResult(fmt.Sprintf("Deleted '%s' from memory", key))
}

// list returns all keys in the store.
func (t *MemoryStoreTool) list() *ToolResult {
	t.storeMu.RLock()
	defer t.storeMu.RUnlock()

	keys := make([]string, 0, len(t.store))
	for key := range t.store {
		keys = append(keys, key)
	}

	result := map[string]any{
		"keys":      keys,
		"count":     len(keys),
		"timestamp": time.Now().Format(time.RFC3339),
	}

	_ = result // For future structured output
	return SilentResult(fmt.Sprintf("Memory store contains %d keys: %v", len(keys), keys))
}

// clear removes all entries from the store.
func (t *MemoryStoreTool) clear() *ToolResult {
	t.storeMu.Lock()
	defer t.storeMu.Unlock()

	count := len(t.store)
	t.store = make(map[string]MemoryEntry)

	// Persist to disk (empty store)
	if err := t.saveToDisk(); err != nil {
		logger.ErrorCF("tool", "Failed to persist memory store",
			map[string]any{
				"tool":  "memory_store",
				"error": err.Error(),
			})
	}

	result := map[string]any{
		"cleared":   count,
		"timestamp": time.Now().Format(time.RFC3339),
	}

	_ = result // For future structured output
	return SilentResult(fmt.Sprintf("Cleared %d entries from memory store", count))
}

// count returns the number of entries in the store.
func (t *MemoryStoreTool) count() *ToolResult {
	t.storeMu.RLock()
	defer t.storeMu.RUnlock()

	result := map[string]any{
		"count":     len(t.store),
		"timestamp": time.Now().Format(time.RFC3339),
	}

	_ = result // For future structured output
	return SilentResult(fmt.Sprintf("Memory store contains %d entries", len(t.store)))
}

// saveToDisk persists the store to disk.
// saveToDisk must be called with storeMu already held by the caller.
func (t *MemoryStoreTool) saveToDisk() error {
	data, err := json.MarshalIndent(t.store, "", "  ")
	if err != nil {
		return err
	}

	// Ensure directory exists
	dir := filepath.Dir(t.dbPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	// Write atomically (temp file + rename)
	tmpPath := t.dbPath + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0o644); err != nil {
		return err
	}

	return os.Rename(tmpPath, t.dbPath)
}

// loadFromDisk loads the store from disk.
func (t *MemoryStoreTool) loadFromDisk() error {
	t.storeMu.Lock()
	defer t.storeMu.Unlock()

	data, err := os.ReadFile(t.dbPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No existing store, start fresh
		}
		return err
	}

	var store map[string]MemoryEntry
	if err := json.Unmarshal(data, &store); err != nil {
		return err
	}

	t.store = store
	return nil
}
