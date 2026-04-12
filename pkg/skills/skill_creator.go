// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT

package skills

import (
	"strings"
)

// SkillCreatorSkill implements a native skill that guides the agent through
// creating new file-based skills (SKILL.md) in workspace/skills/.
// Compiled directly into the binary — no external file dependencies.
//
// OUTPUT: Creates `~/.picoclaw/workspace/skills/{skill-name}/SKILL.md`
// NOT native Go skills (those are developer-only tasks).
type SkillCreatorSkill struct {
	workspace string
}

// NewSkillCreatorSkill creates a new SkillCreatorSkill instance.
func NewSkillCreatorSkill(workspace string) *SkillCreatorSkill {
	return &SkillCreatorSkill{
		workspace: workspace,
	}
}

// Name returns the skill identifier.
func (s *SkillCreatorSkill) Name() string {
	return "skill_creator"
}

// Description returns a concise trigger description.
func (s *SkillCreatorSkill) Description() string {
	return "Create new file-based skills (SKILL.md) in workspace/skills/. Use when user wants to add specialized knowledge, workflows, or tool integrations."
}

// GetInstructions returns the complete step-by-step workflow.
func (s *SkillCreatorSkill) GetInstructions() string {
	return skillCreatorInstructions
}

// GetAntiPatterns returns common anti-patterns.
func (s *SkillCreatorSkill) GetAntiPatterns() string {
	return skillCreatorAntiPatterns
}

// GetExamples returns concrete examples.
func (s *SkillCreatorSkill) GetExamples() string {
	return skillCreatorExamples
}

// BuildSkillContext returns the complete skill context for prompt injection.
func (s *SkillCreatorSkill) BuildSkillContext() string {
	parts := make([]string, 0, 13)
	parts = append(parts, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	parts = append(parts, "🧠 NATIVE SKILL: Skill Creator")
	parts = append(parts, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	parts = append(parts, "")
	parts = append(parts, "**PURPOSE:** Guide creation of file-based skills (SKILL.md) in `workspace/skills/`.")
	parts = append(parts, "**OUTPUT:** `~/.picoclaw/workspace/skills/{skill-name}/SKILL.md`")
	parts = append(parts, "**NOTE:** Creates SKILL.md skills only, NOT native Go skills.")
	parts = append(parts, "")
	parts = append(parts, s.GetInstructions())
	parts = append(parts, "")
	parts = append(parts, s.GetAntiPatterns())
	parts = append(parts, "")
	parts = append(parts, s.GetExamples())
	return strings.Join(parts, "\n")
}

// BuildSummary returns an XML summary for compact context.
func (s *SkillCreatorSkill) BuildSummary() string {
	return `<skill name="skill_creator" type="native">
  <purpose>Guide creation of file-based skills (SKILL.md)</purpose>
  <steps>Understand → Plan → Create → Write → Validate → Iterate</steps>
  <output>~/.picoclaw/workspace/skills/{name}/SKILL.md</output>
  <validation>Name regex, no Node.js, frontmatter check, structure</validation>
</skill>`
}

// ============================================================================
// DOCUMENTATION CONSTANTS
// ============================================================================

// NOTE: `bt` (backtick) is declared in common.go — reuse it here.

const skillCreatorInstructions = `## WHEN TO USE (CRITICAL)

Use this skill **AUTOMATICALLY** when:
- User says "create a skill for X" or "build a new skill"
- User wants to add specialized knowledge or workflows to the agent
- User wants to update an existing skill

**DO NOT** use for:
- Creating native Go skills (pkg/skills/*.go) — developer task
- Installing skills from ClawHub (use install_skill tool)
- General agent configuration

## STEP 1: UNDERSTAND

Ask concrete questions:
1. "What problem does this skill solve?"
2. "Can you give an example of how it would be used?"
3. "Does it need scripts (Python/Bash) or just documentation?"

Conclude when you have clarity.

## STEP 2: PLAN

Analyze what the skill needs:
- Scripts? (Python for APIs, Bash for automation)
- References? (schemas, API docs)
- Assets? (templates, boilerplate)

**Output structure:**
` + bt + bt + `
~/.picoclaw/workspace/skills/{skill-name}/
├── SKILL.md          (required)
├── scripts/          (optional — Python, Bash)
├── references/       (optional — detailed docs)
└── assets/           (optional — templates)
` + bt + bt + `

## STEP 3: CREATE DIRECTORY

Use write_file tool to create:
` + bt + bt + `
write_file("~/.picoclaw/workspace/skills/{name}/SKILL.md", content)
` + bt + bt + `

Create scripts/ and references/ subdirectories as needed.

## STEP 4: WRITE SKILL

### SKILL.md Template:
` + bt + bt + `yaml
---
name: {skill-name}
description: {clear description with WHEN to use triggers}
---

# {Skill Title}

Brief description.

## When to Use
- Case 1
- Case 2

## Usage

` + bt + `bash
command or script example
` + bt + `

## Scripts (if needed)

### script_name.py
` + bt + `python
#!/usr/bin/env python3
# script code
` + bt + `
` + bt + bt + `

**Writing rules:**
- SKILL.md body < 500 lines
- Imperative form ("Do X")
- Only non-obvious context
- NO Node.js, package.json, npm
- NO hardcoded secrets

## STEP 5: VALIDATE

- [ ] ` + bt + `~/.picoclaw/workspace/skills/{name}/SKILL.md` + bt + ` exists
- [ ] Name: lowercase, ` + bt + `^[a-z0-9]+(-[a-z0-9]+)*$` + bt + `, < 64 chars
- [ ] Description includes WHEN to use
- [ ] NO package.json or node_modules
- [ ] NO hardcoded secrets
- [ ] Scripts tested and working

## STEP 6: ITERATE

1. Test skill on real task
2. Notice confusion points
3. Update SKILL.md to clarify
4. Repeat
`

const skillCreatorAntiPatterns = `## ANTI-PATTERNS

### ❌ Node.js Skills
NEVER create skills with package.json, node_modules, or npm.
PicoClaw is <10MB. Node.js + deps = 50-150MB.
**Use:** Python (requests), Bash, or Go binary.

### ❌ Extra Documentation
DO NOT create: README.md, CHANGELOG.md, INSTALLATION_GUIDE.md
Only: SKILL.md + optional scripts/, references/, assets/

### ❌ Bloated SKILL.md
Keep < 500 lines. Move details to references/ files.

### ❌ Hardcoded Secrets
NEVER hardcode API keys, tokens, passwords.
Use env vars or config files.

### ❌ Invalid Names
Always validate: lowercase, no spaces, no path separators.
`

const skillCreatorExamples = `## EXAMPLES

### Example 1: Simple Documentation Skill

**User:** "Create a skill for PostgreSQL backup procedures"

**Agent:**
` + bt + bt + bt + `
Step 1: UNDERSTAND
- "What backup type? Full, incremental, WAL?"
- User: "Full backups with pg_dump"

Step 2: PLAN
- Only needs documentation, no scripts
- Create: SKILL.md with procedures

Step 3-4: CREATE
` + bt + bt + `
write_file("~/.picoclaw/workspace/skills/pg-backup/SKILL.md", """---
name: pg-backup
description: Perform PostgreSQL full backups using pg_dump. Use when user needs to backup database or schedule backups.
---

# PostgreSQL Backup

Full database backup procedures using pg_dump.

## When to Use
- User asks to backup database
- User wants to schedule backups
- User needs disaster recovery procedures

## Full Backup

` + bt + `bash
pg_dump -U postgres -d dbname -F c -f /backups/db_$(date +%Y%m%d).dump
` + bt + `

## Restore

` + bt + `bash
pg_restore -U postgres -d dbname /backups/db_20240101.dump
` + bt + `

## Verify Backup

` + bt + `bash
pg_restore -l /backups/db_20240101.dump
` + bt + `
""")
` + bt + bt + `

Step 5: VALIDATE
- Name: "pg-backup" ✅ (kebab-case, matches regex)
- SKILL.md: 30 lines ✅
- No secrets ✅

### Example 2: Skill with Python Script

**User:** "Create a skill to rotate PDF pages"

**Agent:**
` + bt + bt + bt + `
Step 1: UNDERSTAND
- Rotate PDF pages by 90/180/270 degrees
- Simple task, needs Python script

Step 2: PLAN
- SKILL.md + scripts/rotate_pdf.py

Step 3-4: CREATE
` + bt + bt + `
# Create SKILL.md
write_file("~/.picoclaw/workspace/skills/pdf-rotator/SKILL.md", """---
name: pdf-rotator
description: Rotate PDF pages by 90, 180, or 270 degrees. Use when user needs to fix PDF orientation.
---

# PDF Rotator

Rotate pages in PDF files.

## Usage

` + bt + `bash
python scripts/rotate_pdf.py input.pdf 90
` + bt + `

## Script

` + bt + `python
#!/usr/bin/env python3
from pypdf import PdfReader, PdfWriter
import sys

reader = PdfReader(sys.argv[1])
writer = PdfWriter()
for page in reader.pages:
    page.rotate(int(sys.argv[2]))
    writer.add_page(page)
with open("rotated.pdf", "wb") as f:
    writer.write(f)
` + bt + `
""")
` + bt + bt + `

Step 5: VALIDATE
- Test script with actual PDF ✅
- No Node.js ✅

### Example 3: Updating Existing Skill

**User:** "The pdf-rotator skill is confusing, simplify it"

**Agent:**
` + bt + bt + bt + `
Step 1: Read existing SKILL.md
Step 2: Identify confusion points
Step 3: Rewrite with clearer examples
Step 4: write_file to update
Step 5: Test with real PDF
` + bt + bt + bt + `

## QUICK REFERENCE

| Task | Files Needed | Complexity |
|------|-------------|------------|
| Documentation only | SKILL.md | Low |
| Simple automation | SKILL.md + scripts/ | Medium |
| Complex workflows | SKILL.md + scripts/ + references/ | High |
`
