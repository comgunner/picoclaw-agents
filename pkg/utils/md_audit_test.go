package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAuditMarkdown_NoBrokenLinks(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "README.md"),
		[]byte("# Hello\nSee [docs](docs/guide.md)\n"), 0o644)
	os.MkdirAll(filepath.Join(dir, "docs"), 0o755)
	os.WriteFile(filepath.Join(dir, "docs", "guide.md"), []byte("# Guide\n"), 0o644)

	issues, err := AuditMarkdown(dir)
	if err != nil {
		t.Fatalf("AuditMarkdown error: %v", err)
	}
	if len(issues) != 0 {
		t.Errorf("expected 0 issues, got %d", len(issues))
	}
}

func TestAuditMarkdown_BrokenLinkDetected(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "README.md"),
		[]byte("# Hello\nSee [missing](nonexistent.md)\n"), 0o644)

	issues, err := AuditMarkdown(dir)
	if err != nil {
		t.Fatalf("AuditMarkdown error: %v", err)
	}
	if len(issues) == 0 {
		t.Error("expected broken link, got none")
	}
}
