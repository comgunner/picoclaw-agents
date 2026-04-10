# Generación de Imágenes con Antigravity — Gemini 3.1 Flash Image

Generación de imágenes usando Google Antigravity vía OAuth con **cooldown anti-ban** integrado.

> **💸 GRATIS — No se necesita API key.** Las imágenes generadas vía Antigravity OAuth usan la cuota gratuita de tu cuenta de Google a través de Cloud Code Assist. Sin configuración de facturación, sin tarjeta de crédito, sin cargos.

## 🚀 Inicio Rápido

### 1. Login OAuth

```bash
picoclaw auth login --provider google-antigravity
```

### 2. Verificar Modelos Disponibles

```bash
picoclaw auth status
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
      "cooldown_seconds": 150,
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
⏳ Cooldown: 150s (2.5 min) — BLOQUEO TOTAL
   ↓
✓ Cooldown completo → siguiente imagen permitida
```

### Límites Confirmados

| Límite | Valor |
|--------|-------|
| **Imágenes por 10 min** | ~1-2 |
| **Cooldown obligatorio** | 150s (2.5 min) mínimo |
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

### Mensaje de Cooldown Activo

Cuando intentas generar durante cooldown:

```
⏳ **Cooldown activo** — No se puede generar imagen ahora.

⏱ **Tiempo restante:** 2m 30s
🔐 **Provider:** antigravity

💡 El cooldown es obligatorio (150s) para proteger contra rate limits.
💡 Mientras tanto, puedes:
  - Editar el prompt y reintentar después
  - Usar `image_gen_create` con provider 'gemini' o 'ideogram' si tienes API keys
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

⏳ **Cooldown activado:** 150s antes de la siguiente imagen (anti-ban protection)
```

## 📐 Aspect Ratios Soportados

| Ratio | Uso |
|-------|-----|
| `1:1` | Predeterminado, cuadrado |
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
    "prompt": "Una ciudad cyberpunk de noche con luces de neón",
    "aspect_ratio": "16:9"
  }
}
```

### Parámetros

| Parámetro | Tipo | Requerido | Predeterminado | Descripción |
|-----------|------|-----------|----------------|-------------|
| `prompt` | string | ✅ | — | Descripción de la imagen |
| `model` | string | ❌ | `gemini-3.1-flash-image` | Modelo a usar |
| `aspect_ratio` | string | ❌ | `1:1` | Proporción de la imagen |
| `script_path` | string | ❌ | — | Ruta al script visual (opcional) |

## 🔄 Cadena de Fallback

```
1. antigravity (OAuth + cooldown) ← PREDETERMINADO
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
| **Cooldown** | ✅ Obligatorio (150s) | ❌ No |
| **Rate Limit** | ~1-2 img/10 min | Según cuota de API |
| **Endpoint** | `cloudcode-pa.googleapis.com` | `generativelanguage.googleapis.com` |

## 🛠️ Almacenamiento del Cooldown

El cooldown se persiste en **SQLite** para sobrevivir reinicios del proceso:

```
Archivo: <workspace>/tmp/picoclaw_image_cooldown.db
         Ejemplo: ~/.picoclaw/workspace/tmp/picoclaw_image_cooldown.db
Tabla: cooldown
Campos: id, started_at, duration_seconds, provider, model
```

> ✅ **Ventaja del workspace:** El archivo siempre está dentro del directorio del agent,
> que tiene permisos garantizados y no será borrado por limpieza automática de `/tmp/` del SO.

## ❓ Solución de Problemas

### "No hay credenciales OAuth de Antigravity"
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
- Espera el tiempo indicado (predeterminado 150s)
- El cooldown es obligatorio para proteger tu cuenta
- No se puede saltar (por seguridad)

### "Modelo no encontrado"
Verifica que `gemini-3.1-flash-image` aparezca en:
```bash
picoclaw auth status
```

## 🔗 Referencias

- [Guía de Antigravity Auth](./ANTIGRAVITY.md)
- [Utilidad de Generación de Imágenes](./IMAGE_GEN_util.es.md)
