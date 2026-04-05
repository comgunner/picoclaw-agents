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
	"sync"
	"time"
)

// TokenStatus represents the status of an OAuth token.
type TokenStatus struct {
	Provider    string    `json:"provider"`
	Expiry      time.Time `json:"expiry,omitempty"`
	Status      string    `json:"status"` // "valid" | "expiring_soon" | "expired" | "unknown"
	LastChecked time.Time `json:"last_checked"`
	Email       string    `json:"email,omitempty"`
}

// TokenMonitor monitors OAuth token expiration.
type TokenMonitor struct {
	configPath    string
	CheckInterval time.Duration // Public: check interval (default: 5 minutes)
	stopCh        chan struct{}
	mu            sync.RWMutex
	tokens        map[string]*TokenStatus
}

// NewTokenMonitor creates a new token monitor.
func NewTokenMonitor(configPath string) *TokenMonitor {
	return &TokenMonitor{
		configPath:    configPath,
		CheckInterval: 5 * time.Minute,
		stopCh:        make(chan struct{}),
		tokens:        make(map[string]*TokenStatus),
	}
}

// Start begins background monitoring of token expiration.
// The monitor checks tokens every checkInterval (default: 5 minutes).
func (m *TokenMonitor) Start(ctx context.Context) {
	go m.run(ctx)
}

// Stop halts the background monitoring goroutine.
func (m *TokenMonitor) Stop() {
	m.mu.Lock()
	defer m.mu.Unlock()

	select {
	case <-m.stopCh:
		// Already closed
	default:
		close(m.stopCh)
	}
}

// Status returns a copy of the current token status map.
func (m *TokenMonitor) Status() map[string]*TokenStatus {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Return a copy to prevent race conditions
	tokensCopy := make(map[string]*TokenStatus, len(m.tokens))
	for k, v := range m.tokens {
		tokensCopy[k] = v
	}
	return tokensCopy
}

// run is the main monitoring loop.
func (m *TokenMonitor) run(ctx context.Context) {
	ticker := time.NewTicker(m.CheckInterval)
	defer ticker.Stop()

	// Check immediately on start
	m.CheckTokens()

	for {
		select {
		case <-ctx.Done():
			return
		case <-m.stopCh:
			return
		case <-ticker.C:
			m.CheckTokens()
		}
	}
}

// CheckTokens loads and checks all stored tokens for expiration.
func (m *TokenMonitor) CheckTokens() {
	store, err := LoadStore()
	if err != nil {
		// If we can't load the store, just skip this check
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()

	for provider, cred := range store.Credentials {
		status := &TokenStatus{
			Provider:    provider,
			Expiry:      cred.ExpiresAt,
			LastChecked: now,
			Email:       cred.Email,
		}

		if cred.IsExpired() {
			status.Status = "expired"
		} else if cred.NeedsRefresh() {
			status.Status = "expiring_soon"
		} else {
			status.Status = "valid"
		}

		m.tokens[provider] = status
	}
}

// GetExpiringSoon returns a list of providers with tokens expiring within 1 hour.
func (m *TokenMonitor) GetExpiringSoon() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var expiring []string
	for provider, status := range m.tokens {
		if status.Status == "expiring_soon" || status.Status == "expired" {
			expiring = append(expiring, provider)
		}
	}
	return expiring
}
