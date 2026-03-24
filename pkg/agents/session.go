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
	"fmt"
	"strings"

	"github.com/google/uuid"
)

const (
	SessionKindMain     = "main"
	SessionKindSubagent = "subagent"
)

// SessionKey describes a normalized agent session key.
// Supported formats:
// - agent:<agentId>:main
// - agent:<agentId>:subagent:<uuid>
type SessionKey struct {
	AgentID   string
	Kind      string
	SessionID string
	Depth     int
}

func ParseSessionKey(raw string) (SessionKey, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return SessionKey{}, fmt.Errorf("session key is empty")
	}

	parts := strings.Split(value, ":")
	if len(parts) < 3 || parts[0] != "agent" {
		return SessionKey{}, fmt.Errorf("invalid session key format: %q", raw)
	}

	agentID := strings.TrimSpace(parts[1])
	if agentID == "" {
		return SessionKey{}, fmt.Errorf("agent id is empty")
	}

	if len(parts) == 3 && strings.TrimSpace(parts[2]) == SessionKindMain {
		return SessionKey{
			AgentID: agentID,
			Kind:    SessionKindMain,
			Depth:   0,
		}, nil
	}

	if len(parts) == 4 && strings.TrimSpace(parts[2]) == SessionKindSubagent {
		sessionID := strings.TrimSpace(parts[3])
		if _, err := uuid.Parse(sessionID); err != nil {
			return SessionKey{}, fmt.Errorf("invalid subagent session uuid: %w", err)
		}
		return SessionKey{
			AgentID:   agentID,
			Kind:      SessionKindSubagent,
			SessionID: sessionID,
			Depth:     1,
		}, nil
	}

	return SessionKey{}, fmt.Errorf("unsupported session key format: %q", raw)
}

func GenerateSubagentSessionKey(agentID string) string {
	normalized := strings.TrimSpace(strings.ToLower(agentID))
	if normalized == "" {
		normalized = "main"
	}
	normalized = strings.ReplaceAll(normalized, "/", "_")
	normalized = strings.ReplaceAll(normalized, "\\", "_")
	return fmt.Sprintf("agent:%s:subagent:%s", normalized, uuid.NewString())
}
