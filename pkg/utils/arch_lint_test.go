package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckImports_NoViolations(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "pkg/agent"), 0o755)
	os.MkdirAll(filepath.Join(dir, "pkg/channels"), 0o755)
	os.WriteFile(filepath.Join(dir, "pkg/agent", "agent.go"),
		[]byte(`package agent; import "fmt"; func Hello() { fmt.Println("hi") }`), 0o644)
	os.WriteFile(filepath.Join(dir, "pkg/channels", "ch.go"),
		[]byte(`package channels; import "fmt"; func Start() { fmt.Println("start") }`), 0o644)

	violations, err := CheckImports(dir, nil)
	if err != nil {
		t.Fatalf("CheckImports error: %v", err)
	}
	if len(violations) != 0 {
		t.Errorf("expected 0 violations, got %d", len(violations))
	}
}

func TestCheckImports_ViolationDetected(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "pkg/agent"), 0o755)
	os.WriteFile(filepath.Join(dir, "pkg/agent", "bad.go"),
		[]byte(`package agent; import "github.com/comgunner/picoclaw/pkg/channels"`), 0o644)

	violations, err := CheckImports(dir, nil)
	if err != nil {
		t.Fatalf("CheckImports error: %v", err)
	}
	if len(violations) == 0 {
		t.Error("expected violation, got none")
	}
}
