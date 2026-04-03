// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package auth

import "time"

// QwenOAuthConfig returns the OAuth configuration for Qwen Portal (chat.qwen.ai).
// Uses PKCE flow for public client authentication.
func QwenOAuthConfig() OAuthProviderConfig {
	return OAuthProviderConfig{
		Issuer:     "https://chat.qwen.ai",
		TokenURL:   "https://chat.qwen.ai/api/auth/token",
		ClientID:   "qwen-code",
		Scopes:     "openid profile email offline_access",
		Originator: "picoclaw-agents",
		Port:       18891,
	}
}

// QwenAuthStatusResponse respuesta del endpoint /api/auth/status
type QwenAuthStatusResponse struct {
	Authenticated bool         `json:"authenticated"`
	Error         string       `json:"error,omitempty"`
	User          QwenUserInfo `json:"user,omitempty"`
}

// QwenUserInfo información del usuario autenticado
type QwenUserInfo struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	Plan       string `json:"plan"` // "free" | "plus"
	DailyLimit int    `json:"daily_limit"`
	Remaining  int    `json:"remaining_today"`
}

// QwenRefreshResponse respuesta del endpoint /api/auth/refresh
type QwenRefreshResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	TokenType    string    `json:"token_type"`
}

// QwenDeviceCodeResponse respuesta del endpoint /api/auth/device
type QwenDeviceCodeResponse struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURI string `json:"verification_uri"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
}

// QwenSession representa una sesión OAuth de Qwen Portal
// Se persiste en ~/.picoclaw/workspace/state/qwen_session.json
type QwenSession struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	DeviceCode   string    `json:"device_code"`
	ExpiresAt    time.Time `json:"expires_at"`
	UserID       string    `json:"user_id"`
	Email        string    `json:"email"`
	Plan         string    `json:"plan"` // "free" | "plus"
	DailyLimit   int       `json:"daily_limit"`
	LastVerified time.Time `json:"last_verified"`
}

// IsExpired verifica si la sesión ha expirado
// Considera expirado 1 hora antes para margen de seguridad
func (s *QwenSession) IsExpired() bool {
	if s.ExpiresAt.IsZero() {
		return true
	}
	return time.Now().Add(1 * time.Hour).After(s.ExpiresAt)
}

// NeedsRefresh verifica si la sesión necesita renovación
// Renueva si faltan menos de 3 días para expiración
func (s *QwenSession) NeedsRefresh() bool {
	if s.IsExpired() {
		return true
	}
	return time.Now().Add(3 * 24 * time.Hour).After(s.ExpiresAt)
}

// Validate verifica que la sesión tenga datos mínimos válidos
func (s *QwenSession) Validate() error {
	if s.AccessToken == "" {
		return ErrMissingAccessToken
	}
	if s.Email == "" {
		return ErrMissingEmail
	}
	if s.ExpiresAt.IsZero() {
		return ErrMissingExpiresAt
	}
	return nil
}

// Errores específicos de Qwen
var (
	ErrMissingAccessToken = &QwenError{Message: "access token is missing"}
	ErrMissingEmail       = &QwenError{Message: "email is missing"}
	ErrMissingExpiresAt   = &QwenError{Message: "expires_at is missing"}
	ErrSessionExpired     = &QwenError{Message: "session has expired"}
	ErrInvalidDeviceCode  = &QwenError{Message: "invalid device code format"}
)

// QwenError error específico de autenticación Qwen
type QwenError struct {
	Message string
	Err     error
}

func (e *QwenError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *QwenError) Unwrap() error {
	return e.Err
}
