# Antigravity Provider Guide

**Last Updated:** March 30, 2026  
**Status:** ✅ Production Ready (v1.3.0-alpha)  
**New Feature:** 🎉 Auto-Config on Login - All 15 models added automatically!

---

## Overview

**Antigravity** (Google Cloud Code Assist) is a Google-backed AI provider that offers access to Gemini and Claude models through Google's Cloud infrastructure using **OAuth 2.0 authentication**.

**Key Distinction:** Antigravity uses your **Google One AI Premium** or **Workspace Gemini** plan quotas — NOT a pay-per-use API key.

---

## Quick Start (NEW!)

### One-Command Setup

```bash
# Login and auto-configure all 15 Antigravity models
./picoclaw-agents auth login --provider google-antigravity

# Test with default model (gemini-3-flash)
./picoclaw-agents agent -m "Hello world"
```

**What happens automatically:**
1. ✅ OAuth authentication via browser
2. ✅ **All 15 Antigravity models added to config**
3. ✅ `gemini-3-flash` set as default model
4. ✅ Fallback to `gemini-2.5-flash` configured

**Output:**
```
✓ Google Antigravity login successful!

✓ Added 15 Antigravity models to config

Default model set to: gemini-3-flash (fallback: gemini-2.5-flash)

Available models:
  - gemini-3-flash (default)
  - gemini-3-pro-high, gemini-3-pro-low
  - gemini-3.1-pro-high, gemini-3.1-pro-low, gemini-3.1-flash-lite
  - gemini-3-flash-agent, gemini-3-flash-preview
  - gemini-2.5-flash, gemini-2.5-flash-lite, gemini-2.5-flash-thinking, gemini-2.5-pro
  - claude-sonnet-4-6, claude-opus-4-6-thinking
  - gpt-oss-120b-medium

Try it: picoclaw-agents agent -m "Hello world" --model gemini-3-flash
```

---

## Authentication

### Step 1: Login (Auto-Config)

```bash
./picoclaw-agents auth login --provider google-antigravity
```

**What's NEW (v1.3.0-alpha):**
- 🎉 **Automatically adds ALL 15 Antigravity models** to `~/.picoclaw/config.json`
- 🎉 **Sets `gemini-3-flash` as default model** for all agents
- 🎉 **Configures fallback** to `gemini-2.5-flash`
- 🎉 **No manual config editing required!**

**Alias also works:**
```bash
./picoclaw-agents auth login --provider antigravity
```

### Step 2: Complete OAuth Flow

1. **Browser opens automatically** (local machines)
2. **Sign in** with your Google account (must have Google One AI Premium or Workspace Gemini)
3. **Grant permissions** to PicoClaw
4. **Credentials saved** to `~/.picoclaw/auth.json`
5. **Config updated** with all 15 models ✨

**Headless/Remote (VPS/Coolify/Docker):**
1. Run the command
2. Copy the URL and open in your local browser
3. Complete login
4. Copy the final `localhost:51121` redirect URL from your browser
5. Paste it back into the terminal

### Token Management

| Token Type | Duration | Auto-Refresh |
|------------|----------|--------------|
| `access_token` | ~1 hour | ✅ Yes |
| `refresh_token` | Months/indefinite | N/A |

**Auto-refresh layers:**
1. **Background daemon**: Proactively refreshes every 20 min if <30 min remain
2. **On every request**: Retries with refresh_token even if already expired
3. **`auth models` command**: Also recovers from expired tokens

**Manual re-auth needed only if:**
- You revoked access from `myaccount.google.com > Security > Apps with access`
- You changed your Google password
- The `refresh_token` has been inactive for 6+ months

### Check Status

```bash
./picoclaw-agents auth status
```

### Logout

```bash
./picoclaw-agents auth logout --provider google-antigravity
```

---

## Available Models (OAuth Auth)

### View All Models

```bash
./picoclaw-agents auth models
```

### Complete Model List (15 Models)

**Auto-configured on login (v1.3.0-alpha+):**

| # | Model Name | Description | Best For |
|---|------------|-------------|----------|
| 1 | `gemini-3-flash` ⭐ | **DEFAULT** - Fast, reliable | **Recommended default** |
| 2 | `gemini-3-pro-high` | High reasoning Gemini 3 | Complex reasoning |
| 3 | `gemini-3-pro-low` | Low reasoning Gemini 3 | Simple tasks |
| 4 | `gemini-3.1-pro-high` | High reasoning Gemini 3.1 | Advanced tasks |
| 5 | `gemini-3.1-pro-low` | Low reasoning Gemini 3.1 | Medium tasks |
| 6 | `gemini-3.1-flash-lite` | Lightweight 3.1 model | Fast responses |
| 7 | `gemini-3-flash-agent` | Agent-optimized Flash | Multi-step tasks |
| 8 | `gemini-3-flash-preview` | Preview model | Testing new features |
| 9 | `gemini-2.5-flash` | Gemini 2.5 Flash | Fast responses |
| 10 | `gemini-2.5-flash-lite` | Lightweight 2.5 | Simple tasks |
| 11 | `gemini-2.5-flash-thinking` | Flash with reasoning | Reasoning tasks |
| 12 | `gemini-2.5-pro` | Gemini 2.5 Pro | General purpose |
| 13 | `claude-sonnet-4-6` | Claude Sonnet 4.6 | Writing, analysis |
| 14 | `claude-opus-4-6-thinking` | Claude Opus with thinking | Complex problem solving |
| 15 | `gpt-oss-120b-medium` | Open-source GPT alternative | General use |

> [!NOTE]
> **Model Names Updated (v1.3.0-alpha)**
>
> Model names now match exactly what `auth models` returns:
> - ✅ `gemini-3-flash` (simple, matches API)
> - ❌ `antigravity-gemini-3-flash` (old format, still works but deprecated)

---

## Usage Examples

### Command Line

```bash
# Use default model (gemini-3-flash)
./picoclaw-agents agent -m "Hello"

# Use specific model
./picoclaw-agents agent -m "Hello" --model gemini-3-flash

# Claude for writing
./picoclaw-agents agent -m "Write a poem" --model claude-sonnet-4-6

# Complex reasoning
./picoclaw-agents agent -m "Solve this math problem" --model claude-opus-4-6-thinking

# Fast responses
./picoclaw-agents agent -m "Quick question" --model gemini-3.1-flash-lite
```

### Web UI (NEW!)

**Model Selector Feature (v1.3.0-alpha):**

1. Open Web UI: http://localhost:18800/
2. Click model dropdown in header
3. Select any of the 15 Antigravity models
4. Send message - **model change applies immediately!**

> [!TIP]
> **Web UI Model Override**
>
> The model selected in Web UI now **actually works**! Each user can use different models independently.
>
> See [`PICO_MODEL_OVERRIDE.md`](./PICO_MODEL_OVERRIDE.md) for details.

---

## Configuration

### Automatic Config (v1.3.0-alpha+)

**After login, your config automatically has:**

```json
{
  "agents": {
    "defaults": {
      "model": "gemini-3-flash",
      "fallbacks": ["gemini-2.5-flash"]
    }
  },
  "model_list": [
    {
      "model_name": "gemini-3-flash",
      "model": "antigravity/gemini-3-flash",
      "auth_method": "oauth"
    },
    {
      "model_name": "gemini-3-pro-high",
      "model": "antigravity/gemini-3-pro-high",
      "auth_method": "oauth"
    },
    // ... 13 more models ...
  ]
}
```

### Manual Config (Pre-v1.3.0 or Custom Setup)

If you need to manually add models:

```bash
# Option 1: Re-run login to auto-add all models
./picoclaw-agents auth login --provider google-antigravity

# Option 2: Use sync script
./scripts/sync_antigravity_models.sh

# Option 3: Use fix script (updates model names)
./scripts/fix_antigravity_models.sh
```

### model_list Format

**Correct format (v1.3.0-alpha+):**
```json
{
  "model_name": "gemini-3-flash",
  "model": "antigravity/gemini-3-flash",
  "auth_method": "oauth"
}
```

**Old format (deprecated but still works):**
```json
{
  "model_name": "antigravity-gemini-3-flash",
  "model": "antigravity/gemini-3-flash",
  "auth_method": "oauth"
}
```

> [!IMPORTANT]
> **Use Simple Model Names**
>
> Always use the simple name (e.g., `gemini-3-flash`) in:
> - CLI: `--model gemini-3-flash`
> - Web UI: Select from dropdown
> - Config: `model_name` field
>
> The `antigravity/` prefix is only for the `model` field internally.

---

## Image Generation (API Key Only)

**Antigravity OAuth does NOT support image generation.** For generating images, you must use **Google AI Studio API Key**.

### Supported Image Models (API Key)

| Model | Provider Prefix | Purpose |
|-------|----------------|---------|
| `gemini-2.5-flash-image` | `gemini/` | Nano Banana - image generation |
| `gemini-3-pro-image-preview` | `gemini/` | Nano Banana Pro |
| `gemini-3.1-flash-image-preview` | `gemini/` | Nano Banana 2 |
| `imagen-4.0-generate-001` | `gemini/` | Imagen 4 |
| `imagen-4.0-ultra-generate-001` | `gemini/` | Imagen 4 Ultra |

### Configuration

Add to `~/.picoclaw/config.json`:

```json
{
  "model_list": [
    {
      "model_name": "gemini-2.5-flash-image",
      "model": "gemini/gemini-2.5-flash-image",
      "api_key": "YOUR_GEMINI_API_KEY"  # pragma: allowlist secret
    }
  ],
  "tools": {
    "image_gen": {
      "provider": "gemini",
      "gemini_api_key": "YOUR_GEMINI_API_KEY",  # pragma: allowlist secret
      "gemini_image_model_name": "gemini-2.5-flash-image",
      "output_dir": "~/.picoclaw/workspace/generated_images"
    }
  }
}
```

**Get API Key:** [Google AI Studio](https://aistudio.google.com/app/apikey)

---

## Configuration

### Default Configuration (DeepSeek)

The main `config.example.json` uses **deepseek-chat** as default:

```bash
cp config/config.example.json ~/.picoclaw/config.json
# Add your DeepSeek API key
./picoclaw-agents agent -m "Hello"
```

### Antigravity Configuration

The `config/config.example_antigravity.json` uses `antigravity-gemini-3-flash` for all agents:

```bash
cp config/config.example_antigravity.json ~/.picoclaw/config.json
./picoclaw-agents auth login --provider google-antigravity
./picoclaw-agents agent -m "Hello"
```

### model_list Entries

**Antigravity (OAuth):**
```json
{
  "model_name": "antigravity-gemini-3-flash",
  "model": "antigravity/gemini-3-flash",
  "api_key": "",
  "auth_method": "oauth"
}
```

**Google AI Studio (API Key):**
```json
{
  "model_name": "gemini-2.5-flash",
  "model": "gemini/gemini-2.5-flash",
  "api_key": "YOUR_GEMINI_API_KEY"  # pragma: allowlist secret
}
```

### Comparative Examples

#### 1. Gemini 2.5 Flash

```json
// Antigravity OAuth (uses Google Cloud quota)
{
  "model_name": "ag-gemini-2.5-flash",
  "model": "antigravity/gemini-2.5-flash",
  "api_key": "",
  "auth_method": "oauth"
}

// Google AI Studio API Key (pay-per-use or free tier)
{
  "model_name": "gemini-2.5-flash",
  "model": "gemini/gemini-2.5-flash",
  "api_key": "YOUR_GEMINI_API_KEY"  # pragma: allowlist secret
}
```

#### 2. Gemini 3 Flash

```json
// Antigravity OAuth
{
  "model_name": "ag-gemini-3-flash",
  "model": "antigravity/gemini-3-flash",
  "api_key": "",
  "auth_method": "oauth"
}

// Google AI Studio API Key
{
  "model_name": "gemini-3-flash-preview",
  "model": "gemini/gemini-3-flash-preview",
  "api_key": "YOUR_GEMINI_API_KEY"  # pragma: allowlist secret
}
```

---

## Model Routing Architecture

PicoClaw uses a 3-step pipeline:

### Configuration Fields
- **`model_name`**: Internal alias — the friendly name you use (e.g., `antigravity-gemini-3-flash`)
- **`model`**: Routing instruction — must contain `provider/model-id` (e.g., `antigravity/gemini-3-flash`)

### The Pipeline

1. **Memory Load**: On startup, reads `model_list` from `~/.picoclaw/config.json` into RAM. Changes require restart.

2. **Factory Routing**: The alias is looked up → the `model` field is split by `/` → the `antigravity` prefix selects the Antigravity provider.

3. **Prefix Sanitization**: Before calling the API, the provider strips all prefixes:
   - `antigravity/gemini-3-flash` → `gemini-3-flash` ✅
   - `antigravity-gemini-3-flash` → `gemini-3-flash` ✅ (dash prefix also handled)

> [!TIP]
> Both `antigravity/gemini-3-flash` (slash) and `antigravity-gemini-3-flash` (dash) are valid.

---

## Real-world Usage (Coolify/Docker)

### Option 1: Copy Credentials

```bash
# Authenticate locally first
./picoclaw-agents auth login --provider google-antigravity

# Copy credentials to server
scp ~/.picoclaw/auth.json user@your-server:~/.picoclaw/
```

### Option 2: Authenticate on Server

```bash
# Run on server (headless flow)
./picoclaw-agents auth login --provider google-antigravity
# Copy URL, open locally, paste redirect URL back
```

---

## Troubleshooting

| Error | Cause | Solution |
|-------|-------|----------|
| `403 PERMISSION_DENIED` | Token expired/revoked | `./picoclaw-agents auth login --provider google-antigravity` |
| `ACCESS_TOKEN_SCOPE_INSUFFICIENT` | Token expired/revoked | `./picoclaw-agents auth login --provider google-antigravity` |
| `404 NOT_FOUND` | Model alias not resolved | Verify `model` field has `antigravity/` prefix and `auth_method: "oauth"` |
| `401 invalid_api_key` | Wrong provider used | Check `model` field has `antigravity/` prefix, not OpenAI-style key |
| `429 Rate Limit` | Quota exhausted | Wait for reset time shown by PicoClaw, or switch model |
| Empty response | Model restricted for project | Try `antigravity-gemini-3-flash` or `gemini-2.5-flash` |
| "Gemini for Google Cloud not enabled" | API not enabled | Enable in [Google Cloud Console](https://console.cloud.google.com) |

---

## Requirements

- **Google Account** with:
  - Google One AI Premium plan, OR
  - Workspace Gemini add-on
- **Google Cloud Project** with Gemini API enabled
- **PicoClaw** v3.4.4 or later

---

## Commands Reference

| Command | Description |
|---------|-------------|
| `./picoclaw-agents auth login --provider google-antigravity` | Login with Antigravity |
| `./picoclaw-agents auth status` | Check authentication status |
| `./picoclaw-agents auth models` | List available models |
| `./picoclaw-agents auth logout --provider google-antigravity` | Logout from Antigravity |
| `./picoclaw-agents agent -m "msg" --model <model>` | Use specific model |

---

## Related Documentation

- [ANTIGRAVITY.es.md](./ANTIGRAVITY.es.md) - Spanish version
- [IMAGE_GEN_util.md](./IMAGE_GEN_util.md) - Image generation workflows
- [Google Cloud Console](https://console.cloud.google.com) - Manage quotas and billing
- [Google AI Studio](https://aistudio.google.com) - Get API keys for image generation

---

**Quick Start:** Run `./picoclaw-agents auth login --provider google-antigravity` now! 🚀
