// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package service

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newInstallCommand() *cobra.Command {
	var (
		port   int
		public bool
		config string
		dryRun bool
	)

	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install PicoClaw-Agents as an OS service",
		Long: `Install the PicoClaw-Agents gateway as a system service that starts automatically.

On Linux: creates a systemd user unit (~/.config/systemd/user/picoclaw-agents.service)
On macOS: creates a launchd LaunchAgent (~/Library/LaunchAgents/com.picoclaw.agents.gateway.plist)
On Windows: creates a scheduled task ("PicoClaw-Agents Gateway")

The service will NOT run simultaneously with a manually started gateway to prevent
duplicate connections (Telegram, Discord, etc.).

Examples:
  picoclaw-agents service install              Install with default settings
  picoclaw-agents service install --dry-run    Preview without making changes
  picoclaw-agents service install --port 9999  Install on port 9999
  picoclaw-agents service install --public     Listen on all interfaces
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := ServiceConfig{
				Port:       port,
				Public:     public,
				ConfigPath: config,
				DryRun:     dryRun,
			}

			platform, err := GetServicePlatform(cfg)
			if err != nil {
				return err
			}

			if dryRun {
				fmt.Println("🔍 Dry Run Mode — No changes will be made")
				fmt.Println("==================================================")
				fmt.Println()
				return platform.DryRunInstall()
			}

			fmt.Println("📦 Installing PicoClaw-Agents as OS service...")
			fmt.Println()
			return platform.Install()
		},
	}

	cmd.Flags().IntVar(&port, "port", 18800, "Gateway port")
	cmd.Flags().BoolVar(&public, "public", false, "Listen on all interfaces (0.0.0.0)")
	cmd.Flags().StringVar(&config, "config", "", "Path to config.json (default: ~/.picoclaw/config.json)")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be done without making changes")

	return cmd
}
