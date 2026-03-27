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

// TestEmbeddedSkillsCount verifies that embedded skills are loaded
func TestEmbeddedSkillsCount(t *testing.T) {
	loader := skills.NewSkillsLoader("/tmp/workspace", "/tmp/global", "/tmp/builtin")
	allSkills := loader.ListSkills()

	// Count embedded skills
	embeddedCount := 0
	for _, skill := range allSkills {
		if skill.Source == "embedded" {
			embeddedCount++
		}
	}

	// We expect at least 150 embedded skills (158 converted)
	assert.GreaterOrEqual(t, embeddedCount, 150, "should have at least 150 embedded skills")
	t.Logf("Found %d embedded skills", embeddedCount)
}

// TestEmbeddedSkillLoad verifies that an embedded skill can be loaded
func TestEmbeddedSkillLoad(t *testing.T) {
	loader := skills.NewSkillsLoader("/tmp/workspace", "/tmp/global", "/tmp/builtin")

	// Try to load a known embedded skill
	content, ok := loader.LoadSkill("marketing-seo-specialist")
	require.True(t, ok, "should be able to load marketing-seo-specialist")
	assert.NotEmpty(t, content)
	assert.Contains(t, content, "SEO")
	assert.Contains(t, content, "search engine")
}

// TestEmbeddedSkillsListIncludes verifies that embedded skills appear in list
func TestEmbeddedSkillsListIncludes(t *testing.T) {
	loader := skills.NewSkillsLoader("/tmp/workspace", "/tmp/global", "/tmp/builtin")
	allSkills := loader.ListSkills()

	// Build a map of skill names
	skillNames := make(map[string]bool)
	for _, skill := range allSkills {
		skillNames[skill.Name] = true
	}

	// Check for some expected embedded skills
	expectedSkills := []string{
		"marketing-seo-specialist",
		"marketing-content-creator",
		"design-ux-architect",
		"design-ui-designer",
		"game-designer",
		"unity-architect",
		"product-manager",
		"sales-engineer",
		"support-support-responder",
		"testing-api-tester",
	}

	for _, expected := range expectedSkills {
		assert.True(t, skillNames[expected], "should have skill: %s", expected)
	}
}

// TestEmbeddedSkillsNotDuplicated verifies native skills are not duplicated
func TestEmbeddedSkillsNotDuplicated(t *testing.T) {
	loader := skills.NewSkillsLoader("/tmp/workspace", "/tmp/global", "/tmp/builtin")
	allSkills := loader.ListSkills()

	// Count occurrences of each skill name
	nameCount := make(map[string]int)
	for _, skill := range allSkills {
		nameCount[skill.Name]++
	}

	// Check that no skill appears more than once
	for name, count := range nameCount {
		assert.Equal(t, 1, count, "skill %s should appear only once, but appears %d times", name, count)
	}

	// Specifically check that native skills are not duplicated
	nativeSkills := []string{
		"backend_developer",
		"frontend_developer",
		"devops_engineer",
		"security_engineer",
		"qa_engineer",
		"data_engineer",
		"ml_engineer",
		"fullstack_developer",
		"researcher",
		"odoo_developer",
	}

	for _, native := range nativeSkills {
		assert.Equal(t, 1, nameCount[native], "native skill %s should not be duplicated", native)
	}
}

// TestEmbeddedSkillCategories verifies categories are present
func TestEmbeddedSkillCategories(t *testing.T) {
	loader := skills.NewSkillsLoader("/tmp/workspace", "/tmp/global", "/tmp/builtin")
	allSkills := loader.ListSkills()

	// Count skills by category (from path)
	categories := make(map[string]int)
	for _, skill := range allSkills {
		if skill.Source == "embedded" {
			// Extract category from path: embedded://category/name/SKILL.md
			if strings.HasPrefix(skill.Path, "embedded://") {
				parts := strings.Split(strings.TrimPrefix(skill.Path, "embedded://"), "/")
				if len(parts) >= 2 {
					categories[parts[0]]++
				}
			}
		}
	}

	// Expected categories
	expectedCategories := []string{
		"academic",
		"design",
		"engineering",
		"game-development",
		"marketing",
		"product",
		"sales",
		"support",
		"testing",
	}

	for _, category := range expectedCategories {
		assert.Greater(t, categories[category], 0, "should have skills in category: %s", category)
		t.Logf("Category %s: %d skills", category, categories[category])
	}
}

// TestEmbeddedSkillContent verifies content structure
func TestEmbeddedSkillContent(t *testing.T) {
	loader := skills.NewSkillsLoader("/tmp/workspace", "/tmp/global", "/tmp/builtin")

	// Load a skill and verify structure
	content, ok := loader.LoadSkill("design-ux-architect")
	require.True(t, ok)

	// Should have frontmatter stripped (should not start with ---)
	assert.False(t, strings.HasPrefix(strings.TrimSpace(content), "---"), "frontmatter should be stripped")

	// Should have actual skill content
	assert.Contains(t, content, "#")
	assert.Contains(t, content, "UX")
	assert.Contains(t, content, "Architect")
}

// TestEmbeddedSkillPriority verifies native skills have priority over embedded
func TestEmbeddedSkillPriority(t *testing.T) {
	loader := skills.NewSkillsLoader("/tmp/workspace", "/tmp/global", "/tmp/builtin")
	allSkills := loader.ListSkills()

	// Find backend_developer (native skill)
	var backendDev skills.SkillInfo
	found := false
	for _, skill := range allSkills {
		if skill.Name == "backend_developer" {
			backendDev = skill
			found = true
			break
		}
	}

	require.True(t, found, "backend_developer should exist")
	assert.Equal(t, "native", backendDev.Source, "backend_developer should be from native source")
}

// TestEmbeddedSkillMetadata verifies metadata is parsed correctly
func TestEmbeddedSkillMetadata(t *testing.T) {
	loader := skills.NewSkillsLoader("/tmp/workspace", "/tmp/global", "/tmp/builtin")
	allSkills := loader.ListSkills()

	// Find a specific skill and check metadata
	for _, skill := range allSkills {
		if skill.Name == "marketing-seo-specialist" {
			assert.NotEmpty(t, skill.Description, "should have description")
			assert.Equal(t, "embedded", skill.Source)
			assert.Contains(t, skill.Path, "embedded://")
			return
		}
	}

	t.Fatal("marketing-seo-specialist not found")
}

// TestEmbeddedSkillsBuildSummary verifies BuildSkillsSummary includes embedded
func TestEmbeddedSkillsBuildSummary(t *testing.T) {
	loader := skills.NewSkillsLoader("/tmp/workspace", "/tmp/global", "/tmp/builtin")
	summary := loader.BuildSkillsSummary()

	assert.NotEmpty(t, summary)
	assert.Contains(t, summary, "<skills>")
	assert.Contains(t, summary, "</skills>")

	// Should include embedded skills
	assert.Contains(t, summary, "embedded")
}
