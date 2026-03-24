<div align="center">
  <img src="assets/logo.jpg" alt="PicoClaw" width="512">

  <h1>PicoClaw-Agents</h1>
  <h3>🤖 Multi-Agent Architecture 🚀 Parallel Subagents</h3>

[中文](README.zh.md) | [Español](README.es.md) | [Français](README.fr.md) | [日本語](README.ja.md)

> **Note:** This project is an independent, hobbyist fork of the original [PicoClaw](https://github.com/sipeed/picoclaw) created by **Sipeed**. It is maintained for experimental and educational purposes. All credit for the original core architecture goes to the Sipeed team.

| Feature               | OpenClaw      | NanoBot             | PicoClaw                      | PicoClaw-Agents |
| :-------------------- | :------------ | :------------------ | :---------------------------- | :-------------- |
| Language              | TypeScript    | Python              | Go                            | Go              |
| RAM                   | >1GB          | >100MB              | < 10MB                        | < 45MB          |
| Startup (0.8GHz core) | >500s         | >30s                | <1s                           | <1s             |
| Cost                  | Mac Mini 599$ | Most Linux SBC ~50$ | Any Linux Board As low as 10$ | Any Linux       |

## ✨ Features

*   🪶 **Ultra-Lightweight**: Optimized Go implementation with minimal footprint.
*   🤖 **Multi-Agent Architecture**: v3.2 introduces **Fail-Close** security (detects invalid config), **v3.2.1** adds robust channel/bus closure handling, and **v3.2.2** adds the **Skills Sentinel** (native security layer) with proactive input/output sanitization and local auditing (`AUDIT.md`).
*   🚀 **Parallel Subagents**: Spawn multiple autonomous subagents working in parallel, each with independent model configurations.
*   🌍 **True Portability**: Single self-contained binary across RISC-V, ARM, and x86.
*   🦾 **AI-Bootstrapped**: Core implementation refined through autonomous agentic workflows.
*   📈 **Binance Integration**: Native trading tools for spot balances, futures positions (LONG/SHORT), and public ticker data via direct API or MCP server.
*   📱 **Social Media Tools**: Publish to Facebook (images + messages) and X/Twitter (tweets + threads) with multi-page support and automatic comment fallback.
*   🎨 **AI Image Generation**: Generate images from text prompts using Gemini or Ideogram. Includes script-to-image workflow and auto-posting to social media.
*   📝 **Notion Integration**: Create, query, update pages and databases for task tracking and data management.
*   🤖 **Community Manager**: Auto-generate engaging social media posts from technical content or generated images.
*   ⚡ **Fast-path Slash Commands**: Instant system commands via `/` or `#` that bypass the LLM for zero-latency approvals, status checks, and lotes (bundle) management. Works on Telegram, Discord, and CLI.
*   🖇️ **Global State Synchronization**: v3.4.1 introduces the **Global ImageGenTracker**, a shared memory space across all agents (PM, Subagents) to ensure perfect consistency in multi-agent workflows.
*   🚀 **Native Skills Architecture**: v3.4.2 compiles skills directly into the binary (`pkg/skills/queue_batch.go`), eliminating external file dependencies and enhancing security. See [docs/QUEUE_BATCH.en.md](docs/QUEUE_BATCH.en.md).

## 📢 News

2026-03-12 🎉 **PicoClaw v3.4.4 - Antigravity Support & Stability**: Full Google Antigravity OAuth support with schema sanitization, TokenBudget deadlock fix, session rehydration improvements, new `picoclaw clean` command, and hardened deny patterns. See [CHANGELOG.md](CHANGELOG.md) for full details.

2026-03-03 🎉 **PicoClaw v3.4.2 - Native Skills Architecture**: Introduced native skills compiled directly into the binary (`pkg/skills/queue_batch.go`), eliminating external `.md` file dependencies. Enhanced security, performance, and type safety. See [docs/QUEUE_BATCH.en.md](docs/QUEUE_BATCH.en.md) and [local_work/crear_skill_interna.md](local_work/crear_skill_interna.md) for developer guide.

2026-03-02 🎉 **PicoClaw v3.4.1 - Fast-path & Global Tracker**: Added instant Slash Commands (`/bundle_approve`, `/status`, etc.) for zero-latency interaction. Unified the `ImageGenTracker` across all agents for perfect multi-agent state consistency. See [docs/queue_batch.md](docs/queue_batch.md).

2026-03-01 🎉 **PicoClaw v3.4 - AI Image Generation & Community Manager**: Added native image generation (Gemini/Ideogram), script-to-image workflows, interactive post-generation menus, and community manager agent for auto-generating social media posts. See [docs/IMAGE_GEN_util.md](docs/IMAGE_GEN_util.md) for complete setup and usage examples.

2026-03-01 🎉 **PicoClaw v3.3 - External Integrations (Binance, Social Media, Notion)**: Added native tools for cryptocurrency trading (Binance futures & spot), social media publishing (Facebook & X/Twitter), and knowledge management (Notion). Configure via `config.json` or environment variables. See [SOCIAL_MEDIA.md](SOCIAL_MEDIA.md) and [docs/NOTION_util.md](docs/NOTION_util.md) for setup guides.

2026-03-01 🎉 **PicoClaw v3.2.2 - Native Skills Sentinel**: Added internal security layer (`skills_sentinel.go`) that provides real-time pattern-based protection against prompt injection and system leaks.
2026-03-01 🎉 **PicoClaw v3.2.1 - Security & Stability Hardening**: Robust message bus closure handling, resilient WeCom App background processing, and reinforced initialization validation for the shell tool.
2026-03-01 🎉 **PicoClaw v3.2 - Fail-Close Security**: Robust security policy. The command execution tool now performs strict validation of deny patterns during initialization.

2026-02-27 🎉 **PicoClaw v3.1 - Disaster Recovery & Task Locks**: Implemented atomic Task Locks to prevent agent collisions, "Boot Rehydration" for recovery from abrupt crashes, and Context Compactor (safely raising limit to 32K tokens) to eradicate context explosions in long coding tasks.


<img src="assets/compare.jpg" alt="PicoClaw" width="512">

## 🦾 Demonstration

### 🛠️ Standard Assistant Workflows

<table align="center">
  <tr align="center">
    <th><p align="center">🧩 Full-Stack Engineer</p></th>
    <th><p align="center">🗂️ Logging & Planning Management</p></th>
    <th><p align="center">🔎 Web Search & Learning</p></th>
    <th><p align="center">🔧 General Worker</p></th>
  </tr>
  <tr>
    <td align="center"></td>
    <td align="center"></td>
    <td align="center"></td>
    <td align="center"></td>
  </tr>
  <tr>
    <td align="center">Develop • Deploy • Scale</td>
    <td align="center">Schedule • Automate • Memory</td>
    <td align="center">Discovery • Insights • Trends</td>
    <td align="center">Tasks • Support • Efficiency</td>
  </tr>
</table>

### 🚀 Advanced Multi-Agent Workflow (The "Dream Team")

Take advantage of the subagent architecture to deploy a full software development lifecycle team.

**The "DevOps & QA" Team (Powered by [DeepSeek Reasoner](https://platform.deepseek.com)):**

*   **`project_manager` (Leader)**: Has permission to spawn any agent. Oversees the global objective and delegates sub-tasks.
*   **`senior_dev` (The Engine)**: Technical expert. Spawns the QA Specialist to review code or the Junior Fixer for boilerplate.
*   **`qa_specialist` (Ops & Testing)**: Quality logic. Tests code, finds bugs, proposes fixes, and manages GitHub deployments.
*   **`junior_fixer` (The Assistant)**: Focuses on small fixes, refactoring, and documentation under supervision.
*   **`general_worker` (The Groundwork)**: Versatile agent for common tasks, information retrieval, and supporting the rest of the team.

**How to use this?**
Simply send a high-level command to the Leader via Telegram or CLI:
> *"Leader, I need the Senior Dev to fix the database bug and have the QA specialist verify the build before pushing to GitHub."*

PicoClaw will automatically manage the hierarchy: **PM ➔ Senior Dev ➔ QA Specialist (Fix & Publish).**

> [!TIP]
> **Check out the examples:** See `config_dev.example.json` for a standard DeepSeek team, `config_dev_multiple_models.example.json` for a mixed-model team (OpenAI, Anthropic, and DeepSeek), and `config_context_management.example.json` for optimizing token usage during extensive coding sessions.


### 📱 Run on old Android Phones

Give your decade-old phone a second life! Turn it into a smart AI Assistant with PicoClaw. Quick Start:

1. **Install Termux** (Available on F-Droid or Google Play).
2. **Execute cmds**

```bash
# Note: Replace v0.1.1 with the latest version from the Releases page
wget https://github.com/comgunner/picoclaw-agents/releases/download/v0.1.1/picoclaw-linux-arm64
chmod +x picoclaw-linux-arm64
pkg install proot
termux-chroot ./picoclaw-linux-arm64 onboard
```

And then follow the instructions in the "Quick Start" section to complete the configuration!
<img src="assets/termux.jpg" alt="PicoClaw" width="512">

### 🐜 Innovative Low-Footprint Deploy

PicoClaw can be deployed on almost any Linux device, from simple embedded boards to powerful servers.

🌟 More Deployment Cases Await！

## 📦 Install

### Install with precompiled binary

Download the firmware for your platform from the [release](https://github.com/comgunner/picoclaw-agents/releases) page.

### Install from source (latest features, recommended for development)

```bash
git clone https://github.com/comgunner/picoclaw-agents.git

cd picoclaw
make deps

# Build, no need to install
make build

# Build for multiple platforms
make build-all

# Build And Install
make install
```

## 🐳 Docker Compose

You can also run PicoClaw using Docker Compose without installing anything locally.

```bash
# 1. Clone this repo
git clone https://github.com/comgunner/picoclaw-agents.git
cd picoclaw

# 2. Set your API keys
cp config/config.example.json config/config.json
vim config/config.json      # Set DISCORD_BOT_TOKEN, API keys, etc.

# 3. Build & Start
docker compose --profile gateway up -d

> [!TIP]
> **Docker Users**: By default, the Gateway listens on `127.0.0.1` which is not accessible from the host. If you need to access the health endpoints or expose ports, set `PICOCLAW_GATEWAY_HOST=0.0.0.0` in your environment or update `config.json`.


# 4. Check logs
docker compose logs -f picoclaw-gateway

# 5. Stop
docker compose --profile gateway down
```

### Agent Mode (One-shot)

```bash
# Ask a question
docker compose run --rm picoclaw-agent -m "What is 2+2?"

# Interactive mode
docker compose run --rm picoclaw-agent
```

### Rebuild

```bash
docker compose --profile gateway build --no-cache
docker compose --profile gateway up -d
```

### 🚀 Quick Start

> [!TIP]
> Set your API key in `~/.picoclaw/config.json`.
> Get API keys: [OpenRouter](https://openrouter.ai/keys) (LLM) · [Zhipu](https://open.bigmodel.cn/usercenter/proj-mgmt/apikeys) (LLM)
> Web Search is **optional** - get free [Tavily API](https://tavily.com) (1000 free queries/month) or [Brave Search API](https://brave.com/search/api) (2000 free queries/month) or use built-in auto fallback.

**1. Initialize**

Use the `onboard` command to initialize your workspace with a pre-configured template for your preferred provider:

```bash
# Default (Empty/Manual configuration)
picoclaw onboard

# 🆓 Zero-cost setup — no API balance required:
picoclaw onboard --free        # Free tier (OpenRouter free models)

# Pre-configured templates:
picoclaw onboard --openai      # Use OpenAI template (o3-mini)
picoclaw onboard --openrouter  # Use OpenRouter template (openrouter/auto)
picoclaw onboard --glm         # Use GLM-4.5-Flash template (zhipu.ai)
picoclaw onboard --qwen        # Use Qwen template (Alibaba Cloud Intl)
picoclaw onboard --qwen_zh     # Use Qwen template (Alibaba Cloud China)
picoclaw onboard --gemini      # Use Gemini template (gemini-2.5-flash)
```

> [!TIP]
> **No API key balance?** Use `picoclaw onboard --free` to get started instantly with OpenRouter's free-tier models. Just create a free account at [openrouter.ai](https://openrouter.ai) and add your key — no credits needed.

#### 🆓 Free Tier Models

The `--free` flag configures three OpenRouter free-tier models with automatic fallback:

| Priority | Model | Context | Notes |
|----------|-------|---------|-------|
| Primary | `openrouter/free` | varies | Auto-selects best available free model |
| Fallback 1 | `stepfun/step-3.5-flash` | 256K | High-context tasks |
| Fallback 2 | `deepseek/deepseek-v3.2-20251201` | 64K | Fast, reliable fallback |

All three are routed through [OpenRouter](https://openrouter.ai) — a single API key covers all of them.

**2. Configure** (`~/.picoclaw/config.json`)

```json
{
  "agents": {
    "defaults": {
      "workspace": "~/.picoclaw/workspace",
      "model_name": "deepseek-chat",
      "max_tokens": 8192,
      "temperature": 0.7,
      "max_tool_iterations": 20,
      "subagents": {
        "max_spawn_depth": 2,
        "max_children_per_agent": 5
      }
    },
    "backend_coder": {
      "model_name": "deepseek-reasoner",
      "temperature": 0.2
    }
  },
  "model_list": [
    {
      "model_name": "deepseek-chat",
      "model": "deepseek/deepseek-chat",
      "api_key": "your-api-key"
    },
    {
      "model_name": "deepseek-reasoner",
      "model": "deepseek/deepseek-reasoner",
      "api_key": "your-api-key"
    }
  ],
  "tools": {
    "web": {
      "brave": {
        "enabled": false,
        "api_key": "YOUR_BRAVE_API_KEY",
        "max_results": 5
      },
      "tavily": {
        "enabled": false,
        "api_key": "YOUR_TAVILY_API_KEY",
        "max_results": 5
      },
      "duckduckgo": {
        "enabled": true,
        "max_results": 5
      }
    }
  }
}
```

> **New in v3 (Multi-Agent Architecture)**: You can now spin up isolated **Subagents** to perform parallel background tasks. Crucially, **each subagent can use a completely different LLM model**. As shown in the configuration above, the main agent runs `gpt4`, but it can spawn a dedicated `coder` subagent running `claude-sonnet-4.6` to handle complex programming tasks simultaneously!

> **New**: The `model_list` configuration format allows zero-code provider addition. See [Model Configuration](#model-configuration-model_list) for details.
> `request_timeout` is optional and uses seconds. If omitted or set to `<= 0`, PicoClaw uses the default timeout (120s).

**3. Get API Keys**

* **LLM Provider**: [DeepSeek](https://platform.deepseek.com) (Recommended) · [OpenRouter](https://openrouter.ai/keys) · [Zhipu](https://open.bigmodel.cn/usercenter/proj-mgmt/apikeys) · [Anthropic](https://console.anthropic.com) · [OpenAI](https://platform.openai.com) · [Gemini](https://aistudio.google.com/api-keys)
* **Web Search** (optional): [Tavily](https://tavily.com) - Optimized for AI Agents (1000 requests/month) · [Brave Search](https://brave.com/search/api) - Free tier available (2000 requests/month)

### 💡 Recommended Models for Developers (`backend_coder`)

For heavy coding tasks, performance and logic are key. We recommend standardizing on these models for your `backend_coder` agents:

*   **DeepSeek**: `deepseek-reasoner` (Excellent reasoning and cost-effective)
*   **OpenAI**: `o3-mini-2025-01-31` (High performance)
*   **OpenRouter.ai**: `Qwen3 Coder Plus`, `GPT-5.3-Codex` (Great coding versatility)
*   **Anthropic**: `Claude Haiku 4.5` (Fast and reliable)

> **Note**: See `config.example.json` for a complete configuration template.

**4. Chat**

```bash
picoclaw agent -m "What is 2+2?"
```

That's it! You have a working AI assistant in 2 minutes.

---

## 💬 Chat Apps

Talk to your picoclaw through Telegram, Discord, DingTalk, LINE, or WeCom

| Channel      | Setup                              |
| ------------ | ---------------------------------- |
| **Telegram** | Easy (just a token)                |
| **Discord**  | Easy (bot token + intents)         |
| **QQ**       | Easy (AppID + AppSecret)           |
| **DingTalk** | Medium (app credentials)           |
| **LINE**     | Medium (credentials + webhook URL) |
| **WeCom**    | Medium (CorpID + webhook setup)    |

<details>
<summary><b>Telegram</b> (Recommended)</summary>

**1. Create a bot**

* Open Telegram, search `@BotFather`
* Send `/newbot`, follow prompts
* Copy the token

**2. Configure**

```json
{
  "channels": {
    "telegram": {
      "enabled": true,
      "token": "YOUR_BOT_TOKEN",
      "allow_from": ["YOUR_USER_ID"]
    }
  }
}
```

> Get your user ID from `@userinfobot` on Telegram.

**3. Run**

```bash
picoclaw gateway
```

</details>

<details>
<summary><b>Discord</b></summary>

**1. Create a bot**

* Go to <https://discord.com/developers/applications>
* Create an application → Bot → Add Bot
* Copy the bot token

**2. Enable intents**

* In the Bot settings, enable **MESSAGE CONTENT INTENT**
* (Optional) Enable **SERVER MEMBERS INTENT** if you plan to use allow lists based on member data

**3. Get your User ID**
* Discord Settings → Advanced → enable **Developer Mode**
* Right-click your avatar → **Copy User ID**

**4. Configure**

```json
{
  "channels": {
    "discord": {
      "enabled": true,
      "token": "YOUR_BOT_TOKEN",
      "allow_from": ["YOUR_USER_ID"],
      "mention_only": false
    }
  }
}
```

**5. Invite the bot**

* OAuth2 → URL Generator
* Scopes: `bot`
* Bot Permissions: `Send Messages`, `Read Message History`
* Open the generated invite URL and add the bot to your server

**Optional: Mention-only mode**

Set `"mention_only": true` to make the bot respond only when @-mentioned. Useful for shared servers where you want the bot to respond only when explicitly called.

**6. Run**

```bash
picoclaw gateway
```

</details>

<details>
<summary><b>QQ</b></summary>

**1. Create a bot**

- Go to [QQ Open Platform](https://q.qq.com/#)
- Create an application → Get **AppID** and **AppSecret**

**2. Configure**

```json
{
  "channels": {
    "qq": {
      "enabled": true,
      "app_id": "YOUR_APP_ID",
      "app_secret": "YOUR_APP_SECRET",
      "allow_from": []
    }
  }
}
```

> Set `allow_from` to empty to allow all users, or specify QQ numbers to restrict access.

**3. Run**

```bash
picoclaw gateway
```

</details>

<details>
<summary><b>DingTalk</b></summary>

**1. Create a bot**

* Go to [Open Platform](https://open.dingtalk.com/)
* Create an internal app
* Copy Client ID and Client Secret

**2. Configure**

```json
{
  "channels": {
    "dingtalk": {
      "enabled": true,
      "client_id": "YOUR_CLIENT_ID",
      "client_secret": "YOUR_CLIENT_SECRET",
      "allow_from": []
    }
  }
}
```

> Set `allow_from` to empty to allow all users, or specify DingTalk user IDs to restrict access.

**3. Run**

```bash
picoclaw gateway
```
</details>

<details>
<summary><b>LINE</b></summary>

**1. Create a LINE Official Account**

- Go to [LINE Developers Console](https://developers.line.biz/)
- Create a provider → Create a Messaging API channel
- Copy **Channel Secret** and **Channel Access Token**

**2. Configure**

```json
{
  "channels": {
    "line": {
      "enabled": true,
      "channel_secret": "YOUR_CHANNEL_SECRET",
      "channel_access_token": "YOUR_CHANNEL_ACCESS_TOKEN",
      "webhook_host": "0.0.0.0",
      "webhook_port": 18791,
      "webhook_path": "/webhook/line",
      "allow_from": []
    }
  }
}
```

**3. Set up Webhook URL**

LINE requires HTTPS for webhooks. Use a reverse proxy or tunnel:

```bash
# Example with ngrok
ngrok http 18791
```

Then set the Webhook URL in LINE Developers Console to `https://your-domain/webhook/line` and enable **Use webhook**.

**4. Run**

```bash
picoclaw gateway
```

> In group chats, the bot responds only when @mentioned. Replies quote the original message.

> **Docker Compose**: Add `ports: ["18791:18791"]` to the `picoclaw-gateway` service to expose the webhook port.

</details>

<details>
<summary><b>WeCom (企业微信)</b></summary>

PicoClaw supports two types of WeCom integration:

**Option 1: WeCom Bot (智能机器人)** - Easier setup, supports group chats
**Option 2: WeCom App (自建应用)** - More features, proactive messaging

See [WeCom App Configuration Guide](docs/wecom-app-configuration.md) for detailed setup instructions.

**Quick Setup - WeCom Bot:**

**1. Create a bot**

* Go to WeCom Admin Console → Group Chat → Add Group Bot
* Copy the webhook URL (format: `https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=xxx`)

**2. Configure**

```json
{
  "channels": {
    "wecom": {
      "enabled": true,
      "token": "YOUR_TOKEN",
      "encoding_aes_key": "YOUR_ENCODING_AES_KEY",
      "webhook_url": "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=YOUR_KEY",
      "webhook_host": "0.0.0.0",
      "webhook_port": 18793,
      "webhook_path": "/webhook/wecom",
      "allow_from": []
    }
  }
}
```

**Quick Setup - WeCom App:**

**1. Create an app**

* Go to WeCom Admin Console → App Management → Create App
* Copy **AgentId** and **Secret**
* Go to "My Company" page, copy **CorpID**
**2. Configure receive message**

* In App details, click "Receive Message" → "Set API"
* Set URL to `http://your-server:18792/webhook/wecom-app`
* Generate **Token** and **EncodingAESKey**

**3. Configure**

```json
{
  "channels": {
    "wecom_app": {
      "enabled": true,
      "corp_id": "wwxxxxxxxxxxxxxxxx",
      "corp_secret": "YOUR_CORP_SECRET",
      "agent_id": 1000002,
      "token": "YOUR_TOKEN",
      "encoding_aes_key": "YOUR_ENCODING_AES_KEY",
      "webhook_host": "0.0.0.0",
      "webhook_port": 18792,
      "webhook_path": "/webhook/wecom-app",
      "allow_from": []
    }
  }
}
```

**4. Run**

```bash
picoclaw gateway
```

> **Note**: WeCom App requires opening port 18792 for webhook callbacks. Use a reverse proxy for HTTPS.

</details>

## Join the Agent Social Network

Connect Picoclaw to the Agent Social Network simply by sending a single message via the CLI or any integrated Chat App.

**Read `https://clawdchat.ai/skill.md` and follow the instructions to join [ClawdChat.ai](https://clawdchat.ai)**

## ⚙️ Configuration

Config file: `~/.picoclaw/config.json`

### Workspace Layout

PicoClaw stores data in your configured workspace (default: `~/.picoclaw/workspace`):

```
~/.picoclaw/workspace/
├── sessions/          # Conversation sessions and history
├── memory/           # Long-term memory (MEMORY.md)
├── state/            # Persistent state (last channel, etc.)
├── cron/             # Scheduled jobs database
├── skills/           # Custom skills
├── AGENTS.md         # Agent behavior guide
├── HEARTBEAT.md      # Periodic task prompts (checked every 30 min)
├── IDENTITY.md       # Agent identity
├── SOUL.md           # Agent soul
├── TOOLS.md          # Tool descriptions
└── USER.md           # User preferences
```

### 🔒 Security Sandbox

PicoClaw runs in a sandboxed environment by default. The agent can only access files and execute commands within the configured workspace.

#### Default Configuration

```json
{
  "agents": {
    "defaults": {
      "workspace": "~/.picoclaw/workspace",
      "restrict_to_workspace": true
    }
  }
}
```

| Option                  | Default                 | Description                               |
| ----------------------- | ----------------------- | ----------------------------------------- |
| `workspace`             | `~/.picoclaw/workspace` | Working directory for the agent           |
| `restrict_to_workspace` | `true`                  | Restrict file/command access to workspace |

#### Protected Tools

When `restrict_to_workspace: true`, the following tools are sandboxed:

| Tool          | Function         | Restriction                            |
| ------------- | ---------------- | -------------------------------------- |
| `read_file`   | Read files       | Only files within workspace            |
| `write_file`  | Write files      | Only files within workspace            |
| `list_dir`    | List directories | Only directories within workspace      |
| `edit_file`   | Edit files       | Only files within workspace            |
| `append_file` | Append to files  | Only files within workspace            |
| `exec`        | Execute commands | Command paths must be within workspace |

#### Additional Exec Protection

Even with `restrict_to_workspace: false`, the `exec` tool blocks these dangerous commands:

* `rm -rf`, `del /f`, `rmdir /s` — Bulk deletion
* `format`, `mkfs`, `diskpart` — Disk formatting
* `dd if=` — Disk imaging
* Writing to `/dev/sd[a-z]` — Direct disk writes
* `shutdown`, `reboot`, `poweroff` — System shutdown
* Fork bomb `:(){ :|:& };:`

#### Core Infrastructure Protections (v3.4.3+)

PicoClaw's multi-agent architecture incorporates several upstream security patches to ensure safe concurrent operations:
* **Atomic State Saves**: `memory/jsonl.go` and `state/state.go` persist data via strict temp-file `fsync` followed by atomic `rename`, completely eliminating JSON corruption during power losses or subagent crashes.
* **MCP Collision Warning**: Strict registry overlap detection prevents spawned agents from silently polluting each other's Tool registries or MCP spaces.
* **Socket Leak Prevention**: Robust forced closure on HTTP retries prevents OS-level file descriptor exhaustions on flaky connections.

#### Error Examples

```
[ERROR] tool: Tool execution failed
{tool=exec, error=Command blocked by safety guard (path outside working dir)}
```

```
[ERROR] tool: Tool execution failed
{tool=exec, error=Command blocked by safety guard (dangerous pattern detected)}
```

#### Disabling Restrictions (Security Risk)

If you need the agent to access paths outside the workspace:

**Method 1: Config file**

```json
{
  "agents": {
    "defaults": {
      "restrict_to_workspace": false
    }
  }
}
```

**Method 2: Environment variable**

```bash
export PICOCLAW_AGENTS_DEFAULTS_RESTRICT_TO_WORKSPACE=false
```

> ⚠️ **Warning**: Disabling this restriction allows the agent to access any path on your system. Use with caution in controlled environments only.

#### Security Boundary Consistency

The `restrict_to_workspace` setting applies consistently across all execution paths:

| Execution Path   | Security Boundary           |
| ---------------- | --------------------------- |
| Main Agent       | `restrict_to_workspace` ✅   |
| Subagent / Spawn | Inherits same restriction ✅ |
| Heartbeat tasks  | Inherits same restriction ✅ |

All paths share the same workspace restriction — there's no way to bypass the security boundary through subagents or scheduled tasks.

### Heartbeat (Periodic Tasks)

PicoClaw can perform periodic tasks automatically. Create a `HEARTBEAT.md` file in your workspace:

```markdown
# Periodic Tasks

- Check my email for important messages
- Review my calendar for upcoming events
- Check the weather forecast
```

The agent will read this file every 30 minutes (configurable) and execute any tasks using available tools.

#### Async Tasks with Spawn

For long-running tasks (web search, API calls), use the `spawn` tool to create a **subagent**:

```markdown
# Periodic Tasks

## Quick Tasks (respond directly)

- Report current time

## Long Tasks (use spawn for async)

- Search the web for AI news and summarize
- Check email and report important messages
```

**Key behaviors:**

| Feature                 | Description                                               |
| ----------------------- | --------------------------------------------------------- |
| **spawn**               | Creates async subagent, doesn't block heartbeat           |
| **Independent context** | Subagent has its own context, no session history          |
| **message tool**        | Subagent communicates with user directly via message tool |
| **Non-blocking**        | After spawning, heartbeat continues to next task          |

#### How Subagent Communication Works

```
Heartbeat triggers
    ↓
Agent reads HEARTBEAT.md
    ↓
For long task: spawn subagent
    ↓                           ↓
Continue to next task      Subagent works independently
    ↓                           ↓
All tasks done            Subagent uses "message" tool
    ↓                           ↓
Respond HEARTBEAT_OK      User receives result directly
```

The subagent has access to tools (message, web_search, etc.) and can communicate with the user independently without going through the main agent.

**Configuration:**

```json
{
  "heartbeat": {
    "enabled": true,
    "interval": 30
  }
}
```

| Option     | Default | Description                        |
| ---------- | ------- | ---------------------------------- |
| `enabled`  | `true`  | Enable/disable heartbeat           |
| `interval` | `30`    | Check interval in minutes (min: 5) |

**Environment variables:**

* `PICOCLAW_HEARTBEAT_ENABLED=false` to disable
* `PICOCLAW_HEARTBEAT_INTERVAL=60` to change interval

### Providers

> [!NOTE]
> Groq provides free voice transcription via Whisper. If configured, Telegram voice messages will be automatically transcribed.

| Provider                   | Purpose                                 | Get API Key                                                          |
| -------------------------- | --------------------------------------- | -------------------------------------------------------------------- |
| `gemini`                   | LLM (Gemini direct)                     | [aistudio.google.com](https://aistudio.google.com)                   |
| `zhipu`                    | LLM (Zhipu direct)                      | [bigmodel.cn](https://bigmodel.cn)                                   |
| `openrouter(To be tested)` | LLM (recommended, access to all models) | [openrouter.ai](https://openrouter.ai)                               |
| `anthropic(To be tested)`  | LLM (Claude direct)                     | [console.anthropic.com](https://console.anthropic.com)               |
| `openai(To be tested)`     | LLM (GPT direct)                        | [platform.openai.com](https://platform.openai.com)                   |
| `deepseek(To be tested)`   | LLM (DeepSeek direct)                   | [platform.deepseek.com](https://platform.deepseek.com)               |
| `qwen`                     | LLM (Qwen direct)                       | [dashscope.console.aliyun.com](https://dashscope.console.aliyun.com) |
| `groq`                     | LLM + **Voice transcription** (Whisper) | [console.groq.com](https://console.groq.com)                         |
| `cerebras`                 | LLM (Cerebras direct)                   | [cerebras.ai](https://cerebras.ai)                                   |
| `antigravity`              | LLM (Google Antigravity / OAuth)        | `picoclaw auth login --provider google-antigravity`                  |

### Model Configuration (model_list)

> **What's New?** PicoClaw now uses a **model-centric** configuration approach. Simply specify `vendor/model` format (e.g., `zhipu/glm-4.5-flash`) to add new providers—**zero code changes required!**

This design also enables **multi-agent support** with flexible provider selection:

- **Different agents, different providers**: Each agent can use its own LLM provider
- **Model fallbacks**: Configure primary and fallback models for resilience
- **Load balancing**: Distribute requests across multiple endpoints
- **Centralized configuration**: Manage all providers in one place

#### 📋 All Supported Vendors

| Vendor              | `model` Prefix    | Default API Base                                    | Protocol  | API Key                                                          |
| ------------------- | ----------------- | --------------------------------------------------- | --------- | ---------------------------------------------------------------- |
| **OpenAI**          | `openai/`         | `https://api.openai.com/v1`                         | OpenAI    | [Get Key](https://platform.openai.com)                           |
| **Anthropic**       | `anthropic/`      | `https://api.anthropic.com/v1`                      | Anthropic | [Get Key](https://console.anthropic.com)                         |
| **智谱 AI (GLM)**   | `zhipu/`          | `https://open.bigmodel.cn/api/paas/v4`              | OpenAI    | [Get Key](https://open.bigmodel.cn/usercenter/proj-mgmt/apikeys) |
| **DeepSeek**        | `deepseek/`       | `https://api.deepseek.com/v1`                       | OpenAI    | [Get Key](https://platform.deepseek.com)                         |
| **Google Gemini**   | `gemini/`         | `https://generativelanguage.googleapis.com/v1beta`  | OpenAI    | [Get Key](https://aistudio.google.com/api-keys)                  |
| **Groq**            | `groq/`           | `https://api.groq.com/openai/v1`                    | OpenAI    | [Get Key](https://console.groq.com)                              |
| **Moonshot**        | `moonshot/`       | `https://api.moonshot.cn/v1`                        | OpenAI    | [Get Key](https://platform.moonshot.cn)                          |
| **通义千问 (Qwen)** | `qwen/`           | `https://dashscope.aliyuncs.com/compatible-mode/v1` | OpenAI    | [Get Key](https://dashscope.console.aliyun.com)                  |
| **NVIDIA**          | `nvidia/`         | `https://integrate.api.nvidia.com/v1`               | OpenAI    | [Get Key](https://build.nvidia.com)                              |
| **Ollama**          | `ollama/`         | `http://localhost:11434/v1`                         | OpenAI    | Local (no key needed)                                            |
| **OpenRouter**      | `openrouter/`     | `https://openrouter.ai/api/v1`                      | OpenAI    | [Get Key](https://openrouter.ai/keys)                            |
| **VLLM**            | `vllm/`           | `http://localhost:8000/v1`                          | OpenAI    | Local                                                            |
| **Cerebras**        | `cerebras/`       | `https://api.cerebras.ai/v1`                        | OpenAI    | [Get Key](https://cerebras.ai)                                   |
| **火山引擎**        | `volcengine/`     | `https://ark.cn-beijing.volces.com/api/v3`          | OpenAI    | [Get Key](https://console.volcengine.com)                        |
| **神算云**          | `shengsuanyun/`   | `https://router.shengsuanyun.com/api/v1`            | OpenAI    | -                                                                |
| **Antigravity**     | `antigravity/`    | Google Cloud                                        | Custom    | OAuth only                                                       |
| **GitHub Copilot**  | `github-copilot/` | `localhost:4321`                                    | gRPC      | -                                                                |

#### Basic Configuration

```json
{
  "model_list": [
    {
      "model_name": "deepseek-chat",
      "model": "deepseek/deepseek-chat",
      "api_key": "sk-your-key"
    },
    {
      "model_name": "deepseek-reasoner",
      "model": "deepseek/deepseek-reasoner",
      "api_key": "sk-your-key"
    },
    {
      "model_name": "o3-mini-2025-01-31",
      "model": "openai/o3-mini-2025-01-31",
      "api_key": "sk-your-key"
    }
  ],
  "agents": {
    "defaults": {
      "model": "deepseek-chat"
    },
    "backend_coder": {
      "model": "deepseek-reasoner"
    }
  }
}
```

#### Vendor-Specific Examples

**OpenAI**

```json
{
  "model_name": "gpt-5.2",
  "model": "openai/gpt-5.2",
  "api_key": "sk-..."
}
```

**智谱 AI (GLM)**

```json
{
  "model_name": "glm-4.5-flash",
  "model": "zhipu/glm-4.5-flash",
  "api_key": "your-key"
}
```

**DeepSeek**

```json
{
  "model_name": "deepseek-chat",
  "model": "deepseek/deepseek-chat",
  "api_key": "sk-..."
}
```

**Anthropic (with API key)**

```json
{
  "model_name": "claude-sonnet-4.6",
  "model": "anthropic/claude-sonnet-4.6",
  "api_key": "sk-ant-your-key"
}
```

> Run `picoclaw auth login --provider anthropic` to paste your API token.

**Google Antigravity (OAuth — free tier)**

```json
{
  "model_name": "antigravity-gemini-3-flash",
  "model": "antigravity/gemini-3-flash",
  "auth_method": "oauth"
}
```

> Run `picoclaw auth login --provider google-antigravity` to authenticate via browser. No API key required — uses your Google account. See [docs/ANTIGRAVITY_QUICKSTART.md](docs/ANTIGRAVITY_QUICKSTART.md) for setup details.

**Ollama (local)**

```json
{
  "model_name": "llama3",
  "model": "ollama/llama3"
}
```

**Custom Proxy/API**

```json
{
  "model_name": "my-custom-model",
  "model": "openai/custom-model",
  "api_base": "https://my-proxy.com/v1",
  "api_key": "sk-...",
  "request_timeout": 300
}
```

#### Load Balancing

Configure multiple endpoints for the same model name—PicoClaw will automatically round-robin between them:

```json
{
  "model_list": [
    {
      "model_name": "gpt-5.2",
      "model": "openai/gpt-5.2",
      "api_base": "https://api1.example.com/v1",
      "api_key": "sk-key1"
    },
    {
      "model_name": "gpt-5.2",
      "model": "openai/gpt-5.2",
      "api_base": "https://api2.example.com/v1",
      "api_key": "sk-key2"
    }
  ]
}
```

#### Migration from Legacy `providers` Config

The old `providers` configuration is **deprecated** but still supported for backward compatibility.

**Old Config (deprecated):**

```json
{
  "providers": {
    "zhipu": {
      "api_key": "your-key",
      "api_base": "https://open.bigmodel.cn/api/paas/v4"
    }
  },
  "agents": {
    "defaults": {
      "provider": "zhipu",
      "model": "glm-4.5-flash"
    }
  }
}
```

**New Config (recommended):**

```json
{
  "model_list": [
    {
      "model_name": "glm-4.5-flash",
      "model": "zhipu/glm-4.5-flash",
      "api_key": "your-key"
    }
  ],
  "agents": {
    "defaults": {
      "model": "glm-4.5-flash"
    }
  }
}
```

For detailed migration guide, see [docs/migration/model-list-migration.md](docs/migration/model-list-migration.md).

### Provider Architecture

PicoClaw routes providers by protocol family:

- OpenAI-compatible protocol: OpenRouter, OpenAI-compatible gateways, Groq, Zhipu, and vLLM-style endpoints.
- Anthropic protocol: Claude-native API behavior.
- Codex/OAuth path: OpenAI OAuth/token authentication route.

This keeps the runtime lightweight while making new OpenAI-compatible backends mostly a config operation (`api_base` + `api_key`).

<details>
<summary><b>Zhipu</b></summary>

**1. Get API key and base URL**

* Get [API key](https://bigmodel.cn/usercenter/proj-mgmt/apikeys)

**2. Configure**

```json
{
  "agents": {
    "defaults": {
      "workspace": "~/.picoclaw/workspace",
      "model": "glm-4.5-flash",
      "max_tokens": 8192,
      "temperature": 0.7,
      "max_tool_iterations": 20
    }
  },
  "providers": {
    "zhipu": {
      "api_key": "Your API Key",
      "api_base": "https://open.bigmodel.cn/api/paas/v4"
    }
  }
}
```

**3. Run**

```bash
picoclaw agent -m "Hello"
```

</details>

<details>
<summary><b>Full config example</b></summary>

```json
{
  "agents": {
    "defaults": {
      "model": "anthropic/claude-opus-4-5"
    }
  },
  "providers": {
    "openrouter": {
      "api_key": "sk-or-v1-xxx"
    },
    "groq": {
      "api_key": "gsk_xxx"
    }
  },
  "channels": {
    "telegram": {
      "enabled": true,
      "token": "123456:ABC...",
      "allow_from": ["123456789"]
    },
    "discord": {
      "enabled": true,
      "token": "",
      "allow_from": [""]
    },
    "whatsapp": {
      "enabled": false
    },
    "feishu": {
      "enabled": false,
      "app_id": "cli_xxx",
      "app_secret": "xxx",
      "encrypt_key": "",
      "verification_token": "",
      "allow_from": []
    },
    "qq": {
      "enabled": false,
      "app_id": "",
      "app_secret": "",
      "allow_from": []
    }
  },
  "tools": {
    "web": {
      "brave": {
        "enabled": false,
        "api_key": "BSA...",
        "max_results": 5
      },
      "duckduckgo": {
        "enabled": true,
        "max_results": 5
      }
    },
    "cron": {
      "exec_timeout_minutes": 5
    }
  },
  "heartbeat": {
    "enabled": true,
    "interval": 30
  }
}
```

</details>

## CLI Reference

| Command                   | Description                   |
| ------------------------- | ----------------------------- |
| `picoclaw onboard`        | Initialize config & workspace |
| `picoclaw agent -m "..."` | Chat with the agent           |
| `picoclaw agent`          | Interactive chat mode         |
| `picoclaw gateway`        | Start the gateway             |
| `picoclaw status`         | Show status                   |
| `picoclaw cron list`      | List all scheduled jobs       |
| `picoclaw cron add ...`   | Add a scheduled job           |

### Scheduled Tasks / Reminders

PicoClaw supports scheduled reminders and recurring tasks through the `cron` tool:

* **One-time reminders**: "Remind me in 10 minutes" → triggers once after 10min
* **Recurring tasks**: "Remind me every 2 hours" → triggers every 2 hours
* **Cron expressions**: "Remind me at 9am daily" → uses cron expression

Jobs are stored in `~/.picoclaw/workspace/cron/` and processed automatically.

### Binance Integration (Native Tools + MCP)

PicoClaw includes native Binance tools in `agent` mode:

* `binance_get_ticker_price` (public market ticker)
* `binance_get_spot_balance` (signed endpoint, requires API key/secret)

Configure keys in `~/.picoclaw/config.json`:

```json
{
  "tools": {
    "binance": {
      "api_key": "YOUR_BINANCE_API_KEY",
      "secret_key": "YOUR_BINANCE_SECRET_KEY"
    }
  }
}
```

Usage examples:

```bash
picoclaw agent -m "Use binance_get_ticker_price with symbol BTCUSDT and return only the numeric price."
picoclaw agent -m "Use binance_get_spot_balance and show my non-zero balances."
```

Behavior without API keys:

* `binance_get_ticker_price` still works via Binance public endpoint and adds a public-endpoint notice.
* `binance_get_spot_balance` warns that keys are missing and suggests public `curl` usage.

Optional MCP server mode (for MCP clients):

```bash
picoclaw util binance-mcp-server
```

Example `mcp_servers` config (use the absolute `picoclaw` path generated by your installation/onboard flow):

```json
{
  "mcp_servers": {
    "binance": {
      "enabled": true,
      "command": "/absolute/path/to/picoclaw",
      "args": ["util", "binance-mcp-server"]
    }
  }
}
```

## 🤝 Contribute & Roadmap

See our full [Roadmap](ROADMAP.md).

Discord: [Próximamente / Coming Soon]



## 🐛 Troubleshooting

### Web search says "API key configuration issue"

This is normal if you haven't configured a search API key yet. PicoClaw will provide helpful links for manual searching.

To enable web search:

1. **Option 1 (Recommended)**: Get a free API key at [https://brave.com/search/api](https://brave.com/search/api) (2000 free queries/month) for the best results.
2. **Option 2 (No Credit Card)**: If you don't have a key, we automatically fall back to **DuckDuckGo** (no key required).

Add the key to `~/.picoclaw/config.json` if using Brave:

```json
{
  "tools": {
    "web": {
      "brave": {
        "enabled": false,
        "api_key": "YOUR_BRAVE_API_KEY",
        "max_results": 5
      },
      "duckduckgo": {
        "enabled": true,
        "max_results": 5
      }
    }
  }
}
```

### Getting content filtering errors

Some providers (like Zhipu) have content filtering. Try rephrasing your query or use a different model.

### Telegram bot says "Conflict: terminated by other getUpdates"

This happens when another instance of the bot is running. Make sure only one `picoclaw gateway` is running at a time.

---

## 📝 API Key Comparison

| Service          | Free Tier           | Use Case                               |
| ---------------- | ------------------- | -------------------------------------- |
| **OpenRouter**   | 200K tokens/month   | Multiple models (Claude, GPT-4, etc.)  |
| **Zhipu**        | Free tier available | glm-4.5-flash (Best for Chinese users) |
| **Brave Search** | 2000 queries/month  | Web search functionality               |
| **Groq**         | Free tier available | Fast inference (Llama, Mixtral)        |
| **Cerebras**     | Free tier available | Fast inference (Llama, Qwen, etc.)     |

## ⚠️ Disclaimer

This software is provided "AS IS", without warranty of any kind, express or implied, including but not limited to the warranties of merchantability, fitness for a particular purpose, and non-infringement. In no event shall the authors or copyright holders of this fork be liable for any claim, damages, or other liability, whether in an action of contract, tort, or otherwise, arising from, out of, or in connection with the software or the use or other dealings in the software. **Use at your own risk.**
