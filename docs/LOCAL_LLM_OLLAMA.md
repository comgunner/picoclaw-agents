# Running PicoClaw-Agents with Local LLMs via Ollama

Run AI models **100% offline** on your own hardware — no API keys, no cloud, no data leaving your machine.

---

## What is Ollama?

[Ollama](https://ollama.com) is the easiest way to run large language models locally. It provides a simple API compatible with OpenAI's format, which means picoclaw-agents connects to it out of the box.

---

## 1. Install Ollama

### macOS

```bash
# Option A — Direct download (recommended)
# https://ollama.com/download/mac
# Download and run the .dmg installer

# Option B — One-liner
brew install ollama
```

### Linux

```bash
curl -fsSL https://ollama.com/install.sh | sh
```

### Windows

```powershell
irm https://ollama.com/install.ps1 | iex
```

### Termux (Android)

```bash
pkg update
pkg install ollama
```

---

## 2. Find the Right Model for Your Hardware

Not sure which model fits your machine? Use `llm-checker` to get personalized recommendations:

```bash
# Install
npm install -g llm-checker

# Detect your hardware
llm-checker hw-detect

# Get recommendations by category
llm-checker recommend --category coding
llm-checker recommend --category general

# Run with auto-selection
llm-checker ai-run --category coding --prompt "Write a hello world in Python"
```

**Example output on Mac M1 (8GB RAM):**
```
SYSTEM INFORMATION
│ CPU: M1 (8 cores, 2.4GHz)
│ Architecture: Apple Silicon
│ RAM: 8GB
│ GPU: Apple M1 / VRAM: 4GB (Integrated)
│ Hardware Tier: LOW

INTELLIGENT RECOMMENDATIONS
│ BEST OVERALL: yi:6b          → ollama pull yi:6b
│ Coding:       deepseek-coder:6.7b
│ Reasoning:    deepseek-coder:33b
│ Multimodal:   qwen2.5vl:7b
│ Creative:     yi:6b
│ Chat:         yi:6b
│ Reading:      qwen3:1.7b
```

> **Tip:** For Termux (Android), also install llm-checker: `npm install -g llm-checker`

---

## 3. Recommended Models (Curated by picoclaw-agents)

These models are lightweight, fast, and work well as agent backends:

| Model | Size | Best For | Pull Command |
|-------|------|----------|--------------|
| `llama3.2:1b` | ~800MB | General chat, fast responses | `ollama pull llama3.2:1b` |
| `qwen2.5:0.5b` | ~400MB | Ultra-lightweight, low RAM | `ollama pull qwen2.5:0.5b` |
| `qwen2.5-coder:0.5b` | ~400MB | Code generation, minimal footprint | `ollama pull qwen2.5-coder:0.5b` |

**Pull and run:**
```bash
ollama pull llama3.2:1b
ollama run llama3.2:1b
```

### Top 10 Models from the Ollama Library

Visit **https://ollama.com/library** for the full catalog. Current top models:

| # | Model | Highlights |
|---|-------|-----------|
| 1 | `gemma3` | Google Gemma 3, multilingual |
| 2 | `llama3.2` | Meta Llama 3.2, fast & capable |
| 3 | `qwen2.5` | Alibaba Qwen 2.5, strong at code |
| 4 | `phi4` | Microsoft Phi-4, small but smart |
| 5 | `mistral` | Mistral 7B, strong reasoning |
| 6 | `deepseek-r1` | DeepSeek R1, chain-of-thought |
| 7 | `llava` | Multimodal (text + images) |
| 8 | `codellama` | Meta CodeLlama, code-focused |
| 9 | `deepseek-coder-v2` | DeepSeek Coder V2, top coding |
| 10 | `nomic-embed-text` | Text embeddings, RAG use cases |

---

## 4. Connect picoclaw-agents to Ollama

### Option A — Edit `~/.picoclaw/config.json` directly

Add entries to your `model_list`:

```json
{
  "model_list": [
    {
      "model_name": "llama3.2:1b",
      "model": "llama3.2:1b",
      "api_base": "http://localhost:11434/v1",
      "api_key": "ollama"  # pragma: allowlist secret
    },
    {
      "model_name": "qwen2.5:0.5b",
      "model": "qwen2.5:0.5b",
      "api_base": "http://localhost:11434/v1",
      "api_key": "ollama"  # pragma: allowlist secret
    },
    {
      "model_name": "qwen2.5-coder:0.5b",
      "model": "qwen2.5-coder:0.5b",
      "api_base": "http://localhost:11434/v1",
      "api_key": "ollama"  # pragma: allowlist secret
    }
  ]
}
```

> **Note:** `api_key: "ollama"` `# pragma: allowlist secret` is required by the OpenAI-compat client even though Ollama doesn't use authentication. Any non-empty string works.

### Option B — Web UI (picoclaw-agents-launcher)

1. Start the launcher: `picoclaw-agents-launcher`
2. Open **http://localhost:18800/models**
3. Click **+ Add Model**
4. Fill in the form:
   - **Name:** `llama3.2:1b`
   - **Model:** `llama3.2:1b`
   - **API Base:** `http://localhost:11434/v1`
   - **API Key:** `ollama`
5. Check **Default Model** if you want it as default
6. Click **Save**

---

## 5. Run the Agent

```bash
# Single message
picoclaw-agents agent --model llama3.2:1b -m "Hello, are you running locally?"

# Interactive mode
picoclaw-agents agent --model qwen2.5-coder:0.5b

# Set as default in config and just run
picoclaw-agents agent -m "Write a Python script to list files"
```

---

## 6. Verify Ollama is Running

```bash
# Check Ollama status
ollama list

# Test the API directly
curl http://localhost:11434/v1/models

# Ollama runs on port 11434 by default
# Start manually if needed:
ollama serve
```

---

## Hardware Tips

| RAM | Recommended Models |
|-----|--------------------|
| 4GB | `qwen2.5:0.5b`, `llama3.2:1b` |
| 8GB | `llama3.2:3b`, `qwen2.5:3b`, `deepseek-coder:6.7b` |
| 16GB | `llama3.2:8b`, `qwen2.5:7b`, `mistral:7b` |
| 32GB+ | `llama3.1:70b` (quantized), `deepseek-r1:32b` |

- **Apple Silicon (M1/M2/M3):** Ollama uses Metal GPU acceleration automatically — models run faster than on equivalent Intel hardware
- **NVIDIA GPU:** CUDA is used automatically if available
- **CPU-only:** Works fine for 1B–3B models, slower for larger ones

---

## Troubleshooting

**`connection refused` error:**
```bash
# Ollama is not running — start it:
ollama serve
```

**Model not found:**
```bash
# Pull the model first:
ollama pull llama3.2:1b
```

**Out of memory:**
```bash
# Use a smaller model:
ollama pull qwen2.5:0.5b   # only ~400MB
```

**Slow responses:**
```bash
# Check which GPU backend is active:
ollama ps
```
