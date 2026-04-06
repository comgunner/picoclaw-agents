package utils

import (
	"testing"
)

func TestSocialMediaMCPServer_CreatesWithoutError(t *testing.T) {
	cfg := &SocialMediaMCPConfig{}
	s := NewSocialMediaMCPServer(cfg)
	if s == nil {
		t.Fatal("server should not be nil")
	}
}

func TestSocialMediaMCPServer_WithCredentials(t *testing.T) {
	cfg := &SocialMediaMCPConfig{
		FacebookPageID:    "123456",
		FacebookPageToken: "test-token",
		XAPIKey:           "test-key",
		XAPISecret:        "test-secret",
		XAccessToken:      "test-access",
		XAccessSecret:     "test-access-secret",
	}
	s := NewSocialMediaMCPServer(cfg)
	if s == nil {
		t.Fatal("server should not be nil with credentials")
	}
}

func TestSocialMediaMCPServer_ConfigFromEnv(t *testing.T) {
	cfg := SocialMediaConfigFromEnv()
	// Should not panic even if no env vars are set
	if cfg == nil {
		t.Fatal("config should not be nil")
	}
}

func TestSocialMediaConfigFromEnv_EmptyEnv(t *testing.T) {
	// Even with no env vars, should return valid config struct
	cfg := SocialMediaConfigFromEnv()
	if cfg.FacebookPageID != "" {
		t.Logf("FacebookPageID found in env: %s", cfg.FacebookPageID)
	}
}
