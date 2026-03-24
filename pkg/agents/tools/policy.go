// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package tools

import "strings"

const (
	// SUBAGENT_TOOL_DENY_ALWAYS contains tools that subagents are never allowed to use.
	// These are system-level or interactive tools that could compromise security.
	SUBAGENT_TOOL_DENY_ALWAYS = "gateway,agents_list,memory_search,memory_get,sessions_send"

	// SUBAGENT_TOOL_DENY_LEAF contains tools that only make sense for orchestrator subagents
	// that can spawn children. Leaf subagents (depth >= maxSpawnDepth) should not use these.
	SUBAGENT_TOOL_DENY_LEAF = "sessions_list,sessions_history,sessions_spawn,subagents"
)

// GetToolDenyListForSubagent returns the combined deny list for a subagent at a given depth.
// - denyAlways: tools denied for all subagents regardless of depth
// - denyLeaf: tools denied for leaf subagents (depth >= maxDepth)
func GetToolDenyListForSubagent(depth int, maxDepth int) []string {
	denyList := strings.Split(SUBAGENT_TOOL_DENY_ALWAYS, ",")

	// If this is a leaf node (depth >= maxDepth), add leaf-specific denials
	if depth >= maxDepth {
		leafDenyList := strings.Split(SUBAGENT_TOOL_DENY_LEAF, ",")
		denyList = append(denyList, leafDenyList...)
	}

	return denyList
}

// IsToolAllowedForSubagent checks if a tool is allowed for a subagent at a given depth.
// Returns true if the tool is allowed, false if it's denied by policy.
func IsToolAllowedForSubagent(toolName string, depth int, maxDepth int) bool {
	denyList := GetToolDenyListForSubagent(depth, maxDepth)

	toolLower := strings.ToLower(toolName)
	for _, denied := range denyList {
		if strings.TrimSpace(denied) == toolLower {
			return false
		}
	}

	return true
}

// FilterToolsForSubagent applies policy restrictions to a list of available tools,
// returning only those that are allowed for a subagent at the given depth.
func FilterToolsForSubagent(availableTools []string, depth int, maxDepth int) []string {
	var allowed []string
	for _, tool := range availableTools {
		if IsToolAllowedForSubagent(tool, depth, maxDepth) {
			allowed = append(allowed, tool)
		}
	}
	return allowed
}
