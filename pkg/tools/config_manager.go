// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/comgunner/picoclaw/pkg/logger"
)

// ConfigManagerTool provides configuration read and validation capabilities.
type ConfigManagerTool struct {
	configPath string
}

// NewConfigManagerTool creates a new ConfigManagerTool instance.
func NewConfigManagerTool(configPath string) *ConfigManagerTool {
	return &ConfigManagerTool{
		configPath: configPath,
	}
}

// Name returns the tool name.
func (t *ConfigManagerTool) Name() string {
	return "config_manager"
}

// Description returns the tool description.
func (t *ConfigManagerTool) Description() string {
	return "Read, validate, and backup system configuration (config.json). Use action='read' to view config, 'validate' to check it, 'backup' to create a backup, or 'get' to retrieve a specific key."
}

// Parameters returns the JSON schema for tool parameters.
func (t *ConfigManagerTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"action": map[string]any{
				"type":        "string",
				"description": "Action to perform: 'read', 'validate', 'backup', or 'get'",
				"enum":        []string{"read", "validate", "backup", "get"},
			},
			"key": map[string]any{
				"type":        "string",
				"description": "Configuration key to retrieve (for action='get'). Use dot notation for nested keys (e.g., 'agents.defaults.model_name')",
			},
		},
		"required": []string{"action"},
	}
}

// Execute runs the configuration manager tool.
func (t *ConfigManagerTool) Execute(ctx context.Context, args map[string]any) *ToolResult {
	action, ok := args["action"].(string)
	if !ok {
		return ErrorResult("action is required and must be one of: read, validate, backup, get")
	}

	switch action {
	case "read":
		return t.readConfig()
	case "validate":
		return t.validateConfig()
	case "backup":
		return t.backupConfig()
	case "get":
		key, _ := args["key"].(string)
		return t.getConfigKey(key)
	default:
		return ErrorResult(fmt.Sprintf("unknown action: %s. Valid options: read, validate, backup, get", action))
	}
}

// readConfig reads the entire configuration file.
func (t *ConfigManagerTool) readConfig() *ToolResult {
	data, err := os.ReadFile(t.configPath)
	if err != nil {
		logger.ErrorCF("tool", "Failed to read config",
			map[string]any{
				"tool":  "config_manager",
				"error": err.Error(),
				"path":  t.configPath,
			})
		return ErrorResult(fmt.Sprintf("failed to read config: %v", err))
	}

	var cfg map[string]any
	if err := json.Unmarshal(data, &cfg); err != nil {
		logger.ErrorCF("tool", "Failed to parse config JSON",
			map[string]any{
				"tool":  "config_manager",
				"error": err.Error(),
			})
		return ErrorResult(fmt.Sprintf("failed to parse config JSON: %v", err))
	}

	// Security: Remove sensitive fields before returning
	sanitized := sanitizeConfig(cfg)

	result := map[string]any{
		"config": sanitized,
		"size":   len(data),
		"path":   t.configPath,
	}

	_ = result // For future structured output
	return SilentResult(fmt.Sprintf("Configuration loaded from %s (%d bytes)", t.configPath, len(data)))
}

// validateConfig validates the configuration structure.
func (t *ConfigManagerTool) validateConfig() *ToolResult {
	data, err := os.ReadFile(t.configPath)
	if err != nil {
		return ErrorResult(fmt.Sprintf("failed to read config: %v", err))
	}

	var cfg map[string]any
	if err := json.Unmarshal(data, &cfg); err != nil {
		return ErrorResult(fmt.Sprintf("invalid config JSON: %v", err))
	}

	// Validate required fields
	var issues []string

	// Check agents configuration
	if agents, ok := cfg["agents"].(map[string]any); ok {
		if defaults, ok := agents["defaults"].(map[string]any); ok {
			if _, hasModel := defaults["model_name"]; !hasModel {
				if _, hasProvider := defaults["provider"]; !hasProvider {
					issues = append(issues, "agents.defaults.model_name or agents.defaults.provider is required")
				}
			}
		}
	}

	// Check model_list configuration
	if modelList, ok := cfg["model_list"].([]any); ok {
		if len(modelList) == 0 {
			issues = append(issues, "model_list is empty - at least one model should be configured")
		}
	} else {
		issues = append(issues, "model_list is required - configure at least one LLM provider")
	}

	if len(issues) > 0 {
		return SilentResult(fmt.Sprintf("Configuration validation warnings:\n- %s", strings.Join(issues, "\n- ")))
	}

	return SilentResult("Configuration is valid")
}

// backupConfig creates a timestamped backup of the configuration.
func (t *ConfigManagerTool) backupConfig() *ToolResult {
	data, err := os.ReadFile(t.configPath)
	if err != nil {
		return ErrorResult(fmt.Sprintf("failed to read config: %v", err))
	}

	// Create backup path with timestamp
	dir := filepath.Dir(t.configPath)
	timestamp := time.Now().Format("20060102_150405")
	backupPath := filepath.Join(dir, fmt.Sprintf("config.backup.%s.json", timestamp))

	if err := os.WriteFile(backupPath, data, 0o644); err != nil {
		logger.ErrorCF("tool", "Failed to create config backup",
			map[string]any{
				"tool":  "config_manager",
				"error": err.Error(),
				"path":  backupPath,
			})
		return ErrorResult(fmt.Sprintf("failed to create backup: %v", err))
	}

	result := map[string]any{
		"backup_path": backupPath,
		"size":        len(data),
		"timestamp":   timestamp,
	}

	_ = result // For future structured output
	return SilentResult(fmt.Sprintf("Configuration backed up to %s", backupPath))
}

// getConfigKey retrieves a specific configuration value by key.
func (t *ConfigManagerTool) getConfigKey(key string) *ToolResult {
	if key == "" {
		return ErrorResult("key is required for action='get'. Use dot notation (e.g., 'agents.defaults.model_name')")
	}

	data, err := os.ReadFile(t.configPath)
	if err != nil {
		return ErrorResult(fmt.Sprintf("failed to read config: %v", err))
	}

	var cfg map[string]any
	if err := json.Unmarshal(data, &cfg); err != nil {
		return ErrorResult(fmt.Sprintf("failed to parse config JSON: %v", err))
	}

	// Navigate nested keys using dot notation
	value := navigateNestedKey(cfg, key)
	if value == nil {
		return ErrorResult(fmt.Sprintf("configuration key '%s' not found", key))
	}

	// Sanitize sensitive data
	sanitized := sanitizeValue(value)

	result := map[string]any{
		"key":   key,
		"value": sanitized,
	}

	_ = result // For future structured output
	return SilentResult(fmt.Sprintf("Configuration '%s': %v", key, sanitized))
}

// sanitizeConfig removes sensitive fields from configuration.
func sanitizeConfig(cfg map[string]any) map[string]any {
	// Create a copy to avoid modifying original
	sanitized := make(map[string]any)
	for k, v := range cfg {
		// Skip sensitive keys
		if isSensitiveKey(k) {
			sanitized[k] = "[REDACTED]"
			continue
		}
		if nested, ok := v.(map[string]any); ok {
			sanitized[k] = sanitizeConfig(nested)
		} else if arr, ok := v.([]any); ok {
			sanitizedArr := make([]any, 0, len(arr))
			for _, item := range arr {
				if itemMap, ok := item.(map[string]any); ok {
					sanitizedArr = append(sanitizedArr, sanitizeConfig(itemMap))
				} else {
					sanitizedArr = append(sanitizedArr, item)
				}
			}
			sanitized[k] = sanitizedArr
		} else {
			sanitized[k] = v
		}
	}
	return sanitized
}

// isSensitiveKey checks if a key contains sensitive data.
func isSensitiveKey(key string) bool {
	sensitive := []string{"api_key", "secret", "password", "token", "credential", "private_key"}
	keyLower := strings.ToLower(key)
	for _, s := range sensitive {
		if strings.Contains(keyLower, s) {
			return true
		}
	}
	return false
}

// navigateNestedKey navigates a nested map using dot notation.
func navigateNestedKey(m map[string]any, key string) any {
	parts := strings.Split(key, ".")
	var current any = m

	for _, part := range parts {
		if mp, ok := current.(map[string]any); ok {
			current = mp[part]
			if current == nil {
				return nil
			}
		} else {
			return nil
		}
	}

	return current
}

// sanitizeValue sanitizes a single value for display.
func sanitizeValue(v any) any {
	if s, ok := v.(string); ok {
		// Check if it looks like a sensitive value
		if len(s) > 20 && (strings.Contains(s, "sk-") || strings.Contains(s, "key_") || strings.Contains(s, "token_")) {
			return "[REDACTED]"
		}
		// Truncate long strings
		if len(s) > 200 {
			return s[:200] + "...[truncated]"
		}
	}
	return v
}
