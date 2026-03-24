// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package tasklock

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Manager handles a workspace's TaskLocks in-memory and keeps them synced with the disk.
type Manager struct {
	workspace string
	locksDir  string

	mu    sync.RWMutex
	locks map[string]*TaskLock // indexed by taskID
}

// NewManager creates a manager for a given workspace and hydrates pre-existing locks perfectly.
func NewManager(workspace string) *Manager {
	locksDir := filepath.Join(workspace, ".locks")
	os.MkdirAll(locksDir, 0o755)

	m := &Manager{
		workspace: workspace,
		locksDir:  locksDir,
		locks:     make(map[string]*TaskLock),
	}

	m.hydrateFromDisk()
	return m
}

// hydrateFromDisk populates the in-memory store from existing .lock files.
// Used during reboot to recover state.
func (m *Manager) hydrateFromDisk() {
	entries, err := os.ReadDir(m.locksDir)
	if err != nil {
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".lock") {
			lockPath := filepath.Join(m.locksDir, entry.Name())

			// Initialize empty lock just to hold the path and try loading it
			tl := &TaskLock{filePath: lockPath}
			if err := tl.Load(); err == nil {
				tl.filePath = lockPath // ensure path remains mapped in struct after JSON unmarshal
				m.locks[tl.TaskID] = tl
			}
		}
	}
}

// CreateLock initializes a new lock for a task assignment.
func (m *Manager) CreateLock(taskID, assignedAgent, parentAgent, sessionKey string) (*TaskLock, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.locks[taskID]; exists {
		return nil, fmt.Errorf("task lock '%s' already exists", taskID)
	}

	fileName := fmt.Sprintf("task_%s_%s.lock", taskID, assignedAgent)
	filePath := filepath.Join(m.locksDir, fileName)

	tl := NewTaskLock(taskID, assignedAgent, parentAgent, sessionKey, filePath)
	if err := tl.SaveAtomic(); err != nil {
		return nil, fmt.Errorf("failed to save initial task lock: %w", err)
	}

	m.locks[taskID] = tl
	return tl, nil
}

// GetLock retrieves a specific task lock by its unique ID.
func (m *Manager) GetLock(taskID string) (*TaskLock, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	tl, ok := m.locks[taskID]
	return tl, ok
}

// GetLocksByAgent returns all active locks for a specific agent.
func (m *Manager) GetLocksByAgent(agentID string) []*TaskLock {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var agentLocks []*TaskLock
	for _, tl := range m.locks {
		if tl.AssignedAgent == agentID {
			agentLocks = append(agentLocks, tl)
		}
	}
	return agentLocks
}

// GetAllActiveLocks returns all currently held locks across all agents.
func (m *Manager) GetAllActiveLocks() []*TaskLock {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var all []*TaskLock
	for _, tl := range m.locks {
		all = append(all, tl)
	}
	return all
}

// IsFileLocked checks if any active task lock has laid claim to evaluating/editing a specific file.
// Returns heavily boolean if true, along with the Agent ID that currently holds it.
func (m *Manager) IsFileLocked(filename string) (bool, string) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, tl := range m.locks {
		// Only consider locks actively in progress
		if tl.GetStatus() == StatusInProgress {
			lockedFiles := tl.GetLockedFiles()
			for _, file := range lockedFiles {
				if file == filename || strings.HasSuffix(filename, file) || strings.HasSuffix(file, filename) {
					return true, tl.AssignedAgent
				}
			}
		}
	}
	return false, ""
}

// RemoveLock deletes the lock permanently from both memory and disk.
func (m *Manager) RemoveLock(taskID string) error {
	m.mu.Lock()
	tl, exists := m.locks[taskID]
	if exists {
		delete(m.locks, taskID)
	}
	m.mu.Unlock()

	if exists {
		return tl.Delete()
	}
	return nil
}

// StartWatchdog runs a background routine to clean up dead/ghost locks (agent crashes).
func (m *Manager) StartWatchdog(timeout time.Duration, cleanupInterval time.Duration) {
	go func() {
		ticker := time.NewTicker(cleanupInterval)
		defer ticker.Stop()
		for range ticker.C {
			m.pruneGhostLocks(timeout)
		}
	}()
}

func (m *Manager) pruneGhostLocks(timeout time.Duration) {
	m.mu.Lock()
	var toDelete []string

	for id, tl := range m.locks {
		// Check if it's stuck
		if tl.GetStatus() == StatusInProgress || tl.GetStatus() == StatusNetworkRetry ||
			tl.GetStatus() == StatusRecovering {
			lastActive := tl.GetLastUpdated()
			if time.Since(lastActive) > timeout {
				// Assumed dead / crashed goroutine due to node failure
				toDelete = append(toDelete, id)
			}
		}
	}
	m.mu.Unlock()

	// Safely trigger standard cleanup
	for _, id := range toDelete {
		m.RemoveLock(id)
	}
}
