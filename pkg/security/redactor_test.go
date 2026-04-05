// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package security

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecretRedactor_RegexPatterns(t *testing.T) {
	redactor := NewSecretRedactor()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "No secrets",
			input:    "This is a normal string.",
			expected: "This is a normal string.",
		},
		{
			name:     "OpenAI key",
			input:    "Here is my key: sk-abcdefghijklmnopqrstuvwxyz0123456789101112",
			expected: "Here is my key: [REDACTED_SECRET]",
		},
		{
			name:     "Anthropic key",
			input:    "API_KEY=sk-ant-api03-abcdefghijklmnop-qrstuvwxyz_0123456789",
			expected: "API_KEY=[REDACTED_SECRET]",
		},
		{
			name:     "Telegram Bot Token",
			input:    "TELEGRAM_BOT_TOKEN=1234567890:ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghi",
			expected: "TELEGRAM_BOT_TOKEN=[REDACTED_SECRET]",
		},
		{
			name:     "Facebook Token",
			input:    "Access via EAAK1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ",
			expected: "Access via [REDACTED_SECRET]",
		},
		{
			name:     "Private Key",
			input:    "Here is my key:\n-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQE... (lots of base64) ...xyz\n-----END RSA PRIVATE KEY-----\nKeep it safe.",
			expected: "Here is my key:\n[REDACTED_SECRET]\nKeep it safe.",
		},
		{
			name:     "GitHub Token",
			input:    "export GITHUB_TOKEN=ghp_abcdefghijklmnopqrstuvwxyz0123456789",
			expected: "export GITHUB_TOKEN=[REDACTED_SECRET]",
		},
		{
			name:     "Slack Token",
			input:    "curl -H 'Authorization: Bearer xoxb-1234-5678-abcdef' https://slack.com",
			expected: "curl -H 'Authorization: Bearer [REDACTED_SECRET]' https://slack.com",
		},
		{
			name:     "JWT Token",
			input:    "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
			expected: "Authorization: Bearer [REDACTED_SECRET]",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := redactor.Redact(tc.input)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestSecretRedactor_ExactMatch(t *testing.T) {
	redactor := NewSecretRedactor()

	customKey := "custom_secret_key_12345"
	redactor.AddSecret(customKey)

	shortKey := "short" // Should be ignored (len < 8)
	redactor.AddSecret(shortKey)

	input := "Connecting with custom_secret_key_12345 and short password."
	expected := "Connecting with [REDACTED_SECRET] and short password."

	actual := redactor.Redact(input)
	assert.Equal(t, expected, actual)
	assert.True(t, strings.Contains(actual, "short")) // Verify short was NOT redacted
}

func TestSecretRedactor_GlobalInstance(t *testing.T) {
	input := "Bearer sk-1234567890123456789012345678901234567890"
	expected := "Bearer [REDACTED_SECRET]"

	actual := GlobalRedactor.Redact(input)
	assert.Equal(t, expected, actual)
}
