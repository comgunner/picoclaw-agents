# Guía para Desarrolladores de PicoClaw

> **Última Actualización:** Marzo 2026 | **Versión:** v3.4.5+

## Tabla de Contenidos

- [Introducción](#introducción)
  - [¿Qué es PicoClaw?](#qué-es-picoclaw)
  - [Descripción General de la Arquitectura](#descripción-general-de-la-arquitectura)
  - [Diferenciadores Clave](#diferenciadores-clave)
- [Configuración del Entorno de Desarrollo](#configuración-del-entorno-de-desarrollo)
  - [Requisitos de Versión de Go](#requisitos-de-versión-de-go)
  - [Herramientas Requeridas](#herramientas-requeridas)
  - [Recomendaciones de IDE](#recomendaciones-de-ide)
  - [Extensiones Recomendadas](#extensiones-recomendadas)
- [Compilación desde el Código Fuente](#compilación-desde-el-código-fuente)
  - [Comandos de Compilación](#comandos-de-compilación)
  - [Compilación Cruzada](#compilación-cruzada)
  - [Compilación con GoReleaser](#compilación-con-goreleaser)
- [Estructura del Proyecto](#estructura-del-proyecto)
  - [Diseño de Directorios](#diseño-de-directorios)
  - [Paquetes Principales](#paquetes-principales)
  - [Comandos CLI](#comandos-cli)
- [Arquitectura Multi-Agente](#arquitectura-multi-agente)
  - [Cómo Funcionan los Subagentes](#cómo-funcionan-los-subagentes)
  - [Creación de Subagentes](#creación-de-subagentes)
  - [Diferentes Modelos LLM por Subagente](#diferentes-modelos-llm-por-subagente)
  - [Bloqueos de Tareas y Prevención de Colisiones](#bloqueos-de-tareas-y-prevención-de-colisiones)
- [Arquitectura de Skills Nativos (v3.4.2+)](#arquitectura-de-skills-nativos-v342)
  - [Skills Nativos vs Externos](#skills-nativos-vs-externos)
  - [Creación de Nuevos Skills Nativos](#creación-de-nuevos-skills-nativos)
  - [Estructura del Directorio pkg/skills/](#estructura-del-directorio-pkgskills)
  - [Ejemplo: queue_batch.go](#ejemplo-queue_batchgo)
- [Desarrollo de Herramientas](#desarrollo-de-herramientas)
  - [Creación de Nuevas Herramientas](#creación-de-nuevas-herramientas)
  - [Registro de Herramientas](#registro-de-herramientas)
  - [Validación de Entrada](#validación-de-entrada)
  - [Manejo de Errores](#manejo-de-errores)
  - [Consideraciones de Seguridad](#consideraciones-de-seguridad)
- [Desarrollo de Canales](#desarrollo-de-canales)
  - [Agregar Nuevos Canales de Chat](#agregar-nuevos-canales-de-chat)
  - [Manejo de Webhooks](#manejo-de-webhooks)
  - [Formato de Mensajes](#formato-de-mensajes)
  - [Limitación de Tasa](#limitación-de-tasa)
- [Integración de Proveedores](#integración-de-proveedores)
  - [Agregar Nuevos Proveedores LLM](#agregar-nuevos-proveedores-llm)
  - [Integración OAuth](#integración-oauth)
  - [Gestión de Claves API](#gestión-de-claves-api)
- [Pruebas](#pruebas)
  - [Pruebas Unitarias](#pruebas-unitarias)
  - [Pruebas de Integración](#pruebas-de-integración)
  - [Pruebas de Seguridad](#pruebas-de-seguridad)
  - [Ejecución de Pruebas](#ejecución-de-pruebas)
- [Estilo de Código y Convenciones](#estilo-de-código-y-convenciones)
  - [Formato Go](#formato-go)
  - [Convenciones de Nombres](#convenciones-de-nombres)
  - [Estilo de Comentarios](#estilo-de-comentarios)
  - [Patrones de Manejo de Errores](#patrones-de-manejo-de-errores)
- [Flujo de Trabajo de Git](#flujo-de-trabajo-de-git)
  - [Estrategia de Ramas](#estrategia-de-ramas)
  - [Formato de Mensajes de Commit](#formato-de-mensajes-de-commit)
  - [Proceso de Pull Request](#proceso-de-pull-request)
  - [Lista de Verificación de Code Review](#lista-de-verificación-de-code-review)
- [Actualizaciones del CHANGELOG](#actualizaciones-del-changelog)
  - [Cuándo Actualizar](#cuándo-actualizar)
  - [Formato](#formato)
  - [Ejemplos](#ejemplos)
  - [Lista de Verificación Pre-commit](#lista-de-verificación-pre-commit)
- [Depuración](#depuración)
  - [Configuración de Logging](#configuración-de-logging)
  - [Modo Debug](#modo-debug)
  - [Problemas Comunes y Soluciones](#problemas-comunes-y-soluciones)
- [Optimización del Rendimiento](#optimización-del-rendimiento)
  - [Perfilado](#perfilado)
  - [Gestión de Memoria](#gestión-de-memoria)
  - [Patrones de Concurrencia](#patrones-de-concurrencia)
- [Despliegue](#despliegue)
  - [Despliegue con Docker](#despliegue-con-docker)
  - [Distribución de Binarios](#distribución-de-binarios)
  - [Consideraciones de Producción](#consideraciones-de-producción)
- [Solución de Problemas](#solución-de-problemas)
  - [Errores de Compilación Comunes](#errores-de-compilación-comunes)
  - [Problemas de Ejecución](#problemas-de-ejecución)
  - [Cómo Obtener Ayuda](#cómo-obtener-ayuda)

---

## Introducción

### ¿Qué es PicoClaw?

**PicoClaw** es un framework multi-agente de IA ultra-ligero escrito en Go. Permite ejecutar asistentes de IA personales en hardware mínimo (<10MB RAM, <1s de inicio) mientras soporta múltiples canales de chat (Telegram, Discord, Slack, etc.) y proveedores LLM (OpenAI, Anthropic, DeepSeek, Google Antigravity, etc.).

**Características Clave:**
- 🪶 **Ultra-Ligero**: <10MB de uso de RAM, <1s de tiempo de inicio
- 🤖 **Arquitectura Multi-Agente**: Subagentes paralelos con configuraciones de modelo independientes
- 🌍 **Verdadera Portabilidad**: Un solo binario a través de arquitecturas RISC-V, ARM y x86
- 🛡️ **Seguridad Primero**: Seguridad fail-close, sandboxing de workspace, Skills Sentinel
- 🚀 **Listo para Producción**: Despliegue con Docker, compilaciones GoReleaser, pruebas comprehensivas

### Descripción General de la Arquitectura

```
┌─────────────────────────────────────────────────────────────┐
│                    Arquitectura de PicoClaw                  │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │   Telegram   │  │   Discord    │  │    Slack     │      │
│  │   Canal      │  │   Canal      │  │   Canal      │      │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘      │
│         │                 │                 │               │
│         └─────────────────┼─────────────────┘               │
│                           │                                 │
│                  ┌────────▼────────┐                        │
│                  │  ChannelManager  │                        │
│                  └────────┬────────┘                        │
│                           │                                 │
│         ┌─────────────────┼─────────────────┐              │
│         │                 │                 │               │
│  ┌──────▼───────┐  ┌──────▼───────┐  ┌──────▼───────┐      │
│  │ Agent Loop 1 │  │ Agent Loop 2 │  │ Agent Loop N │      │
│  │ (Agente Prin)│  │ (Subagente 1)│  │ (Subagente N)│      │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘      │
│         │                 │                 │               │
│         └─────────────────┼─────────────────┘               │
│                           │                                 │
│                  ┌────────▼────────┐                        │
│                  │  AgentMessageBus │                        │
│                  └────────┬────────┘                        │
│                           │                                 │
│         ┌─────────────────┼─────────────────┐              │
│         │                 │                 │               │
│  ┌──────▼───────┐  ┌──────▼───────┐  ┌──────▼───────┐      │
│  │  Proveedor   │  │  Herramientas│  │    Skills    │      │
│  │   (LLM)      │  │   (Nativos)  │  │   (Nativos)  │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

**Componentes Principales:**

1. **Canales**: Adaptadores de comunicación (Telegram, Discord, Slack, etc.)
2. **Agent Loop**: Bucle principal de procesamiento para cada agente
3. **AgentMessageBus**: Sistema de comunicación inter-agente
4. **Proveedores**: Integraciones LLM (OpenAI, Anthropic, etc.)
5. **Herramientas**: Capacidades nativas (exec, filesystem, web search, etc.)
6. **Skills**: Conocimiento y workflows compilados

### Diferenciadores Clave

| Característica | PicoClaw | NanoBot | OpenClaw |
|---------|----------|---------|----------|
| **Lenguaje** | Go 1.25.8 | Python 3.11+ | TypeScript |
| **Uso de RAM** | <10MB | ~50MB | ~500MB |
| **Tiempo de Inicio** | <1s | ~2s | ~10s |
| **Tamaño de Código** | ~10K líneas | ~4K líneas | 430K+ líneas |
| **Ideal Para** | Embedded/IoT | Investigación/Aprendizaje | Producción |
| **Apps Móviles** | ❌ | ❌ | ✅ |
| **Tamaño Binario** | ~15MB | N/A | N/A |

---

## Configuración del Entorno de Desarrollo

### Requisitos de Versión de Go

**Versión Mínima de Go:** 1.25.8

PicoClaw requiere Go 1.25.8 o posterior debido a:
- Soporte mejorado de genéricos
- Logging `slog` mejorado
- Mejores patrones de manejo de errores
- Optimizaciones de rendimiento en Go 1.25+

**Verificar Versión de Go:**
```bash
go version
# Esperado: go version go1.25.8 ...
```

**Instalar/Actualizar Go:**
```bash
# macOS (Homebrew)
brew install go@1.25

# Linux (snap)
sudo snap install go --classic

# Instalación manual
wget https://go.dev/dl/go1.25.8.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.25.8.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
```

### Herramientas Requeridas

| Herramienta | Versión | Propósito | Instalación |
|------|---------|---------|--------------|
| **Go** | 1.25.8+ | Lenguaje principal | [go.dev](https://go.dev/dl/) |
| **make** | 4.0+ | Automatización de compilación | `apt install make` / `brew install make` |
| **git** | 2.30+ | Control de versiones | `apt install git` / `brew install git` |
| **docker** | 24.0+ | Contenerización | [docker.com](https://docs.docker.com/get-docker/) |
| **golangci-lint** | 1.60+ | Linting | `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest` |

**Verificar Instalación:**
```bash
go version
make --version
git --version
docker --version
golangci-lint --version
```

### Recomendaciones de IDE

#### VS Code (Recomendado)

**Por qué VS Code:**
- Excelente soporte Go vía extensión oficial
- Ligero y rápido
- Ecosistema rico en extensiones
- Terminal y depurador integrados

**Instalación:**
```bash
# macOS
brew install --cask visual-studio-code

# Linux (snap)
sudo snap install code --classic

# Windows
# Descargar desde https://code.visualstudio.com/
```

#### GoLand (JetBrains)

**Por qué GoLand:**
- IDE profesional con características Go avanzadas
- Herramientas de base de datos integradas
- Excelente soporte de refactoring
- Perfilador integrado

**Instalación:**
```bash
# macOS (Homebrew Cask)
brew install --cask goland

# Linux (snap)
sudo snap install goland --classic

# Licencia: Comercial (gratis para estudiantes)
```

#### Neovim (Usuarios Avanzados)

**Por qué Neovim:**
- Extremadamente ligero
- Altamente personalizable
- Soporte LSP nativo

**Configuración:**
```bash
# Instalar Neovim
brew install neovim  # macOS
sudo apt install neovim  # Linux

# Instalar plugin Go
nvim --headless "+Lazy! sync" +quit
```

### Extensiones Recomendadas

#### Extensiones de VS Code

1. **Go** (Oficial)
   ```json
   {
     "gopls": {
       "ui.semanticTokens": true,
       "ui.diagnostic.analyses": {
         "unusedparams": true,
         "shadow": true,
         "nilness": true,
         "unusedwrite": true,
         "useany": true
       }
     }
   }
   ```

2. **GitLens** — Git superpotenciado
3. **Docker** — Gestión de contenedores
4. **Remote - SSH** — Desarrollo remoto
5. **YAML** — Soporte YAML para archivos de configuración
6. **Markdown All in One** — Escritura de documentación

#### Plugins de GoLand

1. **Go** (Incluido)
2. **Docker** (Incluido)
3. **GitToolBox** — Integración Git mejorada
4. **String Manipulation** — Utilidades de texto
5. **Key Promoter X** — Aprender atajos

---

## Compilación desde el Código Fuente

### Comandos de Compilación

PicoClaw usa un `Makefile` para automatización de compilación con múltiples objetivos:

```bash
# Clonar repositorio
git clone https://github.com/comgunner/picoclaw-agents.git
cd picoclaw

# Descargar dependencias
make deps

# Compilar para plataforma actual (ejecuta go generate primero)
make build

# Compilar e instalar en ~/.local/bin
make install

# Ejecutar todas las verificaciones pre-commit
make check

# Limpiar artefactos de compilación
make clean
```

**Salida de Compilación:**
```
Building picoclaw for darwin/arm64...
Build complete: build/picoclaw-darwin-arm64
```

**Ubicación del Binario:**
```
build/picoclaw-{platform}-{arch}
build/picoclaw  # Enlace simbólico a plataforma actual
```

### Compilación Cruzada

Compilar para todas las plataformas soportadas:

```bash
make build-all
```

**Plataformas Soportadas:**

| Plataforma | Arquitectura | Nombre Binario |
|----------|--------------|-------------|
| Linux | amd64 | `picoclaw-linux-amd64` |
| Linux | arm64 | `picoclaw-linux-arm64` |
| Linux | loong64 | `picoclaw-linux-loong64` |
| Linux | riscv64 | `picoclaw-linux-riscv64` |
| Linux | armv7 | `picoclaw-linux-armv7` |
| macOS | arm64 (Apple Silicon) | `picoclaw-darwin-arm64` |
| macOS | amd64 (Intel) | `picoclaw-darwin-amd64` |
| Windows | amd64 | `picoclaw-windows-amd64.exe` |

**Compilación Cruzada Manual:**
```bash
# Compilar para Linux ARM64
GOOS=linux GOARCH=arm64 go build -o picoclaw-linux-arm64 ./cmd/picoclaw

# Compilar para Windows AMD64
GOOS=windows GOARCH=amd64 go build -o picoclaw-windows-amd64.exe ./cmd/picoclaw
```

### Compilación con GoReleaser

Para releases de producción, PicoClaw usa GoReleaser:

```bash
# Instalar GoReleaser
go install github.com/goreleaser/goreleaser/v2@latest

# Ejecutar GoReleaser (build snapshot)
goreleaser release --snapshot --clean

# Release completo (requiere GITHUB_TOKEN)
goreleaser release --clean
```

**Salidas de GoReleaser:**
- Binarios para todas las plataformas
- Imágenes Docker
- Paquetes RPM/DEB
- Archivos (tar.gz, zip)

**Configuración:** `.goreleaser.yaml`

---

*(Continúa en la siguiente sección debido a la extensión del documento)*

## Estructura del Proyecto

### Diseño de Directorios

```
picoclaw/
├── cmd/picoclaw/              # Punto de entrada CLI
│   ├── main.go                # Función main
│   └── internal/              # Comandos CLI
│       ├── agent/             # Comando agent
│       ├── agents/            # Gestión de agentes
│       ├── auth/              # Autenticación
│       ├── clean/             # Comando clean
│       ├── cron/              # Jobs Cron
│       ├── gateway/           # Comando gateway
│       ├── migrate/           # Utilidades de migración
│       ├── onboard/           # Configuración inicial
│       ├── skills/            # Gestión de skills
│       ├── status/            # Comando status
│       ├── util/              # Comandos de utilidad
│       └── version/           # Comando version
│
├── pkg/                       # Paquetes principales (API pública)
│   ├── agent/                 # Loop y instancia de agente
│   ├── agents/                # Coordinación multi-agente
│   ├── auth/                  # Autenticación (OAuth)
│   ├── bus/                   # Bus de mensajes
│   ├── channels/              # Adaptadores de canales de chat
│   ├── config/                # Configuración
│   ├── constants/             # Constantes
│   ├── context/               # Gestión de contexto
│   ├── cron/                  # Planificador Cron
│   ├── devices/               # Soporte para dispositivos IoT
│   ├── gateway/               # Gateway WebSocket
│   ├── health/                # Checks de salud
│   ├── heartbeat/             # Sistema heartbeat
│   ├── logger/                # Logging
│   ├── memory/                # Almacenamiento de memoria
│   ├── migrate/               # Migraciones
│   ├── providers/             # Proveedores LLM
│   ├── routing/               # Enrutamiento de mensajes
│   ├── security/              # Herramientas de seguridad
│   ├── session/               # Gestión de sesiones
│   ├── skills/                # Sistema de skills
│   ├── state/                 # Gestión de estado
│   ├── tasklock/              # Bloqueo de tareas
│   ├── tools/                 # Herramientas nativas
│   ├── utils/                 # Utilidades
│   └── voice/                 # Procesamiento de voz
│
├── internal/                  # Implementación interna (privada)
│
├── config/                    # Plantillas de configuración
│   ├── config.example.json    # Configuración de ejemplo
│   ├── config_dev.example.json # Configuración de desarrollo
│   └── ...
│
├── docs/                      # Documentación
│   ├── DEVELOPER_GUIDE.md     # Este archivo
│   ├── CONTRIBUTING.md        # Guías de contribución
│   ├── CHANGELOG.md           # Historial de versiones
│   ├── SECURITY.md            # Documentación de seguridad
│   └── ...                    # Docs específicos de características
│
├── workspace/                 # Workspace por defecto del agente
│   ├── sessions/              # Historial de conversaciones
│   ├── memory/                # Memoria a largo plazo
│   ├── state/                 # Estado persistente
│   ├── cron/                  # Jobs programados
│   └── skills/                # Skills personalizados
│
├── scripts/                   # Scripts de utilidad y compilación
├── assets/                    # Imágenes y recursos
├── releases/                  # Artefactos de release
├── local_work/                # Trabajo personal (NO commiteado)
├── Makefile                   # Automatización de compilación
├── go.mod                     # Definición de módulo Go
├── go.sum                     # Checksums de dependencias
├── .goreleaser.yaml           # Configuración de GoReleaser
├── docker-compose.yml         # Servicios Docker
└── Dockerfile                 # Imagen de contenedor
```

### Paquetes Principales

#### pkg/agent/

**Propósito:** Loop del agente, gestión de instancias, manejo de contexto

**Archivos Clave:**
- `loop.go` — Loop principal de procesamiento del agente
- `instance.go` — Representación de instancia de agente
- `context.go` — Construcción y gestión de contexto
- `context_compactor.go` — Compactación de contexto
- `memory.go` — Memoria del agente
- `registry.go` — Registro de agentes

**Ejemplo de Uso:**
```go
import "github.com/comgunner/picoclaw/pkg/agent"

// Crear instancia de agente
agent := agent.NewAgentInstance(config, workspace)

// Ejecutar loop del agente
err := agent.Run(ctx)
```

#### pkg/providers/

**Propósito:** Integraciones de proveedores LLM

**Proveedores Soportados:**
- OpenAI (GPT-4, o3-mini)
- Anthropic (Claude 3, Claude 4)
- DeepSeek (DeepSeek Chat, DeepSeek Reasoner)
- Google Antigravity (Gemini, Claude vía Google)
- GitHub Copilot (Codex, Copilot CLI)
- OpenRouter (Gateway multi-proveedor)
- Zhipu (Modelos GLM)
- Mistral (Mistral, Mixtral)

**Estructura de Directorios:**
```
pkg/providers/
├── factory.go                 # Factory de proveedores
├── types.go                   # Interfaces de proveedor
├── openai_compat/             # Proveedores compatibles con OpenAI
├── antigravity_provider.go    # Google Antigravity
├── claude_provider.go         # Anthropic Claude
├── github_copilot_provider.go # GitHub Copilot
└── ...
```

#### pkg/tools/

**Propósito:** Herramientas nativas (capacidades)

**Herramientas Disponibles:**
- `exec` — Ejecución de comandos shell
- `read_file`, `write_file`, `edit_file` — Operaciones de filesystem
- `web_search`, `web_fetch` — Operaciones web
- `spawn` — Creación de subagentes
- `queue`, `batch_id` — Delegación queue/batch
- `memory_store` — Operaciones de memoria
- `config_manager` — Gestión de configuración
- `system_diagnostics` — Monitoreo del sistema
- `binance` — Trading de criptomonedas
- `social_media` — Publicación en redes sociales
- `notion` — Operaciones de Notion
- `image_gen` — Generación de imágenes

#### pkg/channels/

**Propósito:** Adaptadores de canales de chat

**Canales Soportados:**
- Telegram (`telego`)
- Discord (`discordgo`)
- Slack (`slack-go`)
- QQ (`botgo`)
- DingTalk (`dingtalk-stream-sdk-go`)
- LINE
- WeCom (WeChat Work)
- Feishu/Lark (`oapi-sdk-go`)

#### pkg/skills/

**Propósito:** Sistema de skills (conocimiento compilado)

**Skills Nativos:**
- `queue_batch` — Workflow de delegación queue/batch

**Skills Externos:**
- Skills de workspace (`~/.picoclaw/workspace/skills/`)
- Skills globales (`~/.picoclaw/skills/`)
- Skills builtin (compilados en el binario)

### Comandos CLI

PicoClaw proporciona múltiples comandos CLI vía Cobra:

```bash
# Configuración inicial
picoclaw onboard

# Operaciones de agente
picoclaw agent -m "Hello"           # Query one-shot
picoclaw agent interactive          # Modo interactivo
picoclaw agents list                # Listar agentes
picoclaw agents spawn <name>        # Crear subagente

# Gateway (bot de larga duración)
picoclaw gateway

# Autenticación
picoclaw auth login --provider google-antigravity
picoclaw auth status
picoclaw auth logout --provider google-antigravity

# Gestión de skills
picoclaw skills list
picoclaw skills install <name>
picoclaw skills remove <name>

# Jobs Cron
picoclaw cron list
picoclaw cron add "0 * * * *" "Verificar estado"
picoclaw cron remove <id>

# Utilidades
picoclaw status                     # Estado del sistema
picoclaw clean --all                # Limpiar sesiones
picoclaw version                    # Información de versión
picoclaw migrate                    # Ejecutar migraciones
```

---

*(Nota: Esta es la primera mitad de la guía. Las secciones restantes cubren arquitectura multi-agente, desarrollo de skills nativos, herramientas, canales, proveedores, pruebas, estilo de código, Git workflow, CHANGELOG, depuración, optimización, despliegue y solución de problemas.)*

## Recursos Adicionales

- **Documentación Completa en Inglés:** [`DEVELOPER_GUIDE.md`](./DEVELOPER_GUIDE.md)
- **Guía de Contribución:** [`CONTRIBUTING.es.md`](./CONTRIBUTING.es.md)
- **Documentación de Seguridad:** [`SECURITY.es.md`](./SECURITY.es.md)
- **README en Español:** [`../../README.es.md`](../../README.es.md)

---

*PicoClaw: IA Ultra-Eficiente en Go. Hardware de $10 · 10MB RAM · <1s de Inicio*
