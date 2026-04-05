# Telegram Connection Fix - Cloudflare WARP Proxy

**Last Updated:** March 31, 2026  
**Status:** ✅ Production Ready  
**Severity:** CRITICAL - Complete Telegram bot unresponsiveness  

---

## Overview

This guide resolves **complete Telegram bot unresponsiveness** caused by ISP or network-level blocking of Telegram API servers. The solution uses **Cloudflare WARP** (WireGuard-based tunnel) with a local **SOCKS5 proxy** to bypass network restrictions.

**Problem:** Telegram bot shows connection timeout errors when connecting to Telegram API servers.

**Solution:** Route Telegram traffic through encrypted Cloudflare WARP tunnel with local SOCKS5 proxy.

---

## Symptoms

### Error Logs

```
ERROR Execution error getUpdates: 
  request call: fasthttp do request: error when dialing 149.154.166.110:443: 
  dialing to the given TCP address timed out

ERROR Error on get me: 
  telego: getMe: internal execution: request call: fasthttp do request: 
  error when dialing 149.154.166.110:443: 
  dialing to the given TCP address timed out
```

### Affected Components

| Component | Status | Notes |
|-----------|--------|-------|
| **Telegram Bot** | ❌ FAILED | Connection timeout to `149.154.166.110:443` |
| **Discord Bot** | ✅ Working | Same gateway, different library |
| **Web UI** | ✅ Working | Model configuration functional |
| **CLI** | ✅ Working | Commands functional |

### Root Cause

**Network/ISP Blocking** - System cannot connect to Telegram API servers:

- **Server IP:** `149.154.166.110` (Telegram API)
- **Port:** `443` (HTTPS)
- **Error:** `dialing to the given TCP address timed out`

**Evidence:**
1. Discord works → Gateway is functional
2. Web UI works → Model configuration is correct
3. CLI works → Configuration is correct
4. Telegram fails → **Specific network blocking to Telegram API**

---

## Solution Architecture

### Overview

```
┌─────────────────────────────────────────────────────────────┐
│                    PicoClaw-Agents                          │
│                                                             │
│  Telegram Bot ──→ SOCKS5 Proxy ──→ WARP Tunnel ──→ Internet │
│                     (socat)         (Cloudflare)             │
│                  127.0.0.1:40000   Encrypted                │
└─────────────────────────────────────────────────────────────┘
```

### Benefits

- ✅ **Bypasses network blocking** - Traffic routed through Cloudflare
- ✅ **Encrypted** - WireGuard encryption
- ✅ **Non-disruptive** - Doesn't affect other services
- ✅ **Persistent** - Survives reboots (with systemd)
- ✅ **Local proxy** - Easy to configure in applications

---

## Installation Guide

### Platform-Specific Instructions

Select your operating system:

- [**Linux (Ubuntu/Debian)**](#linux-ubuntudebian)
- [**macOS**](#macos)
- [**Windows**](#windows)

---

### Linux (Ubuntu/Debian)

#### Step 1: Install Cloudflare WARP Client

```bash
# Update package list
sudo apt update

# Install prerequisites
sudo apt install gnupg lsb-release -y

# Add Cloudflare GPG key
curl -fsSL https://pkg.cloudflareclient.com/pubkey.gpg | \
  sudo gpg --yes --dearmor --output /usr/share/keyrings/cloudflare-warp-archive-keyring.gpg

# Add Cloudflare repository
echo "deb [signed-by=/usr/share/keyrings/cloudflare-warp-archive-keyring.gpg] \
  https://pkg.cloudflareclient.com/ $(lsb_release -cs) main" | \
  sudo tee /etc/apt/sources.list.d/cloudflare-client.list

# Install WARP client
sudo apt update && sudo apt install cloudflare-warp -y
```

#### Step 2: Register Device with Cloudflare

```bash
# Register your device (outputs "Success")
warp-cli registration new
```

#### Step 3: Configure Proxy Mode

```bash
# Set to proxy mode (critical for server safety)
warp-cli mode proxy

# This creates a local SOCKS5 proxy on 127.0.0.1:40000
```

#### Step 4: Connect to WARP

```bash
# Enable the tunnel
warp-cli connect

# Verify connection status
warp-cli status

# Expected output: "Status: Connected"
```

#### Step 5: Install and Configure socat

```bash
# Install socat
sudo apt install socat -y

# Create SOCKS5 proxy bridge
# This forwards external:40001 → internal:127.0.0.1:40000
socat TCP-LISTEN:40001,fork,bind=0.0.0.0 TCP:127.0.0.1:40000 &

# Make it persistent (add to /etc/rc.local or create systemd service)
```

---

### macOS

#### Step 1: Install Cloudflare WARP

**Option A: Using Homebrew (Recommended)**

```bash
# Install Homebrew if not installed
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# Install Cloudflare WARP
brew install --cask cloudflare-warp
```

**Option B: Direct Download**

1. Download from: https://1.1.1.1/
2. Open the `.dmg` file
3. Drag "Cloudflare WARP" to Applications
4. Open from Applications folder

#### Step 2: Register Device with Cloudflare

```bash
# Open Terminal and register
warp-cli registration new

# Expected output: "Success"
```

#### Step 3: Configure Proxy Mode

```bash
# Set to proxy mode
warp-cli mode proxy

# This creates a local SOCKS5 proxy on 127.0.0.1:40000
```

#### Step 4: Connect to WARP

```bash
# Enable the tunnel
warp-cli connect

# Verify connection status
warp-cli status

# Expected output: "Status: Connected"
```

#### Step 5: Install and Configure socat

```bash
# Install socat using Homebrew
brew install socat

# Create SOCKS5 proxy bridge
socat TCP-LISTEN:40001,fork,bind=0.0.0.0 TCP:127.0.0.1:40000 &

# Make it persistent (using LaunchAgents)
# See Persistence section below
```

#### macOS Persistence (LaunchAgents)

```bash
# Create LaunchAgent for socat
cat > ~/Library/LaunchAgents/com.user.socat-proxy.plist << 'EOF'
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.user.socat-proxy</string>
    <key>ProgramArguments</key>
    <array>
        <string>/usr/local/bin/socat</string>
        <string>TCP-LISTEN:40001,fork,bind=0.0.0.0</string>
        <string>TCP:127.0.0.1:40000</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
</dict>
</plist>
EOF

# Load the LaunchAgent
launchctl load ~/Library/LaunchAgents/com.user.socat-proxy.plist

# Verify it's running
launchctl list | grep socat
```

---

### Windows

#### Step 1: Install Cloudflare WARP

**Option A: Using Winget (Recommended)**

```powershell
# Install Cloudflare WARP using winget
winget install Cloudflare.CloudflareWARP

# Or using Chocolatey
choco install cloudflare-warp
```

**Option B: Direct Download**

1. Download from: https://1.1.1.1/
2. Run the installer (`Cloudflare_WARP_setup.exe`)
3. Follow installation wizard
4. Launch Cloudflare WARP from Start menu

#### Step 2: Register Device with Cloudflare

```powershell
# Open PowerShell as Administrator
cd "C:\Program Files\Cloudflare\Cloudflare WARP"

# Register your device
.\warp-cli.exe registration new

# Expected output: "Success"
```

#### Step 3: Configure Proxy Mode

```powershell
# Set to proxy mode
.\warp-cli.exe mode proxy

# This creates a local SOCKS5 proxy on 127.0.0.1:40000
```

#### Step 4: Connect to WARP

```powershell
# Enable the tunnel
.\warp-cli.exe connect

# Verify connection status
.\warp-cli.exe status

# Expected output: "Status: Connected"
```

#### Step 5: Install and Configure socat for Windows

**Option A: Using WSL (Recommended)**

```bash
# Install WSL if not installed
wsl --install

# In WSL terminal, install socat
sudo apt update && sudo apt install socat -y

# Create SOCKS5 proxy bridge
socat TCP-LISTEN:40001,fork,bind=0.0.0.0 TCP:127.0.0.1:40000 &
```

**Option B: Using Cygwin**

```bash
# Install Cygwin with socat package
# https://www.cygwin.com/

# In Cygwin terminal
socat TCP-LISTEN:40001,fork,bind=0.0.0.0 TCP:127.0.0.1:40000 &
```

**Option C: Using Alternative (netsh)**

```powershell
# Create port proxy using netsh (built into Windows)
netsh interface portproxy add v4tov4 listenport=40001 listenaddress=0.0.0.0 connectport=40000 connectaddress=127.0.0.0

# Verify proxy
netsh interface portproxy show all

# Remove proxy (if needed)
netsh interface portproxy delete v4tov4 listenport=40001 listenaddress=0.0.0.0
```

#### Windows Persistence

**For WARP:**
- WARP for Windows auto-starts by default
- Verify in Task Manager → Startup tab

**For socat (WSL):**

```bash
# Add to WSL bashrc
echo "socat TCP-LISTEN:40001,fork,bind=0.0.0.0 TCP:127.0.0.1:40000 &" >> ~/.bashrc

# Or create systemd service in WSL
```

**For netsh proxy:**

```powershell
# Create scheduled task for port proxy
$action = New-ScheduledTaskAction -Execute "netsh" -Argument "interface portproxy add v4tov4 listenport=40001 listenaddress=0.0.0.0 connectport=40000 connectaddress=127.0.0.0"
$trigger = New-ScheduledTaskTrigger -AtStartup
Register-ScheduledTask -TaskName "WARP Proxy Bridge" -Action $action -Trigger $trigger -RunLevel Highest
```

---

## Configuration

### PicoClaw-Agents Configuration

Edit `~/.picoclaw/config.json`:

**Option A: Local Proxy (localhost)**

```json
{
  "channels": {
    "telegram": {
      "enabled": true,
      "token": "YOUR_TELEGRAM_BOT_TOKEN",
      "proxy": "socks5://127.0.0.1:40000",
      "allow_from": [
        "YOUR_USER_ID"
      ]
    }
  }
}
```

**Option B: Remote Proxy (another IP on network)**

```json
{
  "channels": {
    "telegram": {
      "enabled": true,
      "token": "YOUR_TELEGRAM_BOT_TOKEN",
      "proxy": "socks5://192.168.0.100:40001",
      "allow_from": [
        "YOUR_USER_ID"
      ]
    }
  }
}
```

**Configuration Notes:**

| Setting | Value | Description |
|---------|-------|-------------|
| **Host** | `127.0.0.1` | Local proxy (WARP + socat on same server) |
| **Port** | `40000` | Default WARP proxy port |
| **Protocol** | `SOCKS5` | Proxy protocol |
| **Full URL** | `socks5://127.0.0.1:40000` | Complete proxy URL |

**Important:**
- **Local proxy:** Use `127.0.0.1:40000` when WARP + socat run on the same server
- **Remote proxy:** Use `192.168.x.x:40001` when proxy runs on another machine in the network
- **Token:** Replace `YOUR_TELEGRAM_BOT_TOKEN` with your actual token from @BotFather
- **User ID:** Replace `YOUR_USER_ID` with your Telegram numeric user ID

### Alternative: Environment Variables

```bash
# Add to ~/.bashrc or ~/.zshrc
export ALL_PROXY=socks5://127.0.0.1:40000
export HTTPS_PROXY=socks5://127.0.0.1:40000

# Or temporary for one session
export ALL_PROXY=socks5://127.0.0.1:40000
./build/picoclaw-agents-launcher
```

---

## Verification

### Step 1: Test SOCKS5 Proxy

```bash
# Test proxy connectivity
curl -v -x socks5://127.0.0.1:40000 https://www.google.com

# Expected: Successful connection through proxy
```

### Step 2: Test Telegram API

```bash
# Replace YOUR_BOT_TOKEN with actual token
curl -v -x socks5://127.0.0.1:40000 \
  "https://api.telegram.org/botYOUR_BOT_TOKEN/getMe"

# Expected output:
# {"ok":true,"result":{"id":123456789,"is_bot":true,"first_name":"...","username":"..."}}
```

### Step 3: Verify Application Logs

```bash
# Restart application
pkill -f picoclaw
./build/picoclaw-agents-launcher &

# Wait 15 seconds
sleep 15

# Check Telegram logs
tail -100 ~/.picoclaw/logs/launcher.log | grep -i telegram

# Expected output:
# [INFO] telegram: Telegram bot connected {username=your_bot_username}
```

### Step 4: Test Bot in Telegram

1. Open Telegram
2. Find your bot
3. Send `/start`
4. Bot should respond immediately

---

## Performance Comparison

| Metric | Before (Blocked) | After (WARP) |
|--------|-----------------|--------------|
| **Telegram API** | ❌ Timeout | ✅ < 500ms |
| **Bot Response** | ❌ No response | ✅ Instant |
| **Connection** | ❌ Blocked | ✅ Routed via Cloudflare |
| **Other Services** | ✅ Working | ✅ Still working |

---

## Troubleshooting

### Problem: WARP Won't Connect

```bash
# Check registration status
warp-cli registration status

# If not registered, register again
warp-cli registration new

# Check firewall
sudo ufw status
sudo ufw allow out to any port 2408 proto udp
```

### Problem: Proxy Not Working

```bash
# Check if WARP is running
warp-cli status

# Check if port 40000 is listening
netstat -tlnp | grep 40000

# Restart WARP
warp-cli disconnect
warp-cli connect
```

### Problem: socat Not Forwarding

```bash
# Check if socat is running
ps aux | grep socat

# Restart socat
pkill socat
socat TCP-LISTEN:40001,fork,bind=0.0.0.0 TCP:127.0.0.1:40000 &

# Verify port is listening
netstat -tlnp | grep 40001
```

### Problem: Still Getting Timeouts

```bash
# Test direct connection (should fail)
curl -v https://api.telegram.org/botYOUR_TOKEN/getMe

# Test through proxy (should work)
curl -v -x socks5://127.0.0.1:40000 https://api.telegram.org/botYOUR_TOKEN/getMe

# If proxy test fails, check WARP status
warp-cli status
warp-cli disconnect
warp-cli connect
```

---

## Security Considerations

### What This Does

- ✅ **Encrypts traffic** via WireGuard
- ✅ **Routes through Cloudflare** network
- ✅ **Bypasses ISP blocking**
- ✅ **Local proxy** (only accessible from localhost)

### What This Doesn't Do

- ❌ **Doesn't anonymize** your identity
- ❌ **Doesn't hide** bot token
- ❌ **Doesn't bypass** Telegram API limits
- ❌ **Doesn't prevent** rate limits

### Best Practices

1. **Keep proxy local** - Don't expose to public network
2. **Use firewall** - Block external access to port 40000
3. **Monitor logs** - Watch for unusual activity
4. **Update regularly** - Keep WARP client updated

```bash
# Firewall rule to block external access
sudo ufw deny from any to any port 40000
sudo ufw allow from 127.0.0.1 to any port 40000
```

---

## Persistence (Auto-start on Boot)

### Option 1: systemd Service for WARP

```bash
# Create systemd service
sudo nano /etc/systemd/system/warp-client.service

# Add content:
[Unit]
Description=Cloudflare WARP Client
After=network.target

[Service]
Type=simple
ExecStart=/opt/cloudflare-warp/warp-client
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target

# Enable and start
sudo systemctl enable warp-client
sudo systemctl start warp-client
```

### Option 2: systemd Service for socat

```bash
# Create systemd service
sudo nano /etc/systemd/system/socks5-proxy.service

# Add content:
[Unit]
Description=SOCKS5 Proxy Bridge for WARP
After=warp-client.service
Requires=warp-client.service

[Service]
Type=simple
ExecStart=/usr/bin/socat TCP-LISTEN:40001,fork,bind=0.0.0.0 TCP:127.0.0.1:40000
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target

# Enable and start
sudo systemctl enable socks5-proxy
sudo systemctl start socks5-proxy
```

---

## References

### Official Documentation

- **Cloudflare WARP:** https://developers.cloudflare.com/cloudflare-one/connections/connect-devices/warp/
- **Telegram Bot API:** https://core.telegram.org/bots/api
- **Telegram Proxy:** https://core.telegram.org/api/proxy

### Related Issues

- **Telegram API Blocking:** Common in some regions/ISPs
- **Cloudflare WARP:** Free tier available, paid for teams
- **SOCKS5 Proxy:** Standard protocol for proxied connections

---

## Summary

### Problem
- Telegram API servers blocked/timed out
- Bot couldn't connect to Telegram
- Other services (Discord, Web UI, CLI) worked fine

### Solution
- Install Cloudflare WARP (WireGuard tunnel)
- Configure in proxy mode (127.0.0.1:40000)
- Use socat to bridge ports if needed
- Configure Telegram channel to use SOCKS5 proxy

### Result
- ✅ Telegram bot connected successfully
- ✅ All commands functional
- ✅ No impact on other services
- ✅ Encrypted, secure connection

---

**Author:** @comgunner  
**Repository:** https://github.com/comgunner/picoclaw-agents  
**Last Updated:** March 31, 2026
