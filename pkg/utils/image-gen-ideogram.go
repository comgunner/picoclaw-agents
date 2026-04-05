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
	"strconv"
	"strings"
	"time"

	"github.com/comgunner/picoclaw/pkg/config"
)

// IdeogramV3Config holds all configuration for Ideogram V3 API
type IdeogramV3Config struct {
	APIKey         string
	APIURL         string
	AspectRatio    string
	RenderingSpeed string
	StyleType      string
	NumImages      int
	NegativePrompt string
}

// IdeogramRequest represents a request to Ideogram API (legacy JSON format)
type IdeogramRequest struct {
	Prompt         string `json:"prompt"`
	AspectRatio    string `json:"aspect_ratio"`
	MagicPrompt    bool   `json:"magic_prompt_option"`
	Model          string `json:"model"`
	NegativePrompt string `json:"negative_prompt,omitempty"`
	Style          string `json:"style_type,omitempty"`
	NumImages      int    `json:"num_images,omitempty"`
}

// IdeogramResponse represents a response from Ideogram API
type IdeogramResponse struct {
	Success bool            `json:"success"`
	Images  []IdeogramImage `json:"images"`
	Error   *IdeogramError  `json:"error,omitempty"`
}

// IdeogramImage represents a generated image
type IdeogramImage struct {
	URL string `json:"url"`
}

// IdeogramError represents an Ideogram API error
type IdeogramError struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

// GenerateImageWithIdeogram generates an image using Ideogram API with full configuration
// Supports both V3 (multipart/form-data) and Legacy (JSON with num_images:1)
func GenerateImageWithIdeogram(ctx context.Context, cfg IdeogramV3Config, prompt, outputPath string) error {
	// Build enriched prompt
	enrichedPrompt := prompt
	if !strings.Contains(strings.ToLower(prompt), "aspect ratio") {
		enrichedPrompt = fmt.Sprintf("%s . Image aspect ratio %s.", prompt, cfg.AspectRatio)
	}

	// Ensure no text
	if !strings.Contains(strings.ToLower(prompt), "no text") {
		enrichedPrompt += " . No text, no typography, no watermarks."
	}

	// Detect if V3 API (URL with ideogram-v3) or Legacy
	isV3 := strings.Contains(cfg.APIURL, "ideogram-v3") || cfg.APIURL == ""

	if isV3 {
		return generateWithIdeogramV3(ctx, cfg, enrichedPrompt, outputPath)
	}
	return generateWithIdeogramLegacy(ctx, cfg, enrichedPrompt, outputPath)
}

// generateWithIdeogramV3 uses Ideogram V3 API with multipart/form-data and configurable parameters
func generateWithIdeogramV3(ctx context.Context, cfg IdeogramV3Config, prompt, outputPath string) error {
	// Build multipart/form-data request
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Required fields (matching Python implementation)
	_ = writer.WriteField("prompt", prompt)
	_ = writer.WriteField("aspect_ratio", cfg.AspectRatio)
	_ = writer.WriteField("rendering_speed", cfg.RenderingSpeed)
	_ = writer.WriteField("style_type", cfg.StyleType)
	_ = writer.WriteField("num_images", strconv.Itoa(cfg.NumImages))

	// Optional fields
	if cfg.NegativePrompt != "" {
		_ = writer.WriteField("negative_prompt", cfg.NegativePrompt)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("error closing multipart writer: %v", err)
	}

	// Use configured API URL or default to V3
	apiURL := cfg.APIURL
	if apiURL == "" {
		apiURL = "https://api.ideogram.ai/v1/ideogram-v3/generate"
	}

	// Generate with retries
	maxRetries := 3
	var imageData []byte

	for attempt := 1; attempt <= maxRetries; attempt++ {
		httpReq, err := http.NewRequestWithContext(ctx, "POST", apiURL, body)
		if err != nil {
			return fmt.Errorf("error creating request: %v", err)
		}
		httpReq.Header.Set("Api-Key", cfg.APIKey)
		httpReq.Header.Set("Content-Type", writer.FormDataContentType())

		client := &http.Client{Timeout: 120 * time.Second}
		resp, err := client.Do(httpReq)
		if err != nil {
			if attempt < maxRetries {
				time.Sleep(time.Duration(attempt*2) * time.Second)
				continue
			}
			return fmt.Errorf("error calling Ideogram V3 API: %v", err)
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error reading response: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			if attempt < maxRetries {
				time.Sleep(time.Duration(attempt*2) * time.Second)
				continue
			}
			return fmt.Errorf("Ideogram V3 API error (%d): %s", resp.StatusCode, string(respBody))
		}

		// V3 returns image directly in body
		imageData = respBody
		break
	}

	if imageData == nil {
		return fmt.Errorf("no image data in response")
	}

	// Save image
	if err := saveImage(imageData, outputPath); err != nil {
		return fmt.Errorf("error saving image: %v", err)
	}

	return nil
}

// generateWithIdeogramLegacy uses Ideogram Legacy API with JSON and configurable parameters
func generateWithIdeogramLegacy(ctx context.Context, cfg IdeogramV3Config, prompt, outputPath string) error {
	// Build JSON request
	req := IdeogramRequest{
		Prompt:         prompt,
		AspectRatio:    cfg.AspectRatio,
		MagicPrompt:    true,
		Model:          "V_2", // Default model for Legacy API
		NegativePrompt: cfg.NegativePrompt,
		Style:          cfg.StyleType,
		NumImages:      cfg.NumImages,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("error marshaling request: %v", err)
	}

	// Use configured API URL or default to Legacy
	apiURL := cfg.APIURL
	if apiURL == "" {
		apiURL = "https://api.ideogram.ai/generate"
	}

	// Generate with retries
	maxRetries := 3
	var imageURL string

	for attempt := 1; attempt <= maxRetries; attempt++ {
		httpReq, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewReader(jsonData))
		if err != nil {
			return fmt.Errorf("error creating request: %v", err)
		}
		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("Api-Key", cfg.APIKey)

		client := &http.Client{Timeout: 120 * time.Second}
		resp, err := client.Do(httpReq)
		if err != nil {
			if attempt < maxRetries {
				time.Sleep(time.Duration(attempt*2) * time.Second)
				continue
			}
			return fmt.Errorf("error calling Ideogram Legacy API: %v", err)
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error reading response: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			if attempt < maxRetries {
				time.Sleep(time.Duration(attempt*2) * time.Second)
				continue
			}
			return fmt.Errorf("Ideogram Legacy API error (%d): %s", resp.StatusCode, string(respBody))
		}

		var apiResp IdeogramResponse
		if err := json.Unmarshal(respBody, &apiResp); err != nil {
			return fmt.Errorf("error unmarshaling response: %v", err)
		}

		if !apiResp.Success {
			if apiResp.Error != nil {
				if attempt < maxRetries {
					time.Sleep(time.Duration(attempt*2) * time.Second)
					continue
				}
				return fmt.Errorf("Ideogram error: %s", apiResp.Error.Message)
			}
		}

		if len(apiResp.Images) > 0 {
			imageURL = apiResp.Images[0].URL
			break
		}

		if attempt < maxRetries {
			time.Sleep(time.Duration(attempt*2) * time.Second)
		}
	}

	if imageURL == "" {
		return fmt.Errorf("no image URL in response")
	}

	// Download image
	if err := downloadImage(ctx, imageURL, outputPath); err != nil {
		return fmt.Errorf("error downloading image: %v", err)
	}

	return nil
}

// downloadImage descarga una imagen desde URL
func downloadImage(ctx context.Context, url, outputPath string) error {
	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download error: %d", resp.StatusCode)
	}

	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Asegurar que el directorio existe
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	return os.WriteFile(outputPath, imageData, 0o644)
}

// Note: saveImage está definida en image-gen-gemini.go (mismo paquete utils)

// GenerateImageWithIdeogramFromScript generates an image from script using configured parameters
func GenerateImageWithIdeogramFromScript(
	ctx context.Context,
	cfg IdeogramV3Config,
	script, topic, outputPath string,
) error {
	// Use default model for prompt generation if not specified
	model := "gemini-2.5-flash"

	// First generate visual prompt from script
	visualPrompt, err := BuildVisualPromptFromScript(ctx, cfg.APIKey, model, script, topic, cfg.AspectRatio, "en", "")
	if err != nil {
		return fmt.Errorf("error building visual prompt: %v", err)
	}

	// Then generate image (auto-detects V3 vs Legacy)
	return GenerateImageWithIdeogram(ctx, cfg, visualPrompt, outputPath)
}

// ============================================================================
// HELPER FUNCTIONS - Configuration from environment and config.json
// ============================================================================

// NewIdeogramV3ConfigFromEnv creates IdeogramV3Config from environment variables and config
// Precedence: env vars > config.json > defaults
func NewIdeogramV3ConfigFromEnv(cfg config.ImageGenToolsConfig) IdeogramV3Config {
	return IdeogramV3Config{
		APIKey:         getEnvOrDefault("IDEOGRAM_API_KEY", cfg.IdeogramAPIKey),
		APIURL:         getEnvOrDefault("IDEOGRAM_ENDPOINT_V3", cfg.IdeogramAPIURL),
		AspectRatio:    getEnvOrDefault("IDEOGRAM_ASPECT_RATIO", cfg.IdeogramAspectRatio),
		RenderingSpeed: getEnvOrDefault("IDEOGRAM_RENDERING_SPEED", cfg.IdeogramRenderingSpeed),
		StyleType:      getEnvOrDefault("IDEOGRAM_STYLE_TYPE", cfg.IdeogramStyleType),
		NumImages:      getEnvIntOrDefault("IDEOGRAM_NUM_IMAGES", cfg.IdeogramNumImages),
		NegativePrompt: getEnvOrDefault("IDEOGRAM_NEGATIVE", cfg.IdeogramNegativePrompt),
	}
}

// getEnvOrDefault returns env var if set, otherwise fallback
func getEnvOrDefault(envKey, fallback string) string {
	if val := os.Getenv(envKey); val != "" {
		return val
	}
	return fallback
}

// getEnvIntOrDefault returns env var as int if set, otherwise fallback
func getEnvIntOrDefault(envKey string, fallback int) int {
	if val := os.Getenv(envKey); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return fallback
}
