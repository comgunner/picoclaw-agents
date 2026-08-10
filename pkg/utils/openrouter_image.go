// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// OpenRouterImageRequest represents a request for image generation via OpenRouter.
type OpenRouterImageRequest struct {
	Prompt      string
	Model       string // e.g. "krea/krea-2-medium-turbo"
	AspectRatio string // e.g. "1:1", "16:9"
	APIKey      string
	OutputDir   string // directory to save the image
}

// OpenRouterImageResult contains the result of a successful image generation.
type OpenRouterImageResult struct {
	Path       string // saved file path
	Model      string // model that succeeded
	DurationMs int64  // time taken in milliseconds
}

// GenerateImageWithOpenRouter generates an image using OpenRouter's image generation API.
// Uses cascade fallback if multiple models are configured.
func GenerateImageWithOpenRouter(ctx context.Context, req OpenRouterImageRequest) (*OpenRouterImageResult, error) {
	if req.APIKey == "" {
		return nil, fmt.Errorf("OpenRouter API key not configured. Run: picoclaw auth login --provider openrouter")
	}

	if req.Model == "" {
		req.Model = "krea/krea-2-medium-turbo"
	}
	if req.AspectRatio == "" {
		req.AspectRatio = "1:1"
	}

	start := time.Now()

	// Use cascade runner for image generation
	runner := NewCascadeRunner(req.APIKey)
	cascade := []string{req.Model}

	log.Printf("[openrouter-image] Generating with model: %s", req.Model)

	imgData, usedModel, err := runner.ImageCascade(ctx, req.Prompt, cascade, req.AspectRatio)
	if err != nil {
		return nil, fmt.Errorf("image generation failed: %w", err)
	}

	// Save image
	outputPath := req.OutputDir
	if outputPath == "" {
		outputPath = "./workspace/image_gen"
	}

	if err := os.MkdirAll(outputPath, 0o755); err != nil {
		return nil, fmt.Errorf("create output dir: %w", err)
	}

	filename := fmt.Sprintf("openrouter_%d.png", time.Now().UnixMilli())
	fullPath := filepath.Join(outputPath, filename)

	if err := os.WriteFile(fullPath, imgData, 0o644); err != nil {
		return nil, fmt.Errorf("save image: %w", err)
	}

	duration := time.Since(start).Milliseconds()
	log.Printf("[openrouter-image] Saved %s (%d bytes, %dms, model: %s)", fullPath, len(imgData), duration, usedModel)

	return &OpenRouterImageResult{
		Path:       fullPath,
		Model:      usedModel,
		DurationMs: duration,
	}, nil
}

// GenerateImageWithOpenRouterCascade generates an image with cascade fallback across multiple models.
func GenerateImageWithOpenRouterCascade(
	ctx context.Context,
	apiKey string,
	prompt string,
	models []string,
	aspectRatio string,
	outputDir string,
) (*OpenRouterImageResult, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("OpenRouter API key not configured. Run: picoclaw auth login --provider openrouter")
	}

	if len(models) == 0 {
		models = []string{"krea/krea-2-medium-turbo"}
	}
	if aspectRatio == "" {
		aspectRatio = "1:1"
	}

	start := time.Now()

	runner := NewCascadeRunner(apiKey)

	log.Printf("[openrouter-image] Cascade: trying %d models", len(models))

	imgData, usedModel, err := runner.ImageCascade(ctx, prompt, models, aspectRatio)
	if err != nil {
		return nil, fmt.Errorf("all cascade image models failed: %w", err)
	}

	// Save image
	if outputDir == "" {
		outputDir = "./workspace/image_gen"
	}

	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return nil, fmt.Errorf("create output dir: %w", err)
	}

	filename := fmt.Sprintf("openrouter_%d.png", time.Now().UnixMilli())
	fullPath := filepath.Join(outputDir, filename)

	if err := os.WriteFile(fullPath, imgData, 0o644); err != nil {
		return nil, fmt.Errorf("save image: %w", err)
	}

	duration := time.Since(start).Milliseconds()
	log.Printf("[openrouter-image] Saved %s (%d bytes, %dms, model: %s)", fullPath, len(imgData), duration, usedModel)

	return &OpenRouterImageResult{
		Path:       fullPath,
		Model:      usedModel,
		DurationMs: duration,
	}, nil
}
