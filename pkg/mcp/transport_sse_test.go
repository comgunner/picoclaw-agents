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

func TestSSETransport_Call_Success(t *testing.T) {
	respBody := jsonRPCResponse{
		JSONRPC: "2.0",
		ID:      1,
		Result:  json.RawMessage(`{"tools":[]}`),
	}
	respData, _ := json.Marshal(respBody)

	sseData := "data: " + string(respData) + "\n\n"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(sseData))
	}))
	defer server.Close()

	transport, err := NewSSETransport(server.URL, nil)
	if err != nil {
		t.Fatalf("NewSSETransport: %v", err)
	}
	defer transport.Close()

	result, err := transport.Call(context.Background(), "tools/list", nil)
	if err != nil {
		t.Fatalf("Call: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestSSETransport_Call_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	}))
	defer server.Close()

	transport, err := NewSSETransport(server.URL, nil)
	if err != nil {
		t.Fatalf("NewSSETransport: %v", err)
	}
	defer transport.Close()

	_, err = transport.Call(context.Background(), "tools/list", nil)
	if err == nil {
		t.Fatal("expected error for HTTP 500, got nil")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("expected error to contain '500', got: %v", err)
	}
}

func TestSSETransport_Call_ContextCancel(t *testing.T) {
	// Server that never responds (hangs)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)
		// Write headers but never send data — simulates hanging SSE stream
		flusher, ok := w.(http.Flusher)
		if ok {
			flusher.Flush()
		}
		// Block forever (or until test timeout)
		<-r.Context().Done()
	}))
	defer server.Close()

	transport, err := NewSSETransport(server.URL, nil)
	if err != nil {
		t.Fatalf("NewSSETransport: %v", err)
	}
	defer transport.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err = transport.Call(ctx, "tools/list", nil)
	if err == nil {
		t.Fatal("expected context cancellation error, got nil")
	}
	// The error may be wrapped (e.g., "SSE scanner error: context deadline exceeded")
	if !strings.Contains(err.Error(), "context deadline exceeded") && !strings.Contains(err.Error(), "canceled") {
		t.Errorf("expected context cancellation error, got: %v", err)
	}
}

func TestSSETransport_Call_ParseSSE(t *testing.T) {
	// Test with full SSE event format including event:, id:, retry: fields
	respBody := jsonRPCResponse{
		JSONRPC: "2.0",
		ID:      1,
		Result:  json.RawMessage(`{"ok":true}`),
	}
	respData, _ := json.Marshal(respBody)

	sseData := "event: message\n" +
		"id: 123\n" +
		"retry: 5000\n" +
		"data: " + string(respData) + "\n\n"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(sseData))
	}))
	defer server.Close()

	transport, err := NewSSETransport(server.URL, nil)
	if err != nil {
		t.Fatalf("NewSSETransport: %v", err)
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

func TestSSETransport_Call_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)
		// Send empty SSE events — no data matching our request ID
		_, _ = w.Write([]byte("\n\n"))
	}))
	defer server.Close()

	transport, err := NewSSETransport(server.URL, nil)
	if err != nil {
		t.Fatalf("NewSSETransport: %v", err)
	}
	defer transport.Close()

	_, err = transport.Call(context.Background(), "tools/list", nil)
	if err == nil {
		t.Fatal("expected error for empty response, got nil")
	}
	if !strings.Contains(err.Error(), "without response") {
		t.Errorf("expected 'without response' error, got: %v", err)
	}
}

func TestSSETransport_Close(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	transport, err := NewSSETransport(server.URL, nil)
	if err != nil {
		t.Fatalf("NewSSETransport: %v", err)
	}

	// Close should not panic and should return nil
	err = transport.Close()
	if err != nil {
		t.Errorf("Close() returned error: %v", err)
	}

	// Double-close should also be safe
	err = transport.Close()
	if err != nil {
		t.Errorf("Close() on closed transport returned error: %v", err)
	}
}

func TestSSETransport_Call_MCPError(t *testing.T) {
	respBody := jsonRPCResponse{
		JSONRPC: "2.0",
		ID:      1,
		Error: &jsonRPCError{
			Code:    -32601,
			Message: "Method not found",
		},
	}
	respData, _ := json.Marshal(respBody)

	sseData := "data: " + string(respData) + "\n\n"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(sseData))
	}))
	defer server.Close()

	transport, err := NewSSETransport(server.URL, nil)
	if err != nil {
		t.Fatalf("NewSSETransport: %v", err)
	}
	defer transport.Close()

	_, err = transport.Call(context.Background(), "invalid_method", nil)
	if err == nil {
		t.Fatal("expected MCP error, got nil")
	}
	if !strings.Contains(err.Error(), "-32601") {
		t.Errorf("expected error to contain error code, got: %v", err)
	}
}

func TestNewSSETransport_EmptyURL(t *testing.T) {
	_, err := NewSSETransport("", nil)
	if err == nil {
		t.Error("expected error for empty URL, got nil")
	}
}
