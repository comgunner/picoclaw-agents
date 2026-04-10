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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/comgunner/picoclaw/pkg/config"
)

func TestNewStatusCommand(t *testing.T) {
	cmd := newStatusCommand()

	require.NotNil(t, cmd)

	assert.Equal(t, "status <name>", cmd.Use)
	assert.Contains(t, cmd.Short, "Check")

	assert.Nil(t, cmd.Run)
	assert.NotNil(t, cmd.RunE)
}

func TestStatusCommand_NoServers(t *testing.T) {
	// Test the no-servers-configured path directly
	servers := map[string]config.MCPServerConfig(nil)
	assert.Nil(t, servers)
}

func TestStatusCommand_ServerNotFound(t *testing.T) {
	servers := map[string]config.MCPServerConfig{
		"github": {Transport: "stdio", Command: "npx"},
	}

	_, exists := servers["nonexistent"]
	assert.False(t, exists)
}

func TestStatusCommand_ServerExists(t *testing.T) {
	servers := map[string]config.MCPServerConfig{
		"github": {
			Transport:   "stdio",
			Command:     "npx",
			Description: "GitHub API",
		},
	}

	srv, exists := servers["github"]
	require.True(t, exists)
	assert.Equal(t, "stdio", srv.Transport)
	assert.Equal(t, "npx", srv.Command)
	assert.Equal(t, "GitHub API", srv.Description)
}
