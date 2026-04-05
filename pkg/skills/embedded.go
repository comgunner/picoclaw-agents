// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package skills

import (
	"embed"
	"io/fs"
)

//go:embed all:data
var embeddedSkillsFS embed.FS

// GetEmbeddedSkillsFS returns the embedded skills filesystem.
// Skills are organized as data/{category}/{skill-name}/SKILL.md
func GetEmbeddedSkillsFS() fs.FS {
	sub, _ := fs.Sub(embeddedSkillsFS, "data")
	return sub
}
