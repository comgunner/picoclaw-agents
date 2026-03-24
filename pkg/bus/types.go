// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package bus

import (
	"errors"
)

var ErrBusClosed = errors.New("message bus is closed")

type InboundMessage struct {
	Channel    string            `json:"channel"`
	SenderID   string            `json:"sender_id"`
	ChatID     string            `json:"chat_id"`
	Content    string            `json:"content"`
	Media      []string          `json:"media,omitempty"`
	SessionKey string            `json:"session_key"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

type Button struct {
	Text string `json:"text"`
	Data string `json:"data"`
}

type OutboundMessage struct {
	Channel string   `json:"channel"`
	ChatID  string   `json:"chat_id"`
	Content string   `json:"content"`
	Media   []string `json:"media,omitempty"`
	Buttons []Button `json:"buttons,omitempty"`
}

type MessageHandler func(InboundMessage) error
