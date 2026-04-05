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
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/comgunner/picoclaw/pkg/auth"
)

func newTokensCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tokens",
		Short: "List stored OAuth tokens and their status",
		Long:  "Display all stored authentication tokens with their expiration status.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runTokens()
		},
	}

	return cmd
}

func runTokens() error {
	store, err := auth.LoadStore()
	if err != nil {
		return fmt.Errorf("loading auth store: %w", err)
	}

	if len(store.Credentials) == 0 {
		fmt.Println("No authentication tokens stored.")
		fmt.Println("\nUse 'picoclaw auth login <provider>' to authenticate.")
		return nil
	}

	// Create monitor to get status
	home, _ := os.UserHomeDir()
	configPath := filepath.Join(home, ".picoclaw", "config.json")
	monitor := auth.NewTokenMonitor(configPath)
	monitor.CheckTokens()

	status := monitor.Status()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Provider\tStatus\tExpiry\tEmail")
	fmt.Fprintln(w, "─────────────────────────────────────────────────────────")

	for provider, cred := range store.Credentials {
		tokenStatus := "unknown"
		expiry := "N/A"
		email := cred.Email

		if s, ok := status[provider]; ok {
			tokenStatus = s.Status
			if !s.Expiry.IsZero() {
				expiry = s.Expiry.Format("2006-01-02 15:04")
			}
		}

		if cred.IsExpired() {
			tokenStatus = "expired"
		} else if cred.NeedsRefresh() {
			tokenStatus = "expiring_soon"
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", provider, tokenStatus, expiry, email)
	}

	w.Flush()
	return nil
}
