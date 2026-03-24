# Generación de Imágenes - Util

Guía rápida para usar herramientas de generación de imágenes en PicoClaw desde terminal y Telegram.

> **PicoClaw v3.4.1**: Incluye **Comandos Slash Fast-path** para gestión instantánea de lotes y **Global Tracker** para consistencia multi-agente.
>
> **⚠️ IMPORTANTE: ESTO ES PARA GENERAR IMÁGENES ESTÁTICAS, NO VIDEOS**
>
> - `text_script_create`: Genera **TEXTO PARA POSTS** de Facebook/Twitter (como el copy de una publicación)
> - `image_gen_create`: Genera **IMÁGENES ESTÁTICAS** desde texto
> - **NO hay generación de video** en esta herramienta

## 📌 Almacenamiento de Imágenes (Arquitectura)

**NOTA IMPORTANTE**: Las imágenes generadas se guardan siempre en el **Global Shared Workspace**:

```
~/.picoclaw/workspace/image_gen/  ← Aquí se guardan TODAS las imágenes
```

**Esto es INTENCIONAL** para que:
- ✅ Las herramientas de publicación (`social_manager`, `social_post_bundle`) puedan acceder a las imágenes
- ✅ Múltiples workers puedan compartir el mismo repositorio de imágenes
- ✅ Las imágenes persistan aunque el worker es se recicle

**Para más detalles**, ver: [Arquitectura de Workspaces en SOCIAL_MEDIA.es.md](SOCIAL_MEDIA.es.md#arquitectura-de-workspaces-crítico)

---

## Requisitos

Configura tus credenciales en `~/.picoclaw/config.json`:

```json
{
  "tools": {
    "image_gen": {
      "provider": "gemini",
      "gemini_api_key": "TU_GEMINI_API_KEY",
      "gemini_image_model": "gemini-2.5-flash-image-preview",
      "ideogram_api_key": "",
      "ideogram_api_url": "https://api.ideogram.ai/v1/ideogram-v3/generate",
      "aspect_ratio": "4:5",
      "output_dir": "./workspace/image_gen"
    }
  }
}
```

O usa variables de entorno:

```bash
# Proveedor
export PICOCLAW_TOOLS_IMAGE_GEN_PROVIDER="gemini"

# Gemini (modelo con capacidad de imagen)
export PICOCLAW_TOOLS_IMAGE_GEN_GEMINI_API_KEY="tu_api_key"
export PICOCLAW_TOOLS_IMAGE_GEN_GEMINI_IMAGE_MODEL="gemini-2.5-flash-image-preview"

# Ideogram
export PICOCLAW_TOOLS_IMAGE_GEN_IDEOGRAM_API_KEY="tu_api_key"
export PICOCLAW_TOOLS_IMAGE_GEN_IDEOGRAM_API_URL="https://api.ideogram.ai/v1/ideogram-v3/generate"

# Aspect Ratio
export PICOCLAW_TOOLS_IMAGE_GEN_ASPECT_RATIO="4:5"

# Directorio de salida
export PICOCLAW_TOOLS_IMAGE_GEN_OUTPUT_DIR="./workspace/image_gen"
```

### Prioridad de Configuración

- Las variables de entorno tienen prioridad sobre `~/.picoclaw/config.json`.
- El nombre válido de variable es `PICOCLAW_TOOLS_IMAGE_GEN_GEMINI_IMAGE_MODEL`.
- `GEMINI_IMAGEN_MODEL` no lo usa PicoClaw.
- Las rutas relativas de `tools.image_gen.output_dir` se resuelven desde el workspace del agente (`agents.defaults.workspace`), no desde el `cwd` del proceso.
- Después de cambiar config/env, reinicia el proceso/servicio de PicoClaw.

---

## Herramientas Disponibles

### Reglas de Enrutamiento (Importante)

- Si el usuario pide solo imagen (ejemplo: `Genera una imagen de...`), usar solo `image_gen_create`.
- Usar `text_script_create` solo cuando el usuario pida explícitamente guion/texto de post, o pida flujo Script -> Imagen.
- `prompt_base_img.txt` se usa para construir un prompt visual desde un guion existente (no para solicitudes directas de solo imagen).

### `text_script_create`

Genera **TEXTO PARA POSTS DE REDES SOCIALES** (Facebook, Twitter, Discord).

**NO es para video** - es el texto que acompaña una publicación con imagen.

**Parámetros:**
- `topic` (requerido): Tema del post
- `category` (opcional): 'noticia', 'historia', 'tutorial', 'anuncio'
- `tone` (opcional): 'professional', 'casual', 'engaging'
- `language` (opcional): 'en', 'es' (auto-detectado)

**Ejemplo de uso:**
```bash
./picoclaw agent -m "Usa text_script_create con topic='Inteligencia Artificial', category='noticia'"
```

**Output:** Texto para post de Facebook (máx 1200 caracteres, estilo viral)

---

### `image_gen_create`

Genera **IMÁGENES ESTÁTICAS** desde un prompt de texto.

**Parámetros:**
- `prompt` (requerido): Descripción de la imagen
- `provider` (opcional): 'gemini' o 'ideogram'
- `aspect_ratio` (opcional): '4:5', '16:9', '1:1'

**Nota Técnica:**
- **Gemini:** Usa `gemini_image_model` desde config/env (default: `gemini-2.5-flash-image-preview`)
- Si Gemini devuelve `NOT_FOUND` para el modelo y tienes Ideogram configurado, PicoClaw cambia automáticamente a Ideogram
- En algunas cuentas/despliegues, ciertos modelos de Gemini pueden devolver salida cuadrada; para `4:5/16:9` estricto, usa Ideogram
- **Ideogram:** V3 API (1 imagen por defecto)

**Ejemplo de uso:**
```bash
./picoclaw agent -m "Usa image_gen_create con prompt='Atardecer cinematográfico en montañas'"
```

**Output:** Imagen JPG estática

---

### `image_gen_workflow`

Muestra opciones después de generar una imagen (publicar en redes, etc.).

**Parámetros:**
- `image_path` (requerido): Ruta de la imagen

---

### `script_to_image_workflow`

Muestra el flujo completo para crear **texto para post + generar imagen** basada en el tema.

**Parámetros:**
- `topic` (requerido): Tema para el guion e imagen
- `category` (opcional): 'historia', 'noticia', 'tutorial', 'anuncio'
- `create_script_first` (opcional): true/false (default: true)

**Ejemplo de uso:**
```bash
# Formato recomendado: espacios entre parámetros
./picoclaw agent -m "Usa script_to_image_workflow topic='Historia de dragones' category=historia"

# También funciona con 'con' (un solo parámetro)
./picoclaw agent -m "Usa script_to_image_workflow con topic='Historia de dragones'"
```

**Output:** 
- Muestra los pasos del workflow Script → Imagen
- **IMPORTANTE**: El agente mostrará las instrucciones pero necesitarás pedirle que ejecute cada paso:

```bash
# Paso 1: El agente te dirá que ejecute text_script_create
# Si no lo hace automáticamente, pídeselo:
./picoclaw agent -m "Ejecuta el paso 1: genera el script"

# Paso 2: Después de generar el script, pide la imagen:
./picoclaw agent -m "Ahora genera la imagen con image_gen_create"

# Paso 3: Opciones de publicación:
./picoclaw agent -m "Muestra opciones para la imagen generada"
```

**Nota:** Los archivos se guardan en la misma carpeta: `./workspace/image_gen/{ID}/`

**Workflow automático alternativo:**
Si quieres ejecutar todo en una sola petición, describe lo que quieres directamente:
```bash
./picoclaw agent -m "Genera un script sobre dragones y luego crea una imagen basada en ese tema"
```

---

### `community_manager_create_draft`

Crea borrador de post para redes sociales. **Soporta múltiples plataformas en una sola llamada**, guardando todas las variantes en la misma carpeta.

**Parámetros:**
- `raw_data` (requerido): Contenido técnico
- `platform` (requerido): 'discord', 'twitter', 'facebook', 'blog'. **Separar por comas para múltiples plataformas** (ej: 'discord,twitter,facebook')
- `include_emojis` (opcional): true/false (default: true)

**Ejemplo de uso (una plataforma):**
```bash
./picoclaw agent -m "Usa community_manager_create_draft raw_data='Contenido del post' platform=twitter"
```

**Ejemplo de uso (múltiples plataformas - RECOMENDADO):**
```bash
# IMPORTANTE: Usar comillas para la lista de plataformas
./picoclaw agent -m "Usa community_manager_create_draft raw_data='Contenido del post' platform='discord,twitter,facebook'"
```

**Output:** 
- **Una sola carpeta** con archivos para cada plataforma solicitada
- **Mismo contenido base** adaptado por plataforma (Twitter se acorta automáticamente a <280 caracteres)
- Estructura: `./text_scripts/{ID}/{ID}.-post_{platform}.txt`

### `community_from_image`

Genera texto para post basado en una imagen.

**Parámetros:**
- `image_path` (requerido): Ruta de la imagen
- `platform` (opcional): 'discord', 'twitter', 'facebook'

---

## Ejemplos de Uso

### Terminal

**Formato de parámetros:**
- ✅ **Recomendado:** `tool_name key1=value1 key2=value2` (espacios entre parámetros)
- ✅ **Alternativa:** `tool_name con key=value` (un solo parámetro con `con`)
- ✅ **Múltiples valores:** `key='valor1,valor2,valor3'` (usar comillas para listas con comas)
- ⚠️ **Evitar:** `key=valor1,valor2` sin comillas (las comas separan argumentos)

```bash
# Generar TEXTO PARA POST de Facebook
./picoclaw agent -m "Usa text_script_create topic='Inteligencia Artificial' category=noticia"

# Generar variantes para múltiples plataformas (RECOMENDADO)
# IMPORTANTE: Usar comillas para la lista de plataformas
./picoclaw agent -m "Usa community_manager_create_draft raw_data='La IA está revolucionando el mundo' platform='discord,twitter,facebook'"

# Generar para una sola plataforma
./picoclaw agent -m "Usa community_manager_create_draft raw_data='Contenido específico' platform=twitter"

# Generar IMAGEN ESTÁTICA
./picoclaw agent -m "Usa image_gen_create prompt='Atardecer cinematográfico en montañas'"

# Generar con aspect ratio
./picoclaw agent -m "Usa image_gen_create prompt='Foto de producto' aspect_ratio=16:9"

# Crear TEXTO POST + IMAGEN (workflow completo)
./picoclaw agent -m "Usa script_to_image_workflow topic='Historia de dragones' category=historia"

# Workflow post-generación
./picoclaw agent -m "Usa image_gen_workflow image_path='./workspace/image_gen/20260301_abc/20260301_abc.-imagen.jpg'"

# Generar texto para post desde imagen
./picoclaw agent -m "Usa community_from_image image_path='./workspace/image_gen/test.jpg' platform=facebook"
```

---

---

### Regla de Oro: Aprobación Social 🛡️

Para evitar publicaciones accidentales o contenido no deseado, PicoClaw tiene una **Regla Crítica**:

- **No Publicación Directa**: El agente nunca publicará en Facebook o Twitter sin antes mostrarte el borrador (texto e imagen) en el chat mediante la herramienta `message`.
- **Interacción Obligatoria**: Verás el contenido y botones para confirmar.
- **Excepción**: Solo si usas palabras como "publica directo" o "sin aprobación", el agente saltará este paso.

---

### Telegram

En Telegram, puedes usar lenguaje natural directamente o comandos explícitos.

**Ejemplos de lenguaje natural:**
- "Genera una imagen de un gato astronauta en Marte"
- "Crea un post para Facebook sobre coches eléctricos y su imagen"
- "Dibuja un paisaje cyber-punk en 16:9"

**Flujo interactivo:**
Cuando generas una imagen, PicoClaw te responderá con la imagen y **botones interactivos**:
- `📖 Ver guion`: Muestra el texto completo generado.
- `📱 Publicar`: Abre el menú de publicación en redes sociales.
- `🔄 Regenerar`: Intenta generar una versión diferente.

**Uso de subagentes (Recomendado):**
- "@picoclaw subagent task='Crea un guion de terror y su imagen'"
- "@picoclaw spawn task='Genera contenido sobre astronomía y publícalo'"

---

### Discord

Discord soporta las mismas funciones que Telegram, pero con la ventaja de una visualización más rica en canales específicos.

**Interacción en canales:**
- `!agent Genera una imagen de un bosque encantado`
- `!agent Crea un script para Twitter sobre el nuevo iPhone y genera su imagen`

**Mensajes con Embeds:**
PicoClaw enviará la imagen dentro de un **Embed**, lo que permite ver la metadata (tema, carpeta, ID) de forma organizada junto a los **botones de acción**.

**Comandos rápidos:**
- `subagent task='Genera imagen de un robot'`
- `image_gen_create prompt='Logo circular de café' aspect_ratio='1:1'`

## Estructura de Archivos

```
./workspace/image_gen/
├── tracker.json
├── 20260301_143022_abc123/
│   ├── 20260301_143022_abc123.-script.txt        # TEXTO PARA POST (Facebook/Twitter)
│   ├── 20260301_143022_abc123.-prompt_visual.txt # Prompt para generar imagen
│   └── 20260301_143022_abc123.-imagen.jpg        # IMAGEN ESTÁTICA generada
└── ...

./text_scripts/
├── tracker.json
├── 20260302_012338_6a2uue/
│   ├── 20260302_012338_6a2uue.-post_twitter.txt   # Versión Twitter (<280 chars)
│   ├── 20260302_012338_6a2uue.-post_facebook.txt  # Versión Facebook (completa)
│   └── 20260302_012338_6a2uue.-post_discord.txt   # Versión Discord (interactiva)
└── ...
```

**Nota:** Cuando se solicitan múltiples plataformas en una sola llamada, **todas las variantes se guardan en la misma carpeta**.

---

## Prompts Personalizados

Puedes usar tus propios prompts:

```json
{
  "tools": {
    "image_gen": {
      "image_script_path": "./workspace/prompt_base.txt",
      "image_gen_script_path": "./workspace/prompt_base_img.txt"
    }
  }
}
```

### Prompts por Defecto

**DEFAULT_IMAGE_SCRIPT** (para texto de posts):
- Estilo viral para Facebook
- Máximo 1200 caracteres
- Estructura: Gancho → Desarrollo → Cierre
- Optimizado para engagement (comentarios, shares)
- Archivo custom típico: `./workspace/prompt_base.txt`

**DEFAULT_IMAGE_GEN_SCRIPT** (para generar prompt visual desde un guion):
- Formato storyboard de Hollywood
- Output en inglés (mejor calidad de imagen)
- Incluye: personajes, atmósfera, iluminación
- Negative prompt: sin texto, sin marcas de agua
- Archivo custom típico: `./workspace/prompt_base_img.txt`

---

## Aspect Ratios

- `"4:5"` - Retrato (Instagram, Facebook)
- `"16:9"` - Paisaje (YouTube, Twitter)
- `"1:1"` - Cuadrado (Instagram feed)
- `"9:16"` - Vertical (Stories, Reels, TikTok)

---

---
321: 
322: ## Automatización Multi-Agente (Recomendado)
323: 
324: PicoClaw permite delegar flujos de trabajo completos a subagentes autónomos. Esto evita tener que ejecutar cada comando manualmente.
325: 
326: ### Cómo funciona
327: Usamos la herramienta `spawn` para darle una tarea compleja al agente. El agente "spawneará" un subagente que ejecutará todas las herramientas necesarias de forma secuencial y te informará cuando termine.
328: 
329: ### Ejemplos de delegación total:
330: 
331: **1. Crear contenido completo:**
332: ```bash
333: ./picoclaw agent -m "Usa spawn task='Escribe una historia de dragones, genera su imagen y dame opciones de publicación' label='dragones_full'"
334: ```
335: 
336: **2. Crear y publicar directamente:**
337: ```bash
338: ./picoclaw agent -m "Usa spawn task='Crea un post sobre IA, genera la imagen y publícalo en Discord' label='ia_post'"
339: ```
340: 
341: **Ventajas:**
342: - ✅ **Un solo comando:** No tienes que esperar a que termine cada paso.
343: - ✅ **Autonomía:** El subagente decide qué herramientas usar según tu descripción.
344: - ✅ **Multitarea:** Puedes seguir usando el agente principal mientras el subagente trabaja.
345: 
346: ---
347: 
348: ## Workflows Manuales

### Workflow 1: Post de Facebook con Imagen

```bash
# Paso 1: Generar TEXTO PARA POST
./picoclaw agent -m "Usa text_script_create con topic='Revolución de la IA'"

# Paso 2: Generar IMAGEN
./picoclaw agent -m "Usa image_gen_create con prompt='Robot futurista'"

# Paso 3: Publicar en Facebook
./picoclaw agent -m "Publica el texto y la imagen en Facebook"
```

### Workflow 2: Generación Directa

```bash
# Todo en uno
./picoclaw agent -m "Genera imagen de producto y publica en Twitter con texto"
```

### Workflow 3: Delegación Multi-Agente (v3.4.1+)

```bash
# Delegar flujo completo a subagente
./picoclaw agent -m "spawn task='Genera imagen sobre IA, crea post para Twitter y publícalo' label='campana_ia'"

# El subagente hará:
# 1. Genera imagen con image_gen_create
# 2. Crea post con community_manager_create_draft
# 3. Publica en Twitter con x_post_tweet
# 4. Informa cuando esté completo
```

El **Global Tracker** asegura que el contenido generado por el subagente esté inmediatamente disponible para el agente principal.

---

## ⚡ Comandos Slash Fast-path (v3.4.1+)

Después de recibir una notificación de generación de imagen (ej: `#IMA_GEN_...`), usa comandos rápidos:

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

Ver [docs/queue_batch.es.md](docs/queue_batch.es.md) para documentación completa.

---

## Límites de Tasa

| Proveedor | Límite     |
| --------- | ---------- |
| Gemini    | 60 req/min |
| Ideogram  | 20 req/min |

---

## Solución de Problemas

### "API key no configurada"
Configura `gemini_api_key` en config.json

### "Generación de imagen falló"
Verifica que el prompt no tenga contenido prohibido

### "Gemini model NOT_FOUND"
- Verifica el modelo en `~/.picoclaw/config.json` dentro de `tools.image_gen.gemini_image_model`
- Verifica que el entorno en runtime no lo esté sobrescribiendo:
  - `PICOCLAW_TOOLS_IMAGE_GEN_GEMINI_IMAGE_MODEL`
  - `PICOCLAW_TOOLS_IMAGE_GEN_PROVIDER`
- Reinicia PicoClaw después de cambios.

### "Múltiples imágenes generadas"
El sistema genera **1 sola imagen** por defecto:
- Gemini: `sampleCount: 1`
- Ideogram V3: 1 imagen por defecto

---

## Notas Técnicas

### Gemini API

**Modelo:** `gemini-2.5-flash-image-preview` (default, configurable)

- ✅ Usa la ruta de API de Gemini según el tipo de modelo
- ✅ Si no está disponible en tu cuenta, fallback automático a Ideogram (si está configurado)

### Ideogram API

**V3 API (Recomendado):**
- Endpoint: `https://api.ideogram.ai/v1/ideogram-v3/generate`
- 1 imagen por defecto

---

## Resumen

| Herramienta                | Output          | ¿Es video? |
| -------------------------- | --------------- | ---------- |
| `text_script_create`       | TEXTO PARA POST | ❌ NO       |
| `image_gen_create`         | IMAGEN ESTÁTICA | ❌ NO       |
| `script_to_image_workflow` | TEXTO + IMAGEN  | ❌ NO       |
| `community_from_image`     | TEXTO PARA POST | ❌ NO       |

**TODO ES PARA REDES SOCIALES (Facebook, Twitter, Discord) - NO HAY VIDEO**
