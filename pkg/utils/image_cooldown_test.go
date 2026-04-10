// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
// Copyright (c) 2026 PicoClaw contributors

package utils

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewImageCooldown_CreatesDBFile(t *testing.T) {
	tmpDir := t.TempDir()
	cd, err := NewImageCooldown(tmpDir)
	if err != nil {
		t.Fatalf("NewImageCooldown() error = %v", err)
	}
	defer cd.Close()

	expectedPath := filepath.Join(tmpDir, "tmp", "picoclaw_image_cooldown.db")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("Expected DB at %s but not found", expectedPath)
	}
}

func TestImageCooldown_SetAndGet(t *testing.T) {
	cd, _ := NewImageCooldown(t.TempDir())
	defer cd.Close()

	err := cd.Set(120, "antigravity", "gemini-3.1-flash-image")
	if err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	if !cd.IsOnCooldown() {
		t.Error("Expected cooldown to be active after Set()")
	}

	remaining := cd.GetRemaining()
	if remaining <= 0 || remaining > 120 {
		t.Errorf("GetRemaining() = %v, want (0, 120]", remaining)
	}
}

func TestImageCooldown_Expired(t *testing.T) {
	cd, _ := NewImageCooldown(t.TempDir())
	defer cd.Close()

	err := cd.Set(1, "antigravity", "gemini-3.1-flash-image")
	if err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	time.Sleep(1100 * time.Millisecond)

	if cd.IsOnCooldown() {
		t.Error("Expected cooldown to be expired")
	}

	if cd.GetRemaining() != 0 {
		t.Errorf("GetRemaining() after expiry = %v, want 0", cd.GetRemaining())
	}
}

// TestImageCooldown_GetInfo_NoDeadlock verifies Fix #1: GetInfo() must not deadlock.
// It runs GetInfo() 100 times concurrently — if there's a deadlock, this test will timeout.
func TestImageCooldown_GetInfo_NoDeadlock(t *testing.T) {
	cd, _ := NewImageCooldown(t.TempDir())
	defer cd.Close()

	cd.Set(300, "antigravity", "gemini-3.1-flash-image")

	done := make(chan bool, 100)
	for i := 0; i < 100; i++ {
		go func() {
			info := cd.GetInfo()
			if info["active"] != true {
				t.Error("Expected active cooldown")
			}
			done <- true
		}()
	}
	for i := 0; i < 100; i++ {
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Fatal("DEADLOCK detected — GetInfo() did not return in 5s")
		}
	}
}

func TestImageCooldown_Clear(t *testing.T) {
	cd, _ := NewImageCooldown(t.TempDir())
	defer cd.Close()

	cd.Set(300, "antigravity", "gemini-3.1-flash-image")
	if !cd.IsOnCooldown() {
		t.Fatal("Expected cooldown active")
	}

	err := cd.Clear()
	if err != nil {
		t.Fatalf("Clear() error = %v", err)
	}

	if cd.IsOnCooldown() {
		t.Error("Expected cooldown cleared")
	}
}

func TestImageCooldown_Persistence(t *testing.T) {
	tmpDir := t.TempDir()

	cd1, _ := NewImageCooldown(tmpDir)
	cd1.Set(300, "antigravity", "gemini-3.1-flash-image")
	cd1.Close()

	cd2, err := NewImageCooldown(tmpDir)
	if err != nil {
		t.Fatalf("Reopen error = %v", err)
	}
	defer cd2.Close()

	if !cd2.IsOnCooldown() {
		t.Error("Expected cooldown persisted after reopen")
	}
}

func TestImageCooldown_GetRemaining_Accuracy(t *testing.T) {
	cd, _ := NewImageCooldown(t.TempDir())
	defer cd.Close()

	cd.Set(10, "antigravity", "test")
	time.Sleep(2 * time.Second)

	remaining := cd.GetRemaining()
	if remaining < 7 || remaining > 9 {
		t.Errorf("GetRemaining() after 2s sleep = %v, want ~8", remaining)
	}
}

func TestImageCooldown_GetInfo_WhenNotSet(t *testing.T) {
	cd, _ := NewImageCooldown(t.TempDir())
	defer cd.Close()

	info := cd.GetInfo()
	if info["active"] != false {
		t.Error("Expected inactive cooldown when nothing is set")
	}
	if info["remaining"] != 0.0 {
		t.Errorf("Expected remaining=0.0, got %v", info["remaining"])
	}
}

func TestImageCooldown_ResolvePath(t *testing.T) {
	// Test empty string fallback
	result := resolveCooldownWorkspacePath("")
	if result == "" {
		t.Error("Expected non-empty fallback path")
	}

	// Test tilde expansion
	result = resolveCooldownWorkspacePath("~")
	if result == "~" {
		t.Error("Expected tilde to be expanded")
	}

	// Test absolute path passthrough
	result = resolveCooldownWorkspacePath("/some/absolute/path")
	if result != "/some/absolute/path" {
		t.Errorf("Expected /some/absolute/path, got %s", result)
	}
}
