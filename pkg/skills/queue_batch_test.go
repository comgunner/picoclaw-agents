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

func TestQueueBatchSkill_Name(t *testing.T) {
	skill := NewQueueBatchSkill("/tmp/test-workspace")

	name := skill.Name()
	if name != "queue_batch" {
		t.Errorf("expected name 'queue_batch', got '%s'", name)
	}
}

func TestQueueBatchSkill_Description(t *testing.T) {
	skill := NewQueueBatchSkill("/tmp/test-workspace")

	desc := skill.Description()
	if desc == "" {
		t.Error("expected non-empty description")
	}

	// Verify description contains key concepts
	if !strings.Contains(desc, "queue") {
		t.Error("description should mention 'queue'")
	}
	if !strings.Contains(desc, "fire and forget") {
		t.Error("description should mention 'fire and forget'")
	}
}

func TestQueueBatchSkill_GetInstructions(t *testing.T) {
	skill := NewQueueBatchSkill("/tmp/test-workspace")

	instructions := skill.GetInstructions()
	if instructions == "" {
		t.Error("expected non-empty instructions")
	}

	// Verify instructions contain critical sections
	requiredSections := []string{
		"WHEN TO USE",
		"Signals to Use Queue",
		"USAGE PATTERN",
		"batch_id",
		"queue(action=",
	}

	for _, section := range requiredSections {
		if !strings.Contains(instructions, section) {
			t.Errorf("instructions missing section: %s", section)
		}
	}
}

func TestQueueBatchSkill_GetAntiPatterns(t *testing.T) {
	skill := NewQueueBatchSkill("/tmp/test-workspace")

	antiPatterns := skill.GetAntiPatterns()
	if antiPatterns == "" {
		t.Error("expected non-empty anti-patterns")
	}

	// Verify anti-patterns contain expected content
	requiredPatterns := []string{
		"Anti-Pattern 1: Unnecessary Polling",
		"Anti-Pattern 2: Not Releasing",
		"Anti-Pattern 3: Manual IDs",
	}

	for _, pattern := range requiredPatterns {
		if !strings.Contains(antiPatterns, pattern) {
			t.Errorf("anti-patterns missing: %s", pattern)
		}
	}
}

func TestQueueBatchSkill_GetExamples(t *testing.T) {
	skill := NewQueueBatchSkill("/tmp/test-workspace")

	examples := skill.GetExamples()
	if examples == "" {
		t.Error("expected non-empty examples")
	}

	// Verify examples contain concrete scenarios
	requiredExamples := []string{
		"Example 1: Batch Image Generation",
		"Example 2: Bulk Social Media Upload",
		"Example 3: Model Training",
	}

	for _, example := range requiredExamples {
		if !strings.Contains(examples, example) {
			t.Errorf("examples missing: %s", example)
		}
	}
}

func TestQueueBatchSkill_BuildSkillContext(t *testing.T) {
	skill := NewQueueBatchSkill("/tmp/test-workspace")

	context := skill.BuildSkillContext()
	if context == "" {
		t.Error("expected non-empty skill context")
	}

	// Verify context has proper structure
	if !strings.HasPrefix(context, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━") {
		t.Error("context should start with separator")
	}
	if !strings.Contains(context, "🚀 NATIVE SKILL: Queue/Batch Delegation") {
		t.Error("context should contain skill title")
	}
}

func TestQueueBatchSkill_BuildSummary(t *testing.T) {
	skill := NewQueueBatchSkill("/tmp/test-workspace")

	summary := skill.BuildSummary()
	if summary == "" {
		t.Error("expected non-empty summary")
	}

	// Verify summary is valid XML-like format
	requiredTags := []string{
		"<skill name=\"queue_batch\"",
		"<purpose>",
		"<pattern>",
		"<tools>",
		"<savings>",
	}

	for _, tag := range requiredTags {
		if !strings.Contains(summary, tag) {
			t.Errorf("summary missing tag: %s", tag)
		}
	}
}

func TestQueueBatchSkill_WorkspaceIndependence(t *testing.T) {
	// Verify skill works with different workspace paths
	workspaces := []string{
		"/tmp/test-ws",
		"/home/user/.picoclaw",
		"",
	}

	for _, ws := range workspaces {
		skill := NewQueueBatchSkill(ws)

		if skill.Name() != "queue_batch" {
			t.Errorf("skill name should be independent of workspace")
		}

		if skill.Description() == "" {
			t.Errorf("skill description should be independent of workspace")
		}
	}
}

func TestQueueBatchSkill_Concurrency(t *testing.T) {
	skill := NewQueueBatchSkill("/tmp/test-workspace")

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
