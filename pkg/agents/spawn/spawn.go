// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package spawn

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/comgunner/picoclaw/pkg/config"
	"github.com/comgunner/picoclaw/pkg/gateway"
	"github.com/comgunner/picoclaw/pkg/routing"
)

type Params struct {
	Task              string
	Label             string
	AgentID           string
	Model             string
	Thinking          string
	RunTimeoutSeconds int
	Thread            bool
	Mode              string
	Cleanup           string
}

type Context struct {
	AgentSessionKey string
	RequesterAgent  string
}

type Result struct {
	Status          string
	RunID           string
	ChildSessionKey string
	Error           string
}

func SpawnSubagentDirect(
	ctx context.Context,
	params Params,
	spawnCtx Context,
	cfg *config.Config,
	registry *Registry,
	gatewayClient gateway.Client,
	allowlistCheck func(parentAgentID, targetAgentID string) bool,
) (Result, error) {
	task := strings.TrimSpace(params.Task)
	if task == "" {
		return Result{Status: "error", Error: "task is required"}, nil
	}
	if cfg == nil {
		return Result{Status: "error", Error: "config is nil"}, nil
	}
	if registry == nil {
		return Result{Status: "error", Error: "registry is nil"}, nil
	}
	if gatewayClient == nil {
		return Result{Status: "error", Error: "gateway client is nil"}, nil
	}

	requesterSession := strings.TrimSpace(spawnCtx.AgentSessionKey)
	if requesterSession == "" {
		requesterSession = "agent:main:main"
	}
	depth := DepthFromSessionKey(requesterSession)

	maxDepth := DefaultSubagentMaxSpawnDepth
	if cfg.Agents.List != nil {
		// Look for requester agent config to get its subagent settings
		requesterAgentID := routing.NormalizeAgentID(spawnCtx.RequesterAgent)
		if requesterAgentID == "" {
			if parsed := routing.ParseAgentSessionKey(requesterSession); parsed != nil {
				requesterAgentID = routing.NormalizeAgentID(parsed.AgentID)
			}
		}

		for _, agentCfg := range cfg.Agents.List {
			if routing.NormalizeAgentID(agentCfg.ID) == requesterAgentID {
				if agentCfg.Subagents != nil && agentCfg.Subagents.MaxSpawnDepth > 0 {
					maxDepth = agentCfg.Subagents.MaxSpawnDepth
				}
				break
			}
		}
	}
	if !AllowedToSpawn(depth, maxDepth) {
		return Result{
			Status: "forbidden",
			Error: fmt.Sprintf(
				"sessions_spawn is not allowed at this depth (current depth: %d, max: %d)",
				depth,
				maxDepth,
			),
		}, nil
	}

	maxChildren := 5
	if cfg.Agents.List != nil {
		// Look for requester agent config to get its subagent settings
		requesterAgentID := routing.NormalizeAgentID(spawnCtx.RequesterAgent)
		if requesterAgentID == "" {
			if parsed := routing.ParseAgentSessionKey(requesterSession); parsed != nil {
				requesterAgentID = routing.NormalizeAgentID(parsed.AgentID)
			}
		}

		for _, agentCfg := range cfg.Agents.List {
			if routing.NormalizeAgentID(agentCfg.ID) == requesterAgentID {
				if agentCfg.Subagents != nil && agentCfg.Subagents.MaxChildrenPerAgent > 0 {
					maxChildren = agentCfg.Subagents.MaxChildrenPerAgent
				}
				break
			}
		}
	}
	active := registry.Count(requesterSession)
	if active >= maxChildren {
		return Result{
			Status: "forbidden",
			Error: fmt.Sprintf(
				"sessions_spawn has reached max active children for this session (%d/%d)",
				active,
				maxChildren,
			),
		}, nil
	}

	requesterAgentID := routing.NormalizeAgentID(spawnCtx.RequesterAgent)
	if requesterAgentID == "" {
		if parsed := routing.ParseAgentSessionKey(requesterSession); parsed != nil {
			requesterAgentID = routing.NormalizeAgentID(parsed.AgentID)
		}
	}
	if requesterAgentID == "" {
		requesterAgentID = routing.DefaultAgentID
	}

	targetAgentID := routing.NormalizeAgentID(params.AgentID)
	if targetAgentID == "" {
		targetAgentID = requesterAgentID
	}

	if targetAgentID != requesterAgentID && allowlistCheck != nil && !allowlistCheck(requesterAgentID, targetAgentID) {
		return Result{
			Status: "forbidden",
			Error:  fmt.Sprintf("agent %s is not allowed to spawn %s", requesterAgentID, targetAgentID),
		}, nil
	}

	childSessionKey := fmt.Sprintf("agent:%s:subagent:%s", targetAgentID, uuid.NewString())
	runTimeoutSeconds := params.RunTimeoutSeconds
	if runTimeoutSeconds < 0 {
		runTimeoutSeconds = 0
	}

	callCtx := ctx
	var cancel context.CancelFunc
	if runTimeoutSeconds > 0 {
		callCtx, cancel = context.WithTimeout(ctx, time.Duration(runTimeoutSeconds)*time.Second)
	} else {
		callCtx, cancel = context.WithTimeout(ctx, 30*time.Second)
	}
	defer cancel()

	resp, err := gatewayClient.Call(callCtx, gateway.Call{
		Method: "agent",
		Params: map[string]any{
			"sessionKey": childSessionKey,
			"message":    task,
			"deliver":    false,
			"lane":       "subagent",
			"label":      strings.TrimSpace(params.Label),
			"model":      strings.TrimSpace(params.Model),
			"thinking":   strings.TrimSpace(params.Thinking),
			"mode":       strings.TrimSpace(params.Mode),
			"cleanup":    strings.TrimSpace(params.Cleanup),
			"thread":     params.Thread,
		},
	})
	if err != nil {
		return Result{Status: "error", Error: err.Error(), ChildSessionKey: childSessionKey}, nil
	}

	runID := strings.TrimSpace(resp.RunID)
	if runID == "" {
		runID = fmt.Sprintf("run-%s", uuid.NewString())
	}

	if err := registry.Register(Run{
		RunID:            runID,
		RequesterSession: requesterSession,
		ChildSessionKey:  childSessionKey,
		AgentID:          targetAgentID,
		ParentDepth:      depth,
		StartedAt:        time.Now(),
	}); err != nil {
		return Result{Status: "error", Error: err.Error(), ChildSessionKey: childSessionKey}, nil
	}

	return Result{
		Status:          "accepted",
		RunID:           runID,
		ChildSessionKey: childSessionKey,
	}, nil
}
