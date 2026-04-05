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
	"github.com/spf13/cobra"
)

// NewConfigCommand creates the config command.
func NewConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage configuration",
		Long:  "Configuration management commands: validate, show, backup.",
	}

	cmd.AddCommand(
		newValidateCommand(),
	)

	return cmd
}
