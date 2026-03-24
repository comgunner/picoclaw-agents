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
	"regexp"
	"strings"
	"sync"
)

// SecretRedactor is responsible for scanning strings and censoring known
// secret patterns (like API keys, tokens) to prevent leaks in logs or chat.
type SecretRedactor struct {
	mu            sync.RWMutex
	exactSecrets  map[string]struct{}
	regexPatterns []*regexp.Regexp
}

// GlobalRedactor is a singleton instance of SecretRedactor for easy access
var GlobalRedactor = NewSecretRedactor()

const redactedPlaceholder = "[REDACTED_SECRET]"

// NewSecretRedactor initializes a new SecretRedactor with common regex patterns.
func NewSecretRedactor() *SecretRedactor {
	r := &SecretRedactor{
		exactSecrets: make(map[string]struct{}),
		regexPatterns: []*regexp.Regexp{
			// OpenAI
			regexp.MustCompile(`sk-[a-zA-Z0-9]{40,}`),
			// Anthropic
			regexp.MustCompile(`sk-ant-[a-zA-Z0-9\-_]+`),
			// Telegram Bot Token
			regexp.MustCompile(`\b[0-9]{8,10}:[a-zA-Z0-9_\-]{35,}\b`),
			// Facebook Page/User Token
			regexp.MustCompile(`\bEAAK[a-zA-Z0-9]+\b`),
			// Private Keys (RSA, OPENSSH, DSA, EC, or generic)
			regexp.MustCompile(
				`(?s)-----BEGIN (?:RSA |OPENSSH |DSA |EC )?PRIVATE KEY-----.*?-----END (?:RSA |OPENSSH |DSA |EC )?PRIVATE KEY-----`,
			),
			// GitHub Tokens (ghp, gho, ghu, ghs, ghr)
			regexp.MustCompile(`(?:ghp|gho|ghu|ghs|ghr)_[a-zA-Z0-9]{36}`),
			// Slack Tokens (xoxb, xapp, xoxp)
			regexp.MustCompile(`(?:xoxb|xapp|xoxp)-[a-zA-Z0-9\-]+`),
			// Zhipu / Generic JWTs (header.payload.signature)
			// A basic JWT pattern for commonly used tokens (be careful of false positives, so we keep it somewhat strict)
			regexp.MustCompile(`\beyJ[a-zA-Z0-9_-]*\.[a-zA-Z0-9_-]*\.[a-zA-Z0-9_-]+\b`),
		},
	}
	return r
}

// AddSecret adds a specific string (like a loaded API key) to be redacted exactly.
// It will ignore empty strings or very short strings (less than 8 chars) to avoid false positives.
func (r *SecretRedactor) AddSecret(secret string) {
	secret = strings.TrimSpace(secret)
	if len(secret) < 8 {
		return // Ignore short strings to prevent accidental redaction of common words
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.exactSecrets[secret] = struct{}{}
}

// ConfigValueToRedact extracts exact string values from a raw struct (like config strings).
// Passing loaded API keys here ensures they are caught even if they don't match a regex.
func (r *SecretRedactor) AddSecrets(secrets ...string) {
	for _, s := range secrets {
		r.AddSecret(s)
	}
}

// Redact scans the input string and replaces matches with [REDACTED_SECRET].
func (r *SecretRedactor) Redact(input string) string {
	if input == "" {
		return input
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	res := input

	// 1. Replace exact secrets first
	for secret := range r.exactSecrets {
		res = strings.ReplaceAll(res, secret, redactedPlaceholder)
	}

	// 2. Replace regex patterns
	for _, pattern := range r.regexPatterns {
		res = pattern.ReplaceAllString(res, redactedPlaceholder)
	}

	return res
}
