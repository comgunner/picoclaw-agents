# Free Tier Providers — Zero-Cost AI with PicoClaw-Agents

**Last Updated:** April 5, 2026  
**Status:** ✅ All providers tested and operational — no credit card required

---

## 🆓 Overview

PicoClaw-Agents supports multiple **100% free** LLM providers. No credit card, no paid plan, no trial expiration. Just sign up and start using powerful AI models.

### Quick Comparison

| Provider | Free Model | Cost | Signup | Best For |
|----------|-----------|------|--------|----------|
| **OpenRouter** | `openrouter/auto` | Free forever | [openrouter.ai](https://openrouter.ai/) | Auto-routing, simplicity |
| **Zhipu AI** | `glm-4.5-flash` | Free forever | [z.ai](https://z.ai/) | Speed, coding tasks |
| **OpenAI** | `gpt-4.1-mini` | Free tier | [chatgpt.com](https://chatgpt.com/#settings/Security) | Quality, reasoning |
| **Qwen** | `qwen-plus` | Free tier | [dashscope.aliyun.com](https://dashscope.aliyun.com/) | Best overall quality ⭐ |

---

## 1. 🌐 OpenRouter Free (`openrouter-free`)

### Why Use It
- **Zero configuration**: Auto-routes to the best free model available
- **No model selection needed**: OpenRouter picks the optimal one automatically
- **Fallback built-in**: If one model goes down, switches to another

### Setup

```bash
# Interactive login
./picoclaw-agents auth login --provider openrouter-free

# Or via onboard wizard
./picoclaw-agents onboard --free
```

### Configuration
```json
{
  "model_list": [
    {
      "model_name": "openrouter-free",
      "model": "openrouter/auto",
      "api_key": "sk-or-v1-..."
    }
  ]
}
```

### Available Free Models
OpenRouter auto-routes to the best available, which may include:
- `stepfun/step-3.5-flash` (256K context)
- `deepseek/deepseek-v3.2` (64K context)
- Various Llama, Gemma, and MiniMax variants

### Links
- **Free Models:** https://openrouter.ai/collections/free-models
- **API Keys:** https://openrouter.ai/keys

---

## 2. 🧠 Zhipu AI (`zhipu`) — 100% Free Forever

### Why Use It
- **Completely free**: No credit card, no paid plan required
- **Fast inference**: `glm-4.5-flash` is optimized for speed
- **Great for coding**: Strong code generation and understanding
- **No usage limits**: Generous free tier for personal use

### Setup

```bash
# Interactive login
./picoclaw-agents auth login --provider zhipu
```

The wizard will:
1. Prompt for your Zhipu API key
2. Auto-configure `glm-4.5-flash` as the default model
3. Add all available GLM models to your config

### Available Models

| Model | Context | Use Case |
|-------|---------|----------|
| `glm-4.5-flash` 🆓 | 128K | Default — fast, free, capable |
| `glm-4.7-flash` | 128K | Latest flash model |
| `glm-5` | 128K | Premium model (may require credits) |
| `glm-5-turbo` | 128K | Optimized for speed |
| `glm-4.5-air` | 128K | Lightweight variant |

### Links
- **Signup:** https://z.ai/
- **API Keys:** https://open.bigmodel.cn/usercenter/proj-mgmt/apikeys
- **Pricing:** https://z.ai/pricing

---

## 3. 🔑 OpenAI (`openai`) — Free Tier via Device Code

### Why Use It
- **High quality**: GPT models are industry-leading
- **Free tier available**: Works with OpenAI's free ChatGPT plan
- **No API key needed**: Uses OAuth Device Code flow

### Setup

```bash
# Device Code OAuth (Browser method NOT supported)
./picoclaw-agents auth login --provider openai --device-code
```

**Important:** OpenAI only supports **Device Code** authentication. Browser OAuth is NOT available.

### How It Works
1. Run the command above
2. Open the URL shown in terminal on any device
3. Enter the 8-character code
4. Authorize with your OpenAI/ChatGPT account
5. Models auto-added to your config

### Available Models

| Model | Free? | Context |
|-------|-------|---------|
| `gpt-4.1-mini` | ✅ Free tier | 32K |
| `gpt-4.1` | ✅ Free tier | 32K |
| `gpt-5` | ✅ Free tier | 32K |
| `o3-mini` | ✅ Free tier | 32K |
| `o3` | ✅ Free tier | 32K |
| `o1` | ✅ Free tier | 32K |

### Requirements
- Enable Device Code authorization at [chatgpt.com/#settings/Security](https://chatgpt.com/#settings/Security)
- Free OpenAI/ChatGPT account

---

## 4. ⭐ Qwen (`qwen`) — Best Overall Free Tier

### Why Use It
- **Highest quality**: Comparable to paid models in benchmarks
- **Generous free tier**: High rate limits for personal use
- **Multilingual**: Excellent in English, Chinese, and many other languages
- **Long context**: Up to 128K tokens on some models

### ⚠️ Important Note
> **How long will Qwen's free tier last?** We don't know. As of April 2026, it's the **best free option available** — highly recommended to use it while it lasts. No official end date has been announced.

### Setup

```bash
# OAuth via browser (recommended)
./picoclaw-agents auth login --provider qwen

# Or paste API key directly
./picoclaw-agents auth login --provider qwen --token
```

### Available Models

| Model | Free? | Context | Use Case |
|-------|-------|---------|----------|
| `qwen-plus` 🆓 | ✅ Free tier | 128K | **Default** — best balance |
| `qwen-max` | ✅ Free tier | 32K | Highest quality |
| `qwen-turbo` | ✅ Free tier | 128K | Fastest inference |
| `qwen-long` | ✅ Free tier | 1M | Ultra-long context |
| `qwen-vl-max` | ✅ Free tier | 32K | Vision/language |
| `qwen-vl-plus` | ✅ Free tier | 32K | Vision (lighter) |

### Regional Endpoints

| Region | URL |
|--------|-----|
| **US (Virginia)** | `https://dashscope-us.aliyuncs.com/compatible-mode/v1` |
| **Singapore** | `https://dashscope-sg.aliyuncs.com/compatible-mode/v1` |
| **China (Beijing)** | `https://dashscope.aliyuncs.com/compatible-mode/v1` |

### Links
- **Signup:** https://dashscope.aliyun.com/
- **API Keys:** https://dashscope.console.aliyun.com/api-key
- **Models:** https://help.aliyun.com/zh/model-studio/getting-started/models

---

## 📊 Recommendation Matrix

| Need | Recommended Provider |
|------|---------------------|
| **Zero setup, just works** | `openrouter-free` |
| **Best quality (use it while it lasts)** | `qwen` ⭐ |
| **Coding tasks, fast inference** | `zhipu` |
| **OpenAI ecosystem, reasoning** | `openai --device-code` |
| **Ultra-long documents** | `qwen` (qwen-long: 1M context) |
| **Vision/image understanding** | `qwen` (qwen-vl-max) |

---

## 🚀 Multi-Provider Setup (Recommended)

For maximum reliability, configure multiple providers as fallbacks:

```bash
# 1. Set up your primary free provider
./picoclaw-agents auth login --provider qwen

# 2. Add OpenRouter as backup
./picoclaw-agents auth login --provider openrouter-free

# 3. Add Zhipu as another backup
./picoclaw-agents auth login --provider zhipu
```

Then switch between them in the WebUI:
```
http://localhost:18800/credentials
```

Or via CLI:
```bash
# List available models
./picoclaw-agents models list

# Switch model/provider
./picoclaw-agents agent --model qwen-plus -m "Hello"
./picoclaw-agents agent --model glm-4.5-flash -m "Hello"
./picoclaw-agents agent --model openrouter/auto -m "Hello"
```

---

## ⚠️ Known Limitations

| Provider | Limitation |
|----------|-----------|
| **OpenRouter Free** | Model may change without notice; rate limits vary |
| **Zhipu** | `glm-5` may require paid credits; stick to `glm-4.5-flash` for free |
| **OpenAI** | Free tier has usage caps; Device Code only (no Browser OAuth) |
| **Qwen** | Free tier duration unknown; may change policy at any time |

---

## 🔗 Related Documentation

- **OpenRouter Guide:** [OPENROUTER_FREE.md](OPENROUTER_FREE.md) / [OPENROUTER_FREE.es.md](OPENROUTER_FREE.es.md)
- **Changelog:** [CHANGELOG.md](../CHANGELOG.md)
- **Token Overflow Fix:** [local_work/openrouter_free_token_fix.md](../local_work/openrouter_free_token_fix.md)
- **Config Reference:** [local_work/CONFIG_FIELD_REFERENCE.md](../local_work/CONFIG_FIELD_REFERENCE.md)

---

*FREE_TIER_PROVIDERS.md — Updated April 5, 2026*
