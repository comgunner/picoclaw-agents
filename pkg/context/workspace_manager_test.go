// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package context

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWorkspaceManager_NewWorkspaceManager(t *testing.T) {
	tmpDir := t.TempDir()
	wm := NewWorkspaceManager(tmpDir)

	assert.Equal(t, tmpDir, wm.BasePath)

	// Verify directory structure was created
	dirs := []string{"active", "memory", "cold", "temp", "sessions", "state", "scripts"}
	for _, dir := range dirs {
		fullPath := filepath.Join(tmpDir, dir)
		_, err := os.Stat(fullPath)
		assert.NoError(t, err, "Directory %s should exist", dir)
	}
}

func TestWorkspaceManager_StoreActive(t *testing.T) {
	tmpDir := t.TempDir()
	wm := NewWorkspaceManager(tmpDir)

	data := []byte("Active data content")
	filename := "test.txt"

	path, err := wm.StoreActive(data, filename)
	require.NoError(t, err)

	// Verify file was created
	expectedPath := filepath.Join(tmpDir, "active", filename)
	assert.Equal(t, expectedPath, path)

	// Verify content
	content, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, data, content)
}

func TestWorkspaceManager_StoreCold(t *testing.T) {
	tmpDir := t.TempDir()
	wm := NewWorkspaceManager(tmpDir)

	data := []byte("Cold storage data")
	filename := "archive.txt"

	path, err := wm.StoreCold(data, filename)
	require.NoError(t, err)

	// Verify file was created
	expectedPath := filepath.Join(tmpDir, "cold", filename)
	assert.Equal(t, expectedPath, path)

	// Verify content
	content, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, data, content)
}

func TestWorkspaceManager_StoreTemp(t *testing.T) {
	tmpDir := t.TempDir()
	wm := NewWorkspaceManager(tmpDir)

	data := []byte("Temporary data")
	filename := "temp.txt"

	path, err := wm.StoreTemp(data, filename)
	require.NoError(t, err)

	// Verify file was created
	expectedPath := filepath.Join(tmpDir, "temp", filename)
	assert.Equal(t, expectedPath, path)

	// Verify content
	content, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, data, content)
}

func TestWorkspaceManager_GetActivePath(t *testing.T) {
	tmpDir := t.TempDir()
	wm := NewWorkspaceManager(tmpDir)

	filename := "session.json"
	path := wm.GetActivePath(filename)

	expectedPath := filepath.Join(tmpDir, "active", filename)
	assert.Equal(t, expectedPath, path)
}

func TestWorkspaceManager_GetColdPath(t *testing.T) {
	tmpDir := t.TempDir()
	wm := NewWorkspaceManager(tmpDir)

	filename := "archive.zip"
	path := wm.GetColdPath(filename)

	expectedPath := filepath.Join(tmpDir, "cold", filename)
	assert.Equal(t, expectedPath, path)
}

func TestWorkspaceManager_EnsureStructureIdempotent(t *testing.T) {
	tmpDir := t.TempDir()

	// Create workspace manager (creates structure)
	wm1 := NewWorkspaceManager(tmpDir)
	assert.NotNil(t, wm1)

	// Create another workspace manager on same directory
	wm2 := NewWorkspaceManager(tmpDir)
	assert.NotNil(t, wm2)

	// Both should have same base path
	assert.Equal(t, wm1.BasePath, wm2.BasePath)

	// All directories should still exist
	dirs := []string{"active", "memory", "cold", "temp", "sessions", "state", "scripts"}
	for _, dir := range dirs {
		fullPath := filepath.Join(tmpDir, dir)
		_, err := os.Stat(fullPath)
		assert.NoError(t, err, "Directory %s should exist", dir)
	}
}

func TestWorkspaceManager_StoreMultipleFiles(t *testing.T) {
	tmpDir := t.TempDir()
	wm := NewWorkspaceManager(tmpDir)

	// Store multiple active files
	files := map[string]string{
		"file1.txt": "Content 1",
		"file2.txt": "Content 2",
		"file3.txt": "Content 3",
	}

	for filename, content := range files {
		path, err := wm.StoreActive([]byte(content), filename)
		require.NoError(t, err)

		// Verify content
		readContent, err := os.ReadFile(path)
		require.NoError(t, err)
		assert.Equal(t, content, string(readContent))
	}
}

func TestWorkspaceManager_StoreBinaryData(t *testing.T) {
	tmpDir := t.TempDir()
	wm := NewWorkspaceManager(tmpDir)

	// Store binary data
	binaryData := []byte{0x00, 0x01, 0x02, 0xFF, 0xFE, 0xFD}
	filename := "binary.bin"

	path, err := wm.StoreActive(binaryData, filename)
	require.NoError(t, err)

	// Verify binary content preserved
	content, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, binaryData, content)
}

func TestWorkspaceManager_StoreEmptyData(t *testing.T) {
	tmpDir := t.TempDir()
	wm := NewWorkspaceManager(tmpDir)

	// Store empty data
	filename := "empty.txt"
	path, err := wm.StoreActive([]byte(""), filename)
	require.NoError(t, err)

	// Verify file exists and is empty
	content, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Empty(t, content)
}

func TestWorkspaceManager_OverwriteFile(t *testing.T) {
	tmpDir := t.TempDir()
	wm := NewWorkspaceManager(tmpDir)

	filename := "overwrite.txt"

	// Store initial content
	_, err := wm.StoreActive([]byte("Initial content"), filename)
	require.NoError(t, err)

	// Overwrite with new content
	_, err = wm.StoreActive([]byte("New content"), filename)
	require.NoError(t, err)

	// Verify new content
	path := filepath.Join(tmpDir, "active", filename)
	content, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, "New content", string(content))
}

func TestWorkspaceManager_DirectoryPermissions(t *testing.T) {
	tmpDir := t.TempDir()
	_ = NewWorkspaceManager(tmpDir)

	// Verify directories have correct permissions (0755)
	dirs := []string{"active", "memory", "cold", "temp", "sessions", "state", "scripts"}
	for _, dir := range dirs {
		fullPath := filepath.Join(tmpDir, dir)
		info, err := os.Stat(fullPath)
		require.NoError(t, err)

		// Check if directory is readable and writable
		assert.True(t, info.IsDir())
		assert.Equal(t, os.FileMode(0o755), info.Mode().Perm())
	}
}

func TestWorkspaceManager_StoreInSubdirectories(t *testing.T) {
	tmpDir := t.TempDir()
	_ = NewWorkspaceManager(tmpDir)

	// Store file with subdirectory path
	filename := "subdir/nested/file.txt"
	data := []byte("Nested content")

	// This would fail because subdirectory doesn't exist
	// Testing that the API handles this gracefully
	assert.NotPanics(t, func() {
		wm := NewWorkspaceManager(tmpDir)
		_, err := wm.StoreActive(data, filename)
		assert.Error(t, err)
	})
}

func TestWorkspaceManager_ConcurrentAccess(t *testing.T) {
	tmpDir := t.TempDir()
	wm := NewWorkspaceManager(tmpDir)

	done := make(chan bool, 10)

	// Concurrent writes to different files
	for i := 0; i < 10; i++ {
		go func(idx int) {
			filename := "concurrent_" + string(rune('A'+idx)) + ".txt"
			data := []byte("Content from goroutine " + string(rune('A'+idx)))

			path, err := wm.StoreActive(data, filename)
			if err != nil {
				t.Errorf("Failed to store file: %v", err)
			}

			// Verify content
			content, readErr := os.ReadFile(path)
			if readErr != nil {
				t.Errorf("Failed to read file: %v", readErr)
			}

			if string(content) != string(data) {
				t.Errorf("Content mismatch")
			}

			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestWorkspaceManager_PathConstruction(t *testing.T) {
	tmpDir := t.TempDir()
	wm := NewWorkspaceManager(tmpDir)

	tests := []struct {
		name     string
		filename string
		verifier func(string) string
	}{
		{"Simple filename", "test.txt", wm.GetActivePath},
		{"Filename with spaces", "my file.txt", wm.GetActivePath},
		{"Filename with special chars", "file-@#$%.txt", wm.GetActivePath},
		{"JSON file", "data.json", wm.GetActivePath},
		{"Markdown file", "README.md", wm.GetActivePath},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.verifier(tt.filename)
			assert.NotEmpty(t, path)

			// Path should be absolute
			assert.True(t, filepath.IsAbs(path))

			// Path should contain the filename
			assert.True(t, filepath.Base(path) == tt.filename)
		})
	}
}
