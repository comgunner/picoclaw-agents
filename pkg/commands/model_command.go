package commands

import (
	"fmt"
	"strings"

	"github.com/comgunner/picoclaw/pkg/config"
	"github.com/comgunner/picoclaw/pkg/logger"
)

// ModelCommandHandler maneja el commando /model
type ModelCommandHandler struct {
	configPath string
}

// NewModelCommandHandler crea un nuevo handler
func NewModelCommandHandler(configPath string, _ any) *ModelCommandHandler {
	return &ModelCommandHandler{
		configPath: configPath,
	}
}

// Handle procesa el commando /model
func (h *ModelCommandHandler) Handle(message string, channelID string) (string, error) {
	parts := strings.Fields(message)

	if len(parts) == 1 {
		// /model sin arguments → listar modelos
		return h.listModels(), nil
	}

	subcommand := parts[1]

	switch subcommand {
	case "info":
		if len(parts) < 3 {
			return "❌ Uso: /model info <nombre_modelo>", nil
		}
		return h.getModelInfo(parts[2]), nil

	case "provider":
		if len(parts) < 3 {
			return "❌ Uso: /model provider <nombre_proveedor>", nil
		}
		return h.filterModelsByProvider(parts[2]), nil

	default:
		// Asumir que es selección de modelo
		modelName := strings.Join(parts[1:], " ")
		return h.selectModel(modelName, channelID), nil
	}
}

// listModels lista todos los modelos configurados
func (h *ModelCommandHandler) listModels() string {
	cfg, err := config.LoadConfig(h.configPath)
	if err != nil {
		return fmt.Sprintf("❌ Error cargando config: %v", err)
	}

	if len(cfg.ModelList) == 0 {
		return "❌ No hay modelos configurados"
	}

	// Obtener modelo actual
	currentModel := cfg.Agents.Defaults.ModelName

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("📦 Modelos disponibles (%d configurados):\n\n", len(cfg.ModelList)))

	for i, model := range cfg.ModelList {
		sanitized := h.sanitizeModelConfig(model)
		marker := "  "
		if model.ModelName == currentModel {
			marker = "👉"
		}

		authMethod := h.getAuthMethodLabel(model.AuthMethod)
		sb.WriteString(fmt.Sprintf("%s %d. %s (%s)\n", marker, i+1, sanitized.Model, authMethod))
	}

	sb.WriteString("\n💡 Usa: /model <nombre> para cambiar\n")
	sb.WriteString("   Ej: /model openai/gpt-5.4\n")
	sb.WriteString("   Ej: /model provider openai\n")
	sb.WriteString("   Ej: /model info antigravity/gemini-3-flash\n")

	return sb.String()
}

// selectModel cambia el modelo actual
func (h *ModelCommandHandler) selectModel(modelName string, channelID string) string {
	cfg, err := config.LoadConfig(h.configPath)
	if err != nil {
		return fmt.Sprintf("❌ Error cargando config: %v", err)
	}

	// Verificar que el modelo existe
	found := false
	var selectedModel config.ModelConfig
	for _, model := range cfg.ModelList {
		if model.Model == modelName || model.ModelName == modelName {
			found = true
			selectedModel = model
			break
		}
	}

	if !found {
		return fmt.Sprintf("❌ Modelo '%s' no encontrado\n\n%s", modelName, h.listModels())
	}

	// Actualizar modelo por defecto
	oldModel := cfg.Agents.Defaults.ModelName
	cfg.Agents.Defaults.ModelName = selectedModel.ModelName

	// Limpiar overrides de agentes para forzar el nuevo modelo global
	for i := range cfg.Agents.List {
		cfg.Agents.List[i].Model = nil
	}

	// Guardar config
	if err := config.SaveConfig(h.configPath, cfg); err != nil {
		return fmt.Sprintf("❌ Error guardando config: %v", err)
	}

	// Log del cambio
	logger.InfoCF("model_command", "Model changed", map[string]any{
		"from":    oldModel,
		"to":      selectedModel.ModelName,
		"channel": channelID,
	})

	sanitized := h.sanitizeModelConfig(selectedModel)
	authMethod := h.getAuthMethodLabel(selectedModel.AuthMethod)

	return fmt.Sprintf("✅ Modelo cambiado a: %s\n"+
		"   Proveedor: %s (%s)\n"+
		"   Anterior: %s",
		sanitized.ModelName,
		h.getProviderName(selectedModel.Model),
		authMethod,
		oldModel)
}

// filterModelsByProvider filtra modelos por proveedor
func (h *ModelCommandHandler) filterModelsByProvider(provider string) string {
	cfg, err := config.LoadConfig(h.configPath)
	if err != nil {
		return fmt.Sprintf("❌ Error cargando config: %v", err)
	}

	var filtered []config.ModelConfig
	providerLower := strings.ToLower(provider)

	for _, model := range cfg.ModelList {
		if strings.Contains(strings.ToLower(model.Model), providerLower) ||
			strings.Contains(strings.ToLower(model.ModelName), providerLower) {
			filtered = append(filtered, model)
		}
	}

	if len(filtered) == 0 {
		return fmt.Sprintf("❌ No se encontraron modelos para '%s'", provider)
	}

	currentModel := cfg.Agents.Defaults.ModelName

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("🔍 Modelos %s (%d):\n\n", provider, len(filtered)))

	for i, model := range filtered {
		sanitized := h.sanitizeModelConfig(model)
		marker := "  "
		if model.ModelName == currentModel {
			marker = "👉"
		}

		authMethod := h.getAuthMethodLabel(model.AuthMethod)
		sb.WriteString(fmt.Sprintf("%s %d. %s (%s)\n", marker, i+1, sanitized.Model, authMethod))
	}

	sb.WriteString("\n💡 Usa: /model <nombre> para cambiar")

	return sb.String()
}

// getModelInfo muestra información detallada de un modelo
func (h *ModelCommandHandler) getModelInfo(modelName string) string {
	cfg, err := config.LoadConfig(h.configPath)
	if err != nil {
		return fmt.Sprintf("❌ Error cargando config: %v", err)
	}

	var found *config.ModelConfig
	for i, model := range cfg.ModelList {
		if model.Model == modelName || model.ModelName == modelName {
			found = &cfg.ModelList[i]
			break
		}
	}

	if found == nil {
		return fmt.Sprintf("❌ Modelo '%s' no encontrado", modelName)
	}

	sanitized := h.sanitizeModelConfig(*found)
	provider := h.getProviderName(found.Model)
	authMethod := h.getAuthMethodLabel(found.AuthMethod)

	var sb strings.Builder
	sb.WriteString("📊 Información del Modelo:\n\n")
	sb.WriteString(fmt.Sprintf("   Nombre: %s\n", sanitized.ModelName))
	sb.WriteString(fmt.Sprintf("   Modelo: %s\n", sanitized.Model))
	sb.WriteString(fmt.Sprintf("   Proveedor: %s\n", provider))
	sb.WriteString(fmt.Sprintf("   Auth: %s\n", authMethod))

	if found.AuthMethod == "oauth" {
		sb.WriteString("   Estado: ✅ OAuth activo\n")
	} else if found.APIKey != "" && found.APIKey != "[REDACTED]" {
		sb.WriteString("   Estado: ✅ API Key configurada\n")
	} else {
		sb.WriteString("   Estado: ⚠️ Sin autenticación\n")
	}

	return sb.String()
}

// sanitizeModelConfig oculta información sensible
func (h *ModelCommandHandler) sanitizeModelConfig(cfg config.ModelConfig) config.ModelConfig {
	if cfg.APIKey != "" {
		cfg.APIKey = "[REDACTED]"
	}
	if strings.Contains(cfg.APIBase, "key=") ||
		strings.Contains(cfg.APIBase, "token=") {
		cfg.APIBase = "[REDACTED]"
	}
	return cfg
}

// getAuthMethodLabel convierte auth_method a label legible
func (h *ModelCommandHandler) getAuthMethodLabel(authMethod string) string {
	switch authMethod {
	case "oauth":
		return "OAuth"
	case "api_key":
		return "API Key"
	case "":
		return "Local"
	default:
		return authMethod
	}
}

// getProviderName extrae nombre del proveedor del modelo
func (h *ModelCommandHandler) getProviderName(model string) string {
	parts := strings.Split(model, "/")
	if len(parts) > 0 {
		return parts[0]
	}
	return "Unknown"
}
