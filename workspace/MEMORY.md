# MEMORY.md - PicoClaw-Agents Long-Term Memory

**Last Updated:** 2026-04-12 (12:30)
**Project:** PicoClaw-Agents (comgunner/picoclaw-agents)
**Maintainer:** @comgunner
**Version:** v1.2.4-20-g19419db-dirty

---

## 🧠 Native Skill: skill_creator (2026-04-12)

### Overview
- **File:** `pkg/skills/skill_creator.go` (340 líneas)
- **Tests:** `pkg/skills/skill_creator_test.go` (8 tests, ALL PASS)
- **Type:** Native compiled-in skill (Go constants)
- **Output:** Creates `~/.picoclaw/workspace/skills/{skill-name}/SKILL.md` (file-based skills)
- **NOT:** Does NOT create native Go skills (.go files in pkg/skills/)

### Registration
- **Registry:** `pkg/skills/loader.go` — added to `nativeSkillsRegistry` struct + `GetSkillCreatorSkill()` getter
- **Loader methods:** `LoadNativeSkillCreatorSkill()`, `BuildNativeSkillCreatorSummary()`
- **Context injection:** `pkg/agent/context.go` — included in `BuildSystemPrompt()` (skipped for Compact level)
- **WebUI:** Listed in `/api/skills` with `source: "native"` after fix to `listNativeSkills()`

### Bug Fixes During Implementation
1. **Race condition:** Documented (same pattern as 14 existing native skills — safe during single-threaded bootstrap)
2. **Token bloat:** Fixed by respecting `promptLevel != PromptLevelCompact`
3. **Description too long:** Reduced from ~200 to ~100 chars
4. **Output confusion:** Clarified skill creates SKILL.md files, NOT Go files
5. **WebUI not listing skill:** Fixed by adding entry to `listNativeSkills()` in loader.go

### Workflow (6 Steps)
1. **Understand** — Ask concrete questions about skill purpose
2. **Plan** — Analyze what scripts/references/assets are needed
3. **Create** — Create directory structure in workspace/skills/
4. **Write** — Write SKILL.md with frontmatter + instructions
5. **Validate** — Check name regex, no secrets, no Node.js, <500 lines
6. **Iterate** — Test on real tasks, improve based on feedback

---

## 🔒 Security Hardening Session (2026-04-12)

### Upstream Audit & Patch Adaptation
- **Audit source:** `sipeed/picoclaw` (original at `/Volumes/UPLOAD/.../picoclaw_original`)
- **Audit scope:** 790 commits since 2026-02-28
- **Issues found:** 12 (4 critical, 6 high, 2 medium)
- **Patches adapted:** 8 (all successfully applied and tested)

### Critical Fixes Applied

1. **GHSA-pv8c-p6jf-3fpp: Channel-based exec access control**
   - File: `pkg/tools/shell.go`
   - Added `allowRemote` field (default: false), `SetContext()`, `SetAllowRemote()`
   - Remote channels blocked from exec by default

2. **GHSA-pv8c-p6jf-3fpp: SSRF prevention in web_fetch**
   - File: `pkg/tools/web.go`
   - Added `isBlockedHostname()` — blocks localhost, .local, .internal, cloud metadata
   - Pre-flight validation in `WebFetchTool.Execute()`

3. **Workspace sandbox bypass (URL path confusion)**
   - File: `pkg/tools/shell.go`
   - 20-char context window (was 10), added ws://, wss:// schemes, safe paths map

4. **Cron tool remote access restriction**
   - File: `pkg/tools/cron.go`
   - Added `allowedChannels` map, `SetAllowedChannels()`, validation in `addJob()`

5. **File permissions hardening**
   - Files: `pkg/session/manager.go`, `pkg/state/state.go`
   - Dirs: 0755 → **0700**, Files: 0644 → **0600**

6. **Session key path traversal prevention**
   - File: `pkg/session/manager.go`
   - Sanitizes `/` and `\` in `GetOrCreate()`

7. **Disk wiping deny pattern accuracy**
   - File: `pkg/tools/shell.go`
   - Unified pattern: `(?i)\b(format|mkfs|diskpart)\b\s`

8. **Dependency updates**
   - `golang.org/x/sys`: 0.42.0 → **0.43.0**
   - `github.com/mymmrac/telego`: 1.6.0 → **1.8.0**
   - `modernc.org/sqlite`: 1.48.1 → **1.48.2**

### Build & Test Results
- **Binaries:** `build/picoclaw-agents-darwin-arm64` (~27MB), launcher (~27MB), TUI (~10MB)
- **Version:** v1.2.4-20-g19419db-dirty
- **Go:** go1.26.0
- **Tests:** 29 total, ALL PASS (pkg/tools: 13, pkg/session: 3, pkg/state: 5, pkg/skills: 8)

### auth.json Investigation
- **Report:** User saw blank auth.json after launching WebUI
- **Result:** File **intact** with all credentials (google-antigravity, openai, openrouter-free, qwen, zhipu)
- **Root cause:** UI display artifact, not file modification
- **Analysis:** Launcher does NOT write auth on startup; `SaveStore()` only called via `SetCredential()` which does `LoadStore()` first

---

## 📌 Project Context

### Fork Identity
- **Based on:** Sipeed/PicoClaw (original)
- **Fork URL:** https://github.com/comgunner/picoclaw-agents
- **Language:** Go 1.25.8+
- **Architecture:** Multi-agent with subagent spawning, task locks, boot rehydration
- **Custom tools:** image_gen_antigravity, social_post_bundle, text_script_create, community_manager_create_draft
- **Custom CLI flags:** --openai, --qwen, --glm, --openrouter, --free

### Native Skills (15 total)
1. queue_batch — Delegate heavy tasks to background queue
2. binance_mcp — Binance MCP server for trading
3. fullstack_developer — Full-stack development assistant
4. n8n_workflow — n8n Workflow Automation
5. agent_team_workflow — Multi-Agent Team Orchestrator
6. researcher — Deep Research Agent
7. backend_developer — Backend development expert
8. frontend_developer — Frontend development expert
9. devops_engineer — DevOps expert
10. security_engineer — Security expert
11. qa_engineer — QA expert
12. data_engineer — Data engineering expert
13. ml_engineer — ML/AI expert
14. odoo_developer — Odoo Architect
15. **skill_creator** — Create file-based skills (SKILL.md) ← NEW

### Config Location
- **Config:** `~/.picoclaw/config.json`
- **Auth:** `~/.picoclaw/auth.json` (OAuth tokens)
- **Security:** `~/.picoclaw/.security.yml` (exists but empty)
- **Workspace:** `~/.picoclaw/workspace/`

### Key Directories
- `pkg/tools/` — Native tools (shell, web, cron, image_gen, social media, etc.)
- `pkg/agent/` — Agent loop, context management, session handling
- `pkg/providers/` — LLM provider implementations
- `pkg/channels/` — Chat channel integrations
- `pkg/auth/` — OAuth credential storage and refresh
- `pkg/skills/` — Native compiled-in skills (15 skills)
- `web/backend/` — WebUI launcher and API
- `local_work/` — Personal work directory (not committed)

---

## 🔧 Useful Commands

```bash
# Build
make build                              # Main binary
make build-launcher                     # WebUI launcher
make build-launcher-tui                 # TUI launcher
make build-all                          # All platforms

# Run
./build/picoclaw-agents-darwin-arm64 agent          # Interactive agent
./build/picoclaw-agents-darwin-arm64 gateway         # Long-running bot
./build/picoclaw-agents-launcher-darwin-arm64 --public  # WebUI (network access)
./build/picoclaw-agents-launcher-tui-darwin-arm64     # TUI menu

# Test
go test ./pkg/tools/... -v
go test ./pkg/session/... -v
go test ./pkg/state/... -v
go test ./pkg/skills/... -v -run SkillCreator

# Auth
./build/picoclaw-agents-darwin-arm64 auth login --provider google-antigravity
./build/picoclaw-agents-darwin-arm64 auth login --provider openai --device-code
./build/picoclaw-agents-darwin-arm64 auth status
```

---

## ⚠️ Important Notes

1. **auth.json is safe** — Launcher does NOT modify it on startup. If UI shows empty credentials, it's a display issue.

2. **Exec tool restricted** — Remote channels (Telegram, Discord, etc.) cannot execute shell commands by default (GHSA patch).

3. **File permissions changed** — Session/state files now use 0700/0600 instead of 0755/0644.

4. **Dependencies updated** — telego 1.6.0→1.8.0 may have breaking API changes. Monitor Telegram behavior.

5. **skill_creator output** — Creates SKILL.md files in workspace/skills/, NOT Go files. Native Go skills are developer-only.

6. **WebUI skills list** — Native skills must be registered in `listNativeSkills()` in loader.go to appear in WebUI.

7. **Workspace directory** — `local_work/` is git-ignored. All audit docs and plans are stored there.

---

## 📋 Pending Work (Phase 2)

- [ ] `.security.yml` architecture — Separate credentials from config.json
- [ ] Sensitive data filtering — Prevent LLM credential exposure
- [ ] Process isolation — OS-level sandboxing for shell commands (2,266 lines)
- [ ] Integrate channel restrictions in agent loop (`execTool.SetContext()`, `cronTool.SetAllowedChannels()`)
- [ ] Monthly upstream security audits

---

*Memory maintained by @comgunner. Update after significant changes.*
