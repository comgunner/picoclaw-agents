// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package mcp

import "github.com/spf13/cobra"

// NewMCPCommand creates the root mcp command.
func NewMCPCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mcp",
		Short: "Manage MCP servers (add, list, status, remove)",
		Long:  "Manage Model Context Protocol (MCP) server configurations: add, list, check status, and remove servers.",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(
		newAddCommand(),
		newListCommand(),
		newStatusCommand(),
		newRemoveCommand(),
	)

	return cmd
}
