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
	"time"

	"github.com/comgunner/picoclaw/pkg/utils"
)

// ============== Discord Webhook Tool ==============

type DiscordWebhookTool struct {
	defaultWebhookURL string
}

func NewDiscordWebhookTool() *DiscordWebhookTool {
	return NewDiscordWebhookToolFromConfig("")
}

func NewDiscordWebhookToolFromConfig(configWebhookURL string) *DiscordWebhookTool {
	webhookURL := strings.TrimSpace(os.Getenv(utils.EnvDiscordWebhookURL))
	if webhookURL == "" {
		webhookURL = strings.TrimSpace(configWebhookURL)
	}
	return &DiscordWebhookTool{
		defaultWebhookURL: webhookURL,
	}
}

func (t *DiscordWebhookTool) Name() string {
	return "discord_post"
}

func (t *DiscordWebhookTool) Description() string {
	return "Publicar mensaje en Discord vía webhook. Soporta texto simple y texto con imagen adjunta."
}

func (t *DiscordWebhookTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"webhook_url": map[string]any{
				"type":        "string",
				"description": "URL del webhook de Discord. Si se omite, usa el default del config o variable de entorno.",
			},
			"message": map[string]any{
				"type":        "string",
				"description": "Contenido del mensaje a publicar",
			},
			"image_path": map[string]any{
				"type":        "string",
				"description": "Ruta absoluta de la imagen a adjuntar (opcional)",
			},
			"username": map[string]any{
				"type":        "string",
				"description": "Nombre que aparece como author del mensaje (opcional, default: PicoClaw)",
			},
		},
		"required": []string{"message"},
	}
}

func (t *DiscordWebhookTool) Execute(ctx context.Context, args map[string]any) *ToolResult {
	webhookURL, _ := args["webhook_url"].(string)
	message, _ := args["message"].(string)
	imagePath, _ := args["image_path"].(string)
	username, _ := args["username"].(string)

	// Usar default si no se proportionó webhook_url
	if webhookURL == "" {
		webhookURL = t.defaultWebhookURL
	}

	// Validar parámetros requeridos
	message = strings.TrimSpace(message)
	if message == "" {
		return ErrorResult("message es requerido")
	}

	// Validar webhook URL
	if webhookURL == "" {
		return UserResult(
			"Discord Webhook URL no configurada. " +
				"Configura en config.json (tools.social_media.discord) o usa variable de entorno:\n" +
				"  DISCORD_WEBHOOK_URL\n" +
				"También puedes pasar webhook_url como parámetro.",
		)
	}

	// Username por defecto
	if username == "" {
		username = "PicoClaw"
	}

	callCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	req := utils.DiscordWebhookRequest{
		WebhookURL: webhookURL,
		Message:    message,
		ImagePath:  strings.TrimSpace(imagePath),
		Username:   username,
	}

	err := utils.DiscordWebhookPost(callCtx, req)
	if err != nil {
		return ErrorResult(fmt.Sprintf("discord post falló: %v", err)).WithError(err)
	}

	preview := message
	if len(preview) > 80 {
		preview = preview[:80] + "..."
	}
	if len(preview) > 0 {
		if idx := strings.Index(preview, "\n"); idx > 0 {
			preview = preview[:idx]
		}
	}

	return UserResult(fmt.Sprintf("Mensaje enviado a Discord exitosamente: %s", preview))
}
