# Cliente MCP — Integración con Servidores MCP Externos

## Descripción General

PicoClaw-Agents ahora puede consumir herramientas de servidores **MCP (Model Context Protocol)** externos. Esto significa que puedes conectarte a cualquier servidor compatible con MCP (GitHub, Filesystem, PostgreSQL, Brave Search, etc.) y sus herramientas aparecerán automáticamente en el registro de herramientas de tu agente.

**Versión del protocolo:** MCP 2024-11-05  
**Transportes soportados:** stdio, SSE, HTTP

## Inicio Rápido

### 1. Instalar un Servidor MCP

Instala cualquier servidor MCP compatible. Ejemplo usando el servidor de GitHub:

```bash
npm install -g @modelcontextprotocol/server-github
```

### 2. Agregar a config.json

Edita `~/.picoclaw/config.json` y agrega la configuración de tu servidor MCP:

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
            "GITHUB_PERSONAL_ACCESS_TOKEN": "ghp_tu_token_aqui"
          },
          "description": "GitHub: repos, issues, PRs"
        }
      }
    }
  }
}
```

### 3. Reiniciar tu Agente

```bash
./build/picoclaw-agents gateway
# o
./build/picoclaw-agents agent -m "Lista mis repositorios de GitHub"
```

El agente se conectará automáticamente al servidor MCP y registrará todas las herramientas disponibles.

## Configuración

### Referencia Completa de Configuración

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
          "description": "Acceso al sistema de archivos (sandboxed)"
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
            "BRAVE_API_KEY": "tu_api_key"
          }
        },
        "sse-ejemplo": {
          "transport": "sse",
          "url": "http://localhost:3001/sse",
          "headers": {
            "Authorization": "Bearer tu_token"
          }
        }
      }
    }
  }
}
```

### Opciones de Configuración

| Campo | Tipo | Requerido | Descripción |
|-------|------|-----------|-------------|
| `enabled` | booleano | Sí | Habilitar/deshabilitar cliente MCP |
| `default_timeout` | duración | No | Timeout por defecto para llamadas (defecto: 30s) |
| `servers` | mapa | Sí | Mapa de nombre del servidor → configuración |
| `transport` | string | Sí | Tipo de transporte: `stdio`, `sse`, `http` |
| `command` | string | stdio | Comando a ejecutar |
| `args` | array | stdio | Argumentos del comando |
| `env` | mapa | No | Variables de entorno para el subproceso |
| `url` | string | SSE/HTTP | URL del servidor |
| `headers` | mapa | No | Encabezados HTTP para transportes SSE/HTTP |
| `enabled_tools` | array | No | Lista blanca de herramientas (`["*"]` = todas) |
| `timeout` | duración | No | Timeout específico por servidor |
| `description` | string | No | Descripción legible por humanos |

## Comandos CLI

Gestiona servidores MCP desde la línea de comandos:

```bash
# Listar todos los servidores MCP configurados
picoclaw-agents mcp list

# Mostrar estado de conexión de un servidor
picoclaw-agents mcp status <nombre_servidor>

# Agregar un nuevo servidor MCP
picoclaw-agents mcp add <nombre> --transport stdio --command npx --args "...args..."

# Eliminar un servidor MCP
picoclaw-agents mcp remove <nombre>
```

## Servidores MCP Disponibles

| Servidor | Comando | Descripción |
|----------|---------|-------------|
| **GitHub** | `npx @modelcontextprotocol/server-github` | Repos, issues, PRs, commits |
| **Filesystem** | `npx @modelcontextprotocol/server-filesystem` | Leer/escribir archivos (sandboxed) |
| **PostgreSQL** | `npx @modelcontextprotocol/server-postgres` | Consultas de base de datos |
| **Brave Search** | `npx @modelcontextprotocol/server-brave-search` | Búsqueda web |
| **SQLite** | `npx @modelcontextprotocol/server-sqlite` | Acceso a base de datos SQLite |
| **Slack** | `npx @modelcontextprotocol/server-slack` | Mensajería Slack |
| **Google Drive** | `npx @modelcontextprotocol/server-google-drive` | Gestión de archivos |
| **Puppeteer** | `npx @modelcontextprotocol/server-puppeteer` | Automatización del navegador |

Encuentra más servidores en el [Registro de Servidores MCP](https://github.com/modelcontextprotocol/servers).

## Seguridad

### Lista Blanca de Comandos

Para el transporte `stdio`, solo un conjunto predefinido de comandos pueden ejecutar subprocesos:

```
npx, node, python, python3, pipx, uvx, go, picoclaw-agents
```

Esto previene la ejecución arbitraria de comandos mediante configuraciones MCP maliciosas. Para permitir un nuevo comando, debes editar la lista blanca en el código fuente (`pkg/config/config.go`) o hacer fork del proyecto.

### Filtrado de Herramientas

Usa `enabled_tools` para restringir qué herramientas de un servidor se registran:

```json
{
  "servers": {
    "filesystem": {
      "enabled_tools": ["read_file", "list_dir"]
    }
  }
}
```

### Fallos No Fatales

Si un servidor MCP falla al conectarse, el agente **continúa con los servidores restantes**. Un solo servidor roto no crasheará el agente.

### Aislamiento del Espacio de Trabajo

Las herramientas MCP que acceden al sistema de archivos están sujetas a las mismas restricciones de espacio de trabajo que las herramientas nativas (`restrict_to_workspace: true` por defecto).

## Solución de Problemas

### "Server failed to connect"

- Verifica que el comando esté en la lista permitida: `npx`, `node`, `python`, etc.
- Comprueba que el servidor MCP esté instalado: `npx @modelcontextprotocol/server-github --help`
- Revisa los logs para ver la salida stderr del subproceso.

### "Tool not found"

- Verifica `enabled_tools` en tu configuración — la herramienta puede estar filtrada.
- Ejecuta `picoclaw-agents mcp list` para ver las herramientas disponibles.

### "Context deadline exceeded"

- Incrementa el `timeout` para ese servidor.
- El servidor MCP puede estar lento o colgado — reinicia el agente.

### "Parse error" / "Malformed JSON"

- El servidor MCP devolvió JSON inválido. Revisa los logs del servidor y la compatibilidad de versiones.
- Asegúrate de usar un servidor compatible con MCP 2024-11-05.

## Detalles del Protocolo

- **Protocolo:** Model Context Protocol (MCP) 2024-11-05
- **JSON-RPC:** Versión 2.0
- **Transportes:** stdio (subproceso), SSE (Server-Sent Events), HTTP (REST)
- **Inicialización:** `initialize` → `notifications/initialized` → `tools/list`
- **Llamadas a herramientas:** `tools/call` con `{name, arguments}`
- **Tamaño máximo de respuesta:** 10 MB (protección MAX_LINE_BYTES)
