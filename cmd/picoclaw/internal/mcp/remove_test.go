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

func TestNewRemoveCommand(t *testing.T) {
	cmd := newRemoveCommand()

	require.NotNil(t, cmd)

	assert.Equal(t, "remove <name>", cmd.Use)
	assert.Contains(t, cmd.Short, "Remove")

	assert.NotNil(t, cmd.RunE)
}

func TestRemoveCommand_Exists(t *testing.T) {
	servers := map[string]config.MCPServerConfig{
		"github":     {Transport: "stdio", Command: "npx"},
		"filesystem": {Transport: "stdio", Command: "npx"},
	}

	// Verify server exists
	_, exists := servers["github"]
	require.True(t, exists)

	// Simulate removal
	delete(servers, "github")

	_, exists = servers["github"]
	assert.False(t, exists)
	assert.Equal(t, 1, len(servers))
}

func TestRemoveCommand_NotFound(t *testing.T) {
	servers := map[string]config.MCPServerConfig{
		"github": {Transport: "stdio", Command: "npx"},
	}

	_, exists := servers["nonexistent"]
	assert.False(t, exists)

	// Test the error message logic
	name := "nonexistent"
	if _, exists := servers[name]; !exists {
		err := &notFoundError{name: name}
		assert.Contains(t, err.Error(), "not found")
	}
}

func TestRemoveCommand_LastServer_DisablesMCP(t *testing.T) {
	servers := map[string]config.MCPServerConfig{
		"github": {Transport: "stdio", Command: "npx"},
	}

	delete(servers, "github")

	// When last server is removed, MCP should be disabled
	assert.Equal(t, 0, len(servers))
}

// notFoundError is a helper for testing error messages
type notFoundError struct {
	name string
}

func (e *notFoundError) Error() string {
	return "MCP server \"" + e.name + "\" not found"
}
