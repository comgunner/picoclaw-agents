// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package context

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLazyLoader_ReferenceFile(t *testing.T) {
	// Create temp directory and test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	err := os.WriteFile(testFile, []byte("Hello, World!"), 0o644)
	require.NoError(t, err)

	loader := NewLazyLoader(tmpDir)

	// Create reference
	ref, err := loader.ReferenceFile(testFile)
	require.NoError(t, err)
	require.NotNil(t, ref)

	// Verify reference
	assert.Equal(t, testFile, ref.Path)
	assert.False(t, ref.Loaded)
	assert.Empty(t, ref.Content)
	assert.Equal(t, int64(13), ref.Size)
	assert.GreaterOrEqual(t, ref.LineCount, 0)

	// Verify ID format
	assert.True(t, strings.HasPrefix(ref.ID, "#FILE_TEST.TXT_"))
	assert.Greater(t, len(ref.ID), 15) // #FILE_TEST.TXT_ + hash
}

func TestLazyLoader_CacheHit(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	err := os.WriteFile(testFile, []byte("Hello, World!"), 0o644)
	require.NoError(t, err)

	loader := NewLazyLoader(tmpDir)

	// First reference
	ref1, err := loader.ReferenceFile(testFile)
	require.NoError(t, err)

	// Second reference (should be cached)
	ref2, err := loader.ReferenceFile(testFile)
	require.NoError(t, err)

	// Should be same reference
	assert.Equal(t, ref1.ID, ref2.ID)
	assert.Same(t, ref1, ref2)
}

func TestLazyLoader_LoadContent(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	content := "Line 1\nLine 2\nLine 3\n"
	err := os.WriteFile(testFile, []byte(content), 0o644)
	require.NoError(t, err)

	loader := NewLazyLoader(tmpDir)

	// Create reference
	ref, err := loader.ReferenceFile(testFile)
	require.NoError(t, err)

	// Load content
	err = loader.LoadContent(ref)
	require.NoError(t, err)

	// Verify content loaded
	assert.True(t, ref.Loaded)
	assert.Equal(t, content, ref.Content)
	assert.GreaterOrEqual(t, ref.LineCount, 3)
}

func TestLazyLoader_UnloadContent(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	content := "Hello, World!"
	err := os.WriteFile(testFile, []byte(content), 0o644)
	require.NoError(t, err)

	loader := NewLazyLoader(tmpDir)

	// Create and load reference
	ref, err := loader.ReferenceFile(testFile)
	require.NoError(t, err)
	err = loader.LoadContent(ref)
	require.NoError(t, err)

	// Unload content
	loader.UnloadContent(ref)

	// Verify content unloaded
	assert.False(t, ref.Loaded)
	assert.Empty(t, ref.Content)
	assert.Equal(t, 0, ref.LineCount)
}

func TestLazyLoader_GetReference(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	err := os.WriteFile(testFile, []byte("Hello, World!"), 0o644)
	require.NoError(t, err)

	loader := NewLazyLoader(tmpDir)

	// Create reference
	ref, err := loader.ReferenceFile(testFile)
	require.NoError(t, err)

	// Get by ID
	foundRef, found := loader.GetReference(ref.ID)
	assert.True(t, found)
	assert.Same(t, ref, foundRef)

	// Get non-existent ID
	_, found = loader.GetReference("#FILE_NONEXISTENT_12345678")
	assert.False(t, found)
}

func TestLazyLoader_FileTooLarge(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "large.bin")

	// Create 2MB file (over 1MB limit)
	largeContent := make([]byte, 2*1024*1024)
	err := os.WriteFile(testFile, largeContent, 0o644)
	require.NoError(t, err)

	loader := NewLazyLoader(tmpDir)

	// Create reference
	ref, err := loader.ReferenceFile(testFile)
	require.NoError(t, err)

	// Try to load (should fail)
	err = loader.LoadContent(ref)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "file too large")
}

func TestLazyLoader_CacheEviction(t *testing.T) {
	tmpDir := t.TempDir()
	loader := NewLazyLoader(tmpDir)
	loader.maxCacheSize = 3 // Small cache for testing

	// Create test files
	files := make([]string, 5)
	for i := 0; i < 5; i++ {
		testFile := filepath.Join(tmpDir, "test"+string(rune('A'+i))+".txt")
		err := os.WriteFile(testFile, []byte("Content"), 0o644)
		require.NoError(t, err)
		files[i] = testFile

		// Small delay to ensure different mod times
		time.Sleep(10 * time.Millisecond)
	}

	// Add references
	for _, file := range files {
		_, err := loader.ReferenceFile(file)
		require.NoError(t, err)
	}

	// Cache should have max 3 entries
	assert.LessOrEqual(t, len(loader.cache), loader.maxCacheSize)
}

func TestLazyLoader_NonExistentFile(t *testing.T) {
	tmpDir := t.TempDir()
	loader := NewLazyLoader(tmpDir)

	// Try to reference non-existent file
	_, err := loader.ReferenceFile(filepath.Join(tmpDir, "nonexistent.txt"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to stat file")
}

func TestGenerateFileID(t *testing.T) {
	// Test ID generation
	path := "/home/user/test.txt"
	id := generateFileID(path)

	assert.True(t, strings.HasPrefix(id, "#FILE_TEST.TXT_"))
	assert.Greater(t, len(id), 15)

	// Same path should generate same ID
	id2 := generateFileID(path)
	assert.Equal(t, id, id2)

	// Different paths should generate different IDs
	id3 := generateFileID("/home/user/other.txt")
	assert.NotEqual(t, id, id3)
}

func TestGenerateFileID_LongBasename(t *testing.T) {
	// Test with long basename (>12 chars)
	path := "/home/user/verylongfilename.txt"
	id := generateFileID(path)

	// Should contain FILE prefix and hash
	assert.Contains(t, id, "#FILE_")
	assert.Greater(t, len(id), 15)
}

func TestFormatReference(t *testing.T) {
	ref := &LazyReference{
		ID:        "#FILE_TEST.TXT_12345678",
		Path:      "/home/user/test.txt",
		Size:      1536,
		LineCount: 25,
	}

	formatted := FormatReference(ref)

	assert.Contains(t, formatted, "📄")
	assert.Contains(t, formatted, "`test.txt`")
	assert.Contains(t, formatted, "1.5KB")
	assert.Contains(t, formatted, "25 lines")
	assert.Contains(t, formatted, "#FILE_TEST.TXT_12345678")
}

func TestFormatSize(t *testing.T) {
	tests := []struct {
		bytes  int64
		expect string
	}{
		{100, "100B"},
		{1024, "1.0KB"},
		{1536, "1.5KB"},
		{1048576, "1.0MB"},
		{1572864, "1.5MB"},
	}

	for _, tt := range tests {
		t.Run(tt.expect, func(t *testing.T) {
			result := formatSize(tt.bytes)
			assert.Equal(t, tt.expect, result)
		})
	}
}

func TestLazyLoader_MultilineFile(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "multiline.txt")

	// Create file with known line count
	lines := []string{
		"Line 1",
		"Line 2",
		"Line 3",
		"Line 4",
		"Line 5",
	}
	content := strings.Join(lines, "\n")
	err := os.WriteFile(testFile, []byte(content), 0o644)
	require.NoError(t, err)

	loader := NewLazyLoader(tmpDir)
	ref, err := loader.ReferenceFile(testFile)
	require.NoError(t, err)

	err = loader.LoadContent(ref)
	require.NoError(t, err)

	// Should count 5 lines (4 newlines + 1)
	assert.Equal(t, 5, ref.LineCount)
}

func TestLazyLoader_EmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "empty.txt")

	// Create empty file
	err := os.WriteFile(testFile, []byte(""), 0o644)
	require.NoError(t, err)

	loader := NewLazyLoader(tmpDir)
	ref, err := loader.ReferenceFile(testFile)
	require.NoError(t, err)

	err = loader.LoadContent(ref)
	require.NoError(t, err)

	assert.Empty(t, ref.Content)
	assert.Equal(t, 1, ref.LineCount) // Empty file = 1 line
}

func TestLazyLoader_LoadContentTwice(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	content := "Hello, World!"
	err := os.WriteFile(testFile, []byte(content), 0o644)
	require.NoError(t, err)

	loader := NewLazyLoader(tmpDir)
	ref, err := loader.ReferenceFile(testFile)
	require.NoError(t, err)

	// Load twice (second should be no-op)
	err = loader.LoadContent(ref)
	require.NoError(t, err)

	err = loader.LoadContent(ref)
	require.NoError(t, err)

	assert.True(t, ref.Loaded)
	assert.Equal(t, content, ref.Content)
}
