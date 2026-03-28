// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package channels

import (
	"context"
	"strings"

	"github.com/comgunner/picoclaw/pkg/bus"
	"github.com/comgunner/picoclaw/pkg/config"
)

type Channel interface {
	Name() string
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Send(ctx context.Context, msg bus.OutboundMessage) error
	IsRunning() bool
	IsAllowed(senderID string) bool
}

// BaseChannelOption is a functional option for configuring a BaseChannel.
type BaseChannelOption func(*BaseChannel)

// WithGroupTrigger sets the group trigger configuration on the channel.
func WithGroupTrigger(gt config.GroupTriggerConfig) BaseChannelOption {
	return func(c *BaseChannel) {
		c.groupTrigger = gt
	}
}

type BaseChannel struct {
	config       any
	bus          *bus.MessageBus
	running      bool
	name         string
	allowList    []string
	groupTrigger config.GroupTriggerConfig
}

func NewBaseChannel(
	name string,
	cfg any,
	msgBus *bus.MessageBus,
	allowList []string,
	opts ...BaseChannelOption,
) *BaseChannel {
	c := &BaseChannel{
		config:    cfg,
		bus:       msgBus,
		name:      name,
		allowList: allowList,
		running:   false,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *BaseChannel) Name() string {
	return c.name
}

func (c *BaseChannel) IsRunning() bool {
	return c.running
}

func (c *BaseChannel) IsAllowed(senderID string) bool {
	if len(c.allowList) == 0 {
		return true
	}

	// Extract parts from compound senderID like "123456|username"
	idPart := senderID
	userPart := ""
	if idx := strings.Index(senderID, "|"); idx > 0 {
		idPart = senderID[:idx]
		userPart = senderID[idx+1:]
	}

	for _, allowed := range c.allowList {
		// Strip leading "@" from allowed value for username matching
		trimmed := strings.TrimPrefix(allowed, "@")
		allowedID := trimmed
		allowedUser := ""
		if idx := strings.Index(trimmed, "|"); idx > 0 {
			allowedID = trimmed[:idx]
			allowedUser = trimmed[idx+1:]
		}

		// Support either side using "id|username" compound form.
		// This keeps backward compatibility with legacy Telegram allowlist entries.
		if senderID == allowed ||
			idPart == allowed ||
			senderID == trimmed ||
			idPart == trimmed ||
			idPart == allowedID ||
			(allowedUser != "" && senderID == allowedUser) ||
			(userPart != "" && (userPart == allowed || userPart == trimmed || userPart == allowedUser)) {
			return true
		}
	}

	return false
}

func (c *BaseChannel) HandleMessage(senderID, chatID, content string, media []string, metadata map[string]string) {
	if !c.IsAllowed(senderID) {
		return
	}

	msg := bus.InboundMessage{
		Channel:  c.name,
		SenderID: senderID,
		ChatID:   chatID,
		Content:  content,
		Media:    media,
		Metadata: metadata,
	}

	c.bus.PublishInbound(msg)
}

func (c *BaseChannel) setRunning(running bool) {
	c.running = running
}

// IsAllowedSender checks whether a structured SenderInfo is in the allowList.
// Supports plain PlatformID, canonical "platform:id", @username, and "id|username" formats.
func (c *BaseChannel) IsAllowedSender(sender bus.SenderInfo) bool {
	if len(c.allowList) == 0 {
		return true
	}
	for _, entry := range c.allowList {
		if entry == sender.CanonicalID || entry == sender.PlatformID {
			return true
		}
		if strings.HasPrefix(entry, "@") && strings.TrimPrefix(entry, "@") == sender.Username {
			return true
		}
		if idx := strings.Index(entry, "|"); idx > 0 {
			id, user := entry[:idx], entry[idx+1:]
			if id == sender.PlatformID || id == sender.CanonicalID || user == sender.Username {
				return true
			}
		}
	}
	return false
}

// ShouldRespondInGroup decides whether the bot should respond to a group message,
// and returns the content to use (with any trigger prefix stripped).
//
// Rules:
//  1. If isMentioned: always respond, content unchanged.
//  2. If content matches a configured prefix: respond, strip prefix.
//  3. If prefixes are configured or MentionOnly is set: do NOT respond (no match).
//  4. Default (no restrictions): respond permissively.
func (c *BaseChannel) ShouldRespondInGroup(isMentioned bool, content string) (bool, string) {
	if isMentioned {
		return true, content
	}

	gt := c.groupTrigger
	for _, prefix := range gt.Prefixes {
		if prefix == "" {
			continue
		}
		if strings.HasPrefix(content, prefix) {
			stripped := strings.TrimLeft(content[len(prefix):], " ")
			return true, stripped
		}
	}

	if len(gt.Prefixes) > 0 || gt.MentionOnly {
		return false, content
	}

	return true, content
}
