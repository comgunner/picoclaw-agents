# Antigravity Models Reference

**Last Updated:** March 30, 2026  
**Version:** v1.3.0-alpha  
**Total Models:** 15

---

## Quick Reference Table

| # | Model Name | Full Model Path | Category | Speed | Quality | Best For |
|---|------------|-----------------|----------|-------|---------|----------|
| 1 | `gemini-3-flash` ⭐ | `antigravity/gemini-3-flash` | Gemini 3 | 🚀🚀🚀 | ⭐⭐⭐ | **DEFAULT - Best overall** |
| 2 | `gemini-3-pro-high` | `antigravity/gemini-3-pro-high` | Gemini 3 | 🚀🚀 | ⭐⭐⭐⭐⭐ | Complex reasoning |
| 3 | `gemini-3-pro-low` | `antigravity/gemini-3-pro-low` | Gemini 3 | 🚀🚀🚀 | ⭐⭐⭐ | Simple tasks |
| 4 | `gemini-3.1-pro-high` | `antigravity/gemini-3.1-pro-high` | Gemini 3.1 | 🚀🚀 | ⭐⭐⭐⭐⭐ | Advanced tasks |
| 5 | `gemini-3.1-pro-low` | `antigravity/gemini-3.1-pro-low` | Gemini 3.1 | 🚀🚀🚀 | ⭐⭐⭐ | Medium tasks |
| 6 | `gemini-3.1-flash-lite` | `antigravity/gemini-3.1-flash-lite` | Gemini 3.1 | 🚀🚀🚀🚀 | ⭐⭐ | Fast responses |
| 7 | `gemini-3-flash-agent` | `antigravity/gemini-3-flash-agent` | Gemini 3 | 🚀🚀🚀 | ⭐⭐⭐⭐ | Multi-step agents |
| 8 | `gemini-3-flash-preview` | `antigravity/gemini-3-flash-preview` | Gemini 3 | 🚀🚀🚀 | ⭐⭐⭐ | Testing features |
| 9 | `gemini-2.5-flash` | `antigravity/gemini-2.5-flash` | Gemini 2.5 | 🚀🚀🚀 | ⭐⭐⭐ | Fast responses |
| 10 | `gemini-2.5-flash-lite` | `antigravity/gemini-2.5-flash-lite` | Gemini 2.5 | 🚀🚀🚀🚀 | ⭐⭐ | Simple tasks |
| 11 | `gemini-2.5-flash-thinking` | `antigravity/gemini-2.5-flash-thinking` | Gemini 2.5 | 🚀🚀 | ⭐⭐⭐⭐ | Reasoning tasks |
| 12 | `gemini-2.5-pro` | `antigravity/gemini-2.5-pro` | Gemini 2.5 | 🚀🚀 | ⭐⭐⭐⭐ | General purpose |
| 13 | `claude-sonnet-4-6` | `antigravity/claude-sonnet-4-6` | Claude | 🚀🚀 | ⭐⭐⭐⭐⭐ | Writing, analysis |
| 14 | `claude-opus-4-6-thinking` | `antigravity/claude-opus-4-6-thinking` | Claude | 🚀 | ⭐⭐⭐⭐⭐⭐ | Complex problems |
| 15 | `gpt-oss-120b-medium` | `antigravity/gpt-oss-120b-medium` | GPT-OSS | 🚀🚀 | ⭐⭐⭐ | General use |

---

## Model Categories

### 🏆 Recommended Models

| Use Case | Model | Why |
|----------|-------|-----|
| **Default / General** | `gemini-3-flash` ⭐ | Best speed/quality balance |
| **Complex Reasoning** | `claude-opus-4-6-thinking` | Highest reasoning capability |
| **Writing & Analysis** | `claude-sonnet-4-6` | Best for natural language |
| **Fast Responses** | `gemini-3.1-flash-lite` | Lowest latency |
| **Agent Workflows** | `gemini-3-flash-agent` | Optimized for multi-step |
| **Budget/Quota Saving** | `gemini-2.5-flash-lite` | Lowest quota usage |

---

## Detailed Model Descriptions

### Gemini 3 Series

#### `gemini-3-flash` ⭐ (DEFAULT)
- **Speed:** 🚀🚀🚀
- **Quality:** ⭐⭐⭐
- **Context:** Large
- **Best For:** Daily tasks, general purpose
- **Quota Usage:** Medium
- **Recommendation:** Use as default for all agents

#### `gemini-3-pro-high`
- **Speed:** 🚀🚀
- **Quality:** ⭐⭐⭐⭐⭐
- **Context:** Very Large
- **Best For:** Complex reasoning, math, code analysis
- **Quota Usage:** High
- **Recommendation:** Use for difficult problems only

#### `gemini-3-pro-low`
- **Speed:** 🚀🚀🚀
- **Quality:** ⭐⭐⭐
- **Context:** Large
- **Best For:** Simple tasks with Gemini 3
- **Quota Usage:** Medium-Low
- **Recommendation:** Good balance for routine tasks

#### `gemini-3.1-pro-high`
- **Speed:** 🚀🚀
- **Quality:** ⭐⭐⭐⭐⭐
- **Context:** Very Large
- **Best For:** Advanced reasoning, latest Gemini 3.1
- **Quota Usage:** Very High
- **Recommendation:** Use for cutting-edge capabilities

#### `gemini-3.1-pro-low`
- **Speed:** 🚀🚀🚀
- **Quality:** ⭐⭐⭐
- **Context:** Large
- **Best For:** Medium complexity with 3.1
- **Quota Usage:** Medium
- **Recommendation:** Good for everyday 3.1 tasks

#### `gemini-3.1-flash-lite`
- **Speed:** 🚀🚀🚀🚀
- **Quality:** ⭐⭐
- **Context:** Medium
- **Best For:** Fast responses, simple queries
- **Quota Usage:** Low
- **Recommendation:** Use for quick lookups

#### `gemini-3-flash-agent`
- **Speed:** 🚀🚀🚀
- **Quality:** ⭐⭐⭐⭐
- **Context:** Large
- **Best For:** Multi-step agent workflows
- **Quota Usage:** Medium-High
- **Recommendation:** Use for autonomous agents

#### `gemini-3-flash-preview`
- **Speed:** 🚀🚀🚀
- **Quality:** ⭐⭐⭐
- **Context:** Large
- **Best For:** Testing new features
- **Quota Usage:** Medium
- **Recommendation:** Use for experimentation

---

### Gemini 2.5 Series

#### `gemini-2.5-flash`
- **Speed:** 🚀🚀🚀
- **Quality:** ⭐⭐⭐
- **Context:** Medium-Large
- **Best For:** Fast responses (fallback model)
- **Quota Usage:** Low-Medium
- **Recommendation:** Use as fallback

#### `gemini-2.5-flash-lite`
- **Speed:** 🚀🚀🚀🚀
- **Quality:** ⭐⭐
- **Context:** Medium
- **Best For:** Simple tasks, quota saving
- **Quota Usage:** Very Low
- **Recommendation:** Use when quota is low

#### `gemini-2.5-flash-thinking`
- **Speed:** 🚀🚀
- **Quality:** ⭐⭐⭐⭐
- **Context:** Large
- **Best For:** Reasoning tasks
- **Quota Usage:** Medium-High
- **Recommendation:** Use for logic puzzles

#### `gemini-2.5-pro`
- **Speed:** 🚀🚀
- **Quality:** ⭐⭐⭐⭐
- **Context:** Large
- **Best For:** General purpose (older gen)
- **Quota Usage:** Medium
- **Recommendation:** Good alternative to Gemini 3

---

### Claude Series

#### `claude-sonnet-4-6`
- **Speed:** 🚀🚀
- **Quality:** ⭐⭐⭐⭐⭐
- **Context:** Very Large
- **Best For:** Writing, analysis, natural language
- **Quota Usage:** High
- **Recommendation:** Best for content creation

#### `claude-opus-4-6-thinking`
- **Speed:** 🚀
- **Quality:** ⭐⭐⭐⭐⭐⭐
- **Context:** Maximum
- **Best For:** Complex problem solving, math, science
- **Quota Usage:** Very High
- **Recommendation:** Use for hardest problems only

---

### GPT-OSS Series

#### `gpt-oss-120b-medium`
- **Speed:** 🚀🚀
- **Quality:** ⭐⭐⭐
- **Context:** Medium
- **Best For:** General use, alternative to Gemini
- **Quota Usage:** Medium
- **Recommendation:** Good variety model

---

## Usage Examples

### Command Line

```bash
# Default model (gemini-3-flash)
./picoclaw-agents agent -m "Hello"

# Specific model
./picoclaw-agents agent -m "Analyze this code" --model gemini-3-pro-high

# Claude for writing
./picoclaw-agents agent -m "Write a blog post" --model claude-sonnet-4-6

# Complex reasoning
./picoclaw-agents agent -m "Solve this math problem" --model claude-opus-4-6-thinking

# Fast response
./picoclaw-agents agent -m "Quick question" --model gemini-3.1-flash-lite
```

### Web UI

1. Open http://localhost:18800/
2. Click model dropdown
3. Select desired model
4. Send message

---

## Model Selection Guide

### By Task Type

| Task Type | Recommended Model | Alternative |
|-----------|-------------------|-------------|
| **General Chat** | `gemini-3-flash` | `gemini-2.5-pro` |
| **Code Review** | `gemini-3-pro-high` | `claude-sonnet-4-6` |
| **Creative Writing** | `claude-sonnet-4-6` | `gemini-3-flash` |
| **Math/Logic** | `claude-opus-4-6-thinking` | `gemini-2.5-flash-thinking` |
| **Quick Answers** | `gemini-3.1-flash-lite` | `gemini-2.5-flash-lite` |
| **Research** | `gemini-3-pro-high` | `gemini-3.1-pro-high` |
| **Data Analysis** | `gemini-3-flash-agent` | `gemini-3-pro-high` |
| **Translation** | `claude-sonnet-4-6` | `gemini-3-flash` |
| **Summarization** | `gemini-3-flash` | `gemini-2.5-flash` |
| **Brainstorming** | `gemini-3-flash` | `gpt-oss-120b-medium` |

### By Quota Budget

| Budget Level | Recommended Models |
|--------------|-------------------|
| **Unlimited** | `claude-opus-4-6-thinking`, `gemini-3.1-pro-high` |
| **High** | `gemini-3-pro-high`, `claude-sonnet-4-6` |
| **Medium** | `gemini-3-flash`, `gemini-2.5-pro` |
| **Low** | `gemini-2.5-flash`, `gemini-3.1-flash-lite` |
| **Very Low** | `gemini-2.5-flash-lite` |

---

## Quota Management

### Estimated Quota Usage (per 1K tokens)

| Model | Input Cost | Output Cost | Relative Cost |
|-------|------------|-------------|---------------|
| `gemini-3-flash` | 1x | 1x | 1x (baseline) |
| `gemini-3-pro-high` | 5x | 5x | 5x |
| `gemini-3.1-pro-high` | 7x | 7x | 7x |
| `claude-sonnet-4-6` | 4x | 4x | 4x |
| `claude-opus-4-6-thinking` | 10x | 10x | 10x |
| `gemini-2.5-flash-lite` | 0.5x | 0.5x | 0.5x |

> **Note:** Actual quota usage depends on your Google One/Workspace plan.

---

## Troubleshooting

### Model Not Found

**Error:** `model 'gemini-3-flash' not found`

**Solution:**
```bash
# Re-run login to add all models
./picoclaw-agents auth login --provider google-antigravity

# Or use sync script
./scripts/sync_antigravity_models.sh
```

### Rate Limit (429)

**Error:** `429 Resource exhausted`

**Solutions:**
1. Switch to lighter model: `gemini-2.5-flash-lite`
2. Wait for quota reset (check `auth status`)
3. Use fallback chain in config

### Model Restricted

**Error:** `Model restricted for this project`

**Solution:** Try different model:
- `gemini-3-flash` (most reliable)
- `gemini-2.5-flash` (fallback)

---

## Related Documentation

- [`ANTIGRAVITY.md`](./ANTIGRAVITY.md) - Full Antigravity guide
- [`ANTIGRAVITY.es.md`](./ANTIGRAVITY.es.md) - Spanish version
- [`PICO_MODEL_OVERRIDE.md`](./PICO_MODEL_OVERRIDE.md) - Web UI model override
- [`../local_work/FIX_WEBUI_MODEL_ANTIGRAVITY.md`](../local_work/FIX_WEBUI_MODEL_ANTIGRAVITY.md) - Implementation details

---

*Quick Reference - Antigravity Models - v1.3.0-alpha*
