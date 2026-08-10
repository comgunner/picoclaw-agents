# OpenCode Zen — Available Models

> **Last Updated:** August 9, 2026
> **API Endpoint:** `https://opencode.ai/zen/v1`
> **API Key:** https://opencode.ai/auth

---

## Overview

OpenCode Zen is a curated AI model gateway maintained by the OpenCode team. It provides access to 50+ tested and verified models with pay-as-you-go pricing. All models use an OpenAI-compatible API.

**Quick login:**
```bash
picoclaw-agents auth login --provider opencode
```

---

## Free Models

These models are free for a limited time. No credit card required.

| Model ID | Name | Provider | Notes |
|----------|------|----------|-------|
| `mimo-v2.5-free` | MiMo V2.5 Free | Xiaomi | Strong coding capabilities |
| `deepseek-v4-flash-free` | DeepSeek V4 Flash Free | DeepSeek | Fast and efficient |
| `nemotron-3-ultra-free` | Nemotron 3 Ultra Free | NVIDIA | 550B parameters |
| `ling-3.0-tiny-free` | Ling 3.0 Tiny Free | InclusionAI | Lightweight, fast |
| `north-mini-code-free` | North Mini Code Free | North | Code-focused |
| `laguna-s-2.1-free` | Laguna S 2.1 Free | Laguna | New |
| `longcat-2.0-free` | LongCat 2.0 Free | LongCat | New |
| `big-pickle` | Big Pickle | Stealth | Free trial |

> ⚠️ **Note:** Free models may have data used for improvement during the trial period.

---

## Paid Models

Pay-as-you-go pricing per 1M tokens.

### GPT Models (OpenAI)

| Model ID | Name | Input | Output |
|----------|------|-------|--------|
| `gpt-5.6-sol` | GPT 5.6 Sol | $5.00 | $30.00 |
| `gpt-5.6-terra` | GPT 5.6 Terra | $2.00 | $12.00 |
| `gpt-5.6-luna` | GPT 5.6 Luna | $0.20 | $1.20 |
| `gpt-5.5` | GPT 5.5 | $5.00 | $30.00 |
| `gpt-5.5-pro` | GPT 5.5 Pro | $30.00 | $180.00 |
| `gpt-5.4` | GPT 5.4 | $2.50 | $15.00 |
| `gpt-5.4-pro` | GPT 5.4 Pro | $30.00 | $180.00 |
| `gpt-5.4-mini` | GPT 5.4 Mini | $0.75 | $4.50 |
| `gpt-5.4-nano` | GPT 5.4 Nano | $0.20 | $1.25 |
| `gpt-5.3-codex` | GPT 5.3 Codex | $1.75 | $14.00 |
| `gpt-5.3-codex-spark` | GPT 5.3 Codex Spark | $1.75 | $14.00 |
| `gpt-5.2` | GPT 5.2 | $1.75 | $14.00 |
| `gpt-5.1` | GPT 5.1 | $1.07 | $8.50 |
| `gpt-5` | GPT 5 | $1.07 | $8.50 |
| `gpt-5-nano` | GPT 5 Nano | $0.05 | $0.40 |

### Claude Models (Anthropic)

| Model ID | Name | Input | Output |
|----------|------|-------|--------|
| `claude-fable-5` | Claude Fable 5 | $10.00 | $50.00 |
| `claude-opus-5` | Claude Opus 5 | $5.00 | $25.00 |
| `claude-opus-4-8` | Claude Opus 4.8 | $5.00 | $25.00 |
| `claude-opus-4-7` | Claude Opus 4.7 | $5.00 | $25.00 |
| `claude-opus-4-6` | Claude Opus 4.6 | $5.00 | $25.00 |
| `claude-opus-4-5` | Claude Opus 4.5 | $5.00 | $25.00 |
| `claude-sonnet-5` | Claude Sonnet 5 | $2.00 | $10.00 |
| `claude-sonnet-4-6` | Claude Sonnet 4.6 | $3.00 | $15.00 |
| `claude-sonnet-4-5` | Claude Sonnet 4.5 | $3.00 | $15.00 |
| `claude-haiku-4-5` | Claude Haiku 4.5 | $1.00 | $5.00 |

### Gemini Models (Google)

| Model ID | Name | Input | Output |
|----------|------|-------|--------|
| `gemini-3.6-flash` | Gemini 3.6 Flash | $1.50 | $7.50 |
| `gemini-3.5-flash` | Gemini 3.5 Flash | $1.50 | $9.00 |
| `gemini-3.5-flash-lite` | Gemini 3.5 Flash Lite | $0.30 | $2.50 |
| `gemini-3.1-pro` | Gemini 3.1 Pro | $2.00 | $12.00 |
| `gemini-3-flash` | Gemini 3 Flash | $0.50 | $3.00 |

### Grok Models (xAI)

| Model ID | Name | Input | Output |
|----------|------|-------|--------|
| `grok-4.5` | Grok 4.5 | $2.00 | $6.00 |
| `grok-build-0.1` | Grok Build 0.1 | $1.00 | $2.00 |

### DeepSeek Models

| Model ID | Name | Input | Output |
|----------|------|-------|--------|
| `deepseek-v4-pro` | DeepSeek V4 Pro | $1.74 | $3.48 |
| `deepseek-v4-flash` | DeepSeek V4 Flash | $0.14 | $0.28 |

### Qwen Models (Alibaba)

| Model ID | Name | Input | Output |
|----------|------|-------|--------|
| `qwen3.7-max` | Qwen3.7 Max | $2.50 | $7.50 |
| `qwen3.7-plus` | Qwen3.7 Plus | $0.40 | $1.60 |
| `qwen3.6-plus` | Qwen3.6 Plus | $0.50 | $3.00 |
| `qwen3.5-plus` | Qwen3.5 Plus | $0.20 | $1.20 |

### Kimi Models (Moonshot)

| Model ID | Name | Input | Output |
|----------|------|-------|--------|
| `kimi-k3` | Kimi K3 | $3.00 | $15.00 |
| `kimi-k2.7-code` | Kimi K2.7 Code | $0.95 | $4.00 |
| `kimi-k2.6` | Kimi K2.6 | $0.95 | $4.00 |
| `kimi-k2.5` | Kimi K2.5 | $0.60 | $3.00 |

### GLM Models (Zhipu)

| Model ID | Name | Input | Output |
|----------|------|-------|--------|
| `glm-5.2` | GLM 5.2 | $1.40 | $4.40 |
| `glm-5.1` | GLM 5.1 | $1.40 | $4.40 |
| `glm-5` | GLM 5 | $1.00 | $3.20 |

### MiniMax Models

| Model ID | Name | Input | Output |
|----------|------|-------|--------|
| `minimax-m3` | MiniMax M3 | $0.30 | $1.20 |
| `minimax-m2.7` | MiniMax M2.7 | $0.30 | $1.20 |

---

## Adding Paid Models to PicoClaw

To use a paid model, add it to your config:

```json
{
  "model_list": [
    {
      "model_name": "opencode-claude-sonnet-5",
      "model": "opencode/claude-sonnet-5",
      "api_base": "https://opencode.ai/zen/v1",
      "api_key": "YOUR_OPENCODE_API_KEY" // pragma: allowlist secret
    }
  ]
}
```

Then select it:
```bash
/model opencode-claude-sonnet-5
```

Or in Telegram:
```
/model opencode/claude-sonnet-5
```

---

## Pricing Notes

- Credit card fees passed along at cost (4.4% + $0.30 per transaction)
- Auto-reload: If balance goes below $5, Zen reloads $20 automatically
- Monthly limits available for workspace and per-member
- All models hosted in the US with zero-retention policy (except free trial models)

---

## References

- [OpenCode Zen Docs](https://opencode.ai/docs/zen/)
- [OpenCode Providers](https://opencode.ai/docs/providers/)
- [GitHub Issue #752](https://github.com/sipeed/picoclaw/issues/752)
- [Get API Key](https://opencode.ai/auth)
