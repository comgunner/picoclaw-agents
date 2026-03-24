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

func newSearchCommand(installerFn func() (*skills.SkillInstaller, error)) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search available skills",
		RunE: func(_ *cobra.Command, _ []string) error {
			installer, err := installerFn()
			if err != nil {
				return err
			}
			skillsSearchCmd(installer)
			return nil
		},
	}

	return cmd
}
