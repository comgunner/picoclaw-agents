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
	"context"
	"testing"

	"github.com/comgunner/picoclaw/pkg/config"
)

func TestMCPClientManager_ConnectAll_Empty(t *testing.T) {
	manager := &MCPClientManager{}
	errs := manager.ConnectAll(nil)
	if len(errs) != 0 {
		t.Errorf("expected no errors for empty servers, got %d", len(errs))
	}
}

func TestMCPClientManager_ConnectAll_InvalidTransport(t *testing.T) {
	manager := &MCPClientManager{}
	servers := map[string]config.MCPServerConfig{
		"test": {Transport: "invalid"},
	}
	errs := manager.ConnectAll(servers)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
}

func TestMCPClientManager_ListServers_Empty(t *testing.T) {
	manager := &MCPClientManager{}
	servers := manager.ListServers()
	if len(servers) != 0 {
		t.Errorf("expected 0 servers, got %d", len(servers))
	}
}

func TestMCPClientManager_GetServer_NotFound(t *testing.T) {
	manager := &MCPClientManager{}
	srv := manager.GetServer("nonexistent")
	if srv != nil {
		t.Error("expected nil server for nonexistent name")
	}
}

func TestMCPClientManager_CallTool_ServerNotFound(t *testing.T) {
	manager := &MCPClientManager{}
	_, err := manager.CallTool(context.Background(), "nonexistent", "tool", nil)
	if err == nil {
		t.Fatal("expected error for nonexistent server, got nil")
	}
}

func TestMCPClientManager_CloseAll_Empty(t *testing.T) {
	manager := &MCPClientManager{}
	// Should not panic
	manager.CloseAll()
}

func TestMCPClientManager_ConnectAll_NonFatal(t *testing.T) {
	manager := &MCPClientManager{}
	servers := map[string]config.MCPServerConfig{
		"bad_server": {
			Transport: "invalid",
		},
	}
	errs := manager.ConnectAll(servers)
	// Should return error but not panic
	if len(errs) == 0 {
		t.Error("expected errors for invalid transport")
	}
}

func TestFilterTools_All(t *testing.T) {
	tools := []ToolInfo{
		{Name: "tool1"},
		{Name: "tool2"},
	}

	result := filterTools(tools, []string{"*"})
	if len(result) != 2 {
		t.Errorf("expected 2 tools, got %d", len(result))
	}
}

func TestFilterTools_Selected(t *testing.T) {
	tools := []ToolInfo{
		{Name: "tool1"},
		{Name: "tool2"},
		{Name: "tool3"},
	}

	result := filterTools(tools, []string{"tool1", "tool3"})
	if len(result) != 2 {
		t.Fatalf("expected 2 tools, got %d", len(result))
	}
	if result[0].Name != "tool1" || result[1].Name != "tool3" {
		t.Errorf("expected tool1 and tool3, got %v", result)
	}
}
