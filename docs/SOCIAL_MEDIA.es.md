# Integracion de Redes Sociales

PicoClaw incluye herramientas nativas para publicar en Facebook, X (Twitter) y Discord.

> **PicoClaw v3.5.0**: ¡Ahora soporta generación de imágenes con **Antigravity OAuth** usando `gemini-3.1-flash-image` — no se necesita API key! `social_post_bundle` ahora genera imágenes vía OAuth por defecto. Ver [IMAGE_GEN_util.es.md](./IMAGE_GEN_util.es.md) y [ANTIGRAVITY_IMAGE_GEN.es.md](./ANTIGRAVITY_IMAGE_GEN.es.md).
>
> **PicoClaw v3.4.1**: Incluye **Comandos Slash Fast-path** para gestión instantánea de lotes y **Global Tracker** para consistencia multi-agente.

## Herramientas Soportadas

- `facebook_post`: Publica en Facebook Page (solo texto o imagen + texto, comentario opcional)
- `x_post_tweet`: Publica en X (solo texto o imagen, reply opcional)
- `discord_post`: Publica en webhook de Discord (solo texto o imagen, username opcional)

## Arquitectura de Workspaces (CRÍTICO)

### Resumen Ejecutivo

PicoClaw utiliza **dos niveles de workspaces**:

| Tipo | Ruta | Propósito | Aislamiento |
|------|------|----------|-------------|
| **Worker Isolated Workspace** | `~/.picoclaw/workspace-{worker_id}/` | Datos privados del worker (logs, caché temporal) | ✅ Exclusivo del worker |
| **Global Shared Workspace** | `~/.picoclaw/workspace/` | Repositorio persistente de activos (imágenes, tracker, scripts) | 🔓 Accesible por todos los workers |

⚠️ **ADVERTENCIA CRÍTICA**: Las herramientas que publican en redes sociales (como `social_manager`, `social_post_bundle`) **DEBEN** acceder a imágenes desde el **Global Shared Workspace**, no del workspace aislado del worker.

### Estructura de Directorios

```
~/.picoclaw/
├── config.json                           ← Configuración global
├── workspace/                            ← 🌍 GLOBAL SHARED
│   ├── image_gen/
│   │   ├── image_id_1/
│   │   │   └── image.jpg
│   │   ├── image_id_2/
│   │   │   └── image.jpg
│   │   └── tracker.json                  ← Tracker compartido
│   ├── text_scripts/
│   │   ├── script_id_1.md
│   │   └── script_id_2.md
│   └── other_assets/
│
├── workspace-general_worker/             ← 🔒 AISLADO (Worker A)
│   ├── logs/
│   ├── cache/
│   └── temp/
│
├── workspace-project_manager/            ← 🔒 AISLADO (Worker B)
│   ├── logs/
│   ├── cache/
│   └── temp/
│
└── workspace-{other_workers}/            ← 🔒 AISLADO (Worker N)
```

### Regla de Oro para Desarrolladores

```
📌 REGLA: Si tu herramienta necesita acceder a ACTIVOS PERSISTENTES 
          (imágenes, tracker.json, scripts), SIEMPRE busca en:
          
          ~/.picoclaw/workspace/  ← Global Shared Workspace
          
          NO uses: ~/.picoclaw/workspace-{worker_id}/
```

### Ejemplo de Implementación Correcta

**❌ INCORRECTO** (buscará en workspace aislado):
```go
// ⚠️ Busca en la ruta aislada del worker, va a fallar
baseDir := filepath.Dir(t.tracker.TrackerPath)  // Puede ser workspace-general_worker
searchPath := filepath.Join(baseDir, record.ID, "*.jpg")
```

**✅ CORRECTO** (normaliza a workspace global):
```go
// ✅ Normaliza la ruta eliminando sufijo de worker
func (t *SocialManagerTool) normalizeTrackerPath(trackerPath string) string {
    dir := filepath.Dir(trackerPath)
    // Detecta "workspace-{id}" y lo reemplaza con "workspace"
    // /home/.picoclaw/workspace-general_worker/image_gen → /home/.picoclaw/workspace/image_gen
    ...
}

baseDir := t.normalizeTrackerPath(t.tracker.TrackerPath)  // ✅ Siempre apunta a global
searchPath := filepath.Join(baseDir, record.ID, "*.jpg")
```

### Casos de Uso Prácticos

#### Caso 1: Publicar Imagen Generada
```
1. Worker "imagen_gen" → Genera imagen
2. Almacena en: ~/.picoclaw/workspace/image_gen/abc123/image.jpg ✅
3. Guarda metadata en: ~/.picoclaw/workspace/image_gen/tracker.json ✅
4. Worker "social_manager" → Lee desde GLOBAL SHARED:
   └─ ~/.picoclaw/workspace/image_gen/abc123/image.jpg ✅
```

#### Caso 2: ❌ Error Común (Se llama después de esta fix)
```
❌ ANTES DEL FIX:
1. Imagen generada: ~/.picoclaw/workspace/image_gen/abc123/image.jpg
2. Worker "social_manager" (con sufijo -general_worker) busca en:
   └─ ~/.picoclaw/workspace-general_worker/image_gen/abc123/image.jpg ❌ (No existe)
3. Falla la publicación sin imagen

✅ DESPUÉS DEL FIX:
1. Imagen generada: ~/.picoclaw/workspace/image_gen/abc123/image.jpg
2. Worker "social_manager" normaliza la ruta:
   └─ ~/.picoclaw/workspace/image_gen/abc123/image.jpg ✅ (Encontrada)
3. Publicación exitosa con imagen
```

### Acceso a Activos Compartidos desde Tools

**Para leer/escribir activos compartidos en tu herramienta:**

```go
// Opción 1: Normalizar desde tracker (recomendado)
globalWorkspace := t.normalizeTrackerPath(t.tracker.TrackerPath)
assetPath := filepath.Join(globalWorkspace, "image_gen", assetID, "image.jpg")

// Opción 2: Si tienes acceso a la configuración
configWorkspace := cfg.Agents.Defaults.Workspace  // ~/.picoclaw/workspace
assetPath := filepath.Join(configWorkspace, "image_gen", assetID, "image.jpg")
```

### Por Qué Existe Esta Arquitectura

1. **Aislamiento**: Cada worker tiene su propio espacio para logs/caché sin interferencia
2. **Persistencia**: Los activos (imágenes) se almacenan en un lugar central accesible
3. **Escalabilidad**: Múltiples workers pueden generar y consumir activos del mismo repositorio
4. **Recovery**: Si un worker falla, sus datos aislados no afectan otros workers ni los activos compartidos

---

## Modelo de Token de Facebook (Importante)

Facebook ya no soporta `publish_actions`.
Debes usar **Page Access Token** con permisos modernos:

- `pages_manage_posts`
- `pages_read_engagement`
- `pages_show_list`
- opcional para moderacion: `pages_manage_engagement`

PicoClaw soporta dos modos:

1. Token de pagina directo
2. Refresh automatico de token (si expira con code `190`) usando:
  - `app_id`
  - `app_secret`
  - `user_token`

## Obtener Credenciales - Instrucciones Paso a Paso

### Facebook: Paso 1 - Crear Meta App

1. Ve a: https://developers.facebook.com/apps
2. Haz clic en "Crear aplicación" (Create App)
3. Selecciona tipo de uso: "Otro" (Other)
4. Elige "Empresa" (Business) o "Consumidor" según tu caso
5. Completa:
   - Nombre de la app
   - Correo electrónico de contacto
6. Ve al menú lateral: Configuración de la app > Básica
7. Copia:
   - **Identificador de la app (App ID)**
   - **Clave secreta de la app (App Secret)**

### Facebook: Paso 2 - Configurar OAuth

1. En el panel, busca "Añadir producto"
2. Agrega: "Inicio de sesión con Facebook" (Facebook Login)
3. Ve a: Inicio de sesión con Facebook > Configuración
4. Campo "URI de redireccionamiento OAuth válidos":
   - Déjalo VACÍO para flujo manual (recomendado)
   - O pon la URL de tu sitio si lo requiere
5. Ve a: Roles de la app > Roles
6. Agrega cuentas como Administradores o Desarrolladores para pruebas

### Facebook: Paso 3 - Generar Page Access Token

1. Ve a: https://developers.facebook.com/tools/explorer
2. Selecciona tu app del dropdown
3. Haz clic en "Generate Access Token"
4. Otorga permisos:
   - `pages_manage_posts` (requerido)
   - `pages_read_engagement` (opcional)
   - `pages_show_list` (opcional)
5. Selecciona la página que quieres administrar
6. Copia el **Page Access Token** (empieza con `EAAB...`)

### Facebook: Paso 4 - Obtener Page ID

1. Ve a tu página de Facebook
2. Haz clic en "Información" (About) en el menú izquierdo
3. Baja hasta encontrar **Page ID** (número de 15-16 dígitos)
4. Alternativa: https://developers.facebook.com/tools/explorer → `GET /me/accounts`

### X (Twitter): Paso 1 - Crear Cuenta Desarrollador

1. Ve a: https://developer.twitter.com/
2. Inicia sesión con tu cuenta de X
3. Solicita cuenta de desarrollador si no tienes

### X (Twitter): Paso 2 - Crear Proyecto y App

1. En developer.twitter.com, haz clic en "Create Project"
2. Ingresa:
   - Nombre del proyecto
   - Descripción
   - Caso de uso (ej: "Posting tweets")
3. Después de crear el proyecto, haz clic en "Create App"
4. Ingresa nombre de la app

### X (Twitter): Paso 3 - Configurar OAuth 2.0

1. En Settings de tu app, busca "User authentication settings"
2. Haz clic en "Set up"
3. Activa OAuth 2.0
4. Elige tipo de app: Web App, Native App, etc.
5. Configura URI de redireccionamiento si es necesario
6. Guarda:
   - **Client ID**
   - **Client Secret**

### X (Twitter): Paso 4 - Obtener API Keys y Tokens

1. En el dashboard de tu app, ve a "Keys and Tokens"
2. En **API Key & Secret**:
   - Haz clic en "Generate" para crear:
     - **API Key**
     - **API Secret**
3. En **Access Token & Secret**:
   - Haz clic en "Generate" para crear:
     - **Access Token**
     - **Access Token Secret**
4. Asegúrate de habilitar permisos de **Lectura y Escritura** (Read and Write)

### X (Twitter): Paso 5 - Verificar Permisos

1. Ve a "App Permissions" en settings de tu app
2. Asegúrate que estén configurados como **Read and Write** (no solo Read)
3. Si cambiaste permisos, regenera tus tokens

### Discord: Obtener Webhook URL

1. Ve a tu canal de Discord
2. Haz clic en el ícono de engranaje (Editar canal)
3. Ve a: Integraciones > Webhooks > Nuevo Webhook
4. Copia la URL del webhook

## Configuracion

Actualiza `~/.picoclaw/config.json`:

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

## Configuración de Generación de Imágenes

Para generar imágenes con posts (`social_post_bundle`), configura el proveedor de imágenes:

### Opción A: Antigravity OAuth (Predeterminado — GRATIS — Sin API Key)

Las imágenes generadas vía Antigravity OAuth **no requieren API key y NO te cuestan un centavo**. Usan la cuota gratuita de tu cuenta de Google.

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

**Login requerido:**
```bash
picoclaw auth login --provider google-antigravity
```

### Opción B: Gemini API Key (Fallback — De Pago)

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

### Opción C: Ideogram API Key (Fallback — De Pago)

```json
{
  "tools": {
    "image_gen": {
      "provider": "ideogram",
      "ideogram_api_key": "TU_API_KEY",
      "ideogram_api_url": "https://api.ideogram.ai/v1/ideogram-v3/generate",
      "ideogram_model_name": "V_3_TURBO",
      "aspect_ratio": "1:1",
      "output_dir": "./workspace/image_gen"
    }
  }
}
```

**Orden de prioridad:** Antigravity OAuth (predeterminado, GRATIS) → Gemini API key (de pago) → Ideogram API key (de pago)

## Variables de Entorno

```bash
# Facebook
export PICOCLAW_TOOLS_SOCIAL_FACEBOOK_DEFAULT_PAGE_ID="tu_page_id"
export PICOCLAW_TOOLS_SOCIAL_FACEBOOK_DEFAULT_PAGE_TOKEN="tu_page_token"
export PICOCLAW_TOOLS_SOCIAL_FACEBOOK_APP_ID="tu_app_id"
export PICOCLAW_TOOLS_SOCIAL_FACEBOOK_APP_SECRET="tu_app_secret"
export PICOCLAW_TOOLS_SOCIAL_FACEBOOK_USER_TOKEN="tu_user_token"

# X
export PICOCLAW_TOOLS_SOCIAL_X_API_KEY="tu_api_key"
export PICOCLAW_TOOLS_SOCIAL_X_API_SECRET="tu_api_secret"
export PICOCLAW_TOOLS_SOCIAL_X_ACCESS_TOKEN="tu_access_token"
export PICOCLAW_TOOLS_SOCIAL_X_ACCESS_TOKEN_SECRET="tu_access_token_secret"

# Discord
export DISCORD_WEBHOOK_URL="tu_webhook_url"
```

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

### Uso en Terminal

```bash
# Facebook solo texto
./picoclaw-agents agent -m "Usa facebook_post con message='¡Hola desde PicoClaw!'"

# Facebook imagen + texto
./picoclaw-agents agent -m "Usa facebook_post con message='Actualización de lanzamiento', image_path='/tmp/post.jpg'"

# Facebook imagen + texto + comentario
./picoclaw-agents agent -m "Usa facebook_post con message='Actualización principal', image_path='/tmp/post.jpg', comment='Detalles extra'"

# Facebook multi-página
./picoclaw-agents agent -m "Usa facebook_post con page_id='123456789', page_token='EAAB...', message='Actualización específica de página'"

# X solo texto
./picoclaw-agents agent -m "Usa x_post_tweet con message='Hola X'"

# X con imagen
./picoclaw-agents agent -m "Usa x_post_tweet con message='Mirá esto', image_path='/tmp/photo.jpg'"

# Discord solo texto
./picoclaw-agents agent -m "Usa discord_post con message='Hola Discord'"

# Discord con imagen
./picoclaw-agents agent -m "Usa discord_post con message='Mirá esta imagen', image_path='/tmp/photo.jpg'"

# Publicar en múltiples plataformas
./picoclaw-agents agent -m "Publica '¡Gran anuncio!' en Facebook, Twitter y Discord"
```

### Uso en Telegram

Envía mensajes directamente a tu bot (con `picoclaw-agents gateway` activo):

```text
# Posts simples
Publica en Facebook: "¡Hola desde PicoClaw!"
Publica en Twitter: "Nuevo lanzamiento #PicoClaw"
Publica en Discord: "Anuncio importante"

# Con imágenes
Publica en Facebook la imagen /tmp/foto.jpg con mensaje "¡Nuevo producto!"
Publica en Twitter la imagen /tmp/photo.jpg con texto "Mirá esto"

# Multi-plataforma
Publica en Facebook y Twitter: "Gran noticia hoy"
Publica en todas las redes: "¡Anuncio importante!"
```

### Uso en Discord

Envía mensajes a tu bot de Discord o vía comandos:

```text
# Mensajes directos al bot
Publica en Facebook: "¡Hola desde nuestra comunidad!"
Publica en Twitter: "Nueva función lanzada #actualización"
Publica imagen /ruta/a/imagen.jpg en Facebook con mensaje "Mirá esto"

# Multi-plataforma
Publica en todas las redes: "¡Anuncio importante!"
```

### Combinado con Generación de Imágenes

```bash
# Generar imagen y publicar
./picoclaw-agents agent -m "Genera imagen de producto nuevo y publica en Facebook con texto atractivo"

# Flujo completo
./picoclaw-agents agent -m "Usa script_to_image_workflow con topic='Lanzamiento de producto', luego publica en Twitter"

# Integración con community manager
./picoclaw-agents agent -m "Genera imagen, crea borrador con community_manager para Discord, luego publica"
```

### Ejemplos de Community Manager

```bash
# Crear borrador desde contenido técnico
./picoclaw-agents agent -m "Usa community_manager_create_draft con raw_data='Nuevos endpoints de API lanzados', platform='discord'"

# Generar texto desde imagen
./picoclaw-agents agent -m "Usa community_from_image con image_path='./workspace/image_gen/abc/abc.-imagen.jpg', platform='twitter'"

# Flujo completo
./picoclaw-agents agent -m "Genera imagen, crea post atractivo con community_manager, publica en Facebook"
```

## Notas

- Si Facebook bloquea comentario con code `368`, PicoClaw fusiona el comentario dentro del cuerpo del post.
- Si Facebook responde code `190` y tienes campos de refresh configurados, PicoClaw refresca token y reintenta.
- Para más ejemplos prácticos, revisa `docs/SOCIAL_MEDIA_util.es.md` y `docs/IMAGE_GEN_util.es.md`.

---

## ⚡ Comandos Slash Fast-path (v3.4.1+)

Después de recibir una notificación de lote completado (ej: `#IMA_GEN_...` o `#SOCIAL_...`), usa comandos rápidos para gestión instantánea:

### Comandos de Gestión de Lotes

```text
/bundle_approve id=20260302_161740_yiia22
/bundle_regen id=20260302_161740_yiia22
/bundle_edit id=20260302_161740_yiia22
/bundle_publish id=20260302_161740_yiia22 platforms=facebook,twitter
```

**Beneficios:**
- ✅ **Latencia cero**: Sin razonamiento del LLM, ejecución instantánea
- ✅ **Sintaxis consistente**: Funciona idéntico en Telegram, Discord, CLI
- ✅ **Seguro**: Validación de ID antes de ejecutar

### Comandos de Utilidad

```text
/list pending          # Mostrar todas las tareas pendientes
/status                # Mostrar uso de tokens y estado del sistema
/help                  # Mostrar ayuda interactiva
/show model            # Mostrar modelo activo
/show channel          # Mostrar canal activo
```

### Global Tracker (v3.4.1+)

El **Global ImageGenTracker** es compartido entre todos los agentes:
- **Subagente genera contenido** → **Agente Principal puede publicar inmediatamente**
- **Sin errores de "ID no encontrado"** entre límites de agentes
- **Consistencia perfecta** en flujos de trabajo multi-agente

Ver [docs/queue_batch.es.md](docs/queue_batch.es.md) para documentación completa.

---

## 🤖 Flujos de Trabajo Multi-Agente

Con v3.4.1, puedes delegar flujos completos de redes sociales a subagentes:

```bash
# Delegar flujo completo a subagente
./picoclaw-agents agent -m "spawn task='Genera imagen sobre IA, crea post para Twitter y publícalo' label='campana_social'"

# El subagente hará:
# 1. Genera imagen con image_gen_create
# 2. Crea post con community_manager_create_draft
# 3. Publica en Twitter con x_post_tweet
# 4. Informa cuando esté completo
```

El Global Tracker asegura que el contenido generado por el subagente esté inmediatamente disponible para el agente principal para aprobación y publicación.
