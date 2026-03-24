// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package agents

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Monitor tracks agent health and detects phantom nodes (agents that have locks but don't report in ledger)
type Monitor struct {
	nodes       []string
	taskPath    string
	lockPattern string
	timeout     time.Duration
}

func NewMonitor(nodes []string, taskPath, lockPattern string, timeoutMinutes int) *Monitor {
	if timeoutMinutes <= 0 {
		timeoutMinutes = 5 // default 5 minutes
	}
	return &Monitor{
		nodes:       nodes,
		taskPath:    taskPath,
		lockPattern: lockPattern,
		timeout:     time.Duration(timeoutMinutes) * time.Minute,
	}
}

// Audit performs a cross-check between locks and ledger to detect inconsistencies
func (m *Monitor) Audit() error {
	log.Println("Starting system health audit...")

	// Scan for lock files
	lockFiles, err := m.scanLockFiles()
	if err != nil {
		return fmt.Errorf("failed to scan lock files: %w", err)
	}

	// Load tasks from ledger
	tasks, err := m.loadTasks()
	if err != nil {
		return fmt.Errorf("failed to load tasks: %w", err)
	}

	// Cross-reference locks vs tasks
	for _, lockFile := range lockFiles {
		nodeName := m.extractNodeName(lockFile)
		if nodeName == "" {
			continue
		}

		// Check if this node has a corresponding task in ledger
		hasTask := m.nodeHasTask(nodeName, tasks)
		if !hasTask {
			log.Printf("WARN: LOCK_WITHOUT_LEDGER - Node %s has lock file %s but no task in ledger", nodeName, lockFile)
		}
	}

	// Check for tasks without locks
	for _, task := range tasks {
		nodeName := m.extractNodeFromTask(task)
		if nodeName == "" {
			continue
		}

		hasLock := m.nodeHasLock(nodeName, lockFiles)
		if !hasLock {
			log.Printf("WARN: LEDGER_WITHOUT_LOCK - Task for node %s exists in ledger but no lock file found", nodeName)
		}
	}

	log.Println("Audit completed")
	return nil
}

// DetectPhantoms identifies nodes that have locks but haven't reported in ledger for > timeout
func (m *Monitor) DetectPhantoms() error {
	log.Println("Checking for phantom nodes...")

	lockFiles, err := m.scanLockFiles()
	if err != nil {
		return err
	}

	for _, lockFile := range lockFiles {
		nodeName := m.extractNodeName(lockFile)
		if nodeName == "" {
			continue
		}

		// Check if node has reported recently in ledger
		lastReport, err := m.getLastReportTime(nodeName)
		if err != nil {
			log.Printf("WARN: Could not get last report time for node %s: %v", nodeName, err)
			continue
		}

		if time.Since(lastReport) > m.timeout {
			log.Printf(
				"FLUSH DETECTADO: Node %s has lock file but hasn't reported in %v",
				nodeName,
				time.Since(lastReport),
			)

			// Attempt to flush the phantom node
			if err := m.flushPhantomNode(nodeName, lockFile); err != nil {
				log.Printf("ERROR: Failed to flush phantom node %s: %v", nodeName, err)
			} else {
				log.Printf("INFO: Successfully flushed phantom node %s", nodeName)
			}
		}
	}

	return nil
}

func (m *Monitor) scanLockFiles() ([]string, error) {
	var lockFiles []string

	err := filepath.Walk(filepath.Dir(m.lockPattern), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if strings.Contains(info.Name(), filepath.Base(m.lockPattern)) {
			lockFiles = append(lockFiles, path)
		}

		return nil
	})

	return lockFiles, err
}

func (m *Monitor) loadTasks() ([]map[string]any, error) {
	// This is a simplified implementation - in reality, this would parse the actual task ledger
	// For now, we'll return an empty slice to simulate the functionality
	return []map[string]any{}, nil
}

func (m *Monitor) extractNodeName(lockFile string) string {
	base := filepath.Base(lockFile)
	// Remove the lock pattern suffix to get node name
	nodeName := strings.TrimSuffix(base, filepath.Ext(base))
	return strings.TrimPrefix(nodeName, ".")
}

func (m *Monitor) extractNodeFromTask(task map[string]any) string {
	// Extract node name from task - implementation depends on task structure
	if node, ok := task["node"].(string); ok {
		return node
	}
	return ""
}

func (m *Monitor) nodeHasTask(nodeName string, tasks []map[string]any) bool {
	for _, task := range tasks {
		if n := m.extractNodeFromTask(task); n == nodeName {
			return true
		}
	}
	return false
}

func (m *Monitor) nodeHasLock(nodeName string, lockFiles []string) bool {
	for _, lockFile := range lockFiles {
		if n := m.extractNodeName(lockFile); n == nodeName {
			return true
		}
	}
	return false
}

func (m *Monitor) getLastReportTime(nodeName string) (time.Time, error) {
	// This would check the ledger for the last report time of the node
	// For now, return current time to simulate recent activity
	return time.Now(), nil
}

func (m *Monitor) flushPhantomNode(nodeName, lockFile string) error {
	log.Printf("Flushing phantom node %s by removing lock file %s", nodeName, lockFile)

	// Remove the lock file
	if err := os.Remove(lockFile); err != nil {
		return fmt.Errorf("failed to remove lock file: %w", err)
	}

	// Here we would typically trigger task reassignment
	// For now, just log the action
	log.Printf("Reassigning tasks from phantom node %s", nodeName)

	return nil
}

// Run continuously monitors the system health
func (m *Monitor) Run(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute) // Check every minute
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Monitor stopped")
			return
		case <-ticker.C:
			if err := m.Audit(); err != nil {
				log.Printf("Audit error: %v", err)
			}
			if err := m.DetectPhantoms(); err != nil {
				log.Printf("Phantom detection error: %v", err)
			}
		}
	}
}
