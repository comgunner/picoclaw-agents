// PicoClaw - Ultra-lightweight personal AI agent
// Agent Management - Runtime Instance Tracking
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package management

import (
	"sync"
	"time"
)

// AgentStatus represents the current lifecycle state of a running agent.
type AgentStatus string

const (
	StatusIdle    AgentStatus = "idle"
	StatusActive  AgentStatus = "active"
	StatusBusy    AgentStatus = "busy"
	StatusError   AgentStatus = "error"
	StatusStopped AgentStatus = "stopped"
)

// TaskInfo describes a single unit of work assigned to an agent.
type TaskInfo struct {
	ID          string     `json:"id"`
	Task        string     `json:"task"`
	Status      string     `json:"status"`   // queued | processing | completed | failed
	Progress    string     `json:"progress"` // e.g. "45%" or "2/5 steps"
	StartedAt   time.Time  `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	Error       string     `json:"error,omitempty"`
}

// AgentQueue holds the task queues for a single agent instance.
type AgentQueue struct {
	ActiveTasks []TaskInfo `json:"active_tasks"`
	QueuedTasks []TaskInfo `json:"queued_tasks"`
	Completed   []TaskInfo `json:"completed"`
}

// AgentMetrics tracks performance counters and resource usage for an agent.
type AgentMetrics struct {
	ActiveTasks    int           `json:"active_tasks"`
	QueuedTasks    int           `json:"queued_tasks"`
	CompletedToday int           `json:"completed_today"`
	FailedToday    int           `json:"failed_today"`
	AvgCompletion  time.Duration `json:"avg_completion_ns"` // stored as nanoseconds for JSON
	CPUUsage       float64       `json:"cpu_usage"`
	MemoryUsage    int64         `json:"memory_usage_bytes"`
	TokenUsage     int64         `json:"token_usage"`
	LastActivity   time.Time     `json:"last_activity"`
}

// AgentInstance tracks runtime state for a single active agent.
// All exported methods are safe for concurrent use.
type AgentInstance struct {
	ID           string
	Status       AgentStatus
	Queue        *AgentQueue
	Metrics      *AgentMetrics
	StartedAt    time.Time
	LastActivity time.Time
	mu           sync.RWMutex
}

// NewAgentInstance creates and initializes a new AgentInstance tracker.
func NewAgentInstance(agentID string) *AgentInstance {
	now := time.Now()
	return &AgentInstance{
		ID:           agentID,
		Status:       StatusIdle,
		Queue:        &AgentQueue{},
		Metrics:      &AgentMetrics{LastActivity: now},
		StartedAt:    now,
		LastActivity: now,
	}
}

// UpdateStatus sets the agent status and refreshes the last-activity timestamp.
func (i *AgentInstance) UpdateStatus(status AgentStatus) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.Status = status
	i.LastActivity = time.Now()
}

// AddActiveTask appends a task to the active queue and updates metrics.
func (i *AgentInstance) AddActiveTask(task TaskInfo) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.Queue.ActiveTasks = append(i.Queue.ActiveTasks, task)
	i.Metrics.ActiveTasks = len(i.Queue.ActiveTasks)
	i.LastActivity = time.Now()
}

// CompleteTask moves a task from active to completed (or increments failure count).
func (i *AgentInstance) CompleteTask(taskID string, success bool) {
	i.mu.Lock()
	defer i.mu.Unlock()

	for idx, task := range i.Queue.ActiveTasks {
		if task.ID != taskID {
			continue
		}

		now := time.Now()
		task.CompletedAt = &now
		if success {
			task.Status = "completed"
		} else {
			task.Status = "failed"
		}

		// Remove from active slice
		i.Queue.ActiveTasks = append(
			i.Queue.ActiveTasks[:idx],
			i.Queue.ActiveTasks[idx+1:]...,
		)

		if success {
			i.Queue.Completed = append(i.Queue.Completed, task)
			i.Metrics.CompletedToday++
		} else {
			i.Metrics.FailedToday++
		}

		i.Metrics.ActiveTasks = len(i.Queue.ActiveTasks)
		i.LastActivity = now
		return
	}
}

// GetStatus returns the current agent status.
func (i *AgentInstance) GetStatus() AgentStatus {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.Status
}

// GetMetrics returns a snapshot copy of the current metrics.
func (i *AgentInstance) GetMetrics() AgentMetrics {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return *i.Metrics
}

// GetQueue returns a snapshot copy of the current queue state.
func (i *AgentInstance) GetQueue() AgentQueue {
	i.mu.RLock()
	defer i.mu.RUnlock()
	q := AgentQueue{
		ActiveTasks: make([]TaskInfo, len(i.Queue.ActiveTasks)),
		QueuedTasks: make([]TaskInfo, len(i.Queue.QueuedTasks)),
		Completed:   make([]TaskInfo, len(i.Queue.Completed)),
	}
	copy(q.ActiveTasks, i.Queue.ActiveTasks)
	copy(q.QueuedTasks, i.Queue.QueuedTasks)
	copy(q.Completed, i.Queue.Completed)
	return q
}

// InstanceRegistry tracks all active AgentInstances.
// All methods are safe for concurrent use.
type InstanceRegistry struct {
	instances map[string]*AgentInstance
	mu        sync.RWMutex
}

// NewInstanceRegistry creates a new, empty InstanceRegistry.
func NewInstanceRegistry() *InstanceRegistry {
	return &InstanceRegistry{
		instances: make(map[string]*AgentInstance),
	}
}

// Register creates a new AgentInstance for agentID if one does not already exist,
// or returns the existing one.
func (r *InstanceRegistry) Register(agentID string) *AgentInstance {
	r.mu.Lock()
	defer r.mu.Unlock()

	if instance, ok := r.instances[agentID]; ok {
		return instance
	}

	instance := NewAgentInstance(agentID)
	r.instances[agentID] = instance
	return instance
}

// Get returns the AgentInstance for agentID, if present.
func (r *InstanceRegistry) Get(agentID string) (*AgentInstance, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	instance, ok := r.instances[agentID]
	return instance, ok
}

// Unregister removes the AgentInstance for agentID.
func (r *InstanceRegistry) Unregister(agentID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.instances, agentID)
}

// List returns all active AgentInstances (order is non-deterministic).
func (r *InstanceRegistry) List() []*AgentInstance {
	r.mu.RLock()
	defer r.mu.RUnlock()

	instances := make([]*AgentInstance, 0, len(r.instances))
	for _, instance := range r.instances {
		instances = append(instances, instance)
	}
	return instances
}

// Count returns the number of registered agent instances.
func (r *InstanceRegistry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.instances)
}
