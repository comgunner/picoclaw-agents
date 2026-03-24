// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package security

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Auditor handles logging of security events to a local file.
type Auditor struct {
	auditPath string
}

// NewAuditor creates a new Auditor that logs to the specified workspace.
func NewAuditor(workspace string) *Auditor {
	return &Auditor{
		auditPath: filepath.Join(workspace, "local_work", "AUDIT.md"),
	}
}

// LogSecurityEvent records a blocked attack in the audit log.
func (a *Auditor) LogSecurityEvent(agentID, sessionKey, eventType, query, reason string) {
	// Ensure directory exists
	dir := filepath.Dir(a.auditPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		log.Printf("security auditor: failed to create audit directory %s: %v", dir, err)
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")

	// Truncate query for log readability
	displayQuery := query
	if len(displayQuery) > 100 {
		displayQuery = displayQuery[:97] + "..."
	}

	entry := fmt.Sprintf("## [%s] SECURITY_EVENT: BLOCKED\n\n"+
		"**Agent:** %s\n"+
		"**Session:** %s\n"+
		"**Event Type:** %s\n"+
		"**Query:** %q\n"+
		"**Reason:** %s\n\n",
		timestamp, agentID, sessionKey, eventType, displayQuery, reason)

	f, err := os.OpenFile(a.auditPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		log.Printf("security auditor: failed to open audit file %s: %v", a.auditPath, err)
		return
	}
	defer f.Close()

	if _, err := f.WriteString(entry); err != nil {
		log.Printf("security auditor: failed to write to audit file %s: %v", a.auditPath, err)
	}
}
