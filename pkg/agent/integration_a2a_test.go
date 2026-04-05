// pkg/agent/integration_a2a_test.go
//
// Integration tests for A2A Orchestrator
// Based on concepts from icueth's picoclaw-agents fork (@icueth)
//
// Credits: @icueth (https://github.com/icueth)
// License: Same as base project (MIT)

package agent

import (
	"testing"
	"time"

	"github.com/comgunner/picoclaw/pkg/config"
	"github.com/comgunner/picoclaw/pkg/mailbox"
)

func ptrFloat64(v float64) *float64 { return &v }

func createIntegrationTestConfig() *config.Config {
	return &config.Config{
		Agents: config.AgentsConfig{
			Defaults: config.AgentDefaults{
				Provider:    "openai",
				Model:       "gpt-4",
				Temperature: ptrFloat64(0.7),
				TopP:        ptrFloat64(0.9),
				MaxTokens:   4096,
			},
			DepartmentModels: map[string]map[string]any{
				"engineering": {
					"provider":        "openai",
					"model":           "gpt-4",
					"temperature":     0.1,
					"top_p":           0.5,
					"enable_thinking": true,
				},
				"testing": {
					"provider":        "openai",
					"model":           "gpt-4",
					"temperature":     0.1,
					"top_p":           0.5,
					"enable_thinking": false,
				},
			},
		},
	}
}

func TestA2A_FullWorkflow(t *testing.T) {
	// Setup
	cfg := createIntegrationTestConfig()

	// Create orchestrator
	orch := NewA2AOrchestrator(100)
	deptRouter := NewDepartmentRouter(cfg)

	// Register agents
	agents := []struct {
		id   string
		role string
		dept string
	}{
		{"pm", "coordinator", "core"},
		{"dev", "developer", "engineering"},
		{"qa", "tester", "testing"},
	}

	for _, a := range agents {
		agent := &AgentInstance{
			ID:   a.id,
			Role: a.role,
		}
		if err := orch.RegisterAgent(agent); err != nil {
			t.Fatalf("Failed to register agent %s: %v", a.id, err)
		}
		deptRouter.RegisterAgent(a.id, a.dept)
	}

	// Start orchestrator cleanup
	go orch.StartCleanup(10 * time.Second)

	// Test Discovery Phase
	if err := orch.WorkflowDiscovery(); err != nil {
		t.Fatalf("Discovery failed: %v", err)
	}

	// Test Planning Phase
	tasks := map[string]string{
		"dev": "Implement feature X",
		"qa":  "Test feature X",
	}

	if err := orch.WorkflowPlanning("pm", tasks); err != nil {
		t.Fatalf("Planning failed: %v", err)
	}

	// Verify messages were sent
	devMb, _ := orch.mailboxHub.Get("dev")
	if devMb.GetUnreadCount() == 0 {
		t.Errorf("Dev should have received tasks")
	}

	qaMb, _ := orch.mailboxHub.Get("qa")
	if qaMb.GetUnreadCount() == 0 {
		t.Errorf("QA should have received tasks")
	}

	// Test model routing
	devModel := deptRouter.GetModelConfig("dev")
	if devModel.Model != "gpt-4" {
		t.Errorf("Dev should have gpt-4 model")
	}

	// Test stats
	stats := orch.GetOrchestrationStats()
	if stats["agents"].(int) != 3 {
		t.Errorf("Should have 3 agents")
	}
}

func TestA2A_MessagePriority(t *testing.T) {
	orch := NewA2AOrchestrator(100)

	orch.RegisterAgent(&AgentInstance{ID: "dev", Role: "developer"})

	// Send multiple messages with different priorities
	orch.SendMessage("pm", "dev", mailbox.MessageTypeStatus, mailbox.PriorityNormal, "Normal priority")
	orch.SendMessage("pm", "dev", mailbox.MessageTypeTask, mailbox.PriorityCritical, "Critical task")
	orch.SendMessage("pm", "dev", mailbox.MessageTypeStatus, mailbox.PriorityLow, "Low priority")

	// Receive should get critical first
	devMb, _ := orch.mailboxHub.Get("dev")

	first, _ := devMb.Receive()
	if first.Priority != mailbox.PriorityCritical {
		t.Errorf("Critical message should be received first")
	}

	second, _ := devMb.Receive()
	if second.Priority != mailbox.PriorityNormal {
		t.Errorf("Normal message should be received second")
	}

	third, _ := devMb.Receive()
	if third.Priority != mailbox.PriorityLow {
		t.Errorf("Low message should be received third")
	}
}

func TestA2A_TaskCompletion(t *testing.T) {
	orch := NewA2AOrchestrator(100)
	orch.RegisterAgent(&AgentInstance{ID: "dev", Role: "developer"})

	// Assign task
	task := "Implement feature"
	if err := orch.AssignTask("pm", "dev", task, mailbox.PriorityHigh); err != nil {
		t.Fatalf("AssignTask failed: %v", err)
	}

	// Find task key in shared context
	var taskKey string
	keys := orch.sharedContext.Keys()
	for _, key := range keys {
		if len(key) > 5 && key[:5] == "task:" {
			taskKey = key
			break
		}
	}

	if taskKey == "" {
		t.Fatalf("Task key not found in shared context")
	}

	// Report completion
	result := "Feature implemented successfully"
	if err := orch.ReportTaskComplete("dev", taskKey, result); err != nil {
		t.Fatalf("ReportTaskComplete failed: %v", err)
	}

	// Verify result stored
	resultKey := taskKey + ":result"
	val, ok := orch.sharedContext.Get(resultKey)
	if !ok {
		t.Errorf("Result not stored in shared context")
	}

	resultMap, ok := val.(map[string]any)
	if !ok {
		t.Errorf("Result is not a map")
	}

	if resultMap["status"] != "completed" {
		t.Errorf("Expected status 'completed'")
	}

	if resultMap["result"] != result {
		t.Errorf("Expected result '%s', got '%v'", result, resultMap["result"])
	}
}

func TestA2A_BroadcastWithDepartmentRouting(t *testing.T) {
	cfg := createIntegrationTestConfig()

	orch := NewA2AOrchestrator(100)
	deptRouter := NewDepartmentRouter(cfg)

	// Register agents in different departments
	orch.RegisterAgent(&AgentInstance{ID: "dev1", Role: "developer"})
	orch.RegisterAgent(&AgentInstance{ID: "dev2", Role: "developer"})
	orch.RegisterAgent(&AgentInstance{ID: "qa1", Role: "tester"})

	deptRouter.RegisterAgent("dev1", "engineering")
	deptRouter.RegisterAgent("dev2", "engineering")
	deptRouter.RegisterAgent("qa1", "testing")

	// Broadcast from dev1
	msg := "Team meeting at 3pm"
	failed := orch.BroadcastMessage("dev1", mailbox.MessageTypeBroadcast, msg)

	if len(failed) != 0 {
		t.Errorf("Broadcast should not fail, failed: %v", failed)
	}

	// Verify dev2 and qa1 received, but not dev1 (sender)
	dev2Mb, _ := orch.mailboxHub.Get("dev2")
	qa1Mb, _ := orch.mailboxHub.Get("qa1")
	dev1Mb, _ := orch.mailboxHub.Get("dev1")

	if dev2Mb.GetUnreadCount() != 1 {
		t.Errorf("dev2 should receive broadcast")
	}

	if qa1Mb.GetUnreadCount() != 1 {
		t.Errorf("qa1 should receive broadcast")
	}

	if dev1Mb.GetUnreadCount() != 0 {
		t.Errorf("dev1 (sender) should not receive broadcast")
	}
}

func TestA2A_AgentStatus(t *testing.T) {
	orch := NewA2AOrchestrator(100)

	agent := &AgentInstance{
		ID:     "dev",
		Role:   "developer",
		Status: "idle",
	}
	orch.RegisterAgent(agent)

	// Get status
	status := orch.GetAgentStatus("dev")

	if status["id"] != "dev" {
		t.Errorf("Expected id 'dev'")
	}

	if status["role"] != "developer" {
		t.Errorf("Expected role 'developer'")
	}

	if status["status"] != "idle" {
		t.Errorf("Expected status 'idle'")
	}

	// Update status
	agent.Status = "running"
	agent.LastActive = time.Now().Unix()

	status = orch.GetAgentStatus("dev")
	if status["status"] != "running" {
		t.Errorf("Expected status 'running'")
	}
}

func TestA2A_ConcurrentOperations(t *testing.T) {
	orch := NewA2AOrchestrator(100)

	// Register multiple agents
	for i := 0; i < 10; i++ {
		agentID := "agent-" + string(rune(i))
		orch.RegisterAgent(&AgentInstance{ID: agentID, Role: "worker"})
	}

	// Concurrent operations
	done := make(chan bool)

	// Send messages concurrently
	for i := 0; i < 10; i++ {
		go func(idx int) {
			from := "agent-" + string(rune(idx))
			to := "agent-" + string(rune((idx+1)%10))
			orch.SendMessage(from, to, mailbox.MessageTypeTask, mailbox.PriorityNormal, "Task")
			done <- true
		}(i)
	}

	// Wait for all sends
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify all agents have messages
	for i := 0; i < 10; i++ {
		agentID := "agent-" + string(rune(i))
		mb, _ := orch.mailboxHub.Get(agentID)
		if mb.Size() == 0 {
			t.Errorf("Agent %s should have messages", agentID)
		}
	}
}

func TestA2A_IntegrationWithDepartmentModels(t *testing.T) {
	cfg := createIntegrationTestConfig()

	orch := NewA2AOrchestrator(100)
	deptRouter := NewDepartmentRouter(cfg)

	// Register agents
	orch.RegisterAgent(&AgentInstance{ID: "eng_lead", Role: "lead"})
	orch.RegisterAgent(&AgentInstance{ID: "senior_dev", Role: "developer"})
	orch.RegisterAgent(&AgentInstance{ID: "qa_lead", Role: "lead"})

	// Assign to departments
	deptRouter.RegisterAgent("eng_lead", "engineering")
	deptRouter.RegisterAgent("senior_dev", "engineering")
	deptRouter.RegisterAgent("qa_lead", "testing")

	// Verify model routing
	engModel := deptRouter.GetModelConfig("eng_lead")
	if engModel.Temperature != 0.1 {
		t.Errorf("Engineering should have temperature 0.1")
	}
	if !engModel.EnableThinking {
		t.Errorf("Engineering should have thinking enabled")
	}

	qaModel := deptRouter.GetModelConfig("qa_lead")
	if qaModel.Temperature != 0.1 {
		t.Errorf("Testing should have temperature 0.1")
	}
	if qaModel.EnableThinking {
		t.Errorf("Testing should not have thinking enabled")
	}

	// Unassigned agent should use default
	unknownModel := deptRouter.GetModelConfig("unknown_agent")
	if unknownModel.Temperature != 0.7 {
		t.Errorf("Unknown agent should use default temperature 0.7")
	}
}

func TestA2A_WorkflowRunFullWorkflow(t *testing.T) {
	orch := NewA2AOrchestrator(100)

	// Register agents
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

	// Verify all phases completed
	stats := orch.GetOrchestrationStats()
	if stats["agents"].(int) != 3 {
		t.Errorf("Expected 3 agents")
	}

	// Verify tasks were assigned
	devMb, _ := orch.mailboxHub.Get("dev")
	if devMb.GetUnreadCount() == 0 {
		t.Errorf("Dev should have tasks")
	}

	qaMb, _ := orch.mailboxHub.Get("qa")
	if qaMb.GetUnreadCount() == 0 {
		t.Errorf("QA should have tasks")
	}
}

func TestA2A_MailboxCleanup(t *testing.T) {
	orch := NewA2AOrchestrator(100)
	orch.RegisterAgent(&AgentInstance{ID: "dev", Role: "developer"})

	// Send message that expires immediately
	msg := mailbox.Message{
		Type:      mailbox.MessageTypeStatus,
		From:      "pm",
		To:        "dev",
		Content:   "Expired message",
		ExpiresAt: time.Now().Add(-1 * time.Second),
	}
	orch.mailboxHub.SendTo("dev", msg)

	// Start cleanup
	orch.StartCleanup(100 * time.Millisecond)

	// Wait for cleanup
	time.Sleep(200 * time.Millisecond)

	// Verify message was cleaned
	devMb, _ := orch.mailboxHub.Get("dev")
	if devMb.Size() != 0 {
		t.Errorf("Expected 0 messages after cleanup, got %d", devMb.Size())
	}
}

func TestA2A_SharedContextIntegration(t *testing.T) {
	orch := NewA2AOrchestrator(100)
	orch.RegisterAgent(&AgentInstance{ID: "dev", Role: "developer"})

	// Store data in shared context
	orch.sharedContext.Set("shared_data", "test_value")
	orch.sharedContext.AddMessageLog("pm", "dev", "task", "Do something")

	// Verify data
	val, ok := orch.sharedContext.Get("shared_data")
	if !ok || val != "test_value" {
		t.Errorf("Shared context data not stored correctly")
	}

	logs := orch.sharedContext.GetMessageLog()
	if len(logs) != 1 {
		t.Errorf("Expected 1 log entry, got %d", len(logs))
	}

	if logs[0].From != "pm" || logs[0].To != "dev" {
		t.Errorf("Log entry incorrect")
	}
}
