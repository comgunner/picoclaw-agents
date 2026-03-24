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
	"os"
	"testing"
)

func TestSummaryCache_StoreAndRetrieve(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cache_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	cache := NewSummaryCache(tmpDir)

	// Test Store
	cache.StoreSummary("session_1", "test_topic", "This is a summary", 10)

	// Test Retrieve
	summary, found := cache.FindSimilarSummary("test_topic")
	if !found {
		t.Errorf("expected to find summary for topic")
	}
	if summary != "This is a summary" {
		t.Errorf("expected summary text to match, got %s", summary)
	}

	// Invalid Topic
	_, found = cache.FindSimilarSummary("invalid_topic")
	if found {
		t.Errorf("did not expect to find summary for invalid topic")
	}
}
