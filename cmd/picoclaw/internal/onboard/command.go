// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package onboard

import (
	"embed"

	"github.com/spf13/cobra"
)

//go:generate cp -r ../../../../workspace .
//go:embed workspace
var embeddedFiles embed.FS

func NewOnboardCommand() *cobra.Command {
	var (
		glm        bool
		qwen       bool
		qwenZh     bool
		openrouter bool
		openai     bool
		gemini     bool
	)

	cmd := &cobra.Command{
		Use:     "onboard",
		Aliases: []string{"o"},
		Short:   "Initialize picoclaw configuration and workspace",
		Run: func(cmd *cobra.Command, args []string) {
			template := ""
			if glm {
				template = "glm"
			} else if qwen {
				template = "qwen"
			} else if qwenZh {
				template = "qwen_zh"
			} else if openrouter {
				template = "openrouter"
			} else if openai {
				template = "openai"
			} else if gemini {
				template = "gemini"
			}
			onboard(template)
		},
	}

	cmd.Flags().BoolVar(&glm, "glm", false, "Use GLM-4.5-Flash template (z.ai)")
	cmd.Flags().BoolVar(&qwen, "qwen", false, "Use Qwen template (Alibaba Cloud Intl)")
	cmd.Flags().BoolVar(&qwenZh, "qwen_zh", false, "Use Qwen template (Alibaba Cloud China)")
	cmd.Flags().BoolVar(&openrouter, "openrouter", false, "Use OpenRouter template (openrouter/auto)")
	cmd.Flags().BoolVar(&openai, "openai", false, "Use OpenAI template (o3-mini)")
	cmd.Flags().BoolVar(&gemini, "gemini", false, "Use Gemini template (gemini-2.5-flash)")

	return cmd
}
