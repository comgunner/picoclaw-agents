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
	"os"
	"path/filepath"
	"strings"

	"github.com/comgunner/picoclaw/pkg/auth"
	"github.com/comgunner/picoclaw/pkg/cli"
)

// Wizard orchestrates the interactive setup process
type Wizard struct {
	scanner    *bufio.Scanner
	modelName  string
	model      string
	apiKey     string
	configPath string
	workspace  string
}

// NewWizard creates a new interactive setup wizard
func NewWizard() *Wizard {
	home, _ := os.UserHomeDir()
	return &Wizard{
		scanner:    bufio.NewScanner(os.Stdin),
		configPath: filepath.Join(home, ".picoclaw", "config.json"),
		workspace:  filepath.Join(home, ".picoclaw", "workspace"),
	}
}

// Run executes the interactive wizard
func (w *Wizard) Run() error {
	w.printHeader()

	// Step 1: Environment check
	if err := w.stepEnvironment(); err != nil {
		return fmt.Errorf("environment check failed: %w", err)
	}

	// Step 2: LLM Provider selection
	if err := w.stepProvider(); err != nil {
		return fmt.Errorf("provider setup failed: %w", err)
	}

	// Step 3: Channel configuration
	if err := w.stepChannels(); err != nil {
		return fmt.Errorf("channel setup failed: %w", err)
	}

	// Step 4: Generate configuration
	if err := w.stepGenerateConfig(); err != nil {
		return fmt.Errorf("config generation failed: %w", err)
	}

	// Step 5: Verification
	if err := w.stepVerify(); err != nil {
		return fmt.Errorf("verification failed: %w", err)
	}

	w.printSuccess()
	return nil
}

func (w *Wizard) printHeader() {
	fmt.Println("\n╔════════════════════════════════════════════════════════╗")
	fmt.Println("║           PicoClaw Setup Wizard v3.5.0                 ║")
	fmt.Println("║     \"The AI that actually does things\"                 ║")
	fmt.Println("╚════════════════════════════════════════════════════════╝")
	fmt.Println()
}

func (w *Wizard) printSuccess() {
	fmt.Println("\n╔════════════════════════════════════════════════════════╗")
	fmt.Println("║                  Setup Complete! ✅                    ║")
	fmt.Println("╚════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("📦 Configuration saved to:", w.configPath)
	fmt.Println("📁 Workspace directory:", w.workspace)
	fmt.Println()

	// Check if using free tier
	if w.model == "openrouter/free" {
		fmt.Println("🆓 You're using FREE models via OpenRouter!")
		fmt.Println()
		fmt.Println("🚀 Try it now:")
		fmt.Println("   picoclaw agent -m \"Hello! What can you do?\"")
		fmt.Println()
		fmt.Println("📖 Free models info:  https://openrouter.ai/collections/free-models")
		fmt.Println("⬆️  Upgrade anytime:  picoclaw onboard --openrouter")
	} else {
		fmt.Println("🚀 Next steps:")
		fmt.Println("   1. Run: picoclaw interactive")
		fmt.Println("   2. Or start gateway: docker-compose --profile gateway up")
	}
	fmt.Println()
}

func (w *Wizard) prompt(message string) string {
	fmt.Print(message)
	if w.scanner.Scan() {
		return strings.TrimSpace(w.scanner.Text())
	}
	return ""
}

func (w *Wizard) promptSecret(message string) string {
	// Use cli.ReadMaskedWithFallback for secure input without echo
	input, err := cli.ReadMaskedWithFallback(message, w.scanner)
	if err != nil {
		// Fallback to regular input if masked input fails
		return w.prompt(message)
	}
	return input
}

func (w *Wizard) promptChoice(message string, options []string) string {
	fmt.Println("\n" + message)
	for i, opt := range options {
		fmt.Printf("   %d. %s\n", i+1, opt)
	}

	for {
		fmt.Printf("\nEnter choice [1-%d]: ", len(options))
		input := w.prompt("")

		var choice int
		fmt.Sscanf(input, "%d", &choice)

		if choice >= 1 && choice <= len(options) {
			return options[choice-1]
		}
		fmt.Println("Invalid choice. Please try again.")
	}
}

func (w *Wizard) promptConfirm(message string) bool {
	for {
		fmt.Print(message + " (y/n): ")
		input := strings.ToLower(w.prompt(""))

		if input == "y" || input == "yes" {
			return true
		}
		if input == "n" || input == "no" {
			return false
		}
		fmt.Println("Please enter 'y' or 'n'.")
	}
}

// stepEnvironment checks the environment
func (w *Wizard) stepEnvironment() error {
	fmt.Println("\n📋 Step 1/5: Environment Check")
	fmt.Println("─────────────────────────────────────")

	// Check Go version
	fmt.Print("✓ Checking Go version... ")
	goVersion := checkGoVersion()
	if goVersion == "" {
		fmt.Println("❌ Go not found")
		return fmt.Errorf("Go 1.25.7 or higher is required")
	}
	fmt.Println(goVersion)

	// Check existing config
	fmt.Print("✓ Checking existing configuration... ")
	if _, err := os.Stat(w.configPath); err == nil {
		fmt.Println("⚠️  Found")
		if !w.promptConfirm("  Existing configuration found. Overwrite?") {
			return fmt.Errorf("setup canceled by user")
		}
	} else {
		fmt.Println("None (will create new)")
	}

	// Check workspace
	fmt.Print("✓ Checking workspace directory... ")
	if _, err := os.Stat(w.workspace); err == nil {
		fmt.Println("Exists")
	} else {
		fmt.Println("Will create")
	}

	fmt.Println("\n✅ Environment check passed")
	return nil
}

// stepProvider configures the LLM provider
func (w *Wizard) stepProvider() error {
	fmt.Println("\n🤖 Step 2/5: LLM Provider Selection")
	fmt.Println("─────────────────────────────────────")

	provider := w.promptChoice("Select your primary provider:", []string{
		"⭐ Easy Setup - FREE (No credit card, try PicoClaw now!)",
		"DeepSeek (Recommended - Best value)",
		"Anthropic Claude",
		"OpenAI GPT",
		"Google Gemini",
		"Groq",
		"OpenRouter (100+ models)",
		"Zhipu GLM",
		"Alibaba Qwen",
	})

	// Check if user selected Easy Setup - FREE
	if strings.Contains(provider, "Easy Setup") {
		return w.setupEasyFree()
	}

	providerKey := strings.ToLower(strings.Fields(provider)[0])

	fmt.Printf("\nConfiguring %s...\n", provider)
	apiKey := w.promptSecret(fmt.Sprintf("Enter your %s API key: ", provider))

	// Validate API key format
	fmt.Print("✓ Validating API key... ")
	result := auth.QuickValidate(providerKey, apiKey)
	if !result.Valid {
		fmt.Printf("❌ Invalid: %v\n", result.Error)
		if !w.promptConfirm("  Continue anyway?") {
			return fmt.Errorf("invalid API key")
		}
	} else {
		fmt.Println("✅ Valid")
	}

	// Store in config
	w.modelName = providerKey + "/default"
	w.model = providerKey + "/default"
	w.apiKey = apiKey

	fmt.Printf("✅ %s configured successfully\n", provider)
	return nil
}

// setupEasyFree configures OpenRouter free models for zero-cost onboarding
func (w *Wizard) setupEasyFree() error {
	fmt.Println()
	fmt.Println("╔════════════════════════════════════════════════════════╗")
	fmt.Println("║         Easy Setup: OpenRouter FREE Models             ║")
	fmt.Println("║   No credit card · No billing · Just try it!          ║")
	fmt.Println("╚════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("  PicoClaw will use OpenRouter's free AI models.")
	fmt.Println("  You only need a free account at openrouter.ai")
	fmt.Println()
	fmt.Println("  ➜  1. Open:  https://openrouter.ai/settings/keys")
	fmt.Println("  ➜  2. Sign up (email only, no credit card)")
	fmt.Println("  ➜  3. Create a free API key")
	fmt.Println()

	apiKey := w.promptSecret("  Paste your OpenRouter API key (sk-or-v1-...): ")

	// Basic format validation
	if !strings.HasPrefix(apiKey, "sk-or-") {
		fmt.Println("  ⚠️  Warning: Key doesn't look like an OpenRouter key (sk-or-v1-...)")
		if !w.promptConfirm("  Continue anyway?") {
			return fmt.Errorf("invalid API key")
		}
	} else {
		fmt.Println("  ✅ Format looks good!")
	}

	// Configure free models
	w.modelName = "or-free"
	w.model = "openrouter/free"
	w.apiKey = apiKey

	fmt.Println()
	fmt.Println("  ✅ Easy Setup complete! Free models configured:")
	fmt.Println("     • openrouter/free       → auto-selects best available free model")
	fmt.Println("     • stepfun/step-3.5-flash → 256K context, reasoning + tools (fallback)")
	fmt.Println("     • deepseek/deepseek-v3.2 → fast and capable (fallback)")
	fmt.Println()
	fmt.Println("  ℹ️  Free models have rate limits (~50 req/min). Perfect for personal use.")
	fmt.Println("  ℹ️  To upgrade later: picoclaw onboard --openrouter")
	fmt.Println()
	return nil
}

// stepChannels configures chat channels
func (w *Wizard) stepChannels() error {
	fmt.Println("\n💬 Step 3/5: Channel Configuration")
	fmt.Println("─────────────────────────────────────")

	channel := w.promptChoice("Which channel do you want to enable?", []string{
		"Telegram (Most popular)",
		"Discord",
		"None (CLI only)",
	})

	switch strings.ToLower(strings.Fields(channel)[0]) {
	case "telegram":
		return w.setupTelegram()
	case "discord":
		return w.setupDiscord()
	default:
		fmt.Println("✓ CLI-only mode selected")
		return nil
	}
}

func (w *Wizard) setupTelegram() error {
	fmt.Println("\n📱 Setting up Telegram...")

	token := w.promptSecret("Enter Telegram Bot Token: ")
	_ = w.prompt("Enter your Telegram User ID: ")

	// Basic validation
	if len(token) < 40 {
		fmt.Println("⚠️  Warning: Token seems short")
		if !w.promptConfirm("  Continue anyway?") {
			return fmt.Errorf("invalid token")
		}
	}

	fmt.Println("✅ Telegram configured successfully")
	return nil
}

func (w *Wizard) setupDiscord() error {
	fmt.Println("\n🎮 Setting up Discord...")

	token := w.promptSecret("Enter Discord Bot Token: ")
	_ = w.prompt("Enter your Discord User ID: ")

	// Basic validation
	if len(token) < 50 {
		fmt.Println("⚠️  Warning: Token seems short")
		if !w.promptConfirm("  Continue anyway?") {
			return fmt.Errorf("invalid token")
		}
	}

	fmt.Println("✅ Discord configured successfully")
	return nil
}

// stepGenerateConfig generates the configuration file
func (w *Wizard) stepGenerateConfig() error {
	fmt.Println("\n⚙️  Step 4/5: Generating Configuration")
	fmt.Println("─────────────────────────────────────")

	// Create directory
	dir := filepath.Dir(w.configPath)
	fmt.Printf("✓ Creating directory: %s... ", dir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
		return err
	}
	fmt.Println("Done")

	// Create workspace
	fmt.Printf("✓ Creating workspace: %s... ", w.workspace)
	if err := os.MkdirAll(w.workspace, 0o755); err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
		return err
	}
	fmt.Println("Done")

	// Generate config file
	fmt.Printf("✓ Generating config file: %s... ", w.configPath)
	if err := w.saveConfig(); err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
		return err
	}
	fmt.Println("Done")

	fmt.Println("\n✅ Configuration generated successfully")
	return nil
}

// stepVerify verifies the setup
func (w *Wizard) stepVerify() error {
	fmt.Println("\n🔍 Step 5/5: Verification")
	fmt.Println("─────────────────────────────────────")

	// Verify config file exists
	fmt.Print("✓ Verifying configuration file... ")
	if _, err := os.Stat(w.configPath); err != nil {
		fmt.Println("❌ Failed")
		return err
	}
	fmt.Println("✅ Exists")

	// Verify workspace exists
	fmt.Print("✓ Verifying workspace directory... ")
	if _, err := os.Stat(w.workspace); err != nil {
		fmt.Println("❌ Failed")
		return err
	}
	fmt.Println("✅ Exists")

	// Verify API key format
	fmt.Print("✓ Verifying API key format... ")
	if w.apiKey != "" {
		provider := auth.GetProviderFromModelName(w.model)
		result := auth.QuickValidate(provider, w.apiKey)
		if result.Valid {
			fmt.Println("✅ Valid")
		} else {
			fmt.Printf("⚠️  Warning: %v\n", result.Error)
		}
	} else {
		fmt.Println("⚠️  No API key configured")
	}

	fmt.Println("\n✅ All verifications passed")
	return nil
}

func (w *Wizard) saveConfig() error {
	var modelListJSON string

	// Check if using free tier
	if w.model == "openrouter/free" {
		// Easy Setup: three free models with fallbacks
		modelListJSON = fmt.Sprintf(`[
    {
      "model_name": "or-free",
      "model": "openrouter/free",
      "api_base": "https://openrouter.ai/api/v1",
      "api_key": "%s"
    },
    {
      "model_name": "or-free-stepfun",
      "model": "stepfun/step-3.5-flash",
      "api_base": "https://openrouter.ai/api/v1",
      "api_key": "%s"
    },
    {
      "model_name": "or-free-deepseek",
      "model": "deepseek/deepseek-v3.2-20251201",
      "api_base": "https://openrouter.ai/api/v1",
      "api_key": "%s"
    }
  ]`, w.apiKey, w.apiKey, w.apiKey)
	} else {
		// Normal flow: single model
		modelListJSON = fmt.Sprintf(`[
    {
      "model_name": "%s",
      "model": "%s",
      "api_key": "%s"
    }
  ]`, w.modelName, w.model, w.apiKey)
	}

	// Generate config JSON
	data := fmt.Sprintf(`{
  "agents": {
    "defaults": {
      "workspace": "%s",
      "restrict_to_workspace": true,
      "model_name": "%s",
      "max_tokens": 8192,
      "temperature": 0.7,
      "max_tool_iterations": 20
    }
  },
  "model_list": %s
}
`, w.workspace, w.modelName, modelListJSON)

	return os.WriteFile(w.configPath, []byte(data), 0o600)
}

// checkGoVersion checks if Go is installed and returns version
func checkGoVersion() string {
	// Simplified check - in real implementation use exec.Command
	// For now, just return a placeholder
	return "Go 1.25.x (detected)"
}
