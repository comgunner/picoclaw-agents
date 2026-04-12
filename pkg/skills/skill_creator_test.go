// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT

package skills

import (
	"strings"
	"testing"
)

func TestSkillCreatorSkill_Name(t *testing.T) {
	s := NewSkillCreatorSkill("/tmp")
	if got := s.Name(); got != "skill_creator" {
		t.Errorf("Name() = %q, want %q", got, "skill_creator")
	}
}

func TestSkillCreatorSkill_Description(t *testing.T) {
	s := NewSkillCreatorSkill("/tmp")
	desc := s.Description()
	if desc == "" {
		t.Fatal("Description() returned empty string")
	}
	if len(desc) > 200 {
		t.Errorf("Description() too long: %d chars", len(desc))
	}
	if !strings.Contains(strings.ToLower(desc), "skill") {
		t.Errorf("Description() should mention skill, got: %q", desc)
	}
}

func TestSkillCreatorSkill_GetInstructions(t *testing.T) {
	s := NewSkillCreatorSkill("/tmp")
	instructions := s.GetInstructions()
	if instructions == "" {
		t.Fatal("GetInstructions() returned empty string")
	}
	for _, step := range []string{"STEP 1", "STEP 2", "STEP 3", "STEP 4", "STEP 5", "STEP 6"} {
		if !strings.Contains(instructions, step) {
			t.Errorf("GetInstructions() should contain %q", step)
		}
	}
}

func TestSkillCreatorSkill_GetAntiPatterns(t *testing.T) {
	s := NewSkillCreatorSkill("/tmp")
	anti := s.GetAntiPatterns()
	if anti == "" {
		t.Fatal("GetAntiPatterns() returned empty string")
	}
	for _, keyword := range []string{"Node.js", "package.json", "Secrets"} {
		if !strings.Contains(anti, keyword) {
			t.Errorf("GetAntiPatterns() should mention %q", keyword)
		}
	}
}

func TestSkillCreatorSkill_GetExamples(t *testing.T) {
	s := NewSkillCreatorSkill("/tmp")
	examples := s.GetExamples()
	if examples == "" {
		t.Fatal("GetExamples() returned empty string")
	}
	count := strings.Count(examples, "Example")
	if count < 2 {
		t.Errorf("GetExamples() should contain at least 2 examples, found %d", count)
	}
}

func TestSkillCreatorSkill_BuildSkillContext(t *testing.T) {
	s := NewSkillCreatorSkill("/tmp")
	ctx := s.BuildSkillContext()
	if ctx == "" {
		t.Fatal("BuildSkillContext() returned empty string")
	}
	if !strings.Contains(ctx, "Skill Creator") {
		t.Errorf("BuildSkillContext() should contain 'Skill Creator'")
	}
	if !strings.Contains(ctx, "SKILL.md") {
		t.Errorf("BuildSkillContext() should mention SKILL.md output")
	}
	for _, section := range []string{"NATIVE SKILL", "PURPOSE", "OUTPUT", "WHEN TO USE", "STEP 1", "ANTI-PATTERNS", "EXAMPLES"} {
		if !strings.Contains(ctx, section) {
			t.Errorf("BuildSkillContext() should contain %q", section)
		}
	}
}

func TestSkillCreatorSkill_BuildSummary(t *testing.T) {
	s := NewSkillCreatorSkill("/tmp")
	summary := s.BuildSummary()
	if summary == "" {
		t.Fatal("BuildSummary() returned empty string")
	}
	if !strings.Contains(summary, `<skill name="skill_creator"`) {
		t.Errorf("BuildSummary() should contain XML skill tag")
	}
	for _, keyword := range []string{"purpose", "steps", "output", "validation"} {
		if !strings.Contains(summary, keyword) {
			t.Errorf("BuildSummary() should mention %q", keyword)
		}
	}
}

func TestSkillCreatorSkill_ConsistentSingleton(t *testing.T) {
	s1 := GetSkillCreatorSkill("/tmp/test1")
	s2 := GetSkillCreatorSkill("/tmp/test2")
	if s1 != s2 {
		t.Error("GetSkillCreatorSkill() should return same singleton instance")
	}
}
