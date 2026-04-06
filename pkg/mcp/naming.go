// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package mcp

import (
	"fmt"
	"hash/fnv"
	"strings"
	"unicode"
)

// MaxToolNameLen is the maximum tool name length for OpenAI compatibility.
const MaxToolNameLen = 64

// sanitizeName converts a string to a safe tool name.
func sanitizeName(s string) string {
	var b strings.Builder
	prevUnderscore := false
	for _, r := range strings.ToLower(s) {
		if r == '_' || r == '-' {
			if !prevUnderscore && b.Len() > 0 {
				b.WriteRune('_')
				prevUnderscore = true
			}
			continue
		}
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
			prevUnderscore = false
		}
	}
	// Trim trailing underscore
	result := strings.TrimRight(b.String(), "_")
	if result == "" {
		result = "unnamed"
	}
	return result
}

// MCPToolName generates a collision-free tool name for an MCP tool.
// Format: mcp_{server}_{tool} with FNV32a hash suffix if > 64 chars.
func MCPToolName(serverName, toolName string) string {
	name := fmt.Sprintf("mcp_%s_%s", sanitizeName(serverName), sanitizeName(toolName))
	if len(name) <= MaxToolNameLen {
		return name
	}
	// Truncate + FNV32a hash suffix (icueth pattern)
	h := fnv.New32a()
	h.Write([]byte(name))
	suffix := fmt.Sprintf("_%x", h.Sum32())
	maxLen := MaxToolNameLen - len(suffix)
	if maxLen < 1 {
		maxLen = 1
	}
	return name[:maxLen] + suffix
}
