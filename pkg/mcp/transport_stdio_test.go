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
	"time"
)

func TestNewStdioTransport_InvalidCommand(t *testing.T) {
	_, err := NewStdioTransport("nonexistent_command_xyz", nil, nil)
	if err == nil {
		t.Error("expected error for nonexistent command, got nil")
	}
}

func TestNewStdioTransport_ValidCommand(t *testing.T) {
	// Use 'echo' as a simple test command that exists on all systems
	transport, err := NewStdioTransport("echo", []string{"test"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer transport.Close()
}

func TestStdioTransport_Call_ContextCancellation(t *testing.T) {
	// Create a transport using a subprocess that doesn't respond
	// Use 'sleep' which will block indefinitely
	transport, err := NewStdioTransport("sleep", []string{"10"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer transport.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err = transport.Call(ctx, "test_method", nil)
	if err == nil {
		t.Fatal("expected context cancellation error, got nil")
	}
	if err != context.DeadlineExceeded {
		t.Errorf("expected context.DeadlineExceeded, got %v", err)
	}
}

func TestStdioTransport_Call_MAX_LINE_BYTES_Protection(t *testing.T) {
	// This test verifies the MAX_LINE_BYTES constant exists and is correct
	if MAX_LINE_BYTES != 10*1024*1024 {
		t.Errorf("expected MAX_LINE_BYTES to be 10MB, got %d", MAX_LINE_BYTES)
	}
}

func TestStdioTransport_Close(t *testing.T) {
	transport, err := NewStdioTransport("echo", []string{"test"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Close should not panic
	transport.Close()
}
