# Changelog

All notable changes to the PicoClaw project will be documented in this file.

> **Current Version:** v1.1.0 (as of March 2026)
>
> This changelog documents all changes by date. Version numbers in internal references (e.g., v3.4.2) refer to feature milestones, not release versions.

---

## 2026-03-26

### 📚 Documentation

#### **New MCP Builder Agent Documentation**
- `docs/MCP_BUILDER_AGENT.md` — Complete guide in English
- `docs/MCP_BUILDER_AGENT.es.md` — Complete guide in Spanish
- Examples added to `README.md` and `README.es.md`

**What's included:**
- What is MCP Builder Agent and when to use it
- 5 use cases with code examples (API integration, database access, workflow automation, file operations, custom tools)
- 3 activation methods (CLI, Chat, Config)
- 2 complete examples: GitHub MCP Server + PostgreSQL Database Server
- MCP server structure and tool anatomy
- Best practices (DO/DON'T) with code examples
- Full API reference (server.tool, server.resource, server.prompt)
- Return types and environment variables
- Links to official MCP resources

**Quick Example from README:**
```bash
# Invoke MCP Builder
picoclaw-agents agent -m "Build an MCP server for GitHub API"
```

#### **New Native Skills Complete List**
- `docs/NATIVE_SKILLS_LIST.md` — Complete list of all 14 native skills (English)
- `docs/NATIVE_SKILLS_LIST.es.md` — Lista completa de 14 skills nativas (Spanish)
- `docs/SKILLS.md` — Updated to  with 14 skills (English)
- `docs/SKILLS.es.md` — Guía actualizada  con 14 skills (Spanish)

**Native Skills Catalog ():**

**Engineering Role Skills (7):**
1. `backend_developer` — REST APIs, databases, microservices
2. `frontend_developer` — React, Vue, performance, accessibility
3. `devops_engineer` — CI/CD, Kubernetes, Terraform, monitoring
4. `security_engineer` — OWASP, penetration testing, compliance
5. `qa_engineer` — Test automation, coverage analysis, quality gates
6. `data_engineer` — ETL pipelines, data warehouses, streaming
7. `ml_engineer` — Model training, deployment, MLOps

**General Purpose Skills (4):**
8. `fullstack_developer` — Full-stack development, architecture
9. `researcher` — Deep research, web search, synthesis
10. `queue_batch` — Background task delegation, fire-and-forget
11. `agent_team_workflow` — Multi-agent orchestration

**Integration Skills (3):**
12. `binance_mcp` — Binance trading, market data
13. `n8n_workflow` — n8n automation, workflow creation
14. `odoo_developer` — Odoo architect, L10n-Mexico, CFDI 4.0

**Documentation Includes:**
- Detailed description of each skill
- Best practices and when to use
- Configuration examples
- Multi-skill combination patterns
- Troubleshooting guide
- Performance considerations (token usage)

---

## 2026-03-26

### ✨ New Features

#### **Comando `sandbox` — Workspaces aislados**
- `picoclaw-agents sandbox init [name]` — crea workspace aislado con permisos restrictivos (700)
- `picoclaw-agents sandbox status [name]` — verifica permisos y estructura del sandbox
- Subdirectorios automáticos: workspace, sessions, memory, state
- README generado automáticamente con instrucciones de uso

#### **Comando `util codegen` — Generador de código boilerplate**
- `picoclaw-agents util codegen --type <api|service|handler|model|config> --name <Name>`
- Genera código Go desde plantillas predefinidas
- Soporta 5 tipos: api, service, handler, model, config
- Integración con agente para generación automática vía tool `codegen`

### 🏗️ Internal

#### **New `pkg/tools/codegen.go`**
- `CodeGeneratorTool` — tool nativo para generación de código
- Plantillas para: API interfaces, services, HTTP handlers, models, configs
- Función `extractBaseName()` para parsear nombres compuestos (UserService → User)

#### **New `cmd/picoclaw/internal/sandbox/`**
- `command.go` — comandos `sandbox init` y `sandbox status`
- `command_test.go` — tests de creación y verificación

#### **Extended `cmd/picoclaw/internal/util/`**
- `codegen.go` — CLI wrapper para codegen tool

### 🧪 Tests
- `TestSandboxCommand_Runs`, `TestSandboxInitCommand_CreatesDirectory`
- Tests de subdirectorios y README

---

## 2026-03-26

### ✨ New Features

#### **Auth token monitor**
- `picoclaw-agents auth tokens` — lista estado de tokens OAuth cacheados
- `picoclaw-agents auth monitor --interval <min>` — monitoreo continuo de expiración
- Detección local de tokens: `valid` / `expiring_soon` / `expired` (sin HTTP en background)
- Umbbral de expiración: 5 minutos antes del expiry

### 🏗️ Internal

#### **New `pkg/auth/monitor.go`**
- `TokenMonitor` con goroutine de monitoreo configurable
- Lee `~/.picoclaw/auth.json`, sin llamadas HTTP automáticas
- `CheckInterval` público para configuración custom
- `CheckTokens()` — verifica expiración local
- `GetExpiringSoon()` — filtra tokens próximos a expirar

#### **New `pkg/auth/monitor_test.go`**
- `TestTokenMonitor_Start_Stop`, `TestTokenMonitor_ExpiringStatus`
- `TestTokenMonitor_GetExpiringSoon`, `TestTokenMonitor_Status_ReturnsCopy`

#### **Extended `cmd/picoclaw/internal/auth/`**
- `tokens.go` — subcomando `auth tokens` con output tabular
- `monitor.go` — subcomando `auth monitor` con watch en tiempo real
- `command.go` — registrados `newTokensCommand()` y `newMonitorCommand()`

### 🧪 Tests
- `TestTokenMonitor_*` — 6 tests nuevos

---

## 2026-03-26

### ✨ New Features

#### **Comando `config validate`**
- `picoclaw-agents config validate` — valida schema y valores semánticos de config.json
- Detecta: API keys faltantes, tokens inválidos, agent IDs duplicados, subagents mal configurados
- Acumula todos los errores (no fail-fast) para mejor UX
- Límite de output: 20 errores máx (evita overflow en configs grandes)

#### **Secret masking en wizard**
- Input de API keys en `onboard` ya no muestra caracteres en pantalla
- Funciona con `golang.org/x/term.ReadPassword`
- Fallback automático a texto plano en entornos no-TTY (CI, pipes)
- Integrado en `wizard.promptSecret()`

### 🏗️ Internal

#### **New `pkg/cli/input.go`**
- `ReadMasked(prompt)` — input sin eco para terminals interactivos
- `ReadMaskedWithFallback(prompt, scanner)` — compatible con tests/pipes
- `ReadLine(scanner)` — lectura normal de líneas
- `Confirm(prompt, scanner)` — confirmaciones y/n

#### **New `pkg/cli/input_test.go`**
- `TestReadLine_*`, `TestConfirm_*`, `TestReadMaskedWithFallback_*`

#### **New `pkg/config/validator.go`**
- `Validator.Validate(cfg)` — validación semántica de struct Config
- `Validator.ValidateFile(path)` — lee y valida archivo config.json
- Validaciones: model_list uniqueness, agent IDs únicos, Telegram token format, Binance keys pares
- `ValidationErrorList` — error acumulativo con formato legible

#### **New `pkg/config/validator_test.go`**
- `TestValidator_ValidConfig`, `TestValidator_MissingAPIKey`
- `TestValidator_InvalidTelegramToken`, `TestValidator_DuplicateAgentID`
- `TestValidator_SubagentsMaxSpawnDepth`, `TestValidator_BinancePartialKeys`

#### **New `cmd/picoclaw/internal/config/`**
- `command.go` — `NewConfigCommand()`
- `validate.go` — `newValidateCommand()` con flag `--config`

#### **Extended `cmd/picoclaw/internal/onboard/wizard.go`**
- Import `pkg/cli`
- `promptSecret()` usa `cli.ReadMaskedWithFallback()`

### 🧪 Tests
- `TestValidator_*` — 13 tests nuevos
- `TestReadMaskedWithFallback_NonTTY`, `TestConfirm_*`

### ⚠️ Upgrade Notes
- **Nueva dependencia:** `golang.org/x/term` (agregada a go.mod)

---

## 2026-03-26

### ✨ New Features

#### **Comando `doctor` — diagnóstico de entorno**
- `picoclaw-agents doctor` — verifica Go, Docker, workspace, WSL y seguridad
- Flag `--json` para output estructurado
- Detección de WSL (Windows Subsystem for Linux)
- Security checks: root, binaries peligrosas en PATH, puertos abiertos

**Output ejemplo:**
```
=== PicoClaw-Agents Doctor ===

System:
  OS/Arch:  darwin/arm64
  WSL:      false
  Shell:    /bin/zsh

Requirements:
  Go:         go version go1.26.0 darwin/arm64 [OK]
  Docker:     installed (not running)
  Workspace:  /Users/user/.picoclaw [OK]
  Config:     exists

Security:
  Root:        false (good)
  Dangerous:   nc
  Open ports:  none

✓ Environment ready!
```

### 🏗️ Internal

#### **Extended `pkg/setup/checker.go`**
- Nuevo campo `WSL bool` en `EnvironmentReport`
- Nuevo struct `SecurityReport` con `RunningAsRoot`, `DangerousBinaries`, `OpenPorts`
- Funciones: `detectWSL()`, `runSecurityChecks()`, `isPortOpen()`
- `detectWSL()` solo activo en Linux (`runtime.GOOS == "linux"`)
- `isPortOpen()` usa `net.DialTimeout` con 500ms timeout

#### **Extended `pkg/setup/checker_test.go`**
- `TestDetectWSL_NonLinux` — en macOS/Windows siempre false
- `TestSecurityChecks_NotRoot` — euid != 0 en tests normales
- `TestSecurityChecks_Ports` — verifica lista de puertos
- `TestSecurityChecks_DangerousBinaries` — verifica binaries detectados
- `TestEnvironmentReport_WithSecurity` — String() incluye security section
- `TestEnvironmentReport_WithoutSecurity` — String() omite section si está limpio

#### **New `cmd/picoclaw/internal/doctor/`**
- `command.go` — `NewDoctorCommand()` registrado en root
- `command_test.go` — 9 tests: `TestDoctorCommand_*`, `TestRunDoctor_*`

### 🧪 Tests
- `TestDetectWSL_NonLinux`, `TestSecurityChecks_NotRoot`, `TestDoctorCommand_Runs`
- Todos los tests pasan: 14 en `pkg/setup/...`, 9 en `cmd/.../doctor/...`

---

## 2026-03-26

### ✨ New Features

#### **158 Embedded Skills via `//go:embed`** (Fases 1-5)
- **158 new skills** embebidos en el binario usando `//go:embed all:data`
- Skills organizados por categoría: academic, design, engineering, game-development, marketing, paid-media, product, project-management, sales, spatial-computing, specialized, support, testing
- Binario self-contained — sin archivos externos, sin instalación adicional
- Aumento de binario: ~750 KB (de 19 MB a ~20 MB) — dentro del límite de 50 MB ✅

**Skills por categoría:**
- **academic** (5): anthropologist, geographer, historian, narratologist, psychologist
- **design** (8): brand-guardian, image-prompt-engineer, inclusive-visuals-specialist, ui-designer, ux-architect, ux-researcher, visual-storyteller, whimsy-injector
- **engineering** (23 nuevos): ai-engineer, backend-architect, code-reviewer, database-optimizer, devops-automator, embedded-firmware-engineer, feishu-integration-developer, git-workflow-master, incident-response-commander, mobile-app-builder, rapid-prototyper, senior-developer, software-architect, solidity-smart-contract-engineer, sre, technical-writer, threat-detection-engineer, wechat-mini-program-developer, etc.
- **game-development** (20): blender-addon-engineer, game-audio-engineer, game-designer, godot-gameplay-scripter, godot-multiplayer-engineer, godot-shader-developer, level-designer, narrative-designer, roblox-avatar-creator, roblox-experience-designer, roblox-systems-scripter, technical-artist, unity-architect, unity-editor-tool-developer, unity-multiplayer-engineer, unity-shader-graph-artist, unreal-multiplayer-architect, unreal-systems-engineer, unreal-technical-artist, unreal-world-builder
- **marketing** (27): ai-citation-strategist, app-store-optimizer, baidu-seo-specialist, bilibili-content-strategist, book-co-author, carousel-growth-engine, china-ecommerce-operator, content-creator, cross-border-ecommerce, douyin-strategist, growth-hacker, instagram-curator, kuaishou-strategist, linkedin-content-creator, livestream-commerce-coach, podcast-strategist, private-domain-operator, reddit-community-builder, seo-specialist, short-video-editing-coach, social-media-strategist, tiktok-strategist, twitter-engager, wechat-official-account, weibo-strategist, xiaohongshu-specialist, zhihu-strategist
- **paid-media** (7): auditor, creative-strategist, paid-social-strategist, ppc-strategist, programmatic-buyer, search-query-analyst, tracking-specialist
- **product** (5): behavioral-nudge-engine, feedback-synthesizer, manager, sprint-prioritizer, trend-researcher
- **project-management** (6): experiment-tracker, jira-workflow-steward, project-shepherd, studio-operations, studio-producer, project-manager-senior
- **sales** (8): account-strategist, coach, deal-strategist, discovery-coach, engineer, outbound-strategist, pipeline-analyst, proposal-strategist
- **spatial-computing** (6): macos-spatial-metal-engineer, terminal-integration-specialist, visionos-spatial-engineer, xr-cockpit-interaction-specialist, xr-immersive-developer, xr-interface-architect
- **specialized** (27): accounts-payable-agent, agentic-identity-trust, agents-orchestrator, automation-governance-architect, blockchain-security-auditor, compliance-auditor, corporate-training-designer, data-consolidation-agent, government-digital-presales-consultant, healthcare-marketing-compliance, identity-graph-operator, lsp-index-engineer, recruitment-specialist, report-distribution-agent, sales-data-extraction-agent, cultural-intelligence-strategist, developer-advocate, document-generator, french-consulting-market, korean-business-navigator, mcp-builder, model-qa, salesforce-architect, workflow-architect, study-abroad-advisor, supply-chain-strategist, zk-steward
- **support** (6): analytics-reporter, executive-summary-generator, finance-tracker, infrastructure-maintainer, legal-compliance-checker, support-responder
- **testing** (8): accessibility-auditor, api-tester, evidence-collector, performance-benchmarker, reality-checker, test-results-analyzer, tool-evaluator, workflow-optimizer

**Skills omitidos (ya existen como native Go):**
- backend_developer, frontend_developer, devops_engineer, security_engineer, qa_engineer, data_engineer, ml_engineer

### 🏗️ Internal

#### **New `pkg/skills/embedded.go`**
- `//go:embed all:data` directive for embedding skills filesystem
- `GetEmbeddedSkillsFS()` function to access embedded FS
- Skills organized as `data/{category}/{skill-name}/SKILL.md`

#### **Extended `pkg/skills/loader.go`**
- Added `embeddedFS fs.FS` field to `SkillsLoader` struct
- Auto-initialized in `NewSkillsLoader()` constructor
- Extended `ListSkills()` to include embedded skills (lowest priority)
- Extended `LoadSkill()` with fallback to embedded FS
- New `addEmbeddedSkills()` helper function
- New `parseSkillFrontmatter()` function for YAML frontmatter parsing
- Priority order: workspace > global > builtin > embedded

#### **Conversion Script**
- `cmd/tools/convert_skills/main.go` — tool to convert skills from `local_work/skills_import/` to embedded format
- Generates frontmatter with `name`, `description`, `category`, `version`
- Strips metadata headers from source files
- Outputs to `pkg/skills/data/{category}/{skill-name}/SKILL.md`
- Skips 7 native Go skills and excluded categories (examples, strategy)

### 🧪 Tests

#### **New `pkg/skills/embedded_skills_test.go`**
- `TestEmbeddedSkillsCount` — verifies ≥150 embedded skills loaded
- `TestEmbeddedSkillLoad` — loads specific skill and verifies content
- `TestEmbeddedSkillsListIncludes` — verifies expected skills present
- `TestEmbeddedSkillsNotDuplicated` — native skills not duplicated
- `TestEmbeddedSkillCategories` — verifies categories present
- `TestEmbeddedSkillContent` — verifies frontmatter stripped correctly
- `TestEmbeddedSkillPriority` — native skills have priority
- `TestEmbeddedSkillMetadata` — metadata parsed correctly
- `TestEmbeddedSkillsBuildSummary` — summary includes embedded skills

### 📝 Documentation

#### **Updated `local_work/plan_integracion_160skills_nativos.md`**
- Complete implementation plan for 178 skills via `//go:embed`
- Architecture decisions and comparisons
- File structure and format specifications
- Code change requirements
- Phase checklist and risk mitigations

### ⚠️ Upgrade Notes

- **No breaking changes**: All existing configurations remain compatible
- **Native skills have priority**: Existing native Go skills take precedence over embedded versions
- **No configuration needed**: Embedded skills auto-load on startup
- **Binary size**: ~750 KB increase (well under 50 MB limit)

---

## 2026-03-26

### ✨ New Features

#### **Bug Fix: Researcher Skill Registration** (Fase 0)
- **Fixed**: `researcher` skill existed in `pkg/skills/researcher.go` but was not registered
- **Added**: Registration in `nativeSkillsRegistry` struct in `pkg/skills/loader.go`
- **Added**: `GetResearcherSkill()`, `LoadNativeResearcherSkill()`, `BuildNativeResearcherSummary()` methods
- **Added**: Entry in `listNativeSkills()` for researcher skill
- **Result**: Researcher skill now available for use in `config.json`

#### **Security: Secret Scanner + Log Sanitizer** (Fase 1)
- **New `pkg/security/scanner.go`**: Static analysis scanner for hardcoded secrets
  - Detects 12 secret types: OpenAI, Anthropic, Google API, GitHub tokens, AWS keys, Slack tokens, Stripe secrets, Telegram bot tokens, DeepSeek keys, JWTs
  - `ScanFile()` and `ScanDir()` methods for file/directory scanning
  - Placeholder detection to avoid false positives on example files
  - Text file filtering (skips binaries, `.git`, `vendor`, `node_modules`)
  
- **New `pkg/security/sanitizer.go`**: Explicit sanitization function
  - `Sanitize(s string) string`: Redacts secrets from arbitrary strings
  - `SanitizeMap(m map[string]any) map[string]any`: Recursive map sanitization
  - Use cases: Tool results, user messages, external API responses
  - Format: `[REDACTED_pattern_name]`

- **Test Suite**: 
  - `pkg/security/scanner_test.go`: 13 tests for scanner functionality
  - `pkg/security/sanitizer_test.go`: 21 tests for sanitization

#### **Native Engineering Role Skills** (Fase 4-5)
Added 7 new native skills for specialized engineering roles, compiled directly into the binary:

- **`backend_developer`**: Backend development expert — REST APIs, databases, microservices, performance, security
- **`frontend_developer`**: Frontend development expert — React, Vue, performance, accessibility, design systems
- **`devops_engineer`**: DevOps expert — CI/CD pipelines, containers, infrastructure as code, monitoring, SRE
- **`security_engineer`**: Security expert — OWASP, penetration testing, hardening, threat modeling, compliance
- **`qa_engineer`**: QA expert — testing strategies, test automation, coverage analysis, quality gates
- **`data_engineer`**: Data engineering expert — ETL pipelines, data warehouses, streaming, data quality
- **`ml_engineer`**: ML/AI expert — model training, deployment, evaluation pipelines, MLOps, feature engineering

Each skill includes:
- Comprehensive role instructions and responsibilities
- Technology stack recommendations
- Best practices and quality checklists
- Anti-patterns to avoid (with code examples)
- Concrete usage examples (with code snippets)

#### **Environment Checker Package** (Fase 2)
- **New `pkg/setup/checker.go`**: Standalone environment validation package
- **`EnvironmentReport` struct**: OS, Arch, Go version, Docker status, workspace validation
- **`CheckEnvironment()` function**: Complete environment diagnostics
- **`IsReady()` method**: Validates minimum requirements (Go + workspace)
- **`String()` method**: Tabular output for terminal display
- Useful for `picoclaw-agents doctor` command and onboarding wizard

#### **Skills Import System** (Fase 3)
- **Python import script**: `local_work/scripts/import_skills_from_agency.py`
- Generates Markdown source files for skill conversion
- Output: `local_work/skills_import/engineering/*.md`
- Automated skill documentation generation

### 🛠️ Core Improvements

#### **Skills Loader Enhancement**
- **Expanded `nativeSkillsRegistry`**: Now holds 13 native skills (was 6)
- **14 new getter methods**: `GetBackendDeveloperSkill()`, `GetFrontendDeveloperSkill()`, etc.
- **14 new loader methods**: `LoadNativeBackendDeveloperSkill()`, `BuildNativeBackendDeveloperSummary()`, etc.
- **Thread-safe lazy initialization**: All skills use singleton pattern

#### **Configuration Examples** (Fase 6)
- **Updated `config/config.example.json`**: Added `_examples` section
- **Single specialized agent example**: Backend developer with skills
- **Orchestrator + subagents example**: Tech lead coordinating 7 engineering specialists
- Demonstrates multi-agent architecture with role-based skills

### 🧪 Tests

#### **Comprehensive Test Suite** (Fase 7)
- **New `pkg/skills/engineering_skills_test.go`**: 45+ test cases
- **Individual skill tests**: Name, Description, Instructions, Context, Summary for each of 7 skills
- **Integration tests**: Consistent structure across all engineering skills
- **Anti-patterns tests**: Verify all skills contain anti-pattern documentation
- **Examples tests**: Verify all skills contain usage examples
- **Structure tests**: Verify XML format and required sections

**Test Coverage:**
```
=== RUN   TestAllEngineeringSkillsHaveConsistentStructure
=== RUN   TestAllEngineeringSkillsHaveAntiPatterns
=== RUN   TestAllEngineeringSkillsHaveExamples
=== RUN   TestEngineeringSkillContextsContainRequiredSections
--- PASS: All tests (100% pass rate)
```

### 📝 Documentation

#### **New Documentation Files** (Fase 8)
- **`docs/SKILLS.md`**: Comprehensive comparison of Skills vs Tools
  - When to use Skills (role injection) vs Tools (action execution)
  - Complete table of all 13 native skills
  - Usage examples and configuration patterns

- **`docs/ADDING_NATIVE_SKILLS.md`**: Developer guide for contributing new skills
  - Step-by-step template for creating native skills
  - Interface requirements and method signatures
  - Registration process in `loader.go`
  - Testing requirements and examples

#### **Updated Documentation**
- **`CHANGELOG.md`**: This file, with complete v3.6.0 release notes
- **`config/config.example.json`**: Added extensive examples section

### 📦 Skills Import Report

Generated skills from agency-agents repository:
```json
{
  "timestamp": "2026-03-26T...",
  "skills_imported": [
    "backend_developer",
    "frontend_developer",
    "devops_engineer",
    "security_engineer",
    "qa_engineer",
    "data_engineer",
    "ml_engineer"
  ],
  "output_directory": "local_work/skills_import/engineering/"
}
```

### 🔧 Technical Details

**Native Skill Pattern:**
```go
type BackendDeveloperSkill struct {
    workspace string
}

func (b *BackendDeveloperSkill) Name() string
func (b *BackendDeveloperSkill) Description() string
func (b *BackendDeveloperSkill) GetInstructions() string
func (b *BackendDeveloperSkill) GetAntiPatterns() string
func (b *BackendDeveloperSkill) GetExamples() string
func (b *BackendDeveloperSkill) BuildSkillContext() string
func (b *BackendDeveloperSkill) BuildSummary() string
```

**Total Lines of Code Added:**
- 7 skill files: ~8,500 lines
- Test file: ~650 lines
- Setup checker: ~180 lines
- Loader updates: ~250 lines
- **Total: ~9,580 lines**

### ⚠️ Upgrade Notes

- **No breaking changes**: All existing configurations remain compatible
- **New skills are opt-in**: Add to `config.json` `skills` array to use
- **Example configurations**: Copy from `_examples` section in `config.example.json`

### 🎯 Usage Example

**Single Specialized Agent:**
```json
{
  "id": "backend_dev",
  "name": "Backend Developer",
  "model": "deepseek-chat",
  "skills": ["backend_developer"],
  "tools_override": ["read_file", "write_file", "edit_file", "exec"]
}
```

**Multi-Agent Team:**
```json
{
  "id": "tech_lead",
  "name": "Technical Lead",
  "skills": ["fullstack_developer", "agent_team_workflow"],
  "subagents": {
    "allow_agents": ["backend_dev", "frontend_dev", "devops_eng", "qa_eng"]
  }
}
```

---

## 2026-03-23

### ✨ New Features
- **Autonomous Agent Runtime (LP-03)**: Introduced a background runtime for each agent that automatically processes internal messages. Agents no longer need to manually call `agent_receive` to check for tasks.
- **Runtime Manager**: A new coordination layer in `AgentLoop` that manages lifecycle and goroutines for all autonomous agents.
- **Enhanced Agent Autonomy**: Agents now automatically switch to `StatusBusy` when processing an internal task and can send auto-responses upon completion.
- **Extended Configuration**: Added `runtime` options to `AgentConfig` and `AgentDefaults` in `config.json`, allowing fine-grained control over which agents have autonomous capabilities enabled.

### 🛠️ Core Improvements
- **Message Bus Integration**: Added `GetChannel()` to `AgentMessageBus` to allow direct, non-blocking subscription to agent-specific inboxes.
- **Agent Instance Updates**: `AgentInstance` now tracks its own `Runtime` configuration for faster access during autonomous execution.

---

## 2026-03-12

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

## 2026-03-04

### 🛡️ Upstream Security Patch Adaptations

Adapted and applied 2 of 6 upstream patches from audit `upstream_audit_2026-03-04.json` (see `local_work/patch_execution_log_2026-03-04.md` for full details).

- **🔒 Registry Collision Warning** (`pkg/tools/registry.go`): Added structured warning via `logger.WarnCF` when `Register()` overwrites an existing tool by name. Critical for multi-agent environments where MCP servers per agent could silently contaminate each other's tool namespace. Upstream ref: [`a2591e0`](https://github.com/sipeed/picoclaw/commit/a2591e03a942ae244b50539d4b9d26da3a0b3d58)

- **📝 JSONL Memory Store** (`pkg/memory/jsonl.go` — *new file*): Introduced append-only JSONL session history store with atomic writes (temp→fsync→rename) to prevent file corruption under concurrent multi-agent writes. Sharded mutex design (`64 shards`) eliminates cross-agent lock contention. Adapted from upstream: [`6d894d6`](https://github.com/sipeed/picoclaw/commit/6d894d6138cb89a8bc714d69b03c9a6a14cb03d7) — `fileutil` dependency replaced by inlined `writeFileAtomic` for fork compatibility.

**Patches confirmed already present in fork (no action needed):**
- `web_fetch` `ForLLM` content pass-through fix (was already at `web.go:666`)
- HTTP retry `resp.Body` close on socket leak (already in `http_retry.go`)
- `state.go` atomic temp-rename saves (already implemented)
- Shell security deny patterns for `.env`/`id_rsa`/AWS credentials (already in `shell.go`)


## 2026-03-03

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

## 2026-03-02

### 🛡️ Security & Stability
- **🛡️ Native Skills Sentinel**: Implemented `skills_sentinel.go` as a native internal security tool. It provides proactive pattern-matching protection against prompt injection (input) and system leaks (output sanitization).
- **📝 Local Auditing**: Integrated a security auditor that records all blocked attacks and suspicious activities in `local_work/AUDIT.md`.

## 2026-03-01

## 2026-03-01

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

## 2026-02-27

### ✨ Core Features
- **🛡️ Task Lock System**: Implemented atomic `.lock` files for robust disaster recovery and concurrency control among subagents.
- **🔄 Boot Rehydration**: The Gateway will now automatically wake up and re-hydrate agents interrupted by system crashes or restarts.
- **🧠 Context Compactor**: Built-in intelligent context pruning and tool-output truncation. Safely elevated default `MaxTokens` to 32,768, permanently eliminating "Context Explosion" silent drops.
- **⚡️ Tool Mutual Exclusion**: `FileLockChecker` integration prevents concurrent agents from editing the same file simultaneously.
- **🤖 o3-mini Support**: Standardized on `o3-mini` for high-performance OpenAI tasks, including automatic `max_completion_tokens` handling.
- **🌍 Qwen Regional Fixes**: Documented and implemented support for Alibaba Cloud Virginia (US-East-1) regional endpoints.

## 2026-02-27

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
