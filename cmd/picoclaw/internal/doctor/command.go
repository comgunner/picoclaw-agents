// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package doctor

import (
	"encoding/json"
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/comgunner/picoclaw/pkg/setup"
)

// NewDoctorCommand creates the doctor command.
func NewDoctorCommand() *cobra.Command {
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Check environment and configuration health",
		Long:  "Diagnose environment issues: Go version, Docker, workspace, WSL, security.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runDoctor(cmd, jsonOutput)
		},
	}

	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")
	return cmd
}

func runDoctor(cmd *cobra.Command, jsonOutput bool) error {
	report := setup.CheckEnvironment()

	if jsonOutput {
		return outputJSON(cmd, report)
	}

	return outputTabular(cmd, report)
}

func outputJSON(cmd *cobra.Command, report *setup.EnvironmentReport) error {
	encoder := json.NewEncoder(cmd.OutOrStdout())
	encoder.SetIndent("", "  ")
	return encoder.Encode(report)
}

func outputTabular(cmd *cobra.Command, report *setup.EnvironmentReport) error {
	w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)

	fmt.Fprintln(w, "=== PicoClaw-Agents Doctor ===")
	fmt.Fprintln(w, "")

	fmt.Fprintln(w, "System:")
	fmt.Fprintf(w, "  OS/Arch:\t%s/%s\n", report.OS, report.Arch)
	fmt.Fprintf(w, "  WSL:\t%v\n", report.WSL)
	fmt.Fprintf(w, "  Shell:\t%s\n", report.Shell)
	fmt.Fprintln(w, "")

	fmt.Fprintln(w, "Requirements:")
	goStatus := statusIcon(report.GoOK)
	fmt.Fprintf(w, "  Go:\t%s %s\n", report.GoVersion, goStatus)

	dockerStatus := "not installed"
	if report.DockerInstalled {
		if report.DockerRunning {
			dockerStatus = "running"
		} else {
			dockerStatus = "installed (not running)"
		}
	}
	fmt.Fprintf(w, "  Docker:\t%s\n", dockerStatus)

	workspaceStatus := statusIcon(report.WorkspaceOK)
	fmt.Fprintf(w, "  Workspace:\t%s %s\n", report.WorkspacePath, workspaceStatus)

	configStatus := "not found"
	if report.ExistingConfig {
		configStatus = "exists"
	}
	fmt.Fprintf(w, "  Config:\t%s\n", configStatus)
	fmt.Fprintln(w, "")

	// Security section
	fmt.Fprintln(w, "Security:")
	rootStatus := "false (good)"
	if report.SecurityChecks.RunningAsRoot {
		rootStatus = "true (WARNING: running as root)"
	}
	fmt.Fprintf(w, "  Root:\t%s\n", rootStatus)

	if len(report.SecurityChecks.DangerousBinaries) > 0 {
		fmt.Fprintf(w, "  Dangerous:\t%s\n", joinStrings(report.SecurityChecks.DangerousBinaries, ", "))
	} else {
		fmt.Fprintln(w, "  Dangerous:\tnone")
	}

	if len(report.SecurityChecks.OpenPorts) > 0 {
		ports := make([]string, len(report.SecurityChecks.OpenPorts))
		for i, port := range report.SecurityChecks.OpenPorts {
			ports[i] = fmt.Sprintf("%d", port)
		}
		fmt.Fprintf(w, "  Open ports:\t%s\n", joinStrings(ports, ", "))
	} else {
		fmt.Fprintln(w, "  Open ports:\tnone")
	}
	fmt.Fprintln(w, "")

	// Final status
	if report.IsReady() {
		fmt.Fprintln(w, "✓ Environment ready!")
	} else {
		fmt.Fprintln(w, "✗ Environment has issues")
		if !report.GoOK {
			fmt.Fprintln(w, "  - Go is not installed or version < 1.21")
		}
		if !report.WorkspaceOK {
			fmt.Fprintln(w, "  - Workspace directory is not accessible")
		}
	}

	return w.Flush()
}

func statusIcon(ok bool) string {
	if ok {
		return "[OK]"
	}
	return "[MISSING]"
}

func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}
