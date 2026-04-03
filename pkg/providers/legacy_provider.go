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
	"fmt"

	"github.com/comgunner/picoclaw/pkg/config"
)

// CreateProvider creates a provider based on the configuration.
// It uses the model_list configuration (new format) to create providers.
// The old providers config is automatically converted to model_list during config loading.
// Returns the provider, the model ID to use, and any error.
func CreateProvider(cfg *config.Config) (LLMProvider, string, error) {
	model := cfg.Agents.Defaults.GetModelName()
	return CreateProviderForModel(cfg, model)
}

// CreateProviderForModel creates a provider based on the configuration and a specific model name.
// It allows for dynamic provider switching when the user selects a different model via WebUI/CLI.
func CreateProviderForModel(cfg *config.Config, model string) (LLMProvider, string, error) {
	// Ensure model_list is populated from providers config if needed
	if cfg.HasProvidersConfig() {
		providerModels := config.ConvertProvidersToModelList(cfg)
		existingModelNames := make(map[string]bool)
		for _, m := range cfg.ModelList {
			existingModelNames[m.ModelName] = true
		}
		for _, pm := range providerModels {
			if !existingModelNames[pm.ModelName] {
				cfg.ModelList = append(cfg.ModelList, pm)
			}
		}
	}

	// Must have model_list at this point
	if len(cfg.ModelList) == 0 {
		return nil, "", fmt.Errorf("no providers configured. Please add entries to model_list in your config")
	}

	// Get model config from model_list
	modelCfg, err := cfg.GetModelConfig(model)
	if err != nil {
		// Fallback: If model is not in list, try to create a dynamic config
		// if it looks like a protocol/model format or we can infer it.
		protocol, modelID := ExtractProtocol(model)

		// Try to find a base config for this protocol to get API keys/base URLs
		// We look for any model in the list that uses the same protocol
		var baseCfg *config.ModelConfig
		for _, m := range cfg.ModelList {
			p, _ := ExtractProtocol(m.Model)
			if p == protocol {
				baseCfg = &m
				break
			}
		}

		if baseCfg != nil {
			// Create a clone with the requested model ID
			dynamicCfg := *baseCfg
			dynamicCfg.ModelName = model
			dynamicCfg.Model = protocol + "/" + modelID
			modelCfg = &dynamicCfg
		} else {
			return nil, "", fmt.Errorf(
				"model %q not found in model_list and no base provider for %q found: %w",
				model,
				protocol,
				err,
			)
		}
	}

	// Inject global workspace if not set in model config
	if modelCfg.Workspace == "" {
		modelCfg.Workspace = cfg.WorkspacePath()
	}

	// Use factory to create provider
	provider, modelID, err := CreateProviderFromConfig(modelCfg)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create provider for model %q: %w", model, err)
	}

	return provider, modelID, nil
}
