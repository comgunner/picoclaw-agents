// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package agents

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestManager_RegisterConcurrent_NoRaces(t *testing.T) {
	m := NewManager()
	const total = 100

	var wg sync.WaitGroup
	errCh := make(chan error, total)
	for i := 0; i < total; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_, err := m.RegisterAgent(fmt.Sprintf("agent-%d", i))
			if err != nil {
				errCh <- err
			}
		}(i)
	}
	wg.Wait()
	close(errCh)

	for err := range errCh {
		t.Fatalf("unexpected register error: %v", err)
	}

	items := m.ListAgents()
	if len(items) != total {
		t.Fatalf("expected %d agents, got %d", total, len(items))
	}
}

func TestManager_StartStopAgent(t *testing.T) {
	m := NewManager()
	if _, err := m.RegisterAgent("main"); err != nil {
		t.Fatalf("register: %v", err)
	}

	done := make(chan struct{})
	err := m.StartAgent(context.Background(), "main", func(ctx context.Context, _ string) {
		defer close(done)
		<-ctx.Done()
	})
	if err != nil {
		t.Fatalf("start: %v", err)
	}

	if stopErr := m.StopAgent("main"); stopErr != nil {
		t.Fatalf("stop: %v", stopErr)
	}

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("runner did not stop in time")
	}
}
