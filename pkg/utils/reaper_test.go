package utils

import (
	"runtime"
	"testing"
)

func TestFindOrphans_NoPanic(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("not supported on Windows")
	}
	orphans, err := FindOrphans()
	if err != nil {
		t.Fatalf("FindOrphans() error: %v", err)
	}
	// orphans may be nil or empty slice — both are valid when no orphans exist
	t.Logf("Found %d orphan(s)", len(orphans))
}
