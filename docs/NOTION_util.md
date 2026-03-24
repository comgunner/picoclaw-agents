# Notion Util

Quick guide to use Notion tools in PicoClaw from terminal and Telegram.

> **PicoClaw v3.4.1**: Now supports **Fast-path Slash Commands** for instant operations and **Global Tracker** for multi-agent consistency.

## Requirements

Configure your credential in `~/.picoclaw/config.json`:

```json
{
  "tools": {
    "notion": {
      "api_key": "YOUR_NOTION_API_KEY"
    }
  }
}
```

Or use environment variables:

```bash
export PICOCLAW_TOOLS_NOTION_API_KEY="your_notion_api_key"
```

### Getting Notion API Key

1. Go to https://notion.so/my-integrations
2. Create a new integration (+ New integration)
3. Copy the API key (starts with `ntn_` or `secret_`)
4. Save the key in `~/.config/notion/api_key` or in config.json

### Connecting Pages/Databases

1. Open the page or database you want to connect
2. Click "..." (more options)
3. Select "Connect to" → your integration
4. Now the integration can read/write to that page/database

## Available Tools

- `notion_create_page` - Create page in a database
- `notion_query_database` - Query a data source (database)
- `notion_search` - Search pages and databases
- `notion_update_page` - Update existing page

## Natural Agent Interaction

You can request operations in natural language:

```text
Create a page in Notion with title "Weekly meeting"
Search Notion for pages about "project"
Query the tasks database and show pending items
Update page XYZ status to "Completed"
```

## Terminal Usage

### Create page in database

```bash
./picoclaw agent -m "Use notion_create_page with database_id='abc123', properties={'Name': {'title': [{'text': {'content': 'My Page'}}]}}"
```

### Query database

```bash
./picoclaw agent -m "Use notion_query_database with data_source_id='xyz789'"
```

### Query with filter

```bash
./picoclaw agent -m "Use notion_query_database with data_source_id='xyz789', filter={'property': 'Status', 'select': {'equals': 'Active'}}"
```

### Search pages

```bash
./picoclaw agent -m "Use notion_search with query='project'"
```

### Update page

```bash
./picoclaw agent -m "Use notion_update_page with page_id='page123', properties={'Status': {'select': {'name': 'Done'}}}"
```

## Telegram Usage

With `picoclaw gateway` running, send these messages to the bot:

### Create page

```text
Create a page in Notion in database 'abc123' named "New task"
```

### Query database

```text
Query the tasks database in Notion
```

### Search

```text
Search Notion for pages about "meeting"
```

## Notion Properties

Common property formats:

| Type | JSON Format |
|------|-------------|
| Title | `{"title": [{"text": {"content": "..."}}]}` |
| Rich text | `{"rich_text": [{"text": {"content": "..."}}]}` |
| Select | `{"select": {"name": "Option"}}` |
| Multi-select | `{"multi_select": [{"name": "A"}, {"name": "B"}]}` |
| Date | `{"date": {"start": "2024-01-15"}}` |
| Checkbox | `{"checkbox": true}` |
| Number | `{"number": 42}` |
| URL | `{"url": "https://..."}` |
| Email | `{"email": "a@b.com"}` |
| Relation | `{"relation": [{"id": "page_id"}]}` |

## API Basics

All requests use:
- Header `Authorization: Bearer YOUR_API_KEY`
- Header `Notion-Version: 2025-09-03`

### Double ID in Databases

Each database has two IDs:
- `database_id` - For creating pages (`parent: {"database_id": "..."}`)
- `data_source_id` - For querying (`POST /v1/data_sources/{id}/query`)

## Complete Example: Create Task Page

```bash
./picoclaw agent -m "Use notion_create_page with:
  database_id='your_database_id',
  properties={
    'Name': {'title': [{'text': {'content': 'Review code'}}]},
    'Status': {'select': {'name': 'Todo'}},
    'Date': {'date': {'start': '2024-01-20'}}
  }"
```

## Important Notes

- Rate limit: ~3 requests per second
- Page/database IDs are UUIDs (with or without hyphens)
- The API cannot set view filters (UI only)
- Use `is_inline: true` for databases embedded in pages

---

## ⚡ Fast-path Slash Commands (v3.4.1+)

Use quick commands for instant Notion operations:

```text
/notion_create database=XYZ title="Meeting notes"
/notion_query database=XYZ
/notion_search query="project"
/notion_update page=ABC status="Done"
```

**Benefits:**
- ✅ **Zero latency**: No LLM reasoning, instant execution
- ✅ **Consistent syntax**: Works identically on Telegram, Discord, CLI

### Global Tracker (v3.4.1+)

The **Global ImageGenTracker** is shared across all agents:
- **Subagent creates/updates Notion pages** → **Main Agent can immediately query**
- **No "ID not found" errors** across agent boundaries

See [docs/queue_batch.md](docs/queue_batch.md) for complete documentation.
