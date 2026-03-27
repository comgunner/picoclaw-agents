// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package tools

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"

	"github.com/comgunner/picoclaw/pkg/logger"
)

// UsagePoint represents a single usage data point.
type UsagePoint struct {
	Timestamp string  `json:"timestamp"`
	CPU       float64 `json:"cpu_percent"`
	RAM       float64 `json:"ram_percent"`
}

// ResourceMonitorTool provides resource monitoring and threshold alerts.
type ResourceMonitorTool struct {
	workspace  string
	history    []UsagePoint
	historyMu  sync.RWMutex
	maxHistory int
}

// NewResourceMonitorTool creates a new ResourceMonitorTool instance.
func NewResourceMonitorTool(workspace string) *ResourceMonitorTool {
	return &ResourceMonitorTool{
		workspace:  workspace,
		history:    make([]UsagePoint, 0),
		maxHistory: 100, // Keep last 100 data points
	}
}

// Name returns the tool name.
func (t *ResourceMonitorTool) Name() string {
	return "resource_monitor"
}

// Description returns the tool description.
func (t *ResourceMonitorTool) Description() string {
	return "Monitor CPU/RAM usage with threshold alerts and automatic throttling recommendations. Track resource usage history and get optimization suggestions."
}

// Parameters returns the JSON schema for tool parameters.
func (t *ResourceMonitorTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"action": map[string]any{
				"type":        "string",
				"description": "Action: 'cpu_threshold', 'ram_threshold', 'throttle', 'history', or 'current'",
				"enum":        []string{"cpu_threshold", "ram_threshold", "throttle", "history", "current"},
			},
			"threshold": map[string]any{
				"type":        "number",
				"description": "Threshold percentage (0-100). Default: 80",
			},
			"hours": map[string]any{
				"type":        "integer",
				"description": "Hours of history to retrieve. Default: 1 (last hour)",
			},
		},
		"required": []string{"action"},
	}
}

// Execute runs the resource monitor tool.
func (t *ResourceMonitorTool) Execute(ctx context.Context, args map[string]any) *ToolResult {
	action, ok := args["action"].(string)
	if !ok {
		return ErrorResult(
			"action is required and must be one of: cpu_threshold, ram_threshold, throttle, history, current",
		)
	}

	threshold := 80.0
	if t, ok := args["threshold"].(float64); ok {
		threshold = t
	}

	hours := 1
	if h, ok := args["hours"].(float64); ok {
		hours = int(h)
	}

	switch action {
	case "cpu_threshold":
		return t.checkCPUThreshold(threshold)
	case "ram_threshold":
		return t.checkRAMThreshold(threshold)
	case "throttle":
		return t.getThrottleRecommendation(threshold)
	case "history":
		return t.getHistory(hours)
	case "current":
		return t.getCurrentUsage()
	default:
		return ErrorResult(
			fmt.Sprintf(
				"unknown action: %s. Valid options: cpu_threshold, ram_threshold, throttle, history, current",
				action,
			),
		)
	}
}

// checkCPUThreshold checks if CPU usage exceeds the threshold.
func (t *ResourceMonitorTool) checkCPUThreshold(threshold float64) *ToolResult {
	percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		logger.ErrorCF("tool", "Failed to get CPU usage",
			map[string]any{
				"tool":  "resource_monitor",
				"error": err.Error(),
			})
		return ErrorResult(fmt.Sprintf("failed to get CPU usage: %v", err))
	}

	currentCPU := percent[0]
	alert := currentCPU > threshold

	// Record in history
	t.recordUsage(currentCPU, getRAMUsagePercent())

	result := map[string]any{
		"cpu_percent": currentCPU,
		"threshold":   threshold,
		"alert":       alert,
		"timestamp":   time.Now().Format(time.RFC3339),
	}

	_ = result // For future structured output
	if alert {
		return SilentResult(fmt.Sprintf("⚠️ CPU ALERT: %.1f%% exceeds threshold of %.1f%%", currentCPU, threshold))
	}
	return SilentResult(fmt.Sprintf("CPU usage: %.1f%% (threshold: %.1f%%, OK)", currentCPU, threshold))
}

// checkRAMThreshold checks if RAM usage exceeds the threshold.
func (t *ResourceMonitorTool) checkRAMThreshold(threshold float64) *ToolResult {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		logger.ErrorCF("tool", "Failed to get RAM usage",
			map[string]any{
				"tool":  "resource_monitor",
				"error": err.Error(),
			})
		return ErrorResult(fmt.Sprintf("failed to get RAM usage: %v", err))
	}

	currentRAM := vmStat.UsedPercent
	alert := currentRAM > threshold

	// Record in history
	t.recordUsage(getCPUUsagePercent(), currentRAM)

	result := map[string]any{
		"ram_percent": currentRAM,
		"threshold":   threshold,
		"alert":       alert,
		"timestamp":   time.Now().Format(time.RFC3339),
	}

	_ = result // For future structured output
	if alert {
		return SilentResult(fmt.Sprintf("⚠️ RAM ALERT: %.1f%% exceeds threshold of %.1f%%", currentRAM, threshold))
	}
	return SilentResult(fmt.Sprintf("RAM usage: %.1f%% (threshold: %.1f%%, OK)", currentRAM, threshold))
}

// getThrottleRecommendation provides throttling recommendations based on resource usage.
func (t *ResourceMonitorTool) getThrottleRecommendation(threshold float64) *ToolResult {
	cpuPercent := getCPUUsagePercent()
	ramPercent := getRAMUsagePercent()

	recommendations := []string{}
	status := "OPTIMAL"

	if cpuPercent > threshold {
		recommendations = append(
			recommendations,
			fmt.Sprintf("- Reduce concurrent agent spawns (CPU at %.1f%%)", cpuPercent),
		)
		status = "THROTTLE_RECOMMENDED"
	}
	if ramPercent > threshold {
		recommendations = append(
			recommendations,
			fmt.Sprintf("- Limit context size or reduce agent count (RAM at %.1f%%)", ramPercent),
		)
		status = "THROTTLE_RECOMMENDED"
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations, "- No throttling needed, resources are healthy")
		status = "OPTIMAL"
	}

	result := map[string]any{
		"cpu_percent":     cpuPercent,
		"ram_percent":     ramPercent,
		"threshold":       threshold,
		"status":          status,
		"recommendations": recommendations,
		"timestamp":       time.Now().Format(time.RFC3339),
	}

	_ = result // For future structured output
	return SilentResult(fmt.Sprintf("Resource status: %s. CPU: %.1f%%, RAM: %.1f%%", status, cpuPercent, ramPercent))
}

// getHistory returns resource usage history.
func (t *ResourceMonitorTool) getHistory(hours int) *ToolResult {
	t.historyMu.RLock()
	defer t.historyMu.RUnlock()

	if len(t.history) == 0 {
		return SilentResult("No resource usage history available yet")
	}

	// Calculate average usage
	var totalCPU, totalRAM float64
	for _, point := range t.history {
		totalCPU += point.CPU
		totalRAM += point.RAM
	}

	avgCPU := totalCPU / float64(len(t.history))
	avgRAM := totalRAM / float64(len(t.history))

	// Find peak usage
	peakCPU := 0.0
	peakRAM := 0.0
	for _, point := range t.history {
		if point.CPU > peakCPU {
			peakCPU = point.CPU
		}
		if point.RAM > peakRAM {
			peakRAM = point.RAM
		}
	}

	result := map[string]any{
		"data_points":   len(t.history),
		"hours_covered": hours,
		"avg_cpu":       avgCPU,
		"avg_ram":       avgRAM,
		"peak_cpu":      peakCPU,
		"peak_ram":      peakRAM,
		"timestamp":     time.Now().Format(time.RFC3339),
	}

	_ = result // For future structured output
	return SilentResult(
		fmt.Sprintf(
			"Resource history: %d data points, Avg CPU: %.1f%%, Avg RAM: %.1f%%, Peak CPU: %.1f%%, Peak RAM: %.1f%%",
			len(t.history),
			avgCPU,
			avgRAM,
			peakCPU,
			peakRAM,
		),
	)
}

// getCurrentUsage returns current CPU and RAM usage.
func (t *ResourceMonitorTool) getCurrentUsage() *ToolResult {
	cpuPercent := getCPUUsagePercent()
	ramPercent := getRAMUsagePercent()

	// Record in history
	t.recordUsage(cpuPercent, ramPercent)

	result := map[string]any{
		"cpu_percent": cpuPercent,
		"ram_percent": ramPercent,
		"timestamp":   time.Now().Format(time.RFC3339),
	}

	_ = result // For future structured output
	return SilentResult(fmt.Sprintf("Current resources: CPU %.1f%%, RAM %.1f%%", cpuPercent, ramPercent))
}

// recordUsage adds a usage point to history.
func (t *ResourceMonitorTool) recordUsage(cpuPercent, ramPercent float64) {
	t.historyMu.Lock()
	defer t.historyMu.Unlock()

	point := UsagePoint{
		Timestamp: time.Now().Format(time.RFC3339),
		CPU:       cpuPercent,
		RAM:       ramPercent,
	}

	t.history = append(t.history, point)

	// Trim history if exceeds max
	if len(t.history) > t.maxHistory {
		t.history = t.history[len(t.history)-t.maxHistory:]
	}
}

// getCPUUsagePercent is a helper to get current CPU usage.
func getCPUUsagePercent() float64 {
	percent, err := cpu.Percent(time.Second, false)
	if err != nil || len(percent) == 0 {
		return 0
	}
	return percent[0]
}

// getRAMUsagePercent is a helper to get current RAM usage.
func getRAMUsagePercent() float64 {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return 0
	}
	return vmStat.UsedPercent
}
