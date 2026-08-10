package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	// SessionDuration is how long a session cookie lasts.
	SessionDuration = 24 * time.Hour
	// SessionCookieName is the name of the session cookie.
	SessionCookieName = "picoclaw_session"
	// PasswordFile is the file storing the hashed password.
	PasswordFile = ".dashboard_password"
	// DefaultPassword is the initial password for first-time setup.
	DefaultPassword = "picoclaw" // pragma: allowlist secret
)

// Session represents an active user session.
type Session struct {
	Token     string
	ExpiresAt time.Time
}

// sessionManager manages active sessions.
type sessionManager struct {
	sessions map[string]*Session
	mu       sync.RWMutex
}

var sessions = &sessionManager{
	sessions: make(map[string]*Session),
}

// CreateSession creates a new session and returns the token.
func CreateSession() string {
	token := generateSecureToken()
	sessions.mu.Lock()
	defer sessions.mu.Unlock()
	sessions.sessions[token] = &Session{
		Token:     token,
		ExpiresAt: time.Now().Add(SessionDuration),
	}
	return token
}

// ValidateSession checks if a session token is valid.
func ValidateSession(token string) bool {
	sessions.mu.RLock()
	defer sessions.mu.RUnlock()
	session, ok := sessions.sessions[token]
	if !ok {
		return false
	}
	if time.Now().After(session.ExpiresAt) {
		delete(sessions.sessions, token)
		return false
	}
	return true
}

// CleanupSessions removes expired sessions.
func CleanupSessions() {
	sessions.mu.Lock()
	defer sessions.mu.Unlock()
	now := time.Now()
	for token, session := range sessions.sessions {
		if now.After(session.ExpiresAt) {
			delete(sessions.sessions, token)
		}
	}
}

// HashPassword hashes a password using bcrypt.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword compares a password with a hash.
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GetPasswordHash returns the stored password hash, or creates a default one.
func GetPasswordHash(configPath string) (string, error) {
	passwordPath := filepath.Join(filepath.Dir(configPath), PasswordFile)

	// Check if password file exists
	data, err := os.ReadFile(passwordPath)
	if err == nil && len(data) > 0 {
		return string(data), nil
	}

	// Create default password hash
	hash, err := HashPassword(DefaultPassword)
	if err != nil {
		return "", fmt.Errorf("failed to hash default password: %w", err)
	}

	// Save the hash
	if err := os.WriteFile(passwordPath, []byte(hash), 0o600); err != nil {
		return "", fmt.Errorf("failed to save password: %w", err)
	}

	return hash, nil
}

// ChangePassword changes the dashboard password.
func ChangePassword(configPath, oldPassword, newPassword string) error {
	currentHash, err := GetPasswordHash(configPath)
	if err != nil {
		return err
	}

	if !CheckPassword(oldPassword, currentHash) {
		return fmt.Errorf("incorrect current password")
	}

	newHash, err := HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	passwordPath := filepath.Join(filepath.Dir(configPath), PasswordFile)
	if err := os.WriteFile(passwordPath, []byte(newHash), 0o600); err != nil {
		return fmt.Errorf("failed to save new password: %w", err)
	}

	return nil
}

// registerAuthRoutes binds authentication endpoints to the ServeMux.
func (h *Handler) registerAuthRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/auth/login", h.handleLogin)
	mux.HandleFunc("POST /api/auth/logout", h.handleLogout)
	mux.HandleFunc("POST /api/auth/change-password", h.handleChangePassword)
	mux.HandleFunc("GET /api/auth/status", h.handleAuthStatus)
}

// handleLogin authenticates a user and creates a session.
func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	hash, err := GetPasswordHash(h.configPath)
	if err != nil {
		http.Error(w, "Failed to load password", http.StatusInternalServerError)
		return
	}

	if !CheckPassword(req.Password, hash) {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	token := CreateSession()

	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(SessionDuration.Seconds()),
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"message": "Login successful",
	})
}

// handleLogout destroys the current session.
func (h *Handler) handleLogout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(SessionCookieName)
	if err == nil {
		sessions.mu.Lock()
		delete(sessions.sessions, cookie.Value)
		sessions.mu.Unlock()
	}

	// Clear the cookie
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"message": "Logged out",
	})
}

// handleChangePassword changes the dashboard password.
func (h *Handler) handleChangePassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := ChangePassword(h.configPath, req.OldPassword, req.NewPassword); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"message": "Password changed",
	})
}

// handleAuthStatus returns the current authentication status.
func (h *Handler) handleAuthStatus(w http.ResponseWriter, r *http.Request) {
	authenticated := false
	cookie, err := r.Cookie(SessionCookieName)
	if err == nil && ValidateSession(cookie.Value) {
		authenticated = true
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"authenticated": authenticated,
	})
}
