// PicoClaw - Ultra-lightweight personal AI agent
// Agent Management - Suite Initializer
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

// Package managementinit wires the management.Suite with the management tools
// and registers all 12 Management Skill tools into a ToolRegistry.
//
// It lives in a separate package to avoid the import cycle that would arise if
// the management package imported its own tools sub-package.
//
// Typical usage in pkg/agent/loop.go:
//
//	managementinit.RegisterManagementTools(cfg, agent.Tools)
package managementinit

import (
	"github.com/comgunner/picoclaw/pkg/agents/management"
	managementtools "github.com/comgunner/picoclaw/pkg/agents/management/tools"
	"github.com/comgunner/picoclaw/pkg/config"
	"github.com/comgunner/picoclaw/pkg/tools"
)

// Suite holds all shared state objects that back the Management Skill tools.
// Callers can keep a reference to Suite to access registry, bus, etc. at runtime.
type Suite struct {
	Registry    *management.AgentRegistry
	Instances   *management.InstanceRegistry
	MessageBus  *management.AgentMessageBus
	RateLimiter *management.RateLimiter
}

// NewSuite creates a fully-initialized Suite from the application config.
func NewSuite(cfg *config.Config) *Suite {
	return &Suite{
		Registry:    management.NewAgentRegistry(cfg),
		Instances:   management.NewInstanceRegistry(),
		MessageBus:  management.NewAgentMessageBus(),
		RateLimiter: management.NewRateLimiter(),
	}
}

// RegisterAllTools registers all 12 management tools into toolRegistry.
func (s *Suite) RegisterAllTools(toolRegistry *tools.ToolRegistry) {
	// Phase 1 – Basic Management
	toolRegistry.Register(managementtools.NewAgentGetTool(s.Registry, s.Instances))
	toolRegistry.Register(managementtools.NewAgentListTool(s.Registry, s.Instances))
	toolRegistry.Register(managementtools.NewAgentCanSpawnTool(s.Registry))
	toolRegistry.Register(managementtools.NewAgentDefaultTool(s.Registry, s.Instances))

	// Phase 2 – Activity Monitoring
	toolRegistry.Register(managementtools.NewAgentMonitorTool(s.Registry, s.Instances))
	toolRegistry.Register(managementtools.NewAgentSpawnTool(s.Registry, s.Instances, s.RateLimiter))

	// Phase 3 – Inter-Agent Communication
	toolRegistry.Register(managementtools.NewAgentSendTool(s.Registry, s.MessageBus, s.RateLimiter))
	toolRegistry.Register(managementtools.NewAgentReceiveTool(s.Registry, s.MessageBus))
	toolRegistry.Register(managementtools.NewAgentRespondTool(s.Registry, s.MessageBus))
	toolRegistry.Register(managementtools.NewAgentBroadcastTool(s.Registry, s.MessageBus, s.RateLimiter))

	// Phase 4 – Debugging
	toolRegistry.Register(managementtools.NewAgentMessageLogTool(s.MessageBus))
	toolRegistry.Register(managementtools.NewAgentMessageStatsTool(s.MessageBus))
}

// RegisterManagementTools is a one-call convenience that creates a Suite and
// immediately registers all 12 tools.  It returns the Suite for callers that
// need runtime access to the shared state (e.g. to monitor agents externally).
func RegisterManagementTools(cfg *config.Config, toolRegistry *tools.ToolRegistry) *Suite {
	s := NewSuite(cfg)
	s.RegisterAllTools(toolRegistry)
	return s
}
