// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/comgunner/picoclaw/cmd/picoclaw/internal"
	"github.com/comgunner/picoclaw/cmd/picoclaw/internal/agent"
	"github.com/comgunner/picoclaw/cmd/picoclaw/internal/agents"
	"github.com/comgunner/picoclaw/cmd/picoclaw/internal/auth"
	"github.com/comgunner/picoclaw/cmd/picoclaw/internal/clean"
	"github.com/comgunner/picoclaw/cmd/picoclaw/internal/config"
	"github.com/comgunner/picoclaw/cmd/picoclaw/internal/cron"
	"github.com/comgunner/picoclaw/cmd/picoclaw/internal/doctor"
	"github.com/comgunner/picoclaw/cmd/picoclaw/internal/gateway"
	"github.com/comgunner/picoclaw/cmd/picoclaw/internal/migrate"
	"github.com/comgunner/picoclaw/cmd/picoclaw/internal/onboard"
	"github.com/comgunner/picoclaw/cmd/picoclaw/internal/sandbox"
	"github.com/comgunner/picoclaw/cmd/picoclaw/internal/skills"
	"github.com/comgunner/picoclaw/cmd/picoclaw/internal/status"
	"github.com/comgunner/picoclaw/cmd/picoclaw/internal/util"
	"github.com/comgunner/picoclaw/cmd/picoclaw/internal/version"
)

func NewPicoclawCommand() *cobra.Command {
	short := fmt.Sprintf("%s picoclaw - Personal AI Assistant v%s\n\n", internal.Logo, internal.GetVersion())

	cmd := &cobra.Command{
		Use:     "picoclaw",
		Short:   short,
		Example: "picoclaw list",
	}

	cmd.AddCommand(
		onboard.NewOnboardCommand(),
		agent.NewAgentCommand(),
		agents.NewAgentsCommand(),
		auth.NewAuthCommand(),
		clean.NewCleanCommand(),
		config.NewConfigCommand(),
		cron.NewCronCommand(),
		doctor.NewDoctorCommand(),
		gateway.NewGatewayCommand(),
		migrate.NewMigrateCommand(),
		sandbox.NewSandboxCommand(),
		skills.NewSkillsCommand(),
		status.NewStatusCommand(),
		util.NewUtilCommand(),
		version.NewVersionCommand(),
	)

	return cmd
}

func main() {
	cmd := NewPicoclawCommand()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
