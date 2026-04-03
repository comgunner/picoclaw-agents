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
	"fmt"
	"os"
	"path/filepath"
	"strings"

	managementinit "github.com/comgunner/picoclaw/pkg/agents/management/init"
	"github.com/comgunner/picoclaw/pkg/config"
	"github.com/comgunner/picoclaw/pkg/logger"
	"github.com/comgunner/picoclaw/pkg/providers"
	"github.com/comgunner/picoclaw/pkg/routing"
	"github.com/comgunner/picoclaw/pkg/session"
	"github.com/comgunner/picoclaw/pkg/tools"
)

// AgentInstance represents a fully configured agent with its own workspace,
// session manager, context builder, and tool registry.
type AgentInstance struct {
	ID             string
	Name           string
	Model          string
	Fallbacks      []string
	Workspace      string
	MaxIterations  int
	MaxTokens      int
	Temperature    float64
	ContextWindow  int
	Provider       providers.LLMProvider
	Sessions       *session.SessionManager
	ContextBuilder *ContextBuilder
	Tools          *tools.ToolRegistry
	Subagents      *config.SubagentsConfig
	Runtime        *config.AgentRuntimeConfig
	SkillsFilter   []string
	Candidates     []providers.FallbackCandidate

	// A2A Components (Agent-to-Agent Communication)
	mailbox          any // *mailbox.Mailbox - assigned by orchestrator
	departmentRouter any //nolint:unused // *DepartmentRouter - reserved for future A2A routing

	Role       string // Agent role (e.g., "coordinator", "developer", "tester")
	Status     string // Current status (e.g., "idle", "running", "waiting")
	LastActive int64  // Last activity timestamp (Unix seconds)
}

// ProviderFactory is a function that returns a provider for a given model.
type ProviderFactory func(model string) (providers.LLMProvider, string, error)

// NewAgentInstance creates an agent instance from config.
func NewAgentInstance(
	agentCfg *config.AgentConfig,
	defaults *config.AgentDefaults,
	cfg *config.Config,
	factory ProviderFactory,
) *AgentInstance {
	workspace := resolveAgentWorkspace(agentCfg, defaults)
	os.MkdirAll(workspace, 0o755)

	model := resolveAgentModel(agentCfg, defaults)
	fallbacks := resolveAgentFallbacks(agentCfg, defaults)

	// Resolve the provider for this specific agent's model
	provider, resolvedModel, err := factory(model)
	if err != nil {
		logger.ErrorCF("agent", "Failed to resolve provider for agent",
			map[string]any{
				"agent_id": agentCfg.ID,
				"model":    model,
				"error":    err.Error(),
			})
		// If factory fails, we might still have a "default" provider from the loop
		// But it's better to log error and let it use whatever was returned.
	}
	if resolvedModel != "" {
		model = resolvedModel
	}

	restrict := defaults.RestrictToWorkspace
	globalWorkspace := ""
	if cfg != nil {
		globalWorkspace = cfg.WorkspacePath()
	}
	toolsRegistry := tools.NewToolRegistryWithWorkspace(globalWorkspace)
	toolsRegistry.Register(tools.NewReadFileTool(workspace, restrict))
	toolsRegistry.Register(tools.NewWriteFileTool(workspace, restrict))
	toolsRegistry.Register(tools.NewListDirTool(workspace, restrict))
	execTool, err := tools.NewExecToolWithConfig(workspace, restrict, cfg)
	if err != nil {
		// Fallback to basic tool or log error - given it's a security tool, we should probably fail initialization
		// or at least log it very loudly. For now, we'll log and skip registration if it fails.
		panic(fmt.Sprintf("Failed to initialize ExecTool: %v", err))
	}
	toolsRegistry.Register(execTool)
	toolsRegistry.Register(tools.NewWorkspaceMaintenanceTool(""))
	toolsRegistry.Register(tools.NewEditFileTool(workspace, restrict))
	toolsRegistry.Register(tools.NewAppendFileTool(workspace, restrict))

	toolsRegistry.Register(tools.NewSystemDiagnosticsTool(workspace))
	toolsRegistry.Register(tools.NewConfigManagerTool(filepath.Join(workspace, "..", "config.json")))
	toolsRegistry.Register(tools.NewResourceMonitorTool(workspace))
	toolsRegistry.Register(tools.NewMemoryStoreTool(workspace))
	toolsRegistry.Register(tools.NewVersionControlTool(workspace))
	toolsRegistry.Register(tools.NewSelfDiagnosticsTool(workspace, filepath.Join(workspace, "..", "config.json")))

	// Register all 12 Agent Management tools (agent_list, agent_get, agent_can_spawn, etc.)
	managementinit.RegisterManagementTools(cfg, toolsRegistry)

	logger.InfoCF("agent", "Native tools registered",
		map[string]any{
			"workspace":          workspace,
			"native_tools_count": 5,
			"tools": []string{
				"system_diagnostics",
				"config_manager",
				"resource_monitor",
				"memory_store",
				"version_control",
			},
		})

	sessionsDir := filepath.Join(workspace, "sessions")
	sessionsManager := session.NewSessionManager(sessionsDir)

	contextBuilder := NewContextBuilder(workspace)

	agentID := routing.DefaultAgentID
	agentName := ""
	var subagents *config.SubagentsConfig
	var skillsFilter []string

	if agentCfg != nil {
		agentID = routing.NormalizeAgentID(agentCfg.ID)
		agentName = agentCfg.Name
		subagents = agentCfg.Subagents
		skillsFilter = agentCfg.Skills
	}

	maxIter := defaults.MaxToolIterations
	if maxIter == 0 {
		maxIter = 20
	}

	maxTokens := defaults.MaxTokens
	if maxTokens == 0 {
		maxTokens = 8192
	}

	temperature := 0.7
	if defaults.Temperature != nil {
		temperature = *defaults.Temperature
	}

	// Resolve fallback candidates.
	// Model names may be model_list aliases (e.g. "openrouter-free") rather than
	// real provider/model IDs (e.g. "openrouter/free").  Resolve them first so that
	// ResolveCandidates receives the actual model strings.
	var modelList []config.ModelConfig
	if cfg != nil {
		modelList = cfg.ModelList
	}
	modelCfg := providers.ModelConfig{
		Primary:   resolveModelAlias(model, modelList),
		Fallbacks: resolveModelAliases(fallbacks, modelList),
	}
	candidates := providers.ResolveCandidates(modelCfg, defaults.Provider)

	var runtimeCfg *config.AgentRuntimeConfig
	if agentCfg != nil && agentCfg.Runtime != nil {
		runtimeCfg = agentCfg.Runtime
	} else if defaults.Runtime != nil {
		runtimeCfg = defaults.Runtime
	}

	return &AgentInstance{
		ID:             agentID,
		Name:           agentName,
		Model:          modelCfg.Primary, // FIX: Use resolved model (e.g. "gpt-4o") not alias (e.g. "chatgpt-gpt-4o")
		Fallbacks:      fallbacks,
		Workspace:      workspace,
		MaxIterations:  maxIter,
		MaxTokens:      maxTokens,
		Temperature:    temperature,
		ContextWindow:  maxTokens,
		Provider:       provider,
		Sessions:       sessionsManager,
		ContextBuilder: contextBuilder,
		Tools:          toolsRegistry,
		Subagents:      subagents,
		Runtime:        runtimeCfg,
		SkillsFilter:   skillsFilter,
		Candidates:     candidates,
	}
}

// resolveAgentWorkspace determines the workspace directory for an agent.
func resolveAgentWorkspace(agentCfg *config.AgentConfig, defaults *config.AgentDefaults) string {
	if agentCfg != nil && strings.TrimSpace(agentCfg.Workspace) != "" {
		return expandHome(strings.TrimSpace(agentCfg.Workspace))
	}
	defaultWorkspace := expandHome(strings.TrimSpace(defaults.Workspace))
	if agentCfg == nil || agentCfg.Default || agentCfg.ID == "" || routing.NormalizeAgentID(agentCfg.ID) == "main" {
		return defaultWorkspace
	}

	id := routing.NormalizeAgentID(agentCfg.ID)
	if defaultWorkspace != "" {
		parent := filepath.Dir(defaultWorkspace)
		base := filepath.Base(defaultWorkspace)
		if parent != "" && base != "." && base != "/" {
			return filepath.Join(parent, base+"-"+id)
		}
	}

	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".picoclaw", "workspace-"+id)
}

// resolveAgentModel resolves the primary model for an agent.
// FIX #901: Normalize model names (openrouter-free → openrouter/auto)
func resolveAgentModel(agentCfg *config.AgentConfig, defaults *config.AgentDefaults) string {
	var model string

	// First, get the raw model name (without normalization)
	if agentCfg != nil && agentCfg.Model != nil && strings.TrimSpace(agentCfg.Model.Primary) != "" {
		model = strings.TrimSpace(agentCfg.Model.Primary)
	} else {
		model = defaults.GetModelName()
	}

	// Then normalize (openrouter-free → openrouter/auto, etc.)
	// FIX #901: Normalize the model name
	return providers.NormalizeModelName(model)
}

// resolveAgentFallbacks resolves the fallback models for an agent.
func resolveAgentFallbacks(agentCfg *config.AgentConfig, defaults *config.AgentDefaults) []string {
	if agentCfg != nil && agentCfg.Model != nil && agentCfg.Model.Fallbacks != nil {
		return agentCfg.Model.Fallbacks
	}
	return defaults.ModelFallbacks
}

func expandHome(path string) string {
	if path == "" {
		return path
	}
	if path[0] == '~' {
		home, _ := os.UserHomeDir()
		if len(path) > 1 && path[1] == '/' {
			return home + path[1:]
		}
		return home
	}
	return path
}

// resolveModelAlias looks up a model_name alias in the model_list and returns
// the actual model string (e.g. "openrouter/auto").  If not found, the original
// name is returned unchanged so existing provider/model strings pass through.
// FIX #901: Also normalizes deprecated model names (openrouter-free → openrouter/auto)
func resolveModelAlias(name string, modelList []config.ModelConfig) string {
	for _, m := range modelList {
		if m.ModelName == name {
			// FIX #901: Normalize the model string from model_list
			return providers.NormalizeModelName(m.Model)
		}
	}
	// FIX #901: Normalize if not found in model_list
	return providers.NormalizeModelName(name)
}

func resolveModelAliases(names []string, modelList []config.ModelConfig) []string {
	resolved := make([]string, len(names))
	for i, n := range names {
		resolved[i] = resolveModelAlias(n, modelList)
	}
	return resolved
}
