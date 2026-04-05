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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// GeminiImageRequest representa una petición para generar imagen
type GeminiImageRequest struct {
	Prompt      string
	AspectRatio string // "4:5", "16:9", "1:1", "9:16"
	Model       string // "gemini-2.5-flash-image"
	APIKey      string
}

// GeminiImageResponse representa la respuesta de Gemini para imágenes (Imagen 3.0)
type GeminiImageResponse struct {
	Predictions []GeminiImagePrediction `json:"predictions,omitempty"`
	Candidates  []GeminiImageCandidate  `json:"candidates,omitempty"` // fallback para formato antiguo
}

// GeminiImagePrediction representa una predicción de imagen (Imagen 3.0)
type GeminiImagePrediction struct {
	BytesBase64Encoded string         `json:"bytesBase64Encoded"`
	Sigs               map[string]any `json:"sigs,omitempty"`
}

// GeminiImageCandidate representa un candidato de respuesta de imagen (formato antiguo)
type GeminiImageCandidate struct {
	Content GeminiImageContent `json:"content"`
}

// GeminiImageContent representa el contenido de imagen
type GeminiImageContent struct {
	Parts []GeminiImagePart `json:"parts"`
}

// GeminiImagePart representa una parte de imagen
type GeminiImagePart struct {
	InlineData *GeminiInlineData `json:"inlineData,omitempty"`
	Text       string            `json:"text,omitempty"`
}

// GeminiInlineData representa datos de imagen en línea
type GeminiInlineData struct {
	MimeType string `json:"mimeType"`
	Data     string `json:"data"` // Base64 encoded
}

// GenerateImageWithGemini genera una imagen usando Gemini API.
// - Para modelos "imagen-*": usa endpoint :predict
// - Para modelos "gemini-*": usa endpoint :generateContent con salida de imagen inline
func GenerateImageWithGemini(ctx context.Context, req GeminiImageRequest, outputPath string) error {
	// Construir prompt con aspect ratio
	prompt := req.Prompt
	if !strings.Contains(strings.ToLower(prompt), "aspect ratio") {
		prompt = fmt.Sprintf("%s . Image aspect ratio %s.", prompt, req.AspectRatio)
	}

	// Asegurar que no haya texto en la imagen
	if !strings.Contains(strings.ToLower(prompt), "no text") {
		prompt += " . No text, no typography, no watermarks."
	}

	model := req.Model
	if model == "" {
		model = "gemini-2.5-flash-image-preview"
	}

	var imageData []byte
	var err error
	if strings.HasPrefix(strings.ToLower(model), "imagen-") {
		imageData, err = generateImageWithImagenPredict(ctx, req.APIKey, model, prompt, req.AspectRatio)
	} else {
		imageData, err = generateImageWithGeminiGenerateContent(ctx, req.APIKey, model, prompt, req.AspectRatio)
	}
	if err != nil {
		return err
	}

	// Guardar imagen
	if err := saveImage(imageData, outputPath); err != nil {
		return fmt.Errorf("error saving image: %v", err)
	}

	return nil
}

func generateImageWithImagenPredict(ctx context.Context, apiKey, model, prompt, aspectRatio string) ([]byte, error) {
	maxRetries := 3
	currentPrompt := prompt

	for attempt := 1; attempt <= maxRetries; attempt++ {
		apiRequest := map[string]any{
			"instances": []map[string]any{
				{
					"prompt": currentPrompt,
				},
			},
			"parameters": map[string]any{
				"sampleCount":   1,
				"aspectRatio":   aspectRatio,
				"safetySetting": "BLOCK_NONE",
			},
		}

		jsonData, err := json.Marshal(apiRequest)
		if err != nil {
			return nil, fmt.Errorf("error marshaling request: %v", err)
		}

		url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:predict?key=%s", model, apiKey)
		httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonData))
		if err != nil {
			return nil, fmt.Errorf("error creating request: %v", err)
		}
		httpReq.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 120 * time.Second}
		resp, err := client.Do(httpReq)
		if err != nil {
			if attempt < maxRetries {
				time.Sleep(time.Duration(attempt*2) * time.Second)
				continue
			}
			return nil, fmt.Errorf("error calling Gemini API: %v", err)
		}

		body, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			return nil, fmt.Errorf("error reading response: %v", readErr)
		}

		if resp.StatusCode != http.StatusOK {
			bodyStr := strings.ToLower(string(body))
			if attempt < maxRetries && (strings.Contains(bodyStr, "safety") || strings.Contains(bodyStr, "filter")) {
				currentPrompt = sanitizePrompt(prompt)
				time.Sleep(time.Duration(attempt*2) * time.Second)
				continue
			}
			return nil, fmt.Errorf("Gemini API error (%d): %s", resp.StatusCode, string(body))
		}

		var apiResp GeminiImageResponse
		if err := json.Unmarshal(body, &apiResp); err != nil {
			return nil, fmt.Errorf("error unmarshaling response: %v", err)
		}

		if len(apiResp.Predictions) > 0 {
			imageData, err := base64.StdEncoding.DecodeString(apiResp.Predictions[0].BytesBase64Encoded)
			if err != nil {
				return nil, fmt.Errorf("error decoding image data: %v", err)
			}
			return imageData, nil
		}

		if len(apiResp.Candidates) > 0 {
			for _, part := range apiResp.Candidates[0].Content.Parts {
				if part.InlineData != nil && strings.HasPrefix(part.InlineData.MimeType, "image/") {
					imageData, err := base64.StdEncoding.DecodeString(part.InlineData.Data)
					if err != nil {
						return nil, fmt.Errorf("error decoding image: %v", err)
					}
					return imageData, nil
				}
			}
		}

		if attempt < maxRetries {
			time.Sleep(time.Duration(attempt*2) * time.Second)
		}
	}

	return nil, fmt.Errorf("no image data in response")
}

func generateImageWithGeminiGenerateContent(
	ctx context.Context,
	apiKey, model, prompt, aspectRatio string,
) ([]byte, error) {
	maxRetries := 3
	currentPrompt := prompt
	includeImageConfig := strings.TrimSpace(aspectRatio) != "" && strings.TrimSpace(aspectRatio) != "1:1"

	for attempt := 1; attempt <= maxRetries; attempt++ {
		generationConfig := map[string]any{
			"responseModalities": []string{"TEXT", "IMAGE"},
		}
		if includeImageConfig {
			generationConfig["imageConfig"] = map[string]any{
				"aspectRatio": aspectRatio,
			}
		}

		apiRequest := map[string]any{
			"contents": []map[string]any{
				{
					"parts": []map[string]any{
						{"text": currentPrompt},
					},
				},
			},
			"generationConfig": generationConfig,
		}

		jsonData, err := json.Marshal(apiRequest)
		if err != nil {
			return nil, fmt.Errorf("error marshaling request: %v", err)
		}

		url := fmt.Sprintf(
			"https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s",
			model,
			apiKey,
		)
		httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonData))
		if err != nil {
			return nil, fmt.Errorf("error creating request: %v", err)
		}
		httpReq.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 120 * time.Second}
		resp, err := client.Do(httpReq)
		if err != nil {
			if attempt < maxRetries {
				time.Sleep(time.Duration(attempt*2) * time.Second)
				continue
			}
			return nil, fmt.Errorf("error calling Gemini API: %v", err)
		}

		body, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			return nil, fmt.Errorf("error reading response: %v", readErr)
		}

		if resp.StatusCode != http.StatusOK {
			bodyStr := strings.ToLower(string(body))
			// Some deployments may not support generationConfig.imageConfig yet.
			// Retry without imageConfig instead of failing hard.
			if includeImageConfig && strings.Contains(bodyStr, "imageconfig") {
				includeImageConfig = false
				if attempt < maxRetries {
					time.Sleep(time.Duration(attempt) * time.Second)
					continue
				}
			}
			if attempt < maxRetries && (strings.Contains(bodyStr, "safety") || strings.Contains(bodyStr, "filter")) {
				currentPrompt = sanitizePrompt(prompt)
				time.Sleep(time.Duration(attempt*2) * time.Second)
				continue
			}
			return nil, fmt.Errorf("Gemini API error (%d): %s", resp.StatusCode, string(body))
		}

		var apiResp GeminiImageResponse
		if err := json.Unmarshal(body, &apiResp); err != nil {
			return nil, fmt.Errorf("error unmarshaling response: %v", err)
		}

		if len(apiResp.Candidates) == 0 {
			if attempt < maxRetries {
				time.Sleep(time.Duration(attempt*2) * time.Second)
				continue
			}
			return nil, fmt.Errorf("no candidates in Gemini response")
		}

		for _, part := range apiResp.Candidates[0].Content.Parts {
			if part.InlineData != nil && strings.HasPrefix(part.InlineData.MimeType, "image/") {
				imageData, err := base64.StdEncoding.DecodeString(part.InlineData.Data)
				if err != nil {
					return nil, fmt.Errorf("error decoding image: %v", err)
				}
				return imageData, nil
			}
		}

		if attempt < maxRetries {
			time.Sleep(time.Duration(attempt*2) * time.Second)
		}
	}

	return nil, fmt.Errorf("no image data in response")
}

// sanitizePrompt suaviza un prompt para evitar filtros de seguridad
func sanitizePrompt(prompt string) string {
	replacements := map[string]string{
		"blood":     "liquid",
		"kill":      "neutralize",
		"murder":    "eliminate",
		"terrorist": "extremist",
		"empire":    "civilization",
		"war":       "conflict",
		"weapon":    "tool",
		"gun":       "device",
		"bomb":      "explosive device",
		"attack":    "action",
	}

	result := prompt
	for bad, good := range replacements {
		result = strings.ReplaceAll(result, bad, good)
	}

	return result
}

// saveImage guarda datos de imagen a archivo
func saveImage(imageData []byte, outputPath string) error {
	// Asegurar que el directorio existe
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	return os.WriteFile(outputPath, imageData, 0o644)
}

// GenerateImageWithGeminiFromScript genera imagen desde un script
func GenerateImageWithGeminiFromScript(
	ctx context.Context,
	apiKey, script, topic, aspectRatio, outputPath string,
) error {
	// Usar un modelo por defecto para la generación del prompt si no se especifica
	model := "gemini-2.5-flash"

	// Primero generar prompt visual desde el script
	visualPrompt, err := BuildVisualPromptFromScript(ctx, apiKey, model, script, topic, aspectRatio, "en", "")
	if err != nil {
		return fmt.Errorf("error building visual prompt: %v", err)
	}

	// Luego generar imagen usando Imagen 3.0
	req := GeminiImageRequest{
		Prompt:      visualPrompt,
		AspectRatio: aspectRatio,
		Model:       "gemini-2.0-flash-exp-image-preview", // Fallback para imagen
		APIKey:      apiKey,
	}

	return GenerateImageWithGemini(ctx, req, outputPath)
}
