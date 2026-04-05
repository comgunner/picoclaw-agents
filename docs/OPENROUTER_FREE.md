# OpenRouter Free Tier with PicoClaw-Agents

**Last Updated:** March 28, 2026
**Version:** v1.3.0-alpha-fix901

---

## 🆓 Quick Start — Free without API Key

PicoClaw-Agents supports **100% free** models from OpenRouter without requiring a credit card.

### Option 1: Interactive Onboard (Recommended)

```bash
# Run setup wizard
picoclaw-agents onboard --free
```

**The wizard will guide you:**
1. Request your OpenRouter API key (free)
2. Configure free models automatically
3. Create your workspace

### Option 2: Manual Config

Create `~/.picoclaw/config.json`:

```json
{
  "agents": {
    "defaults": {
      "workspace": "~/.picoclaw/workspace",
      "model": "openrouter/auto",
      "max_tokens": 8192,
      "max_tool_iterations": 20
    }
  },
  "model_list": [
    {
      "model_name": "or-auto",
      "model": "openrouter/auto",
      "api_key": "sk-or-v1-YOUR_API_KEY_HERE"  // pragma: allowlist secret
    }
  ]
}
```

---

## 🎯 Available Free Models

OpenRouter offers several free models. PicoClaw configures them automatically:

### Default Configuration (`onboard --free`)

| Priority | Model | Context | Recommended Use |
|----------|-------|---------|-----------------|
| **1** | `openrouter/auto` | Variable | Auto-selects best free model |
| **2** | `stepfun/step-3.5-flash` | 256K | Long context, reasoning |
| **3** | `deepseek/deepseek-v3.2-20251201` | 64K | Fast inference |

### What is `openrouter/auto`?

- **Auto-selection:** OpenRouter automatically chooses the best available free model
- **Automatic fallback:** If one model fails, uses the next one
- **No configuration:** No need to specify individual models

---

## 📝 Useful Commands

### Check Configuration

```bash
# Check current model
picoclaw-agents agent --model "openrouter/auto" -m "What model are you using?"

# Check auth status
picoclaw-agents auth status
```

### Test Individual Models

```bash
# StepFun (256K context)
picoclaw-agents agent --model "stepfun/step-3.5-flash" -m "Hello"

# DeepSeek (fast)
picoclaw-agents agent --model "deepseek/deepseek-v3.2-20251201" -m "Hello"
```

### Change Default Model

Edit `~/.picoclaw/config.json`:

```json
{
  "agents": {
    "defaults": {
      "model": "stepfun/step-3.5-flash"  // ← Change here
    }
  }
}
```

---

## 🔑 Get OpenRouter API Key

### Step 1: Create Account

1. Go to https://openrouter.ai
2. Click "Sign Up"
3. Register with email (no card required)

### Step 2: Create API Key

1. Go to https://openrouter.ai/keys
2. Click "Create Key"
3. Copy the key (starts with `sk-or-v1-...`)

### Step 3: Configure in PicoClaw

```bash
# During onboard
picoclaw-agents onboard --free
# → Paste your API key when prompted

# Or edit config.json manually
nano ~/.picoclaw/config.json
# → Paste your API key in "api_key"
```

---

## ⚠️ Common Issues

### Error: "openrouter-free is not a valid model ID"

**Cause:** Old configuration with invalid model name.

**Solution:** Update your config:

```bash
# Option 1: Re-run onboard
picoclaw-agents onboard --free

# Option 2: Edit config manually
sed -i 's|"openrouter/free"|"openrouter/auto"|g' ~/.picoclaw/config.json
sed -i 's|"openrouter-free"|"openrouter/auto"|g' ~/.picoclaw/config.json
```

### Error: "401 Unauthorized"

**Cause:** Invalid or expired API key.

**Solution:**
1. Verify your key at https://openrouter.ai/keys
2. Update in `~/.picoclaw/config.json`
3. Re-run `picoclaw-agents onboard --free`

### Error: "Rate limit exceeded"

**Cause:** Free request limit reached.

**Free Tier Limits:**
- ~50 requests/minute
- ~1000 requests/day (varies by model)

**Solution:**
- Wait a few minutes
- Use slower models with higher limits
- Consider upgrading to paid tier

---

## 📊 Free Tier Limits

### Rate Limits

| Model | Requests/min | Requests/day | Max Context |
|-------|--------------|--------------|-------------|
| `openrouter/auto` | ~50 | ~1000 | Variable |
| `stepfun/step-3.5-flash` | ~20 | ~500 | 256K |
| `deepseek/deepseek-v3.2` | ~30 | ~800 | 64K |

### Best Practices

1. **Use `openrouter/auto`** — Best availability/speed balance
2. **Avoid polling** — Don't make requests too frequently
3. **Batch tasks** — Group tasks when possible
4. **Monitor usage** — Check your dashboard at openrouter.ai

---

## 🚀 Usage Examples

### Simple Chat

```bash
picoclaw-agents agent -m "What is the capital of France?"
```

### Code Task

```bash
picoclaw-agents agent -m "Create a Python function that calculates Fibonacci"
```

### Web Search

```bash
picoclaw-agents agent -m "Search for the latest AI news"
```

### Complex Task (Multi-agent)

```bash
# With team mode configured
picoclaw-agents agent -m "Create a REST API with Node.js and Express"
```

---

## 📚 Additional Resources

### Official Links

- **OpenRouter:** https://openrouter.ai
- **Free Models:** https://openrouter.ai/models?order=-free
- **API Docs:** https://openrouter.ai/docs
- **Keys:** https://openrouter.ai/keys

### PicoClaw Documentation

- **CHANGELOG:** [CHANGELOG.md](../CHANGELOG.md)
- **README:** [README.md](../README.md)
- **Fix #901:** [local_work/fix_901_openrouter_normalization.md](../local_work/fix_901_openrouter_normalization.md)

---

## ❓ FAQ

### Is it really free?

Yes. OpenRouter offers free models without a credit card. There are rate limits but they're sufficient for personal use.

### Do I need to configure anything else?

No. With `picoclaw-agents onboard --free` everything is configured automatically.

### Can I switch to paid models later?

Yes. Just edit `config.json` and change the model or add your card in OpenRouter.

### What happens if free models run out?

OpenRouter auto-selects another available free model. If all are exhausted, you'll receive a rate limit error.

### Does it work with all channels (Telegram, Discord)?

Yes. The free tier works the same for CLI, Telegram, Discord, etc.

---

**Document created:** March 28, 2026
**Version:** v1.3.0-alpha-fix901
**Maintainer:** @comgunner
