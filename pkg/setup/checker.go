// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package setup

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// SecurityReport contains security check results.
type SecurityReport struct {
	RunningAsRoot     bool     // true if running as root (euid == 0)
	DangerousBinaries []string // dangerous binaries found in PATH (nc, netcat, ncat, telnet)
	OpenPorts         []int    // open ports detected (22, 80, 443, 8080, 8443)
}

// EnvironmentReport contains the environment check results.
// All fields are public for easy JSON serialization.
type EnvironmentReport struct {
	OS              string         // runtime.GOOS
	Arch            string         // runtime.GOARCH
	GoVersion       string         // output of `go version` or "unknown"
	GoOK            bool           // true if Go >= 1.21
	DockerInstalled bool           // true if docker binary exists in PATH
	DockerRunning   bool           // true if `docker info` returns exit 0
	WorkspacePath   string         // absolute path to workspace (~/.picoclaw by default)
	WorkspaceOK     bool           // true if directory exists and has rw permissions
	Shell           string         // value of $SHELL
	ExistingConfig  bool           // true if ~/.picoclaw/config.json already exists
	WSL             bool           // true if running in Windows Subsystem for Linux
	SecurityChecks  SecurityReport // security check results
}

// IsReady returns true if the environment meets minimum requirements to run picoclaw-agents.
// Minimum requirements: Go installed and workspace accessible.
// Docker is optional (required only if using --profile gateway/agent).
func (r *EnvironmentReport) IsReady() bool {
	return r.GoOK && r.WorkspaceOK
}

// String returns a tabular representation of the report for terminal display.
func (r *EnvironmentReport) String() string {
	check := func(ok bool) string {
		if ok {
			return "OK"
		}
		return "MISSING"
	}

	lines := []string{
		"=== Environment Check ===",
		fmt.Sprintf("OS/Arch:        %s/%s", r.OS, r.Arch),
		fmt.Sprintf("WSL:            %v", r.WSL),
		fmt.Sprintf("Shell:          %s", r.Shell),
		fmt.Sprintf("Go version:     %s [%s]", r.GoVersion, check(r.GoOK)),
		fmt.Sprintf("Docker:         installed=%v running=%v", r.DockerInstalled, r.DockerRunning),
		fmt.Sprintf("Workspace:      %s [%s]", r.WorkspacePath, check(r.WorkspaceOK)),
		fmt.Sprintf("Config exists:  %v", r.ExistingConfig),
		fmt.Sprintf("Ready:          %v", r.IsReady()),
	}

	// Add security section if relevant
	if r.SecurityChecks.RunningAsRoot || len(r.SecurityChecks.DangerousBinaries) > 0 ||
		len(r.SecurityChecks.OpenPorts) > 0 {
		lines = append(lines, "")
		lines = append(lines, "=== Security Checks ===")
		lines = append(lines, fmt.Sprintf("Running as root: %v", r.SecurityChecks.RunningAsRoot))

		if len(r.SecurityChecks.DangerousBinaries) > 0 {
			lines = append(
				lines,
				fmt.Sprintf("Dangerous binaries: %s", strings.Join(r.SecurityChecks.DangerousBinaries, ", ")),
			)
		} else {
			lines = append(lines, "Dangerous binaries: none")
		}

		if len(r.SecurityChecks.OpenPorts) > 0 {
			ports := make([]string, len(r.SecurityChecks.OpenPorts))
			for i, port := range r.SecurityChecks.OpenPorts {
				ports[i] = fmt.Sprintf("%d", port)
			}
			lines = append(lines, fmt.Sprintf("Open ports: %s", strings.Join(ports, ", ")))
		} else {
			lines = append(lines, "Open ports: none")
		}
	}

	return strings.Join(lines, "\n")
}

// CheckEnvironment executes all checks and returns a complete EnvironmentReport.
// Does not return error: on individual failures, marks the corresponding field as false.
func CheckEnvironment() *EnvironmentReport {
	report := &EnvironmentReport{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}

	report.GoVersion, report.GoOK = checkGoVersion()
	report.DockerInstalled, report.DockerRunning = checkDocker()
	report.WorkspacePath, report.WorkspaceOK = checkWorkspace()
	report.Shell = os.Getenv("SHELL")
	report.ExistingConfig = checkConfigExists(report.WorkspacePath)
	report.WSL = detectWSL()
	report.SecurityChecks = runSecurityChecks()

	return report
}

// checkGoVersion verifies that Go is installed and >= 1.21.
// Returns (version_string, ok).
func checkGoVersion() (string, bool) {
	out, err := exec.Command("go", "version").Output()
	if err != nil {
		return "not found", false
	}
	version := strings.TrimSpace(string(out))
	// Format: "go version go1.22.0 linux/amd64"
	// Basic check: contains "go1." followed by number >= 21
	for major := 21; major <= 99; major++ {
		if strings.Contains(version, fmt.Sprintf("go1.%d", major)) {
			return version, true
		}
	}
	return version, false
}

// checkDocker verifies if Docker is installed and running.
// Returns (installed, running).
func checkDocker() (bool, bool) {
	_, err := exec.LookPath("docker")
	if err != nil {
		return false, false
	}
	// Docker is installed; check if daemon responds
	err = exec.Command("docker", "info").Run()
	return true, err == nil
}

// checkWorkspace verifies that the workspace directory exists and is rw accessible.
// Uses ~/.picoclaw as default if HOME is available.
func checkWorkspace() (string, bool) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", false
	}
	workspace := filepath.Join(home, ".picoclaw")

	// Try to create if not exists
	if err := os.MkdirAll(workspace, 0o700); err != nil {
		return workspace, false
	}

	// Verify write permission with temporary file
	tmp := filepath.Join(workspace, ".check_write")
	if err := os.WriteFile(tmp, []byte("ok"), 0o600); err != nil {
		return workspace, false
	}
	os.Remove(tmp)

	return workspace, true
}

// checkConfigExists returns true if config.json exists in the workspace.
func checkConfigExists(workspacePath string) bool {
	configPath := filepath.Join(workspacePath, "config.json")
	_, err := os.Stat(configPath)
	return err == nil
}

// detectWSL returns true if running inside Windows Subsystem for Linux.
// Only checks on Linux; always returns false on macOS/Windows.
func detectWSL() bool {
	if runtime.GOOS != "linux" {
		return false
	}

	// Read /proc/version to check for Microsoft/WSL indicators
	data, err := os.ReadFile("/proc/version")
	if err != nil {
		return false
	}

	version := strings.ToLower(string(data))
	// WSL1 contains "Microsoft" and WSL2 contains "microsoft" or "WSL"
	return strings.Contains(version, "microsoft") || strings.Contains(version, "wsl")
}

// runSecurityChecks performs security checks and returns a SecurityReport.
func runSecurityChecks() SecurityReport {
	report := SecurityReport{
		RunningAsRoot:     os.Geteuid() == 0,
		DangerousBinaries: []string{},
		OpenPorts:         []int{},
	}

	// Check for dangerous binaries in PATH
	dangerousBins := []string{"nc", "netcat", "ncat", "telnet"}
	for _, bin := range dangerousBins {
		if _, err := exec.LookPath(bin); err == nil {
			report.DangerousBinaries = append(report.DangerousBinaries, bin)
		}
	}

	// Check for common open ports
	commonPorts := []int{22, 80, 443, 8080, 8443}
	for _, port := range commonPorts {
		if isPortOpen(port) {
			report.OpenPorts = append(report.OpenPorts, port)
		}
	}

	return report
}

// isPortOpen checks if a TCP port is open on localhost using dial with timeout.
// Returns true if connection succeeds (port is open), false otherwise.
func isPortOpen(port int) bool {
	addr := fmt.Sprintf(":%d", port)
	conn, err := net.DialTimeout("tcp", addr, 500*time.Millisecond)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}
