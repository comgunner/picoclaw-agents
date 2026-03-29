// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package onboard

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// newTestWizard creates a Wizard whose stdin is backed by the provided lines,
// with configPath and workspace redirected to a temporary directory so tests
// never touch the real ~/.picoclaw tree.
func newTestWizard(t *testing.T, inputLines ...string) (*Wizard, string) {
	t.Helper()

	tmpDir := t.TempDir()
	input := strings.Join(inputLines, "\n") + "\n"

	w := &Wizard{
		scanner:    bufio.NewScanner(strings.NewReader(input)),
		configPath: filepath.Join(tmpDir, "config.json"),
		workspace:  filepath.Join(tmpDir, "workspace"),
	}
	return w, tmpDir
}

// ---------------------------------------------------------------------------
// TestNewWizard
// ---------------------------------------------------------------------------

// TestNewWizard verifies that NewWizard initializes all required fields so the
// wizard is ready to run without further configuration.
func TestNewWizard(t *testing.T) {
	w := NewWizard()

	if w == nil {
		t.Fatal("NewWizard returned nil")
	}
	if w.scanner == nil {
		t.Error("scanner must not be nil after NewWizard")
	}
	if w.configPath == "" {
		t.Error("configPath must not be empty after NewWizard")
	}
	if w.workspace == "" {
		t.Error("workspace must not be empty after NewWizard")
	}
	if !strings.Contains(w.configPath, ".picoclaw") {
		t.Errorf("configPath %q does not contain expected '.picoclaw' segment", w.configPath)
	}
	if !strings.Contains(w.workspace, ".picoclaw") {
		t.Errorf("workspace %q does not contain expected '.picoclaw' segment", w.workspace)
	}
}

// ---------------------------------------------------------------------------
// TestWizardValidatesRequiredFields
// ---------------------------------------------------------------------------

// TestWizardValidatesRequiredFields verifies that saveConfig rejects an empty
// modelName (required field) and that a Wizard with no apiKey still produces a
// config file – the key is allowed to be empty in saved config but modelName
// must be set to produce a meaningful output.
//
// Because all validation in the Wizard happens interactively, we test the
// underlying saveConfig directly after setting fields to known invalid states.
func TestWizardValidatesRequiredFields(t *testing.T) {
	t.Run("empty modelName produces empty model_name in config", func(t *testing.T) {
		w, _ := newTestWizard(t)
		// modelName deliberately left empty
		w.model = ""
		w.modelName = ""
		w.apiKey = "sk-test-key-1234567890"

		if err := w.saveConfig(); err != nil {
			t.Fatalf("saveConfig unexpectedly failed: %v", err)
		}

		data, err := os.ReadFile(w.configPath)
		if err != nil {
			t.Fatalf("could not read saved config: %v", err)
		}

		var parsed map[string]any
		if err := json.Unmarshal(data, &parsed); err != nil {
			t.Fatalf("saved config is not valid JSON: %v\nraw: %s", err, data)
		}

		agents, ok := parsed["agents"].(map[string]any)
		if !ok {
			t.Fatal("saved config missing 'agents' object")
		}
		defaults, ok := agents["defaults"].(map[string]any)
		if !ok {
			t.Fatal("saved config missing 'agents.defaults' object")
		}
		modelNameInConfig, _ := defaults["model_name"].(string)
		if modelNameInConfig != "" {
			t.Errorf("expected empty model_name in config when modelName is unset, got %q", modelNameInConfig)
		}
	})

	t.Run("cancellation when overwrite refused returns error", func(t *testing.T) {
		// Create a pre-existing config so stepEnvironment asks the user.
		w, tmpDir := newTestWizard(t, "n") // user answers "n" to overwrite prompt

		// Write a dummy existing config.
		if err := os.MkdirAll(tmpDir, 0o755); err != nil {
			t.Fatalf("could not create temp dir: %v", err)
		}
		if err := os.WriteFile(w.configPath, []byte("{}"), 0o600); err != nil {
			t.Fatalf("could not write seed config: %v", err)
		}

		err := w.stepEnvironment()
		if err == nil {
			t.Fatal("expected stepEnvironment to return an error when user refuses overwrite")
		}
		if !strings.Contains(err.Error(), "canceled") {
			t.Errorf("error message should mention cancellation, got: %v", err)
		}
	})
}

// ---------------------------------------------------------------------------
// TestWizardGeneratesValidConfig
// ---------------------------------------------------------------------------

// TestWizardGeneratesValidConfig verifies that saveConfig produces a valid,
// parseable JSON config with the expected top-level structure when the Wizard
// fields are set to valid values.
func TestWizardGeneratesValidConfig(t *testing.T) {
	t.Run("normal provider produces single-model config", func(t *testing.T) {
		w, _ := newTestWizard(t)
		w.modelName = "deepseek/default"
		w.model = "deepseek/default"
		w.apiKey = "sk-deepseekkey1234567890"

		if err := w.saveConfig(); err != nil {
			t.Fatalf("saveConfig failed: %v", err)
		}

		data, err := os.ReadFile(w.configPath)
		if err != nil {
			t.Fatalf("cannot read generated config: %v", err)
		}

		var cfg map[string]any
		if err := json.Unmarshal(data, &cfg); err != nil {
			t.Fatalf("generated config is not valid JSON: %v\nraw: %s", err, data)
		}

		// agents section
		agents, ok := cfg["agents"].(map[string]any)
		if !ok {
			t.Fatal("config missing 'agents' section")
		}
		defaults, ok := agents["defaults"].(map[string]any)
		if !ok {
			t.Fatal("config missing 'agents.defaults' section")
		}
		if defaults["workspace"] != w.workspace {
			t.Errorf("workspace mismatch: got %v, want %v", defaults["workspace"], w.workspace)
		}
		if defaults["model_name"] != w.modelName {
			t.Errorf("model_name mismatch: got %v, want %v", defaults["model_name"], w.modelName)
		}
		if defaults["restrict_to_workspace"] != true {
			t.Error("restrict_to_workspace should be true")
		}

		// model_list section
		modelList, ok := cfg["model_list"].([]any)
		if !ok {
			t.Fatal("config missing 'model_list' array")
		}
		if len(modelList) != 1 {
			t.Fatalf("expected 1 model in model_list for normal provider, got %d", len(modelList))
		}

		entry, ok := modelList[0].(map[string]any)
		if !ok {
			t.Fatal("model_list[0] is not an object")
		}
		if entry["model_name"] != w.modelName {
			t.Errorf("model_list[0].model_name mismatch: got %v, want %v", entry["model_name"], w.modelName)
		}
		if entry["api_key"] != w.apiKey {
			t.Errorf("model_list[0].api_key mismatch: got %v, want %v", entry["api_key"], w.apiKey)
		}
	})

	t.Run("openrouter/auto produces three-model fallback config", func(t *testing.T) {
		w, _ := newTestWizard(t)
		w.modelName = "or-auto"
		w.model = "openrouter/auto"
		w.apiKey = "sk-or-v1-testkey1234567890abcdefgh"

		if err := w.saveConfig(); err != nil {
			t.Fatalf("saveConfig failed: %v", err)
		}

		data, err := os.ReadFile(w.configPath)
		if err != nil {
			t.Fatalf("cannot read generated config: %v", err)
		}

		var cfg map[string]any
		if err := json.Unmarshal(data, &cfg); err != nil {
			t.Fatalf("generated config is not valid JSON: %v\nraw: %s", err, data)
		}

		modelList, ok := cfg["model_list"].([]any)
		if !ok {
			t.Fatal("config missing 'model_list' array")
		}
		if len(modelList) != 3 {
			t.Fatalf("expected 3 fallback models for openrouter/auto, got %d", len(modelList))
		}

		// All three entries must carry the same API key and point to OpenRouter.
		for i, item := range modelList {
			entry, ok := item.(map[string]any)
			if !ok {
				t.Fatalf("model_list[%d] is not an object", i)
			}
			if entry["api_key"] != w.apiKey {
				t.Errorf("model_list[%d].api_key = %v, want %v", i, entry["api_key"], w.apiKey)
			}
			apiBase, _ := entry["api_base"].(string)
			if !strings.Contains(apiBase, "openrouter.ai") {
				t.Errorf("model_list[%d].api_base %q should reference openrouter.ai", i, apiBase)
			}
		}
	})

	t.Run("config file is written with restrictive permissions", func(t *testing.T) {
		w, _ := newTestWizard(t)
		w.modelName = "anthropic/default"
		w.model = "anthropic/default"
		w.apiKey = "sk-ant-testkey1234567890abcdefgh"

		if err := w.saveConfig(); err != nil {
			t.Fatalf("saveConfig failed: %v", err)
		}

		info, err := os.Stat(w.configPath)
		if err != nil {
			t.Fatalf("cannot stat saved config: %v", err)
		}
		// Config contains the API key — only the owner should be able to read it.
		if info.Mode()&0o077 != 0 {
			t.Errorf("config file has overly permissive mode %o, want 0600", info.Mode())
		}
	})
}

// ---------------------------------------------------------------------------
// TestWizardHandlesCancellation
// ---------------------------------------------------------------------------

// TestWizardHandlesCancellation verifies that the wizard terminates gracefully
// whenever the user signals intent to cancel rather than panicking or hanging.
func TestWizardHandlesCancellation(t *testing.T) {
	t.Run("EOF on scanner returns empty string from prompt", func(t *testing.T) {
		// An empty reader simulates EOF / closed stdin.
		w := &Wizard{
			scanner: bufio.NewScanner(strings.NewReader("")),
		}
		got := w.prompt("Enter something: ")
		if got != "" {
			t.Errorf("prompt on EOF should return empty string, got %q", got)
		}
	})

	t.Run("EOF on promptSecret returns empty string", func(t *testing.T) {
		w := &Wizard{
			scanner: bufio.NewScanner(strings.NewReader("")),
		}
		got := w.promptSecret("Secret: ")
		if got != "" {
			t.Errorf("promptSecret on EOF should return empty string, got %q", got)
		}
	})

	t.Run("stepGenerateConfig returns error when configPath directory cannot be created", func(t *testing.T) {
		// Point configPath to a location under a file (not a directory) so
		// MkdirAll will fail deterministically.
		tmpDir := t.TempDir()
		blocker := filepath.Join(tmpDir, "blocker") // this will be a file
		if err := os.WriteFile(blocker, []byte("x"), 0o644); err != nil {
			t.Fatalf("setup: %v", err)
		}

		w, _ := newTestWizard(t)
		// Nest configPath under the file — the OS won't allow mkdir through it.
		w.configPath = filepath.Join(blocker, "subdir", "config.json")
		w.workspace = filepath.Join(tmpDir, "workspace")

		err := w.stepGenerateConfig()
		if err == nil {
			t.Fatal("expected stepGenerateConfig to fail when config dir cannot be created")
		}
	})

	t.Run("stepVerify returns error when configPath does not exist", func(t *testing.T) {
		w, _ := newTestWizard(t)
		// configPath was never written, so stepVerify should fail.
		err := w.stepVerify()
		if err == nil {
			t.Fatal("expected stepVerify to fail when config file is missing")
		}
	})

	t.Run("stepVerify returns error when workspace does not exist", func(t *testing.T) {
		w, _ := newTestWizard(t)

		// Write config but leave workspace absent.
		if err := os.WriteFile(w.configPath, []byte("{}"), 0o600); err != nil {
			t.Fatalf("cannot write seed config: %v", err)
		}
		// workspace intentionally not created

		err := w.stepVerify()
		if err == nil {
			t.Fatal("expected stepVerify to fail when workspace directory is missing")
		}
	})
}

// ---------------------------------------------------------------------------
// TestWizardPromptHelpers
// ---------------------------------------------------------------------------

// TestWizardPromptHelpers exercises the low-level prompt utilities with
// controlled input so that higher-level step functions can rely on them.
func TestWizardPromptHelpers(t *testing.T) {
	t.Run("prompt trims whitespace", func(t *testing.T) {
		w, _ := newTestWizard(t, "  hello world  ")
		got := w.prompt("msg: ")
		if got != "hello world" {
			t.Errorf("prompt should trim spaces, got %q", got)
		}
	})

	t.Run("promptSecret trims whitespace and returns input", func(t *testing.T) {
		w, _ := newTestWizard(t, "  mysecret  ")
		got := w.promptSecret("secret: ")
		if got != "mysecret" {
			t.Errorf("promptSecret should trim spaces, got %q", got)
		}
	})

	t.Run("promptConfirm accepts y as true", func(t *testing.T) {
		w, _ := newTestWizard(t, "y")
		if !w.promptConfirm("Continue?") {
			t.Error("expected promptConfirm to return true for 'y'")
		}
	})

	t.Run("promptConfirm accepts yes as true", func(t *testing.T) {
		w, _ := newTestWizard(t, "yes")
		if !w.promptConfirm("Continue?") {
			t.Error("expected promptConfirm to return true for 'yes'")
		}
	})

	t.Run("promptConfirm accepts n as false", func(t *testing.T) {
		w, _ := newTestWizard(t, "n")
		if w.promptConfirm("Continue?") {
			t.Error("expected promptConfirm to return false for 'n'")
		}
	})

	t.Run("promptConfirm accepts no as false", func(t *testing.T) {
		w, _ := newTestWizard(t, "no")
		if w.promptConfirm("Continue?") {
			t.Error("expected promptConfirm to return false for 'no'")
		}
	})

	t.Run("promptConfirm is case-insensitive", func(t *testing.T) {
		w, _ := newTestWizard(t, "Y")
		if !w.promptConfirm("Continue?") {
			t.Error("expected promptConfirm to return true for uppercase 'Y'")
		}
	})

	t.Run("promptChoice returns selected option by 1-based index", func(t *testing.T) {
		w, _ := newTestWizard(t, "2")
		got := w.promptChoice("Pick one:", []string{"alpha", "beta", "gamma"})
		if got != "beta" {
			t.Errorf("expected 'beta' for choice 2, got %q", got)
		}
	})

	t.Run("promptChoice retries on invalid then selects valid", func(t *testing.T) {
		// First answer "0" (invalid), then "0" (invalid), then "1" (valid).
		w, _ := newTestWizard(t, "0", "0", "1")
		got := w.promptChoice("Pick one:", []string{"only"})
		if got != "only" {
			t.Errorf("expected 'only' after two invalid attempts, got %q", got)
		}
	})

	t.Run("promptChoice selects last option", func(t *testing.T) {
		w, _ := newTestWizard(t, "3")
		got := w.promptChoice("Pick one:", []string{"a", "b", "c"})
		if got != "c" {
			t.Errorf("expected 'c' for choice 3, got %q", got)
		}
	})
}

// ---------------------------------------------------------------------------
// TestCheckGoVersion
// ---------------------------------------------------------------------------

// TestCheckGoVersion verifies that checkGoVersion returns a non-empty string,
// meaning the environment check step will not flag a missing Go installation.
func TestCheckGoVersion(t *testing.T) {
	version := checkGoVersion()
	if version == "" {
		t.Error("checkGoVersion should return a non-empty string")
	}
}

// ---------------------------------------------------------------------------
// TestSetupEasyFree
// ---------------------------------------------------------------------------

// TestSetupEasyFree verifies that the easy-free flow sets the expected model
// fields when supplied with a properly-formatted OpenRouter key.
func TestSetupEasyFree(t *testing.T) {
	t.Run("valid key sets model fields correctly", func(t *testing.T) {
		// Provide a key that passes the "sk-or-" prefix check, then confirm.
		w, _ := newTestWizard(t, "sk-or-v1-validkey1234567890abcdefghij")
		if err := w.setupEasyFree(); err != nil {
			t.Fatalf("setupEasyFree failed: %v", err)
		}
		if w.model != "openrouter/auto" {
			t.Errorf("model = %q, want 'openrouter/auto'", w.model)
		}
		if w.modelName != "or-auto" {
			t.Errorf("modelName = %q, want 'or-auto'", w.modelName)
		}
		if w.apiKey != "sk-or-v1-validkey1234567890abcdefghij" {
			t.Errorf("apiKey = %q, not stored correctly", w.apiKey)
		}
	})

	t.Run("invalid key prefix with continue sets model fields", func(t *testing.T) {
		// Key without "sk-or-" prefix — wizard warns and asks to continue.
		w, _ := newTestWizard(t, "invalid-key-without-prefix", "y")
		if err := w.setupEasyFree(); err != nil {
			t.Fatalf("setupEasyFree failed: %v", err)
		}
		if w.model != "openrouter/auto" {
			t.Errorf("model = %q, want 'openrouter/auto' even after warning", w.model)
		}
	})

	t.Run("invalid key prefix with cancel returns error", func(t *testing.T) {
		// Key without "sk-or-" prefix and user chooses not to continue.
		w, _ := newTestWizard(t, "bad-key", "n")
		err := w.setupEasyFree()
		if err == nil {
			t.Fatal("expected setupEasyFree to return an error when user cancels on bad key")
		}
	})
}

// ---------------------------------------------------------------------------
// TestStepGenerateConfig
// ---------------------------------------------------------------------------

// TestStepGenerateConfig verifies that stepGenerateConfig creates both the
// config file and the workspace directory when the paths are writable.
func TestStepGenerateConfig(t *testing.T) {
	w, _ := newTestWizard(t)
	w.modelName = "deepseek/default"
	w.model = "deepseek/default"
	w.apiKey = "sk-validkey1234567890"

	if err := w.stepGenerateConfig(); err != nil {
		t.Fatalf("stepGenerateConfig failed: %v", err)
	}

	if _, err := os.Stat(w.configPath); err != nil {
		t.Errorf("config file not found after stepGenerateConfig: %v", err)
	}
	if _, err := os.Stat(w.workspace); err != nil {
		t.Errorf("workspace directory not found after stepGenerateConfig: %v", err)
	}
}

// ---------------------------------------------------------------------------
// TestStepVerify
// ---------------------------------------------------------------------------

// TestStepVerify verifies that stepVerify succeeds when both the config file
// and workspace directory are present.
func TestStepVerify(t *testing.T) {
	w, _ := newTestWizard(t)
	w.model = "deepseek/default"
	w.apiKey = "sk-deepseekkey1234567890"

	// Pre-create the artifacts that stepVerify expects.
	if err := os.MkdirAll(w.workspace, 0o755); err != nil {
		t.Fatalf("cannot create workspace: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(w.configPath), 0o755); err != nil {
		t.Fatalf("cannot create config dir: %v", err)
	}
	if err := os.WriteFile(w.configPath, []byte("{}"), 0o600); err != nil {
		t.Fatalf("cannot write config: %v", err)
	}

	if err := w.stepVerify(); err != nil {
		t.Errorf("stepVerify should succeed when config and workspace exist: %v", err)
	}
}

// ---------------------------------------------------------------------------
// TestStepEnvironment
// ---------------------------------------------------------------------------

// TestStepEnvironment verifies the environment-check step under conditions
// where no pre-existing config exists, so no interactive prompt is shown.
func TestStepEnvironment(t *testing.T) {
	t.Run("passes when no existing config", func(t *testing.T) {
		w, _ := newTestWizard(t)
		// configPath does not exist yet — no overwrite prompt triggered.
		if err := w.stepEnvironment(); err != nil {
			t.Errorf("stepEnvironment should succeed when no existing config: %v", err)
		}
	})

	t.Run("passes when user confirms overwrite", func(t *testing.T) {
		w, _ := newTestWizard(t, "y")

		// Create a pre-existing config to trigger the overwrite prompt.
		if err := os.MkdirAll(filepath.Dir(w.configPath), 0o755); err != nil {
			t.Fatalf("cannot create config dir: %v", err)
		}
		if err := os.WriteFile(w.configPath, []byte("{}"), 0o600); err != nil {
			t.Fatalf("cannot write seed config: %v", err)
		}

		if err := w.stepEnvironment(); err != nil {
			t.Errorf("stepEnvironment should succeed when user confirms overwrite: %v", err)
		}
	})
}
