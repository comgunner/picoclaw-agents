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
	"path/filepath"
	"strings"

	"github.com/comgunner/picoclaw/pkg/utils"
)

// ============== Image Generation Workflow Tools ==============

// ImageGenWorkflowTool proporciona el workflow interaction post-generación
type ImageGenWorkflowTool struct {
	tracker *utils.ImageGenTracker
}

func NewImageGenWorkflowTool() *ImageGenWorkflowTool {
	return &ImageGenWorkflowTool{}
}

func NewImageGenWorkflowToolWithTracker(tracker *utils.ImageGenTracker) *ImageGenWorkflowTool {
	return &ImageGenWorkflowTool{
		tracker: tracker,
	}
}

func (t *ImageGenWorkflowTool) Name() string {
	return "image_gen_workflow"
}

func (t *ImageGenWorkflowTool) Description() string {
	return "Gestionar el flujo post-generación de imágenes: aprobar, rechazar para regenerar, o preparar publicación."
}

func (t *ImageGenWorkflowTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"image_path": map[string]any{
				"type":        "string",
				"description": "Ruta de la imagen generada",
			},
			"action": map[string]any{
				"type":        "string",
				"description": "Acción a realizar: 'approve' (aprobar), 'reject' (rechazar/regenerar), 'generate_text' (crear copy), 'schedule' (programar)",
				"enum":        []string{"approve", "reject", "generate_text", "schedule"},
			},
			"prompt_used": map[string]any{
				"type":        "string",
				"description": "Prompt usado (opcional, para regeneración)",
			},
		},
		"required": []string{"image_path"},
	}
}

func (t *ImageGenWorkflowTool) Execute(ctx context.Context, args map[string]any) *ToolResult {
	imagePath, _ := args["image_path"].(string)
	action, _ := args["action"].(string)
	promptUsed, _ := args["prompt_used"].(string)

	if imagePath == "" {
		return ErrorResult("image_path es requerido")
	}

	// Extraer ID del path (formato: .../ID/ID.-imagen.jpg)
	id := t.extractIDFromPath(imagePath)

	// Si no hay acción, mostrar el menú (compatibilidad anterior)
	if action == "" {
		menu := t.buildWorkflowMenu(imagePath, promptUsed)
		return UserResult(menu)
	}

	// Procesar acción
	if t.tracker == nil {
		return ErrorResult("Tracker no inicializado en la tool")
	}

	switch action {
	case "approve":
		t.tracker.UpdateMetadata(id, "status", "approved")
		t.tracker.UpdateMetadata(id, "user_approved", "true")
		return UserResult(fmt.Sprintf("✅ **Imagen aprobada**. Se ha marcado como lista para publicar.\nID: `%s`", id))

	case "reject":
		t.tracker.UpdateMetadata(id, "status", "rejected")
		t.tracker.UpdateMetadata(id, "user_approved", "false")
		// Sugerir regeneración
		return UserResult(
			fmt.Sprintf("🔄 **Imagen rechazada**. ¿Quieres regenerarla con un prompt diferente?\nID: `%s`", id),
		)

	case "generate_text":
		t.tracker.UpdateMetadata(id, "status", "text_pending")
		return UserResult(
			"📝 **Solicitud de texto recibida**. Por favor, usa la tool `social_manager` o `community_from_image` para generar el copy.",
		)

	case "schedule":
		return UserResult("⏰ **Programación**: Indica la fecha y hora para programar la publicación.")
	}

	return UserResult(t.buildWorkflowMenu(imagePath, promptUsed))
}

func (t *ImageGenWorkflowTool) extractIDFromPath(path string) string {
	filename := filepath.Base(path)
	if idx := strings.Index(filename, ".-"); idx > 0 {
		return filename[:idx]
	}
	// Fallback: usar nombre del directorio padre
	return filepath.Base(filepath.Dir(path))
}

func (t *ImageGenWorkflowTool) buildWorkflowMenu(imagePath, promptUsed string) string {
	var sb strings.Builder

	sb.WriteString("🎨 **Opciones de Workflow de Imagen**\n\n")
	sb.WriteString(fmt.Sprintf("📁 Archivo: `%s`\n\n", imagePath))

	sb.WriteString("**Selecciona una acción:**\n")
	sb.WriteString("1. ✅ **Aprobar** (`action='approve'`): Marcar como válida.\n")
	sb.WriteString("2. 🔄 **Rechazar** (`action='reject'`): Marcar como no válida/regenerar.\n")
	sb.WriteString("3. 📝 **Generar Texto** (`action='generate_text'`): Crear copy para redes.\n")
	sb.WriteString("4. ⏰ **Programar** (`action='schedule'`): Definir fecha de publicación.\n\n")

	sb.WriteString("💡 *Puedes decir algo como: \"Me gusta, apruébala\" o \"No me gusta, rechaza y regenera\"*")

	return sb.String()
}

// ScriptToImageWorkflowTool proporciona el flujo Script → Imagen
type ScriptToImageWorkflowTool struct{}

func NewScriptToImageWorkflowTool() *ScriptToImageWorkflowTool {
	return &ScriptToImageWorkflowTool{}
}

func (t *ScriptToImageWorkflowTool) Name() string {
	return "script_to_image_workflow"
}

func (t *ScriptToImageWorkflowTool) Description() string {
	return "Flujo completo: crear guion de texto primero, luego generar imagen basada en el guion."
}

func (t *ScriptToImageWorkflowTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"topic": map[string]any{
				"type":        "string",
				"description": "Tema para el guion e imagen",
			},
			"category": map[string]any{
				"type":        "string",
				"description": "Categoría del guion: 'historia', 'noticia', 'tutorial'",
			},
			"create_script_first": map[string]any{
				"type":        "boolean",
				"description": "Crear guion de texto primero (default: true)",
			},
		},
		"required": []string{"topic"},
	}
}

func (t *ScriptToImageWorkflowTool) Execute(ctx context.Context, args map[string]any) *ToolResult {
	topic, _ := args["topic"].(string)
	category, _ := args["category"].(string)

	// Default explícito: create_script_first = true
	createScriptFirst := true
	if val, ok := args["create_script_first"].(bool); ok {
		createScriptFirst = val
	}

	if topic == "" {
		return ErrorResult("topic es requerido")
	}

	if category == "" {
		category = "historia"
	}

	// Workflow description - Instrucciones claras para que el LLM ejecute
	var sb strings.Builder

	sb.WriteString("📋 **Flujo Script → Imagen**\n\n")
	sb.WriteString(fmt.Sprintf("📌 Tema: %s\n", topic))
	sb.WriteString(fmt.Sprintf("📁 Categoría: %s\n\n", category))

	if createScriptFirst {
		sb.WriteString("## EJECUTA ESTOS PASOS EN ORDEN:\n\n")
		sb.WriteString("**Paso 1**: Generar guion de texto\n")
		sb.WriteString(fmt.Sprintf("  - EJECUTA: `Usa text_script_create topic='%s' category='%s'`\n", topic, category))
		sb.WriteString("  - El script se guarda automáticamente en <workspace>/image_gen/<id>/<id>.-script.txt\n\n")

		sb.WriteString("**Paso 2**: Generar imagen desde el tema\n")
		sb.WriteString(fmt.Sprintf("  - EJECUTA: `Usa image_gen_create prompt='%s'`\n", topic))
		sb.WriteString("  - La imagen se guarda en la misma carpeta: <workspace>/image_gen/<id>/<id>.-imagen.jpg\n\n")

		sb.WriteString("**Paso 3**: Mostrar opciones post-generación\n")
		sb.WriteString("  - EJECUTA: `Usa image_gen_workflow image_path=<ruta de la imagen>`\n\n")

		sb.WriteString(
			"💡 **IMPORTANTE**: Ejecuta las herramientas en secuencia. Comienza con `text_script_create` AHORA.\n\n",
		)

		sb.WriteString("🤖 **Opcional (Automatización Multi-Agente)**:\n")
		sb.WriteString("Si quieres que PicoClaw complete todo el flujo solo, EJECUTA:\n")
		sb.WriteString(
			fmt.Sprintf(
				"`Usa spawn task='Crea un guion sobre %s, genera su imagen y muestra el workflow final' label='flujo_imagen'`\n",
				topic,
			),
		)
	} else {
		sb.WriteString("**Modo directo**: Generar imagen sin guion\n\n")
		sb.WriteString(fmt.Sprintf("💡 EJECUTA AHORA: `Usa image_gen_create prompt='%s'`\n", topic))
	}

	return UserResult(sb.String())
}
