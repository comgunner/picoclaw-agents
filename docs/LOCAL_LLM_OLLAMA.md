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

Create a `Modelfile` to build a custom model with permanent resource limits:

```Modelfile
FROM llama3.2:3b
PARAMETER num_thread 6
PARAMETER num_ctx 4096
PARAMETER num_gpu 25
PARAMETER num_batch 256
PARAMETER num_keep 4
SYSTEM You are a helpful assistant for PicoClaw-Agents.
```

Build and use it:

```bash
# Create the custom model
ollama create picoclaw-llama -f Modelfile

# Verify it exists
ollama list

# Run it
ollama run picoclaw-llama

# Connect picoclaw-agents to it
picoclaw-agents agent --model picoclaw-llama -m "Hello"
```

**Now configure picoclaw-agents** to use your custom model by editing `~/.picoclaw/config.json`:

```json
{
  "model_list": [
    {
      "model_name": "picoclaw-llama",
      "model": "picoclaw-llama",
      "api_base": "http://localhost:11434/v1",
      "api_key": "ollama"
    }
  ],
  "agents": {
    "defaults": {
      "model_name": "picoclaw-llama"
    }
  }
}
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
