// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/comgunner/picoclaw/pkg/logger"
)

// SummaryEntry represents a cached summary of a conversation
type SummaryEntry struct {
	ID        string    `json:"id"`
	SessionID string    `json:"session_id"`
	Summary   string    `json:"summary"`
	Tokens    int       `json:"tokens"`
	Topic     string    `json:"topic"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// SummaryCache provides a persistent cache for conversation summaries
type SummaryCache struct {
	mu       sync.RWMutex
	filePath string
	entries  map[string]SummaryEntry
}

// NewSummaryCache initializes the summary cache backed by a JSON file.
// In a full DB setup, this would connect to the DB. For now, it uses a JSON file.
func NewSummaryCache(workspace string) *SummaryCache {
	cacheDir := filepath.Join(workspace, "cache")
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		logger.WarnCF("cache", "Failed to create cache directory", map[string]any{"error": err.Error()})
	}

	cacheFile := filepath.Join(cacheDir, "summaries.json")
	cache := &SummaryCache{
		filePath: cacheFile,
		entries:  make(map[string]SummaryEntry),
	}

	cache.load()
	return cache
}

func (c *SummaryCache) load() {
	c.mu.Lock()
	defer c.mu.Unlock()

	data, err := os.ReadFile(c.filePath)
	if err != nil {
		if !os.IsNotExist(err) {
			logger.WarnCF("cache", "Failed to read summary cache", map[string]any{"error": err.Error()})
		}
		return
	}

	if err := json.Unmarshal(data, &c.entries); err != nil {
		logger.WarnCF("cache", "Failed to unmarshal summary cache", map[string]any{"error": err.Error()})
	}
}

func (c *SummaryCache) save() {
	data, err := json.MarshalIndent(c.entries, "", "  ")
	if err != nil {
		logger.WarnCF("cache", "Failed to marshal summary cache", map[string]any{"error": err.Error()})
		return
	}

	if err := os.WriteFile(c.filePath, data, 0o644); err != nil {
		logger.WarnCF("cache", "Failed to write summary cache", map[string]any{"error": err.Error()})
	}
}

// StoreSummary saves a summary into the cache
func (c *SummaryCache) StoreSummary(sessionID, topic, summary string, tokens int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	id := sessionID + "_" + time.Now().Format("20060102150405")
	entry := SummaryEntry{
		ID:        id,
		SessionID: sessionID,
		Summary:   summary,
		Tokens:    tokens,
		Topic:     topic,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 7 days TTL
	}

	c.entries[id] = entry
	c.save()
}

// FindSimilarSummary looks for a recent summary with the same topic and sessionID.
// BUG-02 FIX: Added sessionID parameter to filter by session - previously returned
// summaries from any session, causing incorrect context injection across different conversations.
func (c *SummaryCache) FindSimilarSummary(sessionID, topic string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	now := time.Now()
	for _, entry := range c.entries {
		// FIX: Now filters by both sessionID AND topic to avoid cross-session contamination
		if entry.SessionID == sessionID && entry.Topic == topic && entry.ExpiresAt.After(now) {
			return entry.Summary, true
		}
	}
	return "", false
}
