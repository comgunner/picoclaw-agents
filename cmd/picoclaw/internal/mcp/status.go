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
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/comgunner/picoclaw/cmd/picoclaw/internal"
	"github.com/comgunner/picoclaw/pkg/config"
	"github.com/comgunner/picoclaw/pkg/mcp"
)

func newStatusCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status <name>",
		Short: "Check MCP server connection status",
		Long:  "Attempts to connect to the specified MCP server and reports its status.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return checkStatus(args[0])
		},
	}

	return cmd
}

func checkStatus(name string) error {
	cfgPath := internal.GetConfigPath()
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	servers := cfg.Tools.MCP.Servers
	if servers == nil {
		fmt.Fprintf(os.Stdout, "❌ MCP server %q not found (no servers configured)\n", name)
		return nil
	}

	srv, exists := servers[name]
	if !exists {
		fmt.Fprintf(os.Stdout, "❌ MCP server %q not found\n", name)
		fmt.Fprintf(os.Stdout, "Use 'picoclaw-agents mcp list' to see configured servers.\n")
		return nil
	}

	// Attempt connection
	timeout := time.Duration(srv.Timeout) * time.Second
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	transport, err := mcp.NewTransport(srv)
	if err != nil {
		fmt.Fprintf(os.Stdout, "❌ MCP server %q — error: %v\n", name, err)
		return nil
	}
	defer transport.Close()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_, err = transport.Call(ctx, "initialize", map[string]any{
		"protocolVersion": mcp.MCPProtocolVersion,
		"clientInfo":      map[string]string{"name": "picoclaw-agents", "version": "dev"},
		"capabilities":    map[string]any{},
	})
	if err != nil {
		fmt.Fprintf(os.Stdout, "❌ MCP server %q — connection failed: %v\n", name, err)
		return nil
	}

	fmt.Fprintf(os.Stdout, "✅ MCP server %q — connected\n", name)
	fmt.Fprintf(os.Stdout, "   Transport:   %s\n", srv.Transport)
	if srv.Command != "" {
		fmt.Fprintf(os.Stdout, "   Command:     %s\n", srv.Command)
	}
	if srv.URL != "" {
		fmt.Fprintf(os.Stdout, "   URL:         %s\n", srv.URL)
	}
	if srv.Description != "" {
		fmt.Fprintf(os.Stdout, "   Description: %s\n", srv.Description)
	}

	return nil
}
