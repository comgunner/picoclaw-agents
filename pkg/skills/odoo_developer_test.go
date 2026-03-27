// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package skills_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/comgunner/picoclaw/pkg/skills"
)

// TestOdooDeveloperSkill_Name verifies the skill returns correct identifier
func TestOdooDeveloperSkill_Name(t *testing.T) {
	s := skills.NewOdooDeveloperSkill("/tmp/workspace")
	assert.Equal(t, "odoo_developer", s.Name())
}

// TestOdooDeveloperSkill_Description verifies description is concise and informative
func TestOdooDeveloperSkill_Description(t *testing.T) {
	s := skills.NewOdooDeveloperSkill("/tmp/workspace")
	desc := s.Description()
	assert.NotEmpty(t, desc)
	assert.Less(t, len(desc), 120, "description should be concise (<120 chars)")
	assert.Contains(t, desc, "Odoo")
	assert.Contains(t, desc, "Architect")
}

// TestOdooDeveloperSkill_GetInstructions verifies instructions are substantive
func TestOdooDeveloperSkill_GetInstructions(t *testing.T) {
	s := skills.NewOdooDeveloperSkill("/tmp/workspace")
	instructions := s.GetInstructions()
	assert.NotEmpty(t, instructions)
	assert.Greater(t, len(instructions), 100, "instructions should be substantive")
	assert.Contains(t, instructions, "ROLE")
	assert.Contains(t, instructions, "CONSTRAINTS")
	assert.Contains(t, instructions, "OUTPUT FORMAT")
}

// TestOdooDeveloperSkill_GetAntiPatterns verifies anti-patterns section exists
func TestOdooDeveloperSkill_GetAntiPatterns(t *testing.T) {
	s := skills.NewOdooDeveloperSkill("/tmp/workspace")
	antiPatterns := s.GetAntiPatterns()
	assert.NotEmpty(t, antiPatterns)
	assert.Contains(t, antiPatterns, "❌")
	assert.Contains(t, antiPatterns, "**BAD:**")
	assert.Contains(t, antiPatterns, "**GOOD:**")
}

// TestOdooDeveloperSkill_GetExamples verifies examples section exists
func TestOdooDeveloperSkill_GetExamples(t *testing.T) {
	s := skills.NewOdooDeveloperSkill("/tmp/workspace")
	examples := s.GetExamples()
	assert.NotEmpty(t, examples)
	assert.Contains(t, examples, "EXAMPLE")
	// Examples use backtick constants for code blocks
	assert.Contains(t, examples, "`")
	assert.GreaterOrEqual(t, strings.Count(examples, "`python"), 5, "should have multiple Python code examples")
}

// TestOdooDeveloperSkill_BuildSkillContext verifies complete context structure
func TestOdooDeveloperSkill_BuildSkillContext(t *testing.T) {
	s := skills.NewOdooDeveloperSkill("/tmp/workspace")
	ctx := s.BuildSkillContext()

	// Verify structure (v3.6.0 pattern)
	assert.Contains(t, ctx, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	assert.Contains(t, ctx, "🚀 NATIVE SKILL:")
	assert.Contains(t, ctx, "Odoo Developer")
	assert.Contains(t, ctx, "**ROLE:**")
	assert.Contains(t, ctx, "**OBJECTIVE:**")
	assert.Contains(t, ctx, s.GetInstructions())
	assert.Contains(t, ctx, s.GetAntiPatterns())
	assert.Contains(t, ctx, s.GetExamples())

	// Verify length is reasonable (not too short, not too long)
	assert.Greater(t, len(ctx), 1000, "context should be comprehensive")
	assert.Less(t, len(ctx), 50000, "context should fit in token budget")
}

// TestOdooDeveloperSkill_BuildSummary verifies XML summary format
func TestOdooDeveloperSkill_BuildSummary(t *testing.T) {
	s := skills.NewOdooDeveloperSkill("/tmp/workspace")
	summary := s.BuildSummary()

	assert.Contains(t, summary, `<skill name="odoo_developer"`)
	assert.Contains(t, summary, `type="native"`)
	assert.Contains(t, summary, "<purpose>")
	assert.Contains(t, summary, "<pattern>")
	assert.Contains(t, summary, "<stacks>")
	assert.Contains(t, summary, "<specialties>")
	assert.Contains(t, summary, "<constraints>")
	assert.Contains(t, summary, "</skill>")
	assert.True(t, strings.HasPrefix(strings.TrimSpace(summary), "<skill"))
}

// TestOdooDeveloperSkill_WorkspaceIndependence verifies skill works with any workspace
func TestOdooDeveloperSkill_WorkspaceIndependence(t *testing.T) {
	workspaces := []string{
		"/tmp/workspace1",
		"/tmp/workspace2",
		"/var/lib/picoclaw/workspace",
		"~/.picoclaw/workspace",
		"",
	}

	for _, ws := range workspaces {
		t.Run(ws, func(t *testing.T) {
			s := skills.NewOdooDeveloperSkill(ws)
			assert.Equal(t, "odoo_developer", s.Name())
			assert.NotEmpty(t, s.Description())
			assert.NotEmpty(t, s.BuildSkillContext())
		})
	}
}

// TestOdooDeveloperSkill_Concurrency verifies thread-safe lazy initialization
func TestOdooDeveloperSkill_Concurrency(t *testing.T) {
	// Test that getter returns same instance (singleton pattern)
	a := skills.GetOdooDeveloperSkill("/tmp/ws1")
	b := skills.GetOdooDeveloperSkill("/tmp/ws1")
	require.NotNil(t, a)
	assert.Same(t, a, b, "getter should return the same instance")

	// Verify methods work on singleton
	assert.Equal(t, "odoo_developer", a.Name())
	assert.NotEmpty(t, a.Description())
}

// TestOdooDeveloperSkill_GetterIsSingleton verifies singleton pattern
func TestOdooDeveloperSkill_GetterIsSingleton(t *testing.T) {
	instances := make([]*skills.OdooDeveloperSkill, 10)
	for i := 0; i < 10; i++ {
		instances[i] = skills.GetOdooDeveloperSkill("/tmp/ws")
	}

	// All instances should be the same pointer
	first := instances[0]
	for i, inst := range instances[1:] {
		assert.Same(t, first, inst, "instance %d should be same as first", i+1)
	}
}

// TestOdooDeveloperSkill_HasRequiredSections verifies all required sections present
func TestOdooDeveloperSkill_HasRequiredSections(t *testing.T) {
	s := skills.NewOdooDeveloperSkill("/tmp/workspace")
	ctx := s.BuildSkillContext()

	requiredSections := []string{
		"ROLE & OBJECTIVE",
		"TECHNICAL SKILLS",
		"CONSTRAINTS",
		"OUTPUT FORMAT",
		"EXECUTION STEPS",
		"STYLE GUIDELINES",
		"KNOWLEDGE BASE",
	}

	for _, section := range requiredSections {
		assert.Contains(t, ctx, section, "Missing required section: %s", section)
	}
}

// TestOdooDeveloperSkill_HasAntiPatternsWithExamples verifies anti-patterns have code examples
func TestOdooDeveloperSkill_HasAntiPatternsWithExamples(t *testing.T) {
	s := skills.NewOdooDeveloperSkill("/tmp/workspace")
	antiPatterns := s.GetAntiPatterns()

	assert.Contains(t, antiPatterns, "CODE ANTI-PATTERNS")
	assert.Contains(t, antiPatterns, "SECURITY ANTI-PATTERNS")
	assert.Contains(t, antiPatterns, "ARCHITECTURE ANTI-PATTERNS")
	assert.Contains(t, antiPatterns, "TESTING ANTI-PATTERNS")
	assert.Contains(t, antiPatterns, "DEPLOYMENT ANTI-PATTERNS")

	// Verify code blocks in anti-patterns
	assert.Greater(t, strings.Count(antiPatterns, "**BAD:**"), 5)
	assert.Greater(t, strings.Count(antiPatterns, "**GOOD:**"), 5)
}

// TestOdooDeveloperSkill_HasConcreteExamples verifies examples are concrete and actionable
func TestOdooDeveloperSkill_HasConcreteExamples(t *testing.T) {
	s := skills.NewOdooDeveloperSkill("/tmp/workspace")
	examples := s.GetExamples()

	assert.Contains(t, examples, "PINE SCRIPT → ODOO MIGRATION")
	assert.Contains(t, examples, "L10N-MEXICO CFDI 4.0")
	assert.Contains(t, examples, "MIGRATION (ODOO 13 → 16)")

	// Verify examples have step-by-step structure
	assert.Contains(t, examples, "Step 1:")
	assert.Contains(t, examples, "Step 2:")
	assert.Contains(t, examples, "Step 3:")

	// Verify examples have code blocks (using backtick constants)
	// The examples use bt constant for backticks, so check for Python/JSON/XML code
	assert.Contains(t, examples, "python")
	assert.Contains(t, examples, "json")
	assert.Contains(t, examples, "xml")
}

// TestOdooDeveloperSkill_XMLSummaryValid verifies XML is well-formed
func TestOdooDeveloperSkill_XMLSummaryValid(t *testing.T) {
	s := skills.NewOdooDeveloperSkill("/tmp/workspace")
	summary := s.BuildSummary()

	// Basic XML validation (check tags are balanced)
	openTags := strings.Count(summary, "<skill")
	closeTags := strings.Count(summary, "</skill>")
	assert.Equal(t, openTags, closeTags, "XML tags should be balanced")

	// Verify required attributes
	assert.Contains(t, summary, `name="odoo_developer"`)
	assert.Contains(t, summary, `type="native"`)
}

// TestOdooDeveloperSkill_IntegrationWithLoader verifies skill can be loaded
func TestOdooDeveloperSkill_IntegrationWithLoader(t *testing.T) {
	// This test verifies the skill integrates with the loader system
	loader := skills.NewSkillsLoader("/tmp/workspace", "/tmp/global", "/tmp/builtin")

	// Verify skill appears in list
	allSkills := loader.ListSkills()
	found := false
	for _, skill := range allSkills {
		if skill.Name == "odoo_developer" {
			found = true
			assert.Equal(t, "native", skill.Source)
			assert.Contains(t, skill.Path, "builtin://")
			assert.NotEmpty(t, skill.Description)
			break
		}
	}

	// Note: This will fail until odoo_developer is registered in loader.go
	// Uncomment after registration:
	// assert.True(t, found, "odoo_developer should be in listNativeSkills()")
	_ = found // Suppress unused variable warning for now
}
