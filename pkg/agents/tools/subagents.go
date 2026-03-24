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

	"github.com/comgunner/picoclaw/pkg/agents/spawn"
	"github.com/comgunner/picoclaw/pkg/gateway"
	"github.com/comgunner/picoclaw/pkg/tools"
)

type SubagentsTool struct {
	agentSessionKey string
	registry        *spawn.Registry
	gatewayClient   gateway.Client
}

func NewSubagentsTool(reg *spawn.Registry, client gateway.Client) *SubagentsTool {
	return &SubagentsTool{
		registry:      reg,
		gatewayClient: client,
	}
}

func (s *SubagentsTool) Name() string {
	return "subagents"
}

func (s *SubagentsTool) Description() string {
	return `List, kill, or steer spawned sub-agents for your session. Use this for sub-agent orchestration.`
}

func (s *SubagentsTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"action": map[string]any{
				"type":        "string",
				"enum":        []string{"list", "kill", "steer"},
				"description": "Action to perform on subagents",
			},
			"target": map[string]any{
				"type":        "string",
				"description": "Target subagent ID for kill/steer actions",
			},
			"message": map[string]any{
				"type":        "string",
				"description": "Message to send when steering",
			},
		},
		"required": []string{"action"},
	}
}

func (s *SubagentsTool) SetContext(channel, chatID string) {
	s.agentSessionKey = fmt.Sprintf("%s:%s", channel, chatID)
}

func (s *SubagentsTool) Execute(ctx context.Context, args map[string]any) *tools.ToolResult {
	action, ok := args["action"].(string)
	if !ok {
		return tools.ErrorResult("action is required").WithError(fmt.Errorf("action parameter is required"))
	}

	switch action {
	case "list":
		return s.handleList()
	case "kill":
		return s.handleKill(args)
	case "steer":
		return s.handleSteer(args)
	default:
		return tools.ErrorResult(fmt.Sprintf("unknown action: %s", action)).
			WithError(fmt.Errorf("unknown action: %s", action))
	}
}

func (s *SubagentsTool) handleList() *tools.ToolResult {
	if s.registry == nil {
		return tools.ErrorResult("registry is nil").WithError(fmt.Errorf("registry is nil"))
	}

	runs := s.registry.Get(s.agentSessionKey)
	result := map[string]any{
		"action": "list",
		"runs":   runs,
		"count":  len(runs),
	}

	resultBytes, _ := json.Marshal(result)
	return tools.NewToolResult(string(resultBytes))
}

func (s *SubagentsTool) handleKill(args map[string]any) *tools.ToolResult {
	if s.registry == nil {
		return tools.ErrorResult("registry is nil").WithError(fmt.Errorf("registry is nil"))
	}

	target, ok := args["target"].(string)
	if !ok || target == "" {
		return tools.ErrorResult("target is required for kill action").
			WithError(fmt.Errorf("target is required for kill action"))
	}

	err := s.registry.Complete(target)
	if err != nil {
		return tools.ErrorResult(fmt.Sprintf("failed to kill subagent %s: %v", target, err)).WithError(err)
	}

	// Also try to kill the run via gateway
	if s.gatewayClient != nil {
		_, gwErr := s.gatewayClient.Call(context.Background(), gateway.Call{
			Method: "kill",
			Params: map[string]any{
				"runId": target,
			},
		})
		if gwErr != nil {
			return tools.ErrorResult(fmt.Sprintf("subagent marked as completed in registry but failed to kill via gateway: %v", gwErr)).
				WithError(gwErr)
		}
	}

	result := map[string]any{
		"action": "kill",
		"target": target,
		"status": "killed",
	}

	resultBytes, _ := json.Marshal(result)
	return tools.NewToolResult(string(resultBytes))
}

func (s *SubagentsTool) handleSteer(args map[string]any) *tools.ToolResult {
	if s.gatewayClient == nil {
		return tools.ErrorResult("gateway client is nil").WithError(fmt.Errorf("gateway client is nil"))
	}

	target, ok := args["target"].(string)
	if !ok || target == "" {
		return tools.ErrorResult("target is required for steer action").
			WithError(fmt.Errorf("target is required for steer action"))
	}

	message, ok := args["message"].(string)
	if !ok || message == "" {
		return tools.ErrorResult("message is required for steer action").
			WithError(fmt.Errorf("message is required for steer action"))
	}

	_, err := s.gatewayClient.Call(context.Background(), gateway.Call{
		Method: "steer",
		Params: map[string]any{
			"target":  target,
			"message": message,
		},
	})
	if err != nil {
		return tools.ErrorResult(fmt.Sprintf("failed to steer subagent %s: %v", target, err)).WithError(err)
	}

	result := map[string]any{
		"action":  "steer",
		"target":  target,
		"message": message,
		"status":  "steered",
	}

	resultBytes, _ := json.Marshal(result)
	return tools.NewToolResult(string(resultBytes))
}
