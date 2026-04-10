// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package tools

import (
	"context"
	"testing"
	"time"

	"github.com/comgunner/picoclaw/pkg/mcp"
)

// mockMCPCaller implements MCPCaller for testing
type mockMCPCaller struct {
	result *mcp.ToolCallResult
	err    error
}

func (m *mockMCPCaller) CallTool(
	ctx context.Context,
	serverName, toolName string,
	args map[string]any,
) (*mcp.ToolCallResult, error) {
	return m.result, m.err
}

func TestMCPOperatorTool_Name(t *testing.T) {
	info := mcp.ToolInfo{
		Name:        "test_tool",
		Description: "A test tool",
		Schema:      map[string]any{},
	}
	caller := &mockMCPCaller{}
	tool := NewMCPOperatorTool("test_server", info, caller, 30*time.Second)

	expected := "mcp_test_server_test_tool"
	if tool.Name() != expected {
		t.Errorf("Name() = %q, want %q", tool.Name(), expected)
	}
}

func TestMCPOperatorTool_Description(t *testing.T) {
	info := mcp.ToolInfo{
		Name:        "test_tool",
		Description: "Test description",
		Schema:      map[string]any{},
	}
	caller := &mockMCPCaller{}
	tool := NewMCPOperatorTool("test_server", info, caller, 30*time.Second)

	if tool.Description() != "Test description" {
		t.Errorf("Description() = %q, want %q", tool.Description(), "Test description")
	}
}

func TestMCPOperatorTool_Parameters(t *testing.T) {
	params := map[string]any{"key": "value"}
	info := mcp.ToolInfo{
		Name:        "test_tool",
		Description: "Test",
		Schema:      params,
	}
	caller := &mockMCPCaller{}
	tool := NewMCPOperatorTool("test_server", info, caller, 30*time.Second)

	if tool.Parameters()["key"] != "value" {
		t.Errorf("Parameters() = %v, want %v", tool.Parameters(), params)
	}
}

func TestMCPOperatorTool_Execute_Success(t *testing.T) {
	info := mcp.ToolInfo{
		Name:        "test_tool",
		Description: "Test",
		Schema:      map[string]any{},
	}
	caller := &mockMCPCaller{
		result: &mcp.ToolCallResult{
			Content: []mcp.ContentBlock{
				{Type: "text", Text: "Success result"},
			},
			IsError: false,
		},
	}
	tool := NewMCPOperatorTool("test_server", info, caller, 30*time.Second)

	result := tool.Execute(context.Background(), nil)
	if result.IsError {
		t.Error("expected success result, got error")
	}
	if result.ForLLM != "Success result" {
		t.Errorf("ForLLM = %q, want %q", result.ForLLM, "Success result")
	}
}
