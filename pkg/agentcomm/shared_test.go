// Package agentcomm provides inter-agent communication primitives
//
// This implementation is adapted from @icueth's picoclaw-agents fork:
// https://github.com/icueth/picoclaw-agents/tree/main/pkg/agentcomm
//
// Credits: @icueth (https://github.com/icueth)
// License: Same as base project (MIT)

package agentcomm

import (
	"sync"
	"testing"
	"time"
)

func TestSharedContext_SetGet(t *testing.T) {
	sc := NewSharedContext(100, 1000)

	// Set a value
	sc.Set("test_key", "test_value")

	// Get it back
	val, ok := sc.Get("test_key")
	if !ok || val != "test_value" {
		t.Errorf("Expected 'test_value', got %v", val)
	}
}

func TestSharedContext_GetString(t *testing.T) {
	sc := NewSharedContext(100, 1000)

	// Set a string value
	sc.Set("string_key", "hello world")

	// Get it back
	val, ok := sc.GetString("string_key")
	if !ok || val != "hello world" {
		t.Errorf("Expected 'hello world', got %v", val)
	}

	// Try to get non-existent key
	_, ok = sc.GetString("nonexistent")
	if ok {
		t.Errorf("Expected false for non-existent key")
	}
}

func TestSharedContext_Delete(t *testing.T) {
	sc := NewSharedContext(100, 1000)

	// Set and delete
	sc.Set("to_delete", "value")
	sc.Delete("to_delete")

	_, ok := sc.Get("to_delete")
	if ok {
		t.Errorf("Expected key to be deleted")
	}
}

func TestSharedContext_Clear(t *testing.T) {
	sc := NewSharedContext(100, 1000)

	// Set multiple values
	sc.Set("key1", "value1")
	sc.Set("key2", "value2")
	sc.Set("key3", "value3")

	// Clear all
	sc.Clear()

	if sc.ContextSize() != 0 {
		t.Errorf("Expected empty context after clear")
	}
}

func TestSharedContext_MessageLog(t *testing.T) {
	sc := NewSharedContext(100, 1000)

	// Add message to log
	sc.AddMessageLog("agent1", "agent2", "request", "do something")

	// Retrieve log
	logs := sc.GetMessageLog()
	if len(logs) != 1 {
		t.Errorf("Expected 1 log entry, got %d", len(logs))
	}

	if logs[0].From != "agent1" || logs[0].To != "agent2" {
		t.Errorf("Message log entry incorrect")
	}
}

func TestSharedContext_GetMessageLogSince(t *testing.T) {
	sc := NewSharedContext(100, 1000)

	// Get timestamp before adding
	before := time.Now().UnixNano() / 1e6

	// Add messages
	time.Sleep(10 * time.Millisecond)
	sc.AddMessageLog("pm", "dev", "task", "review code")
	sc.AddMessageLog("dev", "qa", "status", "done")

	// Get logs since before
	logs := sc.GetMessageLogSince(before)
	if len(logs) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(logs))
	}
}

func TestSharedContext_GetMessagesForAgent(t *testing.T) {
	sc := NewSharedContext(100, 1000)

	// Add messages
	sc.AddMessageLog("pm", "dev", "task", "review code")
	sc.AddMessageLog("dev", "qa", "status", "done")
	sc.AddMessageLog("qa", "pm", "response", "approved")

	// Get messages for "dev"
	devMsgs := sc.GetMessagesForAgent("dev")
	if len(devMsgs) != 2 { // dev is sender and receiver
		t.Errorf("Expected 2 messages for dev, got %d", len(devMsgs))
	}
}

func TestSharedContext_MaxLogSize(t *testing.T) {
	sc := NewSharedContext(10, 1000) // maxLogSize = 10

	// Add 15 messages
	for i := 0; i < 15; i++ {
		sc.AddMessageLog("agent", "agent", "test", "message")
	}

	// Should only keep last 10
	if sc.LogSize() != 10 {
		t.Errorf("Expected 10 messages (maxLogSize), got %d", sc.LogSize())
	}
}

func TestSharedContext_MaxContext(t *testing.T) {
	sc := NewSharedContext(100, 5) // maxContext = 5

	// Add 8 entries
	for i := 0; i < 8; i++ {
		sc.Set("key"+string(rune(i)), i)
	}

	// Should only keep 5
	if sc.ContextSize() != 5 {
		t.Errorf("Expected 5 context entries, got %d", sc.ContextSize())
	}
}

func TestSharedContext_Keys(t *testing.T) {
	sc := NewSharedContext(100, 1000)

	// Add entries
	sc.Set("key1", "value1")
	sc.Set("key2", "value2")
	sc.Set("key3", "value3")

	keys := sc.Keys()
	if len(keys) != 3 {
		t.Errorf("Expected 3 keys, got %d", len(keys))
	}
}

func TestSharedContext_ThreadSafety(t *testing.T) {
	sc := NewSharedContext(100, 1000)

	// Concurrent writes
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(idx int) {
			for j := 0; j < 100; j++ {
				sc.Set("key", idx*100+j)
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Should not panic or deadlock
	keys := sc.Keys()
	if len(keys) == 0 {
		t.Errorf("Expected keys after concurrent writes")
	}
}

func TestSharedContext_ConcurrentReadWrite(t *testing.T) {
	sc := NewSharedContext(100, 1000)

	// Writer goroutines
	doneWrite := make(chan bool)
	for i := 0; i < 5; i++ {
		go func(idx int) {
			for j := 0; j < 50; j++ {
				sc.Set("key"+string(rune(idx)), j)
			}
			doneWrite <- true
		}(i)
	}

	// Reader goroutines
	doneRead := make(chan bool)
	for i := 0; i < 5; i++ {
		go func() {
			for j := 0; j < 50; j++ {
				sc.Keys()
				sc.Get("key")
			}
			doneRead <- true
		}()
	}

	// Wait for all
	for i := 0; i < 5; i++ {
		<-doneWrite
		<-doneRead
	}

	// Should not panic or deadlock
	if sc.ContextSize() == 0 {
		t.Errorf("Expected context to have data")
	}
}

func TestAgentMessage_Creation(t *testing.T) {
	msg := NewAgentMessage("pm", "dev", MsgRequest, "review code", "session-123")

	if msg.From != "pm" {
		t.Errorf("Expected 'pm', got %s", msg.From)
	}

	if msg.To != "dev" {
		t.Errorf("Expected 'dev', got %s", msg.To)
	}

	if msg.Type != MsgRequest {
		t.Errorf("Expected MsgRequest, got %s", msg.Type)
	}

	if msg.SessionID != "session-123" {
		t.Errorf("Expected 'session-123', got %s", msg.SessionID)
	}

	if msg.Timestamp == 0 {
		t.Errorf("Expected timestamp to be set")
	}

	if msg.ID == "" {
		t.Errorf("Expected ID to be generated")
	}
}

func TestAgentMessage_Types(t *testing.T) {
	types := []MessageType{
		MsgRequest,
		MsgResponse,
		MsgBroadcast,
		MsgContextUpdate,
		MsgHeartbeat,
		MsgTerminate,
	}

	for _, msgType := range types {
		msg := NewAgentMessage("agent1", "agent2", msgType, "test", "session")
		if msg.Type != msgType {
			t.Errorf("Expected %s, got %s", msgType, msg.Type)
		}
	}
}

// TestGenerateID_UniquenessUnderConcurrency verifica que generateID no genera
// duplicados bajo alta concurrencia (B-06 fix: contador atómico)
func TestGenerateID_UniquenessUnderConcurrency(t *testing.T) {
	const goroutines = 100
	const idsPerGoroutine = 100
	total := goroutines * idsPerGoroutine

	ids := make(chan string, total)
	var wg sync.WaitGroup

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < idsPerGoroutine; j++ {
				msg := NewAgentMessage("a", "b", MsgRequest, "test", "s")
				ids <- msg.ID
			}
		}()
	}

	wg.Wait()
	close(ids)

	seen := make(map[string]bool, total)
	for id := range ids {
		if seen[id] {
			t.Errorf("Duplicate ID found: %s", id)
		}
		seen[id] = true
	}

	if len(seen) != total {
		t.Errorf("Expected %d unique IDs, got %d", total, len(seen))
	}
}
