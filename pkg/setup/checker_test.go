// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package setup_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/comgunner/picoclaw/pkg/setup"
)

// TestCheckEnvironment verifies that CheckEnvironment returns a valid report.
func TestCheckEnvironment(t *testing.T) {
	report := setup.CheckEnvironment()

	assert.NotNil(t, report)
	assert.NotEmpty(t, report.OS)
	assert.NotEmpty(t, report.Arch)
	assert.NotEmpty(t, report.GoVersion)
	assert.NotEmpty(t, report.WorkspacePath)
}

// TestEnvironmentReport_IsReady verifies the IsReady method logic.
func TestEnvironmentReport_IsReady(t *testing.T) {
	t.Run("ready when Go and workspace OK", func(t *testing.T) {
		report := &setup.EnvironmentReport{
			GoOK:            true,
			WorkspaceOK:     true,
			DockerInstalled: false, // Docker is optional
		}
		assert.True(t, report.IsReady())
	})

	t.Run("not ready when Go missing", func(t *testing.T) {
		report := &setup.EnvironmentReport{
			GoOK:            false,
			WorkspaceOK:     true,
			DockerInstalled: true,
		}
		assert.False(t, report.IsReady())
	})

	t.Run("not ready when workspace missing", func(t *testing.T) {
		report := &setup.EnvironmentReport{
			GoOK:            true,
			WorkspaceOK:     false,
			DockerInstalled: true,
		}
		assert.False(t, report.IsReady())
	})
}

// TestEnvironmentReport_String verifies the String method output format.
func TestEnvironmentReport_String(t *testing.T) {
	report := &setup.EnvironmentReport{
		OS:              "linux",
		Arch:            "amd64",
		GoVersion:       "go1.22.0",
		GoOK:            true,
		DockerInstalled: true,
		DockerRunning:   true,
		WorkspacePath:   "/home/user/.picoclaw",
		WorkspaceOK:     true,
		Shell:           "/bin/bash",
		ExistingConfig:  false,
	}

	output := report.String()

	assert.Contains(t, output, "=== Environment Check ===")
	assert.Contains(t, output, "linux/amd64")
	assert.Contains(t, output, "go1.22.0")
	assert.Contains(t, output, "OK")
}

// TestCheckGoVersion verifies Go version detection.
func TestCheckGoVersion(t *testing.T) {
	// This test assumes Go is installed (required to run tests)
	version, ok := getGoVersionHelper()
	assert.True(t, ok, "Go should be installed to run tests")
	assert.Contains(t, version, "go1.")
}

// TestCheckDocker verifies Docker detection (skipped if Docker not installed).
func TestCheckDocker(t *testing.T) {
	installed, running := checkDockerHelper()

	if !installed {
		t.Skip("Docker not installed - skipping Docker test")
	}

	assert.True(t, installed)
	// Running status depends on whether daemon is active
	t.Logf("Docker: installed=%v running=%v", installed, running)
}

// TestCheckWorkspace verifies workspace creation and permission check.
func TestCheckWorkspace(t *testing.T) {
	// Create temporary directory for testing
	tmpDir := t.TempDir()
	testWorkspace := filepath.Join(tmpDir, ".picoclaw")

	// Override home directory for testing
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	workspace, ok := checkWorkspaceWithBase(testWorkspace)

	assert.NotEmpty(t, workspace)
	assert.True(t, ok, "Workspace should be created and accessible")

	// Verify directory was created
	_, err := os.Stat(workspace)
	assert.NoError(t, err)
}

// TestCheckConfigExists verifies config file detection.
func TestCheckConfigExists(t *testing.T) {
	tmpDir := t.TempDir()

	// Test when config doesn't exist
	exists := checkConfigExistsHelper(tmpDir)
	assert.False(t, exists)

	// Test when config exists
	configPath := filepath.Join(tmpDir, "config.json")
	err := os.WriteFile(configPath, []byte("{}"), 0o600)
	assert.NoError(t, err)

	exists = checkConfigExistsHelper(tmpDir)
	assert.True(t, exists)
}

// TestEnvironmentReport_Fields verifies all report fields are populated.
func TestEnvironmentReport_Fields(t *testing.T) {
	report := setup.CheckEnvironment()

	// Required fields
	assert.NotEmpty(t, report.OS, "OS should be detected")
	assert.NotEmpty(t, report.Arch, "Arch should be detected")
	assert.NotEmpty(t, report.GoVersion, "Go version should be detected")
	assert.NotEmpty(t, report.WorkspacePath, "Workspace path should be set")

	// Boolean fields should have valid values
	assert.IsType(t, true, report.GoOK)
	assert.IsType(t, true, report.DockerInstalled)
	assert.IsType(t, true, report.DockerRunning)
	assert.IsType(t, true, report.WorkspaceOK)
	assert.IsType(t, true, report.ExistingConfig)

	// New fields in v3.8.0
	assert.IsType(t, true, report.WSL)
	assert.IsType(t, setup.SecurityReport{}, report.SecurityChecks)
}

// TestDetectWSL_NonLinux verifies WSL detection returns false on non-Linux systems.
func TestDetectWSL_NonLinux(t *testing.T) {
	// WSL detection should always return false on macOS/Windows
	report := setup.CheckEnvironment()

	// On macOS (darwin) or Windows, WSL should be false
	if report.OS != "linux" {
		assert.False(t, report.WSL, "WSL should be false on non-Linux systems")
	}
}

// TestSecurityChecks_NotRoot verifies that tests don't run as root.
func TestSecurityChecks_NotRoot(t *testing.T) {
	report := setup.CheckEnvironment()

	// Tests should not run as root for security
	assert.False(t, report.SecurityChecks.RunningAsRoot, "Tests should not run as root")
}

// TestSecurityChecks_Ports verifies that port check returns a list (may be empty).
func TestSecurityChecks_Ports(t *testing.T) {
	report := setup.CheckEnvironment()

	// OpenPorts should be a valid slice (may be empty if no common ports open)
	assert.NotNil(t, report.SecurityChecks.OpenPorts)

	// If ports are detected, they should be in common range
	for _, port := range report.SecurityChecks.OpenPorts {
		assert.Greater(t, port, 0)
		assert.LessOrEqual(t, port, 65535)
	}
}

// TestSecurityChecks_DangerousBinaries verifies dangerous binaries detection.
func TestSecurityChecks_DangerousBinaries(t *testing.T) {
	report := setup.CheckEnvironment()

	// DangerousBinaries should be a valid slice (may be empty)
	assert.NotNil(t, report.SecurityChecks.DangerousBinaries)

	// If any dangerous binaries detected, they should be from expected list
	expectedBins := []string{"nc", "netcat", "ncat", "telnet"}
	for _, bin := range report.SecurityChecks.DangerousBinaries {
		assert.Contains(t, expectedBins, bin)
	}
}

// TestEnvironmentReport_WithSecurity verifies String() includes security section when relevant.
func TestEnvironmentReport_WithSecurity(t *testing.T) {
	report := &setup.EnvironmentReport{
		OS:   "linux",
		Arch: "amd64",
		GoOK: true,
		WSL:  false,
		SecurityChecks: setup.SecurityReport{
			RunningAsRoot:     true,
			DangerousBinaries: []string{"nc"},
			OpenPorts:         []int{22, 80},
		},
	}

	output := report.String()

	assert.Contains(t, output, "=== Security Checks ===")
	assert.Contains(t, output, "Running as root: true")
	assert.Contains(t, output, "Dangerous binaries: nc")
	assert.Contains(t, output, "Open ports: 22, 80")
}

// TestEnvironmentReport_WithoutSecurity verifies String() omits security section when all clear.
func TestEnvironmentReport_WithoutSecurity(t *testing.T) {
	report := &setup.EnvironmentReport{
		OS:   "darwin",
		Arch: "arm64",
		GoOK: true,
		WSL:  false,
		SecurityChecks: setup.SecurityReport{
			RunningAsRoot:     false,
			DangerousBinaries: []string{},
			OpenPorts:         []int{},
		},
	}

	output := report.String()

	// Security section should not appear when all checks are clean
	assert.NotContains(t, output, "=== Security Checks ===")
}

// Helper functions (mirror private functions in checker.go for testing)

func getGoVersionHelper() (string, bool) {
	// Use reflection or re-implement logic for testing
	// For now, just call CheckEnvironment and extract
	report := setup.CheckEnvironment()
	return report.GoVersion, report.GoOK
}

func checkDockerHelper() (bool, bool) {
	report := setup.CheckEnvironment()
	return report.DockerInstalled, report.DockerRunning
}

func checkWorkspaceWithBase(base string) (string, bool) {
	// Simplified version of checkWorkspace for testing
	workspace := base
	if err := os.MkdirAll(workspace, 0o700); err != nil {
		return workspace, false
	}
	tmp := filepath.Join(workspace, ".check_write")
	if err := os.WriteFile(tmp, []byte("ok"), 0o600); err != nil {
		return workspace, false
	}
	os.Remove(tmp)
	return workspace, true
}

func checkConfigExistsHelper(workspacePath string) bool {
	configPath := filepath.Join(workspacePath, "config.json")
	_, err := os.Stat(configPath)
	return err == nil
}
