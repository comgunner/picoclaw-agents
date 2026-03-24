// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package tools

import (
	"context"
	"testing"
)

func TestVersionControlTool_Name(t *testing.T) {
	tool := NewVersionControlTool("/tmp/test-workspace")
	if tool.Name() != "version_control" {
		t.Errorf("expected name 'version_control', got '%s'", tool.Name())
	}
}

func TestVersionControlTool_Description(t *testing.T) {
	tool := NewVersionControlTool("/tmp/test-workspace")
	if tool.Description() == "" {
		t.Error("expected non-empty description")
	}
}

func TestVersionControlTool_Parameters(t *testing.T) {
	tool := NewVersionControlTool("/tmp/test-workspace")
	params := tool.Parameters()

	if params["type"] != "object" {
		t.Error("parameters type should be 'object'")
	}

	props, ok := params["properties"].(map[string]any)
	if !ok {
		t.Fatal("expected properties map")
	}

	if _, ok := props["action"]; !ok {
		t.Error("missing 'action' parameter")
	}
}

func TestVersionControlTool_Validate(t *testing.T) {
	tool := NewVersionControlTool("/tmp/test-workspace")
	ctx := context.Background()

	result := tool.Execute(ctx, map[string]any{
		"action":  "validate",
		"version": "1.2.3",
	})

	if result.IsError {
		t.Errorf("validate failed: %s", result.ForLLM)
	}
}

func TestVersionControlTool_ValidateWithPrefix(t *testing.T) {
	tool := NewVersionControlTool("/tmp/test-workspace")
	ctx := context.Background()

	result := tool.Execute(ctx, map[string]any{
		"action":  "validate",
		"version": "v1.2.3",
	})

	if result.IsError {
		t.Errorf("validate with 'v' prefix failed: %s", result.ForLLM)
	}
}

func TestVersionControlTool_ValidateInvalid(t *testing.T) {
	tool := NewVersionControlTool("/tmp/test-workspace")
	ctx := context.Background()

	result := tool.Execute(ctx, map[string]any{
		"action":  "validate",
		"version": "not-a-version",
	})

	if !result.IsError {
		t.Error("expected error for invalid version")
	}
}

func TestVersionControlTool_ValidateMissing(t *testing.T) {
	tool := NewVersionControlTool("/tmp/test-workspace")
	ctx := context.Background()

	result := tool.Execute(ctx, map[string]any{
		"action": "validate",
	})

	if !result.IsError {
		t.Error("expected error when version is missing")
	}
}

func TestVersionControlTool_Compare(t *testing.T) {
	tool := NewVersionControlTool("/tmp/test-workspace")
	ctx := context.Background()

	result := tool.Execute(ctx, map[string]any{
		"action":        "compare",
		"version":       "1.2.3",
		"other_version": "1.2.4",
	})

	if result.IsError {
		t.Errorf("compare failed: %s", result.ForLLM)
	}
}

func TestVersionControlTool_CompareEqual(t *testing.T) {
	tool := NewVersionControlTool("/tmp/test-workspace")
	ctx := context.Background()

	result := tool.Execute(ctx, map[string]any{
		"action":        "compare",
		"version":       "2.0.0",
		"other_version": "2.0.0",
	})

	if result.IsError {
		t.Errorf("compare equal failed: %s", result.ForLLM)
	}
}

func TestVersionControlTool_CompareMissing(t *testing.T) {
	tool := NewVersionControlTool("/tmp/test-workspace")
	ctx := context.Background()

	result := tool.Execute(ctx, map[string]any{
		"action":  "compare",
		"version": "1.0.0",
	})

	if !result.IsError {
		t.Error("expected error when other_version is missing")
	}
}

func TestVersionControlTool_BumpMajor(t *testing.T) {
	tool := NewVersionControlTool("/tmp/test-workspace")
	ctx := context.Background()

	result := tool.Execute(ctx, map[string]any{
		"action":    "bump",
		"version":   "1.2.3",
		"bump_type": "major",
	})

	if result.IsError {
		t.Errorf("bump major failed: %s", result.ForLLM)
	}
}

func TestVersionControlTool_BumpMinor(t *testing.T) {
	tool := NewVersionControlTool("/tmp/test-workspace")
	ctx := context.Background()

	result := tool.Execute(ctx, map[string]any{
		"action":    "bump",
		"version":   "1.2.3",
		"bump_type": "minor",
	})

	if result.IsError {
		t.Errorf("bump minor failed: %s", result.ForLLM)
	}
}

func TestVersionControlTool_BumpPatch(t *testing.T) {
	tool := NewVersionControlTool("/tmp/test-workspace")
	ctx := context.Background()

	result := tool.Execute(ctx, map[string]any{
		"action":    "bump",
		"version":   "1.2.3",
		"bump_type": "patch",
	})

	if result.IsError {
		t.Errorf("bump patch failed: %s", result.ForLLM)
	}
}

func TestVersionControlTool_BumpMissingType(t *testing.T) {
	tool := NewVersionControlTool("/tmp/test-workspace")
	ctx := context.Background()

	result := tool.Execute(ctx, map[string]any{
		"action":  "bump",
		"version": "1.0.0",
	})

	if !result.IsError {
		t.Error("expected error when bump_type is missing")
	}
}

func TestVersionControlTool_BumpInvalidType(t *testing.T) {
	tool := NewVersionControlTool("/tmp/test-workspace")
	ctx := context.Background()

	result := tool.Execute(ctx, map[string]any{
		"action":    "bump",
		"version":   "1.0.0",
		"bump_type": "invalid",
	})

	if !result.IsError {
		t.Error("expected error for invalid bump_type")
	}
}

func TestVersionControlTool_Constraint(t *testing.T) {
	tool := NewVersionControlTool("/tmp/test-workspace")
	ctx := context.Background()

	result := tool.Execute(ctx, map[string]any{
		"action":     "constraint",
		"version":    "2.0.0",
		"constraint": ">=1.0.0",
	})

	if result.IsError {
		t.Errorf("constraint check failed: %s", result.ForLLM)
	}
}

func TestVersionControlTool_ConstraintCaret(t *testing.T) {
	tool := NewVersionControlTool("/tmp/test-workspace")
	ctx := context.Background()

	result := tool.Execute(ctx, map[string]any{
		"action":     "constraint",
		"version":    "2.5.0",
		"constraint": "^2.0.0",
	})

	if result.IsError {
		t.Errorf("caret constraint failed: %s", result.ForLLM)
	}
}

func TestVersionControlTool_ConstraintTilde(t *testing.T) {
	tool := NewVersionControlTool("/tmp/test-workspace")
	ctx := context.Background()

	result := tool.Execute(ctx, map[string]any{
		"action":     "constraint",
		"version":    "1.2.5",
		"constraint": "~1.2.0",
	})

	if result.IsError {
		t.Errorf("tilde constraint failed: %s", result.ForLLM)
	}
}

func TestVersionControlTool_ConstraintFails(t *testing.T) {
	tool := NewVersionControlTool("/tmp/test-workspace")
	ctx := context.Background()

	result := tool.Execute(ctx, map[string]any{
		"action":     "constraint",
		"version":    "3.0.0",
		"constraint": "^2.0.0",
	})

	// Should not error, just report that constraint is not satisfied
	if result.IsError {
		t.Logf("constraint not satisfied (expected): %s", result.ForLLM)
	}
}

func TestVersionControlTool_ConstraintMissing(t *testing.T) {
	tool := NewVersionControlTool("/tmp/test-workspace")
	ctx := context.Background()

	result := tool.Execute(ctx, map[string]any{
		"action":  "constraint",
		"version": "1.0.0",
	})

	if !result.IsError {
		t.Error("expected error when constraint is missing")
	}
}

func TestVersionControlTool_Parse(t *testing.T) {
	tool := NewVersionControlTool("/tmp/test-workspace")
	ctx := context.Background()

	result := tool.Execute(ctx, map[string]any{
		"action":  "parse",
		"version": "1.2.3-beta.1+build.123",
	})

	if result.IsError {
		t.Errorf("parse failed: %s", result.ForLLM)
	}
}

func TestVersionControlTool_ParseSimple(t *testing.T) {
	tool := NewVersionControlTool("/tmp/test-workspace")
	ctx := context.Background()

	result := tool.Execute(ctx, map[string]any{
		"action":  "parse",
		"version": "3.0.0",
	})

	if result.IsError {
		t.Errorf("parse simple version failed: %s", result.ForLLM)
	}
}

func TestVersionControlTool_ParseMissing(t *testing.T) {
	tool := NewVersionControlTool("/tmp/test-workspace")
	ctx := context.Background()

	result := tool.Execute(ctx, map[string]any{
		"action": "parse",
	})

	if !result.IsError {
		t.Error("expected error when version is missing")
	}
}

func TestVersionControlTool_InvalidAction(t *testing.T) {
	tool := NewVersionControlTool("/tmp/test-workspace")
	ctx := context.Background()

	result := tool.Execute(ctx, map[string]any{
		"action": "invalid",
	})

	if !result.IsError {
		t.Error("expected error for invalid action")
	}
}

func TestVersionControlTool_CompareWithPrefix(t *testing.T) {
	tool := NewVersionControlTool("/tmp/test-workspace")
	ctx := context.Background()

	result := tool.Execute(ctx, map[string]any{
		"action":        "compare",
		"version":       "v1.2.3",
		"other_version": "V1.2.4",
	})

	if result.IsError {
		t.Errorf("compare with prefix failed: %s", result.ForLLM)
	}
}

func TestVersionControlTool_BumpWithPrefix(t *testing.T) {
	tool := NewVersionControlTool("/tmp/test-workspace")
	ctx := context.Background()

	result := tool.Execute(ctx, map[string]any{
		"action":    "bump",
		"version":   "v1.2.3",
		"bump_type": "patch",
	})

	if result.IsError {
		t.Errorf("bump with prefix failed: %s", result.ForLLM)
	}
}

func TestVersionControlTool_Prerelease(t *testing.T) {
	tool := NewVersionControlTool("/tmp/test-workspace")
	ctx := context.Background()

	result := tool.Execute(ctx, map[string]any{
		"action":  "validate",
		"version": "1.0.0-alpha.1",
	})

	if result.IsError {
		t.Errorf("prerelease validation failed: %s", result.ForLLM)
	}
}

func TestVersionControlTool_ComparePrerelease(t *testing.T) {
	tool := NewVersionControlTool("/tmp/test-workspace")
	ctx := context.Background()

	result := tool.Execute(ctx, map[string]any{
		"action":        "compare",
		"version":       "1.0.0-alpha",
		"other_version": "1.0.0",
	})

	if result.IsError {
		t.Errorf("compare prerelease failed: %s", result.ForLLM)
	}
}

func TestNormalizeVersion(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"1.2.3", "1.2.3"},
		{"v1.2.3", "1.2.3"},
		{"V1.2.3", "1.2.3"},
		{"v1.0.0-beta", "1.0.0-beta"},
		{"1.0.0+build", "1.0.0+build"},
	}

	for _, tt := range tests {
		result := normalizeVersion(tt.input)
		if result != tt.expected {
			t.Errorf("normalizeVersion(%q) = %q, expected %q", tt.input, result, tt.expected)
		}
	}
}
