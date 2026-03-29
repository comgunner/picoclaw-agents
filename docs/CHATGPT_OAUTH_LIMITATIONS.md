# ChatGPT OAuth - Feature Removed

**Date:** March 28, 2026  
**Status:** ❌ **REMOVED - Not Available**

---

## 🚨 Feature Removal Notice

**The `--provider chatgpt` OAuth feature has been removed from PicoClaw as of v3.4.5.**

This feature was experimental and had fundamental limitations that prevented it from working correctly with the OpenAI API.

---

## ⚠️ Why Was It Removed?

### 1. **Token Incompatibility**

OAuth tokens obtained from ChatGPT web (`chat.openai.com`) **do not work** with the OpenAI REST API (`api.openai.com/v1`).

**Error when attempting to use ChatGPT OAuth with API:**
```json
{
  "error": {
    "message": "invalid model ID",
    "type": "invalid_request_error"
  }
}
```

### 2. **Separate Systems**

| System | Token Type | Endpoint | Purpose |
|--------|-----------|----------|---------|
| **ChatGPT Web** | OAuth token | `chat.openai.com` | Web interface only |
| **OpenAI API** | API key | `api.openai.com/v1` | REST API access |

These are **completely separate authentication systems** that cannot be interchanged.

### 3. **OAuth Flow Issues**

- Browser OAuth flow returned `unknown_error` from OpenAI
- Required CAPTCHA and additional verification
- Not suitable for automated CLI authentication

---

## ✅ Recommended Alternatives

### Option 1: OpenAI Codex OAuth (Recommended)

Use the official OpenAI Codex OAuth flow, which works with the Codex backend:

```bash
# Login with OpenAI OAuth
picoclaw-agents auth login --provider openai

# This connects to the Codex backend (chatgpt.com/backend-api/codex)
# Optimized for coding tasks
picoclaw-agents agent -m "Hello"
```

**Configuration:**
```json
{
  "model_list": [{
    "model_name": "gpt-5.2",
    "model": "openai/gpt-5.2",
    "auth_method": "oauth"
  }]
}
```

### Option 2: OpenAI API Key

Use a standard OpenAI API key:

1. Get an API key from https://platform.openai.com/api-keys
2. Configure in `~/.picoclaw/config.json`:

```json
{
  "model_list": [{
    "model_name": "gpt-4o-mini",
    "model": "openai/gpt-4o-mini",
    "api_key": "sk-proj-..."  // pragma: allowlist secret
  }]
}
```

### Option 3: OpenRouter Free Tier

Use OpenRouter's free tier for testing:

```bash
picoclaw-agents onboard --free
```

---

## 📚 Supported OAuth Providers

The following OAuth providers are currently supported:

| Provider | Command | Status |
|----------|---------|--------|
| **OpenAI Codex** | `auth login --provider openai` | ✅ Active |
| **Anthropic** | `auth login --provider anthropic` | ✅ Active (token paste) |
| **Google Antigravity** | `auth login --provider google-antigravity` | ✅ Active |

**ChatGPT web OAuth is NOT available.**

---

## 🔍 Technical Details

### Why ChatGPT OAuth Cannot Work

1. **No Public OAuth API**: OpenAI does not expose a public OAuth API for ChatGPT web
2. **Token Scope Limitation**: ChatGPT OAuth tokens are scoped only for web interface access
3. **Different Backend Systems**: ChatGPT web and OpenAI API use completely separate authentication infrastructures

### Previous Implementation Issues

The previous implementation attempted to:
- Use the same OAuth credentials as Codex but with a different `originator` parameter
- Connect to `api.openai.com/v1` using ChatGPT OAuth tokens
- Store credentials under a separate `chatgpt` provider key

**All of these approaches failed** because the fundamental architecture doesn't support it.

---

## 📝 Migration Guide

If you previously configured ChatGPT OAuth:

### 1. Remove old configuration

```bash
# Logout from chatgpt provider (if still configured)
picoclaw-agents auth logout --provider chatgpt
```

### 2. Clean up config file

Edit `~/.picoclaw/config.json` and remove any entries with:
- `"model": "chatgpt/..."`
- `"auth_method": "oauth"` for chatgpt models

### 3. Use OpenAI Codex instead

```bash
# Login with OpenAI OAuth
picoclaw-agents auth login --provider openai

# Update your config to use openai models
# Example: "openai/gpt-4o" instead of "chatgpt/gpt-4o"
```

---

## 📞 Support

If you have questions or need help migrating:

- **GitHub Issues**: https://github.com/comgunner/picoclaw-agents/issues
- **Documentation**: See `docs/` for supported provider guides
- **Community**: Check Discord/Telegram channels

---

**Feature Removed:** March 28, 2026  
**Version:** v3.4.5  
**Author:** @comgunner
