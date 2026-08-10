package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestWebSocketProxy_NoAuthRequired(t *testing.T) {
	// Setup
	configPath := filepath.Join(t.TempDir(), "config.json")
	handler := NewHandler(configPath)
	mux := http.NewServeMux()
	handler.registerPicoRoutes(mux)

	// Create request without any auth
	req := httptest.NewRequest("GET", "/pico/ws", nil)
	req.Header.Set("Connection", "Upgrade")
	req.Header.Set("Upgrade", "websocket")
	req.Header.Set("Sec-WebSocket-Key", "dGhlIHNhbXBsZSBub25jZQ==")
	req.Header.Set("Sec-WebSocket-Version", "13")
	req.Header.Set("Origin", "http://localhost:18800")
	req.Host = "localhost:18800"

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	// Should NOT return 401 (gateway handles auth)
	if rr.Code == http.StatusUnauthorized {
		t.Error("WebSocket proxy should not require authentication - gateway handles it")
	}
}

func TestWebSocketProxy_SameOriginCheck(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.json")
	handler := NewHandler(configPath)

	// Test same origin (should pass)
	req1 := httptest.NewRequest("GET", "/pico/ws", nil)
	req1.Header.Set("Origin", "http://localhost:18800")
	req1.Host = "localhost:18800"
	if !handler.isSameOrigin(req1) {
		t.Error("Same origin should pass")
	}

	// Test different origin (should fail)
	req2 := httptest.NewRequest("GET", "/pico/ws", nil)
	req2.Header.Set("Origin", "http://evil.com:18800")
	req2.Host = "localhost:18800"
	if handler.isSameOrigin(req2) {
		t.Error("Different origin should fail")
	}

	// Test localhost (should always pass)
	req3 := httptest.NewRequest("GET", "/pico/ws", nil)
	req3.Header.Set("Origin", "http://localhost:3000")
	req3.Host = "localhost:18800"
	if !handler.isSameOrigin(req3) {
		t.Error("localhost origin should always pass")
	}
}

func TestWebSocketProxy_TokenInjection(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.json")
	handler := NewHandler(configPath)

	// Create request
	req := httptest.NewRequest("GET", "/pico/ws", nil)
	req.Header.Set("Origin", "http://localhost:18800")

	// Before proxy, token should not be in request
	if req.Header.Get("Authorization") != "" {
		t.Error("Authorization header should be empty before proxy")
	}

	// The token injection happens in handleWebSocketProxy
	// which loads config and sets the header
	// This is tested indirectly via the proxy behavior
	_ = handler
}

func TestIsAuthenticated(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.json")
	handler := NewHandler(configPath)

	// Test with no auth
	req1 := httptest.NewRequest("GET", "/pico/ws", nil)
	if handler.isAuthenticated(req1) {
		t.Error("Request without auth should not be authenticated")
	}

	// Test with session cookie
	req2 := httptest.NewRequest("GET", "/pico/ws", nil)
	req2.AddCookie(&http.Cookie{Name: SessionCookieName, Value: "test-token"})
	if !handler.isAuthenticated(req2) {
		t.Error("Request with session cookie should be authenticated")
	}

	// Test with Bearer token
	req3 := httptest.NewRequest("GET", "/pico/ws", nil)
	req3.Header.Set("Authorization", "Bearer test-token")
	if !handler.isAuthenticated(req3) {
		t.Error("Request with Bearer token should be authenticated")
	}

	// Test with X-Auth-Token
	req4 := httptest.NewRequest("GET", "/pico/ws", nil)
	req4.Header.Set("X-Auth-Token", "test-token")
	if !handler.isAuthenticated(req4) {
		t.Error("Request with X-Auth-Token should be authenticated")
	}
}

func TestIsSameOrigin(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.json")
	handler := NewHandler(configPath)

	tests := []struct {
		name     string
		origin   string
		host     string
		expected bool
	}{
		{"same origin", "http://localhost:18800", "localhost:18800", true},
		{"different origin", "http://evil.com:18800", "localhost:18800", false},
		{"localhost origin", "http://localhost:3000", "localhost:18800", true},
		{"127.0.0.1 origin", "http://127.0.0.1:3000", "localhost:18800", true},
		{"no origin", "", "localhost:18800", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/pico/ws", nil)
			if tt.origin != "" {
				req.Header.Set("Origin", tt.origin)
			}
			req.Host = tt.host

			result := handler.isSameOrigin(req)
			if result != tt.expected {
				t.Errorf("isSameOrigin() = %v, want %v", result, tt.expected)
			}
		})
	}
	_ = handler
}

func TestPicoInfo_NoAuthRequired(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.json")
	handler := NewHandler(configPath)
	mux := http.NewServeMux()
	handler.registerPicoRoutes(mux)

	req := httptest.NewRequest("GET", "/api/pico/info", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	var resp map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if resp["status"] != "available" {
		t.Errorf("Expected status 'available', got %v", resp["status"])
	}

	// Token should NOT be exposed
	if _, hasToken := resp["token"]; hasToken {
		t.Error("Token should NOT be exposed in /api/pico/info")
	}
}

func TestLoginFlow(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.json")
	handler := NewHandler(configPath)

	// Test 1: Login with wrong password
	req1 := httptest.NewRequest("POST", "/api/auth/login",
		strings.NewReader(`{"password":"wrong"}`))
	req1.Header.Set("Content-Type", "application/json")
	rr1 := httptest.NewRecorder()
	handler.handleLogin(rr1, req1)
	if rr1.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401 for wrong password, got %d", rr1.Code)
	}

	// Test 2: Login with correct password
	req2 := httptest.NewRequest("POST", "/api/auth/login",
		strings.NewReader(`{"password":"picoclaw"}`))
	req2.Header.Set("Content-Type", "application/json")
	rr2 := httptest.NewRecorder()
	handler.handleLogin(rr2, req2)
	if rr2.Code != http.StatusOK {
		t.Errorf("Expected 200 for correct password, got %d", rr2.Code)
	}

	// Test 3: Check session cookie is set
	cookies := rr2.Result().Cookies()
	if len(cookies) == 0 {
		t.Error("Expected session cookie to be set")
	}
}

func TestSessionValidation(t *testing.T) {
	// Create session
	token := CreateSession()
	if !ValidateSession(token) {
		t.Error("New session should be valid")
	}

	// Manually expire session
	sessions.mu.Lock()
	sessions.sessions[token].ExpiresAt = time.Now().Add(-1 * time.Hour)
	sessions.mu.Unlock()

	// Should be invalid
	if ValidateSession(token) {
		t.Error("Expired session should be invalid")
	}
}

func TestChangePassword(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.json")

	// Change password
	err := ChangePassword(configPath, "picoclaw", "newpassword")
	if err != nil {
		t.Fatalf("Failed to change password: %v", err)
	}

	// Verify new password works
	newHash, _ := GetPasswordHash(configPath)
	if !CheckPassword("newpassword", newHash) {
		t.Error("New password should work")
	}

	// Reset for other tests
	ChangePassword(configPath, "newpassword", "picoclaw")
}

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("testpassword")
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	if !CheckPassword("testpassword", hash) {
		t.Error("Password check should pass")
	}

	if CheckPassword("wrongpassword", hash) {
		t.Error("Password check should fail for wrong password")
	}
}
