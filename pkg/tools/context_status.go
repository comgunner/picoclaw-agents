// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package tools

import (
	"context"
	"fmt"

	pcontext "github.com/comgunner/picoclaw/pkg/context"
)

// ContextStatusTool provides visibility into token usage and context state.
type ContextStatusTool struct {
	middleware *pcontext.ContextMiddleware
}

// NewContextStatusTool creates a new context status tool.
func NewContextStatusTool(middleware *pcontext.ContextMiddleware) *ContextStatusTool {
	return &ContextStatusTool{
		middleware: middleware,
	}
}

// Name returns the tool identifier.
func (t *ContextStatusTool) Name() string {
	return "context_status"
}

// Description returns a brief description.
func (t *ContextStatusTool) Description() string {
	return "Check current token usage, context budget, and workspace state"
}

// Parameters returns the tool parameters.
func (t *ContextStatusTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"detail": map[string]any{
				"type":        "string",
				"enum":        []string{"summary", "detailed"},
				"default":     "summary",
				"description": "Level of detail: summary (quick stats) or detailed (full breakdown)",
			},
		},
		"required": []string{},
	}
}

// Execute executes the tool.
func (t *ContextStatusTool) Execute(ctx context.Context, args map[string]any) *ToolResult {
	detail, _ := args["detail"].(string)
	if detail == "" {
		detail = "summary"
	}

	current, maxTokens, percentage := t.middleware.GetUsage()

	var output string
	if detail == "detailed" {
		output = fmt.Sprintf(
			"📊 **Context Status (Detailed)**\n\n"+
				"**Token Budget:**\n"+
				"- Current: %d / %d tokens (%.1f%%)\n"+
				"- Threshold (80%%): %d tokens\n"+
				"- Status: %s\n\n"+
				"**Lazy Loader Cache:**\n"+
				"- Max cache size: %d\n\n"+
				"**GC Statistics:**\n"+
				"Last run: %s\n",
			current, maxTokens, percentage,
			maxTokens*8/10,
			getStatus(percentage),
			100,   // loader max cache size
			"N/A", // GC info
		)
	} else {
		output = fmt.Sprintf(
			"📊 **Context Status:** %d/%d tokens (%.1f%%) - %s",
			current, maxTokens, percentage,
			getStatus(percentage),
		)
	}

	return UserResult(output)
}

func getStatus(percentage float64) string {
	switch {
	case percentage >= 90:
		return "🔴 CRITICAL - Emergency GC triggered"
	case percentage >= 80:
		return "🟠 WARNING - Near threshold"
	case percentage >= 60:
		return "🟡 MODERATE - Monitor usage"
	default:
		return "🟢 HEALTHY"
	}
}
