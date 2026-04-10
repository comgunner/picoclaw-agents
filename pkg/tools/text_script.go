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

	"github.com/comgunner/picoclaw/pkg/auth"
	"github.com/comgunner/picoclaw/pkg/utils"
)

// ============== Text Script Tools ==============

type TextScriptCreateTool struct {
	apiKey         string
	useAntigravity bool
	outputDir      string
	tracker        *utils.ImageGenTracker
	templatePath   string
	model          string
	aspectRatio    string
	provider       string
}

func NewTextScriptCreateTool() *TextScriptCreateTool {
	return NewTextScriptCreateToolFromConfig("", "", "", "", "", "", "")
}

func NewTextScriptCreateToolFromConfig(
	configAPIKey, configOutputDir, configTemplatePath, workspace, configModel, configAspectRatio, configProvider string,
) *TextScriptCreateTool {
	apiKey := strings.TrimSpace(os.Getenv(utils.EnvGeminiAPIKey))
	outputDir := strings.TrimSpace(os.Getenv(utils.EnvImageGenOutputDir))
	templatePath := strings.TrimSpace(os.Getenv(utils.EnvImageScriptPath))
	model := strings.TrimSpace(os.Getenv(utils.EnvGeminiImageModel))
	aspectRatio := strings.TrimSpace(os.Getenv(utils.EnvAspectRatio))
	provider := strings.TrimSpace(os.Getenv(utils.EnvImageGenProvider))

	if apiKey == "" {
		apiKey = strings.TrimSpace(configAPIKey)
	}
	if outputDir == "" {
		outputDir = strings.TrimSpace(configOutputDir)
	}
	outputDir = resolveOutputDir(outputDir, workspace)
	if templatePath == "" {
		templatePath = strings.TrimSpace(configTemplatePath)
	}
	templatePath = resolvePathInWorkspace(templatePath, workspace)

	if model == "" {
		model = strings.TrimSpace(configModel)
	}
	if model == "" {
		model = "gemini-2.5-flash"
	}
	if aspectRatio == "" {
		aspectRatio = strings.TrimSpace(configAspectRatio)
	}
	if provider == "" {
		provider = strings.TrimSpace(configProvider)
	}
	if provider == "" {
		provider = "gemini"
	}

	// Check if Antigravity OAuth is available (PRIMARY method).
	useAntigravity := false
	cred, err := auth.GetCredential("google-antigravity")
	if err == nil && cred != nil && !cred.IsExpired() {
		useAntigravity = true
	}

	// Initialize tracker
	trackerPath := filepath.Join(outputDir, "tracker.json")
	tracker, _ := utils.NewImageGenTracker(trackerPath)

	return &TextScriptCreateTool{
		apiKey:         apiKey,
		useAntigravity: useAntigravity,
		outputDir:      outputDir,
		tracker:        tracker,
		templatePath:   templatePath,
		model:          model,
		aspectRatio:    aspectRatio,
		provider:       provider,
	}
}

func (t *TextScriptCreateTool) Name() string {
	return "text_script_create"
}

func (t *TextScriptCreateTool) Description() string {
	return "Generate text script/screenplay from a topic (news, history, tutorial, etc.). Uses prompt_base.txt as template. Use only when user requests script/post text or Script -> Image workflow."
}

func (t *TextScriptCreateTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"topic": map[string]any{
				"type":        "string",
				"description": "Main topic of the script",
			},
			"category": map[string]any{
				"type":        "string",
				"description": "Category: 'news', 'history', 'tutorial', 'announcement', etc.",
			},
			"duration": map[string]any{
				"type":        "string",
				"description": "Estimated duration: '30s', '60s', '5min'",
			},
			"tone": map[string]any{
				"type":        "string",
				"description": "Tone: 'professional', 'casual', 'engaging'",
			},
			"language": map[string]any{
				"type":        "string",
				"description": "Language: 'en', 'es' (default: auto-detected from topic)",
			},
		},
		"required": []string{"topic"},
	}
}

func (t *TextScriptCreateTool) Execute(ctx context.Context, args map[string]any) *ToolResult {
	topic, _ := args["topic"].(string)
	category, _ := args["category"].(string)
	duration, _ := args["duration"].(string)
	tone, _ := args["tone"].(string)
	language, _ := args["language"].(string)

	// Validate required parameters
	topic = strings.TrimSpace(topic)
	if topic == "" {
		return ErrorResult("topic is required")
	}

	// Ensure output directory
	if t.outputDir == "" {
		t.outputDir = "./workspace/image_gen"
	}
	if err := os.MkdirAll(t.outputDir, 0o755); err != nil {
		return ErrorResult(fmt.Sprintf("error creating directory: %v", err))
	}

	callCtx, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	req := utils.TextScriptRequest{
		Topic:        topic,
		Category:     category,
		Language:     language,
		Duration:     duration,
		Tone:         tone,
		TemplatePath: t.templatePath,
	}

	var result *utils.TextScriptResult
	var err error

	// PRIMARY: Try Antigravity OAuth (no API key needed)
	// Includes automatic retry with backoff on 429 rate limit (5s→15s→30s)
	if t.useAntigravity {
		result, err = GenerateTextScriptAntigravity(callCtx, t.model, req)
		if err != nil {
			return ErrorResult(fmt.Sprintf("text_script_create failed: %v", err)).WithError(err)
		}
		return t.saveAndReturn(result, topic, category)
	}

	// FALLBACK: Use Gemini API key (only if OAuth not available)
	if t.apiKey == "" {
		return UserResult(
			"No text generation method available. " +
				"Configure Antigravity OAuth (`picoclaw auth login --provider google-antigravity`) " +
				"or Gemini API key in config.json (tools.image_gen.gemini_api_key)",
		)
	}

	result, err = utils.GenerateTextScript(callCtx, t.apiKey, t.model, req)
	if err != nil {
		return ErrorResult(fmt.Sprintf("text_script_create failed: %v", err)).WithError(err)
	}

	return t.saveAndReturn(result, topic, category)
}

// saveAndReturn saves the script and returns the result.
func (t *TextScriptCreateTool) saveAndReturn(result *utils.TextScriptResult, topic, category string) *ToolResult {
	id := utils.GenerateID()
	scriptDir := filepath.Join(t.outputDir, id)
	if err := os.MkdirAll(scriptDir, 0o755); err != nil {
		return ErrorResult(fmt.Sprintf("error creating directory: %v", err))
	}

	scriptPath := filepath.Join(scriptDir, fmt.Sprintf("%s.-script.txt", id))
	if err := os.WriteFile(scriptPath, []byte(result.Script), 0o644); err != nil {
		return ErrorResult(fmt.Sprintf("error saving script: %v", err))
	}

	// Register in tracker
	if t.tracker != nil {
		record := utils.ImageGenRecord{
			ID:          id,
			DateTime:    time.Now().Format("2006-01-02 15:04:05"),
			Prompt:      topic,
			Provider:    t.provider,
			ScriptPath:  scriptPath,
			AspectRatio: t.aspectRatio,
			Model:       t.model,
			Language:    result.Language,
			Metadata: map[string]string{
				"word_count":         fmt.Sprintf("%d", result.WordCount),
				"estimated_duration": result.EstimatedDuration,
				"category":           category,
			},
		}
		t.tracker.Add(record)
	}

	return UserResult(
		fmt.Sprintf("Script generated successfully.\nPath: %s\nWords: %d\nEstimated Duration: %s\nLanguage: %s",
			scriptPath, result.WordCount, result.EstimatedDuration, result.Language),
	)
}
