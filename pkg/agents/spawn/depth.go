// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package spawn

import "strings"

const DefaultSubagentMaxSpawnDepth = 2

func DepthFromSessionKey(sessionKey string) int {
	key := strings.TrimSpace(sessionKey)
	if key == "" {
		return 0
	}
	if strings.Contains(strings.ToLower(key), ":subagent:") {
		return 1
	}
	return 0
}

func AllowedToSpawn(currentDepth, maxDepth int) bool {
	if maxDepth <= 0 {
		maxDepth = DefaultSubagentMaxSpawnDepth
	}
	return currentDepth < maxDepth
}
