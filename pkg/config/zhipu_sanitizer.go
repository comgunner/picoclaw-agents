// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package config

import (
	"fmt"
	"strings"

	"github.com/comgunner/picoclaw/pkg/logger"
)

// ZhipuValidModels lista de modelos válidos de Zhipu AI (confirmados por API)
var ZhipuValidModels = map[string]bool{
	"glm-5":         true,
	"glm-5-turbo":   true,
	"glm-5.1":       true,
	"glm-4.7":       true,
	"glm-4.7-flash": true,
	"glm-4.6":       true,
	"glm-4.5":       true,
	"glm-4.5-flash": true,
	"glm-4.5-air":   true,
}

// ZhipuObsoleteModels mapea modelos obsoletos a sus reemplazos más cercanos
var ZhipuObsoleteModels = map[string]string{
	// glm-4-flash → glm-4.6 (similar velocidad)
	"glm-4-flash": "glm-4.6",

	// glm-4-air → glm-4.5-air (mismo nombre)
	"glm-4-air": "glm-4.5-air",

	// glm-4-airx → glm-5 (mejor calidad)
	"glm-4-airx": "glm-5",

	// glm-4-long → glm-5 (no hay long disponible, usar default)
	"glm-4-long": "glm-5",

	// glm-4v-flash → glm-5 (no hay vision disponible, usar default)
	"glm-4v-flash": "glm-5",
}

// SanitizeZhipuModels verifica y corrige modelos de Zhipu en la configuración
// Retorna true si se hicieron cambios, false si todo estaba correcto
func SanitizeZhipuModels(cfg *Config) bool {
	if cfg == nil {
		return false
	}

	changesMade := false

	// 1. Sanitizar modelo por defecto
	if cfg.Agents.Defaults.GetModelName() != "" {
		modelName := cfg.Agents.Defaults.GetModelName()
		if isZhipuObsoleteModel(modelName) {
			replacement := getZhipuReplacement(modelName)
			logger.WarnC("config", fmt.Sprintf(
				"Modelo de Zhipu obsoleto detectado: %s → actualizando a %s",
				modelName, replacement,
			))
			cfg.Agents.Defaults.ModelName = replacement
			cfg.Agents.Defaults.Model = replacement
			changesMade = true
		}
	}

	// 2. Sanitizar lista de modelos
	for i := range cfg.ModelList {
		model := &cfg.ModelList[i]
		if isZhipuModel(model.Model) {
			// Verificar si el modelo es obsoleto
			if isZhipuObsoleteModel(model.ModelName) {
				replacement := getZhipuReplacement(model.ModelName)
				logger.WarnC("config", fmt.Sprintf(
					"Modelo de Zhipu obsoleto en ModelList: %s → actualizando a %s",
					model.ModelName, replacement,
				))
				model.ModelName = replacement
				model.Model = strings.Replace(model.Model, model.ModelName, replacement, 1)
				changesMade = true
			}

			// Verificar si el modelo existe realmente
			if !isZhipuValidModel(model.ModelName) {
				// Intentar encontrar un reemplazo
				replacement := findBestZhipuReplacement(model.ModelName)
				if replacement != "" {
					logger.WarnC("config", fmt.Sprintf(
						"Modelo de Zhipu no existe: %s → actualizando a %s",
						model.ModelName, replacement,
					))
					model.ModelName = replacement
					model.Model = strings.Replace(model.Model, model.ModelName, replacement, 1)
					changesMade = true
				}
			}
		}
	}

	if changesMade {
		logger.InfoC("config", "Configuración de Zhipu sanitizada exitosamente")
	}

	return changesMade
}

// isZhipuModel verifica si un modelo pertenece a Zhipu
func isZhipuModel(model string) bool {
	lower := strings.ToLower(strings.TrimSpace(model))
	return strings.HasPrefix(lower, "glm-") ||
		strings.HasPrefix(lower, "zhipu/") ||
		lower == "zhipu" ||
		lower == "z.ai"
}

// isZhipuValidModel verifica si un modelo de Zhipu es válido
func isZhipuValidModel(model string) bool {
	modelName := strings.TrimSpace(model)
	// Remover prefijo si existe
	if strings.Contains(modelName, "/") {
		parts := strings.Split(modelName, "/")
		if len(parts) > 1 {
			modelName = parts[1]
		}
	}
	return ZhipuValidModels[modelName]
}

// isZhipuObsoleteModel verifica si un modelo de Zhipu es obsoleto
func isZhipuObsoleteModel(model string) bool {
	modelName := strings.TrimSpace(model)
	// Remover prefijo si existe
	if strings.Contains(modelName, "/") {
		parts := strings.Split(modelName, "/")
		if len(parts) > 1 {
			modelName = parts[1]
		}
	}
	_, exists := ZhipuObsoleteModels[modelName]
	return exists
}

// getZhipuReplacement obtiene el reemplazo para un modelo obsoleto
func getZhipuReplacement(model string) string {
	modelName := strings.TrimSpace(model)
	// Remover prefijo si existe
	if strings.Contains(modelName, "/") {
		parts := strings.Split(modelName, "/")
		if len(parts) > 1 {
			modelName = parts[1]
		}
	}
	if replacement, exists := ZhipuObsoleteModels[modelName]; exists {
		return replacement
	}
	return "glm-5" // Default fallback
}

// findBestZhipuReplacement intenta encontrar el mejor reemplazo para un modelo inválido
func findBestZhipuReplacement(model string) string {
	modelName := strings.TrimSpace(model)

	// Si es un modelo glm-4.x, sugerir glm-4.5 o glm-5
	if strings.HasPrefix(modelName, "glm-4") {
		return "glm-4.5"
	}

	// Si es un modelo glm-5.x, sugerir glm-5
	if strings.HasPrefix(modelName, "glm-5") {
		return "glm-5"
	}

	// Default a glm-5
	return "glm-5"
}

// GetZhipuModelsInfo retorna información sobre modelos Zhipu
func GetZhipuModelsInfo() map[string]any {
	return map[string]any{
		"valid_models":    getMapKeys(ZhipuValidModels),
		"obsolete_models": ZhipuObsoleteModels,
		"default_model":   "glm-5",
		"total_valid":     len(ZhipuValidModels),
		"total_obsolete":  len(ZhipuObsoleteModels),
	}
}

// getMapKeys helper para obtener keys de un mapa
func getMapKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
