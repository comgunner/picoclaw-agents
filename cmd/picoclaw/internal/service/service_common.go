// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package service

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// ServiceConfig holds configuration for the service installation.
type ServiceConfig struct {
	Port       int    // Gateway port (default: 18800)
	Public     bool   // Listen on 0.0.0.0 (default: false)
	ConfigPath string // Path to config.json (default: ~/.picoclaw/config.json)
	DryRun     bool   // Show what would be done without executing (default: false)
}

// ServiceStatus represents the current state of the service.
type ServiceStatus struct {
	Installed bool   // Whether the service is installed
	Active    bool   // Whether the service is currently running
	PID       int    // Process ID (0 if not running)
	Port      int    // Configured port
	Error     string // Error message if any
}

// ServicePlatform defines the interface for platform-specific service management.
type ServicePlatform interface {
	Install() error
	DryRunInstall() error
	Uninstall() error
	Start() error
	Stop() error
	Status() (ServiceStatus, error)
	Logs(lines int, follow bool) error
}

// GetServicePlatform returns the appropriate platform implementation for the current OS.
func GetServicePlatform(cfg ServiceConfig) (ServicePlatform, error) {
	exePath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("cannot determine binary path: %w", err)
	}
	if _, err := os.Stat(exePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("binary not found at %s — service would fail to start", exePath)
	}

	// Resolve config path if not specified
	if cfg.ConfigPath == "" {
		home, _ := os.UserHomeDir()
		cfg.ConfigPath = filepath.Join(home, ".picoclaw", "config.json")
	}

	// Default port
	if cfg.Port == 0 {
		cfg.Port = 18800
	}

	// Use newPlatform which is defined per-platform via build tags
	return newPlatform(cfg, exePath)
}

// isPortAvailable checks if a TCP port is free on localhost.
func isPortAvailable(port int) bool {
	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return false
	}
	ln.Close()
	return true
}

// isGatewayRunning checks if something is already listening on the given port.
func isGatewayRunning(port int) bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", port), 2*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// writeServiceFileAtomic writes content to a file using atomic rename pattern.
// This prevents file corruption if the process crashes mid-write.
func writeServiceFileAtomic(path, content string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create service directory: %w", err)
	}

	// Write to temp file first
	tmpFile, err := os.CreateTemp(dir, "picoclaw-service-*.tmp")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()

	cleanup := true
	defer func() {
		if cleanup {
			tmpFile.Close()
			os.Remove(tmpPath)
		}
	}()

	if _, err := tmpFile.WriteString(content); err != nil {
		return fmt.Errorf("failed to write service file: %w", err)
	}

	if err := tmpFile.Chmod(0o644); err != nil {
		return fmt.Errorf("failed to set permissions: %w", err)
	}

	if err := tmpFile.Sync(); err != nil {
		return fmt.Errorf("failed to sync service file: %w", err)
	}

	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to close temp file: %w", err)
	}

	// Atomic rename
	if err := os.Rename(tmpPath, path); err != nil {
		return fmt.Errorf("failed to rename service file: %w", err)
	}

	cleanup = false
	return nil
}

// xmlEscape escapes special characters for XML/plist content.
func xmlEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	return s
}

// runCmdSilent executes a command and returns combined output + error.
func runCmdSilent(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// homeDir returns the user's home directory.
func homeDir() string {
	home, _ := os.UserHomeDir()
	return home
}

// serviceLabel returns the unique service label for this platform.
func serviceLabel() string {
	return "com.picoclaw.agents.gateway"
}
