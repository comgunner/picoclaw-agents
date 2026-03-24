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
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/comgunner/picoclaw/pkg/providers"
)

type TaskStatus string

const (
	StatusInProgress   TaskStatus = "in_progress"
	StatusNetworkRetry TaskStatus = "network_retry"
	StatusRecovering   TaskStatus = "recovering"
	StatusDone         TaskStatus = "done"
)

// TaskLock represents the persistent state (flight recorder) for a task lock and hydration check.
type TaskLock struct {
	mu sync.RWMutex `json:"-"`

	TaskID            string              `json:"task_id"`
	AssignedAgent     string              `json:"assigned_agent"`
	ParentAgent       string              `json:"parent_agent"`
	Status            TaskStatus          `json:"status"`
	StartedAt         time.Time           `json:"started_at"`
	LastUpdated       time.Time           `json:"last_updated"`
	Objective         string              `json:"objective"`
	ModelUsed         string              `json:"model_used"`
	SessionKey        string              `json:"session_key"`
	LockedFiles       []string            `json:"locked_files"`
	ContextCheckpoint []providers.Message `json:"context_checkpoint"`
	PendingAction     string              `json:"pending_action"`
	RetryCount        int                 `json:"retry_count"`

	// filePath stores the absolute path to the `.lock` file for this task.
	filePath string `json:"-"`
}

func NewTaskLock(taskID, assignedAgent, parentAgent, sessionKey, filePath string) *TaskLock {
	now := time.Now()
	return &TaskLock{
		TaskID:        taskID,
		AssignedAgent: assignedAgent,
		ParentAgent:   parentAgent,
		Status:        StatusInProgress,
		StartedAt:     now,
		LastUpdated:   now,
		SessionKey:    sessionKey,
		filePath:      filePath,
	}
}

// UpdateState safely updates the status, action, and context, then flushes.
func (tl *TaskLock) UpdateState(status TaskStatus, pendingAction string, context []providers.Message) error {
	tl.mu.Lock()
	tl.Status = status
	tl.PendingAction = pendingAction
	if context != nil {
		tl.ContextCheckpoint = context
	}
	tl.LastUpdated = time.Now()
	tl.mu.Unlock()

	return tl.SaveAtomic()
}

// ClaimFiles safely updates the locked files slice and flushes.
func (tl *TaskLock) ClaimFiles(files []string) error {
	tl.mu.Lock()
	tl.LockedFiles = files
	tl.LastUpdated = time.Now()
	tl.mu.Unlock()

	return tl.SaveAtomic()
}

// IncrementRetry safely bumps the error retry counter.
func (tl *TaskLock) IncrementRetry() error {
	tl.mu.Lock()
	tl.RetryCount++
	tl.LastUpdated = time.Now()
	tl.mu.Unlock()

	return tl.SaveAtomic()
}

func (tl *TaskLock) GetStatus() TaskStatus {
	tl.mu.RLock()
	defer tl.mu.RUnlock()
	return tl.Status
}

func (tl *TaskLock) GetLastUpdated() time.Time {
	tl.mu.RLock()
	defer tl.mu.RUnlock()
	return tl.LastUpdated
}

func (tl *TaskLock) GetLockedFiles() []string {
	tl.mu.RLock()
	defer tl.mu.RUnlock()
	// Return a copy to avoid external mutation
	files := make([]string, len(tl.LockedFiles))
	copy(files, tl.LockedFiles)
	return files
}

func (tl *TaskLock) GetContext() []providers.Message {
	tl.mu.RLock()
	defer tl.mu.RUnlock()
	ctx := make([]providers.Message, len(tl.ContextCheckpoint))
	copy(ctx, tl.ContextCheckpoint)
	return ctx
}

func (tl *TaskLock) GetSessionKey() string {
	tl.mu.RLock()
	defer tl.mu.RUnlock()
	return tl.SessionKey
}

func (tl *TaskLock) GetFilePath() string {
	tl.mu.RLock()
	defer tl.mu.RUnlock()
	return tl.filePath
}

// SaveAtomic writes the JSON lock structure to a temporary file, then renames it over the actual file.
func (tl *TaskLock) SaveAtomic() error {
	tl.mu.RLock()
	defer tl.mu.RUnlock()

	if tl.filePath == "" {
		return fmt.Errorf("task lock has no file path configured")
	}

	tempFile := tl.filePath + ".tmp"
	data, err := json.MarshalIndent(tl, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal task lock '%s': %w", tl.TaskID, err)
	}

	if err := os.WriteFile(tempFile, data, 0o644); err != nil {
		return fmt.Errorf("failed to write temp file for task lock '%s': %w", tl.TaskID, err)
	}

	// Atomic rename from temp to target (POSIX atomic guarantee)
	if err := os.Rename(tempFile, tl.filePath); err != nil {
		os.Remove(tempFile) // cleanup
		return fmt.Errorf("failed to rename temp file for task lock '%s': %w", tl.TaskID, err)
	}

	return nil
}

// Load reads the JSON data from disk. Not safe for concurrent reads matching writes,
// assumed to be called during manager's initialization/hydration phase.
func (tl *TaskLock) Load() error {
	tl.mu.Lock()
	defer tl.mu.Unlock()

	if tl.filePath == "" {
		return fmt.Errorf("task lock has no file path configured")
	}

	data, err := os.ReadFile(tl.filePath)
	if err != nil {
		return fmt.Errorf("failed to read task lock file '%s': %w", tl.filePath, err)
	}

	if err := json.Unmarshal(data, tl); err != nil {
		return fmt.Errorf("failed to unmarshal task lock '%s': %w", tl.TaskID, err)
	}

	return nil
}

// Delete permanently removes the active task lock files from the disk (e.g. on clean finish).
func (tl *TaskLock) Delete() error {
	tl.mu.Lock()
	defer tl.mu.Unlock()

	if tl.filePath != "" {
		err := os.Remove(tl.filePath)
		if err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	return nil
}
