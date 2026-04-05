// pkg/agent/orchestrator_test.go
//
// Tests for A2A Orchestrator
// Based on concepts from icueth's picoclaw-agents fork (@icueth)
//
// Credits: @icueth (https://github.com/icueth)
// License: Same as base project (MIT)

package agent

import (
	"testing"
	"time"

	"github.com/comgunner/picoclaw/pkg/mailbox"
)

func TestOrchestrator_RegisterAgent(t *testing.T) {
	orch := NewA2AOrchestrator(100)

	// Create a mock agent
	agent := &AgentInstance{
		ID:   "test_dev",
		Role: "developer",
	}

	// Register
	if err := orch.RegisterAgent(agent); err != nil {
		t.Fatalf("Failed to register: %v", err)
	}

	// Verify
	if agent.mailbox == nil {
		t.Errorf("Agent mailbox not initialized")
	}

	// Check agent count
	if orch.GetAgentCount() != 1 {
		t.Errorf("Expected 1 agent, got %d", orch.GetAgentCount())
	}
}

func TestOrchestrator_RegisterDuplicate(t *testing.T) {
	orch := NewA2AOrchestrator(100)

	agent := &AgentInstance{ID: "dev", Role: "developer"}
	orch.RegisterAgent(agent)

	// Try to register again
	agent2 := &AgentInstance{ID: "dev", Role: "developer"}
	if err := orch.RegisterAgent(agent2); err == nil {
		t.Errorf("Should reject duplicate registration")
	}
}

func TestOrchestrator_UnregisterAgent(t *testing.T) {
	orch := NewA2AOrchestrator(100)

	agent := &AgentInstance{ID: "dev", Role: "developer"}
	orch.RegisterAgent(agent)

	// Unregister
	if err := orch.UnregisterAgent("dev"); err != nil {
		t.Fatalf("Unregister failed: %v", err)
	}

	// Check agent count
	if orch.GetAgentCount() != 0 {
		t.Errorf("Expected 0 agents after unregister, got %d", orch.GetAgentCount())
	}
}

func TestOrchestrator_Unregister_NotFound(t *testing.T) {
	orch := NewA2AOrchestrator(100)

	// Try to unregister non-existent agent
	if err := orch.UnregisterAgent("nonexistent"); err == nil {
		t.Errorf("Expected error for non-existent agent")
	}
}

func TestOrchestrator_SendMessage(t *testing.T) {
	orch := NewA2AOrchestrator(100)

	// Register agents
	orch.RegisterAgent(&AgentInstance{ID: "pm", Role: "coordinator"})
	orch.RegisterAgent(&AgentInstance{ID: "dev", Role: "developer"})

	// Send message
	if err := orch.SendMessage("pm", "dev", mailbox.MessageTypeTask, mailbox.PriorityHigh, "Review PR"); err != nil {
		t.Fatalf("SendMessage failed: %v", err)
	}

	// Verify reception
	devMb, _ := orch.mailboxHub.Get("dev")
	if devMb.GetUnreadCount() != 1 {
		t.Errorf("Dev should have 1 unread message")
	}
}

func TestOrchestrator_BroadcastMessage(t *testing.T) {
	orch := NewA2AOrchestrator(100)

	// Register 3 agents
	orch.RegisterAgent(&AgentInstance{ID: "pm", Role: "coordinator"})
	orch.RegisterAgent(&AgentInstance{ID: "dev", Role: "developer"})
	orch.RegisterAgent(&AgentInstance{ID: "qa", Role: "tester"})

	// Broadcast from pm
	failed := orch.BroadcastMessage("pm", mailbox.MessageTypeBroadcast, "Meeting at 3pm")

	if len(failed) != 0 {
		t.Errorf("Broadcast should not fail")
	}

	// Verify all except pm received
	devMb, _ := orch.mailboxHub.Get("dev")
	qaMb, _ := orch.mailboxHub.Get("qa")

	if devMb.GetUnreadCount() != 1 || qaMb.GetUnreadCount() != 1 {
		t.Errorf("All agents should receive broadcast")
	}

	// Sender should not receive
	pmMb, _ := orch.mailboxHub.Get("pm")
	if pmMb.GetUnreadCount() != 0 {
		t.Errorf("Sender should not receive broadcast")
	}
}

func TestOrchestrator_AssignTask(t *testing.T) {
	orch := NewA2AOrchestrator(100)

	orch.RegisterAgent(&AgentInstance{ID: "pm", Role: "coordinator"})
	orch.RegisterAgent(&AgentInstance{ID: "dev", Role: "developer"})

	task := "Implement feature X"
	if err := orch.AssignTask("pm", "dev", task, mailbox.PriorityHigh); err != nil {
		t.Fatalf("AssignTask failed: %v", err)
	}

	// Verify task stored in shared context
	keys := orch.sharedContext.Keys()
	found := false
	for _, key := range keys {
		if key[:5] == "task:" {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Task not stored in shared context")
	}

	// Verify dev received message
	devMb, _ := orch.mailboxHub.Get("dev")
	if devMb.GetUnreadCount() != 1 {
		t.Errorf("Dev should have 1 unread message")
	}
}

func TestOrchestrator_ReportTaskComplete(t *testing.T) {
	orch := NewA2AOrchestrator(100)
	orch.RegisterAgent(&AgentInstance{ID: "dev", Role: "developer"})

	taskKey := "task:dev:123456"
	result := "Feature implemented"

	if err := orch.ReportTaskComplete("dev", taskKey, result); err != nil {
		t.Fatalf("ReportTaskComplete failed: %v", err)
	}

	// Verify result stored
	keys := orch.sharedContext.Keys()
	found := false
	for _, key := range keys {
		if len(key) >= 7 && key[len(key)-7:] == ":result" {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Result not stored in shared context")
	}
}

func TestOrchestrator_GetAgentStatus(t *testing.T) {
	orch := NewA2AOrchestrator(100)

	agent := &AgentInstance{
		ID:     "dev",
		Role:   "developer",
		Status: "idle",
	}
	orch.RegisterAgent(agent)

	status := orch.GetAgentStatus("dev")

	if status["id"] != "dev" {
		t.Errorf("Expected id 'dev', got %v", status["id"])
	}

	if status["role"] != "developer" {
		t.Errorf("Expected role 'developer', got %v", status["role"])
	}

	if status["status"] != "idle" {
		t.Errorf("Expected status 'idle', got %v", status["status"])
	}
}

func TestOrchestrator_GetAgentStatus_NotFound(t *testing.T) {
	orch := NewA2AOrchestrator(100)

	status := orch.GetAgentStatus("nonexistent")

	if _, ok := status["error"]; !ok {
		t.Errorf("Expected error for non-existent agent")
	}
}

func TestOrchestrator_Stats(t *testing.T) {
	orch := NewA2AOrchestrator(100)

	orch.RegisterAgent(&AgentInstance{ID: "agent1", Role: "dev"})
	orch.RegisterAgent(&AgentInstance{ID: "agent2", Role: "qa"})

	stats := orch.GetOrchestrationStats()

	if stats["agents"].(int) != 2 {
		t.Errorf("Expected 2 agents in stats")
	}

	if _, ok := stats["mailbox_stats"]; !ok {
		t.Errorf("Mailbox stats missing")
	}

	if _, ok := stats["context_size"]; !ok {
		t.Errorf("Context size missing")
	}
}

func TestOrchestrator_ListAgents(t *testing.T) {
	orch := NewA2AOrchestrator(100)

	orch.RegisterAgent(&AgentInstance{ID: "pm", Role: "coordinator"})
	orch.RegisterAgent(&AgentInstance{ID: "dev", Role: "developer"})
	orch.RegisterAgent(&AgentInstance{ID: "qa", Role: "tester"})

	agents := orch.ListAgents()

	if len(agents) != 3 {
		t.Errorf("Expected 3 agents, got %d", len(agents))
	}

	// Check all agents are listed
	expected := map[string]bool{
		"pm":  true,
		"dev": true,
		"qa":  true,
	}

	for _, agent := range agents {
		if !expected[agent] {
			t.Errorf("Unexpected agent: %s", agent)
		}
	}
}

func TestOrchestrator_GetSharedContext(t *testing.T) {
	orch := NewA2AOrchestrator(100)

	ctx := orch.GetSharedContext()
	if ctx == nil {
		t.Errorf("Shared context should not be nil")
	}

	// Use context
	ctx.Set("test_key", "test_value")
	val, ok := ctx.Get("test_key")
	if !ok || val != "test_value" {
		t.Errorf("Shared context not working")
	}
}

func TestOrchestrator_GetMailboxHub(t *testing.T) {
	orch := NewA2AOrchestrator(100)

	hub := orch.GetMailboxHub()
	if hub == nil {
		t.Errorf("Mailbox hub should not be nil")
	}

	// Use hub
	hub.Register("test_agent")
	_, ok := hub.Get("test_agent")
	if !ok {
		t.Errorf("Mailbox hub not working")
	}
}

func TestOrchestrator_Stop(t *testing.T) {
	orch := NewA2AOrchestrator(100)
	orch.RegisterAgent(&AgentInstance{ID: "dev", Role: "developer"})

	// Start cleanup
	orch.StartCleanup(100 * time.Millisecond)

	// Stop
	if err := orch.Stop(); err != nil {
		t.Fatalf("Stop failed: %v", err)
	}

	// Should not panic after stop
	time.Sleep(150 * time.Millisecond)
}

func TestOrchestrator_WorkflowDiscovery(t *testing.T) {
	orch := NewA2AOrchestrator(100)

	// Register agents
	orch.RegisterAgent(&AgentInstance{ID: "pm", Role: "coordinator"})
	orch.RegisterAgent(&AgentInstance{ID: "dev", Role: "developer"})
	orch.RegisterAgent(&AgentInstance{ID: "qa", Role: "tester"})

	// Run discovery
	if err := orch.WorkflowDiscovery(); err != nil {
		t.Fatalf("WorkflowDiscovery failed: %v", err)
	}

	// Check capabilities stored
	keys := orch.sharedContext.Keys()
	capabilitiesFound := 0
	for _, key := range keys {
		if len(key) > 6 && key[:6] == "agent:" && len(key) >= 13 && key[len(key)-13:] == ":capabilities" {
			capabilitiesFound++
		}
	}

	if capabilitiesFound != 3 {
		t.Errorf("Expected 3 capability entries, got %d", capabilitiesFound)
	}
}

func TestOrchestrator_WorkflowPlanning(t *testing.T) {
	orch := NewA2AOrchestrator(100)

	orch.RegisterAgent(&AgentInstance{ID: "pm", Role: "coordinator"})
	orch.RegisterAgent(&AgentInstance{ID: "dev", Role: "developer"})
	orch.RegisterAgent(&AgentInstance{ID: "qa", Role: "tester"})

	// Assign tasks
	tasks := map[string]string{
		"dev": "Implement feature",
		"qa":  "Test feature",
	}

	if err := orch.WorkflowPlanning("pm", tasks); err != nil {
		t.Fatalf("WorkflowPlanning failed: %v", err)
	}

	// Verify tasks assigned
	devMb, _ := orch.mailboxHub.Get("dev")
	qaMb, _ := orch.mailboxHub.Get("qa")

	if devMb.GetUnreadCount() != 1 {
		t.Errorf("Dev should have 1 task")
	}
	if qaMb.GetUnreadCount() != 1 {
		t.Errorf("QA should have 1 task")
	}
}

func TestOrchestrator_RunFullWorkflow(t *testing.T) {
	orch := NewA2AOrchestrator(100)

	orch.RegisterAgent(&AgentInstance{ID: "pm", Role: "coordinator"})
	orch.RegisterAgent(&AgentInstance{ID: "dev", Role: "developer"})
	orch.RegisterAgent(&AgentInstance{ID: "qa", Role: "tester"})

	// Run full workflow
	tasks := map[string]string{
		"dev": "Implement feature X",
		"qa":  "Test feature X",
	}

	if err := orch.RunFullWorkflow("pm", tasks); err != nil {
		t.Fatalf("RunFullWorkflow failed: %v", err)
	}

	// Verify workflow completed
	stats := orch.GetOrchestrationStats()
	if stats["agents"].(int) != 3 {
		t.Errorf("Expected 3 agents")
	}
}

func TestOrchestrator_ConcurrentRegister(t *testing.T) {
	orch := NewA2AOrchestrator(100)

	// Concurrent registrations
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(idx int) {
			for j := 0; j < 10; j++ {
				agent := &AgentInstance{
					ID:   "agent-" + string(rune(idx*10+j)),
					Role: "worker",
				}
				orch.RegisterAgent(agent)
			}
			done <- true
		}(i)
	}

	// Wait for all
	for i := 0; i < 10; i++ {
		<-done
	}

	// Should have 100 agents
	if orch.GetAgentCount() != 100 {
		t.Errorf("Expected 100 agents, got %d", orch.GetAgentCount())
	}
}
