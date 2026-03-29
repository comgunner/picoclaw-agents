// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package agent

import (
	"fmt"

	"github.com/comgunner/picoclaw/pkg/config"
	"github.com/comgunner/picoclaw/pkg/providers"
)

// PruneMessages recorta tool results voluminosos en memoria.
// No modifica el historial persistido — solo lo que se envía al LLM.
//
// SPRINT 1 FEATURE: Tool result pruning before LLM call
func PruneMessages(messages []providers.Message, cfg config.ContextPruningConfig) []providers.Message {
	if !cfg.Enabled || cfg.MaxToolResultChars <= 0 {
		return messages
	}

	// Build exclusion and aggression sets for O(1) lookup
	excludeSet := make(map[string]bool, len(cfg.ExcludeTools))
	for _, t := range cfg.ExcludeTools {
		excludeSet[t] = true
	}

	aggressiveSet := make(map[string]bool, len(cfg.AggressiveTools))
	for _, t := range cfg.AggressiveTools {
		aggressiveSet[t] = true
	}

	result := make([]providers.Message, len(messages))
	copy(result, messages)

	for i, msg := range result {
		if msg.Role != "tool" {
			continue
		}

		// Determinar si esta tool está excluida
		toolName := extractToolName(messages, i)
		if excludeSet[toolName] {
			continue
		}

		// Determinar límite: agresivo = mitad del límite normal
		limit := cfg.MaxToolResultChars
		if aggressiveSet[toolName] {
			limit = limit / 2
		}

		// Truncar si excede el límite
		if len(msg.Content) > limit {
			// Preservar últimos 200 chars para contexto final
			tailLen := 200
			if len(msg.Content) < tailLen {
				tailLen = len(msg.Content)
			}

			// Calculator head length: reservar espacio para tail y marcador
			markerLen := 100 // espacio para "... [truncated X chars, tool: Y] ..."
			headLen := limit - tailLen - markerLen
			if headLen < 0 {
				headLen = 0
			}

			head := msg.Content[:headLen]
			tail := msg.Content[len(msg.Content)-tailLen:]

			result[i].Content = fmt.Sprintf(
				"%s\n\n... [truncated %d chars, tool: %s] ...\n\n%s",
				head,
				len(msg.Content)-limit,
				toolName,
				tail,
			)
		}
	}

	return result
}

// extractToolName busca el nombre de la tool en el mensaje assistant precedente
// que tiene tool_calls con el toolCallID correspondiente
func extractToolName(messages []providers.Message, toolMsgIdx int) string {
	if toolMsgIdx == 0 {
		return ""
	}

	toolCallID := messages[toolMsgIdx].ToolCallID
	if toolCallID == "" {
		return ""
	}

	// Buscar hacia atrás el assistant con tool_calls que tenga este ID
	for j := toolMsgIdx - 1; j >= 0; j-- {
		msg := messages[j]
		if msg.Role != "assistant" {
			continue
		}
		for _, tc := range msg.ToolCalls {
			if tc.ID == toolCallID {
				return tc.Function.Name
			}
		}
	}

	return ""
}
