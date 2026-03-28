// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package channels

import "net/http"

// WebhookHandler is implemented by channels that expose HTTP endpoints (e.g., WebSocket).
// The gateway mounts each WebhookHandler onto the shared health-server mux so that all
// HTTP traffic (health checks, WebSocket upgrades, webhooks) is served from a single port.
type WebhookHandler interface {
	// WebhookPath returns the URL path prefix this channel handles (e.g., "/pico/").
	WebhookPath() string
	// ServeHTTP handles incoming HTTP requests routed to WebhookPath.
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}
