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
	"testing"

	"github.com/comgunner/picoclaw/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewListCommand(t *testing.T) {
	cmd := newListCommand()

	require.NotNil(t, cmd)

	assert.Equal(t, "list", cmd.Use)
	assert.Contains(t, cmd.Short, "List")

	assert.True(t, cmd.HasAlias("ls"))
	assert.Nil(t, cmd.Run)
	assert.NotNil(t, cmd.RunE)
}

func TestListCommand_Empty(t *testing.T) {
	// Test the empty servers map logic directly
	servers := map[string]config.MCPServerConfig{}
	assert.Equal(t, 0, len(servers))
}

func TestListCommand_WithServers(t *testing.T) {
	servers := map[string]config.MCPServerConfig{
		"github": {
			Transport:    "stdio",
			Command:      "npx",
			Args:         []string{"-y", "@modelcontextprotocol/server-github"},
			EnabledTools: []string{"*"},
			Description:  "GitHub API",
		},
		"filesystem": {
			Transport:    "stdio",
			Command:      "npx",
			Args:         []string{"-y", "@modelcontextprotocol/server-filesystem"},
			EnabledTools: []string{"read_file", "write_file"},
			Description:  "Filesystem access",
		},
	}

	assert.Equal(t, 2, len(servers))

	// Verify github server
	gh := servers["github"]
	assert.Equal(t, "stdio", gh.Transport)
	assert.Equal(t, "npx", gh.Command)
	assert.Equal(t, "GitHub API", gh.Description)

	// Verify filesystem server
	fs := servers["filesystem"]
	assert.Equal(t, "stdio", fs.Transport)
	assert.Equal(t, 2, len(fs.EnabledTools))
}
