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
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/comgunner/picoclaw/cmd/picoclaw/internal"
	"github.com/comgunner/picoclaw/pkg/skills"
)

type deps struct {
	workspace    string
	installer    *skills.SkillInstaller
	skillsLoader *skills.SkillsLoader
}

func NewSkillsCommand() *cobra.Command {
	var d deps

	cmd := &cobra.Command{
		Use:   "skills",
		Short: "Manage skills",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := internal.LoadConfig()
			if err != nil {
				return fmt.Errorf("error loading config: %w", err)
			}

			d.workspace = cfg.WorkspacePath()
			d.installer = skills.NewSkillInstaller(d.workspace)

			// get global config directory and builtin skills directory
			globalDir := filepath.Dir(internal.GetConfigPath())
			globalSkillsDir := filepath.Join(globalDir, "skills")
			builtinSkillsDir := filepath.Join(globalDir, "picoclaw", "skills")
			d.skillsLoader = skills.NewSkillsLoader(d.workspace, globalSkillsDir, builtinSkillsDir)

			return nil
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}

	installerFn := func() (*skills.SkillInstaller, error) {
		if d.installer == nil {
			return nil, fmt.Errorf("skills installer is not initialized")
		}
		return d.installer, nil
	}

	loaderFn := func() (*skills.SkillsLoader, error) {
		if d.skillsLoader == nil {
			return nil, fmt.Errorf("skills loader is not initialized")
		}
		return d.skillsLoader, nil
	}

	workspaceFn := func() (string, error) {
		if d.workspace == "" {
			return "", fmt.Errorf("workspace is not initialized")
		}
		return d.workspace, nil
	}

	cmd.AddCommand(
		newListCommand(loaderFn),
		newInstallCommand(installerFn),
		newInstallBuiltinCommand(workspaceFn),
		newListBuiltinCommand(),
		newRemoveCommand(installerFn),
		newSearchCommand(installerFn),
		newShowCommand(loaderFn),
	)

	return cmd
}
