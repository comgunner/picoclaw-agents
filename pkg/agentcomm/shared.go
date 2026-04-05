// Package agentcomm provides inter-agent communication primitives
//
// This implementation is adapted from @icueth's picoclaw-agents fork:
// https://github.com/icueth/picoclaw-agents/tree/main/pkg/agentcomm
//
// Original source maintains the same purpose: thread-safe shared context
// for agent coordination and message logging.
//
// Credits: @icueth (https://github.com/icueth)
// License: Same as base project (MIT)

package agentcomm

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// MessageType defines the type of message being sent between agents
type MessageType string

const (
	MsgRequest       MessageType = "request"
	MsgResponse      MessageType = "response"
	MsgBroadcast     MessageType = "broadcast"
	MsgContextUpdate MessageType = "context_update"
	MsgHeartbeat     MessageType = "heartbeat"
	MsgTerminate     MessageType = "terminate"
)

// AgentStatus represents the current state of an agent
type AgentStatus string

const (
	StatusIdle      AgentStatus = "idle"
	StatusRunning   AgentStatus = "running"
	StatusWaiting   AgentStatus = "waiting"
	StatusCompleted AgentStatus = "completed"
	StatusFailed    AgentStatus = "failed"
)

// AgentMessage represents a message sent between agents
type AgentMessage struct {
	ID        string      `json:"id"`
	Type      MessageType `json:"type"`
	From      string      `json:"from"`
	To        string      `json:"to"`
	Content   string      `json:"content"`
	SessionID string      `json:"session_id"`
	ReplyTo   string      `json:"reply_to,omitempty"`
	Timestamp int64       `json:"timestamp"` // milliseconds since epoch
}

// NewAgentMessage creates a new agent message with timestamp
func NewAgentMessage(from, to string, msgType MessageType, content, sessionID string) AgentMessage {
	return AgentMessage{
		ID:        generateID(),
		Type:      msgType,
		From:      from,
		To:        to,
		Content:   content,
		SessionID: sessionID,
		Timestamp: time.Now().UnixNano() / 1e6,
	}
}

// SharedContext provides a thread-safe key-value store for agent coordination
type SharedContext struct {
	context    map[string]any
	msgLog     []MessageLogEntry
	maxLogSize int
	maxContext int
	mu         sync.RWMutex
}

// MessageLogEntry represents a logged message between agents
type MessageLogEntry struct {
	From      string `json:"from"`
	To        string `json:"to"`
	Type      string `json:"type"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
}

// NewSharedContext creates a new shared context with limits
func NewSharedContext(maxLogSize, maxContext int) *SharedContext {
	return &SharedContext{
		context:    make(map[string]any),
		msgLog:     make([]MessageLogEntry, 0),
		maxLogSize: maxLogSize,
		maxContext: maxContext,
	}
}

// Set stores a value in the shared context
func (sc *SharedContext) Set(key string, value any) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	if len(sc.context) >= sc.maxContext {
		// Remove oldest entry (first key)
		for k := range sc.context {
			delete(sc.context, k)
			break
		}
	}
	sc.context[key] = value
}

// Get retrieves a value from the shared context
func (sc *SharedContext) Get(key string) (any, bool) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	val, ok := sc.context[key]
	return val, ok
}

// GetString retrieves a string value from the shared context
func (sc *SharedContext) GetString(key string) (string, bool) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	val, ok := sc.context[key]
	if !ok {
		return "", false
	}
	str, ok := val.(string)
	return str, ok
}

// Delete removes a value from the shared context
func (sc *SharedContext) Delete(key string) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	delete(sc.context, key)
}

// Clear removes all values from the shared context
func (sc *SharedContext) Clear() {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.context = make(map[string]any)
}

// Keys returns all keys in the shared context
func (sc *SharedContext) Keys() []string {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	keys := make([]string, 0, len(sc.context))
	for k := range sc.context {
		keys = append(keys, k)
	}
	return keys
}

// AddMessageLog adds a message to the log
func (sc *SharedContext) AddMessageLog(from, to, msgType, content string) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	entry := MessageLogEntry{
		From:      from,
		To:        to,
		Type:      msgType,
		Content:   content,
		Timestamp: time.Now().UnixNano() / 1e6,
	}

	sc.msgLog = append(sc.msgLog, entry)

	// Trim if exceeds max size
	if len(sc.msgLog) > sc.maxLogSize {
		sc.msgLog = sc.msgLog[len(sc.msgLog)-sc.maxLogSize:]
	}
}

// GetMessageLog returns the message log
func (sc *SharedContext) GetMessageLog() []MessageLogEntry {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	// Return a copy to prevent race conditions
	logCopy := make([]MessageLogEntry, len(sc.msgLog))
	copy(logCopy, sc.msgLog)
	return logCopy
}

// GetMessageLogSince returns messages logged after a specific timestamp
func (sc *SharedContext) GetMessageLogSince(timestamp int64) []MessageLogEntry {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	result := make([]MessageLogEntry, 0)
	for _, entry := range sc.msgLog {
		if entry.Timestamp >= timestamp {
			result = append(result, entry)
		}
	}
	return result
}

// GetMessagesForAgent returns all messages where the agent is sender or receiver
func (sc *SharedContext) GetMessagesForAgent(agentID string) []MessageLogEntry {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	result := make([]MessageLogEntry, 0)
	for _, entry := range sc.msgLog {
		if entry.From == agentID || entry.To == agentID {
			result = append(result, entry)
		}
	}
	return result
}

// ContextSize returns the number of items in the shared context
func (sc *SharedContext) ContextSize() int {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	return len(sc.context)
}

// LogSize returns the number of messages in the log
func (sc *SharedContext) LogSize() int {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	return len(sc.msgLog)
}

// idCounter ensures unique IDs under concurrent generation (B-06 fix)
var idCounter atomic.Uint64

// generateID creates a unique message ID (timestamp + atomic counter)
func generateID() string {
	count := idCounter.Add(1)
	return fmt.Sprintf("%s-%d", time.Now().Format("20060102150405.000000000"), count)
}
