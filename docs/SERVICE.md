# OS Service Management

Install PicoClaw-Agents as a system service that starts automatically on boot.

## Supported Platforms

| Platform | Service Type | Auto-start | Requires sudo |
|----------|-------------|------------|--------------|
| **Linux** | systemd (user) | Yes (on login) | No |
| **macOS** | launchd (LaunchAgent) | Yes (on login) | No |
| **Windows** | Scheduled Task (ONLOGON) | Yes (on login) | No |

> **Note:** User-level services start when you log in. For system-level services (start before login), use `sudo` with system-level configuration (not currently supported).

## Quick Start

### Install

```bash
# Install with defaults
picoclaw-agents service install

# Install on a custom port, listening on all interfaces
picoclaw-agents service install --port 9999 --public

# Preview what would be installed (dry run)
picoclaw-agents service install --dry-run
```

### Manage the Service

```bash
# Check status
picoclaw-agents service status

# View logs
picoclaw-agents service logs              # Last 50 lines
picoclaw-agents service logs -n 100       # Last 100 lines
picoclaw-agents service logs -f           # Follow in real-time

# Stop and start
picoclaw-agents service stop
picoclaw-agents service start
picoclaw-agents service restart

# Uninstall
picoclaw-agents service uninstall
```

## Install Options

| Flag | Default | Description |
|------|---------|-------------|
| `--port` | `18800` | Gateway port number |
| `--public` | `false` | Listen on all interfaces (0.0.0.0) instead of localhost |
| `--config` | `~/.picoclaw/config.json` | Path to configuration file |
| `--dry-run` | `false` | Show what would be done without making changes |

## Platform Details

### Linux (systemd)

**Unit file location:** `~/.config/systemd/user/picoclaw-agents.service`

```bash
# Manual systemd commands (if needed)
systemctl --user status picoclaw-agents.service
systemctl --user restart picoclaw-agents.service
journalctl --user -u picoclaw-agents.service -f
```

**Requirements:** systemd user session must be available. Check with:
```bash
systemctl --user is-system-running
```

### macOS (launchd)

**Plist location:** `~/Library/LaunchAgents/com.picoclaw.agents.gateway.plist`

**Log files:**
- stdout: `~/Library/Logs/picoclaw-agents/gateway.out.log`
- stderr: `~/Library/Logs/picoclaw-agents/gateway.err.log`

```bash
# Manual launchctl commands (if needed)
launchctl list | grep com.picoclaw.agents.gateway
launchctl bootout gui/$(id -u)/com.picoclaw.agents.gateway  # stop
launchctl bootstrap gui/$(id -u)/ ~/Library/LaunchAgents/com.picoclaw.agents.gateway.plist  # start
```

**Note:** launchd services only work in a GUI session. If connecting via SSH, the service will not auto-start.

### Windows (Scheduled Task)

**Task name:** `PicoClaw-Agents Gateway`

**Files created:**
- Wrapper script: `%APPDATA%\PicoClaw\gateway.cmd`
- PID file: `%APPDATA%\PicoClaw\gateway.pid`
- Logs: `%APPDATA%\PicoClaw\logs\gateway.log`

```powershell
# Manual schtasks commands (if needed)
schtasks /Query /TN "PicoClaw-Agents Gateway" /FO LIST
schtasks /End /TN "PicoClaw-Agents Gateway"    # stop
schtasks /Run /TN "PicoClaw-Agents Gateway"     # start
schtasks /Delete /TN "PicoClaw-Agents Gateway" /F  # uninstall
```

**Note:** The scheduled task triggers on login (ONLOGON). For startup before login, use `schtasks /SC ONSTART` manually.

## ⚠️ Important Warnings

### Prevent Duplicate Connections

**The service and a manually-started gateway MUST NOT run simultaneously.** If both are active:

- **Telegram:** Two connections to the same bot token → messages processed twice, API rate limits, potential bans
- **Discord:** Two bot instances → intent conflicts, duplicate messages
- **Port conflict:** Both processes try to bind to the same port

The service installer checks if the gateway port is already in use and will refuse to install if it detects a running gateway.

**If you see this error:**
```
Error: gateway is already running on port 18800.
```

**Solution:** Stop the manual gateway first:
```bash
# If running via CLI, press Ctrl+C
# Or find and kill the process:
picoclaw-agents gateway stop   # if available
# Or:
kill <PID>
```

### Service vs Manual Gateway

| | Service | Manual (`picoclaw-agents gateway`) |
|--|---------|-------------------------------------|
| Auto-start | ✅ On login | ❌ No |
| Survives logout | ✅ Yes (macOS/Linux) | ❌ No |
| Logs to file | ✅ Yes | ❌ stdout only |
| Restart on crash | ✅ Yes | ❌ No |
| Can run simultaneously | ❌ NO | ❌ NO |

**Rule:** Use ONE or the other, never both at the same time.

## Troubleshooting

### Service won't start

1. Check if the port is in use:
   ```bash
   lsof -i :18800    # macOS/Linux
   netstat -ano | findstr 18800  # Windows
   ```

2. Check the binary path hasn't changed:
   ```bash
   picoclaw-agents service status
   ```

3. View error logs:
   ```bash
   # Linux
   journalctl --user -u picoclaw-agents.service -n 50
   
   # macOS
   tail -50 ~/Library/Logs/picoclaw-agents/gateway.err.log
   
   # Windows
   type %APPDATA%\PicoClaw\logs\gateway.log
   ```

### Service installed but not running

```bash
# Try starting it
picoclaw-agents service start

# Check status
picoclaw-agents service status
```

### Moved or reinstalled PicoClaw

If you moved the binary to a different path after installing the service, the service will point to the old path. Uninstall and reinstall:

```bash
picoclaw-agents service uninstall
picoclaw-agents service install
```
