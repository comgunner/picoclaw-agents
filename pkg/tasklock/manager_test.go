// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package tasklock

import (
	"os"
	"testing"
	"time"

	"github.com/comgunner/picoclaw/pkg/providers"
)

func TestTaskLockManager_Lifecycle(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewManager(tempDir)

	// Create a new lock
	tl, err := manager.CreateLock("task-123", "qa_agent", "parent_mgr", "ses-abc")
	if err != nil {
		t.Fatalf("failed to create lock: %v", err)
	}

	if tl.GetStatus() != StatusInProgress {
		t.Errorf("expected status %s, got %s", StatusInProgress, tl.GetStatus())
	}

	// Test claiming files
	files := []string{"main.go", "utils.go"}
	if err := tl.ClaimFiles(files); err != nil {
		t.Fatalf("failed to claim files: %v", err)
	}

	locked, owner := manager.IsFileLocked("main.go")
	if !locked || owner != "qa_agent" {
		t.Errorf("expected main.go to be locked by qa_agent")
	}

	// Update Context
	msg := []providers.Message{{Role: "user", Content: "Hello"}}
	if err := tl.UpdateState(StatusInProgress, "waiting", msg); err != nil {
		t.Fatalf("failed to update state: %v", err)
	}

	if len(tl.GetContext()) != 1 {
		t.Errorf("expected context length 1")
	}

	// Remove Lock
	if err := manager.RemoveLock("task-123"); err != nil {
		t.Fatalf("failed to remove lock: %v", err)
	}

	// Verify File is unlocked
	locked, _ = manager.IsFileLocked("main.go")
	if locked {
		t.Errorf("expected main.go to be unlocked")
	}

	// Verify File is deleted
	if _, err := os.Stat(tl.GetFilePath()); !os.IsNotExist(err) {
		t.Errorf("expected lock file to be deleted from disk, but it still exists")
	}
}

func TestTaskLockManager_Hydrate(t *testing.T) {
	tempDir := t.TempDir()
	m1 := NewManager(tempDir)

	tl1, err := m1.CreateLock("task-456", "dev_agent", "mgr", "ses-123")
	if err != nil {
		t.Fatalf("failed to create lock: %v", err)
	}

	// Flush some data
	err = tl1.ClaimFiles([]string{"test.go"})
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	// Wait briefly to ensure file stamp settles if needed, though atomic guarantees immediate available
	time.Sleep(10 * time.Millisecond)

	// Simulate a crash/restart by creating a new manager pointed to the same workspace dir
	m2 := NewManager(tempDir)

	tl2, exists := m2.GetLock("task-456")
	if !exists {
		t.Fatalf("expected lock to be hydrated from disk by new manager instance")
	}

	if tl2.AssignedAgent != "dev_agent" {
		t.Errorf("expected agent 'dev_agent', got '%s'", tl2.AssignedAgent)
	}

	if len(tl2.GetLockedFiles()) != 1 || tl2.GetLockedFiles()[0] != "test.go" {
		t.Errorf("expected locked files to persist through hydration")
	}
}
