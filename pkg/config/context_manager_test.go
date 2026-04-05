package config

import (
	"testing"
)

func TestDefaultConfig_ContextManager(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Agents.Defaults.ContextManager != "seahorse" {
		t.Errorf("DefaultConfig().Agents.Defaults.ContextManager = %q, want %q",
			cfg.Agents.Defaults.ContextManager, "seahorse")
	}
	if cfg.Agents.Defaults.ContextManagerConfig == nil {
		t.Error("DefaultConfig().Agents.Defaults.ContextManagerConfig is nil, expected non-nil")
	}
}

func TestTemplateDefaultConfig_ContextManager(t *testing.T) {
	cfg := TemplateDefaultConfig()
	if cfg.Agents.Defaults.ContextManager != "seahorse" {
		t.Errorf("TemplateDefaultConfig().Agents.Defaults.ContextManager = %q, want %q",
			cfg.Agents.Defaults.ContextManager, "seahorse")
	}
}

func TestOpenRouterFreeDefaultConfig_ContextManager(t *testing.T) {
	cfg := OpenRouterFreeDefaultConfig()
	if cfg.Agents.Defaults.ContextManager != "seahorse" {
		t.Errorf("OpenRouterFreeDefaultConfig().Agents.Defaults.ContextManager = %q, want %q",
			cfg.Agents.Defaults.ContextManager, "seahorse")
	}
}

func TestGLMDefaultConfig_ContextManager(t *testing.T) {
	cfg := GLMDefaultConfig()
	if cfg.Agents.Defaults.ContextManager != "seahorse" {
		t.Errorf("GLMDefaultConfig().Agents.Defaults.ContextManager = %q, want %q",
			cfg.Agents.Defaults.ContextManager, "seahorse")
	}
}

func TestOpenAIDefaultConfig_ContextManager(t *testing.T) {
	cfg := OpenAIDefaultConfig()
	if cfg.Agents.Defaults.ContextManager != "seahorse" {
		t.Errorf("OpenAIDefaultConfig().Agents.Defaults.ContextManager = %q, want %q",
			cfg.Agents.Defaults.ContextManager, "seahorse")
	}
}

func TestQwenDefaultConfig_ContextManager(t *testing.T) {
	cfg := QwenDefaultConfig(false)
	if cfg.Agents.Defaults.ContextManager != "seahorse" {
		t.Errorf("QwenDefaultConfig().Agents.Defaults.ContextManager = %q, want %q",
			cfg.Agents.Defaults.ContextManager, "seahorse")
	}
}

func TestGeminiDefaultConfig_ContextManager(t *testing.T) {
	cfg := GeminiDefaultConfig()
	if cfg.Agents.Defaults.ContextManager != "seahorse" {
		t.Errorf("GeminiDefaultConfig().Agents.Defaults.ContextManager = %q, want %q",
			cfg.Agents.Defaults.ContextManager, "seahorse")
	}
}

func TestOpenRouterDefaultConfig_ContextManager(t *testing.T) {
	cfg := OpenRouterDefaultConfig()
	if cfg.Agents.Defaults.ContextManager != "seahorse" {
		t.Errorf("OpenRouterDefaultConfig().Agents.Defaults.ContextManager = %q, want %q",
			cfg.Agents.Defaults.ContextManager, "seahorse")
	}
}
