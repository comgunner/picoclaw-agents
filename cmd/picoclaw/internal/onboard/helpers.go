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
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/comgunner/picoclaw/cmd/picoclaw/internal"
	"github.com/comgunner/picoclaw/pkg/config"
)

func onboard(template string) {
	configPath := internal.GetConfigPath()

	if _, err := os.Stat(configPath); err == nil {
		fmt.Printf("Config already exists at %s\n", configPath)
		fmt.Print("Overwrite? (y/n): ")
		var response string
		fmt.Scanln(&response)
		if response != "y" {
			fmt.Println("Aborted.")
			return
		}
	}

	var cfg *config.Config
	switch template {
	case "glm":
		cfg = config.GLMDefaultConfig()
	case "qwen":
		cfg = config.QwenDefaultConfig(false)
	case "qwen_zh":
		cfg = config.QwenDefaultConfig(true)
	case "openrouter":
		cfg = config.OpenRouterDefaultConfig()
	case "openai":
		cfg = config.OpenAIDefaultConfig()
	case "gemini":
		cfg = config.GeminiDefaultConfig()
	case "free":
		cfg = config.OpenRouterFreeDefaultConfig()
	default:
		cfg = config.TemplateDefaultConfig()
	}

	if err := config.SaveConfig(configPath, cfg); err != nil {
		fmt.Printf("Error saving config: %v\n", err)
		os.Exit(1)
	}

	workspace := cfg.WorkspacePath()
	createWorkspaceTemplates(workspace)

	fmt.Printf("%s picoclaw is ready!\n", internal.Logo)
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Add your API key to", configPath)
	fmt.Println("")
	fmt.Println("     Selected Template:", template)
	switch template {
	case "glm":
		fmt.Println("     - Platform: https://z.ai/")
		fmt.Println("     - Keys:     https://z.ai/manage-apikey/apikey-list")
	case "qwen":
		fmt.Println("     - Platform: https://dashscope.aliyun.com/")
		fmt.Println("     - Keys:     https://dashscope.console.aliyun.com/apiKey")
	case "qwen_zh":
		fmt.Println("     - Platform: https://dashscope.aliyun.com/")
		fmt.Println("     - Keys:     https://dashscope.console.aliyun.com/apiKey")
	case "openrouter":
		fmt.Println("     - Platform: https://openrouter.ai/")
		fmt.Println("     - Keys:     https://openrouter.ai/settings/keys")
	case "openai":
		fmt.Println("     - Platform: https://platform.openai.com/")
		fmt.Println("     - Keys:     https://platform.openai.com/settings/organization/api-keys")
	case "gemini":
		fmt.Println("     - Platform: https://aistudio.google.com/")
		fmt.Println("     - Keys:     https://aistudio.google.com/app/apikey")
	case "free":
		fmt.Println("     - Platform: https://openrouter.ai/")
		fmt.Println("     - Keys:     https://openrouter.ai/settings/keys")
		fmt.Println("     - Note:     Free tier models — no balance required, just sign up")
	default:
		fmt.Println("     Recommended:")
		fmt.Println("     - OpenRouter: https://openrouter.ai/keys (access 100+ models)")
		fmt.Println("     - Ollama:     https://ollama.com (local, free)")
	}

	fmt.Println("\n  2. Chat: picoclaw agent -m \"Hello!\"")
}

func createWorkspaceTemplates(workspace string) {
	err := copyEmbeddedToTarget(workspace)
	if err != nil {
		fmt.Printf("Error copying workspace templates: %v\n", err)
	}
}

func copyEmbeddedToTarget(targetDir string) error {
	// Ensure target directory exists
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return fmt.Errorf("Failed to create target directory: %w", err)
	}

	// Walk through all files in embed.FS
	err := fs.WalkDir(embeddedFiles, "workspace", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if d.IsDir() {
			return nil
		}

		// Read embedded file
		data, err := embeddedFiles.ReadFile(path)
		if err != nil {
			return fmt.Errorf("Failed to read embedded file %s: %w", path, err)
		}

		new_path, err := filepath.Rel("workspace", path)
		if err != nil {
			return fmt.Errorf("Failed to get relative path for %s: %v\n", path, err)
		}

		// Build target file path
		targetPath := filepath.Join(targetDir, new_path)

		// Ensure target file's directory exists
		if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
			return fmt.Errorf("Failed to create directory %s: %w", filepath.Dir(targetPath), err)
		}

		// Write file
		if err := os.WriteFile(targetPath, data, 0o644); err != nil {
			return fmt.Errorf("Failed to write file %s: %w", targetPath, err)
		}

		return nil
	})

	return err
}
