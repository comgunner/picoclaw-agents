# Plan de Fix: WebUI "Not Configured" para Modelos Modelfile de Ollama

**Fecha:** 2026-04-05
**Problema:** Modelos Ollama personalizados sin prefijo `ollama/` aparecen como "Not Configured" en la WebUI
**Severidad:** Media (solo UI visual, el agente funciona correctamente)

---

## Diagnóstico

### Flujo actual (roto)

```
config.json → model: "picoclaw-qwen25-min" (sin "/")
    ↓
splitModel("picoclaw-qwen25-min") → ("openai", "picoclaw-qwen25-min")
    ↓
modelProtocol → "openai"
    ↓
probeOpenAICompatibleModel("http://localhost:11434/v1", "picoclaw-qwen25-min", "ollama")
    ↓
GET http://localhost:11434/v1/models
    ↓
Ollama retorna: "qwen2.5:0.5b" (nombre base, NO "picoclaw-qwen25-min")
    ↓
No match → probe falla → isModelConfigured = false → WebUI: "Not Configured" ❌
```

### Flujo esperado

```
config.json → model: "picoclaw-qwen25-min", api_base: "http://localhost:11434/v1", auth_method: "local"
    ↓
Detectar Ollama por api_base (localhost:11434)
    ↓
probeOllamaModel("http://localhost:11434", "picoclaw-qwen25-min")
    ↓
GET http://localhost:11434/api/tags → retorna "picoclaw-qwen25-min:latest"
    ↓
Match encontrado → isModelConfigured = true → WebUI: "Configured" ✅
```

---

## Archivos a Modificar

### 1. `web/backend/api/model_status.go`

**Función a modificar:** `probeLocalModelAvailability()`

#### Cambio propuesto

**Antes (línea ~82):**
```go
func probeLocalModelAvailability(m *config.ModelConfig) bool {
    apiBase := modelProbeAPIBase(m)
    protocol, modelID := splitModel(m.Model)
    switch protocol {
    case "ollama":
        return probeOllamaModel(apiBase, modelID)
    case "vllm":
        return probeOpenAICompatibleModelFunc(apiBase, modelID, m.APIKey)
    case "github-copilot", "copilot":
        return probeTCPServiceFunc(apiBase)
    case "claude-cli", "claudecli", "codex-cli", "codexcli":
        return true
    default:
        if hasLocalAPIBase(apiBase) {
            return probeOpenAICompatibleModelFunc(apiBase, modelID, m.APIKey)
        }
        return false
    }
}
```

**Después:**
```go
func probeLocalModelAvailability(m *config.ModelConfig) bool {
    apiBase := modelProbeAPIBase(m)
    protocol, modelID := splitModel(m.Model)
    switch protocol {
    case "ollama":
        return probeOllamaModel(apiBase, modelID)
    case "vllm":
        return probeOpenAICompatibleModelFunc(apiBase, modelID, m.APIKey)
    case "github-copilot", "copilot":
        return probeTCPServiceFunc(apiBase)
    case "claude-cli", "claudecli", "codex-cli", "codexcli":
        return true
    default:
        // NEW: Check if this is an Ollama local model without the ollama/ prefix
        // When auth_method is "local" and api_base points to Ollama (port 11434),
        // use the Ollama probe even if the model name doesn't have "ollama/" prefix.
        if isOllamaAPIBase(apiBase) {
            return probeOllamaModel(apiBase, modelID)
        }
        if hasLocalAPIBase(apiBase) {
            return probeOpenAICompatibleModelFunc(apiBase, modelID, m.APIKey)
        }
        return false
    }
}
```

**Nueva función auxiliar a agregar:**
```go
// isOllamaAPIBase checks if the api_base points to a local Ollama instance.
func isOllamaAPIBase(apiBase string) bool {
    u, err := url.Parse(strings.TrimSpace(apiBase))
    if err != nil || u.Hostname() == "" {
        return false
    }
    port := u.Port()
    // Ollama default port is 11434
    return (strings.ToLower(u.Hostname()) == "localhost" ||
        u.Hostname() == "127.0.0.1" ||
        u.Hostname() == "::1" ||
        u.Hostname() == "0.0.0.0") &&
        (port == "11434" || port == "")
}
```

### 2. `web/backend/api/model_status.go` — Función `requiresRuntimeProbe()`

**Cambio menor (opcional pero recomendado):**

Agregar `auth_method == "local"` como condición explícita para usar probe de Ollama cuando el `api_base` apunta a localhost:11434.

**No es estrictamente necesario** si el fix en `probeLocalModelAvailability` funciona, pero hace el código más explícito.

---

## Impacto del Cambio

### Modelos que se benefician
Cualquier modelo Modelfile de Ollama sin prefijo `ollama/`:
- `picoclaw-qwen25-min`
- `picoclaw-qwen3-tiny`
- `picoclaw-gemma4-8gb`
- `picoclaw-coder-minimal`
- `mi-modelo-personalizado`

### Modelos que NO se ven afectados
- Modelos con prefijo: `ollama/qwen2.5:0.5b` → ya funcionan
- Modelos remotos: `openai/gpt-4`, `anthropic/claude-sonnet-4-6` → no usan probe local
- Modelos vllm: `vllm/mi-modelo` → siguen funcionando igual

---

## Plan de Implementación

### Paso 1: Agregar función `isOllamaAPIBase()`

En `web/backend/api/model_status.go`, después de `hasLocalAPIBase()`:

```go
// isOllamaAPIBase checks if the api_base points to a local Ollama instance
// (port 11434 on localhost).
func isOllamaAPIBase(apiBase string) bool {
    apiBase = strings.TrimSpace(apiBase)
    if apiBase == "" {
        return false
    }
    u, err := url.Parse(apiBase)
    if err != nil || u.Hostname() == "" {
        return false
    }
    port := u.Port()
    return (strings.ToLower(u.Hostname()) == "localhost" ||
        u.Hostname() == "127.0.0.1" ||
        u.Hostname() == "::1" ||
        u.Hostname() == "0.0.0.0") &&
        (port == "11434" || port == "")
}
```

### Paso 2: Modificar `probeLocalModelAvailability()`

En el bloque `default:` de la función, agregar la verificación de Ollama antes del probe genérico:

```go
default:
    // Ollama Modelfile models without ollama/ prefix
    if isOllamaAPIBase(apiBase) {
        return probeOllamaModel(apiBase, modelID)
    }
    if hasLocalAPIBase(apiBase) {
        return probeOpenAICompatibleModelFunc(apiBase, modelID, m.APIKey)
    }
    return false
```

### Paso 3: Agregar test

En `web/backend/api/model_status_test.go`:

```go
func TestIsOllamaAPIBase(t *testing.T) {
    tests := []struct {
        name     string
        apiBase  string
        expected bool
    }{
        {"localhost port 11434", "http://localhost:11434/v1", true},
        {"localhost port 11434 no v1", "http://localhost:11434", true},
        {"127.0.0.1 port 11434", "http://127.0.0.1:11434/v1", true},
        {"0.0.0.0 port 11434", "http://0.0.0.0:11434/v1", true},
        {"different port", "http://localhost:8000/v1", false},
        {"remote host", "http://remote-server:11434/v1", false},
        {"empty", "", false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := isOllamaAPIBase(tt.apiBase)
            if result != tt.expected {
                t.Errorf("isOllamaAPIBase(%q) = %v, want %v", tt.apiBase, result, tt.expected)
            }
        })
    }
}
```

### Paso 4: Verificación

```bash
# 1. Build
go build ./web/backend/

# 2. Tests
go test ./web/backend/api/ -run "Ollama|ModelStatus|Configured" -v

# 3. Manual test
# Agregar modelo Modelfile al config.json y verificar que la WebUI muestra "Configured"
```

---

## Riesgos y Mitigación

| Riesgo | Probabilidad | Impacto | Mitigación |
|--------|-------------|---------|------------|
| Detección falsa de Ollama | Baja | Medio | `isOllamaAPIBase` es estricto: requiere localhost + port 11434 |
| Rompe modelos OpenAI locales | Baja | Alto | Solo afecta el `default` case; OpenAI usa `openai/` prefix → entra en otro branch |
| Ollama en puerto custom | Media | Bajo | Si Ollama corre en otro puerto, el usuario debe usar `ollama/` prefix |

---

## Líneas de Código Afectadas

| Archivo | Líneas a agregar | Líneas a modificar |
|---------|-----------------|-------------------|
| `web/backend/api/model_status.go` | ~15 (nueva función + mod switch) | ~4 (bloque default del switch) |
| `web/backend/api/model_status_test.go` | ~25 (nuevos tests) | 0 |
| **Total** | **~40** | **~4** |

---

*fix_webui_modelfiles.md — 2026-04-05*
