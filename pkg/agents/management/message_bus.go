// PicoClaw - Ultra-lightweight personal AI agent
// Agent Management - Message Bus (inter-agent communication)
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package management

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// AgentMessage represents a directed message sent between two named agents.
type AgentMessage struct {
	ID               string          `json:"id"`
	SenderID         string          `json:"sender_id"`
	RecipientID      string          `json:"recipient_id"`
	MessageType      string          `json:"message_type"` // info | task | result | broadcast
	Payload          json.RawMessage `json:"payload"`
	RequiresResponse bool            `json:"requires_response"`
	SentAt           time.Time       `json:"sent_at"`
	ExpiresAt        *time.Time      `json:"expires_at,omitempty"`
}

// MessageDeliveryStatus records the lifecycle of a sent message.
type MessageDeliveryStatus string

const (
	StatusSent      MessageDeliveryStatus = "sent"
	StatusDelivered MessageDeliveryStatus = "delivered"
	StatusRead      MessageDeliveryStatus = "read"
	StatusResponded MessageDeliveryStatus = "responded"
)

// MessageLog records an audit entry for a single delivered message.
type MessageLog struct {
	MessageID   string                `json:"message_id"`
	SenderID    string                `json:"sender_id"`
	RecipientID string                `json:"recipient_id"`
	Type        string                `json:"type"`
	Status      MessageDeliveryStatus `json:"status"`
	Timestamp   time.Time             `json:"timestamp"`
}

// MessageStats aggregates statistics across the entire message bus lifetime.
type MessageStats struct {
	TotalMessages   int            `json:"total_messages"`
	MessagesToday   int            `json:"messages_today"`
	AvgResponseTime time.Duration  `json:"avg_response_time_ns"`
	BusiestAgent    string         `json:"busiest_agent"`
	MostCommonType  string         `json:"most_common_type"`
	DeliveryRate    float64        `json:"delivery_rate"`
	ResponseRate    float64        `json:"response_rate"`
	ByType          map[string]int `json:"by_type"`
	ByHour          map[int]int    `json:"by_hour"`
}

// AgentMessageBus is a lightweight, in-process message broker for inter-agent communication.
// Channels are created lazily on the first Send to a given recipient.
// All methods are safe for concurrent use.
type AgentMessageBus struct {
	channels map[string]chan AgentMessage
	history  []MessageLog
	stats    *MessageStats
	mu       sync.RWMutex
}

// channelBuffer is the per-agent inbox capacity.
const channelBuffer = 100

// NewAgentMessageBus creates a new, empty AgentMessageBus.
func NewAgentMessageBus() *AgentMessageBus {
	return &AgentMessageBus{
		channels: make(map[string]chan AgentMessage),
		history:  make([]MessageLog, 0),
		stats: &MessageStats{
			ByType: make(map[string]int),
			ByHour: make(map[int]int),
		},
	}
}

// Send enqueues a message into the recipient's inbox channel.
// Returns an error when the recipient's channel is full (non-blocking send).
func (mb *AgentMessageBus) Send(msg AgentMessage) error {
	mb.mu.Lock()
	defer mb.mu.Unlock()

	// Lazy channel creation
	if _, ok := mb.channels[msg.RecipientID]; !ok {
		mb.channels[msg.RecipientID] = make(chan AgentMessage, channelBuffer)
	}

	select {
	case mb.channels[msg.RecipientID] <- msg:
		// Audit log
		mb.history = append(mb.history, MessageLog{
			MessageID:   msg.ID,
			SenderID:    msg.SenderID,
			RecipientID: msg.RecipientID,
			Type:        msg.MessageType,
			Status:      StatusSent,
			Timestamp:   msg.SentAt,
		})

		// Stats
		mb.stats.TotalMessages++
		mb.stats.ByType[msg.MessageType]++
		mb.stats.ByHour[msg.SentAt.Hour()]++

		// Track busiest agent
		mb.updateBusiestAgent()

		return nil

	default:
		return fmt.Errorf("inbox full for agent %q: drop message %s", msg.RecipientID, msg.ID)
	}
}

// GetChannel returns the specific inbound channel for an agent, creating it lazily if it doesn't exist.
// This allows an autonomous runtime to listen directly to the agent's inbox.
func (mb *AgentMessageBus) GetChannel(agentID string) <-chan AgentMessage {
	mb.mu.Lock()
	defer mb.mu.Unlock()

	if ch, ok := mb.channels[agentID]; ok {
		return ch
	}

	ch := make(chan AgentMessage, channelBuffer)
	mb.channels[agentID] = ch
	return ch
}

// Receive drains pending messages from agentID's inbox.
// When types is non-empty only messages whose MessageType matches an entry are returned;
// non-matching messages are re-queued to avoid losing them.
func (mb *AgentMessageBus) Receive(agentID string, types []string) ([]AgentMessage, error) {
	mb.mu.Lock()
	defer mb.mu.Unlock()

	ch, ok := mb.channels[agentID]
	if !ok {
		return nil, nil
	}

	var matched []AgentMessage
	var unmatched []AgentMessage

	// Drain all pending messages
	for {
		select {
		case msg := <-ch:
			if matchesType(msg.MessageType, types) {
				matched = append(matched, msg)
			} else {
				unmatched = append(unmatched, msg)
			}
		default:
			goto done
		}
	}

done:
	// Re-queue unmatched messages (best-effort; drops if channel fills up)
	for _, m := range unmatched {
		select {
		case ch <- m:
		default:
			// channel is full: discard rather than deadlock
		}
	}

	return matched, nil
}

// matchesType returns true when filter is empty (accept all) or msgType appears in filter.
func matchesType(msgType string, filter []string) bool {
	if len(filter) == 0 {
		return true
	}
	for _, t := range filter {
		if t == msgType {
			return true
		}
	}
	return false
}

// MarkAsRead updates the delivery status for messageID in the audit history.
func (mb *AgentMessageBus) MarkAsRead(messageID string) error {
	mb.mu.Lock()
	defer mb.mu.Unlock()

	for i := range mb.history {
		if mb.history[i].MessageID == messageID {
			mb.history[i].Status = StatusRead
			return nil
		}
	}

	return fmt.Errorf("message not found: %s", messageID)
}

// GetLogs returns audit log entries filtered by agentID (sender or recipient) and/or a since
// time cutoff. A limit of 0 or negative returns all matching entries.
func (mb *AgentMessageBus) GetLogs(agentID string, limit int, since *time.Time) []MessageLog {
	mb.mu.RLock()
	defer mb.mu.RUnlock()

	logs := make([]MessageLog, 0)

	for _, entry := range mb.history {
		if agentID != "" && entry.RecipientID != agentID && entry.SenderID != agentID {
			continue
		}
		if since != nil && entry.Timestamp.Before(*since) {
			continue
		}

		logs = append(logs, entry)

		if limit > 0 && len(logs) >= limit {
			break
		}
	}

	return logs
}

// GetStats returns a copy of the current aggregate statistics.
func (mb *AgentMessageBus) GetStats() MessageStats {
	mb.mu.RLock()
	defer mb.mu.RUnlock()

	// Return a deep-ish copy so callers cannot mutate internal maps
	byType := make(map[string]int, len(mb.stats.ByType))
	for k, v := range mb.stats.ByType {
		byType[k] = v
	}
	byHour := make(map[int]int, len(mb.stats.ByHour))
	for k, v := range mb.stats.ByHour {
		byHour[k] = v
	}

	return MessageStats{
		TotalMessages:   mb.stats.TotalMessages,
		MessagesToday:   mb.stats.MessagesToday,
		AvgResponseTime: mb.stats.AvgResponseTime,
		BusiestAgent:    mb.stats.BusiestAgent,
		MostCommonType:  mb.stats.MostCommonType,
		DeliveryRate:    mb.stats.DeliveryRate,
		ResponseRate:    mb.stats.ResponseRate,
		ByType:          byType,
		ByHour:          byHour,
	}
}

// updateBusiestAgent refreshes the BusiestAgent counter.
// Must be called while mb.mu is held for writing.
func (mb *AgentMessageBus) updateBusiestAgent() {
	type entry struct {
		id    string
		count int
	}
	var best entry
	for id, ch := range mb.channels {
		cnt := len(ch)
		if cnt > best.count {
			best = entry{id: id, count: cnt}
		}
	}
	if best.count > 0 {
		mb.stats.BusiestAgent = best.id
	}
}

// InboxSize returns the number of pending messages waiting for agentID.
func (mb *AgentMessageBus) InboxSize(agentID string) int {
	mb.mu.RLock()
	defer mb.mu.RUnlock()
	if ch, ok := mb.channels[agentID]; ok {
		return len(ch)
	}
	return 0
}
