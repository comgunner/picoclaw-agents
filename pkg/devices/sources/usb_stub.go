// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

//go:build !linux

package sources

import (
	"context"

	"github.com/comgunner/picoclaw/pkg/devices/events"
)

type USBMonitor struct{}

func NewUSBMonitor() *USBMonitor {
	return &USBMonitor{}
}

func (m *USBMonitor) Kind() events.Kind {
	return events.KindUSB
}

func (m *USBMonitor) Start(ctx context.Context) (<-chan *events.DeviceEvent, error) {
	ch := make(chan *events.DeviceEvent)
	close(ch) // Immediately close, no events
	return ch, nil
}

func (m *USBMonitor) Stop() error {
	return nil
}
