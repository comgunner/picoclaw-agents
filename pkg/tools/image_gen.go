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
	"path/filepath"
	"strings"
	"time"

	"github.com/comgunner/picoclaw/pkg/utils"
)

// ============== Image Generation Tools ==============

type ImageGenCreateTool struct {
	provider           string
	geminiAPIKey       string
	geminiImageModel   string
	geminiTextModel    string
	ideogramAPIKey     string
	ideogramAPIURL     string
	aspectRatio        string
	outputDir          string
	tracker            *utils.ImageGenTracker
	imageScriptPath    string
	imageGenScriptPath string
}

func NewImageGenCreateTool() *ImageGenCreateTool {
	return NewImageGenCreateToolFromConfig("", "", "", "", "", "", "", "", "", "", "")
}

func NewImageGenCreateToolFromConfig(
	configProvider, configGeminiKey, configGeminiImageModel, configGeminiTextModel, configIdeogramKey, configIdeogramURL,
	configAspectRatio, configOutputDir, configImageScriptPath, configImageGenScriptPath, workspace string,
) *ImageGenCreateTool {
	provider := strings.TrimSpace(os.Getenv(utils.EnvImageGenProvider))
	geminiKey := strings.TrimSpace(os.Getenv(utils.EnvGeminiAPIKey))
	geminiImageModel := strings.TrimSpace(os.Getenv(utils.EnvGeminiImageModel))
	geminiTextModel := strings.TrimSpace(os.Getenv(utils.EnvGeminiTextModel))
	ideogramKey := strings.TrimSpace(os.Getenv(utils.EnvIdeogramAPIKey))
	ideogramURL := strings.TrimSpace(os.Getenv(utils.EnvIdeogramAPIURL))
	aspectRatio := strings.TrimSpace(os.Getenv(utils.EnvAspectRatio))
	outputDir := strings.TrimSpace(os.Getenv(utils.EnvImageGenOutputDir))
	imageScriptPath := strings.TrimSpace(os.Getenv(utils.EnvImageScriptPath))
	imageGenScriptPath := strings.TrimSpace(os.Getenv(utils.EnvImageGenScriptPath))

	if provider == "" {
		provider = strings.TrimSpace(configProvider)
	}
	if provider == "" {
		provider = "gemini" // Default
	}
	if geminiKey == "" {
		geminiKey = strings.TrimSpace(configGeminiKey)
	}
	if geminiImageModel == "" {
		geminiImageModel = strings.TrimSpace(configGeminiImageModel)
	}
	if geminiImageModel == "" {
		// Fallback model for Gemini generateContent image responses.
		geminiImageModel = "gemini-2.0-flash-exp" // More generic fallback
	}
	if geminiTextModel == "" {
		geminiTextModel = strings.TrimSpace(configGeminiTextModel)
	}
	if geminiTextModel == "" {
		geminiTextModel = "gemini-2.5-flash"
	}
	if ideogramKey == "" {
		ideogramKey = strings.TrimSpace(configIdeogramKey)
	}
	if ideogramURL == "" {
		ideogramURL = strings.TrimSpace(configIdeogramURL)
	}
	if aspectRatio == "" {
		aspectRatio = strings.TrimSpace(configAspectRatio)
	}
	if aspectRatio == "" {
		aspectRatio = "4:5" // Default
	}
	if outputDir == "" {
		outputDir = strings.TrimSpace(configOutputDir)
	}
	outputDir = resolveOutputDir(outputDir, workspace)
	if imageScriptPath == "" {
		imageScriptPath = strings.TrimSpace(configImageScriptPath)
	}
	if imageGenScriptPath == "" {
		imageGenScriptPath = strings.TrimSpace(configImageGenScriptPath)
	}
	imageScriptPath = resolvePathInWorkspace(imageScriptPath, workspace)
	imageGenScriptPath = resolvePathInWorkspace(imageGenScriptPath, workspace)

	// Initialize tracker
	trackerPath := filepath.Join(outputDir, "tracker.json")
	tracker, _ := utils.NewImageGenTracker(trackerPath)

	return &ImageGenCreateTool{
		provider:           provider,
		geminiAPIKey:       geminiKey,
		geminiImageModel:   geminiImageModel,
		geminiTextModel:    geminiTextModel,
		ideogramAPIKey:     ideogramKey,
		ideogramAPIURL:     ideogramURL,
		aspectRatio:        aspectRatio,
		outputDir:          outputDir,
		tracker:            tracker,
		imageScriptPath:    imageScriptPath,
		imageGenScriptPath: imageGenScriptPath,
	}
}

func (t *ImageGenCreateTool) GetTracker() *utils.ImageGenTracker {
	return t.tracker
}

func (t *ImageGenCreateTool) Name() string {
	return "image_gen_create"
}

func (t *ImageGenCreateTool) Description() string {
	return "Generar imagen desde prompt de texto usando Gemini o Ideogram. Para solicitudes directas de imagen, usar esta tool sin pasar por text_script_create."
}

func (t *ImageGenCreateTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"prompt": map[string]any{
				"type":        "string",
				"description": "Prompt de texto para generar la imagen",
			},
			"provider": map[string]any{
				"type":        "string",
				"description": "Proveedor: 'gemini' o 'ideogram' (default: config)",
				"enum":        []string{"gemini", "ideogram"},
			},
			"aspect_ratio": map[string]any{
				"type":        "string",
				"description": "Aspect ratio: '4:5', '16:9', '1:1', etc. (default: config)",
			},
			"script_path": map[string]any{
				"type":        "string",
				"description": "Ruta al archivo de guion generado previamente. ÚSALO SIEMPRE si vienes de 'text_script_create' para que la imagen se guarde en la misma carpeta.",
			},
		},
		"required": []string{"prompt"},
	}
}

func (t *ImageGenCreateTool) Execute(ctx context.Context, args map[string]any) *ToolResult {
	prompt, _ := args["prompt"].(string)
	provider, _ := args["provider"].(string)
	aspectRatio, _ := args["aspect_ratio"].(string)
	scriptPath, _ := args["script_path"].(string)

	// Validar parámetros requeridos
	prompt = strings.TrimSpace(prompt)
	if prompt == "" {
		return ErrorResult("prompt es requerido")
	}

	// Usar defaults de config si no se proporcionan
	if provider == "" {
		provider = t.provider
	}
	if aspectRatio == "" {
		aspectRatio = t.aspectRatio
	} else if t.aspectRatio != "" && aspectRatio != t.aspectRatio && !promptMentionsAspectRatio(prompt, aspectRatio) {
		aspectRatio = t.aspectRatio
	}

	// 1. Generar ID único ANTES de nada
	id := utils.GenerateID()

	// 2. Crear Tracker Entry INMEDIATAMENTE (Estado: pending)
	// Esto asegura trazabilidad incluso si falla el setup posterior.
	if t.tracker != nil {
		t.tracker.Add(utils.ImageGenRecord{
			ID:          id,
			DateTime:    time.Now().Format("2006-01-02 15:04:05"),
			Prompt:      prompt,
			Provider:    provider,
			AspectRatio: aspectRatio,
			Model:       fmt.Sprintf("%s-image", provider),
			Language:    utils.DetectLanguage(prompt),
			Metadata: map[string]string{
				"status": "pending",
			},
		})
	}

	// Validar API keys
	if provider == "gemini" && t.geminiAPIKey == "" {
		if t.tracker != nil {
			t.tracker.UpdateMetadata(id, "status", "failed")
			t.tracker.UpdateMetadata(id, "error", "Gemini API Key no configurada")
		}
		return UserResult("Gemini API Key no configurada. Configure en config.json.")
	}
	if provider == "ideogram" && t.ideogramAPIKey == "" {
		if t.tracker != nil {
			t.tracker.UpdateMetadata(id, "status", "failed")
			t.tracker.UpdateMetadata(id, "error", "Ideogram API Key no configurada")
		}
		return UserResult("Ideogram API Key no configurada. Configure en config.json.")
	}

	// Asegurar output directory
	if err := os.MkdirAll(t.outputDir, 0o755); err != nil {
		if t.tracker != nil {
			t.tracker.UpdateMetadata(id, "status", "failed")
			t.tracker.UpdateMetadata(id, "error", err.Error())
		}
		return ErrorResult(fmt.Sprintf("error creando directorio: %v", err))
	}

	callCtx, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	// 3. Update Status: in_progress
	if t.tracker != nil {
		t.tracker.UpdateMetadata(id, "status", "in_progress")
	}

	imageDir := filepath.Join(t.outputDir, id)

	// Si hay script_path, leer script y generar prompt visual
	var finalPrompt string
	var scriptContent string
	if scriptPath != "" {
		data, err := os.ReadFile(scriptPath)
		if err != nil {
			if t.tracker != nil {
				t.tracker.UpdateMetadata(id, "status", "failed")
				t.tracker.UpdateMetadata(id, "error", err.Error())
			}
			return ErrorResult(fmt.Sprintf("error leyendo script: %v", err))
		}
		scriptContent = string(data)
		scriptDir := filepath.Dir(scriptPath)
		if strings.TrimSpace(scriptDir) != "" && scriptDir != "." {
			imageDir = scriptDir
			if base := filepath.Base(scriptDir); strings.TrimSpace(base) != "" {
				id = base
				// El ID ya está en el tracker pero lo sobreescribimos con el del script para consistencia?
				// En este punto es mejor usar el ID del script si existe.
			}
		}

		// Generar prompt visual desde script
		finalPrompt, err = utils.BuildVisualPromptFromScript(
			callCtx, t.geminiAPIKey, t.geminiTextModel, scriptContent, prompt, aspectRatio, "en", t.imageGenScriptPath)
		if err != nil {
			if t.tracker != nil {
				t.tracker.UpdateMetadata(id, "status", "failed")
				t.tracker.UpdateMetadata(id, "error", err.Error())
			}
			return ErrorResult(fmt.Sprintf("error generando prompt visual: %v", err))
		}
	} else {
		finalPrompt = prompt
	}

	// 4. Generar imagen
	imagePath := filepath.Join(imageDir, fmt.Sprintf("%s.-imagen.jpg", id))
	var err error

	if provider == "gemini" {
		req := utils.GeminiImageRequest{
			Prompt:      finalPrompt,
			AspectRatio: aspectRatio,
			Model:       t.geminiImageModel,
			APIKey:      t.geminiAPIKey,
		}
		err = utils.GenerateImageWithGemini(callCtx, req, imagePath)
		if err != nil && isGeminiModelNotFoundError(err) && t.ideogramAPIKey != "" {
			ideogramCfg := utils.IdeogramV3Config{
				APIKey:         t.ideogramAPIKey,
				APIURL:         t.ideogramAPIURL,
				AspectRatio:    aspectRatio,
				RenderingSpeed: "TURBO",
				StyleType:      "REALISTIC",
				NumImages:      1,
			}
			err = utils.GenerateImageWithIdeogram(callCtx, ideogramCfg, finalPrompt, imagePath)
			if err == nil {
				provider = "ideogram"
				if t.tracker != nil {
					t.tracker.UpdateMetadata(id, "provider", "ideogram")
				}
			}
		}
	} else {
		ideogramCfg := utils.IdeogramV3Config{
			APIKey:         t.ideogramAPIKey,
			APIURL:         t.ideogramAPIURL,
			AspectRatio:    aspectRatio,
			RenderingSpeed: "TURBO",
			StyleType:      "REALISTIC",
			NumImages:      1,
		}
		err = utils.GenerateImageWithIdeogram(callCtx, ideogramCfg, finalPrompt, imagePath)
	}

	if err != nil {
		if t.tracker != nil {
			t.tracker.UpdateMetadata(id, "status", "failed")
			t.tracker.UpdateMetadata(id, "error", err.Error())
		}
		return ErrorResult(fmt.Sprintf("image_gen_create falló: %v", err)).WithError(err)
	}

	// 5. Guardar archivos de soporte
	promptPath := filepath.Join(imageDir, fmt.Sprintf("%s.-prompt_visual.txt", id))
	_ = os.WriteFile(promptPath, []byte(finalPrompt), 0o644)

	var savedScriptPath string
	if scriptContent != "" {
		if scriptPath != "" {
			savedScriptPath = scriptPath
		} else {
			savedScriptPath = filepath.Join(imageDir, fmt.Sprintf("%s.-script.txt", id))
			_ = os.WriteFile(savedScriptPath, []byte(scriptContent), 0o644)
		}
	}

	// 6. Actualizar Tracker con Resultado Exitoso
	if t.tracker != nil {
		t.tracker.UpdateMetadata(id, "status", "generated")
		t.tracker.UpdateMetadata(id, "image_path", imagePath)
		t.tracker.UpdateMetadata(id, "prompt_path", promptPath)
		if savedScriptPath != "" {
			t.tracker.UpdateMetadata(id, "script_path", savedScriptPath)
		}
		if finalPrompt != prompt {
			t.tracker.UpdateMetadata(id, "visual_prompt", finalPrompt)
		}
	}

	// 7. Devolver Resultado DETALLADO (Feedback Visual Explícito)
	trackerFile := "tracker.json"
	if t.tracker != nil {
		trackerFile = t.tracker.TrackerPath
	}

	return UserResult(fmt.Sprintf(`✅ **Imagen generada exitosamente**

📁 **Archivos creados:**
┌─────────────────────────────────────────────────────────────────┐
│ Imagen:   %s
│ Prompt:   %s
│ Script:   %s
└─────────────────────────────────────────────────────────────────┘

🎨 **Prompt visual usado:** 
%s

📊 **Estado registrado en:** 
%s (status: generated)

**¿Qué quieres hacer con esta imagen?**
[1] ✅ **Me gusta** - Aprobada para publicación
[2] 🔄 **No me gusta** - Regenerar (límite 3)
[3] 📝 **Generar texto** - Crear copy para redes
[4] ⏰ **Programar** - Publicar más tarde
[5] ❌ **Finalizar** - Terminar sin publicar`,
		imagePath, promptPath, savedScriptPath, finalPrompt, trackerFile))
}

func isGeminiModelNotFoundError(err error) bool {
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "not_found") ||
		(strings.Contains(msg, "not found") && strings.Contains(msg, "models/"))
}

func promptMentionsAspectRatio(prompt, ratio string) bool {
	p := strings.ToLower(prompt)
	r := strings.ToLower(strings.TrimSpace(ratio))
	if r != "" && strings.Contains(p, r) {
		return true
	}
	keywords := []string{
		"aspect ratio", "relación de aspecto", "formato",
		"16:9", "9:16", "4:5", "1:1",
		"vertical", "portrait", "retrato",
		"horizontal", "landscape", "panoramic", "widescreen",
		"square", "cuadrado",
	}
	for _, k := range keywords {
		if strings.Contains(p, k) {
			return true
		}
	}
	return false
}

func resolveOutputDir(raw, workspace string) string {
	ws := strings.TrimSpace(workspace)
	if ws == "" {
		ws = "./workspace"
	}
	ws = expandHomePath(ws)
	if ws == "" {
		ws = "./workspace"
	}
	if raw == "" {
		return filepath.Join(ws, "image_gen")
	}
	resolved := resolvePathInWorkspace(raw, ws)
	if resolved == "" {
		return filepath.Join(ws, "image_gen")
	}
	return resolved
}

func resolvePathInWorkspace(raw, workspace string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	p := expandHomePath(raw)
	if filepath.IsAbs(p) {
		return filepath.Clean(p)
	}
	ws := expandHomePath(strings.TrimSpace(workspace))
	if ws == "" {
		return filepath.Clean(p)
	}
	clean := filepath.Clean(strings.TrimPrefix(p, "./"))
	if clean == "workspace" {
		return filepath.Clean(ws)
	}
	if strings.HasPrefix(clean, "workspace/") {
		return filepath.Join(ws, strings.TrimPrefix(clean, "workspace/"))
	}
	return filepath.Join(ws, clean)
}

func expandHomePath(p string) string {
	p = strings.TrimSpace(p)
	if p == "" || p[0] != '~' {
		return p
	}
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return p
	}
	if p == "~" {
		return home
	}
	if len(p) > 1 && p[1] == '/' {
		return filepath.Join(home, p[2:])
	}
	return p
}
