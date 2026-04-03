// Package mailbox provides message queue infrastructure for agents
//
// This implementation is based on @icueth's picoclaw-agents:
// https://github.com/icueth/picoclaw-agents/tree/main/pkg/mailbox
//
// Features: priority queue, message history, subscriptions, cleanup
//
// Credits: @icueth (https://github.com/icueth)
// License: Same as base project (MIT)

package mailbox

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

// MessagePriority defines priority levels for messages
type MessagePriority int

const (
	PriorityCritical MessagePriority = 1
	PriorityHigh     MessagePriority = 2
	PriorityNormal   MessagePriority = 3
	PriorityLow      MessagePriority = 4
)

// MessageType defines the type of message
type MessageType string

const (
	MessageTypeTask      MessageType = "task"
	MessageTypeQuestion  MessageType = "question"
	MessageTypeResponse  MessageType = "response"
	MessageTypeStatus    MessageType = "status"
	MessageTypeBroadcast MessageType = "broadcast"
	MessageTypeCommand   MessageType = "command"
)

// Message represents a message in the mailbox
type Message struct {
	ID        string          `json:"id"`
	Type      MessageType     `json:"type"`
	From      string          `json:"from"`
	To        string          `json:"to"`
	Priority  MessagePriority `json:"priority"`
	Content   string          `json:"content"`
	Read      bool            `json:"read"`
	CreatedAt time.Time       `json:"created_at"`
	ExpiresAt time.Time       `json:"expires_at,omitempty"`
}

// Mailbox represents a message queue for an agent
type Mailbox struct {
	agentID     string
	messages    []Message
	capacity    int
	mu          sync.RWMutex
	subscribers []chan Message
}

// NewMailbox creates a new mailbox for an agent
func NewMailbox(agentID string, capacity int) *Mailbox {
	return &Mailbox{
		agentID:     agentID,
		messages:    make([]Message, 0),
		capacity:    capacity,
		subscribers: make([]chan Message, 0),
	}
}

// Send adds a message to the mailbox
func (mb *Mailbox) Send(msg Message) error {
	mb.mu.Lock()
	defer mb.mu.Unlock()

	// Check capacity
	if len(mb.messages) >= mb.capacity {
		// Remove oldest message
		mb.messages = mb.messages[1:]
	}

	// Set defaults
	if msg.ID == "" {
		msg.ID = fmt.Sprintf("%s-%d", mb.agentID, time.Now().UnixNano())
	}
	if msg.CreatedAt.IsZero() {
		msg.CreatedAt = time.Now()
	}
	if msg.To == "" {
		msg.To = mb.agentID
	}

	mb.messages = append(mb.messages, msg)

	// Sort by priority (lower number = higher priority)
	sort.Slice(mb.messages, func(i, j int) bool {
		if mb.messages[i].Priority != mb.messages[j].Priority {
			return mb.messages[i].Priority < mb.messages[j].Priority
		}
		return mb.messages[i].CreatedAt.Before(mb.messages[j].CreatedAt)
	})

	// Notify subscribers
	for _, ch := range mb.subscribers {
		select {
		case ch <- msg:
		default:
			// Channel full, skip
		}
	}

	return nil
}

// Receive retrieves and marks as read the highest priority message
func (mb *Mailbox) Receive() (Message, error) {
	mb.mu.Lock()
	defer mb.mu.Unlock()

	if len(mb.messages) == 0 {
		return Message{}, fmt.Errorf("no messages")
	}

	// Get first message (highest priority)
	msg := mb.messages[0]
	msg.Read = true
	mb.messages[0] = msg

	// Remove from queue
	mb.messages = mb.messages[1:]

	return msg, nil
}

// Peek retrieves the highest priority message without removing it
func (mb *Mailbox) Peek() (Message, error) {
	mb.mu.RLock()
	defer mb.mu.RUnlock()

	if len(mb.messages) == 0 {
		return Message{}, fmt.Errorf("no messages")
	}

	return mb.messages[0], nil
}

// Subscribe returns a channel that receives message notifications
func (mb *Mailbox) Subscribe() chan Message {
	mb.mu.Lock()
	defer mb.mu.Unlock()

	ch := make(chan Message, 10)
	mb.subscribers = append(mb.subscribers, ch)
	return ch
}

// Unsubscribe removes a subscriber channel
func (mb *Mailbox) Unsubscribe(ch chan Message) {
	mb.mu.Lock()
	defer mb.mu.Unlock()

	for i, sub := range mb.subscribers {
		if sub == ch {
			mb.subscribers = append(mb.subscribers[:i], mb.subscribers[i+1:]...)
			close(ch)
			return
		}
	}
}

// GetUnreadCount returns the number of unread messages
func (mb *Mailbox) GetUnreadCount() int {
	mb.mu.RLock()
	defer mb.mu.RUnlock()

	count := 0
	for _, msg := range mb.messages {
		if !msg.Read {
			count++
		}
	}
	return count
}

// Size returns the total number of messages
func (mb *Mailbox) Size() int {
	mb.mu.RLock()
	defer mb.mu.RUnlock()
	return len(mb.messages)
}

// Clear removes all messages from the mailbox
func (mb *Mailbox) Clear() {
	mb.mu.Lock()
	defer mb.mu.Unlock()
	mb.messages = make([]Message, 0)
}

// CleanupExpired removes expired messages
func (mb *Mailbox) CleanupExpired() int {
	mb.mu.Lock()
	defer mb.mu.Unlock()

	now := time.Now()
	removed := 0

	// Filter out expired messages
	active := make([]Message, 0, len(mb.messages))
	for _, msg := range mb.messages {
		if !msg.ExpiresAt.IsZero() && msg.ExpiresAt.Before(now) {
			removed++
		} else {
			active = append(active, msg)
		}
	}

	mb.messages = active
	return removed
}

// GetMessagesByType returns all messages of a specific type
func (mb *Mailbox) GetMessagesByType(msgType MessageType) []Message {
	mb.mu.RLock()
	defer mb.mu.RUnlock()

	result := make([]Message, 0)
	for _, msg := range mb.messages {
		if msg.Type == msgType {
			result = append(result, msg)
		}
	}
	return result
}

// GetMessagesByPriority returns all messages with a specific priority
func (mb *Mailbox) GetMessagesByPriority(priority MessagePriority) []Message {
	mb.mu.RLock()
	defer mb.mu.RUnlock()

	result := make([]Message, 0)
	for _, msg := range mb.messages {
		if msg.Priority == priority {
			result = append(result, msg)
		}
	}
	return result
}

// SetExpiration sets expiration time for a message by ID
func (mb *Mailbox) SetExpiration(msgID string, expiresAt time.Time) error {
	mb.mu.Lock()
	defer mb.mu.Unlock()

	for i, msg := range mb.messages {
		if msg.ID == msgID {
			mb.messages[i].ExpiresAt = expiresAt
			return nil
		}
	}

	return fmt.Errorf("message not found")
}

// MarkAsRead marks a message as read by ID
func (mb *Mailbox) MarkAsRead(msgID string) error {
	mb.mu.Lock()
	defer mb.mu.Unlock()

	for i, msg := range mb.messages {
		if msg.ID == msgID {
			mb.messages[i].Read = true
			return nil
		}
	}

	return fmt.Errorf("message not found")
}

// GetAgentID returns the agent ID this mailbox belongs to
func (mb *Mailbox) GetAgentID() string {
	return mb.agentID
}

// CloseSubscribers closes all subscriber channels safely (B-05 fix)
func (mb *Mailbox) CloseSubscribers() {
	mb.mu.Lock()
	defer mb.mu.Unlock()
	for _, sub := range mb.subscribers {
		close(sub)
	}
	mb.subscribers = nil
}
