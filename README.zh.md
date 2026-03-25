<div align="center">
  <img src="assets/logo.jpg" alt="PicoClaw" width="512">

  <h1>PicoClaw-Agents</h1>
  <h3>🤖 多智能体架构 🚀 并行子智能体</h3>

**中文** | [English](README.md) | [Español](README.es.md) | [Français](README.fr.md) | [日本語](README.ja.md)

> **注意：** 本项目是 **Sipeed** 创建的原始 [PicoClaw](https://github.com/sipeed/picoclaw) 的独立业余分叉（fork）。维护目的是为了实验和教育。原始核心架构的所有功劳归 Sipeed 团队所有。

| 特性                   | OpenClaw      | NanoBot               | PicoClaw                   | PicoClaw-Agents |
| :--------------------- | :------------ | :-------------------- | :------------------------- | :-------------- |
| 语言                   | TypeScript    | Python                | Go                         | Go              |
| RAM                    | >1GB          | >100MB                | < 10MB                     | < 45MB          |
| 启动时间 (0.8GHz core) | >500s         | >30s                  | <1s                        | <1s             |
| 成本                   | Mac Mini $599 | 大多数 Linux SBC ~$50 | 任何 Linux 开发板 低至 $10 | 任何 Linux      |

## ✨ 特性

*   🪶 **极致轻量**: 优化的 Go 语言实现，极低的内存占用。
*   🤖 **多智能体架构**: v3.2 引入 **Fail-Close** 安全机制，**v3.2.1** 优化稳定性，**v3.2.2** 增加原生 **Skills Sentinel**（安全哨兵）层，提供主动输入/输出清理及本地审计功能 (`AUDIT.md`)。
*   🚀 **并行子智能体**: 支持同时运行多个自主子智能体，每个子智能体可独立配置模型。
*   🌍 **真正可移植**: 跨 RISC-V、ARM 和 x86 架构的单二进制文件。
*   🦾 **AI 自举**: 核心代码通过自主 Agent 工作流不断精简和优化。

## 📢 新闻

2026-03-01 🎉 **PicoClaw v3.2.2 - 原生技能哨兵 (Skills Sentinel)**: 增加了其内部安全层 (`skills_sentinel.go`)，提供基于模式的实时保护，防止提示注入和系统泄漏。
2026-03-01 🎉 **PicoClaw v3.2 - Fail-Close 安全机制与稳定性**: 引入了更强大的安全策略。ExecTool 在启动时会严格验证安全规则。

2026-02-27 🎉 **PicoClaw v3.1 - 灾后恢复与任务锁**: 引入了原子级任务锁机制，防止 Agent 冲突；支持“启动再水化”，能从异常重启中快速恢复；优化了上下文压缩逻辑（安全提升至 32K token），彻底解决长代码任务中的上下文爆炸问题。


<img src="assets/compare.jpg" alt="PicoClaw" width="512">

## 🦾 演示

### 🛠️ 标准助手工作流

<table align="center">
  <tr align="center">
    <th><p align="center">🧩 全栈工程师模式</p></th>
    <th><p align="center">🗂️ 日志与规划管理</p></th>
    <th><p align="center">🔎 网络搜索与学习</p></th>
    <th><p align="center">🔧 通用事务处理</p></th>
  </tr>
  <tr>
    <td align="center"></td>
    <td align="center"></td>
    <td align="center"></td>
    <td align="center"></td>
  </tr>
  <tr>
    <td align="center">开发 • 部署 • 扩展</td>
    <td align="center">日程 • 自动化 • 记忆</td>
    <td align="center">发现 • 洞察 • 趋势</td>
    <td align="center">任务 • 支援 • 效率</td>
  </tr>
</table>

### 🚀 高级多智能体工作流（软件开发“梦想团队”）

利用子智能体架构部署一个完整的软件开发生命周期团队。

**“DevOps & QA”团队（由 [DeepSeek Reasoner](https://platform.deepseek.com) 提供支持）：**

*   **`project_manager` (项目经理)**：有权调用任何代理。监督全局目标并委派子任务。
*   **`senior_dev` (高级开发)**：技术核心。调用 QA 专家审核代码，或调用初级开发处理基础任务。
*   **`qa_specialist` (质量与运维)**：质量保证逻辑。测试代码、发现漏洞、提出修复方案并管理 GitHub 部署。
*   **`junior_fixer` (初级助手)**：在监督下负责小型修复、代码重构和文档编写。
*   **`general_worker` (基础勤务)**：多功能代理，负责通用任务、信息采集并为整个团队提供基础支撑。

**如何使用？**
只需通过 Telegram 或 CLI 向项目经理发送高级指令：
> *"PM，请让高级开发修复数据库 bug，并让 QA 在推送到 GitHub 之前验证构建。"*

PicoClaw 将自动管理层级结构：**项目经理 ➔ 高级开发 ➔ QA 专家（修复与发布）。**

> [!TIP]
> **查看示例：** 访问 `config_dev.example.json` 查看标准 DeepSeek 团队配置，`config_dev_multiple_models.example.json` 查看混合模型团队（OpenAI、Anthropic 和 DeepSeek），`config_context_management.example.json` 查看长代码任务中的上下文管理优化。


### 📱 在手机上轻松运行

picoclaw-agents 可以将你10年前的老旧手机废物利用，变身成为你的AI助理！快速指南:

1. **先去应用商店下载安装Termux**
2. **打开后执行指令**

```bash
# 注意: 下面的v0.1.1 可以换为你实际看到的最新版本
wget https://github.com/comgunner/picoclaw-agents/releases/download/v0.1.1/picoclaw-agents_Linux_arm64
chmod +x picoclaw-agents_Linux_arm64
pkg install proot
termux-chroot ./picoclaw-agents_Linux_arm64 onboard
```

然后跟随下面的“快速开始”章节继续配置picoclaw-agents即可使用！
<img src="assets/termux.jpg" alt="PicoClaw" width="512">

### 🐜 创新的低占用部署

PicoClaw 几乎可以部署在任何 Linux 设备上，从嵌入式板卡到高性能服务器。

🌟 更多部署案例敬请期待！

## 📦 安装

### 使用预编译二进制文件安装

从 [release](https://github.com/comgunner/picoclaw-agents/releases) 页面下载适用于您平台的固件。

### 从源码安装（最新特性，推荐开发使用）

```bash
git clone https://github.com/comgunner/picoclaw-agents.git

cd picoclaw-agents
make deps

# 编译，无需安装
make build

# 跨平台编译
make build-all

# 编译并安装
make install
```

## 🐳 Docker Compose

您也可以使用 Docker Compose 运行 PicoClaw，无需在本地安装任何软件。

```bash
# 1. 克隆此仓库
git clone https://github.com/comgunner/picoclaw-agents.git
cd picoclaw-agents

# 2. 设置您的 API 密钥
cp config/config.example.json config/config.json
vim config/config.json      # 设置 DISCORD_BOT_TOKEN, API 密钥等。

# 3. 编译并启动
docker compose --profile gateway up -d

> [!TIP]
> **Docker 用户**：默认情况下，Gateway 监听 `127.0.0.1`，宿主机无法直接访问。如果您需要访问健康检查端点或暴露端口，请在环境中设置 `PICOCLAW_GATEWAY_HOST=0.0.0.0` 或更新 `config.json`。


# 4. 查看日志
docker compose logs -f picoclaw-gateway

# 5. 停止
docker compose --profile gateway down
```

### 智能体模式 (单次运行)

```bash
# 提问
docker compose run --rm picoclaw-agents-agent -m "2+2等于几？"

# 交互模式
docker compose run --rm picoclaw-agents-agent
```

### 重新编译

```bash
docker compose --profile gateway build --no-cache
docker compose --profile gateway up -d
```

### 🚀 快速开始

> [!TIP]
> 在 `~/.picoclaw/config.json` 中设置您的 API 密钥。
> 获取 API 密钥：[OpenRouter](https://openrouter.ai/keys) (LLM) · [智谱 AI](https://open.bigmodel.cn/usercenter/proj-mgmt/apikeys) (LLM)
> 网页搜索是**可选**的——可获取免费的 [Tavily API](https://tavily.com) (每月1000次) 或 [Brave Search API](https://brave.com/search/api) (每月2000次)，或使用内置的自动回退方案。

**1. 初始化**

使用 `onboard` 命令并选择您喜欢的供应商模板来初始化工作区：

```bash
# 默认 (空/手动配置)
picoclaw-agents onboard

# 预配置模板：
picoclaw-agents onboard --openai      # 使用 OpenAI 模板 (o3-mini)
picoclaw-agents onboard --openrouter  # 使用 OpenRouter 模板 (openrouter/auto)
picoclaw-agents onboard --glm         # 使用 GLM-4.5-Flash 模板 (zhipu.ai)
picoclaw-agents onboard --qwen        # 使用通义千问模板 (阿里云国际版)
picoclaw-agents onboard --qwen_zh     # 使用通义千问模板 (阿里云国内版)
picoclaw-agents onboard --gemini      # 使用 Gemini 模板 (gemini-2.5-flash)
```

> [!TIP]
> **没有 API 余额？** 使用 `picoclaw-agents onboard --free` 即可立即开始使用 OpenRouter 的免费模型。只需在 [openrouter.ai](https://openrouter.ai) 创建免费账号并添加密钥 — 无需充值。

#### 🆓 免费模型

`--free` 选项配置三个 OpenRouter 免费模型，支持自动故障转移：

| 优先级 | 模型 | 上下文 | 说明 |
|--------|------|--------|------|
| 主要 | `openrouter/free` | 动态 | 自动选择当前最优免费模型 |
| 备用 1 | `stepfun/step-3.5-flash` | 256K | 适合长上下文任务 |
| 备用 2 | `deepseek/deepseek-v3.2-20251201` | 64K | 快速可靠的兜底模型 |

三个模型均通过 [OpenRouter](https://openrouter.ai) 路由 — 一个 API 密钥即可覆盖全部。

**2. 配置** (`~/.picoclaw/config.json`)

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

> **v3 核心 (多智能体架构)**：您现在可以启动隔离的**子智能体 (Subagents)** 来执行并行的后台任务。至关重要的是，**每个子智能体都可以使用完全不同的 LLM 模型**。如上配置所示，主智能体运行 `gpt4`，但它可以创建一个专门运行 `claude-sonnet-4.6` 的 `coder` 子智能体来同时处理复杂的编程任务！

> **新特性**：`model_list` 配置格式支持零代码添加供应商。详见 [模型配置 (model_list)](#模型配置-model_list)。
> `request_timeout` 是可选的，单位为秒。如果省略或设置为 `<= 0`，PicoClaw 将使用默认超时时间 (120s)。

**3. 获取 API 密钥**

* **LLM 供应商**: [DeepSeek](https://platform.deepseek.com) (推荐) · [OpenRouter](https://openrouter.ai/keys) · [智谱 AI](https://open.bigmodel.cn/usercenter/proj-mgmt/apikeys) · [Anthropic](https://console.anthropic.com) · [OpenAI](https://platform.openai.com) · [Gemini](https://aistudio.google.com/api-keys)
* **网页搜索** (可选): [Tavily](https://tavily.com) - 为 AI Agent 优化 (每月1000次) · [Brave Search](https://brave.com/search/api) - 提供免费档 (每月2000次)

### 💡 开发者推荐模型 (`backend_coder`)

对于重度编码任务，性能和逻辑是关键。我们建议在您的 `backend_coder` 代理中标准化使用以下模型：

*   **DeepSeek**: `deepseek-reasoner` (极致逻辑与性价比)
*   **OpenAI**: `o3-mini-2025-01-31` (高性能)
*   **OpenRouter.ai**: `Qwen3 Coder Plus`, `GPT-5.3-Codex` (编码通用性极佳)
*   **Anthropic**: `Claude Haiku 4.5` (快速且可靠)

> **注意**：完整配置模板请参考 `config.example.json`。

**4. 对话**

```bash
picoclaw-agents agent -m "2+2等于几？"
```

这就完成了！您在 2 分钟内就拥有了一个可以工作的 AI 助手。

---

## 💬 聊天应用

通过 Telegram, Discord, 钉钉, LINE 或 企业微信 与您的 PicoClaw 对话。

| 频道         | 设置                           |
| ------------ | ------------------------------ |
| **Telegram** | ⭐ 简单 (只需一个 token)        |
| **Discord**  | ⭐ 简单 (bot token + intents)   |
| **QQ**       | ⭐ 简单 (AppID + AppSecret)     |
| **钉钉**     | ⚙️ 中等 (应用凭据)              |
| **LINE**     | ⚙️ 中等 (凭据 + webhook URL)    |
| **企业微信** | ⚙️ 中等 (CorpID + webhook 设置) |

<details>
<summary><b>Telegram</b> (推荐)</summary>

**1. 创建机器人**

* 在 Telegram 中搜索 `@BotFather`
* 发送 `/newbot`，按提示操作
* 复制 token

**2. 配置**

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

> 通过 Telegram 的 `@userinfobot` 获取您的用户 ID。

**3. 启动**

```bash
picoclaw-agents gateway
```

</details>

<details>
<summary><b>Discord</b></summary>

**1. 创建机器人**

* 访问 <https://discord.com/developers/applications>
* 创建应用 → Bot → Add Bot
* 复制机器人 token

**2. 开启 Intents**

* 在 Bot 设置中，开启 **MESSAGE CONTENT INTENT**
* (可选) 开启 **SERVER MEMBERS INTENT** (如果您计划基于成员数据进行白名单验证)

**3. 获取您的 User ID**
* Discord 设置 → 高级 → 开启 **开发者模式 (Developer Mode)**
* 右键点击您的头像 → **复制用户 ID**

**4. 配置**

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

**5. 邀请机器人**

* OAuth2 → URL Generator
* 作用域 (Scopes): `bot`
* 机器人权限: `Send Messages`, `Read Message History`
* 打开生成的邀请链接，将机器人添加到您的服务器

**可选：仅提醒模式 (Mention-only mode)**

设置 `"mention_only": true` 使机器人仅在被 @ 提到时响应。适用于您希望机器人仅在被明确调用时才响应的共享服务器。

**6. 启动**

```bash
picoclaw-agents gateway
```

</details>

<details>
<summary><b>QQ</b></summary>

**1. 创建机器人**

- 访问 [QQ 开放平台](https://q.qq.com/#)
- 创建应用 → 获取 **AppID** 和 **AppSecret**

**2. 配置**

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

> 将 `allow_from` 留空以允许所有用户，或指定 QQ 号以限制访问。

**3. 启动**

```bash
picoclaw-agents gateway
```

</details>

<details>
<summary><b>钉钉</b></summary>

**1. 创建机器人**

* 访问 [开放平台](https://open.dingtalk.com/)
* 创建内部应用 (Internal app)
* 复制 Client ID 和 Client Secret

**2. 配置**

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

> 将 `allow_from` 留空以允许所有用户，或指定钉钉用户 ID 以限制访问。

**3. 启动**

```bash
picoclaw-agents gateway
```
</details>

<details>
<summary><b>LINE</b></summary>

**1. 创建 LINE 官方帐号**

- 访问 [LINE Developers Console](https://developers.line.biz/)
- 创建供应商 (provider) → 创建 Messaging API 频道
- 复制 **Channel Secret** 和 **Channel Access Token**

**2. 配置**

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

**3. 设置 Webhook URL**

LINE 要求 Webhook 使用 HTTPS。请使用反向代理或隧道：

```bash
# ngrok 示例
ngrok http 18791
```

然后在 LINE Developers Console 中将 Webhook URL 设置为 `https://your-domain/webhook/line` 并开启 **Use webhook**。

**4. 启动**

```bash
picoclaw-agents gateway
```

> 在群聊中，机器人仅在被 @ 提到时响应。回复将引用原始消息。

> **Docker Compose**：在 `picoclaw-gateway` 服务中添加 `ports: ["18791:18791"]` 以暴露 Webhook 端口。

</details>

<details>
<summary><b>企业微信 (WeCom)</b></summary>

PicoClaw 支持两类企业微信集成：

**方案 1：企业微信机器人 (智能机器人)** - 设置更简单，支持群聊。
**方案 2：企业微信自建应用 (自建应用)** - 功能更多，支持主动推送。

详见 [企业微信应用配置指南](docs/wecom-app-configuration.md) 了解详细设置步骤。

**快速设置 - 企业微信机器人：**

**1. 创建机器人**

* 进入企业微信管理后台 → 群聊 → 添加群机器人
* 复制 webhook URL (格式: `https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=xxx`)

**2. 配置**

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

**快速设置 - 企业微信自建应用：**

**1. 创建应用**

* 进行企业微信管理后台 → 应用管理 → 创建应用
* 复制 **AgentId** 和 **Secret**
* 进入“我公司”页面，复制 **CorpID**
**2. 配置接收消息**

* 在应用详情中，点击“接收消息” → “设置 API (Set API)”
* 将 URL 设置为 `http://your-server:18792/webhook/wecom-app`
* 生成 **Token** 和 **EncodingAESKey**

**3. 配置**

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

**4. 启动**

```bash
picoclaw-agents gateway
```

> **注意**：企业微信应用需要为 Webhook 回调开放 18792 端口。请使用反向代理启用 HTTPS。

</details>

## 加入智能体社交网络

只需通过 CLI 或任何集成的聊天应用发送一条消息，即可将您的 PicoClaw 连接到智能体社交网络。

**阅读 `https://clawdchat.ai/skill.md` 并按照说明加入 [ClawdChat.ai](https://clawdchat.ai)**

## ⚙️ 配置

配置文件：`~/.picoclaw/config.json`

### 工作区布局

PicoClaw 将数据存储在您配置的工作区中（默认：`~/.picoclaw/workspace`）：

```
~/.picoclaw/workspace/
├── sessions/          # 对话会话与历史记录
├── memory/           # 长期记忆 (MEMORY.md)
├── state/            # 持久状态 (上次使用的频道等)
├── cron/             # 定时任务数据库
├── skills/           # 自定义技能
├── AGENTS.md         # 智能体行为指南
├── HEARTBEAT.md      # 定期任务提示词 (每 30 分钟检查一次)
├── IDENTITY.md       # 智能体身份
├── SOUL.md           # 智能体灵魂
├── TOOLS.md          # 工具描述
└── USER.md           # 用户偏好
```

### 🔒 安全沙箱

PicoClaw 默认在沙箱环境中运行。Agent 只能访问已配置工作区内的文件并执行其中的命令。

#### 默认配置

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

| 选项                    | 默认值                  | 描述                          |
| ----------------------- | ----------------------- | ----------------------------- |
| `workspace`             | `~/.picoclaw/workspace` | Agent 的工作目录              |
| `restrict_to_workspace` | `true`                  | 限制文件/命令访问仅限于工作区 |

#### 受保护的工具

当 `restrict_to_workspace: true` 时，以下工具将受到沙箱限制：

| 工具          | 功能       | 限制                   |
| ------------- | ---------- | ---------------------- |
| `read_file`   | 读取文件   | 仅限工作区内的文件     |
| `write_file`  | 写入文件   | 仅限工作区内的文件     |
| `list_dir`    | 列出目录   | 仅限工作区内的目录     |
| `edit_file`   | 编辑文件   | 仅限工作区内的文件     |
| `append_file` | 追加至文件 | 仅限工作区内的文件     |
| `exec`        | 执行命令   | 命令路径必须在工作区内 |

#### 额外的 Exec 保护

即使 `restrict_to_workspace: false`，`exec` 工具也会拦截以下危险命令：

* `rm -rf`, `del /f`, `rmdir /s` — 批量删除
* `format`, `mkfs`, `diskpart` — 磁盘格式化
* `dd if=` — 磁盘镜像写入
* 写入至 `/dev/sd[a-z]` — 直接写入磁盘
* `shutdown`, `reboot`, `poweroff` — 系统关机/重启
* Fork 炸弹 `:(){ :|:& };:`

#### 错误示例

```
[ERROR] tool: Tool execution failed
{tool=exec, error=Command blocked by safety guard (path outside working dir)}
```

```
[ERROR] tool: Tool execution failed
{tool=exec, error=Command blocked by safety guard (dangerous pattern detected)}
```

#### 禁用限制 (安全风险)

如果您需要 Agent 访问工作区外的路径：

**方法 1：配置文件**

```json
{
  "agents": {
    "defaults": {
      "restrict_to_workspace": false
    }
  }
}
```

**方法 2：环境变量**

```bash
export PICOCLAW_AGENTS_DEFAULTS_RESTRICT_TO_WORKSPACE=false
```

> ⚠️ **警告**：禁用此限制将允许 Agent 访问您系统上的任何路径。请仅在受控环境下谨慎使用。

#### 安全边界一致性

`restrict_to_workspace` 设置在所有执行路径中均一致生效：

| 执行路径                | 安全边界                  |
| ----------------------- | ------------------------- |
| 主智能体 (Main Agent)   | `restrict_to_workspace` ✅ |
| 子智能体 / 派生 (Spawn) | 继承相同的限制 ✅          |
| 定期任务 (Heartbeat)    | 继承相同的限制 ✅          |

所有路径共享相同的工作区限制——无法通过子智能体或定时任务绕过安全边界。

### 定期任务 (Heartbeat)

PicoClaw 可以自动执行定期任务。在工作区中创建 `HEARTBEAT.md` 文件：

```markdown
# 定期任务

- 检查我的电子邮件是否有重要消息
- 查看我的日历是否有即将参加的活动
- 检查天气预报
```

Agent 将每 30 分钟 (可配置) 读取一次此文件，并使用可用工具执行任务。

#### 使用 Spawn 执行异步任务

对于耗时任务 (如网页搜索、API 调用)，请使用 `spawn` 工具创建**子智能体**：

```markdown
# 定期任务

## 快速任务 (直接响)

- 报告当前时间

## 耗时任务 (使用 spawn 执行异步任务)

- 在网页上搜索 AI 新闻并进行总结
- 检查电子邮件并报告重要消息
```

**关键行为：**

| 特性                        | 描述                                         |
| --------------------------- | -------------------------------------------- |
| **spawn**                   | 创建异步子智能体，不阻塞定期任务 (Heartbeat) |
| **独立上下文**              | 子智能体拥有自己的上下文，无会话历史记录     |
| **消息工具 (message tool)** | 子智能体通过消息工具直接与用户沟通           |
| **非阻塞**                  | 派生后，定期任务将继续执行下一项任务         |

#### 子智能体通信原理

```
定期任务触发
    ↓
Agent 读取 HEARTBEAT.md
    ↓
耗时任务：启动子智能体 (spawn)
    ↓                           ↓
继续执行下一项任务           子智能体独立工作
    ↓                           ↓
所有任务完成               子智能体使用“消息”工具
    ↓                           ↓
回复 HEARTBEAT_OK          用户直接收到结果
```

子智能体可以访问工具 (消息、网页搜索等)，并能独立与用户沟通，无需经过主智能体。

**配置：**

```json
{
  "heartbeat": {
    "enabled": true,
    "interval": 30
  }
}
```

| 选项       | 默认值 | 描述                                 |
| ---------- | ------ | ------------------------------------ |
| `enabled`  | `true` | 开启/关闭定期任务                    |
| `interval` | `30`   | 检查间隔时间，单位为分钟 (最小值: 5) |

**环境变量：**

* `PICOCLAW_HEARTBEAT_ENABLED=false` 禁用
* `PICOCLAW_HEARTBEAT_INTERVAL=60` 修改间隔

### 供应商

> [!NOTE]
> Groq 通过 Whisper 提供免费的语音转录。如果配置了 Groq，Telegram 语音消息将被自动转录为文字。

| 供应商               | 用途                         | 获取 API 密钥                                                        |
| -------------------- | ---------------------------- | -------------------------------------------------------------------- |
| `gemini`             | LLM (Gemini 直连)            | [aistudio.google.com](https://aistudio.google.com)                   |
| `zhipu`              | LLM (智谱 AI 直连)           | [bigmodel.cn](https://bigmodel.cn)                                   |
| `openrouter(待测试)` | LLM (推荐，可访问所有模型)   | [openrouter.ai](https://openrouter.ai)                               |
| `anthropic(待测试)`  | LLM (Claude 直连)            | [console.anthropic.com](https://console.anthropic.com)               |
| `openai(待测试)`     | LLM (GPT 直连)               | [platform.openai.com](https://platform.openai.com)                   |
| `deepseek(待测试)`   | LLM (DeepSeek 直连)          | [platform.deepseek.com](https://platform.deepseek.com)               |
| `qwen`               | LLM (通义千问直连)           | [dashscope.console.aliyun.com](https://dashscope.console.aliyun.com) |
| `groq`               | LLM + **语音转录** (Whisper) | [console.groq.com](https://console.groq.com)                         |
| `cerebras`           | LLM (Cerebras 直连)          | [cerebras.ai](https://cerebras.ai)                                   |

### 模型配置 (model_list)

> **新特性**：PicoClaw 现在采用以**模型为中心 (model-centric)** 的配置方式。只需通过 `vendor/model` 格式 (例如 `zhipu/glm-4.5-flash`) 即可添加新供应商——**无需修改任何代码！**

此设计还支持**多智能体 (multi-agent)** 及灵活的供应商选择：

- **不同智能体，不同供应商**：每个 Agent 都可以使用各自的 LLM 供应商。
- **模型回退 (fallback)**：配置主模型和备用模型以增强弹性。
- **负载均衡**：在多个端点间分配请求。
- **集中管理**：在一个地方管理所有供应商。

#### 📋 所有支持的厂商

| 厂商                | `model` 前缀      | 默认 API 基础路径                                   | 协议      | API 密钥                                                          |
| ------------------- | ----------------- | --------------------------------------------------- | --------- | ----------------------------------------------------------------- |
| **OpenAI**          | `openai/`         | `https://api.openai.com/v1`                         | OpenAI    | [获取 Key](https://platform.openai.com)                           |
| **Anthropic**       | `anthropic/`      | `https://api.anthropic.com/v1`                      | Anthropic | [获取 Key](https://console.anthropic.com)                         |
| **智谱 AI (GLM)**   | `zhipu/`          | `https://open.bigmodel.cn/api/paas/v4`              | OpenAI    | [获取 Key](https://open.bigmodel.cn/usercenter/proj-mgmt/apikeys) |
| **DeepSeek**        | `deepseek/`       | `https://api.deepseek.com/v1`                       | OpenAI    | [获取 Key](https://platform.deepseek.com)                         |
| **Google Gemini**   | `gemini/`         | `https://generativelanguage.googleapis.com/v1beta`  | OpenAI    | [获取 Key](https://aistudio.google.com/api-keys)                  |
| **Groq**            | `groq/`           | `https://api.groq.com/openai/v1`                    | OpenAI    | [获取 Key](https://console.groq.com)                              |
| **Moonshot**        | `moonshot/`       | `https://api.moonshot.cn/v1`                        | OpenAI    | [获取 Key](https://platform.moonshot.cn)                          |
| **通义千问 (Qwen)** | `qwen/`           | `https://dashscope.aliyuncs.com/compatible-mode/v1` | OpenAI    | [获取 Key](https://dashscope.console.aliyun.com)                  |
| **NVIDIA**          | `nvidia/`         | `https://integrate.api.nvidia.com/v1`               | OpenAI    | [获取 Key](https://build.nvidia.com)                              |
| **Ollama**          | `ollama/`         | `http://localhost:11434/v1`                         | OpenAI    | 本地 (无需 Key)                                                   |
| **OpenRouter**      | `openrouter/`     | `https://openrouter.ai/api/v1`                      | OpenAI    | [获取 Key](https://openrouter.ai/keys)                            |
| **VLLM**            | `vllm/`           | `http://localhost:8000/v1`                          | OpenAI    | 本地                                                              |
| **Cerebras**        | `cerebras/`       | `https://api.cerebras.ai/v1`                        | OpenAI    | [获取 Key](https://cerebras.ai)                                   |
| **火山引擎**        | `volcengine/`     | `https://ark.cn-beijing.volces.com/api/v3`          | OpenAI    | [获取 Key](https://console.volcengine.com)                        |
| **神算云**          | `shengsuanyun/`   | `https://router.shengsuanyun.com/api/v1`            | OpenAI    | -                                                                 |
| **Antigravity**     | `antigravity/`    | Google Cloud                                        | 自定义    | 仅限 OAuth                                                        |
| **GitHub Copilot**  | `github-copilot/` | `localhost:4321`                                    | gRPC      | -                                                                 |

#### 基础配置

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

#### 厂商特定示例

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

**Anthropic (带 API 密钥)**

```json
{
  "model_name": "claude-sonnet-4.6",
  "model": "anthropic/claude-sonnet-4.6",
  "api_key": "sk-ant-your-key"
}
```

> 运行 `picoclaw-agents auth login --provider anthropic` 粘贴您的 API token。

**Ollama (本地)**

```json
{
  "model_name": "llama3",
  "model": "ollama/llama3"
}
```

**自定义代理 (Proxy)/API**

```json
{
  "model_name": "my-custom-model",
  "model": "openai/custom-model",
  "api_base": "https://my-proxy.com/v1",
  "api_key": "sk-...",
  "request_timeout": 300
}
```

#### 负载均衡

为同一个模型名称配置多个端点——PicoClaw 会在它们之间自动执行轮询 (Round-robin)：

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

#### 从旧版 `providers` 配置迁移

旧版 `providers` 配置已**弃用 (deprecated)**，但为了向后兼容仍受支持。

**旧版配置 (不推荐)：**

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

**新版配置 (推荐)：**

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

详细迁移指南请参考 [docs/migration/model-list-migration.md](docs/migration/model-list-migration.md)。

### 供应商架构

PicoClaw 按协议族路由供应商：

- OpenAI 兼容协议：OpenRouter, 兼容 OpenAI 的网关, Groq, 智谱 AI 和 vLLM 风格端点。
- Anthropic 协议：Claude 原生 API 行为。
- Codex/OAuth 路径：OpenAI OAuth/token 认证路径。

这使得运行环境保持轻量，同时添加新的兼容 OpenAI 的后端仅需进行配置操作 (`api_base` + `api_key`)。

<details>
<summary><b>智谱 AI</b></summary>

**1. 获取 API 密钥和基础 URL**

* 获取 [API key](https://bigmodel.cn/usercenter/proj-mgmt/apikeys)

**2. 配置**

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

**3. 运行**

```bash
picoclaw-agents agent -m "Hello"
```

</details>

<details>
<summary><b>完整配置示例</b></summary>

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

## CLI 参考

| 命令                      | 描述               |
| ------------------------- | ------------------ |
| `picoclaw-agents onboard`        | 初始化配置与工作区 |
| `picoclaw-agents agent -m "..."` | 与 Agent 对话      |
| `picoclaw-agents agent`          | 交互式对话模式     |
| `picoclaw-agents gateway`        | 启动网关 (Gateway) |
| `picoclaw-agents status`         | 查看状态           |
| `picoclaw-agents cron list`      | 列出所有定时任务   |
| `picoclaw-agents cron add ...`   | 添加定时任务       |

### 定时任务 / 提醒

PicoClaw 通过 `cron` 工具支持定时提醒和循环任务：

* **单次提醒**："10分钟后提醒我" → 10分钟后触发一次
* **循环任务**："每2小时提醒我一次" → 每2小时触发一次
* **Cron 表达式**："每天上午9点提醒我" → 使用 cron 表达式

任务存储在 `~/.picoclaw/workspace/cron/` 并自动处理。

### Binance 集成 (原生工具 + MCP)

PicoClaw 在 `agent` 模式内置 Binance 原生工具：

* `binance_get_ticker_price`（公开市场行情）
* `binance_get_spot_balance`（签名端点，需要 API key/secret）

在 `~/.picoclaw/config.json` 中配置密钥：

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

使用示例：

```bash
picoclaw-agents agent -m "Use binance_get_ticker_price with symbol BTCUSDT and return only the numeric price."
picoclaw-agents agent -m "Use binance_get_spot_balance and show my non-zero balances."
```

未配置 API keys 时：

* `binance_get_ticker_price` 仍会通过 Binance 公共端点返回结果，并附带公共端点提示。
* `binance_get_spot_balance` 会提示缺少密钥，并建议使用公共 `curl` 查询。

可选 MCP 服务模式（供 MCP 客户端使用）：

```bash
picoclaw-agents util binance-mcp-server
```

`mcp_servers` 配置示例（请使用安装/onboard 生成的 `picoclaw-agents` 绝对路径）：

```json
{
  "mcp_servers": {
    "binance": {
      "enabled": true,
      "command": "/picoclaw/absolute/path",
      "args": ["util", "binance-mcp-server"]
    }
  }
}
```

## 🤝 贡献与路线图

查看完整的 [路线图 (Roadmap)](ROADMAP.md)。

Discord: [即将推出 / Coming Soon]


## 🐛 故障排除

### 网页搜索提示 \"API key configuration issue\"

如果您尚未配置搜索 API 密钥，这是正常现象。PicoClaw 将提供有用的链接供手动搜索。

要开启网页搜索：

1. **方案 1 (推荐)**：在 [https://brave.com/search/api](https://brave.com/search/api) 获取免费 API 密钥 (每月2000次)，以获得最佳结果。
2. **方案 2 (无信用卡)**：如果您没有密钥，我们将自动回退到 **DuckDuckGo** (无需密钥)。

如果使用 Brave，请将密钥添加到 `~/.picoclaw/config.json`：

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

### 遇到内容过滤错误

某些供应商 (如智谱 AI) 设有内容过滤。请尝试重新描述您的问题或切换其他模型。

### Telegram 机器人提示 \"Conflict: terminated by other getUpdates\"

当运行了多个机器人实例时会发生这种情况。请确保一次只运行一个 `picoclaw-agents gateway`。

---

## 📝 API 密钥对比

| 服务             | 免费档额度    | 用途 (Use Case)                |
| ---------------- | ------------- | ------------------------------ |
| **OpenRouter**   | 20万 token/月 | 多模型访问 (Claude, GPT-4 等)  |
| **智谱 AI**      | 提供免费档    | glm-4.5-flash (最适合中文用户) |
| **Brave Search** | 2000次查询/月 | 网页搜索功能                   |
| **Groq**         | 提供免费档    | 极速推理 (Llama, Mixtral)      |
| **Cerebras**     | 提供免费档    | 极速推理 (Llama, Qwen 等)      |

## ⚠️ 免责声明

本软件按“原样”提供，不提供任何形式的明示或暗示保证，包括但不限于对适销性、特定用途的适用性和非侵权性的保证。在任何情况下，对于因本软件或本软件的使用或其他交易而产生的、由本软件引起的或与本软件有关的任何索赔、损害或其他责任，无论是在合同诉讼、侵权诉讼还是其他诉讼中，本分叉的作者或版权所有者均不承担任何责任。**使用风险自负。**
