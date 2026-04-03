// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package agent

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/comgunner/picoclaw/pkg/config"
)

func TestNewAgentInstance_UsesDefaultsTemperatureAndMaxTokens(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "agent-instance-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	cfg := &config.Config{
		Agents: config.AgentsConfig{
			Defaults: config.AgentDefaults{
				Workspace:         tmpDir,
				Model:             "test-model",
				MaxTokens:         1234,
				MaxToolIterations: 5,
			},
		},
	}

	configuredTemp := 1.0
	cfg.Agents.Defaults.Temperature = &configuredTemp

	agent := NewAgentInstance(nil, &cfg.Agents.Defaults, cfg, testFactory)

	if agent.MaxTokens != 1234 {
		t.Fatalf("MaxTokens = %d, want %d", agent.MaxTokens, 1234)
	}
	if agent.Temperature != 1.0 {
		t.Fatalf("Temperature = %f, want %f", agent.Temperature, 1.0)
	}
}

func TestResolveAgentWorkspace_NonDefaultDerivedFromConfiguredDefaultWorkspace(t *testing.T) {
	defaults := &config.AgentDefaults{
		Workspace: "/opt/picoclaw/workspace",
	}
	agentCfg := &config.AgentConfig{
		ID: "qa_specialist",
	}

	got := resolveAgentWorkspace(agentCfg, defaults)
	want := "/opt/picoclaw/workspace-qa_specialist"
	if got != want {
		t.Fatalf("resolveAgentWorkspace() = %q, want %q", got, want)
	}
}

func TestResolveAgentWorkspace_NonDefaultStillUsesExplicitWorkspace(t *testing.T) {
	defaults := &config.AgentDefaults{
		Workspace: "/opt/picoclaw/workspace",
	}
	agentCfg := &config.AgentConfig{
		ID:        "qa_specialist",
		Workspace: "/srv/agents/qa-specialist",
	}

	got := resolveAgentWorkspace(agentCfg, defaults)
	want := "/srv/agents/qa-specialist"
	if got != want {
		t.Fatalf("resolveAgentWorkspace() = %q, want %q", got, want)
	}
}

func TestResolveAgentWorkspace_NonDefaultFallsBackToHomeWhenNoDefaultWorkspace(t *testing.T) {
	defaults := &config.AgentDefaults{
		Workspace: "",
	}
	agentCfg := &config.AgentConfig{
		ID: "qa_specialist",
	}

	home, _ := os.UserHomeDir()
	want := filepath.Join(home, ".picoclaw", "workspace-qa_specialist")
	got := resolveAgentWorkspace(agentCfg, defaults)
	if got != want {
		t.Fatalf("resolveAgentWorkspace() = %q, want %q", got, want)
	}
}

func TestNewAgentInstance_DefaultsTemperatureWhenZero(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "agent-instance-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	cfg := &config.Config{
		Agents: config.AgentsConfig{
			Defaults: config.AgentDefaults{
				Workspace:         tmpDir,
				Model:             "test-model",
				MaxTokens:         1234,
				MaxToolIterations: 5,
			},
		},
	}

	configuredTemp := 0.0
	cfg.Agents.Defaults.Temperature = &configuredTemp

	agent := NewAgentInstance(nil, &cfg.Agents.Defaults, cfg, testFactory)

	if agent.Temperature != 0.0 {
		t.Fatalf("Temperature = %f, want %f", agent.Temperature, 0.0)
	}
}

func TestNewAgentInstance_DefaultsTemperatureWhenUnset(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "agent-instance-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	cfg := &config.Config{
		Agents: config.AgentsConfig{
			Defaults: config.AgentDefaults{
				Workspace:         tmpDir,
				Model:             "test-model",
				MaxTokens:         1234,
				MaxToolIterations: 5,
			},
		},
	}

	agent := NewAgentInstance(nil, &cfg.Agents.Defaults, cfg, testFactory)

	if agent.Temperature != 0.7 {
		t.Fatalf("Temperature = %f, want %f", agent.Temperature, 0.7)
	}
}
