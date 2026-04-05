// PicoClaw - Ultra-lightweight personal AI agent
// Agent Management - Autonomous Runtime
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/comgunner/picoclaw/pkg/agents/management"
	"github.com/comgunner/picoclaw/pkg/logger"
)

// AgentRuntime handles autonomous execution of tasks for a single agent.
type AgentRuntime struct {
	AgentID      string
	Instance     *AgentInstance
	MessageBus   *management.AgentMessageBus
	InstancesReg *management.InstanceRegistry
	loop         *AgentLoop
}

// RuntimeManager coordinates all AgentRuntimes.
type RuntimeManager struct {
	runtimes []*AgentRuntime
	loop     *AgentLoop
}

// NewRuntimeManager creates runtimes for all registered agents.
func NewRuntimeManager(
	loop *AgentLoop,
	msgBus *management.AgentMessageBus,
	instReg *management.InstanceRegistry,
) *RuntimeManager {
	rm := &RuntimeManager{
		loop:     loop,
		runtimes: make([]*AgentRuntime, 0),
	}

	for _, agentID := range loop.registry.ListAgentIDs() {
		agent, ok := loop.registry.GetAgent(agentID)
		if !ok {
			continue
		}

		// Check if runtime is enabled
		runtimeEnabled := false
		if agent.Runtime != nil {
			runtimeEnabled = agent.Runtime.Enabled
		}

		if !runtimeEnabled {
			continue
		}

		rt := &AgentRuntime{
			AgentID:      agentID,
			Instance:     agent,
			MessageBus:   msgBus,
			InstancesReg: instReg,
			loop:         loop,
		}
		rm.runtimes = append(rm.runtimes, rt)
	}

	return rm
}

// StartAll starts all autonomous agent runtimes in separate goroutines.
func (rm *RuntimeManager) StartAll(ctx context.Context) {
	for _, rt := range rm.runtimes {
		go rt.Start(ctx)
	}
}

// Start runs the infinite loop for this agent's autonomous processing.
func (r *AgentRuntime) Start(ctx context.Context) {
	ch := r.MessageBus.GetChannel(r.AgentID)
	if ch == nil {
		logger.ErrorC("runtime", "No channel found for agent: "+r.AgentID)
		return
	}

	logger.InfoCF("runtime", "Agent runtime started", map[string]any{"agent_id": r.AgentID})

	for {
		select {
		case <-ctx.Done():
			logger.InfoCF("runtime", "Agent runtime stopped", map[string]any{"agent_id": r.AgentID})
			return
		case msg := <-ch:
			r.processMessage(ctx, msg)
		}
	}
}

// processMessage handles a single inbound management message autonomously.
func (r *AgentRuntime) processMessage(ctx context.Context, msg management.AgentMessage) {
	logger.InfoCF("runtime", "Processing autonomous message", map[string]any{
		"agent_id":   r.AgentID,
		"message_id": msg.ID,
		"type":       msg.MessageType,
	})

	// Update instance status
	if inst, ok := r.InstancesReg.Get(r.AgentID); ok {
		inst.UpdateStatus(management.StatusBusy)
		defer inst.UpdateStatus(management.StatusIdle)
	}

	// Construct a synthetic user message for the LLM
	payloadStr := "{}"
	if len(msg.Payload) > 0 {
		payloadStr = string(msg.Payload)
	}

	systemPrompt := fmt.Sprintf(
		"[Internal Task from %s]\nTask Type: %s\nPayload: %s\n\nPlease process this task using your available tools. You are working autonomously. Provide a clear summary of your results.",
		msg.SenderID,
		msg.MessageType,
		payloadStr,
	)

	// Since this is autonomous internal communication, we use a special internal channel and session.
	// Session key separates this from standard chat sessions.
	sessionKey := fmt.Sprintf("auto_%s_%s", r.AgentID, msg.SenderID)

	opts := processOptions{
		SessionKey:      sessionKey,
		Channel:         "system", // Internal channel
		ChatID:          "autonomous",
		UserMessage:     systemPrompt,
		DefaultResponse: "Autonomous task completed.",
		EnableSummary:   false,
		SendResponse:    false, // Do not send to external bus (Telegram/Discord)
	}

	response, err := r.loop.runAgentLoop(ctx, r.Instance, opts)
	if err != nil {
		logger.ErrorCF("runtime", "Error in autonomous processing", map[string]any{
			"agent_id": r.AgentID,
			"error":    err.Error(),
		})
		response = fmt.Sprintf("Error processing task: %v", err)
	}

	// Mark as read in the message bus
	_ = r.MessageBus.MarkAsRead(msg.ID)

	// Auto-respond if required or if it's an explicit task
	if msg.RequiresResponse || strings.Contains(strings.ToLower(response), "error") || msg.MessageType == "task" {
		respPayload := map[string]any{
			"result": response,
			"status": "completed",
		}
		if err != nil {
			respPayload["status"] = "error"
		}

		// Attempt to parse response as JSON just in case the LLM returned structured data
		var jsonRes map[string]any
		if err := json.Unmarshal([]byte(response), &jsonRes); err == nil {
			respPayload["result"] = jsonRes
		}

		respBytes, _ := json.Marshal(respPayload)
		err = r.MessageBus.Send(management.AgentMessage{
			ID:               fmt.Sprintf("RESP_%s_%d", msg.ID, time.Now().UnixNano()),
			SenderID:         r.AgentID,
			RecipientID:      msg.SenderID,
			MessageType:      "result",
			Payload:          respBytes,
			RequiresResponse: false,
			SentAt:           time.Now(),
		})
		if err != nil {
			logger.ErrorCF("runtime", "Failed to send auto-response", map[string]any{
				"agent_id": r.AgentID,
				"error":    err.Error(),
			})
		}
	}
}
