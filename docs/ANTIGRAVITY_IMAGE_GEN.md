# Antigravity Image Generation — Gemini 3.1 Flash Image

Image generation using Google Antigravity via OAuth with built-in **anti-ban cooldown**.

> **💸 FREE — No API key required.** Images generated via Antigravity OAuth use your Google account's free quota through Cloud Code Assist. No billing setup, no credit card, no charges.

## 🚀 Quick Start

### 1. Login OAuth

```bash
picoclaw auth login --provider google-antigravity
```

### 2. Verificar Modelos Disponibles

```bash
picoclaw auth status --provider google-antigravity
```

Debes ver `gemini-3.1-flash-image` en la lista con estado ✓.

### 3. Configurar

En `~/.picoclaw/config.json`:

```json
{
  "tools": {
    "image_gen": {
      "provider": "antigravity",
      "antigravity_model": "gemini-3.1-flash-image",
      "aspect_ratio": "1:1",
      "cooldown_seconds": 300,
      "output_dir": "./workspace/image_gen"
    }
  }
}
```

O automáticamente después de hacer login:

```bash
picoclaw auth login --provider google-antigravity
# Configura automáticamente: provider="antigravity", model="gemini-3.1-flash-image"
```

## ⏳ Cooldown Anti-Ban

### ¿Por qué Cooldown?

El API de Antigravity tiene un límite de **~1-2 imágenes cada 10 minutos**. Sin protección:

1. El agent intenta generar múltiples imágenes seguidas
2. Recibe errores 429 (Rate Limited) repetidamente
3. Puede ser baneado temporalmente por abuso de API

### Cómo Funciona

```
✅ Imagen generada exitosamente
   ↓
⏳ Cooldown: 300s (5 min) — BLOQUEO TOTAL
   ↓
✓ Cooldown completo → siguiente imagen permitida
```

### Límites Confirmados

| Límite | Valor |
|--------|-------|
| **Imágenes por 10 min** | ~1-2 |
| **Cooldown obligatorio** | 300s (5 min) mínimo |
| **Reintentos en 429** | ✅ Backoff exponencial |

### Retry Delays (Backoff Exponencial)

```
Intento 1 → espera 30s  (0.5 min)
Intento 2 → espera 60s  (1 min)
Intento 3 → espera 120s (2 min)
Intento 4 → espera 300s (5 min)
Intento 5 → espera 600s (10 min)
```

**Total máximo de espera:** 1110s = 18.5 minutos

### Configurar Cooldown

**En config.json:**

```json
{
  "tools": {
    "image_gen": {
      "cooldown_seconds": 300
    }
  }
}
```

**Variables de entorno:**

```bash
# Default: 300 (5 minutos)
PICOCLAW_IMAGE_COOLDOWN_SECONDS=300

# Ruta del archivo SQLite (default: <workspace>/tmp/picoclaw_image_cooldown.db)
PICOCLAW_IMAGE_COOLDOWN_DB=~/.picoclaw/workspace/tmp/picoclaw_image_cooldown.db
```

**Mínimo:** 60 segundos (por debajo de esto se usa 60 automáticamente)

### Mensaje de Cooldown Activo

Cuando intentas generar durante cooldown:

```
⏳ **Cooldown activo** — No se puede generar imagen ahora.

⏱ **Tiempo restante:** 2m 30s
🔐 **Provider:** antigravity

💡 El cooldown es obligatorio (300s) para proteger contra rate limits.
💡 Mientras tanto, puedes:
  - Editar el prompt y reintentar después
  - Usar `image_gen_create` con provider 'gemini' o 'ideogram' si tienes API keys
```

### Mensaje de Rate Limit (429)

```
⏳ **Rate Limited (429)** — La API de Antigravity tiene un límite de ~1-2 imágenes cada 10 minutos.

❌ Error: 429 Rate Limited

💡 Reintentos automáticos con delays: 30s, 60s, 120s, 300s, 600s
💡 Espera 5-10 minutos y vuelve a intentar
💡 O usa `image_gen_create` con provider 'gemini' o 'ideogram'
```

### Mensaje de Éxito

```
✅ **Imagen generada exitosamente (Antigravity OAuth)**

📁 **Archivos creados:**
┌─────────────────────────────────────────────────────────────────┐
│ Imagen:   ./workspace/image_gen/abc123/abc123.-imagen.jpg
│ Prompt:   ./workspace/image_gen/abc123/abc123.-prompt_visual.txt
│ Script:   ./workspace/image_gen/abc123/abc123.-script.txt
└─────────────────────────────────────────────────────────────────┘

🎨 **Modelo:** gemini-3.1-flash-image
📐 **Aspect Ratio:** 1:1
🔐 **Auth:** OAuth (google-antigravity)
📊 **Tamaño:** 1.2 MB

⏳ **Cooldown activado:** 300s antes de la siguiente imagen (anti-ban protection)
```

## 📐 Aspect Ratios Soportados

| Ratio | Uso |
|-------|-----|
| `1:1` | Default, cuadrado |
| `16:9` | Widescreen, desktop |
| `9:16` | Stories, móvil vertical |
| `4:5` | Instagram portrait |
| `3:4` | Retrato clásico |

## 🤖 Uso desde Agent

### Tool: `image_gen_antigravity`

```json
{
  "name": "image_gen_antigravity",
  "arguments": {
    "prompt": "A cyberpunk city at night with neon lights",
    "aspect_ratio": "16:9"
  }
}
```

### Parámetros

| Parámetro | Tipo | Requerido | Default | Descripción |
|-----------|------|-----------|---------|-------------|
| `prompt` | string | ✅ | — | Descripción de la imagen |
| `model` | string | ❌ | `gemini-3.1-flash-image` | Modelo a usar |
| `aspect_ratio` | string | ❌ | `1:1` | Proporción de la imagen |
| `script_path` | string | ❌ | — | Ruta al script visual (opcional) |

## 🔄 Fallback Chain

```
1. antigravity (OAuth + cooldown) ← DEFAULT
   ↓ si no hay credenciales OAuth
2. gemini (API key, sin cooldown)
   ↓ si falla o modelo no encontrado
3. ideogram (API key, sin cooldown)
```

> ⚠️ **Nota:** Solo antigravity tiene cooldown obligatorio. Gemini e Ideogram no tienen cooldown.

## 📋 Modelos de Imagen Disponibles vía Antigravity

| Modelo | Descripción |
|--------|-------------|
| `gemini-3.1-flash-image` | **Nuevo** — Modelo de imagen más reciente |
| `gemini-2.5-flash-image` | Legacy — Funciona como fallback |

> ⚠️ **Importante:** Los modelos de imagen deben aparecer en `picoclaw auth status` con estado ✓ para funcionar.

## 🆚 Comparación: Antigravity vs Gemini API Key

| Característica | Antigravity OAuth | Gemini API Key |
|----------------|-------------------|----------------|
| **Autenticación** | OAuth (sesión Google) | API Key |
| **Setup** | `auth login` | Configurar key en config.json |
| **Costo** | Gratuito (incluido en cuenta Google) | Requiere billing en Google Cloud |
| **Modelo** | `gemini-3.1-flash-image` | `gemini-2.5-flash-image` |
| **Cooldown** | ✅ Obligatorio (300s) | ❌ No |
| **Rate Limit** | ~1-2 img/10 min | Según cuota de API |
| **Endpoint** | `cloudcode-pa.googleapis.com` | `generativelanguage.googleapis.com` |

## 🛠️ Cooldown Storage

El cooldown se persiste en **SQLite** para sobrevivir reinicios del proceso:

```
Archivo: <workspace>/tmp/picoclaw_image_cooldown.db
         Ejemplo: ~/.picoclaw/workspace/tmp/picoclaw_image_cooldown.db
Tabla: cooldown
Campos: id, started_at, duration_seconds, provider, model
```

> ✅ **Ventaja del workspace:** El archivo siempre está dentro del directorio del agent,
> que tiene permisos garantizados y no será borrado por limpieza automática de `/tmp/` del SO.

### Comandos útiles

**Ver estado del cooldown:**
```bash
# Via tool (desde el agent)
# El agent reporta automáticamente: "⏳ Cooldown: 120s restantes"
```

**Limpiar cooldown manualmente** (no recomendado):
```bash
rm ~/.picoclaw/workspace/tmp/picoclaw_image_cooldown.db
# O el path de tu workspace:
rm <tu_workspace>/tmp/picoclaw_image_cooldown.db
```

## ❓ Troubleshooting

### "No Antigravity OAuth credentials"
```bash
picoclaw auth login --provider google-antigravity
```

### "Token OAuth expirado"
```bash
picoclaw auth login --provider google-antigravity
# O el auto-refresh lo hace automáticamente
```

### "Rate Limited (429)"
- Espera 5-10 minutos
- El sistema hace retry automático con delays exponenciales
- O usa `image_gen_create` con provider 'gemini' o 'ideogram'

### "Cooldown activo"
- Espera el tiempo indicado (default 300s)
- El cooldown es obligatorio para proteger tu cuenta
- No se puede saltar (por seguridad)

### "Model not found"
Verifica que `gemini-3.1-flash-image` aparezca en:
```bash
picoclaw auth status --provider google-antigravity
```

## 🔗 Referencias

- [Antigravity Auth Guide](./ANTIGRAVITY.md)
- [Image Gen Utility](./IMAGE_GEN_util.md)
