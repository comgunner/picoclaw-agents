// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package skills

import (
	"fmt"
	"strings"
	"sync"
)

// SkillsManager manages skills with lazy loading.
// Instead of loading all 157 skills into the system prompt,
// it only loads skills when they are actually needed.
type SkillsManager struct {
	loader       *SkillsLoader
	allSkills    []SkillInfo
	loadedSkills map[string]*LoadedSkill
	activeSkills map[string]bool
	mu           sync.RWMutex
}

// LoadedSkill represents a skill that has been loaded into memory
type LoadedSkill struct {
	Info    SkillInfo
	Content string
}

// NewSkillsManager creates a new SkillsManager
func NewSkillsManager(loader *SkillsLoader) *SkillsManager {
	return &SkillsManager{
		loader:       loader,
		loadedSkills: make(map[string]*LoadedSkill),
		activeSkills: make(map[string]bool),
	}
}

// Init loads the list of available skills but doesn't load their content
func (sm *SkillsManager) Init() {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.allSkills = sm.loader.ListSkills()
}

// GetSkillsSummary returns only skill names (not descriptions)
// This is much lighter than BuildSkillsSummary which includes full descriptions
func (sm *SkillsManager) GetSkillsSummary() string {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	if len(sm.allSkills) == 0 {
		return ""
	}

	var names []string
	for _, s := range sm.allSkills {
		names = append(names, s.Name)
	}

	return fmt.Sprintf("Available skills (%d): %s", len(sm.allSkills), strings.Join(names, ", "))
}

// GetSkillsSummaryWithDescriptions returns skill names and short descriptions
// Lighter than full BuildSkillsSummary but more informative than just names
func (sm *SkillsManager) GetSkillsSummaryWithDescriptions() string {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	if len(sm.allSkills) == 0 {
		return ""
	}

	var lines []string
	lines = append(lines, "<skills>")
	for _, s := range sm.allSkills {
		// Truncate description to save tokens
		desc := s.Description
		if len(desc) > 100 {
			desc = desc[:100] + "..."
		}
		lines = append(lines, fmt.Sprintf("  <skill name=\"%s\" desc=\"%s\"/>", s.Name, desc))
	}
	lines = append(lines, "</skills>")

	return strings.Join(lines, "\n")
}

// LoadSkill loads a specific skill by name (lazy loading)
func (sm *SkillsManager) LoadSkill(name string) (*LoadedSkill, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// Check if already loaded
	if skill, ok := sm.loadedSkills[name]; ok {
		sm.activeSkills[name] = true
		return skill, nil
	}

	// Find skill info
	var skillInfo *SkillInfo
	for _, s := range sm.allSkills {
		if s.Name == name {
			skillInfo = &s
			break
		}
	}

	if skillInfo == nil {
		return nil, fmt.Errorf("skill %q not found", name)
	}

	// Load skill content
	content, ok := sm.loader.LoadSkill(name)
	if !ok {
		return nil, fmt.Errorf("failed to load skill %q", name)
	}

	loaded := &LoadedSkill{
		Info:    *skillInfo,
		Content: content,
	}

	sm.loadedSkills[name] = loaded
	sm.activeSkills[name] = true

	return loaded, nil
}

// LoadSkillsForNames loads multiple skills by name
func (sm *SkillsManager) LoadSkillsForNames(names []string) ([]*LoadedSkill, error) {
	var loaded []*LoadedSkill
	for _, name := range names {
		skill, err := sm.LoadSkill(name)
		if err != nil {
			continue // Skip skills that can't be loaded
		}
		loaded = append(loaded, skill)
	}
	return loaded, nil
}

// GetActiveSkillsContent returns the content of all active skills
func (sm *SkillsManager) GetActiveSkillsContent() string {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	var parts []string
	for name := range sm.activeSkills {
		if skill, ok := sm.loadedSkills[name]; ok {
			parts = append(parts, fmt.Sprintf("### Skill: %s\n\n%s", skill.Info.Name, skill.Content))
		}
	}

	if len(parts) == 0 {
		return ""
	}

	return strings.Join(parts, "\n\n---\n\n")
}

// DeactivateSkill removes a skill from active set
func (sm *SkillsManager) DeactivateSkill(name string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.activeSkills, name)
}

// GetLoadedCount returns number of loaded skills
func (sm *SkillsManager) GetLoadedCount() int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return len(sm.loadedSkills)
}

// GetActiveCount returns number of active skills
func (sm *SkillsManager) GetActiveCount() int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return len(sm.activeSkills)
}

// GetAllCount returns total number of available skills
func (sm *SkillsManager) GetAllCount() int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return len(sm.allSkills)
}
