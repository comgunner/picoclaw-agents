// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package tools

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/comgunner/picoclaw/pkg/config"
)

type ExecTool struct {
	workingDir          string
	timeout             time.Duration
	denyPatterns        []*regexp.Regexp
	allowPatterns       []*regexp.Regexp
	restrictToWorkspace bool
}

// DefaultDenyPatterns son los patrones de comandos peligrosos bloqueados por defecto
var DefaultDenyPatterns = []string{
	`(?i)\brm\s+(-[rf]+\s+)?/`,            // rm -rf /
	`(?i)\bdel\s+/f\s`,                    // Windows del
	`(?i)\bformat\s`,                      // format
	`(?i)\bmkfs`,                          // mkfs
	`(?i)\bdd\s+if=`,                      // dd if=
	`(?i)\bshutdown\b`,                    // shutdown
	`(?i)\breboot\b`,                      // reboot
	`(?i)\bpoweroff\b`,                    // poweroff
	`(?i):\(\)\s*{\s*:\|:&\s*}`,           // Fork bomb
	`(?i)/dev/sd[a-z]`,                    // Disk writes
	`(?i)\bchmod\s+[0-7]{3,4}\s+/`,        // chmod dangerous paths
	`(?i)\bchown\s+.*:/`,                  // chown dangerous paths
	`(?i)(tail|journalctl).*?\s+-f\b`,     // journalctl -f, tail -f (infinite loop)
	`(?i)\b(top|htop|nano|vim|vi|less)\b`, // Interactive TTY commands (infinite block)
}

func NewExecTool(workingDir string, restrict bool) (*ExecTool, error) {
	return NewExecToolWithConfig(workingDir, restrict, nil)
}

func NewExecToolWithConfig(workingDir string, restrict bool, config *config.Config) (*ExecTool, error) {
	rawPatterns := DefaultDenyPatterns

	if config != nil && len(config.Tools.Exec.CustomDenyPatterns) > 0 {
		rawPatterns = config.Tools.Exec.CustomDenyPatterns
	}

	if len(rawPatterns) == 0 {
		return nil, fmt.Errorf("deny patterns cannot be empty - security risk")
	}

	denyPatterns := make([]*regexp.Regexp, 0, len(rawPatterns))
	for _, pattern := range rawPatterns {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("invalid deny pattern %q: %w", pattern, err)
		}
		denyPatterns = append(denyPatterns, re)
	}

	// Validate working directory
	resolvedWD := workingDir
	if workingDir != "" {
		abs, err := filepath.Abs(workingDir)
		if err != nil {
			return nil, fmt.Errorf("invalid working directory: %w", err)
		}
		if isPathBlacklisted(abs) {
			return nil, fmt.Errorf("invalid working directory: access to sensitive system file blocked")
		}
		resolvedWD = abs
	}

	return &ExecTool{
		workingDir:          resolvedWD,
		timeout:             60 * time.Second,
		denyPatterns:        denyPatterns,
		allowPatterns:       nil,
		restrictToWorkspace: restrict,
	}, nil
}

func (t *ExecTool) Name() string {
	return "exec"
}

func (t *ExecTool) Description() string {
	return "Execute a shell command and return its output. Use with caution."
}

func (t *ExecTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"command": map[string]any{
				"type":        "string",
				"description": "The shell command to execute",
			},
			"working_dir": map[string]any{
				"type":        "string",
				"description": "Optional working directory for the command",
			},
		},
		"required": []string{"command"},
	}
}

// blockedToolNames son nombres de tools internas que NUNCA deben ejecutarse como shell commands.
var blockedToolNames = map[string]bool{
	"image_gen_antigravity":          true,
	"imagegenantigravity":            true, // variant sin underscores
	"image_gen_create":               true,
	"image_gen_workflow":             true,
	"text_script_create":             true,
	"social_post_bundle":             true,
	"social_manager":                 true,
	"community_manager_create_draft": true,
	"subagent":                       true,
	"spawn":                          true,
	"queue":                          true,
	"batch_id":                       true,
	"find_skills":                    true,
	"install_skill":                  true,
	"agent_list":                     true,
	"agent_default":                  true,
	"agent_receive":                  true,
	"self_diagnostics":               true,
	"system_diagnostics":             true,
	"resource_monitor":               true,
	"memory_store":                   true,
	"version_control":                true,
	"config_manager":                 true,
}

func (t *ExecTool) Execute(ctx context.Context, args map[string]any) *ToolResult {
	command, ok := args["command"].(string)
	if !ok {
		return ErrorResult("command is required")
	}

	// Block internal tool names from being executed as shell commands
	fields := strings.Fields(command)
	if len(fields) > 0 {
		firstWord := fields[0]
		if blockedToolNames[firstWord] {
			return ErrorResult(fmt.Sprintf(
				"Command '%s' is an internal tool, not a shell command. Use the tool directly, not via exec.",
				firstWord))
		}
	}

	cwd := t.workingDir
	if wd, ok := args["working_dir"].(string); ok && wd != "" {
		if t.restrictToWorkspace && t.workingDir != "" {
			resolvedWD, err := validatePath(wd, t.workingDir, true)
			if err != nil {
				return ErrorResult("Command blocked by safety guard (" + err.Error() + ")")
			}
			cwd = resolvedWD
		} else {
			cwd = wd
		}
	}

	if cwd == "" {
		wd, err := os.Getwd()
		if err == nil {
			cwd = wd
		}
	}

	if guardError := t.guardCommand(command, cwd); guardError != "" {
		return ErrorResult(guardError)
	}

	// timeout == 0 means no timeout
	var cmdCtx context.Context
	var cancel context.CancelFunc
	if t.timeout > 0 {
		cmdCtx, cancel = context.WithTimeout(ctx, t.timeout)
	} else {
		cmdCtx, cancel = context.WithCancel(ctx)
	}
	defer cancel()

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.CommandContext(cmdCtx, "powershell", "-NoProfile", "-NonInteractive", "-Command", command)
	} else {
		cmd = exec.CommandContext(cmdCtx, "sh", "-c", command)
	}
	if cwd != "" {
		cmd.Dir = cwd
	}

	prepareCommandForTermination(cmd)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		return ErrorResult(fmt.Sprintf("failed to start command: %v", err))
	}

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	var err error
	select {
	case err = <-done:
	case <-cmdCtx.Done():
		_ = terminateProcessTree(cmd)
		select {
		case err = <-done:
		case <-time.After(2 * time.Second):
			if cmd.Process != nil {
				_ = cmd.Process.Kill()
			}
			err = <-done
		}
	}

	output := stdout.String()
	if stderr.Len() > 0 {
		output += "\nSTDERR:\n" + stderr.String()
	}

	if err != nil {
		if errors.Is(cmdCtx.Err(), context.DeadlineExceeded) {
			msg := fmt.Sprintf("Command timed out after %v", t.timeout)
			return &ToolResult{
				ForLLM:  msg,
				ForUser: msg,
				IsError: true,
			}
		}
		output += fmt.Sprintf("\nExit code: %v", err)
	}

	if output == "" {
		output = "(no output)"
	}

	maxLen := 10000
	if len(output) > maxLen {
		output = output[:maxLen] + fmt.Sprintf("\n... (truncated, %d more chars)", len(output)-maxLen)
	}

	if err != nil {
		return &ToolResult{
			ForLLM:  output,
			ForUser: output,
			IsError: true,
		}
	}

	return &ToolResult{
		ForLLM:  output,
		ForUser: output,
		IsError: false,
	}
}

func (t *ExecTool) guardCommand(command, cwd string) string {
	cmd := strings.TrimSpace(command)
	lower := strings.ToLower(cmd)

	if (strings.HasPrefix(lower, "ping ") || lower == "ping") && !strings.Contains(lower, " -c") {
		return "Command blocked by safety guard (ping must use -c flag for infinite loop prevention)"
	}

	for _, pattern := range t.denyPatterns {
		if pattern.MatchString(lower) {
			return "Command blocked by safety guard (dangerous pattern detected)"
		}
	}

	if len(t.allowPatterns) > 0 {
		allowed := false
		for _, pattern := range t.allowPatterns {
			if pattern.MatchString(lower) {
				allowed = true
				break
			}
		}
		if !allowed {
			return "Command blocked by safety guard (not in allowlist)"
		}
	}

	if t.restrictToWorkspace {
		if strings.Contains(cmd, "..\\") || strings.Contains(cmd, "../") {
			return "Command blocked by safety guard (path traversal detected)"
		}

		cwdPath, err := filepath.Abs(cwd)
		if err != nil {
			return ""
		}

		// Match potential file paths: Windows paths (C:\...) and Unix paths (/...)
		// Also match file:// URIs which need special handling
		pathPattern := regexp.MustCompile(`(?:file://[^\s\"']+|[A-Za-z]:\\[^\\\"'\s]+|/[^\s\"']+)`)
		matches := pathPattern.FindAllStringIndex(cmd, -1)

		for _, match := range matches {
			raw := cmd[match[0]:match[1]]

			// Skip URL path components (e.g., https://github.com, ftp://...)
			// Only skip if preceded by a web URL scheme
			if strings.HasPrefix(raw, "//") {
				// Check if this is part of a web URL (http, https, ftp)
				precedingContext := ""
				if match[0] > 0 {
					start := match[0] - 10
					if start < 0 {
						start = 0
					}
					precedingContext = strings.ToLower(cmd[start:match[0]])
				}
				if strings.HasSuffix(precedingContext, "http:") ||
					strings.HasSuffix(precedingContext, "https:") ||
					strings.HasSuffix(precedingContext, "ftp:") {
					continue
				}
			}

			// Handle file:// URIs - extract the actual path and check it
			if strings.HasPrefix(raw, "file://") {
				raw = raw[7:] // Remove "file://" prefix
			}

			p, err := filepath.Abs(raw)
			if err != nil {
				continue
			}

			rel, err := filepath.Rel(cwdPath, p)
			if err != nil {
				continue
			}

			if strings.HasPrefix(rel, "..") {
				return "Command blocked by safety guard (path outside working dir)"
			}
		}
	}

	return ""
}

func (t *ExecTool) SetTimeout(timeout time.Duration) {
	t.timeout = timeout
}

func (t *ExecTool) SetRestrictToWorkspace(restrict bool) {
	t.restrictToWorkspace = restrict
}

func (t *ExecTool) SetAllowPatterns(patterns []string) error {
	t.allowPatterns = make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return fmt.Errorf("invalid allow pattern %q: %w", p, err)
		}
		t.allowPatterns = append(t.allowPatterns, re)
	}
	return nil
}
