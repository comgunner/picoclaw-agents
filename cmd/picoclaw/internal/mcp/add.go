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
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"

	"github.com/comgunner/picoclaw/cmd/picoclaw/internal"
	"github.com/comgunner/picoclaw/pkg/config"
)

func newAddCommand() *cobra.Command {
	var (
		transport    string
		command      string
		cmdArgs      []string
		envVars      []string
		enabledTools []string
		description  string
		url          string
		headers      []string
		timeout      int
	)

	cmd := &cobra.Command{
		Use:   "add <name>",
		Short: "Add a new MCP server",
		Long: `Add a new MCP server configuration to config.json.

Examples:
  # Add a stdio server (using npx)
  picoclaw-agents mcp add github --transport stdio --command npx \
    --args "-y,@modelcontextprotocol/server-github" \
    --env "GITHUB_PERSONAL_ACCESS_TOKEN=ghp_xxx" \
    --description "GitHub API"

  # Add an HTTP server
  picoclaw-agents mcp add myserver --transport http --url http://localhost:3000`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			return addServer(
				name,
				transport,
				command,
				cmdArgs,
				envVars,
				enabledTools,
				description,
				url,
				headers,
				timeout,
			)
		},
	}

	cmd.Flags().StringVar(&transport, "transport", "stdio", "Transport type (stdio, http, sse)")
	cmd.Flags().StringVar(&command, "command", "", "Command to execute (for stdio transport)")
	cmd.Flags().StringSliceVar(&cmdArgs, "args", nil, "Command arguments (comma-separated)")
	cmd.Flags().StringSliceVar(&envVars, "env", nil, "Environment variables (KEY=VALUE, comma-separated)")
	cmd.Flags().StringSliceVar(&enabledTools, "tools", []string{"*"}, "Enabled tools (* for all)")
	cmd.Flags().StringVar(&description, "description", "", "Server description")
	cmd.Flags().StringVar(&url, "url", "", "Server URL (for http/sse transport)")
	cmd.Flags().StringSliceVar(&headers, "headers", nil, "HTTP headers (KEY=VALUE, comma-separated)")
	cmd.Flags().IntVar(&timeout, "timeout", 0, "Request timeout in seconds (0 = use default)")

	return cmd
}

func addServer(
	name, transport, cmd_str string,
	cmdArgs, envVars, enabledTools []string,
	description, url string,
	headers []string,
	timeout int,
) error {
	if name == "" {
		return fmt.Errorf("server name is required")
	}

	if transport == "" {
		transport = "stdio"
	}

	// Validate transport
	switch transport {
	case "stdio", "http", "sse":
	default:
		return fmt.Errorf("invalid transport %q (allowed: stdio, http, sse)", transport)
	}

	// Parse env vars
	envMap := make(map[string]string)
	for _, e := range envVars {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid env variable format: %q (expected KEY=VALUE)", e)
		}
		envMap[parts[0]] = parts[1]
	}

	// Parse headers
	headerMap := make(map[string]string)
	for _, h := range headers {
		parts := strings.SplitN(h, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid header format: %q (expected KEY=VALUE)", h)
		}
		headerMap[parts[0]] = parts[1]
	}

	// Validate command whitelist for stdio transport
	if transport == "stdio" {
		srvCfg := config.MCPServerConfig{
			Transport: transport,
			Command:   cmd_str,
		}
		if err := srvCfg.ValidateCommand(); err != nil {
			return err
		}
	}

	// Load config
	cfgPath := internal.GetConfigPath()
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	// Check for duplicate
	if cfg.Tools.MCP.Servers != nil {
		if _, exists := cfg.Tools.MCP.Servers[name]; exists {
			return fmt.Errorf("MCP server %q already exists (use 'mcp remove %s' first)", name, name)
		}
	}

	// Initialize servers map if nil
	if cfg.Tools.MCP.Servers == nil {
		cfg.Tools.MCP.Servers = make(map[string]config.MCPServerConfig)
	}

	// Resolve command to absolute path if stdio
	resolvedCmd := cmd_str
	if transport == "stdio" && cmd_str != "" {
		if absPath, err := exec.LookPath(cmd_str); err == nil {
			resolvedCmd = absPath
		}
	}

	// Add server config
	cfg.Tools.MCP.Servers[name] = config.MCPServerConfig{
		Transport:    transport,
		Command:      resolvedCmd,
		Args:         cmdArgs,
		Env:          envMap,
		URL:          url,
		Headers:      headerMap,
		EnabledTools: enabledTools,
		Timeout:      timeout,
		Description:  description,
	}

	// Enable MCP if not already
	cfg.Tools.MCP.Enabled = true

	// Save config
	if err := config.SaveConfig(cfgPath, cfg); err != nil {
		return fmt.Errorf("save config: %w", err)
	}

	fmt.Fprintf(os.Stdout, "✅ MCP server %q added successfully\n", name)
	fmt.Fprintf(os.Stdout, "   Transport: %s\n", transport)
	if cmd_str != "" {
		fmt.Fprintf(os.Stdout, "   Command:   %s\n", resolvedCmd)
	}
	if len(cmdArgs) > 0 {
		fmt.Fprintf(os.Stdout, "   Args:      %s\n", strings.Join(cmdArgs, " "))
	}
	if description != "" {
		fmt.Fprintf(os.Stdout, "   Description: %s\n", description)
	}

	return nil
}
