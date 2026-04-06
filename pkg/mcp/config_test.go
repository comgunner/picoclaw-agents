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
)

func TestMCPServerConfig_ValidateCommand_StdioValid(t *testing.T) {
	cfg := config.MCPServerConfig{
		Transport: "stdio",
		Command:   "npx",
		Args:      []string{"-y", "@modelcontextprotocol/server-filesystem"},
	}

	if err := cfg.ValidateCommand(); err != nil {
		t.Errorf("expected no error for valid command, got %v", err)
	}
}

func TestMCPServerConfig_ValidateCommand_StdioInvalid(t *testing.T) {
	cfg := config.MCPServerConfig{
		Transport: "stdio",
		Command:   "evil_command",
	}

	if err := cfg.ValidateCommand(); err == nil {
		t.Error("expected error for invalid command, got nil")
	}
}

func TestMCPServerConfig_ValidateCommand_NonStdio(t *testing.T) {
	cfg := config.MCPServerConfig{
		Transport: "http",
		URL:       "http://example.com",
	}

	if err := cfg.ValidateCommand(); err != nil {
		t.Errorf("expected no error for non-stdio transport, got %v", err)
	}
}
