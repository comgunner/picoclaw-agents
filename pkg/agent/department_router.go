// pkg/agent/department_router.go
//
// Department Model Router - Routes models by department with fallback to default
// Based on concepts from icueth's picoclaw-agents fork (@icueth)
// Enhanced with fallback to default model for unassigned departments
//
// Original source: https://github.com/icueth/picoclaw-agents
// License: Same as base project (MIT)

package agent

import (
	"fmt"
	"sync"

	"github.com/comgunner/picoclaw/pkg/config"
	"github.com/comgunner/picoclaw/pkg/logger"
)

// DepartmentModelConfig holds full model configuration for a department
type DepartmentModelConfig struct {
	Model          string
	Temperature    float64
	TopP           float64
	EnableThinking bool
}

// DepartmentRouter handles model routing by department with fallback
type DepartmentRouter struct {
	defaultModel       string
	defaultTemperature float64
	defaultTopP        float64
	departmentModels   map[string]DepartmentModelConfig // department -> model config
	agentDepartments   map[string]string                // agent_id -> department
	mu                 sync.RWMutex
}

// NewDepartmentRouter creates router with fallback to default model
func NewDepartmentRouter(cfg *config.Config) *DepartmentRouter {
	router := &DepartmentRouter{
		departmentModels:   make(map[string]DepartmentModelConfig),
		agentDepartments:   make(map[string]string),
		defaultModel:       cfg.Agents.Defaults.GetModelName(),
		defaultTemperature: 0.7,
	}

	// Read default temperature/topP from config if set
	if cfg.Agents.Defaults.Temperature != nil {
		router.defaultTemperature = *cfg.Agents.Defaults.Temperature
	}
	if cfg.Agents.Defaults.TopP != nil {
		router.defaultTopP = *cfg.Agents.Defaults.TopP
	}

	// Load department_models from config (B-07 fix)
	for dept, rawCfg := range cfg.Agents.DepartmentModels {
		dc := DepartmentModelConfig{
			Model:       router.defaultModel,
			Temperature: router.defaultTemperature,
			TopP:        router.defaultTopP,
		}
		if model, ok := rawCfg["model"].(string); ok && model != "" {
			dc.Model = model
		}
		if t, ok := rawCfg["temperature"].(float64); ok {
			dc.Temperature = t
		}
		if tp, ok := rawCfg["top_p"].(float64); ok {
			dc.TopP = tp
		}
		if et, ok := rawCfg["enable_thinking"].(bool); ok {
			dc.EnableThinking = et
		}
		router.departmentModels[dept] = dc
		logger.InfoCF("department_router", "Loaded department model from config", map[string]any{
			"department": dept,
			"model":      dc.Model,
		})
	}

	logger.InfoCF("department_router", "Initialized with default model", map[string]any{
		"default_model":      router.defaultModel,
		"departments_loaded": len(router.departmentModels),
	})

	return router
}

// RegisterAgent registers an agent to a department (B-09 fix)
func (r *DepartmentRouter) RegisterAgent(agentID string, department string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.agentDepartments[agentID] = department
	logger.InfoCF("department_router", "Agent registered to department", map[string]any{
		"agent":      agentID,
		"department": department,
	})
}

// GetModelConfig returns full model configuration for agent (with fallback)
func (r *DepartmentRouter) GetModelConfig(agentID string) DepartmentModelConfig {
	r.mu.RLock()
	dept, exists := r.agentDepartments[agentID]
	r.mu.RUnlock()

	defaultCfg := DepartmentModelConfig{
		Model:       r.defaultModel,
		Temperature: r.defaultTemperature,
		TopP:        r.defaultTopP,
	}

	if !exists {
		logger.DebugCF("department_router", "No department for agent, using default", map[string]any{
			"agent":       agentID,
			"fallback_to": r.defaultModel,
		})
		return defaultCfg
	}

	r.mu.RLock()
	dc, ok := r.departmentModels[dept]
	r.mu.RUnlock()

	if !ok {
		logger.DebugCF("department_router", "Department has no model, using default", map[string]any{
			"agent":       agentID,
			"department":  dept,
			"fallback_to": r.defaultModel,
		})
		return defaultCfg
	}

	logger.DebugCF("department_router", "Using department model", map[string]any{
		"agent":      agentID,
		"department": dept,
		"model":      dc.Model,
	})
	return dc
}

// UpdateDepartmentModel updates the model name for a department
func (r *DepartmentRouter) UpdateDepartmentModel(department string, model string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Preserve existing config, only update model name
	dc := r.departmentModels[department]
	dc.Model = model
	r.departmentModels[department] = dc

	logger.InfoCF("department_router", "Department model updated", map[string]any{
		"department": department,
		"model":      model,
	})

	return nil
}

// GetDefaultModel returns the fallback/default model name
func (r *DepartmentRouter) GetDefaultModel() string {
	return r.defaultModel
}

// GetDepartmentModel returns model name for specific department
func (r *DepartmentRouter) GetDepartmentModel(department string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if dc, ok := r.departmentModels[department]; ok {
		return dc.Model, nil
	}

	return "", fmt.Errorf("department %s not found, using default", department)
}

// ListDepartments returns all departments with their model names
func (r *DepartmentRouter) ListDepartments() map[string]string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make(map[string]string)
	for dept, dc := range r.departmentModels {
		result[dept] = dc.Model
	}

	return result
}

// GetAgentDepartment returns department for agent
func (r *DepartmentRouter) GetAgentDepartment(agentID string) string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if dept, ok := r.agentDepartments[agentID]; ok {
		return dept
	}

	return "unknown"
}

// RemoveAgent removes agent from department routing
func (r *DepartmentRouter) RemoveAgent(agentID string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.agentDepartments, agentID)
	logger.InfoCF("department_router", "Agent removed from department routing", map[string]any{
		"agent": agentID,
	})
}

// GetDepartmentAgents returns all agents in a department
func (r *DepartmentRouter) GetDepartmentAgents(department string) []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	agents := make([]string, 0)
	for agentID, dept := range r.agentDepartments {
		if dept == department {
			agents = append(agents, agentID)
		}
	}

	return agents
}

// GetDepartmentCount returns the number of configured departments
func (r *DepartmentRouter) GetDepartmentCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.departmentModels)
}

// GetAgentCount returns the number of agents registered to departments
func (r *DepartmentRouter) GetAgentCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.agentDepartments)
}
