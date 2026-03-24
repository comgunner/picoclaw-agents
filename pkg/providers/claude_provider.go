// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package providers

import (
	"context"
	"fmt"

	anthropicprovider "github.com/comgunner/picoclaw/pkg/providers/anthropic"
)

type ClaudeProvider struct {
	delegate *anthropicprovider.Provider
}

func NewClaudeProvider(token string) *ClaudeProvider {
	return &ClaudeProvider{
		delegate: anthropicprovider.NewProvider(token),
	}
}

func NewClaudeProviderWithBaseURL(token, apiBase string) *ClaudeProvider {
	return &ClaudeProvider{
		delegate: anthropicprovider.NewProviderWithBaseURL(token, apiBase),
	}
}

func NewClaudeProviderWithTokenSource(token string, tokenSource func() (string, error)) *ClaudeProvider {
	return &ClaudeProvider{
		delegate: anthropicprovider.NewProviderWithTokenSource(token, tokenSource),
	}
}

func NewClaudeProviderWithTokenSourceAndBaseURL(
	token string, tokenSource func() (string, error), apiBase string,
) *ClaudeProvider {
	return &ClaudeProvider{
		delegate: anthropicprovider.NewProviderWithTokenSourceAndBaseURL(token, tokenSource, apiBase),
	}
}

func newClaudeProviderWithDelegate(delegate *anthropicprovider.Provider) *ClaudeProvider {
	return &ClaudeProvider{delegate: delegate}
}

func (p *ClaudeProvider) Chat(
	ctx context.Context, messages []Message, tools []ToolDefinition, model string, options map[string]any,
) (*LLMResponse, error) {
	resp, err := p.delegate.Chat(ctx, messages, tools, model, options)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (p *ClaudeProvider) GetDefaultModel() string {
	return p.delegate.GetDefaultModel()
}

func createClaudeTokenSource() func() (string, error) {
	return func() (string, error) {
		cred, err := getCredential("anthropic")
		if err != nil {
			return "", fmt.Errorf("loading auth credentials: %w", err)
		}
		if cred == nil {
			return "", fmt.Errorf("no credentials for anthropic. Run: picoclaw auth login --provider anthropic")
		}
		return cred.AccessToken, nil
	}
}
