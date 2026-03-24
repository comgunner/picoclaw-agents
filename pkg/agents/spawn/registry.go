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
	"fmt"
	"sync"
	"time"
)

type Run struct {
	RunID             string
	RequesterSession  string
	ChildSessionKey   string
	AgentID           string
	ParentDepth       int
	StartedAt         time.Time
	ExpectedCompleted bool
}

type Registry struct {
	mu             sync.RWMutex
	runsByReq      map[string][]Run
	runToRequester map[string]string
	maxConcurrent  int
}

func NewRegistry(maxConcurrent int) *Registry {
	if maxConcurrent <= 0 {
		maxConcurrent = 1
	}
	return &Registry{
		runsByReq:      make(map[string][]Run),
		runToRequester: make(map[string]string),
		maxConcurrent:  maxConcurrent,
	}
}

func (r *Registry) Register(run Run) error {
	if run.RunID == "" || run.RequesterSession == "" {
		return fmt.Errorf("invalid run metadata")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.runToRequester[run.RunID]; exists {
		return fmt.Errorf("run already registered: %s", run.RunID)
	}

	active := 0
	for _, runs := range r.runsByReq {
		active += len(runs)
	}
	if active >= r.maxConcurrent {
		return fmt.Errorf("max concurrent subagents reached (%d)", r.maxConcurrent)
	}

	r.runsByReq[run.RequesterSession] = append(r.runsByReq[run.RequesterSession], run)
	r.runToRequester[run.RunID] = run.RequesterSession
	return nil
}

func (r *Registry) Count(requesterSession string) int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.runsByReq[requesterSession])
}

func (r *Registry) Get(requesterSession string) []Run {
	r.mu.RLock()
	defer r.mu.RUnlock()
	runs := r.runsByReq[requesterSession]
	out := make([]Run, len(runs))
	copy(out, runs)
	return out
}

func (r *Registry) Complete(runID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	requester, ok := r.runToRequester[runID]
	if !ok {
		return fmt.Errorf("run not found: %s", runID)
	}

	runs := r.runsByReq[requester]
	idx := -1
	for i := range runs {
		if runs[i].RunID == runID {
			idx = i
			break
		}
	}
	if idx < 0 {
		return fmt.Errorf("run index missing for %s", runID)
	}

	r.runsByReq[requester] = append(runs[:idx], runs[idx+1:]...)
	delete(r.runToRequester, runID)
	if len(r.runsByReq[requester]) == 0 {
		delete(r.runsByReq, requester)
	}
	return nil
}

func (r *Registry) TotalActive() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	total := 0
	for _, runs := range r.runsByReq {
		total += len(runs)
	}
	return total
}
