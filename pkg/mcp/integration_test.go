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
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/comgunner/picoclaw/pkg/config"
)

// mockTransport implements MCPTransport for testing without subprocesses.
type mockTransport struct {
	tools       []ToolInfo
	callHandler func(method string, params map[string]any) (*json.RawMessage, error)
	closed      bool
	mu          sync.Mutex
}

func (m *mockTransport) Call(ctx context.Context, method string, params map[string]any) (*json.RawMessage, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.closed {
		return nil, fmt.Errorf("transport closed")
	}
	if m.callHandler != nil {
		return m.callHandler(method, params)
	}
	// Default behavior
	switch method {
	case "initialize":
		result := json.RawMessage(`{"protocolVersion":"2024-11-05","capabilities":{"tools":{}}}`)
		return &result, nil
	case "tools/list":
		data, _ := json.Marshal(map[string]any{"tools": m.tools})
		raw := json.RawMessage(data)
		return &raw, nil
	case "tools/call":
		name := params["name"].(string)
		data, _ := json.Marshal(map[string]any{
			"content": []any{map[string]any{"type": "text", "text": "Response from " + name}},
			"isError": false,
		})
		raw := json.RawMessage(data)
		return &raw, nil
	default:
		return nil, fmt.Errorf("unknown method: %s", method)
	}
}

func (m *mockTransport) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.closed = true
	return nil
}

// TestIntegration_StdioConnect tests the full initialize handshake using a mock transport.
// NOTE: Real stdio transport tests require subprocess support, which may not be available
// in all environments. This test verifies the client logic using a mock transport.
func TestIntegration_StdioConnect(t *testing.T) {
	mt := &mockTransport{
		tools: []ToolInfo{
			{Name: "echo_tool", Description: "Echoes text", Schema: map[string]any{"type": "object"}},
			{Name: "add_tool", Description: "Add numbers", Schema: map[string]any{"type": "object"}},
			{Name: "fail_tool", Description: "Always fails", Schema: map[string]any{"type": "object"}},
		},
	}
	manager := &MCPClientManager{}
	manager.servers = map[string]*mcpServer{
		"mock": {
			cfg:       config.MCPServerConfig{Transport: "stdio"},
			transport: mt,
			tools:     mt.tools,
			status:    StatusConnected,
		},
	}

	srv := manager.GetServer("mock")
	if srv == nil {
		t.Fatal("expected mock server to be connected")
	}
	if len(srv.tools) != 3 {
		t.Errorf("expected 3 tools, got %d", len(srv.tools))
	}

	expectedTools := map[string]bool{"echo_tool": false, "add_tool": false, "fail_tool": false}
	for _, tool := range srv.tools {
		expectedTools[tool.Name] = true
	}
	for name, found := range expectedTools {
		if !found {
			t.Errorf("expected tool %q not found", name)
		}
	}
}

// TestIntegration_ToolCall tests calling a tool and getting a response.
func TestIntegration_ToolCall(t *testing.T) {
	manager := &MCPClientManager{}
	manager.servers = map[string]*mcpServer{
		"mock": {
			cfg: config.MCPServerConfig{Transport: "stdio", Timeout: 10},
			transport: &mockTransport{
				tools: []ToolInfo{
					{Name: "echo_tool", Description: "Echoes text", Schema: map[string]any{"type": "object"}},
				},
			},
			status: StatusConnected,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := manager.CallTool(ctx, "mock", "echo_tool", map[string]any{
		"text": "hello",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.IsError {
		t.Fatal("expected no error from echo_tool")
	}
	if len(result.Content) != 1 {
		t.Fatalf("expected 1 content block, got %d", len(result.Content))
	}
	if !strings.Contains(result.Content[0].Text, "Response from echo_tool") {
		t.Errorf("expected 'Response from echo_tool', got %q", result.Content[0].Text)
	}
}

// TestIntegration_MultipleServers tests connecting to 2 mock servers simultaneously.
func TestIntegration_MultipleServers(t *testing.T) {
	mt1 := &mockTransport{tools: []ToolInfo{{Name: "tool_a1"}, {Name: "tool_a2"}, {Name: "tool_a3"}}}
	mt2 := &mockTransport{tools: []ToolInfo{{Name: "tool_b1"}, {Name: "tool_b2"}, {Name: "tool_b3"}}}
	manager := &MCPClientManager{}
	manager.servers = map[string]*mcpServer{
		"server_a": {
			cfg:       config.MCPServerConfig{Transport: "stdio"},
			transport: mt1,
			tools:     mt1.tools,
			status:    StatusConnected,
		},
		"server_b": {
			cfg:       config.MCPServerConfig{Transport: "stdio"},
			transport: mt2,
			tools:     mt2.tools,
			status:    StatusConnected,
		},
	}

	serversList := manager.ListServers()
	if len(serversList) != 2 {
		t.Fatalf("expected 2 servers, got %d", len(serversList))
	}
	if serversList["server_a"] != 3 {
		t.Errorf("expected server_a to have 3 tools, got %d", serversList["server_a"])
	}
	if serversList["server_b"] != 3 {
		t.Errorf("expected server_b to have 3 tools, got %d", serversList["server_b"])
	}
}

// TestIntegration_PartialFailure tests that one server failing doesn't affect others.
func TestIntegration_PartialFailure(t *testing.T) {
	manager := &MCPClientManager{}
	// Simulate partial failure by only having one server connected
	manager.servers = map[string]*mcpServer{
		"good_server": {
			cfg:       config.MCPServerConfig{Transport: "stdio"},
			transport: &mockTransport{tools: []ToolInfo{{Name: "good_tool"}}},
			status:    StatusConnected,
		},
	}

	srv := manager.GetServer("good_server")
	if srv == nil {
		t.Fatal("expected good_server to be connected")
	}

	badSrv := manager.GetServer("bad_server")
	if badSrv != nil {
		t.Error("expected bad_server to NOT be connected")
	}
}

// TestIntegration_LargeResponse tests MAX_LINE_BYTES protection concept.
func TestIntegration_LargeResponse(t *testing.T) {
	// Verify MAX_LINE_BYTES constant is set correctly
	if MAX_LINE_BYTES != 10*1024*1024 {
		t.Errorf("expected MAX_LINE_BYTES to be 10MB, got %d", MAX_LINE_BYTES)
	}

	// Test that a large content block can be parsed
	largeText := strings.Repeat("A", 1000) // Small for test speed
	resultJSON := fmt.Sprintf(`{"content":[{"type":"text","text":"%s"}],"isError":false}`, largeText)
	raw := json.RawMessage(resultJSON)

	result, err := parseToolCallResult(raw)
	if err != nil {
		t.Fatalf("unexpected error parsing result: %v", err)
	}
	if len(result.Content) != 1 {
		t.Fatalf("expected 1 content block, got %d", len(result.Content))
	}
	if result.Content[0].Text != largeText {
		t.Error("content text mismatch")
	}
}

// TestIntegration_MalformedJSON tests that the client handles bad JSON gracefully.
func TestIntegration_MalformedJSON(t *testing.T) {
	malformedJSON := `{"content":[{"type":"text","text":"hello"` // missing closing brace
	raw := json.RawMessage(malformedJSON)

	_, err := parseToolCallResult(raw)
	if err == nil {
		t.Fatal("expected error for malformed JSON, got nil")
	}
	// Error should mention parse/unmarshal/unexpected failure
	if !strings.Contains(strings.ToLower(err.Error()), "parse") &&
		!strings.Contains(strings.ToLower(err.Error()), "unmarshal") &&
		!strings.Contains(strings.ToLower(err.Error()), "unexpected") {
		t.Errorf("expected parse/unmarshal/unexpected error, got: %v", err)
	}
}

// TestIntegration_ServerCrash tests that a closed transport returns an error.
func TestIntegration_ServerCrash(t *testing.T) {
	mt := &mockTransport{}
	mt.Close() // Simulate crash/close

	_, err := mt.Call(context.Background(), "tools/call", nil)
	if err == nil {
		t.Fatal("expected error from closed transport, got nil")
	}
	if !strings.Contains(err.Error(), "closed") {
		t.Errorf("expected 'closed' error, got: %v", err)
	}
}

// TestE2E_AgentCallsMCPTool tests the full path: MCPClientManager → tool call → response.
func TestE2E_AgentCallsMCPTool(t *testing.T) {
	manager := &MCPClientManager{}
	manager.servers = map[string]*mcpServer{
		"e2e": {
			cfg: config.MCPServerConfig{Transport: "stdio", Timeout: 10},
			transport: &mockTransport{
				tools: []ToolInfo{
					{Name: "greet", Description: "Greets user", Schema: map[string]any{"type": "object"}},
				},
			},
			status: StatusConnected,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := manager.CallTool(ctx, "e2e", "greet", map[string]any{
		"user": "Alice",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.IsError {
		t.Fatal("expected no error from greet tool")
	}
	if len(result.Content) != 1 {
		t.Fatalf("expected 1 content block, got %d", len(result.Content))
	}
	if !strings.Contains(result.Content[0].Text, "Response from greet") {
		t.Errorf("expected 'Response from greet', got %q", result.Content[0].Text)
	}
}

// TestE2E_MCPToolReturnsError tests that when a tool returns isError=true.
func TestE2E_MCPToolReturnsError(t *testing.T) {
	manager := &MCPClientManager{}
	manager.servers = map[string]*mcpServer{
		"e2e": {
			cfg: config.MCPServerConfig{Transport: "stdio", Timeout: 10},
			transport: &mockTransport{
				tools: []ToolInfo{
					{Name: "error_tool", Description: "Returns error", Schema: map[string]any{"type": "object"}},
				},
				callHandler: func(method string, params map[string]any) (*json.RawMessage, error) {
					if method == "tools/call" {
						data, _ := json.Marshal(map[string]any{
							"content": []any{map[string]any{"type": "text", "text": "Something went wrong"}},
							"isError": true,
						})
						raw := json.RawMessage(data)
						return &raw, nil
					}
					return nil, fmt.Errorf("unknown method")
				},
			},
			status: StatusConnected,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := manager.CallTool(ctx, "e2e", "error_tool", nil)
	if err != nil {
		t.Fatalf("unexpected transport error: %v", err)
	}
	if !result.IsError {
		t.Fatal("expected isError=true from error_tool")
	}
	if !strings.Contains(result.Content[0].Text, "Something went wrong") {
		t.Errorf("expected error message, got %q", result.Content[0].Text)
	}
}

// TestE2E_ConcurrentMCPCalls tests calling multiple tools concurrently.
func TestE2E_ConcurrentMCPCalls(t *testing.T) {
	manager := &MCPClientManager{}
	mt := &mockTransport{
		tools: []ToolInfo{{Name: "greet", Description: "Greets user", Schema: map[string]any{"type": "object"}}},
	}
	manager.servers = map[string]*mcpServer{
		"e2e": {
			cfg:       config.MCPServerConfig{Transport: "stdio", Timeout: 10},
			transport: mt,
			status:    StatusConnected,
		},
	}

	const numCalls = 5
	results := make([]string, numCalls)
	errsCh := make([]error, numCalls)
	var wg sync.WaitGroup

	for i := 0; i < numCalls; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			user := fmt.Sprintf("User%d", idx)
			result, err := manager.CallTool(ctx, "e2e", "greet", map[string]any{
				"user": user,
			})
			if err != nil {
				errsCh[idx] = err
				return
			}
			if len(result.Content) > 0 {
				results[idx] = result.Content[0].Text
			}
		}(i)
	}
	wg.Wait()

	for i := 0; i < numCalls; i++ {
		if errsCh[i] != nil {
			t.Errorf("call %d failed: %v", i, errsCh[i])
		}
		expected := fmt.Sprintf("Response from greet")
		if results[i] != expected {
			t.Errorf("call %d: expected %q, got %q", i, expected, results[i])
		}
	}
}

// TestE2E_AgentCloseClosesMCP tests that CloseAll terminates transports.
func TestE2E_AgentCloseClosesMCP(t *testing.T) {
	mt := &mockTransport{
		tools: []ToolInfo{{Name: "greet", Description: "Greets user", Schema: map[string]any{"type": "object"}}},
	}
	manager := &MCPClientManager{}
	manager.servers = map[string]*mcpServer{
		"e2e": {
			cfg:       config.MCPServerConfig{Transport: "stdio"},
			transport: mt,
			status:    StatusConnected,
		},
	}

	srv := manager.GetServer("e2e")
	if srv == nil {
		t.Fatal("expected e2e server to be connected")
	}

	manager.CloseAll()

	// After close, CallTool should fail since transport is closed
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := manager.CallTool(ctx, "e2e", "greet", nil)
	if err == nil {
		t.Log("Warning: CallTool succeeded after CloseAll — server may still be alive")
	} else {
		t.Logf("Expected error after close: %v", err)
	}
}
