// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package tools

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/comgunner/picoclaw/pkg/utils"
)

// ============== Text Approval Tool ==============

type TextApprovalTool struct {
	tracker *utils.ImageGenTracker
}

func NewTextApprovalTool() *TextApprovalTool {
	return &TextApprovalTool{}
}

func NewTextApprovalToolWithTracker(tracker *utils.ImageGenTracker) *TextApprovalTool {
	return &TextApprovalTool{tracker: tracker}
}

func (t *TextApprovalTool) Name() string {
	return "text_approval"
}

func (t *TextApprovalTool) Description() string {
	return "Aprobar, rechazar o editar texto borrador para redes sociales. Usar después de community_manager_create_draft para obtener aprobación humana antes de generar imagen o publicar."
}

func (t *TextApprovalTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"text_path": map[string]any{
				"type":        "string",
				"description": "Path al archivo de texto borrador (ej: ./workspace/text_scripts/20260302_.../...-post_facebook.txt)",
			},
			"action": map[string]any{
				"type":        "string",
				"description": "Acción a realizar: 'approve', 'reject', 'edit'",
				"enum":        []string{"approve", "reject", "edit"},
			},
			"platform": map[string]any{
				"type":        "string",
				"description": "Plataforma para la que se aprobó el texto: 'facebook', 'twitter', 'discord', etc.",
			},
			"edit_instructions": map[string]any{
				"type":        "string",
				"description": "Si action='edit', describir los cambios necesarios",
			},
		},
		"required": []string{"text_path", "action"},
	}
}

func (t *TextApprovalTool) Execute(ctx context.Context, args map[string]any) *ToolResult {
	textPath, ok := args["text_path"].(string)
	if !ok || strings.TrimSpace(textPath) == "" {
		return ErrorResult("text_path es requerido y debe ser un string válido")
	}

	action, ok := args["action"].(string)
	if !ok || strings.TrimSpace(action) == "" {
		return ErrorResult("action es requerido: 'approve', 'reject', o 'edit'")
	}

	platform, _ := args["platform"].(string)
	editInstructions, _ := args["edit_instructions"].(string)

	// Leer el texto para mostrarlo
	text, err := os.ReadFile(textPath)
	if err != nil {
		return ErrorResult(fmt.Sprintf("No se pudo leer el texto: %v", err))
	}

	// Extraer ID del path
	id := extractIDFromPath(textPath)

	// Actualizar tracker si está configurado
	if t.tracker != nil {
		switch action {
		case "approve":
			t.tracker.UpdateMetadata(id, "status", "text_approved")
			t.tracker.UpdateMetadata(id, "user_approved_text", "true")
			if platform != "" {
				t.tracker.UpdateMetadata(id, "platforms_target", platform)
			}
		case "reject":
			t.tracker.UpdateMetadata(id, "status", "text_rejected")
			t.tracker.UpdateMetadata(id, "user_approved_text", "false")
			count := t.incrementRegenerationCount(id, "text")
			t.tracker.UpdateMetadata(id, "regeneration_count_text", fmt.Sprintf("%d", count))
		case "edit":
			t.tracker.UpdateMetadata(id, "status", "text_editing")
			if editInstructions != "" {
				t.tracker.UpdateMetadata(id, "edit_instructions", editInstructions)
			}
		}
	}

	// Construir respuesta según la acción
	switch action {
	case "approve":
		response := fmt.Sprintf(`✅ **Texto APROBADO**

📄 Texto aprobado para %s:
%s

📁 Path: %s
📊 Tracker actualizado: status="text_approved"

✅ Listo para generar imagen o publicar.`,
			platform, string(text), textPath)

		return UserResult(response)

	case "reject":
		response := fmt.Sprintf(`❌ **Texto RECHAZADO**

📄 Texto original:
%s

📁 Path: %s
📊 Tracker actualizado: status="text_rejected"

🔄 ¿Quieres que regenerate el texto? Si es así, indica qué cambios necesitas.`,
			string(text), textPath)

		return UserResult(response)

	case "edit":
		if editInstructions == "" {
			return ErrorResult("Si action='edit', debes proporcionar edit_instructions describiendo los cambios")
		}

		response := fmt.Sprintf(`✏️ **Texto en EDICIÓN**

📄 Texto original:
%s

✍️ Cambios solicitados:
%s

📁 Path: %s
📊 Tracker actualizado: status="text_editing"

🔄 Regenerando texto con los cambios...`,
			string(text), editInstructions, textPath)

		return UserResult(response)

	default:
		return ErrorResult("action debe ser 'approve', 'reject', o 'edit'")
	}
}

// extractIDFromPath extrae el ID único de un path de archivo
func extractIDFromPath(path string) string {
	// Path típico: ./workspace/text_scripts/20260302_104101_j9q1yf/20260302_104101_j9q1yf.-post_facebook.txt
	// Queremos extraer: 20260302_104101_j9q1yf

	parts := strings.Split(path, "/")
	for i := len(parts) - 1; i >= 0; i-- {
		if strings.Contains(parts[i], ".-") {
			// Encontramos el archivo, extraer ID del nombre
			filename := parts[i]
			idx := strings.Index(filename, ".-")
			if idx > 0 {
				return filename[:idx]
			}
		}
	}

	// Fallback: intentar extraer del directorio
	for i := len(parts) - 2; i >= 0; i-- {
		if strings.Contains(parts[i], "_") && len(parts[i]) >= 15 {
			return parts[i]
		}
	}

	return ""
}

// incrementRegenerationCount incrementa el contador de regeneraciones
func (t *TextApprovalTool) incrementRegenerationCount(id string, kind string) int {
	if t.tracker == nil {
		return 0
	}

	record, ok := t.tracker.Get(id)
	if !ok {
		return 0
	}

	key := "regeneration_count"
	if kind == "text" {
		key = "regeneration_count_text"
	} else if kind == "image" {
		key = "regeneration_count_image"
	}

	count := 0
	if val, exists := record.Metadata[key]; exists {
		fmt.Sscanf(val, "%d", &count)
	}
	count++

	t.tracker.UpdateMetadata(id, key, fmt.Sprintf("%d", count))
	return count
}
