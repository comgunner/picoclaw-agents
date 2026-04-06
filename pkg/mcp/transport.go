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
	"context"
	"encoding/json"
	"fmt"

	"github.com/comgunner/picoclaw/pkg/config"
)

// MCPTransport abstracts communication with an MCP server.
type MCPTransport interface {
	// Call sends a JSON-RPC request and returns the response.
	// Must respect context cancellation.
	Call(ctx context.Context, method string, params map[string]any) (*json.RawMessage, error)

	// Close terminates the transport and releases resources.
	Close() error
}

// NewTransport creates a transport based on the server config.
func NewTransport(cfg config.MCPServerConfig) (MCPTransport, error) {
	switch cfg.Transport {
	case "stdio":
		return NewStdioTransport(cfg.Command, cfg.Args, cfg.Env)
	case "sse":
		return NewSSETransport(cfg.URL, cfg.Headers)
	case "http":
		return NewHTTPTransport(cfg.URL, cfg.Headers)
	default:
		return nil, fmt.Errorf("unsupported MCP transport: %q", cfg.Transport)
	}
}
