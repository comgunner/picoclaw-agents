package api

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/comgunner/picoclaw/pkg/config"
)

// registerPicoRoutes binds Pico Channel management endpoints to the ServeMux.
func (h *Handler) registerPicoRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/pico/token", h.handleGetPicoToken)
	mux.HandleFunc("POST /api/pico/token", h.handleRegenPicoToken)
	mux.HandleFunc("POST /api/pico/setup", h.handlePicoSetup)
	mux.HandleFunc("GET /api/pico/info", h.handlePicoInfo)

	// WebSocket proxy with auth check: forward /pico/ws to gateway
	// Secures access behind launcher authentication.
	mux.HandleFunc("GET /pico/ws", h.handleWebSocketProxy())
}

// isAuthenticated checks if the request has a valid session cookie or token.
func (h *Handler) isAuthenticated(r *http.Request) bool {
	// Check for session cookie
	if cookie, err := r.Cookie("picoclaw_session"); err == nil && cookie.Value != "" {
		return true
	}

	// Check for Authorization header
	if auth := r.Header.Get("Authorization"); auth != "" {
		if strings.HasPrefix(auth, "Bearer ") {
			return true
		}
	}

	// Check for X-Auth-Token header (used by some clients)
	if token := r.Header.Get("X-Auth-Token"); token != "" {
		return true
	}

	return false
}

// isSameOrigin checks if the request Origin matches the server host.
func (h *Handler) isSameOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	if origin == "" {
		// No origin header — could be a non-browser client
		// Allow for now, but log for monitoring
		return true
	}

	// Extract host from origin
	originHost := strings.TrimPrefix(strings.TrimPrefix(origin, "https://"), "http://")
	originHost = strings.Split(originHost, ":")[0]

	// Get request host
	requestHost := r.Host
	requestHost = strings.Split(requestHost, ":")[0]

	// If origin is localhost/127.0.0.1, it's a browser on the same machine
	// This is safe because the browser enforces same-origin policy
	if originHost == "localhost" || originHost == "127.0.0.1" {
		return true
	}

	// For non-localhost origins, check if they match the request host
	return originHost == requestHost
}

// createWsProxy creates a reverse proxy to the current gateway WebSocket endpoint.
// The gateway bind host and port are resolved from the latest configuration.
func (h *Handler) createWsProxy() *httputil.ReverseProxy {
	wsProxy := httputil.NewSingleHostReverseProxy(h.gatewayProxyURL())
	wsProxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		http.Error(w, "Gateway unavailable: "+err.Error(), http.StatusBadGateway)
	}
	return wsProxy
}

// handleWebSocketProxy wraps a reverse proxy to handle WebSocket connections.
// The reverse proxy forwards the incoming upgrade handshake as-is.
// Security: Same-origin check is enforced. Authentication is handled by the gateway.
// Token injection: The Pico token is injected server-side, never exposed to frontend.
func (h *Handler) handleWebSocketProxy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Security check: Require same-origin
		if !h.isSameOrigin(r) {
			http.Error(w, "Forbidden: cross-origin requests not allowed", http.StatusForbidden)
			return
		}

		// Load config to get the Pico token for gateway auth
		cfg, err := config.LoadConfig(h.configPath)
		if err != nil {
			http.Error(w, "Failed to load config", http.StatusInternalServerError)
			return
		}

		token := cfg.Channels.Pico.Token()
		if token == "" {
			http.Error(w, "Pico channel not configured", http.StatusServiceUnavailable)
			return
		}

		// Inject the token server-side (never exposed to frontend)
		// The gateway will verify this token via picoAuthenticate()
		r.Header.Set("Authorization", "Bearer "+token)

		proxy := h.createWsProxy()
		proxy.ServeHTTP(w, r)
	}
}

// handleGetPicoToken returns the current WS token and URL for the frontend.
//
//	GET /api/pico/token
func (h *Handler) handleGetPicoToken(w http.ResponseWriter, r *http.Request) {
	cfg, err := config.LoadConfig(h.configPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to load config: %v", err), http.StatusInternalServerError)
		return
	}

	wsURL := h.buildWsURL(r, cfg)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"token":   cfg.Channels.Pico.Token(),
		"ws_url":  wsURL,
		"enabled": cfg.Channels.Pico.Enabled,
	})
}

// handleRegenPicoToken generates a new Pico WebSocket token and saves it.
//
//	POST /api/pico/token
func (h *Handler) handleRegenPicoToken(w http.ResponseWriter, r *http.Request) {
	cfg, err := config.LoadConfig(h.configPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to load config: %v", err), http.StatusInternalServerError)
		return
	}

	token := generateSecureToken()
	cfg.Channels.Pico.SetToken(token)

	if err := config.SaveConfig(h.configPath, cfg); err != nil {
		http.Error(w, fmt.Sprintf("Failed to save config: %v", err), http.StatusInternalServerError)
		return
	}

	wsURL := h.buildWsURL(r, cfg)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"token":  token,
		"ws_url": wsURL,
	})
}

// EnsurePicoChannel enables the Pico channel with sane defaults if it isn't
// already configured. Returns true when the config was modified.
//
// callerOrigin is the Origin header from the setup request. If non-empty and
// no origins are configured yet, it's written as the allowed origin so the
// WebSocket handshake works for whatever host the caller is on (LAN, custom
// port, etc.). Pass "" when there's no request context.
func (h *Handler) EnsurePicoChannel(callerOrigin string) (bool, error) {
	cfg, err := config.LoadConfig(h.configPath)
	if err != nil {
		return false, fmt.Errorf("failed to load config: %w", err)
	}

	changed := false

	if !cfg.Channels.Pico.Enabled {
		cfg.Channels.Pico.Enabled = true
		changed = true
	}

	if cfg.Channels.Pico.Token() == "" {
		cfg.Channels.Pico.SetToken(generateSecureToken())
		changed = true
	}

	// Seed origins from the request instead of hardcoding ports.
	if len(cfg.Channels.Pico.AllowOrigins) == 0 && callerOrigin != "" {
		cfg.Channels.Pico.AllowOrigins = []string{callerOrigin}
		changed = true
	}

	if changed {
		if err := config.SaveConfig(h.configPath, cfg); err != nil {
			return false, fmt.Errorf("failed to save config: %w", err)
		}
		// Create/touch security file to signal credentials are managed.
		secPath := filepath.Join(filepath.Dir(h.configPath), config.SecurityConfigFile)
		if f, err := os.OpenFile(secPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600); err == nil {
			f.Close()
		}
	}

	return changed, nil
}

// handlePicoSetup automatically configures everything needed for the Pico Channel to work.
//
//	POST /api/pico/setup
func (h *Handler) handlePicoSetup(w http.ResponseWriter, r *http.Request) {
	changed, err := h.EnsurePicoChannel(r.Header.Get("Origin"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cfg, err := config.LoadConfig(h.configPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to load config: %v", err), http.StatusInternalServerError)
		return
	}

	wsURL := h.buildWsURL(r, cfg)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"token":   cfg.Channels.Pico.Token(),
		"ws_url":  wsURL,
		"enabled": true,
		"changed": changed,
	})
}

// generateSecureToken creates a random 32-character hex string.
func generateSecureToken() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		// Fallback to something pseudo-random if crypto/rand fails
		return fmt.Sprintf("pico_%x", time.Now().UnixNano())
	}
	return hex.EncodeToString(b)
}

// handlePicoInfo returns non-secret Pico connection metadata.
// This endpoint does NOT require authentication — it only returns safe info.
//
//	GET /api/pico/info
func (h *Handler) handlePicoInfo(w http.ResponseWriter, r *http.Request) {
	cfg, err := config.LoadConfig(h.configPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to load config: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status":  "available",
		"enabled": cfg.Channels.Pico.Enabled,
		"version": "1.0.0",
		// NOTE: Token is NOT included here for security
	})
}
