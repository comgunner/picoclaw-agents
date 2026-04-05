// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/comgunner/picoclaw/pkg/config"
)

func newValidateCommand() *cobra.Command {
	var configPath string

	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate config.json schema and values",
		Long:  "Validate config.json for syntax, required fields, and semantic correctness.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runValidate(configPath)
		},
	}

	cmd.Flags().StringVarP(&configPath, "config", "c", "", "path to config.json (default: ~/.picoclaw/config.json)")
	return cmd
}

func runValidate(configPath string) error {
	// Use default path if not specified
	if configPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("getting home directory: %w", err)
		}
		configPath = filepath.Join(home, ".picoclaw", "config.json")
	}

	// Validate
	v := config.NewValidator()
	if err := v.ValidateFile(configPath); err != nil {
		return err
	}

	fmt.Printf("✓ Config is valid: %s\n", configPath)
	return nil
}
