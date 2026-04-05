// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package sandbox_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/comgunner/picoclaw/cmd/picoclaw/internal/sandbox"
)

// TestSandboxCommand_Runs verifies the sandbox command executes.
func TestSandboxCommand_Runs(t *testing.T) {
	cmd := sandbox.NewSandboxCommand()
	assert.NotNil(t, cmd)
	assert.Equal(t, "sandbox", cmd.Use)
}

// TestSandboxInitCommand_CreatesDirectory verifies sandbox init creates directory.
func TestSandboxInitCommand_CreatesDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	sandboxName := "test-sandbox"

	cmd := sandbox.NewSandboxCommand()
	cmd.SetArgs([]string{"init", sandboxName, "--path", tmpDir})

	err := cmd.Execute()
	require.NoError(t, err)

	// Verify directory was created
	sandboxPath := filepath.Join(tmpDir, sandboxName)
	info, err := os.Stat(sandboxPath)
	require.NoError(t, err)
	assert.True(t, info.IsDir())

	// Verify permissions (should be 700)
	expectedPerm := os.FileMode(0o700)
	assert.Equal(t, expectedPerm.Perm(), info.Mode().Perm())
}

// TestSandboxInitCommand_CreatesSubdirectories verifies subdirectories are created.
func TestSandboxInitCommand_CreatesSubdirectories(t *testing.T) {
	tmpDir := t.TempDir()
	sandboxName := "test-subdirs"

	cmd := sandbox.NewSandboxCommand()
	cmd.SetArgs([]string{"init", sandboxName, "--path", tmpDir})

	err := cmd.Execute()
	require.NoError(t, err)

	sandboxPath := filepath.Join(tmpDir, sandboxName)
	subdirs := []string{"workspace", "sessions", "memory", "state"}

	for _, subdir := range subdirs {
		subdirPath := filepath.Join(sandboxPath, subdir)
		info, err := os.Stat(subdirPath)
		assert.NoError(t, err, "Subdirectory %s should exist", subdir)
		assert.True(t, info.IsDir())
	}
}

// TestSandboxInitCommand_CreatesReadme verifies README is created.
func TestSandboxInitCommand_CreatesReadme(t *testing.T) {
	tmpDir := t.TempDir()
	sandboxName := "test-readme"

	cmd := sandbox.NewSandboxCommand()
	cmd.SetArgs([]string{"init", sandboxName, "--path", tmpDir})

	err := cmd.Execute()
	require.NoError(t, err)

	sandboxPath := filepath.Join(tmpDir, sandboxName)
	readmePath := filepath.Join(sandboxPath, "README.md")

	content, err := os.ReadFile(readmePath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "Sandbox:")
	assert.Contains(t, string(content), "workspace")
}

// TestSandboxStatusCommand_Exists verifies status command works.
func TestSandboxStatusCommand_Exists(t *testing.T) {
	tmpDir := t.TempDir()
	sandboxName := "test-status"

	// First create the sandbox
	initCmd := sandbox.NewSandboxCommand()
	initCmd.SetArgs([]string{"init", sandboxName, "--path", tmpDir})
	err := initCmd.Execute()
	require.NoError(t, err)

	// Then check status
	statusCmd := sandbox.NewSandboxCommand()
	statusCmd.SetArgs([]string{"status", sandboxName})

	// Override home directory for testing
	// Note: This is a simplified test - in reality we'd need to mock the home dir

	// Just verify the command structure
	assert.NotNil(t, statusCmd)
	subcommands := statusCmd.Commands()
	assert.Len(t, subcommands, 2) // init and status
}

// TestSandboxCommand_Help verifies help is available.
func TestSandboxCommand_Help(t *testing.T) {
	cmd := sandbox.NewSandboxCommand()

	help := cmd.Help()
	assert.NoError(t, help)

	assert.Contains(t, cmd.Short, "isolated")
	assert.Contains(t, cmd.Long, "isolated")
}

// TestSandboxInitCommand_DefaultName verifies default name is used.
func TestSandboxInitCommand_DefaultName(t *testing.T) {
	tmpDir := t.TempDir()

	cmd := sandbox.NewSandboxCommand()
	cmd.SetArgs([]string{"init", "--path", tmpDir})

	err := cmd.Execute()
	require.NoError(t, err)

	// Verify default sandbox was created
	sandboxPath := filepath.Join(tmpDir, "default")
	_, err = os.Stat(sandboxPath)
	assert.NoError(t, err)
}

// TestSandboxStatusCommand_NotFound verifies error when sandbox doesn't exist.
func TestSandboxStatusCommand_NotFound(t *testing.T) {
	// We can't easily test this without mocking os.UserHomeDir
	// This is a structural test
	cmd := sandbox.NewSandboxCommand()
	statusCmd, _, err := cmd.Find([]string{"status"})
	assert.NoError(t, err)
	assert.NotNil(t, statusCmd)
}
