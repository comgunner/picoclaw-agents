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
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/comgunner/picoclaw/cmd/picoclaw/internal"
	"github.com/comgunner/picoclaw/pkg/config"
)

func newRemoveCommand() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "remove <name>",
		Short: "Remove an MCP server",
		Long:  "Remove an MCP server configuration from config.json.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return removeServer(args[0], force)
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Skip confirmation prompt")

	return cmd
}

func removeServer(name string, force bool) error {
	cfgPath := internal.GetConfigPath()
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	servers := cfg.Tools.MCP.Servers
	if servers == nil || len(servers) == 0 {
		return fmt.Errorf("no MCP servers configured")
	}

	if _, exists := servers[name]; !exists {
		return fmt.Errorf("MCP server %q not found (use 'mcp list' to see configured servers)", name)
	}

	// Confirm removal unless --force is set
	if !force {
		fmt.Fprintf(os.Stdout, "Are you sure you want to remove MCP server %q? [y/N] ", name)
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))
		if input != "y" && input != "yes" {
			fmt.Fprintln(os.Stdout, "Aborted.")
			return nil
		}
	}

	delete(servers, name)

	// If no servers remain, disable MCP
	if len(servers) == 0 {
		cfg.Tools.MCP.Enabled = false
		cfg.Tools.MCP.Servers = nil
	}

	if err := config.SaveConfig(cfgPath, cfg); err != nil {
		return fmt.Errorf("save config: %w", err)
	}

	fmt.Fprintf(os.Stdout, "✅ MCP server %q removed successfully\n", name)
	return nil
}
