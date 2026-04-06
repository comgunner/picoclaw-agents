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
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMCPCommand(t *testing.T) {
	cmd := NewMCPCommand()

	require.NotNil(t, cmd)

	assert.Equal(t, "mcp", cmd.Use)
	assert.Contains(t, cmd.Short, "MCP servers")

	assert.Nil(t, cmd.Run)
	assert.NotNil(t, cmd.RunE)

	assert.False(t, cmd.HasSubCommands() == false) // has subcommands

	allowedCommands := []string{
		"add",
		"list",
		"status",
		"remove",
	}

	subcommands := cmd.Commands()
	for _, subcmd := range subcommands {
		found := slices.Contains(allowedCommands, subcmd.Name())
		assert.True(t, found, "unexpected subcommand %q", subcmd.Name())
	}
}

func TestMCPCommand_Help(t *testing.T) {
	cmd := NewMCPCommand()
	err := cmd.Execute()
	// Should not error — help is displayed
	require.NoError(t, err)
}
