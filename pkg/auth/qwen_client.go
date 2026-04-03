// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	// QwenBaseURL es la URL base para la API de Qwen Portal
	QwenBaseURL = "https://chat.qwen.ai"

	// QwenAPIAuthStatus endpoint para verificar estado de autenticación
	QwenAPIAuthStatus = "/api/auth/status"

	// QwenAPIAuthRefresh endpoint para renovar tokens
	QwenAPIAuthRefresh = "/api/auth/refresh"

	// QwenHTTPTimeout timeout para operaciones HTTP
	QwenHTTPTimeout = 30 * time.Second
)

// QwenClient cliente HTTP para interactuar con Qwen Portal API
type QwenClient struct {
	baseURL    string
	httpClient *http.Client
	session    *QwenSession
}

// NewQwenClient crea un nuevo cliente Qwen
// session: sesión Qwen válida (puede ser nil para operaciones públicas)
func NewQwenClient(session *QwenSession) *QwenClient {
	return &QwenClient{
		baseURL: QwenBaseURL,
		httpClient: &http.Client{
			Timeout: QwenHTTPTimeout,
		},
		session: session,
	}
}

// NewQwenClientWithTimeout crea un cliente Qwen con timeout personalizado
func NewQwenClientWithTimeout(session *QwenSession, timeout time.Duration) *QwenClient {
	return &QwenClient{
		baseURL: QwenBaseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		session: session,
	}
}

// AuthStatus verifica el estado de autenticación con el servidor Qwen
// Retorna información del usuario autenticado
func (c *QwenClient) AuthStatus(ctx context.Context) (*QwenAuthStatusResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		c.baseURL+QwenAPIAuthStatus, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Agregar header de autorización si hay sesión
	if c.session != nil && c.session.AccessToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.session.AccessToken)
	}

	// Agregar headers específicos de Qwen
	req.Header.Set("User-Agent", "picoclaw-agents/1.0")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Manejar errores HTTP
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("auth status request failed: HTTP %d - %s",
			resp.StatusCode, string(body))
	}

	var status QwenAuthStatusResponse
	if err := json.Unmarshal(body, &status); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &status, nil
}

// RefreshToken renueva el access token usando el refresh token
// Retorna una nueva sesión con tokens actualizados
func (c *QwenClient) RefreshToken(ctx context.Context) (*QwenSession, error) {
	if c.session == nil {
		return nil, fmt.Errorf("session is nil")
	}

	if c.session.RefreshToken == "" {
		return nil, fmt.Errorf("refresh token is empty")
	}

	payload := map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": c.session.RefreshToken,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		c.baseURL+QwenAPIAuthRefresh, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "picoclaw-agents/1.0")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Manejar errores HTTP
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("token refresh failed: HTTP %d - %s",
			resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var refreshResp QwenRefreshResponse
	if err := json.Unmarshal(body, &refreshResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Crear nueva sesión con tokens actualizados
	newSession := &QwenSession{
		AccessToken:  refreshResp.AccessToken,
		RefreshToken: refreshResp.RefreshToken,
		ExpiresAt:    refreshResp.ExpiresAt,
		UserID:       c.session.UserID,
		Email:        c.session.Email,
		Plan:         c.session.Plan,
		DailyLimit:   c.session.DailyLimit,
		LastVerified: time.Now(),
	}

	return newSession, nil
}

// GetUserInfo obtiene información del usuario autenticado
// Wrapper convenience sobre AuthStatus
func (c *QwenClient) GetUserInfo(ctx context.Context) (*QwenUserInfo, error) {
	status, err := c.AuthStatus(ctx)
	if err != nil {
		return nil, err
	}

	if !status.Authenticated {
		return nil, fmt.Errorf("not authenticated: %s", status.Error)
	}

	return &status.User, nil
}

// VerifySession verifica que la sesión sea válida y no esté expirada
// Combina validación local y verificación remota
func (c *QwenClient) VerifySession(ctx context.Context) error {
	if c.session == nil {
		return fmt.Errorf("session is nil")
	}

	// Validación local
	if err := c.session.Validate(); err != nil {
		return fmt.Errorf("local validation failed: %w", err)
	}

	if c.session.IsExpired() {
		return ErrSessionExpired
	}

	// Verificación remota (opcional, puede fallar por red)
	_, err := c.AuthStatus(ctx)
	if err != nil {
		// Si falla la verificación remota pero la sesión es válida localmente,
		// permitimos continuar (puede ser problema temporal de red)
		//nolint:nilerr // Intentional: remote verification failure is acceptable
		return nil
	}

	return nil
}

// UpdateSession actualiza la sesión con información fresca del servidor
// Útil para refresh de daily_limit y remaining
func (c *QwenClient) UpdateSession(ctx context.Context) (*QwenSession, error) {
	if c.session == nil {
		return nil, fmt.Errorf("session is nil")
	}

	userInfo, err := c.GetUserInfo(ctx)
	if err != nil {
		return nil, err
	}

	// Actualizar sesión con información fresca
	c.session.UserID = userInfo.ID
	c.session.Email = userInfo.Email
	c.session.Plan = userInfo.Plan
	c.session.DailyLimit = userInfo.DailyLimit
	c.session.LastVerified = time.Now()

	return c.session, nil
}

// IsHealthy verifica si el servicio Qwen está disponible
// Health check simple sin autenticación
func (c *QwenClient) IsHealthy(ctx context.Context) bool {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodHead, c.baseURL, nil)
	if err != nil {
		return false
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode >= 200 && resp.StatusCode < 400
}
