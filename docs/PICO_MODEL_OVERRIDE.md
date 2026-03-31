# Pico Model Override Protocol

**Version:** 1.0.0  
**Created:** March 30, 2026  
**Component:** Web UI → Pico Channel → Agent Loop

---

## 📋 Overview

The Pico protocol supports client-specified model selection, allowing Web UI users to override the agent's configured default model on a per-message basis.

**Key Features:**
- ✅ **Dynamic Model Switching** - Change models without restarting agents
- ✅ **Per-User Isolation** - Each user can use different models independently
- ✅ **Zero Configuration** - Works automatically with Web UI model selector
- ✅ **Backward Compatible** - Falls back to agent's default if no model specified

---

## 🔄 Data Flow

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│  Web UI         │     │  Frontend       │     │  Backend        │     │  Agent Loop     │
│  ModelSelector  │     │  controller.ts  │     │  pico.go        │     │  loop.go        │
└────────┬────────┘     └────────┬────────┘     └────────┬────────┘     └────────┬────────┘
         │                       │                       │                       │
         │ Select model          │                       │                       │
         │ "gpt-4"               │                       │                       │
         ├──────────────────────>│                       │                       │
         │                       │                       │                       │
         │                       │ Send message          │                       │
         │                       │ {content, model_name} │                       │
         │                       ├──────────────────────>│                       │
         │                       │                       │                       │
         │                       │                       │ Extract model_name    │
         │                       │                       │ Add to metadata       │
         │                       │                       ├──────────────────────>│
         │                       │                       │                       │
         │                       │                       │                       │ Use model from
         │                       │                       │                       │ metadata if present
         │                       │                       │                       │
```

---

## 🔧 Implementation Details

### 1. Frontend (TypeScript)

**File:** `web/frontend/src/features/chat/controller.ts`

```typescript
export function sendChatMessage(content: string, modelName?: string) {
  socket.send(
    JSON.stringify({
      type: "message.send",
      id,
      payload: {
        content,
        model_name: modelName,  // ← Client-specified model
      },
    }),
  )
}
```

**Behavior:**
- If `modelName` is provided, it's included in the WebSocket payload
- If `modelName` is `undefined`, the field is omitted (backward compatible)

---

### 2. Backend (Go)

**File:** `pkg/channels/pico.go`

```go
func (c *PicoChannel) picoHandleMessageSend(pc *picoConn, msg PicoMessage) {
  content, _ := msg.Payload["content"].(string)
  modelName, _ := msg.Payload["model_name"].(string)

  metadata := map[string]string{
    "platform":   "pico",
    "session_id": sessionID,
    "conn_id":    pc.id,
    "peer_kind":  "direct",
    "peer_id":    sessionID,
  }

  // Include model_name in metadata if provided by client
  if modelName != "" {
    metadata["model_name"] = modelName
  }

  c.HandleMessage(senderID, chatID, content, nil, metadata)
}
```

**Behavior:**
- Extracts `model_name` from payload
- Adds to metadata map only if non-empty
- Logs model_name for debugging

---

### 3. Agent Loop (Go)

**File:** `pkg/agent/loop.go`

```go
// Determine model to use: client-specified model from metadata takes priority
modelToUse := agent.Model
if clientModel, ok := opts.Metadata["model_name"]; ok && clientModel != "" {
  modelToUse = clientModel
  logger.DebugCF("agent", "Using client-specified model", map[string]any{
    "model":         modelToUse,
    "agent_id":      agent.ID,
    "session_key":   opts.SessionKey,
  })
}

// Use modelToUse in LLM call
return agent.Provider.Chat(ctx, prunedMessages, providerToolDefs, modelToUse, ...)
```

**Behavior:**
- Checks metadata for `model_name`
- Uses client-specified model if present
- Falls back to agent's configured model otherwise
- Logs debug information

---

## 📝 Usage Examples

### Web UI Model Selection

1. Open Web UI: `http://localhost:18800/`
2. Click model selector in header
3. Choose desired model (e.g., `gpt-4`, `claude-3-5-sonnet`, `gemini-pro`)
4. Send message - model change takes effect immediately

### Programmatic Usage (WebSocket)

```javascript
const ws = new WebSocket('ws://localhost:18800/pico/ws?session_id=abc-123')

ws.send(JSON.stringify({
  type: "message.send",
  id: "msg-123",
  payload: {
    content: "Hello, world!",
    model_name: "anthropic/claude-3-5-sonnet"  // Override model
  }
}))
```

---

## 🔍 Debugging

### Enable Debug Logs

```bash
# Set log level to debug
export PICOCLAW_LOG_LEVEL=debug
./build/picoclaw-agents gateway
```

### Expected Log Output

**When model override is active:**
```
[DEBUG] pico: Received message {session_id=abc-123, preview=Hello world, model_name=claude-3-5-sonnet}
[DEBUG] agent: Using client-specified model {model=claude-3-5-sonnet, agent_id=engineering_manager}
```

**When using default model:**
```
[DEBUG] pico: Received message {session_id=abc-123, preview=Hello world, model_name=}
[INFO]  agent: Processing message from pico:user: Hello world
```

---

## ⚠️ Important Notes

### Model Name Format

Model names must match the format in `~/.picoclaw/config.json`:

```json
{
  "model_list": [
    {
      "model_name": "gpt-4",
      "model": "openai/gpt-4"
    },
    {
      "model_name": "claude-3-5-sonnet",
      "model": "anthropic/claude-3-5-sonnet"
    }
  ]
}
```

**Supported formats:**
- ✅ `gpt-4` (model_name alias)
- ✅ `openai/gpt-4` (full model path)
- ✅ `anthropic/claude-3-5-sonnet` (provider-prefixed)

**Invalid formats:**
- ❌ `openrouter-free` (deprecated, use `openrouter/auto`)
- ❌ `free` (ambiguous, use `openrouter/auto`)

---

### Priority Order

Model resolution follows this priority:

1. **Client-specified model** (from Web UI or WebSocket payload)
2. **Agent's configured model** (from `config.json`)
3. **Default model** (from agent defaults)

---

### Backward Compatibility

The protocol is fully backward compatible:

- **Old clients** (no `model_name` field) → Agent uses configured default
- **New clients** (with `model_name` field) → Agent uses client-specified model

No breaking changes to existing integrations.

---

## 🧪 Testing

### Manual Test Steps

1. **Setup:**
   ```bash
   cd picoclaw-agents
   make build
   ./build/picoclaw-agents gateway
   ```

2. **Open Web UI:**
   - Navigate to `http://localhost:18800/`
   - Verify model selector shows available models

3. **Test Model Override:**
   - Select `gpt-4` from model selector
   - Send message: "What model are you using?"
   - Check logs for: `Using client-specified model {model=gpt-4}`

4. **Test Fallback:**
   - Select different model (e.g., `claude-3-5-sonnet`)
   - Send another message
   - Verify logs show new model being used

5. **Verify Isolation:**
   - Open Web UI in incognito window
   - Select different model
   - Send message
   - Confirm first session still uses its selected model

---

## 📊 Performance Impact

| Metric | Impact | Notes |
|--------|--------|-------|
| **Latency** | None | Model resolution is O(1) map lookup |
| **Memory** | None | No additional allocations |
| **CPU** | Negligible | Single string comparison |
| **Compatibility** | 100% | Fully backward compatible |

---

## 🔗 Related Files

| Component | File | Purpose |
|-----------|------|---------|
| **Frontend** | `web/frontend/src/features/chat/controller.ts` | Send model_name in payload |
| **Frontend** | `web/frontend/src/components/chat/chat-page.tsx` | Pass model to sendMessage |
| **Frontend** | `web/frontend/src/hooks/use-chat-models.ts` | Model selection state |
| **Backend** | `pkg/channels/pico.go` | Extract model_name from payload |
| **Agent** | `pkg/agent/loop.go` | Use model from metadata |
| **Agent** | `pkg/agent/instance.go` | Model resolution (fallback) |

---

## 🐛 Troubleshooting

### Issue: Model change doesn't take effect

**Symptoms:**
- Web UI shows model changed
- Logs show agent using different model

**Causes:**
1. **Frontend not sending model_name**
   - Check browser console for errors
   - Verify WebSocket payload includes `model_name`

2. **Backend not extracting model_name**
   - Check debug logs: `model_name` should appear in "Received message"
   - Verify `pico.go` changes are deployed

3. **Agent not using metadata**
   - Check debug logs: "Using client-specified model" should appear
   - Verify `loop.go` changes are deployed

**Solution:**
```bash
# Rebuild and restart
make build
./build/picoclaw-agents gateway

# Check logs
tail -f logs/picoclaw.log | grep "model_name"
```

---

### Issue: 404 Model Not Found

**Symptoms:**
- Error: `API request failed: Status: 404, model 'xxx' not found`

**Causes:**
1. **Invalid model name**
   - Model name doesn't exist in provider
   - Typo in model name

2. **Model not configured**
   - Model not in `model_list` in config.json
   - API key missing for provider

**Solution:**
```bash
# Check available models
./build/picoclaw-agents models list

# Verify config
cat ~/.picoclaw/config.json | jq '.model_list'
```

---

### Issue: 402 Payment Required (OpenRouter)

**Symptoms:**
- Error: `This request requires more credits`

**Causes:**
- Switching to premium model without sufficient credits
- Model selector set to paid model

**Solution:**
1. Select free model (e.g., `openrouter/auto`)
2. Or add credits to OpenRouter account
3. Or use different provider

---

## 📚 References

- **Original Issue:** `local_work/problemas_encontrados.md` - "CRÍTICO: Selección de modelo en Web UI no se aplica"
- **Pico Protocol:** `pkg/channels/pico.go` - WebSocket message format
- **Agent Loop:** `pkg/agent/loop.go` - Message processing logic
- **Model Resolution:** `pkg/agent/instance.go` - Model alias resolution

---

*Pico Model Override Protocol - Dynamic model selection for PicoClaw Web UI*
