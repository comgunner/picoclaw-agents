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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// CascadeRunner implements free→paid model fallback for OpenRouter.
// Ported from openrouter-studio's text_cascade() pattern.
type CascadeRunner struct {
	apiKey     string
	httpClient *http.Client
}

// NewCascadeRunner creates a new cascade runner with the given API key.
func NewCascadeRunner(apiKey string) *CascadeRunner {
	return &CascadeRunner{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

// TextCascade tries text models in order (free→free1→free2→free3).
// Returns the first successful response content and the model that succeeded.
func (c *CascadeRunner) TextCascade(
	ctx context.Context,
	prompt string,
	cascade []string,
	temperature float64,
	maxTokens int,
) (string, string, error) {
	if len(cascade) == 0 {
		return "", "", fmt.Errorf("empty cascade model list")
	}

	var lastErr error
	for i, model := range cascade {
		content, err := c.tryTextModel(ctx, model, prompt, temperature, maxTokens)
		if err == nil {
			return content, model, nil
		}

		lastErr = err
		errMsg := err.Error()

		if i < len(cascade)-1 {
			if strings.Contains(errMsg, "429") || strings.Contains(errMsg, "rate") {
				log.Printf("[cascade] Rate limited on %s, trying %s...", model, cascade[i+1])
				time.Sleep(5 * time.Second)
				continue
			}
			log.Printf("[cascade] Model %s failed (%v), trying %s...", model, err, cascade[i+1])
			continue
		}
	}

	return "", "", fmt.Errorf("all cascade models failed, last error: %w", lastErr)
}

// tryTextModel makes a single chat completion request to the given model.
func (c *CascadeRunner) tryTextModel(
	ctx context.Context,
	model string,
	prompt string,
	temperature float64,
	maxTokens int,
) (string, error) {
	body := fmt.Sprintf(`{
		"model": %q,
		"messages": [{"role": "user", "content": %q}],
		"temperature": %v,
		"max_tokens": %d
	}`, model, prompt, temperature, maxTokens)

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		"https://openrouter.ai/api/v1/chat/completions",
		strings.NewReader(body),
	)
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Http-Referer", "https://picoclaw.ai")
	req.Header.Set("X-Openrouter-Title", "PicoClaw Agents")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == 429 {
		return "", fmt.Errorf("429 rate limited")
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("API error %d: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("parse response: %w", err)
	}

	if result.Error != nil {
		return "", fmt.Errorf("API error: %s", result.Error.Message)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return result.Choices[0].Message.Content, nil
}

// ImageCascade tries image models in order, returns first success.
// Returns image bytes, the model that succeeded, or an error.
func (c *CascadeRunner) ImageCascade(
	ctx context.Context,
	prompt string,
	cascade []string,
	aspectRatio string,
) ([]byte, string, error) {
	if len(cascade) == 0 {
		return nil, "", fmt.Errorf("empty cascade model list")
	}

	var lastErr error
	for i, model := range cascade {
		imgData, err := c.tryImageModel(ctx, model, prompt, aspectRatio)
		if err == nil {
			return imgData, model, nil
		}

		lastErr = err
		errMsg := err.Error()

		if i < len(cascade)-1 {
			if strings.Contains(errMsg, "429") || strings.Contains(errMsg, "rate") {
				log.Printf("[cascade-image] Rate limited on %s, trying %s...", model, cascade[i+1])
				time.Sleep(30 * time.Second)
				continue
			}
			log.Printf("[cascade-image] Model %s failed (%v), trying %s...", model, err, cascade[i+1])
			continue
		}
	}

	return nil, "", fmt.Errorf("all cascade image models failed, last error: %w", lastErr)
}

// tryImageModel makes a single image generation request.
func (c *CascadeRunner) tryImageModel(
	ctx context.Context,
	model string,
	prompt string,
	aspectRatio string,
) ([]byte, error) {
	width, height := 1024, 1024
	switch aspectRatio {
	case "16:9":
		width, height = 1408, 704
	case "9:16":
		width, height = 704, 1408
	case "4:3":
		width, height = 1248, 936
	case "3:4":
		width, height = 936, 1248
	}

	body := fmt.Sprintf(`{
		"model": %q,
		"prompt": %q,
		"response_format": "b64_json",
		"width": %d,
		"height": %d
	}`, model, prompt, width, height)

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		"https://openrouter.ai/api/v1/images/generations",
		strings.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Http-Referer", "https://picoclaw.ai")
	req.Header.Set("X-Openrouter-Title", "PicoClaw Agents")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == 429 {
		return nil, fmt.Errorf("429 rate limited")
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Data []struct {
			B64JSON string `json:"b64_json"`
			URL     string `json:"url"`
		} `json:"data"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	if result.Error != nil {
		return nil, fmt.Errorf("API error: %s", result.Error.Message)
	}

	if len(result.Data) == 0 {
		return nil, fmt.Errorf("no images in response")
	}

	if result.Data[0].B64JSON != "" {
		return base64.StdEncoding.DecodeString(result.Data[0].B64JSON)
	}

	if result.Data[0].URL != "" {
		return downloadImageFromURL(ctx, c.httpClient, result.Data[0].URL)
	}

	return nil, fmt.Errorf("no image data in response")
}

// downloadImageFromURL fetches an image from a URL.
func downloadImageFromURL(ctx context.Context, client *http.Client, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
