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
	"strings"

	"github.com/comgunner/picoclaw/pkg/utils"
)

// ============== Image Approval Tool ==============

type ImageApprovalTool struct {
	tracker *utils.ImageGenTracker
}

func NewImageApprovalTool() *ImageApprovalTool {
	return &ImageApprovalTool{}
}

func NewImageApprovalToolWithTracker(tracker *utils.ImageGenTracker) *ImageApprovalTool {
	return &ImageApprovalTool{tracker: tracker}
}

func (t *ImageApprovalTool) Name() string {
	return "image_approval"
}

func (t *ImageApprovalTool) Description() string {
	return "Aprobar, rechazar o regenerar imagen generada. Usar después de image_gen_create para obtener aprobación humana antes de publicar."
}

func (t *ImageApprovalTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"image_path": map[string]any{
				"type":        "string",
				"description": "Path al archivo de imagen generada (ej: ./workspace/image_gen/20260302_.../...-imagen.jpg)",
			},
			"action": map[string]any{
				"type":        "string",
				"description": "Acción a realizar: 'approve', 'reject', 'regenerate_with_changes'",
				"enum":        []string{"approve", "reject", "regenerate_with_changes"},
			},
			"new_prompt": map[string]any{
				"type":        "string",
				"description": "Si action='regenerate_with_changes', nuevo prompt para regenerar",
			},
		},
		"required": []string{"image_path", "action"},
	}
}

func (t *ImageApprovalTool) Execute(ctx context.Context, args map[string]any) *ToolResult {
	imagePath, ok := args["image_path"].(string)
	if !ok || strings.TrimSpace(imagePath) == "" {
		return ErrorResult("image_path es requerido y debe ser un string válido")
	}

	action, ok := args["action"].(string)
	if !ok || strings.TrimSpace(action) == "" {
		return ErrorResult("action es requerido: 'approve', 'reject', o 'regenerate_with_changes'")
	}

	newPrompt, _ := args["new_prompt"].(string)

	// Extraer ID del path
	id := extractIDFromPath(imagePath)

	// Actualizar tracker si está configurado
	if t.tracker != nil {
		switch action {
		case "approve":
			t.tracker.UpdateMetadata(id, "status", "image_approved")
			t.tracker.UpdateMetadata(id, "user_approved_image", "true")
		case "reject":
			t.tracker.UpdateMetadata(id, "status", "image_rejected")
			t.tracker.UpdateMetadata(id, "user_approved_image", "false")
			count := t.incrementRegenerationCount(id, "image")
			t.tracker.UpdateMetadata(id, "regeneration_count_image", fmt.Sprintf("%d", count))
		case "regenerate_with_changes":
			t.tracker.UpdateMetadata(id, "status", "image_regenerating")
			t.tracker.UpdateMetadata(id, "user_approved_image", "false")
			if newPrompt != "" {
				t.tracker.UpdateMetadata(id, "prompt", newPrompt)
			}
			count := t.incrementRegenerationCount(id, "image")
			t.tracker.UpdateMetadata(id, "regeneration_count_image", fmt.Sprintf("%d", count))
		}
	}

	// Construir respuesta según la acción
	switch action {
	case "approve":
		response := fmt.Sprintf(`✅ **Imagen APROBADA**

📸 Imagen: %s
📊 Tracker actualizado: status="image_approved"

✅ Listo para publicar. ¿Dónde quieres publicar?
- Usa social_manager o la herramienta de publicación correspondiente`,
			imagePath)

		return UserResult(response)

	case "reject":
		response := fmt.Sprintf(`❌ **Imagen RECHAZADA**

📸 Imagen: %s
📊 Tracker actualizado: status="image_rejected"

🔄 ¿Quieres que regenerate la imagen?
- Si es así, indica qué cambios necesitas o usa "regenerar" para intentar con variación del mismo prompt`,
			imagePath)

		return UserResult(response)

	case "regenerate_with_changes":
		if newPrompt == "" {
			return ErrorResult("Si action='regenerate_with_changes', debes proporcionar new_prompt")
		}

		response := fmt.Sprintf(`🔄 **Regenerando Imagen**

📸 Imagen original: %s

🎨 Nuevo prompt:
%s

📊 Tracker actualizado: status="image_regenerating"

⏳ Generando nueva imagen...`,
			imagePath, newPrompt)

		return UserResult(response)

	default:
		return ErrorResult("action debe ser 'approve', 'reject', o 'regenerate_with_changes'")
	}
}

// incrementRegenerationCount incrementa el contador de regeneraciones
func (t *ImageApprovalTool) incrementRegenerationCount(id string, kind string) int {
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
