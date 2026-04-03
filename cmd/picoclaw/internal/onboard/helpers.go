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
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/comgunner/picoclaw/cmd/picoclaw/internal"
	"github.com/comgunner/picoclaw/pkg/cli"
	"github.com/comgunner/picoclaw/pkg/config"
)

func onboard(template string, force bool) {
	configPath := internal.GetConfigPath()

	if _, err := os.Stat(configPath); err == nil {
		if !force {
			// BUG-001 FIX: Skip automático si config existe y no se usa --force
			fmt.Printf("✓ Config file already exists at %s\n", configPath)
			fmt.Println("✓ Skipping onboard to preserve existing configuration")
			fmt.Println("ℹ️  Use --force or -f to overwrite (not recommended)")
			return
		}

		// Solo si force=true, mostrar advertencia y procedure
		fmt.Println("⚠️  Existing config found, will overwrite...")
		fmt.Println("⚠️  WARNING: This will delete your existing configuration!")
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
		apiKey := promptFreeAPIKey()
		if apiKey != "" {
			for i := range cfg.ModelList {
				cfg.ModelList[i].APIKey = apiKey // pragma: allowlist secret
			}
		}
	default:
		cfg = config.TemplateDefaultConfig()
	}

	if err := config.SaveConfig(configPath, cfg); err != nil {
		fmt.Printf("Error saving config: %v\n", err)
		os.Exit(1)
	}

	workspace := cfg.WorkspacePath()
	createWorkspaceTemplates(workspace)

	fmt.Printf("%s picoclaw-agents is ready!\n", internal.Logo)
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
		fmt.Println("     - API key already saved ✅")
		fmt.Println("     - Platform: https://openrouter.ai/")
		fmt.Println("     - Model: openrouter/auto (routes to best free model with tool support)")
	default:
		fmt.Println("     Recommended:")
		fmt.Println("     - OpenRouter: https://openrouter.ai/keys (access 100+ models)")
		fmt.Println("     - Ollama:     https://ollama.com (local, free)")
	}

	fmt.Println("\n  2. Chat: picoclaw-agents agent -m \"Hello!\"")
}

// promptFreeAPIKey asks for an OpenRouter API key interactively.
// Returns the key, or empty string if the user skips.
func promptFreeAPIKey() string {
	fmt.Println()
	fmt.Println("╔════════════════════════════════════════════════════════╗")
	fmt.Println("║         Easy Setup: OpenRouter FREE Models             ║")
	fmt.Println("║   No credit card · No billing · Just try it!          ║")
	fmt.Println("╚════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("  OpenRouter requires an API key even for free models.")
	fmt.Println("  Sign up at https://openrouter.ai (email only, no card).")
	fmt.Println()
	fmt.Println("  ➜  1. Open:  https://openrouter.ai/settings/keys")
	fmt.Println("  ➜  2. Sign up (email only, no credit card)")
	fmt.Println("  ➜  3. Create a free API key (sk-or-v1-...)")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)
	apiKey, err := cli.ReadMaskedWithFallback("  Paste your OpenRouter API key (or Enter to skip): ", scanner)
	fmt.Println() // newline after masked input
	if err != nil || strings.TrimSpace(apiKey) == "" {
		fmt.Println("  ⚠️  No key provided. Edit ~/.picoclaw/config.json later to add it.")
		return ""
	}

	apiKey = strings.TrimSpace(apiKey)
	fmt.Printf("  Key received: %s***%s\n", apiKey[:min(6, len(apiKey))], apiKey[max(0, len(apiKey)-4):])
	if !strings.HasPrefix(apiKey, "sk-or-") {
		fmt.Println("  ⚠️  Key doesn't look like an OpenRouter key (expected sk-or-v1-...).")
		fmt.Print("  Continue anyway? (y/n): ")
		var resp string
		fmt.Scanln(&resp)
		if strings.ToLower(strings.TrimSpace(resp)) != "y" {
			return ""
		}
	} else {
		fmt.Println("  ✅ Key format looks good!")
	}
	fmt.Println()
	return apiKey
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
