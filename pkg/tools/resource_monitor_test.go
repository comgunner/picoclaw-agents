// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package tools

import (
	"context"
	"testing"
)

func TestResourceMonitorTool_Name(t *testing.T) {
	tool := NewResourceMonitorTool("/tmp/test-workspace")
	if tool.Name() != "resource_monitor" {
		t.Errorf("expected name 'resource_monitor', got '%s'", tool.Name())
	}
}

func TestResourceMonitorTool_Description(t *testing.T) {
	tool := NewResourceMonitorTool("/tmp/test-workspace")
	if tool.Description() == "" {
		t.Error("expected non-empty description")
	}
}

func TestResourceMonitorTool_Parameters(t *testing.T) {
	tool := NewResourceMonitorTool("/tmp/test-workspace")
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

func TestResourceMonitorTool_ExecuteCurrent(t *testing.T) {
	tool := NewResourceMonitorTool("/tmp/test-workspace")
	ctx := context.Background()

	result := tool.Execute(ctx, map[string]any{
		"action": "current",
	})

	if result.IsError {
		t.Errorf("current usage failed: %s", result.ForLLM)
	}
}

func TestResourceMonitorTool_ExecuteCPUThreshold(t *testing.T) {
	tool := NewResourceMonitorTool("/tmp/test-workspace")
	ctx := context.Background()

	result := tool.Execute(ctx, map[string]any{
		"action":    "cpu_threshold",
		"threshold": 80.0,
	})

	if result.IsError {
		t.Errorf("CPU threshold check failed: %s", result.ForLLM)
	}
}

func TestResourceMonitorTool_ExecuteRAMThreshold(t *testing.T) {
	tool := NewResourceMonitorTool("/tmp/test-workspace")
	ctx := context.Background()

	result := tool.Execute(ctx, map[string]any{
		"action":    "ram_threshold",
		"threshold": 80.0,
	})

	if result.IsError {
		t.Errorf("RAM threshold check failed: %s", result.ForLLM)
	}
}

func TestResourceMonitorTool_ExecuteThrottle(t *testing.T) {
	tool := NewResourceMonitorTool("/tmp/test-workspace")
	ctx := context.Background()

	result := tool.Execute(ctx, map[string]any{
		"action":    "throttle",
		"threshold": 80.0,
	})

	if result.IsError {
		t.Errorf("throttle recommendation failed: %s", result.ForLLM)
	}
}

func TestResourceMonitorTool_ExecuteHistory(t *testing.T) {
	tool := NewResourceMonitorTool("/tmp/test-workspace")
	ctx := context.Background()

	// First, record some data
	tool.Execute(ctx, map[string]any{"action": "current"})

	result := tool.Execute(ctx, map[string]any{
		"action": "history",
		"hours":  1,
	})

	if result.IsError {
		t.Errorf("history retrieval failed: %s", result.ForLLM)
	}
}

func TestResourceMonitorTool_ExecuteInvalidAction(t *testing.T) {
	tool := NewResourceMonitorTool("/tmp/test-workspace")
	ctx := context.Background()

	result := tool.Execute(ctx, map[string]any{
		"action": "invalid",
	})

	if !result.IsError {
		t.Error("expected error for invalid action")
	}
}

func TestResourceMonitorTool_ExecuteMissingAction(t *testing.T) {
	tool := NewResourceMonitorTool("/tmp/test-workspace")
	ctx := context.Background()

	result := tool.Execute(ctx, map[string]any{})

	if !result.IsError {
		t.Error("expected error when action is missing")
	}
}

func TestResourceMonitorTool_ThresholdAlert(t *testing.T) {
	tool := NewResourceMonitorTool("/tmp/test-workspace")
	ctx := context.Background()

	// Set very low threshold to trigger alert
	result := tool.Execute(ctx, map[string]any{
		"action":    "cpu_threshold",
		"threshold": 0.001, // Almost certainly will exceed
	})

	// Should not error, may return alert message
	if result.IsError {
		t.Errorf("threshold check should not error: %s", result.ForLLM)
	}
}

func TestResourceMonitorTool_Concurrency(t *testing.T) {
	tool := NewResourceMonitorTool("/tmp/test-workspace")
	ctx := context.Background()

	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			_ = tool.Execute(ctx, map[string]any{"action": "current"})
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestResourceMonitorTool_HistoryRecording(t *testing.T) {
	tool := NewResourceMonitorTool("/tmp/test-workspace")
	ctx := context.Background()

	// Record multiple data points
	for i := 0; i < 5; i++ {
		tool.Execute(ctx, map[string]any{"action": "current"})
	}

	// Check history
	result := tool.Execute(ctx, map[string]any{
		"action": "history",
		"hours":  1,
	})

	if result.IsError {
		t.Errorf("history should contain recorded data: %s", result.ForLLM)
	}
}

func TestResourceMonitorTool_HistoryLimit(t *testing.T) {
	tool := NewResourceMonitorTool("/tmp/test-workspace")

	// Record more than maxHistory (100)
	for i := 0; i < 150; i++ {
		tool.recordUsage(float64(i), float64(i))
	}

	// Verify history is trimmed
	tool.historyMu.RLock()
	historyLen := len(tool.history)
	tool.historyMu.RUnlock()

	if historyLen > tool.maxHistory {
		t.Errorf("history should be limited to %d, got %d", tool.maxHistory, historyLen)
	}
}

func TestGetCPUUsagePercent(t *testing.T) {
	cpu := getCPUUsagePercent()
	if cpu < 0 || cpu > 100 {
		t.Errorf("CPU usage should be between 0 and 100, got %.2f", cpu)
	}
}

func TestGetRAMUsagePercent(t *testing.T) {
	ram := getRAMUsagePercent()
	if ram < 0 || ram > 100 {
		t.Errorf("RAM usage should be between 0 and 100, got %.2f", ram)
	}
}
