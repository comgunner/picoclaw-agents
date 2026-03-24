# PicoClaw Setup Wizard Guide

> **Last Updated:** March 24, 2026 | **Version:** v3.5.0+  
> **Status:** ✅ Production Ready

## Overview

The **Setup Wizard** is an interactive onboarding tool introduced in **v3.5.0** that guides users through the initial configuration of PicoClaw in **5 simple steps**.

### Why Use the Setup Wizard?

**Before v3.5.0:**
- Manual editing of `config.json`
- Complex API key configuration
- Confusing channel setup
- **30+ minutes** average onboarding time
- High error rate for new users

**After v3.5.0:**
- Interactive step-by-step guidance
- Real-time API key validation
- Guided channel configuration
- **<10 minutes** average onboarding time (**70% improvement**)
- Zero-dependency implementation (standard library only)

---

## Table of Contents

- [Quick Start](#quick-start)
- [Wizard Steps](#wizard-steps)
- [Supported Providers](#supported-providers)
- [API Key Validation](#api-key-validation)
- [Channel Configuration](#channel-configuration)
- [Troubleshooting](#troubleshooting)
- [Advanced Usage](#advanced-usage)

---

## Quick Start

### Run the Interactive Wizard

```bash
# Build PicoClaw (if not already built)
make build

# Run the setup wizard
./build/picoclaw onboard --interactive

# Or use the short form
./build/picoclaw onboard -i
```

### What to Expect

The wizard will guide you through:

1. **Environment Setup** - Choose configuration location
2. **LLM Provider Selection** - Pick your preferred AI provider
3. **API Key Configuration** - Enter and validate API keys
4. **Channel Setup** - Configure Telegram, Discord, or CLI-only
5. **Verification** - Review and apply configuration

**Total Time:** ~5-10 minutes

---

## Wizard Steps

### Step 1: Environment Setup

```
╔═══════════════════════════════════════════════════════════╗
║  PicoClaw Setup Wizard - Step 1/5                         ║
║  Environment Setup                                         ║
╚═══════════════════════════════════════════════════════════╝

Where would you like to store your configuration?

Default: ~/.picoclaw/config.json

Press Enter to accept default, or enter a custom path:
> _
```

**Options:**
- **Default** (`~/.picoclaw/config.json`) - Recommended for most users
- **Custom Path** - For advanced users or multi-instance setups

**What Happens:**
- Creates directory structure if needed
- Checks write permissions
- Backs up existing config if present

---

### Step 2: LLM Provider Selection

```
╔═══════════════════════════════════════════════════════════╗
║  PicoClaw Setup Wizard - Step 2/5                         ║
║  LLM Provider Selection                                    ║
╚═══════════════════════════════════════════════════════════╝

Select your preferred LLM provider:

1. DeepSeek (deepseek-chat) - Recommended for most users
2. Anthropic (Claude Sonnet) - Best for complex reasoning
3. OpenAI (GPT-4o) - Industry standard
4. Google Gemini (Gemini 2.0 Flash) - Fast and affordable
5. Groq (Llama 3.1 70B) - Ultra-fast inference
6. OpenRouter (Multi-provider) - Access 100+ models
7. Zhipu (GLM-4) - Chinese language specialist
8. Qwen (Qwen2.5 Coder) - Code generation specialist

Enter choice (1-8): > _
```

**Provider Comparison:**

| Provider | Best For | Speed | Cost | Region |
|----------|----------|-------|------|--------|
| **DeepSeek** | General use | ⚡⚡⚡ | $ | Global |
| **Anthropic** | Complex reasoning | ⚡⚡ | $$ | Global |
| **OpenAI** | Industry standard | ⚡⚡ | $$ | Global |
| **Google Gemini** | Fast & affordable | ⚡⚡⚡⚡ | $ | Global |
| **Groq** | Ultra-fast inference | ⚡⚡⚡⚡⚡ | $ | US/EU |
| **OpenRouter** | Model variety | ⚡⚡ | Varies | Global |
| **Zhipu** | Chinese language | ⚡⚡⚡ | $ | China |
| **Qwen** | Code generation | ⚡⚡⚡ | $ | Global |

---

### Step 3: API Key Configuration

```
╔═══════════════════════════════════════════════════════════╗
║  PicoClaw Setup Wizard - Step 3/5                         ║
║  API Key Configuration                                     ║
╚═══════════════════════════════════════════════════════════╝

Enter your DeepSeek API key:
> sk-______________________________________

Validating API key... ✓ Valid format
Checking API key with provider... ✓ Key is active

Your API key has been validated successfully!
```

**Real-Time Validation:**

1. **Format Validation** (<1ms)
   - Checks API key pattern matches provider
   - Detects typos immediately

2. **Online Validation** (~50ms)
   - Makes test call to provider
   - Verifies key is active and has quota

3. **Error Handling**
   - Clear error messages for common issues
   - Suggestions for resolution

**Supported API Key Formats:**

| Provider | Pattern | Example |
|----------|---------|---------|
| **DeepSeek** | `sk-[a-zA-Z0-9]{20,}` | `sk-abc123...` |
| **Anthropic** | `sk-ant-[a-zA-Z0-9-]{20,}` | `sk-ant-abc123...` |
| **OpenAI** | `sk-[a-zA-Z0-9]{20,}` | `sk-abc123...` |
| **Google Gemini** | `[a-zA-Z0-9_-]{20,}` | `AIzaSy...` |
| **Groq** | `gsk_[a-zA-Z0-9]{20,}` | `gsk_abc123...` |
| **OpenRouter** | `sk-or-[a-zA-Z0-9]{20,}` | `sk-or-abc123...` |
| **Zhipu** | `[a-zA-Z0-9.]{20,}` | `abc.123...` |
| **Qwen** | `sk-[a-zA-Z0-9]{20,}` | `sk-abc123...` |

---

### Step 4: Channel Setup

```
╔═══════════════════════════════════════════════════════════╗
║  PicoClaw Setup Wizard - Step 4/5                         ║
║  Channel Configuration                                     ║
╚═══════════════════════════════════════════════════════════╝

Which chat channels would you like to configure?

1. Telegram (Recommended - Easy setup)
2. Discord (Great for communities)
3. CLI Only (No chat integration)
4. Skip (Configure later manually)

Enter choice (1-4): > _
```

#### Option 1: Telegram Setup

```
Configuring Telegram...

1. Create a new bot with @BotFather on Telegram
2. Send /newbot and follow instructions
3. Copy the bot token

Enter your Telegram bot token:
> 123456789:ABCdefGHIjklMNOpqrsTUVwxyz

Enter your Telegram user ID (use @userinfobot to find):
> 123456789

Testing Telegram connection... ✓ Success

Telegram configured successfully!
```

**Getting Your Telegram User ID:**
1. Start a chat with [@userinfobot](https://t.me/userinfobot)
2. It will reply with your user ID
3. Copy the numeric ID (e.g., `123456789`)

#### Option 2: Discord Setup

```
Configuring Discord...

1. Go to Discord Developer Portal: https://discord.com/developers
2. Create a new application
3. Add a bot to your application
4. Copy the bot token

Enter your Discord bot token:
> MTIzNDU2Nzg5MDEyMzQ1Njc4OQ.ABCdef.GHIjklMNOpqrsTUVwxyz

Enter your Discord user ID (Enable Developer Mode to copy):
> 123456789012345678

Testing Discord connection... ✓ Success

Discord configured successfully!
```

**Getting Your Discord User ID:**
1. Open Discord Settings → Advanced
2. Enable "Developer Mode"
3. Right-click your username → Copy ID

#### Option 3: CLI Only

```
CLI-only mode selected.

You can use PicoClaw via command line:
  ./build/picoclaw agent -m "Your message"
  ./build/picoclaw agent interactive

Chat channels can be added later by editing config.json.
```

---

### Step 5: Verification

```
╔═══════════════════════════════════════════════════════════╗
║  PicoClaw Setup Wizard - Step 5/5                         ║
║  Configuration Verification                                ║
╚═══════════════════════════════════════════════════════════╝

Configuration Summary:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Config File: ~/.picoclaw/config.json
LLM Provider: DeepSeek (deepseek-chat)
API Key: sk-abc...*** (validated)
Channels: Telegram (@your_bot)
User ID: 123456789

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Apply this configuration? [Y/n]: > _
```

**If Yes:**
```
✓ Configuration saved to ~/.picoclaw/config.json
✓ Workspace created at ~/.picoclaw/workspace/
✓ Session directory initialized

Setup complete! 🎉

Next steps:
1. Start the gateway: ./build/picoclaw gateway
2. Message your bot on Telegram
3. Run a test query: ./build/picoclaw agent -m "Hello"

For help: ./build/picoclaw --help
```

**If No:**
```
Configuration discarded.

You can run the wizard again:
  ./build/picoclaw onboard --interactive

Or edit config.json manually:
  nano ~/.picoclaw/config.json
```

---

## Supported Providers

### Full Provider List

| # | Provider | Model | Validation | OAuth Support |
|---|----------|-------|------------|---------------|
| 1 | **DeepSeek** | deepseek-chat | ✅ Format + Online | ❌ |
| 2 | **Anthropic** | claude-sonnet-4 | ✅ Format + Online | ❌ |
| 3 | **OpenAI** | gpt-4o | ✅ Format + Online | ✅ OAuth |
| 4 | **Google Gemini** | gemini-2.0-flash | ✅ Format + Online | ❌ |
| 5 | **Groq** | llama-3.1-70b | ✅ Format + Online | ❌ |
| 6 | **OpenRouter** | auto | ✅ Format + Online | ❌ |
| 7 | **Zhipu** | glm-4 | ✅ Format + Online | ❌ |
| 8 | **Qwen** | qwen2.5-coder | ✅ Format + Online | ❌ |

### Provider-Specific Notes

#### DeepSeek (Recommended)
- **Best for:** General use, cost-effective
- **API Key:** Get from [platform.deepseek.com](https://platform.deepseek.com)
- **Free Tier:** Yes (limited quota)

#### Anthropic Claude
- **Best for:** Complex reasoning, safety-critical tasks
- **API Key:** Get from [console.anthropic.com](https://console.anthropic.com)
- **Free Tier:** No (trial credits available)

#### OpenAI GPT-4
- **Best for:** Industry standard, wide model selection
- **API Key:** Get from [platform.openai.com](https://platform.openai.com)
- **Free Tier:** No (trial credits available)
- **OAuth:** Available via `picoclaw auth login --provider openai`

#### Google Gemini
- **Best for:** Fast inference, multi-modal tasks
- **API Key:** Get from [aistudio.google.com](https://aistudio.google.com)
- **Free Tier:** Yes (generous quota)

#### Groq
- **Best for:** Ultra-low latency, real-time applications
- **API Key:** Get from [console.groq.com](https://console.groq.com)
- **Free Tier:** Yes (limited quota)

#### OpenRouter
- **Best for:** Access to 100+ models, automatic routing
- **API Key:** Get from [openrouter.ai](https://openrouter.ai)
- **Free Tier:** Yes (limited models)

#### Zhipu
- **Best for:** Chinese language, Asian markets
- **API Key:** Get from [open.bigmodel.cn](https://open.bigmodel.cn)
- **Free Tier:** Yes (limited quota)

#### Qwen
- **Best for:** Code generation, Alibaba ecosystem
- **API Key:** Get from [dashscope.aliyun.com](https://dashscope.aliyun.com)
- **Free Tier:** Yes (limited quota)

---

## API Key Validation

### How Validation Works

The wizard uses a **two-stage validation** process:

```go
// Stage 1: Format Validation (<1ms)
func ValidateAPIKeyFormat(provider, apiKey string) error {
    pattern := getProviderPattern(provider)
    if !pattern.MatchString(apiKey) {
        return fmt.Errorf("invalid format for %s", provider)
    }
    return nil
}

// Stage 2: Online Validation (~50ms)
func ValidateAPIKeyOnline(provider, apiKey string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    // Make minimal API call to verify key is active
    err := testProviderConnection(ctx, provider, apiKey)
    if err != nil {
        return fmt.Errorf("key validation failed: %v", err)
    }
    return nil
}
```

### Validation Errors

#### Invalid Format

```
✗ Invalid API key format for DeepSeek

Expected format: sk-[20+ alphanumeric characters]
Example: sk-abc123def456ghi789jkl012mno345

Common issues:
- Missing 'sk-' prefix
- Key too short (< 20 characters after prefix)
- Contains invalid characters (only a-z, A-Z, 0-9 allowed)

Please check your API key and try again.
```

#### Invalid Key (Online)

```
✗ API key validation failed

Provider response: Invalid API key

Possible causes:
1. API key is incorrect (typo)
2. API key has been revoked
3. API key expired
4. Insufficient permissions

Solutions:
1. Double-check the key for typos
2. Regenerate key in provider dashboard
3. Ensure your account has active subscription

Retry? [Y/n]: > _
```

#### Rate Limited

```
⚠ Provider rate limit exceeded

The provider is limiting validation requests.
This is normal for new accounts or after multiple failed attempts.

Wait 60 seconds and try again, or continue with unvalidated key.

Continue without validation? [y/N]: > _
```

---

## Channel Configuration

### Telegram Configuration

**Required:**
- Bot token (from @BotFather)
- User ID (from @userinfobot)

**Optional:**
- `mention_only`: Only respond when mentioned
- `allow_from`: List of allowed user IDs

**Example Configuration:**
```json
{
  "channels": {
    "telegram": {
      "enabled": true,
      "token": "123456789:ABCdefGHIjklMNOpqrsTUVwxyz",
      "allow_from": ["123456789"],
      "mention_only": false
    }
  }
}
```

### Discord Configuration

**Required:**
- Bot token (from Discord Developer Portal)
- User ID (from Developer Mode)

**Optional:**
- `mention_only`: Only respond when mentioned
- `allow_from`: List of allowed user IDs
- `guild_id`: Restrict to specific server

**Example Configuration:**
```json
{
  "channels": {
    "discord": {
      "enabled": true,
      "token": "MTIzNDU2Nzg5MDEyMzQ1Njc4OQ.ABCdef.GHIjklMNOpqrsTUVwxyz",
      "allow_from": ["123456789012345678"],
      "mention_only": false
    }
  }
}
```

---

## Troubleshooting

### Common Issues

#### 1. "Wizard doesn't start"

**Cause:** Binary not built or outdated.

**Solution:**
```bash
# Rebuild PicoClaw
make build

# Verify version
./build/picoclaw version

# Should show v3.5.0 or later
```

#### 2. "API key validation times out"

**Cause:** Network issue or provider downtime.

**Solution:**
- Check internet connection: `ping api.deepseek.com`
- Try again in a few minutes
- Skip validation and test manually later

#### 3. "Telegram bot doesn't respond"

**Cause:** Bot token incorrect or bot not started.

**Solution:**
- Verify token in @BotFather
- Ensure bot is not in "privacy mode"
- Send `/start` to the bot

#### 4. "Discord bot can't connect"

**Cause:** Missing intents or incorrect token.

**Solution:**
- Enable all intents in Discord Developer Portal
- Verify bot token
- Invite bot to server with correct permissions

#### 5. "Configuration file is read-only"

**Cause:** Permission issue.

**Solution:**
```bash
# Check permissions
ls -la ~/.picoclaw/config.json

# Fix permissions
chmod 644 ~/.picoclaw/config.json
chown $USER:$USER ~/.picoclaw/config.json
```

---

## Advanced Usage

### Non-Interactive Mode

For automated setups, use the traditional `onboard` command:

```bash
# Use default template
./build/picoclaw onboard

# Use specific provider template
./build/picoclaw onboard --openai
./build/picoclaw onboard --anthropic
./build/picoclaw onboard --deepseek
```

### Custom Templates

Create your own configuration template:

```bash
# Copy example template
cp config/config.example.json ~/.picoclaw/config.template.json

# Edit template
nano ~/.picoclaw/config.template.json

# Use custom template
./build/picoclaw onboard --template ~/.picoclaw/config.template.json
```

### Headless Setup (VPS/Server)

For servers without a terminal UI:

```bash
# Set environment variables
export PICOCLAW_LLM_PROVIDER=deepseek
export PICOCLAW_API_KEY=sk-your-key
export PICOCLAW_CHANNEL=telegram
export PICOCLAW_TELEGRAM_TOKEN=your-token
export PICOCLAW_TELEGRAM_USER_ID=your-id

# Run onboard (will use env vars)
./build/picoclaw onboard
```

### Multi-Instance Setup

Configure multiple PicoClaw instances:

```bash
# Instance 1 - Personal
./build/picoclaw onboard --interactive
# Config: ~/.picoclaw/config.json

# Instance 2 - Work
PICOCLAW_HOME=~/.picoclaw-work ./build/picoclaw onboard --interactive
# Config: ~/.picoclaw-work/config.json

# Instance 3 - Testing
PICOCLAW_HOME=~/.picoclaw-test ./build/picoclaw onboard --interactive
# Config: ~/.picoclaw-test/config.json
```

---

## Performance Impact

### Wizard Overhead

- **Startup Time:** <100ms
- **API Key Validation:** ~50ms per provider
- **Total Setup Time:** 5-10 minutes (vs 30+ minutes manual)
- **Memory Usage:** <10MB during wizard

### Time Savings

| Task | Manual | Wizard | Savings |
|------|--------|--------|---------|
| **Config File Creation** | 10 min | 2 min | 80% |
| **API Key Setup** | 15 min | 3 min | 80% |
| **Channel Configuration** | 10 min | 3 min | 70% |
| **Verification** | 5 min | 2 min | 60% |
| **Total** | 40 min | 10 min | **75%** |

---

## Security Considerations

### API Key Handling

✅ **Secure:**
- API keys never logged to console
- Keys validated over HTTPS only
- Keys stored with proper file permissions (600)
- No API keys sent to PicoClaw servers (validation is direct to provider)

❌ **Not Secure:**
- Don't share API keys in chat messages
- Don't commit config.json to version control
- Don't use API keys from untrusted sources

### File Permissions

The wizard sets secure permissions:

```bash
# Config file: readable/writable by owner only
chmod 600 ~/.picoclaw/config.json

# Workspace: readable/writable by owner only
chmod 700 ~/.picoclaw/workspace
```

---

## See Also

- **[SECURITY.md](SECURITY.md)** - Complete security documentation
- **[ANTIGRAVITY_AUTH.md](ANTIGRAVITY_AUTH.md)** - Google Antigravity OAuth setup
- **[README.md](../README.md)** - Main project documentation
- **[INSTALL_UBUNTU_SERVER.MD](../INSTALL_UBUNTU_SERVER.MD)** - Ubuntu server installation

---

**Last Updated:** March 24, 2026  
**Version:** v3.5.0+  
**Maintained By:** @comgunner

---

*PicoClaw: Ultra-Efficient AI in Go. $10 Hardware · <45MB RAM · <1s Startup*
