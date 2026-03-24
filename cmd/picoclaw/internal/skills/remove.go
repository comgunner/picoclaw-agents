// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package skills

import (
	"github.com/spf13/cobra"

	"github.com/comgunner/picoclaw/pkg/skills"
)

func newRemoveCommand(installerFn func() (*skills.SkillInstaller, error)) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove",
		Aliases: []string{"rm", "uninstall"},
		Short:   "Remove installed skill",
		Args:    cobra.ExactArgs(1),
		Example: `picoclaw skills remove weather`,
		RunE: func(_ *cobra.Command, args []string) error {
			installer, err := installerFn()
			if err != nil {
				return err
			}
			skillsRemoveCmd(installer, args[0])
			return nil
		},
	}

	return cmd
}
