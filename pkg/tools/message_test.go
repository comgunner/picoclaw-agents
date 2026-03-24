// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package tools

import (
	"context"
	"errors"
	"testing"

	"github.com/comgunner/picoclaw/pkg/bus"
)

func TestMessageTool_Execute_Success(t *testing.T) {
	tool := NewMessageTool()
	tool.SetContext("test-channel", "test-chat-id")

	var sentChannel, sentChatID, sentContent string
	var sentMedia []string
	var sentButtons []bus.Button
	tool.SetSendCallback(func(channel, chatID, content string, media []string, buttons []bus.Button) error {
		sentChannel = channel
		sentChatID = chatID
		sentContent = content
		sentMedia = media
		sentButtons = buttons
		return nil
	})

	ctx := context.Background()
	args := map[string]any{
		"content": "Hello, world!",
	}

	result := tool.Execute(ctx, args)

	// Verify message was sent with correct parameters
	if sentChannel != "test-channel" {
		t.Errorf("Expected channel 'test-channel', got '%s'", sentChannel)
	}
	if sentChatID != "test-chat-id" {
		t.Errorf("Expected chatID 'test-chat-id', got '%s'", sentChatID)
	}
	if sentContent != "Hello, world!" {
		t.Errorf("Expected content 'Hello, world!', got '%s'", sentContent)
	}

	// Verify context media
	if len(sentMedia) != 0 {
		t.Error("Expected 0 media files")
	}

	// Test with media
	argsMedia := map[string]any{
		"content": "Message with image",
		"media":   []any{"/path/to/image.jpg"},
	}
	tool.Execute(ctx, argsMedia)
	if len(sentMedia) != 1 || sentMedia[0] != "/path/to/image.jpg" {
		t.Errorf("Expected 1 media file '/path/to/image.jpg', got %v", sentMedia)
	}

	// Test with buttons
	argsButtons := map[string]any{
		"content": "Confirm action",
		"buttons": []any{
			map[string]any{"text": "Yes", "data": "confirm"},
			map[string]any{"text": "No", "data": "cancel"},
		},
	}
	tool.Execute(ctx, argsButtons)
	if len(sentButtons) != 2 || sentButtons[0].Text != "Yes" || sentButtons[1].Data != "cancel" {
		t.Errorf("Expected 2 buttons [Yes/confirm, No/cancel], got %v", sentButtons)
	}

	// Verify ToolResult meets US-011 criteria:
	// - Send success returns SilentResult (Silent=true)
	if !result.Silent {
		t.Error("Expected Silent=true for successful send")
	}

	// - ForLLM contains send status description
	if result.ForLLM != "Message sent to test-channel:test-chat-id" {
		t.Errorf("Expected ForLLM 'Message sent to test-channel:test-chat-id', got '%s'", result.ForLLM)
	}

	// - ForUser is empty (user already received message directly)
	if result.ForUser != "" {
		t.Errorf("Expected ForUser to be empty, got '%s'", result.ForUser)
	}

	// - IsError should be false
	if result.IsError {
		t.Error("Expected IsError=false for successful send")
	}
}

func TestMessageTool_Execute_WithCustomChannel(t *testing.T) {
	tool := NewMessageTool()
	tool.SetContext("default-channel", "default-chat-id")

	var sentChannel, sentChatID string
	tool.SetSendCallback(func(channel, chatID, content string, media []string, buttons []bus.Button) error {
		sentChannel = channel
		sentChatID = chatID
		return nil
	})

	ctx := context.Background()
	args := map[string]any{
		"content": "Test message",
		"channel": "custom-channel",
		"chat_id": "custom-chat-id",
	}

	result := tool.Execute(ctx, args)

	// Verify custom channel/chatID were used instead of defaults
	if sentChannel != "custom-channel" {
		t.Errorf("Expected channel 'custom-channel', got '%s'", sentChannel)
	}
	if sentChatID != "custom-chat-id" {
		t.Errorf("Expected chatID 'custom-chat-id', got '%s'", sentChatID)
	}

	if !result.Silent {
		t.Error("Expected Silent=true")
	}
	if result.ForLLM != "Message sent to custom-channel:custom-chat-id" {
		t.Errorf("Expected ForLLM 'Message sent to custom-channel:custom-chat-id', got '%s'", result.ForLLM)
	}
}

func TestMessageTool_Execute_SendFailure(t *testing.T) {
	tool := NewMessageTool()
	tool.SetContext("test-channel", "test-chat-id")

	sendErr := errors.New("network error")
	tool.SetSendCallback(func(channel, chatID, content string, media []string, buttons []bus.Button) error {
		return sendErr
	})

	ctx := context.Background()
	args := map[string]any{
		"content": "Test message",
	}

	result := tool.Execute(ctx, args)

	// Verify ToolResult for send failure:
	// - Send failure returns ErrorResult (IsError=true)
	if !result.IsError {
		t.Error("Expected IsError=true for failed send")
	}

	// - ForLLM contains error description
	expectedErrMsg := "sending message: network error"
	if result.ForLLM != expectedErrMsg {
		t.Errorf("Expected ForLLM '%s', got '%s'", expectedErrMsg, result.ForLLM)
	}

	// - Err field should contain original error
	if result.Err == nil {
		t.Error("Expected Err to be set")
	}
	if result.Err != sendErr {
		t.Errorf("Expected Err to be sendErr, got %v", result.Err)
	}
}

func TestMessageTool_Execute_MissingContent(t *testing.T) {
	tool := NewMessageTool()
	tool.SetContext("test-channel", "test-chat-id")

	ctx := context.Background()
	args := map[string]any{} // content missing

	result := tool.Execute(ctx, args)

	// Verify error result for missing content
	if !result.IsError {
		t.Error("Expected IsError=true for missing content")
	}
	if result.ForLLM != "content is required" {
		t.Errorf("Expected ForLLM 'content is required', got '%s'", result.ForLLM)
	}
}

func TestMessageTool_Execute_NoTargetChannel(t *testing.T) {
	tool := NewMessageTool()
	// No SetContext called, so defaultChannel and defaultChatID are empty

	tool.SetSendCallback(func(channel, chatID, content string, media []string, buttons []bus.Button) error {
		return nil
	})

	ctx := context.Background()
	args := map[string]any{
		"content": "Test message",
	}

	result := tool.Execute(ctx, args)

	// Verify error when no target channel specified
	if !result.IsError {
		t.Error("Expected IsError=true when no target channel")
	}
	if result.ForLLM != "No target channel/chat specified" {
		t.Errorf("Expected ForLLM 'No target channel/chat specified', got '%s'", result.ForLLM)
	}
}

func TestMessageTool_Execute_NotConfigured(t *testing.T) {
	tool := NewMessageTool()
	tool.SetContext("test-channel", "test-chat-id")
	// No SetSendCallback called

	ctx := context.Background()
	args := map[string]any{
		"content": "Test message",
	}

	result := tool.Execute(ctx, args)

	// Verify error when send callback not configured
	if !result.IsError {
		t.Error("Expected IsError=true when send callback not configured")
	}
	if result.ForLLM != "Message sending not configured" {
		t.Errorf("Expected ForLLM 'Message sending not configured', got '%s'", result.ForLLM)
	}
}

func TestMessageTool_Name(t *testing.T) {
	tool := NewMessageTool()
	if tool.Name() != "message" {
		t.Errorf("Expected name 'message', got '%s'", tool.Name())
	}
}

func TestMessageTool_Description(t *testing.T) {
	tool := NewMessageTool()
	desc := tool.Description()
	if desc == "" {
		t.Error("Description should not be empty")
	}
}

func TestMessageTool_Parameters(t *testing.T) {
	tool := NewMessageTool()
	params := tool.Parameters()

	// Verify parameters structure
	typ, ok := params["type"].(string)
	if !ok || typ != "object" {
		t.Error("Expected type 'object'")
	}

	props, ok := params["properties"].(map[string]any)
	if !ok {
		t.Fatal("Expected properties to be a map")
	}

	// Check required properties
	required, ok := params["required"].([]string)
	if !ok || len(required) != 1 || required[0] != "content" {
		t.Error("Expected 'content' to be required")
	}

	// Check content property
	contentProp, ok := props["content"].(map[string]any)
	if !ok {
		t.Error("Expected 'content' property")
	}
	if contentProp["type"] != "string" {
		t.Error("Expected content type to be 'string'")
	}

	// Check channel property (optional)
	channelProp, ok := props["channel"].(map[string]any)
	if !ok {
		t.Error("Expected 'channel' property")
	}
	if channelProp["type"] != "string" {
		t.Error("Expected channel type to be 'string'")
	}

	// Check chat_id property (optional)
	chatIDProp, ok := props["chat_id"].(map[string]any)
	if !ok {
		t.Error("Expected 'chat_id' property")
	}
	if chatIDProp["type"] != "string" {
		t.Error("Expected chat_id type to be 'string'")
	}
}

func TestMessageToolGhostConversationPrevention(t *testing.T) {
	tool := NewMessageTool()

	// Mock send callback
	sentContent := ""
	tool.SetSendCallback(func(channel, chatID, content string, media []string, buttons []bus.Button) error {
		sentContent = content
		return nil
	})

	// Set context
	tool.SetContext("telegram", "123456789")

	// Test normal message execution
	args := map[string]any{
		"content": "Hello, this is a legitimate message!",
	}

	result := tool.Execute(context.Background(), args)

	// Verify the result
	if result.IsError {
		t.Errorf("Expected successful execution, got error: %v", result.Err)
	}

	if sentContent != "Hello, this is a legitimate message!" {
		t.Errorf("Expected content to be sent, got: %s", sentContent)
	}

	// Verify that the tool was marked as having sent in this round
	if !tool.HasSentInRound() {
		t.Error("Expected tool to be marked as having sent in this round")
	}

	// Test with reminder-like content to see if it still processes
	args2 := map[string]any{
		"content": "Remember me to take out trash at 8 PM",
	}

	result2 := tool.Execute(context.Background(), args2)

	// The tool should still execute normally - the prevention happens at the LLM level
	if result2.IsError {
		t.Errorf("Expected successful execution for reminder-like content, got error: %v", result2.Err)
	}

	if sentContent != "Remember me to take out trash at 8 PM" {
		t.Errorf("Expected reminder content to be sent, got: %s", sentContent)
	}

	t.Log("Message tool correctly handles both regular and reminder-like content")
}
