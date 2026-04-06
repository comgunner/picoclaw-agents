package tools

import (
	"context"
	"fmt"

	"github.com/comgunner/picoclaw/pkg/utils"
)

// ReaperTool allows the agent to find and kill orphaned picoclaw-agents processes.
type ReaperTool struct{}

func NewReaperTool() *ReaperTool { return &ReaperTool{} }

func (t *ReaperTool) Name() string { return "reaper" }

func (t *ReaperTool) Description() string {
	return "Find and kill orphaned picoclaw-agents processes. " +
		"Use this when old agent processes are consuming memory after being stopped. " +
		"Call with action='find' to list orphans, or action='kill' to terminate them."
}

func (t *ReaperTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"action": map[string]any{
				"type":        "string",
				"enum":        []string{"find", "kill"},
				"description": "Action: 'find' lists orphans, 'kill' terminates them",
			},
		},
		"required": []string{"action"},
	}
}

func (t *ReaperTool) Execute(ctx context.Context, args map[string]any) *ToolResult {
	action, _ := args["action"].(string)
	if action == "" {
		return ErrorResult("action is required (find or kill)")
	}

	orphans, err := utils.FindOrphans()
	if err != nil {
		return ErrorResult(fmt.Sprintf("find orphans: %v", err))
	}

	if len(orphans) == 0 {
		return UserResult("✅ No orphan processes found.")
	}

	if action == "find" {
		result := fmt.Sprintf("Found %d orphan(s):\n", len(orphans))
		for _, o := range orphans {
			result += fmt.Sprintf("  PID %d: %s\n", o.PID, o.Cmd)
		}
		return UserResult(result)
	}

	killed, err := utils.KillOrphans()
	if err != nil {
		return ErrorResult(fmt.Sprintf("kill orphans: %v", err))
	}
	result := fmt.Sprintf("✅ Killed %d orphan(s):\n", len(killed))
	for _, o := range killed {
		result += fmt.Sprintf("  PID %d: %s\n", o.PID, o.Cmd)
	}
	return UserResult(result)
}
