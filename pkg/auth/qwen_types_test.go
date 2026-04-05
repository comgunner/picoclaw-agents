// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package auth

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// ============================================
// Tests para qwen_types.go
// ============================================

func TestQwenSession_IsExpired(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name        string
		session     *QwenSession
		wantExpired bool
	}{
		{
			name: "session not expired",
			session: &QwenSession{
				ExpiresAt: now.Add(24 * time.Hour),
			},
			wantExpired: false,
		},
		{
			name: "session expired yesterday",
			session: &QwenSession{
				ExpiresAt: now.Add(-24 * time.Hour),
			},
			wantExpired: true,
		},
		{
			name: "session expires in 30 minutes (within 1h margin)",
			session: &QwenSession{
				ExpiresAt: now.Add(30 * time.Minute),
			},
			wantExpired: true,
		},
		{
			name: "session with zero ExpiresAt",
			session: &QwenSession{
				ExpiresAt: time.Time{},
			},
			wantExpired: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.session.IsExpired()
			if got != tt.wantExpired {
				t.Errorf("IsExpired() = %v, want %v", got, tt.wantExpired)
			}
		})
	}
}

func TestQwenSession_NeedsRefresh(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name        string
		session     *QwenSession
		wantRefresh bool
	}{
		{
			name: "session valid for 10 days",
			session: &QwenSession{
				ExpiresAt: now.Add(10 * 24 * time.Hour),
			},
			wantRefresh: false,
		},
		{
			name: "session expires in 2 days (within 3 day margin)",
			session: &QwenSession{
				ExpiresAt: now.Add(2 * 24 * time.Hour),
			},
			wantRefresh: true,
		},
		{
			name: "session already expired",
			session: &QwenSession{
				ExpiresAt: now.Add(-24 * time.Hour),
			},
			wantRefresh: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.session.NeedsRefresh()
			if got != tt.wantRefresh {
				t.Errorf("NeedsRefresh() = %v, want %v", got, tt.wantRefresh)
			}
		})
	}
}

func TestQwenSession_Validate(t *testing.T) {
	tests := []struct {
		name    string
		session *QwenSession
		wantErr bool
		errType error
	}{
		{
			name: "valid session",
			session: &QwenSession{
				AccessToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test",
				Email:       "user@example.com",
				ExpiresAt:   time.Now().Add(24 * time.Hour),
			},
			wantErr: false,
		},
		{
			name: "missing access token",
			session: &QwenSession{
				Email:     "user@example.com",
				ExpiresAt: time.Now().Add(24 * time.Hour),
			},
			wantErr: true,
			errType: ErrMissingAccessToken,
		},
		{
			name: "missing email",
			session: &QwenSession{
				AccessToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test",
				ExpiresAt:   time.Now().Add(24 * time.Hour),
			},
			wantErr: true,
			errType: ErrMissingEmail,
		},
		{
			name: "missing expires at",
			session: &QwenSession{
				AccessToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test",
				Email:       "user@example.com",
			},
			wantErr: true,
			errType: ErrMissingExpiresAt,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.session.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.errType != nil && err != tt.errType {
				t.Errorf("Validate() error type = %v, want %v", err, tt.errType)
			}
		})
	}
}

func TestQwenError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  *QwenError
		want string
	}{
		{
			name: "error without wrapped error",
			err:  &QwenError{Message: "test error"},
			want: "test error",
		},
		{
			name: "error with wrapped error",
			err: &QwenError{
				Message: "outer error",
				Err:     os.ErrNotExist,
			},
			want: "outer error: file does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			if got != tt.want {
				t.Errorf("Error() = %q, want %q", got, tt.want)
			}
		})
	}
}

// ============================================
// Tests para qwen_session.go
// ============================================

func TestSessionManager_SaveLoad(t *testing.T) {
	// Crear directorio temporal para tests
	tmpDir := t.TempDir()
	sessionPath := filepath.Join(tmpDir, "state", "qwen_session.json")

	sm := &SessionManager{
		sessionPath: sessionPath,
		workspace:   tmpDir,
	}

	session := &QwenSession{
		AccessToken:  "test_access_token",
		RefreshToken: "test_refresh_token",
		Email:        "test@example.com",
		Plan:         "free",
		DailyLimit:   2000,
		ExpiresAt:    time.Now().Add(24 * time.Hour),
		LastVerified: time.Now(),
	}

	// Test Save
	err := sm.Save(session)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Verificar que el archivo existe
	if _, statErr := os.Stat(sessionPath); os.IsNotExist(statErr) {
		t.Fatal("Save() did not create session file")
	}

	// Verificar permisos (0600)
	info, err := os.Stat(sessionPath)
	if err != nil {
		t.Fatalf("Stat() error = %v", err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Errorf("File permissions = %o, want %o", info.Mode().Perm(), 0o600)
	}

	// Test Load
	loaded, err := sm.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if loaded == nil {
		t.Fatal("Load() returned nil")
	}

	// Verificar campos cargados
	if loaded.AccessToken != session.AccessToken {
		t.Errorf("AccessToken = %q, want %q", loaded.AccessToken, session.AccessToken)
	}
	if loaded.Email != session.Email {
		t.Errorf("Email = %q, want %q", loaded.Email, session.Email)
	}
	if loaded.Plan != session.Plan {
		t.Errorf("Plan = %q, want %q", loaded.Plan, session.Plan)
	}
}

func TestSessionManager_LoadNonExistent(t *testing.T) {
	tmpDir := t.TempDir()
	sm := NewSessionManager(tmpDir)

	session, err := sm.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if session != nil {
		t.Error("Load() should return nil for non-existent session")
	}
}

func TestSessionManager_Delete(t *testing.T) {
	tmpDir := t.TempDir()
	sm := NewSessionManager(tmpDir)

	session := &QwenSession{
		AccessToken: "test",
		Email:       "test@example.com",
		ExpiresAt:   time.Now().Add(24 * time.Hour),
	}

	// Guardar sesión
	if err := sm.Save(session); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Eliminar sesión
	if err := sm.Delete(); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	// Verificar que fue eliminada
	if sm.Exists() {
		t.Error("Session should be deleted")
	}

	// Delete en sesión no existente no debe fallar
	if err := sm.Delete(); err != nil {
		t.Errorf("Delete() on non-existent session should not error: %v", err)
	}
}

func TestSessionManager_Exists(t *testing.T) {
	tmpDir := t.TempDir()
	sm := NewSessionManager(tmpDir)

	// Inicialmente no debe existir
	if sm.Exists() {
		t.Error("Session should not exist initially")
	}

	// Guardar sesión
	session := &QwenSession{
		AccessToken: "test",
		Email:       "test@example.com",
		ExpiresAt:   time.Now().Add(24 * time.Hour),
	}
	if err := sm.Save(session); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Ahora debe existir
	if !sm.Exists() {
		t.Error("Session should exist after Save()")
	}
}

func TestSessionInfo_RedactEmail(t *testing.T) {
	tests := []struct {
		name string
		info *SessionInfo
		want string
	}{
		{
			name: "normal email",
			info: &SessionInfo{Email: "user@example.com"},
			want: "us***@example.com",
		},
		{
			name: "short email",
			info: &SessionInfo{Email: "a@b.com"},
			want: "***",
		},
		{
			name: "empty email",
			info: &SessionInfo{Email: ""},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.info.RedactEmail()
			if got != tt.want {
				t.Errorf("RedactEmail() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSessionInfo_TimeUntilExpiry(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name string
		info *SessionInfo
		want time.Duration
	}{
		{
			name: "expires in 1 hour",
			info: &SessionInfo{ExpiresAt: now.Add(time.Hour)},
			want: time.Hour,
		},
		{
			name: "already expired",
			info: &SessionInfo{ExpiresAt: now.Add(-time.Hour), IsExpired: true},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.info.TimeUntilExpiry()
			// Allow 1 second tolerance for timing
			if got < tt.want-time.Second || got > tt.want+time.Second {
				t.Errorf("TimeUntilExpiry() = %v, want %v", got, tt.want)
			}
		})
	}
}

// ============================================
// Tests para qwen_oauth.go
// ============================================

func TestQwenOAuth_extractOAuthInfo(t *testing.T) {
	tmpDir := t.TempDir()
	oauth := NewQwenOAuth(tmpDir)

	tests := []struct {
		name        string
		output      string
		wantURL     bool
		wantCode    string
		wantSuccess bool
	}{
		{
			name: "valid output with URL and code",
			output: `
🔐 Starting Qwen Portal OAuth authentication...
🔗 OAuth Link: https://chat.qwen.ai/authorize?user_code=M17WU0SC&client=qwen-code
📱 Device Code: M17WU0SC
			`,
			wantURL:     true,
			wantCode:    "M17WU0SC",
			wantSuccess: true,
		},
		{
			name: "output without URL",
			output: `
🔐 Starting Qwen Portal OAuth authentication...
Waiting for output...
			`,
			wantURL:     false,
			wantCode:    "",
			wantSuccess: false,
		},
		{
			name: "output with invalid code format (too short)",
			output: `
🔗 OAuth Link: https://chat.qwen.ai/authorize?user_code=ABC123
📱 Device Code: ABC123
			`,
			wantURL:     true,
			wantCode:    "", // Invalid code (not 8 chars)
			wantSuccess: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := oauth.extractOAuthInfo(tt.output)

			if result.Success != tt.wantSuccess {
				t.Errorf("Success = %v, want %v", result.Success, tt.wantSuccess)
			}

			hasURL := result.AuthURL != ""
			if hasURL != tt.wantURL {
				t.Errorf("AuthURL present = %v, want %v", hasURL, tt.wantURL)
			}

			if result.DeviceCode != tt.wantCode {
				t.Errorf("DeviceCode = %q, want %q", result.DeviceCode, tt.wantCode)
			}
		})
	}
}

func TestDependencyError_Error(t *testing.T) {
	err := &DependencyError{
		Name:    "tmux",
		Message: "tmux not found",
		Install: "brew install tmux",
	}

	got := err.Error()
	want := "tmux not found: brew install tmux"

	if got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}
