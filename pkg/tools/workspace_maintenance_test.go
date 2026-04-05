package tools

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestWorkspaceMaintenanceTool_DryRunDoesNotDeleteFiles(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("unix-specific assertion for dry-run execution")
	}

	workspace := t.TempDir()
	sessionsDir := filepath.Join(workspace, "sessions")
	if err := os.MkdirAll(sessionsDir, 0o755); err != nil {
		t.Fatalf("mkdir sessions: %v", err)
	}

	oldSession := filepath.Join(sessionsDir, "heartbeat-old.json")
	if err := os.WriteFile(oldSession, []byte(`{"alive":true}`), 0o644); err != nil {
		t.Fatalf("write old session: %v", err)
	}
	oldTime := time.Now().Add(-10 * 24 * time.Hour)
	if err := os.Chtimes(oldSession, oldTime, oldTime); err != nil {
		t.Fatalf("chtime old session: %v", err)
	}

	tool := NewWorkspaceMaintenanceTool(workspace)
	result := tool.Execute(context.Background(), map[string]any{
		"dry_run":        true,
		"clean_sessions": true,
	})

	if result.IsError {
		t.Fatalf("unexpected error: %s", result.ForLLM)
	}
	if _, err := os.Stat(oldSession); err != nil {
		t.Fatalf("dry run removed file unexpectedly: %v", err)
	}
}

func TestWorkspaceMaintenanceTool_ParseOutput(t *testing.T) {
	tool := NewWorkspaceMaintenanceTool("/tmp/workspace")
	output := strings.Join([]string{
		"SESSIONS_ARCHIVED: 5",
		"LOGS_COMPRESSED: 3",
		"TEMP_FILES_DELETED: 10",
		"SPACE_FREED_BYTES: 1048576",
	}, "\n")

	got := tool.parseOutput(output, "/tmp/workspace", 120, "linux")
	if got.SessionsArchived != 5 {
		t.Fatalf("sessions: want 5, got %d", got.SessionsArchived)
	}
	if got.LogsCompressed != 3 {
		t.Fatalf("logs: want 3, got %d", got.LogsCompressed)
	}
	if got.TempFilesDeleted != 10 {
		t.Fatalf("temp: want 10, got %d", got.TempFilesDeleted)
	}
	if got.SpaceFreedBytes != 1048576 {
		t.Fatalf("space: want 1048576, got %d", got.SpaceFreedBytes)
	}
}

func TestFormatBytes(t *testing.T) {
	cases := []struct {
		in   int64
		want string
	}{
		{0, "0 bytes"},
		{512, "512 bytes"},
		{1024, "1.00 KB"},
		{1536, "1.50 KB"},
		{1024 * 1024, "1.00 MB"},
		{1024 * 1024 * 1024, "1.00 GB"},
	}

	for _, tc := range cases {
		got := FormatBytes(tc.in)
		if got != tc.want {
			t.Fatalf("FormatBytes(%d): want %q, got %q", tc.in, tc.want, got)
		}
	}
}
