// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package skills_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/comgunner/picoclaw/pkg/skills"
)

// ============================================================================
// BackendDeveloperSkill Tests
// ============================================================================

func TestBackendDeveloperSkillName(t *testing.T) {
	skill := skills.NewBackendDeveloperSkill("/tmp/workspace")
	assert.Equal(t, "backend_developer", skill.Name())
}

func TestBackendDeveloperSkillDescription(t *testing.T) {
	skill := skills.NewBackendDeveloperSkill("/tmp/workspace")
	desc := skill.Description()
	assert.NotEmpty(t, desc)
	assert.Contains(t, desc, "Backend")
}

func TestBackendDeveloperSkillGetInstructions(t *testing.T) {
	skill := skills.NewBackendDeveloperSkill("/tmp/workspace")
	instructions := skill.GetInstructions()
	assert.NotEmpty(t, instructions)
	assert.Contains(t, instructions, "API")
	assert.Contains(t, instructions, "Database")
}

func TestBackendDeveloperSkillBuildSkillContext(t *testing.T) {
	skill := skills.NewBackendDeveloperSkill("/tmp/workspace")
	context := skill.BuildSkillContext()
	assert.NotEmpty(t, context)
	assert.Contains(t, context, "Backend Developer")
	assert.Contains(t, context, "ROLE:")
}

func TestBackendDeveloperSkillBuildSummary(t *testing.T) {
	skill := skills.NewBackendDeveloperSkill("/tmp/workspace")
	summary := skill.BuildSummary()
	assert.NotEmpty(t, summary)
	assert.Contains(t, summary, `<skill`)
	assert.Contains(t, summary, "backend_developer")
}

// ============================================================================
// FrontendDeveloperSkill Tests
// ============================================================================

func TestFrontendDeveloperSkillName(t *testing.T) {
	skill := skills.NewFrontendDeveloperSkill("/tmp/workspace")
	assert.Equal(t, "frontend_developer", skill.Name())
}

func TestFrontendDeveloperSkillDescription(t *testing.T) {
	skill := skills.NewFrontendDeveloperSkill("/tmp/workspace")
	desc := skill.Description()
	assert.NotEmpty(t, desc)
	assert.Contains(t, desc, "Frontend")
}

func TestFrontendDeveloperSkillGetInstructions(t *testing.T) {
	skill := skills.NewFrontendDeveloperSkill("/tmp/workspace")
	instructions := skill.GetInstructions()
	assert.NotEmpty(t, instructions)
	assert.Contains(t, instructions, "Component")
	assert.Contains(t, instructions, "React")
}

func TestFrontendDeveloperSkillBuildSkillContext(t *testing.T) {
	skill := skills.NewFrontendDeveloperSkill("/tmp/workspace")
	context := skill.BuildSkillContext()
	assert.NotEmpty(t, context)
	assert.Contains(t, context, "Frontend Developer")
	assert.Contains(t, context, "ROLE:")
}

func TestFrontendDeveloperSkillBuildSummary(t *testing.T) {
	skill := skills.NewFrontendDeveloperSkill("/tmp/workspace")
	summary := skill.BuildSummary()
	assert.NotEmpty(t, summary)
	assert.Contains(t, summary, `<skill`)
	assert.Contains(t, summary, "frontend_developer")
}

// ============================================================================
// DevOpsEngineerSkill Tests
// ============================================================================

func TestDevOpsEngineerSkillName(t *testing.T) {
	skill := skills.NewDevOpsEngineerSkill("/tmp/workspace")
	assert.Equal(t, "devops_engineer", skill.Name())
}

func TestDevOpsEngineerSkillDescription(t *testing.T) {
	skill := skills.NewDevOpsEngineerSkill("/tmp/workspace")
	desc := skill.Description()
	assert.NotEmpty(t, desc)
	assert.Contains(t, desc, "DevOps")
}

func TestDevOpsEngineerSkillGetInstructions(t *testing.T) {
	skill := skills.NewDevOpsEngineerSkill("/tmp/workspace")
	instructions := skill.GetInstructions()
	assert.NotEmpty(t, instructions)
	assert.Contains(t, instructions, "CI/CD")
	assert.Contains(t, instructions, "Kubernetes")
}

func TestDevOpsEngineerSkillBuildSkillContext(t *testing.T) {
	skill := skills.NewDevOpsEngineerSkill("/tmp/workspace")
	context := skill.BuildSkillContext()
	assert.NotEmpty(t, context)
	assert.Contains(t, context, "DevOps Engineer")
	assert.Contains(t, context, "ROLE:")
}

func TestDevOpsEngineerSkillBuildSummary(t *testing.T) {
	skill := skills.NewDevOpsEngineerSkill("/tmp/workspace")
	summary := skill.BuildSummary()
	assert.NotEmpty(t, summary)
	assert.Contains(t, summary, `<skill`)
	assert.Contains(t, summary, "devops_engineer")
}

// ============================================================================
// SecurityEngineerSkill Tests
// ============================================================================

func TestSecurityEngineerSkillName(t *testing.T) {
	skill := skills.NewSecurityEngineerSkill("/tmp/workspace")
	assert.Equal(t, "security_engineer", skill.Name())
}

func TestSecurityEngineerSkillDescription(t *testing.T) {
	skill := skills.NewSecurityEngineerSkill("/tmp/workspace")
	desc := skill.Description()
	assert.NotEmpty(t, desc)
	assert.Contains(t, desc, "Security")
}

func TestSecurityEngineerSkillGetInstructions(t *testing.T) {
	skill := skills.NewSecurityEngineerSkill("/tmp/workspace")
	instructions := skill.GetInstructions()
	assert.NotEmpty(t, instructions)
	assert.Contains(t, instructions, "OWASP")
	assert.Contains(t, instructions, "threat")
}

func TestSecurityEngineerSkillBuildSkillContext(t *testing.T) {
	skill := skills.NewSecurityEngineerSkill("/tmp/workspace")
	context := skill.BuildSkillContext()
	assert.NotEmpty(t, context)
	assert.Contains(t, context, "Security Engineer")
	assert.Contains(t, context, "ROLE:")
}

func TestSecurityEngineerSkillBuildSummary(t *testing.T) {
	skill := skills.NewSecurityEngineerSkill("/tmp/workspace")
	summary := skill.BuildSummary()
	assert.NotEmpty(t, summary)
	assert.Contains(t, summary, `<skill`)
	assert.Contains(t, summary, "security_engineer")
}

// ============================================================================
// QAEngineerSkill Tests
// ============================================================================

func TestQAEngineerSkillName(t *testing.T) {
	skill := skills.NewQAEngineerSkill("/tmp/workspace")
	assert.Equal(t, "qa_engineer", skill.Name())
}

func TestQAEngineerSkillDescription(t *testing.T) {
	skill := skills.NewQAEngineerSkill("/tmp/workspace")
	desc := skill.Description()
	assert.NotEmpty(t, desc)
	assert.Contains(t, desc, "QA")
}

func TestQAEngineerSkillGetInstructions(t *testing.T) {
	skill := skills.NewQAEngineerSkill("/tmp/workspace")
	instructions := skill.GetInstructions()
	assert.NotEmpty(t, instructions)
	assert.Contains(t, instructions, "Test")
	// Check for automation (case-insensitive)
	assert.True(
		t,
		strings.Contains(strings.ToLower(instructions), "automation"),
		"instructions should mention automation",
	)
}

func TestQAEngineerSkillBuildSkillContext(t *testing.T) {
	skill := skills.NewQAEngineerSkill("/tmp/workspace")
	context := skill.BuildSkillContext()
	assert.NotEmpty(t, context)
	assert.Contains(t, context, "QA Engineer")
	assert.Contains(t, context, "ROLE:")
}

func TestQAEngineerSkillBuildSummary(t *testing.T) {
	skill := skills.NewQAEngineerSkill("/tmp/workspace")
	summary := skill.BuildSummary()
	assert.NotEmpty(t, summary)
	assert.Contains(t, summary, `<skill`)
	assert.Contains(t, summary, "qa_engineer")
}

// ============================================================================
// DataEngineerSkill Tests
// ============================================================================

func TestDataEngineerSkillName(t *testing.T) {
	skill := skills.NewDataEngineerSkill("/tmp/workspace")
	assert.Equal(t, "data_engineer", skill.Name())
}

func TestDataEngineerSkillDescription(t *testing.T) {
	skill := skills.NewDataEngineerSkill("/tmp/workspace")
	desc := skill.Description()
	assert.NotEmpty(t, desc)
	assert.Contains(t, desc, "Data")
}

func TestDataEngineerSkillGetInstructions(t *testing.T) {
	skill := skills.NewDataEngineerSkill("/tmp/workspace")
	instructions := skill.GetInstructions()
	assert.NotEmpty(t, instructions)
	assert.Contains(t, instructions, "ETL")
	assert.Contains(t, instructions, "pipeline")
}

func TestDataEngineerSkillBuildSkillContext(t *testing.T) {
	skill := skills.NewDataEngineerSkill("/tmp/workspace")
	context := skill.BuildSkillContext()
	assert.NotEmpty(t, context)
	assert.Contains(t, context, "Data Engineer")
	assert.Contains(t, context, "ROLE:")
}

func TestDataEngineerSkillBuildSummary(t *testing.T) {
	skill := skills.NewDataEngineerSkill("/tmp/workspace")
	summary := skill.BuildSummary()
	assert.NotEmpty(t, summary)
	assert.Contains(t, summary, `<skill`)
	assert.Contains(t, summary, "data_engineer")
}

// ============================================================================
// MLEngineerSkill Tests
// ============================================================================

func TestMLEngineerSkillName(t *testing.T) {
	skill := skills.NewMLEngineerSkill("/tmp/workspace")
	assert.Equal(t, "ml_engineer", skill.Name())
}

func TestMLEngineerSkillDescription(t *testing.T) {
	skill := skills.NewMLEngineerSkill("/tmp/workspace")
	desc := skill.Description()
	assert.NotEmpty(t, desc)
	assert.Contains(t, desc, "ML")
}

func TestMLEngineerSkillGetInstructions(t *testing.T) {
	skill := skills.NewMLEngineerSkill("/tmp/workspace")
	instructions := skill.GetInstructions()
	assert.NotEmpty(t, instructions)
	assert.Contains(t, instructions, "model")
	assert.Contains(t, instructions, "training")
}

func TestMLEngineerSkillBuildSkillContext(t *testing.T) {
	skill := skills.NewMLEngineerSkill("/tmp/workspace")
	context := skill.BuildSkillContext()
	assert.NotEmpty(t, context)
	assert.Contains(t, context, "ML Engineer")
	assert.Contains(t, context, "ROLE:")
}

func TestMLEngineerSkillBuildSummary(t *testing.T) {
	skill := skills.NewMLEngineerSkill("/tmp/workspace")
	summary := skill.BuildSummary()
	assert.NotEmpty(t, summary)
	assert.Contains(t, summary, `<skill`)
	assert.Contains(t, summary, "ml_engineer")
}

// ============================================================================
// Integration Tests: All Skills Have Consistent Structure
// ============================================================================

func TestAllEngineeringSkillsHaveConsistentStructure(t *testing.T) {
	workspace := "/tmp/workspace"

	skills := []struct {
		name  string
		skill interface {
			Name() string
			Description() string
			BuildSkillContext() string
			BuildSummary() string
		}
	}{
		{"backend_developer", skills.NewBackendDeveloperSkill(workspace)},
		{"frontend_developer", skills.NewFrontendDeveloperSkill(workspace)},
		{"devops_engineer", skills.NewDevOpsEngineerSkill(workspace)},
		{"security_engineer", skills.NewSecurityEngineerSkill(workspace)},
		{"qa_engineer", skills.NewQAEngineerSkill(workspace)},
		{"data_engineer", skills.NewDataEngineerSkill(workspace)},
		{"ml_engineer", skills.NewMLEngineerSkill(workspace)},
	}

	for _, s := range skills {
		t.Run(s.name, func(t *testing.T) {
			// Name should match expected
			assert.Equal(t, s.name, s.skill.Name())

			// Description should not be empty
			assert.NotEmpty(t, s.skill.Description())

			// BuildSkillContext should return non-empty string
			context := s.skill.BuildSkillContext()
			assert.NotEmpty(t, context)
			assert.True(t, len(context) > 100, "Skill context should be substantial")

			// BuildSummary should contain XML skill tag
			summary := s.skill.BuildSummary()
			assert.Contains(t, summary, "<skill")
			assert.Contains(t, summary, s.name)
		})
	}
}

func TestAllEngineeringSkillsHaveAntiPatterns(t *testing.T) {
	workspace := "/tmp/workspace"

	// Test that anti-patterns exist and contain expected content
	skillGetAntiPatterns := []struct {
		name            string
		getAntiPatterns func() string
	}{
		{"backend_developer", func() string { return skills.NewBackendDeveloperSkill(workspace).GetAntiPatterns() }},
		{"frontend_developer", func() string { return skills.NewFrontendDeveloperSkill(workspace).GetAntiPatterns() }},
		{"devops_engineer", func() string { return skills.NewDevOpsEngineerSkill(workspace).GetAntiPatterns() }},
		{"security_engineer", func() string { return skills.NewSecurityEngineerSkill(workspace).GetAntiPatterns() }},
		{"qa_engineer", func() string { return skills.NewQAEngineerSkill(workspace).GetAntiPatterns() }},
		{"data_engineer", func() string { return skills.NewDataEngineerSkill(workspace).GetAntiPatterns() }},
		{"ml_engineer", func() string { return skills.NewMLEngineerSkill(workspace).GetAntiPatterns() }},
	}

	for _, s := range skillGetAntiPatterns {
		t.Run(s.name, func(t *testing.T) {
			antiPatterns := s.getAntiPatterns()
			assert.NotEmpty(t, antiPatterns)
			assert.Contains(t, antiPatterns, "❌") // Should contain anti-pattern markers
		})
	}
}

func TestAllEngineeringSkillsHaveExamples(t *testing.T) {
	workspace := "/tmp/workspace"

	// Test that examples exist and contain expected content
	skillGetExamples := []struct {
		name        string
		getExamples func() string
	}{
		{"backend_developer", func() string { return skills.NewBackendDeveloperSkill(workspace).GetExamples() }},
		{"frontend_developer", func() string { return skills.NewFrontendDeveloperSkill(workspace).GetExamples() }},
		{"devops_engineer", func() string { return skills.NewDevOpsEngineerSkill(workspace).GetExamples() }},
		{"security_engineer", func() string { return skills.NewSecurityEngineerSkill(workspace).GetExamples() }},
		{"qa_engineer", func() string { return skills.NewQAEngineerSkill(workspace).GetExamples() }},
		{"data_engineer", func() string { return skills.NewDataEngineerSkill(workspace).GetExamples() }},
		{"ml_engineer", func() string { return skills.NewMLEngineerSkill(workspace).GetExamples() }},
	}

	for _, s := range skillGetExamples {
		t.Run(s.name, func(t *testing.T) {
			examples := s.getExamples()
			assert.NotEmpty(t, examples)
			assert.Contains(t, examples, "**Request:**") // Should contain request/response format
		})
	}
}

func TestEngineeringSkillContextsContainRequiredSections(t *testing.T) {
	workspace := "/tmp/workspace"

	skills := []interface {
		BuildSkillContext() string
	}{
		skills.NewBackendDeveloperSkill(workspace),
		skills.NewFrontendDeveloperSkill(workspace),
		skills.NewDevOpsEngineerSkill(workspace),
		skills.NewSecurityEngineerSkill(workspace),
		skills.NewQAEngineerSkill(workspace),
		skills.NewDataEngineerSkill(workspace),
		skills.NewMLEngineerSkill(workspace),
	}

	for _, skill := range skills {
		context := skill.BuildSkillContext()

		// Should contain header
		assert.Contains(t, context, "━━━━━━━━")

		// Should contain ROLE section
		assert.Contains(t, context, "**ROLE:**")

		// Should contain instructions
		assert.True(t, strings.Contains(context, "## ") || strings.Contains(context, "### "),
			"Context should contain markdown headers")
	}
}
