// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package channels

import "errors"

// BuildMediaScope builds a media scope string for channel/media lifecycle management.
// Used by Weixin channel and potentially other channels.
func BuildMediaScope(channel, chatID, messageID string) string {
	return channel + ":" + chatID + ":" + messageID
}

// Channel errors - exported for use by channel implementations
var (
	// ErrNotRunning indicates the channel is not running
	ErrNotRunning = errors.New("channel not running")

	// ErrTemporary indicates a transient failure (e.g. network timeout, 5xx)
	ErrTemporary = errors.New("temporary failure")

	// ErrSendFailed indicates a message send operation failed
	ErrSendFailed = errors.New("send failed")

	// ErrRateLimit indicates rate limit exceeded
	ErrRateLimit = errors.New("rate limited")
)
