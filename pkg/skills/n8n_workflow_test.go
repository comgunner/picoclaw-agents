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

func TestN8NWorkflowSkill_Name(t *testing.T) {
	skill := NewN8NWorkflowSkill("/tmp/test-workspace")

	name := skill.Name()
	if name != "n8n_workflow" {
		t.Errorf("expected name 'n8n_workflow', got '%s'", name)
	}
}

func TestN8NWorkflowSkill_Description(t *testing.T) {
	skill := NewN8NWorkflowSkill("/tmp/test-workspace")

	desc := skill.Description()
	if desc == "" {
		t.Error("expected non-empty description")
	}

	// Verify description contains key concepts
	if !strings.Contains(desc, "n8n") {
		t.Error("description should mention 'n8n'")
	}
	if !strings.Contains(desc, "Workflow") {
		t.Error("description should mention 'Workflow'")
	}
	if !strings.Contains(desc, "JSON") {
		t.Error("description should mention 'JSON'")
	}
}

func TestN8NWorkflowSkill_GetInstructions(t *testing.T) {
	skill := NewN8NWorkflowSkill("/tmp/test-workspace")

	instructions := skill.GetInstructions()
	if instructions == "" {
		t.Error("expected non-empty instructions")
	}

	// Verify instructions contain critical sections
	requiredSections := []string{
		"ROLE & OBJECTIVE",
		"WORKFLOW JSON STRUCTURE",
		"NODE STRUCTURE",
		"CONNECTION STRUCTURE",
		"NODE LIBRARY",
		"Webhook",
		"HTTP Request",
		"Function",
		"EXPRESSION SYNTAX",
		"IMPORT/EXPORT METHODS",
		"BEST PRACTICES",
		"SECURITY CONSIDERATIONS",
	}

	for _, section := range requiredSections {
		if !strings.Contains(instructions, section) {
			t.Errorf("instructions missing section: %s", section)
		}
	}
}

func TestN8NWorkflowSkill_GetAntiPatterns(t *testing.T) {
	skill := NewN8NWorkflowSkill("/tmp/test-workspace")

	antiPatterns := skill.GetAntiPatterns()
	if antiPatterns == "" {
		t.Error("expected non-empty anti-patterns")
	}

	// Verify anti-patterns contain expected content
	requiredPatterns := []string{
		"WORKFLOW DESIGN ANTI-PATTERNS",
		"SECURITY ANTI-PATTERNS",
		"PERFORMANCE ANTI-PATTERNS",
		"CONNECTION ANTI-PATTERNS",
		"No Trigger Node",
		"Hardcoded Credentials",
		"Exposed API Keys",
		"Missing Connections",
	}

	for _, pattern := range requiredPatterns {
		if !strings.Contains(antiPatterns, pattern) {
			t.Errorf("anti-patterns missing: %s", pattern)
		}
	}
}

func TestN8NWorkflowSkill_GetExamples(t *testing.T) {
	skill := NewN8NWorkflowSkill("/tmp/test-workspace")

	examples := skill.GetExamples()
	if examples == "" {
		t.Error("expected non-empty examples")
	}

	// Verify examples contain concrete scenarios
	requiredExamples := []string{
		"EXAMPLE 1: BASIC WORKFLOW",
		"EXAMPLE 2: INTERMEDIATE WORKFLOW",
		"EXAMPLE 3: ADVANCED WORKFLOW",
		"EXAMPLE 4: ERROR HANDLING",
		"EXAMPLE 5: BATCH PROCESSING",
		"Google Sheets",
		"Gmail",
		"PostgreSQL",
		"Slack",
		"HubSpot",
	}

	for _, example := range requiredExamples {
		if !strings.Contains(examples, example) {
			t.Errorf("examples missing: %s", example)
		}
	}
}

func TestN8NWorkflowSkill_BuildSkillContext(t *testing.T) {
	skill := NewN8NWorkflowSkill("/tmp/test-workspace")

	context := skill.BuildSkillContext()
	if context == "" {
		t.Error("expected non-empty skill context")
	}

	// Verify context has proper structure
	if !strings.HasPrefix(context, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━") {
		t.Error("context should start with separator")
	}
	if !strings.Contains(context, "🚀 NATIVE SKILL: n8n Workflow Automation Expert") {
		t.Error("context should contain skill title")
	}
	if !strings.Contains(context, "**ROLE:**") {
		t.Error("context should contain role section")
	}
}

func TestN8NWorkflowSkill_BuildSummary(t *testing.T) {
	skill := NewN8NWorkflowSkill("/tmp/test-workspace")

	summary := skill.BuildSummary()
	if summary == "" {
		t.Error("expected non-empty summary")
	}

	// Verify summary is valid XML-like format
	requiredTags := []string{
		`<skill name="n8n_workflow"`,
		"<purpose>",
		"<pattern>",
		"<nodes>",
		"<features>",
	}

	for _, tag := range requiredTags {
		if !strings.Contains(summary, tag) {
			t.Errorf("summary missing tag: %s", tag)
		}
	}
}

func TestN8NWorkflowSkill_WorkspaceIndependence(t *testing.T) {
	// Verify skill works with different workspace paths
	workspaces := []string{
		"/tmp/test-ws",
		"/home/user/.picoclaw",
		"",
	}

	for _, ws := range workspaces {
		skill := NewN8NWorkflowSkill(ws)

		if skill.Name() != "n8n_workflow" {
			t.Errorf("skill name should be independent of workspace")
		}

		if skill.Description() == "" {
			t.Errorf("skill description should be independent of workspace")
		}
	}
}

func TestN8NWorkflowSkill_Concurrency(t *testing.T) {
	skill := NewN8NWorkflowSkill("/tmp/test-workspace")

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

func TestN8NWorkflowSkill_NodeLibraryCoverage(t *testing.T) {
	skill := NewN8NWorkflowSkill("/tmp/test-workspace")
	instructions := skill.GetInstructions()

	// Verify node library covers all major categories
	nodeCategories := []string{
		"Trigger Nodes",
		"Action Nodes",
		"App Nodes",
		"Logic Nodes",
		"Webhook",
		"Schedule",
		"Manual",
		"HTTP Request",
		"Function",
		"Gmail",
		"Google Sheets",
		"Slack",
		"Discord",
		"Telegram",
		"PostgreSQL",
		"MySQL",
		"HubSpot",
		"IF",
		"Switch",
		"Wait",
		"Error Trigger",
		"Loop",
	}

	for _, category := range nodeCategories {
		if !strings.Contains(instructions, category) {
			t.Errorf("node library missing: %s", category)
		}
	}
}

func TestN8NWorkflowSkill_JSONStructureValidation(t *testing.T) {
	skill := NewN8NWorkflowSkill("/tmp/test-workspace")
	instructions := skill.GetInstructions()

	// Verify JSON structure documentation is present
	jsonStructureElements := []string{
		"Root Object",
		"Required Fields",
		"Optional Fields",
		"NODE STRUCTURE",
		"Required Node Fields",
		"Optional Node Fields",
		"CONNECTION STRUCTURE",
		"nodes",
		"connections",
		"parameters",
		"typeVersion",
		"position",
	}

	for _, element := range jsonStructureElements {
		if !strings.Contains(instructions, element) {
			t.Errorf("JSON structure missing: %s", element)
		}
	}
}

func TestN8NWorkflowSkill_ExpressionSyntaxCoverage(t *testing.T) {
	skill := NewN8NWorkflowSkill("/tmp/test-workspace")
	instructions := skill.GetInstructions()

	// Verify expression syntax is documented
	expressionElements := []string{
		"EXPRESSION SYNTAX",
		"Accessing Data",
		"$json",
		"$node",
		"$now",
		"Date Helpers",
		"String Functions",
		"Math Functions",
	}

	for _, element := range expressionElements {
		if !strings.Contains(instructions, element) {
			t.Errorf("expression syntax missing: %s", element)
		}
	}
}

func TestN8NWorkflowSkill_ImportExportMethods(t *testing.T) {
	skill := NewN8NWorkflowSkill("/tmp/test-workspace")
	instructions := skill.GetInstructions()

	// Verify import/export methods are documented
	importExportElements := []string{
		"IMPORT/EXPORT METHODS",
		"Import Methods",
		"Export Methods",
		"Import from File",
		"Import from URL",
		"Copy/Paste",
		"Export to File",
	}

	for _, element := range importExportElements {
		if !strings.Contains(instructions, element) {
			t.Errorf("import/export methods missing: %s", element)
		}
	}
}

func TestN8NWorkflowSkill_SecurityConsiderations(t *testing.T) {
	skill := NewN8NWorkflowSkill("/tmp/test-workspace")
	instructions := skill.GetInstructions()
	antiPatterns := skill.GetAntiPatterns()

	// Verify security topics are covered
	securityTopics := []string{
		"SECURITY CONSIDERATIONS",
		"Credentials",
		"API Keys",
		"Webhook",
		"Sanitization",
		"Exposed API Keys",
		"Hardcoded Credentials",
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

func TestN8NWorkflowSkill_BestPractices(t *testing.T) {
	skill := NewN8NWorkflowSkill("/tmp/test-workspace")
	instructions := skill.GetInstructions()

	// Verify best practices are documented
	bestPracticeTopics := []string{
		"BEST PRACTICES",
		"Workflow Design",
		"Error Handling",
		"Performance",
		"Version Control",
		"Testing",
	}

	for _, topic := range bestPracticeTopics {
		if !strings.Contains(instructions, topic) {
			t.Errorf("best practices missing topic: %s", topic)
		}
	}
}

func TestN8NWorkflowSkill_WorkflowExamples(t *testing.T) {
	skill := NewN8NWorkflowSkill("/tmp/test-workspace")
	examples := skill.GetExamples()

	// Verify workflow examples contain JSON
	if !strings.Contains(examples, "json") {
		t.Error("examples should contain JSON code blocks")
	}

	// Verify examples include import instructions
	if !strings.Contains(examples, "Import") {
		t.Error("examples should include import instructions")
	}

	// Verify examples have expression explanations
	if !strings.Contains(examples, "Expression") {
		t.Error("examples should include expression explanations")
	}
}
