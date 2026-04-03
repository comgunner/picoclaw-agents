// pkg/agent/department_router_test.go
//
// Tests for Department Router
// Based on concepts from icueth's picoclaw-agents fork (@icueth)
//
// Credits: @icueth (https://github.com/icueth)
// License: Same as base project (MIT)

package agent

import (
	"testing"

	"github.com/comgunner/picoclaw/pkg/config"
)

func createTestConfig() *config.Config {
	return &config.Config{
		Agents: config.AgentsConfig{
			Defaults: config.AgentDefaults{
				ModelName: "gpt-4",
			},
			DepartmentModels: map[string]map[string]any{
				"engineering": {
					"model": "gpt-4",
				},
				"design": {
					"model": "claude-3-sonnet",
				},
			},
		},
	}
}

func TestDepartmentRouter_DefaultFallback(t *testing.T) {
	cfg := createTestConfig()

	router := NewDepartmentRouter(cfg)

	// Agent without department should use default
	model := router.GetModelConfig("unknown_agent")

	if model.Model != "gpt-4" {
		t.Errorf("Expected gpt-4, got %s", model.Model)
	}
}

func TestDepartmentRouter_DepartmentOverride(t *testing.T) {
	cfg := createTestConfig()

	router := NewDepartmentRouter(cfg)

	// Register agent to engineering department
	router.RegisterAgent("senior_dev", "engineering")
	router.RegisterAgent("designer", "design")

	// Engineering agent should get engineering model
	engModel := router.GetModelConfig("senior_dev")
	if engModel.Model != "gpt-4" {
		t.Errorf("Engineering agent got wrong model: %s", engModel.Model)
	}

	// Design agent should get design model
	designModel := router.GetModelConfig("designer")
	if designModel.Model != "claude-3-sonnet" {
		t.Errorf("Design agent got wrong model: %s", designModel.Model)
	}
}

func TestDepartmentRouter_UnassignedDepartment(t *testing.T) {
	cfg := createTestConfig()

	router := NewDepartmentRouter(cfg)

	// Agent assigned to undefined department should fallback
	router.RegisterAgent("agent1", "undefined_dept")

	model := router.GetModelConfig("agent1")
	if model.Model != "gpt-4" {
		t.Errorf("Should fallback to default for undefined department")
	}
}

func TestDepartmentRouter_GetDefaultModel(t *testing.T) {
	cfg := createTestConfig()

	router := NewDepartmentRouter(cfg)

	defaultModel := router.GetDefaultModel()

	if defaultModel != "gpt-4" {
		t.Errorf("Expected gpt-4, got %s", defaultModel)
	}
}

func TestDepartmentRouter_GetDepartmentModel(t *testing.T) {
	cfg := createTestConfig()

	router := NewDepartmentRouter(cfg)

	// Get existing department
	engModel, err := router.GetDepartmentModel("engineering")
	if err != nil {
		t.Fatalf("Expected no error for existing department")
	}
	if engModel != "gpt-4" {
		t.Errorf("Expected gpt-4 for engineering")
	}

	// Get non-existent department
	_, err = router.GetDepartmentModel("nonexistent")
	if err == nil {
		t.Errorf("Expected error for non-existent department")
	}
}

func TestDepartmentRouter_ListDepartments(t *testing.T) {
	cfg := createTestConfig()

	router := NewDepartmentRouter(cfg)

	depts := router.ListDepartments()

	if len(depts) != 2 {
		t.Errorf("Expected 2 departments, got %d", len(depts))
	}

	if _, ok := depts["engineering"]; !ok {
		t.Errorf("Expected engineering department")
	}

	if _, ok := depts["design"]; !ok {
		t.Errorf("Expected design department")
	}
}

func TestDepartmentRouter_GetAgentDepartment(t *testing.T) {
	cfg := createTestConfig()

	router := NewDepartmentRouter(cfg)
	router.RegisterAgent("dev1", "engineering")
	router.RegisterAgent("designer1", "design")

	dept := router.GetAgentDepartment("dev1")
	if dept != "engineering" {
		t.Errorf("Expected engineering, got %s", dept)
	}

	dept = router.GetAgentDepartment("designer1")
	if dept != "design" {
		t.Errorf("Expected design, got %s", dept)
	}

	dept = router.GetAgentDepartment("unknown")
	if dept != "unknown" {
		t.Errorf("Expected unknown, got %s", dept)
	}
}

func TestDepartmentRouter_RemoveAgent(t *testing.T) {
	cfg := createTestConfig()

	router := NewDepartmentRouter(cfg)
	router.RegisterAgent("dev1", "engineering")

	// Verify registered
	dept := router.GetAgentDepartment("dev1")
	if dept != "engineering" {
		t.Errorf("Expected engineering before removal")
	}

	// Remove
	router.RemoveAgent("dev1")

	// Verify removed
	dept = router.GetAgentDepartment("dev1")
	if dept != "unknown" {
		t.Errorf("Expected unknown after removal")
	}
}

func TestDepartmentRouter_GetDepartmentAgents(t *testing.T) {
	cfg := createTestConfig()

	router := NewDepartmentRouter(cfg)
	router.RegisterAgent("dev1", "engineering")
	router.RegisterAgent("dev2", "engineering")
	router.RegisterAgent("designer1", "design")

	agents := router.GetDepartmentAgents("engineering")
	if len(agents) != 2 {
		t.Errorf("Expected 2 engineering agents, got %d", len(agents))
	}

	agents = router.GetDepartmentAgents("design")
	if len(agents) != 1 {
		t.Errorf("Expected 1 design agent, got %d", len(agents))
	}

	agents = router.GetDepartmentAgents("nonexistent")
	if len(agents) != 0 {
		t.Errorf("Expected 0 agents for nonexistent department")
	}
}

func TestDepartmentRouter_GetDepartmentCount(t *testing.T) {
	cfg := createTestConfig()

	router := NewDepartmentRouter(cfg)

	count := router.GetDepartmentCount()
	if count != 2 {
		t.Errorf("Expected 2 departments, got %d", count)
	}
}

func TestDepartmentRouter_GetAgentCount(t *testing.T) {
	cfg := createTestConfig()

	router := NewDepartmentRouter(cfg)
	router.RegisterAgent("dev1", "engineering")
	router.RegisterAgent("dev2", "engineering")
	router.RegisterAgent("designer1", "design")

	count := router.GetAgentCount()
	if count != 3 {
		t.Errorf("Expected 3 agents, got %d", count)
	}
}

func TestDepartmentRouter_UpdateDepartmentModel(t *testing.T) {
	cfg := createTestConfig()

	router := NewDepartmentRouter(cfg)

	// Update engineering model
	if err := router.UpdateDepartmentModel("engineering", "claude-3-opus"); err != nil {
		t.Fatalf("UpdateDepartmentModel failed: %v", err)
	}

	// Verify update
	engModel, _ := router.GetDepartmentModel("engineering")
	if engModel != "claude-3-opus" {
		t.Errorf("Expected claude-3-opus, got %s", engModel)
	}
}

func TestDepartmentRouter_EmptyConfig(t *testing.T) {
	cfg := &config.Config{
		Agents: config.AgentsConfig{
			Defaults: config.AgentDefaults{
				ModelName: "gpt-3.5-turbo",
			},
			DepartmentModels: nil,
		},
	}

	router := NewDepartmentRouter(cfg)

	// Should still work with default
	model := router.GetModelConfig("any_agent")
	if model.Model != "gpt-3.5-turbo" {
		t.Errorf("Expected gpt-3.5-turbo, got %s", model.Model)
	}

	// No departments
	if router.GetDepartmentCount() != 0 {
		t.Errorf("Expected 0 departments")
	}
}

func TestDepartmentRouter_ConcurrentAccess(t *testing.T) {
	cfg := createTestConfig()

	router := NewDepartmentRouter(cfg)

	// Concurrent registrations
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(idx int) {
			for j := 0; j < 10; j++ {
				agentID := "agent-" + string(rune(idx*10+j))
				router.RegisterAgent(agentID, "engineering")
			}
			done <- true
		}(i)
	}

	// Wait for all
	for i := 0; i < 10; i++ {
		<-done
	}

	// Should have 100 agents
	if router.GetAgentCount() != 100 {
		t.Errorf("Expected 100 agents, got %d", router.GetAgentCount())
	}

	// Concurrent reads
	for i := 0; i < 100; i++ {
		agentID := "agent-" + string(rune(i))
		model := router.GetModelConfig(agentID)
		if model.Model != "gpt-4" {
			t.Errorf("Expected gpt-4 for %s", agentID)
		}
	}
}
