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
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// ImageGenRecord representa un registro de imagen generada
type ImageGenRecord struct {
	ID          string            `json:"id"`
	DateTime    string            `json:"fecha_hora"`
	Prompt      string            `json:"prompt"`
	Provider    string            `json:"provider"`
	ImagePath   string            `json:"image_path"`
	ScriptPath  string            `json:"script_path"`
	PromptPath  string            `json:"prompt_path"`
	AspectRatio string            `json:"aspect_ratio"`
	Model       string            `json:"model"`
	Language    string            `json:"language"` // "en" o "es"
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// ImageGenTracker controla el estado de imágenes generadas
type ImageGenTracker struct {
	mu          sync.RWMutex
	TrackerPath string                    `json:"-"`
	Records     map[string]ImageGenRecord `json:"records"`
}

// NewImageGenTracker crea un nuevo tracker
func NewImageGenTracker(trackerPath string) (*ImageGenTracker, error) {
	tracker := &ImageGenTracker{
		TrackerPath: trackerPath,
		Records:     make(map[string]ImageGenRecord),
	}

	// Cargar existente si existe
	if err := tracker.Load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	return tracker, nil
}

// Load carga el tracker desde archivo JSON
func (t *ImageGenTracker) Load() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	data, err := os.ReadFile(t.TrackerPath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, t)
}

// Save guarda el tracker a archivo JSON
func (t *ImageGenTracker) Save() error {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.saveLocked()
}

// saveLocked persists tracker data.
// Caller must hold either read or write lock.
func (t *ImageGenTracker) saveLocked() error {
	// Asegurar que el directorio existe
	dir := filepath.Dir(t.TrackerPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(t.TrackerPath, data, 0o644)
}

// Add agrega un nuevo registro
func (t *ImageGenTracker) Add(record ImageGenRecord) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.Records[record.ID] = record
	return t.saveLocked()
}

// Get obtiene un registro por ID
func (t *ImageGenTracker) Get(id string) (ImageGenRecord, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	record, ok := t.Records[id]
	return record, ok
}

// Update actualiza un registro existente
func (t *ImageGenTracker) Update(id string, record ImageGenRecord) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if _, ok := t.Records[id]; !ok {
		return fmt.Errorf("record not found: %s", id)
	}

	t.Records[id] = record
	return t.saveLocked()
}

// UpdateMetadata actualiza solo los metadatos de un registro
func (t *ImageGenTracker) UpdateMetadata(id string, key, value string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	record, ok := t.Records[id]
	if !ok {
		return fmt.Errorf("record not found: %s", id)
	}

	if record.Metadata == nil {
		record.Metadata = make(map[string]string)
	}
	record.Metadata[key] = value
	t.Records[id] = record
	return t.saveLocked()
}

// Exists verifica si un registro existe
func (t *ImageGenTracker) Exists(id string) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()

	_, ok := t.Records[id]
	return ok
}

// List retorna todos los registros
func (t *ImageGenTracker) List() []ImageGenRecord {
	t.mu.RLock()
	defer t.mu.RUnlock()

	records := make([]ImageGenRecord, 0, len(t.Records))
	for _, record := range t.Records {
		records = append(records, record)
	}
	return records
}

// ListByDate retorna registros de una fecha específica
func (t *ImageGenTracker) ListByDate(date string) []ImageGenRecord {
	t.mu.RLock()
	defer t.mu.RUnlock()

	records := make([]ImageGenRecord, 0)
	for _, record := range t.Records {
		if len(record.DateTime) >= 10 && record.DateTime[:10] == date {
			records = append(records, record)
		}
	}
	return records
}

// CountByDate cuenta registros de una fecha específica
func (t *ImageGenTracker) CountByDate(date string) int {
	t.mu.RLock()
	defer t.mu.RUnlock()

	count := 0
	for _, record := range t.Records {
		if len(record.DateTime) >= 10 && record.DateTime[:10] == date {
			count++
		}
	}
	return count
}

// Delete elimina un registro
func (t *ImageGenTracker) Delete(id string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if _, ok := t.Records[id]; !ok {
		return fmt.Errorf("record not found: %s", id)
	}

	delete(t.Records, id)
	return t.saveLocked()
}

// GenerateID genera un ID único basado en timestamp
func GenerateID() string {
	return time.Now().Format("20060102_150405") + "_" + randomString(6)
}

// randomString genera una string aleatorio para IDs únicos
func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
		time.Sleep(time.Nanosecond)
	}
	return string(b)
}
