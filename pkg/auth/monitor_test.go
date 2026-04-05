// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package auth_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/comgunner/picoclaw/pkg/auth"
)

// TestTokenMonitor_Start_Stop verifies the monitor can start and stop cleanly.
func TestTokenMonitor_Start_Stop(t *testing.T) {
	monitor := auth.NewTokenMonitor("")
	ctx := context.Background()

	monitor.Start(ctx)

	// Give it a moment to initialize
	time.Sleep(10 * time.Millisecond)

	// Should not panic
	monitor.Stop()

	// Stopping twice should also not panic
	monitor.Stop()

	// Verify default interval
	assert.Equal(t, 5*time.Minute, monitor.CheckInterval)
}

// TestTokenMonitor_ExpiringStatus verifies token status detection.
func TestTokenMonitor_ExpiringStatus(t *testing.T) {
	// Create a store with test credentials
	store := &auth.AuthStore{
		Credentials: make(map[string]*auth.AuthCredential),
	}

	// Add a valid token (expires in 24 hours)
	store.Credentials["valid_provider"] = &auth.AuthCredential{
		Provider:  "valid_provider",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	// Add an expiring soon token (expires in 3 minutes - within 5 minute threshold)
	store.Credentials["expiring_provider"] = &auth.AuthCredential{
		Provider:  "expiring_provider",
		ExpiresAt: time.Now().Add(3 * time.Minute),
	}

	// Add an expired token
	store.Credentials["expired_provider"] = &auth.AuthCredential{
		Provider:  "expired_provider",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}

	// Save the store
	err := auth.SaveStore(store)
	assert.NoError(t, err)
	defer func() {
		// Cleanup: save empty store
		empty := &auth.AuthStore{Credentials: make(map[string]*auth.AuthCredential)}
		auth.SaveStore(empty)
	}()

	// Create monitor and check tokens
	monitor := auth.NewTokenMonitor("")
	monitor.CheckTokens()

	status := monitor.Status()

	// Verify valid token
	if valid, ok := status["valid_provider"]; assert.True(t, ok) {
		assert.Equal(t, "valid", valid.Status)
	}

	// Verify expiring soon token
	if expiring, ok := status["expiring_provider"]; assert.True(t, ok) {
		assert.Equal(t, "expiring_soon", expiring.Status)
	}

	// Verify expired token
	if expired, ok := status["expired_provider"]; assert.True(t, ok) {
		assert.Equal(t, "expired", expired.Status)
	}
}

// TestTokenMonitor_GetExpiringSoon verifies filtering of expiring tokens.
func TestTokenMonitor_GetExpiringSoon(t *testing.T) {
	// Create a store with test credentials
	store := &auth.AuthStore{
		Credentials: make(map[string]*auth.AuthCredential),
	}

	// Add tokens with different statuses
	store.Credentials["valid"] = &auth.AuthCredential{
		Provider:  "valid",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	// Token expiring in 3 minutes (within 5 minute threshold)
	store.Credentials["expiring"] = &auth.AuthCredential{
		Provider:  "expiring",
		ExpiresAt: time.Now().Add(3 * time.Minute),
	}
	store.Credentials["expired"] = &auth.AuthCredential{
		Provider:  "expired",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}

	err := auth.SaveStore(store)
	assert.NoError(t, err)
	defer func() {
		empty := &auth.AuthStore{Credentials: make(map[string]*auth.AuthCredential)}
		auth.SaveStore(empty)
	}()

	monitor := auth.NewTokenMonitor("")
	monitor.CheckTokens()

	expiring := monitor.GetExpiringSoon()
	assert.Len(t, expiring, 2)
	assert.Contains(t, expiring, "expiring")
	assert.Contains(t, expiring, "expired")
}

// TestTokenMonitor_Status_ReturnsCopy verifies Status returns a copy of the map.
func TestTokenMonitor_Status_ReturnsCopy(t *testing.T) {
	// Create a store with test credentials
	store := &auth.AuthStore{
		Credentials: make(map[string]*auth.AuthCredential),
	}
	store.Credentials["test"] = &auth.AuthCredential{
		Provider:  "test",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	err := auth.SaveStore(store)
	assert.NoError(t, err)
	defer func() {
		empty := &auth.AuthStore{Credentials: make(map[string]*auth.AuthCredential)}
		auth.SaveStore(empty)
	}()

	monitor := auth.NewTokenMonitor("")
	monitor.CheckTokens()

	status1 := monitor.Status()
	status2 := monitor.Status()

	// Verify both maps have the same keys but are different instances
	_, ok1 := status1["test"]
	_, ok2 := status2["test"]
	assert.True(t, ok1, "status1 should have test key")
	assert.True(t, ok2, "status2 should have test key")

	// Verify they're different map instances by checking length after deletion
	delete(status1, "test")
	_, stillExists := status2["test"]
	assert.True(t, stillExists, "status2 should still have test key after deleting from status1")
}

// TestTokenMonitor_EmptyStore verifies handling of empty store.
func TestTokenMonitor_EmptyStore(t *testing.T) {
	// Ensure store is empty
	empty := &auth.AuthStore{Credentials: make(map[string]*auth.AuthCredential)}
	auth.SaveStore(empty)

	monitor := auth.NewTokenMonitor("")
	monitor.CheckTokens()

	status := monitor.Status()
	assert.Empty(t, status, "Should have no tokens when store is empty")
}

// TestTokenMonitor_NewTokenMonitor verifies constructor.
func TestTokenMonitor_NewTokenMonitor(t *testing.T) {
	monitor := auth.NewTokenMonitor("/custom/path")

	assert.NotNil(t, monitor)
	// Verify the monitor is properly initialized
	status := monitor.Status()
	assert.NotNil(t, status)
	assert.Empty(t, status) // Should be empty initially
}
