// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package config

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Validator validates semantically a parsed Config.
type Validator struct {
	errors []ValidationError
}

// ValidationError represents a single validation error.
type ValidationError struct {
	Field   string
	Message string
}

// NewValidator creates a new Validator instance.
func NewValidator() *Validator {
	return &Validator{
		errors: make([]ValidationError, 0),
	}
}

// Validate executes all validations and returns error if there are failures.
// Accumulates all errors (not fail-fast) for better UX.
func (v *Validator) Validate(cfg *Config) error {
	v.errors = make([]ValidationError, 0)

	// Validate model_list
	v.validateModelList(cfg)

	// Validate agents
	v.validateAgents(cfg)

	// Validate channels
	v.validateChannels(cfg)

	// Validate tools
	v.validateTools(cfg)

	if len(v.errors) > 0 {
		return &MultiValidationError{errors: v.errors}
	}

	return nil
}

// ValidateFile reads and validates a config.json file.
func (v *Validator) ValidateFile(path string) error {
	// Read file
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("config file not found: %s", path)
		}
		return fmt.Errorf("reading config file: %w", err)
	}

	// Parse JSON
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("parsing config JSON: %w", err)
	}

	// Store config path
	cfg.configPath = path

	// Validate semantically
	return v.Validate(&cfg)
}

// validateModelList validates the model_list configuration.
func (v *Validator) validateModelList(cfg *Config) {
	modelNames := make(map[string]bool)

	for i, model := range cfg.ModelList {
		field := fmt.Sprintf("model_list[%d]", i)

		// Check model_name is present
		if model.ModelName == "" {
			v.addError(field+".model_name", "model_name is required")
			continue
		}

		// Check for duplicate model_names
		if modelNames[model.ModelName] {
			v.addError(field+".model_name", "duplicate model_name: "+model.ModelName)
		}
		modelNames[model.ModelName] = true

		// Check model reference format (provider/model)
		if model.Model != "" && !strings.Contains(model.Model, "/") {
			v.addWarning(field+".model", "model should be in format 'provider/model': "+model.Model)
		}

		// Check API key if model requires it
		if model.APIKey == "" && !isFreeModel(model.Model) {
			v.addError(field+".api_key", "api_key is required for non-free models")
		}
	}
}

// validateAgents validates agent configurations.
func (v *Validator) validateAgents(cfg *Config) {
	agentIDs := make(map[string]bool)

	// Validate defaults
	if cfg.Agents.Defaults.Workspace == "" {
		v.addWarning("agents.defaults.workspace", "workspace not set, using default")
	}

	// Validate agent list
	for i, agent := range cfg.Agents.List {
		field := fmt.Sprintf("agents.list[%d]", i)

		// Check ID is present
		if agent.ID == "" {
			v.addError(field+".id", "agent id is required")
			continue
		}

		// Check for duplicate IDs
		if agentIDs[agent.ID] {
			v.addError(field+".id", "duplicate agent id: "+agent.ID)
		}
		agentIDs[agent.ID] = true

		// Validate subagents config
		if agent.Subagents != nil {
			subField := field + ".subagents"

			// If allow_agents is not empty, max_spawn_depth should be > 0
			if len(agent.Subagents.AllowAgents) > 0 && agent.Subagents.MaxSpawnDepth <= 0 {
				v.addError(subField+".max_spawn_depth", "must be > 0 when allow_agents is set")
			}

			// Validate allow_agents references exist
			for _, allowedID := range agent.Subagents.AllowAgents {
				if allowedID != "*" && !agentIDs[allowedID] {
					v.addWarning(subField+".allow_agents", "references unknown agent: "+allowedID)
				}
			}
		}

		// Validate model reference
		if agent.Model != nil && agent.Model.Primary != "" {
			if !v.modelExists(agent.Model.Primary, cfg.ModelList) {
				v.addWarning(field+".model", "model not found in model_list: "+agent.Model.Primary)
			}
		}
	}
}

// validateChannels validates channel configurations.
func (v *Validator) validateChannels(cfg *Config) {
	// Validate Telegram
	if cfg.Channels.Telegram.Enabled {
		field := "channels.telegram"
		if cfg.Channels.Telegram.Token == "" {
			v.addError(field+".token", "token is required when telegram is enabled")
		} else if !validateTelegramToken(cfg.Channels.Telegram.Token) {
			v.addError(field+".token", "invalid telegram token format")
		}
		if len(cfg.Channels.Telegram.AllowFrom) == 0 {
			v.addWarning(field+".allow_from", "no user IDs in allow_from - bot won't respond to anyone")
		}
	}

	// Validate Discord
	if cfg.Channels.Discord.Enabled {
		field := "channels.discord"
		if cfg.Channels.Discord.Token == "" {
			v.addError(field+".token", "token is required when discord is enabled")
		}
	}

	// Validate other channels similarly...
}

// validateTools validates tool configurations.
func (v *Validator) validateTools(cfg *Config) {
	// Validate web search tools
	if cfg.Tools.Web.Brave.Enabled && cfg.Tools.Web.Brave.APIKey == "" {
		v.addError("tools.web.brave.api_key", "api_key required when brave search is enabled")
	}

	if cfg.Tools.Web.Tavily.Enabled && cfg.Tools.Web.Tavily.APIKey == "" {
		v.addError("tools.web.tavily.api_key", "api_key required when tavily is enabled")
	}

	// Validate image generation
	if cfg.Tools.ImageGen.Provider == "ideogram" && cfg.Tools.ImageGen.IdeogramAPIKey == "" {
		v.addError("tools.image_gen.ideogram_api_key", "api_key required for ideogram provider")
	}

	// Validate Binance (check if API key is set when secret is also set, or vice versa)
	if (cfg.Tools.Binance.APIKey != "" && cfg.Tools.Binance.SecretKey == "") ||
		(cfg.Tools.Binance.APIKey == "" && cfg.Tools.Binance.SecretKey != "") {
		v.addError("tools.binance", "both api_key and secret_key must be set together")
	}
}

// modelExists checks if a model_name exists in the model_list.
func (v *Validator) modelExists(modelName string, modelList []ModelConfig) bool {
	for _, m := range modelList {
		if m.ModelName == modelName {
			return true
		}
	}
	return false
}

// isFreeModel checks if a model is free and doesn't require an API key.
func isFreeModel(model string) bool {
	freeModels := []string{
		"openrouter/free",
		"free",
	}
	for _, free := range freeModels {
		if strings.Contains(model, free) {
			return true
		}
	}
	return false
}

// validateTelegramToken validates the format of a Telegram bot token.
// Format: digits:alphanumeric_underscore_dash (e.g., 123456789:ABCdefGHIjklMNOpqrsTUVwxyz)
func validateTelegramToken(token string) bool {
	// Telegram token format: numeric_part:alphanumeric_part
	// Numeric part: 8-12 digits
	// Alphanumeric part: 35+ characters including letters, numbers, underscore, dash
	pattern := regexp.MustCompile(`^\d{8,12}:[a-zA-Z0-9_-]{35,}$`)
	return pattern.MatchString(token)
}

// addError adds a validation error.
func (v *Validator) addError(field, message string) {
	v.errors = append(v.errors, ValidationError{
		Field:   field,
		Message: message,
	})
}

// addWarning adds a validation warning (non-blocking).
func (v *Validator) addWarning(field, message string) {
	// Warnings are stored but don't cause validation to fail
	// Could be logged or displayed separately
	v.errors = append(v.errors, ValidationError{
		Field:   field,
		Message: "WARNING: " + message,
	})
}

// MultiValidationError is a collection of validation errors.
type MultiValidationError struct {
	errors []ValidationError
}

// Error implements the error interface.
func (e *MultiValidationError) Error() string {
	var sb strings.Builder
	sb.WriteString("config validation failed with ")
	sb.WriteString(fmt.Sprintf("%d", len(e.errors)))
	sb.WriteString(" error(s):\n")

	for i, err := range e.errors {
		isWarning := strings.Contains(err.Message, "WARNING:")
		prefix := "  - "
		if isWarning {
			prefix = "  ⚠️  "
		}
		sb.WriteString(fmt.Sprintf("%s[%s] %s\n", prefix, err.Field, strings.TrimPrefix(err.Message, "WARNING: ")))

		// Limit output to first 20 errors
		if i >= 19 && len(e.errors) > 20 {
			sb.WriteString(fmt.Sprintf("\n... and %d more errors", len(e.errors)-20))
			break
		}
	}

	return sb.String()
}
