// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package gateway

import (
	"context"
	"fmt"
	"sync"
)

type Call struct {
	Method    string
	Params    map[string]any
	TimeoutMs int
}

type Response struct {
	RunID  string
	Status string
	Data   map[string]any
}

type Client interface {
	Call(ctx context.Context, request Call) (Response, error)
}

type MockClient struct {
	mu        sync.RWMutex
	handlers  map[string]func(ctx context.Context, request Call) (Response, error)
	callCount map[string]int
}

func NewMockClient() *MockClient {
	return &MockClient{
		handlers:  make(map[string]func(ctx context.Context, request Call) (Response, error)),
		callCount: make(map[string]int),
	}
}

func (m *MockClient) SetHandler(method string, h func(ctx context.Context, request Call) (Response, error)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.handlers[method] = h
}

func (m *MockClient) Call(ctx context.Context, request Call) (Response, error) {
	if err := ctx.Err(); err != nil {
		return Response{}, err
	}
	m.mu.Lock()
	m.callCount[request.Method]++
	h := m.handlers[request.Method]
	m.mu.Unlock()

	if h == nil {
		return Response{}, fmt.Errorf("no handler for method: %s", request.Method)
	}
	return h(ctx, request)
}

func (m *MockClient) Count(method string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.callCount[method]
}
