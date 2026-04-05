// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)

package skills

import (
	"strings"
	"testing"
)

func TestBinanceMCPSkill_Name(t *testing.T) {
	skill := NewBinanceMCPSkill("/tmp/test-workspace")

	name := skill.Name()
	if name != "binance_mcp" {
		t.Errorf("expected name 'binance_mcp', got '%s'", name)
	}
}

func TestBinanceMCPSkill_Description(t *testing.T) {
	skill := NewBinanceMCPSkill("/tmp/test-workspace")

	desc := skill.Description()
	if desc == "" {
		t.Error("expected non-empty description")
	}

	// Verify description contains key concepts
	if !strings.Contains(desc, "MCP") {
		t.Error("description should mention 'MCP'")
	}
	if !strings.Contains(desc, "Binance") {
		t.Error("description should mention 'Binance'")
	}
	if !strings.Contains(desc, "trading") {
		t.Error("description should mention 'trading'")
	}
}

func TestBinanceMCPSkill_GetInstructions(t *testing.T) {
	skill := NewBinanceMCPSkill("/tmp/test-workspace")

	instructions := skill.GetInstructions()
	if instructions == "" {
		t.Error("expected non-empty instructions")
	}

	// Verify instructions contain critical sections
	requiredSections := []string{
		"WHEN TO USE",
		"Public Data",
		"Private Trading",
		"CONNECTION MODES",
		"Public Mode",
		"Private Mode",
		"CONFIGURATION",
		"get_ticker_price",
		"open_futures_position",
		"confirm=true",
	}

	for _, section := range requiredSections {
		if !strings.Contains(instructions, section) {
			t.Errorf("instructions missing section: %s", section)
		}
	}
}

func TestBinanceMCPSkill_GetAntiPatterns(t *testing.T) {
	skill := NewBinanceMCPSkill("/tmp/test-workspace")

	antiPatterns := skill.GetAntiPatterns()
	if antiPatterns == "" {
		t.Error("expected non-empty anti-patterns")
	}

	// Verify anti-patterns contain expected content
	requiredPatterns := []string{
		"Anti-Pattern 1: Trading Without Confirmation",
		"Anti-Pattern 2: Assuming API Keys Exist",
		"Anti-Pattern 3: Hardcoding Symbols",
		"Anti-Pattern 4: Ignoring Rate Limits",
		"Anti-Pattern 5: Unsafe Position Sizing",
		"confirm=true",
	}

	for _, pattern := range requiredPatterns {
		if !strings.Contains(antiPatterns, pattern) {
			t.Errorf("anti-patterns missing: %s", pattern)
		}
	}
}

func TestBinanceMCPSkill_GetExamples(t *testing.T) {
	skill := NewBinanceMCPSkill("/tmp/test-workspace")

	examples := skill.GetExamples()
	if examples == "" {
		t.Error("expected non-empty examples")
	}

	// Verify examples contain concrete scenarios
	requiredExamples := []string{
		"Example 1: Public Price Query",
		"Example 2: Market Depth Analysis",
		"Example 3: Check Balances",
		"Example 4: Open Futures Position",
		"Example 5: Close Position",
		"Example 6: Futures Volume Analysis",
	}

	for _, example := range requiredExamples {
		if !strings.Contains(examples, example) {
			t.Errorf("examples missing: %s", example)
		}
	}
}

func TestBinanceMCPSkill_BuildSkillContext(t *testing.T) {
	skill := NewBinanceMCPSkill("/tmp/test-workspace")

	context := skill.BuildSkillContext()
	if context == "" {
		t.Error("expected non-empty skill context")
	}

	// Verify context has proper structure
	if !strings.HasPrefix(context, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━") {
		t.Error("context should start with separator")
	}
	if !strings.Contains(context, "🚀 NATIVE SKILL: Binance MCP Connection") {
		t.Error("context should contain skill title")
	}
	if !strings.Contains(context, "**PURPOSE:**") {
		t.Error("context should contain purpose section")
	}
}

func TestBinanceMCPSkill_BuildSummary(t *testing.T) {
	skill := NewBinanceMCPSkill("/tmp/test-workspace")

	summary := skill.BuildSummary()
	if summary == "" {
		t.Error("expected non-empty summary")
	}

	// Verify summary is valid XML-like format
	requiredTags := []string{
		`<skill name="binance_mcp"`,
		"<purpose>",
		"<pattern>",
		"<tools>",
		"<modes>",
	}

	for _, tag := range requiredTags {
		if !strings.Contains(summary, tag) {
			t.Errorf("summary missing tag: %s", tag)
		}
	}
}

func TestBinanceMCPSkill_WorkspaceIndependence(t *testing.T) {
	// Verify skill works with different workspace paths
	workspaces := []string{
		"/tmp/test-ws",
		"/home/user/.picoclaw",
		"",
	}

	for _, ws := range workspaces {
		skill := NewBinanceMCPSkill(ws)

		if skill.Name() != "binance_mcp" {
			t.Errorf("skill name should be independent of workspace")
		}

		if skill.Description() == "" {
			t.Errorf("skill description should be independent of workspace")
		}
	}
}

func TestBinanceMCPSkill_Concurrency(t *testing.T) {
	skill := NewBinanceMCPSkill("/tmp/test-workspace")

	// Test concurrent access to skill methods
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func() {
			_ = skill.Name()
			_ = skill.Description()
			_ = skill.GetInstructions()
			_ = skill.GetAntiPatterns()
			_ = skill.GetExamples()
			_ = skill.BuildSkillContext()
			_ = skill.BuildSummary()
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestBinanceMCPSkill_ContentValidation(t *testing.T) {
	skill := NewBinanceMCPSkill("/tmp/test-workspace")

	// Verify instructions mention both connection modes
	instructions := skill.GetInstructions()
	if !strings.Contains(instructions, "Public Mode") {
		t.Error("instructions should describe Public Mode")
	}
	if !strings.Contains(instructions, "Private Mode") {
		t.Error("instructions should describe Private Mode")
	}
	if !strings.Contains(instructions, "BINANCE_API_KEY") {
		t.Error("instructions should mention API key configuration")
	}

	// Verify safety checklist is present
	if !strings.Contains(instructions, "SAFETY CHECKLIST") && !strings.Contains(instructions, "Safety Checklist") {
		// Check if safety items are present even without exact header
		hasSafetyItems := strings.Contains(instructions, "Verify API Mode") ||
			strings.Contains(instructions, "Check Balance") ||
			strings.Contains(instructions, "User Approval")
		if !hasSafetyItems {
			t.Error("instructions should include safety checklist or safety items")
		}
	}

	// Verify quick reference tables exist
	examples := skill.GetExamples()
	if !strings.Contains(examples, "Public Commands") {
		t.Error("examples should include public commands reference")
	}
	if !strings.Contains(examples, "Private Commands") {
		t.Error("examples should include private commands reference")
	}
	if !strings.Contains(examples, "Fast-path Commands") {
		t.Error("examples should include fast-path commands section")
	}
	// Verify it mentions current limitations
	if !strings.Contains(examples, "Currently Available Fast-paths") {
		t.Error("examples should list currently available fast-paths")
	}
}

func TestBinanceMCPSkill_SecurityEmphasis(t *testing.T) {
	skill := NewBinanceMCPSkill("/tmp/test-workspace")

	instructions := skill.GetInstructions()
	antiPatterns := skill.GetAntiPatterns()

	// Verify security is emphasized
	if !strings.Contains(instructions, "confirm=true") {
		t.Error("instructions should emphasize confirm parameter")
	}
	if !strings.Contains(instructions, "API keys never exposed") {
		t.Error("instructions should mention API key security")
	}

	// Verify trading warnings exist
	if !strings.Contains(antiPatterns, "Trading Without Confirmation") {
		t.Error("anti-patterns should warn about trading without confirmation")
	}
	if !strings.Contains(antiPatterns, "Unsafe Position Sizing") {
		t.Error("anti-patterns should warn about unsafe position sizing")
	}
}
