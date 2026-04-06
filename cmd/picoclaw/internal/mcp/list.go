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
	"text/tabwriter"

	"github.com/comgunner/picoclaw/cmd/picoclaw/internal"
	"github.com/comgunner/picoclaw/pkg/config"
	"github.com/spf13/cobra"
)

func newListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List configured MCP servers",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return listServers()
		},
	}

	return cmd
}

func listServers() error {
	cfgPath := internal.GetConfigPath()
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	servers := cfg.Tools.MCP.Servers
	if len(servers) == 0 {
		fmt.Fprintln(os.Stdout, "No MCP servers configured.")
		fmt.Fprintln(os.Stdout, "Use 'picoclaw-agents mcp add <name>' to add a server.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tTRANSPORT\tCOMMAND/URL\tTOOLS\tDESCRIPTION")
	fmt.Fprintln(w, "----\t---------\t-----------\t-----\t-----------")

	for name, srv := range servers {
		commandOrURL := srv.Command
		if srv.URL != "" {
			commandOrURL = srv.URL
		}
		if commandOrURL == "" {
			commandOrURL = "-"
		}

		toolCount := "-"
		if len(srv.EnabledTools) > 0 {
			if len(srv.EnabledTools) == 1 && srv.EnabledTools[0] == "*" {
				toolCount = "all"
			} else {
				toolCount = fmt.Sprintf("%d", len(srv.EnabledTools))
			}
		}

		desc := srv.Description
		if desc == "" {
			desc = "-"
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", name, srv.Transport, commandOrURL, toolCount, desc)
	}

	w.Flush()

	if !cfg.Tools.MCP.Enabled {
		fmt.Fprintln(os.Stdout, "\n⚠️  MCP is disabled in config. Set tools.mcp.enabled to true to enable.")
	}

	return nil
}
