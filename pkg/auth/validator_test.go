// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package auth

import (
	"context"
	"testing"
	"time"
)

func TestValidateAPIKeyFormat_OpenAI(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		wantError bool
	}{
		{"valid key", "sk-projabcdefghijklmnopqrstuvwxyz1234567890", false},
		{"valid key short", "sk-abcdefghijklmnopqrstuvwxyz", false},
		{"invalid prefix", "pk-abcdefghijklmnopqrstuvwxyz1234567890", true},
		{"too short", "sk-abc123", true},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAPIKeyFormat("openai", tt.key)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateAPIKeyFormat() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateAPIKeyFormat_Anthropic(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		wantError bool
	}{
		{"valid key", "sk-ant-api03-abcdefghijklmnopqrstuvwxyz1234567890", false},
		{"valid key short", "sk-ant-abcdefghijklmnopqrstuvwxyz", false},
		{"invalid prefix", "sk-ant-abc", true},
		{"too short", "sk-ant-", true},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAPIKeyFormat("anthropic", tt.key)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateAPIKeyFormat() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateAPIKeyFormat_Groq(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		wantError bool
	}{
		{"valid key", "gsk_abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", false},
		{"valid key 50chars", "gsk_abcdefghijklmnopqrstuvwxyz1234567890ABCDEFGHIJKLMN", false},
		{"too short", "gsk_abc123", true},
		{"invalid prefix", "gsk_abc", true},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAPIKeyFormat("groq", tt.key)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateAPIKeyFormat() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateAPIKeyFormat_OpenRouter(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		wantError bool
	}{
		{"valid key", "sk-or-abcdefghijklmnopqrstuvwxyz1234567890", false},
		{"valid key with dash", "sk-or-abc-def-ghi-jkl-mno-pqr-stu-vwx-yz", false},
		{"too short", "sk-or-abc", true},
		{"invalid prefix", "sk-or-", true},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAPIKeyFormat("openrouter", tt.key)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateAPIKeyFormat() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateAPIKeyFormat_Gemini(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		wantError bool
	}{
		{"valid key", "AIzaSyabcdefghijklmnopqrstuvwxyz1234567", false},
		{"valid key with dash", "AIzaSyA-BCDEFGHIJKLMNOPQRSTUVWXYZ123456", false},
		{"invalid prefix", "AIzbSyabcdefghijklmnopqrstuvwxyz1234567", true},
		{"too short", "AIza123", true},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAPIKeyFormat("gemini", tt.key)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateAPIKeyFormat() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateAPIKeyFormat_GitHubCopilot(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		wantError bool
	}{
		{"valid ghu", "ghu_abcdefghijklmnopqrstuvwxyz0123456789", false},
		{"valid gho", "gho_abcdefghijklmnopqrstuvwxyz0123456789", false},
		{"valid ghs", "ghs_abcdefghijklmnopqrstuvwxyz0123456789", false},
		{"valid ghr", "ghr_abcdefghijklmnopqrstuvwxyz0123456789", false},
		{"invalid prefix", "ghx_abcdefghijklmnopqrstuvwxyz0123456789", true},
		{"too short", "ghu_abc", true},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAPIKeyFormat("github-copilot", tt.key)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateAPIKeyFormat() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestQuickValidate(t *testing.T) {
	tests := []struct {
		name      string
		provider  string
		key       string
		wantValid bool
	}{
		{"valid openai", "openai", "sk-projabcdefghijklmnopqrstuvwxyz1234567890", true},
		{"invalid openai", "openai", "invalid-key", false},
		{"valid anthropic", "anthropic", "sk-ant-api03-abcdefghijklmnopqrstuvwxyz", true},
		{"valid groq", "groq", "gsk_abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456", true},
		{"unknown provider", "unknown", "any-key-1234567890", true}, // Unknown providers pass format check
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := QuickValidate(tt.provider, tt.key)
			if result.Valid != tt.wantValid {
				t.Errorf("QuickValidate() valid = %v, want %v", result.Valid, tt.wantValid)
			}
		})
	}
}

func TestGetProviderFromModelName(t *testing.T) {
	tests := []struct {
		name      string
		modelName string
		want      string
	}{
		{"openai gpt-4", "openai/gpt-4", "openai"},
		{"anthropic claude", "anthropic/claude-3", "anthropic"},
		{"deepseek chat", "deepseek/deepseek-chat", "deepseek"},
		{"gemini flash", "gemini/gemini-2.0-flash", "gemini"},
		{"groq llama", "groq/llama-3-70b", "groq"},
		{"openrouter auto", "openrouter/auto", "openrouter"},
		{"github copilot", "github-copilot/gpt-4", "github-copilot"},
		{"claude alias", "claude/claude-3-sonnet", "anthropic"},
		{"gpt alias", "gpt/gpt-4-turbo", "openai"},
		{"no slash", "gpt-4", ""},
		{"empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetProviderFromModelName(tt.modelName)
			if got != tt.want {
				t.Errorf("GetProviderFromModelName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateAllKeys(t *testing.T) {
	// Skip online validation in CI/CD unless explicitly enabled
	if testing.Short() {
		t.Skip("Skipping online validation test")
	}

	providers := map[string]string{
		"openai":    "sk-proj-invalid-key-for-testing",
		"anthropic": "sk-ant-invalid-key-for-testing",
		"groq":      "gsk_invalidkeyfortesting1234567890abcdefghijklmnopqrstuvwxyz",
	}

	results := ValidateAllKeys(providers)

	if len(results) != len(providers) {
		t.Errorf("ValidateAllKeys() returned %d results, expected %d", len(results), len(providers))
	}

	// All should be invalid since we're using fake keys
	for _, result := range results {
		if result.Valid {
			t.Errorf("ValidateAllKeys() provider %s should be invalid with fake key", result.Provider)
		}
	}
}

func TestValidateAPIKeyOnline_Timeout(t *testing.T) {
	// Test timeout behavior with a very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	result := ValidateAPIKeyOnline(ctx, "openai", "sk-proj-testkey1234567890abcdefghijklmnop")

	// Should either timeout or fail with connection error
	if result.Valid {
		t.Error("ValidateAPIKeyOnline() should not validate successfully with 1ms timeout")
	}
}

func TestAPIKeyFormat_Comprehensive(t *testing.T) {
	// Test all providers with valid and invalid keys
	providers := map[string]struct {
		valid   string
		invalid string
	}{
		"openai": {
			valid:   "sk-projabcdefghijklmnopqrstuvwxyz1234567890",
			invalid: "pk-invalid-key",
		},
		"anthropic": {
			valid:   "sk-ant-api03-abcdefghijklmnopqrstuvwxyz",
			invalid: "sk-invalid-key",
		},
		"deepseek": {
			valid:   "sk-abcdefghijklmnopqrstuvwxyz123456",
			invalid: "ds-invalid-key",
		},
		"zhipu": {
			valid:   "abc123def456ghi789jkl012mno345pqr678",
			invalid: "zhipu",
		},
		"gemini": {
			valid:   "AIzaSyabcdefghijklmnopqrstuvwxyz1234567",
			invalid: "AIzb-invalid-key",
		},
		"groq": {
			valid:   "gsk_abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456",
			invalid: "groq-invalid",
		},
		"openrouter": {
			valid:   "sk-or-abcdefghijklmnopqrstuvwxyz123456",
			invalid: "or-invalid",
		},
		"qwen": {
			valid:   "abcdefghijklmnopqrstuvwxyz123456",
			invalid: "qw",
		},
		"github-copilot": {
			valid:   "ghu_abcdefghijklmnopqrstuvwxyz0123456789",
			invalid: "ghx-invalid-key",
		},
	}

	for provider, keys := range providers {
		t.Run(provider+"_valid", func(t *testing.T) {
			err := ValidateAPIKeyFormat(provider, keys.valid)
			if err != nil {
				t.Errorf("ValidateAPIKeyFormat(%s, valid) error = %v", provider, err)
			}
		})

		t.Run(provider+"_invalid", func(t *testing.T) {
			err := ValidateAPIKeyFormat(provider, keys.invalid)
			if err == nil {
				t.Errorf("ValidateAPIKeyFormat(%s, invalid) expected error, got nil", provider)
			}
		})
	}
}
