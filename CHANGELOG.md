# Changelog

All notable changes to the PicoClaw project will be documented in this file.

## [3.4.5] - 2026-03-23

### ✨ New Features
- **Autonomous Agent Runtime (LP-03)**: Introduced a background runtime for each agent that automatically processes internal messages. Agents no longer need to manually call `agent_receive` to check for tasks.
- **Runtime Manager**: A new coordination layer in `AgentLoop` that manages lifecycle and goroutines for all autonomous agents.
- **Enhanced Agent Autonomy**: Agents now automatically switch to `StatusBusy` when processing an internal task and can send auto-responses upon completion.
- **Extended Configuration**: Added `runtime` options to `AgentConfig` and `AgentDefaults` in `config.json`, allowing fine-grained control over which agents have autonomous capabilities enabled.

### 🛠️ Core Improvements
- **Message Bus Integration**: Added `GetChannel()` to `AgentMessageBus` to allow direct, non-blocking subscription to agent-specific inboxes.
- **Agent Instance Updates**: `AgentInstance` now tracks its own `Runtime` configuration for faster access during autonomous execution.

---

## [3.4.4] - 2026-03-12

### 🛡️ Security
- **Deny Patterns (MP-01)**: Added `DefaultDenyPatterns` to `pkg/tools/shell.go` with 12 patterns blocking dangerous commands (`rm -rf /`, `shutdown`, `dd if=`, fork bombs, disk writes, etc.). `NewExecToolWithConfig` now fails closed if deny patterns are empty. Warning no longer appears at startup.
- **Gemini/Antigravity Schema Fix**: Added `sanitizeSchemaForGemini()` to handle JSON Schema types incompatible with Google AI Platform. Replaces `"type": "any"` and invalid types with `"type": "object"`.

### 🐛 Bug Fixes
- **Model Naming (MP-02)**: Fixed auto-generated config from `picoclaw-agents auth login --provider google-antigravity` using incorrect model name `"gemini-flash"`. Now generates `"antigravity-gemini-3-flash"` consistently.
- **Tool Response Parsing**: Improved tool response parsing in Antigravity provider with better JSON handling and name resolution from call IDs.
- **TokenBudget Deadlock (Problema 9)**: Fixed agent blocking indefinitely when token budget exceeded 80%. Implemented Hard Limit (100%) in `CanAfford` and Soft Limit (80%) in `Charge` for preventive GC. Agent now self-recovers automatically.
- **Rehydration Diagnostic Loop (Problema 10)**: Fixed agent entering a prolonged tool-calling diagnostic loop after crash recovery. Added explicit suppressor in rehydration message to prevent LLM from invoking internal diagnostic tools before confirming availability to the user. Heartbeat stranded locks are now silently discarded instead of triggering full recovery flow.

### ✨ New Features
- **Clean Command (LP-02)**: New `picoclaw-agents clean` command to remove old or corrupt session files. Supports `--all`, `--older-than <duration>`, and `--dry-run` flags.
- **Native Tools Logging (LP-01)**: Added explicit startup log when the 5 native tools register (`system_diagnostics`, `config_manager`, `resource_monitor`, `memory_store`, `version_control`).

### 🧪 Tests
- **Antigravity Provider Tests (LP-03)**: Added `TestSanitizeSchemaForGemini_ReplacesAnyType`, `TestSanitizeSchemaForGemini_InvalidTypes`, and `TestBuildRequest_ToolResponse` in `pkg/providers/antigravity_provider_test.go`.
- **TokenBudget Tests**: New tests in `pkg/context/token_budget_test.go` verifying Hard/Soft Limit behavior and GC trigger.

### 📝 Documentation
- Added `docs/ANTIGRAVITY_QUICKSTART.md` — Quick start guide for Google Antigravity OAuth login.
- Updated `docs/ANTIGRAVITY_AUTH.md` with comprehensive troubleshooting section.
- Added `docs/ANTIGRAVITY_USAGE.md` with usage examples and config reference.

### ⚠️ Upgrade Notes
- If you logged in with `google-antigravity` before this release, update your `model_name` in `~/.picoclaw/config.json` from `"gemini-flash"` to `"antigravity-gemini-3-flash"`.
- Sessions created before the schema fix may be corrupt. Run `picoclaw-agents clean --all` to clear them.

---

## [3.4.3] - 2026-03-04

### 🛡️ Upstream Security Patch Adaptations

Adapted and applied 2 of 6 upstream patches from audit `upstream_audit_2026-03-04.json` (see `local_work/patch_execution_log_2026-03-04.md` for full details).

- **🔒 Registry Collision Warning** (`pkg/tools/registry.go`): Added structured warning via `logger.WarnCF` when `Register()` overwrites an existing tool by name. Critical for multi-agent environments where MCP servers per agent could silently contaminate each other's tool namespace. Upstream ref: [`a2591e0`](https://github.com/sipeed/picoclaw/commit/a2591e03a942ae244b50539d4b9d26da3a0b3d58)

- **📝 JSONL Memory Store** (`pkg/memory/jsonl.go` — *new file*): Introduced append-only JSONL session history store with atomic writes (temp→fsync→rename) to prevent file corruption under concurrent multi-agent writes. Sharded mutex design (`64 shards`) eliminates cross-agent lock contention. Adapted from upstream: [`6d894d6`](https://github.com/sipeed/picoclaw/commit/6d894d6138cb89a8bc714d69b03c9a6a14cb03d7) — `fileutil` dependency replaced by inlined `writeFileAtomic` for fork compatibility.

**Patches confirmed already present in fork (no action needed):**
- `web_fetch` `ForLLM` content pass-through fix (was already at `web.go:666`)
- HTTP retry `resp.Body` close on socket leak (already in `http_retry.go`)
- `state.go` atomic temp-rename saves (already implemented)
- Shell security deny patterns for `.env`/`id_rsa`/AWS credentials (already in `shell.go`)


## [3.4.2] - 2026-03-03

### ✨ Native Skills Architecture

- **🚀 Native Queue/Batch Skill**: Migrated `queue_batch` skill from external Markdown file to native Go code (`pkg/skills/queue_batch.go`). All documentation is now compiled into the binary, eliminating external file dependencies at runtime.
- **📦 Skills Loader Refactoring**: Updated `pkg/skills/loader.go` with native skills registry pattern. Added `GetQueueBatchSkill()`, `LoadNativeQueueBatchSkill()`, and `BuildNativeQueueBatchSummary()` methods.
- **🎯 Context Builder Integration**: Modified `pkg/agent/context.go` to use native skill injection via `LoadNativeQueueBatchSkill()` instead of hardcoded strings.
- **🧪 Comprehensive Test Suite**: Added `pkg/skills/queue_batch_test.go` with 9 test cases covering all public methods, concurrency, and workspace independence.
- **📚 Developer Documentation**: Created `local_work/crear_skill_interna.md` - complete guide for developing native skills with code templates and integration steps.
- **🌍 Documentation Updates**: Updated `docs/QUEUE_BATCH.en.md` and `docs/QUEUE_BATCH.es.md` with native skill architecture details and developer integration guide.

### 🔧 Technical Details

**Native Skill Pattern:**
```go
type QueueBatchSkill struct {
    workspace string
}

func (q *QueueBatchSkill) BuildSkillContext() string
func (q *QueueBatchSkill) BuildSummary() string
```

**Benefits:**
- Zero runtime dependencies on external `.md` files
- Enhanced security (skill cannot be tampered with)
- Automatic updates with binary releases
- Maximum performance (embedded documentation strings)

### 📝 Migration Notes

If you have custom integrations relying on `pkg/skills/queue_batch/SKILL.md`, update to use:
- `loader.LoadNativeQueueBatchSkill()` for full skill context
- `loader.BuildNativeQueueBatchSummary()` for XML summary

---

## [3.4.1] - 2026-03-02

### 🛡️ Security & Stability
- **🛡️ Native Skills Sentinel**: Implemented `skills_sentinel.go` as a native internal security tool. It provides proactive pattern-matching protection against prompt injection (input) and system leaks (output sanitization).
- **📝 Local Auditing**: Integrated a security auditor that records all blocked attacks and suspicious activities in `local_work/AUDIT.md`.

## [3.2.1] - 2026-03-01

## [3.2.0] - 2026-03-01

### 🛡️ Security & Stability
- **🔒 Fail-Close ExecTool**: Robust security policy. The command execution tool now performs strict validation of deny patterns during initialization. Invalid regex will prevent the agent from starting, eliminating "fail-open" vulnerabilities.
- **🚦 Robust Startup**: Improved `ChannelManager` checks. The system now error-outs early if no communication channels (Telegram, Discord, etc.) are enabled, preventing silent hangs.
- **🔄 Improved Agent Loop**: Enhanced `AgentLoop` with proactive context cancellation checks. Reduces log noise and ensures clean resource release during shutdown or bus disconnection.

### 🔧 Configuration & Agents
- **🤖 General Worker Agent**: Added a versatile `general_worker` to the default multi-agent suite for general-purpose tasks.
- **📄 Expanded Provider Templates**: `model_list` expanded to include comprehensive templates for OpenAI, Anthropic, DeepSeek, Google Gemini, Alibaba Qwen, Mistral, and more.
- **🧠 DeepSeek Default**: Standardized on `deepseek-chat` as the primary model across all default agents for optimal reasoning and cost efficiency.

### 📦 Dependencies
- **🖥️ TUI Foundation**: Added `tcell/v2` and `tview` dependencies to support the upcoming terminal management dashboard.

## [3.1.0] - 2026-02-27

### ✨ Core Features
- **🛡️ Task Lock System**: Implemented atomic `.lock` files for robust disaster recovery and concurrency control among subagents.
- **🔄 Boot Rehydration**: The Gateway will now automatically wake up and re-hydrate agents interrupted by system crashes or restarts.
- **🧠 Context Compactor**: Built-in intelligent context pruning and tool-output truncation. Safely elevated default `MaxTokens` to 32,768, permanently eliminating "Context Explosion" silent drops.
- **⚡️ Tool Mutual Exclusion**: `FileLockChecker` integration prevents concurrent agents from editing the same file simultaneously.
- **🤖 o3-mini Support**: Standardized on `o3-mini` for high-performance OpenAI tasks, including automatic `max_completion_tokens` handling.
- **🌍 Qwen Regional Fixes**: Documented and implemented support for Alibaba Cloud Virginia (US-East-1) regional endpoints.

## [3.0.0] - 2026-02-27

### ✨ Core Features
- **🚀 Advanced Multi-Agent Architecture**: Full support for isolated subagent sessions and the ability to execute different LLM models in parallel.
- **👥 The "Dream Team" Workflow**: New documented use case for a complete software development lifecycle, including `project_manager`, `senior_dev`, `qa_specialist`, and `junior_fixer` roles.
- **🧠 DeepSeek Standardization**: **DeepSeek** (`deepseek-chat` and `deepseek-reasoner`) is now established as the default model suite due to excellent reasoning and API efficiency.

### 📝 Documentation
- **🌍 Multilingual Support**: Updated and synchronized `README` across 7 languages (EN, ES, ZH, JA, FR, PT-BR, VI).
- **🛠 Installation Guides**: New detailed server installation guides for Ubuntu (`install_ubuntu_server.md`).
- **💡 Recommended Models**: New section with specific model recommendations for technical development tasks (`backend_coder`).

### 🔧 Configuration
- **📄 config_dev.example.json**: Created advanced config showcasing the potential of a multi-agent dev team.
- **📄 config.example.json**: Updated with new agent standards and payload cleanup.

### 🛡 Security & Maintenance
- **🔒 API Scrubbing**: Purged all real API keys from standard configurations, replacing them with safe placeholders.
- **🧹 Repository Cleanup**: Cleaned up the Git history, `.git` garbage, and temporary files (`.DS_Store`, bins) for a clean open-source release.
- **🤖 Telegram Fix**: Re-implemented the `isMessageAllowed` security check to ensure only authorized users can interact with the bot.

---
*PicoClaw: Ultra-Efficient AI in Go. $10 Hardware · 10MB RAM.*
