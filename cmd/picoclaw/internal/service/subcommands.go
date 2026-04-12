// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT

package service

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newStartCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start the PicoClaw-Agents service",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := ServiceConfig{}
			platform, err := GetServicePlatform(cfg)
			if err != nil {
				return err
			}
			return platform.Start()
		},
	}
	return cmd
}

func newStopCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop the PicoClaw-Agents service",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := ServiceConfig{}
			platform, err := GetServicePlatform(cfg)
			if err != nil {
				return err
			}
			return platform.Stop()
		},
	}
	return cmd
}

func newRestartCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "restart",
		Short: "Restart the PicoClaw-Agents service",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := ServiceConfig{}
			platform, err := GetServicePlatform(cfg)
			if err != nil {
				return err
			}

			fmt.Println("🔄 Restarting service...")
			// Use platform-specific restart via type assertion
			switch p := platform.(type) {
			case interface{ restart() error }:
				return p.restart()
			default:
				if err := platform.Stop(); err != nil {
					return err
				}
				return platform.Start()
			}
		},
	}
	return cmd
}

func newStatusCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Check the PicoClaw-Agents service status",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := ServiceConfig{}
			platform, err := GetServicePlatform(cfg)
			if err != nil {
				return err
			}

			status, err := platform.Status()
			if err != nil {
				return fmt.Errorf("failed to get status: %w", err)
			}

			fmt.Println("📊 PicoClaw-Agents Service Status")
			fmt.Println("====================================")

			if status.Installed {
				fmt.Println("  Installed:   Yes")
			} else {
				fmt.Println("  Installed:   No")
				fmt.Println()
				fmt.Println("  The service is not installed.")
				fmt.Println("  Run 'picoclaw-agents service install' to install it.")
				return nil
			}

			if status.Active {
				fmt.Println("  Active:      Yes")
			} else {
				fmt.Println("  Active:      No (installed but stopped)")
			}

			if status.PID > 0 {
				fmt.Printf("  PID:         %d\n", status.PID)
			}
			fmt.Printf("  Port:        %d\n", status.Port)

			if status.Error != "" {
				fmt.Printf("  Error:       %s\n", status.Error)
			}

			return nil
		},
	}
	return cmd
}

func newLogsCommand() *cobra.Command {
	var (
		lines  int
		follow bool
	)

	cmd := &cobra.Command{
		Use:   "logs",
		Short: "View PicoClaw-Agents service logs",
		Long: `View service logs. Use -f to follow logs in real-time.

Examples:
  picoclaw-agents service logs          Show last 50 lines
  picoclaw-agents service logs -n 100   Show last 100 lines
  picoclaw-agents service logs -f       Follow logs in real-time
  picoclaw-agents service logs -n 10 -f Show last 10 lines and follow
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := ServiceConfig{}
			platform, err := GetServicePlatform(cfg)
			if err != nil {
				return err
			}
			return platform.Logs(lines, follow)
		},
	}

	cmd.Flags().IntVarP(&lines, "lines", "n", 50, "Number of lines to show")
	cmd.Flags().BoolVarP(&follow, "follow", "f", false, "Follow log output")

	return cmd
}
