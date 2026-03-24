# User Documentation: Queue & Batch System

PicoClaw's queue and batch processing system allows performing heavy background tasks without blocking the agent or unnecessarily consuming tokens during prolonged waits.

> **PicoClaw v3.4.2**: Features **Native Queue/Batch Skill** compiled directly into the binary for maximum performance and zero external dependencies.

## Included Tools

### 1. `batch_id`
Generates a unique identifier to reference tasks.
- **Usage**: `batch_id(prefix="SOCIAL")`
- **Result**: `#SOCIAL_02_03_26_1505`

### 2. `queue`
Allows listing and checking the status of ongoing processes.
- **Usage**:
    - `queue(action="list")`: Shows all active tasks.
    - `queue(task_id="#ID")`: Shows the progress of a specific task.

---

## Recommended Workflow

1. **Initiation**: The user or agent initiates a "Macro" tool (such as `social_post_bundle`).
2. **Decoupling**: The system immediately returns a tracking ID and releases the LLM.
3. **Tracking**: The user can check the status at any time using the ID.
4. **Completion**: Upon finishing, the system sends a direct notification with action buttons (**Approve**, **Redo**, **Publish**).

## User Benefits
- **Speed**: No need to wait for the LLM to process each intermediate step.
- **Organization**: Each creation has its own identification "plate" (#...).
- **Economy**: Saves thousands of tokens in automated tasks.

---

## Quick Commands (Slash Commands)

PicoClaw includes a **Fast-path** system that interceptes commands starting with `/` or `#` to execute them instantly without consulting the AI. This ensures immediate responses and data consistency.

### Social Bundle Commands
Used after receiving a task completion notification (e.g., `#IMA_GEN_02_03_26_1500`):

- `/bundle_approve id=ID`: Approves the batch and proceeds to publication/saving according to the workflow.
- `/bundle_regen id=ID`: Requests complete batch regeneration (image and text).
- `/bundle_edit id=ID`: Allows editing the batch text before approving.
- `/bundle_publish id=ID platforms=PLATFORMS`: Publishes the approved bundle to the specified platforms (e.g., `facebook,twitter`).

**⚠️ IMPORTANT: ID must exist in tracker**

Before publishing, verify the ID exists using:
```
/list pending
```
or query the tracker directly:
```
queue(action="list")
```

**Valid IDs** follow this format: `AAAAMMDD_HHMMSS_XXXXXX` (e.g., `20260302_161740_yiia22`)

### Utility Commands
- `/show [model|channel]`: Shows the active model or communication channel.
- `/list [models|channels|agents]`: Lists the options configured in the system.
- `/status`: Shows current token usage and system status.
- `/help`: Shows interactive help with all available commands.

### Binance Trading Commands (Tool-Based)

**Note:** Binance trading operations use tool functions directly (not fast-path commands yet):

| Tool Function | Description | Example |
|---------------|-------------|---------|
| `get_ticker_price` | Get crypto price | `get_ticker_price(symbol="BTCUSDT")` |
| `get_order_book` | Spot order book | `get_order_book(symbol="ETHUSDT", limit=10)` |
| `get_futures_order_book` | Futures order book | `get_futures_order_book(symbol="BTCUSDT")` |
| `list_futures_volume` | Volume ranking | `list_futures_volume(top=10)` |
| `get_spot_balance` | Spot balances | `get_spot_balance()` (requires API) |
| `get_futures_balance` | Futures balances | `get_futures_balance()` (requires API) |
| `open_futures_position` | Open position | `open_futures_position(symbol="BTCUSDT", side="LONG", quantity="0.001", leverage=5, confirm=true)` |
| `close_futures_position` | Close position | `close_futures_position(symbol="BTCUSDT", confirm=true)` |

**Fast-path commands for Binance** (e.g., `/binance_price`, `/binance_open`) are documented for future implementation.

### Supported Channels
These commands work identically across:
1. **Telegram**: Via the command menu (the `/` icon).
2. **Discord**: As native application commands (Slash Commands).
3. **Terminal (CLI)**: By typing them directly in interactive mode (`./picoclaw agent -m ""`).

---

## Architecture: Global Tracker
To ensure commands work in multi-agent installations, PicoClaw uses a unique, shared `ImageGenTracker`. This allows that if a **Subagent** generates an image, the **Main Agent** (Project Manager) can immediately approve it using its ID, without "ID not found" errors.

## Developer Guide (Go/Native Skills)

The Queue/Batch system is now a **Native Skill** compiled directly into the PicoClaw binary. This provides:

- **Zero External Dependencies**: No need for external `.md` files at runtime
- **Maximum Performance**: Documentation strings are embedded in the binary
- **Enhanced Security**: Skill cannot be modified or tampered with externally
- **Automatic Updates**: Skill updates with each new PicoClaw release

### Integration Architecture

The native skill is implemented in `pkg/skills/queue_batch.go` and registered in:
1. `pkg/skills/loader.go` - Native skills registry
2. `pkg/agent/context.go` - System prompt injection

To integrate new native tools:
1. Create `pkg/skills/{name}.go` with skill struct and documentation constants
2. Register in `loader.go` native skills registry
3. Inject via `context.go` BuildSystemPrompt() method
4. Use `tools.GetGlobalQueueManager()` for task management
5. Execute heavy logic in separate goroutine
6. Notify via `MessageBus` on completion

**See:** `local_work/crear_skill_interna.md` for complete native skill development guide.

---

## Troubleshooting

### Error: "ID not found in tracker"

**Symptom:**
```
/bundle_publish id=20260302_163848_qqqia2 platforms=facebook
❌ Error: Imagen con ID 20260302_163848_qqqia2 no encontrada en el tracker
```

**Possible Causes:**

1. **Incorrect ID or typo**: The ID doesn't match any record in the tracker
2. **Task already published**: The ID was archived after successful publication
3. **Task expired**: The tracker automatically cleans old records
4. **Different session**: The ID belongs to another session/PicoClaw instance

**Solution:**

**Step 1: Verify existing IDs**
```bash
# List all pending tasks
/list pending

# Or query the queue
queue(action="list")
```

**Step 2: Identify the correct ID**
Valid IDs follow this format:
```
AAAAMMDD_HHMMSS_XXXXXX
││││││││ ││││││ ││││││
││││││││ ││││││ └─── Random (6 chars)
││││││││ │││││└───── Seconds (163848 = 16:38:48)
││││││││ ││││└────── Minutes (38)
││││││││ ││└──────── Hours (16)
││││││││ └────────── Day (02)
││││││└───────────── Month (03)
││││└─────────────── Year (2026)
```

**Step 3: Use the correct ID**
```bash
# ✅ Correct (ID exists in tracker)
/bundle_publish id=20260302_161740_yiia22 platforms=facebook

# ❌ Incorrect (ID doesn't exist)
/bundle_publish id=20260302_163848_qqqia2 platforms=facebook
```

**Prevention:**

1. **Copy and paste IDs** from the original notification instead of typing manually
2. **Use autocomplete** from Telegram/Discord when available
3. **Verify before publishing** with `/list pending`
4. **Save important IDs** in a safe place until publication is complete

**Example of Correct Workflow:**

```
1. Agent generates image:
   Bot: "✅ Image generated. ID: 20260302_161740_yiia22"

2. User approves:
   User: "/bundle_approve id=20260302_161740_yiia22"
   Bot: "✅ Bundle approved"

3. User publishes:
   User: "/bundle_publish id=20260302_161740_yiia22 platforms=facebook"
   Bot: "🚀 Publication completed"
```
