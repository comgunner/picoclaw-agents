package tools

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/comgunner/picoclaw/pkg/utils"
)

// MdAuditTool allows the agent to audit markdown files for broken links.
type MdAuditTool struct {
	workspace string
}

func NewMdAuditTool(workspace string) *MdAuditTool { return &MdAuditTool{workspace: workspace} }

func (t *MdAuditTool) Name() string { return "md_audit" }

func (t *MdAuditTool) Description() string {
	return "Scan Markdown files for broken internal links. " +
		"Use this to verify that documentation links are valid. " +
		"Scans the docs/ directory by default."
}

func (t *MdAuditTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"dir": map[string]any{
				"type":        "string",
				"description": "Directory to scan for .md files (default: docs/)",
			},
		},
	}
}

func (t *MdAuditTool) Execute(ctx context.Context, args map[string]any) *ToolResult {
	dir := "docs"
	if d, ok := args["dir"].(string); ok && d != "" {
		dir = d
	}
	dir, err := filepath.Abs(dir)
	if err != nil {
		return ErrorResult(fmt.Sprintf("resolve dir: %v", err))
	}

	issues, err := utils.AuditMarkdown(dir)
	if err != nil {
		return ErrorResult(fmt.Sprintf("audit: %v", err))
	}

	if len(issues) == 0 {
		return UserResult("✅ No broken internal links found in docs/.")
	}

	result := fmt.Sprintf("⚠️ Found %d link issue(s):\n", len(issues))
	for _, iss := range issues {
		result += fmt.Sprintf("  %s:%d — %s\n", iss.File, iss.Line, iss.Link)
	}
	return UserResult(result)
}
