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

func TestAgentTeamWorkflowSkill_Name(t *testing.T) {
	skill := NewAgentTeamWorkflowSkill("/tmp/test-workspace")

	name := skill.Name()
	if name != "agent_team_workflow" {
		t.Errorf("expected name 'agent_team_workflow', got '%s'", name)
	}
}

func TestAgentTeamWorkflowSkill_Description(t *testing.T) {
	skill := NewAgentTeamWorkflowSkill("/tmp/test-workspace")

	desc := skill.Description()
	if desc == "" {
		t.Error("expected non-empty description")
	}

	// Verify description contains key concepts
	if !strings.Contains(desc, "Multi-Agent") {
		t.Error("description should mention 'Multi-Agent'")
	}
	if !strings.Contains(desc, "Orchestrator") {
		t.Error("description should mention 'Orchestrator'")
	}
	if !strings.Contains(desc, "config.json") {
		t.Error("description should mention 'config.json'")
	}
}

func TestAgentTeamWorkflowSkill_GetInstructions(t *testing.T) {
	skill := NewAgentTeamWorkflowSkill("/tmp/test-workspace")

	instructions := skill.GetInstructions()
	if instructions == "" {
		t.Error("expected non-empty instructions")
	}

	// Verify instructions contain critical sections
	requiredSections := []string{
		"ROLE & OBJECTIVE",
		"READING CONFIG.JSON",
		"RESOURCE RULES",
		"SPAWN OPTIMIZATION",
		"TEAM PATTERNS",
		"LOAD BALANCING",
		"INTERNAL METHODS",
		"SECURITY & CONSTRAINTS",
		"allow_agents",
		"max_spawn_depth",
		"max_concurrent",
	}

	for _, section := range requiredSections {
		if !strings.Contains(instructions, section) {
			t.Errorf("instructions missing section: %s", section)
		}
	}
}

func TestAgentTeamWorkflowSkill_GetAntiPatterns(t *testing.T) {
	skill := NewAgentTeamWorkflowSkill("/tmp/test-workspace")

	antiPatterns := skill.GetAntiPatterns()
	if antiPatterns == "" {
		t.Error("expected non-empty anti-patterns")
	}

	// Verify anti-patterns contain expected content
	requiredPatterns := []string{
		"TEAM ORGANIZATION ANTI-PATTERNS",
		"RESOURCE ANTI-PATTERNS",
		"COMMUNICATION ANTI-PATTERNS",
		"MODEL SELECTION ANTI-PATTERNS",
		"SECURITY ANTI-PATTERNS",
		"Spawning Without Config Check",
		"No Concurrency Control",
		"Cross-Session Contamination",
	}

	for _, pattern := range requiredPatterns {
		if !strings.Contains(antiPatterns, pattern) {
			t.Errorf("anti-patterns missing: %s", pattern)
		}
	}
}

func TestAgentTeamWorkflowSkill_GetExamples(t *testing.T) {
	skill := NewAgentTeamWorkflowSkill("/tmp/test-workspace")

	examples := skill.GetExamples()
	if examples == "" {
		t.Error("expected non-empty examples")
	}

	// Verify examples contain concrete scenarios
	requiredExamples := []string{
		"EXAMPLE 1: DEVELOPMENT TEAM",
		"EXAMPLE 2: CONTENT CREATION TEAM",
		"EXAMPLE 3: IMAGE GENERATION BATCH",
		"EXAMPLE 4: RESEARCH TEAM",
		"EXAMPLE 5: CONFIG-BASED AGENT SELECTION",
		"EXAMPLE 6: RESOURCE-AWARE SPAWNING",
		"project_manager",
		"senior_dev",
		"general_worker",
	}

	for _, example := range requiredExamples {
		if !strings.Contains(examples, example) {
			t.Errorf("examples missing: %s", example)
		}
	}
}

func TestAgentTeamWorkflowSkill_BuildSkillContext(t *testing.T) {
	skill := NewAgentTeamWorkflowSkill("/tmp/test-workspace")

	context := skill.BuildSkillContext()
	if context == "" {
		t.Error("expected non-empty skill context")
	}

	// Verify context has proper structure
	if !strings.HasPrefix(context, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━") {
		t.Error("context should start with separator")
	}
	if !strings.Contains(context, "🚀 NATIVE SKILL: Agent Team Workflow Orchestrator") {
		t.Error("context should contain skill title")
	}
	if !strings.Contains(context, "**ROLE:**") {
		t.Error("context should contain role section")
	}
}

func TestAgentTeamWorkflowSkill_BuildSummary(t *testing.T) {
	skill := NewAgentTeamWorkflowSkill("/tmp/test-workspace")

	summary := skill.BuildSummary()
	if summary == "" {
		t.Error("expected non-empty summary")
	}

	// Verify summary is valid XML-like format
	requiredTags := []string{
		`<skill name="agent_team_workflow"`,
		"<purpose>",
		"<pattern>",
		"<config>",
		"<resources>",
		"<patterns>",
	}

	for _, tag := range requiredTags {
		if !strings.Contains(summary, tag) {
			t.Errorf("summary missing tag: %s", tag)
		}
	}
}

func TestAgentTeamWorkflowSkill_WorkspaceIndependence(t *testing.T) {
	// Verify skill works with different workspace paths
	workspaces := []string{
		"/tmp/test-ws",
		"/home/user/.picoclaw",
		"",
	}

	for _, ws := range workspaces {
		skill := NewAgentTeamWorkflowSkill(ws)

		if skill.Name() != "agent_team_workflow" {
			t.Errorf("skill name should be independent of workspace")
		}

		if skill.Description() == "" {
			t.Errorf("skill description should be independent of workspace")
		}
	}
}

func TestAgentTeamWorkflowSkill_Concurrency(t *testing.T) {
	skill := NewAgentTeamWorkflowSkill("/tmp/test-workspace")

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

func TestAgentTeamWorkflowSkill_TeamPatternsCoverage(t *testing.T) {
	skill := NewAgentTeamWorkflowSkill("/tmp/test-workspace")
	instructions := skill.GetInstructions()

	// Verify all team patterns are documented
	teamPatterns := []string{
		"Development Team",
		"Content Team",
		"Image Team",
		"Social Media Team",
		"Research Team",
		"General Team",
	}

	for _, pattern := range teamPatterns {
		if !strings.Contains(instructions, pattern) {
			t.Errorf("team pattern missing: %s", pattern)
		}
	}
}

func TestAgentTeamWorkflowSkill_ResourceRules(t *testing.T) {
	skill := NewAgentTeamWorkflowSkill("/tmp/test-workspace")
	instructions := skill.GetInstructions()

	// Verify resource rules are documented
	resourceRules := []string{
		"CPU Limits",
		"RAM Limits",
		"Concurrency Limits",
		"semaphore",
	}

	for _, rule := range resourceRules {
		if !strings.Contains(instructions, rule) {
			t.Errorf("resource rule missing: %s", rule)
		}
	}
}

func TestAgentTeamWorkflowSkill_SpawnOptimization(t *testing.T) {
	skill := NewAgentTeamWorkflowSkill("/tmp/test-workspace")
	instructions := skill.GetInstructions()

	// Verify spawn optimization concepts are documented
	spawnConcepts := []string{
		"allow_agents",
		"max_spawn_depth",
		"max_children_per_agent",
		"max_concurrent",
		"Wildcard",
	}

	for _, concept := range spawnConcepts {
		if !strings.Contains(instructions, concept) {
			t.Errorf("spawn optimization concept missing: %s", concept)
		}
	}
}

func TestAgentTeamWorkflowSkill_InternalMethods(t *testing.T) {
	skill := NewAgentTeamWorkflowSkill("/tmp/test-workspace")
	instructions := skill.GetInstructions()

	// Verify internal methods are documented
	internalMethods := []string{
		"GetAgent",
		"CanSpawnSubagent",
		"ListAgentIDs",
		"GetDefaultAgent",
		"registry",
	}

	for _, method := range internalMethods {
		if !strings.Contains(instructions, method) {
			t.Errorf("internal method missing: %s", method)
		}
	}
}

func TestAgentTeamWorkflowSkill_SecurityConstraints(t *testing.T) {
	skill := NewAgentTeamWorkflowSkill("/tmp/test-workspace")
	instructions := skill.GetInstructions()

	// Verify security constraints are documented
	securityConstraints := []string{
		"restrict_to_workspace",
		"max_tool_iterations",
		"Token Budget",
		"Session Isolation",
		"dm_scope",
	}

	for _, constraint := range securityConstraints {
		if !strings.Contains(instructions, constraint) {
			t.Errorf("security constraint missing: %s", constraint)
		}
	}
}

func TestAgentTeamWorkflowSkill_LoadBalancing(t *testing.T) {
	skill := NewAgentTeamWorkflowSkill("/tmp/test-workspace")
	instructions := skill.GetInstructions()

	// Verify load balancing concepts are documented
	loadBalancingConcepts := []string{
		"Round-Robin",
		"Fallback",
		"Model Selection",
		"task type",
	}

	for _, concept := range loadBalancingConcepts {
		if !strings.Contains(instructions, concept) {
			t.Errorf("load balancing concept missing: %s", concept)
		}
	}
}

func TestAgentTeamWorkflowSkill_ConfigReading(t *testing.T) {
	skill := NewAgentTeamWorkflowSkill("/tmp/test-workspace")
	instructions := skill.GetInstructions()

	// Verify config reading is documented
	configReading := []string{
		"~/.picoclaw/config.json",
		"agents",
		"list",
		"subagents",
		"defaults",
	}

	for _, item := range configReading {
		if !strings.Contains(instructions, item) {
			t.Errorf("config reading item missing: %s", item)
		}
	}
}

func TestAgentTeamWorkflowSkill_CodeExamples(t *testing.T) {
	skill := NewAgentTeamWorkflowSkill("/tmp/test-workspace")
	examples := skill.GetExamples()

	// Verify examples contain Go code
	if !strings.Contains(examples, "go") {
		t.Error("examples should contain Go code blocks")
	}

	// Verify examples have spawn flow
	if !strings.Contains(examples, "Spawn Flow") {
		t.Error("examples should include spawn flow")
	}

	// Verify examples have resource allocation
	if !strings.Contains(examples, "Resource Allocation") {
		t.Error("examples should include resource allocation")
	}
}
