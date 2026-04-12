// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package service

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newUninstallCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uninstall",
		Short: "Remove the PicoClaw-Agents OS service",
		Long: `Stop and remove the PicoClaw-Agents service from the OS.

On Linux: removes the systemd user unit
On macOS: removes the launchd LaunchAgent
On Windows: removes the scheduled task and wrapper script
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := ServiceConfig{}
			platform, err := GetServicePlatform(cfg)
			if err != nil {
				return err
			}

			fmt.Println("🗑️  Uninstalling PicoClaw-Agents service...")
			fmt.Println()
			return platform.Uninstall()
		},
	}

	return cmd
}
