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
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// APIKeyFormat define el formato esperado para cada proveedor
var APIKeyFormat = map[string]*regexp.Regexp{
	// OpenAI: sk- followed by 20+ alphanumeric chars
	"openai": regexp.MustCompile(`^sk-[a-zA-Z0-9]{20,}$`),

	// Anthropic: sk-ant- followed by alphanumeric and dashes
	"anthropic": regexp.MustCompile(`^sk-ant-[a-zA-Z0-9-]{20,}$`),

	// DeepSeek: sk- followed by 20+ alphanumeric chars
	"deepseek": regexp.MustCompile(`^sk-[a-zA-Z0-9]{20,}$`),

	// Zhipu/GLM: alphanumeric with dots, 20+ chars
	"zhipu": regexp.MustCompile(`^[a-zA-Z0-9.]{20,}$`),

	// Gemini: AIza followed by 35+ alphanumeric/underscore/dash
	"gemini": regexp.MustCompile(`^AIza[a-zA-Z0-9_-]{35,}$`),

	// Groq: gsk_ followed by 50+ alphanumeric chars
	"groq": regexp.MustCompile(`^gsk_[a-zA-Z0-9]{50,}$`),

	// OpenRouter: sk-or- followed by 30+ alphanumeric/dash/underscore
	"openrouter": regexp.MustCompile(`^sk-or-[a-zA-Z0-9_-]{30,}$`),

	// Qwen (Alibaba): alphanumeric, 20+ chars
	"qwen": regexp.MustCompile(`^[a-zA-Z0-9]{20,}$`),

	// GitHub Copilot: ghu_, gho_, ghs_, ghr_ followed by 36 chars
	"github-copilot": regexp.MustCompile(`^(ghu|gho|ghs|ghr)_[a-zA-Z0-9]{36}$`),
}

// ProviderEndpoints define URLs mínimas para validación de API keys
var ProviderEndpoints = map[string]string{
	"openai":         "https://api.openai.com/v1/models",
	"anthropic":      "https://api.anthropic.com/v1/models",
	"deepseek":       "https://api.deepseek.com/models",
	"zhipu":          "https://open.bigmodel.cn/api/paas/v4/models",
	"gemini":         "https://generativelanguage.googleapis.com/v1beta/models",
	"groq":           "https://api.groq.com/openai/v1/models",
	"openrouter":     "https://openrouter.ai/api/v1/models",
	"qwen":           "https://dashscope.aliyuncs.com/api/v1/models",
	"github-copilot": "https://api.githubcopilot.com/models",
}

// ValidationResult contiene el resultado de la validación de una API key
type ValidationResult struct {
	Provider string
	Valid    bool
	Error    error
	Message  string
}

// ValidateAPIKeyFormat valida el formato de una API key sin hacer llamadas HTTP
func ValidateAPIKeyFormat(provider, key string) error {
	if key == "" {
		return fmt.Errorf("API key for %s is empty", provider)
	}

	// Trim whitespace
	key = strings.TrimSpace(key)

	// Verificar longitud mínima
	if len(key) < 10 {
		return fmt.Errorf("API key for %s too short (min 10 chars, got %d)", provider, len(key))
	}

	pattern, exists := APIKeyFormat[provider]
	if !exists {
		// Si no hay patrón conocido, al menos verificar longitud
		return nil
	}

	if !pattern.MatchString(key) {
		return fmt.Errorf("API key for %s has invalid format (expected pattern: %s)", provider, pattern.String())
	}

	return nil
}

// ValidateAPIKeyOnline valida una API key haciendo una llamada HTTP mínima al proveedor
func ValidateAPIKeyOnline(ctx context.Context, provider, key string) *ValidationResult {
	result := &ValidationResult{
		Provider: provider,
		Valid:    false,
	}

	// Validar formato primero
	if err := ValidateAPIKeyFormat(provider, key); err != nil {
		result.Error = err
		return result
	}

	endpoint, exists := ProviderEndpoints[provider]
	if !exists {
		result.Error = fmt.Errorf("no validation endpoint configured for provider %s", provider)
		return result
	}

	// Crear cliente HTTP con timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Construir request
	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		result.Error = fmt.Errorf("failed to create request: %w", err)
		return result
	}

	// Setear headers según proveedor
	switch provider {
	case "openai", "deepseek", "groq", "openrouter":
		req.Header.Set("Authorization", "Bearer "+key)
	case "anthropic":
		req.Header.Set("X-Api-Key", key)
		req.Header.Set("Anthropic-Version", "2023-06-01")
	case "gemini", "qwen":
		// Estas APIs usan query params o headers especiales
		req.Header.Set("Authorization", "Bearer "+key)
	case "github-copilot":
		req.Header.Set("Authorization", "Bearer "+key)
	default:
		req.Header.Set("Authorization", "Bearer "+key)
	}

	// Ejecutar request
	resp, err := client.Do(req)
	if err != nil {
		result.Error = fmt.Errorf("connection error: %w", err)
		return result
	}
	defer resp.Body.Close()

	// Verificar status code
	if resp.StatusCode == 200 || resp.StatusCode == 201 {
		result.Valid = true
		result.Message = "API key validated successfully"
		return result
	}

	// Interpretar errores comunes
	switch resp.StatusCode {
	case 401, 403:
		result.Error = fmt.Errorf("API key rejected (HTTP %d) - invalid or expired", resp.StatusCode)
	case 429:
		result.Error = fmt.Errorf("rate limited by provider (HTTP %d)", resp.StatusCode)
	default:
		result.Error = fmt.Errorf("API key validation failed (HTTP %d)", resp.StatusCode)
	}

	return result
}

// ValidateAllKeys valida todas las keys de proveedores en una configuración
// Devuelve una lista de errores encontrados
func ValidateAllKeys(providers map[string]string) []ValidationResult {
	var results []ValidationResult

	ctx := context.Background()

	for provider, key := range providers {
		if key == "" {
			continue // Skip empty keys
		}

		result := ValidateAPIKeyOnline(ctx, provider, key)
		results = append(results, *result)
	}

	return results
}

// QuickValidate valida rápidamente una API key (solo formato, sin HTTP)
// Útil para validación inicial antes de hacer llamadas HTTP
func QuickValidate(provider, key string) *ValidationResult {
	result := &ValidationResult{
		Provider: provider,
		Valid:    false,
	}

	if err := ValidateAPIKeyFormat(provider, key); err != nil {
		result.Error = err
		return result
	}

	result.Valid = true
	result.Message = "Format validation passed (online validation not performed)"
	return result
}

// GetProviderFromModelName extrae el nombre del proveedor desde un model_name
// Ej: "openai/gpt-4" → "openai", "anthropic/claude-3" → "anthropic"
func GetProviderFromModelName(modelName string) string {
	if !strings.Contains(modelName, "/") {
		return ""
	}

	parts := strings.SplitN(modelName, "/", 2)
	if len(parts) < 2 {
		return ""
	}

	provider := strings.ToLower(parts[0])

	// Mapear aliases comunes
	providerAliases := map[string]string{
		"claude":         "anthropic",
		"gpt":            "openai",
		"gemini":         "gemini",
		"deepseek":       "deepseek",
		"glm":            "zhipu",
		"qwen":           "qwen",
		"groq":           "groq",
		"openrouter":     "openrouter",
		"github-copilot": "github-copilot",
	}

	if mapped, exists := providerAliases[provider]; exists {
		return mapped
	}

	return provider
}
