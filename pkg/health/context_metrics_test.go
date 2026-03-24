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
	"testing"
)

func TestContextMetrics(t *testing.T) {
	metrics := &ContextMetricsStruct{}

	metrics.RecordCompaction()
	metrics.RecordTokens(100)
	metrics.RecordCacheHit()
	metrics.RecordCacheMiss()
	metrics.RecordError()

	data := metrics.GetMetrics()

	if data["compaction_count"] != 1 {
		t.Errorf("expected 1 compaction, got %d", data["compaction_count"])
	}
	if data["total_tokens"] != 100 {
		t.Errorf("expected 100 tokens, got %d", data["total_tokens"])
	}
	if data["cache_hits"] != 1 {
		t.Errorf("expected 1 cache hit, got %d", data["cache_hits"])
	}
	if data["cache_misses"] != 1 {
		t.Errorf("expected 1 cache miss, got %d", data["cache_misses"])
	}
	if data["errors"] != 1 {
		t.Errorf("expected 1 error, got %d", data["errors"])
	}
}
