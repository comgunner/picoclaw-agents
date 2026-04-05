// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)

package security_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/comgunner/picoclaw/pkg/security"
)

func TestSanitizeAPIKeys(t *testing.T) {
	input := "The API key is sk-X9kL2mNpQ4rS7tUv0WxY3zA6bC8dE1fG5hI9jK2lM4nO7pR for this request"
	result := security.Sanitize(input)
	assert.NotContains(t, result, "sk-X9kL2mNpQ4rS7tUv0WxY3zA6bC8dE1fG5hI9jK2lM4nO7pR")
	assert.Contains(t, result, "[REDACTED_")
}

func TestSanitizeAnthropicKey(t *testing.T) {
	input := "Using sk-ant-api03-abcdefghijklmnopqrstuvwxyz12345678901234 for Claude"
	result := security.Sanitize(input)
	assert.NotContains(t, result, "sk-ant-api03")
	assert.Contains(t, result, "[REDACTED_anthropic_key]")
}

func TestSanitizeJWT(t *testing.T) {
	jwt := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	result := security.Sanitize(jwt)
	assert.NotContains(t, result, jwt)
	assert.Contains(t, result, "[REDACTED_jwt]")
}

func TestSanitizeGitHubToken(t *testing.T) {
	input := "Token: ghp_abcdefghijklmnopqrstuvwxyz1234567890"
	result := security.Sanitize(input)
	assert.NotContains(t, result, "ghp_")
	assert.Contains(t, result, "[REDACTED_github_token]")
}

func TestSanitizeMultipleSecrets(t *testing.T) {
	input := `Config:
  openai: sk-X9kL2mNpQ4rS7tUv0WxY3zA6bC8dE1fG5hI9jK2lM4nO7pR
  anthropic: sk-ant-api03-abcdefghijklmnopqrstuvwxyz12345678901234
  jwt: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.abc123
`
	result := security.Sanitize(input)
	assert.NotContains(t, result, "sk-X9kL2mNpQ4rS7tUv0WxY3zA6bC8dE1fG5hI9jK2lM4nO7pR")
	assert.NotContains(t, result, "sk-ant-api03")
	assert.NotContains(t, result, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9")
	assert.Contains(t, result, "[REDACTED_openai_key]")
	assert.Contains(t, result, "[REDACTED_anthropic_key]")
	assert.Contains(t, result, "[REDACTED_jwt]")
}

func TestSanitizePreservesNonSecrets(t *testing.T) {
	input := "Hello, the response was successful and no secrets here"
	result := security.Sanitize(input)
	assert.Equal(t, input, result)
}

func TestSanitizeEmptyString(t *testing.T) {
	result := security.Sanitize("")
	assert.Equal(t, "", result)
}

func TestSanitizeMap(t *testing.T) {
	m := map[string]any{
		"message": "key is sk-X9kL2mNpQ4rS7tUv0WxY3zA6bC8dE1fG5hI9jK2lM4nO7pR",
		"count":   42,
		"nested": map[string]any{
			"token": "sk-ant-api03-abcdefghijklmnopqrstuvwxyz12345",
		},
	}

	result := security.SanitizeMap(m)
	assert.NotContains(t, result["message"], "sk-X9kL2mNpQ4rS7tUv0WxY3zA6bC8dE1fG5hI9jK2lM4nO7pR")
	assert.Equal(t, 42, result["count"]) // int sin cambios
	nested, ok := result["nested"].(map[string]any)
	assert.True(t, ok)
	assert.NotContains(t, nested["token"], "sk-ant-api03")
}

func TestSanitizeMapNil(t *testing.T) {
	result := security.SanitizeMap(nil)
	assert.Nil(t, result)
}

func TestSanitizeMapEmpty(t *testing.T) {
	m := map[string]any{}
	result := security.SanitizeMap(m)
	assert.NotNil(t, result)
	assert.Empty(t, result)
}

func TestSanitizeMapWithSlice(t *testing.T) {
	m := map[string]any{
		"items": []any{"value1", "sk-abcdefghijklmnopqrstuvwxyz123456", 123},
	}
	result := security.SanitizeMap(m)
	// Los slices no se sanitizan recursivamente en esta implementación
	assert.Equal(t, m["items"], result["items"])
}

func TestSanitizeMapNestedMaps(t *testing.T) {
	m := map[string]any{
		"level1": map[string]any{
			"level2": map[string]any{
				"secret": "sk-X9kL2mNpQ4rS7tUv0WxY3zA6bC8dE1fG5hI9jK2lM4nO7pR",
				"normal": "value",
			},
		},
	}

	result := security.SanitizeMap(m)
	level1, ok := result["level1"].(map[string]any)
	assert.True(t, ok)
	level2, ok := level1["level2"].(map[string]any)
	assert.True(t, ok)
	assert.NotContains(t, level2["secret"], "sk-X9kL2mNpQ4rS7tUv0WxY3zA6bC8dE1fG5hI9jK2lM4nO7pR")
	assert.Contains(t, level2["secret"], "[REDACTED_")
	assert.Equal(t, "value", level2["normal"])
}

func TestSanitizeTelegramBotToken(t *testing.T) {
	input := "Bot token: 1234567890:ABCdefGHIjklMNOpqrsTUVwxyz123456789"
	result := security.Sanitize(input)
	assert.NotContains(t, result, "1234567890:ABCdefGHIjklMNOpqrsTUVwxyz123456789")
	assert.Contains(t, result, "[REDACTED_telegram_bot]")
}

func TestSanitizeAWSCredentials(t *testing.T) {
	input := `AWS config:
  access_key: AKIAIOSFODNN7EXAMPLE
  secret: wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
`
	result := security.Sanitize(input)
	assert.NotContains(t, result, "AKIAIOSFODNN7EXAMPLE")
	assert.Contains(t, result, "[REDACTED_aws_access_key]")
}

func TestSanitizeSlackToken(t *testing.T) {
	input := "Slack token: " + "xoxb-" + "000000000000-0000000000000-XXXXXXXXXXXXXXXXXXXXXXXXXXXX"
	result := security.Sanitize(input)
	assert.NotContains(t, result, "xoxb-")
	assert.Contains(t, result, "[REDACTED_slack_token]")
}

func TestSanitizeStripeKey(t *testing.T) {
	input := "Stripe: " + "sk_live_" + "xxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	result := security.Sanitize(input)
	assert.NotContains(t, result, "sk_live_")
	assert.Contains(t, result, "[REDACTED_stripe_secret]")
}

func TestSanitizeDeepSeekKey(t *testing.T) {
	// DeepSeek keys son solo hex (a-f, 0-9)
	input := "DeepSeek: sk-abcdef1234567890abcdef1234567890ab"
	result := security.Sanitize(input)
	assert.NotContains(t, result, "sk-abcdef1234567890")
	assert.Contains(t, result, "[REDACTED_deepseek_key]")
}

func TestSanitizeGoogleAPIKey(t *testing.T) {
	input := "Google: AIzaSyDaGmWKa4JsXZ-HjGw7ISLn_3namBGewQe"
	result := security.Sanitize(input)
	assert.NotContains(t, result, "AIzaSy")
	assert.Contains(t, result, "[REDACTED_google_api_key]")
}
