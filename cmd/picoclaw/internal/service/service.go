// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package service

import (
	"github.com/spf13/cobra"
)

// NewServiceCommand creates the root "service" command.
func NewServiceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "service",
		Short: "Manage PicoClaw-Agents as an OS service (systemd/launchd/schtasks)",
		Long: `Install, uninstall, start, stop, and check the status of PicoClaw-Agents
as a system service. Supports Linux (systemd), macOS (launchd), and Windows (schtasks).

Examples:
  picoclaw-agents service install              Install as OS service
  picoclaw-agents service install --dry-run    Preview what would be installed
  picoclaw-agents service install --port 9999  Install on custom port
  picoclaw-agents service status               Check service status
  picoclaw-agents service logs -f              Follow service logs
  picoclaw-agents service uninstall            Remove the service
`,
	}

	cmd.AddCommand(newInstallCommand())
	cmd.AddCommand(newUninstallCommand())
	cmd.AddCommand(newStartCommand())
	cmd.AddCommand(newStopCommand())
	cmd.AddCommand(newStatusCommand())
	cmd.AddCommand(newRestartCommand())
	cmd.AddCommand(newLogsCommand())

	return cmd
}
