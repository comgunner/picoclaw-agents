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
