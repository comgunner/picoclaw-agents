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

func TestFullStackDeveloperSkill_Name(t *testing.T) {
	skill := NewFullStackDeveloperSkill("/tmp/test-workspace")

	name := skill.Name()
	if name != "fullstack_developer" {
		t.Errorf("expected name 'fullstack_developer', got '%s'", name)
	}
}

func TestFullStackDeveloperSkill_Description(t *testing.T) {
	skill := NewFullStackDeveloperSkill("/tmp/test-workspace")

	desc := skill.Description()
	if desc == "" {
		t.Error("expected non-empty description")
	}

	// Verify description contains key concepts
	if !strings.Contains(desc, "full-stack") {
		t.Error("description should mention 'full-stack'")
	}
	if !strings.Contains(desc, "development") {
		t.Error("description should mention 'development'")
	}
}

func TestFullStackDeveloperSkill_GetInstructions(t *testing.T) {
	skill := NewFullStackDeveloperSkill("/tmp/test-workspace")

	instructions := skill.GetInstructions()
	if instructions == "" {
		t.Error("expected non-empty instructions")
	}

	// Verify instructions contain critical sections
	requiredSections := []string{
		"DEVELOPMENT WORKFLOW",
		"FRONTEND PATTERNS",
		"BACKEND PATTERNS",
		"DATABASE PATTERNS",
		"TESTING PATTERNS",
		"SECURITY CHECKLIST",
		"GIT WORKFLOW",
		"CODE QUALITY",
	}

	for _, section := range requiredSections {
		if !strings.Contains(instructions, section) {
			t.Errorf("instructions missing section: %s", section)
		}
	}
}

func TestFullStackDeveloperSkill_GetAntiPatterns(t *testing.T) {
	skill := NewFullStackDeveloperSkill("/tmp/test-workspace")

	antiPatterns := skill.GetAntiPatterns()
	if antiPatterns == "" {
		t.Error("expected non-empty anti-patterns")
	}

	// Verify anti-patterns contain expected content
	requiredPatterns := []string{
		"CODE SMELLS",
		"SECURITY ANTI-PATTERNS",
		"TESTING ANTI-PATTERNS",
		"GIT ANTI-PATTERNS",
		"API ANTI-PATTERNS",
		"SQL Injection",
		"XSS",
		"Long Functions",
	}

	for _, pattern := range requiredPatterns {
		if !strings.Contains(antiPatterns, pattern) {
			t.Errorf("anti-patterns missing: %s", pattern)
		}
	}
}

func TestFullStackDeveloperSkill_GetExamples(t *testing.T) {
	skill := NewFullStackDeveloperSkill("/tmp/test-workspace")

	examples := skill.GetExamples()
	if examples == "" {
		t.Error("expected non-empty examples")
	}

	// Verify examples contain concrete scenarios
	requiredExamples := []string{
		"EXAMPLE 1: CREATE REACT COMPONENT",
		"EXAMPLE 2: CREATE REST API ENDPOINT",
		"EXAMPLE 3: DATABASE MIGRATION",
		"EXAMPLE 4: DOCKER CONFIGURATION",
		"EXAMPLE 5: GITHUB ACTIONS CI/CD",
	}

	for _, example := range requiredExamples {
		if !strings.Contains(examples, example) {
			t.Errorf("examples missing: %s", example)
		}
	}
}

func TestFullStackDeveloperSkill_BuildSkillContext(t *testing.T) {
	skill := NewFullStackDeveloperSkill("/tmp/test-workspace")

	context := skill.BuildSkillContext()
	if context == "" {
		t.Error("expected non-empty skill context")
	}

	// Verify context has proper structure
	if !strings.HasPrefix(context, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━") {
		t.Error("context should start with separator")
	}
	if !strings.Contains(context, "🚀 NATIVE SKILL: Full-Stack Developer") {
		t.Error("context should contain skill title")
	}
	if !strings.Contains(context, "**PURPOSE:**") {
		t.Error("context should contain purpose section")
	}
}

func TestFullStackDeveloperSkill_BuildSummary(t *testing.T) {
	skill := NewFullStackDeveloperSkill("/tmp/test-workspace")

	summary := skill.BuildSummary()
	if summary == "" {
		t.Error("expected non-empty summary")
	}

	// Verify summary is valid XML-like format
	requiredTags := []string{
		`<skill name="fullstack_developer"`,
		"<purpose>",
		"<pattern>",
		"<stacks>",
		"<practices>",
	}

	for _, tag := range requiredTags {
		if !strings.Contains(summary, tag) {
			t.Errorf("summary missing tag: %s", tag)
		}
	}
}

func TestFullStackDeveloperSkill_WorkspaceIndependence(t *testing.T) {
	// Verify skill works with different workspace paths
	workspaces := []string{
		"/tmp/test-ws",
		"/home/user/.picoclaw",
		"",
	}

	for _, ws := range workspaces {
		skill := NewFullStackDeveloperSkill(ws)

		if skill.Name() != "fullstack_developer" {
			t.Errorf("skill name should be independent of workspace")
		}

		if skill.Description() == "" {
			t.Errorf("skill description should be independent of workspace")
		}
	}
}

func TestFullStackDeveloperSkill_Concurrency(t *testing.T) {
	skill := NewFullStackDeveloperSkill("/tmp/test-workspace")

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

func TestFullStackDeveloperSkill_CodeExamplesPresent(t *testing.T) {
	skill := NewFullStackDeveloperSkill("/tmp/test-workspace")
	examples := skill.GetExamples()

	// Verify code examples are present for different stacks
	if !strings.Contains(examples, "jsx") && !strings.Contains(examples, "function Component") {
		t.Error("examples should include React/JSX code")
	}
	if !strings.Contains(examples, "express") && !strings.Contains(examples, "router.post") {
		t.Error("examples should include Express/Node.js code")
	}
	if !strings.Contains(examples, "CREATE TABLE") && !strings.Contains(examples, "SQL") {
		t.Error("examples should include SQL code")
	}
	if !strings.Contains(examples, "Dockerfile") && !strings.Contains(examples, "FROM node") {
		t.Error("examples should include Docker configuration")
	}
	if !strings.Contains(examples, "github") && !strings.Contains(examples, "actions") {
		t.Error("examples should include GitHub Actions CI/CD")
	}
}

func TestFullStackDeveloperSkill_SecurityCoverage(t *testing.T) {
	skill := NewFullStackDeveloperSkill("/tmp/test-workspace")
	instructions := skill.GetInstructions()
	antiPatterns := skill.GetAntiPatterns()

	// Verify security topics are covered
	securityTopics := []string{
		"SQL Injection",
		"XSS",
		"Authentication",
		"Authorization",
		"Rate Limiting",
		"HTTPS",
		"Input validation",
		"Password",
	}

	for _, topic := range securityTopics {
		topicLower := strings.ToLower(topic)
		foundInInstructions := strings.Contains(strings.ToLower(instructions), topicLower)
		foundInAntiPatterns := strings.Contains(strings.ToLower(antiPatterns), topicLower)

		if !foundInInstructions && !foundInAntiPatterns {
			t.Errorf("security topic not covered: %s", topic)
		}
	}
}

func TestFullStackDeveloperSkill_TestingCoverage(t *testing.T) {
	skill := NewFullStackDeveloperSkill("/tmp/test-workspace")
	instructions := skill.GetInstructions()
	antiPatterns := skill.GetAntiPatterns()

	// Verify testing topics are covered
	testingTopics := []string{
		"Unit Test",
		"Integration",
		"E2E",
		"TDD",
		"AAA",
		"jest",
		"mock",
	}

	for _, topic := range testingTopics {
		topicLower := strings.ToLower(topic)
		foundInInstructions := strings.Contains(strings.ToLower(instructions), topicLower)
		foundInAntiPatterns := strings.Contains(strings.ToLower(antiPatterns), topicLower)

		if !foundInInstructions && !foundInAntiPatterns {
			t.Errorf("testing topic not covered: %s", topic)
		}
	}
}

func TestFullStackDeveloperSkill_GitWorkflowCoverage(t *testing.T) {
	skill := NewFullStackDeveloperSkill("/tmp/test-workspace")
	instructions := skill.GetInstructions()
	antiPatterns := skill.GetAntiPatterns()

	// Verify Git topics are covered
	gitTopics := []string{
		"branch",
		"commit",
		"Conventional Commits",
		"Pull Request",
		"code review",
	}

	for _, topic := range gitTopics {
		topicLower := strings.ToLower(topic)
		foundInInstructions := strings.Contains(strings.ToLower(instructions), topicLower)
		foundInAntiPatterns := strings.Contains(strings.ToLower(antiPatterns), topicLower)

		if !foundInInstructions && !foundInAntiPatterns {
			t.Errorf("Git topic not covered: %s", topic)
		}
	}
}

func TestFullStackDeveloperSkill_FrontendPatterns(t *testing.T) {
	skill := NewFullStackDeveloperSkill("/tmp/test-workspace")
	instructions := skill.GetInstructions()

	// Verify frontend patterns are covered
	frontendTopics := []string{
		"React",
		"useState",
		"useEffect",
		"Component",
		"State Management",
		"Event Handling",
	}

	for _, topic := range frontendTopics {
		if !strings.Contains(instructions, topic) {
			t.Errorf("frontend topic not covered: %s", topic)
		}
	}
}

func TestFullStackDeveloperSkill_BackendPatterns(t *testing.T) {
	skill := NewFullStackDeveloperSkill("/tmp/test-workspace")
	instructions := skill.GetInstructions()

	// Verify backend patterns are covered
	backendTopics := []string{
		"REST API",
		"Error Handling",
		"Middleware",
		"Express",
	}

	for _, topic := range backendTopics {
		if !strings.Contains(instructions, topic) {
			t.Errorf("backend topic not covered: %s", topic)
		}
	}
}

func TestFullStackDeveloperSkill_DatabasePatterns(t *testing.T) {
	skill := NewFullStackDeveloperSkill("/tmp/test-workspace")
	instructions := skill.GetInstructions()

	// Verify database patterns are covered
	databaseTopics := []string{
		"Repository",
		"SQL",
		"NoSQL",
		"transactions",
		"migrations",
	}

	for _, topic := range databaseTopics {
		if !strings.Contains(instructions, topic) {
			t.Errorf("database topic not covered: %s", topic)
		}
	}
}

func TestFullStackDeveloperSkill_CodeQuality(t *testing.T) {
	skill := NewFullStackDeveloperSkill("/tmp/test-workspace")
	instructions := skill.GetInstructions()

	// Verify code quality topics are covered
	qualityTopics := []string{
		"Lint",
		"Format",
		"Type",
		"Documentation",
		"ESLint",
		"Prettier",
	}

	for _, topic := range qualityTopics {
		if !strings.Contains(instructions, topic) {
			t.Errorf("code quality topic not covered: %s", topic)
		}
	}
}
