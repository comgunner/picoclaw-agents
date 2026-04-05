// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)

package tools

import (
	"context"

	"github.com/comgunner/picoclaw/pkg/utils"
)

// BatchIDTool allows agents to generate unique readable identifiers for non-LLM tasks.
type BatchIDTool struct{}

func NewBatchIDTool() *BatchIDTool {
	return &BatchIDTool{}
}

func (t *BatchIDTool) Name() string {
	return "batch_id"
}

func (t *BatchIDTool) Description() string {
	return "Genera un ID único y legible (ej: #IMA_GEN_02_03_26_1500) para identificar tareas batch sin quema de tokens."
}

func (t *BatchIDTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"prefix": map[string]any{
				"type":        "string",
				"description": "Prefijo para el ID (ej: 'IMA_GEN', 'TEXT_GEN', 'SOCIAL')",
			},
		},
		"required": []string{"prefix"},
	}
}

func (t *BatchIDTool) Execute(ctx context.Context, args map[string]any) *ToolResult {
	prefix, _ := args["prefix"].(string)
	if prefix == "" {
		prefix = "BATCH"
	}

	id := utils.GenerateBatchID(prefix)
	return UserResult(id)
}
