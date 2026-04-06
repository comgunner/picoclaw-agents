// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package mcp

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestHTTPTransport_Call_Success(t *testing.T) {
	respBody := jsonRPCResponse{
		JSONRPC: "2.0",
		ID:      1,
		Result:  json.RawMessage(`{"ok":true}`),
	}
	respData, _ := json.Marshal(respBody)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(respData)
	}))
	defer server.Close()

	transport, err := NewHTTPTransport(server.URL, nil)
	if err != nil {
		t.Fatalf("NewHTTPTransport: %v", err)
	}
	defer transport.Close()

	result, err := transport.Call(context.Background(), "initialize", nil)
	if err != nil {
		t.Fatalf("Call: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestHTTPTransport_Call_SessionID(t *testing.T) {
	var requestCount int
	sessionID := "test-session-abc123"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		w.Header().Set("Content-Type", "application/json")

		// First request: return session ID
		if requestCount == 1 {
			w.Header().Set(McpSessionIDHeader, sessionID)
		} else {
			// Subsequent requests: verify session ID is present
			got := r.Header.Get(McpSessionIDHeader)
			if got != sessionID {
				t.Errorf("expected session ID %q, got %q", sessionID, got)
			}
		}

		respBody := jsonRPCResponse{
			JSONRPC: "2.0",
			ID:      int64(requestCount),
			Result:  json.RawMessage(`{"ok":true}`),
		}
		respData, _ := json.Marshal(respBody)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(respData)
	}))
	defer server.Close()

	transport, err := NewHTTPTransport(server.URL, nil)
	if err != nil {
		t.Fatalf("NewHTTPTransport: %v", err)
	}
	defer transport.Close()

	// First call — should receive session ID
	_, err = transport.Call(context.Background(), "initialize", nil)
	if err != nil {
		t.Fatalf("first Call: %v", err)
	}

	// Verify session ID was saved
	if transport.sessionID != sessionID {
		t.Errorf("expected sessionID %q, got %q", sessionID, transport.sessionID)
	}

	// Second call — should include session ID
	_, err = transport.Call(context.Background(), "tools/list", nil)
	if err != nil {
		t.Fatalf("second Call: %v", err)
	}
}

func TestHTTPTransport_Call_ContextCancel(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate slow server
		time.Sleep(500 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	transport, err := NewHTTPTransport(server.URL, nil)
	if err != nil {
		t.Fatalf("NewHTTPTransport: %v", err)
	}
	defer transport.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_, err = transport.Call(ctx, "tools/list", nil)
	if err == nil {
		t.Fatal("expected context cancellation error, got nil")
	}
}

func TestHTTPTransport_Call_SessionExpired(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"session not found"}`))
	}))
	defer server.Close()

	transport, err := NewHTTPTransport(server.URL, nil)
	if err != nil {
		t.Fatalf("NewHTTPTransport: %v", err)
	}
	defer transport.Close()

	_, err = transport.Call(context.Background(), "tools/list", nil)
	if err == nil {
		t.Fatal("expected session expired error, got nil")
	}
	if !strings.Contains(err.Error(), "404") && !strings.Contains(err.Error(), "session expired") {
		t.Errorf("expected session expired error, got: %v", err)
	}
}

func TestHTTPTransport_Call_NonFatalError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "bad request", http.StatusBadRequest)
	}))
	defer server.Close()

	transport, err := NewHTTPTransport(server.URL, nil)
	if err != nil {
		t.Fatalf("NewHTTPTransport: %v", err)
	}
	defer transport.Close()

	_, err = transport.Call(context.Background(), "tools/list", nil)
	if err == nil {
		t.Fatal("expected HTTP error, got nil")
	}
	if !strings.Contains(err.Error(), "400") {
		t.Errorf("expected error to contain '400', got: %v", err)
	}
}

func TestHTTPTransport_Close(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	transport, err := NewHTTPTransport(server.URL, nil)
	if err != nil {
		t.Fatalf("NewHTTPTransport: %v", err)
	}

	// Close should not panic
	err = transport.Close()
	if err != nil {
		t.Errorf("Close() returned error: %v", err)
	}

	// After close, calls should fail
	_, err = transport.Call(context.Background(), "tools/list", nil)
	if err == nil {
		t.Error("expected error after close, got nil")
	}
	if !strings.Contains(err.Error(), "closed") {
		t.Errorf("expected 'closed' error, got: %v", err)
	}
}

func TestHTTPTransport_Call_MCPError(t *testing.T) {
	respBody := jsonRPCResponse{
		JSONRPC: "2.0",
		ID:      1,
		Error: &jsonRPCError{
			Code:    -32600,
			Message: "Invalid Request",
		},
	}
	respData, _ := json.Marshal(respBody)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(respData)
	}))
	defer server.Close()

	transport, err := NewHTTPTransport(server.URL, nil)
	if err != nil {
		t.Fatalf("NewHTTPTransport: %v", err)
	}
	defer transport.Close()

	_, err = transport.Call(context.Background(), "invalid", nil)
	if err == nil {
		t.Fatal("expected MCP error, got nil")
	}
	if !strings.Contains(err.Error(), "-32600") {
		t.Errorf("expected error to contain error code, got: %v", err)
	}
}

func TestNewHTTPTransport_EmptyURL(t *testing.T) {
	_, err := NewHTTPTransport("", nil)
	if err == nil {
		t.Error("expected error for empty URL, got nil")
	}
}
