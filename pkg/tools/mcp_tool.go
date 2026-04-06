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
	"fmt"
	"time"

	"github.com/comgunner/picoclaw/pkg/mcp"
)

// MCPOperatorTool wraps a remote MCP tool as a local Tool interface implementation.
// It lives in pkg/tools to avoid circular dependency: pkg/mcp → pkg/tools → pkg/mcp.
type MCPOperatorTool struct {
	serverName  string
	toolName    string
	displayName string
	description string
	parameters  map[string]any
	client      MCPCaller
	timeout     time.Duration
}

// MCPCaller is the interface MCPClientManager exposes for tool calls.
type MCPCaller interface {
	CallTool(ctx context.Context, serverName, toolName string, args map[string]any) (*mcp.ToolCallResult, error)
}

// NewMCPOperatorTool creates a new MCP operator tool from a tool info struct.
func NewMCPOperatorTool(serverName string, info mcp.ToolInfo, client MCPCaller, defaultTimeout time.Duration) *MCPOperatorTool {
	return &MCPOperatorTool{
		serverName:  serverName,
		toolName:    info.Name,
		displayName: mcp.MCPToolName(serverName, info.Name),
		description: info.Description,
		parameters:  info.Schema,
		client:      client,
		timeout:     defaultTimeout,
	}
}

// Name returns the display name of the tool.
func (t *MCPOperatorTool) Name() string {
	return t.displayName
}

// Description returns the tool description.
func (t *MCPOperatorTool) Description() string {
	return t.description
}

// Parameters returns the tool's parameter schema.
func (t *MCPOperatorTool) Parameters() map[string]any {
	return t.parameters
}

// Execute calls the remote MCP tool and returns the result.
func (t *MCPOperatorTool) Execute(ctx context.Context, args map[string]any) *ToolResult {
	ctx, cancel := context.WithTimeout(ctx, t.timeout)
	defer cancel()

	result, err := t.client.CallTool(ctx, t.serverName, t.toolName, args)
	if err != nil {
		return ErrorResult(fmt.Sprintf("MCP tool %s/%s error: %v", t.serverName, t.toolName, err))
	}

	content := mcp.FormatContent(result.Content)
	if result.IsError {
		return ErrorResult(content)
	}
	return &ToolResult{ForLLM: content}
}
