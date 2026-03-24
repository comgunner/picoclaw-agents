// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package gateway

import (
	"context"
	"testing"
)

func TestMockClient_CallAndCount(t *testing.T) {
	m := NewMockClient()
	m.SetHandler("agent", func(_ context.Context, req Call) (Response, error) {
		return Response{RunID: "run-1", Status: "accepted", Data: req.Params}, nil
	})

	resp, err := m.Call(context.Background(), Call{Method: "agent", Params: map[string]any{"message": "hello"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.RunID != "run-1" || resp.Status != "accepted" {
		t.Fatalf("unexpected response: %+v", resp)
	}
	if m.Count("agent") != 1 {
		t.Fatalf("expected count=1, got %d", m.Count("agent"))
	}
}
