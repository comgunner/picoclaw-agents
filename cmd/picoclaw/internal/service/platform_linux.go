// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

//go:build linux

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

type linuxSystemd struct {
	cfg     ServiceConfig
	exePath string
}

// newPlatform returns the linux platform implementation.
func newPlatform(cfg ServiceConfig, exePath string) (ServicePlatform, error) {
	return &linuxSystemd{cfg: cfg, exePath: exePath}, nil
}

func (p *linuxSystemd) unitDir() string {
	return filepath.Join(homeDir(), ".config", "systemd", "user")
}

func (p *linuxSystemd) unitPath() string {
	return filepath.Join(p.unitDir(), "picoclaw-agents.service")
}

func (p *linuxSystemd) buildSystemdUnit() string {
	publicFlag := ""
	if p.cfg.Public {
		publicFlag = "--public"
	}

	return `[Unit]
Description=PicoClaw-Agents Gateway Service
After=network.target

[Service]
Type=simple
ExecStart=` + p.exePath + ` gateway --port ` + strconv.Itoa(p.cfg.Port) + ` ` + publicFlag + `
WorkingDirectory=` + homeDir() + `
Restart=always
RestartSec=5
RestartPreventExitStatus=78
TimeoutStopSec=30
KillMode=control-group
Environment=HOME=` + homeDir() + `
Environment=PICOCLAW_HOME=` + filepath.Join(homeDir(), ".picoclaw") + `
PassEnvironment=DISPLAY XDG_RUNTIME_DIR

[Install]
WantedBy=default.target
`
}

func (p *linuxSystemd) preFlightChecks() error {
	if isGatewayRunning(p.cfg.Port) {
		return fmt.Errorf(
			"gateway is already running on port %d. Stop it first with 'picoclaw-agents gateway stop' or 'kill <PID>' to avoid duplicate connections (Telegram, Discord, etc.)",
			p.cfg.Port,
		)
	}
	if _, err := os.Stat(p.unitPath()); err == nil {
		return fmt.Errorf(
			"service is already installed at %s. Use 'picoclaw-agents service uninstall' first",
			p.unitPath(),
		)
	}
	if !isPortAvailable(p.cfg.Port) {
		return fmt.Errorf("port %d is already in use by another process", p.cfg.Port)
	}
	return nil
}

func (p *linuxSystemd) DryRunInstall() error {
	if err := p.preFlightChecks(); err != nil {
		return err
	}

	fmt.Printf("Platform:       Linux (systemd)\n")
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
	fmt.Printf("  1. Create directory: %s\n", p.unitDir())
	fmt.Printf("  2. Write unit file: %s\n", p.unitPath())
	fmt.Println("     (backup existing to .bak if present)")
	fmt.Println("  3. Run: systemctl --user daemon-reload")
	fmt.Println("  4. Run: systemctl --user enable picoclaw-agents.service")
	fmt.Println("  5. Run: systemctl --user start picoclaw-agents.service")
	fmt.Println("  6. Verify: systemctl --user is-active picoclaw-agents.service")
	fmt.Println()
	fmt.Println("Unit file content (would be written):")
	fmt.Println("──────────────────────────────────────")
	fmt.Print(p.buildSystemdUnit())
	fmt.Println("──────────────────────────────────────")
	fmt.Println()
	fmt.Println("ℹ No changes were made. To install, run without --dry-run:")
	fmt.Printf("  picoclaw-agents service install --port %d\n", p.cfg.Port)
	return nil
}

func (p *linuxSystemd) Install() error {
	if err := p.preFlightChecks(); err != nil {
		return err
	}

	// Backup existing unit
	if _, err := os.Stat(p.unitPath()); err == nil {
		backup := p.unitPath() + ".bak"
		if err := os.Rename(p.unitPath(), backup); err != nil {
			fmt.Printf("⚠ Warning: could not backup existing unit: %v\n", err)
		} else {
			fmt.Println("✓ Backed up existing unit to .bak")
		}
	}

	// Write unit atomically
	unitContent := p.buildSystemdUnit()
	if err := writeServiceFileAtomic(p.unitPath(), unitContent); err != nil {
		return fmt.Errorf("failed to write unit file: %w", err)
	}
	fmt.Println("✓ Created systemd unit:", p.unitPath())

	// Daemon reload
	if err := runCmd("systemctl", "--user", "daemon-reload"); err != nil {
		fmt.Printf("⚠ Warning: daemon-reload failed: %v\n", err)
	}
	fmt.Println("✓ Reloaded systemd daemon")

	// Enable
	if err := runCmd("systemctl", "--user", "enable", "picoclaw-agents.service"); err != nil {
		fmt.Printf("⚠ Warning: enable failed: %v\n", err)
	}
	fmt.Println("✓ Enabled service")

	// Start
	if err := runCmd("systemctl", "--user", "start", "picoclaw-agents.service"); err != nil {
		return fmt.Errorf("failed to start service: %w", err)
	}
	fmt.Println("✓ Started service")

	// Verify (poll up to 5 seconds)
	verified := false
	for i := 0; i < 10; i++ {
		out, err := runCmdSilent("systemctl", "--user", "is-active", "picoclaw-agents.service")
		if err == nil && strings.TrimSpace(out) == "active" {
			verified = true
			break
		}
		time.Sleep(500 * time.Millisecond)
	}

	if verified {
		fmt.Println("✓ Service is running (active)")
	} else {
		fmt.Println("⚠ Service started but verification timed out (may need a moment to initialize)")
		fmt.Println("  Check status with: picoclaw-agents service status")
	}

	fmt.Println()
	fmt.Println("🎉 PicoClaw-Agents service installed successfully!")
	fmt.Printf("   Port: %d\n", p.cfg.Port)
	fmt.Println()
	fmt.Println("Manage the service with:")
	fmt.Println("  picoclaw-agents service status   Check status")
	fmt.Println("  picoclaw-agents service logs -f  Follow logs")
	fmt.Println("  picoclaw-agents service stop     Stop service")
	fmt.Println("  picoclaw-agents service uninstall Remove service")
	return nil
}

func (p *linuxSystemd) Uninstall() error {
	// Stop first
	runCmdSilent("systemctl", "--user", "stop", "picoclaw-agents.service")

	// Disable
	runCmdSilent("systemctl", "--user", "disable", "picoclaw-agents.service")

	// Remove unit
	if _, err := os.Stat(p.unitPath()); err == nil {
		if err := os.Remove(p.unitPath()); err != nil {
			return fmt.Errorf("failed to remove unit file: %w", err)
		}
		fmt.Println("✓ Removed unit file")
	}

	// Daemon reload
	runCmdSilent("systemctl", "--user", "daemon-reload")
	fmt.Println("✓ Reloaded systemd daemon")

	// Ask about journal cleanup
	fmt.Print("Clear service journal logs? [y/N]: ")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		answer := strings.ToLower(scanner.Text())
		if answer == "y" || answer == "yes" {
			runCmdSilent("journalctl", "--user", "--unit=picoclaw-agents.service", "--rotate")
			runCmdSilent("journalctl", "--user", "--unit=picoclaw-agents.service", "--vacuum-time=1s")
			fmt.Println("✓ Cleared journal logs")
		}
	}

	fmt.Println("✓ Service uninstalled")
	return nil
}

func (p *linuxSystemd) Start() error {
	if isGatewayRunning(p.cfg.Port) {
		return fmt.Errorf("gateway is already running on port %d", p.cfg.Port)
	}
	if err := runCmd("systemctl", "--user", "start", "picoclaw-agents.service"); err != nil {
		return fmt.Errorf("failed to start service: %w", err)
	}
	fmt.Println("✓ Service started")
	return nil
}

func (p *linuxSystemd) Stop() error {
	if err := runCmd("systemctl", "--user", "stop", "picoclaw-agents.service"); err != nil {
		fmt.Printf("⚠ Warning: stop failed: %v\n", err)
	}
	fmt.Println("✓ Service stopped")
	return nil
}

func (p *linuxSystemd) Status() (ServiceStatus, error) {
	status := ServiceStatus{Port: p.cfg.Port}

	// Check if unit exists
	if _, err := os.Stat(p.unitPath()); err != nil {
		return status, nil // Not installed
	}
	status.Installed = true

	// Check if active
	out, err := runCmdSilent("systemctl", "--user", "is-active", "picoclaw-agents.service")
	if err == nil && strings.TrimSpace(out) == "active" {
		status.Active = true
	}

	// Get main PID
	out, err = runCmdSilent("systemctl", "--user", "show", "picoclaw-agents.service", "--property=MainPID")
	if err == nil && strings.HasPrefix(out, "MainPID=") {
		pidStr := strings.TrimPrefix(out, "MainPID=")
		if pid, err := strconv.Atoi(strings.TrimSpace(pidStr)); err == nil && pid > 0 {
			status.PID = pid
		}
	}

	return status, nil
}

func (p *linuxSystemd) Logs(lines int, follow bool) error {
	args := []string{"--user", "--unit=picoclaw-agents.service", "--no-pager", "-n", strconv.Itoa(lines)}
	if follow {
		args = append(args, "-f")
	}
	cmd := exec.Command("journalctl", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (p *linuxSystemd) restart() error {
	if err := p.Stop(); err != nil {
		return err
	}
	time.Sleep(1 * time.Second)
	return p.Start()
}
