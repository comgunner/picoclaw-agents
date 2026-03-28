// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package auth

// DeviceCodeInfo holds the information returned when initiating a device code flow.
type DeviceCodeInfo struct {
	DeviceAuthID string
	UserCode     string
	VerifyURL    string
	Interval     int
}

// NOTE: GenerateState, RequestDeviceCode, PollDeviceCodeOnce, y ExchangeCodeForTokens
// ahora son funciones públicas en oauth.go - los wrappers aquí fueron eliminados para evitar duplicación
