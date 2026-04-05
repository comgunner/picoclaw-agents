// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package onboard

import (
	"encoding/json"
)

// agentEntry representa un agente en agents.list
type agentEntry struct {
	ID            string        `json:"id"`
	Name          string        `json:"name"`
	Default       bool          `json:"default,omitempty"`
	Model         string        `json:"model,omitempty"`
	Skills        []string      `json:"skills,omitempty"`
	ToolsOverride []string      `json:"tools_override,omitempty"`
	Subagents     *subagentsCfg `json:"subagents,omitempty"`
}

type subagentsCfg struct {
	AllowAgents         []string `json:"allow_agents,omitempty"`
	MaxSpawnDepth       int      `json:"max_spawn_depth,omitempty"`
	MaxChildrenPerAgent int      `json:"max_children_per_agent,omitempty"`
}

// buildAgentListJSON genera el JSON de agents.list según el modo y template elegidos
func buildAgentListJSON(mode, template, modelName string, customSkills []string) string {
	var agents []agentEntry

	switch mode {
	case "solo":
		agents = buildSoloAgent(modelName, customSkills)
	case "team":
		agents = buildTeamTemplate(template, modelName)
	case "custom":
		agents = buildCustomAgent(modelName, customSkills)
	default:
		agents = buildSoloAgent(modelName, nil)
	}

	b, err := json.MarshalIndent(agents, "    ", "  ")
	if err != nil {
		return "[]"
	}
	return string(b)
}

// buildSoloAgent construye un único agente con las skills seleccionadas
func buildSoloAgent(modelName string, skills []string) []agentEntry {
	if len(skills) == 0 {
		// Default solo agent sin skills específicas
		return []agentEntry{
			{
				ID:      "main",
				Name:    "Main Agent",
				Default: true,
				Model:   modelName,
			},
		}
	}

	return []agentEntry{
		{
			ID:      "main",
			Name:    "Main Agent",
			Default: true,
			Model:   modelName,
			Skills:  skills,
		},
	}
}

// buildCustomAgent construye un agente custom con skills seleccionadas
func buildCustomAgent(modelName string, skills []string) []agentEntry {
	return buildSoloAgent(modelName, skills)
}

// buildTeamTemplate genera una lista de agentes según el template
func buildTeamTemplate(template, modelName string) []agentEntry {
	switch template {
	case "dev":
		return devTeamAgents(modelName)
	case "research":
		return researchTeamAgents(modelName)
	case "general":
		return generalTeamAgents(modelName)
	default:
		return devTeamAgents(modelName)
	}
}

// devTeamAgents genera el equipo de desarrollo (equivalente a config_dev.example.json)
func devTeamAgents(model string) []agentEntry {
	specialists := []string{
		"backend_dev",
		"frontend_dev",
		"devops_eng",
		"qa_eng",
		"security_eng",
		"data_eng",
		"ml_eng",
		"researcher",
	}

	return []agentEntry{
		{
			ID:      "engineering_manager",
			Name:    "Engineering Manager",
			Default: true,
			Model:   model,
			Skills:  []string{"fullstack_developer", "agent_team_workflow"},
			Subagents: &subagentsCfg{
				AllowAgents:         specialists,
				MaxSpawnDepth:       3,
				MaxChildrenPerAgent: 3,
			},
		},
		{
			ID:            "backend_dev",
			Name:          "Backend Developer",
			Model:         model,
			Skills:        []string{"backend_developer"},
			ToolsOverride: []string{"read_file", "write_file", "edit_file", "exec", "web_search"},
			Subagents:     &subagentsCfg{},
		},
		{
			ID:            "frontend_dev",
			Name:          "Frontend Developer",
			Model:         model,
			Skills:        []string{"frontend_developer"},
			ToolsOverride: []string{"read_file", "write_file", "edit_file"},
			Subagents:     &subagentsCfg{},
		},
		{
			ID:            "devops_eng",
			Name:          "DevOps Engineer",
			Model:         model,
			Skills:        []string{"devops_engineer"},
			ToolsOverride: []string{"read_file", "write_file", "exec"},
			Subagents:     &subagentsCfg{},
		},
		{
			ID:            "qa_eng",
			Name:          "QA Engineer",
			Model:         model,
			Skills:        []string{"qa_engineer"},
			ToolsOverride: []string{"read_file", "write_file", "exec"},
			Subagents:     &subagentsCfg{},
		},
		{
			ID:            "security_eng",
			Name:          "Security Engineer",
			Model:         model,
			Skills:        []string{"security_engineer"},
			ToolsOverride: []string{"read_file", "web_search"},
			Subagents:     &subagentsCfg{},
		},
		{
			ID:            "data_eng",
			Name:          "Data Engineer",
			Model:         model,
			Skills:        []string{"data_engineer"},
			ToolsOverride: []string{"read_file", "write_file", "exec", "web_search"},
			Subagents:     &subagentsCfg{},
		},
		{
			ID:            "ml_eng",
			Name:          "ML Engineer",
			Model:         model,
			Skills:        []string{"ml_engineer"},
			ToolsOverride: []string{"read_file", "write_file", "exec", "web_search"},
			Subagents:     &subagentsCfg{},
		},
		{
			ID:            "researcher",
			Name:          "Researcher",
			Model:         model,
			Skills:        []string{"researcher"},
			ToolsOverride: []string{"web_search", "web_fetch", "read_file"},
			Subagents:     &subagentsCfg{},
		},
	}
}

// researchTeamAgents genera un equipo de investigación
func researchTeamAgents(model string) []agentEntry {
	return []agentEntry{
		{
			ID:      "coordinator",
			Name:    "Research Coordinator",
			Default: true,
			Model:   model,
			Skills:  []string{"agent_team_workflow", "researcher"},
			Subagents: &subagentsCfg{
				AllowAgents:         []string{"researcher", "analyst"},
				MaxSpawnDepth:       2,
				MaxChildrenPerAgent: 3,
			},
		},
		{
			ID:            "researcher",
			Name:          "Researcher",
			Model:         model,
			Skills:        []string{"researcher"},
			ToolsOverride: []string{"web_search", "web_fetch", "read_file", "write_file"},
			Subagents:     &subagentsCfg{},
		},
		{
			ID:            "analyst",
			Name:          "Data Analyst",
			Model:         model,
			Skills:        []string{"data_engineer"},
			ToolsOverride: []string{"read_file", "write_file", "exec"},
			Subagents:     &subagentsCfg{},
		},
	}
}

// generalTeamAgents genera un equipo multi-agente general
func generalTeamAgents(model string) []agentEntry {
	return []agentEntry{
		{
			ID:      "orchestrator",
			Name:    "Orchestrator",
			Default: true,
			Model:   model,
			Skills:  []string{"agent_team_workflow", "fullstack_developer"},
			Subagents: &subagentsCfg{
				AllowAgents:         []string{"worker_a", "worker_b"},
				MaxSpawnDepth:       2,
				MaxChildrenPerAgent: 2,
			},
		},
		{
			ID:        "worker_a",
			Name:      "Worker A",
			Model:     model,
			Subagents: &subagentsCfg{},
		},
		{
			ID:        "worker_b",
			Name:      "Worker B",
			Model:     model,
			Subagents: &subagentsCfg{},
		},
	}
}

// getNativeSkills devuelve la lista de skills nativas disponibles
func getNativeSkills() []string {
	return []string{
		"fullstack_developer",
		"agent_team_workflow",
		"researcher",
		"backend_developer",
		"frontend_developer",
		"devops_engineer",
		"qa_engineer",
		"security_engineer",
		"data_engineer",
		"ml_engineer",
		"odoo_developer",
		"n8n_workflow",
		"binance_mcp",
		"queue_batch",
	}
}

// getSkillDescription devuelve una descripción corta para cada skill
func getSkillDescription(skill string) string {
	descriptions := map[string]string{
		"fullstack_developer": "Full-stack web development",
		"agent_team_workflow": "Coordinate multi-agent tasks",
		"researcher":          "Web research and analysis",
		"backend_developer":   "Backend / API development",
		"frontend_developer":  "UI / React / Vue development",
		"devops_engineer":     "CI/CD, Docker, infrastructure",
		"qa_engineer":         "Testing and quality assurance",
		"security_engineer":   "Security auditing",
		"data_engineer":       "Data pipelines and SQL",
		"ml_engineer":         "Machine learning workflows",
		"odoo_developer":      "Odoo ERP development",
		"n8n_workflow":        "n8n automation workflows",
		"binance_mcp":         "Binance trading via MCP",
		"queue_batch":         "Batch job queue management",
	}

	if desc, ok := descriptions[skill]; ok {
		return desc
	}
	return ""
}
