# PicoClaw Installation on Ubuntu Server

This guide details the steps required to compile and install PicoClaw in an Ubuntu Server environment.

## Prerequisites

Make sure you have the following components installed before starting:

- **Git:** To clone the repository.
- **Go:** PicoClaw is written in Go, so you need the Go compiler installed and configured on your system.
- **Make:** Tool used to manage the build process.

You can install the basic OS dependencies with:

```bash
sudo apt update
sudo apt install -y build-essential git wget tar
```

### Installing Go (Golang)

PicoClaw requires Go 1.25.7 or higher. The recommended way to install it directly on Ubuntu Server is:

```bash
# Download and install Go 1.25.7
wget https://go.dev/dl/go1.25.7.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.25.7.linux-amd64.tar.gz
rm go1.25.7.linux-amd64.tar.gz

# Add to your PATH permanently
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile
source ~/.profile

# Verify installation
go version
```

## Installation Steps

1. **Clone the repository and enter the directory:**
   (If you are installing globally in `/opt`, make sure to grant your user permissions).

   ```bash
   cd /opt
   sudo git clone https://github.com/comgunner/picoclaw.git
   sudo chown -R $USER:$USER /opt/picoclaw
   cd picoclaw
   ```

2. **Download and install Go dependencies:**
   This will download all necessary libraries specified in the `go.mod` file.

   ```bash
   make deps
   ```

3. **Compile the project:**
   We will generate the binary directly in the `/opt/picoclaw` root folder.

   ```bash
   go build -v -o picoclaw ./cmd/picoclaw
   chmod +x picoclaw
   ls -lh picoclaw
   ```
Once the installation process is complete, you can verify that PicoClaw is working correctly by running:

```bash
./picoclaw --help
```

## Server Configuration

The main configuration file (`config.json`) must be placed in the home directory of the user executing PicoClaw.

1. **Location on Ubuntu Server:**
   ```bash
   ~/.picoclaw/config.json
   ```
   *(This is equivalent to `/home/YOUR_USER/.picoclaw/config.json` or `/root/.picoclaw/config.json` if running as root).*

2. **Copying your config from Mac to Ubuntu:**
   Since you have already configured your bots (like Telegram) and API Keys (Deepseek) locally on your Mac, you can simply transfer your current file to the server using `scp`:

   ```bash
   # Execute this from your Mac:
   scp ~/.picoclaw/config.json user@your_server_ip:~/.picoclaw/config.json
   ```

   *Note: You may need to create the hidden folder on the server before copying it (`mkdir -p ~/.picoclaw`).*

3. **Or use the `onboard` command to auto-generate it:**
   If you prefer not to copy the file, you can let PicoClaw generate a new one with default values (including the multi-agent `deepseek-chat` structure configured in the source code). Just run:

   ```bash
   ./picoclaw onboard
   ```
   This will create the `~/.picoclaw/config.json` file with the entire structure ready. From there, you just need to edit the file using `nano` or `vim` to paste:
   - Your Deepseek `api_key` in the `model_list` section.
   - Your Telegram `token` in the `channels > telegram` section.
   - Your User ID in the `allow_from` list.

## Run as a Service (Systemd)

For PicoClaw to run in the background and restart automatically if the server shuts down or crashes, you can configure it as a Systemd service.

1. **Create the service file:**
   This command will automatically create the file using your current user (`$USER`):

   ```bash
   cat <<EOF | sudo tee /etc/systemd/system/picoclaw.service
   [Unit]
   Description=PicoClaw Gateway Service
   After=network.target

   [Service]
   # Directory where the code lives
   WorkingDirectory=/opt/picoclaw

   # Full path to the compiled binary
   ExecStart=/opt/picoclaw/picoclaw gateway

   # The user executing it
   User=$USER
   Group=$(id -gn)

   # Restart automatically on failure
   Restart=always
   RestartSec=5

   [Install]
   WantedBy=multi-user.target
   EOF
   ```

2. **Enable and start the service:**
   ```bash
   sudo systemctl daemon-reload
   sudo systemctl enable picoclaw.service
   sudo systemctl start picoclaw.service
   ```

4. **View status and logs:**
   ```bash
   sudo systemctl status picoclaw.service
   journalctl -u picoclaw.service -f
   ```
