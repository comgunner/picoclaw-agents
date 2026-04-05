package channels

import "time"

// Pico Protocol message type constants.
const (
	// picoTypeMsgSend is sent from WebSocket client to server (user message).
	picoTypeMsgSend   = "message.send"
	picoTypeMediaSend = "media.send"
	picoTypePing      = "ping"

	// Server-to-client message types.
	picoTypeMsgCreate   = "message.create"
	picoTypeMsgUpdate   = "message.update"
	picoTypeMediaCreate = "media.create"
	picoTypeTypingStart = "typing.start"
	picoTypeTypingStop  = "typing.stop"
	picoTypeError       = "error"
	picoTypePong        = "pong"
)

// PicoMessage is the wire format for all Pico Protocol WebSocket messages.
type PicoMessage struct {
	Type      string         `json:"type"`
	ID        string         `json:"id,omitempty"`
	SessionID string         `json:"session_id,omitempty"`
	Timestamp int64          `json:"timestamp,omitempty"`
	Payload   map[string]any `json:"payload,omitempty"`
}

// newPicoMessage creates a PicoMessage with the current timestamp.
func newPicoMessage(msgType string, payload map[string]any) PicoMessage {
	return PicoMessage{
		Type:      msgType,
		Timestamp: time.Now().UnixMilli(),
		Payload:   payload,
	}
}

// newPicoError creates an error PicoMessage.
func newPicoError(code, message string) PicoMessage {
	return newPicoMessage(picoTypeError, map[string]any{
		"code":    code,
		"message": message,
	})
}
