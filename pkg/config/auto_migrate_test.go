package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestSaveConfig_AutoMigrate verifies that SaveConfig auto-migrates
// configs that don't have context_manager set.
func TestSaveConfig_AutoMigrate(t *testing.T) {
	// Create a minimal config without context_manager
	cfg := &Config{
		Agents: AgentsConfig{
			Defaults: AgentDefaults{
				Workspace:         "~/.picoclaw/workspace",
				MaxTokens:         8192,
				MaxToolIterations: 20,
				// ContextManager intentionally NOT set
			},
		},
	}

	// Save to temp file
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "config.json")
	if err := SaveConfig(path, cfg); err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	// Read the file and check it contains context_manager
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}

	if !strings.Contains(string(data), `"context_manager": "seahorse"`) {
		t.Errorf("Saved config does not contain context_manager: %s", string(data))
	}

	// Load it back and verify
	loaded, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if loaded.Agents.Defaults.ContextManager != "seahorse" {
		t.Errorf("Loaded config ContextManager = %q, want %q",
			loaded.Agents.Defaults.ContextManager, "seahorse")
	}
}

// TestLoadConfig_RuntimeDefault verifies that LoadConfig injects
// context_manager default when it's absent from the file.
func TestLoadConfig_RuntimeDefault(t *testing.T) {
	// Write a minimal legacy config (without context_manager)
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "config.json")
	legacyConfig := `{
  "agents": {
    "defaults": {
      "workspace": "~/.picoclaw/workspace",
      "max_tokens": 8192,
      "max_tool_iterations": 20
    },
    "list": []
  },
  "model_list": []
}`
	if err := os.WriteFile(path, []byte(legacyConfig), 0o644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	// LoadConfig should inject the default
	loaded, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if loaded.Agents.Defaults.ContextManager != "seahorse" {
		t.Errorf("LoadConfig did not inject default: got %q, want %q",
			loaded.Agents.Defaults.ContextManager, "seahorse")
	}
}
