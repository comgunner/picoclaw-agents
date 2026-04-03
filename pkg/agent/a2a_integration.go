// pkg/agent/a2a_integration.go
//
// A2A Integration - Agent-to-Agent Communication Integration
// Based on concepts from icueth's picoclaw-agents fork (@icueth)
//
// This file provides integration helpers for A2A communication
// without modifying the core loop.go extensively
//
// Credits: @icueth (https://github.com/icueth)
// License: Same as base project (MIT)

package agent

import (
	"fmt"
	"time"

	"github.com/comgunner/picoclaw/pkg/agentcomm"
	"github.com/comgunner/picoclaw/pkg/logger"
	"github.com/comgunner/picoclaw/pkg/mailbox"
)

// A2AIntegration provides A2A communication capabilities for an agent
type A2AIntegration struct {
	agentInstance    *AgentInstance
	mailbox          *mailbox.Mailbox
	hub              *mailbox.Hub // B-01/B-02 fix: hub for routing to other agents
	sharedContext    *agentcomm.SharedContext
	departmentRouter *DepartmentRouter
}

// NewA2AIntegration creates A2A integration for an agent
func NewA2AIntegration(
	agent *AgentInstance,
	mb *mailbox.Mailbox,
	hub *mailbox.Hub,
	ctx *agentcomm.SharedContext,
	router *DepartmentRouter,
) *A2AIntegration {
	return &A2AIntegration{
		agentInstance:    agent,
		mailbox:          mb,
		hub:              hub,
		sharedContext:    ctx,
		departmentRouter: router,
	}
}

// ProcessMailboxMessages processes pending messages in the mailbox
func (a *A2AIntegration) ProcessMailboxMessages() error {
	if a.mailbox == nil {
		return nil // No mailbox configured
	}

	for {
		msg, err := a.mailbox.Receive()
		if err != nil {
			// No more messages
			break
		}

		if handleErr := a.handleMessage(msg); handleErr != nil {
			logger.ErrorCF("a2a", "Failed to handle A2A message", map[string]any{
				"from":    msg.From,
				"type":    msg.Type,
				"error":   handleErr.Error(),
				"content": truncateString(msg.Content, 50),
			})
			return fmt.Errorf("failed to handle message from %s: %w", msg.From, handleErr)
		}
	}

	return nil //nolint:nilerr // All errors in loop are handled and returned immediately
}

// handleMessage processes a single A2A message
func (a *A2AIntegration) handleMessage(msg mailbox.Message) error {
	logger.InfoCF("a2a", "Processing A2A message", map[string]any{
		"from":    msg.From,
		"type":    msg.Type,
		"content": truncateString(msg.Content, 50),
	})

	// Update last active timestamp
	a.agentInstance.LastActive = time.Now().Unix()

	switch msg.Type {
	case mailbox.MessageTypeTask:
		return a.handleTask(msg)
	case mailbox.MessageTypeQuestion:
		return a.handleQuestion(msg)
	case mailbox.MessageTypeStatus:
		return a.handleStatus(msg)
	case mailbox.MessageTypeBroadcast:
		return a.handleBroadcast(msg)
	default:
		return fmt.Errorf("unknown message type: %s", msg.Type)
	}
}

// handleTask processes a task message
func (a *A2AIntegration) handleTask(msg mailbox.Message) error {
	logger.InfoCF("a2a", "Task assigned", map[string]any{
		"from": msg.From,
		"task": truncateString(msg.Content, 100),
	})

	// Store task in shared context
	taskKey := fmt.Sprintf("task:%s:%d", a.agentInstance.ID, time.Now().UnixNano())
	a.sharedContext.Set(taskKey, map[string]any{
		"from":      msg.From,
		"content":   msg.Content,
		"status":    "received",
		"timestamp": time.Now().Unix(),
		"priority":  msg.Priority,
	})

	// Update agent status
	a.agentInstance.Status = "running"

	return nil
}

// handleQuestion processes a question message
func (a *A2AIntegration) handleQuestion(msg mailbox.Message) error {
	logger.InfoCF("a2a", "Question received", map[string]any{
		"from":     msg.From,
		"question": truncateString(msg.Content, 100),
		"priority": msg.Priority,
	})

	// Store question in shared context
	questionKey := fmt.Sprintf("question:%s:%d", a.agentInstance.ID, time.Now().UnixNano())
	a.sharedContext.Set(questionKey, map[string]any{
		"from":      msg.From,
		"question":  msg.Content,
		"status":    "pending",
		"timestamp": time.Now().Unix(),
	})

	return nil
}

// handleStatus processes a status message
func (a *A2AIntegration) handleStatus(msg mailbox.Message) error {
	logger.InfoCF("a2a", "Status update received", map[string]any{
		"from":   msg.From,
		"status": truncateString(msg.Content, 100),
	})

	// Store status in shared context
	statusKey := fmt.Sprintf("status:%s:%d", msg.From, time.Now().UnixNano())
	a.sharedContext.Set(statusKey, map[string]any{
		"content":   msg.Content,
		"timestamp": time.Now().Unix(),
	})

	return nil
}

// handleBroadcast processes a broadcast message
func (a *A2AIntegration) handleBroadcast(msg mailbox.Message) error {
	logger.InfoCF("a2a", "Broadcast received", map[string]any{
		"from":      msg.From,
		"broadcast": truncateString(msg.Content, 100),
	})

	// Store broadcast in shared context
	broadcastKey := fmt.Sprintf("broadcast:%d", time.Now().UnixNano())
	a.sharedContext.Set(broadcastKey, map[string]any{
		"from":      msg.From,
		"content":   msg.Content,
		"timestamp": time.Now().Unix(),
	})

	return nil
}

// SendTask sends a task to another agent via the hub (B-01 fix)
func (a *A2AIntegration) SendTask(toAgentID, task string, priority mailbox.MessagePriority) error {
	if a.hub == nil {
		return fmt.Errorf("hub not configured")
	}

	msg := mailbox.Message{
		Type:      mailbox.MessageTypeTask,
		From:      a.agentInstance.ID,
		To:        toAgentID,
		Priority:  priority,
		Content:   task,
		CreatedAt: time.Now(),
	}

	// B-01 fix: route via hub to recipient's mailbox, not sender's
	if err := a.hub.SendTo(toAgentID, msg); err != nil {
		return fmt.Errorf("failed to send task: %w", err)
	}

	logger.InfoCF("a2a", "Task sent", map[string]any{
		"to":   toAgentID,
		"task": truncateString(task, 50),
	})

	return nil
}

// SendStatus sends a status update to another agent via the hub (B-01 fix)
func (a *A2AIntegration) SendStatus(toAgentID, status string) error {
	if a.hub == nil {
		return fmt.Errorf("hub not configured")
	}

	msg := mailbox.Message{
		Type:      mailbox.MessageTypeStatus,
		From:      a.agentInstance.ID,
		To:        toAgentID,
		Priority:  mailbox.PriorityNormal,
		Content:   status,
		CreatedAt: time.Now(),
	}

	// B-01 fix: route via hub to recipient's mailbox
	if err := a.hub.SendTo(toAgentID, msg); err != nil {
		return fmt.Errorf("failed to send status: %w", err)
	}

	return nil
}

// BroadcastStatus broadcasts a status update to all agents via the hub (B-02 fix)
func (a *A2AIntegration) BroadcastStatus(status string) error {
	if a.hub == nil {
		return fmt.Errorf("hub not configured")
	}

	msg := mailbox.Message{
		Type:      mailbox.MessageTypeBroadcast,
		From:      a.agentInstance.ID,
		Content:   status,
		CreatedAt: time.Now(),
	}

	// B-02 fix: actually broadcast via hub (was no-op before)
	failed := a.hub.Broadcast(msg)
	if len(failed) > 0 {
		logger.WarnCF("a2a", "Broadcast partially failed", map[string]any{
			"failed_count": len(failed),
			"failed":       failed,
		})
	}

	logger.InfoCF("a2a", "Status broadcast sent", map[string]any{
		"status":       truncateString(status, 50),
		"failed_count": len(failed),
	})

	return nil
}

// GetPendingTasks returns all pending tasks from shared context
func (a *A2AIntegration) GetPendingTasks() []map[string]any {
	tasks := make([]map[string]any, 0)

	keys := a.sharedContext.Keys()
	for _, key := range keys {
		if len(key) > 5 && key[:5] == "task:" {
			if val, ok := a.sharedContext.Get(key); ok {
				if taskMap, ok := val.(map[string]any); ok {
					if taskMap["status"] == "received" || taskMap["status"] == "assigned" {
						tasks = append(tasks, taskMap)
					}
				}
			}
		}
	}

	return tasks
}

// ReportTaskComplete reports a task as completed
func (a *A2AIntegration) ReportTaskComplete(taskKey string, result string) error {
	resultKey := fmt.Sprintf("%s:result", taskKey)
	a.sharedContext.Set(resultKey, map[string]any{
		"status":    "completed",
		"result":    result,
		"timestamp": time.Now().Unix(),
	})

	// Update original task status
	if task, ok := a.sharedContext.Get(taskKey); ok {
		if taskMap, ok := task.(map[string]any); ok {
			taskMap["status"] = "completed"
			a.sharedContext.Set(taskKey, taskMap)
		}
	}

	logger.InfoCF("a2a", "Task completed", map[string]any{
		"task_key": taskKey,
		"result":   truncateString(result, 100),
	})

	// Update agent status
	a.agentInstance.Status = "idle"

	return nil
}

// GetModelForAgent returns the appropriate model name for this agent based on department
func (a *A2AIntegration) GetModelForAgent() string {
	if a.departmentRouter == nil {
		return a.agentInstance.Model
	}

	// B-10 fix: GetModelConfig now returns DepartmentModelConfig struct
	return a.departmentRouter.GetModelConfig(a.agentInstance.ID).Model
}

// Helper function to truncate strings for logging
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
