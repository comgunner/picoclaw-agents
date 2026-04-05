// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package onboard

import (
	"encoding/json"
	"testing"
)

func TestBuildAgentListJSON_Solo_NoSkills(t *testing.T) {
	result := buildAgentListJSON("solo", "", "anthropic/claude-sonnet-4-5-20251001", nil)

	var agents []agentEntry
	if err := json.Unmarshal([]byte(result), &agents); err != nil {
		t.Fatalf("Failed to unmarshal result: %v", err)
	}

	if len(agents) != 1 {
		t.Fatalf("Expected 1 agent, got %d", len(agents))
	}

	if agents[0].ID != "main" {
		t.Errorf("Expected ID 'main', got %q", agents[0].ID)
	}

	if agents[0].Default != true {
		t.Error("Expected Default=true")
	}

	if len(agents[0].Skills) != 0 {
		t.Errorf("Expected no skills, got %d", len(agents[0].Skills))
	}
}

func TestBuildAgentListJSON_Solo_WithSkills(t *testing.T) {
	skills := []string{"fullstack_developer", "agent_team_workflow"}
	result := buildAgentListJSON("solo", "", "anthropic/claude-sonnet-4-5-20251001", skills)

	var agents []agentEntry
	if err := json.Unmarshal([]byte(result), &agents); err != nil {
		t.Fatalf("Failed to unmarshal result: %v", err)
	}

	if len(agents) != 1 {
		t.Fatalf("Expected 1 agent, got %d", len(agents))
	}

	if len(agents[0].Skills) != 2 {
		t.Errorf("Expected 2 skills, got %d", len(agents[0].Skills))
	}

	if agents[0].Skills[0] != "fullstack_developer" {
		t.Errorf("Expected first skill 'fullstack_developer', got %q", agents[0].Skills[0])
	}
}

func TestBuildAgentListJSON_DevTeam(t *testing.T) {
	result := buildAgentListJSON("team", "dev", "anthropic/claude-sonnet-4-5-20251001", nil)

	var agents []agentEntry
	if err := json.Unmarshal([]byte(result), &agents); err != nil {
		t.Fatalf("Failed to unmarshal result: %v", err)
	}

	// Dev team should have 9 agents: 1 manager + 8 specialists
	if len(agents) != 9 {
		t.Fatalf("Expected 9 agents for dev team, got %d", len(agents))
	}

	// First agent should be engineering_manager with Default=true
	if agents[0].ID != "engineering_manager" {
		t.Errorf("Expected first agent ID 'engineering_manager', got %q", agents[0].ID)
	}

	if agents[0].Default != true {
		t.Error("Expected engineering_manager to have Default=true")
	}

	// Engineering manager should have subagents configured
	if agents[0].Subagents == nil {
		t.Fatal("Expected engineering_manager to have subagents configured")
	}

	if len(agents[0].Subagents.AllowAgents) != 8 {
		t.Errorf("Expected 8 allowed subagents, got %d", len(agents[0].Subagents.AllowAgents))
	}

	if agents[0].Subagents.MaxSpawnDepth != 3 {
		t.Errorf("Expected MaxSpawnDepth=3, got %d", agents[0].Subagents.MaxSpawnDepth)
	}
}

func TestBuildAgentListJSON_ResearchTeam(t *testing.T) {
	result := buildAgentListJSON("team", "research", "anthropic/claude-sonnet-4-5-20251001", nil)

	var agents []agentEntry
	if err := json.Unmarshal([]byte(result), &agents); err != nil {
		t.Fatalf("Failed to unmarshal result: %v", err)
	}

	// Research team should have 3 agents
	if len(agents) != 3 {
		t.Fatalf("Expected 3 agents for research team, got %d", len(agents))
	}

	// First agent should be coordinator
	if agents[0].ID != "coordinator" {
		t.Errorf("Expected first agent ID 'coordinator', got %q", agents[0].ID)
	}

	// Should have researcher and analyst
	foundResearcher := false
	foundAnalyst := false
	for _, a := range agents {
		if a.ID == "researcher" {
			foundResearcher = true
		}
		if a.ID == "analyst" {
			foundAnalyst = true
		}
	}

	if !foundResearcher {
		t.Error("Expected to find 'researcher' agent")
	}
	if !foundAnalyst {
		t.Error("Expected to find 'analyst' agent")
	}
}

func TestBuildAgentListJSON_GeneralTeam(t *testing.T) {
	result := buildAgentListJSON("team", "general", "anthropic/claude-sonnet-4-5-20251001", nil)

	var agents []agentEntry
	if err := json.Unmarshal([]byte(result), &agents); err != nil {
		t.Fatalf("Failed to unmarshal result: %v", err)
	}

	// General team should have 3 agents
	if len(agents) != 3 {
		t.Fatalf("Expected 3 agents for general team, got %d", len(agents))
	}

	// Should have orchestrator, worker_a, worker_b
	foundOrchestrator := false
	foundWorkerA := false
	foundWorkerB := false
	for _, a := range agents {
		if a.ID == "orchestrator" {
			foundOrchestrator = true
		}
		if a.ID == "worker_a" {
			foundWorkerA = true
		}
		if a.ID == "worker_b" {
			foundWorkerB = true
		}
	}

	if !foundOrchestrator {
		t.Error("Expected to find 'orchestrator' agent")
	}
	if !foundWorkerA {
		t.Error("Expected to find 'worker_a' agent")
	}
	if !foundWorkerB {
		t.Error("Expected to find 'worker_b' agent")
	}
}

func TestDevTeamAgents_HasOrchestrator(t *testing.T) {
	agents := devTeamAgents("test-model")

	if len(agents) == 0 {
		t.Fatal("Expected at least one agent")
	}

	// First agent should be orchestrator with subagents
	if agents[0].Subagents == nil {
		t.Fatal("Expected first agent to have subagents configured")
	}

	// Should allow specialists
	specialists := []string{
		"backend_dev", "frontend_dev", "devops_eng", "qa_eng",
		"security_eng", "data_eng", "ml_eng", "researcher",
	}

	for _, spec := range specialists {
		found := false
		for _, allowed := range agents[0].Subagents.AllowAgents {
			if allowed == spec {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected specialist %q to be in AllowAgents", spec)
		}
	}
}

func TestDevTeamAgents_SubagentsConfigured(t *testing.T) {
	agents := devTeamAgents("test-model")

	// Check that specialists have empty subagents (no nested spawning)
	for i := 1; i < len(agents); i++ {
		if agents[i].Subagents == nil {
			t.Errorf("Expected agent %s to have subagents configured (even if empty)", agents[i].ID)
		} else if len(agents[i].Subagents.AllowAgents) != 0 {
			t.Errorf("Expected specialist %s to have no allowed subagents, got %d",
				agents[i].ID, len(agents[i].Subagents.AllowAgents))
		}
	}
}

func TestGetNativeSkills_ReturnsAllSkills(t *testing.T) {
	skills := getNativeSkills()

	expectedCount := 14 // queue_batch, binance_mcp, fullstack_developer, etc.
	if len(skills) != expectedCount {
		t.Errorf("Expected %d skills, got %d", expectedCount, len(skills))
	}

	// Check for some key skills
	expectedSkills := []string{
		"fullstack_developer",
		"agent_team_workflow",
		"backend_developer",
		"binance_mcp",
	}

	for _, expected := range expectedSkills {
		found := false
		for _, skill := range skills {
			if skill == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected to find skill %q", expected)
		}
	}
}

func TestGetSkillDescription_ReturnsDescriptions(t *testing.T) {
	skills := getNativeSkills()

	for _, skill := range skills {
		desc := getSkillDescription(skill)
		if desc == "" {
			t.Errorf("Expected description for skill %q", skill)
		}
	}
}
