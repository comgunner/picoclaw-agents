// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
// Copyright (c) 2026 PicoClaw contributors

package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/comgunner/picoclaw/pkg/logger"
	"github.com/comgunner/picoclaw/pkg/utils"
)

// SocialManagerTool orquesta el flujo de borradores y publicación en redes sociales
type SocialManagerTool struct {
	tracker *utils.ImageGenTracker
	// ToolRegistry para acceder a las herramientas de publicación registradas
	registry *ToolRegistry
}

func NewSocialManagerTool(tracker *utils.ImageGenTracker) *SocialManagerTool {
	return &SocialManagerTool{
		tracker: tracker,
		// registry se setea después con SetRegistry
	}
}

// SetRegistry establece el registry de herramientas para acceder a facebook_post
func (t *SocialManagerTool) SetRegistry(registry *ToolRegistry) {
	t.registry = registry
}

func (t *SocialManagerTool) Name() string {
	return "social_manager"
}

func (t *SocialManagerTool) Description() string {
	return "Gestionar borradores y publicación en redes sociales (Facebook, X/Twitter, Discord) basado en imágenes del tracker."
}

func (t *SocialManagerTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"image_id": map[string]any{
				"type":        "string",
				"description": "ID de la imagen en el tracker",
			},
			"action": map[string]any{
				"type":        "string",
				"description": "Acción: 'draft' (crear borrador), 'publish' (publicar aprobado), 'list_pending' (ver pendientes), 'approve' (aprobar imagen/texto), 'regenerate' (volver a crear)",
				"enum":        []string{"draft", "publish", "list_pending", "approve", "regenerate", "edit", "cancel"},
			},
			"platforms": map[string]any{
				"type":        "string",
				"description": "Plataformas (separadas por coma): 'facebook,twitter,discord'",
			},
			"custom_text": map[string]any{
				"type":        "string",
				"description": "Texto personalizado para el post (opcional)",
			},
		},
		"required": []string{"action"},
	}
}

func (t *SocialManagerTool) Execute(ctx context.Context, args map[string]any) *ToolResult {
	action, _ := args["action"].(string)
	imageID, _ := args["image_id"].(string)
	if imageID == "" {
		imageID, _ = args["id"].(string) // Alias
	}
	platformsArg, _ := args["platforms"].(string)
	customText, _ := args["custom_text"].(string)

	if action == "list_pending" {
		return t.listPending()
	}

	if imageID == "" {
		return ErrorResult("image_id es requerido para esta acción")
	}

	record, ok := t.tracker.Get(imageID)
	if !ok {
		// ID no encontrado - sugerir IDs similares
		similarIDs := t.findSimilarIDs(imageID)

		if len(similarIDs) > 0 {
			// Construir mensaje con sugerencias
			suggestions := "IDs similares encontrados:\n"
			for _, id := range similarIDs {
				suggestions += fmt.Sprintf("  • `%s`\n", id)
			}
			suggestions += "\n💡 Usa `/list pending` para ver todos los IDs válidos."

			return ErrorResult(fmt.Sprintf(
				"❌ **ID no encontrado**: `%s`\n\n%s",
				imageID, suggestions))
		}

		// No hay IDs similares - listar todos los disponibles
		allRecords := t.tracker.List()
		if len(allRecords) > 0 {
			availableIDs := "IDs disponibles en el tracker:\n"
			for _, r := range allRecords {
				status := r.Metadata["status"]
				if status == "" {
					status = "unknown"
				}
				availableIDs += fmt.Sprintf("  • `%s` (%s)\n", r.ID, status)
			}
			availableIDs += "\n💡 Usa `/list pending` para ver solo los pendientes."

			return ErrorResult(fmt.Sprintf(
				"❌ **ID no encontrado**: `%s`\n\n%s",
				imageID, availableIDs))
		}

		return ErrorResult(fmt.Sprintf(
			"❌ **ID no encontrado**: `%s`\n\nNo hay imágenes en el tracker. Genera una imagen primero.",
			imageID))
	}

	platforms := parsePlatforms(platformsArg)
	if len(platforms) == 0 && action == "draft" {
		platforms = []string{"discord"} // Default
	}

	switch action {
	case "draft":
		return t.handleDraft(ctx, record, platforms, customText)
	case "publish":
		return t.handlePublish(ctx, record, platforms)
	case "approve":
		t.tracker.UpdateMetadata(record.ID, "user_approved_text", "true")
		t.tracker.UpdateMetadata(record.ID, "user_approved_image", "true")
		t.tracker.UpdateMetadata(record.ID, "status", "approved")
		return UserResult(
			fmt.Sprintf(
				"✅ **ID %s Aprobado.**\n\n🚀 Para publicar ahora mismo, usa:\n`/bundle_publish id=%s platforms=facebook`\n\n❌ O para cancelar:\n`/bundle_cancel id=%s`",
				record.ID,
				record.ID,
				record.ID,
			),
		)
	case "regenerate":
		return UserResult(
			fmt.Sprintf(
				"🔄 Solicitada regeneración para %s.\n\n💡 Usa `social_post_bundle` de nuevo con el mismo tópico para generar una nueva versión.",
				record.ID,
			),
		)
	case "edit":
		return UserResult(
			fmt.Sprintf(
				"✏️ Modo edición para %s.\n\n💡 Dime qué cambios quieres en el texto o usa:\n`Usa social_manager action='draft' image_id='%s' custom_text='Tu nuevo texto aquí'`",
				record.ID,
				record.ID,
			),
		)
	case "cancel":
		t.tracker.UpdateMetadata(record.ID, "status", "canceled")
		return UserResult(
			fmt.Sprintf(
				"❌ **ID %s Cancelado.**\n\nEl bundle ha sido marcado como cancelado y no se publicará.",
				record.ID,
			),
		)
	}

	return ErrorResult("Acción no válida")
}

func (t *SocialManagerTool) listPending() *ToolResult {
	records := t.tracker.List()
	var pending []string
	for _, r := range records {
		status := r.Metadata["status"]
		if status == "generated" || status == "approved" || status == "text_generated" {
			pending = append(pending, fmt.Sprintf("- `%s`: %s (%s)", r.ID, utils.Truncate(r.Prompt, 50), status))
		}
	}

	if len(pending) == 0 {
		return UserResult("No hay imágenes pendientes de publicación.")
	}

	return UserResult(
		"📋 **Imágenes pendientes de publicación**:\n\n" + strings.Join(
			pending,
			"\n",
		) + "\n\n💡 **Siguiente paso:**\n`Usa social_manager action='draft' image_id='ID'`",
	)
}

func (t *SocialManagerTool) handleDraft(
	ctx context.Context,
	record utils.ImageGenRecord,
	platforms []string,
	customText string,
) *ToolResult {
	var results []string
	var allDrafts []string

	// Usar el prompt visual o el original para el contexto del texto
	contextPrompt := record.Prompt
	if visual, ok := record.Metadata["visual_prompt"]; ok {
		contextPrompt = visual
	}

	// Generar/guardar texto para cada plataforma
	for _, p := range platforms {
		var draft string
		if customText != "" {
			draft = customText
		} else {
			// Simular generación de texto proporcional a la plataforma
			draft = t.generateDraftLocal(contextPrompt, p)
		}

		allDrafts = append(allDrafts, draft)
		results = append(results, fmt.Sprintf("✅ **Borrador para %s**:\n%s\n", strings.Title(p), draft))
	}

	// Guardar el texto del borrador en metadata (usar el primero como default)
	draftText := ""
	if len(allDrafts) > 0 {
		draftText = allDrafts[0]
	}
	t.tracker.UpdateMetadata(record.ID, "draft_text", draftText)
	if customText != "" {
		t.tracker.UpdateMetadata(record.ID, "custom_text", customText)
	}
	t.tracker.UpdateMetadata(record.ID, "status", "text_generated")
	t.tracker.UpdateMetadata(record.ID, "platforms_target", strings.Join(platforms, ","))

	response := fmt.Sprintf(
		"📝 **Borradores generados para la imagen `%s`**\n\n%s\n\n🚀 **¿Publicar ahora?** Usa:\n`Usa social_manager action='publish' image_id='%s'`",
		record.ID,
		strings.Join(results, "\n---\n"),
		record.ID,
	)

	return UserResult(response)
}

func (t *SocialManagerTool) handlePublish(
	ctx context.Context,
	record utils.ImageGenRecord,
	platforms []string,
) *ToolResult {
	// Telemetría inicial
	logger.InfoF("handlePublish iniciado", map[string]any{
		"id":         record.ID,
		"platforms":  platforms,
		"has_image":  record.ImagePath != "",
		"has_script": record.ScriptPath != "",
	})

	// Verificar aprobación de TEXTO
	if record.Metadata["user_approved_text"] != "true" {
		return ErrorResult("🛑 El TEXTO no ha sido aprobado todavía. Usa `text_approval action='approve'` primero.")
	}

	// Verificar aprobación de IMAGEN
	if record.Metadata["user_approved_image"] != "true" {
		return ErrorResult("🛑 La IMAGEN no ha sido aprobada todavía. Usa `image_approval action='approve'` primero.")
	}

	// Si no se especificaron plataformas, usar las del tracker
	if len(platforms) == 0 {
		if target, ok := record.Metadata["platforms_target"]; ok {
			platforms = parsePlatforms(target)
		}
	}

	if len(platforms) == 0 {
		return ErrorResult("No se especificaron plataformas para publicar.")
	}

	var success []string
	var errors []string

	// ==========================================
	// FASE 2: Lectura robusta de TEXTO (3 niveles de prioridad)
	// ==========================================
	postText := ""

	// Prioridad 1: Leer desde archivo de script
	if record.ScriptPath != "" {
		if data, err := os.ReadFile(record.ScriptPath); err == nil {
			postText = string(data)
			logger.InfoF("Leyendo texto desde archivo de script", map[string]any{
				"id":          record.ID,
				"script_path": record.ScriptPath,
				"text_length": len(postText),
			})
		}
	}

	// Prioridad 2: Fallback a metadata draft_text
	if postText == "" {
		if text, ok := record.Metadata["draft_text"]; ok {
			postText = text
			logger.InfoF("Leyendo texto desde metadata draft_text", map[string]any{
				"id":          record.ID,
				"text_length": len(postText),
			})
		}
	}

	// Prioridad 3: Fallback a custom_text
	if postText == "" {
		if text, ok := record.Metadata["custom_text"]; ok {
			postText = text
			logger.InfoF("Leyendo texto desde metadata custom_text", map[string]any{
				"id":          record.ID,
				"text_length": len(postText),
			})
		}
	}

	// Prioridad 4: Fallback final al prompt original
	if postText == "" {
		postText = record.Prompt
		logger.WarnF("Usando prompt como fallback para texto", map[string]any{
			"id":            record.ID,
			"prompt_length": len(postText),
		})
	}

	// ==========================================
	// FASE 3: Lectura robusta de IMAGEN (con filesystem fallback)
	// ==========================================
	imagePath := ""

	// Intento 1: ImagePath del registro
	if record.ImagePath != "" {
		imagePath = record.ImagePath
		logger.InfoF("Usando ImagePath del registro", map[string]any{
			"id":         record.ID,
			"image_path": imagePath,
		})
	}

	// Intento 2: Metadata image_path
	if imagePath == "" {
		if imageFile, ok := record.Metadata["image_path"]; ok {
			imagePath = imageFile
			logger.InfoF("Usando imagen desde metadata", map[string]any{
				"id":         record.ID,
				"image_path": imagePath,
			})
		}
	}

	// Intento 3: Fallback - buscar en filesystem
	if imagePath == "" && record.ID != "" {
		// Normalizar ruta base eliminando sufijos de workers
		baseDir := t.normalizeTrackerPath(t.tracker.TrackerPath)
		if baseDir == "." || baseDir == "" {
			baseDir = "./workspace/image_gen"
		}
		searchPath := filepath.Join(baseDir, record.ID, "*.jpg")

		if matches, err := filepath.Glob(searchPath); err == nil && len(matches) > 0 {
			imagePath = matches[0]
			logger.InfoF("Imagen encontrada via filesystem glob", map[string]any{
				"id":          record.ID,
				"image_path":  imagePath,
				"search_path": searchPath,
			})
		} else {
			// Si no encuentra en la ruta normalizada, intentar también en la ruta original como fallback
			originalBaseDir := filepath.Dir(t.tracker.TrackerPath)
			if originalBaseDir != baseDir && originalBaseDir != "." && originalBaseDir != "" {
				originalSearchPath := filepath.Join(originalBaseDir, record.ID, "*.jpg")
				if matches, err := filepath.Glob(originalSearchPath); err == nil && len(matches) > 0 {
					imagePath = matches[0]
					logger.InfoF("Imagen encontrada en ruta original (fallback)", map[string]any{
						"id":          record.ID,
						"image_path":  imagePath,
						"search_path": originalSearchPath,
					})
				}
			}

			if imagePath == "" {
				logger.WarnF("No se encontró imagen en filesystem", map[string]any{
					"id":          record.ID,
					"search_path": searchPath,
					"error":       err,
				})
			}
		}
	}

	if imagePath == "" {
		logger.WarnF("No se encontró ruta de imagen - publicación solo texto", map[string]any{
			"id": record.ID,
		})
	}

	// Telemetría pre-publicación
	logger.InfoF("Preparando publicación", map[string]any{
		"id":               record.ID,
		"platforms_count":  len(platforms),
		"post_text_length": len(postText),
		"has_image":        imagePath != "",
	})

	// ==========================================
	// Publicar en cada plataforma
	// ==========================================
	for _, p := range platforms {
		var publishErr error

		switch strings.ToLower(p) {
		case "facebook", "fb":
			// Usar la herramienta facebook_post del registry (con credenciales configuradas)
			if t.registry != nil {
				fbTool, ok := t.registry.Get("facebook_post")
				if ok && fbTool != nil {
					// Ejecutar publicación real con facebook_post tool
					result := fbTool.Execute(ctx, map[string]any{
						"message":           postText,
						"image_path":        imagePath,
						"confirmed_by_user": true, // Ya fue aprobado por el usuario
					})

					// Telemetría post-publicación Facebook
					if result.IsError {
						logger.ErrorF("Facebook publicación falló", map[string]any{
							"id":               record.ID,
							"error":            result.ForLLM,
							"post_text_length": len(postText),
							"image_path":       imagePath,
						})
						publishErr = fmt.Errorf("%s", result.ForLLM)
					} else {
						logger.InfoF("Facebook publicación exitosa", map[string]any{
							"id":               record.ID,
							"post_text_length": len(postText),
							"image_path":       imagePath,
						})
						success = append(success, "facebook")
					}
				} else {
					publishErr = fmt.Errorf("herramienta facebook_post no disponible en registry")
				}
			} else {
				publishErr = fmt.Errorf("registry no configurado - no se puede acceder a facebook_post")
			}

		case "twitter", "x":
			// TODO: Implementar cuando tengamos X tool
			publishErr = fmt.Errorf("publicación en X/Twitter aún no implementada")

		case "discord":
			// TODO: Implementar cuando tengamos Discord webhook tool
			publishErr = fmt.Errorf("publicación en Discord aún no implementada")

		default:
			publishErr = fmt.Errorf("plataforma no soportada: %s", p)
		}

		if publishErr != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", p, publishErr))
		}
	}

	// Actualizar estado solo si hubo éxitos
	if len(success) > 0 {
		t.tracker.UpdateMetadata(record.ID, "status", "posted")
		t.tracker.UpdateMetadata(record.ID, "posted_at", time.Now().Format("2006-01-02 15:04:05"))
		t.tracker.UpdateMetadata(record.ID, "posted_platforms", strings.Join(success, ","))
	}

	// Construir respuesta
	if len(success) > 0 && len(errors) == 0 {
		return UserResult(fmt.Sprintf("🚀 **Publicación completada**\n\n✅ Éxito en: %s\nID Imagen: `%s`",
			strings.Join(success, ", "), record.ID))
	}

	if len(success) > 0 && len(errors) > 0 {
		return UserResult(fmt.Sprintf("⚠️ **Publicación parcial**\n\n✅ Éxito en: %s\n❌ Fallos en:\n%s\nID Imagen: `%s`",
			strings.Join(success, ", "),
			strings.Join(errors, "\n"),
			record.ID))
	}

	return ErrorResult(fmt.Sprintf("🛑 **Publicación fallida**\n\n❌ Errores:\n%s",
		strings.Join(errors, "\n")))
}

func (t *SocialManagerTool) generateDraftLocal(prompt, platform string) string {
	// Tono adaptado por plataforma
	switch platform {
	case "twitter", "x":
		return fmt.Sprintf(
			"✨ Look at this amazing AI creation! \n\n%s \n\n#AIArt #PicoClaw #Creative",
			utils.Truncate(prompt, 200),
		)
	case "facebook":
		return fmt.Sprintf(
			"🎨 Capturing the essence of digital art.\n\nInspired by: %s\n\nCreated with PicoClaw AI assistance. What do you think? \n\n#ArtificialIntelligence #DigitalArt #PicoClaw",
			prompt,
		)
	case "discord":
		return fmt.Sprintf(
			"🤖 **New Asset Generated**\n\n> %s\n\nCheck out the full res in the workspace! #picoclaw-art",
			prompt,
		)
	default:
		return prompt
	}
}

// findSimilarIDs encuentra IDs en el tracker que son similares al ID buscado
// Usa una heurística simple basada en:
// 1. Coincidencia parcial (substring)
// 2. Distancia de Levenshtein (caracteres diferentes)
// 3. Coincidencia de fecha (mismo día)
func (t *SocialManagerTool) findSimilarIDs(searchID string) []string {
	allRecords := t.tracker.List()
	var similar []string

	// Normalizar el ID buscado (quitar guiones bajos para comparación)
	searchNormalized := strings.ReplaceAll(searchID, "_", "")

	for _, record := range allRecords {
		recordID := record.ID
		recordNormalized := strings.ReplaceAll(recordID, "_", "")

		// Criterio 1: Coincidencia parcial (el buscado contiene parte del ID real o viceversa)
		if strings.Contains(recordNormalized, searchNormalized) ||
			strings.Contains(searchNormalized, recordNormalized) {
			similar = append(similar, recordID)
			continue
		}

		// Criterio 2: Mismo prefijo de fecha (AAAAMMDD)
		if len(searchID) >= 8 && len(recordID) >= 8 {
			searchDate := searchID[:8]
			recordDate := recordID[:8]

			if searchDate == recordDate {
				// Mismo día - verificar si los caracteres restantes son similares
				searchTime := searchID[9:]
				recordTime := recordID[9:]

				// Si al menos los primeros 4 caracteres del tiempo coinciden
				if len(searchTime) >= 4 && len(recordTime) >= 4 {
					if searchTime[:4] == recordTime[:4] {
						similar = append(similar, recordID)
						continue
					}
				}
			}
		}

		// Criterio 3: Distancia de Levenshtein <= 3 (máximo 3 caracteres diferentes)
		if levenshteinDistance(searchNormalized, recordNormalized) <= 3 {
			similar = append(similar, recordID)
		}
	}

	// Limitar a 5 sugerencias más relevantes
	if len(similar) > 5 {
		similar = similar[:5]
	}

	return similar
}

// levenshteinDistance calcula la distancia de Levenshtein entre dos strings
// (número mínimo de ediciones: inserciones, eliminaciones o sustituciones)
func levenshteinDistance(s1, s2 string) int {
	if len(s1) == 0 {
		return len(s2)
	}
	if len(s2) == 0 {
		return len(s1)
	}

	// Crear matriz de distancia
	matrix := make([][]int, len(s1)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(s2)+1)
		matrix[i][0] = i
	}
	for j := range matrix[0] {
		matrix[0][j] = j
	}

	// Llenar matriz
	for i := 1; i <= len(s1); i++ {
		for j := 1; j <= len(s2); j++ {
			cost := 1
			if s1[i-1] == s2[j-1] {
				cost = 0
			}
			matrix[i][j] = min(
				matrix[i-1][j]+1, // eliminación
				min(
					matrix[i][j-1]+1,      // inserción
					matrix[i-1][j-1]+cost, // sustitución
				),
			)
		}
	}

	return matrix[len(s1)][len(s2)]
}

// normalizeTrackerPath elimina sufijos de workers de la ruta del tracker
// Ejemplo: /home/.picoclaw/workspace-general_worker/image_gen -> /home/.picoclaw/workspace/image_gen
func (t *SocialManagerTool) normalizeTrackerPath(trackerPath string) string {
	if trackerPath == "" {
		return ""
	}

	// Obtener el directorio del tracker
	dir := filepath.Dir(trackerPath)

	// Patrón: si el path contiene "workspace-" seguido de un worker ID, eliminar el sufijo
	// Ejemplo: /home/.picoclaw/workspace-general_worker/image_gen
	parts := strings.Split(dir, string(filepath.Separator))

	for i, part := range parts {
		// Buscar la parte que contiene "workspace-"
		if strings.HasPrefix(part, "workspace-") {
			// Reemplazar "workspace-{id}" con "workspace"
			cleaned := "workspace"
			parts[i] = cleaned

			// Reconstruir la ruta
			normalized := filepath.Join(parts...)
			if !strings.HasPrefix(normalized, string(filepath.Separator)) &&
				strings.HasPrefix(dir, string(filepath.Separator)) {
				normalized = string(filepath.Separator) + normalized
			}

			logger.DebugF("Ruta normalizada (eliminado sufijo de worker)", map[string]any{
				"original":   dir,
				"normalized": normalized,
			})

			return normalized
		}
	}

	// Si no hay sufijo de worker, devolver la ruta original
	return dir
}

// min devuelve el mínimo de dos enteros
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
