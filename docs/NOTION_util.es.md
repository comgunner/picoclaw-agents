# Notion Util

Guía rápida para usar herramientas de Notion en PicoClaw desde terminal y Telegram.

> **PicoClaw v3.4.1**: Ahora soporta **Comandos Slash Fast-path** para operaciones instantáneas y **Global Tracker** para consistencia multi-agente.

## Requisitos

Configura tu credencial en `~/.picoclaw/config.json`:

```json
{
  "tools": {
    "notion": {
      "api_key": "TU_NOTION_API_KEY"
    }
  }
}
```

También puedes usar variables de entorno:

```bash
export PICOCLAW_TOOLS_NOTION_API_KEY="tu_notion_api_key"
```

### Obtener API Key de Notion

1. Ve a https://notion.so/my-integrations
2. Crea una nueva integración (+ New integration)
3. Copia el API key (empieza con `ntn_` o `secret_`)
4. Guarda el key en `~/.config/notion/api_key` o en el config.json

### Conectar páginas/bases de datos

1. Abre la página o database que quieres conectar
2. Click en "..." (más opciones)
3. Selecciona "Connect to" → tu integración
4. Ahora la integración puede leer/escribir en esa página/database

## Herramientas disponibles

- `notion_create_page` - Crear página en un database
- `notion_query_database` - Consultar un data source (database)
- `notion_search` - Buscar páginas y databases
- `notion_update_page` - Actualizar página existente

## Interacción natural con el agente

Puedes pedir operaciones en lenguaje natural:

```text
Crea una página en Notion con el título "Reunión semanal"
Busca en Notion páginas sobre "proyecto"
Consulta el database de tareas y muestra las pendientes
Actualiza el estado de la página XYZ a "Completado"
```

## Uso desde terminal

### Crear página en database

```bash
./picoclaw-agents agent -m "Usa notion_create_page con database_id='abc123', properties={'Name': {'title': [{'text': {'content': 'Mi Página'}}]}}"
```

### Consultar database

```bash
./picoclaw-agents agent -m "Usa notion_query_database con data_source_id='xyz789'"
```

### Consultar con filtro

```bash
./picoclaw-agents agent -m "Usa notion_query_database con data_source_id='xyz789', filter={'property': 'Status', 'select': {'equals': 'Active'}}"
```

### Buscar páginas

```bash
./picoclaw-agents agent -m "Usa notion_search con query='proyecto'"
```

### Actualizar página

```bash
./picoclaw-agents agent -m "Usa notion_update_page con page_id='page123', properties={'Status': {'select': {'name': 'Done'}}}"
```

## Uso desde Telegram

Con `picoclaw-agents gateway` activo, envía estos mensajes al bot:

### Crear página

```text
Crea una página en Notion en el database 'abc123' llamada "Nueva tarea"
```

### Consultar database

```text
Consulta el database de tareas en Notion
```

### Buscar

```text
Busca en Notion páginas sobre "reunión"
```

## Propiedades de Notion

Formatos comunes para propiedades:

| Tipo | Formato JSON |
|------|-------------|
| Title | `{"title": [{"text": {"content": "..."}}]}` |
| Rich text | `{"rich_text": [{"text": {"content": "..."}}]}` |
| Select | `{"select": {"name": "Opción"}}` |
| Multi-select | `{"multi_select": [{"name": "A"}, {"name": "B"}]}` |
| Date | `{"date": {"start": "2024-01-15"}}` |
| Checkbox | `{"checkbox": true}` |
| Number | `{"number": 42}` |
| URL | `{"url": "https://..."}` |
| Email | `{"email": "a@b.com"` |
| Relation | `{"relation": [{"id": "page_id"}]}` |

## API Basics

Todas las peticiones usan:
- Header `Authorization: Bearer TU_API_KEY`
- Header `Notion-Version: 2025-09-03`

### Doble ID en databases

Cada database tiene dos IDs:
- `database_id` - Para crear páginas (`parent: {"database_id": "..."}`)
- `data_source_id` - Para consultar (`POST /v1/data_sources/{id}/query`)

## Ejemplo completo: Crear página de tarea

```bash
./picoclaw-agents agent -m "Usa notion_create_page con:
  database_id='tu_database_id',
  properties={
    'Name': {'title': [{'text': {'content': 'Revisar código'}}]},
    'Status': {'select': {'name': 'Todo'}},
    'Date': {'date': {'start': '2024-01-20'}}
  }"
```

## Notas importantes

- Rate limit: ~3 peticiones por segundo
- Los IDs de página/database son UUIDs (con o sin guiones)
- La API no puede configurar filtros de vista (solo UI)
- Usa `is_inline: true` para databases embebidos en páginas

---

## ⚡ Comandos Slash Fast-path (v3.4.1+)

Usa comandos rápidos para operaciones instantáneas de Notion:

```text
/notion_create database=XYZ title="Notas de reunión"
/notion_query database=XYZ
/notion_search query="proyecto"
/notion_update page=ABC status="Completado"
```

**Beneficios:**
- ✅ **Latencia cero**: Sin razonamiento del LLM, ejecución instantánea
- ✅ **Sintaxis consistente**: Funciona idéntico en Telegram, Discord, CLI

### Global Tracker (v3.4.1+)

El **Global ImageGenTracker** es compartido entre todos los agentes:
- **Subagente crea/actualiza páginas de Notion** → **Agente Principal puede consultar inmediatamente**
- **Sin errores de "ID no encontrado"** entre límites de agentes

Ver [docs/queue_batch.es.md](docs/queue_batch.es.md) para documentación completa.
