// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package mcp

import (
	"encoding/json"
	"testing"
)

func TestMarshalRequest(t *testing.T) {
	req := jsonRPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "initialize",
		Params:  map[string]any{"key": "value"},
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify round-trip
	var parsed jsonRPCRequest
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if parsed.JSONRPC != "2.0" {
		t.Errorf("expected jsonrpc 2.0, got %q", parsed.JSONRPC)
	}
	if parsed.ID != 1 {
		t.Errorf("expected id 1, got %d", parsed.ID)
	}
	if parsed.Method != "initialize" {
		t.Errorf("expected method initialize, got %q", parsed.Method)
	}
	if parsed.Params["key"] != "value" {
		t.Errorf("expected params key=value, got %v", parsed.Params)
	}
}

func TestParseResponse(t *testing.T) {
	raw := `{"jsonrpc":"2.0","id":1,"result":{"tools":[{"name":"test"}]}}`

	var resp jsonRPCResponse
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if resp.JSONRPC != "2.0" {
		t.Errorf("expected jsonrpc 2.0, got %q", resp.JSONRPC)
	}
	if resp.ID != 1 {
		t.Errorf("expected id 1, got %d", resp.ID)
	}
	if resp.Error != nil {
		t.Errorf("expected no error, got %v", resp.Error)
	}
	if len(resp.Result) == 0 {
		t.Error("expected result to be non-empty")
	}
}

func TestParseError(t *testing.T) {
	raw := `{"jsonrpc":"2.0","id":1,"error":{"code":-32600,"message":"Invalid Request"}}`

	var resp jsonRPCResponse
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if resp.Error == nil {
		t.Fatal("expected error, got nil")
	}
	if resp.Error.Code != -32600 {
		t.Errorf("expected error code -32600, got %d", resp.Error.Code)
	}
	if resp.Error.Message != "Invalid Request" {
		t.Errorf("expected error message 'Invalid Request', got %q", resp.Error.Message)
	}
}

func TestProtocolVersion(t *testing.T) {
	if MCPProtocolVersion != "2024-11-05" {
		t.Errorf("expected protocol version 2024-11-05, got %q", MCPProtocolVersion)
	}
}

func TestMaxLineBytes(t *testing.T) {
	if MAX_LINE_BYTES != 10*1024*1024 {
		t.Errorf("expected MAX_LINE_BYTES to be 10MB, got %d", MAX_LINE_BYTES)
	}
}
