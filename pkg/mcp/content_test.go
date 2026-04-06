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
	"strings"
	"testing"
)

func TestFormatContent_Text(t *testing.T) {
	content := []ContentBlock{
		{Type: "text", Text: "Hello, world!"},
	}

	result := FormatContent(content)
	if result != "Hello, world!" {
		t.Errorf("FormatContent(%v) = %q, want %q", content, result, "Hello, world!")
	}
}

func TestFormatContent_MultipleText(t *testing.T) {
	content := []ContentBlock{
		{Type: "text", Text: "Line 1"},
		{Type: "text", Text: "Line 2"},
	}

	result := FormatContent(content)
	expected := "Line 1Line 2"
	if result != expected {
		t.Errorf("FormatContent(%v) = %q, want %q", content, result, expected)
	}
}

func TestFormatContent_Image(t *testing.T) {
	data := strings.Repeat("a", 100)
	content := []ContentBlock{
		{Type: "image", Data: data, MIMEType: "image/png"},
	}

	result := FormatContent(content)
	if !strings.HasPrefix(result, "[image: image/png, ") {
		t.Errorf("FormatContent(%v) = %q, should start with [image: image/png,", content, result)
	}
}

func TestFormatContent_Unknown(t *testing.T) {
	content := []ContentBlock{
		{Type: "video"},
	}

	result := FormatContent(content)
	expected := "[unknown content type: video]"
	if result != expected {
		t.Errorf("FormatContent(%v) = %q, want %q", content, result, expected)
	}
}

func TestParseToolCallResult(t *testing.T) {
	raw := []byte(`{"content":[{"type":"text","text":"result"}],"isError":false}`)

	result, err := parseToolCallResult(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Content) != 1 {
		t.Fatalf("expected 1 content block, got %d", len(result.Content))
	}
	if result.Content[0].Text != "result" {
		t.Errorf("expected content text 'result', got %q", result.Content[0].Text)
	}
	if result.IsError {
		t.Error("expected IsError to be false")
	}
}
