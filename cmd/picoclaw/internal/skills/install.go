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
	"fmt"

	"github.com/spf13/cobra"

	"github.com/comgunner/picoclaw/cmd/picoclaw/internal"
	"github.com/comgunner/picoclaw/pkg/skills"
)

func newInstallCommand(installerFn func() (*skills.SkillInstaller, error)) *cobra.Command {
	var registry string

	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install skill from GitHub",
		Example: `
picoclaw skills install comgunner/picoclaw-skills/weather
picoclaw skills install --registry clawhub github
`,
		Args: func(cmd *cobra.Command, args []string) error {
			if registry != "" {
				if len(args) != 2 {
					return fmt.Errorf("when --registry is set, exactly 2 arguments are required: <name> <slug>")
				}
				return nil
			}

			if len(args) != 1 {
				return fmt.Errorf("exactly 1 argument is required: <github>")
			}

			return nil
		},
		RunE: func(_ *cobra.Command, args []string) error {
			installer, err := installerFn()
			if err != nil {
				return err
			}

			if registry != "" {
				cfg, err := internal.LoadConfig()
				if err != nil {
					return err
				}

				return skillsInstallFromRegistry(cfg, args[0], args[1])
			}

			return skillsInstallCmd(installer, args[0])
		},
	}

	cmd.Flags().StringVar(&registry, "registry", "", "Install from registry: --registry <name> <slug>")

	return cmd
}
