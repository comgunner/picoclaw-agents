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
	"testing"
	"time"
)

func TestHub_Register_Get(t *testing.T) {
	hub := NewHub(100)

	// Register agent
	mb := hub.Register("agent1")
	if mb == nil {
		t.Fatalf("Failed to register agent")
	}

	// Get registered agent
	retrieved, ok := hub.Get("agent1")
	if !ok {
		t.Fatalf("Agent not found after registration")
	}

	if retrieved != mb {
		t.Errorf("Retrieved mailbox differs from registered")
	}

	// Get non-existent agent
	_, ok = hub.Get("nonexistent")
	if ok {
		t.Errorf("Should not find non-existent agent")
	}
}

func TestHub_Register_Duplicate(t *testing.T) {
	hub := NewHub(100)

	// Register same agent twice
	mb1 := hub.Register("agent1")
	mb2 := hub.Register("agent1")

	if mb1 != mb2 {
		t.Errorf("Should return same mailbox for duplicate registration")
	}
}

func TestHub_SendTo(t *testing.T) {
	hub := NewHub(100)
	hub.Register("agent1")
	hub.Register("agent2")

	msg := Message{
		ID:      "msg-1",
		Type:    MessageTypeTask,
		From:    "agent1",
		To:      "agent2",
		Content: "Task for agent2",
	}

	// Send to agent2
	if err := hub.SendTo("agent2", msg); err != nil {
		t.Fatalf("SendTo failed: %v", err)
	}

	// Verify agent2 received
	mb, _ := hub.Get("agent2")
	if mb.GetUnreadCount() != 1 {
		t.Errorf("Expected 1 unread message")
	}
}

func TestHub_SendTo_NotFound(t *testing.T) {
	hub := NewHub(100)

	msg := Message{
		ID:   "msg-1",
		Type: MessageTypeTask,
	}

	// Send to non-existent agent
	if err := hub.SendTo("nonexistent", msg); err == nil {
		t.Errorf("Expected error for non-existent agent")
	} else {
		if mailboxErr, ok := err.(*MailboxError); !ok || mailboxErr.Code != "AGENT_NOT_FOUND" {
			t.Errorf("Expected AGENT_NOT_FOUND error")
		}
	}
}

func TestHub_Broadcast(t *testing.T) {
	hub := NewHub(100)
	hub.Register("pm")
	hub.Register("dev")
	hub.Register("qa")

	msg := Message{
		ID:      "msg-1",
		Type:    MessageTypeBroadcast,
		From:    "pm",
		Content: "All hands on deck",
	}

	// Broadcast
	failed := hub.Broadcast(msg)

	if len(failed) != 0 {
		t.Errorf("Expected no failures, got %d", len(failed))
	}

	// Verify all except sender received
	devMb, _ := hub.Get("dev")
	qaMb, _ := hub.Get("qa")

	if devMb.GetUnreadCount() != 1 {
		t.Errorf("Dev should receive broadcast")
	}
	if qaMb.GetUnreadCount() != 1 {
		t.Errorf("QA should receive broadcast")
	}

	// Sender should not receive
	pmMb, _ := hub.Get("pm")
	if pmMb.GetUnreadCount() != 0 {
		t.Errorf("Sender should not receive broadcast")
	}
}

func TestHub_Broadcast_WithFailures(t *testing.T) {
	hub := NewHub(100)
	hub.Register("agent1")

	// Fill agent1 mailbox to capacity
	for i := 0; i < 100; i++ {
		hub.SendTo("agent1", Message{ID: "msg-" + string(rune(i))})
	}

	// Register agent2 and broadcast from agent1
	hub.Register("agent2")
	msg := Message{
		ID:   "broadcast-1",
		Type: MessageTypeBroadcast,
		From: "agent1",
	}

	failed := hub.Broadcast(msg)

	// agent2 should receive, agent1 is sender
	if len(failed) != 0 {
		t.Errorf("Expected no failures, got %v", failed)
	}
}

func TestHub_Stats(t *testing.T) {
	hub := NewHub(100)
	hub.Register("agent1")
	hub.Register("agent2")

	// Send some messages
	msg := Message{
		ID:      "msg-1",
		Type:    MessageTypeTask,
		Content: "test",
	}
	hub.SendTo("agent1", msg)
	hub.SendTo("agent2", msg)

	// Get stats
	stats := hub.GetStats()

	if stats["total_agents"].(int) != 2 {
		t.Errorf("Expected 2 agents in stats")
	}

	agents := stats["agents"].(map[string]any)
	if len(agents) != 2 {
		t.Errorf("Expected 2 agents in details")
	}

	agent1Stats := agents["agent1"].(map[string]any)
	if agent1Stats["size"].(int) != 1 {
		t.Errorf("Expected agent1 size to be 1")
	}
}

func TestHub_ListAgents(t *testing.T) {
	hub := NewHub(100)
	hub.Register("agent1")
	hub.Register("agent2")
	hub.Register("agent3")

	agents := hub.ListAgents()

	if len(agents) != 3 {
		t.Errorf("Expected 3 agents, got %d", len(agents))
	}

	// Check all agents are listed
	expected := map[string]bool{
		"agent1": true,
		"agent2": true,
		"agent3": true,
	}

	for _, agent := range agents {
		if !expected[agent] {
			t.Errorf("Unexpected agent: %s", agent)
		}
	}
}

func TestHub_Unregister(t *testing.T) {
	hub := NewHub(100)
	hub.Register("agent1")

	// Unregister
	if err := hub.Unregister("agent1"); err != nil {
		t.Fatalf("Unregister failed: %v", err)
	}

	// Should not find agent
	_, ok := hub.Get("agent1")
	if ok {
		t.Errorf("Agent should be unregistered")
	}
}

func TestHub_Unregister_NotFound(t *testing.T) {
	hub := NewHub(100)

	// Unregister non-existent agent
	if err := hub.Unregister("nonexistent"); err == nil {
		t.Errorf("Expected error for non-existent agent")
	}
}

func TestHub_Cleanup(t *testing.T) {
	hub := NewHub(100)
	hub.Register("agent1")

	// Send message that expires immediately
	msg := Message{
		ID:        "msg-1",
		ExpiresAt: time.Now().Add(-1 * time.Second),
	}
	hub.SendTo("agent1", msg)

	// Create context for cleanup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start cleanup with short interval
	hub.StartCleanup(ctx, 100*time.Millisecond)

	// Wait for cleanup to run
	time.Sleep(200 * time.Millisecond)

	// Check message was cleaned
	mb, _ := hub.Get("agent1")
	if mb.Size() != 0 {
		t.Errorf("Expected 0 messages after cleanup, got %d", mb.Size())
	}
}

func TestHub_CleanupCancellation(t *testing.T) {
	hub := NewHub(100)
	hub.Register("agent1")

	// Create context for cleanup
	ctx, cancel := context.WithCancel(context.Background())

	// Start cleanup
	hub.StartCleanup(ctx, 50*time.Millisecond)

	// Cancel immediately
	cancel()

	// Wait a bit
	time.Sleep(100 * time.Millisecond)

	// Should not panic or cause issues
	hub.SendTo("agent1", Message{ID: "msg-1"})

	mb, _ := hub.Get("agent1")
	if mb.Size() != 1 {
		t.Errorf("Expected 1 message after cleanup cancellation")
	}
}

func TestHub_ConcurrentRegister(t *testing.T) {
	hub := NewHub(100)

	// Concurrent registrations
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(idx int) {
			for j := 0; j < 10; j++ {
				agentID := "agent-" + string(rune(idx*10+j))
				hub.Register(agentID)
			}
			done <- true
		}(i)
	}

	// Wait for all
	for i := 0; i < 10; i++ {
		<-done
	}

	// Should have 100 agents
	agents := hub.ListAgents()
	if len(agents) != 100 {
		t.Errorf("Expected 100 agents, got %d", len(agents))
	}
}

func TestHub_ConcurrentSend(t *testing.T) {
	hub := NewHub(1000)
	hub.Register("agent1")

	// Concurrent sends to same agent
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(idx int) {
			for j := 0; j < 50; j++ {
				hub.SendTo("agent1", Message{ID: "msg-" + string(rune(idx*50+j))})
			}
			done <- true
		}(i)
	}

	// Wait for all
	for i := 0; i < 10; i++ {
		<-done
	}

	// Should have 500 messages (or capacity limit)
	mb, _ := hub.Get("agent1")
	if mb.Size() > 500 {
		t.Errorf("Expected max 500 messages, got %d", mb.Size())
	}
}

func TestHub_Broadcast_Empty(t *testing.T) {
	hub := NewHub(100)

	msg := Message{
		ID:      "msg-1",
		Type:    MessageTypeBroadcast,
		From:    "sender",
		Content: "test",
	}

	// Broadcast with no registered agents
	failed := hub.Broadcast(msg)

	if len(failed) != 0 {
		t.Errorf("Expected no failures for empty broadcast")
	}
}

// --- Tests agregados en QA A2A (Fase 4) ---

func TestHub_SendTo_VerifyDelivery(t *testing.T) {
	hub := NewHub(100)
	hub.Register("sender")
	hub.Register("receiver")

	msg := Message{
		Type:    MessageTypeTask,
		From:    "sender",
		To:      "receiver",
		Content: "test message",
	}

	if err := hub.SendTo("receiver", msg); err != nil {
		t.Fatalf("SendTo failed: %v", err)
	}

	// El REMITENTE no debe tener mensajes
	senderMb, _ := hub.Get("sender")
	if senderMb.Size() != 0 {
		t.Errorf("Sender should have 0 messages, not %d", senderMb.Size())
	}

	// El RECEPTOR debe tener el mensaje
	receiverMb, _ := hub.Get("receiver")
	if receiverMb.Size() != 1 {
		t.Errorf("Receiver should have 1 message, not %d", receiverMb.Size())
	}
}

func TestHub_Unregister_WithActiveSubscribers(t *testing.T) {
	hub := NewHub(100)
	hub.Register("agent1")

	mb, _ := hub.Get("agent1")
	ch := mb.Subscribe()

	// Unregister debe cerrar el canal del suscriptor sin panic
	if err := hub.Unregister("agent1"); err != nil {
		t.Fatalf("Unregister failed: %v", err)
	}

	// Canal debe estar cerrado
	select {
	case _, ok := <-ch:
		if ok {
			t.Errorf("Expected closed channel after unregister")
		}
	case <-time.After(100 * time.Millisecond):
		t.Errorf("Channel not closed after unregister")
	}
}

func TestHub_Unregister_ConcurrentSubscribe(t *testing.T) {
	// Verifica que no hay data race entre Unregister y Subscribe concurrentes (B-05)
	hub := NewHub(100)
	hub.Register("agent1")

	done := make(chan struct{})

	// Goroutine que suscribe continuamente
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				mb, ok := hub.Get("agent1")
				if !ok {
					return
				}
				ch := mb.Subscribe()
				// Consumir el canal para no bloquear
				go func() {
					for range ch {
					}
				}()
			}
		}
	}()

	// Dar tiempo para que la goroutine se inicie
	time.Sleep(1 * time.Millisecond)

	// Unregister concurrente — con B-05 corregido no debe haber data race
	hub.Unregister("agent1")
	close(done)
}

// Verifica que el context cancelado detiene el cleanup (regresión)
func TestHub_StartCleanup_StopsOnCancel(t *testing.T) {
	hub := NewHub(100)
	hub.Register("agent1")

	ctx, cancel := context.WithCancel(context.Background())
	hub.StartCleanup(ctx, 50*time.Millisecond)

	// Cancelar y verificar que no hay goroutine leak (no hay assertion directa,
	// pero go test -race detectaría accesos post-cancel)
	cancel()
	time.Sleep(100 * time.Millisecond)
}
