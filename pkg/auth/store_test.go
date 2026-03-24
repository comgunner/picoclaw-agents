package auth

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestAuthCredentialIsExpired(t *testing.T) {
	tests := []struct {
		name      string
		expiresAt time.Time
		want      bool
	}{
		{"zero time", time.Time{}, false},
		{"future", time.Now().Add(time.Hour), false},
		{"past", time.Now().Add(-time.Hour), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &AuthCredential{ExpiresAt: tt.expiresAt}
			if got := c.IsExpired(); got != tt.want {
				t.Errorf("IsExpired() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthCredentialNeedsRefresh(t *testing.T) {
	tests := []struct {
		name      string
		expiresAt time.Time
		want      bool
	}{
		{"zero time", time.Time{}, false},
		{"far future", time.Now().Add(time.Hour), false},
		{"within 5 min", time.Now().Add(3 * time.Minute), true},
		{"already expired", time.Now().Add(-time.Minute), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &AuthCredential{ExpiresAt: tt.expiresAt}
			if got := c.NeedsRefresh(); got != tt.want {
				t.Errorf("NeedsRefresh() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStoreRoundtrip(t *testing.T) {
	tmpDir := t.TempDir()
	origHome := os.Getenv("HOME")
	t.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	cred := &AuthCredential{
		AccessToken:  "test-access-token",
		RefreshToken: "test-refresh-token",
		AccountID:    "acct-123",
		ExpiresAt:    time.Now().Add(time.Hour).Truncate(time.Second),
		Provider:     "openai",
		AuthMethod:   "oauth",
	}

	if err := SetCredential("openai", cred); err != nil {
		t.Fatalf("SetCredential() error: %v", err)
	}

	loaded, err := GetCredential("openai")
	if err != nil {
		t.Fatalf("GetCredential() error: %v", err)
	}
	if loaded == nil {
		t.Fatal("GetCredential() returned nil")
	}
	if loaded.AccessToken != cred.AccessToken {
		t.Errorf("AccessToken = %q, want %q", loaded.AccessToken, cred.AccessToken)
	}
	if loaded.RefreshToken != cred.RefreshToken {
		t.Errorf("RefreshToken = %q, want %q", loaded.RefreshToken, cred.RefreshToken)
	}
	if loaded.Provider != cred.Provider {
		t.Errorf("Provider = %q, want %q", loaded.Provider, cred.Provider)
	}
}

func TestStoreFilePermissions(t *testing.T) {
	tmpDir := t.TempDir()
	origHome := os.Getenv("HOME")
	t.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	cred := &AuthCredential{
		AccessToken: "secret-token",
		Provider:    "openai",
		AuthMethod:  "oauth",
	}
	if err := SetCredential("openai", cred); err != nil {
		t.Fatalf("SetCredential() error: %v", err)
	}

	path := filepath.Join(tmpDir, ".picoclaw", "auth.json")
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat() error: %v", err)
	}
	perm := info.Mode().Perm()
	if perm != 0o600 {
		t.Errorf("file permissions = %o, want 0600", perm)
	}
}

func TestStoreMultiProvider(t *testing.T) {
	tmpDir := t.TempDir()
	origHome := os.Getenv("HOME")
	t.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	openaiCred := &AuthCredential{AccessToken: "openai-token", Provider: "openai", AuthMethod: "oauth"}
	anthropicCred := &AuthCredential{AccessToken: "anthropic-token", Provider: "anthropic", AuthMethod: "token"}

	if err := SetCredential("openai", openaiCred); err != nil {
		t.Fatalf("SetCredential(openai) error: %v", err)
	}
	if err := SetCredential("anthropic", anthropicCred); err != nil {
		t.Fatalf("SetCredential(anthropic) error: %v", err)
	}

	loaded, err := GetCredential("openai")
	if err != nil {
		t.Fatalf("GetCredential(openai) error: %v", err)
	}
	if loaded.AccessToken != "openai-token" {
		t.Errorf("openai token = %q, want %q", loaded.AccessToken, "openai-token")
	}

	loaded, err = GetCredential("anthropic")
	if err != nil {
		t.Fatalf("GetCredential(anthropic) error: %v", err)
	}
	if loaded.AccessToken != "anthropic-token" {
		t.Errorf("anthropic token = %q, want %q", loaded.AccessToken, "anthropic-token")
	}
}

func TestDeleteCredential(t *testing.T) {
	tmpDir := t.TempDir()
	origHome := os.Getenv("HOME")
	t.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	cred := &AuthCredential{AccessToken: "to-delete", Provider: "openai", AuthMethod: "oauth"}
	if err := SetCredential("openai", cred); err != nil {
		t.Fatalf("SetCredential() error: %v", err)
	}

	if err := DeleteCredential("openai"); err != nil {
		t.Fatalf("DeleteCredential() error: %v", err)
	}

	loaded, err := GetCredential("openai")
	if err != nil {
		t.Fatalf("GetCredential() error: %v", err)
	}
	if loaded != nil {
		t.Error("expected nil after delete")
	}
}

func TestLoadStoreEmpty(t *testing.T) {
	tmpDir := t.TempDir()
	origHome := os.Getenv("HOME")
	t.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	store, err := LoadStore()
	if err != nil {
		t.Fatalf("LoadStore() error: %v", err)
	}
	if store == nil {
		t.Fatal("LoadStore() returned nil")
	}
	if len(store.Credentials) != 0 {
		t.Errorf("expected empty credentials, got %d", len(store.Credentials))
	}
}

// TestExpiredTokenShouldAttemptRefresh validates that the combined condition
// (NeedsRefresh || IsExpired) used in the Antigravity provider covers tokens
// that expired during periods of inactivity.
func TestExpiredTokenShouldAttemptRefresh(t *testing.T) {
	tests := []struct {
		name       string
		expiresAt  time.Time
		hasRefresh bool
		wantRetry  bool
	}{
		{"fresh token", time.Now().Add(time.Hour), true, false},
		{"about to expire", time.Now().Add(3 * time.Minute), true, true},
		{"just expired", time.Now().Add(-5 * time.Minute), true, true},
		{"expired hours ago", time.Now().Add(-3 * time.Hour), true, true},
		{"expired but no refresh_token", time.Now().Add(-1 * time.Hour), false, false},
		{"zero time (never expires)", time.Time{}, true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cred := &AuthCredential{
				ExpiresAt:    tt.expiresAt,
				RefreshToken: "",
			}
			if tt.hasRefresh {
				cred.RefreshToken = "some-refresh-token"
			}

			// This is the exact condition from createAntigravityTokenSource()
			shouldRetry := (cred.NeedsRefresh() || cred.IsExpired()) && cred.RefreshToken != ""
			if shouldRetry != tt.wantRetry {
				t.Errorf("shouldRetry = %v, want %v (NeedsRefresh=%v, IsExpired=%v, hasRefresh=%v)",
					shouldRetry, tt.wantRetry, cred.NeedsRefresh(), cred.IsExpired(), tt.hasRefresh)
			}
		})
	}
}
