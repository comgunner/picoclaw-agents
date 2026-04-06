package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/comgunner/picoclaw/pkg/config"
)

func TestProbeLocalModelAvailability_OpenAICompatibleIncludesAPIKey(t *testing.T) {
	const apiKey = "test-api-key"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/models" {
			t.Fatalf("path = %q, want %q", r.URL.Path, "/v1/models")
		}
		if got := r.Header.Get("Authorization"); got != "Bearer "+apiKey {
			http.Error(w, "missing auth", http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"id":"custom-model"}]}`))
	}))
	defer srv.Close()

	model := &config.ModelConfig{
		Model:   "openai/custom-model",
		APIBase: srv.URL + "/v1",
	}
	model.SetAPIKey(apiKey)

	if !probeLocalModelAvailability(model) {
		t.Fatal("probeLocalModelAvailability() = false, want true when api_key is configured")
	}
}

func TestIsOllamaAPIBase(t *testing.T) {
	tests := []struct {
		name     string
		apiBase  string
		expected bool
	}{
		{"localhost port 11434", "http://localhost:11434/v1", true},
		{"localhost port 11434 no v1", "http://localhost:11434", true},
		{"127.0.0.1 port 11434", "http://127.0.0.1:11434/v1", true},
		{"0.0.0.0 port 11434", "http://0.0.0.0:11434/v1", true},
		{"::1 port 11434", "http://[::1]:11434/v1", true},
		{"localhost no port", "http://localhost", true},
		{"localhost default port (empty)", "http://localhost:11434", true},
		{"different port", "http://localhost:8000/v1", false},
		{"remote host", "http://remote-server:11434/v1", false},
		{"empty string", "", false},
		{"whitespace only", "   ", false},
		{"invalid url", "not-a-url", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isOllamaAPIBase(tt.apiBase)
			if result != tt.expected {
				t.Errorf("isOllamaAPIBase(%q) = %v, want %v", tt.apiBase, result, tt.expected)
			}
		})
	}
}
