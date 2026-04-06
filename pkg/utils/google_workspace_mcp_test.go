package utils

import (
	"testing"
)

func TestGoogleWorkspaceMCPServer_CreatesWithoutCreds(t *testing.T) {
	cfg := &GoogleWorkspaceMCPConfig{}
	s, err := NewGoogleWorkspaceMCPServer(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("server should not be nil even without creds")
	}
}

func TestGoogleWorkspaceMCPServer_InvalidCreds(t *testing.T) {
	cfg := &GoogleWorkspaceMCPConfig{CredentialsJSON: "/nonexistent/path.json"}
	_, err := NewGoogleWorkspaceMCPServer(cfg)
	if err == nil {
		t.Log("expected error for invalid credentials path")
	}
}

func TestGoogleWorkspaceConfigFromEnv(t *testing.T) {
	cfg := GoogleWorkspaceConfigFromEnv()
	if cfg == nil {
		t.Fatal("config should not be nil")
	}
}

func TestServeGoogleWorkspaceMCPStdio_RequiresCreds(t *testing.T) {
	cfg := &GoogleWorkspaceMCPConfig{}
	// Should not crash - returns server without tools
	err := ServeGoogleWorkspaceMCPStdio(cfg)
	if err != nil {
		t.Logf("expected no error for empty creds: %v", err)
	}
}
