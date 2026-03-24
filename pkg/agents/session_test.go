// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package agents

import (
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestParseSessionKey_Main(t *testing.T) {
	key, err := ParseSessionKey("agent:tech_lead:main")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if key.AgentID != "tech_lead" || key.Kind != SessionKindMain || key.Depth != 0 {
		t.Fatalf("unexpected key: %+v", key)
	}
}

func TestParseSessionKey_Subagent(t *testing.T) {
	id := uuid.NewString()
	key, err := ParseSessionKey("agent:coder:subagent:" + id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if key.AgentID != "coder" || key.Kind != SessionKindSubagent || key.Depth != 1 || key.SessionID != id {
		t.Fatalf("unexpected key: %+v", key)
	}
}

func TestParseSessionKey_Invalid(t *testing.T) {
	invalid := []string{"", "main", "agent::main", "agent:x:subagent:not-uuid"}
	for _, item := range invalid {
		if _, err := ParseSessionKey(item); err == nil {
			t.Fatalf("expected error for %q", item)
		}
	}
}

func TestGenerateSubagentSessionKey(t *testing.T) {
	key := GenerateSubagentSessionKey("Coder")
	if !strings.HasPrefix(key, "agent:coder:subagent:") {
		t.Fatalf("unexpected format: %s", key)
	}

	parsed, err := ParseSessionKey(key)
	if err != nil {
		t.Fatalf("generated key should parse: %v", err)
	}
	if parsed.Kind != SessionKindSubagent {
		t.Fatalf("expected subagent kind, got %+v", parsed)
	}
}
