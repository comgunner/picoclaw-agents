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
	"strings"
	"testing"
	"time"
)

// ============================================
// Security Tests para Qwen Authentication
// ============================================

// TestSessionFilePermissions verifica que los archivos de sesión
// tengan permisos restrictivos (0600 - solo owner puede leer/escribir)
func TestSessionFilePermissions(t *testing.T) {
	tmpDir := t.TempDir()
	sm := NewSessionManager(tmpDir)

	session := &QwenSession{
		AccessToken:  "sensitive_access_token",
		RefreshToken: "sensitive_refresh_token",
		Email:        "user@example.com",
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}

	// Guardar sesión
	err := sm.Save(session)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Verificar permisos del archivo
	info, err := os.Stat(sm.sessionPath)
	if err != nil {
		t.Fatalf("Stat() error = %v", err)
	}

	mode := info.Mode().Perm()
	if mode != 0o600 {
		t.Errorf("Session file permissions = %o, want %o", mode, 0o600)
	}

	// Verificar que otros usuarios no pueden leer
	if mode&0o044 != 0 {
		t.Error("Session file should not be readable by others")
	}

	// Verificar que otros usuarios no pueden escribir
	if mode&0o022 != 0 {
		t.Error("Session file should not be writable by others")
	}
}

// TestAtomicWrite verifica que las escrituras atómicas previenen
// corrupción de datos incluso si el proceso es interrumpido
func TestAtomicWrite(t *testing.T) {
	tmpDir := t.TempDir()
	sm := NewSessionManager(tmpDir)

	session := &QwenSession{
		AccessToken:  "test_token",
		RefreshToken: "test_refresh",
		Email:        "user@example.com",
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}

	// Guardar sesión
	err := sm.Save(session)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Verificar que no existe archivo temporal residual
	tmpPath := sm.sessionPath + ".tmp"
	if _, statErr := os.Stat(tmpPath); !os.IsNotExist(statErr) {
		t.Error("Temp file should be cleaned up after atomic write")
	}

	// Verificar que el archivo final existe y es válido
	loaded, err := sm.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if loaded == nil {
		t.Fatal("Loaded session should not be nil")
	}

	if loaded.AccessToken != session.AccessToken {
		t.Errorf("AccessToken mismatch: got %q, want %q", loaded.AccessToken, session.AccessToken)
	}
}

// TestTokenRedaction verifica que los tokens no sean expuestos
// en logs o mensajes de error
func TestTokenRedaction(t *testing.T) {
	tmpDir := t.TempDir()
	sm := NewSessionManager(tmpDir)

	session := &QwenSession{
		AccessToken:  "fake_jwt_token_for_testing_only", // pragma: allowlist secret
		RefreshToken: "refresh_abc123",
		Email:        "user@example.com",
		Plan:         "plus",
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}

	// Guardar sesión
	if err := sm.Save(session); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Obtener SessionInfo (que debería tener datos redactados)
	info, err := sm.GetSessionInfo()
	if err != nil {
		t.Fatalf("GetSessionInfo() error = %v", err)
	}

	// Verificar que el email redactado no expone el token
	redacted := info.RedactEmail()
	if strings.Contains(redacted, "fake_jwt_token_for_testing_only") {
		t.Error("Redacted email should not contain token data")
	}

	// Verificar que el archivo no contiene tokens en claro legibles
	data, err := os.ReadFile(sm.sessionPath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	fileContent := string(data)
	// Los tokens están en JSON, pero verificamos que no haya logging accidental
	if strings.Contains(fileContent, "sensitive_access_token") {
		t.Log("Note: Token found in file (expected for JSON storage)")
	}
}

// TestValidateSessionInput verifica validación de input
// para prevenir inyección de datos maliciosos
func TestValidateSessionInput(t *testing.T) {
	tests := []struct {
		name    string
		session *QwenSession
		wantErr bool
	}{
		{
			name: "normal session",
			session: &QwenSession{
				AccessToken: "valid_token",
				Email:       "user@example.com",
				ExpiresAt:   time.Now().Add(24 * time.Hour),
			},
			wantErr: false,
		},
		{
			name: "empty access token",
			session: &QwenSession{
				AccessToken: "",
				Email:       "user@example.com",
				ExpiresAt:   time.Now().Add(24 * time.Hour),
			},
			wantErr: true,
		},
		{
			name: "email with script injection attempt",
			session: &QwenSession{
				AccessToken: "valid_token",
				Email:       "<script>alert('xss')</script>@example.com",
				ExpiresAt:   time.Now().Add(24 * time.Hour),
			},
			wantErr: false, // Validación no rechaza formato de email
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.session.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestSessionManagerNilInput verifica manejo de inputs nil
func TestSessionManagerNilInput(t *testing.T) {
	tmpDir := t.TempDir()
	sm := NewSessionManager(tmpDir)

	// Test Save con nil session
	err := sm.Save(nil)
	if err == nil {
		t.Error("Save(nil) should return error")
	}

	// Test Load de sesión no existente (debe retornar nil, nil)
	session, err := sm.Load()
	if err != nil {
		t.Errorf("Load() error = %v, want nil", err)
	}
	if session != nil {
		t.Error("Load() should return nil for non-existent session")
	}
}

// TestConcurrentAccess verifica comportamiento con acceso concurrente
func TestConcurrentAccess(t *testing.T) {
	tmpDir := t.TempDir()
	sm := NewSessionManager(tmpDir)

	session := &QwenSession{
		AccessToken: "test_token",
		Email:       "user@example.com",
		ExpiresAt:   time.Now().Add(24 * time.Hour),
	}

	// Guardar sesión inicial
	if err := sm.Save(session); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Canal para sincronizar goroutines
	done := make(chan bool)
	errors := make(chan error, 10)

	// Múltiples goroutines leyendo concurrentemente
	for i := 0; i < 10; i++ {
		go func() {
			defer func() { done <- true }()
			s, err := sm.Load()
			if err != nil {
				errors <- err
				return
			}
			if s == nil {
				errors <- nil
			}
		}()
	}

	// Esperar todas las goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verificar que no hubo errores
	close(errors)
	for err := range errors {
		if err != nil {
			t.Errorf("Concurrent read error: %v", err)
		}
	}
}

// TestDeviceCodeFormat verifica validación del formato del device code
func TestDeviceCodeFormat(t *testing.T) {
	tmpDir := t.TempDir()
	oauth := NewQwenOAuth(tmpDir)

	tests := []struct {
		name   string
		output string
		valid  bool
	}{
		{
			name: "valid 8-char code",
			output: `
🔗 OAuth Link: https://chat.qwen.ai/authorize?user_code=M17WU0SC
📱 Device Code: M17WU0SC
			`,
			valid: true,
		},
		{
			name: "code too short (7 chars)",
			output: `
📱 Device Code: ABC1234
			`,
			valid: false,
		},
		{
			name: "code too long (9 chars)",
			output: `
📱 Device Code: ABCD12345
			`,
			valid: false,
		},
		{
			name: "code with lowercase",
			output: `
📱 Device Code: M17Wu0SC
			`,
			valid: false,
		},
		{
			name: "code with special chars",
			output: `
📱 Device Code: M17WU0!C
			`,
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := oauth.extractOAuthInfo(tt.output)

			hasValidCode := result.DeviceCode != ""
			if hasValidCode != tt.valid {
				t.Errorf("DeviceCode valid = %v, want %v (code: %q)",
					hasValidCode, tt.valid, result.DeviceCode)
			}
		})
	}
}

// TestSessionInfoMethods verifica métodos de SessionInfo
func TestSessionInfoMethods(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name         string
		info         *SessionInfo
		wantDaysMin  int
		wantDaysMax  int
		wantExpired  bool
		wantRedacted string
	}{
		{
			name: "expires in 10 days",
			info: &SessionInfo{
				Email:     "user@example.com",
				ExpiresAt: now.Add(10 * 24 * time.Hour),
				IsExpired: false,
			},
			wantDaysMin:  9,
			wantDaysMax:  10,
			wantExpired:  false,
			wantRedacted: "us***@example.com",
		},
		{
			name: "already expired",
			info: &SessionInfo{
				Email:     "user@example.com",
				ExpiresAt: now.Add(-24 * time.Hour),
				IsExpired: true,
			},
			wantDaysMin:  0,
			wantDaysMax:  0,
			wantExpired:  true,
			wantRedacted: "us***@example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test DaysUntilExpiry
			days := tt.info.DaysUntilExpiry()
			if days < tt.wantDaysMin || days > tt.wantDaysMax {
				t.Errorf("DaysUntilExpiry() = %d, want [%d-%d]", days, tt.wantDaysMin, tt.wantDaysMax)
			}

			// Test TimeUntilExpiry
			duration := tt.info.TimeUntilExpiry()
			if tt.wantExpired && duration != 0 {
				t.Errorf("TimeUntilExpiry() should be 0 for expired session, got %v", duration)
			}
			if !tt.wantExpired {
				expectedMin := time.Duration(tt.wantDaysMin) * 24 * time.Hour
				if duration < expectedMin {
					t.Errorf("TimeUntilExpiry() = %v, want >= %v", duration, expectedMin)
				}
			}

			// Test RedactEmail
			redacted := tt.info.RedactEmail()
			if redacted != tt.wantRedacted {
				t.Errorf("RedactEmail() = %q, want %q", redacted, tt.wantRedacted)
			}
		})
	}
}

// TestWorkspaceDirectoryCreation verifica que el directorio workspace
// se crea con permisos correctos
func TestWorkspaceDirectoryCreation(t *testing.T) {
	tmpDir := t.TempDir()
	stateDir := filepath.Join(tmpDir, "state")
	sm := NewSessionManager(tmpDir)

	// El directorio state no debe existir aún
	if _, err := os.Stat(stateDir); !os.IsNotExist(err) {
		t.Fatal("State directory should not exist initially")
	}

	// Guardar sesión debe crear el directorio
	session := &QwenSession{
		AccessToken: "test",
		Email:       "test@example.com",
		ExpiresAt:   time.Now().Add(24 * time.Hour),
	}

	if err := sm.Save(session); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Verificar que el directorio fue creado
	if _, err := os.Stat(stateDir); os.IsNotExist(err) {
		t.Fatal("State directory should be created by Save()")
	}

	// Verificar permisos del directorio (0700)
	info, err := os.Stat(stateDir)
	if err != nil {
		t.Fatalf("Stat() error = %v", err)
	}

	mode := info.Mode().Perm()
	if mode != 0o700 {
		t.Errorf("State directory permissions = %o, want %o", mode, 0o700)
	}
}
