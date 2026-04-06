// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package mcp

import (
	"strings"
	"testing"
)

func TestSanitizeName_Normal(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello_world", "hello_world"},
		{"HelloWorld", "helloworld"},
		{"my-tool-name", "my_tool_name"},
		{"Test123", "test123"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := sanitizeName(tt.input)
			if result != tt.expected {
				t.Errorf("sanitizeName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSanitizeName_SpecialChars(t *testing.T) {
	result := sanitizeName("hello@world!test#123")
	if result != "helloworldtest123" {
		t.Errorf("sanitizeName(%q) = %q, want %q", "hello@world!test#123", result, "helloworldtest123")
	}
}

func TestSanitizeName_Empty(t *testing.T) {
	result := sanitizeName("")
	if result != "unnamed" {
		t.Errorf("sanitizeName(%q) = %q, want %q", "", result, "unnamed")
	}
}

func TestSanitizeName_OnlySpecialChars(t *testing.T) {
	result := sanitizeName("@#$%")
	if result != "unnamed" {
		t.Errorf("sanitizeName(%q) = %q, want %q", "@#$%", result, "unnamed")
	}
}

func TestMCPToolName_Short(t *testing.T) {
	result := MCPToolName("github", "list_issues")
	expected := "mcp_github_list_issues"
	if result != expected {
		t.Errorf("MCPToolName(%q, %q) = %q, want %q", "github", "list_issues", result, expected)
	}
}

func TestMCPToolName_LongWithHash(t *testing.T) {
	// Create a name that exceeds MaxToolNameLen
	serverName := "very_long_server_name_for_testing"
	toolName := "very_very_very_very_very_very_very_very_very_long_tool_name"
	result := MCPToolName(serverName, toolName)

	if len(result) > MaxToolNameLen {
		t.Errorf("MCPToolName result length %d exceeds max %d", len(result), MaxToolNameLen)
	}
	if !strings.HasPrefix(result, "mcp_") {
		t.Errorf("MCPToolName result %q doesn't start with mcp_", result)
	}
	// Should have hash suffix
	if !strings.Contains(result, "_") {
		t.Errorf("MCPToolName result %q should contain hash suffix", result)
	}
}

func TestMCPToolName_Deterministic(t *testing.T) {
	serverName := "github"
	toolName := "list_issues"

	result1 := MCPToolName(serverName, toolName)
	result2 := MCPToolName(serverName, toolName)

	if result1 != result2 {
		t.Errorf("MCPToolName is not deterministic: %q != %q", result1, result2)
	}
}
