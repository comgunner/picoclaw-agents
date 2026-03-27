// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package sandbox

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// NewSandboxCommand creates the sandbox command.
func NewSandboxCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sandbox",
		Short: "Manage isolated workspaces",
		Long:  "Create and manage isolated sandbox workspaces for safe experimentation.",
	}

	cmd.AddCommand(
		newInitCommand(),
		newStatusCommand(),
	)

	return cmd
}

func newInitCommand() *cobra.Command {
	var sandboxPath string

	cmd := &cobra.Command{
		Use:   "init [name]",
		Short: "Initialize a new sandbox workspace",
		Long:  "Create a new isolated sandbox directory with restricted permissions.",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := "default"
			if len(args) > 0 {
				name = args[0]
			}
			return runInit(name, sandboxPath)
		},
	}

	cmd.Flags().StringVarP(&sandboxPath, "path", "p", "", "Base path for sandboxes (default: ~/.picoclaw/sandboxes)")
	return cmd
}

func runInit(name, basePath string) error {
	if basePath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("getting home directory: %w", err)
		}
		basePath = filepath.Join(home, ".picoclaw", "sandboxes")
	}

	sandboxDir := filepath.Join(basePath, name)

	// Create directory with restricted permissions (700 = rwx------)
	if err := os.MkdirAll(sandboxDir, 0o700); err != nil {
		return fmt.Errorf("creating sandbox directory: %w", err)
	}

	// Create subdirectories
	subdirs := []string{"workspace", "sessions", "memory", "state"}
	for _, subdir := range subdirs {
		subdirPath := filepath.Join(sandboxDir, subdir)
		if err := os.MkdirAll(subdirPath, 0o700); err != nil {
			return fmt.Errorf("creating %s directory: %w", subdir, err)
		}
	}

	// Create a README
	readmePath := filepath.Join(sandboxDir, "README.md")
	readme := fmt.Sprintf(`# Sandbox: %s

This is an isolated workspace with restricted permissions.

## Structure
- **workspace/** - Safe directory for file operations
- **sessions/** - Conversation history
- **memory/** - Long-term memory
- **state/** - Persistent state

## Security
- Directory permissions: 700 (owner read/write/execute only)
- Isolated from other sandboxes and main workspace

`, name)

	if err := os.WriteFile(readmePath, []byte(readme), 0o600); err != nil {
		return fmt.Errorf("creating README: %w", err)
	}

	fmt.Printf("✓ Sandbox '%s' initialized at: %s\n", name, sandboxDir)
	fmt.Println("\nSubdirectories created:")
	for _, subdir := range subdirs {
		fmt.Printf("  - %s\n", filepath.Join(sandboxDir, subdir))
	}
	fmt.Println("\nUsage:")
	fmt.Printf("  picoclaw agent --workspace %s\n", filepath.Join(sandboxDir, "workspace"))

	return nil
}

func newStatusCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status [name]",
		Short: "Check sandbox status and permissions",
		Long:  "Verify sandbox directory exists and has correct permissions.",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := "default"
			if len(args) > 0 {
				name = args[0]
			}
			return runStatus(name)
		},
	}

	return cmd
}

func runStatus(name string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("getting home directory: %w", err)
	}

	sandboxDir := filepath.Join(home, ".picoclaw", "sandboxes", name)

	// Check if exists
	info, err := os.Stat(sandboxDir)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("sandbox '%s' does not exist", name)
		}
		return fmt.Errorf("checking sandbox: %w", err)
	}

	// Check permissions
	mode := info.Mode()
	fmt.Printf("Sandbox: %s\n", name)
	fmt.Printf("Path: %s\n", sandboxDir)
	fmt.Printf("Permissions: %s (%o)\n", mode.String(), mode.Perm())
	fmt.Printf("Is Dir: %v\n", info.IsDir())

	// Check if permissions are restrictive (700)
	if mode.Perm()&0o077 == 0 {
		fmt.Println("✓ Permissions are restrictive (no group/other access)")
	} else {
		fmt.Println("⚠️  WARNING: Permissions allow group/other access!")
	}

	// Check subdirectories
	subdirs := []string{"workspace", "sessions", "memory", "state"}
	fmt.Println("\nSubdirectories:")
	for _, subdir := range subdirs {
		subdirPath := filepath.Join(sandboxDir, subdir)
		if _, err := os.Stat(subdirPath); err == nil {
			fmt.Printf("  ✓ %s\n", subdir)
		} else {
			fmt.Printf("  ✗ %s (missing)\n", subdir)
		}
	}

	return nil
}
