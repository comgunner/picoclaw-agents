# Arquitectura de Comunicación de PicoClaw 🦞

Este documento describe el mapa técnico de cómo viajan los datos y cómo los agentes se comunican entre sí, desde el momento en que un usuario envía un mensaje hasta que recibe la respuesta final.

> **PicoClaw v3.4.1**: Incluye **Comandos Rápidos (Fast-path Slash Commands)** para operaciones instantáneas y **Rastreador Global (Global Tracker)** para consistencia multi-agente.

## 1. Entrada de Datos (Usuario ➔ Canal)

Cuando envías un mensaje (ej. a través de Telegram), el flujo comienza en la capa de **Canales**.

- **Archivo**: `pkg/channels/telegram.go` (u otros canales).
- **Acción**: Recibe el webhook/mensaje de la API de Telegram.
- **Formato**: Transforma el objeto nativo de Telegram en una estructura interna `InboundMessage`.

```json
// Representación lógica de InboundMessage
{
  "channel": "telegram",
  "chat_id": "12345678",
  "sender_id": "@user",
  "content": "Analiza este código...",
  "session_key": "telegram:12345678",
  "metadata": { ... }
}
```

## 2. Orquestación Central (MessageBus ➔ AgentLoop)

PicoClaw utiliza un **MessageBus** interno para desacoplar la recepción de la ejecución.

- **Archivo**: `pkg/agent/loop.go`
- **Componente**: `AgentLoop.Run()` consume mensajes del bus.
- **Enrutamiento**: Determina qué agente debe responder (Project Manager, Senior Dev, etc.) consultando el `AgentRegistry`.

## 3. Construcción del Contexto (El "JSON" que ve el LLM)

Antes de invocar la IA, PicoClaw construye un mensaje masivo que contiene toda la "memoria" necesaria.

- **Archivo**: `pkg/agent/context.go` (`BuildMessages`)
- **Acción**: Combina múltiples fuentes en un único arreglo de mensajes.

### Estructura del Contexto Enviado:
1. **System Prompt (Estático)**: Instrucciones de `IDENTITY.md`, `SOUL.md`, etc. (Este bloque se cachea para ahorrar tokens/tiempo).
2. **Contexto Dinámico**: Fecha y hora actuales, y entorno (`runtime`).
3. **Resumen de Contexto**: Si la conversación es muy larga, se inyecta un resumen compacto de los mensajes antiguos.
4. **Historial**: Los últimos mensajes reales de la sesión (ej. últimos 10-15 turnos).
5. **Mensaje Actual**: El mensaje recién ingresado por el usuario.

## 4. Comunicación Inter-Agentes (Subagentes / Spawn)

Esta es la parte más dinámica. Cuando un agente decide que necesita ayuda, usa la herramienta `spawn`.

- **Formato de Petición**: El LLM genera una llamada a herramienta (`tool_call`) de tipo `spawn`.
- **Ejecución**: El `SubagentManager` crea una nueva instancia de ejecución para el agente objetivo.
- **"Chat" Interno**:
    1. Agente 1 (Padre) llama a `spawn(agent_id="senior_dev", task="...")`.
    2. El "fix" o proceso se ejecuta en segundo plano.
    3. Al completarse, el resultado se envía al canal especial `system`.
    4. **Archivo**: `pkg/agent/loop.go` (`processSystemMessage`).
    5. Este método captura el resultado y lo reinyecta en el contexto del Agente 1 como un informe técnico.

```json
// Cómo viaja el resultado interno del subagente
{
  "channel": "system",
  "chat_id": "telegram:12345678", // Mapea al chat original
  "sender_id": "senior_dev",
  "content": "[Sistema: senior_dev] Tarea completada. Fix aplicado en pkg/db..."
}
```

## 5. Salida de Datos (Agente ➔ Salida)

- **Acción**: Una vez el Project Manager (u otro agente principal) genera el texto final para el usuario.
- **Archivo**: `pkg/agent/loop.go` publica un `OutboundMessage` hacia el bus.
- **Entrega**: El mánager de Telegram escucha este mensaje y lo envía directamente al teléfono del usuario.

## Conceptos Clave de Eficiencia

- **Compactación**: Si el historial de mensajes excede el 75% del límite del modelo, el `ContextCompactor` entra en acción, resume los mensajes viejos y libera espacio, evitando errores de "Límite de Contexto Excedido".
- **Memoria Basada en Disco**: Todas las sesiones se guardan en archivos `.json` dentro de la carpeta `workspace/sessions/`, permitiendo que el bot "recuerde" incluso si el servidor se reinicia.

---

## Nuevas Funcionalidades v3.4.1

### Comandos Rápidos (Fast-path Slash Commands)

Los comandos que comienzan con `/` o `#` son interceptados antes de alcanzar el LLM:

1. **El usuario envía**: `/bundle_approve id=20260302_1600`
2. **Intercepción rápida**: No se consulta con el LLM.
3. **Ejecución directa**: Validación de ID y ejecución inmediata.
4. **Respuesta instantánea**: Confirmación provista en <100ms.

**Ventajas:**
- ✅ Latencia cero (sin necesidad del poder de razonamiento del LLM).
- ✅ Consistencia de datos garantizada.
- ✅ Mismo comportamiento para Telegram, Discord y la CLI local.

### Rastreador Global (Global Tracker)

El **ImageGenTracker** ahora se comparte entre todos los agentes:

```
Main Agent (PM)
    │
    ├─► Subagent 1 (Imágenes) ──► Guarda en el Espacio Local Global
    ├─► Subagent 2 (Textos) ────► Guarda en el Espacio Local Global
    └─► Accede inmediatamente a ambos resultados ✅
```

**Beneficios:**
- ✅ Se acabaron los errores de archivos no encontrados ("ID not found") entre divisiones de agentes.
- ✅ Estado de variables a tiempo real (Shared State).
- ✅ Soporte multi-agente sumamente escalable.

---

## Configuración de Modelos

### Modelo por Defecto (deepseek-chat)

El `config.example.json` nativo usa **deepseek-chat** como modelo principal de inicialización en agentes. Esto provee un punto de partida poderoso y costo-eficiente:

```json
{
  "agents": {
    "defaults": {
      "model": "deepseek-chat"
    },
    "list": [
      {
        "id": "project_manager",
        "model": "deepseek-chat"
      },
      {
        "id": "senior_dev",
        "model": "deepseek-chat"
      }
      // ... todos los agentes asumen usar deepseek-chat por defecto
    ]
  }
}
```

### Alternativa: Proveedor Antigravity

Para usuarios con planes como *Google One AI Premium* o *Workspace Gemini*, el archivo pre-elaborado `config.example_antigravity.json` invoca eficientemente a **antigravity-gemini-3-flash** para todos los roles:

```bash
cp config/config.example_antigravity.json ~/.picoclaw/config.json
./picoclaw-agents auth antigravity
```

### Tubería de Resolución de Modelos (Pipeline)

1. **Carga Inicial**: El arreglo virtual (`model_list`) de `config.json` se carga a RAM durante el booteo.
2. **Enrutamiento Fabricante-API**: El esquema del atributo `model` (ej. `deepseek/deepseek-chat`) determina el proveedor en uso.
3. **Llamada de Red**: Las funciones delegadas y exclusivas al proveedor toman y manejan el request local.

> **Nota Crítica**: Los cambios a `model_list` implican un reinicio total de la arquitectura. Para cambiar los modelos subyacentes, apaga tu instancia de PicoClaw, modifica el valor de `model` en `config.json` y vuelve a inicializarlo.
