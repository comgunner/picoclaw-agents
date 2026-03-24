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

func newShowCommand(loaderFn func() (*skills.SkillsLoader, error)) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "show",
		Short:   "Show skill details",
		Args:    cobra.ExactArgs(1),
		Example: `picoclaw skills show weather`,
		RunE: func(_ *cobra.Command, args []string) error {
			loader, err := loaderFn()
			if err != nil {
				return err
			}
			skillsShowCmd(loader, args[0])
			return nil
		},
	}

	return cmd
}
