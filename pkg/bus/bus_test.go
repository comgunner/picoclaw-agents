// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package bus

import (
	"context"
	"testing"
	"time"
)

func TestMessageBus_Close(t *testing.T) {
	mb := NewMessageBus()
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// 1. Initial state
	mb.PublishInbound(InboundMessage{Content: "hello"})
	msg, ok := mb.ConsumeInbound(ctx)
	if !ok || msg.Content != "hello" {
		t.Fatalf("Expected hello, got %v (ok=%v)", msg.Content, ok)
	}

	// 2. Close bus
	mb.Close()

	// 3. Verify ConsumeInbound returns ok=false
	msg, ok = mb.ConsumeInbound(ctx)
	if ok {
		t.Error("Expected ok=false from closed bus, got true")
	}

	// 4. Verify SubscribeOutbound returns ok=false
	_, ok = mb.SubscribeOutbound(ctx)
	if ok {
		t.Error("Expected ok=false for outbound from closed bus, got true")
	}
}
