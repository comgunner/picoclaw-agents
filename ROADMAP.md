
# 🦐 PicoClaw Roadmap

> **Vision**: To build the ultimate lightweight, secure, and fully autonomous AI Agent infrastructure.automate the mundane, unleash your creativity

> **Current Version:** v1.0.0 (as of March 2026)
>
> Version numbers in roadmap items refer to feature milestones, not release versions.

---

## 🎉 Latest: Native Skills Architecture

**2026-03-03** - Introduces **Native Skills** - compiled-in documentation with zero external dependencies:

### 🚀 Native Skills Architecture

Following the security-first design of previous versions, Native Skills takes performance and reliability to the next level by compiling skills directly into the binary:

* **✅ Zero Runtime Dependencies**: No external `.md` files required - all skill documentation embedded in binary
* **✅ Enhanced Security**: Skills cannot be modified or tampered with externally
* **✅ Maximum Performance**: Direct memory access to documentation strings, no file I/O
* **✅ Automatic Updates**: Skills update with each binary release
* **✅ Type Safety**: Full Go type checking for skill interfaces and methods

### 📦 Native Skills Implementation

| Component | File | Purpose |
|-----------|------|---------|
| **Queue/Batch Skill** | `pkg/skills/queue_batch.go` | Fire-and-forget task delegation (351 lines) |
| **Test Suite** | `pkg/skills/queue_batch_test.go` | 9 comprehensive test cases (191 lines) |
| **Skills Loader** | `pkg/skills/loader.go` | Native registry with singleton pattern |
| **Context Builder** | `pkg/agent/context.go` | Dynamic skill injection into system prompt |

### 📚 Developer Documentation

* **[local_work/crear_skill_interna.md](local_work/crear_skill_interna.md)** - Complete guide for developing native skills with code templates
* **[docs/QUEUE_BATCH.en.md](docs/QUEUE_BATCH.en.md)** - Updated user documentation with native skill architecture
* **[docs/QUEUE_BATCH.es.md](docs/QUEUE_BATCH.es.md)** - Documentación de usuario actualizada con arquitectura de skills nativas

### 🔧 Native Skill Pattern

```go
type QueueBatchSkill struct {
    workspace string
}

func (q *QueueBatchSkill) Name() string
func (q *QueueBatchSkill) Description() string
func (q *QueueBatchSkill) GetInstructions() string
func (q *QueueBatchSkill) GetAntiPatterns() string
func (q *QueueBatchSkill) GetExamples() string
func (q *QueueBatchSkill) BuildSkillContext() string
func (q *QueueBatchSkill) BuildSummary() string
```

### 📊 Impact

- **Binary Size**: +35KB (negligible for modern systems)
- **RAM Usage**: No change (skills already loaded in memory)
- **Startup Time**: Improved (no file I/O for skill loading)
- **Security**: Significantly enhanced (no external file manipulation)

---

## 🚀 1. Core Optimization: Extreme Lightweight

*Our defining characteristic. We fight software bloat to ensure PicoClaw runs smoothly on the smallest embedded devices.*

* [**Memory Footprint Reduction**](https://github.com/comgunner/picoclaw/issues/346) 
  * **Goal**: Run smoothly on 64MB RAM embedded boards (e.g., low-end RISC-V SBCs) with the core process consuming < 20MB.
  * **Context**: RAM is expensive and scarce on edge devices. Memory optimization takes precedence over storage size.
  * **Action**: Analyze memory growth between releases, remove redundant dependencies, and optimize data structures.


## 🛡️ 2. Security Hardening: Defense in Depth

*Paying off early technical debt. We invite security experts to help build a "Secure-by-Default" agent.*

* **Input Defense & Permission Control**
  * **Prompt Injection Defense**: Harden JSON extraction logic to prevent LLM manipulation.
  * **Tool Abuse Prevention**: Strict parameter validation to ensure generated commands stay within safe boundaries.
  * **SSRF Protection**: Built-in blocklists for network tools to prevent accessing internal IPs (LAN/Metadata services).


* **Sandboxing & Isolation**
  * **Filesystem Sandbox**: Restrict file R/W operations to specific directories only.
  * **Context Isolation**: Prevent data leakage between different user sessions or channels.
  * **Privacy Redaction**: Auto-redact sensitive info (API Keys, PII) from logs and standard outputs.


* **Authentication & Secrets**
  * **Crypto Upgrade**: Adopt modern algorithms like `ChaCha20-Poly1305` for secret storage.
  * **OAuth 2.0 Flow**: Deprecate hardcoded API keys in the CLI; move to secure OAuth flows.



## 🔌 3. Connectivity: Protocol-First Architecture

*Connect every model, reach every platform.*

* **Provider**
  * [**Architecture Upgrade**](https://github.com/comgunner/picoclaw/issues/283): Refactor from "Vendor-based" to "Protocol-based" classification (e.g., OpenAI-compatible, Ollama-compatible). *(Status: In progress by @Daming, ETA 5 days)*
  * **Local Models**: Deep integration with **Ollama**, **vLLM**, **LM Studio**, and **Mistral** (local inference).
  * **Online Models**: Continued support for frontier closed-source models.


* **Channel**
  * **IM Matrix**: QQ, WeChat (Work), DingTalk, Feishu (Lark), Telegram, Discord, WhatsApp, LINE, Slack, Email, KOOK, Signal, ...
  * **Standards**: Support for the **OneBot** protocol.
  * [**attachment**](https://github.com/comgunner/picoclaw/issues/348): Native handling of images, audio, and video attachments.


* **Skill Marketplace**
  * [**Discovery skills**](https://github.com/comgunner/picoclaw/issues/287): Implement `find_skill` to automatically discover and install skills from the [GitHub Skills Repo] or other registries.



## 🧠 4. Advanced Capabilities: From Chatbot to Agentic AI

*Beyond conversation—focusing on action and collaboration.*

* **Operations**
  * [**MCP Support**](https://github.com/comgunner/picoclaw/issues/290): Native support for the **Model Context Protocol (MCP)**.
  * [**Browser Automation**](https://github.com/comgunner/picoclaw/issues/293): Headless browser control via CDP (Chrome DevTools Protocol) or ActionBook.
  * [**Mobile Operation**](https://github.com/comgunner/picoclaw/issues/292): Android device control (similar to BotDrop).



* **Multi-Agent Collaboration (Current Focus)**
  * ✅ **Multi-Agent Suite (Multi-Agent)**: Full implementation of subagent architecture with autonomous sessions.
  * ✅ **Atomic Task Locks**: Prevention of multi-agent file collisions and state corruption.
  * ✅ **Disaster Recovery (Boot Rehydration)**: Ability to restore agent state and context after crashes.
  * ✅ **Context Compactor**: Automated token management (Safe 32K limit) to maintain long-running reasoning tasks.
  * ✅ **External Integrations (External Integrations)**: Secure API integrations for Binance (trading), Facebook, X/Twitter, Discord (social media), and Notion (knowledge management) with credential isolation and fail-close security.
  * [ ] **Knowledge Sharing Protocol**: Selective RAG sharing between parallel subagents.
  * [ ] **MCP Server Support**: Implement MCP servers for Social Media and Notion tools (currently native-only).


## 📚 5. Developer Experience (DevEx) & Documentation

*Lowering the barrier to entry so anyone can deploy in minutes.*

* [**QuickGuide (Zero-Config Start)**](https://github.com/comgunner/picoclaw/issues/350)
  * Interactive CLI Wizard: If launched without config, automatically detect the environment and guide the user through Token/Network setup step-by-step.


* **Comprehensive Documentation**
  * ✅ [**Multilingual Parity**](https://github.com/comgunner/picoclaw/issues/351): Synchronized READMEs in 7+ languages (EN, ES, ZH, JA, FR, PT, VI).
  * ✅ [**Advanced Use-Case Gallery**](https://github.com/comgunner/picoclaw/issues/352): Documented roles and workflows for enterprise-grade development teams.
  * **Platform Guides**: Dedicated guides for Windows, macOS, Linux, and Android.
  * **Step-by-Step Tutorials**: "Babysitter-level" guides for configuring Providers and Channels.



## 🤖 6. Engineering: AI-Powered Open Source

*Born from Vibe Coding, we continue to use AI to accelerate development.*

* **AI-Enhanced CI/CD**
  * Integrate AI for automated Code Review, Linting, and PR Labeling.
  * **Bot Noise Reduction**: Optimize bot interactions to keep PR timelines clean.
  * **Issue Triage**: AI agents to analyze incoming issues and suggest preliminary fixes.



## 🎨 7. Brand & Community

* [**Logo Design**](https://github.com/comgunner/picoclaw/issues/297): We are looking for a **Mantis Shrimp (Stomatopoda)** logo design!
  * *Concept*: Needs to reflect "Small but Mighty" and "Lightning Fast Strikes."



## 🔌 8. External Integrations  - Pending Work

*Building on External Integrations's secure integration foundation.*

* **MCP Server Implementation**
  * [ ] **Social Media MCP Server**: Implement `ServeSocialMediaMCPStdio` for external MCP clients (currently placeholder)
  * [ ] **Notion MCP Server**: Implement `ServeNotionMCPStdio` for external MCP clients (currently placeholder)
  * [ ] **Binance MCP Enhancement**: Expand existing MCP server with additional trading pairs and order types


* **Testing & Quality**
  * [ ] **Unit Tests**: Add comprehensive tests for `social_facebook.go`, `social_x.go`, `notion_api.go`, `social_discord.go`
  * [ ] **Integration Tests**: End-to-end testing with real API credentials (sandbox mode)
  * [ ] **Rate Limit Persistence**: Store rate limit state between calls to prevent quota exhaustion across sessions


* **Platform Extensions**
  * [ ] **Instagram Integration**: Post to Instagram via Facebook Graph API (same token)
  * [ ] **LinkedIn Integration**: Share posts and articles via LinkedIn API
  * [ ] **Telegram Bot Integration**: Native posting tool (separate from channel receiver)
  * [ ] **YouTube Integration**: Upload videos, manage playlists


* **Security Enhancements**
  * [ ] **Credential Rotation**: Automatic token refresh for OAuth-based platforms
  * [ ] **Audit Logging**: Log all external API calls to `AUDIT.md` for compliance
  * [ ] **Permission Scoping**: Granular control over which agents can access which platforms

---

We welcome contributions to any item on this roadmap! Please submit a PR or open a discussion. Let's build the best Edge AI Agent together!
