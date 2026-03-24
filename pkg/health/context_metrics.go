// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package health

import (
	"sync"
)

type ContextMetricsStruct struct {
	CompactionCount int
	TotalMessages   int
	TotalTokens     int
	CacheHits       int
	CacheMisses     int
	Errors          int
	mu              sync.RWMutex
}

var ContextMetrics = &ContextMetricsStruct{}

func (m *ContextMetricsStruct) RecordCompaction() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.CompactionCount++
}

func (m *ContextMetricsStruct) RecordTokens(tokens int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.TotalTokens += tokens
}

func (m *ContextMetricsStruct) RecordCacheHit() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.CacheHits++
}

func (m *ContextMetricsStruct) RecordCacheMiss() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.CacheMisses++
}

func (m *ContextMetricsStruct) RecordError() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Errors++
}

func (m *ContextMetricsStruct) GetMetrics() map[string]int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return map[string]int{
		"compaction_count": m.CompactionCount,
		"total_tokens":     m.TotalTokens,
		"cache_hits":       m.CacheHits,
		"cache_misses":     m.CacheMisses,
		"errors":           m.Errors,
	}
}
