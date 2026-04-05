// pkg/agent/orchestrator.go
//
// A2A Orchestrator - Agent-to-Agent Coordination System
// Based on concepts from icueth's picoclaw-agents fork (@icueth)
// Adapted for integration with your fork's security model
//
// Original source: https://github.com/icueth/picoclaw-agents
// License: Same as base project (MIT)
//
// Workflow phases: Discovery → Planning → Execution → Integration → Validation

package agent

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/comgunner/picoclaw/pkg/agentcomm"
	"github.com/comgunner/picoclaw/pkg/logger"
	"github.com/comgunner/picoclaw/pkg/mailbox"
)

// A2AOrchestrator coordinates permanent agents with identity
type A2AOrchestrator struct {
	agents        map[string]*AgentInstance // Permanent agents (with IDENTITY.md, SOUL.md)
	mailboxHub    *mailbox.Hub              // Mailbox hub for agent communication
	sharedContext *agentcomm.SharedContext  // Global shared context
	mu            sync.RWMutex
	ctx           context.Context
	cancel        context.CancelFunc
}

// NewA2AOrchestrator creates new orchestrator for agent coordination
func NewA2AOrchestrator(mailboxCapacity int) *A2AOrchestrator {
	ctx, cancel := context.WithCancel(context.Background())

	return &A2AOrchestrator{
		agents:        make(map[string]*AgentInstance),
		mailboxHub:    mailbox.NewHub(mailboxCapacity),
		sharedContext: agentcomm.NewSharedContext(100, 1000),
		ctx:           ctx,
		cancel:        cancel,
	}
}

// RegisterAgent registers a permanent agent with the orchestrator
func (o *A2AOrchestrator) RegisterAgent(agent *AgentInstance) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	agentID := agent.ID
	if _, exists := o.agents[agentID]; exists {
		return fmt.Errorf("agent %s already registered", agentID)
	}

	// Register mailbox
	agent.mailbox = o.mailboxHub.Register(agentID)

	// Store agent
	o.agents[agentID] = agent

	logger.InfoCF("a2a", "Agent registered", map[string]any{
		"agent_id": agentID,
		"role":     agent.Role,
	})

	return nil
}

// UnregisterAgent removes an agent from the orchestrator
func (o *A2AOrchestrator) UnregisterAgent(agentID string) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	if _, exists := o.agents[agentID]; !exists {
		return fmt.Errorf("agent %s not found", agentID)
	}

	// Unregister mailbox
	if err := o.mailboxHub.Unregister(agentID); err != nil {
		return err
	}

	delete(o.agents, agentID)
	logger.InfoCF("a2a", "Agent unregistered", map[string]any{
		"agent_id": agentID,
	})

	return nil
}

// SendMessage sends message from one agent to another
func (o *A2AOrchestrator) SendMessage(
	from, to string,
	msgType mailbox.MessageType,
	priority mailbox.MessagePriority,
	content string,
) error {
	msg := mailbox.Message{
		Type:      msgType,
		From:      from,
		To:        to,
		Priority:  priority,
		Content:   content,
		CreatedAt: time.Now(),
	}

	return o.mailboxHub.SendTo(to, msg)
}

// BroadcastMessage broadcasts to all agents except sender
func (o *A2AOrchestrator) BroadcastMessage(from string, msgType mailbox.MessageType, content string) []string {
	msg := mailbox.Message{
		Type:      msgType,
		From:      from,
		Content:   content,
		CreatedAt: time.Now(),
	}

	return o.mailboxHub.Broadcast(msg)
}

// DiscoverCapabilities broadcasts agent capabilities
func (o *A2AOrchestrator) DiscoverCapabilities(agentID string) error {
	o.mu.RLock()
	agent, ok := o.agents[agentID]
	o.mu.RUnlock()

	if !ok {
		return fmt.Errorf("agent %s not found", agentID)
	}

	// Store capabilities in shared context
	key := fmt.Sprintf("agent:%s:capabilities", agentID)
	o.sharedContext.Set(key, agent.Tools)

	// Broadcast discovery message
	msg := mailbox.Message{
		Type:    mailbox.MessageTypeBroadcast,
		From:    agentID,
		Content: fmt.Sprintf("Agent %s (%s) is ready", agentID, agent.Role),
	}
	o.mailboxHub.Broadcast(msg)

	logger.InfoCF("a2a", "Agent capabilities discovered", map[string]any{
		"agent_id": agentID,
		"role":     agent.Role,
	})

	return nil
}

// AssignTask assigns task to agent with priority
func (o *A2AOrchestrator) AssignTask(fromID, toID string, task string, priority mailbox.MessagePriority) error {
	logger.InfoCF("a2a", "Task assignment", map[string]any{
		"from":     fromID,
		"to":       toID,
		"task":     truncate(task, 50),
		"priority": priority,
	})

	// Store in shared context
	taskKey := fmt.Sprintf("task:%s:%d", toID, time.Now().UnixNano())
	o.sharedContext.Set(taskKey, map[string]any{
		"from":    fromID,
		"content": task,
		"status":  "assigned",
		"time":    time.Now().Unix(),
	})

	// Send via mailbox
	msg := mailbox.Message{
		Type:     mailbox.MessageTypeTask,
		From:     fromID,
		To:       toID,
		Priority: priority,
		Content:  task,
	}

	return o.mailboxHub.SendTo(toID, msg)
}

// ReportTaskComplete reports task completion
func (o *A2AOrchestrator) ReportTaskComplete(agentID string, taskKey string, result string) error {
	// Update in shared context
	resultKey := fmt.Sprintf("%s:result", taskKey)
	o.sharedContext.Set(resultKey, map[string]any{
		"status": "completed",
		"result": result,
		"time":   time.Now().Unix(),
	})

	// Log message
	o.sharedContext.AddMessageLog(agentID, "orchestrator", "response", result)

	logger.InfoCF("a2a", "Task completed", map[string]any{
		"agent_id": agentID,
		"task_key": taskKey,
	})

	return nil
}

// GetAgentStatus returns agent status
func (o *A2AOrchestrator) GetAgentStatus(agentID string) map[string]any {
	o.mu.RLock()
	agent, ok := o.agents[agentID]
	o.mu.RUnlock()

	if !ok {
		return map[string]any{"error": "agent not found"}
	}

	mb, _ := o.mailboxHub.Get(agentID)

	return map[string]any{
		"id":           agentID,
		"role":         agent.Role,
		"status":       agent.Status,
		"mailbox_size": mb.Size(),
		"unread_count": mb.GetUnreadCount(),
		"last_active":  agent.LastActive,
	}
}

// GetOrchestrationStats returns full orchestration statistics
func (o *A2AOrchestrator) GetOrchestrationStats() map[string]any {
	o.mu.RLock()
	agentCount := len(o.agents)
	o.mu.RUnlock()

	return map[string]any{
		"agents":        agentCount,
		"mailbox_stats": o.mailboxHub.GetStats(),
		"context_size":  o.sharedContext.ContextSize(),
		"message_log":   o.sharedContext.LogSize(),
		"timestamp":     time.Now().Unix(),
	}
}

// GetSharedContext returns the shared context for direct access
func (o *A2AOrchestrator) GetSharedContext() *agentcomm.SharedContext {
	return o.sharedContext
}

// GetMailboxHub returns the mailbox hub for direct access
func (o *A2AOrchestrator) GetMailboxHub() *mailbox.Hub {
	return o.mailboxHub
}

// StartCleanup starts periodic cleanup
func (o *A2AOrchestrator) StartCleanup(interval time.Duration) {
	o.mailboxHub.StartCleanup(o.ctx, interval)

	logger.InfoCF("a2a", "A2A cleanup started", map[string]any{
		"interval": interval.String(),
	})
}

// Stop stops the orchestrator
func (o *A2AOrchestrator) Stop() error {
	o.cancel()
	logger.InfoCF("a2a", "A2A orchestrator stopped", nil)
	return nil
}

// WorkflowDiscovery broadcasts all agent capabilities
func (o *A2AOrchestrator) WorkflowDiscovery() error {
	logger.InfoCF("a2a", "Starting workflow: Discovery phase", nil)

	// Phase 1: Discovery - agents broadcast capabilities
	o.mu.RLock()
	agents := make([]*AgentInstance, 0, len(o.agents))
	for _, agent := range o.agents {
		agents = append(agents, agent)
	}
	o.mu.RUnlock()

	for _, agent := range agents {
		if err := o.DiscoverCapabilities(agent.ID); err != nil {
			logger.ErrorCF("a2a", "Discovery failed", map[string]any{
				"agent": agent.ID,
				"error": err.Error(),
			})
		}
	}

	return nil
}

// WorkflowPlanning assigns tasks to agents
func (o *A2AOrchestrator) WorkflowPlanning(coordinator string, taskAssignments map[string]string) error {
	logger.InfoCF("a2a", "Starting workflow: Planning phase", map[string]any{
		"coordinator": coordinator,
		"tasks":       len(taskAssignments),
	})

	for agentID, task := range taskAssignments {
		if err := o.AssignTask(coordinator, agentID, task, mailbox.PriorityHigh); err != nil {
			logger.ErrorCF("a2a", "Task assignment failed", map[string]any{
				"agent": agentID,
				"error": err.Error(),
			})
		}
	}

	return nil
}

// WorkflowExecution monitors task execution (placeholder for future implementation)
func (o *A2AOrchestrator) WorkflowExecution() error {
	logger.InfoCF("a2a", "Starting workflow: Execution phase", nil)
	// Future: Monitor task progress via shared context
	return nil
}

// WorkflowIntegration collects results from all agents (placeholder)
func (o *A2AOrchestrator) WorkflowIntegration() error {
	logger.InfoCF("a2a", "Starting workflow: Integration phase", nil)
	// Future: Collect and integrate results from all agents
	return nil
}

// WorkflowValidation validates final results (placeholder)
func (o *A2AOrchestrator) WorkflowValidation() error {
	logger.InfoCF("a2a", "Starting workflow: Validation phase", nil)
	// Future: Validate integrated results
	return nil
}

// RunFullWorkflow executes complete workflow: Discovery → Planning → Execution → Integration → Validation
func (o *A2AOrchestrator) RunFullWorkflow(coordinator string, taskAssignments map[string]string) error {
	logger.InfoCF("a2a", "Running full A2A workflow", map[string]any{
		"coordinator": coordinator,
		"tasks":       len(taskAssignments),
	})

	// Phase 1: Discovery
	if err := o.WorkflowDiscovery(); err != nil {
		return fmt.Errorf("discovery failed: %w", err)
	}

	// Phase 2: Planning
	if err := o.WorkflowPlanning(coordinator, taskAssignments); err != nil {
		return fmt.Errorf("planning failed: %w", err)
	}

	// Phase 3: Execution
	if err := o.WorkflowExecution(); err != nil {
		return fmt.Errorf("execution failed: %w", err)
	}

	// Phase 4: Integration
	if err := o.WorkflowIntegration(); err != nil {
		return fmt.Errorf("integration failed: %w", err)
	}

	// Phase 5: Validation
	if err := o.WorkflowValidation(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	logger.InfoCF("a2a", "Full A2A workflow completed", nil)
	return nil
}

// ListAgents returns all registered agent IDs
func (o *A2AOrchestrator) ListAgents() []string {
	o.mu.RLock()
	defer o.mu.RUnlock()

	agents := make([]string, 0, len(o.agents))
	for id := range o.agents {
		agents = append(agents, id)
	}
	return agents
}

// GetAgentCount returns the number of registered agents
func (o *A2AOrchestrator) GetAgentCount() int {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return len(o.agents)
}

// Helper function to truncate strings for logging
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
