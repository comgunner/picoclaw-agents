// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package context

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/comgunner/picoclaw/pkg/logger"
)

// LazyReference represents a file reference without loading its content.
// Content is loaded on-demand only when explicitly requested.
type LazyReference struct {
	ID        string    // Unique identifier: #FILE_ABCD_12345678
	Path      string    // Absolute path to file
	Size      int64     // File size in bytes
	ModTime   time.Time // Last modification time
	Loaded    bool      // Whether content is currently loaded
	Content   string    // Cached content (empty if not loaded)
	LineCount int       // Number of lines (for reference)
}

// LazyLoader manages lazy loading of workspace files.
type LazyLoader struct {
	workspacePath string
	cache         map[string]*LazyReference
	maxCacheSize  int

	mu sync.RWMutex
}

// NewLazyLoader creates a new lazy loader for the workspace.
func NewLazyLoader(workspacePath string) *LazyLoader {
	return &LazyLoader{
		workspacePath: workspacePath,
		cache:         make(map[string]*LazyReference),
		maxCacheSize:  100, // Max 100 files in cache
	}
}

// ReferenceFile creates a lazy reference to a file without loading content.
// This is the PREFERRED way to reference files in prompts.
func (ll *LazyLoader) ReferenceFile(path string) (*LazyReference, error) {
	// Convert to absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve path: %w", err)
	}

	// Check cache first
	ll.mu.RLock()
	if ref, ok := ll.cache[absPath]; ok {
		ll.mu.RUnlock()
		return ref, nil
	}
	ll.mu.RUnlock()

	// Get file info
	info, err := os.Stat(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}

	// Generate unique ID
	id := generateFileID(absPath)

	// Create reference
	ref := &LazyReference{
		ID:        id,
		Path:      absPath,
		Size:      info.Size(),
		ModTime:   info.ModTime(),
		Loaded:    false,
		Content:   "",
		LineCount: 0,
	}

	// Add to cache
	ll.mu.Lock()
	// Evict oldest if cache is full
	if len(ll.cache) >= ll.maxCacheSize {
		ll.evictOldest()
	}
	ll.cache[absPath] = ref
	ll.mu.Unlock()

	logger.DebugCF("context", "Created lazy reference",
		map[string]any{
			"file_id": id,
			"path":    absPath,
			"size":    info.Size(),
		})

	return ref, nil
}

// LoadContent loads the file content into memory.
// Use this only when content is actually needed.
func (ll *LazyLoader) LoadContent(ref *LazyReference) error {
	ll.mu.Lock()
	defer ll.mu.Unlock()

	if ref.Loaded {
		return nil // Already loaded
	}

	// Check file size limit (1MB max)
	if ref.Size > 1024*1024 {
		return fmt.Errorf("file too large for lazy loading: %d bytes", ref.Size)
	}

	// Read content
	data, err := os.ReadFile(ref.Path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	ref.Content = string(data)
	ref.Loaded = true
	ref.LineCount = strings.Count(ref.Content, "\n") + 1

	logger.InfoCF("context", "Loaded file content",
		map[string]any{
			"file_id": ref.ID,
			"path":    ref.Path,
			"size":    ref.Size,
			"lines":   ref.LineCount,
		})

	return nil
}

// UnloadContent unloads file content to free memory.
func (ll *LazyLoader) UnloadContent(ref *LazyReference) {
	ll.mu.Lock()
	defer ll.mu.Unlock()

	ref.Content = ""
	ref.Loaded = false
	ref.LineCount = 0

	logger.DebugCF("context", "Unloaded file content",
		map[string]any{
			"file_id": ref.ID,
		})
}

// GetReference returns a reference by ID.
func (ll *LazyLoader) GetReference(id string) (*LazyReference, bool) {
	ll.mu.RLock()
	defer ll.mu.RUnlock()

	for _, ref := range ll.cache {
		if ref.ID == id {
			return ref, true
		}
	}
	return nil, false
}

// evictOldest removes the oldest entry from cache. (Must be called with lock held)
func (ll *LazyLoader) evictOldest() {
	var oldestPath string
	var oldestTime time.Time
	first := true

	for path, ref := range ll.cache {
		if first || ref.ModTime.Before(oldestTime) {
			oldestPath = path
			oldestTime = ref.ModTime
			first = false
		}
	}

	if oldestPath != "" {
		delete(ll.cache, oldestPath)
		logger.DebugCF("context", "Evicted oldest cache entry",
			map[string]any{
				"path": oldestPath,
			})
	}
}

// generateFileID creates a unique, readable identifier for a file.
// Format: #FILE_BASENAME_HASH8
func generateFileID(path string) string {
	base := filepath.Base(path)
	if len(base) > 12 {
		base = base[:12]
	}

	hash := sha256.Sum256([]byte(path))
	hashStr := hex.EncodeToString(hash[:])[:8]

	return fmt.Sprintf("#FILE_%s_%s", strings.ToUpper(base), hashStr)
}

// FormatReference returns a formatted reference string for prompts.
// This is what should be included in LLM prompts instead of full content.
func FormatReference(ref *LazyReference) string {
	sizeStr := formatSize(ref.Size)
	return fmt.Sprintf(
		"📄 `%s` (%s, %d lines) - Use lazy_load(file_id=\"%s\") to read content",
		filepath.Base(ref.Path),
		sizeStr,
		ref.LineCount,
		ref.ID,
	)
}

// formatSize formats file size in human-readable format.
func formatSize(bytes int64) string {
	const (
		KB = 1024
		MB = 1024 * KB
	)

	switch {
	case bytes >= MB:
		return fmt.Sprintf("%.1fMB", float64(bytes)/MB)
	case bytes >= KB:
		return fmt.Sprintf("%.1fKB", float64(bytes)/KB)
	default:
		return fmt.Sprintf("%dB", bytes)
	}
}
