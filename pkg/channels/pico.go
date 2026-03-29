package channels

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/comgunner/picoclaw/pkg/bus"
	"github.com/comgunner/picoclaw/pkg/config"
	"github.com/comgunner/picoclaw/pkg/logger"
)

// Sentinel errors for the Pico WebSocket channel.
var (
	// ErrPicoNotRunning is returned when Send is attempted on a stopped channel.
	ErrPicoNotRunning = errors.New("pico channel is not running")
	// ErrPicoSendFailed is returned when no active WebSocket connections exist for a session.
	ErrPicoSendFailed = errors.New("pico send failed")
	// ErrPicoTemporary signals a transient rejection (e.g. connection limit reached).
	ErrPicoTemporary = errors.New("pico temporary error")
)

// picoConn represents a single active WebSocket connection.
type picoConn struct {
	id        string
	conn      *websocket.Conn
	sessionID string
	writeMu   sync.Mutex
	closed    atomic.Bool
	cancel    context.CancelFunc // cancels per-connection goroutines (e.g. pingLoop)
}

// writeJSON sends a JSON-encoded message with write-lock protection.
func (pc *picoConn) writeJSON(v any) error {
	if pc.closed.Load() {
		return fmt.Errorf("connection closed")
	}
	pc.writeMu.Lock()
	defer pc.writeMu.Unlock()
	return pc.conn.WriteJSON(v)
}

// close cleanly shuts down the connection exactly once.
func (pc *picoConn) close() {
	if pc.closed.CompareAndSwap(false, true) {
		if pc.cancel != nil {
			pc.cancel()
		}
		pc.conn.Close()
	}
}

// PicoChannel implements the Pico Protocol WebSocket server channel.
// The WebUI frontend connects to this channel to chat with agents.
type PicoChannel struct {
	*BaseChannel
	cfg                config.PicoConfig
	upgrader           websocket.Upgrader
	connections        map[string]*picoConn            // connID -> *picoConn
	sessionConnections map[string]map[string]*picoConn // sessionID -> connID -> *picoConn
	connsMu            sync.RWMutex
	running            atomic.Bool // separate from BaseChannel.running (which is unexported)
	ctx                context.Context
	cancel             context.CancelFunc
}

// NewPicoChannel creates a new Pico Protocol WebSocket server channel.
func NewPicoChannel(cfg config.PicoConfig, messageBus *bus.MessageBus) (*PicoChannel, error) {
	if cfg.Token() == "" {
		return nil, fmt.Errorf("pico token is required")
	}

	base := NewBaseChannel("pico", cfg, messageBus, nil)

	allowOrigins := cfg.AllowOrigins
	checkOrigin := func(r *http.Request) bool {
		if len(allowOrigins) == 0 {
			return true // allow all origins when not configured
		}
		origin := r.Header.Get("Origin")
		for _, allowed := range allowOrigins {
			if allowed == "*" || allowed == origin {
				return true
			}
		}
		return false
	}

	return &PicoChannel{
		BaseChannel: base,
		cfg:         cfg,
		upgrader: websocket.Upgrader{
			CheckOrigin:     checkOrigin,
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		connections:        make(map[string]*picoConn),
		sessionConnections: make(map[string]map[string]*picoConn),
	}, nil
}

// IsRunning overrides BaseChannel.IsRunning using this channel's own atomic bool.
func (c *PicoChannel) IsRunning() bool {
	return c.running.Load()
}

// Start implements Channel.
func (c *PicoChannel) Start(ctx context.Context) error {
	logger.InfoC("pico", "Starting Pico Protocol channel")
	c.ctx, c.cancel = context.WithCancel(ctx)
	c.running.Store(true)
	logger.InfoC("pico", "Pico Protocol channel started")
	return nil
}

// Stop implements Channel.
func (c *PicoChannel) Stop(ctx context.Context) error {
	logger.InfoC("pico", "Stopping Pico Protocol channel")
	c.running.Store(false)

	for _, pc := range c.takeAllConnections() {
		pc.close()
	}

	if c.cancel != nil {
		c.cancel()
	}

	logger.InfoC("pico", "Pico Protocol channel stopped")
	return nil
}

// Send implements Channel — delivers a message to all WebSocket connections for the session.
func (c *PicoChannel) Send(ctx context.Context, msg bus.OutboundMessage) error {
	if !c.IsRunning() {
		return ErrPicoNotRunning
	}

	outMsg := newPicoMessage(picoTypeMsgCreate, map[string]any{
		"content": msg.Content,
	})

	return c.broadcastToSession(msg.ChatID, outMsg)
}

// WebhookPath returns the base HTTP path handled by this channel.
func (c *PicoChannel) WebhookPath() string { return "/pico/" }

// ServeHTTP implements http.Handler — routes WebSocket upgrade requests.
func (c *PicoChannel) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/pico")

	switch path {
	case "/ws", "/ws/":
		c.handleWebSocket(w, r)
	default:
		http.NotFound(w, r)
	}
}

// EditMessage sends a message.update to all connections in a session.
func (c *PicoChannel) EditMessage(ctx context.Context, chatID, messageID, content string) error {
	outMsg := newPicoMessage(picoTypeMsgUpdate, map[string]any{
		"message_id": messageID,
		"content":    content,
	})
	return c.broadcastToSession(chatID, outMsg)
}

// StartTyping sends typing.start to the session and returns a stop function.
func (c *PicoChannel) StartTyping(ctx context.Context, chatID string) (func(), error) {
	startMsg := newPicoMessage(picoTypeTypingStart, nil)
	if err := c.broadcastToSession(chatID, startMsg); err != nil {
		return func() {}, err
	}
	return func() {
		_ = c.broadcastToSession(chatID, newPicoMessage(picoTypeTypingStop, nil))
	}, nil
}

// broadcastToSession fans a message out to every connection sharing the given session.
func (c *PicoChannel) broadcastToSession(chatID string, msg PicoMessage) error {
	// chatID is in "pico:<sessionID>" format.
	sessionID := strings.TrimPrefix(chatID, "pico:")
	msg.SessionID = sessionID

	// BUG-04 FIX: Distinguish between "no connections" (expected, not an error) and
	// "all connections failed" (actual error). Previously, both cases returned an error.
	connections := c.sessionConnectionsSnapshot(sessionID)
	if len(connections) == 0 {
		// No active WebSocket connections for this session — this is expected behavior
		// when the WebUI is not open. The message was already enqueued upstream.
		logger.DebugCF("pico", "No active connections for session", map[string]any{
			"session_id": sessionID,
		})
		return nil
	}

	var sent bool
	for _, pc := range connections {
		if err := pc.writeJSON(msg); err != nil {
			logger.DebugCF("pico", "Write to connection failed", map[string]any{
				"conn_id": pc.id,
				"error":   err.Error(),
			})
		} else {
			sent = true
		}
	}

	if !sent {
		// All connections failed to receive the message — this is an actual error
		return fmt.Errorf("all connections failed for session %s: %w", sessionID, ErrPicoSendFailed)
	}
	return nil
}

// handleWebSocket upgrades an HTTP connection to WebSocket and starts the read loop.
func (c *PicoChannel) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	if !c.IsRunning() {
		http.Error(w, "channel not running", http.StatusServiceUnavailable)
		return
	}

	if !c.picoAuthenticate(r) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	const defaultMaxConns = 100
	if c.currentConnCount() >= defaultMaxConns {
		http.Error(w, "too many connections", http.StatusServiceUnavailable)
		return
	}

	// Echo the matched subprotocol back so the browser accepts the upgrade.
	var responseHeader http.Header
	if proto := c.matchedSubprotocol(r); proto != "" {
		responseHeader = http.Header{"Sec-WebSocket-Protocol": {proto}}
	}

	conn, err := c.upgrader.Upgrade(w, r, responseHeader)
	if err != nil {
		logger.ErrorCF("pico", "WebSocket upgrade failed", map[string]any{
			"error": err.Error(),
		})
		return
	}

	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" {
		sessionID = uuid.New().String()
	}

	pc, err := c.createAndAddConnection(conn, sessionID, defaultMaxConns)
	if err != nil {
		_ = conn.WriteControl(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseTryAgainLater, "too many connections"),
			time.Now().Add(2*time.Second),
		)
		_ = conn.Close()
		return
	}

	logger.InfoCF("pico", "WebSocket client connected", map[string]any{
		"conn_id":    pc.id,
		"session_id": sessionID,
	})

	go c.picoReadLoop(pc)
}

// picoAuthenticate validates the request via Bearer header, WebSocket subprotocol, or query param.
func (c *PicoChannel) picoAuthenticate(r *http.Request) bool {
	token := c.cfg.Token()
	if token == "" {
		return false
	}

	// 1. Authorization: Bearer <token>
	if after, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer "); ok && after == token {
		return true
	}

	// 2. Sec-WebSocket-Protocol "token.<value>"
	if c.matchedSubprotocol(r) != "" {
		return true
	}

	// 3. Query parameter (only when explicitly enabled)
	if c.cfg.AllowTokenQuery && r.URL.Query().Get("token") == token {
		return true
	}

	return false
}

// matchedSubprotocol returns the "token.<value>" subprotocol that matches the configured token, or "".
func (c *PicoChannel) matchedSubprotocol(r *http.Request) string {
	token := c.cfg.Token()
	for _, proto := range websocket.Subprotocols(r) {
		if after, ok := strings.CutPrefix(proto, "token."); ok && after == token {
			return proto
		}
	}
	return ""
}

// picoReadLoop drives the WebSocket read pump for a single connection.
func (c *PicoChannel) picoReadLoop(pc *picoConn) {
	defer func() {
		pc.close()
		if removed := c.removeConnection(pc.id); removed != nil {
			logger.InfoCF("pico", "WebSocket client disconnected", map[string]any{
				"conn_id":    removed.id,
				"session_id": removed.sessionID,
			})
		}
	}()

	const readTimeout = 60 * time.Second
	const pingInterval = 30 * time.Second

	_ = pc.conn.SetReadDeadline(time.Now().Add(readTimeout))
	pc.conn.SetPongHandler(func(string) error {
		return pc.conn.SetReadDeadline(time.Now().Add(readTimeout))
	})

	go c.picoPingLoop(pc, pingInterval)

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
		}

		_, rawMsg, err := pc.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				logger.DebugCF("pico", "WebSocket read error", map[string]any{
					"conn_id": pc.id,
					"error":   err.Error(),
				})
			}
			return
		}

		_ = pc.conn.SetReadDeadline(time.Now().Add(readTimeout))

		var msg PicoMessage
		if err := json.Unmarshal(rawMsg, &msg); err != nil {
			_ = pc.writeJSON(newPicoError("invalid_message", "failed to parse message"))
			continue
		}

		c.picoHandleMessage(pc, msg)
	}
}

// picoPingLoop sends periodic WebSocket ping frames to keep the connection alive.
func (c *PicoChannel) picoPingLoop(pc *picoConn, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			if pc.closed.Load() {
				return
			}
			pc.writeMu.Lock()
			err := pc.conn.WriteMessage(websocket.PingMessage, nil)
			pc.writeMu.Unlock()
			if err != nil {
				return
			}
		}
	}
}

// picoHandleMessage dispatches a decoded Pico Protocol message to the appropriate handler.
func (c *PicoChannel) picoHandleMessage(pc *picoConn, msg PicoMessage) {
	switch msg.Type {
	case picoTypePing:
		pong := newPicoMessage(picoTypePong, nil)
		pong.ID = msg.ID
		_ = pc.writeJSON(pong)

	case picoTypeMsgSend:
		c.picoHandleMessageSend(pc, msg)

	default:
		_ = pc.writeJSON(newPicoError("unknown_type", fmt.Sprintf("unknown message type: %s", msg.Type)))
	}
}

// picoHandleMessageSend processes an inbound message.send frame from a WebSocket client.
func (c *PicoChannel) picoHandleMessageSend(pc *picoConn, msg PicoMessage) {
	content, _ := msg.Payload["content"].(string)
	if strings.TrimSpace(content) == "" {
		_ = pc.writeJSON(newPicoError("empty_content", "message content is empty"))
		return
	}

	sessionID := msg.SessionID
	if sessionID == "" {
		sessionID = pc.sessionID
	}

	chatID := "pico:" + sessionID
	senderID := "pico-user"

	metadata := map[string]string{
		"platform":   "pico",
		"session_id": sessionID,
		"conn_id":    pc.id,
		"peer_kind":  "direct",
		"peer_id":    sessionID,
	}

	logger.DebugCF("pico", "Received message", map[string]any{
		"session_id": sessionID,
		"preview":    picoTruncate(content, 50),
	})

	c.HandleMessage(senderID, chatID, content, nil, metadata)
}

// --- Connection index helpers ---

// createAndAddConnection atomically checks the connection limit and registers a new conn.
func (c *PicoChannel) createAndAddConnection(conn *websocket.Conn, sessionID string, maxConns int) (*picoConn, error) {
	c.connsMu.Lock()
	defer c.connsMu.Unlock()
	if len(c.connections) >= maxConns {
		return nil, ErrPicoTemporary
	}

	var connID string
	for {
		connID = uuid.New().String()
		if _, exists := c.connections[connID]; !exists {
			break
		}
	}

	pc := &picoConn{id: connID, conn: conn, sessionID: sessionID}

	c.connections[pc.id] = pc
	bySession, ok := c.sessionConnections[pc.sessionID]
	if !ok {
		bySession = make(map[string]*picoConn)
		c.sessionConnections[pc.sessionID] = bySession
	}
	bySession[pc.id] = pc

	return pc, nil
}

// removeConnection removes a connection from both indexes and returns it when found.
func (c *PicoChannel) removeConnection(connID string) *picoConn {
	c.connsMu.Lock()
	defer c.connsMu.Unlock()

	pc, ok := c.connections[connID]
	if !ok {
		return nil
	}

	delete(c.connections, connID)
	if bySession, ok := c.sessionConnections[pc.sessionID]; ok {
		delete(bySession, connID)
		if len(bySession) == 0 {
			delete(c.sessionConnections, pc.sessionID)
		}
	}

	return pc
}

// takeAllConnections snapshots and clears all connection indexes atomically.
func (c *PicoChannel) takeAllConnections() []*picoConn {
	c.connsMu.Lock()
	defer c.connsMu.Unlock()

	all := make([]*picoConn, 0, len(c.connections))
	for _, pc := range c.connections {
		all = append(all, pc)
	}
	clear(c.connections)
	clear(c.sessionConnections)

	return all
}

// sessionConnectionsSnapshot returns a read-locked snapshot of connections for a session.
func (c *PicoChannel) sessionConnectionsSnapshot(sessionID string) []*picoConn {
	c.connsMu.RLock()
	defer c.connsMu.RUnlock()

	bySession, ok := c.sessionConnections[sessionID]
	if !ok || len(bySession) == 0 {
		return nil
	}

	conns := make([]*picoConn, 0, len(bySession))
	for _, pc := range bySession {
		conns = append(conns, pc)
	}
	return conns
}

// currentConnCount returns the current connection count under a read lock.
func (c *PicoChannel) currentConnCount() int {
	c.connsMu.RLock()
	defer c.connsMu.RUnlock()
	return len(c.connections)
}

// picoTruncate shortens s to at most maxLen runes, appending "..." if truncated.
func picoTruncate(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen]) + "..."
}
