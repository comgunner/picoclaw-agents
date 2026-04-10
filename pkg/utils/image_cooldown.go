// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package utils

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	_ "modernc.org/sqlite" // Pure Go SQLite driver, no CGO required
)

// ImageCooldown manages a global cooldown for image generation.
// Persists in SQLite to survive process restarts.
// The database file is stored in <workspace>/tmp/picoclaw_image_cooldown.db
// to ensure correct permissions and avoid OS /tmp cleanup.
type ImageCooldown struct {
	dbPath string
	mu     sync.Mutex
	db     *sql.DB
}

// NewImageCooldown creates a new cooldown instance.
// workspace: the agent workspace directory (e.g. "~/.picoclaw/workspace").
// The database is created at <workspace>/tmp/picoclaw_image_cooldown.db.
// If PICOCLAW_IMAGE_COOLDOWN_DB env var is set, it overrides the path.
func NewImageCooldown(workspace string) (*ImageCooldown, error) {
	dbPath := os.Getenv(EnvImageCooldownDB)
	if dbPath == "" {
		ws := resolveCooldownWorkspacePath(workspace)
		tmpDir := filepath.Join(ws, "tmp")
		dbPath = filepath.Join(tmpDir, "picoclaw_image_cooldown.db")
	}

	c := &ImageCooldown{dbPath: dbPath}
	if err := c.initDB(); err != nil {
		return nil, fmt.Errorf("init cooldown DB: %w", err)
	}
	return c, nil
}

// resolveCooldownWorkspacePath safely resolves the workspace path.
func resolveCooldownWorkspacePath(workspace string) string {
	ws := strings.TrimSpace(workspace)
	if ws == "" {
		home, err := os.UserHomeDir()
		if err == nil {
			return filepath.Join(home, ".picoclaw", "workspace")
		}
		return "./workspace"
	}
	if strings.HasPrefix(ws, "~") {
		home, err := os.UserHomeDir()
		if err == nil {
			ws = filepath.Join(home, ws[1:])
		}
	}
	return filepath.Clean(ws)
}

func (c *ImageCooldown) initDB() error {
	// Ensure parent directory exists
	dir := filepath.Dir(c.dbPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create cooldown dir: %w", err)
	}

	db, err := sql.Open("sqlite", c.dbPath)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS cooldown (
			id INTEGER PRIMARY KEY CHECK (id = 1),
			started_at REAL NOT NULL,
			duration_seconds REAL NOT NULL,
			provider TEXT DEFAULT 'antigravity',
			model TEXT DEFAULT 'gemini-3.1-flash-image'
		)
	`)
	if err != nil {
		db.Close()
		return err
	}

	c.db = db
	return nil
}

// IsOnCooldown returns true if the cooldown is currently active.
func (c *ImageCooldown) IsOnCooldown() bool {
	return c.GetRemaining() > 0
}

// GetRemaining returns the remaining cooldown seconds (0 if not active).
func (c *ImageCooldown) GetRemaining() float64 {
	c.mu.Lock()
	defer c.mu.Unlock()

	var startedAt, duration float64
	err := c.db.QueryRow(
		"SELECT started_at, duration_seconds FROM cooldown WHERE id = 1",
	).Scan(&startedAt, &duration)

	if err == sql.ErrNoRows {
		return 0
	}
	if err != nil {
		return 0
	}

	elapsed := float64(time.Now().UnixMilli())/1000.0 - startedAt
	remaining := duration - elapsed
	if remaining < 0 {
		// Cooldown expired, clear it
		_, _ = c.db.Exec("DELETE FROM cooldown WHERE id = 1")
		return 0
	}
	return remaining
}

// GetInfo returns detailed cooldown information.
// FIX: Does NOT call GetRemaining() to avoid deadlock.
// sync.Mutex is NOT reentrant — calling GetRemaining() from GetInfo()
// would cause a deadlock because both try to acquire c.mu.
func (c *ImageCooldown) GetInfo() map[string]any {
	c.mu.Lock()
	defer c.mu.Unlock()

	var startedAt, duration float64
	var provider, model string
	err := c.db.QueryRow(
		"SELECT started_at, duration_seconds, provider, model FROM cooldown WHERE id = 1",
	).Scan(&startedAt, &duration, &provider, &model)

	if err == sql.ErrNoRows {
		return map[string]any{
			"active":    false,
			"remaining": 0.0,
		}
	}
	if err != nil {
		return map[string]any{
			"active":    false,
			"remaining": 0.0,
			"error":     err.Error(),
		}
	}

	elapsed := float64(time.Now().UnixMilli())/1000.0 - startedAt
	remaining := duration - elapsed
	if remaining < 0 {
		_, _ = c.db.Exec("DELETE FROM cooldown WHERE id = 1")
		return map[string]any{
			"active":    false,
			"remaining": 0.0,
		}
	}

	return map[string]any{
		"active":     true,
		"remaining":  remaining,
		"started_at": time.Unix(int64(startedAt), 0).Format(time.RFC3339),
		"duration":   duration,
		"provider":   provider,
		"model":      model,
	}
}

// Set activates the cooldown for the specified duration in seconds.
func (c *ImageCooldown) Set(durationSeconds float64, provider, model string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if provider == "" {
		provider = "antigravity"
	}
	if model == "" {
		model = "gemini-3.1-flash-image"
	}

	now := float64(time.Now().UnixMilli()) / 1000.0
	_, err := c.db.Exec(
		"INSERT OR REPLACE INTO cooldown (id, started_at, duration_seconds, provider, model) VALUES (1, ?, ?, ?, ?)",
		now, durationSeconds, provider, model,
	)
	return err
}

// Wait blocks until the cooldown expires or the timeout is reached.
func (c *ImageCooldown) Wait(timeout time.Duration) bool {
	remaining := c.GetRemaining()
	if remaining <= 0 {
		return true
	}

	waitTime := time.Duration(remaining * float64(time.Second))
	if timeout > 0 && waitTime > timeout {
		waitTime = timeout
	}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	deadline := time.Now().Add(waitTime)
	for time.Now().Before(deadline) {
		<-ticker.C
		if !c.IsOnCooldown() {
			return true
		}
	}
	return !c.IsOnCooldown()
}

// Clear manually clears the cooldown (use with caution).
func (c *ImageCooldown) Clear() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, err := c.db.Exec("DELETE FROM cooldown WHERE id = 1")
	return err
}

// Close closes the database connection.
func (c *ImageCooldown) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}
