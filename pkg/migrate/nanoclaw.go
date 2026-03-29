// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package migrate

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ResolveNanoClawHome finds nanoclaw home directory
func ResolveNanoClawHome(override string) (string, error) {
	if override != "" {
		return override, nil
	}
	if envHome := os.Getenv("NANOCLAW_HOME"); envHome != "" {
		return expandHome(envHome), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	// nanoclaw may use ~/.nanoclaw or ~/.config/nanoclaw
	candidates := []string{
		filepath.Join(home, ".nanoclaw"),
		filepath.Join(home, ".config", "nanoclaw"),
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c, nil
		}
	}
	return "", fmt.Errorf("nanoclaw home not found (tried: %v)", candidates)
}

// NanoClawRawConfig represents the raw nanoclaw config structure
type NanoClawRawConfig struct {
	Providers []json.RawMessage `json:"providers,omitempty"`
	Agents    []json.RawMessage `json:"agents,omitempty"`
	Channels  []json.RawMessage `json:"channels,omitempty"`
}

// LoadNanoClawConfig reads nanoclaw config.json
func LoadNanoClawConfig(home string) (map[string]any, error) {
	paths := []string{
		filepath.Join(home, "config.json"),
		filepath.Join(home, "nanoclaw.json"),
	}
	for _, p := range paths {
		data, err := os.ReadFile(p)
		if err == nil {
			var raw map[string]any
			if err := json.Unmarshal(data, &raw); err != nil {
				return nil, fmt.Errorf("invalid nanoclaw config JSON: %w", err)
			}
			return raw, nil
		}
	}
	return nil, fmt.Errorf("nanoclaw config not found in %s", home)
}

// ConvertNanoClawConfig converts nanoclaw config to picoclaw format
func ConvertNanoClawConfig(raw map[string]any) (map[string]any, []string, error) {
	var warnings []string
	result := map[string]any{}

	// Extraer providers
	if providersRaw, ok := raw["providers"].([]any); ok {
		picoProv := map[string]any{}
		for _, p := range providersRaw {
			pm, ok := p.(map[string]any)
			if !ok {
				continue
			}
			pType, _ := pm["type"].(string)
			apiKey, _ := pm["apiKey"].(string)
			if pType != "" && apiKey != "" {
				picoProv[pType] = map[string]any{"api_key": apiKey}
			}
		}
		if len(picoProv) > 0 {
			result["providers"] = picoProv
		}
	}

	// Extraer primer agente como defaults
	if agentsRaw, ok := raw["agents"].([]any); ok && len(agentsRaw) > 0 {
		if a, ok := agentsRaw[0].(map[string]any); ok {
			model, _ := a["model"].(string)
			result["agents"] = map[string]any{
				"defaults": map[string]any{
					"model_name": model,
				},
			}
		}
		if len(agentsRaw) > 1 {
			warnings = append(warnings, fmt.Sprintf(
				"%d agents found — only first agent defaults migrated; review manually", len(agentsRaw),
			))
		}
	}

	// Extraer canal telegram
	if channelsRaw, ok := raw["channels"].([]any); ok {
		for _, ch := range channelsRaw {
			chm, ok := ch.(map[string]any)
			if !ok {
				continue
			}
			if tg, ok := chm["telegram"].(map[string]any); ok {
				token, _ := tg["token"].(string)
				if token != "" {
					result["channels"] = map[string]any{
						"telegram": map[string]any{
							"token": token,
						},
					}
				}
			}
		}
	}

	return result, warnings, nil
}

// PlanNanoClawMigration produces migration actions for nanoclaw → picoclaw
func PlanNanoClawMigration(opts Options, nanoHome, picoHome string) ([]Action, []string, error) {
	var actions []Action
	var warnings []string

	configPath := filepath.Join(nanoHome, "config.json")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		configPath = filepath.Join(nanoHome, "nanoclaw.json")
	}

	if !opts.WorkspaceOnly {
		actions = append(actions, Action{
			Type:        ActionConvertConfig,
			Source:      configPath,
			Destination: filepath.Join(picoHome, "config.json"),
			Description: "convert NanoClaw config to PicoClaw format",
		})
	}

	if !opts.ConfigOnly {
		srcWS := filepath.Join(nanoHome, "workspace")
		dstWS := filepath.Join(picoHome, "workspace")
		if _, err := os.Stat(srcWS); err == nil {
			wsActions, err := PlanWorkspaceMigration(srcWS, dstWS, opts.Force)
			if err != nil {
				return nil, nil, fmt.Errorf("planning nanoclaw workspace migration: %w", err)
			}
			actions = append(actions, wsActions...)
		} else {
			warnings = append(warnings, "nanoclaw workspace not found, skipping workspace migration")
		}

		// groups/default/CLAUDE.md → workspace/AGENTS.md
		groupsPath := filepath.Join(nanoHome, "groups")
		if _, err := os.Stat(groupsPath); err == nil {
			claudePath := filepath.Join(groupsPath, "default", "CLAUDE.md")
			if _, err := os.Stat(claudePath); err == nil {
				actions = append(actions, Action{
					Type:        ActionCopy,
					Source:      claudePath,
					Destination: filepath.Join(dstWS, "AGENTS.md"),
					Description: "migrate default group CLAUDE.md → AGENTS.md",
				})
			}
		}
	}

	return actions, warnings, nil
}

// PrintNanoClawPlan shows a dry-run plan for nanoclaw migration
func PrintNanoClawPlan(actions []Action, warnings []string) {
	fmt.Println("📋 NanoClaw → PicoClaw Migration Plan")
	fmt.Println()

	for i, action := range actions {
		icon := "📄"
		switch action.Type {
		case ActionConvertConfig:
			icon = "⚙️"
		case ActionCreateDir:
			icon = "📁"
		case ActionCopy:
			icon = "📦"
		}
		fmt.Printf("%d. %s %s\n", i+1, icon, action.Description)
		fmt.Printf("   From:  %s\n", action.Source)
		fmt.Printf("   To:    %s\n", action.Destination)
		fmt.Println()
	}

	if len(warnings) > 0 {
		fmt.Println("⚠️  Warnings:")
		for _, w := range warnings {
			fmt.Printf("   - %s\n", w)
		}
		fmt.Println()
	}
}
