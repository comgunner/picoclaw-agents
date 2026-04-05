package agent

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAgentState_SaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	statePath := filepath.Join(tmpDir, "state.json")

	state := &AgentState{
		CurrentModel: "openai/gpt-5.4",
		AgentID:      "engineering_manager",
		Channel:      "telegram",
	}

	// Save
	err := state.Save(statePath)
	assert.NoError(t, err)

	// Load
	loaded, err := LoadAgentState(statePath)
	assert.NoError(t, err)
	assert.Equal(t, "openai/gpt-5.4", loaded.CurrentModel)
	assert.Equal(t, "engineering_manager", loaded.AgentID)
	assert.Equal(t, "telegram", loaded.Channel)
	assert.False(t, loaded.UpdatedAt.IsZero())
}

func TestAgentState_LoadNonExistent_UsesDefaults(t *testing.T) {
	tmpDir := t.TempDir()
	statePath := filepath.Join(tmpDir, "nonexistent.json")

	state, err := LoadAgentState(statePath)

	assert.NoError(t, err)
	assert.Equal(t, "", state.CurrentModel) // Default vacío
}

func TestAgentState_SavePermissionDenied(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping on Windows")
	}

	tmpDir := t.TempDir()
	statePath := filepath.Join(tmpDir, "state.json")

	// Crear directorio sin permisos de escritura
	os.Chmod(tmpDir, 0o555)
	defer os.Chmod(tmpDir, 0o755)

	state := &AgentState{CurrentModel: "test"}
	err := state.Save(statePath)

	assert.Error(t, err)
}

func TestAgentState_AtomicSave(t *testing.T) {
	tmpDir := t.TempDir()
	statePath := filepath.Join(tmpDir, "state.json")

	state := &AgentState{CurrentModel: "openai/gpt-5.4"}
	err := state.Save(statePath)
	assert.NoError(t, err)

	// Verificar que no hay archivos temporales residuales
	files, _ := filepath.Glob(filepath.Join(tmpDir, "*.tmp"))
	assert.Empty(t, files)

	// Verificar que el archivo principal existe y es válido
	loaded, err := LoadAgentState(statePath)
	assert.NoError(t, err)
	assert.Equal(t, "openai/gpt-5.4", loaded.CurrentModel)
}
