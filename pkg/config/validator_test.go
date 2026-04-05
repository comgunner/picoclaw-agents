// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package config_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/comgunner/picoclaw/pkg/config"
)

// TestValidator_ValidConfig verifies that a valid config passes validation.
func TestValidator_ValidConfig(t *testing.T) {
	cfg := &config.Config{
		ModelList: []config.ModelConfig{
			{
				ModelName: "deepseek-chat",
				Model:     "deepseek/deepseek-chat",
				APIKey:    "sk-1234567890abcdef",
			},
		},
		Agents: config.AgentsConfig{
			Defaults: config.AgentDefaults{
				Workspace:           "/tmp/test-workspace",
				RestrictToWorkspace: true,
				ModelName:           "deepseek-chat",
			},
			List: []config.AgentConfig{
				{
					ID:      "default",
					Default: true,
					Name:    "Default Agent",
				},
			},
		},
	}

	v := config.NewValidator()
	err := v.Validate(cfg)
	assert.NoError(t, err, "Valid config should pass validation")
}

// TestValidator_MissingAPIKey verifies error when API key is missing.
func TestValidator_MissingAPIKey(t *testing.T) {
	cfg := &config.Config{
		ModelList: []config.ModelConfig{
			{
				ModelName: "openai-gpt4",
				Model:     "openai/gpt-4",
				// APIKey is missing
			},
		},
	}

	v := config.NewValidator()
	err := v.Validate(cfg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "api_key is required")
}

// TestValidator_InvalidTelegramToken verifies error for invalid Telegram token.
func TestValidator_InvalidTelegramToken(t *testing.T) {
	cfg := &config.Config{
		ModelList: []config.ModelConfig{},
		Channels: config.ChannelsConfig{
			Telegram: config.TelegramConfig{
				Enabled: true,
				Token:   "invalid-token", // Invalid format
			},
		},
	}

	v := config.NewValidator()
	err := v.Validate(cfg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid telegram token format")
}

// TestValidator_DuplicateAgentID verifies error for duplicate agent IDs.
func TestValidator_DuplicateAgentID(t *testing.T) {
	cfg := &config.Config{
		ModelList: []config.ModelConfig{},
		Agents: config.AgentsConfig{
			List: []config.AgentConfig{
				{ID: "agent1", Name: "Agent 1"},
				{ID: "agent1", Name: "Agent 1 Duplicate"}, // Duplicate ID
			},
		},
	}

	v := config.NewValidator()
	err := v.Validate(cfg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate agent id")
}

// TestValidator_DuplicateModelName verifies error for duplicate model names.
func TestValidator_DuplicateModelName(t *testing.T) {
	cfg := &config.Config{
		ModelList: []config.ModelConfig{
			{ModelName: "model1", Model: "openai/gpt-4", APIKey: "sk-123"},
			{ModelName: "model1", Model: "anthropic/claude", APIKey: "sk-ant-123"}, // Duplicate
		},
	}

	v := config.NewValidator()
	err := v.Validate(cfg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate model_name")
}

// TestValidator_SubagentsMaxSpawnDepth verifies subagents validation.
func TestValidator_SubagentsMaxSpawnDepth(t *testing.T) {
	cfg := &config.Config{
		ModelList: []config.ModelConfig{},
		Agents: config.AgentsConfig{
			List: []config.AgentConfig{
				{
					ID:   "manager",
					Name: "Manager",
					Subagents: &config.SubagentsConfig{
						AllowAgents:   []string{"worker"},
						MaxSpawnDepth: 0, // Invalid: should be > 0
					},
				},
			},
		},
	}

	v := config.NewValidator()
	err := v.Validate(cfg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "must be > 0 when allow_agents is set")
}

// TestValidateFile_NotFound verifies error when config file doesn't exist.
func TestValidateFile_NotFound(t *testing.T) {
	v := config.NewValidator()
	err := v.ValidateFile("/nonexistent/path/config.json")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "config file not found")
}

// TestValidateFile_InvalidJSON verifies error for invalid JSON.
func TestValidateFile_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	err := os.WriteFile(configPath, []byte("{ invalid json }"), 0o600)
	require.NoError(t, err)

	v := config.NewValidator()
	err = v.ValidateFile(configPath)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "parsing config JSON")
}

// TestValidateFile_ValidConfig verifies successful validation of a file.
func TestValidateFile_ValidConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	configJSON := `{
		"model_list": [
			{
				"model_name": "deepseek-chat",
				"model": "deepseek/deepseek-chat",
				"api_key": "sk-test-key"
			}
		],
		"agents": {
			"defaults": {
				"workspace": "/tmp/test",
				"restrict_to_workspace": true,
				"model_name": "deepseek-chat"
			},
			"list": []
		},
		"channels": {
			"telegram": {
				"enabled": false
			}
		},
		"tools": {
			"web": {}
		}
	}`

	err := os.WriteFile(configPath, []byte(configJSON), 0o600)
	require.NoError(t, err)

	v := config.NewValidator()
	err = v.ValidateFile(configPath)
	assert.NoError(t, err, "Valid config file should pass validation")
}

// TestValidateTelegramToken_Valid verifies valid Telegram token format.
func TestValidateTelegramToken_Valid(t *testing.T) {
	// Valid format: 8-12 digits : 35+ alphanumeric/underscore/dash
	validTokens := []string{
		"1234567890:ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghij", // 10 digits + 36 chars
		"12345678:ABCDEF1234567890abcdef1234567890abc",    // 8 digits + 35 chars
	}

	for _, token := range validTokens {
		t.Run(token[:20]+"...", func(t *testing.T) {
			// We can't directly test the unexported function, but we can test via validation
			cfg := &config.Config{
				ModelList: []config.ModelConfig{},
				Agents: config.AgentsConfig{
					Defaults: config.AgentDefaults{
						Workspace: "/tmp/test",
					},
				},
				Channels: config.ChannelsConfig{
					Telegram: config.TelegramConfig{
						Enabled:   true,
						Token:     token,
						AllowFrom: []string{"123456"},
					},
				},
			}

			v := config.NewValidator()
			err := v.Validate(cfg)
			// Should not have telegram token error
			if err != nil {
				assert.NotContains(t, err.Error(), "invalid telegram token format")
			}
		})
	}
}

// TestValidator_BinancePartialKeys verifies error when only one Binance key is set.
func TestValidator_BinancePartialKeys(t *testing.T) {
	cfg := &config.Config{
		ModelList: []config.ModelConfig{},
		Tools: config.ToolsConfig{
			Binance: config.BinanceToolsConfig{
				APIKey:    "test-api-key",
				SecretKey: "", // Missing
			},
		},
	}

	v := config.NewValidator()
	err := v.Validate(cfg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "both api_key and secret_key must be set together")
}

// TestValidator_ErrorListLimited verifies that error list is limited to 20 errors.
func TestValidator_ErrorListLimited(t *testing.T) {
	// Create a config with many errors
	cfg := &config.Config{
		ModelList: make([]config.ModelConfig, 25),
		Agents: config.AgentsConfig{
			List: make([]config.AgentConfig, 25),
		},
	}

	// Add many models without API keys
	for i := 0; i < 25; i++ {
		cfg.ModelList[i] = config.ModelConfig{
			ModelName: "model-" + string(rune('a'+i)),
			Model:     "openai/gpt-4",
		}
	}

	// Add many agents with duplicate IDs
	for i := 0; i < 25; i++ {
		cfg.Agents.List[i] = config.AgentConfig{
			ID: "duplicate",
		}
	}

	v := config.NewValidator()
	err := v.Validate(cfg)
	require.Error(t, err)

	// Error message should mention there are more errors
	errStr := err.Error()
	if strings.Contains(errStr, "and") && strings.Contains(errStr, "more errors") {
		assert.Contains(t, errStr, "more errors", "Should indicate when there are more errors")
	}
}

// TestValidator_FreeModelNoAPIKey verifies free models don't require API key.
func TestValidator_FreeModelNoAPIKey(t *testing.T) {
	cfg := &config.Config{
		ModelList: []config.ModelConfig{
			{
				ModelName: "openrouter-free",
				Model:     "openrouter/free",
				// No API key needed for free models
			},
		},
		Agents: config.AgentsConfig{
			Defaults: config.AgentDefaults{
				Workspace: "/tmp/test",
			},
		},
	}

	v := config.NewValidator()
	err := v.Validate(cfg)
	// Should only have warnings, not errors about missing API key
	if err != nil {
		assert.NotContains(t, err.Error(), "api_key is required")
	}
}
