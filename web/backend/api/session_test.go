package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/comgunner/picoclaw/pkg/config"
	"github.com/comgunner/picoclaw/pkg/memory"
	"github.com/comgunner/picoclaw/pkg/providers"
	"github.com/comgunner/picoclaw/pkg/session"
)

func sessionsTestDir(t *testing.T, configPath string) string {
	t.Helper()

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	workspace := cfg.Agents.Defaults.Workspace
	if len(workspace) > 0 && workspace[0] == '~' {
		home, _ := os.UserHomeDir()
		if len(workspace) > 1 && workspace[1] == '/' {
			workspace = home + workspace[1:]
		} else {
			workspace = home
		}
	}

	dir := filepath.Join(workspace, "sessions")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	return dir
}

func TestHandleListSessions_JSONLStorage(t *testing.T) {
	configPath, cleanup := setupOAuthTestEnv(t)
	defer cleanup()

	dir := sessionsTestDir(t, configPath)
	store, err := memory.NewJSONLStore(dir)
	if err != nil {
		t.Fatalf("NewJSONLStore() error = %v", err)
	}

	sessionKey := picoSessionPrefix + "history-jsonl"
	if err := store.AddFullMessage(nil, sessionKey, providers.Message{
		Role:    "user",
		Content: "Explain why the history API is empty after migration.",
	}); err != nil {
		t.Fatalf("AddFullMessage(user) error = %v", err)
	}
	if err := store.AddFullMessage(nil, sessionKey, providers.Message{
		Role:    "assistant",
		Content: "Because the API still reads only legacy JSON session files.",
	}); err != nil {
		t.Fatalf("AddFullMessage(assistant) error = %v", err)
	}
	if err := store.AddFullMessage(nil, sessionKey, providers.Message{
		Role:    "tool",
		Content: "ignored",
	}); err != nil {
		t.Fatalf("AddFullMessage(tool) error = %v", err)
	}
	if err := store.SetSummary(nil, sessionKey, "JSONL-backed session"); err != nil {
		t.Fatalf("SetSummary() error = %v", err)
	}

	h := NewHandler(configPath)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/sessions", nil)
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var items []sessionListItem
	if err := json.Unmarshal(rec.Body.Bytes(), &items); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("len(items) = %d, want 1", len(items))
	}
	if items[0].ID != "history-jsonl" {
		t.Fatalf("items[0].ID = %q, want %q", items[0].ID, "history-jsonl")
	}
	if items[0].MessageCount != 2 {
		t.Fatalf("items[0].MessageCount = %d, want 2", items[0].MessageCount)
	}
	if items[0].Title != "JSONL-backed session" {
		t.Fatalf("items[0].Title = %q, want %q", items[0].Title, "JSONL-backed session")
	}
	if items[0].Preview != "Explain why the history API is empty after migration." {
		t.Fatalf("items[0].Preview = %q", items[0].Preview)
	}
}

func TestHandleListSessions_TitleUsesTrimmedSummary(t *testing.T) {
	configPath, cleanup := setupOAuthTestEnv(t)
	defer cleanup()

	dir := sessionsTestDir(t, configPath)
	store, err := memory.NewJSONLStore(dir)
	if err != nil {
		t.Fatalf("NewJSONLStore() error = %v", err)
	}

	sessionKey := picoSessionPrefix + "summary-title"
	if err := store.AddFullMessage(nil, sessionKey, providers.Message{
		Role:    "user",
		Content: "fallback preview",
	}); err != nil {
		t.Fatalf("AddFullMessage() error = %v", err)
	}
	if err := store.SetSummary(
		nil,
		sessionKey,
		"  This summary is intentionally longer than sixty characters so it must be truncated in the history menu.  ",
	); err != nil {
		t.Fatalf("SetSummary() error = %v", err)
	}

	h := NewHandler(configPath)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/sessions", nil)
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var items []sessionListItem
	if err := json.Unmarshal(rec.Body.Bytes(), &items); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("len(items) = %d, want 1", len(items))
	}
	expectedTitle := truncateRunes(
		"This summary is intentionally longer than sixty characters so it must be truncated in the history menu.",
		maxSessionTitleRunes,
	)
	if items[0].Title != expectedTitle {
		t.Fatalf("items[0].Title = %q", items[0].Title)
	}
	if items[0].Preview != "fallback preview" {
		t.Fatalf("items[0].Preview = %q, want %q", items[0].Preview, "fallback preview")
	}
}

func TestHandleGetSession_JSONLStorage(t *testing.T) {
	configPath, cleanup := setupOAuthTestEnv(t)
	defer cleanup()

	dir := sessionsTestDir(t, configPath)
	store, err := memory.NewJSONLStore(dir)
	if err != nil {
		t.Fatalf("NewJSONLStore() error = %v", err)
	}

	sessionKey := picoSessionPrefix + "detail-jsonl"
	for _, msg := range []providers.Message{
		{Role: "user", Content: "first"},
		{Role: "assistant", Content: "second"},
		{Role: "tool", Content: "ignored"},
	} {
		if err := store.AddFullMessage(nil, sessionKey, msg); err != nil {
			t.Fatalf("AddFullMessage() error = %v", err)
		}
	}
	if err := store.SetSummary(nil, sessionKey, "detail summary"); err != nil {
		t.Fatalf("SetSummary() error = %v", err)
	}

	h := NewHandler(configPath)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/sessions/detail-jsonl", nil)
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var resp struct {
		ID       string `json:"id"`
		Summary  string `json:"summary"`
		Messages []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"messages"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if resp.ID != "detail-jsonl" {
		t.Fatalf("resp.ID = %q, want %q", resp.ID, "detail-jsonl")
	}
	if resp.Summary != "detail summary" {
		t.Fatalf("resp.Summary = %q, want %q", resp.Summary, "detail summary")
	}
	if len(resp.Messages) != 2 {
		t.Fatalf("len(resp.Messages) = %d, want 2", len(resp.Messages))
	}
	if resp.Messages[0].Role != "user" || resp.Messages[0].Content != "first" {
		t.Fatalf("first message = %#v, want user/first", resp.Messages[0])
	}
	if resp.Messages[1].Role != "assistant" || resp.Messages[1].Content != "second" {
		t.Fatalf("second message = %#v, want assistant/second", resp.Messages[1])
	}
}

func TestHandleDeleteSession_JSONLStorage(t *testing.T) {
	configPath, cleanup := setupOAuthTestEnv(t)
	defer cleanup()

	dir := sessionsTestDir(t, configPath)
	store, err := memory.NewJSONLStore(dir)
	if err != nil {
		t.Fatalf("NewJSONLStore() error = %v", err)
	}

	sessionKey := picoSessionPrefix + "delete-jsonl"
	if err := store.AddFullMessage(nil, sessionKey, providers.Message{
		Role:    "user",
		Content: "delete me",
	}); err != nil {
		t.Fatalf("AddFullMessage() error = %v", err)
	}
	if err := store.SetSummary(nil, sessionKey, "delete summary"); err != nil {
		t.Fatalf("SetSummary() error = %v", err)
	}

	h := NewHandler(configPath)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/sessions/delete-jsonl", nil)
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusNoContent, rec.Body.String())
	}

	base := filepath.Join(dir, sanitizeSessionKey(sessionKey))
	for _, path := range []string{base + ".jsonl", base + ".meta.json"} {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Fatalf("expected %s to be removed, stat err = %v", path, err)
		}
	}
}

func TestHandleGetSession_LegacyJSONFallback(t *testing.T) {
	configPath, cleanup := setupOAuthTestEnv(t)
	defer cleanup()

	dir := sessionsTestDir(t, configPath)
	manager := session.NewSessionManager(dir)
	sessionKey := picoSessionPrefix + "legacy-json"
	manager.AddMessage(sessionKey, "user", "legacy user")
	manager.AddMessage(sessionKey, "assistant", "legacy assistant")
	if err := manager.Save(sessionKey); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	h := NewHandler(configPath)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/sessions/legacy-json", nil)
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
}

func TestHandleSessions_FiltersEmptyJSONLFiles(t *testing.T) {
	configPath, cleanup := setupOAuthTestEnv(t)
	defer cleanup()

	dir := sessionsTestDir(t, configPath)
	base := filepath.Join(dir, sanitizeSessionKey(picoSessionPrefix+"empty-jsonl"))
	if err := os.WriteFile(base+".jsonl", []byte{}, 0o644); err != nil {
		t.Fatalf("WriteFile(jsonl) error = %v", err)
	}

	h := NewHandler(configPath)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	listRec := httptest.NewRecorder()
	listReq := httptest.NewRequest(http.MethodGet, "/api/sessions", nil)
	mux.ServeHTTP(listRec, listReq)

	if listRec.Code != http.StatusOK {
		t.Fatalf("list status = %d, want %d, body=%s", listRec.Code, http.StatusOK, listRec.Body.String())
	}

	var items []sessionListItem
	if err := json.Unmarshal(listRec.Body.Bytes(), &items); err != nil {
		t.Fatalf("Unmarshal(list) error = %v", err)
	}
	if len(items) != 0 {
		t.Fatalf("len(items) = %d, want 0", len(items))
	}

	detailRec := httptest.NewRecorder()
	detailReq := httptest.NewRequest(http.MethodGet, "/api/sessions/empty-jsonl", nil)
	mux.ServeHTTP(detailRec, detailReq)

	if detailRec.Code != http.StatusNotFound {
		t.Fatalf("detail status = %d, want %d, body=%s", detailRec.Code, http.StatusNotFound, detailRec.Body.String())
	}
}

// TestPicoSessionPrefix_MatchesNewFormat verifies the updated prefix matches
// the format generated by routing with DMScopePerChannelPeer for the pico channel.
func TestPicoSessionPrefix_MatchesNewFormat(t *testing.T) {
	sessionUUID := "550e8400-e29b-41d4-a716-446655440000"

	// Simulate what BuildAgentPeerSessionKey produces with:
	//   agentID="main", channel="pico", DMScope=per-channel-peer, peer.ID=sessionUUID
	expectedKey := "agent:main:pico:direct:" + sessionUUID
	expectedSanitized := "agent_main_pico_direct_" + sessionUUID

	if got := picoSessionPrefix + sessionUUID; got != expectedKey {
		t.Errorf("picoSessionPrefix + uuid = %q, want %q", got, expectedKey)
	}

	sanitized := strings.ReplaceAll(picoSessionPrefix+sessionUUID, ":", "_")
	if sanitized != expectedSanitized {
		t.Errorf("sanitized key = %q, want %q", sanitized, expectedSanitized)
	}

	if !strings.HasPrefix(sanitized, sanitizedPicoSessionPrefix) {
		t.Errorf("sanitized %q does not have prefix %q", sanitized, sanitizedPicoSessionPrefix)
	}
}

// TestExtractPicoSessionID_UUID verifies extraction of a UUID from the new key format.
func TestExtractPicoSessionID_UUID(t *testing.T) {
	uuid := "550e8400-e29b-41d4-a716-446655440000"
	fullKey := picoSessionPrefix + uuid

	got, ok := extractPicoSessionID(fullKey)
	if !ok {
		t.Fatalf("extractPicoSessionID(%q) ok=false, want true", fullKey)
	}
	if got != uuid {
		t.Errorf("extractPicoSessionID = %q, want %q", got, uuid)
	}
}

// TestExtractPicoSessionIDFromSanitizedKey_UUID verifies extraction from sanitized filename.
func TestExtractPicoSessionIDFromSanitizedKey_UUID(t *testing.T) {
	uuid := "550e8400-e29b-41d4-a716-446655440000"
	sanitizedKey := sanitizedPicoSessionPrefix + uuid

	got, ok := extractPicoSessionIDFromSanitizedKey(sanitizedKey)
	if !ok {
		t.Fatalf("extractPicoSessionIDFromSanitizedKey(%q) ok=false, want true", sanitizedKey)
	}
	if got != uuid {
		t.Errorf("got = %q, want %q", got, uuid)
	}
}

// TestExtractPicoSessionID_OldFormatExtractsWrongUUID shows that the old format
// with "pico:" infix will extract incorrectly (includes "pico:" prefix in UUID).
// This demonstrates why sessions from before the fix are incompatible.
func TestExtractPicoSessionID_OldFormatExtractsWrongUUID(t *testing.T) {
	oldStyleKey := "agent:main:pico:direct:pico:some-uuid"
	// The new prefix DOES match (it's a substring), but the extracted "UUID" is wrong
	got, ok := extractPicoSessionID(oldStyleKey)
	if !ok {
		t.Fatalf("extractPicoSessionID(%q) ok=false, want true (prefix matches)", oldStyleKey)
	}
	// The extracted value includes the unwanted "pico:" infix
	want := "pico:some-uuid"
	if got != want {
		t.Errorf("extracted UUID = %q, want %q (shows old format is incompatible)", got, want)
	}
}

// TestExtractPicoSessionID_NonMainAgent verifies extraction works for any agentID.
func TestExtractPicoSessionID_NonMainAgent(t *testing.T) {
	uuid := "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
	fullKey := "agent:engineering_manager:pico:direct:" + uuid

	got, ok := extractPicoSessionID(fullKey)
	if !ok {
		t.Fatalf("extractPicoSessionID(%q) ok=false, want true", fullKey)
	}
	if got != uuid {
		t.Errorf("extractPicoSessionID = %q, want %q", got, uuid)
	}
}

// TestExtractPicoSessionIDFromSanitizedKey_NonMainAgent verifies extraction from a
// sanitized filename where the agentID itself contains underscores.
func TestExtractPicoSessionIDFromSanitizedKey_NonMainAgent(t *testing.T) {
	uuid := "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
	// agentID = "engineering_manager" which contains an underscore
	sanitizedKey := "agent_engineering_manager_pico_direct_" + uuid

	got, ok := extractPicoSessionIDFromSanitizedKey(sanitizedKey)
	if !ok {
		t.Fatalf("extractPicoSessionIDFromSanitizedKey(%q) ok=false, want true", sanitizedKey)
	}
	if got != uuid {
		t.Errorf("got = %q, want %q", got, uuid)
	}
}

// TestHandleListSessions_NonMainAgentIDSession verifies that sessions written with a
// non-"main" agentID (e.g., "engineering_manager") appear in the session list.
func TestHandleListSessions_NonMainAgentIDSession(t *testing.T) {
	configPath, cleanup := setupOAuthTestEnv(t)
	defer cleanup()

	dir := sessionsTestDir(t, configPath)

	// Simulate what the agent writes when agentID = "engineering_manager"
	sessionUUID := "b2c3d4e5-f6a7-8901-bcde-f23456789012"
	sanitizedBase := "agent_engineering_manager_pico_direct_" + sessionUUID
	jsonlPath := filepath.Join(dir, sanitizedBase+".jsonl")

	msgs := []providers.Message{
		{Role: "user", Content: "Prueba con engineering_manager"},
		{Role: "assistant", Content: "Respuesta del agente"},
	}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	for _, msg := range msgs {
		_ = enc.Encode(msg)
	}
	if err := os.WriteFile(jsonlPath, buf.Bytes(), 0o644); err != nil {
		t.Fatalf("WriteFile(jsonl) error = %v", err)
	}

	h := NewHandler(configPath)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/sessions", nil)
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var items []sessionListItem
	if err := json.Unmarshal(rec.Body.Bytes(), &items); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("len(items) = %d, want 1 — session with engineering_manager agentID not found", len(items))
	}
	if items[0].ID != sessionUUID {
		t.Fatalf("items[0].ID = %q, want %q", items[0].ID, sessionUUID)
	}
	if items[0].MessageCount != 2 {
		t.Fatalf("items[0].MessageCount = %d, want 2", items[0].MessageCount)
	}
}

// TestHandleListSessions_UUIDSession verifies that a session stored with a UUID
// key (as generated by the fixed Pico channel routing) appears in the session list.
func TestHandleListSessions_UUIDSession(t *testing.T) {
	configPath, cleanup := setupOAuthTestEnv(t)
	defer cleanup()

	dir := sessionsTestDir(t, configPath)
	store, err := memory.NewJSONLStore(dir)
	if err != nil {
		t.Fatalf("NewJSONLStore() error = %v", err)
	}

	sessionUUID := "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
	sessionKey := picoSessionPrefix + sessionUUID

	messages := []providers.Message{
		{Role: "user", Content: "Hola, prueba de historial con UUID"},
		{Role: "assistant", Content: "Historial funcionando correctamente"},
	}
	for _, msg := range messages {
		if err := store.AddFullMessage(nil, sessionKey, msg); err != nil {
			t.Fatalf("AddFullMessage(%s) error = %v", msg.Role, err)
		}
	}
	if err := store.SetSummary(nil, sessionKey, "Prueba UUID"); err != nil {
		t.Fatalf("SetSummary() error = %v", err)
	}

	h := NewHandler(configPath)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/sessions", nil)
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var items []sessionListItem
	if err := json.Unmarshal(rec.Body.Bytes(), &items); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("len(items) = %d, want 1 — session UUID %s not found in list",
			len(items), sessionUUID)
	}
	if items[0].ID != sessionUUID {
		t.Fatalf("items[0].ID = %q, want %q", items[0].ID, sessionUUID)
	}
	if items[0].MessageCount != 2 {
		t.Fatalf("items[0].MessageCount = %d, want 2", items[0].MessageCount)
	}
	if items[0].Title != "Prueba UUID" {
		t.Fatalf("items[0].Title = %q, want %q", items[0].Title, "Prueba UUID")
	}
}
