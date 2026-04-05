// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)

package tools

import (
	"context"
	"fmt"
	"runtime"
	"sort"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"

	"github.com/comgunner/picoclaw/pkg/logger"
)

// SystemDiagnosticsTool provides real-time system status information.
type SystemDiagnosticsTool struct {
	workspace string
}

// NewSystemDiagnosticsTool creates a new SystemDiagnosticsTool instance.
func NewSystemDiagnosticsTool(workspace string) *SystemDiagnosticsTool {
	return &SystemDiagnosticsTool{
		workspace: workspace,
	}
}

// Name returns the tool name.
func (t *SystemDiagnosticsTool) Name() string {
	return "system_diagnostics"
}

// Description returns the tool description.
func (t *SystemDiagnosticsTool) Description() string {
	return "Get real-time system diagnostics: CPU, RAM, disk usage, processes, and recent logs."
}

// Parameters returns the JSON schema for tool parameters.
func (t *SystemDiagnosticsTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"metric": map[string]any{
				"type":        "string",
				"description": "Metric to retrieve: 'cpu', 'ram', 'disk', 'processes', 'logs', or 'all'",
				"enum":        []string{"cpu", "ram", "disk", "processes", "logs", "all"},
			},
			"limit": map[string]any{
				"type":        "integer",
				"description": "Limit results (e.g., top N processes, last N log lines). Default: 10",
			},
		},
		"required": []string{"metric"},
	}
}

// Execute runs the system diagnostics tool.
func (t *SystemDiagnosticsTool) Execute(ctx context.Context, args map[string]any) *ToolResult {
	metric, ok := args["metric"].(string)
	if !ok {
		return ErrorResult("metric is required and must be one of: cpu, ram, disk, processes, logs, all")
	}

	limit := 10
	if l, ok := args["limit"].(float64); ok {
		limit = int(l)
	}

	switch metric {
	case "cpu":
		return t.getCPUUsage()
	case "ram":
		return t.getRAMUsage()
	case "disk":
		return t.getDiskUsage()
	case "processes":
		return t.getProcesses(limit)
	case "logs":
		return t.getRecentLogs(limit)
	case "all":
		return t.getAllDiagnostics()
	default:
		return ErrorResult(
			fmt.Sprintf("unknown metric: %s. Valid options: cpu, ram, disk, processes, logs, all", metric),
		)
	}
}

// getCPUUsage returns current CPU usage percentage.
func (t *SystemDiagnosticsTool) getCPUUsage() *ToolResult {
	percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		logger.ErrorCF("tool", "Failed to get CPU usage",
			map[string]any{
				"tool":  "system_diagnostics",
				"error": err.Error(),
			})
		return ErrorResult(fmt.Sprintf("failed to get CPU usage: %v", err))
	}

	// Handle empty result (can happen on some systems like macOS ARM64)
	if len(percent) == 0 {
		logger.WarnCF("tool", "CPU percent returned empty slice",
			map[string]any{
				"tool": "system_diagnostics",
			})
		return SilentResult("CPU Usage: unavailable (system does not provide CPU metrics)")
	}

	result := map[string]any{
		"cpu_percent": percent[0],
		"cores":       runtime.NumCPU(),
		"timestamp":   time.Now().Format(time.RFC3339),
	}

	_ = result // For future structured output
	return SilentResult(fmt.Sprintf("CPU Usage: %.2f%% (%d cores)", percent[0], runtime.NumCPU()))
}

// getRAMUsage returns current RAM usage statistics.
func (t *SystemDiagnosticsTool) getRAMUsage() *ToolResult {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		logger.ErrorCF("tool", "Failed to get RAM usage",
			map[string]any{
				"tool":  "system_diagnostics",
				"error": err.Error(),
			})
		return ErrorResult(fmt.Sprintf("failed to get RAM usage: %v", err))
	}

	result := map[string]any{
		"total_mb":     float64(vmStat.Total / 1024 / 1024),
		"used_mb":      float64(vmStat.Used / 1024 / 1024),
		"available_mb": float64(vmStat.Available / 1024 / 1024),
		"used_percent": vmStat.UsedPercent,
		"timestamp":    time.Now().Format(time.RFC3339),
	}

	_ = result // For future structured output
	return SilentResult(fmt.Sprintf("RAM: %.0fMB / %.0fMB (%.1f%% used)",
		float64(vmStat.Used/1024/1024), float64(vmStat.Total/1024/1024), vmStat.UsedPercent))
}

// getDiskUsage returns disk usage for all mount points.
func (t *SystemDiagnosticsTool) getDiskUsage() *ToolResult {
	partitions, err := disk.Partitions(false)
	if err != nil {
		logger.ErrorCF("tool", "Failed to get disk partitions",
			map[string]any{
				"tool":  "system_diagnostics",
				"error": err.Error(),
			})
		return ErrorResult(fmt.Sprintf("failed to get disk info: %v", err))
	}

	diskInfo := make([]map[string]any, 0, len(partitions))
	for _, p := range partitions {
		usage, err := disk.Usage(p.Mountpoint)
		if err != nil {
			continue // Skip inaccessible mount points
		}
		diskInfo = append(diskInfo, map[string]any{
			"mountpoint":   p.Mountpoint,
			"total_gb":     usage.Total / 1024 / 1024 / 1024,
			"used_gb":      usage.Used / 1024 / 1024 / 1024,
			"free_gb":      usage.Free / 1024 / 1024 / 1024,
			"used_percent": usage.UsedPercent,
		})
	}

	result := map[string]any{
		"partitions": diskInfo,
		"timestamp":  time.Now().Format(time.RFC3339),
	}

	_ = result // For future structured output
	return SilentResult(fmt.Sprintf("Disk usage: %d partitions found", len(diskInfo)))
}

// getProcesses returns top N processes by CPU/memory usage.
func (t *SystemDiagnosticsTool) getProcesses(limit int) *ToolResult {
	procs, err := process.Processes()
	if err != nil {
		logger.ErrorCF("tool", "Failed to get processes",
			map[string]any{
				"tool":  "system_diagnostics",
				"error": err.Error(),
			})
		return ErrorResult(fmt.Sprintf("failed to get processes: %v", err))
	}

	type ProcInfo struct {
		PID    int32   `json:"pid"`
		Name   string  `json:"name"`
		CPU    float64 `json:"cpu_percent"`
		Memory float64 `json:"memory_percent"`
		Status string  `json:"status"`
	}

	procList := make([]ProcInfo, 0, len(procs))
	for _, p := range procs {
		name, _ := p.Name()
		cpuPercent, _ := p.CPUPercent()
		memPercent, _ := p.MemoryPercent()
		statuses, _ := p.Status()

		var status string
		if len(statuses) > 0 {
			status = statuses[0]
		}

		procList = append(procList, ProcInfo{
			PID:    p.Pid,
			Name:   name,
			CPU:    float64(cpuPercent),
			Memory: float64(memPercent),
			Status: status,
		})
	}

	// Sort by CPU usage (descending)
	sort.Slice(procList, func(i, j int) bool {
		return procList[i].CPU > procList[j].CPU
	})

	// Limit results
	if len(procList) > limit {
		procList = procList[:limit]
	}

	result := map[string]any{
		"processes": procList,
		"count":     len(procList),
		"timestamp": time.Now().Format(time.RFC3339),
	}

	_ = result // For future structured output
	return SilentResult(fmt.Sprintf("Top %d processes by CPU usage", len(procList)))
}

// getRecentLogs reads recent system logs (implementation depends on OS).
func (t *SystemDiagnosticsTool) getRecentLogs(limit int) *ToolResult {
	// Placeholder: In production, this would read /var/log/syslog, journalctl, etc.
	// For security, this may require elevated privileges or be disabled by default
	return SilentResult("Log reading requires elevated privileges. Configure log access in config.json.")
}

// getAllDiagnostics returns all diagnostics in a single call.
func (t *SystemDiagnosticsTool) getAllDiagnostics() *ToolResult {
	// Combine all metrics into a single result
	cpuPercent, _ := cpu.Percent(time.Second, false)
	vmStat, _ := mem.VirtualMemory()

	// Handle empty CPU result (can happen on some systems like macOS ARM64)
	var cpuData map[string]any
	var cpuStr string
	if len(cpuPercent) > 0 {
		cpuData = map[string]any{
			"percent": cpuPercent[0],
			"cores":   runtime.NumCPU(),
		}
		cpuStr = fmt.Sprintf("%.2f%%", cpuPercent[0])
	} else {
		cpuData = map[string]any{
			"percent": "unavailable",
			"cores":   runtime.NumCPU(),
		}
		cpuStr = "unavailable"
	}

	result := map[string]any{
		"cpu": cpuData,
		"ram": map[string]any{
			"total_mb":     vmStat.Total / 1024 / 1024,
			"used_mb":      vmStat.Used / 1024 / 1024,
			"used_percent": vmStat.UsedPercent,
		},
		"timestamp": time.Now().Format(time.RFC3339),
	}

	_ = result // For future structured output
	return SilentResult(fmt.Sprintf("System diagnostics: CPU %s, RAM %.1f%% used",
		cpuStr, vmStat.UsedPercent))
}
