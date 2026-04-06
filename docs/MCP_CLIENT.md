# MCP Client — External MCP Server Integration

## Overview

PicoClaw-Agents can now consume tools from external **MCP (Model Context Protocol)** servers. This means you can connect to any MCP-compatible server (GitHub, Filesystem, PostgreSQL, Brave Search, etc.) and their tools will automatically appear in your agent's tool registry.

**Protocol version:** MCP 2024-11-05  
**Supported transports:** stdio, SSE, HTTP

## Quick Start

### 1. Install an MCP Server

Install any MCP-compatible server. Example using the GitHub MCP server:

```bash
npm install -g @modelcontextprotocol/server-github
```

### 2. Add to config.json

Edit `~/.picoclaw/config.json` and add your MCP server configuration:

```json
{
  "tools": {
    "mcp": {
      "enabled": true,
      "default_timeout": "30s",
      "servers": {
        "github": {
          "transport": "stdio",
          "command": "npx",
          "args": ["-y", "@modelcontextprotocol/server-github"],
          "env": {
            "GITHUB_PERSONAL_ACCESS_TOKEN": "ghp_your_token_here"
          },
          "description": "GitHub: repos, issues, PRs"
        }
      }
    }
  }
}
```

### 3. Restart Your Agent

```bash
./build/picoclaw-agents gateway
# or
./build/picoclaw-agents agent -m "List my GitHub repos"
```

The agent will automatically connect to the MCP server and register all available tools.

## Configuration

### Full Configuration Reference

```json
{
  "tools": {
    "mcp": {
      "enabled": true,
      "default_timeout": "30s",
      "servers": {
        "filesystem": {
          "transport": "stdio",
          "command": "npx",
          "args": ["-y", "@modelcontextprotocol/server-filesystem", "/home/user/projects"],
          "enabled_tools": ["read_file", "write_file", "list_dir"],
          "timeout": "60s",
          "description": "Filesystem access (sandboxed)"
        },
        "github": {
          "transport": "stdio",
          "command": "npx",
          "args": ["-y", "@modelcontextprotocol/server-github"],
          "env": {
            "GITHUB_PERSONAL_ACCESS_TOKEN": "ghp_..."
          },
          "description": "GitHub: repos, issues, PRs"
        },
        "postgres": {
          "transport": "stdio",
          "command": "npx",
          "args": ["-y", "@modelcontextprotocol/server-postgres", "postgresql://localhost/mydb"]
        },
        "brave-search": {
          "transport": "stdio",
          "command": "npx",
          "args": ["-y", "@modelcontextprotocol/server-brave-search"],
          "env": {
            "BRAVE_API_KEY": "your_api_key"
          }
        },
        "sse-example": {
          "transport": "sse",
          "url": "http://localhost:3001/sse",
          "headers": {
            "Authorization": "Bearer your_token"
          }
        }
      }
    }
  }
}
```

### Configuration Options

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `enabled` | boolean | Yes | Enable/disable MCP client |
| `default_timeout` | duration | No | Default timeout for tool calls (default: 30s) |
| `servers` | map | Yes | Map of server name → configuration |
| `transport` | string | Yes | Transport type: `stdio`, `sse`, `http` |
| `command` | string | stdio only | Command to execute |
| `args` | array | stdio only | Command arguments |
| `env` | map | No | Environment variables for the subprocess |
| `url` | string | SSE/HTTP only | Server URL |
| `headers` | map | No | HTTP headers for SSE/HTTP transports |
| `enabled_tools` | array | No | Whitelist of tools to register (`["*"]` = all) |
| `timeout` | duration | No | Per-server timeout override |
| `description` | string | No | Human-readable description |

## CLI Commands

Manage MCP servers from the command line:

```bash
# List all configured MCP servers
picoclaw-agents mcp list

# Show connection status for a server
picoclaw-agents mcp status <server_name>

# Add a new MCP server
picoclaw-agents mcp add <name> --transport stdio --command npx --args "...args..."

# Remove an MCP server
picoclaw-agents mcp remove <name>
```

## Available MCP Servers

| Server | Command | Description |
|--------|---------|-------------|
| **GitHub** | `npx @modelcontextprotocol/server-github` | Repos, issues, PRs, commits |
| **Filesystem** | `npx @modelcontextprotocol/server-filesystem` | Read/write files (sandboxed) |
| **PostgreSQL** | `npx @modelcontextprotocol/server-postgres` | Database queries |
| **Brave Search** | `npx @modelcontextprotocol/server-brave-search` | Web search |
| **SQLite** | `npx @modelcontextprotocol/server-sqlite` | SQLite database access |
| **Slack** | `npx @modelcontextprotocol/server-slack` | Slack messaging |
| **Google Drive** | `npx @modelcontextprotocol/server-google-drive` | File management |
| **Puppeteer** | `npx @modelcontextprotocol/server-puppeteer` | Browser automation |

Find more servers at the [MCP Server Registry](https://github.com/modelcontextprotocol/servers).

## Security

### Command Whitelist

For `stdio` transports, only a predefined set of commands are allowed to spawn subprocesses:

```
npx, node, python, python3, pipx, uvx, go, picoclaw-agents
```

This prevents arbitrary shell execution via malicious MCP server configurations. To allow a new command, you must edit the source code whitelist (`pkg/config/config.go`) or fork the project.

### Tool Filtering

Use `enabled_tools` to restrict which tools from a server are registered:

```json
{
  "servers": {
    "filesystem": {
      "enabled_tools": ["read_file", "list_dir"]
    }
  }
}
```

### Non-Fatal Failures

If an MCP server fails to connect, the agent **continues with the remaining servers**. A single broken server won't crash the agent.

### Workspace Sandboxing

MCP tools that access the filesystem are subject to the same workspace restrictions as native tools (`restrict_to_workspace: true` by default).

## Troubleshooting

### "Server failed to connect"

- Verify the command is in the allowed list: `npx`, `node`, `python`, etc.
- Check that the MCP server is installed: `npx @modelcontextprotocol/server-github --help`
- Review logs for stderr output from the subprocess.

### "Tool not found"

- Check `enabled_tools` in your config — the tool may be filtered out.
- Run `picoclaw-agents mcp list` to see available tools.

### "Context deadline exceeded"

- Increase the `timeout` for that server.
- The MCP server may be slow or hung — restart the agent.

### "Parse error" / "Malformed JSON"

- The MCP server returned invalid JSON. Check server logs and version compatibility.
- Ensure you're using an MCP 2024-11-05 compatible server.

## Protocol Details

- **Protocol:** Model Context Protocol (MCP) 2024-11-05
- **JSON-RPC:** Version 2.0
- **Transports:** stdio (subprocess), SSE (Server-Sent Events), HTTP (REST)
- **Initialization:** `initialize` → `notifications/initialized` → `tools/list`
- **Tool calls:** `tools/call` with `{name, arguments}`
- **Max response size:** 10 MB (MAX_LINE_BYTES protection)
