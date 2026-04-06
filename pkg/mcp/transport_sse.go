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
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/comgunner/picoclaw/pkg/logger"
)

// SSETransport communicates with an MCP server via Server-Sent Events (HTTP EventStream).
// Requests are sent as JSON-RPC POST, responses are parsed from SSE events (data: lines).
type SSETransport struct {
	url     string
	headers map[string]string
	mu      sync.Mutex
	id      int64
	client  *http.Client
}

// NewSSETransport creates an SSE transport for MCP communication.
// The URL points to the SSE endpoint of the MCP server.
func NewSSETransport(url string, headers map[string]string) (*SSETransport, error) {
	if url == "" {
		return nil, fmt.Errorf("SSE transport requires URL")
	}
	return &SSETransport{
		url:     url,
		headers: headers,
		client:  &http.Client{},
	}, nil
}

// Call sends a JSON-RPC request via HTTP POST and reads the response from the SSE stream.
// Respects context cancellation — no blocking reads.
func (t *SSETransport) Call(ctx context.Context, method string, params map[string]any) (*json.RawMessage, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

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

	// Create HTTP POST request
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, t.url, bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "text/event-stream")

	// Apply custom headers
	for k, v := range t.headers {
		httpReq.Header.Set(k, v)
	}

	// Send request
	resp, err := t.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("HTTP request: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			logger.WarnCF("mcp", "SSE response body close error",
				map[string]any{"error": closeErr.Error()})
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	// Check content type — should be text/event-stream
	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "text/event-stream") {
		// Some servers may return JSON directly; try parsing as JSON-RPC response
		return t.tryParseJSON(resp.Body, reqID)
	}

	// Parse SSE events from the response body
	return t.parseSSEEvents(ctx, resp.Body, reqID)
}

// tryParseJSON attempts to parse the response body as a direct JSON-RPC response
// (fallback for servers that don't return proper SSE streams).
func (t *SSETransport) tryParseJSON(body io.Reader, reqID int64) (*json.RawMessage, error) {
	raw, err := io.ReadAll(io.LimitReader(body, MAX_LINE_BYTES))
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}
	if len(raw) == 0 {
		return nil, fmt.Errorf("empty response")
	}

	var resp jsonRPCResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	if resp.Error != nil {
		return nil, fmt.Errorf("MCP error %d: %s", resp.Error.Code, resp.Error.Message)
	}
	return &resp.Result, nil
}

// parseSSEEvents reads Server-Sent Events from the body, looking for a response
// matching the given request ID. Respects context cancellation.
func (t *SSETransport) parseSSEEvents(ctx context.Context, body io.Reader, reqID int64) (*json.RawMessage, error) {
	scanner := bufio.NewScanner(body)
	scanner.Buffer(make([]byte, MAX_LINE_BYTES), MAX_LINE_BYTES)

	var eventData strings.Builder
	var eventType string

	for scanner.Scan() {
		// Check context cancellation on each line
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		line := scanner.Text()

		// Empty line = end of event block
		if line == "" {
			if eventData.Len() > 0 {
				result, found, err := t.processSSEEvent(eventType, eventData.String(), reqID)
				if err != nil {
					return nil, err
				}
				if found {
					return result, nil
				}
				eventData.Reset()
				eventType = ""
			}
			continue
		}

		// Parse SSE fields
		if strings.HasPrefix(line, "data:") {
			data := strings.TrimPrefix(line, "data:")
			data = strings.TrimSpace(data)
			if eventData.Len() > 0 {
				eventData.WriteString("\n")
			}
			eventData.WriteString(data)
		} else if strings.HasPrefix(line, "event:") {
			eventType = strings.TrimSpace(strings.TrimPrefix(line, "event:"))
		} else if strings.HasPrefix(line, "id:") {
			// Event ID — logged for debugging but not used for matching
			_ = strings.TrimSpace(strings.TrimPrefix(line, "id:"))
		} else if strings.HasPrefix(line, "retry:") {
			// Retry field — not applicable for our use case
		} else if strings.HasPrefix(line, ":") {
			// Comment line — ignore
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("SSE scanner error: %w", err)
	}

	return nil, fmt.Errorf("SSE stream ended without response for request ID %d", reqID)
}

// processSSEEvent attempts to parse an SSE event as a JSON-RPC response matching reqID.
// Returns (result, found, error).
func (t *SSETransport) processSSEEvent(eventType, eventData string, reqID int64) (*json.RawMessage, bool, error) {
	if eventData == "" {
		return nil, false, nil
	}

	var resp jsonRPCResponse
	if err := json.Unmarshal([]byte(eventData), &resp); err != nil {
		// Not a valid JSON-RPC response — skip
		logger.WarnCF("mcp", "SSE event not a valid JSON-RPC response",
			map[string]any{"event": eventType, "error": err.Error()})
		return nil, false, nil
	}

	// Match by request ID
	if resp.ID != reqID {
		return nil, false, nil
	}

	if resp.Error != nil {
		return nil, true, fmt.Errorf("MCP error %d: %s", resp.Error.Code, resp.Error.Message)
	}
	return &resp.Result, true, nil
}

// Close terminates the transport and releases resources.
// For SSE, this is a no-op since connections are per-request.
func (t *SSETransport) Close() error {
	return nil
}
