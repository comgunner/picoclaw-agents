package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/comgunner/picoclaw/pkg/utils"
)

// BenchTool allows the agent to benchmark command startup time and memory.
type BenchTool struct{}

func NewBenchTool(_ string) *BenchTool { return &BenchTool{} }

func (t *BenchTool) Name() string { return "bench" }

func (t *BenchTool) Description() string {
	return "Benchmark the startup time and memory usage of a command. " +
		"Use this to measure how fast and lightweight a binary or script is. " +
		"Example: bench the picoclaw-agents binary itself."
}

func (t *BenchTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"command": map[string]any{
				"type":        "string",
				"description": "Command to benchmark (e.g., './build/picoclaw-agents')",
			},
			"args": map[string]any{
				"type":        "string",
				"description": "Space-separated arguments (e.g., '--help')",
			},
		},
		"required": []string{"command"},
	}
}

func (t *BenchTool) Execute(ctx context.Context, args map[string]any) *ToolResult {
	command, _ := args["command"].(string)
	argStr, _ := args["args"].(string)

	if command == "" {
		return ErrorResult("command is required")
	}

	var cmdArgs []string
	if argStr != "" {
		cmdArgs = strings.Fields(argStr)
	}

	elapsed, snap, err := utils.BenchmarkStartup(command, cmdArgs)
	if err != nil {
		return ErrorResult(fmt.Sprintf("benchmark failed: %v", err))
	}

	result := fmt.Sprintf("✅ Benchmark complete:\n"+
		"  Startup time: %v\n"+
		"  Peak RSS: %d MB\n"+
		"  Alloc: %d MB\n"+
		"  TotalAlloc: %d MB\n"+
		"  GC cycles: %d\n"+
		"  Goroutines: %d",
		elapsed, snap.SysMB, snap.AllocMB, snap.TotalAllocMB, snap.NumGC, snap.Goroutines)
	return UserResult(result)
}
