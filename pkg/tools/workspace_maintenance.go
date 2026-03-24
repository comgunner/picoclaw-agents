package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/comgunner/picoclaw/pkg/logger"
)

// WorkspaceMaintenanceTool performs workspace cleanup in a single tool call.
type WorkspaceMaintenanceTool struct {
	workspace string
}

// MaintenanceArgs defines workspace maintenance options.
type MaintenanceArgs struct {
	DryRun        bool `json:"dry_run,omitempty"`
	CleanSessions bool `json:"clean_sessions,omitempty"`
	CleanLogs     bool `json:"clean_logs,omitempty"`
	CleanTemp     bool `json:"clean_temp,omitempty"`
	SessionAge    int  `json:"session_age_days,omitempty"`
	LogAge        int  `json:"log_age_days,omitempty"`
	TempAge       int  `json:"temp_age_days,omitempty"`
}

// MaintenanceResult contains structured maintenance execution results.
type MaintenanceResult struct {
	SessionsArchived int      `json:"sessions_archived"`
	LogsCompressed   int      `json:"logs_compressed"`
	TempFilesDeleted int      `json:"temp_files_deleted"`
	SpaceFreedBytes  int64    `json:"space_freed_bytes"`
	ExecutionTimeMs  int64    `json:"execution_time_ms"`
	OperatingSystem  string   `json:"operating_system"`
	WorkspacePath    string   `json:"workspace_path"`
	Warnings         []string `json:"warnings,omitempty"`
}

func NewWorkspaceMaintenanceTool(workspace string) *WorkspaceMaintenanceTool {
	return &WorkspaceMaintenanceTool{workspace: workspace}
}

func (t *WorkspaceMaintenanceTool) SetWorkspacePath(workspace string) {
	if strings.TrimSpace(workspace) != "" {
		t.workspace = workspace
	}
}

func (t *WorkspaceMaintenanceTool) Name() string {
	return "workspace_maintenance"
}

func (t *WorkspaceMaintenanceTool) Description() string {
	return "Run workspace cleanup in one operation: archive old sessions to workspace/cold with gzip, compress old logs, and prune old temp files. Use this instead of exec for maintenance."
}

func (t *WorkspaceMaintenanceTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"dry_run": map[string]any{
				"type":        "boolean",
				"description": "Preview cleanup actions without changing files",
			},
			"clean_sessions": map[string]any{
				"type":        "boolean",
				"description": "Move old session files to workspace/cold and gzip them",
			},
			"clean_logs": map[string]any{
				"type":        "boolean",
				"description": "Compress old .log files to .gz",
			},
			"clean_temp": map[string]any{
				"type":        "boolean",
				"description": "Delete old files from workspace/temp",
			},
			"session_age_days": map[string]any{
				"type":        "integer",
				"description": "Age in days for session archival (default: 7)",
			},
			"log_age_days": map[string]any{
				"type":        "integer",
				"description": "Age in days for log compression (default: 1)",
			},
			"temp_age_days": map[string]any{
				"type":        "integer",
				"description": "Age in days for temp cleanup (default: 3)",
			},
		},
	}
}

func (t *WorkspaceMaintenanceTool) Execute(ctx context.Context, args map[string]any) *ToolResult {
	start := time.Now()

	parsedArgs, err := parseMaintenanceArgs(args)
	if err != nil {
		return ErrorResult(err.Error())
	}

	workspace, err := t.resolveWorkspace()
	if err != nil {
		return ErrorResult(err.Error())
	}

	if err := os.MkdirAll(workspace, 0o755); err != nil {
		return ErrorResult(fmt.Sprintf("unable to create workspace: %v", err))
	}

	var output string
	switch runtime.GOOS {
	case "windows":
		output, err = t.executeWindows(ctx, workspace, parsedArgs)
	default:
		output, err = t.executeUnix(ctx, workspace, parsedArgs)
	}

	durationMs := time.Since(start).Milliseconds()
	result := t.parseOutput(output, workspace, durationMs, runtime.GOOS)

	logger.InfoCF("workspace_maintenance", "workspace maintenance cycle finished",
		map[string]any{
			"workspace":         workspace,
			"dry_run":           parsedArgs.DryRun,
			"sessions_archived": result.SessionsArchived,
			"logs_compressed":   result.LogsCompressed,
			"temp_deleted":      result.TempFilesDeleted,
			"space_freed_bytes": result.SpaceFreedBytes,
			"duration_ms":       durationMs,
			"under_1s":          durationMs < 1000,
		},
	)

	if err != nil {
		return ErrorResult(fmt.Sprintf("workspace maintenance failed: %v\n%s", err, output))
	}

	return SilentResult(result.FormatReport(parsedArgs.DryRun) + "\n\nJSON:\n" + result.JSON())
}

func parseMaintenanceArgs(args map[string]any) (MaintenanceArgs, error) {
	parsed := MaintenanceArgs{
		SessionAge: 7,
		LogAge:     1,
		TempAge:    3,
	}

	if len(args) == 0 {
		parsed.CleanSessions = true
		parsed.CleanLogs = true
		parsed.CleanTemp = true
		return parsed, nil
	}

	raw, err := json.Marshal(args)
	if err != nil {
		return parsed, fmt.Errorf("invalid args: %w", err)
	}
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return parsed, fmt.Errorf("invalid args: %w", err)
	}

	if parsed.SessionAge <= 0 {
		parsed.SessionAge = 7
	}
	if parsed.LogAge <= 0 {
		parsed.LogAge = 1
	}
	if parsed.TempAge <= 0 {
		parsed.TempAge = 3
	}
	if !parsed.CleanSessions && !parsed.CleanLogs && !parsed.CleanTemp {
		parsed.CleanSessions = true
		parsed.CleanLogs = true
		parsed.CleanTemp = true
	}

	return parsed, nil
}

func (t *WorkspaceMaintenanceTool) resolveWorkspace() (string, error) {
	workspace := expandHome(strings.TrimSpace(t.workspace))
	if workspace == "" {
		return "", fmt.Errorf("workspace path is not configured")
	}

	abs, err := filepath.Abs(workspace)
	if err != nil {
		return "", fmt.Errorf("failed to resolve workspace path: %w", err)
	}

	return filepath.Clean(abs), nil
}

func (t *WorkspaceMaintenanceTool) executeUnix(ctx context.Context, workspace string, args MaintenanceArgs) (string, error) {
	cmd := exec.CommandContext(ctx, "bash", "-c", t.generateBashScript())
	cmd.Dir = workspace
	cmd.Env = append(os.Environ(),
		"WORKSPACE="+workspace,
		"DRY_RUN="+boolEnv(args.DryRun),
		"CLEAN_SESSIONS="+boolEnv(args.CleanSessions),
		"CLEAN_LOGS="+boolEnv(args.CleanLogs),
		"CLEAN_TEMP="+boolEnv(args.CleanTemp),
		"SESSION_AGE="+strconv.Itoa(args.SessionAge),
		"LOG_AGE="+strconv.Itoa(args.LogAge),
		"TEMP_AGE="+strconv.Itoa(args.TempAge),
	)

	out, err := cmd.CombinedOutput()
	return string(out), err
}

func (t *WorkspaceMaintenanceTool) executeWindows(ctx context.Context, workspace string, args MaintenanceArgs) (string, error) {
	cmd := exec.CommandContext(ctx, "powershell", "-NoProfile", "-NonInteractive", "-Command", t.generatePowerShellScript())
	cmd.Dir = workspace
	cmd.Env = append(os.Environ(),
		"WORKSPACE="+workspace,
		"DRY_RUN="+boolEnv(args.DryRun),
		"CLEAN_SESSIONS="+boolEnv(args.CleanSessions),
		"CLEAN_LOGS="+boolEnv(args.CleanLogs),
		"CLEAN_TEMP="+boolEnv(args.CleanTemp),
		"SESSION_AGE="+strconv.Itoa(args.SessionAge),
		"LOG_AGE="+strconv.Itoa(args.LogAge),
		"TEMP_AGE="+strconv.Itoa(args.TempAge),
	)

	out, err := cmd.CombinedOutput()
	return string(out), err
}

func (t *WorkspaceMaintenanceTool) generateBashScript() string {
	return `set -eu
WORKSPACE="${WORKSPACE:-.}"
DRY_RUN="${DRY_RUN:-0}"
CLEAN_SESSIONS="${CLEAN_SESSIONS:-1}"
CLEAN_LOGS="${CLEAN_LOGS:-1}"
CLEAN_TEMP="${CLEAN_TEMP:-1}"
SESSION_AGE="${SESSION_AGE:-7}"
LOG_AGE="${LOG_AGE:-1}"
TEMP_AGE="${TEMP_AGE:-3}"

SESSIONS_ARCHIVED=0
LOGS_COMPRESSED=0
TEMP_FILES_DELETED=0
SPACE_FREED_BYTES=0

COLD_DIR="$WORKSPACE/cold"
SESSIONS_DIR="$WORKSPACE/sessions"
TEMP_DIR="$WORKSPACE/temp"
mkdir -p "$COLD_DIR" "$TEMP_DIR"

if [ "$CLEAN_SESSIONS" = "1" ] && [ -d "$SESSIONS_DIR" ]; then
  while IFS= read -r -d '' file; do
    if [ "$DRY_RUN" = "1" ]; then
      SESSIONS_ARCHIVED=$((SESSIONS_ARCHIVED+1))
      continue
    fi
    base="$(basename "$file")"
    target="$COLD_DIR/$base"
    i=0
    while [ -e "$target" ] || [ -e "$target.gz" ]; do
      i=$((i+1))
      target="$COLD_DIR/$base.$i"
    done
    mv "$file" "$target" || continue
    if gzip -f "$target"; then
      SESSIONS_ARCHIVED=$((SESSIONS_ARCHIVED+1))
    fi
  done < <(find "$SESSIONS_DIR" -type f -name "*.json" -mtime +"$SESSION_AGE" -print0 2>/dev/null)
fi

if [ "$CLEAN_LOGS" = "1" ]; then
  while IFS= read -r -d '' file; do
    if [ "$DRY_RUN" = "1" ]; then
      LOGS_COMPRESSED=$((LOGS_COMPRESSED+1))
      continue
    fi
    if gzip -f "$file"; then
      LOGS_COMPRESSED=$((LOGS_COMPRESSED+1))
    fi
  done < <(find "$WORKSPACE" -type f -name "*.log" -mtime +"$LOG_AGE" -print0 2>/dev/null)
fi

if [ "$CLEAN_TEMP" = "1" ] && [ -d "$TEMP_DIR" ]; then
  while IFS= read -r -d '' file; do
    if [ "$DRY_RUN" = "1" ]; then
      TEMP_FILES_DELETED=$((TEMP_FILES_DELETED+1))
      continue
    fi
    size="$(wc -c < "$file" 2>/dev/null || echo 0)"
    if rm -f "$file"; then
      TEMP_FILES_DELETED=$((TEMP_FILES_DELETED+1))
      SPACE_FREED_BYTES=$((SPACE_FREED_BYTES+size))
    fi
  done < <(find "$TEMP_DIR" -type f -mtime +"$TEMP_AGE" -print0 2>/dev/null)
fi

echo "SESSIONS_ARCHIVED: $SESSIONS_ARCHIVED"
echo "LOGS_COMPRESSED: $LOGS_COMPRESSED"
echo "TEMP_FILES_DELETED: $TEMP_FILES_DELETED"
echo "SPACE_FREED_BYTES: $SPACE_FREED_BYTES"
`
}

func (t *WorkspaceMaintenanceTool) generatePowerShellScript() string {
	return `$ErrorActionPreference = "Stop"
$WORKSPACE = if ($env:WORKSPACE) { $env:WORKSPACE } else { "." }
$DRY_RUN = [int]$env:DRY_RUN
$CLEAN_SESSIONS = [int]$env:CLEAN_SESSIONS
$CLEAN_LOGS = [int]$env:CLEAN_LOGS
$CLEAN_TEMP = [int]$env:CLEAN_TEMP
$SESSION_AGE = [int]$env:SESSION_AGE
$LOG_AGE = [int]$env:LOG_AGE
$TEMP_AGE = [int]$env:TEMP_AGE

$SESSIONS_ARCHIVED = 0
$LOGS_COMPRESSED = 0
$TEMP_FILES_DELETED = 0
$SPACE_FREED_BYTES = 0

$COLD_DIR = Join-Path $WORKSPACE "cold"
$SESSIONS_DIR = Join-Path $WORKSPACE "sessions"
$TEMP_DIR = Join-Path $WORKSPACE "temp"
New-Item -ItemType Directory -Path $COLD_DIR -Force | Out-Null
New-Item -ItemType Directory -Path $TEMP_DIR -Force | Out-Null

function Compress-GzipFile {
	param([string]$Path)
	$GzPath = "$Path.gz"
	$Input = [System.IO.File]::OpenRead($Path)
	$Output = [System.IO.File]::Create($GzPath)
	$Gzip = New-Object System.IO.Compression.GzipStream($Output, [System.IO.Compression.CompressionMode]::Compress)
	$Input.CopyTo($Gzip)
	$Gzip.Dispose()
	$Output.Dispose()
	$Input.Dispose()
	Remove-Item -LiteralPath $Path -Force
}

if ($CLEAN_SESSIONS -eq 1 -and (Test-Path -LiteralPath $SESSIONS_DIR)) {
	Get-ChildItem -LiteralPath $SESSIONS_DIR -File -Filter "*.json" |
	Where-Object { $_.LastWriteTime -lt (Get-Date).AddDays(-$SESSION_AGE) } |
	ForEach-Object {
		if ($DRY_RUN -eq 1) {
			$SESSIONS_ARCHIVED++
		} else {
			$Target = Join-Path $COLD_DIR $_.Name
			Move-Item -LiteralPath $_.FullName -Destination $Target -Force
			Compress-GzipFile -Path $Target
			$SESSIONS_ARCHIVED++
		}
	}
}

if ($CLEAN_LOGS -eq 1) {
	Get-ChildItem -LiteralPath $WORKSPACE -File -Filter "*.log" -Recurse |
	Where-Object { $_.LastWriteTime -lt (Get-Date).AddDays(-$LOG_AGE) } |
	ForEach-Object {
		if ($DRY_RUN -eq 1) {
			$LOGS_COMPRESSED++
		} else {
			Compress-GzipFile -Path $_.FullName
			$LOGS_COMPRESSED++
		}
	}
}

if ($CLEAN_TEMP -eq 1 -and (Test-Path -LiteralPath $TEMP_DIR)) {
	Get-ChildItem -LiteralPath $TEMP_DIR -File |
	Where-Object { $_.LastWriteTime -lt (Get-Date).AddDays(-$TEMP_AGE) } |
	ForEach-Object {
		if ($DRY_RUN -eq 1) {
			$TEMP_FILES_DELETED++
		} else {
			$SPACE_FREED_BYTES += $_.Length
			Remove-Item -LiteralPath $_.FullName -Force
			$TEMP_FILES_DELETED++
		}
	}
}

Write-Output "SESSIONS_ARCHIVED: $SESSIONS_ARCHIVED"
Write-Output "LOGS_COMPRESSED: $LOGS_COMPRESSED"
Write-Output "TEMP_FILES_DELETED: $TEMP_FILES_DELETED"
Write-Output "SPACE_FREED_BYTES: $SPACE_FREED_BYTES"
`
}

func (t *WorkspaceMaintenanceTool) parseOutput(output, workspace string, durationMs int64, osType string) *MaintenanceResult {
	result := &MaintenanceResult{
		WorkspacePath:   workspace,
		ExecutionTimeMs: durationMs,
		OperatingSystem: osType,
	}

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "SESSIONS_ARCHIVED":
			result.SessionsArchived = parseInt(value)
		case "LOGS_COMPRESSED":
			result.LogsCompressed = parseInt(value)
		case "TEMP_FILES_DELETED":
			result.TempFilesDeleted = parseInt(value)
		case "SPACE_FREED_BYTES":
			result.SpaceFreedBytes = parseInt64(value)
		case "WARNING":
			result.Warnings = append(result.Warnings, value)
		}
	}

	return result
}

func (r *MaintenanceResult) FormatReport(dryRun bool) string {
	var b strings.Builder

	if dryRun {
		b.WriteString("Workspace maintenance dry run\n\n")
	} else {
		b.WriteString("Workspace maintenance completed\n\n")
	}

	b.WriteString(fmt.Sprintf("OS: %s\n", r.OperatingSystem))
	b.WriteString(fmt.Sprintf("Workspace: %s\n", r.WorkspacePath))
	b.WriteString(fmt.Sprintf("Execution time: %dms\n", r.ExecutionTimeMs))
	b.WriteString(fmt.Sprintf("Sessions archived: %d\n", r.SessionsArchived))
	b.WriteString(fmt.Sprintf("Logs compressed: %d\n", r.LogsCompressed))
	b.WriteString(fmt.Sprintf("Temp files deleted: %d\n", r.TempFilesDeleted))
	b.WriteString(fmt.Sprintf("Space freed: %s\n", FormatBytes(r.SpaceFreedBytes)))

	if len(r.Warnings) > 0 {
		b.WriteString("Warnings:\n")
		for _, warning := range r.Warnings {
			b.WriteString("- " + warning + "\n")
		}
	}

	return strings.TrimSpace(b.String())
}

func (r *MaintenanceResult) JSON() string {
	raw, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return "{}"
	}
	return string(raw)
}

func FormatBytes(size int64) string {
	const (
		kb = 1024
		mb = kb * 1024
		gb = mb * 1024
	)

	switch {
	case size >= gb:
		return fmt.Sprintf("%.2f GB", float64(size)/float64(gb))
	case size >= mb:
		return fmt.Sprintf("%.2f MB", float64(size)/float64(mb))
	case size >= kb:
		return fmt.Sprintf("%.2f KB", float64(size)/float64(kb))
	default:
		return fmt.Sprintf("%d bytes", size)
	}
}

func parseInt(value string) int {
	n, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil {
		return 0
	}
	return n
}

func parseInt64(value string) int64 {
	n, err := strconv.ParseInt(strings.TrimSpace(value), 10, 64)
	if err != nil {
		return 0
	}
	return n
}

func boolEnv(v bool) string {
	if v {
		return "1"
	}
	return "0"
}

func expandHome(path string) string {
	if path == "" {
		return path
	}
	if path[0] == '~' {
		home, _ := os.UserHomeDir()
		if len(path) > 1 && (path[1] == '/' || path[1] == '\\') {
			return home + path[1:]
		}
		return home
	}
	return path
}
