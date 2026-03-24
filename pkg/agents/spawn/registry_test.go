// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package spawn

import (
	"sync"
	"testing"
	"time"
)

func TestRegistry_RegisterConcurrent(t *testing.T) {
	reg := NewRegistry(100) // Increase max concurrent for test
	const total = 50

	var wg sync.WaitGroup
	errCh := make(chan error, total)
	for i := 0; i < total; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			err := reg.Register(Run{
				RunID:            "run-" + string(rune(i+'0')),
				RequesterSession: "session-main",
				ChildSessionKey:  "agent:child:subagent:uuid",
				AgentID:          "child",
				ParentDepth:      0,
				StartedAt:        time.Now(),
			})
			if err != nil {
				errCh <- err
			}
		}(i)
	}
	wg.Wait()
	close(errCh)

	for err := range errCh {
		t.Fatalf("unexpected error: %v", err)
	}

	if count := reg.Count("session-main"); count != total {
		t.Fatalf("expected %d runs, got %d", total, count)
	}
}

func TestRegistry_CompleteRace(t *testing.T) {
	reg := NewRegistry(5)

	run := Run{
		RunID:            "run-test",
		RequesterSession: "session-main",
		ChildSessionKey:  "agent:child:subagent:uuid",
		AgentID:          "child",
		ParentDepth:      0,
		StartedAt:        time.Now(),
	}
	if err := reg.Register(run); err != nil {
		t.Fatalf("register: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		_ = reg.Complete("run-test")
	}()
	go func() {
		defer wg.Done()
		_ = reg.Complete("run-test")
	}()
	wg.Wait()

	// Should be completed without error (second call should return "not found")
}
