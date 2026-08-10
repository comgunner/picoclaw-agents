// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTemplateDefaultConfigMatchesExample(t *testing.T) {
	// 1. Get TemplateDefaultConfig
	defaultCfg := TemplateDefaultConfig()
	defaultData, err := json.MarshalIndent(defaultCfg, "", "  ")
	require.NoError(t, err)

	// 2. Read config.example.json from the project root
	// We assume the test is running in pkg/config, so we go up twice
	examplePath := filepath.Join("..", "..", "config", "config.example.json")
	exampleData, err := os.ReadFile(examplePath)
	if err != nil {
		t.Skipf("Skipping test: could not read %s", examplePath)
	}

	// 3. Unmarshal example into a Config struct to normalize formatting/fields
	var exampleCfg Config
	err = json.Unmarshal(exampleData, &exampleCfg)
	require.NoError(t, err)

	normalizedExampleData, err := json.MarshalIndent(exampleCfg, "", "  ")
	require.NoError(t, err)

	// 4. Compare
	assert.JSONEq(
		t,
		string(normalizedExampleData),
		string(defaultData),
		"DefaultConfig() in pkg/config/defaults.go must match config/config.example.json. Please update defaults.go to match the example.",
	)
}

func TestDefaultConfig_HasOpenRouterFreeModels(t *testing.T) {
	cfg := DefaultConfig()

	// Check that free models exist in model list
	freeModels := make(map[string]bool)
	for _, m := range cfg.ModelList {
		freeModels[m.ModelName] = true
	}

	assert.True(t, freeModels["openrouter-free-1"], "missing openrouter-free-1 model")
	assert.True(t, freeModels["openrouter-free-2"], "missing openrouter-free-2 model")
	assert.True(t, freeModels["openrouter-free-3"], "missing openrouter-free-3 model")
}

func TestDefaultConfig_OpenRouterFreeModelIDs(t *testing.T) {
	cfg := DefaultConfig()

	for _, m := range cfg.ModelList {
		switch m.ModelName {
		case "openrouter-free-1":
			assert.Equal(t, "nvidia/nemotron-3-ultra-550b-a55b:free", m.Model)
			assert.Equal(t, "https://openrouter.ai/api/v1", m.APIBase)
		case "openrouter-free-2":
			assert.Equal(t, "openai/gpt-oss-20b:free", m.Model)
			assert.Equal(t, "https://openrouter.ai/api/v1", m.APIBase)
		case "openrouter-free-3":
			assert.Equal(t, "meta-llama/llama-3.1-8b-instruct", m.Model)
			assert.Equal(t, "https://openrouter.ai/api/v1", m.APIBase)
		}
	}
}

func TestDefaultConfig_ImageGenDefaults(t *testing.T) {
	cfg := DefaultConfig()

	assert.Equal(t, "openrouter", cfg.Tools.ImageGen.Provider)
	assert.Equal(t, "krea/krea-2-medium-turbo", cfg.Tools.ImageGen.OpenRouterImageModel)
	assert.Equal(t, "openrouter/free", cfg.Tools.ImageGen.OpenRouterTextModel)
	assert.Equal(t, "1:1", cfg.Tools.ImageGen.AspectRatio)
}

func TestDefaultConfig_CascadeDefaults(t *testing.T) {
	cfg := DefaultConfig()

	assert.True(t, cfg.Tools.ImageGen.Cascade.Enabled)
	assert.Len(t, cfg.Tools.ImageGen.Cascade.TextModels, 6)
	assert.Equal(t, "opencode/mimo-v2.5-free", cfg.Tools.ImageGen.Cascade.TextModels[0])
	assert.Equal(t, "opencode/deepseek-v4-flash-free", cfg.Tools.ImageGen.Cascade.TextModels[1])
	assert.Equal(t, "opencode/nemotron-3-ultra-free", cfg.Tools.ImageGen.Cascade.TextModels[2])
	assert.Equal(t, "nvidia/nemotron-3-ultra-550b-a55b:free", cfg.Tools.ImageGen.Cascade.TextModels[3])
	assert.Equal(t, "openai/gpt-oss-20b:free", cfg.Tools.ImageGen.Cascade.TextModels[4])
	assert.Equal(t, "meta-llama/llama-3.1-8b-instruct", cfg.Tools.ImageGen.Cascade.TextModels[5])
	assert.Len(t, cfg.Tools.ImageGen.Cascade.ImageModels, 1)
	assert.Equal(t, "krea/krea-2-medium-turbo", cfg.Tools.ImageGen.Cascade.ImageModels[0])
}

func TestCascadeConfig_JSON(t *testing.T) {
	cfg := CascadeConfig{
		Enabled:     true,
		TextModels:  []string{"model-a", "model-b"},
		ImageModels: []string{"img-a"},
	}

	data, err := json.Marshal(cfg)
	require.NoError(t, err)

	var decoded CascadeConfig
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, cfg.Enabled, decoded.Enabled)
	assert.Equal(t, cfg.TextModels, decoded.TextModels)
	assert.Equal(t, cfg.ImageModels, decoded.ImageModels)
}
