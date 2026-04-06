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

## Popular Models Organized by RAM

Visit **https://ollama.com/library** for the full catalog. Models grouped by minimum system requirements:

### Models for ≤ 4 GB RAM

| Model | Size | Best For | Pull Command |
|-------|------|----------|--------------|
| `qwen2.5:0.5b` | ~400MB | Ultra-lightweight, minimal RAM | `ollama pull qwen2.5:0.5b` |
| `qwen2.5-coder:0.5b` | ~400MB | Code generation, tiny footprint | `ollama pull qwen2.5-coder:0.5b` |
| `llama3.2:1b` | ~800MB | General chat, fast responses | `ollama pull llama3.2:1b` |
| `qwen2.5:1.5b` | ~1GB | Multilingual, general use | `ollama pull qwen2.5:1.5b` |
| `deepseek-r1:1.5b` | ~1GB | Chain-of-thought reasoning | `ollama pull deepseek-r1:1.5b` |
| `qwen3:0.6b` | ~400MB | Ultra-lightweight, fast chat | `ollama pull qwen3:0.6b` |
| `qwen3:1.7b` | ~1GB | Balanced speed and capability | `ollama pull qwen3:1.7b` |
| `nomic-embed-text` | ~270MB | Text embeddings, RAG | `ollama pull nomic-embed-text` |

### Models for 8 GB RAM

| Model | Size | Best For | Pull Command |
|-------|------|----------|--------------|
| `llama3.2:3b` | ~2GB | Balanced speed and capability | `ollama pull llama3.2:3b` |
| `qwen2.5:3b` | ~2GB | Strong all-around performance | `ollama pull qwen2.5:3b` |
| `qwen3:4b` | ~3GB | Latest Qwen, excellent reasoning | `ollama pull qwen3:4b` |
| `qwen3-coder:4b` | ~3GB | Latest Qwen Coder, code generation | `ollama pull qwen3-coder:4b` |
| `deepseek-coder:6.7b` | ~4GB | Code generation and understanding | `ollama pull deepseek-coder:6.7b` |
| `mistral:7b` | ~4GB | Strong reasoning, general purpose | `ollama pull mistral:7b` |
| `llava:7b` | ~4GB | Multimodal (text + images) | `ollama pull llava:7b` |
| `codellama:7b` | ~4GB | Code-focused, Meta quality | `ollama pull codellama:7b` |
| `gemma4:e2b` | 7.2 GB | Google's latest, multilingual | `ollama pull gemma4:e2b` |

### Models for 16 GB RAM or More

| Model | Size | Best For | Pull Command |
|-------|------|----------|--------------|
| `llama3.2:8b` | ~5GB | High-quality general purpose | `ollama pull llama3.2:8b` |
| `qwen2.5:7b` | ~5GB | Excellent at code and reasoning | `ollama pull qwen2.5:7b` |
| `qwen3:8b` | ~5GB | Latest Qwen, strong reasoning | `ollama pull qwen3:8b` |
| `qwen3-coder:8b` | ~5GB | Latest Qwen Coder, top coding | `ollama pull qwen3-coder:8b` |
| `phi4:14b` | ~9GB | Microsoft's smartest small model | `ollama pull phi4:14b` |
| `qwen3:14b` | ~9GB | High-quality Qwen, multilingual | `ollama pull qwen3:14b` |
| `qwen3-coder:14b` | ~9GB | Advanced Qwen Coder | `ollama pull qwen3-coder:14b` |
| `deepseek-r1:7b` | ~5GB | Deep reasoning, chain-of-thought | `ollama pull deepseek-r1:7b` |
| `deepseek-coder-v2:16b` | ~10GB | Top-tier coding model | `ollama pull deepseek-coder-v2:16b` |
| `qwen3:32b` | ~20GB | Maximum Qwen quality | `ollama pull qwen3:32b` |
| `gemma4:e4b` | 9.6 GB | Google's latest, latest variant | `ollama pull gemma4:e4b` |
| `gemma4:26b` | 18 GB | Best quality Gemma 4 | `ollama pull gemma4:26b` |
| `gemma4:31b` | 20 GB | Maximum quality Gemma 4 | `ollama pull gemma4:31b` |
| `llama3.1:70b` (Q4) | ~40GB | Maximum capability (quantized) | `ollama pull llama3.1:70b` |

### 🆕 Google Gemma 4 — Detailed Requirements

Gemma 4 is Google's latest open model family, offering excellent multilingual capabilities. Visit **[https://ollama.com/library/gemma4](https://ollama.com/library/gemma4)** for the full catalog.

| Variant | Weight | Min RAM/VRAM | Ideal RAM/VRAM | Pull Command |
|---------|--------|--------------|----------------|--------------|
| `gemma4:e2b` | 7.2 GB | 8 GB | 12 GB | `ollama pull gemma4:e2b` |
| `gemma4:e4b` (latest) | 9.6 GB | 12 GB | 16 GB | `ollama pull gemma4:e4b` |
| `gemma4:26b` | 18 GB | 24 GB | 24 GB+ | `ollama pull gemma4:26b` |
| `gemma4:31b` | 20 GB | 24 GB | 32 GB | `ollama pull gemma4:31b` |

#### CPU vs GPU for Gemma 4

- **CPU:** Ollama works with any modern CPU from recent years (requires AVX instruction support). However, processing these models using only CPU is slow.
- **GPU (Recommended):** Ollama automatically detects if you have a graphics card (Nvidia, AMD, or Apple Silicon chip) and sends the model there for much faster performance. The key is that your GPU has enough VRAM to hold the chosen model's weight.

#### Hardware Tips for Gemma 4

| Your Hardware | Recommended Variant |
|---------------|--------------------|
| 8 GB RAM | `gemma4:e2b` (will be slow, 12 GB ideal) |
| 12 GB RAM | `gemma4:e2b` or `gemma4:e4b` |
| 16 GB RAM | `gemma4:e4b` (latest) |
| 24 GB VRAM (GPU) | `gemma4:26b` |
| 32 GB VRAM (GPU) | `gemma4:31b` (best quality) |

> **Note for PicoClaw-Agents users:** Gemma 4 variants have moderate context windows. For best results with the agent loop, use `gemma4:e2b` or `gemma4:e4b` on systems with limited RAM. The larger variants (26b, 31b) require significant VRAM and work best on dedicated GPU setups.

### 🆕 Qwen 3 & Qwen 3 Coder — Detailed Requirements

Qwen 3 is Alibaba's latest open model family with excellent multilingual and reasoning capabilities. Qwen 3 Coder is the specialized variant for code generation. Visit **[https://ollama.com/library/qwen3](https://ollama.com/library/qwen3)** and **[https://ollama.com/library/qwen3-coder](https://ollama.com/library/qwen3-coder)** for the full catalogs.

#### Qwen 3 (General Purpose)

| Variant | Weight | Min RAM/VRAM | Ideal RAM/VRAM | Pull Command |
|---------|--------|--------------|----------------|--------------|
| `qwen3:0.6b` | ~400MB | 2 GB | 4 GB | `ollama pull qwen3:0.6b` |
| `qwen3:1.7b` | ~1GB | 4 GB | 4 GB | `ollama pull qwen3:1.7b` |
| `qwen3:4b` | ~3GB | 6 GB | 8 GB | `ollama pull qwen3:4b` |
| `qwen3:8b` | ~5GB | 8 GB | 16 GB | `ollama pull qwen3:8b` |
| `qwen3:14b` | ~9GB | 12 GB | 16 GB | `ollama pull qwen3:14b` |
| `qwen3:32b` | ~20GB | 24 GB | 32 GB | `ollama pull qwen3:32b` |

#### Qwen 3 Coder (Code Generation)

| Variant | Weight | Min RAM/VRAM | Ideal RAM/VRAM | Pull Command |
|---------|--------|--------------|----------------|--------------|
| `qwen3-coder:4b` | ~3GB | 6 GB | 8 GB | `ollama pull qwen3-coder:4b` |
| `qwen3-coder:8b` | ~5GB | 8 GB | 16 GB | `ollama pull qwen3-coder:8b` |
| `qwen3-coder:14b` | ~9GB | 12 GB | 16 GB | `ollama pull qwen3-coder:14b` |
| `qwen3-coder:32b` | ~20GB | 24 GB | 32 GB | `ollama pull qwen3-coder:32b` |

#### CPU vs GPU for Qwen 3

- **CPU:** Works with any modern CPU (AVX support required). Smaller variants (0.6b, 1.7b) run reasonably on CPU alone. Larger variants will be slow.
- **GPU (Recommended):** Ollama automatically detects Nvidia, AMD, or Apple Silicon GPUs and offloads computation for much faster inference. The key is having enough VRAM for your chosen variant.

#### Hardware Tips for Qwen 3

| Your Hardware | Recommended Variant |
|---------------|--------------------|
| 4 GB RAM | `qwen3:0.6b` or `qwen3:1.7b` |
| 8 GB RAM | `qwen3:4b` or `qwen3-coder:4b` |
| 16 GB RAM | `qwen3:8b`, `qwen3:14b`, or coder variants |
| 32 GB RAM/VRAM | `qwen3:32b` or `qwen3-coder:32b` (best quality) |

> **Note for PicoClaw-Agents users:** Qwen 3 models have good context windows and work well with the agent loop. For coding tasks, prefer `qwen3-coder` variants. For general use, `qwen3:4b` on 8GB systems or `qwen3:8b` on 16GB systems provide the best balance.

---

## 3b. Limiting RAM, CPU, and GPU Usage

> **Verified:** All parameters below are **official Ollama settings** documented in the [Ollama Modelfile](https://github.com/ollama/ollama/blob/main/docs/modelfile.md) specification.

### Where Ollama Configuration Lives (Per OS)

Ollama does not use a `.env` file. Configuration is managed through environment variables when starting the server. Here's where and how to set them:

| OS | Config File / Location | Example |
|----|----------------------|---------|
| **macOS** | `~/.config/ollama/config.json` (rarely used)<br>Or launchd plist: `~/Library/LaunchAgents/com.ollama.ollama.plist` | See macOS section below |
| **Linux** | systemd service override: `systemctl edit ollama.service`<br>Or `~/.bashrc` / `~/.zshrc` | See Linux section below |
| **Windows** | System Environment Variables (Settings → System → About → Advanced → Environment Variables)<br>Or PowerShell profile: `$PROFILE` | See Windows section below |
| **Termux** | `~/.bashrc` or `~/.zshrc` | See Termux section below |

**Common environment variables:**

| Variable | Purpose | Default |
|----------|---------|---------|
| `OLLAMA_HOST` | Bind address (e.g., `0.0.0.0:11434` for network access) | `127.0.0.1:11434` |
| `OLLAMA_KEEP_ALIVE` | How long model stays in RAM after use (`5m`, `1h`, `-1` forever) | `5m` |
| `OLLAMA_NUM_PARALLEL` | Max concurrent requests | `1` |
| `OLLAMA_MAX_LOADED_MODELS` | Max models loaded simultaneously | `1` |
| `OLLAMA_GPU_ENABLED` | Set to `0` to disable GPU entirely | `1` |
| `OLLAMA_TMPDIR` | Temp directory for model loading | System temp |

### Three Ways to Apply Resource Limits

### Method A: Via `/set` in the CLI (Interactive, Session Only)

When running a model interactively (`ollama run llama3`), adjust parameters on the fly:

```
/set parameter num_thread 4
/set parameter num_ctx 2048
/set parameter num_gpu 10
```

Changes take effect immediately but **are lost when you exit the session**.

### Method B: Via Modelfile (Permanent — Recommended for picoclaw-agents)

#### Where to Save the Modelfile

The Modelfile is a plain text file that **you create and keep wherever you want**. Ollama reads it once to build your custom model, then the model is stored permanently in Ollama's internal storage. The Modelfile itself can be deleted after building — but it's good practice to keep it for reference.

| OS | Recommended Location | Internal Model Storage |
|----|---------------------|----------------------|
| **macOS** | `~/ollama-modelfiles/` | `~/.ollama/models/` |
| **Linux** | `~/ollama-modelfiles/` | `~/.ollama/models/` |
| **Windows** | `C:\Users\You\ollama-modelfiles\` | `C:\Users\You\.ollama\models\` |
| **Termux** | `~/ollama-modelfiles/` | `~/.ollama/models/` |

**Workflow:**

```bash
# 1. Create a directory for your Modelfiles (anywhere you want)
mkdir -p ~/ollama-modelfiles
cd ~/ollama-modelfiles

# 2. Create the Modelfile (use any text editor)
nano Modelfile    # or: vim Modelfile, code Modelfile

# 3. Build the custom model (Ollama reads the file once)
ollama create my-custom-model -f Modelfile

# 4. The model is now stored in Ollama's internal storage
#    You can delete the Modelfile if you want, but keeping it is useful
ollama list

# 5. Use the model
ollama run my-custom-model
```

**Key concept:** The Modelfile is like a **recipe**. Once you bake the cake (`ollama create`), you don't need the recipe anymore — but it's handy if you want to bake another one later.

#### Example 1: Gemma 4:e2b — Limited to 8GB RAM

For a Mac or laptop with 8GB RAM, keeping RAM usage under control while still running Gemma 4:

```bash
mkdir -p ~/ollama-modelfiles
nano ~/ollama-modelfiles/Modelfile-gemma4-8gb
```

Paste the following content:

```Modelfile
FROM gemma4:e2b

# CPU threads — minimum that still works (4 is safe for most CPUs)
PARAMETER num_thread 4

# Context window — reduced from default 8192 to 2048
# This saves ~600MB+ of RAM during inference
PARAMETER num_ctx 2048

# GPU layers — on 8GB unified memory, leave room for OS + picoclaw
# gemma4:e2b has ~40 layers total; 20 on GPU leaves ~20 on CPU
PARAMETER num_gpu 20

# Batch size — smaller batches = less RAM at once
PARAMETER num_batch 128

# Keep first 4 tokens (system prompt) always in context
PARAMETER num_keep 4

# Optional: custom system prompt for picoclaw-agents
SYSTEM You are PicoClaw, a helpful AI assistant. Be concise and action-oriented.
```

Save with `Ctrl+O`, then `Enter`, then `Ctrl+X`.

Build and connect:

```bash
ollama create picoclaw-gemma4-8gb -f ~/ollama-modelfiles/Modelfile-gemma4-8gb

# Verify
ollama list | grep gemma

# Test it
ollama run picoclaw-gemma4-8gb "Hello, how much RAM do you use?"

# Configure picoclaw-agents — add to ~/.picoclaw/config.json:
```

```json
{
  "model_list": [
    {
      "model_name": "picoclaw-gemma4-8gb",
      "model": "picoclaw-gemma4-8gb",
      "api_base": "http://localhost:11434/v1",
      "api_key": "ollama"
    }
  ],
  "agents": {
    "defaults": {
      "model_name": "picoclaw-gemma4-8gb",
      "max_tokens": 2048
    }
  }
}
```

**Expected RAM usage:** ~5-6GB total (model ~3GB + context + overhead)

**Use it:**
```bash
./build/picoclaw-agents agent --model picoclaw-gemma4-8gb -m "Hello"
```

#### Example 2: Qwen 3:8b — Minimal Settings for 16GB System

For a desktop/laptop with 16GB RAM, running Qwen 3:8b efficiently:

```bash
nano ~/ollama-modelfiles/Modelfile-qwen3-16gb
```

Paste the following content:

```Modelfile
FROM qwen3:8b

# CPU threads — match your CPU core count (6 is conservative)
PARAMETER num_thread 6

# Context window — 4096 tokens, balanced for agent loop
PARAMETER num_ctx 4096

# GPU layers — qwen3:8b has ~35 layers; offload 25 to GPU, 10 stay on CPU
PARAMETER num_gpu 25

# Batch size — moderate for good throughput without spiking RAM
PARAMETER num_batch 256

# Keep system prompt tokens
PARAMETER num_keep 4
```

Save with `Ctrl+O`, then `Enter`, then `Ctrl+X`.

Build and connect:

```bash
ollama create picoclaw-qwen3-16gb -f ~/ollama-modelfiles/Modelfile-qwen3-16gb

# Test
ollama run picoclaw-qwen3-16gb "Write a Python function to reverse a string"

# Configure picoclaw-agents:
```

```json
{
  "model_list": [
    {
      "model_name": "picoclaw-qwen3-16gb",
      "model": "picoclaw-qwen3-16gb",
      "api_base": "http://localhost:11434/v1",
      "api_key": "ollama"
    }
  ],
  "agents": {
    "defaults": {
      "model_name": "picoclaw-qwen3-16gb",
      "max_tokens": 4096
    }
  }
}
```

**Expected RAM/VRAM usage:** ~6-8GB total

**Use it:**
```bash
./build/picoclaw-agents agent --model picoclaw-qwen3-16gb -m "Hello"
```

#### Example 3: Qwen 2.5-Coder:0.5b — Absolute Minimum (Ultra-Low RAM)

For the smallest possible footprint — Termux, Raspberry Pi, or any constrained system. This is the **minimum viable configuration**:

```bash
nano ~/ollama-modelfiles/Modelfile-qwen-coder-minimal
```

Paste the following content:

```Modelfile
FROM qwen2.5-coder:0.5b

# Minimum threads — 2 is the lowest that still works
PARAMETER num_thread 2

# Minimum context — 512 tokens is the absolute floor
# (below this the model may error or produce garbage)
PARAMETER num_ctx 512

# No GPU offload — this model is so small it fits entirely in CPU RAM
# Setting num_gpu 0 ensures zero VRAM usage
PARAMETER num_gpu 0

# Smallest batch size — minimum RAM during prompt processing
PARAMETER num_batch 64

# No keep — save every token of the tiny context window
PARAMETER num_keep 0
```

Save with `Ctrl+O`, then `Enter`, then `Ctrl+X`.

Build and connect:

```bash
ollama create picoclaw-coder-minimal -f ~/ollama-modelfiles/Modelfile-qwen-coder-minimal

# Test — note: short responses due to tiny context window
ollama run picoclaw-coder-minimal "def hello():"

# Configure picoclaw-agents:
```

```json
{
  "model_list": [
    {
      "model_name": "picoclaw-coder-minimal",
      "model": "picoclaw-coder-minimal",
      "api_base": "http://localhost:11434/v1",
      "api_key": "ollama"
    }
  ],
  "agents": {
    "defaults": {
      "model_name": "picoclaw-coder-minimal",
      "max_tokens": 512,
      "max_tool_iterations": 5
    }
  }
}
```

**Expected RAM usage:** ~600MB total (model ~400MB + context + overhead)
**VRAM usage:** 0MB (CPU-only)

**Use it:**
```bash
./build/picoclaw-agents agent --model picoclaw-coder-minimal -m "def hello():"
```

This configuration works on:
- Termux (Android, 2GB+ RAM)
- Raspberry Pi 4 (2GB RAM)
- Old laptops (2GB+ RAM)
- Any system where you need AI with minimal footprint

---

#### 🍎 Bonus: Mac Mini M1 — 1GB RAM Limit (4 Ultra-Low Configurations)

For Mac Mini M1 running other workloads, where you want to cap Ollama at **1GB RAM maximum**.

> **⚠️ Important:** First create the Modelfiles directory:
> ```bash
> mkdir -p ~/ollama-modelfiles
> cd ~/ollama-modelfiles
> ```

##### Configuration 1: Qwen 2.5:0.5b — Bare Minimum Agent

The absolute smallest Qwen model. Runs on virtually nothing:

```bash
nano ~/ollama-modelfiles/Modelfile-qwen25-min
```

Paste the following content:

```Modelfile
FROM qwen2.5:0.5b

# Absolute minimum threads (M1 has 8 cores, but we use 2 to save RAM)
PARAMETER num_thread 2

# Minimum viable context — 512 tokens
PARAMETER num_ctx 512

# No GPU offload — keep everything in CPU to control memory precisely
PARAMETER num_gpu 0

# Smallest batch size
PARAMETER num_batch 32

# No keep
PARAMETER num_keep 0
```

Save with `Ctrl+O`, then `Enter`, then `Ctrl+X`.

Build:

```bash
ollama create picoclaw-qwen25-min -f ~/ollama-modelfiles/Modelfile-qwen25-min
```

**Expected RAM:** ~500MB | **VRAM:** 0MB | **Speed:** ~5-8 tokens/sec on M1

##### Configuration 2: Qwen 3:0.6b — Tiny but Capable

Slightly larger than 0.5b but still very lightweight:

```bash
nano ~/ollama-modelfiles/Modelfile-qwen3-tiny
```

Paste the following content:

```Modelfile
FROM qwen3:0.6b

# 2 threads — barely uses CPU
PARAMETER num_thread 2

# Tiny context — enough for simple Q&A
PARAMETER num_ctx 512

# CPU-only for predictable memory
PARAMETER num_gpu 0

# Minimal batch
PARAMETER num_batch 32

PARAMETER num_keep 0
```

Save with `Ctrl+O`, then `Enter`, then `Ctrl+X`.

Build:

```bash
ollama create picoclaw-qwen3-tiny -f ~/ollama-modelfiles/Modelfile-qwen3-tiny
```

**Expected RAM:** ~550MB | **VRAM:** 0MB | **Speed:** ~6-10 tokens/sec on M1

##### Configuration 3: Qwen 2.5-Coder:0.5b — Minimal Coding Assistant

The smallest model with coding capability:

```bash
nano ~/ollama-modelfiles/Modelfile-coder-tiny
```

Paste the following content:

```Modelfile
FROM qwen2.5-coder:0.5b

# 2 threads — minimum
PARAMETER num_thread 2

# 512 tokens context — enough for short code snippets
PARAMETER num_ctx 512

# CPU-only
PARAMETER num_gpu 0

# Tiny batch
PARAMETER num_batch 32

PARAMETER num_keep 0
```

Save with `Ctrl+O`, then `Enter`, then `Ctrl+X`.

Build:

```bash
ollama create picoclaw-coder-tiny -f ~/ollama-modelfiles/Modelfile-coder-tiny
```

**Expected RAM:** ~500MB | **VRAM:** 0MB | **Speed:** ~5-8 tokens/sec on M1

##### Configuration 4: Gemma 2:2b — Smallest Gemma Available

The smallest Gemma model available is `gemma2:2b` (~1.7GB model weight). To fit within ~1GB RAM, we use **aggressive limits**:

```bash
nano ~/ollama-modelfiles/Modelfile-gemma2-tiny
```

Paste the following content:

```Modelfile
FROM gemma2:2b

# 2 threads — minimum CPU usage
PARAMETER num_thread 2

# Smallest possible context for this model
PARAMETER num_ctx 512

# No GPU offload — prevents VRAM spikes
PARAMETER num_gpu 0

# Absolute minimum batch
PARAMETER num_batch 32

PARAMETER num_keep 0
```

Save with `Ctrl+O`, then `Enter`, then `Ctrl+X`.

Build:

```bash
ollama create picoclaw-gemma2-tiny -f ~/ollama-modelfiles/Modelfile-gemma2-tiny
```

**Expected RAM:** ~900MB-1.1GB | **VRAM:** 0MB | **Speed:** ~3-5 tokens/sec on M1

> **⚠️ Note on Gemma 2:2b:** This model's base size is ~1.7GB. Even with aggressive parameter limits, RAM usage will hover around 900MB-1.1GB during inference. If you need strictly under 1GB, stick with the Qwen 0.5b/0.6b variants above.

##### Quick Comparison: Mac Mini M1 1GB Limit

| Model | RAM Usage | VRAM | Speed | Best For |
|-------|-----------|------|-------|----------|
| `qwen2.5:0.5b` | ~500MB | 0MB | ~5-8 t/s | General Q&A |
| `qwen3:0.6b` | ~550MB | 0MB | ~6-10 t/s | Better reasoning |
| `qwen2.5-coder:0.5b` | ~500MB | 0MB | ~5-8 t/s | Simple coding |
| `gemma2:2b` | ~900MB-1.1GB | 0MB | ~3-5 t/s | Best quality in this tier |

##### picoclaw-agents config for all 4:

```json
{
  "model_list": [
    {
      "model_name": "picoclaw-qwen25-min",
      "model": "picoclaw-qwen25-min",
      "api_base": "http://localhost:11434/v1",
      "api_key": "ollama"
    },
    {
      "model_name": "picoclaw-qwen3-tiny",
      "model": "picoclaw-qwen3-tiny",
      "api_base": "http://localhost:11434/v1",
      "api_key": "ollama"
    },
    {
      "model_name": "picoclaw-coder-tiny",
      "model": "picoclaw-coder-tiny",
      "api_base": "http://localhost:11434/v1",
      "api_key": "ollama"
    },
    {
      "model_name": "picoclaw-gemma2-tiny",
      "model": "picoclaw-gemma2-tiny",
      "api_base": "http://localhost:11434/v1",
      "api_key": "ollama"
    }
  ],
  "agents": {
    "defaults": {
      "model_name": "picoclaw-qwen3-tiny",
      "max_tokens": 512,
      "max_tool_iterations": 3
    }
  }
}
```

##### Verify RAM Usage on Mac Mini M1

```bash
# Check Ollama process memory
ps aux | grep ollama | grep -v grep | awk '{printf "PID %s: %s MB RSS\n", $2, $6/1024}'

# Or use Activity Monitor → Memory tab → filter "ollama"
# You should see all 4 models under ~1GB when idle
```

##### Using Your Models with picoclaw-agents

After adding models to `~/.picoclaw/config.json`, use them from any platform:

**CLI (One-shot message):**
```bash
./build/picoclaw-agents agent --model picoclaw-qwen25-min -m "hola"
./build/picoclaw-agents agent --model picoclaw-qwen3-tiny -m "Write a Python function"
./build/picoclaw-agents agent --model picoclaw-coder-tiny -m "def hello():"
./build/picoclaw-agents agent --model picoclaw-gemma2-tiny -m "Explain recursion"
```

**CLI (Interactive mode):**
```bash
./build/picoclaw-agents agent --model picoclaw-qwen25-min
# Type your messages interactively
```

**Telegram/Discord (via gateway):**
Once the gateway is running (`./build/picoclaw-agents gateway`), switch models in conversation:
```
/model picoclaw-qwen25-min
hola
```

**WebUI (picoclaw-agents-launcher):**
1. Open http://localhost:18800/chat
2. Select `picoclaw-qwen25-min` from the model dropdown
3. Start chatting

#### Quick Reference: Modelfile Parameter Ranges

| Parameter | Absolute Min | Typical | Max |
|-----------|-------------|---------|-----|
| `num_thread` | 1 | 4-8 | CPU cores |
| `num_ctx` | 512 | 2048-4096 | Model max (32K-128K) |
| `num_gpu` | 0 (CPU-only) | 20-35 | All layers |
| `num_batch` | 32 | 128-512 | 2048 |
| `num_keep` | 0 | 4-8 | num_ctx / 2 |

#### Rebuilding After Changes

If you edit the Modelfile, rebuild the model:

```bash
# Edit the Modelfile
nano ~/ollama-modelfiles/Modelfile-gemma4-8gb

# Rebuild (overwrites the previous version)
ollama create picoclaw-gemma4-8gb -f ~/ollama-modelfiles/Modelfile-gemma4-8gb

# The old model is replaced — no need to delete first
ollama list
```

### Method C: Via Environment Variables (Server-Level)

#### macOS — Using launchd (Desktop App)

If you installed Ollama via the `.dmg`, it runs as a launchd service:

```bash
# Stop Ollama
launchctl stop com.ollama.ollama

# Edit the launchd plist (if it exists)
nano ~/Library/LaunchAgents/com.ollama.ollama.plist

# Or use environment variables inline when starting manually:
# First quit the menu bar app, then:
OLLAMA_NUM_PARALLEL=2 OLLAMA_KEEP_ALIVE=1h ollama serve
```

**Using terminal app instead of background service:**

```bash
# Quit the menu bar app first
# Then start with env vars:
OLLAMA_KEEP_ALIVE=1h OLLAMA_NUM_PARALLEL=2 ollama serve
```

#### Linux — Using systemd

```bash
# Edit the systemd service override
systemctl edit ollama.service

# Add these lines in the editor:
[Service]
Environment="OLLAMA_NUM_PARALLEL=2"
Environment="OLLAMA_KEEP_ALIVE=1h"
Environment="OLLAMA_HOST=0.0.0.0:11434"

# Reload and restart
systemctl daemon-reload
systemctl restart ollama.service

# Verify
systemctl status ollama.service
env | grep OLLAMA
```

#### Windows — Using System Environment Variables

```powershell
# Method 1: Set permanently via PowerShell (requires restart)
[System.Environment]::SetEnvironmentVariable("OLLAMA_KEEP_ALIVE", "1h", "User")
[System.Environment]::SetEnvironmentVariable("OLLAMA_NUM_PARALLEL", "2", "User")

# Method 2: Set for current session only
$env:OLLAMA_KEEP_ALIVE = "1h"
$env:OLLAMA_NUM_PARALLEL = "2"

# Restart Ollama service
Restart-Service Ollama

# Or if running manually, start with env vars:
$env:OLLAMA_KEEP_ALIVE = "1h"
ollama serve
```

#### Termux — Using Shell Config

```bash
# Add to ~/.bashrc or ~/.zshrc (persistent)
echo 'export OLLAMA_KEEP_ALIVE=30m' >> ~/.bashrc
echo 'export OLLAMA_NUM_PARALLEL=1' >> ~/.bashrc
source ~/.bashrc

# Or set for current session only
export OLLAMA_KEEP_ALIVE=30m
export OLLAMA_NUM_PARALLEL=1

# Start Ollama
ollama serve
```

### All Official Resource Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `num_thread` | int | Auto (CPU cores) | CPU threads for inference |
| `num_ctx` | int | 2048 | Context window size (tokens) |
| `num_gpu` | int | All layers | Number of model layers to offload to GPU |
| `num_batch` | int | 512 | Batch size for prompt processing |
| `num_keep` | int | 0 | Initial tokens to keep in context |
| `main_gpu` | int | 0 | Primary GPU ID (multi-GPU setups) |
| `use_mmap` | bool | true | Use memory mapping for model loading |
| `numa` | bool | false | Enable NUMA memory optimization |

### Platform-Specific Examples for picoclaw-agents

#### Windows — Limit GPU VRAM

```powershell
# Modelfile: limit to 20 layers on GPU (rest on CPU/RAM)
FROM gemma4:26b
PARAMETER num_gpu 20
PARAMETER num_ctx 4096
PARAMETER num_thread 8

# Build and run
ollama create gemma4-limited -f Modelfile
ollama run gemma4-limited

# Or via /set during interactive session
ollama run gemma4:26b
/set parameter num_gpu 20
/set parameter num_ctx 4096
```

**Check GPU usage in Windows:**
```powershell
# Task Manager → Performance tab → GPU → Dedicated GPU Memory
# Or via PowerShell:
Get-Counter '\GPU Process Memory(*)\Local Usage'
```

#### macOS (Apple Silicon) — Limit Unified Memory

```bash
# Modelfile: limit threads and context for an 8GB Mac
FROM qwen3:8b
PARAMETER num_thread 6
PARAMETER num_ctx 2048
PARAMETER num_gpu 30

# Build and run
ollama create qwen3-lite -f Modelfile
ollama run qwen3-lite
```

**Check memory usage on macOS:**
```bash
# Monitor Ollama process memory
ps aux | grep ollama | awk '{print $6/1024 " MB", $11}'

# Or use Activity Monitor → Memory tab → filter "ollama"
```

#### Linux (NVIDIA GPU) — Limit GPU Layers + RAM

```bash
# Modelfile: partial GPU offload for limited VRAM
FROM llama3.1:70b
PARAMETER num_gpu 35
PARAMETER num_ctx 4096
PARAMETER num_batch 256
PARAMETER num_thread 8

ollama create llama70b-limited -f Modelfile
ollama run llama70b-limited
```

**Check GPU memory on Linux:**
```bash
# NVIDIA GPU memory
nvidia-smi --query-gpu=memory.used,memory.total --format=csv

# Or watch in real-time
watch -n 1 nvidia-smi
```

#### Termux (Android) — CPU-Only, Minimal Footprint

```bash
# Modelfile: ultra-lightweight for mobile
FROM qwen3:0.6b
PARAMETER num_thread 4
PARAMETER num_ctx 1024
PARAMETER num_batch 128
PARAMETER num_gpu 0

ollama create qwen-mobile -f Modelfile
ollama run qwen-mobile
```

**Check memory on Termux:**
```bash
# Process memory
ps -o pid,rss,comm | grep ollama | awk '{print $1, $2/1024 " MB", $3}'

# Or use htop
pkg install htop && htop
```

### Quick Reference: Parameter Values by Hardware

| Hardware | `num_thread` | `num_ctx` | `num_gpu` | `num_batch` |
|----------|-------------|-----------|-----------|-------------|
| Termux (Android, 4GB) | 4 | 1024 | 0 | 128 |
| Laptop (8GB RAM, no GPU) | 6 | 2048 | 0 | 256 |
| Mac M1 (8GB) | 6 | 2048 | 30 | 256 |
| Desktop (16GB + 8GB VRAM) | 8 | 4096 | 35 | 512 |
| Workstation (32GB + 24GB VRAM) | 12 | 8192 | auto | 512 |

> **Tip:** `num_gpu = 0` forces CPU-only mode. `num_gpu = auto` (or omitted) lets Ollama decide based on available VRAM.

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
