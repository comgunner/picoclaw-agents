# Generación de Imágenes — Guía de Uso

Guía rápida para usar las herramientas de generación de imágenes en PicoClaw desde terminal y Telegram.

> **PicoClaw v3.5.0**: ¡Ahora soporta generación de imágenes con **Antigravity OAuth** usando `gemini-3.1-flash-image` — **no se necesita API key, completamente GRATIS!** Solo login con tu cuenta de Google. Incluye **cooldown obligatorio de 150s** para protección anti-ban. Ver [ANTIGRAVITY_IMAGE_GEN.es.md](./ANTIGRAVITY_IMAGE_GEN.es.md).
>
> **PicoClaw v3.4.1**: Incluye **Fast-path Slash Commands** para gestión instantánea de bundles y **Global Tracker** para consistencia multi-agente.
>
> **⚠️ IMPORTANTE: ESTO ES PARA GENERAR IMÁGENES ESTÁTICAS, NO VIDEOS**
>
> - `text_script_create`: Genera **TEXTO PARA POSTS** de Facebook/Twitter (como copy de post)
> - `image_gen_create`: Genera **IMÁGENES ESTÁTICAS** desde texto
> - `image_gen_antigravity`: Genera imágenes vía **Antigravity OAuth** (predeterminado, sin API key, GRATIS)
> - **NO hay generación de video** en esta herramienta

---

## 💸 Generación de Imágenes GRATIS con Antigravity OAuth

**Las imágenes generadas vía Antigravity OAuth NO requieren API key y NO te cuestan un centavo.** Usan la cuota gratuita de tu cuenta de Google a través del servicio Cloud Code Assist. Sin configuración de facturación, sin tarjeta de crédito, sin cargos.

| Método | Costo | ¿Requiere API Key? |
|--------|-------|-------------------|
| **Antigravity OAuth** (predeterminado) | ✅ **GRATIS** | ❌ No — solo login de Google |
| Gemini API | 💰 Facturación por token | ✅ Sí — requiere facturación en Google Cloud |
| Ideogram API | 💰 Facturación por imagen | ✅ Sí — requiere plan pagado |

**Recomendación:** Usa Antigravity OAuth como predeterminado. Es gratis y funciona inmediatamente.

---

## 🆕 Generación de Imágenes con Antigravity OAuth (Recomendado — Predeterminado — GRATIS)

Desde v3.5.0, el método **predeterminado** usa **Google Antigravity** vía OAuth:
- **No se necesita API key** — solo login con tu cuenta de Google
- **Modelo:** `gemini-3.1-flash-image`
- **Costo:** ✅ **GRATIS** — sin API key, sin facturación, sin cargos
- **Cooldown:** 150s (2.5 min) obligatorio después de cada generación (anti-ban)

### Configuración

```bash
picoclaw auth login --provider google-antigravity
```

Esto configura automáticamente:
- `provider: "antigravity"`
- `antigravity_model: "gemini-3.1-flash-image"`
- `cooldown_seconds: 150`

---

## Opciones de Configuración

### Opción A: Antigravity OAuth (Predeterminado — Recomendado — GRATIS)

**No se necesita API key. Sin costo.** Solo login OAuth con tu cuenta de Google.

```json
{
  "tools": {
    "image_gen": {
      "provider": "antigravity",
      "antigravity_model": "gemini-3.1-flash-image",
      "cooldown_seconds": 150,
      "aspect_ratio": "1:1",
      "output_dir": "./workspace/image_gen"
    }
  }
}
```

**Comando de login:**
```bash
picoclaw auth login --provider google-antigravity
```

---

### Opción B: Gemini API Key (Fallback — De Pago)

Usa si tienes una API key de Gemini y prefieres acceso directo por API. **Requiere facturación en Google Cloud Console.**

```json
{
  "tools": {
    "image_gen": {
      "provider": "gemini",
      "gemini_api_key": "TU_API_KEY",
      "gemini_text_model_name": "gemini-3-flash-agent",
      "gemini_image_model_name": "gemini-2.5-flash-image",
      "aspect_ratio": "1:1",
      "output_dir": "./workspace/image_gen"
    }
  }
}
```

**Costo:** Facturación por token vía Google Cloud Console.

---

### Opción C: Ideogram API Key (Fallback — De Pago)

Usa si tienes una API key de Ideogram. **Requiere plan pagado de Ideogram.**

```json
{
  "tools": {
    "image_gen": {
      "provider": "ideogram",
      "ideogram_api_key": "TU_API_KEY",
      "ideogram_api_url": "https://api.ideogram.ai/v1/ideogram-v3/generate",
      "ideogram_model_name": "V_3_TURBO",
      "ideogram_aspect_ratio": "4x5",
      "ideogram_rendering_speed": "TURBO",
      "ideogram_style_type": "REALISTIC",
      "aspect_ratio": "1:1",
      "output_dir": "./workspace/image_gen"
    }
  }
}
```

**Costo:** Facturación por imagen vía suscripción de Ideogram.

---

### Configuración Completa (Todos los Proveedores)

```json
{
  "tools": {
    "image_gen": {
      "provider": "antigravity",
      "antigravity_model": "gemini-3.1-flash-image",
      "cooldown_seconds": 150,
      "gemini_api_key": "",
      "gemini_text_model_name": "gemini-3-flash-agent",
      "gemini_image_model_name": "gemini-2.5-flash-image",
      "ideogram_api_key": "",
      "ideogram_api_url": "https://api.ideogram.ai/v1/ideogram-v3/generate",
      "aspect_ratio": "1:1",
      "output_dir": "./workspace/image_gen"
    }
  }
}
```

**Orden de prioridad:** Antigravity OAuth (predeterminado, GRATIS) → Gemini API key (de pago) → Ideogram API key (de pago)

---

## Ejemplos de Uso

### Ejemplo 1: Generar una imagen (GRATIS vía OAuth)

**Usuario:** `Generate an image of a cute cat wearing sunglasses`

**Qué pasa:** El agent llama `image_gen_antigravity` → genera imagen vía OAuth (GRATIS) → envía como foto adjunta a Telegram/Discord.

### Ejemplo 2: Post de Facebook con imagen (GRATIS vía OAuth)

**Usuario (Español):** `genera un post para facebook con imagen sobre peligro nuclear y reloj del juicio final adjunta la imagen`

**Usuario (Inglés):** `Generate a Facebook post with image about nuclear danger and doomsday clock, attach the image`

**Qué pasa:** El agent llama `social_post_bundle` → genera texto → genera imagen vía Antigravity OAuth (GRATIS) → copia imagen al directorio del bundle → envía post con imagen adjunta.

### Ejemplo 3: Generación simple de imagen (GRATIS vía OAuth)

**Usuario (Español):** `genera una imagen de un pajaro con lentes de sol estilo matrix`

**Usuario (Inglés):** `Generate an image of a bird with sunglasses, Matrix style`

**Qué pasa:** El agent llama `image_gen_antigravity` → genera imagen vía OAuth (GRATIS) → envía como foto adjunta.

---

## Herramientas Disponibles

### Reglas de Enrutamiento (Importante)

- Si el usuario solo pide una imagen (ejemplo: `Generate an image of...`), usa solo `image_gen_create` o `image_gen_antigravity`.
- Usa `text_script_create` solo cuando el usuario pida explícitamente texto para post, o pida flujo Script → Imagen.
- `prompt_base_img.txt` se usa para construir un prompt visual desde un script existente (no para solicitudes solo de imagen).

### `text_script_create`

Genera **TEXTO PARA POSTS DE REDES SOCIALES** (Facebook, Twitter, Discord).

**NO es para video** — este es el texto que acompaña una imagen en un post.

### `image_gen_antigravity`

Genera imágenes vía **Antigravity OAuth** (predeterminado). **GRATIS — sin API key.**

### `image_gen_create`

Genera imágenes vía **Gemini API** o **Ideogram API** (fallback si OAuth no está configurado).

### `image_gen_workflow`

Flujo de script-a-imagen. Primero genera un script de texto, luego crea una imagen coincidente.

---

## Aspect Ratios Soportados

| Ratio | Uso |
|-------|-----|
| `1:1` | Predeterminado, cuadrado |
| `16:9` | Widescreen, desktop |
| `9:16` | Stories, móvil vertical |
| `4:5` | Instagram portrait |
| `3:4` | Retrato clásico |

---

## Solución de Problemas

### "No hay credenciales OAuth de Antigravity"

Ejecuta el comando de login:
```bash
picoclaw auth login --provider google-antigravity
```

### "Rate Limited (429)"

Antigravity tiene una cuota compartida para TODAS las llamadas (chat + imágenes). Espera ~5 minutos y reintenta. El sistema tiene retry automático con backoff exponencial (5s → 15s → 30s → 60s → 120s).

### "Method doesn't allow unregistered callers" (403)

Esto significa que estás intentando usar Gemini API sin API key. Cambia a Antigravity OAuth (predeterminado) que es GRATIS y no requiere API key.

### "Tool definitions exceed budget threshold"

Configura `context_window: 128000` en `~/.picoclaw/config.json` y reinicia el gateway.

### "no such file or directory" para la imagen

Este es un bug de mismatch de workspaces — corregido en v3.5.0 con `copyFile` automático. Asegúrate de usar el binario más reciente.

---

## Referencias

- [ANTIGRAVITY_IMAGE_GEN.es.md](./ANTIGRAVITY_IMAGE_GEN.es.md) — Guía completa de Antigravity OAuth
- [SOCIAL_MEDIA.es.md](./SOCIAL_MEDIA.es.md) — Publicación en redes sociales
- [SOCIAL_MEDIA_util.es.md](./SOCIAL_MEDIA_util.es.md) — Utilidades de redes sociales
