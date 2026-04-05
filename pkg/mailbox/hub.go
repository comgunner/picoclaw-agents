// Package mailbox provides message queue infrastructure for agents
//
// This implementation is based on @icueth's picoclaw-agents:
// https://github.com/icueth/picoclaw-agents/tree/main/pkg/mailbox
//
// Credits: @icueth (https://github.com/icueth)
// License: Same as base project (MIT)

package mailbox

import (
	"context"
	"sync"
	"time"

	"github.com/comgunner/picoclaw/pkg/logger"
)

// Hub manages all agent mailboxes globally
type Hub struct {
	mailboxes map[string]*Mailbox
	mu        sync.RWMutex
	capacity  int
}

// NewHub creates a new mailbox hub for agent coordination
func NewHub(mailboxCapacity int) *Hub {
	return &Hub{
		mailboxes: make(map[string]*Mailbox),
		capacity:  mailboxCapacity,
	}
}

// Register creates or returns mailbox for an agent
func (h *Hub) Register(agentID string) *Mailbox {
	h.mu.Lock()
	defer h.mu.Unlock()

	if mb, exists := h.mailboxes[agentID]; exists {
		return mb
	}

	mb := NewMailbox(agentID, h.capacity)
	h.mailboxes[agentID] = mb
	logger.InfoCF("mailbox", "Registered agent mailbox", map[string]any{"agent_id": agentID})
	return mb
}

// Get retrieves mailbox for an agent
func (h *Hub) Get(agentID string) (*Mailbox, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	mb, ok := h.mailboxes[agentID]
	return mb, ok
}

// SendTo sends message to specific agent
func (h *Hub) SendTo(agentID string, msg Message) error {
	mb, ok := h.Get(agentID)
	if !ok {
		return ErrAgentNotFound(agentID)
	}
	return mb.Send(msg)
}

// Broadcast sends message to all agents except sender
func (h *Hub) Broadcast(msg Message) (failed []string) {
	h.mu.RLock()
	mailboxes := make(map[string]*Mailbox)
	for id, mb := range h.mailboxes {
		mailboxes[id] = mb
	}
	h.mu.RUnlock()

	for id, mb := range mailboxes {
		if id == msg.From {
			continue
		}
		if err := mb.Send(msg); err != nil {
			failed = append(failed, id)
			logger.WarnCF("mailbox", "Broadcast failed", map[string]any{
				"agent_id": id,
				"error":    err.Error(),
			})
		}
	}
	return failed
}

// Unregister removes a mailbox for an agent
func (h *Hub) Unregister(agentID string) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, exists := h.mailboxes[agentID]; !exists {
		return ErrAgentNotFound(agentID)
	}

	// Close all subscribers safely via mailbox's own lock (B-05 fix)
	mb := h.mailboxes[agentID]
	mb.CloseSubscribers()

	delete(h.mailboxes, agentID)
	logger.InfoCF("mailbox", "Unregistered agent mailbox", map[string]any{"agent_id": agentID})
	return nil
}

// GetStats returns hub statistics
func (h *Hub) GetStats() map[string]any {
	h.mu.RLock()
	defer h.mu.RUnlock()

	stats := map[string]any{
		"total_agents": len(h.mailboxes),
		"agents":       make(map[string]any),
	}

	agentStats := stats["agents"].(map[string]any)
	for id, mb := range h.mailboxes {
		agentStats[id] = map[string]any{
			"size":         mb.Size(),
			"unread_count": mb.GetUnreadCount(),
		}
	}

	return stats
}

// ListAgents returns all registered agent IDs
func (h *Hub) ListAgents() []string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	agents := make([]string, 0, len(h.mailboxes))
	for id := range h.mailboxes {
		agents = append(agents, id)
	}
	return agents
}

// StartCleanup starts periodic cleanup of expired messages
func (h *Hub) StartCleanup(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				h.mu.RLock()
				mailboxes := make([]*Mailbox, 0, len(h.mailboxes))
				for _, mb := range h.mailboxes {
					mailboxes = append(mailboxes, mb)
				}
				h.mu.RUnlock()

				for _, mb := range mailboxes {
					removed := mb.CleanupExpired()
					if removed > 0 {
						logger.DebugCF("mailbox", "Cleaned expired messages", map[string]any{
							"mailbox": mb.agentID,
							"removed": removed,
						})
					}
				}
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
}

// StopCleanup stops the cleanup goroutine (called via context cancellation)

// Error types

// MailboxError represents a mailbox-related error
type MailboxError struct {
	Code    string
	Message string
}

func (e *MailboxError) Error() string {
	return e.Code + ": " + e.Message
}

// ErrAgentNotFound returns an error when agent is not registered
func ErrAgentNotFound(agentID string) error {
	return &MailboxError{
		Code:    "AGENT_NOT_FOUND",
		Message: "Agent not registered: " + agentID,
	}
}

// ErrMailboxFull returns an error when mailbox is at capacity
func ErrMailboxFull(agentID string) error {
	return &MailboxError{
		Code:    "MAILBOX_FULL",
		Message: "Mailbox for agent " + agentID + " is at capacity",
	}
}
