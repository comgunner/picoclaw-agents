// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// SessionManager gestiona el ciclo de vida de sesiones Qwen
// Persiste sesiones en ~/.picoclaw/workspace/state/qwen_session.json
type SessionManager struct {
	sessionPath string
	workspace   string
}

// NewSessionManager crea un gestor de sesiones Qwen
// workspace: directorio base del workspace (ej: ~/.picoclaw/workspace)
func NewSessionManager(workspace string) *SessionManager {
	return &SessionManager{
		sessionPath: filepath.Join(workspace, "state", "qwen_session.json"),
		workspace:   workspace,
	}
}

// Save persiste la sesión en disco con atomic write
// Usa el patrón: temp file + fsync + rename para prevenir corrupción
func (m *SessionManager) Save(session *QwenSession) error {
	if session == nil {
		return fmt.Errorf("session cannot be nil")
	}

	// Validar sesión antes de persistir
	if err := session.Validate(); err != nil {
		return fmt.Errorf("invalid session: %w", err)
	}

	// Crear directorio si no existe
	dir := filepath.Dir(m.sessionPath)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("failed to create state directory: %w", err)
	}

	// Serializar a JSON con indentación
	data, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	// Atomic write: temp file + fsync + rename
	tmpPath := m.sessionPath + ".tmp"

	// Escribir archivo temporal con permisos restrictivos (0600)
	if writeErr := os.WriteFile(tmpPath, data, 0o600); writeErr != nil {
		return fmt.Errorf("failed to write temp file: %w", writeErr)
	}

	// Fsync para garantizar persistencia en disco
	f, openErr := os.OpenFile(tmpPath, os.O_RDWR, 0o600)
	if openErr != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("failed to open temp file: %w", openErr)
	}

	if syncErr := f.Sync(); syncErr != nil {
		f.Close()
		os.Remove(tmpPath)
		return fmt.Errorf("failed to sync file: %w", syncErr)
	}

	if closeErr := f.Close(); closeErr != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("failed to close file: %w", closeErr)
	}

	// Atomic rename
	if renameErr := os.Rename(tmpPath, m.sessionPath); renameErr != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("failed to rename file: %w", renameErr)
	}

	// Verificar permisos después de rename (algunos sistemas pueden cambiarlos)
	if permErr := os.Chmod(m.sessionPath, 0o600); permErr != nil {
		return fmt.Errorf("failed to set file permissions: %w", permErr)
	}

	return nil
}

// Load carga la sesión desde disco
// Retorna nil, nil si no existe sesión
func (m *SessionManager) Load() (*QwenSession, error) {
	data, err := os.ReadFile(m.sessionPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // No session exists
		}
		return nil, fmt.Errorf("failed to read session file: %w", err)
	}

	var session QwenSession
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}

	return &session, nil
}

// Delete elimina la sesión (logout)
// No retorna error si el archivo no existe
func (m *SessionManager) Delete() error {
	if _, err := os.Stat(m.sessionPath); os.IsNotExist(err) {
		return nil
	}
	return os.Remove(m.sessionPath)
}

// Exists verifica si existe una sesión persistida
func (m *SessionManager) Exists() bool {
	_, err := os.Stat(m.sessionPath)
	return err == nil
}

// GetSessionPath retorna la ruta completa al archivo de sesión
func (m *SessionManager) GetSessionPath() string {
	return m.sessionPath
}

// GetSessionInfo retorna información sobre la sesión sin cargarla completa
// Útil para logging y debugging sin exponer tokens
func (m *SessionManager) GetSessionInfo() (*SessionInfo, error) {
	session, err := m.Load()
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, nil
	}

	return &SessionInfo{
		Email:        session.Email,
		Plan:         session.Plan,
		DailyLimit:   session.DailyLimit,
		ExpiresAt:    session.ExpiresAt,
		IsExpired:    session.IsExpired(),
		NeedsRefresh: session.NeedsRefresh(),
		LastVerified: session.LastVerified,
	}, nil
}

// SessionInfo información pública de la sesión (sin tokens)
type SessionInfo struct {
	Email        string    `json:"email"`
	Plan         string    `json:"plan"`
	DailyLimit   int       `json:"daily_limit"`
	ExpiresAt    time.Time `json:"expires_at"`
	IsExpired    bool      `json:"is_expired"`
	NeedsRefresh bool      `json:"needs_refresh"`
	LastVerified time.Time `json:"last_verified"`
}

// RedactEmail returns a redacted version of the email for logging.
// Example: "user@example.com" → "us***@example.com"
func (s *SessionInfo) RedactEmail() string {
	if s.Email == "" {
		return ""
	}

	parts := filepath.SplitList(s.Email)
	if len(parts) == 0 {
		return "***"
	}

	email := s.Email
	atIndex := -1
	for i, c := range email {
		if c == '@' {
			atIndex = i
			break
		}
	}

	if atIndex <= 2 {
		return "***"
	}

	return email[:2] + "***" + email[atIndex:]
}

// TimeUntilExpiry retorna el tiempo restante hasta expiración
// Retorna 0 si ya expiró
func (s *SessionInfo) TimeUntilExpiry() time.Duration {
	if s.IsExpired {
		return 0
	}
	return time.Until(s.ExpiresAt)
}

// DaysUntilExpiry retorna los días restantes hasta expiración
// Retorna 0 si ya expiró
func (s *SessionInfo) DaysUntilExpiry() int {
	duration := s.TimeUntilExpiry()
	return int(duration.Hours() / 24)
}
