# Antigravity Quick Start Guide

**Last Updated:** March 7, 2026  
**Status:** ✅ Production Ready

---

## 🚀 Quick Start

### Step 1: Login with Antigravity

```bash
# Recommended command (works in all cases)
./picoclaw-agents auth login --provider google-antigravity

# Alternative short name (also works)
./picoclaw-agents auth login --provider antigravity
```

> [!IMPORTANT]
> ❌ **DO NOT use:** `./picoclaw-agents auth antigravity` (this command does NOT exist)  
> ✅ **Always use:** `./picoclaw-agents auth login --provider google-antigravity`

---

### Step 2: Complete OAuth Flow

1. **Browser opens automatically** (on local machines)
2. **Sign in with your Google account** (the one with Google One AI Premium or Workspace Gemini)
3. **Grant permissions** to PicoClaw
4. **Browser redirects back** and credentials are saved automatically

**Headless/Remote:** If browser doesn't open, copy the URL from terminal and paste it in your browser manually.

---

### Step 3: Verify Authentication

```bash
# Check auth status
./picoclaw-agents auth status

# View available models
./picoclaw-agents auth models
```

Expected output:
```
✓ Google Antigravity login successful!
Email: your-email@gmail.com
Project ID: your-project-id

Available Antigravity Models:
  - gemini-3.1-pro-high
  - gemini-3-flash
  - gemini-2.5-pro
  - claude-sonnet-4-6
  - claude-opus-4-6-thinking
```

---

### Step 4: Start Using PicoClaw

```bash
# Run the gateway (for chat channels like Telegram, Discord)
./picoclaw-agents gateway

# Or run a one-shot query
./picoclaw-agents agent -m "Hello, how are you?"

# Use a specific Antigravity model
./picoclaw-agents agent -m "Write a poem" --model antigravity/gemini-3-flash
```

---

## 📋 Available Commands

| Command | Description |
|---------|-------------|
| `./picoclaw-agents auth login --provider google-antigravity` | Login with Antigravity |
| `./picoclaw-agents auth status` | Check authentication status |
| `./picoclaw-agents auth models` | List available models |
| `./picoclaw-agents auth logout --provider google-antigravity` | Logout from Antigravity |

---

## 🎯 Available Antigravity Models

| Model | Description | Best For |
|-------|-------------|----------|
| `antigravity/gemini-3.1-pro-high` | Most powerful Gemini | Complex reasoning, coding |
| `antigravity/gemini-3-flash` | Fast & economical | Quick responses, simple tasks |
| `antigravity/gemini-2.5-pro` | Balanced Gemini Pro | General purpose |
| `antigravity/claude-sonnet-4-6` | Claude Sonnet | Writing, analysis |
| `antigravity/claude-opus-4-6-thinking` | Claude Opus with reasoning | Complex problem solving |

---

## ⚙️ Configuration

After successful login, your `~/.picoclaw/config.json` is automatically updated:

```json
{
  "model_list": [
    {
      "model_name": "gemini-flash",
      "model": "antigravity/gemini-3-flash",
      "auth_method": "oauth"
    }
  ],
  "agents": {
    "defaults": {
      "model": "gemini-flash"
    }
  }
}
```

**Credentials are stored in:** `~/.picoclaw/auth/credentials.json`

---

## 🔧 Troubleshooting

### ❌ "required flag(s) 'provider' not set"

**Cause:** Missing `--provider` flag.

**Solution:**
```bash
./picoclaw-agents auth login --provider google-antigravity
```

---

### ❌ "unrecognized command: antigravity"

**Cause:** `./picoclaw-agents auth antigravity` is **NOT a valid command**.

**Solution:**
```bash
./picoclaw-agents auth login --provider google-antigravity
```

---

### ❌ "Token expired"

**Cause:** Access token expired and auto-refresh failed.

**Solution:**
```bash
# Re-authenticate
./picoclaw-agents auth login --provider google-antigravity
```

---

### ❌ "Gemini for Google Cloud is not enabled"

**Cause:** Gemini API not enabled in your Google Cloud project.

**Solution:**
1. Go to [Google Cloud Console](https://console.cloud.google.com)
2. Select your project
3. Enable "Gemini for Google Cloud" API

---

### ❌ Browser doesn't open

**Cause:** Running on headless server or WSL2.

**Solution:** Manually copy the authorization URL from terminal to your browser.

---

## 📚 Requirements

- **Google Account** with one of:
  - Google One AI Premium plan
  - Workspace Gemini add-on
- **Google Cloud Project** with Gemini API enabled
- **PicoClaw** v3.4.4 or later

---

## 🔗 Related Documentation

- [ANTIGRAVITY_AUTH.md](./ANTIGRAVITY_AUTH.md) - Full authentication guide
- [ANTIGRAVITY_USAGE.md](./ANTIGRAVITY_USAGE.md) - Usage examples and best practices
- [Google Cloud Console](https://console.cloud.google.com) - Manage quotas and billing

---

## 💡 Tips

1. **Free with Subscription:** Antigravity uses your Google One/Workspace quota (no separate billing)
2. **Auto-Refresh:** Tokens refresh automatically - no need to re-login frequently
3. **Multiple Models:** Access both Gemini and Claude models through one authentication
4. **Project Quotas:** Monitor usage in Google Cloud Console to avoid rate limits

---

**Ready to start?** Run `./picoclaw-agents auth login --provider google-antigravity` now! 🚀
