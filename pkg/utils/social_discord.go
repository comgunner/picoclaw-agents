// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	EnvDiscordWebhookURL = "DISCORD_WEBHOOK_URL"
)

// DiscordWebhookRequest define los parámetros para publicar en Discord
type DiscordWebhookRequest struct {
	WebhookURL string
	Message    string
	ImagePath  string
	Username   string
}

// DiscordWebhookResponse representa la respuesta de Discord
type DiscordWebhookResponse struct {
	Success bool
	Error   string
}

// DiscordWebhookPost publica un mensaje en Discord vía webhook
// Soporta texto simple y texto con imagen adjunta
// Diseñado para soft-fail: nunca lanza excepciones, retorna error
func DiscordWebhookPost(ctx context.Context, req DiscordWebhookRequest) error {
	webhookURL := req.WebhookURL
	if webhookURL == "" {
		webhookURL = strings.TrimSpace(os.Getenv(EnvDiscordWebhookURL))
	}

	if webhookURL == "" {
		return fmt.Errorf("discord webhook URL no configurada")
	}

	client := &http.Client{Timeout: 30 * time.Second}

	var err error
	if req.ImagePath != "" {
		err = discordPostWithImage(ctx, client, webhookURL, req.Message, req.ImagePath, req.Username)
	} else {
		err = discordPostTextOnly(ctx, client, webhookURL, req.Message, req.Username)
	}

	return err
}

// discordPostWithImage publica un mensaje con imagen adjunta
func discordPostWithImage(
	ctx context.Context,
	client *http.Client,
	webhookURL, message, imagePath, username string,
) error {
	file, err := os.Open(imagePath)
	if err != nil {
		return fmt.Errorf("error abriendo imagen: %v", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Agregar campo de texto
	if err := writer.WriteField("content", message); err != nil {
		return err
	}
	if username != "" {
		if err := writer.WriteField("username", username); err != nil {
			return err
		}
	}

	// Agregar archivo de imagen
	filename := filepath.Base(imagePath)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return err
	}
	if _, err := io.Copy(part, file); err != nil {
		return fmt.Errorf("error copiando archivo: %v", err)
	}

	if err := writer.Close(); err != nil {
		return err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", webhookURL, body)
	if err != nil {
		return fmt.Errorf("error creando request: %v", err)
	}
	httpReq.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("error de red: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error Discord (%d): %s", resp.StatusCode, string(bodyBytes)[:200])
	}

	return nil
}

// discordPostTextOnly publica solo texto sin imagen
func discordPostTextOnly(ctx context.Context, client *http.Client, webhookURL, message, username string) error {
	payload := map[string]string{
		"content": message,
	}
	if username != "" {
		payload["username"] = username
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error serializando payload: %v", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", webhookURL, bytes.NewReader(jsonData))
	if err != nil {
		return fmt.Errorf("error creando request: %v", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	client.Timeout = 30 * time.Second
	resp, err := client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("error de red: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error Discord (%d): %s", resp.StatusCode, string(bodyBytes)[:200])
	}

	return nil
}

// DiscordWebhookPostSimple función simplificada para pruebas rápidas
func DiscordWebhookPostSimple(webhookURL, message string) error {
	ctx := context.Background()
	return DiscordWebhookPost(ctx, DiscordWebhookRequest{
		WebhookURL: webhookURL,
		Message:    message,
		Username:   "PicoClaw",
	})
}
