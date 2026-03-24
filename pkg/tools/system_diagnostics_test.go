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

func TestSystemDiagnosticsTool_Name(t *testing.T) {
	tool := NewSystemDiagnosticsTool("/tmp/test-workspace")
	if tool.Name() != "system_diagnostics" {
		t.Errorf("expected name 'system_diagnostics', got '%s'", tool.Name())
	}
}

func TestSystemDiagnosticsTool_Description(t *testing.T) {
	tool := NewSystemDiagnosticsTool("/tmp/test-workspace")
	if tool.Description() == "" {
		t.Error("expected non-empty description")
	}
}

func TestSystemDiagnosticsTool_Parameters(t *testing.T) {
	tool := NewSystemDiagnosticsTool("/tmp/test-workspace")
	params := tool.Parameters()

	if params["type"] != "object" {
		t.Error("parameters type should be 'object'")
	}

	props, ok := params["properties"].(map[string]any)
	if !ok {
		t.Fatal("expected properties map")
	}

	if _, ok := props["metric"]; !ok {
		t.Error("missing 'metric' parameter")
	}
}

func TestSystemDiagnosticsTool_ExecuteCPU(t *testing.T) {
	tool := NewSystemDiagnosticsTool("/tmp/test-workspace")
	ctx := context.Background()

	result := tool.Execute(ctx, map[string]any{
		"metric": "cpu",
	})

	if result.IsError {
		t.Errorf("CPU diagnostics failed: %s", result.ForLLM)
	}
}

func TestSystemDiagnosticsTool_ExecuteRAM(t *testing.T) {
	tool := NewSystemDiagnosticsTool("/tmp/test-workspace")
	ctx := context.Background()

	result := tool.Execute(ctx, map[string]any{
		"metric": "ram",
	})

	if result.IsError {
		t.Errorf("RAM diagnostics failed: %s", result.ForLLM)
	}
}

func TestSystemDiagnosticsTool_ExecuteDisk(t *testing.T) {
	tool := NewSystemDiagnosticsTool("/tmp/test-workspace")
	ctx := context.Background()

	result := tool.Execute(ctx, map[string]any{
		"metric": "disk",
	})

	if result.IsError {
		t.Errorf("Disk diagnostics failed: %s", result.ForLLM)
	}
}

func TestSystemDiagnosticsTool_ExecuteProcesses(t *testing.T) {
	tool := NewSystemDiagnosticsTool("/tmp/test-workspace")
	ctx := context.Background()

	result := tool.Execute(ctx, map[string]any{
		"metric":  "processes",
		"limit":   5,
	})

	if result.IsError {
		t.Errorf("Processes diagnostics failed: %s", result.ForLLM)
	}
}

func TestSystemDiagnosticsTool_ExecuteAll(t *testing.T) {
	tool := NewSystemDiagnosticsTool("/tmp/test-workspace")
	ctx := context.Background()

	result := tool.Execute(ctx, map[string]any{
		"metric": "all",
	})

	if result.IsError {
		t.Errorf("All diagnostics failed: %s", result.ForLLM)
	}
}

func TestSystemDiagnosticsTool_ExecuteInvalidMetric(t *testing.T) {
	tool := NewSystemDiagnosticsTool("/tmp/test-workspace")
	ctx := context.Background()

	result := tool.Execute(ctx, map[string]any{
		"metric": "invalid",
	})

	if !result.IsError {
		t.Error("expected error for invalid metric")
	}
}

func TestSystemDiagnosticsTool_ExecuteMissingMetric(t *testing.T) {
	tool := NewSystemDiagnosticsTool("/tmp/test-workspace")
	ctx := context.Background()

	result := tool.Execute(ctx, map[string]any{})

	if !result.IsError {
		t.Error("expected error when metric is missing")
	}
}

func TestSystemDiagnosticsTool_ExecuteLogs(t *testing.T) {
	tool := NewSystemDiagnosticsTool("/tmp/test-workspace")
	ctx := context.Background()

	result := tool.Execute(ctx, map[string]any{
		"metric": "logs",
		"limit":  10,
	})

	// Log reading may return a message about privileges
	if result.IsError {
		t.Logf("Log reading returned error (expected): %s", result.ForLLM)
	}
}

func TestSystemDiagnosticsTool_Concurrency(t *testing.T) {
	tool := NewSystemDiagnosticsTool("/tmp/test-workspace")
	ctx := context.Background()

	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			_ = tool.Execute(ctx, map[string]any{"metric": "cpu"})
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}
