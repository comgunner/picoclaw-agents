// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package auth

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"github.com/comgunner/picoclaw/pkg/auth"
)

func newMonitorCommand() *cobra.Command {
	var interval int

	cmd := &cobra.Command{
		Use:   "monitor",
		Short: "Continuously monitor OAuth token expiration",
		Long:  "Watch token expiration status in real-time. Press Ctrl+C to stop.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runMonitor(interval)
		},
	}

	cmd.Flags().IntVarP(&interval, "interval", "i", 5, "Check interval in minutes")
	return cmd
}

func runMonitor(intervalMinutes int) error {
	home, _ := os.UserHomeDir()
	configPath := filepath.Join(home, ".picoclaw", "config.json")

	monitor := auth.NewTokenMonitor(configPath)

	// Set custom interval
	monitor.CheckInterval = time.Duration(intervalMinutes) * time.Minute

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle Ctrl+C
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		fmt.Println("\n\nStopping token monitor...")
		cancel()
		monitor.Stop()
	}()

	fmt.Println("Starting token expiration monitor...")
	fmt.Printf("Checking every %d minutes. Press Ctrl+C to stop.\n\n", intervalMinutes)

	// Start monitoring
	monitor.Start(ctx)

	// Display loop
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			displayStatus(monitor)
		}
	}
}

func displayStatus(monitor *auth.TokenMonitor) {
	status := monitor.Status()
	now := time.Now().Format("2006-01-02 15:04:05")

	fmt.Printf("\n[%s] Token Status:\n", now)

	if len(status) == 0 {
		fmt.Println("  No tokens found.")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "  Provider\tStatus\tExpiry")
	fmt.Fprintln(w, "  ────────────────────────────────────────")

	hasExpiring := false
	for provider, s := range status {
		expiry := "N/A"
		if !s.Expiry.IsZero() {
			expiry = s.Expiry.Format("2006-01-02 15:04")
		}

		statusIcon := "✓"
		if s.Status == "expiring_soon" {
			statusIcon = "⚠️"
			hasExpiring = true
		} else if s.Status == "expired" {
			statusIcon = "✗"
			hasExpiring = true
		}

		fmt.Fprintf(w, "  %s\t%s %s\t%s\n", provider, statusIcon, s.Status, expiry)
	}
	w.Flush()

	if hasExpiring {
		fmt.Println("\n⚠️  Some tokens are expiring or expired. Run 'picoclaw auth login <provider>' to refresh.")
	}
}
