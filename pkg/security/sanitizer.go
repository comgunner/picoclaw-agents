// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)

package security

import (
	"fmt"
	"regexp"
)

// sanitizerPattern agrupa nombre y regexp para el sanitizador.
type sanitizerPattern struct {
	name    string
	pattern *regexp.Regexp
}

// defaultSanitizerPatterns son los patrons usados por Sanitize().
// Se inicializan una vez en init() para no recompilar en cada llamada.
var defaultSanitizerPatterns []sanitizerPattern

func init() {
	rawPatterns := []struct {
		name    string
		pattern string
	}{
		{"anthropic_key", `sk-ant-[a-zA-Z0-9\-]{32,}`},
		{"stripe_secret", `sk_live_[0-9a-zA-Z]{24}`},
		{"deepseek_key", `sk-[a-f0-9]{32,}`},
		{"openai_key", `sk-[a-zA-Z0-9\-]{40,}`},
		{"google_api_key", `AIza[0-9A-Za-z\-_]{35}`},
		{"github_token", `gh[pos]_[a-zA-Z0-9]{36}`},
		{"aws_access_key", `AKIA[0-9A-Z]{16}`},
		{"slack_token", `xox[baprs]-[0-9a-zA-Z\-]{10,}`},
		{"telegram_bot", `[0-9]{8,10}:[a-zA-Z0-9_\-]{35}`},
		{"jwt", `eyJ[a-zA-Z0-9_\-]+\.eyJ[a-zA-Z0-9_\-]+\.[a-zA-Z0-9_\-]+`},
	}

	defaultSanitizerPatterns = make([]sanitizerPattern, 0, len(rawPatterns))
	for _, rp := range rawPatterns {
		compiled, err := regexp.Compile(rp.pattern)
		if err != nil {
			continue
		}
		defaultSanitizerPatterns = append(defaultSanitizerPatterns, sanitizerPattern{
			name:    rp.name,
			pattern: compiled,
		})
	}
}

// Sanitize redacta secrets conocidos de un string arbitrario.
//
// A diferencia de redactor.go (que hookea el logger automáticamente),
// esta función se usa de forma EXPLÍCITA en puntos donde un string
// controlado externamente puede llegar al LLM:
//   - Resultados de tools (ToolResult.Content)
//   - Mensajes del usuario antes de enviarlos al provider
//   - Respuestas de APIs externas
//
// El reemplazo sigue el formato: [REDACTED_NOMBRE_PATRON]
func Sanitize(s string) string {
	for _, sp := range defaultSanitizerPatterns {
		s = sp.pattern.ReplaceAllStringFunc(s, func(match string) string {
			return fmt.Sprintf("[REDACTED_%s]", sp.name)
		})
	}
	return s
}

// SanitizeMap aplica Sanitize a todos los valores string de un map.
// Los valores no-string se copian sin modificación.
// Útil para sanitizar args de tools antes de loguearlos.
func SanitizeMap(m map[string]any) map[string]any {
	if m == nil {
		return nil
	}
	result := make(map[string]any, len(m))
	for k, v := range m {
		switch val := v.(type) {
		case string:
			result[k] = Sanitize(val)
		case map[string]any:
			result[k] = SanitizeMap(val) // recursivo
		default:
			result[k] = v
		}
	}
	return result
}
