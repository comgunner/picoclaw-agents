# PicoClaw-Agents Installation on Ubuntu Server

This guide details the steps required to install PicoClaw-Agents in an Ubuntu Server environment.

## Installation Methods

**IMPORTANT:** There are two installation methods:

1. **From GitHub Releases** (Recommended if releases are available)
2. **From Source Code** (If no releases or for development)

### Check if Releases are Available

Before starting, verify if releases exist:
```bash
curl -s https://api.github.com/repos/comgunner/picoclaw-agents/releases/latest | grep '"tag_name"'
# If it returns something like "v3.7.0", releases are available
# If it returns "Not Found", use Method 2 (from source)
```

---

## Method 1: From GitHub Releases (If available)

**Prerequisites:**
- `wget` or `curl` to download

```bash
sudo apt update
sudo apt install -y wget curl
```

### Installation Steps

1. **Create installation directory:**

   ```bash
   sudo mkdir -p /opt/picoclaw-agents
   sudo chown -R $USER:$USER /opt/picoclaw-agents
   cd /opt/picoclaw-agents
   ```

2. **Download the latest release:**

   ```bash
   # Download latest tar.gz
   wget https://github.com/comgunner/picoclaw-agents/releases/latest/download/picoclaw-agents_Linux_arm64.tar.gz
   
   # Or for specific version (e.g., v3.7.0)
   # wget https://github.com/comgunner/picoclaw-agents/releases/download/v3.7.0/picoclaw-agents_Linux_arm64.tar.gz
   ```

3. **Extract and setup:**

   ```bash
   # Extract tar.gz
   tar -xzf picoclaw-agents_Linux_arm64.tar.gz
   
   # Binary extracts as 'picoclaw-agents'
   chmod +x picoclaw-agents
   
   # Create symlink for global command
   sudo ln -sf /opt/picoclaw-agents/picoclaw-agents /usr/local/bin/picoclaw-agents
   ```

4. **Verify installation:**

   ```bash
   picoclaw-agents --version
   # or
   ./picoclaw-agents --help
   ```

---

## Method 2: From Source Code (Development or if no releases)

**Prerequisites:**
- Go 1.25.7 or higher
- Git
- Make

### 1. Install Go

```bash
# Download and install Go 1.25.7+
wget https://go.dev/dl/go1.25.7.linux-arm64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.25.7.linux-arm64.tar.gz
rm go1.25.7.linux-arm64.tar.gz

# Add to PATH permanently
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile
source ~/.profile

# Verify installation
go version
# Should show: go version go1.25.7 linux/arm64
```

### 2. Clone and Build

```bash
# Clone repository
cd /opt
git clone https://github.com/comgunner/picoclaw-agents.git
sudo chown -R $USER:$USER /opt/picoclaw-agents
cd picoclaw-agents

# Download dependencies
make deps

# Build binary
make build
# Binary is generated at: build/picoclaw-agents
```

### 3. Install Binary

```bash
# Copy binary to installation directory
sudo mkdir -p /opt/picoclaw-agents
cp build/picoclaw-agents /opt/picoclaw-agents/
chmod +x /opt/picoclaw-agents/picoclaw-agents

# Create symlink for global command
sudo ln -sf /opt/picoclaw-agents/picoclaw-agents /usr/local/bin/picoclaw-agents
```

### 4. Verify Installation

```bash
picoclaw-agents --version
# or
./picoclaw-agents --help
```

---

## Server Configuration

The main configuration file (`config.json`) must be placed in the home directory of the user executing PicoClaw-Agents.

### 1. Location on Ubuntu Server

```bash
~/.picoclaw/config.json
```
*(This is equivalent to `/home/YOUR_USER/.picoclaw/config.json` or `/root/.picoclaw/config.json` if running as root).*

### 2. Copying your config from Mac/Local

If you have already configured your bots (like Telegram) and API Keys locally, you can transfer your current file to the server using `scp`:

```bash
# Execute this from your Mac/Local:
scp ~/.picoclaw/config.json user@your_server_ip:~/.picoclaw/config.json
```

*Note: You may need to create the hidden folder on the server before copying it (`mkdir -p ~/.picoclaw`).*

### 3. Or use the `onboard` command to auto-generate

If you prefer not to copy the file, you can let PicoClaw-Agents generate a new one with default values:

```bash
picoclaw-agents onboard
```

This will create the `~/.picoclaw/config.json` file with the entire structure ready. Then edit the file with `nano` or `vim` to add:
- Your Deepseek/OpenAI/Anthropic `api_key` in `model_list`
- Your Telegram `token` in `channels > telegram`
- Your User ID in `allow_from`

---

## Run as a Service (Systemd)

For PicoClaw-Agents to run in the background and restart automatically:

### 1. Create the service file

```bash
cat <<EOF | sudo tee /etc/systemd/system/picoclaw-agents.service
[Unit]
Description=PicoClaw-Agents Gateway Service
After=network.target

[Service]
# Directory where the binary lives
WorkingDirectory=/opt/picoclaw-agents

# Full path to the binary
ExecStart=/opt/picoclaw-agents/picoclaw-agents gateway

# The user executing it
User=$USER
Group=$(id -gn)

# Restart automatically on failure
Restart=always
RestartSec=5

# Logs to journal
StandardOutput=journal
StandardError=journal
SyslogIdentifier=picoclaw-agents

[Install]
WantedBy=multi-user.target
EOF
```

### 2. Enable and start

```bash
sudo systemctl daemon-reload
sudo systemctl enable picoclaw-agents.service
sudo systemctl start picoclaw-agents.service
```

### 3. View status and logs

```bash
# Check status
sudo systemctl status picoclaw-agents.service

# View logs in real-time
journalctl -u picoclaw-agents.service -f

# View logs from last 2 hours
journalctl -u picoclaw-agents.service --since "2 hours ago"
```

### 4. Useful commands

```bash
# Stop service
sudo systemctl stop picoclaw-agents.service

# Restart service
sudo systemctl restart picoclaw-agents.service

# View logs without follow
journalctl -u picoclaw-agents.service -n 50
```

---

## Update

### If using Releases (Method 1)

```bash
# 1. Stop service
sudo systemctl stop picoclaw-agents.service

# 2. Download new version
cd /opt/picoclaw-agents
wget https://github.com/comgunner/picoclaw-agents/releases/latest/download/picoclaw-agents_Linux_arm64.tar.gz

# 3. Extract (replaces old binary)
tar -xzf picoclaw-agents_Linux_arm64.tar.gz
chmod +x picoclaw-agents

# 4. Start service
sudo systemctl start picoclaw-agents.service

# 5. Verify version
picoclaw-agents --version
```

### If using Source Code (Method 2)

```bash
# 1. Stop service
sudo systemctl stop picoclaw-agents.service

# 2. Update code
cd /opt/picoclaw-agents
git pull

# 3. Rebuild
make build

# 4. Replace binary
cp build/picoclaw-agents /opt/picoclaw-agents/picoclaw-agents
chmod +x /opt/picoclaw-agents/picoclaw-agents

# 5. Start service
sudo systemctl start picoclaw-agents.service

# 6. Verify version
picoclaw-agents --version
```

---

## Available Commands

```bash
# One-shot query
picoclaw-agents agent -m "What's the weather?"

# Interactive mode
picoclaw-agents agent

# Gateway (Telegram, Discord bots, etc.)
picoclaw-agents gateway

# List available skills
picoclaw-agents skills list

# Install skill from ClawHub
picoclaw-agents skills install --slug skill-name

# Check version
picoclaw-agents --version

# Initial setup
picoclaw-agents onboard
```

---

## Troubleshooting

### Service won't start

```bash
# Check error logs
journalctl -u picoclaw-agents.service -n 100 --no-pager

# Verify binary exists
ls -lh /opt/picoclaw-agents/picoclaw-agents

# Check permissions
sudo chmod +x /opt/picoclaw-agents/picoclaw-agents
```

### Configuration error

```bash
# Validate config.json
picoclaw-agents onboard --dry-run

# Regenerate config
rm ~/.picoclaw/config.json
picoclaw-agents onboard
```

### Verify gateway is running

```bash
# Check process
ps aux | grep picoclaw-agents

# Check listening port (if using webhook)
sudo netstat -tulpn | grep picoclaw
```

### Error: "command not found: picoclaw-agents"

```bash
# Verify symlink
ls -lh /usr/local/bin/picoclaw-agents

# If it doesn't exist, create manually
sudo ln -sf /opt/picoclaw-agents/picoclaw-agents /usr/local/bin/picoclaw-agents

# Or use full path
/opt/picoclaw-agents/picoclaw-agents --version
```

---

## Useful Links

- **Releases:** https://github.com/comgunner/picoclaw-agents/releases
- **Documentation:** https://github.com/comgunner/picoclaw-agents/tree/main/docs
- **Issues:** https://github.com/comgunner/picoclaw-agents/issues
- **Discussions:** https://github.com/comgunner/picoclaw-agents/discussions

---

## Important Notes

### Binary Name

- **After extracting tar.gz:** The binary is named `picoclaw-agents` (no platform suffix)
- **In commands:** Always use `picoclaw-agents` (the symlink points to the correct binary)

### Server Architecture

- **ARM64 (Raspberry Pi, AWS Graviton):** Use `picoclaw-agents_Linux_arm64.tar.gz`
- **AMD64 (Intel/AMD traditional):** Use `picoclaw-agents_Linux_amd64.tar.gz`

To check your architecture:
```bash
uname -m
# aarch64 or arm64 = ARM
# x86_64 = AMD64
```

### Go Not Required (Method 1)

If you use the pre-compiled binary from Releases, you **DO NOT need to install Go**. Go is only necessary if you compile from source code (Method 2).
