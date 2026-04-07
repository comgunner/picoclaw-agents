# Plan de Migración: n8n-mcp → Go-Native (PicoClaw-Agents)

**Fecha:** 2026-04-06
**Origen:** `/Volumes/UPLOAD/AGENTE_BKP/scripts_gunner/GIT_CLONE/MCP/n8n-mcp/` (TypeScript, ~200 archivos, ~20 herramientas MCP)
**Destino:** `pkg/mcp/n8n/` dentro de PicoClaw-Agents (Go puro)
**Objetivo:** Implementar el mismo set de herramientas MCP de n8n como herramienta nativa de Go, compilada directamente en el binario de picoclaw-agents.

---

## 1. Resumen del Proyecto Original (n8n-mcp)

### Qué hace
Proporciona un servidor MCP que expone herramientas para que agentes de IA gestionen instancias de n8n via su API REST:
- Crear, leer, actualizar, eliminar workflows
- Validar workflows (estructura, conexiones, expresiones)
- Ejecutar y testear workflows
- Buscar y documentar nodos de n8n
- Buscar y deployar templates
- Auditoría de seguridad de instancia
- Gestión de credenciales
- Gestión de datos (datatables)
- Generación de workflows con IA
- Auto-fix de workflows

### Stack actual
| Tecnología | Uso |
|-----------|-----|
| TypeScript 5.x | Lenguaje principal |
| @modelcontextprotocol/sdk | Servidor MCP (stdio + HTTP) |
| SQLite (better-sqlite3) | Base de datos de nodos |
| zod | Validación de esquemas |
| n8n REST API | Backend que gestiona |
| ~200 archivos TypeScript | Código fuente |
| ~20 herramientas MCP | Superficie de API |

### Modo de ejecución
- **Stdio** (por defecto): MCP server via stdin/stdout
- **HTTP**: Servidor HTTP con autenticación por token
- Soporta sesiones (single-session)

---

## 2. Herramientas MCP a Portar (20 herramientas)

### Grupo A: Gestión de Workflows (7 herramientas)

| # | Herramienta | Descripción | Complejidad |
|---|------------|-------------|-------------|
| 1 | `n8n_create_workflow` | Crear workflow (name, nodes[], connections{}) | **Alta** |
| 2 | `n8n_get_workflow` | Obtener workflow por ID | **Baja** |
| 3 | `n8n_update_full_workflow` | Reemplazar workflow completo | **Alta** |
| 4 | `n8n_update_partial_workflow` | Update diff de nodos/propiedades | **Muy Alta** |
| 5 | `n8n_delete_workflow` | Eliminar workflow | **Baja** |
| 6 | `n8n_list_workflows` | Listar workflows con filtros | **Media** |
| 7 | `n8n_validate_workflow` | Validar estructura + expresiones | **Muy Alta** |

### Grupo B: Ejecución y Testing (3 herramientas)

| # | Herramienta | Descripción | Complejidad |
|---|------------|-------------|-------------|
| 8 | `n8n_test_workflow` | Ejecutar workflow de prueba | **Alta** |
| 9 | `n8n_executions` | Listar/buscar ejecuciones | **Media** |
| 10 | `n8n_autofix_workflow` | Auto-corregir errores | **Muy Alta** |

### Grupo C: Nodos y Documentación (3 herramientas)

| # | Herramienta | Descripción | Complejidad |
|---|------------|-------------|-------------|
| 11 | `search_nodes` | Buscar nodos en DB local (SQLite) | **Media** |
| 12 | `get_node` | Info de nodo con múltiples modos | **Alta** |
| 13 | `tools_documentation` | Documentación de herramientas | **Baja** |

### Grupo D: Templates (2 herramientas)

| # | Herramienta | Descripción | Complejidad |
|---|------------|-------------|-------------|
| 14 | `search_templates` | Buscar templates de workflow | **Media** |
| 15 | `get_template` | Obtener template por ID | **Baja** |

### Grupo E: Sistema y Auditoría (4 herramientas)

| # | Herramienta | Descripción | Complejidad |
|---|------------|-------------|-------------|
| 16 | `n8n_health_check` | Health check de instancia n8n | **Baja** |
| 17 | `n8n_audit_instance` | Auditoría de seguridad completa | **Muy Alta** |
| 18 | `n8n_manage_credentials` | CRUD de credenciales en n8n | **Alta** |
| 19 | `n8n_manage_datatable` | CRUD de tablas de datos | **Alta** |
| 20 | `n8n_generate_workflow` | Generar workflow desde descripción | **Alta** |

### Grupo F: Versioning y Deploy (2 herramientas)

| # | Herramienta | Descripción | Complejidad |
|---|------------|-------------|-------------|
| 21 | `n8n_workflow_versions` | Gestión de versiones | **Media** |
| 22 | `n8n_deploy_template` | Deploy template como workflow | **Media** |

---

## 3. Configuración Requerida

### Variables de Entorno (o config.json)

```json
{
  "tools": {
    "n8n": {
      "enabled": false,
      "api_url": "http://localhost:5678",
      "api_key": "eyJhbGciOi...",
      "timeout": 30000,
      "max_retries": 3,
      "database_path": "~/.picoclaw/data/nodes.db"
    }
  }
}
```

### Campos en `pkg/config/config.go` (agregar):

```go
type N8nConfig struct {
    Enabled      bool   `json:"enabled"`
    APIURL       string `json:"api_url"`
    APIKey       string `json:"api_key"`
    Timeout      int    `json:"timeout"`
    MaxRetries   int    `json:"max_retries"`
    DatabasePath string `json:"database_path"`
}
```

---

## 4. Archivos a Crear (Go-Native)

### Estructura propuesta

```
pkg/mcp/n8n/
├── client.go              # N8nApiClient (HTTP client con retry, auth)
├── client_test.go         # Tests del cliente HTTP
├── types.go               # Tipos: Workflow, WorkflowNode, Execution, etc.
├── tools.go               # Tool definitions (20 herramientas)
├── handler_workflows.go   # Handlers: create, get, update, delete, list
├── handler_validation.go  # Handler: validate_workflow
├── handler_execution.go   # Handlers: test_workflow, executions
├── handler_autofix.go     # Handler: autofix_workflow
├── handler_nodes.go       # Handlers: search_nodes, get_node
├── handler_templates.go   # Handlers: search_templates, get_template
├── handler_system.go      # Handlers: health_check, audit_instance
├── handler_credentials.go # Handler: manage_credentials
├── handler_datatable.go   # Handler: manage_datatable
├── handler_generate.go    # Handler: generate_workflow
├── handler_versions.go    # Handler: workflow_versions
├── handler_deploy.go      # Handler: deploy_template
├── validation.go          # WorkflowValidator (estructura + conexiones)
├── security_scanner.go    # SecurityScanner (credenciales expuestas, etc.)
├── node_db.go             # SQLite adapter para búsqueda de nodos
├── registrar.go           # RegisterN8nTools() — integra con ToolRegistry
└── registrar_test.go      # Tests de integración
```

### Archivos a modificar (proyecto existente)

| Archivo | Cambio |
|---------|--------|
| `pkg/config/config.go` | Agregar `N8nConfig` struct y campo en `Tools` |
| `pkg/config/defaults.go` | Default `N8nConfig{Enabled: false}` |
| `pkg/agent/instance.go` | Agregar `toolsRegistry.Register(tools.NewN8nTool(...))` si `cfg.Tools.N8n.Enabled` |
| `cmd/picoclaw/internal/util/command.go` | Agregar `newN8nMCPServerCommand()` |
| `config/config.example.json` | Agregar sección `tools.n8n` |

---

## 5. Dependencias Externas Necesarias

### Nuevas dependencias en `go.mod`

| Paquete | Uso | Justificación |
|---------|-----|---------------|
| `modernc.org/sqlite` | **Ya incluida** | Node database (ya usada por seahorse) |
| `github.com/go-resty/resty/v2` | **Ya incluida** | HTTP client con retry para n8n API |

**Resultado: CERO dependencias nuevas.** Todo lo necesario ya está en el proyecto.

---

## 6. Implementación Detallada por Fase

### Fase 1: Cliente HTTP + Tipos (8sp)

**Archivos:** `client.go`, `client_test.go`, `types.go`

```go
// pkg/mcp/n8n/client.go
package n8n

import (
    "context"
    "fmt"
    "net/http"
    "time"

    "github.com/go-resty/resty/v2"
)

type N8nClient struct {
    baseURL    string
    apiKey     string
    httpClient *resty.Client
}

func NewN8nClient(baseURL, apiKey string, timeout time.Duration, maxRetries int) *N8nClient {
    client := resty.New().
        SetBaseURL(baseURL).
        SetHeader("X-N8N-API-KEY", apiKey).
        SetHeader("Content-Type", "application/json").
        SetTimeout(timeout).
        SetRetryCount(maxRetries).
        SetRetryWaitTime(1 * time.Second).
        SetRetryMaxWaitTime(5 * time.Second)
    return &N8nClient{baseURL: baseURL, apiKey: apiKey, httpClient: client}
}

func (c *N8nClient) Get(ctx context.Context, path string, result any) error {
    resp, err := c.httpClient.R().SetContext(ctx).SetResult(result).Get(path)
    if err != nil {
        return fmt.Errorf("n8n GET %s: %w", path, err)
    }
    if resp.StatusCode() >= 400 {
        return fmt.Errorf("n8n GET %s: status %d: %s", path, resp.StatusCode(), resp.Body())
    }
    return nil
}

func (c *N8nClient) Post(ctx context.Context, path string, body, result any) error {
    resp, err := c.httpClient.R().SetContext(ctx).SetBody(body).SetResult(result).Post(path)
    if err != nil {
        return fmt.Errorf("n8n POST %s: %w", path, err)
    }
    if resp.StatusCode() >= 400 {
        return fmt.Errorf("n8n POST %s: status %d: %s", path, resp.StatusCode(), resp.Body())
    }
    return nil
}

func (c *N8nClient) Delete(ctx context.Context, path string) error {
    resp, err := c.httpClient.R().SetContext(ctx).Delete(path)
    if err != nil {
        return fmt.Errorf("n8n DELETE %s: %w", path, err)
    }
    if resp.StatusCode() >= 400 {
        return fmt.Errorf("n8n DELETE %s: status %d: %s", path, resp.StatusCode(), resp.Body())
    }
    return nil
}

func (c *N8nClient) Put(ctx context.Context, path string, body, result any) error {
    resp, err := c.httpClient.R().SetContext(ctx).SetBody(body).SetResult(result).Put(path)
    if err != nil {
        return fmt.Errorf("n8n PUT %s: %w", path, err)
    }
    if resp.StatusCode() >= 400 {
        return fmt.Errorf("n8n PUT %s: status %d: %s", path, resp.StatusCode(), resp.Body())
    }
    return nil
}

func (c *N8nClient) HealthCheck(ctx context.Context) (map[string]any, error) {
    var result map[string]any
    if err := c.Get(ctx, "/api/health", &result); err != nil {
        return nil, err
    }
    return result, nil
}
```

```go
// pkg/mcp/n8n/types.go
package n8n

type Workflow struct {
    ID         string               `json:"id,omitempty"`
    Name       string               `json:"name"`
    Nodes      []WorkflowNode       `json:"nodes"`
    Connections map[string]any       `json:"connections"`
    Settings   *WorkflowSettings    `json:"settings,omitempty"`
    Active     bool                 `json:"active"`
    CreatedAt  string               `json:"createdAt,omitempty"`
    UpdatedAt  string               `json:"updatedAt,omitempty"`
}

type WorkflowNode struct {
    ID          string `json:"id"`
    Name        string `json:"name"`
    Type        string `json:"type"`
    TypeVersion int    `json:"typeVersion"`
    Position    []int  `json:"position"`
    Parameters  map[string]any `json:"parameters,omitempty"`
    Credentials map[string]any `json:"credentials,omitempty"`
    Disabled    bool   `json:"disabled,omitempty"`
}

type WorkflowSettings struct {
    ExecutionOrder             string `json:"executionOrder,omitempty"`
    Timezone                   string `json:"timezone,omitempty"`
    SaveDataErrorExecution     string `json:"saveDataErrorExecution,omitempty"`
    SaveDataSuccessExecution   string `json:"saveDataSuccessExecution,omitempty"`
    SaveManualExecutions       bool   `json:"saveManualExecutions,omitempty"`
    SaveExecutionProgress      bool   `json:"saveExecutionProgress,omitempty"`
    ExecutionTimeout           int    `json:"executionTimeout,omitempty"`
    ErrorWorkflow              string `json:"errorWorkflow,omitempty"`
}

type Execution struct {
    ID          string    `json:"id"`
    WorkflowID  string    `json:"workflowId"`
    Mode        string    `json:"mode"`
    Status      string    `json:"status"`
    StartedAt   string    `json:"startedAt"`
    StoppedAt   string    `json:"stoppedAt,omitempty"`
    Error       *any      `json:"error,omitempty"`
}

type Credential struct {
    ID   string `json:"id,omitempty"`
    Name string `json:"name"`
    Type string `json:"type"`
    Data map[string]any `json:"data,omitempty"`
}
```

### Fase 2: Workflow CRUD (8sp)

**Archivos:** `handler_workflows.go`, `handler_workflows_test.go`

Implementar los handlers que usan `N8nClient`:
- `CreateWorkflow` → POST /api/workflows
- `GetWorkflow` → GET /api/workflows/{id}
- `UpdateFullWorkflow` → PUT /api/workflows/{id}
- `DeleteWorkflow` → DELETE /api/workflows/{id}
- `ListWorkflows` → GET /api/workflows (con query params)

### Fase 3: Validación (8sp)

**Archivos:** `handler_validation.go`, `validation.go`, `security_scanner.go`

Portar la lógica de:
- `src/services/workflow-validator.ts` → `validation.go`
- `src/services/workflow-security-scanner.ts` → `security_scanner.go`
- `src/services/enhanced-config-validator.ts` → parte de `validation.go`

Validaciones a portar:
- Estructura de nodos (campos requeridos)
- Conexiones válidas (origen/destino existen)
- Expresiones n8n válidas (`{{ $json.xxx }}`)
- Credenciales referenciadas existen
- Nodes sin ciclos infinitos
- Webhook nodes con URL configurada

### Fase 4: Ejecución y Auto-fix (8sp)

**Archivos:** `handler_execution.go`, `handler_autofix.go`

- `TestWorkflow` → POST /api/workflows/{id}/execute
- `ListExecutions` → GET /api/executions (con filtros)
- `AutofixWorkflow` → lógica de auto-corrección (portar de `workflow-auto-fixer.ts`)

### Fase 5: Nodos y Templates (5sp)

**Archivos:** `handler_nodes.go`, `handler_templates.go`, `node_db.go`

- `SearchNodes` → SQLite query en `nodes.db`
- `GetNode` → SQLite query + procesamiento
- `SearchTemplates` → SQLite o HTTP a n8n templates API
- `GetTemplate` → SQLite o HTTP

### Fase 6: Sistema y Auditoría (5sp)

**Archivos:** `handler_system.go`, `handler_credentials.go`, `handler_datatable.go`

- `HealthCheck` → GET /api/health
- `AuditInstance` → análisis completo de la instancia
- `ManageCredentials` → CRUD /api/credentials
- `ManageDatatable` → CRUD de tablas internas

### Fase 7: Generación y Versioning (5sp)

**Archivos:** `handler_generate.go`, `handler_versions.go`, `handler_deploy.go`

- `GenerateWorkflow` → IA-assisted workflow generation
- `WorkflowVersions` → GET /api/workflows/{id}/versions
- `DeployTemplate` → crear workflow desde template

### Fase 8: Tool Definitions + Registrar (5sp)

**Archivos:** `tools.go`, `registrar.go`, `registrar_test.go`

Definir las 20+ herramientas como `Tool` interface y registrarlas.

### Fase 9: CLI Command + Config (3sp)

- `cmd/picoclaw/internal/util/n8n_mcp.go` → subcommand `util n8n-mcp-server`
- `pkg/config/config.go` → `N8nConfig`
- `config/config.example.json` → sección `tools.n8n`

---

## 7. Plan de Tests QA

### Tests Unitarios (por fase)

| Fase | Archivo de Test | Qué prueba |
|------|----------------|------------|
| 1 | `client_test.go` | HTTP requests, auth, retry, timeout |
| 2 | `handler_workflows_test.go` | CRUD con mock server |
| 3 | `validation_test.go` | Workflows válidos/inválidos |
| 3 | `security_scanner_test.go` | Detección de credenciales expuestas |
| 4 | `handler_execution_test.go` | Ejecución con mock server |
| 5 | `handler_nodes_test.go` | Búsqueda en SQLite |
| 6 | `handler_system_test.go` | Health check, audit |

### Dry-Run: Mock n8n Server

Crear un servidor HTTP mock que simule la API de n8n:

```go
// tests/n8n_mock_server.go
func NewMockN8nServer() *httptest.Server {
    return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        switch {
        case r.URL.Path == "/api/health" && r.Method == "GET":
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(200)
            w.Write([]byte(`{"status": "ok", "version": "1.80.0"}`))
        case r.URL.Path == "/api/workflows" && r.Method == "GET":
            w.WriteHeader(200)
            w.Write([]byte(`[{"id": "1", "name": "Test", "nodes": [], "connections": {}}]`))
        case strings.HasPrefix(r.URL.Path, "/api/workflows/") && r.Method == "GET":
            w.WriteHeader(200)
            w.Write([]byte(`{"id": "1", "name": "Test", "nodes": [], "connections": {}}`))
        case r.URL.Path == "/api/workflows" && r.Method == "POST":
            w.WriteHeader(201)
            w.Write([]byte(`{"id": "2", "name": "New", "nodes": [], "connections": {}}`))
        // ... más casos
        default:
            http.Error(w, "not found", 404)
        }
    }))
}
```

### Test de Integración End-to-End

1. Levantar mock n8n server
2. Crear N8nClient apuntando al mock
3. Registrar tools en ToolRegistry
4. Ejecutar cada tool via ToolRegistry.Execute()
5. Verificar respuesta esperada
6. Verificar que las requests HTTP fueron correctas

---

## 8. Métricas de Equivalencia

Para garantizar que la versión Go hace **exactamente lo mismo**, compararemos:

| Métrica | n8n-mcp (TS) | Go-Native | Cómo verificar |
|---------|-------------|-----------|----------------|
| Herramientas expuestas | 20 | 20 | `ListTools` request → comparar nombres |
| Schema de cada tool | JSON Schema | `Parameters()` | Comparar field por field |
| Respuestas de error | Formato n8n | Mismo formato | Comparar mensajes de error |
| Health check response | JSON con version | Mismo JSON | Comparar estructura |
| Workflow create | POST /api/workflows | Mismo endpoint | Capturar HTTP request |
| Validación | mismos errores | mismos errores | Enviar workflow inválido |

---

## 9. Riesgos y Mitigación

| Riesgo | Probabilidad | Impacto | Mitigación |
|--------|-------------|---------|------------|
| n8n API cambia entre versiones | Media | Medio | Usar version header + test contra n8n 1.70+ |
| SQLite DB de nodos cambia formato | Baja | Alto | Usar schema migration |
| Tool schemas no coinciden 100% | Media | Alto | Test de equivalencia contra outputs de TS |
| Auto-fix logic complejo de portar | Alta | Medio | Implementar versión simplificada primero |
| Security scanner pierde casos | Media | Alto | Portar tests de security scanner del TS |

---

## 10. Estimación Total

| Categoría | Story Points |
|-----------|-------------|
| Fase 1: Cliente + Tipos | 8sp |
| Fase 2: Workflow CRUD | 8sp |
| Fase 3: Validación | 8sp |
| Fase 4: Ejecución + Auto-fix | 8sp |
| Fase 5: Nodos + Templates | 5sp |
| Fase 6: Sistema + Auditoría | 5sp |
| Fase 7: Generación + Versioning | 5sp |
| Fase 8: Tools + Registrar | 5sp |
| Fase 9: CLI + Config | 3sp |
| QA + Tests | 10sp |
| **Total** | **65sp** |

---

## 11. Referencia Rápida de la API de n8n

| Método | Endpoint | Uso |
|--------|----------|-----|
| GET | `/api/health` | Health check |
| GET | `/api/workflows` | Listar workflows |
| POST | `/api/workflows` | Crear workflow |
| GET | `/api/workflows/{id}` | Obtener workflow |
| PUT | `/api/workflows/{id}` | Actualizar workflow |
| DELETE | `/api/workflows/{id}` | Eliminar workflow |
| POST | `/api/workflows/{id}/execute` | Ejecutar workflow |
| GET | `/api/executions` | Listar ejecuciones |
| GET | `/api/credentials` | Listar credenciales |
| POST | `/api/credentials` | Crear credencial |
| DELETE | `/api/credentials/{id}` | Eliminar credencial |

---

## 12. Dry-Run: Pasos de Verificación

Una vez implementado:

```bash
# 1. Compilar
go build ./...

# 2. Tests unitarios
go test ./pkg/mcp/n8n/... -v -count=1

# 3. Equivalencia de schemas (comparar con TS)
# Ejecutar n8n-mcp TS y Go, comparar ListTools output
node /Volumes/.../n8n-mcp/dist/mcp/index.js &
# Enviar ListTools request a ambos, diff outputs

# 4. Integración con picoclaw-agents
./build/picoclaw-agents util n8n-mcp-server

# 5. Verificar que el agente descubre las tools
# Iniciar gateway y verificar tools en /api/agent/tools
```

---

*mcp_n8n_migracion_golang.md — 2026-04-06*
*Objetivo: Portar ~20 herramientas MCP de n8n-mcp a Go nativo, cero dependencias nuevas*
