# Social Media Util

Guia operativa rapida para herramientas de redes sociales en PicoClaw.

> **PicoClaw v3.5.0**: ¡Ahora soporta generación de imágenes con **Antigravity OAuth** usando `gemini-3.1-flash-image` — no se necesita API key! `social_post_bundle` ahora genera imágenes vía OAuth por defecto. Ver [IMAGE_GEN_util.es.md](./IMAGE_GEN_util.es.md) y [ANTIGRAVITY_IMAGE_GEN.es.md](./ANTIGRAVITY_IMAGE_GEN.es.md).
>
> **PicoClaw v3.4.1**: Incluye **Comandos Slash Fast-path** para gestión instantánea de lotes y **Global Tracker** para consistencia multi-agente.

## Ejemplos de Uso

### Generar Post de Facebook con Imagen (vía social_post_bundle)

**Usuario (Español):** `genera un post para facebook con imagen sobre peligro nuclear y reloj del juicio final adjunta la imagen`

**Usuario (Inglés):** `Generate a Facebook post with image about nuclear danger and doomsday clock, attach the image`

**Qué pasa:**
1. El agent llama `social_post_bundle` → genera texto vía Antigravity OAuth
2. Genera prompt visual desde el script
3. Genera imagen vía `image_gen_antigravity` (OAuth, sin API key)
4. Copia imagen al directorio del bundle
5. Envía post con imagen adjunta a Telegram/Discord

### Generar Imagen Simple

**Usuario (Español):** `genera una imagen de un pajaro con lentes de sol estilo matrix`

**Usuario (Inglés):** `Generate an image of a bird with sunglasses, Matrix style`

**Qué pasa:**
1. El agent llama `image_gen_antigravity` → genera imagen vía OAuth
2. Envía imagen como foto adjunta a Telegram/Discord

## Configuracion Minima

```json
{
  "tools": {
    "social_media": {
      "facebook": {
        "default_page_id": "TU_FB_PAGE_ID",
        "default_page_token": "TU_FB_PAGE_TOKEN",
        "app_id": "TU_FB_APP_ID",
        "app_secret": "TU_FB_APP_SECRET",
        "user_token": "TU_FB_USER_TOKEN"
      },
      "x": {
        "api_key": "TU_X_API_KEY",
        "api_secret": "TU_X_API_SECRET",
        "access_token": "TU_X_ACCESS_TOKEN",
        "access_token_secret": "TU_X_ACCESS_TOKEN_SECRET"
      },
      "discord": {
        "webhook_url": "https://discord.com/api/webhooks/TU_WEBHOOK_ID/TU_WEBHOOK_TOKEN"
      }
    }
  }
}
```

## Configuracion de Generacion de Imagenes

Para generar imagenes con posts (`social_post_bundle`):

### Antigravity OAuth (Predeterminado — GRATIS — Sin API Key)

Las imagenes via Antigravity OAuth **no requieren API key y NO cuestan un centavo**.

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

Login: `picoclaw auth login --provider google-antigravity`

### Gemini API Key (Fallback — De Pago)

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

### Ideogram API Key (Fallback — De Pago)

```json
{
  "tools": {
    "image_gen": {
      "provider": "ideogram",
      "ideogram_api_key": "TU_API_KEY",
      "ideogram_api_url": "https://api.ideogram.ai/v1/ideogram-v3/generate",
      "aspect_ratio": "1:1",
      "output_dir": "./workspace/image_gen"
    }
  }
}
```

**Prioridad:** Antigravity OAuth (GRATIS) → Gemini API → Ideogram API

## Comportamiento de Facebook

- `facebook_post` soporta:
  - publicacion solo texto
  - publicacion imagen + texto
  - primer comentario opcional
- Si el comentario falla con code `368`, el texto se fusiona al cuerpo del post.
- Si el token expira con code `190` y tienes `app_id/app_secret/user_token`, PicoClaw refresca y reintenta.

## Ejemplos CLI

```bash
# Facebook solo texto
./picoclaw-agents agent -m "Usa facebook_post con message='Hola desde PicoClaw'"

# Facebook imagen + texto
./picoclaw-agents agent -m "Usa facebook_post con message='Actualizacion', image_path='/tmp/post.jpg'"

# Facebook imagen + texto + comentario
./picoclaw-agents agent -m "Usa facebook_post con message='Post principal', image_path='/tmp/post.jpg', comment='Contexto adicional'"

# X solo texto
./picoclaw-agents agent -m "Usa x_post_tweet con message='Hola X'"

# Discord solo texto
./picoclaw-agents agent -m "Usa discord_post con message='Hola Discord'"
```

## Gestión de Lotes (Bundles)

Tras recibir una notificación de lote completado (ej: `#IMA_GEN_...`), puedes usar los comandos rápidos para gestionarlos instantáneamente:

### Comandos de Lotes

- `/bundle_approve id=ID`: Aprueba la publicación y procede a publicar/guardar.
- `/bundle_regen id=ID`: Solicita regeneración completa (imagen y texto).
- `/bundle_edit id=ID`: Edita el texto antes de aprobar.
- `/bundle_publish id=ID platforms=facebook,twitter`: Publica en las plataformas especificadas.

### Comandos de Utilidad

- `/list pending`: Muestra todas las tareas pendientes.
- `/status`: Muestra uso de tokens y estado del sistema.
- `/help`: Muestra ayuda interactiva.
- `/show model`: Muestra modelo activo.
- `/show channel`: Muestra canal activo.

**Beneficios:**
- ✅ **Latencia cero**: Sin razonamiento del LLM, ejecución instantánea
- ✅ **Sintaxis consistente**: Funciona idéntico en Telegram, Discord, CLI
- ✅ **Seguro**: Validación de ID antes de ejecutar

### Global Tracker (v3.4.1+)

El **Global ImageGenTracker** es compartido entre todos los agentes:
- **Subagente genera contenido** → **Agente Principal puede publicar inmediatamente**
- **Sin errores de "ID no encontrado"** entre límites de agentes

Ver [docs/queue_batch.es.md](docs/queue_batch.es.md) para documentación completa.

---

## Prompts en Telegram / Discord / CLI
```text
Publica en Facebook: "Hola desde el bot"
Publica en Facebook con imagen /tmp/post.jpg y mensaje "Nueva actualizacion"
/bundle_approve id=20260302_161740_yiia22  <-- Comando Rápido
```

## Permisos de Facebook

Usa permisos modernos de pagina, no `publish_actions`:

- `pages_manage_posts`
- `pages_read_engagement`
- `pages_show_list`
- opcional: `pages_manage_engagement`
