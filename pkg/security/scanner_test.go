// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)

package security_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/comgunner/picoclaw/pkg/security"
)

func TestScanFileDetectsOpenAIKey(t *testing.T) {
	f, err := os.CreateTemp("", "*.env")
	require.NoError(t, err)
	defer os.Remove(f.Name())

	// OpenAI keys reales tienen formato sk-... con 40+ caracteres
	// Usar clave que no contenga placeholders como "1234567890", "abcdefgh", "example", etc.
	// Clave válida: 48 caracteres después de sk-, sin secuencias obvias
	_, err = f.WriteString("OPENAI_API_KEY=sk-X9kL2mNpQ4rS7tUv0WxY3zA6bC8dE1fG5hI9jK2lM4nO7pR\n")
	require.NoError(t, err)
	f.Close()

	s := security.NewScanner()
	results, err := s.ScanFile(f.Name())
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(results), 1, "Should detect at least 1 secret")

	// Verify OpenAI pattern was detected
	foundOpenAI := false
	for _, r := range results {
		if r.Pattern == "OpenAI" {
			foundOpenAI = true
			assert.Equal(t, 1, r.Line)
			assert.Contains(t, r.Match, "****")
			assert.NotContains(t, r.Match, "X9kL2mNpQ4rS7tUv0WxY3zA6bC8dE1fG5hI9jK2lM4nO7pR")
		}
	}
	assert.True(t, foundOpenAI, "Should detect OpenAI pattern")
}

func TestScanFileDetectsAnthropicKey(t *testing.T) {
	f, err := os.CreateTemp("", "*.go")
	require.NoError(t, err)
	defer os.Remove(f.Name())

	// Clave Anthropic válida sin placeholders (sin "abcdefgh" ni "1234567890")
	_, err = f.WriteString(`apiKey := "sk-ant-api03-X9kL2mNpQ4rS7tUv0WxY3zA6bC8dE1fG5hI9jK2l"` + "\n")
	require.NoError(t, err)
	f.Close()

	s := security.NewScanner()
	results, err := s.ScanFile(f.Name())
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(results), 1, "Should detect at least 1 secret")

	foundAnthropic := false
	for _, r := range results {
		if r.Pattern == "Anthropic" {
			foundAnthropic = true
			break
		}
	}
	assert.True(t, foundAnthropic, "Should detect Anthropic pattern")
}

func TestScanFileDetectsGitHubToken(t *testing.T) {
	f, err := os.CreateTemp("", "*.env")
	require.NoError(t, err)
	defer os.Remove(f.Name())

	// GitHub token válido sin placeholders
	_, err = f.WriteString("GITHUB_TOKEN=ghp_X9kL2mNpQ4rS7tUv0WxY3zA6bC8dE1fG5hI9\n")
	require.NoError(t, err)
	f.Close()

	s := security.NewScanner()
	results, err := s.ScanFile(f.Name())
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(results), 1, "Should detect at least 1 secret")

	foundGitHub := false
	for _, r := range results {
		if r.Pattern == "GitHub Token" {
			foundGitHub = true
			break
		}
	}
	assert.True(t, foundGitHub, "Should detect GitHub Token pattern")
}

func TestScanFileDetectsJWT(t *testing.T) {
	f, err := os.CreateTemp("", "*.txt")
	require.NoError(t, err)
	defer os.Remove(f.Name())

	// JWT de ejemplo (header.payload.signature)
	jwt := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U"
	_, err = f.WriteString("token: " + jwt + "\n")
	require.NoError(t, err)
	f.Close()

	s := security.NewScanner()
	results, err := s.ScanFile(f.Name())
	require.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "JWT", results[0].Pattern)
}

func TestScanFileNoFalsePositives(t *testing.T) {
	f, err := os.CreateTemp("", "*.md")
	require.NoError(t, err)
	defer os.Remove(f.Name())

	// Contenido sin secrets reales
	content := `# Documentation

This is a sample configuration file.
No secrets here, just plain text.
`
	_, err = f.WriteString(content)
	require.NoError(t, err)
	f.Close()

	s := security.NewScanner()
	results, err := s.ScanFile(f.Name())
	require.NoError(t, err)
	assert.Empty(t, results, "No debería detectar falsos positivos")
}

func TestScanFilePlaceholderIgnored(t *testing.T) {
	f, err := os.CreateTemp("", "*.example")
	require.NoError(t, err)
	defer os.Remove(f.Name())

	// Placeholders comunes que deberían ignorarse
	content := `# Example config
OPENAI_API_KEY=your_key_here
ANTHROPIC_KEY=sk-test-placeholder
GITHUB_TOKEN=example_token
`
	_, err = f.WriteString(content)
	require.NoError(t, err)
	f.Close()

	s := security.NewScanner()
	results, err := s.ScanFile(f.Name())
	require.NoError(t, err)
	assert.Empty(t, results, "Los placeholders deberían ignorarse")
}

func TestScanDirRecursive(t *testing.T) {
	dir := t.TempDir()
	subdir := filepath.Join(dir, "subdir")
	require.NoError(t, os.MkdirAll(subdir, 0o700))

	// Secret en subdirectorio
	f := filepath.Join(subdir, "secrets.env")
	err := os.WriteFile(f, []byte(`API_KEY=sk-X9kL2mNpQ4rS7tUv0WxY3zA6bC8dE1fG5hI9jK2lM4nO7pR`), 0o600)
	require.NoError(t, err)

	s := security.NewScanner()
	results, err := s.ScanDir(dir)
	require.NoError(t, err)
	// The value may match multiple patterns (e.g. OpenAI + AWS Secret overlap)
	assert.GreaterOrEqual(t, len(results), 1)
	assert.Equal(t, f, results[0].File)
}

func TestScanDirSkipsHiddenDirectories(t *testing.T) {
	dir := t.TempDir()
	hiddenDir := filepath.Join(dir, ".git")
	require.NoError(t, os.MkdirAll(hiddenDir, 0o700))

	// Secret en directorio oculto - debería omitirse
	f := filepath.Join(hiddenDir, "config")
	err := os.WriteFile(f, []byte(`API_KEY=sk-X9kL2mNpQ4rS7tUv0WxY3zA6bC8dE1fG5hI9jK2lM4nO7pR`), 0o600)
	require.NoError(t, err)

	s := security.NewScanner()
	results, err := s.ScanDir(dir)
	require.NoError(t, err)
	assert.Empty(t, results, "Los directors ocultos deberían omitirse")
}

func TestScanDirSkipsVendorAndNodeModules(t *testing.T) {
	dir := t.TempDir()
	vendorDir := filepath.Join(dir, "vendor")
	nodeModulesDir := filepath.Join(dir, "node_modules")
	require.NoError(t, os.MkdirAll(vendorDir, 0o700))
	require.NoError(t, os.MkdirAll(nodeModulesDir, 0o700))

	// Secrets en directors que deberían omitirse
	err := os.WriteFile(
		filepath.Join(vendorDir, "config.go"),
		[]byte(`const key = "sk-X9kL2mNpQ4rS7tUv0WxY3zA6bC8dE1fG5hI9jK2lM4nO7pR"`),
		0o600,
	)
	require.NoError(t, err)
	err = os.WriteFile(
		filepath.Join(nodeModulesDir, "config.js"),
		[]byte(`const key = "sk-X9kL2mNpQ4rS7tUv0WxY3zA6bC8dE1fG5hI9jK2lM4nO7pR"`),
		0o600,
	)
	require.NoError(t, err)

	s := security.NewScanner()
	results, err := s.ScanDir(dir)
	require.NoError(t, err)
	assert.Empty(t, results, "vendor y node_modules deberían omitirse")
}

func TestScanFileNonExistent(t *testing.T) {
	s := security.NewScanner()
	_, err := s.ScanFile("/nonexistent/path/file.go")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "scanner: open")
}

func TestScanResultString(t *testing.T) {
	result := security.ScanResult{
		File:    "/path/to/file.go",
		Line:    42,
		Pattern: "OpenAI",
		Match:   "sk-ab****",
	}

	str := result.String()
	assert.Contains(t, str, "/path/to/file.go")
	assert.Contains(t, str, "42")
	assert.Contains(t, str, "OpenAI")
	assert.Contains(t, str, "sk-ab****")
}

func TestScanDirMultipleFiles(t *testing.T) {
	dir := t.TempDir()

	// Múltiples archivos con secrets
	// Note: values must not trigger isPlaceholder() — avoid "abcdefgh", "example", etc.
	// Note: broad AWS Secret pattern ([0-9a-zA-Z/+]{40}) may also match OpenAI keys.
	files := map[string]string{
		"config.env":    "OPENAI_KEY=sk-X9kL2mNpQ4rS7tUv0WxY3zA6bC8dE1fG5hI9jK2lM4nO7pR",
		"settings.json": `{"anthropic_key": "sk-ant-api03-Xp7mNt4rS9vW2yZ6bC1dF8gH3jK5lM0nO7pQ"}`,
		"code.go":       `const token = "ghp_Xp7mNt4rS9vW2yZ6bC1dF8gH3jK5lM0nO7pQ"`,
	}

	for filename, content := range files {
		err := os.WriteFile(filepath.Join(dir, filename), []byte(content), 0o600)
		require.NoError(t, err)
	}

	s := security.NewScanner()
	results, err := s.ScanDir(dir)
	require.NoError(t, err)
	// Broad patterns (e.g. AWS Secret) may produce extra matches — assert at least 3
	assert.GreaterOrEqual(t, len(results), 3)

	// Verificar que todos los patrons fueron detectados
	patterns := make(map[string]bool)
	for _, r := range results {
		patterns[r.Pattern] = true
	}
	assert.True(t, patterns["OpenAI"], "Debería detectar OpenAI key")
	assert.True(t, patterns["Anthropic"], "Debería detectar Anthropic key")
	assert.True(t, patterns["GitHub Token"], "Debería detectar GitHub token")
}
