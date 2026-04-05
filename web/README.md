# Picoclaw Web — WebUI Launcher

This directory contains `picoclaw-agents-launcher`, the standalone web service for `picoclaw-agents`.
It provides a complete unified web interface acting as a dashboard, configuration center, and
interactive console (channel client) for the core `picoclaw-agents` engine.

## Architecture

The service is structured as a monorepo containing both the backend and frontend code:

- **`backend/`**: Go-based web server. Provides RESTful APIs, manages WebSocket connections for chat
  via `PicoChannel`, handles the lifecycle of the `picoclaw-agents` gateway process, and embeds
  compiled frontend assets into a single executable.
- **`frontend/`**: Vite + React + TanStack Router single-page application (SPA). Interactive UI for
  chat, skills management, tools, channels configuration, OAuth, models, and system settings.

### Binaries produced

| Binary | Entry point | Description |
|--------|------------|-------------|
| `picoclaw-agents-launcher` | `./web/backend` | WebUI launcher (this package) |
| `picoclaw-agents-launcher-tui` | `./cmd/picoclaw-launcher-tui` | TUI launcher (terminal UI) |
| `picoclaw-agents` | `./cmd/picoclaw` | Core CLI agent |

## Prerequisites

- Go 1.25+
- Node.js 20+ with pnpm

## Running

```bash
# Default: listens on localhost:18899, opens browser automatically
./build/picoclaw-agents-launcher

# Allow access from other devices on the network (0.0.0.0)
./build/picoclaw-agents-launcher -public

# Custom port
./build/picoclaw-agents-launcher -port 8080

# Console mode (no systray, shows logs in terminal)
./build/picoclaw-agents-launcher -console

# No auto-open browser
./build/picoclaw-agents-launcher -no-browser
```

The WebUI is available at: `http://localhost:18899`
For network access, get your local IP: `ipconfig getifaddr en0` (macOS) / `hostname -I` (Linux)

## Development

Run both the frontend dev server and the Go backend simultaneously:

```bash
make dev-launcher
# Frontend: http://localhost:5173 (proxied to http://localhost:18800)
# Backend:  http://localhost:18800
```

Or run them separately:

```bash
# Frontend only (Vite hot reload)
cd web/frontend && pnpm dev

# Backend only
cd web/backend && go run -ldflags "-X github.com/comgunner/picoclaw/pkg/config.Version=dev" .
```

## Build

Build the frontend and embed it into a single Go binary:

```bash
# Build WebUI launcher (builds frontend assets first if needed)
make build-launcher
# Output: build/picoclaw-agents-launcher

# Build TUI launcher
make build-launcher-tui
# Output: build/picoclaw-agents-launcher-tui
```

## Tests

```bash
# Run all Go tests (from repository root)
make test

# Run backend API tests only
go test ./web/backend/api/...

# Run frontend lint
cd web/frontend && pnpm lint
```

## WebSocket Chat (PicoChannel)

The frontend connects to the gateway via WebSocket through a reverse proxy:

```
Browser WebSocket
  → ws://localhost:18899/pico/ws     (WebUI launcher proxy)
    → http://127.0.0.1:18790/pico/  (gateway PicoChannel)
      → WebSocket upgrade + auth
        → messages flow to AgentLoop
```

Authentication uses `Sec-WebSocket-Protocol: token.<value>` subprotocol header.

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/gateway/status` | Gateway process state |
| POST | `/api/gateway/start` | Start gateway subprocess |
| POST | `/api/gateway/stop` | Stop gateway subprocess |
| GET | `/api/skills` | List all skills (native + embedded + workspace + global) |
| GET | `/api/skills/{name}` | Get skill content (works for native skills too) |
| POST | `/api/skills/import` | Import skill from `.md` file |
| DELETE | `/api/skills/{name}` | Delete workspace skill |
| GET | `/api/tools` | List tools with status |
| PUT | `/api/tools/{name}` | Enable/disable tool |
| GET | `/api/models` | List configured models |
| GET | `/api/oauth/providers` | OAuth provider statuses |
| GET | `/api/channels/catalog` | Supported channel list |
| GET | `/api/sessions` | Chat session history |
| GET | `/api/system/autostart` | Launch-at-login status |
| GET | `/api/system/launcher-config` | Launcher port/public config |
| GET | `/api/config` | Full `config.json` contents |
| PATCH | `/api/config` | Patch config fields |
