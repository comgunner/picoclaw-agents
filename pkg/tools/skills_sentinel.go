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
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/comgunner/picoclaw/pkg/logger"
	"github.com/comgunner/picoclaw/pkg/utils"
)

// SkillsSentinelTool provides internal, pattern-based security checks
// against prompt injection, system prompt extraction, and malicious code patterns.
type SkillsSentinelTool struct {
	blacklist        []*regexp.Regexp
	disabledUntil    time.Time
	disabledMu       sync.RWMutex
	onAutoReactivate func() // Callback for auto-reactivation notification
	workspace        string // Path to the picoClaw workspace
}

// SetAutoReactivateCallback sets a callback function to be called when sentinel auto-reactivates
func (t *SkillsSentinelTool) SetAutoReactivateCallback(fn func()) {
	t.onAutoReactivate = fn
}

// NewSkillsSentinelTool creates a new SkillsSentinelTool with a default set of malicious patterns.
func NewSkillsSentinelTool() *SkillsSentinelTool {
	// Patterns refined to reduce false positives while maintaining defense against common techniques
	patterns := []string{
		// Prompt Injection & System Extraction
		`(?i)ignore\s+previous\s+instructions`,
		`(?i)ignore\s+all\s+prior`,
		`(?i)forget\s+everything\s+above`,
		`(?i)disregard\s+above`,
		`(?i)override\s+system`,
		`(?i)(reveal|leak|print|show|output)\s+(me\s+)?(your\s+|the\s+)?system\s+(prompt|instructions)`,
		`(?i)(reveal|leak|print|show|output)\s+(me\s+)?(your\s+|the\s+)?instructions`,
		`(?i)(dump|leak|print|show|output)\s+(the\s+)?configuration`,
		`(?i)you\s+are\s+now\s+DAN`,
		`(?i)developer\s+mode`,
		`(?i)unrestricted\s+mode`,

		// ClickFix / Social Engineering (Downloads)
		`(?i)curl\s+.*\s*\|\s*(bash|sh)`,
		`(?i)wget\s+.*\s*\|\s*(bash|sh)`,
		`(?i)powershell\s+.*iex\s*\(`,
		`(?i)iwr\s+.*\s*-useb\s*\|\s*iex`,

		// RAT / Reverse Shell
		`(?i)bash\s+-i\s*>\s*&\s*/dev/tcp/`,
		`(?i)nc\s+.*\s*-e\s+/(bin/)?(bash|sh)`,
		`(?i)python\s+.*\s*socket.*connect`,

		// Info Stealer / Exfiltration
		`(?i)cat\s+.*\.ssh/id_rsa`,
		`(?i)security\s+find-(generic|internet)-password`,
		`(?i)history\s*\|\s*grep`,
		`(?i)env\s*\|\s*curl`,
	}

	compiled := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		re := regexp.MustCompile(p)
		compiled = append(compiled, re)
	}

	return &SkillsSentinelTool{
		blacklist:     compiled,
		disabledUntil: time.Time{},
	}
}

// SetWorkspace sets the workspace path for the sentinel
func (t *SkillsSentinelTool) SetWorkspace(path string) {
	t.workspace = path
}

// IsDisabled checks if the sentinel is currently disabled
func (t *SkillsSentinelTool) IsDisabled() bool {
	t.disabledMu.RLock()
	defer t.disabledMu.RUnlock()

	if t.disabledUntil.IsZero() {
		return false
	}

	return time.Now().Before(t.disabledUntil)
}

// Disable temporarily disables the sentinel for a specified duration
func (t *SkillsSentinelTool) Disable(duration time.Duration) {
	t.disabledMu.Lock()
	defer t.disabledMu.Unlock()

	t.disabledUntil = time.Now().Add(duration)
}

// Enable immediately enables the sentinel
func (t *SkillsSentinelTool) Enable() {
	t.disabledMu.Lock()
	defer t.disabledMu.Unlock()

	t.disabledUntil = time.Time{}
}

// GetStatus returns the current status of the sentinel
func (t *SkillsSentinelTool) GetStatus() string {
	t.disabledMu.RLock()
	defer t.disabledMu.RUnlock()

	if t.disabledUntil.IsZero() {
		return "active"
	}

	if time.Now().Before(t.disabledUntil) {
		remaining := time.Until(t.disabledUntil)
		return "disabled_" + remaining.Round(time.Second).String()
	}

	// Auto-reenable if time passed
	wasDisabled := !t.disabledUntil.IsZero()
	t.disabledUntil = time.Time{}

	// Trigger callback if set
	if wasDisabled && t.onAutoReactivate != nil {
		go t.onAutoReactivate()
	}

	return "active"
}

func (t *SkillsSentinelTool) Name() string {
	return "skills_sentinel"
}

func (t *SkillsSentinelTool) Description() string {
	return "Internal security tool that validates input against malicious patterns and performs security scans on skills."
}

func (t *SkillsSentinelTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"action": map[string]any{
				"type":        "string",
				"enum":        []string{"validate", "scan"},
				"description": "The action to perform: 'validate' for text check (default), 'scan' for skill scanning.",
			},
			"input": map[string]any{
				"type":        "string",
				"description": "The text to validate (required for 'validate' action).",
			},
		},
		"required": []string{},
	}
}

// Execute checks the input against the blacklist or performs a skill scan.
func (t *SkillsSentinelTool) Execute(ctx context.Context, args map[string]any) *ToolResult {
	action, _ := args["action"].(string)
	if action == "" {
		action = "validate"
	}

	if action == "scan" {
		return t.handleScan(ctx)
	}

	// Default: validate
	// Check if sentinel is disabled
	if t.IsDisabled() {
		return SilentResult("Sentinel temporarily disabled for configuration tasks.")
	}

	input, _ := args["input"].(string)
	inputLower := strings.ToLower(input)

	// "Self-Aware" Exception: Allow queries containing PicoClaw-specific keywords
	if strings.Contains(inputLower, "picoclaw") ||
		strings.Contains(inputLower, "herramienta") ||
		strings.Contains(inputLower, "tool") ||
		strings.Contains(inputLower, "sentinel") ||
		strings.Contains(inputLower, "skill") {
		// If it looks like a question or informational statement, it's likely a false positive.
		if strings.Contains(inputLower, "?") ||
			strings.Contains(inputLower, "¿") ||
			strings.Contains(inputLower, "saber") ||
			strings.Contains(inputLower, "know") ||
			strings.Contains(inputLower, "tienes") ||
			strings.Contains(inputLower, "have") ||
			strings.Contains(inputLower, "has") ||
			strings.Contains(inputLower, "como") ||
			strings.Contains(inputLower, "how") ||
			strings.Contains(inputLower, "qué") ||
			strings.Contains(inputLower, "what") {
			return SilentResult("Input identified as a legitimate system query.")
		}
	}

	for _, re := range t.blacklist {
		if re.MatchString(input) {
			return ErrorResult("Security Violation: Malicious pattern detected by internal sentinel.")
		}
	}

	return SilentResult("Input verified as safe.")
}

func (t *SkillsSentinelTool) handleScan(ctx context.Context) *ToolResult {
	logger.InfoCF("sentinel", "Starting security scan of installed skills", nil)

	var scanPaths []string
	homeHost, _ := os.UserHomeDir()

	// 1. Workspace skills (from config)
	if t.workspace != "" {
		wsPath := utils.ExpandPath(t.workspace)
		scanPaths = append(scanPaths, filepath.Join(wsPath, "skills"))
	}

	// 2. PicoClaw standard paths
	openClawBase := filepath.Join(homeHost, ".picoclaw")
	scanPaths = append(scanPaths, filepath.Join(openClawBase, "skills"))
	scanPaths = append(scanPaths, filepath.Join(openClawBase, "extensions"))

	// Node modules for skills (recursive glob-like search)
	nmPath := filepath.Join(openClawBase, "node_modules")

	findings := []string{}
	scannedCount := 0

	for _, basePath := range scanPaths {
		if _, err := os.Stat(basePath); os.IsNotExist(err) {
			continue
		}

		err := filepath.WalkDir(basePath, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			if !d.IsDir() {
				// Scan files in skill directories
				if t.isRelevantFile(path) {
					scannedCount++
					if warning := t.scanFile(path); warning != "" {
						findings = append(findings, warning)
					}
				}
			}
			return nil
		})
		if err != nil {
			logger.ErrorCF("sentinel", "Error walking path", map[string]any{"path": basePath, "error": err.Error()})
		}
	}

	// Special check for node_modules/@*/
	if _, err := os.Stat(nmPath); err == nil {
		filepath.WalkDir(nmPath, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			// Only go 2 levels deep for @scope/package
			rel, _ := filepath.Rel(nmPath, path)
			parts := strings.Split(rel, string(os.PathSeparator))
			if len(parts) > 2 {
				return fs.SkipDir
			}

			if !d.IsDir() && t.isRelevantFile(path) {
				scannedCount++
				if warning := t.scanFile(path); warning != "" {
					findings = append(findings, warning)
				}
			}
			return nil
		})
	}

	if len(findings) == 0 {
		return SilentResult(fmt.Sprintf("Scan complete. Scanned %d files. No malicious patterns found.", scannedCount))
	}

	report := fmt.Sprintf(
		"🛡️ **Skill Security Scan Report**\n\nScanned %d files. Found %d issues.\n\n",
		scannedCount,
		len(findings),
	)
	report += strings.Join(findings, "\n\n")
	report += "\n\n**Recommendation:** Review these files and remove any suspicious skills using `picoclaw skills remove <name>`."

	return ErrorResult(report)
}

func (t *SkillsSentinelTool) isRelevantFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	name := strings.ToLower(filepath.Base(path))
	return ext == ".js" || ext == ".ts" || ext == ".sh" || ext == ".py" || name == "skill.md" || name == "package.json"
}

func (t *SkillsSentinelTool) scanFile(path string) string {
	content, err := os.ReadFile(path)
	if err != nil {
		return ""
	}

	for _, re := range t.blacklist {
		if re.Match(content) {
			relPath := path
			if t.workspace != "" {
				if rel, err := filepath.Rel(t.workspace, path); err == nil {
					relPath = rel
				}
			}
			return fmt.Sprintf("🔴 **UNSAFE PATTERN DETECTED** in `%s`\nPattern: `%s`", relPath, re.String())
		}
	}
	return ""
}
