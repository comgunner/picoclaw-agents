package tools

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/comgunner/picoclaw/pkg/utils"
)

// ArchLintTool allows the agent to check for forbidden import patterns.
type ArchLintTool struct {
	workspace string
}

func NewArchLintTool(workspace string) *ArchLintTool { return &ArchLintTool{workspace: workspace} }

func (t *ArchLintTool) Name() string { return "arch_lint" }

func (t *ArchLintTool) Description() string {
	return "Check for forbidden import patterns between packages. " +
		"Use this to ensure architectural boundaries are respected. " +
		"Detects cycles and cross-layer violations like agent importing channels."
}

func (t *ArchLintTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"root": map[string]any{
				"type":        "string",
				"description": "Root directory to scan (default: project root)",
			},
		},
	}
}

func (t *ArchLintTool) Execute(ctx context.Context, args map[string]any) *ToolResult {
	root := "."
	if r, ok := args["root"].(string); ok && r != "" {
		root = r
	}
	root, err := filepath.Abs(root)
	if err != nil {
		return ErrorResult(fmt.Sprintf("resolve root: %v", err))
	}

	violations, err := utils.CheckImports(root, nil)
	if err != nil {
		return ErrorResult(fmt.Sprintf("check imports: %v", err))
	}

	if len(violations) == 0 {
		return UserResult("✅ No import violations found. Architecture is clean.")
	}

	result := fmt.Sprintf("⚠️ Found %d import violation(s):\n", len(violations))
	for _, v := range violations {
		result += fmt.Sprintf("  %s → must not import %s\n    File: %s\n",
			v.Rule.From, v.Rule.MustNot, v.File)
	}
	return UserResult(result)
}
