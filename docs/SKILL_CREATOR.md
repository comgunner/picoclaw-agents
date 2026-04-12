# Skill Creator Tutorial

**Skill ID:** `skill_creator`  
**Type:** Native (compiled into binary)  
**Version:** 1.0.0  
**Last Updated:** April 12, 2026

---

## Overview

The **Skill Creator** (`skill_creator`) is a native skill built into PicoClaw-Agents that guides the agent through creating new **file-based skills** (SKILL.md) in your workspace. It follows a structured 6-step workflow to ensure every skill you create is well-designed, secure, and follows PicoClaw conventions.

### What It Does

The Skill Creator teaches the agent how to:

1. Understand what skill the user needs
2. Plan the skill structure and content
3. Create the directory and files in `~/.picoclaw/workspace/skills/`
4. Write a valid SKILL.md with proper YAML frontmatter
5. Validate the skill against naming and security rules
6. Iterate based on real-world usage

### What It Does NOT Do

- ❌ It does **not** create native Go skills (those are developer-only tasks requiring code changes and recompilation)
- ❌ It does **not** install skills from ClawHub (use the `install_skill` tool for that)
- ❌ It does **not** modify agent configuration

---

## Quick Start

Simply ask the agent to create a skill:

```
Create a skill for rotating PDF pages using Python.
```

The agent will automatically invoke the Skill Creator workflow and guide you through the process.

---

## The 6-Step Workflow

### Step 1: Understand

The agent asks concrete questions to understand the skill's purpose:

- "What problem does this skill solve?"
- "Can you give an example of how it would be used?"
- "Does it need scripts (Python/Bash) or just documentation?"

**Example interaction:**

```
User: "I need a skill to backup my database."

Agent: "Sure! Let me understand:
  1. What database system? (PostgreSQL, MySQL, etc.)
  2. Should it be a full or incremental backup?
  3. Any retention policy for old backups?"
```

### Step 2: Plan

The agent analyzes the use case and determines the skill structure:

| Need | Solution |
|------|----------|
| API calls or data processing | Python scripts in `scripts/` |
| Automation or shell commands | Bash scripts in `scripts/` |
| Reference material, schemas | Documentation in `references/` |
| Templates, boilerplate files | Assets in `assets/` |

**Decision tree:**

```
Does the skill need executable code?
├─ YES → Use scripts/ (Python preferred, Bash for simple tasks)
└─ NO  → Documentation-only skill (just SKILL.md)
```

### Step 3: Create

The agent creates the directory structure:

```
~/.picoclaw/workspace/skills/{skill-name}/
├── SKILL.md          ← Required: instructions + YAML frontmatter
├── scripts/          ← Optional: Python, Bash scripts
├── references/       ← Optional: detailed docs, schemas
└── assets/           ← Optional: templates, boilerplate
```

### Step 4: Write

#### SKILL.md Template

Every skill starts with a SKILL.md file containing YAML frontmatter and markdown instructions:

```yaml
---
name: my-skill-name
description: Clear description with usage triggers. Use when X, Y, or Z.
---
```

```markdown
# Skill Title

Brief description of what it does.

## When to Use
- Case 1: when the user asks to...
- Case 2: when you need to...

## Usage

```bash
python scripts/myscript.py arg1 arg2
```

## Scripts (if needed)

### myscript.py
Description of what the script does.

```python
#!/usr/bin/env python3
"""Script description."""
import sys

def main():
    pass

if __name__ == "__main__":
    main()
```
```

**Writing rules:**

| Rule | Description |
|------|-------------|
| Keep it under 500 lines | Move detailed content to `references/` if needed |
| Use imperative form | "Run the script" not "You should run the script" |
| Only add non-obvious context | The agent already knows common knowledge |
| No hardcoded secrets | Use environment variables or config files |

### Step 5: Validate

The skill creator enforces these rules:

| Check | Rule |
|-------|------|
| Name format | `^[a-z0-9]+(-[a-z0-9]+)*$` (lowercase, hyphens) |
| Name length | Less than 64 characters |
| Path separators | No `/`, `\`, or `..` in name |
| Node.js | ❌ No `package.json`, `node_modules`, or npm dependencies |
| Secrets | ❌ No hardcoded API keys, tokens, or passwords |
| SKILL.md size | Under 500 lines (move extras to `references/`) |
| Reference depth | Maximum 1 level deep from SKILL.md |

### Step 6: Iterate

After the skill is created, test it on a real task:

1. Ask the agent to use the skill
2. Notice any confusion or errors
3. Update the SKILL.md or scripts to clarify
4. Test again

---

## Complete Example: PDF Rotator Skill

### User Request

```
Create a skill to rotate PDF pages by 90, 180, or 270 degrees.
```

### Agent Workflow

**Step 1: Understand**
- User wants to rotate PDF pages
- Needs a simple Python script using `pypdf`
- No complex API or external service

**Step 2: Plan**
- SKILL.md with usage instructions
- `scripts/rotate_pdf.py` for the rotation logic

**Step 3: Create**
```bash
mkdir -p ~/.picoclaw/workspace/skills/pdf-rotator/scripts/
```

**Step 4: Write**

SKILL.md:
```yaml
---
name: pdf-rotator
description: Rotate PDF pages by 90, 180, or 270 degrees. Use when user needs to fix PDF orientation or rotate scanned documents.
---

# PDF Rotator

Rotate pages in PDF files using Python.

## When to Use
- User asks to rotate a PDF
- User needs to fix page orientation
- User has scanned documents upside-down

## Requirements

```bash
pip install pypdf
```

## Usage

```bash
python scripts/rotate_pdf.py input.pdf 90
```

Arguments:
- `input.pdf` — Path to the PDF file
- `90` — Rotation angle (90, 180, or 270)

Output: Creates `input_rotated.pdf` in the same directory.

## Script

### scripts/rotate_pdf.py

```python
#!/usr/bin/env python3
"""Rotate all pages in a PDF by the specified angle."""
import sys
from pypdf import PdfReader, PdfWriter

def rotate_pdf(input_path, angle):
    reader = PdfReader(input_path)
    writer = PdfWriter()

    for page in reader.pages:
        page.rotate(int(angle))
        writer.add_page(page)

    output_path = input_path.replace(".pdf", "_rotated.pdf")
    with open(output_path, "wb") as f:
        writer.write(f)

    print(f"Rotated {len(reader.pages)} pages by {angle}°")
    print(f"Saved to: {output_path}")
    return output_path

if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("Usage: rotate_pdf.py <input.pdf> <angle>")
        sys.exit(1)
    rotate_pdf(sys.argv[1], sys.argv[2])
```
```

**Step 5: Validate**
- ✅ Name: `pdf-rotator` (kebab-case, < 64 chars)
- ✅ Description includes usage triggers
- ✅ No secrets
- ✅ Python only, no Node.js
- ✅ Under 500 lines

**Step 6: Test**
```
User: "Rotate my contract.pdf by 90 degrees"
Agent: [Uses pdf-rotator skill to rotate the PDF]
```

---

## Anti-Patterns to Avoid

### ❌ Never Use Node.js

PicoClaw is designed to be lightweight (< 10MB RAM). Node.js + dependencies can easily reach 50-150MB.

**Bad:**
```
skills/my-skill/package.json
skills/my-skill/node_modules/
```

**Good:**
```
skills/my-skill/scripts/myscript.py     # Python
skills/my-skill/scripts/myscript.sh     # Bash
```

### ❌ No Extra Documentation

Do not create `README.md`, `CHANGELOG.md`, `INSTALLATION_GUIDE.md`, etc. These files are only for human developers — the agent only needs SKILL.md.

**Only create:**
- `SKILL.md` (required)
- `scripts/` (optional)
- `references/` (optional)
- `assets/` (optional)

### ❌ Don't Bloat SKILL.md

If your SKILL.md exceeds 500 lines, move detailed content to `references/`:

```
skills/my-skill/
├── SKILL.md              ← Quick-reference guide (< 500 lines)
└── references/
    ├── api_docs.md       ← Detailed API documentation
    ├── schemas.md        ← Database schemas
    └── workflows.md      ← Complex workflow guides
```

Reference from SKILL.md:
```markdown
## Advanced Configuration
See [references/api_docs.md](references/api_docs.md) for complete API reference.
```

### ❌ No Deep Reference Nesting

Maximum 1 level deep from SKILL.md:

```
✅ SKILL.md → references/api_docs.md          (OK)
❌ SKILL.md → references/api_docs.md → sub/   (BAD)
```

### ❌ No Hardcoded Secrets

**Bad:**
```python
API_KEY = "sk-abc123..."  # DON'T DO THIS
```

**Good:**
```python
import os
API_KEY = os.environ.get("MY_API_KEY")  # pragma: allowlist secret
```

---

## Naming Conventions

| Convention | Rule | Examples |
|-----------|------|----------|
| Case | Lowercase only | `pdf-rotator` ✅, `PDF-Rotator` ❌ |
| Separator | Hyphens or underscores | `pdf-rotator` ✅, `pdf_rotator` ✅ |
| Length | < 64 characters | `backup-postgres` ✅, `backup-postgres-database-daily-with-retention-policy` ❌ |
| Characters | `a-z`, `0-9`, `-`, `_` only | `my-skill-1` ✅, `my skill!` ❌ |
| No paths | No `/`, `\`, `..` | `pdf-rotator` ✅, `skills/pdf-rotator` ❌ |

---

## File-Based Skills vs Native Go Skills

| Aspect | File-Based Skills (skill_creator) | Native Go Skills |
|--------|-----------------------------------|-----------------|
| **Who creates** | Agent (via skill_creator) | Developer (writes Go code) |
| **Location** | `~/.picoclaw/workspace/skills/` | `pkg/skills/` in source code |
| **Compilation** | No — loaded at runtime | Yes — compiled into binary |
| **Format** | SKILL.md + optional scripts | Go structs + methods |
| **Use case** | Custom knowledge, workflows | Built-in roles, tool integrations |
| **Requires rebuild** | No | Yes |

**Rule of thumb:** Use the Skill Creator for user-level custom skills. Native Go skills are for framework-level features.

---

## Troubleshooting

### Skill not recognized

1. **Check the name:**
   ```bash
   ls ~/.picoclaw/workspace/skills/
   # Should show: your-skill-name/
   ```

2. **Verify SKILL.md exists:**
   ```bash
   cat ~/.picoclaw/workspace/skills/your-skill-name/SKILL.md
   ```

3. **Check YAML frontmatter:**
   ```yaml
   ---
   name: your-skill-name
   description: A description with usage triggers.
   ---
   ```

### Script won't run

1. **Make it executable:**
   ```bash
   chmod +x ~/.picoclaw/workspace/skills/your-skill-name/scripts/myscript.py
   ```

2. **Check dependencies:**
   ```bash
   pip install -r requirements.txt  # if applicable
   ```

3. **Test manually:**
   ```bash
   python ~/.picoclaw/workspace/skills/your-skill-name/scripts/myscript.py
   ```

---

## Related Documentation

- [ADDING_NATIVE_SKILLS.md](ADDING_NATIVE_SKILLS.md) — Developer guide for creating native Go skills
- [NATIVE_SKILLS_LIST.md](NATIVE_SKILLS_LIST.md) — Complete list of built-in native skills
- [SKILLS.md](SKILLS.md) — General skills documentation
- [SERVICE.md](SERVICE.md) — OS service management

---

*Skill Creator is built into the binary — no installation required. Just ask the agent to create a skill!*
