// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package config

// Build-time variables injected via -ldflags.
// Targets: make build-launcher, make build-launcher-tui
var (
	Version   = "dev"
	GitCommit = "dev"
	BuildTime = ""
	GoVersion = ""
)
