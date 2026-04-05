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
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/comgunner/picoclaw/pkg/utils"
)

// ============== Community Manager Tools ==============

type CommunityManagerTool struct {
	outputDir string
}

func NewCommunityManagerTool() *CommunityManagerTool {
	return NewCommunityManagerToolWithWorkspace("")
}

func NewCommunityManagerToolWithWorkspace(workspace string) *CommunityManagerTool {
	out := resolvePathInWorkspace("./text_scripts", workspace)
	if strings.TrimSpace(out) == "" {
		out = "./text_scripts"
	}
	return &CommunityManagerTool{outputDir: out}
}

func (t *CommunityManagerTool) Name() string {
	return "community_manager_create_draft"
}

func (t *CommunityManagerTool) Description() string {
	return "Crear borrador de comunicado público desde contenido técnico. Adapta tono para redes sociales (Discord, Twitter, Facebook). Soporta múltiples plataformas en una sola llamada."
}

func (t *CommunityManagerTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"raw_data": map[string]any{
				"type":        "string",
				"description": "Contenido técnico o resumen de cambios",
			},
			"platform": map[string]any{
				"type":        "string",
				"description": "Plataforma(s) de destino: 'discord', 'twitter', 'facebook', 'blog'. Separar por comas para múltiples (ej: 'discord,twitter,facebook')",
			},
			"include_emojis": map[string]any{
				"type":        "boolean",
				"description": "Incluir emojis (default: true)",
			},
		},
		"required": []string{"raw_data", "platform"},
	}
}

func (t *CommunityManagerTool) Execute(ctx context.Context, args map[string]any) *ToolResult {
	rawData, _ := args["raw_data"].(string)
	platformArg, _ := args["platform"].(string)
	includeEmojis, _ := args["include_emojis"].(bool)
	_ = args["image_path"] // Imagen opcional para adjuntar (no se usa aún)

	rawData = strings.TrimSpace(rawData)
	if rawData == "" {
		return ErrorResult("raw_data es requerido")
	}

	// Parsear múltiples plataformas (separadas por comas)
	platforms := parsePlatforms(platformArg)
	if len(platforms) == 0 {
		platforms = []string{"discord"} // Default
	}

	// Generar contenido base una sola vez
	baseContent := t.generateBaseContent(rawData, includeEmojis)

	// Generar ID único para esta petición (todas las plataformas van en misma carpeta)
	id := utils.GenerateID()

	var results []string
	var allDraftPaths []string
	for _, platform := range platforms {
		draft := t.adaptContentForPlatform(baseContent, platform, includeEmojis)
		savedPath, saveErr := t.saveDraftInDir(id, platform, rawData, draft)
		if saveErr != nil {
			return ErrorResult(fmt.Sprintf("error guardando borrador para %s: %v", platform, saveErr))
		}

		results = append(results, fmt.Sprintf("**%s**: %s", platform, savedPath))
		allDraftPaths = append(allDraftPaths, savedPath)
	}

	// Construir mensaje de respuesta
	response := fmt.Sprintf(
		"📝 Borradores generados para %d plataforma(s):\n\n%s\n\n📁 Guardados en: %s/%s",
		len(platforms), strings.Join(results, "\n"), t.outputDir, id)

	// Si hay imagen, NO la adjuntar aquí - se enviará en mensaje separado
	// El texto largo debe ir SIEMPRE sin imagen para evitar límite de caption
	result := UserResult(response)

	// Agregar paths al resultado para que el agente pueda leerlos
	result.ForLLM += fmt.Sprintf("\n\n📎 Paths: %s", strings.Join(allDraftPaths, ", "))

	return result
}

// parsePlatforms separa plataformas por coma
func parsePlatforms(platformArg string) []string {
	if platformArg == "" {
		return []string{}
	}
	parts := strings.Split(platformArg, ",")
	var platforms []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			platforms = append(platforms, p)
		}
	}
	return platforms
}

// generateBaseContent genera el contenido base una sola vez
func (t *CommunityManagerTool) generateBaseContent(rawData string, includeEmojis bool) string {
	// Contenido base sin adaptaciones específicas de plataforma
	return rawData
}

// adaptContentForPlatform adapta el contenido base para cada plataforma
func (t *CommunityManagerTool) adaptContentForPlatform(baseContent string, platform string, includeEmojis bool) string {
	emojis := map[string]string{
		"discord":  "🎉 🚀 💬",
		"twitter":  "📢 🚀",
		"facebook": "📱 ✨",
		"blog":     "📝 📖",
	}

	emojiSet := ""
	if includeEmojis {
		emojiSet = emojis[platform]
	}

	switch platform {
	case "discord":
		return fmt.Sprintf("%s ¡Nueva Actualización!\n\n%s\n\n#Actualización #PicoClaw",
			emojiSet, baseContent)
	case "twitter":
		// Twitter requiere < 280 caracteres - acortar del contenido base
		text := fmt.Sprintf("%s %s #PicoClaw", emojiSet, baseContent)
		if len(text) > 280 {
			text = text[:277] + "..."
		}
		return text
	case "facebook":
		return fmt.Sprintf("%s ¡Gran Noticia!\n\n%s\n\n#PicoClaw #Actualización",
			emojiSet, baseContent)
	case "blog":
		return fmt.Sprintf("# Anuncio de Actualización\n\n%s\n\n---\n*Publicado por el equipo de PicoClaw*",
			baseContent)
	default:
		return baseContent
	}
}

type communityDraftRecord struct {
	ID        string `json:"id"`
	DateTime  string `json:"date_time"`
	Platform  string `json:"platform"`
	RawData   string `json:"raw_data"`
	DraftPath string `json:"draft_path"`
	Language  string `json:"language"`
}

type communityDraftTracker struct {
	Records []communityDraftRecord `json:"records"`
}

func (t *CommunityManagerTool) saveDraftInDir(dirID, platform, rawData, draft string) (string, error) {
	if strings.TrimSpace(t.outputDir) == "" {
		t.outputDir = "./text_scripts"
	}
	if err := os.MkdirAll(t.outputDir, 0o755); err != nil {
		return "", err
	}

	// Usar el directorio existente (misma petición = misma carpeta)
	draftDir := filepath.Join(t.outputDir, dirID)
	if err := os.MkdirAll(draftDir, 0o755); err != nil {
		return "", err
	}

	draftPath := filepath.Join(draftDir, fmt.Sprintf("%s.-post_%s.txt", dirID, platform))
	if err := os.WriteFile(draftPath, []byte(draft), 0o644); err != nil {
		return "", err
	}

	// Registrar en tracker (cada plataforma se registra por separado)
	trackerPath := filepath.Join(t.outputDir, "tracker.json")
	if err := appendCommunityDraftRecord(trackerPath, communityDraftRecord{
		ID:        dirID,
		DateTime:  time.Now().Format("2006-01-02 15:04:05"),
		Platform:  platform,
		RawData:   rawData,
		DraftPath: draftPath,
		Language:  utils.DetectLanguage(rawData),
	}); err != nil {
		return "", err
	}

	return draftPath, nil
}

func appendCommunityDraftRecord(trackerPath string, record communityDraftRecord) error {
	var tracker communityDraftTracker
	if data, err := os.ReadFile(trackerPath); err == nil && len(data) > 0 {
		if err := json.Unmarshal(data, &tracker); err != nil {
			return err
		}
	}
	tracker.Records = append(tracker.Records, record)
	data, err := json.MarshalIndent(tracker, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(trackerPath, data, 0o644)
}

// CommunityFromImage genera texto para publicación desde una imagen
type CommunityFromImageTool struct{}

func NewCommunityFromImageTool() *CommunityFromImageTool {
	return &CommunityFromImageTool{}
}

func (t *CommunityFromImageTool) Name() string {
	return "community_from_image"
}

func (t *CommunityFromImageTool) Description() string {
	return "Generar texto atractivo para publicación en redes sociales basado en una imagen generada."
}

func (t *CommunityFromImageTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"image_path": map[string]any{
				"type":        "string",
				"description": "Ruta de la imagen generada",
			},
			"platform": map[string]any{
				"type":        "string",
				"description": "Plataforma: 'discord', 'twitter', 'facebook'",
			},
			"prompt_used": map[string]any{
				"type":        "string",
				"description": "Prompt usado para generar la imagen (para contexto)",
			},
		},
		"required": []string{"image_path"},
	}
}

func (t *CommunityFromImageTool) Execute(ctx context.Context, args map[string]any) *ToolResult {
	imagePath, _ := args["image_path"].(string)
	platform, _ := args["platform"].(string)
	promptUsed, _ := args["prompt_used"].(string)

	if imagePath == "" {
		return ErrorResult("image_path es requerido")
	}

	if platform == "" {
		platform = "discord"
	}

	// Generar texto desde el prompt (simulado - en producción usaría visión por computadora)
	text := t.generateFromPrompt(promptUsed, platform)

	response := fmt.Sprintf(
		"📝 Texto generado para imagen:\n\n%s\n\n📁 Imagen: %s\n\n💡 ¿Publicar ahora?",
		text,
		imagePath,
	)
	return UserResult(response)
}

func (t *CommunityFromImageTool) generateFromPrompt(prompt, platform string) string {
	// Simplificado: usa el prompt como base
	switch platform {
	case "twitter":
		text := fmt.Sprintf("🎨 Nueva imagen generada: %s #AI #PicoClaw", prompt)
		if len(text) > 280 {
			text = text[:277] + "..."
		}
		return text
	case "facebook":
		return fmt.Sprintf("✨ Nueva imagen creada con IA\n\n%s\n\n#PicoClaw #AIArt", prompt)
	default:
		return fmt.Sprintf("🎉 ¡Nueva imagen!\n\n%s\n\n#PicoClaw", prompt)
	}
}
