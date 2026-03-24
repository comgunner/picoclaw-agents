// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package tools

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/comgunner/picoclaw/pkg/health"
	"github.com/comgunner/picoclaw/pkg/security"
)

// selfStartTime captures process start for uptime reporting.
var selfStartTime = time.Now()

// SelfDiagnosticsTool provides NP-04 Level 1 read-only self-diagnosis:
// read_self_logs, get_self_status, list_self_config, check_self_health.
// It has no auto-modification capability.
type SelfDiagnosticsTool struct {
	workspace  string
	configPath string
}

// NewSelfDiagnosticsTool creates a SelfDiagnosticsTool.
// workspace is the agent workspace directory; configPath is optional (path to config.json).
func NewSelfDiagnosticsTool(workspace, configPath string) *SelfDiagnosticsTool {
	return &SelfDiagnosticsTool{workspace: workspace, configPath: configPath}
}

func (t *SelfDiagnosticsTool) Name() string { return "self_diagnostics" }

func (t *SelfDiagnosticsTool) Description() string {
	return "Read-only self-diagnosis for the agent: read its own logs, check runtime status, " +
		"list configuration, and perform a health check. Cannot modify any files or settings."
}

func (t *SelfDiagnosticsTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"mode": map[string]any{
				"type": "string",
				"description": "Diagnosis mode: " +
					"'read_self_logs' (recent agent/service logs), " +
					"'get_self_status' (runtime metrics, uptime, memory), " +
					"'list_self_config' (redacted configuration summary), " +
					"'check_self_health' (workspace, memory, goroutine health check).",
				"enum": []string{"read_self_logs", "get_self_status", "list_self_config", "check_self_health"},
			},
			"lines": map[string]any{
				"type":        "integer",
				"description": "Number of log lines to return (read_self_logs only). Default: 50.",
			},
			"level": map[string]any{
				"type":        "string",
				"description": "Filter by log level (read_self_logs only): 'ERROR', 'WARN', 'INFO', 'DEBUG'. Default: all levels.",
			},
		},
		"required": []string{"mode"},
	}
}

func (t *SelfDiagnosticsTool) Execute(_ context.Context, args map[string]any) *ToolResult {
	mode, ok := args["mode"].(string)
	if !ok {
		return ErrorResult("mode is required: read_self_logs | get_self_status | list_self_config | check_self_health")
	}

	switch mode {
	case "read_self_logs":
		lines := 50
		if l, ok := args["lines"].(float64); ok && l > 0 {
			lines = int(l)
		}
		levelFilter := ""
		if lv, ok := args["level"].(string); ok {
			levelFilter = strings.ToUpper(strings.TrimSpace(lv))
		}
		return t.readSelfLogs(lines, levelFilter)
	case "get_self_status":
		return t.getSelfStatus()
	case "list_self_config":
		return t.listSelfConfig()
	case "check_self_health":
		return t.checkSelfHealth()
	default:
		return ErrorResult(fmt.Sprintf("unknown mode: %q", mode))
	}
}

func (t *SelfDiagnosticsTool) readSelfLogs(lines int, levelFilter string) *ToolResult {
	var sb strings.Builder

	// 1. Try journalctl (Linux + systemd)
	if out, err := exec.Command("journalctl", "-u", "picoclaw", "-n", fmt.Sprintf("%d", lines), "--no-pager", "-o", "short").Output(); err == nil {
		content := string(out)
		if levelFilter != "" {
			content = selfFilterLines(content, levelFilter)
		}
		sb.WriteString("📋 **Logs (journalctl -u picoclaw):**\n```\n")
		sb.WriteString(content)
		sb.WriteString("\n```")
		return UserResult(sb.String())
	}

	// 2. Try workspace sibling picoclaw.log
	logPath := filepath.Join(t.workspace, "..", "picoclaw.log")
	if _, err := os.Stat(logPath); err == nil {
		if content, err := selfTailFile(logPath, lines); err == nil {
			if levelFilter != "" {
				content = selfFilterLines(content, levelFilter)
			}
			sb.WriteString(fmt.Sprintf("📋 **Logs (%s):**\n```\n", logPath))
			sb.WriteString(content)
			sb.WriteString("\n```")
			return UserResult(sb.String())
		}
	}

	// 3. Scan workspace for .log files
	if matches, _ := filepath.Glob(filepath.Join(t.workspace, "*.log")); len(matches) > 0 {
		if content, err := selfTailFile(matches[0], lines); err == nil {
			if levelFilter != "" {
				content = selfFilterLines(content, levelFilter)
			}
			sb.WriteString(fmt.Sprintf("📋 **Logs (%s):**\n```\n", matches[0]))
			sb.WriteString(content)
			sb.WriteString("\n```")
			return UserResult(sb.String())
		}
	}

	return UserResult("ℹ️ No log source found. PicoClaw logs to stdout by default. " +
		"To enable file logging, configure LOG_FILE in your environment or systemd service unit.")
}

func (t *SelfDiagnosticsTool) getSelfStatus() *ToolResult {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	uptime := time.Since(selfStartTime).Round(time.Second)
	metrics := health.ContextMetrics.GetMetrics()

	return UserResult(fmt.Sprintf(
		"🤖 **Agent Self-Status**\n\n"+
			"**Runtime:**\n"+
			"- PID: %d\n"+
			"- Uptime: %s\n"+
			"- Go: %s (%s/%s)\n"+
			"- Goroutines: %d\n\n"+
			"**Memory:**\n"+
			"- Heap alloc: %.1f MB\n"+
			"- Heap in-use: %.1f MB\n"+
			"- GC cycles: %d\n\n"+
			"**Context Metrics:**\n"+
			"- Compactions: %d\n"+
			"- Total tokens processed: %d\n"+
			"- Cache hits/misses: %d/%d\n"+
			"- Errors: %d\n",
		os.Getpid(),
		uptime,
		runtime.Version(), runtime.GOOS, runtime.GOARCH,
		runtime.NumGoroutine(),
		float64(ms.HeapAlloc)/1024/1024,
		float64(ms.HeapInuse)/1024/1024,
		ms.NumGC,
		metrics["compaction_count"],
		metrics["total_tokens"],
		metrics["cache_hits"], metrics["cache_misses"],
		metrics["errors"],
	))
}

func (t *SelfDiagnosticsTool) listSelfConfig() *ToolResult {
	path := t.configPath
	if path == "" {
		// Search common locations relative to workspace
		for _, candidate := range []string{
			filepath.Join(t.workspace, "..", "config.json"),
			filepath.Join(t.workspace, "..", "..", "config", "config.json"),
		} {
			if _, err := os.Stat(candidate); err == nil {
				path = candidate
				break
			}
		}
	}

	if path == "" {
		return UserResult("ℹ️ Config file not found. Showing compiled defaults:\n" +
			"- model: (empty — uses provider default)\n" +
			"- max_tokens: 32768\n" +
			"- restrict_to_workspace: true")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return ErrorResult(fmt.Sprintf("failed to read config: %v", err))
	}

	redacted := security.GlobalRedactor.Redact(string(data))
	return UserResult(fmt.Sprintf("📄 **Config** (`%s`):\n```json\n%s\n```", path, redacted))
}

func (t *SelfDiagnosticsTool) checkSelfHealth() *ToolResult {
	var sb strings.Builder
	sb.WriteString("🏥 **Self-Health Check**\n\n")

	// Workspace accessible
	if _, err := os.Stat(t.workspace); err == nil {
		sb.WriteString("✅ Workspace accessible\n")
	} else {
		sb.WriteString(fmt.Sprintf("❌ Workspace inaccessible: %v\n", err))
	}

	// Workspace writable
	probe := filepath.Join(t.workspace, ".self_health_probe")
	if err := os.WriteFile(probe, []byte("ok"), 0o600); err == nil {
		os.Remove(probe)
		sb.WriteString("✅ Workspace writable\n")
	} else {
		sb.WriteString(fmt.Sprintf("❌ Workspace not writable: %v\n", err))
	}

	// Context compaction error rate
	metrics := health.ContextMetrics.GetMetrics()
	if metrics["errors"] == 0 {
		sb.WriteString("✅ No context compaction errors\n")
	} else {
		sb.WriteString(fmt.Sprintf("⚠️  Context compaction errors: %d (out of %d compactions)\n",
			metrics["errors"], metrics["compaction_count"]))
	}

	// Goroutine count
	goroutines := runtime.NumGoroutine()
	if goroutines < 500 {
		sb.WriteString(fmt.Sprintf("✅ Goroutines: %d (healthy)\n", goroutines))
	} else {
		sb.WriteString(fmt.Sprintf("⚠️  High goroutine count: %d — possible leak\n", goroutines))
	}

	// Heap memory
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	heapMB := float64(ms.HeapAlloc) / 1024 / 1024
	if heapMB < 500 {
		sb.WriteString(fmt.Sprintf("✅ Heap memory: %.0f MB (normal)\n", heapMB))
	} else {
		sb.WriteString(fmt.Sprintf("⚠️  High heap memory: %.0f MB\n", heapMB))
	}

	return UserResult(sb.String())
}

// selfTailFile reads the last n lines of a text file.
func selfTailFile(path string, n int) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
		if len(lines) > n*2 {
			lines = lines[len(lines)-n:]
		}
	}
	if len(lines) > n {
		lines = lines[len(lines)-n:]
	}
	return strings.Join(lines, "\n"), scanner.Err()
}

// selfFilterLines keeps only lines containing the given substring (case-insensitive via ToUpper).
func selfFilterLines(content, filter string) string {
	var out []string
	for _, line := range strings.Split(content, "\n") {
		if strings.Contains(strings.ToUpper(line), filter) {
			out = append(out, line)
		}
	}
	return strings.Join(out, "\n")
}
