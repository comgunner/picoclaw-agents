// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/comgunner/picoclaw/pkg/agents/spawn"
	"github.com/comgunner/picoclaw/pkg/config"
	"github.com/comgunner/picoclaw/pkg/gateway"
	"github.com/comgunner/picoclaw/pkg/tools"
)

type SessionsSpawnTool struct {
	agentSessionKey string
	agentID         string
	config          *config.Config
	registry        *spawn.Registry
	gatewayClient   gateway.Client
	allowlistCheck  func(parentAgentID, targetAgentID string) bool
}

func NewSessionsSpawnTool(
	cfg *config.Config,
	reg *spawn.Registry,
	client gateway.Client,
	allowlistCheck func(parentAgentID, targetAgentID string) bool,
) *SessionsSpawnTool {
	return &SessionsSpawnTool{
		config:         cfg,
		registry:       reg,
		gatewayClient:  client,
		allowlistCheck: allowlistCheck,
	}
}

func (s *SessionsSpawnTool) Name() string {
	return "sessions_spawn"
}

func (s *SessionsSpawnTool) Description() string {
	return `Spawn a sub-agent run in an isolated session. The sub-agent executes the task asynchronously and announces the result to the requester channel. Non-blocking: returns immediately with {status, runId, childSessionKey}.`
}

func (s *SessionsSpawnTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"task": map[string]any{
				"type":        "string",
				"description": "The task to execute (required)",
			},
			"agentId": map[string]any{
				"type":        "string",
				"description": "Target agent ID (optional; default: requester agent)",
			},
			"model": map[string]any{
				"type":        "string",
				"description": "Model override (optional)",
			},
			"runTimeoutSeconds": map[string]any{
				"type":        "integer",
				"description": "Timeout in seconds (optional)",
			},
			"label": map[string]any{
				"type":        "string",
				"description": "Label for logs/UI (optional)",
			},
		},
		"required": []string{"task"},
	}
}

func (s *SessionsSpawnTool) SetContext(channel, chatID string) {
	// Extract agent session key from channel/chatID if possible
	// For now, we'll use a placeholder approach
	s.agentSessionKey = fmt.Sprintf("%s:%s", channel, chatID)
}

func (s *SessionsSpawnTool) SetAgentID(agentID string) {
	s.agentID = agentID
}

func (s *SessionsSpawnTool) Execute(ctx context.Context, args map[string]any) *tools.ToolResult {
	task, ok := args["task"].(string)
	if !ok || strings.TrimSpace(task) == "" {
		return tools.ErrorResult("task is required").WithError(fmt.Errorf("task parameter is required"))
	}

	params := spawn.Params{
		Task:              task,
		Label:             getStringArg(args, "label"),
		AgentID:           getStringArg(args, "agentId"),
		Model:             getStringArg(args, "model"),
		RunTimeoutSeconds: getIntArg(args, "runTimeoutSeconds"),
	}

	spawnCtx := spawn.Context{
		AgentSessionKey: s.agentSessionKey,
		RequesterAgent:  s.agentID,
	}

	result, err := spawn.SpawnSubagentDirect(
		ctx,
		params,
		spawnCtx,
		s.config,
		s.registry,
		s.gatewayClient,
		s.allowlistCheck,
	)
	if err != nil {
		return tools.ErrorResult(fmt.Sprintf("spawn error: %v", err)).WithError(err)
	}

	// Return JSON result as string
	resultBytes, _ := json.Marshal(map[string]any{
		"status":          result.Status,
		"runId":           result.RunID,
		"childSessionKey": result.ChildSessionKey,
		"error":           result.Error,
	})

	if result.Status == "accepted" {
		return tools.AsyncResult(string(resultBytes))
	} else if result.Status == "forbidden" || result.Status == "error" {
		return tools.ErrorResult(string(resultBytes)).WithError(fmt.Errorf("spawn failed: %s", result.Error))
	}

	return tools.NewToolResult(string(resultBytes))
}

func getStringArg(args map[string]any, key string) string {
	if val, ok := args[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func getIntArg(args map[string]any, key string) int {
	if val, ok := args[key]; ok {
		if num, ok := val.(float64); ok {
			return int(num)
		}
		if num, ok := val.(int); ok {
			return num
		}
	}
	return 0
}
