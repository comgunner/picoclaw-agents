// Package mailbox provides message queue infrastructure for agents
//
// This implementation is based on @icueth's picoclaw-agents:
// https://github.com/icueth/picoclaw-agents/tree/main/pkg/mailbox
//
// Credits: @icueth (https://github.com/icueth)
// License: Same as base project (MIT)

package mailbox

import (
	"testing"
	"time"
)

func TestMailbox_Send_Receive(t *testing.T) {
	mb := NewMailbox("test_agent", 100)

	// Send message
	msg := Message{
		ID:        "msg-1",
		Type:      MessageTypeTask,
		From:      "pm",
		To:        "dev",
		Priority:  PriorityHigh,
		Content:   "Review PR",
		CreatedAt: time.Now(),
	}

	if err := mb.Send(msg); err != nil {
		t.Fatalf("Send failed: %v", err)
	}

	// Receive message
	received, err := mb.Receive()
	if err != nil {
		t.Fatalf("Receive failed: %v", err)
	}

	if received.Content != "Review PR" {
		t.Errorf("Expected 'Review PR', got %s", received.Content)
	}

	// Should be marked as read
	if !received.Read {
		t.Errorf("Message should be marked as read")
	}

	// Mailbox should be empty now
	if mb.Size() != 0 {
		t.Errorf("Expected 0 messages after receive, got %d", mb.Size())
	}
}

func TestMailbox_Priority(t *testing.T) {
	mb := NewMailbox("test_agent", 100)

	// Send low priority first
	msg1 := Message{
		ID:       "msg-1",
		Type:     MessageTypeStatus,
		Priority: PriorityLow,
		Content:  "Low priority",
	}
	mb.Send(msg1)

	// Send critical priority
	msg2 := Message{
		ID:       "msg-2",
		Type:     MessageTypeStatus,
		Priority: PriorityCritical,
		Content:  "Critical",
	}
	mb.Send(msg2)

	// Receive should get critical first
	first, _ := mb.Receive()
	if first.Priority != PriorityCritical {
		t.Errorf("Expected critical priority first, got %d", first.Priority)
	}

	// Second should be low
	second, _ := mb.Receive()
	if second.Priority != PriorityLow {
		t.Errorf("Expected low priority second, got %d", second.Priority)
	}
}

func TestMailbox_Unread_Count(t *testing.T) {
	mb := NewMailbox("test_agent", 100)

	// Send 3 messages
	for i := 1; i <= 3; i++ {
		msg := Message{
			ID:   "msg-" + string(rune(i)),
			Type: MessageTypeStatus,
		}
		mb.Send(msg)
	}

	if mb.GetUnreadCount() != 3 {
		t.Errorf("Expected 3 unread, got %d", mb.GetUnreadCount())
	}

	// Receive one
	mb.Receive()

	if mb.GetUnreadCount() != 2 {
		t.Errorf("Expected 2 unread after receive, got %d", mb.GetUnreadCount())
	}
}

func TestMailbox_Subscribe(t *testing.T) {
	mb := NewMailbox("test_agent", 100)

	// Subscribe to messages
	msgChan := mb.Subscribe()

	// Send message
	msg := Message{
		ID:      "msg-1",
		Type:    MessageTypeTask,
		Content: "New task",
	}

	go func() {
		mb.Send(msg)
	}()

	// Should receive notification
	select {
	case received := <-msgChan:
		if received.Content != "New task" {
			t.Errorf("Expected 'New task', got %s", received.Content)
		}
	case <-time.After(1 * time.Second):
		t.Errorf("Timeout waiting for message notification")
	}
}

func TestMailbox_Peek(t *testing.T) {
	mb := NewMailbox("test_agent", 100)

	// Send message
	msg := Message{
		ID:      "msg-1",
		Type:    MessageTypeTask,
		Content: "Peek test",
	}
	mb.Send(msg)

	// Peek should not remove
	peeked, err := mb.Peek()
	if err != nil {
		t.Fatalf("Peek failed: %v", err)
	}

	if peeked.Content != "Peek test" {
		t.Errorf("Expected 'Peek test', got %s", peeked.Content)
	}

	// Size should still be 1
	if mb.Size() != 1 {
		t.Errorf("Expected size 1 after peek, got %d", mb.Size())
	}
}

func TestMailbox_Clear(t *testing.T) {
	mb := NewMailbox("test_agent", 100)

	// Send messages
	for i := 0; i < 5; i++ {
		mb.Send(Message{ID: "msg-" + string(rune(i))})
	}

	// Clear
	mb.Clear()

	if mb.Size() != 0 {
		t.Errorf("Expected 0 messages after clear, got %d", mb.Size())
	}
}

func TestMailbox_Capacity(t *testing.T) {
	mb := NewMailbox("test_agent", 5)

	// Send 7 messages (exceeds capacity)
	for i := 0; i < 7; i++ {
		mb.Send(Message{ID: "msg-" + string(rune(i))})
	}

	// Should only keep last 5
	if mb.Size() != 5 {
		t.Errorf("Expected 5 messages (capacity), got %d", mb.Size())
	}
}

func TestMailbox_CleanupExpired(t *testing.T) {
	mb := NewMailbox("test_agent", 100)

	// Send message that expires immediately
	expired := Message{
		ID:        "msg-1",
		ExpiresAt: time.Now().Add(-1 * time.Second),
	}
	mb.Send(expired)

	// Send message that doesn't expire
	active := Message{
		ID:        "msg-2",
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}
	mb.Send(active)

	// Cleanup
	removed := mb.CleanupExpired()

	if removed != 1 {
		t.Errorf("Expected 1 removed, got %d", removed)
	}

	if mb.Size() != 1 {
		t.Errorf("Expected 1 message remaining, got %d", mb.Size())
	}
}

func TestMailbox_GetMessagesByType(t *testing.T) {
	mb := NewMailbox("test_agent", 100)

	// Send different types
	mb.Send(Message{ID: "msg-1", Type: MessageTypeTask})
	mb.Send(Message{ID: "msg-2", Type: MessageTypeStatus})
	mb.Send(Message{ID: "msg-3", Type: MessageTypeTask})

	// Get task messages
	tasks := mb.GetMessagesByType(MessageTypeTask)
	if len(tasks) != 2 {
		t.Errorf("Expected 2 task messages, got %d", len(tasks))
	}
}

func TestMailbox_GetMessagesByPriority(t *testing.T) {
	mb := NewMailbox("test_agent", 100)

	// Send different priorities
	mb.Send(Message{ID: "msg-1", Priority: PriorityHigh})
	mb.Send(Message{ID: "msg-2", Priority: PriorityNormal})
	mb.Send(Message{ID: "msg-3", Priority: PriorityHigh})

	// Get high priority messages
	high := mb.GetMessagesByPriority(PriorityHigh)
	if len(high) != 2 {
		t.Errorf("Expected 2 high priority messages, got %d", len(high))
	}
}

func TestMailbox_SetExpiration(t *testing.T) {
	mb := NewMailbox("test_agent", 100)

	// Send message
	mb.Send(Message{ID: "msg-1"})

	// Set expiration
	expiresAt := time.Now().Add(1 * time.Hour)
	if err := mb.SetExpiration("msg-1", expiresAt); err != nil {
		t.Fatalf("SetExpiration failed: %v", err)
	}

	// Verify
	msgs := mb.GetMessagesByType("")
	found := false
	for _, msg := range msgs {
		if msg.ID == "msg-1" {
			if msg.ExpiresAt.IsZero() || msg.ExpiresAt.Before(expiresAt) {
				t.Errorf("Expiration not set correctly")
			}
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Message not found")
	}
}

func TestMailbox_MarkAsRead(t *testing.T) {
	mb := NewMailbox("test_agent", 100)

	// Send message
	mb.Send(Message{ID: "msg-1"})

	// Mark as read
	if err := mb.MarkAsRead("msg-1"); err != nil {
		t.Fatalf("MarkAsRead failed: %v", err)
	}

	// Unread count should be 0
	if mb.GetUnreadCount() != 0 {
		t.Errorf("Expected 0 unread after marking as read")
	}
}

func TestMailbox_MarkAsRead_NotFound(t *testing.T) {
	mb := NewMailbox("test_agent", 100)

	// Try to mark non-existent message
	if err := mb.MarkAsRead("nonexistent"); err == nil {
		t.Errorf("Expected error for non-existent message")
	}
}

func TestMailbox_SetExpiration_NotFound(t *testing.T) {
	mb := NewMailbox("test_agent", 100)

	// Try to set expiration for non-existent message
	if err := mb.SetExpiration("nonexistent", time.Now()); err == nil {
		t.Errorf("Expected error for non-existent message")
	}
}

func TestMailbox_GetAgentID(t *testing.T) {
	mb := NewMailbox("test_agent_123", 100)

	if mb.GetAgentID() != "test_agent_123" {
		t.Errorf("Expected 'test_agent_123', got %s", mb.GetAgentID())
	}
}

func TestMailbox_Unsubscribe(t *testing.T) {
	mb := NewMailbox("test_agent", 100)

	// Subscribe
	ch := mb.Subscribe()

	// Unsubscribe
	mb.Unsubscribe(ch)

	// Send message (should not panic)
	mb.Send(Message{ID: "msg-1", Content: "test"})

	// Channel should be closed
	select {
	case _, ok := <-ch:
		if ok {
			t.Errorf("Channel should be closed")
		}
	default:
		t.Errorf("Channel should be closed")
	}
}

func TestMailbox_ConcurrentSend(t *testing.T) {
	mb := NewMailbox("test_agent", 1000)

	// Concurrent sends
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(idx int) {
			for j := 0; j < 50; j++ {
				mb.Send(Message{ID: "msg-" + string(rune(idx*50+j))})
			}
			done <- true
		}(i)
	}

	// Wait for all
	for i := 0; i < 10; i++ {
		<-done
	}

	// Should have all messages (or capacity limit)
	expected := 500
	if mb.Size() > expected {
		t.Errorf("Expected max %d messages, got %d", expected, mb.Size())
	}
}
