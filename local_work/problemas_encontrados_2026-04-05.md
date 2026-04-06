# Problemas Encontrados y Resueltos — 2026-04-05

**Fuente:** `upstream_audit_2026-04-05.json`  
**Total:** 4 issues  
**Resueltos:** 4/4 ✅

---

## 1. ✅ `secret-placeholder.ts` — Exposición de secretos cortos (ce161905)

**Resuelto:** Ya aplicado en versión previa del frontend. Verificado idéntico al original.

---

## 2. ✅ `filesystem.go` — Escape JSON en write_file (71337b6f)

**Resuelto hoy:** Actualizada descripción de `WriteFileTool` y parámetro `content` para clarificar semántica de escape `\n` vs `\\n`.

**Archivo:** `pkg/tools/filesystem.go`
**Cambio:** Descripción + content description actualizados

---

## 3. ✅ `loop.go` — Detección de overflow de contexto (97dec167)

**Resuelto:** Cubierto por integración de ContextManager (Fases 0-5). Solución más robusta que el parche original: budget-aware assembly + multi-level summarization + 3-level system prompt.

---

## 4. ✅ `context_budget.go` — Doble conteo de tokens del sistema (1a44752d)

**Resuelto:** Portado en Fase 0.2 como `pkg/tokenizer/estimator.go`. Usa `max(Content, SystemParts)` para evitar sobreestimación.

---

## Build Verification

- ✅ `make build` exitoso
- ✅ Sin regresiones de seguridad
- ✅ Características custom preservadas

---

## 5. ⚠️ WebUI muestra "Not Configured" para modelos Modelfile de Ollama

**Fecha:** 2026-04-05
**Estado:** Pendiente — ver `fix_webui_modelfiles.md`
**Severidad:** Media (solo WebUI, el agente funciona)

**Problema:** Modelos Ollama creados con Modelfile (ej: `picoclaw-qwen25-min`) aparecen como "Not Configured" en `http://localhost:18800/models` a pesar de estar configurados correctamente en `config.json`.

**Causa:** La función `probeLocalModelAvailability()` en `web/backend/api/model_status.go` interpreta `"picoclaw-qwen25-min"` (sin prefijo `ollama/`) como modelo OpenAI y usa el probe incorrecto (`probeOpenAICompatibleModel` en vez de `probeOllamaModel`). Ollama retorna el nombre base del modelo (`qwen2.5:0.5b`), no el nombre del Modelfile, por lo que el probe falla.

**Solución:** Modificar `model_status.go` para detectar Ollama local cuando `auth_method: "local"` y `api_base` apunta a `localhost:11434`, sin depender del prefijo `ollama/` en el nombre del modelo.

**Workaround actual:** El modelo funciona desde CLI (`./build/picoclaw-agents agent --model picoclaw-qwen25-min -m "hola"`) y desde el WebUI chat (seleccionando el modelo), solo el status muestra "Not Configured".

---

*problemas_encontrados_2026-04-05.md — Actualizado el 5 de abril de 2026*
