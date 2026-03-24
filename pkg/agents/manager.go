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
	"sort"
	"strings"
	"sync"
)

type AgentRuntime struct {
	ID             string
	MainSessionKey string
	Running        bool
}

type AgentRunner func(ctx context.Context, agentID string)

type CompletionEvent struct {
	RunID           string
	ChildSessionKey string
	ResultSummary   string
}

type CompletionHook func(ctx context.Context, event CompletionEvent) error

// Manager keeps in-memory lifecycle state for multiple agent instances.
// It is concurrency-safe and intended as the foundation for multi-agent orchestration.
type Manager struct {
	mu         sync.RWMutex
	agents     map[string]*AgentRuntime
	cancellers map[string]context.CancelFunc
	hooks      []CompletionHook
}

func NewManager() *Manager {
	return &Manager{
		agents:     make(map[string]*AgentRuntime),
		cancellers: make(map[string]context.CancelFunc),
		hooks:      make([]CompletionHook, 0),
	}
}

func (m *Manager) RegisterAgent(agentID string) (string, error) {
	id := strings.TrimSpace(strings.ToLower(agentID))
	if id == "" {
		return "", fmt.Errorf("agent id is empty")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.agents[id]; exists {
		return "", fmt.Errorf("agent already registered: %s", id)
	}

	mainSession := fmt.Sprintf("agent:%s:main", id)
	m.agents[id] = &AgentRuntime{
		ID:             id,
		MainSessionKey: mainSession,
		Running:        false,
	}

	return mainSession, nil
}

func (m *Manager) StartAgent(parentCtx context.Context, agentID string, runner AgentRunner) error {
	id := strings.TrimSpace(strings.ToLower(agentID))
	if id == "" {
		return fmt.Errorf("agent id is empty")
	}
	if runner == nil {
		return fmt.Errorf("runner is nil")
	}

	m.mu.Lock()
	agent, ok := m.agents[id]
	if !ok {
		m.mu.Unlock()
		return fmt.Errorf("agent not found: %s", id)
	}
	if agent.Running {
		m.mu.Unlock()
		return fmt.Errorf("agent already running: %s", id)
	}

	ctx, cancel := context.WithCancel(parentCtx)
	m.cancellers[id] = cancel
	agent.Running = true
	m.mu.Unlock()

	go func() {
		defer func() {
			m.mu.Lock()
			if registered, exists := m.agents[id]; exists {
				registered.Running = false
			}
			delete(m.cancellers, id)
			m.mu.Unlock()
		}()
		runner(ctx, id)
	}()

	return nil
}

func (m *Manager) StopAgent(agentID string) error {
	id := strings.TrimSpace(strings.ToLower(agentID))
	if id == "" {
		return fmt.Errorf("agent id is empty")
	}

	m.mu.RLock()
	cancel, ok := m.cancellers[id]
	m.mu.RUnlock()
	if !ok {
		return fmt.Errorf("agent not running: %s", id)
	}

	cancel()
	return nil
}

func (m *Manager) GetAgent(agentID string) (AgentRuntime, bool) {
	id := strings.TrimSpace(strings.ToLower(agentID))
	m.mu.RLock()
	defer m.mu.RUnlock()
	if runtime, ok := m.agents[id]; ok {
		return *runtime, true
	}
	return AgentRuntime{}, false
}

func (m *Manager) ListAgents() []AgentRuntime {
	m.mu.RLock()
	defer m.mu.RUnlock()
	items := make([]AgentRuntime, 0, len(m.agents))
	for _, runtime := range m.agents {
		items = append(items, *runtime)
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].ID < items[j].ID
	})
	return items
}

func (m *Manager) RegisterCompletionHook(h CompletionHook) {
	if h == nil {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.hooks = append(m.hooks, h)
}

func (m *Manager) EmitCompletion(ctx context.Context, event CompletionEvent) error {
	m.mu.RLock()
	hooks := make([]CompletionHook, len(m.hooks))
	copy(hooks, m.hooks)
	m.mu.RUnlock()

	for _, hook := range hooks {
		if err := hook(ctx, event); err != nil {
			return err
		}
	}
	return nil
}
