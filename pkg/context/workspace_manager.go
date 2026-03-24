// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package context

import (
	"os"
	"path/filepath"

	"github.com/comgunner/picoclaw/pkg/logger"
)

// WorkspaceManager manages the workspace directory structure.
type WorkspaceManager struct {
	BasePath string
}

// NewWorkspaceManager creates a new workspace manager.
func NewWorkspaceManager(basePath string) *WorkspaceManager {
	wm := &WorkspaceManager{BasePath: basePath}
	wm.ensureStructure()
	return wm
}

// ensureStructure creates the required directory structure.
func (wm *WorkspaceManager) ensureStructure() {
	dirs := []string{
		"active",   // Current session (loaded into context)
		"memory",   // Persistent but not loaded
		"cold",     // Archived, never auto-loaded
		"temp",     // Auto-delete after 1h
		"sessions", // Per-user/channel sessions
		"state",    // System state
		"scripts",  // User scripts
	}

	for _, dir := range dirs {
		fullPath := filepath.Join(wm.BasePath, dir)
		err := os.MkdirAll(fullPath, 0o755)
		if err != nil {
			logger.ErrorCF("context", "Failed to create workspace directory",
				map[string]any{
					"path":  fullPath,
					"error": err,
				})
		}
	}

	logger.InfoCF("context", "Workspace structure ensured",
		map[string]any{
			"base_path": wm.BasePath,
		})
}

// StoreActive stores data in the active directory.
func (wm *WorkspaceManager) StoreActive(data []byte, filename string) (string, error) {
	path := filepath.Join(wm.BasePath, "active", filename)
	err := os.WriteFile(path, data, 0o644)
	if err != nil {
		return "", err
	}
	return path, nil
}

// StoreCold stores data in the cold storage directory.
func (wm *WorkspaceManager) StoreCold(data []byte, filename string) (string, error) {
	path := filepath.Join(wm.BasePath, "cold", filename)
	err := os.WriteFile(path, data, 0o644)
	if err != nil {
		return "", err
	}
	return path, nil
}

// StoreTemp stores temporary data.
func (wm *WorkspaceManager) StoreTemp(data []byte, filename string) (string, error) {
	path := filepath.Join(wm.BasePath, "temp", filename)
	err := os.WriteFile(path, data, 0o644)
	if err != nil {
		return "", err
	}
	return path, nil
}

// GetActivePath returns the path to a file in the active directory.
func (wm *WorkspaceManager) GetActivePath(filename string) string {
	return filepath.Join(wm.BasePath, "active", filename)
}

// GetColdPath returns the path to a file in cold storage.
func (wm *WorkspaceManager) GetColdPath(filename string) string {
	return filepath.Join(wm.BasePath, "cold", filename)
}
