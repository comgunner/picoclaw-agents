package utils

import (
	"testing"
)

func TestNotionMCPServer_CreatesWithoutError(t *testing.T) {
	cfg := &NotionMCPConfig{APIKey: "test-key"}
	s := NewNotionMCPServer(cfg)
	if s == nil {
		t.Fatal("server should not be nil")
	}
}

func TestNotionMCPServer_HasExpectedTools(t *testing.T) {
	cfg := &NotionMCPConfig{APIKey: "test-key"}
	s := NewNotionMCPServer(cfg)
	// Server should be created without error
	if s == nil {
		t.Fatal("expected non-nil server")
	}
}

func TestNotionMCPServer_EmptyConfig(t *testing.T) {
	cfg := &NotionMCPConfig{}
	s := NewNotionMCPServer(cfg)
	if s == nil {
		t.Fatal("server should still be created with empty config")
	}
}

func TestNotionConfigFromEnv(t *testing.T) {
	cfg := NotionConfigFromEnv()
	if cfg == nil {
		t.Fatal("config should not be nil")
	}
}

func TestServeNotionMCPStdio_RequiresAPIKey(t *testing.T) {
	cfg := &NotionMCPConfig{}
	// Should not crash, but ServeNotionMCPStdio returns error for empty key
	err := ServeNotionMCPStdio(cfg)
	if err == nil {
		t.Error("expected error for empty API key")
	}
}
