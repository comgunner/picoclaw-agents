// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package constants

// DefaultMaxContextTokens defines the typical maximum number of tokens most
// large language models currently support.  It is used for pre-flight checks
// and to drive user-facing alerts.  Deployments may override this by setting
// a configuration value later on.
const DefaultMaxContextTokens = 131072
