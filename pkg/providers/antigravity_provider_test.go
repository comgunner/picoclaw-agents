package providers

import "testing"

func TestBuildRequestUsesFunctionFieldsWhenToolCallNameMissing(t *testing.T) {
	p := &AntigravityProvider{}

	messages := []Message{
		{
			Role: "assistant",
			ToolCalls: []ToolCall{{
				ID: "call_read_file_123",
				Function: &FunctionCall{
					Name:      "read_file",
					Arguments: `{"path":"README.md"}`,
				},
			}},
		},
		{
			Role:       "tool",
			ToolCallID: "call_read_file_123",
			Content:    "ok",
		},
	}

	req := p.buildRequest(messages, nil, "", nil)
	if len(req.Contents) != 2 {
		t.Fatalf("expected 2 contents, got %d", len(req.Contents))
	}

	modelPart := req.Contents[0].Parts[0]
	if modelPart.FunctionCall == nil {
		t.Fatal("expected functionCall in assistant message")
	}
	if modelPart.FunctionCall.Name != "read_file" {
		t.Fatalf("expected functionCall name read_file, got %q", modelPart.FunctionCall.Name)
	}
	if got := modelPart.FunctionCall.Args["path"]; got != "README.md" {
		t.Fatalf("expected functionCall args[path] to be README.md, got %v", got)
	}

	toolPart := req.Contents[1].Parts[0]
	if toolPart.FunctionResponse == nil {
		t.Fatal("expected functionResponse in tool message")
	}
	if toolPart.FunctionResponse.Name != "read_file" {
		t.Fatalf("expected functionResponse name read_file, got %q", toolPart.FunctionResponse.Name)
	}
}

func TestResolveToolResponseNameInfersNameFromGeneratedCallID(t *testing.T) {
	got := resolveToolResponseName("call_search_docs_999", map[string]string{})
	if got != "search_docs" {
		t.Fatalf("expected inferred tool name search_docs, got %q", got)
	}
}

// Test 1: Schema sanitization
func TestSanitizeSchemaForGemini_ReplacesAnyType(t *testing.T) {
	input := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"value": map[string]any{
				"type": "any", // ❌ Debería ser reemplazado
			},
		},
	}

	output := sanitizeSchemaForGemini(input)

	props := output["properties"].(map[string]any)
	value := props["value"].(map[string]any)

	if value["type"] != "object" {
		t.Errorf("expected 'object', got '%v'", value["type"])
	}
}

// Test 2: Invalid types
func TestSanitizeSchemaForGemini_InvalidTypes(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"any", "object"},
		{"", "object"},
		{"invalid", "object"},
		{"string", "string"},
		{"number", "number"},
	}

	for _, tt := range tests {
		input := map[string]any{"type": tt.input}
		output := sanitizeSchemaForGemini(input)

		if output["type"] != tt.expected {
			t.Errorf("type %q: expected %q, got %q", tt.input, tt.expected, output["type"])
		}
	}
}

// Test 3: Tool response parsing
func TestBuildRequest_ToolResponse(t *testing.T) {
	messages := []Message{
		{Role: "user", Content: "test"},
		{Role: "assistant", ToolCalls: []ToolCall{
			{ID: "call_test_tool_1", Name: "test_tool"},
		}},
		{Role: "tool", Content: "result", ToolCallID: "call_test_tool_1"},
	}

	provider := NewAntigravityProvider()
	req := provider.buildRequest(messages, []ToolDefinition{}, "gemini-3-flash", nil)

	// Verify tool response is in request
	if len(req.Contents) < 3 {
		t.Errorf("expected at least 3 contents, got %d", len(req.Contents))
	}

	// Last content should be tool response
	lastContent := req.Contents[len(req.Contents)-1]
	if lastContent.Role != "user" {
		t.Errorf("expected tool response role 'user', got '%s'", lastContent.Role)
	}
	if lastContent.Parts[0].FunctionResponse == nil {
		t.Error("expected FunctionResponse in tool response")
	}
}
