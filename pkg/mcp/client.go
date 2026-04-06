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
	"sync"
	"time"

	"github.com/comgunner/picoclaw/pkg/config"
	"github.com/comgunner/picoclaw/pkg/logger"
)

// MCPClientManager manages connections to multiple external MCP servers.
type MCPClientManager struct {
	servers map[string]*mcpServer
	mu      sync.RWMutex
}

// mcpServer represents a connected MCP server with its transport and tools.
type mcpServer struct {
	cfg       config.MCPServerConfig
	transport MCPTransport
	tools     []ToolInfo
	status    ServerStatus
}

// ServerStatus represents the connection status of an MCP server.
type ServerStatus string

const (
	// StatusConnected indicates the server is connected and ready.
	StatusConnected ServerStatus = "connected"
	// StatusConnecting indicates the server is in the process of connecting.
	StatusConnecting ServerStatus = "connecting"
	// StatusError indicates the server failed to connect.
	StatusError ServerStatus = "error"
)

// ToolInfo describes a tool exposed by an MCP server.
type ToolInfo struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Schema      map[string]any `json:"inputSchema"`
}

// ConnectAll connects to all configured MCP servers.
// Returns errors for servers that failed — non-fatal, agent continues.
func (m *MCPClientManager) ConnectAll(servers map[string]config.MCPServerConfig) []error {
	m.servers = make(map[string]*mcpServer)
	var errs []error

	for name, cfg := range servers {
		if err := m.connectServer(name, cfg); err != nil {
			errs = append(errs, fmt.Errorf("MCP server %q: %w", name, err))
			logger.WarnCF("mcp", "MCP server failed to connect (continuing with others)",
				map[string]any{"server": name, "error": err.Error()})
		}
	}
	return errs
}

func (m *MCPClientManager) connectServer(name string, cfg config.MCPServerConfig) error {
	transport, err := NewTransport(cfg)
	if err != nil {
		return fmt.Errorf("create transport: %w", err)
	}

	// Initialize handshake
	_, err = transport.Call(context.Background(), "initialize", map[string]any{
		"protocolVersion": MCPProtocolVersion,
		"clientInfo":      map[string]string{"name": "picoclaw-agents", "version": "1.2.5"},
		"capabilities":    map[string]any{},
	})
	if err != nil {
		transport.Close()
		return fmt.Errorf("initialize: %w", err)
	}

	// Notify initialized
	_, _ = transport.Call(context.Background(), "notifications/initialized", nil)

	// Fetch tool list
	result, err := transport.Call(context.Background(), "tools/list", nil)
	if err != nil {
		transport.Close()
		return fmt.Errorf("tools/list: %w", err)
	}

	var toolsResp struct {
		Tools []ToolInfo `json:"tools"`
	}
	if err := json.Unmarshal(*result, &toolsResp); err != nil {
		transport.Close()
		return fmt.Errorf("parse tools: %w", err)
	}

	// Filter tools if enabled_tools is configured
	tools := filterTools(toolsResp.Tools, cfg.EnabledTools)

	m.mu.Lock()
	m.servers[name] = &mcpServer{
		cfg:       cfg,
		transport: transport,
		tools:     tools,
		status:    StatusConnected,
	}
	m.mu.Unlock()

	logger.InfoCF("mcp", "MCP server connected with tools",
		map[string]any{"server": name, "tools": len(tools)})
	return nil
}

// GetServer returns a server by name, or nil if not found.
func (m *MCPClientManager) GetServer(name string) *mcpServer {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.servers[name]
}

// CallTool calls a tool on a specific MCP server.
func (m *MCPClientManager) CallTool(ctx context.Context, serverName, toolName string, args map[string]any) (*ToolCallResult, error) {
	srv := m.GetServer(serverName)
	if srv == nil {
		return nil, fmt.Errorf("MCP server %q not found", serverName)
	}

	timeoutSec := srv.cfg.Timeout
	if timeoutSec == 0 {
		timeoutSec = 30 // default 30 seconds
	}
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSec)*time.Second)
	defer cancel()

	result, err := srv.transport.Call(ctx, "tools/call", map[string]any{
		"name":      toolName,
		"arguments": args,
	})
	if err != nil {
		return nil, fmt.Errorf("call %s/%s: %w", serverName, toolName, err)
	}

	return parseToolCallResult(*result)
}

// ListServers returns a map of server name to tool count.
func (m *MCPClientManager) ListServers() map[string]int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make(map[string]int)
	for name, srv := range m.servers {
		out[name] = len(srv.tools)
	}
	return out
}

// CloseAll closes all MCP server connections.
func (m *MCPClientManager) CloseAll() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for name, srv := range m.servers {
		if err := srv.transport.Close(); err != nil {
			logger.WarnCF("mcp", "MCP server close error",
				map[string]any{"server": name, "error": err.Error()})
		}
	}
}

// Close implements io.Closer for MCPClientManager (delegates to CloseAll).
func (m *MCPClientManager) Close() error {
	m.CloseAll()
	return nil
}

// Servers returns a copy of the connected servers map (for inspection).
func (m *MCPClientManager) Servers() map[string]*mcpServer {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make(map[string]*mcpServer, len(m.servers))
	for name, srv := range m.servers {
		out[name] = srv
	}
	return out
}

// Tools returns the list of tools for a given server.
func (s *mcpServer) Tools() []ToolInfo {
	return s.tools
}

func filterTools(tools []ToolInfo, enabled []string) []ToolInfo {
	if len(enabled) == 0 || (len(enabled) == 1 && enabled[0] == "*") {
		return tools // all tools enabled
	}
	set := make(map[string]bool, len(enabled))
	for _, t := range enabled {
		set[t] = true
	}
	var out []ToolInfo
	for _, t := range tools {
		if set[t.Name] {
			out = append(out, t)
		}
	}
	return out
}
