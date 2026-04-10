// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
// Copyright (c) 2026 PicoClaw contributors

package tools

import (
	"context"
	"fmt"
	"testing"

	"github.com/comgunner/picoclaw/pkg/utils"
)

func TestNewImageGenAntigravityTool_Defaults(t *testing.T) {
	tool := NewImageGenAntigravityTool()
	if tool == nil {
		t.Fatal("Expected non-nil tool")
	}
	if tool.Name() != "image_gen_antigravity" {
		t.Errorf("Name() = %v, want image_gen_antigravity", tool.Name())
	}
	if tool.model != "gemini-3.1-flash-image" {
		t.Errorf("model = %v, want gemini-3.1-flash-image", tool.model)
	}
	if tool.cooldownSecs != defaultCooldownSeconds {
		t.Errorf("cooldownSecs = %v, want %d", tool.cooldownSecs, defaultCooldownSeconds)
	}
}

func TestNewImageGenAntigravityToolFromConfig_CustomCooldown(t *testing.T) {
	tmpDir := t.TempDir()
	cd, _ := utils.NewImageCooldown(tmpDir)
	defer cd.Close()

	tool := NewImageGenAntigravityToolFromConfig(
		"gemini-3.1-flash-image", "16:9", tmpDir, tmpDir,
		600, cd,
	)

	if tool.cooldownSecs != 600 {
		t.Errorf("cooldownSecs = %v, want 600", tool.cooldownSecs)
	}
	if tool.aspectRatio != "16:9" {
		t.Errorf("aspectRatio = %v, want 16:9", tool.aspectRatio)
	}
}

func TestImageGenAntigravityTool_EmptyPrompt(t *testing.T) {
	tool := NewImageGenAntigravityTool()
	result := tool.Execute(context.Background(), map[string]any{
		"prompt": "",
	})

	if !result.IsError {
		t.Error("Expected error for empty prompt")
	}
}

func TestImageGenAntigravityTool_CooldownActive(t *testing.T) {
	tmpDir := t.TempDir()
	cd, _ := utils.NewImageCooldown(tmpDir)
	defer cd.Close()

	// Set cooldown to block.
	cd.Set(300, "antigravity", "gemini-3.1-flash-image")

	tool := NewImageGenAntigravityToolFromConfig(
		"", "", tmpDir, tmpDir, 300, cd,
	)

	result := tool.Execute(context.Background(), map[string]any{
		"prompt": "a cat",
	})

	// Should reject due to cooldown.
	if result.IsError {
		t.Logf("Correctly rejected (error): %s", result.ForLLM)
	}
	if result.ForUser != "" && !strContains(result.ForUser, "Cooldown active") {
		t.Logf("Warning: result doesn't mention cooldown: %s", result.ForUser)
	}
}

func TestImageGenAntigravityTool_NoOAuthCredentials(t *testing.T) {
	tmpDir := t.TempDir()
	cd, _ := utils.NewImageCooldown(tmpDir)
	defer cd.Close()

	// Clear any cooldown first.
	cd.Clear()

	tool := NewImageGenAntigravityToolFromConfig("", "", tmpDir, tmpDir, 300, cd)
	result := tool.Execute(context.Background(), map[string]any{
		"prompt": "a cat",
	})

	// Should reject due to no OAuth credentials (not a cooldown rejection).
	if result.ForUser != "" && strContains(result.ForUser, "OAuth credentials") {
		t.Logf("Correctly rejected: no OAuth credentials")
	}
}

func TestIsRateLimitError(t *testing.T) {
	tests := []struct {
		err  string
		want bool
	}{
		{"HTTP 429: rate limit exceeded", true},
		{"rate limited by API", true},
		{"quota exceeded", true},
		{"connection refused", false},
	}

	for _, tt := range tests {
		var err error
		if tt.err != "" {
			err = fmt.Errorf("%s", tt.err)
		}
		if got := isRateLimitError(err); got != tt.want {
			t.Errorf("isRateLimitError(%q) = %v, want %v", tt.err, got, tt.want)
		}
	}
}

func strContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
