// pkg/agent/a2a_integration_test.go
//
// Unit tests for A2AIntegration
// Verifica que SendTask y BroadcastStatus ruteen correctamente via hub (fixes B-01, B-02)
//
// QA A2A - Fase 4

package agent

import (
	"testing"

	"github.com/comgunner/picoclaw/pkg/agentcomm"
	"github.com/comgunner/picoclaw/pkg/mailbox"
)

func TestA2AIntegration_SendTask_RoutesToRecipient(t *testing.T) {
	// B-01: SendTask debe enviar al mailbox del DESTINATARIO, no al del remitente
	hub := mailbox.NewHub(100)

	senderMb := hub.Register("sender")
	hub.Register("recipient")

	ctx := agentcomm.NewSharedContext(100, 1000)
	senderAgent := &AgentInstance{ID: "sender", Status: "idle"}

	integration := NewA2AIntegration(senderAgent, senderMb, hub, ctx, nil)

	if err := integration.SendTask("recipient", "do something important", mailbox.PriorityHigh); err != nil {
		t.Fatalf("SendTask failed: %v", err)
	}

	// El REMITENTE no debe tener el mensaje en su propio mailbox
	if senderMb.Size() != 0 {
		t.Errorf("Sender mailbox should be empty, got %d messages", senderMb.Size())
	}

	// El DESTINATARIO debe tener el mensaje
	recipientMb, ok := hub.Get("recipient")
	if !ok {
		t.Fatalf("Recipient mailbox not found")
	}
	if recipientMb.Size() != 1 {
		t.Errorf("Recipient mailbox should have 1 message, got %d", recipientMb.Size())
	}

	// Verificar contenido del mensaje
	msg, err := recipientMb.Receive()
	if err != nil {
		t.Fatalf("Failed to receive message: %v", err)
	}
	if msg.Content != "do something important" {
		t.Errorf("Wrong message content: %s", msg.Content)
	}
	if msg.Priority != mailbox.PriorityHigh {
		t.Errorf("Wrong priority: %v", msg.Priority)
	}
}

func TestA2AIntegration_SendStatus_RoutesToRecipient(t *testing.T) {
	// B-01: SendStatus debe enviar al mailbox del DESTINATARIO
	hub := mailbox.NewHub(100)

	senderMb := hub.Register("sender")
	hub.Register("recipient")

	ctx := agentcomm.NewSharedContext(100, 1000)
	senderAgent := &AgentInstance{ID: "sender", Status: "idle"}

	integration := NewA2AIntegration(senderAgent, senderMb, hub, ctx, nil)

	if err := integration.SendStatus("recipient", "task completed"); err != nil {
		t.Fatalf("SendStatus failed: %v", err)
	}

	// Verificar que el mensaje llegó al destinatario
	recipientMb, _ := hub.Get("recipient")
	if recipientMb.Size() != 1 {
		t.Errorf("Recipient mailbox should have 1 message, got %d", recipientMb.Size())
	}

	// El remitente no debe tener el mensaje
	if senderMb.Size() != 0 {
		t.Errorf("Sender mailbox should be empty, got %d messages", senderMb.Size())
	}
}

func TestA2AIntegration_BroadcastStatus_ReachesAll(t *testing.T) {
	// B-02: BroadcastStatus debe llegar a todos los agentes EXCEPTO el remitente
	hub := mailbox.NewHub(100)

	senderMb := hub.Register("sender")
	hub.Register("agent1")
	hub.Register("agent2")

	ctx := agentcomm.NewSharedContext(100, 1000)
	senderAgent := &AgentInstance{ID: "sender", Status: "idle"}

	integration := NewA2AIntegration(senderAgent, senderMb, hub, ctx, nil)

	if err := integration.BroadcastStatus("system online"); err != nil {
		t.Fatalf("BroadcastStatus failed: %v", err)
	}

	// agent1 y agent2 deben recibir el broadcast
	mb1, _ := hub.Get("agent1")
	mb2, _ := hub.Get("agent2")

	if mb1.Size() != 1 {
		t.Errorf("agent1 should have 1 message, got %d", mb1.Size())
	}
	if mb2.Size() != 1 {
		t.Errorf("agent2 should have 1 message, got %d", mb2.Size())
	}

	// El REMITENTE no debe recibir su propio broadcast
	if senderMb.Size() != 0 {
		t.Errorf("Sender should not receive own broadcast, got %d messages", senderMb.Size())
	}

	// Verificar contenido
	msg, _ := mb1.Receive()
	if msg.Content != "system online" {
		t.Errorf("Wrong broadcast content: %s", msg.Content)
	}
}

func TestA2AIntegration_ProcessMailboxMessages(t *testing.T) {
	// Verifica que ProcessMailboxMessages procesa mensajes y almacena en SharedContext
	hub := mailbox.NewHub(100)
	mb := hub.Register("receiver")
	ctx := agentcomm.NewSharedContext(100, 1000)
	agent := &AgentInstance{ID: "receiver", Status: "idle"}

	integration := NewA2AIntegration(agent, mb, hub, ctx, nil)

	// Enviar un mensaje de tarea al receptor
	hub.SendTo("receiver", mailbox.Message{
		Type:     mailbox.MessageTypeTask,
		From:     "sender",
		To:       "receiver",
		Content:  "complete task X",
		Priority: mailbox.PriorityHigh,
	})

	if err := integration.ProcessMailboxMessages(); err != nil {
		t.Fatalf("ProcessMailboxMessages failed: %v", err)
	}

	// Verificar que el mensaje fue almacenado en SharedContext
	keys := ctx.Keys()
	found := false
	for _, key := range keys {
		if len(key) >= 5 && key[:5] == "task:" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Task message not stored in SharedContext after processing")
	}

	// El agente debe estar en estado "running" después de recibir una tarea
	if agent.Status != "running" {
		t.Errorf("Agent status should be 'running' after task, got %s", agent.Status)
	}
}

func TestA2AIntegration_SendTask_NoHub(t *testing.T) {
	// Sin hub configurado debe retornar error
	mb := mailbox.NewMailbox("sender", 100)
	ctx := agentcomm.NewSharedContext(100, 1000)
	agent := &AgentInstance{ID: "sender"}

	integration := NewA2AIntegration(agent, mb, nil, ctx, nil)

	err := integration.SendTask("recipient", "task", mailbox.PriorityNormal)
	if err == nil {
		t.Errorf("Expected error when hub is nil")
	}
}

func TestA2AIntegration_BroadcastStatus_NoHub(t *testing.T) {
	// Sin hub configurado debe retornar error
	mb := mailbox.NewMailbox("sender", 100)
	ctx := agentcomm.NewSharedContext(100, 1000)
	agent := &AgentInstance{ID: "sender"}

	integration := NewA2AIntegration(agent, mb, nil, ctx, nil)

	err := integration.BroadcastStatus("status update")
	if err == nil {
		t.Errorf("Expected error when hub is nil")
	}
}
