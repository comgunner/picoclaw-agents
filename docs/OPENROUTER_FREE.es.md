# OpenRouter Free Tier con PicoClaw-Agents

**Última actualización:** 28 de marzo de 2026
**Versión:** v1.3.0-alpha-fix901

---

## 🆓 Inicio Rápido — Gratis sin API Key

PicoClaw-Agents soporta modelos **100% gratuitos** de OpenRouter sin necesidad de tarjeta de crédito.

### Opción 1: Onboard Interactivo (Recomendado)

```bash
# Ejecutar wizard de configuración
picoclaw-agents onboard --free
```

**El wizard te guiará:**
1. Solicitará tu API key de OpenRouter (gratis)
2. Configurarará modelos free automáticamente
3. Creará tu espacio de trabajo

### Opción 2: Config Manual

Crea `~/.picoclaw/config.json`:

```json
{
  "agents": {
    "defaults": {
      "workspace": "~/.picoclaw/workspace",
      "model": "openrouter/auto",
      "max_tokens": 8192,
      "max_tool_iterations": 20
    }
  },
  "model_list": [
    {
      "model_name": "or-auto",
      "model": "openrouter/auto",
      "api_key": "sk-or-v1-TU_API_KEY_AQUI"  // pragma: allowlist secret
    }
  ]
}
```

---

## 🎯 Modelos Gratuitos Disponibles

OpenRouter ofrece varios modelos gratuitos. PicoClaw los configura automáticamente:

### Configuración por Defecto (`onboard --free`)

| Prioridad | Modelo | Contexto | Uso Recomendado |
|-----------|--------|----------|-----------------|
| **1** | `openrouter/auto` | Variable | Auto-selecciona el mejor free |
| **2** | `stepfun/step-3.5-flash` | 256K | Contexto largo, razonamiento |
| **3** | `deepseek/deepseek-v3.2-20251201` | 64K | Inferencia rápida |

### ¿Qué es `openrouter/auto`?

- **Auto-selección:** OpenRouter elige automáticamente el mejor modelo free disponible
- **Fallback automático:** Si un modelo falla, usa el siguiente
- **Sin configuración:** No necesitas especificar modelos individuales

---

## 📝 Comandos Útiles

### Verificar Configuración

```bash
# Ver modelo actual
picoclaw-agents agent --model "openrouter/auto" -m "¿Qué modelo estás usando?"

# Ver estado de autenticación
picoclaw-agents auth status
```

### Probar Modelos Individuales

```bash
# StepFun (256K contexto)
picoclaw-agents agent --model "stepfun/step-3.5-flash" -m "Hello"

# DeepSeek (rápido)
picoclaw-agents agent --model "deepseek/deepseek-v3.2-20251201" -m "Hello"
```

### Cambiar Modelo por Defecto

Edita `~/.picoclaw/config.json`:

```json
{
  "agents": {
    "defaults": {
      "model": "stepfun/step-3.5-flash"  // ← Cambiar aquí
    }
  }
}
```

---

## 🔑 Obtener API Key de OpenRouter

### Paso 1: Crear Cuenta

1. Ve a https://openrouter.ai
2. Click en "Sign Up"
3. Regístrate con email (no requiere tarjeta)

### Paso 2: Crear API Key

1. Ve a https://openrouter.ai/keys
2. Click "Create Key"
3. Copia la key (empieza con `sk-or-v1-...`)

### Paso 3: Configurar en PicoClaw

```bash
# Durante el onboard
picoclaw-agents onboard --free
# → Pega tu API key cuando la solicite

# O edita config.json manualmente
nano ~/.picoclaw/config.json
# → Pega tu API key en "api_key"
```

---

## ⚠️ Problemas Comunes

### Error: "openrouter-free is not a valid model ID"

**Causa:** Configuración antigua con nombre de modelo inválido.

**Solución:** Actualiza tu config:

```bash
# Opción 1: Re-ejecutar onboard
picoclaw-agents onboard --free

# Opción 2: Editar config manualmente
sed -i 's|"openrouter/free"|"openrouter/auto"|g' ~/.picoclaw/config.json
sed -i 's|"openrouter-free"|"openrouter/auto"|g' ~/.picoclaw/config.json
```

### Error: "401 Unauthorized"

**Causa:** API key inválida o expirada.

**Solución:**
1. Verifica tu key en https://openrouter.ai/keys
2. Actualiza en `~/.picoclaw/config.json`
3. Re-ejecuta `picoclaw-agents onboard --free`

### Error: "Rate limit exceeded"

**Causa:** Límite de requests gratuitos alcanzado.

**Límites Free Tier:**
- ~50 requests/minuto
- ~1000 requests/día (varía por modelo)

**Solución:**
- Espera unos minutos
- Usa modelos más lentos pero con más límite
- Considera upgrade a paid tier

---

## 📊 Límites del Free Tier

### Rate Limits

| Modelo | Requests/min | Requests/día | Contexto Max |
|--------|--------------|--------------|--------------|
| `openrouter/auto` | ~50 | ~1000 | Variable |
| `stepfun/step-3.5-flash` | ~20 | ~500 | 256K |
| `deepseek/deepseek-v3.2` | ~30 | ~800 | 64K |

### Mejores Prácticas

1. **Usa `openrouter/auto`** — Mejor balance disponibilidad/velocidad
2. **Evita polling** — No hagas requests muy frecuentes
3. **Batch tasks** — Agrupa tareas cuando sea posible
4. **Monitorea uso** — Revisa tu dashboard en openrouter.ai

---

## 🚀 Ejemplos de Uso

### Chat Simple

```bash
picoclaw-agents agent -m "¿Cuál es la capital de Francia?"
```

### Tarea de Código

```bash
picoclaw-agents agent -m "Crea una función Python que calcule Fibonacci"
```

### Búsqueda Web

```bash
picoclaw-agents agent -m "Busca las últimas noticias de IA"
```

### Tarea Compleja (Multi-agente)

```bash
# Con team mode configurado
picoclaw-agents agent -m "Crea una API REST con Node.js y Express"
```

---

## 📚 Recursos Adicionales

### Enlaces Oficiales

- **OpenRouter:** https://openrouter.ai
- **Modelos Free:** https://openrouter.ai/models?order=-free
- **API Docs:** https://openrouter.ai/docs
- **Keys:** https://openrouter.ai/keys

### Documentación PicoClaw

- **CHANGELOG:** [CHANGELOG.md](../CHANGELOG.md)
- **README:** [README.md](../README.md)
- **Fix #901:** [local_work/fix_901_openrouter_normalization.md](../local_work/fix_901_openrouter_normalization.md)

---

## ❓ FAQ

### ¿Es realmente gratis?

Sí. OpenRouter ofrece modelos free sin tarjeta de crédito. Hay límites de rate pero son suficientes para uso personal.

### ¿Necesito configurar algo más?

No. Con `picoclaw-agents onboard --free` todo se configura automáticamente.

### ¿Puedo cambiar a modelos paid después?

Sí. Solo edita `config.json` y cambia el modelo o agrega tu tarjeta en OpenRouter.

### ¿Qué pasa si se agotan los free models?

OpenRouter auto-selecciona otro modelo free disponible. Si todos están agotados, recibirás un error de rate limit.

### ¿Funciona con todos los canales (Telegram, Discord)?

Sí. El free tier funciona igual para CLI, Telegram, Discord, etc.

---

**Documento creado:** 28 de marzo de 2026
**Versión:** v1.3.0-alpha-fix901
**Mantenimiento:** @comgunner
