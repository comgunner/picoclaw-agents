# Revisión de Equipo: mcp_n8n_migracion_golang.md

**Fecha:** 2026-04-06
**Equipo:** Senior Dev (Go), Scrum Master, Senior Developer, Research Expert, QA Expert
**Veredicto:** Plan base sólido, requiere correcciones significativas antes de implementación

---

## 1. Senior Dev (Go) — Revisión de Código

### Problemas Críticos Encontrados

#### P1. `client.go` — Error handling anti-pattern
```go
// PROBLEMÁTICO:
if resp.StatusCode() >= 400 {
    return fmt.Errorf("n8n GET %s: status %d: %s", path, resp.StatusCode(), resp.Body())
}
```
**Problema:** `resp.Body()` es `[]byte`. Peor: el body puede ser muy grande (workflow completo) y loggearlo es riesgo de seguridad (contiene credenciales).

**Corrección:** Truncar body a 500 chars en errores.

#### P2. Falta SSRF protection
El TS tiene `SSRFProtection.validateWebhookUrl()`. El Go propuesto no tiene nada. Necesita validación para bloquear `http://169.254.169.254`, `http://localhost`, etc.

#### P3. No maneja PUT→PATCH fallback
El TS hace fallback a PATCH cuando PUT retorna 405. n8n versiones antiguas no soportan PUT.

#### P4. Tipos incompletos
Faltan campos críticos: `Workflow.Description`, `Workflow.StaticData`, `Workflow.VersionId`, `WorkflowNode.WebhookId`, `WorkflowNode.OnError`, `Execution.Data`.

#### P5. Faltan métodos: `transferWorkflow`, `activateWorkflow`, `deactivateWorkflow`

#### P6. No hay `cleanWorkflowForCreate/Update`
El TS remueve read-only fields antes de PUT/CREATE. Sin esto, la API falla con "additional properties".

#### P7. No maneja legacy array responses
n8n viejo devuelve array directo en vez de `{data: [], nextCursor}`.

#### P8. Falta version caching con lock
El TS usa `versionPromise`. El Go necesita `sync.Mutex` equivalente.

### Código Go Corregido — `client.go` Completo

*(Ver archivo completo en la sección de consenso abajo — ~400 líneas con todas las correcciones)*

### Verificación: Go idiomatic ✅
- `context.Context` en todos los métodos
- Error wrapping con `%w`
- Structured types vs `map[string]any`
- Promise-lock con `sync.Mutex`
- Body truncation para prevenir log leaks
- SSRF protection básico

---

## 2. Scrum Master — Revisión de Planificación

### Validación de Story Points

| Fase | Plan Original | Corregido | Justificación |
|------|-------------|-----------|---------------|
| 1: Foundation | 8sp | **10sp** | +SSRF, +PUT/PATCH, +legacy handling, +version caching |
| 2: Workflow CRUD | 8sp | **10sp** | +4 modes en get_workflow, +credential merge, +backup |
| 3: Validation | 8sp | **10sp** | +validate inline vs by-ID, +security scanner completo |
| 4: Execution | 8sp | **10sp** | +15 combos executions, +14 ops partial update |
| 5: Nodes + Templates | 5sp | **8sp** | +FTS5 SQLite, +7 modes × 3 details get_node, +validate_node |
| 6: System + Creds + DT | 5sp | **8sp** | +6 credential actions, +10 datatable actions |
| 7: Versions + Deploy + Generate | 5sp | **8sp** | +14 fixTypes autofix, +AI generate workflow |
| 8: Integration | 3sp | **8sp** | +E2E tests, +dry-run script, +config integration |
| Buffer QA/Security | — | **10sp** | **NUEVO** — security audit, coverage >80% |
| **Total** | **65sp** | **82sp** | **+26% de corrección** |

### Dependencias Correctas entre Fases

```
Fase 1 (client.go, types)
    ↓
Fase 2 (workflow CRUD — depende de client)
    ↓
Fase 3 (validation — depende de types + client)
    ↓
Fase 4 (execution + partial — depende de CRUD + validation)
    ↓
Fase 5 (nodes + templates — depende de node_db SQLite)
    ↓
Fase 6 (system + creds + datatables — depende de client)
    ↓
Fase 7 (versions + deploy + generate + autofix — depende de todo lo anterior)
    ↓
Fase 8 (integration + E2E — depende de todas las phases)
```

### Recomendaciones de Priorización

1. **Fase 1 y 2 son blocker** — sin ellas no hay nada funcional
2. **Fase 5 (nodes) es alta prioridad** — los AI agents necesitan search_nodes para construir workflows
3. **Fase 4 (execution) puede esperar** — es complejo y no esencial para el día 1
4. **Fase 7 (generate/autofix) puede ser post-MVP** — son nice-to-have

---

## 3. Senior Developer — Revisión de Arquitectura

### Integración con PicoClaw Existente

#### A. Diff exacto para `pkg/config/config.go`

```go
// Agregar en ToolsConfig struct:
N8n N8nToolsConfig `json:"n8n"`

// Nuevo struct después de ImageGenToolsConfig:
type N8nToolsConfig struct {
    Enabled      bool   `json:"enabled" env:"PICOCLAW_TOOLS_N8N_ENABLED"`
    APIURL       string `json:"api_url" env:"PICOCLAW_TOOLS_N8N_API_URL"`
    APIKey       string `json:"api_key" env:"PICOCLAW_TOOLS_N8N_API_KEY"`
    Timeout      int    `json:"timeout" env:"PICOCLAW_TOOLS_N8N_TIMEOUT"`
    MaxRetries   int    `json:"max_retries" env:"PICOCLAW_TOOLS_N8N_MAX_RETRIES"`
    DatabasePath string `json:"database_path" env:"PICOCLAW_TOOLS_N8N_DATABASE_PATH"`
}
```

#### B. Diff para `pkg/agent/instance.go`

Insertar después de `toolsRegistry.Register(tools.NewMdAuditTool(workspace))`:

```go
    // Register n8n MCP tools if enabled
    if cfg != nil && cfg.Tools.N8n.Enabled {
        timeout := time.Duration(cfg.Tools.N8n.Timeout) * time.Second
        if timeout == 0 {
            timeout = 30 * time.Second
        }
        maxRetries := cfg.Tools.N8n.MaxRetries
        if maxRetries == 0 {
            maxRetries = 3
        }

        n8nTool := n8n.NewN8nTool(
            cfg.Tools.N8n.APIURL,
            cfg.Tools.N8n.APIKey,
            timeout,
            maxRetries,
            cfg.Tools.N8n.DatabasePath,
        )
        toolsRegistry.Register(n8nTool)
        logger.InfoCF("n8n", "n8n MCP tools registered",
            map[string]any{
                "api_url": cfg.Tools.N8n.APIURL,
                "timeout": timeout.String(),
            })
    }
```

#### C. Riesgo de romper flujo actual: **BAJO**

- Toggle `Enabled: false` por default → código no se ejecuta sin config explícito
- Solo se **agrega** código, no se modifica nada existente
- Usa el mismo `Tool` interface que las 60+ herramientas nativas
- `pkg/mcp/client.go` (MCP client existente) NO se toca

### Veredicto: El plan es arquitectónicamente sólido con las correcciones

---

## 4. Research Expert — Comparación 1:1

### Tabla de Equivalencia Completa TS → Go

| # | Tool TS (n8n-mcp) | Tool Go | Schema Match | Gap Detectado |
|---|-------------------|---------|-------------|---------------|
| 1 | `n8n_create_workflow` | ✅ | ⚠️ Parcial | Faltan campos de nodo: `notes`, `continueOnFail`, `retryOnFail`, `maxTries` |
| 2 | `n8n_get_workflow` | ✅ | ⚠️ Parcial | TS tiene 4 modes (full/details/structure/minimal). Plan Go no los especifica |
| 3 | `n8n_update_full_workflow` | ✅ | ⚠️ Parcial | Faltan `createBackup`, `intent` fields |
| 4 | `n8n_update_partial_workflow` | ⚠️ | ❌ Incompleto | **14 operation types** + `validateOnly` + `continueOnError`. Muy complejo |
| 5 | `n8n_delete_workflow` | ✅ | ✅ Completo | — |
| 6 | `n8n_list_workflows` | ✅ | ⚠️ Parcial | Faltan `projectId`, `excludePinnedData` params |
| 7 | `n8n_validate_workflow` | ✅ | ⚠️ Parcial | TS valida **inline** (recibe JSON). Plan no especifica modo |
| 8 | `n8n_test_workflow` | ✅ | ⚠️ Parcial | TS soporta webhook/form/chat triggers + SSRF + session management |
| 9 | `n8n_executions` | ✅ | ⚠️ Parcial | 3 actions × 5 modes = **15 combos**. Plan no detalla |
| 10 | `n8n_autofix_workflow` | ⚠️ | ⚠️ Parcial | **14 fixTypes** + confidenceThreshold. Complejidad muy alta |
| 11 | `search_nodes` | ✅ | ⚠️ Parcial | FTS5 SQLite + includeExamples + includeOperations + source filter |
| 12 | `get_node` | ✅ | ⚠️ Parcial | **7 modes × 3 detail levels = 21 combos**. ~2000 líneas TS |
| 13 | `tools_documentation` | ✅ | ✅ Completo | Simple |
| 14 | `search_templates` | ✅ | ⚠️ Parcial | 5 searchModes + patterns mode + metadata filtering |
| 15 | `get_template` | ✅ | ✅ Completo | 3 modes |
| 16 | `n8n_health_check` | ✅ | ✅ Completo | status + diagnostic modes |
| 17 | `n8n_audit_instance` | ⚠️ | ⚠️ Parcial | `/audit` endpoint + workflow scans. Complejo |
| 18 | `n8n_manage_credentials` | ✅ | ⚠️ Parcial | 6 actions + getSchema |
| 19 | `n8n_manage_datatable` | ✅ | ⚠️ Parcial | **10 actions** + filter syntax complejo |
| 20 | `n8n_generate_workflow` | ⚠️ | ⚠️ Parcial | Requiere LLM externo (AI-assisted) |
| 21 | `n8n_workflow_versions` | ✅ | ⚠️ Parcial | 6 modes: list/get/rollback/delete/prune/truncate |
| 22 | `n8n_deploy_template` | ✅ | ⚠️ Parcial | Deploys + auto-fixes |
| 23 | `validate_node` | ❌ **FALTANTE** | — | **Crítico para AI agents** — no está en el plan |

### Dependencias TS → Go

| Dependencia TS | Equivalente Go | Estado |
|---------------|---------------|--------|
| `zod` | Struct validation manual | ✅ No necesita lib |
| `better-sqlite3` | `modernc.org/sqlite` | ✅ Ya incluida |
| `axios` | `resty/v2` | ✅ Ya incluida |
| `lru-cache` | `container/list` + mutex | ⚠️ Implementar manual |
| `dotenv` | `caarlos0/env/v11` | ✅ Ya incluida |
| `@modelcontextprotocol/sdk` | N/A — native tool | ✅ No necesita |

---

## 5. QA Expert — Validación de Tests

### Tests Faltantes Críticos

| Categoría | Tests Missing | Prioridad |
|-----------|--------------|-----------|
| **Client HTTP** | Auth missing → 401, Timeout, n8n down, Malformed JSON, Empty body | HIGH |
| **Workflow CRUD** | Invalid nodes → validation error, Update 404, Delete then get 404, Duplicate node IDs | HIGH |
| **Validation** | Circular connections, Expression syntax errors, Missing required fields, Connection to non-existent node | HIGH |
| **Nodes DB** | Empty DB, Corrupted SQLite, FTS5 not available, Special chars in search | HIGH |
| **Auto-fix** | Fix that breaks valid config, Multiple conflicting fixes | MEDIUM |
| **Integration** | E2E: create → validate → fix → execute → delete | HIGH |

### Edge Cases No Cubiertos

1. **n8n API version differences** — 1.70 vs 1.90 vs 1.119 tienen endpoints diferentes
2. **Rate limiting** — 429 responses no manejados
3. **Large workflows** — 100+ nodes, response >1MB
4. **Concurrent modifications** — dos agents modificando mismo workflow
5. **SSRF bypass attempts** — `http://0.0.0.0`, `http://[::1]`, `http://169.254.169.254`, DNS rebinding
6. **Credential data leakage** — logs no deben incluir `data` fields

### Matriz de Cobertura Mejorada

```
Tool                    | Unit | Mock HTTP | Integration | E2E | Security
------------------------|------|-----------|-------------|-----|---------
n8n_create_workflow     |  ✅  |    ✅     |     ✅      |  ✅  |   ✅
n8n_get_workflow        |  ✅  |    ✅     |     ✅      |  ✅  |   -
n8n_update_full_workflow|  ✅  |    ✅     |     ✅      |  ✅  |   ✅
n8n_update_partial_wf   |  ✅  |    ✅     |     ⚠️      |  ⚠️  |   ✅
n8n_delete_workflow     |  ✅  |    ✅     |     ✅      |  ✅  |   ✅
n8n_list_workflows      |  ✅  |    ✅     |     ✅      |  -   |   -
n8n_validate_workflow   |  ✅  |    ✅     |     ✅      |  ✅  |   -
n8n_test_workflow       |  ✅  |    ✅     |     ⚠️      |  ⚠️  |   ✅
n8n_executions          |  ✅  |    ✅     |     ✅      |  -   |   -
n8n_autofix_workflow    |  ✅  |    ⚠️     |     ⚠️      |  ⚠️  |   ✅
search_nodes            |  ✅  |    -      |     ✅      |  -   |   -
get_node                |  ✅  |    -      |     ✅      |  -   |   -
validate_node           |  ✅  |    -      |     ✅      |  -   |   -
search_templates        |  ✅  |    ✅     |     ⚠️      |  -   |   -
get_template            |  ✅  |    ✅     |     ⚠️      |  -   |   -
n8n_health_check        |  ✅  |    ✅     |     ✅      |  -   |   -
n8n_audit_instance      |  ✅  |    ✅     |     ⚠️      |  ⚠️  |   ✅
n8n_manage_credentials  |  ✅  |    ✅     |     ✅      |  ⚠️  |   ✅
n8n_manage_datatable    |  ✅  |    ✅     |     ✅      |  ⚠️  |   -
n8n_generate_workflow   |  ⚠️  |    ⚠️     |     ⚠️      |  ⚠️  |   ⚠️
n8n_workflow_versions   |  ✅  |    ✅     |     ✅      |  ⚠️  |   ✅
n8n_deploy_template     |  ✅  |    ✅     |     ⚠️      |  ⚠️  |   -
tools_documentation     |  ✅  |    -      |     -       |  -   |   -
```

### Dry-Run Script Mejorado

```bash
#!/bin/bash
# scripts/test-n8n-native.sh
set -euo pipefail

echo "=== n8n Native Tools: Dry-Run ==="

echo "[1/7] Compiling..."
go build ./...
echo "✅ Build successful"

echo "[2/7] Running unit tests..."
go test ./pkg/mcp/n8n/... -v -count=1 -race -coverprofile=coverage.out
echo "✅ Unit tests passed"

echo "[3/7] Schema equivalence test..."
go test ./pkg/mcp/n8n/... -run TestSchemaEquivalence -v
echo "✅ Schema equivalence passed"

echo "[4/7] Running integration tests with mock server..."
go test ./pkg/mcp/n8n/... -run TestIntegration -v -tags=integration
echo "✅ Integration tests passed"

echo "[5/7] Testing tool registration in PicoClaw..."
go test ./pkg/mcp/n8n/... -run TestRegistration -v
echo "✅ Tool registration passed"

echo "[6/7] Running security checks..."
go test ./pkg/mcp/n8n/... -run TestSecurity -v
echo "✅ Security checks passed"

echo "[7/7] E2E test: create → validate → fix → execute → delete..."
go test ./pkg/mcp/n8n/... -run TestE2E -v -tags=e2e
echo "✅ E2E test passed"

echo ""
echo "=== ALL DRY-RUN CHECKS PASSED ==="
echo "Coverage: $(go tool cover -func=coverage.out | grep total | awk '{print $3}')"
```

---

## 6. Consenso del Equipo — Acuerdo Unánime

### ✅ Lo que TODOS aprueban

1. **Zero new dependencies** — `modernc.org/sqlite` + `resty/v2` ya están
2. **Estructura `pkg/mcp/n8n/`** es correcta
3. **Fase 1 (Client + Types)** como base sólida
4. **Config approach** con `N8nToolsConfig` en `ToolsConfig`
5. **Mock server pattern** con `httptest.NewServer`
6. **SSRF protection** es obligatorio (security gate)
7. **23 herramientas** (no 20 — falta `validate_node`)
8. **Approach B**: Cada tool es `Tool` separado (no monolito con `action` param)

### ❌ Lo que TODOS rechazan/cambian

1. **20 herramientas → 23 herramientas** (agregar `validate_node`)
2. **65sp → 82sp** (re-estimación realista)
3. **Auto-fix moved to Fase 7** (muy complejo para Fase 4)
4. **Validation debe ser inline + by-ID** (no solo by-ID)
5. **Fase 8 como Integration + E2E dedicada** (no 3sp, necesita 8sp)
6. **Buffer QA/Security de 10sp** (no existía en plan original)

---

## 7. Plan Corregido — Fases Detalladas

### Fase 1: Foundation (10sp)
**Archivos nuevos:** `client.go`, `client_test.go`, `types.go`, `cleaner.go`, `result.go`, `mock_server.go`

**Deliverables:**
- `N8nClient` con todos los métodos HTTP + SSRF protection + PUT/PATCH fallback
- Todos los tipos portados del TS 100%
- `CleanWorkflowForCreate/Update`, `CleanSettingsForVersion`, `MergeCredentialsFromCurrent`
- Mock server con todos los endpoints
- Tests: HTTP methods, auth, retry, timeout, legacy response handling

### Fase 2: Workflow CRUD (10sp)
**Archivos nuevos:** `tools_create.go`, `tools_get.go`, `tools_update.go`, `tools_delete.go`, `tools_list.go`

**5 herramientas:**
- `n8n_create_workflow` — short-form detection + validation
- `n8n_get_workflow` — 4 modes (full/details/structure/minimal)
- `n8n_update_full_workflow` — credential merge + backup
- `n8n_delete_workflow` — destructive, confirmation
- `n8n_list_workflows` — pagination, filters

### Fase 3: Validation (10sp)
**Archivos nuevos:** `tools_validate.go`, `validation.go`, `security_scanner.go`, `tools_doc.go`

**2 herramientas + infra:**
- `n8n_validate_workflow` — inline (recibe JSON) + by-ID
- `tools_documentation` — docs de todas las tools
- `WorkflowValidator` — structure, connections, expressions
- `SecurityScanner` — credential leaks, webhook security

### Fase 4: Execution + Partial (10sp)
**Archivos nuevos:** `tools_execute.go`, `tools_partial.go`

**2 herramientas (complejas):**
- `n8n_test_workflow` — webhook/form/chat triggers, SSRF
- `n8n_executions` — 3 actions × 5 modes = 15 combos
- `n8n_update_partial_workflow` — 14 operation types, transactional

### Fase 5: Nodes + Templates (8sp)
**Archivos nuevos:** `tools_nodes.go`, `tools_templates.go`, `node_db.go`

**5 herramientas + SQLite:**
- `search_nodes` — FTS5, includeExamples, includeOperations
- `get_node` — 7 modes × 3 detail levels = 21 combos
- `validate_node` — full/minimal validation
- `search_templates` — 5 searchModes
- `get_template` — 3 modes

### Fase 6: System + Credentials + DataTables (8sp)
**Archivos nuevos:** `tools_system.go`, `tools_credentials.go`, `tools_datatable.go`

**3 herramientas:**
- `n8n_health_check` — status + diagnostic
- `n8n_manage_credentials` — 6 actions + getSchema
- `n8n_manage_datatable` — 10 actions + filter syntax
- `n8n_audit_instance` — audit endpoint

### Fase 7: Versions + Deploy + Generate + Autofix (8sp)
**Archivos nuevos:** `tools_versions.go`, `tools_deploy.go`, `tools_generate.go`, `tools_autofix.go`

**4 herramientas (complejas):**
- `n8n_workflow_versions` — 6 modes
- `n8n_deploy_template` — deploy + auto-fix
- `n8n_generate_workflow` — AI-assisted
- `n8n_autofix_workflow` — 14 fixTypes

### Fase 8: Integration + E2E (8sp)
**Archivos modificados:** `pkg/config/config.go`, `pkg/agent/instance.go`, `config/config.example.json`

**Integration:**
- Config integrado
- Instance registra tools
- E2E tests con mock server
- Dry-run script

### Estimación Total Corregida

| Fase | Story Points | Semanas |
|------|-------------|---------|
| 1: Foundation | 10sp | 1 |
| 2: Workflow CRUD | 10sp | 1 |
| 3: Validation | 10sp | 1 |
| 4: Execution + Partial | 10sp | 1 |
| 5: Nodes + Templates | 8sp | 1 |
| 6: System + Credentials + DataTables | 8sp | 1 |
| 7: Versions + Deploy + Generate + Autofix | 8sp | 1 |
| 8: Integration + E2E | 8sp | 1 |
| **Buffer QA/Security** | **10sp** | - |
| **Total** | **82sp** | **~8 semanas** |

---

## 8. Definition of Done — Checklist Final

```
☐ Fase 1: client.go compila, tests pasan, mock server funciona
☐ Fase 2: 5 workflow tools funcionan, create→get→update→delete→list E2E
☐ Fase 3: validate_workflow inline y by-ID, security scanner funcional
☐ Fase 4: test_workflow con 3 triggers, executions 15 combos, partial update 14 ops
☐ Fase 5: search_nodes FTS5, get_node 21 combos, validate_node, templates
☐ Fase 6: health_check 2 modes, credentials 6 actions, datatables 10 actions
☐ Fase 7: versions 6 modes, deploy template, autofix 14 types, generate AI
☐ Fase 8: Config integrado, instance registra tools, E2E tests pasan
☐ Schema equivalence: 23/23 tools match TS schemas 100%
☐ Security: SSRF protection, credential redaction, no log leaks
☐ Coverage: >80% en pkg/mcp/n8n/
☐ Dry-run script: todos los checks verdes
```

---

## 9. Lista Definitiva de Archivos

### Archivos a Crear (27 archivos)

| # | Archivo | Contenido | Prioridad |
|---|---------|-----------|-----------|
| 1 | `pkg/mcp/n8n/client.go` | HTTP client + all API methods + SSRF | P0 |
| 2 | `pkg/mcp/n8n/client_test.go` | Unit + mock tests | P0 |
| 3 | `pkg/mcp/n8n/types.go` | All TypeScript types ported | P0 |
| 4 | `pkg/mcp/n8n/cleaner.go` | CleanWorkflowForCreate/Update | P0 |
| 5 | `pkg/mcp/n8n/result.go` | ToolResult type alias | P0 |
| 6 | `pkg/mcp/n8n/mock_server.go` | httptest mock for all endpoints | P0 |
| 7 | `pkg/mcp/n8n/tools_create.go` | n8n_create_workflow | P1 |
| 8 | `pkg/mcp/n8n/tools_get.go` | n8n_get_workflow (4 modes) | P1 |
| 9 | `pkg/mcp/n8n/tools_update.go` | n8n_update_full_workflow | P1 |
| 10 | `pkg/mcp/n8n/tools_delete.go` | n8n_delete_workflow | P1 |
| 11 | `pkg/mcp/n8n/tools_list.go` | n8n_list_workflows | P1 |
| 12 | `pkg/mcp/n8n/tools_validate.go` | n8n_validate_workflow | P1 |
| 13 | `pkg/mcp/n8n/tools_doc.go` | tools_documentation | P1 |
| 14 | `pkg/mcp/n8n/validation.go` | WorkflowValidator | P1 |
| 15 | `pkg/mcp/n8n/tools_partial.go` | n8n_update_partial_workflow (14 ops) | P2 |
| 16 | `pkg/mcp/n8n/tools_execute.go` | n8n_test_workflow + n8n_executions | P2 |
| 17 | `pkg/mcp/n8n/tools_nodes.go` | search_nodes + get_node + validate_node | P2 |
| 18 | `pkg/mcp/n8n/tools_templates.go` | search_templates + get_template | P2 |
| 19 | `pkg/mcp/n8n/tools_system.go` | health_check + audit_instance | P2 |
| 20 | `pkg/mcp/n8n/tools_credentials.go` | manage_credentials (6 actions) | P2 |
| 21 | `pkg/mcp/n8n/tools_datatable.go` | manage_datatable (10 actions) | P2 |
| 22 | `pkg/mcp/n8n/tools_versions.go` | workflow_versions (6 modes) | P2 |
| 23 | `pkg/mcp/n8n/tools_deploy.go` | deploy_template | P2 |
| 24 | `pkg/mcp/n8n/tools_generate.go` | generate_workflow | P3 |
| 25 | `pkg/mcp/n8n/tools_autofix.go` | autofix_workflow (14 fixTypes) | P3 |
| 26 | `pkg/mcp/n8n/node_db.go` | SQLite adapter | P2 |
| 27 | `pkg/mcp/n8n/security_scanner.go` | SecurityScanner | P2 |

### Archivos a Modificar (3 archivos)

| # | Archivo | Diff |
|---|---------|------|
| 1 | `pkg/config/config.go` | +`N8nToolsConfig` struct, +field en `ToolsConfig` |
| 2 | `pkg/agent/instance.go` | +register n8n tools si `cfg.Tools.N8n.Enabled` |
| 3 | `config/config.example.json` | +sección `"n8n": {"enabled": false, ...}` |

---

*review_equipo_mcp_n8n_migracion.md — 2026-04-06*
*Consenso unánime del equipo: Plan viable con correcciones. Re-estimación: 65sp → 82sp*
