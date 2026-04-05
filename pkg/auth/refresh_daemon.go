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
	"log"
	"sync"
	"time"
)

var daemonOnce sync.Once

// StartRefreshDaemon launches a background goroutine that proactively refreshes
// the OAuth token before it expires, preventing expiration during idle periods.
//
// Parameters:
//   - provider: credential key in auth.json (e.g. "google-antigravity")
//   - cfg: OAuth provider config with scopes, client ID, etc.
//   - checkInterval: how often to check the token (recommended: 20 min)
//   - refreshThreshold: refresh if less than this duration remains (recommended: 30 min)
//
// The daemon runs exactly once per process (guarded by sync.Once).
// It silently handles errors — the next real request will also attempt refresh.
func StartRefreshDaemon(provider string, cfg OAuthProviderConfig, checkInterval, refreshThreshold time.Duration) {
	daemonOnce.Do(func() {
		go func() {
			ticker := time.NewTicker(checkInterval)
			defer ticker.Stop()

			for range ticker.C {
				cred, err := GetCredential(provider)
				if err != nil || cred == nil || cred.RefreshToken == "" {
					continue
				}

				// Skip if token has plenty of time left
				if !cred.ExpiresAt.IsZero() && time.Until(cred.ExpiresAt) > refreshThreshold {
					continue
				}

				// Token is about to expire or already expired — refresh it
				refreshed, err := RefreshAccessToken(cred, cfg)
				if err != nil {
					log.Printf("[WARN] refresh_daemon: proactive token refresh failed for %s: %v", provider, err)
					continue
				}

				// Preserve fields that Google doesn't return in refresh responses
				refreshed.Email = cred.Email
				if refreshed.ProjectID == "" {
					refreshed.ProjectID = cred.ProjectID
				}

				if err := SetCredential(provider, refreshed); err != nil {
					log.Printf("[WARN] refresh_daemon: failed to save refreshed token for %s: %v", provider, err)
					continue
				}

				log.Printf("[INFO] refresh_daemon: proactively refreshed token for %s (new expiry: %s)",
					provider, refreshed.ExpiresAt.Format(time.RFC3339))
			}
		}()
	})
}
