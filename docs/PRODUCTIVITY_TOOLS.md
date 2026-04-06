# Productivity Tools

**Date:** 2026-04-05
**Version:** v1.4.0
**Language:** Go (stdlib only, zero external dependencies)

---

## Overview

PicoClaw-Agents includes a suite of built-in productivity utilities accessible via the `picoclaw-agents util` command. These tools help you benchmark performance, manage processes, validate architecture, and audit documentation — all without leaving the CLI.

All productivity tools are implemented in pure Go using only the standard library. No external dependencies are required.

---

## Available Commands

### 1. `bench` — Benchmark Performance

Measure the startup time and memory usage of any command.

```bash
# Benchmark the picoclaw-agents binary itself
picoclaw-agents util bench --self

# Benchmark an external command
picoclaw-agents util bench ./build/picoclaw-agents -- --help

# Benchmark with arguments
picoclaw-agents util bench /usr/bin/python3 -- "-c print('hello')"
```

**Output:**
```
Startup time: 12.4ms
Peak RSS: 24 MB
Alloc: 5 MB
```

**Use cases:**
- Verify the binary stays under the 45 MB RAM target
- Compare startup times across different builds
- Profile external tools during development

---

### 2. `reaper` — Clean Orphan Processes

Find and terminate orphaned `picoclaw-agents` processes whose parent has exited (PPID = 1).

```bash
# List orphans without killing them
picoclaw-agents util reaper --dry-run

# Find and kill orphans
picoclaw-agents util reaper
```

**Output:**
```
Found 3 orphan(s):
  PID 12345: /usr/local/bin/picoclaw-agents gateway
  PID 12389: /usr/local/bin/picoclaw-agents agent -m hello
  PID 12401: picoclaw-agents-launcher

✅ Killed 3 orphan(s).
```

**Use cases:**
- Clean up after crashed agents
- Free memory on resource-constrained devices (Termux, Raspberry Pi)
- Routine maintenance before restarting services

**Platform support:** Linux, macOS, FreeBSD. Not supported on Windows.

---

### 3. `arch-lint` — Validate Import Boundaries

Check for forbidden import patterns between packages. Detects architectural violations like `pkg/agent` importing `pkg/channels` or `pkg/tools` importing `pkg/agent`.

```bash
# Check the current project
picoclaw-agents util arch-lint .

# Check a specific directory
picoclaw-agents util arch-lint /path/to/project
```

**Default forbidden rules:**

| From package | Must not import | Reason |
|-------------|-----------------|--------|
| `pkg/agent` | `pkg/channels` | Agents must not depend on specific channel implementations |
| `pkg/tools` | `pkg/agent` | Tools must not depend on the agent core (prevents cycles) |
| `pkg/mcp` | `pkg/agent` | MCP client must be agent-independent |
| `pkg/mcp` | `pkg/providers` | MCP client must be provider-independent |

**Output (clean):**
```
✅ No import violations found.
```

**Output (violation):**
```
Found 1 violation(s):
  pkg/agent → must not import pkg/channels
    File: pkg/agent/loop.go
```

**Use cases:**
- Pre-commit architectural validation
- Refactoring safety checks
- Onboarding new contributors

---

### 4. `md-audit` — Audit Documentation Links

Scan Markdown files for broken internal links. Ensures that documentation references are valid.

```bash
# Scan the docs/ directory (default)
picoclaw-agents util md-audit

# Scan a specific directory
picoclaw-agents util md-audit /path/to/markdown-files
```

**Output (clean):**
```
✅ No broken internal links found.
```

**Output (issues):**
```
Found 2 link issue(s):
  README.md:45 — docs/MISSING_FILE.md
  docs/SETUP.md:12 — ../old-guide.md
```

**Behavior:**
- Checks **relative internal links** only
- Skips `http://` and `https://` external links
- Skips anchor links (`#section-name`)
- Validates that the target file exists relative to the source file's directory

**Use cases:**
- Documentation quality assurance
- Pre-release link validation
- Catching typos in cross-references

---

## Agent Integration

All four tools are registered as **native agent tools**. The AI agent can discover and invoke them automatically during conversation:

| Agent Tool Name | CLI Equivalent | Description |
|-----------------|----------------|-------------|
| `bench` | `util bench` | Benchmark command startup time and memory |
| `reaper` | `util reaper` | Find and kill orphaned processes |
| `arch_lint` | `util arch-lint` | Check forbidden import patterns |
| `md_audit` | `util md-audit` | Scan for broken markdown links |

The agent will use these tools when appropriate — for example, running `md_audit` when asked to verify documentation, or `arch_lint` when checking code quality.

---

## Technical Details

### Zero External Dependencies

All productivity tools use only the Go standard library:

| Tool | Packages Used |
|------|--------------|
| `bench` | `runtime`, `os/exec`, `io`, `time`, `fmt` |
| `reaper` | `os`, `os/exec`, `runtime`, `strconv`, `strings`, `fmt` |
| `arch-lint` | `go/parser`, `go/token`, `os`, `path/filepath`, `strings`, `fmt` |
| `md-audit` | `os`, `path/filepath`, `regexp`, `strings`, `fmt` |

**Result:** `go.mod` and `go.sum` are unchanged. Cross-platform builds work without installing additional packages.

### Cross-Platform Support

| Platform | `bench` | `reaper` | `arch-lint` | `md-audit` |
|----------|---------|----------|-------------|------------|
| Linux (amd64/arm64) | ✅ | ✅ | ✅ | ✅ |
| macOS (arm64/amd64) | ✅ | ✅ | ✅ | ✅ |
| Windows (amd64) | ✅ | ❌ (skips) | ✅ | ✅ |
| FreeBSD (amd64) | ✅ | ✅ | ✅ | ✅ |

---

*PRODUCTIVITY_TOOLS.md — 2026-04-05*
