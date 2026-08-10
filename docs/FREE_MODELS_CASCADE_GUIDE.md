# Free Models Cascade Guide

> **Survival mechanism** for your bot — automatic fallback between free models to keep your bot online 24/7.

---

## What is Cascade Fallback?

Cascade fallback is a **survival mechanism** for your bot. When the primary free model fails (rate limit, timeout, unavailable), the system automatically tries the next free model in the cascade. This ensures your bot **stays online 24/7** without manual intervention.

## Why Free Models (No Credit Card Required)

We support **two providers** with free tiers:
- **OpenCode Zen** — 8 curated free models, tested and verified
- **OpenRouter** — 200+ models with automatic fallback

Both offer:
- ✅ **Free tier models** without requiring a credit card
- ✅ **Automatic fallback** between multiple free models
- ✅ **OpenAI-compatible API** (drop-in replacement)
- ✅ **No vendor lock-in** (switch models anytime)

## Available Free Models (August 2026)

### OpenCode Zen Free Models (Priority)

| Priority | Model | Provider | Use Case |
|----------|-------|----------|----------|
| 1 | `opencode/mimo-v2.5-free` | Xiaomi | MiMo V2.5, strong coding |
| 2 | `opencode/deepseek-v4-flash-free` | DeepSeek | Fast, efficient |
| 3 | `opencode/nemotron-3-ultra-free` | NVIDIA | 550B parameters |
| 4 | `opencode/ling-3.0-tiny-free` | InclusionAI | Lightweight, fast |
| 5 | `opencode/north-mini-code-free` | North | Code-focused |
| 6 | `opencode/laguna-s-2.1-free` | Laguna | New |
| 7 | `opencode/longcat-2.0-free` | LongCat | New |
| 8 | `opencode/big-pickle` | Stealth | Free trial |

### OpenRouter Free Models (Fallback)

| Priority | Model | Provider | Use Case |
|----------|-------|----------|----------|
| 9 | `openrouter/free` | OpenRouter | Auto-select best free model |
| 10 | `nvidia/nemotron-3-ultra-550b-a55b:free` | NVIDIA | Large context, complex reasoning |
| 11 | `openai/gpt-oss-20b:free` | OpenAI | Balanced performance |
| 12 | `inclusionai/ling-3.0-tiny:free` | InclusionAI | Fast responses, low latency |
| 13 | `meta-llama/llama-3.1-8b-instruct:free` | Meta | General purpose |

## Paid Fallback Models

If all free models fail, the cascade can fall back to paid models:

| Model | Provider | Cost | Use Case |
|-------|----------|------|----------|
| `deepseek/deepseek-v4-flash-0731` | DeepSeek | $0.27/M tokens | High quality, cost-effective |
| `xiaomi/mimo-v2.5` | Xiaomi | $0.10/M tokens | Budget option |

## Configuration

### Basic Configuration (Free Only)

```json
{
  "agents": {
    "defaults": {
      "provider": "opencode",
      "model_name": "opencode-mimo-free",
      "model": "opencode/mimo-v2.5-free",
      "cascade": {
        "enabled": true,
        "text_models": [
          "opencode/mimo-v2.5-free",
          "opencode/deepseek-v4-flash-free",
          "opencode/nemotron-3-ultra-free",
          "opencode/ling-3.0-tiny-free",
          "openrouter/free",
          "nvidia/nemotron-3-ultra-550b-a55b:free",
          "openai/gpt-oss-20b:free",
          "meta-llama/llama-3.1-8b-instruct:free"
        ]
      }
    }
  }
}
```

### With Paid Fallback

```json
{
  "agents": {
    "defaults": {
      "provider": "openrouter",
      "model_name": "openrouter-free",
      "model": "openrouter/free",
      "cascade": {
        "enabled": true,
        "text_models": [
          "openrouter/free",
          "nvidia/nemotron-3-ultra-550b-a55b:free",
          "openai/gpt-oss-20b:free",
          "inclusionai/ling-3.0-tiny:free",
          "meta-llama/llama-3.1-8b-instruct:free"
        ],
        "paid_fallbacks": [
          "deepseek/deepseek-v4-flash-0731",
          "xiaomi/mimo-v2.5"
        ]
      }
    }
  }
}
```

## How Cascade Works

```
User Request
    ↓
Try Model 1: openrouter/free
    ↓ (if fails)
Try Model 2: nvidia/nemotron-3-ultra-550b-a55b:free
    ↓ (if fails)
Try Model 3: openai/gpt-oss-20b:free
    ↓ (if fails)
Try Model 4: inclusionai/ling-3.0-tiny:free
    ↓ (if fails)
Try Model 5: meta-llama/llama-3.1-8b-instruct:free
    ↓ (if fails)
Try Paid Fallback: deepseek/deepseek-v4-flash-0731
    ↓ (if fails)
Try Paid Fallback: xiaomi/mimo-v2.5
    ↓
Return error if all fail
```

## Importance for Bot Uptime

### Without Cascade
- Single model fails → bot goes offline
- Manual intervention required
- User experience degraded

### With Cascade
- Automatic failover between models
- Bot stays online 24/7
- No manual intervention needed
- Seamless user experience

## Image Generation Cascade

The same cascade pattern applies to image generation:

```json
{
  "image_gen": {
    "provider": "openrouter",
    "openrouter_image_model": "krea/krea-2-medium-turbo",
    "openrouter_text_model": "openrouter/free"
  }
}
```

### Image Models
| Priority | Model | Use Case |
|----------|-------|----------|
| 1 | `krea/krea-2-medium-turbo` | Primary (cheapest) |
| 2 | `bytedance-seed/seedream-4.5` | Alternative |
| 3 | `x-ai/grok-imagine-image-quality` | High quality |

## Monitoring

### Check Cascade Status
```bash
# List available models
picoclaw models

# Test specific model
picoclaw chat --model openrouter/free "Hello"

# Check OpenRouter status
curl http://localhost:18792/health
```

### Logs
```bash
# Watch cascade logs
journalctl -u picoclaw.service -f | grep cascade

# Check model failures
journalctl -u picoclaw.service | grep "model.*failed"
```

## Troubleshooting

### All Free Models Fail
1. Check OpenRouter status: https://status.openrouter.ai
2. Check API key: `picoclaw auth login`
3. Check rate limits: Free models have per-minute limits
4. Consider adding paid fallbacks

### High Latency
1. Cascade tries multiple models sequentially
2. First successful response is returned
3. Consider reducing cascade list for faster response

### Rate Limited
1. Free models have rate limits
2. Cascade automatically tries next model
3. Wait for rate limit reset or use paid fallback

## Best Practices

1. **Always enable cascade** for production bots
2. **Start with free models** to minimize cost
3. **Add paid fallbacks** for critical applications
4. **Monitor logs** for model failures
5. **Test cascade** regularly with `picoclaw chat`
6. **Keep models updated** as providers add/remove free tiers

## Related Documentation

- [OpenRouter Free Models](OPENROUTER_FREE.md)
- [Configuration Guide](CONFIGURATION.md)
- [Provider Setup](PROVIDERS.md)
- [Troubleshooting](TROUBLESHOOTING.md)
