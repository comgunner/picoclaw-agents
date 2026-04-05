# Documentación de Usuario: Queue & Batch System

El sistema de colas y procesamiento batch de PicoClaw permite realizar tareas pesadas en segundo plano sin bloquear al agente ni consumir tokens de forma innecesaria en esperas prolongadas.

> **PicoClaw v3.4.2**: Incluye **Skill Nativa Queue/Batch** compilada directamente en el binario para máximo rendimiento y cero dependencias externas.

## Herramientas Incluidas

### 1. `batch_id`
Genera un identificador único para referenciar tareas.
- **Uso**: `batch_id(prefix="SOCIAL")`
- **Resultado**: `#SOCIAL_02_03_26_1505`

### 2. `queue`
Permite listar y consultar el estado de procesos en curso.
- **Uso**:
    - `queue(action="list")`: Muestra todas las tareas activas.
    - `queue(task_id="#ID")`: Muestra el progreso de una tarea específica.

---

## Flujo de Trabajo Recomendado

1. **Iniciación**: El usuario o el agente inician una herramienta de "Macro" (como `social_post_bundle`).
2. **Desacoplamiento**: El sistema devuelve inmediatamente un ID de seguimiento y libera al LLM.
3. **Seguimiento**: El usuario puede consultar el estado usando el ID en cualquier momento.
4. **Finalización**: Al terminar, el sistema envía una notificación directa con botones de acción (**Aprobar**, **Rehacer**, **Publicar**).

## Ventajas para el Usuario
- **Velocidad**: No hay que esperar a que el LLM procese cada paso intermedio.
- **Orden**: Cada creación tiene su propia "placa" de identificación (#...).
- **Economía**: Ahorra miles de tokens en tareas automatizadas.

---

## Comandos Rápidos (Slash Commands)

PicoClaw incluye un sistema de **Fast-path** que intercepta comandos iniciados con `/` o `#` para ejecutarlos instantáneamente sin consultar a la IA. Esto garantiza respuestas inmediatas y consistencia de datos.

### Comandos de Lotes Social (Bundles)
Se usan tras recibir una notificación de tarea completada (ej: `#IMA_GEN_02_03_26_1500`):

- `/bundle_approve id=ID`: Aprueba el lote y procede a la publicación/guardado según el flujo.
- `/bundle_regen id=ID`: Solicita la regeneración completa del lote (imagen y texto).
- `/bundle_edit id=ID`: Permite editar el texto del lote.
- `/bundle_publish id=ID platforms=PLATFORMS`: Publica el lote aprobado en las plataformas especificadas (ej: `facebook,twitter`).

**⚠️ IMPORTANTE: El ID debe existir en el tracker**

Antes de publicar, verifica que el ID exista usando:
```
/list pending
```
o consulta el tracker directamente:
```
queue(action="list")
```

**IDs válidos** siguen este formato: `AAAAMMDD_HHMMSS_XXXXXX` (ej: `20260302_161740_yiia22`)

### Comandos de Utilidad
- `/show [model|channel]`: Muestra el modelo activo o el canal de comunicación.
- `/list [models|channels|agents]`: Lista las opciones configuradas en el sistema.
- `/status`: Muestra el uso actual de tokens y ocupación de la ventana de contexto.
- `/help`: Muestra la ayuda interactiva con todos los comandos disponibles.

### Comandos de Trading Binance (Basados en Herramientas)

**Nota:** Las operaciones de trading de Binance usan funciones de herramienta directamente (aún no son comandos fast-path):

| Función de Herramienta | Descripción | Ejemplo |
|------------------------|-------------|---------|
| `get_ticker_price` | Obtener precio crypto | `get_ticker_price(symbol="BTCUSDT")` |
| `get_order_book` | Libro de órdenes spot | `get_order_book(symbol="ETHUSDT", limit=10)` |
| `get_futures_order_book` | Libro de órdenes futures | `get_futures_order_book(symbol="BTCUSDT")` |
| `list_futures_volume` | Ranking por volumen | `list_futures_volume(top=10)` |
| `get_spot_balance` | Balances spot | `get_spot_balance()` (requiere API) |
| `get_futures_balance` | Balances futures | `get_futures_balance()` (requiere API) |
| `open_futures_position` | Abrir posición | `open_futures_position(symbol="BTCUSDT", side="LONG", quantity="0.001", leverage=5, confirm=true)` |
| `close_futures_position` | Cerrar posición | `close_futures_position(symbol="BTCUSDT", confirm=true)` |

**Comandos fast-path para Binance** (ej: `/binance_price`, `/binance_open`) están documentados para futura implementación.

### Canales Soportados
Estos comandos funcionan de forma idéntica en:
1. **Telegram**: A través del menú de comandos (el icono `/`).
2. **Discord**: Como comandos de aplicación (Slash Commands nativos).
3. **Terminal (CLI)**: Escribiéndolos directamente en el modo interactivo (`./picoclaw-agents agent -m ""`).

---

## Arquitectura: Global Tracker
Para asegurar que los comandos funcionen en instalaciones multi-agente, PicoClaw utiliza un `ImageGenTracker` único y compartido. Esto permite que si un **Subagente** genera una imagen, el **Agente Principal** (Project Manager) pueda aprobarla inmediatamente usando su ID, sin errores de "ID no encontrado".

## Guía para Desarrolladores (Go/Skills Nativas)

El sistema Queue/Batch es ahora una **Skill Nativa** compilada directamente en el binario de PicoClaw. Esto proporciona:

- **Cero Dependencias Externas**: No requiere archivos `.md` externos en runtime
- **Máximo Rendimiento**: Strings de documentación embebidos en el binario
- **Seguridad Mejorada**: La skill no puede ser modificada externamente
- **Actualizaciones Automáticas**: La skill se actualiza con cada release de PicoClaw

### Arquitectura de Integración

La skill nativa está implementada en `pkg/skills/queue_batch.go` y registrada en:
1. `pkg/skills/loader.go` - Registry de skills nativas
2. `pkg/agent/context.go` - Inyección en el system prompt

Para integrar nuevas herramientas nativas:
1. Crear `pkg/skills/{nombre}.go` con struct de skill y constantes de documentación
2. Registrar en `loader.go` native skills registry
3. Inyectar vía `context.go` BuildSystemPrompt()
4. Usar `tools.GetGlobalQueueManager()` para gestión de tareas
5. Ejecutar lógica pesada en goroutine separada
6. Notificar vía `MessageBus` al completar

**Ver:** `local_work/crear_skill_interna.md` para guía completa de desarrollo de skills nativas.

---

## Solución de Problemas

### Error: "ID no encontrada en el tracker"

**Síntoma:**
```
/bundle_publish id=20260302_163848_qqqia2 platforms=facebook
❌ Error: Imagen con ID 20260302_163848_qqqia2 no encontrada en el tracker
```

**Causas Posibles:**

1. **ID incorrecto o typo**: El ID no coincide con ningún registro en el tracker
2. **Tarea ya publicada**: El ID fue archivado después de publicación exitosa
3. **Tarea expirada**: El tracker limpia registros antiguos automáticamente
4. **Sesión diferente**: El ID pertenece a otra sesión/instancia de PicoClaw

**Solución:**

**Paso 1: Verificar IDs existentes**
```bash
# Lista todas las tareas pendientes
/list pending

# O consulta el queue
queue(action="list")
```

**Paso 2: Identificar el ID correcto**
Los IDs válidos tienen este formato:
```
AAAAMMDD_HHMMSS_XXXXXX
││││││││ ││││││ ││││││
││││││││ ││││││ └─── Aleatorio (6 chars)
││││││││ │││││└───── Segundos (163848 = 16:38:48)
││││││││ ││││└────── Minutos (38)
││││││││ ││└──────── Horas (16)
││││││││ └────────── Día (02)
││││││└───────────── Mes (03)
││││└─────────────── Año (2026)
```

**Paso 3: Usar el ID correcto**
```bash
# ✅ Correcto (ID existe en tracker)
/bundle_publish id=20260302_161740_yiia22 platforms=facebook

# ❌ Incorrecto (ID no existe)
/bundle_publish id=20260302_163848_qqqia2 platforms=facebook
```

**Prevención:**

1. **Copiar y pegar IDs** desde la notificación original en lugar de escribirlos manualmente
2. **Usar autocompletado** de Telegram/Discord cuando esté disponible
3. **Verificar antes de publicar** con `/list pending`
4. **Guardar IDs importantes** en un lugar seguro hasta completar la publicación

**Ejemplo de Flujo Correcto:**

```
1. Agente genera imagen:
   Bot: "✅ Imagen generada. ID: 20260302_161740_yiia22"

2. Usuario aprueba:
   Usuario: "/bundle_approve id=20260302_161740_yiia22"
   Bot: "✅ Bundle aprobado"

3. Usuario publica:
   Usuario: "/bundle_publish id=20260302_161740_yiia22 platforms=facebook"
   Bot: "🚀 Publicación completada"
```
