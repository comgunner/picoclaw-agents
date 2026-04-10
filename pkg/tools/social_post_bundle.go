// PicoClaw - Ultra-lightweight personal AI agent
package tools

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/comgunner/picoclaw/pkg/bus"
	"github.com/comgunner/picoclaw/pkg/logger"
	"github.com/comgunner/picoclaw/pkg/utils"
)

// SocialPostBundleTool is a macro-tool that generates both a text script and an image
// in a single orchestrated execution to reduce LLM round-trips.
type SocialPostBundleTool struct {
	geminiAPIKey       string
	geminiTextModel    string
	geminiImageModel   string
	ideogramAPIKey     string
	ideogramAPIURL     string
	outputDir          string
	aspectRatio        string
	imageScriptPath    string
	imageGenScriptPath string
	tracker            *utils.ImageGenTracker
	queue              *QueueManager
	bus                *bus.MessageBus
	channel            string
	chatID             string
	imageGenTool       *ImageGenAntigravityTool // ← Antigravity OAuth for images
}

func NewSocialPostBundleTool(
	geminiKey, geminiTextModel, geminiImageModel, ideogramKey, ideogramURL,
	aspectRatio, outputDir, imageScriptPath, imageGenScriptPath, workspace string,
	queue *QueueManager, bus *bus.MessageBus, tracker *utils.ImageGenTracker,
	imageGenTool *ImageGenAntigravityTool,
) *SocialPostBundleTool {
	if geminiTextModel == "" {
		geminiTextModel = "gemini-2.5-flash"
	}
	if geminiImageModel == "" {
		geminiImageModel = "gemini-2.5-flash-image"
	}
	if aspectRatio == "" {
		aspectRatio = "4:5"
	}

	resolvedOutputDir := resolveOutputDir(outputDir, workspace)
	resolvedImageScriptPath := resolvePathInWorkspace(imageScriptPath, workspace)
	resolvedImageGenScriptPath := resolvePathInWorkspace(imageGenScriptPath, workspace)

	return &SocialPostBundleTool{
		geminiAPIKey:       geminiKey,
		geminiTextModel:    geminiTextModel,
		geminiImageModel:   geminiImageModel,
		ideogramAPIKey:     ideogramKey,
		ideogramAPIURL:     ideogramURL,
		outputDir:          resolvedOutputDir,
		aspectRatio:        aspectRatio,
		imageScriptPath:    resolvedImageScriptPath,
		imageGenScriptPath: resolvedImageGenScriptPath,
		tracker:            tracker,
		queue:              queue,
		bus:                bus,
		imageGenTool:       imageGenTool,
	}
}

func (t *SocialPostBundleTool) SetContext(channel, chatID string) {
	t.channel = channel
	t.chatID = chatID
}

func (t *SocialPostBundleTool) Name() string {
	return "social_post_bundle"
}

func (t *SocialPostBundleTool) Description() string {
	return "Macro-Tool: Genera un post completo (Guion + Imagen) en background. Ahorra un 90% de tiempo y tokens comparado con el flujo manual. Úsala SIEMPRE para publicaciones sociales."
}

func (t *SocialPostBundleTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"topic": map[string]any{
				"type":        "string",
				"description": "Tema central del post y la imagen",
			},
			"category": map[string]any{
				"type":        "string",
				"description": "Categoría: 'news', 'history', 'mystery', etc.",
			},
			"language": map[string]any{
				"type":        "string",
				"description": "Idioma: 'en', 'es' (por defecto auto-detectado)",
			},
			"provider": map[string]any{
				"type":        "string",
				"description": "Proveedor de imagen: 'gemini' o 'ideogram'",
				"enum":        []string{"gemini", "ideogram"},
			},
		},
		"required": []string{"topic"},
	}
}

func (t *SocialPostBundleTool) Execute(ctx context.Context, args map[string]any) *ToolResult {
	topic, _ := args["topic"].(string)
	category, _ := args["category"].(string)
	language, _ := args["language"].(string)
	provider, _ := args["provider"].(string)

	topic = strings.TrimSpace(topic)
	if topic == "" {
		return ErrorResult("topic es requerido")
	}

	if t.queue == nil {
		return ErrorResult("QueueManager no inicializado")
	}

	// Crear ID único legibe
	taskID := t.queue.AddTask("POST_BUNDLE", args)

	// Iniciar proceso en background (Batch)
	go func() {
		logger.InfoF("SocialPostBundle: Starting background generation", map[string]any{
			"task_id":  taskID,
			"topic":    topic,
			"category": category,
			"provider": provider,
		})

		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		t.queue.UpdateStatus(taskID, StatusProcessing, nil)

		id := utils.GenerateID()
		idDir := filepath.Join(t.outputDir, id)
		_ = os.MkdirAll(idDir, 0o755)

		logger.InfoF("SocialPostBundle: Created output directory", map[string]any{
			"task_id": taskID,
			"id_dir":  idDir,
		})

		// Check cancellation before starting
		if t.queue != nil && t.queue.IsCancelled(taskID) {
			logger.InfoF("SocialPostBundle: Task canceled, stopping", map[string]any{"task_id": taskID})
			return
		}

		// 1. Generar Guion — PRIMARY: Antigravity OAuth con retry automático
		logger.InfoF("SocialPostBundle: Starting script generation", map[string]any{
			"task_id":  taskID,
			"topic":    topic,
			"language": language,
			"category": category,
		})

		req := utils.TextScriptRequest{
			Topic:        topic,
			Category:     category,
			Language:     language,
			TemplatePath: t.imageScriptPath,
		}

		// Try Antigravity OAuth with automatic retry on 429
		scriptRes, scriptErr := GenerateTextScriptAntigravity(bgCtx, t.geminiTextModel, req)
		if scriptErr != nil {
			// Check if it's a rate limit error — retry once more after delay
			if isTextScriptRateLimit(scriptErr) {
				logger.WarnF("SocialPostBundle: Rate limited, retrying once after 30s", map[string]any{
					"task_id": taskID,
				})
				select {
				case <-time.After(30 * time.Second):
					scriptRes, scriptErr = GenerateTextScriptAntigravity(bgCtx, t.geminiTextModel, req)
				case <-bgCtx.Done():
					scriptErr = fmt.Errorf("context canceled during retry wait: %w", bgCtx.Err())
				}
			}
		}

		var finalResult *ToolResult
		if scriptErr != nil {
			logger.ErrorF("SocialPostBundle: Script generation failed", map[string]any{
				"task_id": taskID,
				"error":   scriptErr.Error(),
			})
			finalResult = ErrorResult(t.buildUserFriendlyError(scriptErr, taskID))
		} else {
			logger.InfoF("SocialPostBundle: Script generated successfully", map[string]any{
				"task_id":           taskID,
				"script_length":     len(scriptRes.Script),
				"detected_language": scriptRes.Language,
			})

			scriptPath := filepath.Join(idDir, fmt.Sprintf("%s.-script.txt", id))
			_ = os.WriteFile(scriptPath, []byte(scriptRes.Script), 0o644)

			logger.InfoF("SocialPostBundle: Script saved to file", map[string]any{
				"task_id":     taskID,
				"script_path": scriptPath,
			})

			// 2. Generar Prompt Visual
			if t.queue != nil && t.queue.IsCancelled(taskID) {
				logger.InfoF("SocialPostBundle: Task canceled before visual prompt", map[string]any{"task_id": taskID})
				return
			}
			logger.InfoF("SocialPostBundle: Starting visual prompt generation", map[string]any{
				"task_id":       taskID,
				"script_length": len(scriptRes.Script),
				"aspect_ratio":  t.aspectRatio,
			})

			visualPrompt, err := utils.BuildVisualPromptFromScript(
				bgCtx,
				t.geminiAPIKey,
				t.geminiTextModel,
				scriptRes.Script,
				topic,
				t.aspectRatio,
				scriptRes.Language,
				t.imageGenScriptPath,
			)
			if err != nil {
				logger.WarnF("SocialPostBundle: Visual prompt generation failed, using fallback", map[string]any{
					"task_id":  taskID,
					"error":    err.Error(),
					"fallback": topic,
				})
				visualPrompt = topic
			} else {
				logger.InfoF("SocialPostBundle: Visual prompt generated successfully", map[string]any{
					"task_id":              taskID,
					"visual_prompt_length": len(visualPrompt),
				})
			}

			// 3. Generar Imagen — usa Antigravity OAuth (predeterminado)
			if t.queue != nil && t.queue.IsCancelled(taskID) {
				logger.InfoF(
					"SocialPostBundle: Task canceled before image generation",
					map[string]any{"task_id": taskID},
				)
				return
			}
			logger.InfoF("SocialPostBundle: Starting image generation", map[string]any{
				"task_id":       taskID,
				"provider":      "antigravity",
				"prompt_length": len(visualPrompt),
				"aspect_ratio":  t.aspectRatio,
			})

			imagePath := filepath.Join(idDir, fmt.Sprintf("%s.-imagen.jpg", id))
			var imageErr error

			// Use Antigravity OAuth tool for image generation
			if t.imageGenTool != nil {
				imgCtx, imgCancel := context.WithTimeout(bgCtx, 120*time.Second)
				defer imgCancel()

				imgResult := t.imageGenTool.Execute(imgCtx, map[string]any{
					"prompt":       visualPrompt,
					"aspect_ratio": t.aspectRatio,
				})

				if imgResult.IsError {
					imageErr = fmt.Errorf("%s", imgResult.ForLLM)
					logger.ErrorF("SocialPostBundle: Image generation failed", map[string]any{
						"task_id":  taskID,
						"provider": "antigravity",
						"error":    imgResult.ForLLM,
					})
				} else {
					// Extract image path from the result text
					generatedPath := extractImagePathFromResult(imgResult.ForLLM)
					if generatedPath != "" {
						// Copy image to the bundle's output directory
						expectedPath := filepath.Join(idDir, fmt.Sprintf("%s.-imagen.jpg", id))
						if copyErr := copyFile(generatedPath, expectedPath); copyErr == nil {
							imagePath = expectedPath
							logger.InfoF("SocialPostBundle: Image copied to bundle directory", map[string]any{
								"task_id":     taskID,
								"source":      generatedPath,
								"destination": expectedPath,
							})
						} else {
							logger.WarnF("SocialPostBundle: Failed to copy image, using original path", map[string]any{
								"task_id": taskID,
								"error":   copyErr.Error(),
							})
							imagePath = generatedPath
						}
					}
					logger.InfoF("SocialPostBundle: Image generated successfully (Antigravity OAuth)", map[string]any{
						"task_id":    taskID,
						"image_path": imagePath,
					})
				}
			} else if provider == "ideogram" && t.ideogramAPIKey != "" {
				// Fallback: Ideogram API
				logger.InfoF("SocialPostBundle: Generating image with Ideogram (fallback)", map[string]any{
					"task_id":    taskID,
					"image_path": imagePath,
				})
				ideogramCfg := utils.IdeogramV3Config{
					APIKey:         t.ideogramAPIKey,
					APIURL:         t.ideogramAPIURL,
					AspectRatio:    t.aspectRatio,
					RenderingSpeed: "TURBO",
					StyleType:      "REALISTIC",
					NumImages:      1,
				}
				imageErr = utils.GenerateImageWithIdeogram(bgCtx, ideogramCfg, visualPrompt, imagePath)
				if imageErr != nil {
					logger.ErrorF("SocialPostBundle: Image generation failed", map[string]any{
						"task_id":  taskID,
						"provider": "ideogram",
						"error":    imageErr.Error(),
					})
				} else {
					logger.InfoF("SocialPostBundle: Image generated successfully (Ideogram)", map[string]any{
						"task_id":    taskID,
						"image_path": imagePath,
					})
				}
			} else {
				// No image generation method available
				imageErr = fmt.Errorf("no image generation method configured")
				logger.WarnF("SocialPostBundle: No image generation method available", map[string]any{
					"task_id": taskID,
				})
			}

			if imageErr != nil {
				finalResult = ErrorResult(t.buildUserFriendlyError(imageErr, taskID))
			} else {
				// Check cancellation before saving
				if t.queue != nil && t.queue.IsCancelled(taskID) {
					logger.InfoF("SocialPostBundle: Task canceled before saving", map[string]any{"task_id": taskID})
					return
				}
				// 4. Éxito - Guardar en tracker
				logger.InfoF("SocialPostBundle: Saving record to tracker", map[string]any{
					"task_id": taskID,
					"id":      id,
				})

				if t.tracker != nil {
					metadata := map[string]string{
						"status":     "generated",
						"batch_id":   taskID,
						"draft_text": scriptRes.Script, // Save complete script to metadata
					}

					err = t.tracker.Add(utils.ImageGenRecord{
						ID:          id,
						DateTime:    time.Now().Format("2006-01-02 15:04:05"),
						Prompt:      topic,
						Provider:    "antigravity",
						ScriptPath:  scriptPath,
						AspectRatio: t.aspectRatio,
						Model:       "gemini-3.1-flash-image",
						Language:    scriptRes.Language,
						Metadata:    metadata,
					})

					if err != nil {
						logger.ErrorF("SocialPostBundle: Failed to save image record to tracker", map[string]any{
							"task_id":     taskID,
							"id":          id,
							"script_path": scriptPath,
							"image_path":  imagePath,
							"error":       err.Error(),
						})
					} else {
						logger.InfoF("SocialPostBundle: Successfully saved image record to tracker", map[string]any{
							"task_id":     taskID,
							"id":          id,
							"script_path": scriptPath,
							"image_path":  imagePath,
						})
					}
				} else {
					logger.WarnF("SocialPostBundle: Tracker is nil, skipping record save", map[string]any{
						"task_id": taskID,
					})
				}

				finalResult = UserResult(
					fmt.Sprintf(
						"✅ **Post %s Listo**\n\n📝 **Guion:**\n```\n%s\n```\n\n🖼 **Imagen:** %s\n\n💡 **Opciones (Copia y pega):**\n1) `/bundle_approve id=%s`\n2) `/bundle_regen id=%s`\n3) `/bundle_edit id=%s`\n4) `/bundle_cancel id=%s`",
						taskID,
						scriptRes.Script,
						imagePath,
						id,
						id,
						id,
						id,
					),
				)
			}
		}

		status := StatusCompleted
		if finalResult.IsError {
			status = StatusFailed
		}

		logger.InfoF("SocialPostBundle: Updating queue status", map[string]any{
			"task_id":  taskID,
			"status":   status,
			"is_error": finalResult.IsError,
		})

		t.queue.UpdateStatus(taskID, status, finalResult)

		// NOTIFICACIÓN DIRECTA (Ahorro de Tokens)
		// Skip notification if task was canceled
		if t.queue.IsCancelled(taskID) {
			logger.InfoF("SocialPostBundle: Task canceled, skipping notification", map[string]any{"task_id": taskID})
			return
		}

		logger.InfoF("SocialPostBundle: Publishing notification via MessageBus", map[string]any{
			"task_id": taskID,
			"channel": t.channel,
			"chat_id": t.chatID,
		})

		if t.bus != nil && t.channel != "" && t.chatID != "" {
			var media []string
			if !finalResult.IsError {
				media = []string{filepath.Join(idDir, fmt.Sprintf("%s.-imagen.jpg", id))}
			}
			t.bus.PublishOutbound(bus.OutboundMessage{
				Channel: t.channel,
				ChatID:  t.chatID,
				Content: finalResult.ForUser,
				Buttons: finalResult.Buttons,
				Media:   media,
			})
			logger.InfoF("SocialPostBundle: Notification published successfully", map[string]any{
				"task_id": taskID,
			})
		} else {
			logger.WarnF("SocialPostBundle: Cannot publish notification, missing bus or context", map[string]any{
				"task_id":     taskID,
				"has_bus":     t.bus != nil,
				"has_channel": t.channel != "",
				"has_chat_id": t.chatID != "",
			})
		}

		logger.InfoF("SocialPostBundle: Background generation complete", map[string]any{
			"task_id": taskID,
			"status":  status,
		})
	}()

	return UserResult(
		fmt.Sprintf(
			"⏳ Tarea iniciada con ID: `%s`\nTe avisaré en cuanto el post esté listo para tu aprobación.",
			taskID,
		),
	)
}

// buildUserFriendlyError convierte errores técnicos en mensajes accionables
func (t *SocialPostBundleTool) buildUserFriendlyError(err error, taskID string) string {
	errStr := err.Error()

	// Error 403 - API Key faltante (Gemini u otros)
	if strings.Contains(errStr, "403") && strings.Contains(errStr, "PERMISSION_DENIED") {
		return fmt.Sprintf(`❌ Error: No se pudo generar la imagen.

🔑 Motivo: Falta configurar la API Key de Gemini o permisos denegados.

📝 Solución:
1. Obtén tu API Key gratis en: https://aistudio.google.com/app/apikey
2. Agrega tu key en ~/.picoclaw/config.json:

   {
     "tools": {
       "image_gen": {
         "provider": "gemini",
         "gemini_api_key": "TU_API_KEY_AQUI"
       }
     }
   }

3. Reinicia el agente y vuelve a intentar.

💡 ¿Necesitas ayuda? Revisa docs/IMAGE_GEN_util.md

ID de tarea: %s`, taskID)
	}

	// Error de rate limit (429)
	if strings.Contains(errStr, "429") {
		return fmt.Sprintf(`⏳ Error: Límite de peticiones alcanzado.

🕐 Motivo: Demasiadas solicitudes en poco tiempo.

💡 Solución:
- Espera unos minutos y vuelve a intentar

ID de tarea: %s`, taskID)
	}

	// Error genérico (fallback)
	return fmt.Sprintf(`❌ Error: No se pudo completar la tarea.

🔍 Detalle: %s

💡 Si el problema persiste, revisa los logs o consulta la documentación.

ID de tarea: %s`, errStr, taskID)
}

// extractImagePathFromResult extracts the image path from the tool result text.
func extractImagePathFromResult(text string) string {
	// Look for pattern: "Image:    /path/to/image.jpg"
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		if strings.Contains(line, "Image:") || strings.Contains(line, "imagen.jpg") {
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				path := strings.TrimSpace(parts[1])
				if strings.HasSuffix(path, ".jpg") || strings.HasSuffix(path, ".jpeg") ||
					strings.HasSuffix(path, ".png") {
					return path
				}
			}
		}
	}
	return ""
}

// copyFile copies a file from src to dst.
func copyFile(src, dst string) error {
	sourceFile, openErr := os.Open(src)
	if openErr != nil {
		return fmt.Errorf("open source: %w", openErr)
	}
	defer sourceFile.Close()

	// Ensure destination directory exists
	if mkdirErr := os.MkdirAll(filepath.Dir(dst), 0o755); mkdirErr != nil {
		return fmt.Errorf("create destination dir: %w", mkdirErr)
	}

	destFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("create destination: %w", err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return fmt.Errorf("copy content: %w", err)
	}
	return nil
}
