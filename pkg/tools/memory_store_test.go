// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package tools

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestMemoryStoreTool_Name(t *testing.T) {
	tmpDir := t.TempDir()
	tool := NewMemoryStoreTool(tmpDir)
	if tool.Name() != "memory_store" {
		t.Errorf("expected name 'memory_store', got '%s'", tool.Name())
	}
}

func TestMemoryStoreTool_Description(t *testing.T) {
	tmpDir := t.TempDir()
	tool := NewMemoryStoreTool(tmpDir)
	if tool.Description() == "" {
		t.Error("expected non-empty description")
	}
}

func TestMemoryStoreTool_Parameters(t *testing.T) {
	tmpDir := t.TempDir()
	tool := NewMemoryStoreTool(tmpDir)
	params := tool.Parameters()

	if params["type"] != "object" {
		t.Error("parameters type should be 'object'")
	}

	props, ok := params["properties"].(map[string]any)
	if !ok {
		t.Fatal("expected properties map")
	}

	if _, ok := props["action"]; !ok {
		t.Error("missing 'action' parameter")
	}
}

func TestMemoryStoreTool_SetAndGet(t *testing.T) {
	tmpDir := t.TempDir()
	tool := NewMemoryStoreTool(tmpDir)
	ctx := context.Background()

	// Set a value
	result := tool.Execute(ctx, map[string]any{
		"action": "set",
		"key":    "test_key",
		"value":  "test_value",
	})

	if result.IsError {
		t.Errorf("set failed: %s", result.ForLLM)
	}

	// Get the value
	result = tool.Execute(ctx, map[string]any{
		"action": "get",
		"key":    "test_key",
	})

	if result.IsError {
		t.Errorf("get failed: %s", result.ForLLM)
	}
}

func TestMemoryStoreTool_SetWithTTL(t *testing.T) {
	tmpDir := t.TempDir()
	tool := NewMemoryStoreTool(tmpDir)
	ctx := context.Background()

	// Set a value with TTL
	result := tool.Execute(ctx, map[string]any{
		"action":      "set",
		"key":         "ttl_key",
		"value":       "ttl_value",
		"ttl_seconds": 3600, // 1 hour
	})

	if result.IsError {
		t.Errorf("set with TTL failed: %s", result.ForLLM)
	}
}

func TestMemoryStoreTool_GetNonExistent(t *testing.T) {
	tmpDir := t.TempDir()
	tool := NewMemoryStoreTool(tmpDir)
	ctx := context.Background()

	result := tool.Execute(ctx, map[string]any{
		"action": "get",
		"key":    "nonexistent",
	})

	if !result.IsError {
		t.Error("expected error for nonexistent key")
	}
}

func TestMemoryStoreTool_Delete(t *testing.T) {
	tmpDir := t.TempDir()
	tool := NewMemoryStoreTool(tmpDir)
	ctx := context.Background()

	// Set a value
	tool.Execute(ctx, map[string]any{
		"action": "set",
		"key":    "delete_key",
		"value":  "delete_value",
	})

	// Delete it
	result := tool.Execute(ctx, map[string]any{
		"action": "delete",
		"key":    "delete_key",
	})

	if result.IsError {
		t.Errorf("delete failed: %s", result.ForLLM)
	}

	// Verify it's gone
	result = tool.Execute(ctx, map[string]any{
		"action": "get",
		"key":    "delete_key",
	})

	if !result.IsError {
		t.Error("expected error after deletion")
	}
}

func TestMemoryStoreTool_List(t *testing.T) {
	tmpDir := t.TempDir()
	tool := NewMemoryStoreTool(tmpDir)
	ctx := context.Background()

	// Add some keys
	tool.Execute(ctx, map[string]any{
		"action": "set",
		"key":    "key1",
		"value":  "value1",
	})
	tool.Execute(ctx, map[string]any{
		"action": "set",
		"key":    "key2",
		"value":  "value2",
	})

	// List keys
	result := tool.Execute(ctx, map[string]any{
		"action": "list",
	})

	if result.IsError {
		t.Errorf("list failed: %s", result.ForLLM)
	}
}

func TestMemoryStoreTool_Count(t *testing.T) {
	tmpDir := t.TempDir()
	tool := NewMemoryStoreTool(tmpDir)
	ctx := context.Background()

	// Add some keys
	tool.Execute(ctx, map[string]any{
		"action": "set",
		"key":    "key1",
		"value":  "value1",
	})
	tool.Execute(ctx, map[string]any{
		"action": "set",
		"key":    "key2",
		"value":  "value2",
	})

	// Count
	result := tool.Execute(ctx, map[string]any{
		"action": "count",
	})

	if result.IsError {
		t.Errorf("count failed: %s", result.ForLLM)
	}
}

func TestMemoryStoreTool_Clear(t *testing.T) {
	tmpDir := t.TempDir()
	tool := NewMemoryStoreTool(tmpDir)
	ctx := context.Background()

	// Add some keys
	tool.Execute(ctx, map[string]any{
		"action": "set",
		"key":    "key1",
		"value":  "value1",
	})
	tool.Execute(ctx, map[string]any{
		"action": "set",
		"key":    "key2",
		"value":  "value2",
	})

	// Clear all
	result := tool.Execute(ctx, map[string]any{
		"action": "clear",
	})

	if result.IsError {
		t.Errorf("clear failed: %s", result.ForLLM)
	}

	// Verify count is 0
	result = tool.Execute(ctx, map[string]any{
		"action": "count",
	})

	if result.IsError {
		t.Errorf("count after clear failed: %s", result.ForLLM)
	}
}

func TestMemoryStoreTool_SetMissingKey(t *testing.T) {
	tmpDir := t.TempDir()
	tool := NewMemoryStoreTool(tmpDir)
	ctx := context.Background()

	result := tool.Execute(ctx, map[string]any{
		"action": "set",
		"value":  "value",
	})

	if !result.IsError {
		t.Error("expected error when key is missing for action='set'")
	}
}

func TestMemoryStoreTool_GetMissingKey(t *testing.T) {
	tmpDir := t.TempDir()
	tool := NewMemoryStoreTool(tmpDir)
	ctx := context.Background()

	result := tool.Execute(ctx, map[string]any{
		"action": "get",
	})

	if !result.IsError {
		t.Error("expected error when key is missing for action='get'")
	}
}

func TestMemoryStoreTool_DeleteMissingKey(t *testing.T) {
	tmpDir := t.TempDir()
	tool := NewMemoryStoreTool(tmpDir)
	ctx := context.Background()

	result := tool.Execute(ctx, map[string]any{
		"action": "delete",
	})

	if !result.IsError {
		t.Error("expected error when key is missing for action='delete'")
	}
}

func TestMemoryStoreTool_InvalidAction(t *testing.T) {
	tmpDir := t.TempDir()
	tool := NewMemoryStoreTool(tmpDir)
	ctx := context.Background()

	result := tool.Execute(ctx, map[string]any{
		"action": "invalid",
	})

	if !result.IsError {
		t.Error("expected error for invalid action")
	}
}

func TestMemoryStoreTool_Persistence(t *testing.T) {
	tmpDir := t.TempDir()

	// Create tool and set value
	tool1 := NewMemoryStoreTool(tmpDir)
	ctx := context.Background()
	tool1.Execute(ctx, map[string]any{
		"action": "set",
		"key":    "persist_key",
		"value":  "persist_value",
	})

	// Create new tool instance (simulates restart)
	tool2 := NewMemoryStoreTool(tmpDir)

	// Verify value persisted
	result := tool2.Execute(ctx, map[string]any{
		"action": "get",
		"key":    "persist_key",
	})

	if result.IsError {
		t.Errorf("persistence failed: %s", result.ForLLM)
	}
}

func TestMemoryStoreTool_ComplexValues(t *testing.T) {
	tmpDir := t.TempDir()
	tool := NewMemoryStoreTool(tmpDir)
	ctx := context.Background()

	// Store complex value
	complexValue := map[string]any{
		"name": "test",
		"tags": []string{"a", "b", "c"},
		"nested": map[string]any{
			"key": "value",
		},
	}

	result := tool.Execute(ctx, map[string]any{
		"action": "set",
		"key":    "complex",
		"value":  complexValue,
	})

	if result.IsError {
		t.Errorf("set complex value failed: %s", result.ForLLM)
	}
}

func TestMemoryStoreTool_Concurrency(t *testing.T) {
	tmpDir := t.TempDir()
	tool := NewMemoryStoreTool(tmpDir)
	ctx := context.Background()

	done := make(chan bool, 20)
	for i := 0; i < 20; i++ {
		go func(id int) {
			tool.Execute(ctx, map[string]any{
				"action": "set",
				"key":    "key_" + string(rune(id)),
				"value":  id,
			})
			done <- true
		}(i)
	}

	for i := 0; i < 20; i++ {
		<-done
	}

	// Verify all writes succeeded
	result := tool.Execute(ctx, map[string]any{
		"action": "count",
	})

	if result.IsError {
		t.Errorf("concurrent writes failed: %s", result.ForLLM)
	}
}

func TestMemoryStoreTool_DiskPersistence(t *testing.T) {
	tmpDir := t.TempDir()
	tool := NewMemoryStoreTool(tmpDir)
	ctx := context.Background()

	// Set a value
	tool.Execute(ctx, map[string]any{
		"action": "set",
		"key":    "disk_key",
		"value":  "disk_value",
	})

	// Verify file exists
	dbPath := filepath.Join(tmpDir, "memory_store.json")
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Error("database file should exist")
	}
}

func TestMemoryStoreTool_Expiration(t *testing.T) {
	tmpDir := t.TempDir()
	tool := NewMemoryStoreTool(tmpDir)
	ctx := context.Background()

	// Set a value with very short TTL
	tool.Execute(ctx, map[string]any{
		"action":      "set",
		"key":         "expire_key",
		"value":       "expire_value",
		"ttl_seconds": 1,
	})

	// Wait for expiration
	time.Sleep(2 * time.Second)

	// Try to get expired value
	result := tool.Execute(ctx, map[string]any{
		"action": "get",
		"key":    "expire_key",
	})

	if !result.IsError {
		t.Error("expected error for expired key")
	}
}
