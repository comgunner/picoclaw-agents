package tools

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/comgunner/picoclaw/pkg/utils"
)

// SubagentListTool lists current subagent tasks and their status.
type SubagentListTool struct {
	manager *SubagentManager
}

func NewSubagentListTool(manager *SubagentManager) *SubagentListTool {
	return &SubagentListTool{manager: manager}
}

func (t *SubagentListTool) Name() string {
	return "subagent_list"
}

func (t *SubagentListTool) Description() string {
	return "List all current subagent tasks and their status (running, completed, failed)."
}

func (t *SubagentListTool) Parameters() map[string]any {
	return map[string]any{
		"type":       "object",
		"properties": map[string]any{},
	}
}

func (t *SubagentListTool) Execute(ctx context.Context, args map[string]any) *ToolResult {
	if t.manager == nil {
		return ErrorResult("Subagent manager not configured")
	}

	tasks := t.manager.ListTasks()
	if len(tasks) == 0 {
		return UserResult("No hay tareas de subagente activas o registradas.")
	}

	var sb strings.Builder
	sb.WriteString("📋 **Tareas de Subagente:**\n\n")

	for _, task := range tasks {
		created := time.UnixMilli(task.Created).Format("15:04:05")
		statusEmoji := "⏳"
		if task.Status == "completed" {
			statusEmoji = "✅"
		} else if task.Status == "failed" {
			statusEmoji = "❌"
		}

		sb.WriteString(fmt.Sprintf("%s **%s** [%s] - %s\n", statusEmoji, task.ID, created, task.Label))
		sb.WriteString(fmt.Sprintf("  - Tarea: %s\n", utils.Truncate(task.Task, 60)))
		sb.WriteString(fmt.Sprintf("  - Estado: %s\n\n", task.Status))
	}

	return UserResult(sb.String())
}
