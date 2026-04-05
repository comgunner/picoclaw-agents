# Guía del Proveedor Antigravity

**Última Actualización:** 30 de marzo de 2026  
**Estado:** ✅ Production Ready (v1.3.0-alpha)  
**Novedad:** 🎉 ¡Auto-Config en Login - Los 15 modelos se agregan automáticamente!

---

## Descripción General

**Antigravity** (Google Cloud Code Assist) es un proveedor de IA respaldado por Google que ofrece acceso a modelos Gemini y Claude a través de la infraestructura de Google Cloud usando **autenticación OAuth 2.0**.

**Distinción Clave:** Antigravity usa las cuotas de tu plan **Google One AI Premium** o **Workspace Gemini** — NO una clave API de pago por uso.

---

## Inicio Rápido (¡NUEVO!)

### Configuración en Un Comando

```bash
# Login y auto-configurar los 15 modelos de Antigravity
./picoclaw-agents auth login --provider google-antigravity

# Probar con modelo por defecto (gemini-3-flash)
./picoclaw-agents agent -m "Hola mundo"
```

**Lo que sucede automáticamente:**
1. ✅ Autenticación OAuth vía navegador
2. ✅ **¡Los 15 modelos de Antigravity se agregan al config!**
3. ✅ `gemini-3-flash` configurado como modelo por defecto
4. ✅ Fallback a `gemini-2.5-flash` configurado

**Salida:**
```
✓ ¡Login de Google Antigravity exitoso!

✓ Se agregaron 15 modelos de Antigravity al config

Modelo por defecto: gemini-3-flash (fallback: gemini-2.5-flash)

Modelos disponibles:
  - gemini-3-flash (por defecto)
  - gemini-3-pro-high, gemini-3-pro-low
  - gemini-3.1-pro-high, gemini-3.1-pro-low, gemini-3.1-flash-lite
  - gemini-3-flash-agent, gemini-3-flash-preview
  - gemini-2.5-flash, gemini-2.5-flash-lite, gemini-2.5-flash-thinking, gemini-2.5-pro
  - claude-sonnet-4-6, claude-opus-4-6-thinking
  - gpt-oss-120b-medium

Pruébalo: picoclaw-agents agent -m "Hola mundo" --model gemini-3-flash
```

---

## Autenticación

### Paso 1: Login (Auto-Config)

```bash
./picoclaw-agents auth login --provider google-antigravity
```

**Novedades (v1.3.0-alpha):**
- 🎉 **Agrega AUTOMÁTICAMENTE los 15 modelos de Antigravity** a `~/.picoclaw/config.json`
- 🎉 **Configura `gemini-3-flash` como modelo por defecto** para todos los agentes
- 🎉 **Configura fallback** a `gemini-2.5-flash`
- 🎉 **¡No requiere edición manual del config!**

**Alias también funciona:**
```bash
./picoclaw-agents auth login --provider antigravity
```

### Paso 2: Completar Flujo OAuth

1. **El navegador se abre automáticamente** (máquinas locales)
2. **Inicia sesión** con tu cuenta de Google (debe tener Google One AI Premium o Workspace Gemini)
3. **Otorga permisos** a PicoClaw
4. **Credenciales guardadas** en `~/.picoclaw/auth.json`
5. **Config actualizado** con los 15 modelos ✨

**Headless/Remoto (VPS/Coolify/Docker):**
1. Ejecuta el comando
2. Copia la URL y ábrela en tu navegador local
3. Completa el login
4. Copia la URL final de redireccionamiento `localhost:51121` de tu navegador
5. Pégala de vuelta en la terminal

### Gestión de Tokens

| Tipo de Token | Duración | Auto-Refresh |
|---------------|----------|--------------|
| `access_token` | ~1 hora | ✅ Sí |
| `refresh_token` | Meses/indefinido | N/A |

**Capas de auto-refresh:**
1. **Daemon en background**: Renueva proactivamente cada 20 min si quedan <30 min
2. **En cada request**: Reintenta con refresh_token incluso si ya expiró
3. **Comando `auth models`**: También se recupera de tokens expirados

**Re-autenticación manual solo si:**
- Revocaste el acceso desde `myaccount.google.com > Seguridad > Apps con acceso`
- Cambiaste tu contraseña de Google
- El `refresh_token` ha estado inactivo por 6+ meses

### Verificar Estado

```bash
./picoclaw-agents auth status
```

### Cerrar Sesión

```bash
./picoclaw-agents auth logout --provider google-antigravity
```

---

## Modelos Disponibles (OAuth Auth)

### Ver Todos los Modelos

```bash
./picoclaw-agents auth models
```

### Lista Completa de Modelos (15 Modelos)

**Auto-configurados al hacer login (v1.3.0-alpha+):**

| # | Nombre del Modelo | Descripción | Mejor Para |
|---|-------------------|-------------|------------|
| 1 | `gemini-3-flash` ⭐ | **POR DEFECTO** - Rápido, confiable | **Recomendado por defecto** |
| 2 | `gemini-3-pro-high` | Alto razonamiento Gemini 3 | Razonamiento complejo |
| 3 | `gemini-3-pro-low` | Bajo razonamiento Gemini 3 | Tareas simples |
| 4 | `gemini-3.1-pro-high` | Alto razonamiento Gemini 3.1 | Tareas avanzadas |
| 5 | `gemini-3.1-pro-low` | Bajo razonamiento Gemini 3.1 | Tareas medias |
| 6 | `gemini-3.1-flash-lite` | Modelo ligero 3.1 | Respuestas rápidas |
| 7 | `gemini-3-flash-agent` | Flash optimizado para agentes | Tareas multi-paso |
| 8 | `gemini-3-flash-preview` | Modelo preview | Testing nuevas features |
| 9 | `gemini-2.5-flash` | Gemini 2.5 Flash | Respuestas rápidas |
| 10 | `gemini-2.5-flash-lite` | Modelo ligero 2.5 | Tareas simples |
| 11 | `gemini-2.5-flash-thinking` | Flash con razonamiento | Tareas de razonamiento |
| 12 | `gemini-2.5-pro` | Gemini 2.5 Pro | Propósito general |
| 13 | `claude-sonnet-4-6` | Claude Sonnet 4.6 | Escritura, análisis |
| 14 | `claude-opus-4-6-thinking` | Claude Opus con thinking | Problemas complejos |
| 15 | `gpt-oss-120b-medium` | Alternativa GPT open-source | Uso general |

> [!NOTE]
> **Nombres de Modelos Actualizados (v1.3.0-alpha)**
>
> Los nombres ahora coinciden exactamente con lo que devuelve `auth models`:
> - ✅ `gemini-3-flash` (simple, coincide con API)
> - ❌ `antigravity-gemini-3-flash` (formato antiguo, aún funciona pero obsoleto)

---

## Ejemplos de Uso

### Línea de Comandos

```bash
# Usar modelo por defecto (gemini-3-flash)
./picoclaw-agents agent -m "Hola"

# Usar modelo específico
./picoclaw-agents agent -m "Hola" --model gemini-3-flash

# Claude para escribir
./picoclaw-agents agent -m "Escribe un poema" --model claude-sonnet-4-6

# Razonamiento complejo
./picoclaw-agents agent -m "Resuelve este problema matemático" --model claude-opus-4-6-thinking

# Respuestas rápidas
./picoclaw-agents agent -m "Pregunta rápida" --model gemini-3.1-flash-lite
```

### Web UI (¡NUEVO!)

**Selector de Modelos (v1.3.0-alpha):**

1. Abre Web UI: http://localhost:18800/
2. Haz clic en el dropdown de modelos en el header
3. Selecciona cualquiera de los 15 modelos de Antigravity
4. Envía mensaje - **¡el cambio de modelo aplica inmediatamente!**

> [!TIP]
> **Web UI Model Override**
>
> ¡El modelo seleccionado en Web UI **ahora sí funciona**! Cada usuario puede usar modelos diferentes independientemente.
>
> Ver [`PICO_MODEL_OVERRIDE.md`](./PICO_MODEL_OVERRIDE.md) para detalles.

---

## Configuración

### Config Automático (v1.3.0-alpha+)

**Después del login, tu config tiene automáticamente:**

```json
{
  "agents": {
    "defaults": {
      "model": "gemini-3-flash",
      "fallbacks": ["gemini-2.5-flash"]
    }
  },
  "model_list": [
    {
      "model_name": "gemini-3-flash",
      "model": "antigravity/gemini-3-flash",
      "auth_method": "oauth"
    },
    {
      "model_name": "gemini-3-pro-high",
      "model": "antigravity/gemini-3-pro-high",
      "auth_method": "oauth"
    },
    // ... 13 modelos más ...
  ]
}
```

### Config Manual (Pre-v1.3.0 o Setup Personalizado)

Si necesitas agregar modelos manualmente:

```bash
# Opción 1: Re-ejecutar login para auto-agregar todos los modelos
./picoclaw-agents auth login --provider google-antigravity

# Opción 2: Usar script de sync
./scripts/sync_antigravity_models.sh

# Opción 3: Usar script de fix (actualiza nombres)
./scripts/fix_antigravity_models.sh
```

### Formato model_list

**Formato correcto (v1.3.0-alpha+):**
```json
{
  "model_name": "gemini-3-flash",
  "model": "antigravity/gemini-3-flash",
  "auth_method": "oauth"
}
```

**Formato antiguo (obsoleto pero funciona):**
```json
{
  "model_name": "antigravity-gemini-3-flash",
  "model": "antigravity/gemini-3-flash",
  "auth_method": "oauth"
}
```

> [!IMPORTANT]
> **Usa Nombres Simples de Modelos**
>
> Siempre usa el nombre simple (ej: `gemini-3-flash`) en:
> - CLI: `--model gemini-3-flash`
> - Web UI: Seleccionar del dropdown
> - Config: campo `model_name`
>
> El prefijo `antigravity/` es solo para el campo `model` internamente.

---

## Generación de Imágenes (Solo API Key)

**Antigravity OAuth NO soporta generación de imágenes.** Para generar imágenes, debes usar **Google AI Studio API Key**.

### Modelos de Imagen Soportados (API Key)

| Modelo | Prefijo Provider | Propósito |
|--------|------------------|-----------|
| `gemini-2.5-flash-image` | `gemini/` | Nano Banana - generación de imágenes |
| `gemini-3-pro-image-preview` | `gemini/` | Nano Banana Pro |
| `gemini-3.1-flash-image-preview` | `gemini/` | Nano Banana 2 |
| `imagen-4.0-generate-001` | `gemini/` | Imagen 4 |
| `imagen-4.0-ultra-generate-001` | `gemini/` | Imagen 4 Ultra |

### Configuración

Agregar a `~/.picoclaw/config.json`:

```json
{
  "model_list": [
    {
      "model_name": "gemini-2.5-flash-image",
      "model": "gemini/gemini-2.5-flash-image",
      "api_key": "TU_GEMINI_API_KEY"  # pragma: allowlist secret
    }
  ],
  "tools": {
    "image_gen": {
      "provider": "gemini",
      "gemini_api_key": "TU_GEMINI_API_KEY",
      "gemini_image_model_name": "gemini-2.5-flash-image",
      "output_dir": "~/.picoclaw/workspace/generated_images"
    }
  }
}
```

**Obtener API Key:** [Google AI Studio](https://aistudio.google.com/app/apikey)

---

## Comandos de Referencia

| Comando | Descripción |
|---------|-------------|
| `./picoclaw-agents auth login --provider google-antigravity` | Login con Antigravity |
| `./picoclaw-agents auth status` | Verificar estado de autenticación |
| `./picoclaw-agents auth models` | Listar modelos disponibles |
| `./picoclaw-agents auth logout --provider google-antigravity` | Cerrar sesión de Antigravity |
| `./picoclaw-agents agent -m "msg" --model <modelo>` | Usar modelo específico |

---

## Solución de Problemas

| Error | Causa | Solución |
|-------|-------|----------|
| `403 PERMISSION_DENIED` | Token expirado/revocado | `./picoclaw-agents auth login --provider google-antigravity` |
| `ACCESS_TOKEN_SCOPE_INSUFFICIENT` | Token expirado/revocado | `./picoclaw-agents auth login --provider google-antigravity` |
| `404 NOT_FOUND` | Alias de modelo no resuelto | Verificar que `model` tiene prefijo `antigravity/` y `auth_method: "oauth"` |
| `401 invalid_api_key` | Proveedor incorrecto | Verificar que `model` tiene prefijo `antigravity/`, no clave OpenAI |
| `429 Rate Limit` | Cuota agotada | Esperar reset mostrado por PicoClaw, o cambiar modelo |
| Respuesta vacía | Modelo restringido para proyecto | Probar `gemini-3-flash` o `gemini-2.5-flash` |
| "Gemini for Google Cloud not enabled" | API no habilitada | Habilitar en [Google Cloud Console](https://console.cloud.google.com) |

---

## Requisitos

- **Cuenta de Google** con:
  - Plan Google One AI Premium, O
  - Add-on Workspace Gemini
- **Proyecto de Google Cloud** con Gemini API habilitado
- **PicoClaw** v1.3.0-alpha o posterior

---

## Enlaces Relacionados

- [`ANTIGRAVITY.md`](./ANTIGRAVITY.md) - Versión en inglés
- [`PICO_MODEL_OVERRIDE.md`](./PICO_MODEL_OVERRIDE.md) - Protocolo de override de modelo
- [`IMAGE_GEN_util.md`](./IMAGE_GEN_util.md) - Guías de generación de imágenes
- [Google Cloud Console](https://console.cloud.google.com) - Gestionar cuotas y billing
- [Google AI Studio](https://aistudio.google.com) - Obtener API keys para generación de imágenes

---

## Historial de Cambios

### v1.3.0-alpha (30 de marzo de 2026)

**Novedades:**
- 🎉 Auto-config de los 15 modelos en login
- 🎉 `gemini-3-flash` como modelo por defecto
- 🎉 Web UI model override funcional
- 🎉 Nombres de modelos simplificados

**Cambios:**
- `auth login --provider google-antigravity` ahora agrega automáticamente:
  - Los 15 modelos de Antigravity al config
  - `gemini-3-flash` como default para todos los agentes
  - Fallback a `gemini-2.5-flash`

---

**Inicio Rápido:** ¡Ejecuta `./picoclaw-agents auth login --provider google-antigravity` ahora! 🚀
