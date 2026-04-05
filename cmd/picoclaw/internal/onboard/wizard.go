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
	"github.com/comgunner/picoclaw/pkg/config"
)

// skillInfo represents a native skill with its description
type skillInfo struct {
	id   string
	name string
	desc string
}

// nativeSkillsList returns the list of native skills available for selection
func nativeSkillsList() []skillInfo {
	allSkills := getNativeSkills()
	result := make([]skillInfo, len(allSkills))
	for i, skill := range allSkills {
		result[i] = skillInfo{
			id:   skill,
			name: skill,
			desc: getSkillDescription(skill),
		}
	}
	return result
}

// Wizard orchestrates the interactive setup process
type Wizard struct {
	scanner    *bufio.Scanner
	modelName  string
	model      string
	apiKey     string
	configPath string
	workspace  string

	// Channel configuration — BUG-03 FIX: persist channel config
	channelType   string // "telegram" | "discord" | ""
	channelToken  string // bot token
	channelUserID string // user/server ID for allowed_users

	// Agent mode — SPRINT 2 ONBOARD: team mode and skills
	agentMode     string   // "solo" | "team"
	agentTemplate string   // "dev" | "research" | "general" (for team mode)
	customSkills  []string // skills chosen in solo/custom mode
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

	// Step 4: Agent mode selection
	if err := w.stepMode(); err != nil {
		return fmt.Errorf("agent mode setup failed: %w", err)
	}

	// Step 5: Generate configuration
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
	fmt.Printf("║           PicoClaw Setup Wizard %-22s ║\n", config.Version)
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

	// BUG-03 FIX: Show channel status in success message
	if w.channelType == "telegram" && w.channelToken != "" {
		fmt.Println("💬 Channel: Telegram ✅")
	} else if w.channelType == "discord" && w.channelToken != "" {
		fmt.Println("🎮 Channel: Discord ✅")
	} else {
		fmt.Println("💬 Channel: CLI-only mode")
	}

	// SPRINT 2 ONBOARD: Show agent mode and skills
	if w.agentMode == "team" {
		templateName := w.agentTemplate
		switch w.agentTemplate {
		case "dev":
			templateName = "Dev Team (9 agents)"
		case "research":
			templateName = "Research Team (3 agents)"
		case "general":
			templateName = "General Team (3 agents)"
		}
		fmt.Printf("🤖 Mode:      %s\n", templateName)
	} else {
		if len(w.customSkills) > 0 {
			fmt.Printf("🤖 Mode:      Solo Agent with skills: %s\n", strings.Join(w.customSkills, ", "))
		} else {
			fmt.Println("🤖 Mode:      Solo Agent")
		}
	}
	fmt.Println()

	// Check if using free tier
	if w.model == "openrouter/auto" || w.model == "or-auto" {
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
		fmt.Println()
		fmt.Println("⚠️  WARNING: Existing configuration found!")
		fmt.Println("   To preserve your config:")
		fmt.Println("     1. Backup: cp ~/.picoclaw/config.json ~/.picoclaw/config.json.bak")
		fmt.Println("     2. Edit:   Edit the existing file instead")
		fmt.Println("     3. Skip:   Remove ~/.picoclaw/config.json to start fresh")
		fmt.Println()
		return fmt.Errorf("setup canceled: config already exists")
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
	w.modelName = "or-auto"
	w.model = "openrouter/auto"
	w.apiKey = apiKey

	fmt.Println()
	fmt.Println("  ✅ Easy Setup complete! Free model configured:")
	fmt.Println("     • openrouter/auto → auto-selects best free model with tool support")
	fmt.Println()
	fmt.Println("  ℹ️  Free models have rate limits. Perfect for personal use.")
	fmt.Println("  ℹ️  To upgrade later: picoclaw-agents onboard --openrouter")
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

	// BUG-03 FIX: Store token and userID in wizard struct instead of local variables
	w.channelToken = w.promptSecret("Enter Telegram Bot Token: ")
	w.channelUserID = w.prompt("Enter your Telegram User ID: ")
	w.channelType = "telegram"

	// Basic validation
	if len(w.channelToken) < 40 {
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

	// BUG-03 FIX: Store token and userID in wizard struct instead of local variables
	w.channelToken = w.promptSecret("Enter Discord Bot Token: ")
	w.channelUserID = w.prompt("Enter your Discord User/Server ID: ")
	w.channelType = "discord"

	// Basic validation
	if len(w.channelToken) < 50 {
		fmt.Println("⚠️  Warning: Token seems short")
		if !w.promptConfirm("  Continue anyway?") {
			return fmt.Errorf("invalid token")
		}
	}

	fmt.Println("✅ Discord configured successfully")
	return nil
}

// stepMode lets the user choose between solo agent and team templates.
func (w *Wizard) stepMode() error {
	fmt.Println("\n🤖 Step 4/6: Agent Mode")
	fmt.Println("─────────────────────────────────────")

	mode := w.promptChoice("How do you want to use PicoClaw?", []string{
		"Solo Agent       — One general-purpose agent (simple, recommended to start)",
		"Dev Team         — Engineering team: manager + 8 specialists",
		"Research Team    — Coordinator + researcher + analyst",
		"General Team     — Orchestrator + 2 workers",
	})

	switch {
	case strings.Contains(mode, "Dev Team"):
		w.agentMode = "team"
		w.agentTemplate = "dev"
		fmt.Println("✅ Dev Team selected (engineering_manager + 8 specialists)")
	case strings.Contains(mode, "Research"):
		w.agentMode = "team"
		w.agentTemplate = "research"
		fmt.Println("✅ Research Team selected (coordinator + researcher + analyst)")
	case strings.Contains(mode, "General"):
		w.agentMode = "team"
		w.agentTemplate = "general"
		fmt.Println("✅ General Team selected (orchestrator + 2 workers)")
	default:
		w.agentMode = "solo"
		w.agentTemplate = ""
		w.customSkills = w.promptSkills()
		if len(w.customSkills) > 0 {
			fmt.Printf("✅ Solo Agent with skills: %s\n", strings.Join(w.customSkills, ", "))
		} else {
			fmt.Println("✅ Solo Agent (no extra skills)")
		}
	}
	return nil
}

// promptSkills shows the native skills list and lets the user pick.
func (w *Wizard) promptSkills() []string {
	skills := nativeSkillsList()
	fmt.Println("\n  Available native skills (optional):")
	for i, s := range skills {
		fmt.Printf("    %2d. %-28s — %s\n", i+1, s.name, s.desc)
	}
	fmt.Print("\n  Enter numbers (e.g. \"1 3\") or Enter to skip: ")
	input := w.prompt("")
	if strings.TrimSpace(input) == "" {
		return nil
	}
	var chosen []string
	for _, part := range strings.Fields(input) {
		var n int
		if _, err := fmt.Sscanf(part, "%d", &n); err == nil && n >= 1 && n <= len(skills) {
			chosen = append(chosen, skills[n-1].id)
		}
	}
	return chosen
}

// stepGenerateConfig generates the configuration file
func (w *Wizard) stepGenerateConfig() error {
	fmt.Println("\n⚙️  Step 5/6: Generating Configuration")
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
	fmt.Println("\n🔍 Step 6/6: Verification")
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
	if w.model == "openrouter/auto" {
		// Easy Setup: three free models with fallbacks
		// FIX: Use "openrouter/auto" instead of "openrouter/free" (issue #901)
		modelListJSON = fmt.Sprintf(`[
    {
      "model_name": "or-auto",
      "model": "openrouter/auto",
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
      "model": "deepseek/deepseek-v3.2",
      "api_base": "https://openrouter.ai/api/v1",
      "api_key": "%s"
    }
  ]`, w.apiKey, w.apiKey, w.apiKey)
		w.model = "or-auto" // Use the alias as default
	} else {
		modelListJSON = fmt.Sprintf(`[
    {
      "model_name": "%s",
      "model": "%s",
      "api_key": "%s"
    }
  ]`, w.modelName, w.model, w.apiKey)
	}

	// BUG-03 FIX: Generate channels JSON if channel is configured
	channelsJSON := ""
	if w.channelType == "telegram" && w.channelToken != "" {
		channelsJSON = fmt.Sprintf(`,
  "channels": {
    "telegram": {
      "token": "%s",
      "allowed_users": ["%s"]
    }
  }`, w.channelToken, w.channelUserID)
	} else if w.channelType == "discord" && w.channelToken != "" {
		channelsJSON = fmt.Sprintf(`,
  "channels": {
    "discord": {
      "token": "%s",
      "allowed_users": ["%s"]
    }
  }`, w.channelToken, w.channelUserID)
	}

	// SPRINT 2 ONBOARD: Generate agents.list with skills and subagents
	agentsListJSON := buildAgentListJSON(w.agentMode, w.agentTemplate, w.modelName, w.customSkills)

	// Generate config JSON
	data := fmt.Sprintf(`{
  "agents": {
    "defaults": {
      "workspace": "%s",
      "restrict_to_workspace": true,
      "model_name": "%s",
      "max_tokens": 8192,
      "temperature": 0.7,
      "max_tool_iterations": 20,
      "context_manager": "seahorse",
      "context_manager_config": {
        "context_threshold": 0.75,
        "fresh_tail_count": 16,
        "leaf_target_tokens": 1200,
        "condensed_target_tokens": 2000,
        "max_compact_iterations": 20
      }
    },
    "list": %s
  },
  "model_list": %s%s
}
`, w.workspace, w.modelName, agentsListJSON, modelListJSON, channelsJSON)

	return os.WriteFile(w.configPath, []byte(data), 0o600)
}

// checkGoVersion checks if Go is installed and returns version
func checkGoVersion() string {
	// Simplified check - in real implementation use exec.Command
	// For now, just return a placeholder
	return "Go 1.25.x (detected)"
}
