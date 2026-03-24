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
	"testing"
)

func TestGetToolDenyListForSubagent(t *testing.T) {
	tests := []struct {
		name     string
		depth    int
		maxDepth int
		wantLen  int
		wantHas  []string
	}{
		{
			name:     "depth 0 with max 2",
			depth:    0,
			maxDepth: 2,
			wantLen:  5, // gateway,agents_list,memory_search,memory_get,sessions_send
			wantHas:  []string{"gateway", "memory_search"},
		},
		{
			name:     "depth 2 with max 2 (leaf)",
			depth:    2,
			maxDepth: 2,
			wantLen:  9, // all + leaf tools
			wantHas:  []string{"gateway", "sessions_spawn", "subagents"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetToolDenyListForSubagent(tt.depth, tt.maxDepth)
			if len(got) != tt.wantLen {
				t.Errorf("GetToolDenyListForSubagent() length = %d, want %d", len(got), tt.wantLen)
			}
			for _, wantItem := range tt.wantHas {
				found := false
				for _, gotItem := range got {
					if gotItem == wantItem {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("GetToolDenyListForSubagent() missing expected item %s in %v", wantItem, got)
				}
			}
		})
	}
}

func TestIsToolAllowedForSubagent(t *testing.T) {
	tests := []struct {
		name     string
		tool     string
		depth    int
		maxDepth int
		want     bool
	}{
		{
			name:     "gateway always denied",
			tool:     "gateway",
			depth:    0,
			maxDepth: 2,
			want:     false,
		},
		{
			name:     "sessions_spawn allowed at depth 0",
			tool:     "sessions_spawn",
			depth:    0,
			maxDepth: 2,
			want:     true,
		},
		{
			name:     "sessions_spawn denied at leaf depth",
			tool:     "sessions_spawn",
			depth:    2,
			maxDepth: 2,
			want:     false,
		},
		{
			name:     "read_file allowed",
			tool:     "read_file",
			depth:    1,
			maxDepth: 2,
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsToolAllowedForSubagent(tt.tool, tt.depth, tt.maxDepth)
			if got != tt.want {
				t.Errorf("IsToolAllowedForSubagent() = %v, want %v", got, tt.want)
			}
		})
	}
}
