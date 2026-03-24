// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package announce

import (
	"context"
	"fmt"

	"github.com/comgunner/picoclaw/pkg/agents/spawn"
	"github.com/comgunner/picoclaw/pkg/gateway"
)

type AnnounceStep struct {
	gatewayClient gateway.Client
}

func NewAnnounceStep(client gateway.Client) *AnnounceStep {
	return &AnnounceStep{
		gatewayClient: client,
	}
}

func (a *AnnounceStep) AnnounceSubagentCompletion(
	ctx context.Context,
	childSessionKey string,
	result any,
) error {
	if a.gatewayClient == nil {
		return fmt.Errorf("gateway client is nil")
	}

	// Extract summary from result
	summary := extractSummary(result)

	// Build announce message
	announceMsg := fmt.Sprintf(
		"Subagent completion from %s:\n%s",
		childSessionKey, summary,
	)

	// Extract requester session key from child session key
	requesterSessionKey := getRequesterFromChild(childSessionKey)

	// Publish announcement to requester
	_, err := a.gatewayClient.Call(ctx, gateway.Call{
		Method: "announce",
		Params: map[string]any{
			"requesterSessionKey": requesterSessionKey,
			"message":             announceMsg,
			"sourceRunID":         childSessionKey,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to post announcement: %w", err)
	}

	return nil
}

func extractSummary(result any) string {
	if result == nil {
		return "No result returned"
	}

	switch v := result.(type) {
	case string:
		return v
	case *spawn.Run:
		return fmt.Sprintf("Run completed: %s", v.RunID)
	default:
		return fmt.Sprintf("Result: %v", v)
	}
}

func getRequesterFromChild(childSessionKey string) string {
	// Implementation would extract requester from child session key
	// For now, return a placeholder
	return "agent:requester:main"
}
