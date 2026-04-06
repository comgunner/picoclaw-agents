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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/comgunner/picoclaw/pkg/logger"
)

// McpSessionIDHeader is the header name used by MCP Streamable HTTP transport
// to carry the session identifier across requests (MCP 2025-03-26).
const McpSessionIDHeader = "Mcp-Session-Id"

// HTTPTransport communicates with an MCP server via Streamable HTTP (MCP 2025-03-26).
// It preserves the Mcp-Session-Id header across calls and handles session termination.
type HTTPTransport struct {
	url       string
	headers   map[string]string
	mu        sync.Mutex
	id        int64
	client    *http.Client
	sessionID string // Mcp-Session-Id from first response
	closed    bool
}

// NewHTTPTransport creates a Streamable HTTP transport for MCP communication.
// The URL points to the MCP server endpoint.
func NewHTTPTransport(url string, headers map[string]string) (*HTTPTransport, error) {
	if url == "" {
		return nil, fmt.Errorf("HTTP transport requires URL")
	}
	return &HTTPTransport{
		url:     url,
		headers: headers,
		client:  &http.Client{},
	}, nil
}

// Call sends a JSON-RPC request via HTTP POST and returns the response.
// It preserves the Mcp-Session-Id header across calls for session management.
// Respects context cancellation — no blocking reads.
func (t *HTTPTransport) Call(ctx context.Context, method string, params map[string]any) (*json.RawMessage, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.closed {
		return nil, fmt.Errorf("HTTP transport is closed")
	}

	atomic.AddInt64(&t.id, 1)
	reqID := t.id

	req := jsonRPCRequest{
		JSONRPC: "2.0",
		ID:      reqID,
		Method:  method,
		Params:  params,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, t.url, bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json, text/event-stream")

	// Apply custom headers
	for k, v := range t.headers {
		httpReq.Header.Set(k, v)
	}

	// Attach session ID if we have one from a previous response
	if t.sessionID != "" {
		httpReq.Header.Set(McpSessionIDHeader, t.sessionID)
	}

	resp, err := t.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("HTTP request: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			logger.WarnCF("mcp", "HTTP response body close error",
				map[string]any{"error": closeErr.Error()})
		}
	}()

	// Save session ID from first response (MCP 2025-03-26 spec)
	if sessionID := resp.Header.Get(McpSessionIDHeader); sessionID != "" && t.sessionID == "" {
		t.sessionID = sessionID
		logger.InfoCF("mcp", "MCP session established",
			map[string]any{"session_id": sessionID})
	}

	// Handle session termination: 404 means session expired
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("MCP session expired (HTTP 404)")
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	// Read response body
	body, err := io.ReadAll(io.LimitReader(resp.Body, MAX_LINE_BYTES))
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}
	if len(body) == 0 {
		return nil, fmt.Errorf("empty response body")
	}

	// Try parsing as JSON-RPC response
	var respJSON jsonRPCResponse
	if err := json.Unmarshal(body, &respJSON); err != nil {
		return nil, fmt.Errorf("parse JSON response: %w", err)
	}

	if respJSON.Error != nil {
		return nil, fmt.Errorf("MCP error %d: %s", respJSON.Error.Code, respJSON.Error.Message)
	}

	return &respJSON.Result, nil
}

// Close terminates the transport and releases resources.
// For HTTP transport, this is a no-op since connections are per-request.
func (t *HTTPTransport) Close() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.closed = true
	return nil
}
