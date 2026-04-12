// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

//go:build darwin

package service

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type darwinLaunchd struct {
	cfg     ServiceConfig
	exePath string
}

// newPlatform returns the darwin platform implementation.
func newPlatform(cfg ServiceConfig, exePath string) (ServicePlatform, error) {
	return &darwinLaunchd{cfg: cfg, exePath: exePath}, nil
}

func (p *darwinLaunchd) plistPath() string {
	return filepath.Join(homeDir(), "Library", "LaunchAgents", serviceLabel()+".plist")
}

func (p *darwinLaunchd) logDir() string {
	return filepath.Join(homeDir(), "Library", "Logs", "picoclaw-agents")
}

func (p *darwinLaunchd) buildLaunchdPlist() string {
	publicArg := ""
	if p.cfg.Public {
		publicArg = "<string>--public</string>"
	}

	return `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN"
  "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
  <key>Label</key>
  <string>` + serviceLabel() + `</string>
  <key>ProgramArguments</key>
  <array>
    <string>` + xmlEscape(p.exePath) + `</string>
    <string>gateway</string>
    <string>--port</string>
    <string>` + strconv.Itoa(p.cfg.Port) + `</string>
    ` + publicArg + `
  </array>
  <key>RunAtLoad</key>
  <true/>
  <key>KeepAlive</key>
  <dict>
    <key>SuccessfulExit</key>
    <false/>
  </dict>
  <key>ThrottleInterval</key>
  <integer>3</integer>
  <key>StandardOutPath</key>
  <string>` + xmlEscape(filepath.Join(p.logDir(), "gateway.out.log")) + `</string>
  <key>StandardErrorPath</key>
  <string>` + xmlEscape(filepath.Join(p.logDir(), "gateway.err.log")) + `</string>
  <key>WorkingDirectory</key>
  <string>` + xmlEscape(homeDir()) + `</string>
  <key>EnvironmentVariables</key>
  <dict>
    <key>HOME</key>
    <string>` + xmlEscape(homeDir()) + `</string>
    <key>PICOCLAW_HOME</key>
    <string>` + xmlEscape(filepath.Join(homeDir(), ".picoclaw")) + `</string>
  </dict>
</dict>
</plist>`
}

func (p *darwinLaunchd) preFlightChecks() error {
	if isGatewayRunning(p.cfg.Port) {
		return fmt.Errorf(
			"gateway is already running on port %d. Stop it first with 'picoclaw-agents gateway stop' or 'kill <PID>' to avoid duplicate connections (Telegram, Discord, etc.)",
			p.cfg.Port,
		)
	}
	if _, err := os.Stat(p.plistPath()); err == nil {
		return fmt.Errorf(
			"service is already installed at %s. Use 'picoclaw-agents service uninstall' first",
			p.plistPath(),
		)
	}
	if !isPortAvailable(p.cfg.Port) {
		return fmt.Errorf("port %d is already in use by another process", p.cfg.Port)
	}
	return nil
}

func (p *darwinLaunchd) DryRunInstall() error {
	if err := p.preFlightChecks(); err != nil {
		return err
	}

	fmt.Printf("Platform:       macOS (launchd)\n")
	fmt.Printf("Port:           %d\n", p.cfg.Port)
	fmt.Printf("Public:         %v\n", p.cfg.Public)
	fmt.Printf("Config:         %s\n", p.cfg.ConfigPath)
	fmt.Printf("Binary:         %s\n", p.exePath)
	fmt.Println()
	fmt.Println("Pre-flight checks:")
	fmt.Println("  ✓ Port", p.cfg.Port, "is available")
	fmt.Println("  ✓ Service is not already installed")
	fmt.Println("  ✓ Binary exists at", p.exePath)
	fmt.Println()
	fmt.Println("What would be done:")
	fmt.Printf("  1. Create directory: %s\n", p.logDir())
	fmt.Printf("  2. Write plist file: %s\n", p.plistPath())
	fmt.Println("     (backup existing to .bak if present)")
	fmt.Println("  3. Run: launchctl bootstrap gui/$(id -u)/", p.plistPath())
	fmt.Println("  4. Verify: launchctl list | grep", serviceLabel())
	fmt.Println()
	fmt.Println("Plist file content (would be written):")
	fmt.Println("──────────────────────────────────────")
	fmt.Println(p.buildLaunchdPlist())
	fmt.Println("──────────────────────────────────────")
	fmt.Println()
	fmt.Println("ℹ No changes were made. To install, run without --dry-run:")
	fmt.Printf("  picoclaw-agents service install --port %d\n", p.cfg.Port)
	return nil
}

func (p *darwinLaunchd) Install() error {
	if err := p.preFlightChecks(); err != nil {
		return err
	}

	// Create log directory
	if err := os.MkdirAll(p.logDir(), 0o755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	// Backup existing plist
	if _, err := os.Stat(p.plistPath()); err == nil {
		backup := p.plistPath() + ".bak"
		if err := os.Rename(p.plistPath(), backup); err != nil {
			fmt.Printf("⚠ Warning: could not backup existing plist: %v\n", err)
		} else {
			fmt.Println("✓ Backed up existing plist to .bak")
		}
	}

	// Write plist atomically
	plistContent := p.buildLaunchdPlist()
	if err := writeServiceFileAtomic(p.plistPath(), plistContent); err != nil {
		return fmt.Errorf("failed to write plist file: %w", err)
	}
	fmt.Println("✓ Created LaunchAgent:", p.plistPath())

	// Bootstrap service (try modern, fallback to legacy)
	uid := os.Getuid()
	if _, err := runCmdSilent("launchctl", "bootstrap", fmt.Sprintf("gui/%d/", uid), p.plistPath()); err != nil {
		// Fallback to legacy load -w
		if _, err2 := runCmdSilent("launchctl", "load", "-w", p.plistPath()); err2 != nil {
			return fmt.Errorf("failed to start service (bootstrap: %v, load: %v)", err, err2)
		}
	}
	fmt.Println("✓ Bootstrapped service")

	// Verify
	if _, err := runCmdSilent("launchctl", "list", serviceLabel()); err == nil {
		fmt.Println("✓ Service is running")
	} else {
		fmt.Println("⚠ Service started but verification failed (may need a moment to initialize)")
	}

	fmt.Println()
	fmt.Println("🎉 PicoClaw-Agents service installed successfully!")
	fmt.Printf("   Port: %d\n", p.cfg.Port)
	fmt.Printf("   Logs: %s\n", filepath.Join(p.logDir(), "gateway.out.log"))
	fmt.Println()
	fmt.Println("Manage the service with:")
	fmt.Println("  picoclaw-agents service status   Check status")
	fmt.Println("  picoclaw-agents service logs -f  Follow logs")
	fmt.Println("  picoclaw-agents service stop     Stop service")
	fmt.Println("  picoclaw-agents service uninstall Remove service")
	return nil
}

func (p *darwinLaunchd) Uninstall() error {
	// Stop first
	p.Stop()

	// Remove plist
	if _, err := os.Stat(p.plistPath()); err == nil {
		if err := os.Remove(p.plistPath()); err != nil {
			return fmt.Errorf("failed to remove plist: %w", err)
		}
		fmt.Println("✓ Removed plist file")
	}

	// Ask about log cleanup
	fmt.Print("Remove service logs? [y/N]: ")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		answer := strings.ToLower(scanner.Text())
		if answer == "y" || answer == "yes" {
			if err := os.RemoveAll(p.logDir()); err != nil {
				fmt.Printf("⚠ Warning: could not remove logs: %v\n", err)
			} else {
				fmt.Println("✓ Removed service logs")
			}
		}
	}

	fmt.Println("✓ Service uninstalled")
	return nil
}

func (p *darwinLaunchd) Start() error {
	if isGatewayRunning(p.cfg.Port) {
		return fmt.Errorf("gateway is already running on port %d", p.cfg.Port)
	}

	uid := os.Getuid()
	if _, err := runCmdSilent("launchctl", "bootstrap", fmt.Sprintf("gui/%d/", uid), p.plistPath()); err != nil {
		if _, err2 := runCmdSilent("launchctl", "load", "-w", p.plistPath()); err2 != nil {
			return fmt.Errorf("failed to start service: %v", err)
		}
	}
	fmt.Println("✓ Service started")
	return nil
}

func (p *darwinLaunchd) Stop() error {
	uid := os.Getuid()
	runCmdSilent("launchctl", "bootout", fmt.Sprintf("gui/%d/", uid), p.plistPath())
	runCmdSilent("launchctl", "stop", serviceLabel())
	fmt.Println("✓ Service stopped")
	return nil
}

func (p *darwinLaunchd) Status() (ServiceStatus, error) {
	status := ServiceStatus{Port: p.cfg.Port}

	// Check if installed (plist file exists)
	if _, statErr := os.Stat(p.plistPath()); os.IsNotExist(statErr) {
		return status, nil
	}
	status.Installed = true

	// Check if running
	out, err := runCmdSilent("launchctl", "list", serviceLabel())
	if err == nil && strings.Contains(out, serviceLabel()) {
		status.Active = true
		// Try to find PID
		if lines := strings.Split(out, "\n"); len(lines) > 1 {
			fields := strings.Fields(lines[0])
			if len(fields) >= 2 {
				if pid, err := strconv.Atoi(fields[0]); err == nil && pid > 0 {
					status.PID = pid
				}
			}
		}
	}

	return status, nil
}

func (p *darwinLaunchd) Logs(lines int, follow bool) error {
	logFile := filepath.Join(p.logDir(), "gateway.out.log")
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		return fmt.Errorf("no logs found at %s", logFile)
	}

	if follow {
		cmd := exec.Command("tail", "-f", "-n", strconv.Itoa(lines), logFile)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	cmd := exec.Command("tail", "-n", strconv.Itoa(lines), logFile)
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

// restart is not exposed via cobra but used internally
func (p *darwinLaunchd) restart() error {
	if err := p.Stop(); err != nil {
		return err
	}
	// Small delay to ensure port is released
	time.Sleep(1 * time.Second)
	return p.Start()
}
