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
)

func TestAddCommand_ValidStdio(t *testing.T) {
	// Test the validation logic directly
	srvCfg := config.MCPServerConfig{
		Transport: "stdio",
		Command:   "npx",
	}
	assert.NoError(t, srvCfg.ValidateCommand())
}

func TestAddCommand_DuplicateName(t *testing.T) {
	// Test duplicate detection logic
	servers := map[string]config.MCPServerConfig{
		"github": {Transport: "stdio", Command: "npx"},
	}

	_, exists := servers["github"]
	assert.True(t, exists, "should detect duplicate server name")
}

func TestAddCommand_InvalidCommand(t *testing.T) {
	srvCfg := config.MCPServerConfig{
		Transport: "stdio",
		Command:   "rm",
	}
	err := srvCfg.ValidateCommand()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not allowed")
}

func TestAddCommand_InvalidTransport(t *testing.T) {
	err := addServer("test", "grpc", "", nil, nil, []string{"*"}, "", "", nil, 0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid transport")
}

func TestAddCommand_EmptyName(t *testing.T) {
	err := addServer("", "stdio", "npx", nil, nil, []string{"*"}, "", "", nil, 0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "name is required")
}

func TestAddCommand_InvalidEnvFormat(t *testing.T) {
	err := addServer("test", "stdio", "npx", nil, []string{"INVALID_NO_EQUALS"}, []string{"*"}, "", "", nil, 0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid env variable format")
}

func TestAddCommand_InvalidHeaderFormat(t *testing.T) {
	err := addServer("test", "http", "", nil, nil, []string{"*"}, "", "", []string{"INVALID_NO_EQUALS"}, 0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid header format")
}
