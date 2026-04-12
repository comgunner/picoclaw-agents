// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

//go:build windows

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

type windowsSchtasks struct {
	cfg     ServiceConfig
	exePath string
}

// newPlatform returns the windows platform implementation.
func newPlatform(cfg ServiceConfig, exePath string) (ServicePlatform, error) {
	return &windowsSchtasks{cfg: cfg, exePath: exePath}, nil
}

func (p *windowsSchtasks) taskName() string {
	return "PicoClaw-Agents Gateway"
}

func (p *windowsSchtasks) appDataDir() string {
	appData := os.Getenv("APPDATA")
	if appData == "" {
		appData = filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Roaming")
	}
	return filepath.Join(appData, "PicoClaw")
}

func (p *windowsSchtasks) wrapperPath() string {
	return filepath.Join(p.appDataDir(), "gateway.cmd")
}

func (p *windowsSchtasks) pidPath() string {
	return filepath.Join(p.appDataDir(), "gateway.pid")
}

func (p *windowsSchtasks) logDir() string {
	return filepath.Join(p.appDataDir(), "logs")
}

func (p *windowsSchtasks) buildWrapperScript() string {
	publicFlag := ""
	if p.cfg.Public {
		publicFlag = " --public"
	}

	return `@echo off
setlocal EnableDelayedExpansion
set "PICOCLOW_PID=` + p.pidPath() + `"
set "PICOCLOW_LOG=` + filepath.Join(p.logDir(), "gateway.log") + `"

REM Check if already running
if exist "!PICOCLOW_PID!" (
    set /p EXISTING_PID=<"!PICOCLOW_PID!"
    tasklist /FI "PID eq !EXISTING_PID!" 2>NUL | find /I "picoclaw" >NUL
    if not ERRORLEVEL 1 (
        echo Gateway is already running (PID: !EXISTING_PID!)
        exit /b 1
    )
)

REM Start gateway and save PID
start "PicoClaw Gateway" /B "` + p.exePath + `" gateway --port ` + strconv.Itoa(p.cfg.Port) + publicFlag + ` >> "!PICOCLOW_LOG!" 2>&1
timeout /t 2 /nobreak >NUL

REM Get the PID of the most recently started picoclaw process
for /f "tokens=2 delims=," %%a in ('tasklist /FI "IMAGENAME eq picoclaw-agents*" /FO CSV /NH 2^>NUL') do (
    echo %%~a > "!PICOCLOW_PID!"
    goto :done
)
:done
`
}

func (p *windowsSchtasks) preFlightChecks() error {
	if isGatewayRunning(p.cfg.Port) {
		return fmt.Errorf(
			"gateway is already running on port %d. Stop it first with 'picoclaw-agents gateway stop' or 'taskkill /PID <PID>' to avoid duplicate connections (Telegram, Discord, etc.)",
			p.cfg.Port,
		)
	}

	// Check if task already exists
	out, _ := runCmdSilent("schtasks", "/Query", "/TN", p.taskName(), "/FO", "LIST")
	if strings.Contains(out, "TaskName:") {
		return fmt.Errorf(
			"service is already installed as scheduled task '%s'. Use 'picoclaw-agents service uninstall' first",
			p.taskName(),
		)
	}

	if !isPortAvailable(p.cfg.Port) {
		return fmt.Errorf("port %d is already in use by another process", p.cfg.Port)
	}
	return nil
}

func (p *windowsSchtasks) DryRunInstall() error {
	if err := p.preFlightChecks(); err != nil {
		return err
	}

	fmt.Printf("Platform:       Windows (schtasks)\n")
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
	fmt.Printf("  1. Create directory: %s\n", p.appDataDir())
	fmt.Printf("  2. Create log directory: %s\n", p.logDir())
	fmt.Printf("  3. Write wrapper script: %s\n", p.wrapperPath())
	fmt.Println("  4. Register scheduled task:", p.taskName())
	fmt.Println("  5. Start task immediately")
	fmt.Println()
	fmt.Println("Wrapper script content (would be written):")
	fmt.Println("──────────────────────────────────────")
	fmt.Print(p.buildWrapperScript())
	fmt.Println("──────────────────────────────────────")
	fmt.Println()
	fmt.Println("ℹ No changes were made. To install, run without --dry-run:")
	fmt.Printf("  picoclaw-agents service install --port %d\n", p.cfg.Port)
	return nil
}

func (p *windowsSchtasks) Install() error {
	if err := p.preFlightChecks(); err != nil {
		return err
	}

	// Create directories
	if err := os.MkdirAll(p.appDataDir(), 0o755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	if err := os.MkdirAll(p.logDir(), 0o755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	// Write wrapper script
	wrapperContent := p.buildWrapperScript()
	if err := writeServiceFileAtomic(p.wrapperPath(), wrapperContent); err != nil {
		return fmt.Errorf("failed to write wrapper script: %w", err)
	}
	fmt.Println("✓ Created wrapper script:", p.wrapperPath())

	// Register scheduled task (ONLOGON with LIMITED rights)
	createCmd := exec.Command("schtasks", "/Create", "/F",
		"/TN", p.taskName(),
		"/TR", p.wrapperPath(),
		"/SC", "ONLOGON",
		"/RL", "LIMITED",
	)
	createCmd.Stdout = os.Stdout
	createCmd.Stderr = os.Stderr
	if err := createCmd.Run(); err != nil {
		return fmt.Errorf("failed to create scheduled task: %w", err)
	}
	fmt.Println("✓ Created scheduled task:", p.taskName())

	// Start immediately
	runCmdSilent("schtasks", "/Run", "/TN", p.taskName())
	fmt.Println("✓ Started service")

	// Verify (poll PID file)
	verified := false
	for i := 0; i < 5; i++ {
		if _, err := os.Stat(p.pidPath()); err == nil {
			verified = true
			break
		}
		time.Sleep(1 * time.Second)
	}

	if verified {
		fmt.Println("✓ Service is running")
	} else {
		fmt.Println("⚠ Service started but PID verification pending (may need a moment to initialize)")
	}

	fmt.Println()
	fmt.Println("🎉 PicoClaw-Agents service installed successfully!")
	fmt.Printf("   Port: %d\n", p.cfg.Port)
	fmt.Printf("   Logs: %s\n", filepath.Join(p.logDir(), "gateway.log"))
	fmt.Println()
	fmt.Println("Manage the service with:")
	fmt.Println("  picoclaw-agents service status   Check status")
	fmt.Println("  picoclaw-agents service logs     View logs")
	fmt.Println("  picoclaw-agents service stop     Stop service")
	fmt.Println("  picoclaw-agents service uninstall Remove service")
	return nil
}

func (p *windowsSchtasks) Uninstall() error {
	// Stop first
	p.Stop()

	// Delete task
	runCmdSilent("schtasks", "/Delete", "/TN", p.taskName(), "/F")
	fmt.Println("✓ Deleted scheduled task")

	// Remove wrapper script
	if _, err := os.Stat(p.wrapperPath()); err == nil {
		os.Remove(p.wrapperPath())
		fmt.Println("✓ Removed wrapper script")
	}

	// Remove PID file
	os.Remove(p.pidPath())

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

func (p *windowsSchtasks) Start() error {
	if isGatewayRunning(p.cfg.Port) {
		return fmt.Errorf("gateway is already running on port %d", p.cfg.Port)
	}
	if err := runCmd("schtasks", "/Run", "/TN", p.taskName()); err != nil {
		return fmt.Errorf("failed to start task: %w", err)
	}
	fmt.Println("✓ Service started")
	return nil
}

func (p *windowsSchtasks) Stop() error {
	runCmdSilent("schtasks", "/End", "/TN", p.taskName())

	// Also kill process from PID file
	if pidStr, err := os.ReadFile(p.pidPath()); err == nil {
		pid, err := strconv.Atoi(strings.TrimSpace(string(pidStr)))
		if err == nil && pid > 0 {
			runCmdSilent("taskkill", "/PID", strconv.Itoa(pid), "/F")
		}
	}
	fmt.Println("✓ Service stopped")
	return nil
}

func (p *windowsSchtasks) Status() (ServiceStatus, error) {
	status := ServiceStatus{Port: p.cfg.Port}

	// Check if task exists
	out, err := runCmdSilent("schtasks", "/Query", "/TN", p.taskName(), "/FO", "LIST")
	if err != nil || !strings.Contains(out, "TaskName:") {
		return status, nil // Not installed
	}
	status.Installed = true

	// Check if running
	if strings.Contains(out, "Running") {
		status.Active = true
	}

	// Get PID
	if pidStr, err := os.ReadFile(p.pidPath()); err == nil {
		if pid, err := strconv.Atoi(strings.TrimSpace(string(pidStr))); err == nil && pid > 0 {
			status.PID = pid
		}
	}

	return status, nil
}

func (p *windowsSchtasks) Logs(lines int, follow bool) error {
	logFile := filepath.Join(p.logDir(), "gateway.log")
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		return fmt.Errorf("no logs found at %s", logFile)
	}

	if follow {
		cmd := exec.Command("powershell", "-Command",
			fmt.Sprintf("Get-Content -Path '%s' -Wait -Tail %d", logFile, lines))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	cmd := exec.Command("powershell", "-Command",
		fmt.Sprintf("Get-Content -Path '%s' -Tail %d", logFile, lines))
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func (p *windowsSchtasks) restart() error {
	if err := p.Stop(); err != nil {
		return err
	}
	time.Sleep(2 * time.Second)
	return p.Start()
}
