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
	"encoding/json"
	"fmt"
	"strings"
)

// ToolCallResult is the parsed result from an MCP tool call.
type ToolCallResult struct {
	Content []ContentBlock `json:"content"`
	IsError bool           `json:"isError"`
}

// ContentBlock represents a single content item in an MCP tool response.
type ContentBlock struct {
	Type     string `json:"type"` // "text" or "image"
	Text     string `json:"text,omitempty"`
	Data     string `json:"data,omitempty"` // base64
	MIMEType string `json:"mimeType,omitempty"`
}

// parseToolCallResult parses raw JSON into a ToolCallResult.
func parseToolCallResult(raw json.RawMessage) (*ToolCallResult, error) {
	var result ToolCallResult
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("parse tool call result: %w", err)
	}
	return &result, nil
}

// FormatContent formats MCP content as a readable string.
func FormatContent(content []ContentBlock) string {
	var sb strings.Builder
	for _, block := range content {
		switch block.Type {
		case "text":
			sb.WriteString(block.Text)
		case "image":
			sb.WriteString(fmt.Sprintf("[image: %s, %s]", block.MIMEType, block.Data[:minInt(50, len(block.Data))]))
		default:
			sb.WriteString(fmt.Sprintf("[unknown content type: %s]", block.Type))
		}
	}
	return sb.String()
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
