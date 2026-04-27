# Changelog - 26 de Abril de 2026

## [1.4.2] - 2026-04-26

### 🚀 Nuevas Funcionalidades
- **Comando `/clear` Universal**: Implementada limpieza profunda de historial de chat en Telegram, Discord y CLI.
  - Soporte para **Legacy Context** (archivos JSONL).
  - Soporte para **Seahorse Engine** (base de datos SQLite).
  - Resuelve errores críticos de `corrupted thought signature` en modelos de Antigravity/Google.
- **Mejora en Comandos Rápidos**: El sistema de comandos (Fast-Path) ahora resuelve correctamente la sesión y el agente antes de ejecutarse, garantizando que `/status`, `/compact` y `/clear` afecten al contexto correcto en grupos de Telegram y Discord.

### 🔧 Mejoras Técnicas
- **Session Manager**: Añadido método `Clear()` para purgar mensajes y resúmenes de forma atómica.
- **Seahorse Engine**: Implementado `DeleteConversationHistory()` con eliminación en cascada (mensajes, partes, sumarios y FTS5).
- **Telegram**: Actualizada lista de comandos sugeridos (Bot UI) y ayuda dinámica.
- **Discord**: Registrado nuevo Slash Command `/clear`.

### 🐞 Corrección de Errores
- **Context Routing**: Corregido bug donde los comandos rápidos fallaban al no heredar el `SessionKey` ruteado por el `AgentLoop`.
- **Antigravity Stability**: Añadida mitigación manual para contextos expirados mediante el comando de limpieza.

### 🏗️ Compilación
- Generados binarios nativos para macOS Apple Silicon (ARM64):
  - `picoclaw-agents`
  - `picoclaw-agents-launcher-darwin-arm64`
