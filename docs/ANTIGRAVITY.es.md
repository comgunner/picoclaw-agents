# Guía del Proveedor Antigravity

**Última Actualización:** 29 de marzo de 2026  
**Estado:** ✅ Production Ready (v3.4.4)

---

## Descripción General

**Antigravity** (Google Cloud Code Assist) es un proveedor de IA respaldado por Google que ofrece acceso a modelos Gemini y Claude a través de la infraestructura de Google Cloud usando **autenticación OAuth 2.0**.

**Distinción Clave:** Antigravity usa las cuotas de tu plan **Google One AI Premium** o **Workspace Gemini** — NO una clave API de pago por uso.

---

## Autenticación

### Paso 1: Login

```bash
./picoclaw-agents auth login --provider google-antigravity
```

**Alias también funciona:**
```bash
./picoclaw-agents auth login --provider antigravity
```

### Paso 2: Completar Flujo OAuth

1. **El navegador se abre automáticamente** (máquinas locales)
2. **Inicia sesión** con tu cuenta de Google (debe tener Google One AI Premium o Workspace Gemini)
3. **Otorga permisos** a PicoClaw
4. **Credenciales guardadas** en `~/.picoclaw/auth.json`

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

Estos modelos están disponibles **solo vía Antigravity OAuth** — usan tu cuota de Google Cloud:

```bash
./picoclaw-agents auth models
```

### Modelos Validados (a marzo 2026)

| Nombre del Modelo | Descripción | Mejor Para |
|-------------------|-------------|------------|
| `antigravity-gemini-3-flash` | Rápido, confiable | **Recomendado por defecto** |
| `gemini-3-flash` | Igual al anterior (auto-resuelve) | Tareas rápidas |
| `gemini-3-pro-high` | Alto razonamiento Gemini 3 | Razonamiento complejo |
| `gemini-3.1-pro-high` | Alto razonamiento Gemini 3.1 | Tareas avanzadas |
| `gemini-3.1-flash-image` | Multimodal (solo **entrada** de imagen) | Análisis de imágenes |
| `gemini-2.5-pro` | Gemini 2.5 Pro | Propósito general |
| `gemini-2.5-flash` | Gemini 2.5 Flash | Respuestas rápidas |
| `gemini-2.5-flash-thinking` | Flash con razonamiento | Tareas de razonamiento |
| `gemini-2.5-flash-lite` | Modelo ligero | Tareas simples |
| `claude-sonnet-4-6` | Claude Sonnet | Escritura, análisis |
| `claude-opus-4-6-thinking` | Claude Opus con thinking | Resolución de problemas complejos |
| `gpt-oss-120b-medium` | Alternativa GPT open-source | Uso general |
| `chat_20706` | Modelo interno de Google | Testing |
| `chat_23310` | Modelo interno de Google | Testing |
| `tab_flash_lite_preview` | Modelo preview | Testing |
| `tab_jump_flash_lite_preview` | Modelo preview | Testing |

> [!IMPORTANT]
> **Generación de Imágenes NO Soportada vía OAuth**
> 
> Los modelos con sufijo `-image` (ej. `gemini-3.1-flash-image`) soportan solo **entrada/análisis de imágenes** — NO generación de imágenes.
> 
> Para **generación de imágenes**, debes usar **Clave API de Google AI Studio** (ver abajo).

### Ejemplos de Uso

```bash
# Usar modelo específico
./picoclaw-agents agent -m "Hola" --model antigravity-gemini-3-flash

# Claude con thinking
./picoclaw-agents agent -m "Resuelve este problema" --model claude-opus-4-6-thinking

# Análisis de imagen (NO generación)
./picoclaw-agents agent -m "Describe esta imagen" --model gemini-3.1-flash-image
```

---

## Generación de Imágenes (Solo Clave API)

**Antigravity OAuth NO soporta generación de imágenes.** Para generar imágenes, debes usar **Clave API de Google AI Studio**.

### Modelos de Imagen Soportados (Clave API)

| Modelo | Prefijo Provider | Propósito |
|--------|------------------|-----------|
| `gemini-2.5-flash-image` | `gemini/` | Nano Banana - generación de imágenes |
| `gemini-3-pro-image-preview` | `gemini/` | Nano Banana Pro |
| `gemini-3.1-flash-image-preview` | `gemini/` | Nano Banana 2 |
| `imagen-4.0-generate-001` | `gemini/` | Imagen 4 |
| `imagen-4.0-ultra-generate-001` | `gemini/` | Imagen 4 Ultra |

### Configuración

Agrega a `~/.picoclaw/config.json`:

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
      "gemini_api_key": "TU_GEMINI_API_KEY",  # pragma: allowlist secret
      "gemini_image_model_name": "gemini-2.5-flash-image",
      "output_dir": "~/.picoclaw/workspace/generated_images"
    }
  }
}
```

**Obtener Clave API:** [Google AI Studio](https://aistudio.google.com/app/apikey)

---

## Configuración

### Configuración Predeterminada (DeepSeek)

El archivo principal `config.example.json` usa **deepseek-chat** por defecto:

```bash
cp config/config.example.json ~/.picoclaw/config.json
# Agrega tu clave de API de DeepSeek
./picoclaw-agents agent -m "Hola"
```

### Configuración de Antigravity

El archivo `config/config.example_antigravity.json` usa `antigravity-gemini-3-flash` para todos los agentes:

```bash
cp config/config.example_antigravity.json ~/.picoclaw/config.json
./picoclaw-agents auth login --provider google-antigravity
./picoclaw-agents agent -m "Hola"
```

### Entradas de model_list

**Antigravity (OAuth):**
```json
{
  "model_name": "antigravity-gemini-3-flash",
  "model": "antigravity/gemini-3-flash",
  "api_key": "",
  "auth_method": "oauth"
}
```

**Google AI Studio (Clave API):**
```json
{
  "model_name": "gemini-2.5-flash",
  "model": "gemini/gemini-2.5-flash",
  "api_key": "TU_GEMINI_API_KEY"  # pragma: allowlist secret
}
```

### Ejemplos Comparativos

#### 1. Gemini 2.5 Flash

```json
// Antigravity OAuth (usa cuota de Google Cloud)
{
  "model_name": "ag-gemini-2.5-flash",
  "model": "antigravity/gemini-2.5-flash",
  "api_key": "",
  "auth_method": "oauth"
}

// Google AI Studio Clave API (pago por uso o capa gratuita)
{
  "model_name": "gemini-2.5-flash",
  "model": "gemini/gemini-2.5-flash",
  "api_key": "TU_GEMINI_API_KEY"  # pragma: allowlist secret
}
```

#### 2. Gemini 3 Flash

```json
// Antigravity OAuth
{
  "model_name": "ag-gemini-3-flash",
  "model": "antigravity/gemini-3-flash",
  "api_key": "",
  "auth_method": "oauth"
}

// Google AI Studio Clave API
{
  "model_name": "gemini-3-flash-preview",
  "model": "gemini/gemini-3-flash-preview",
  "api_key": "TU_GEMINI_API_KEY"  # pragma: allowlist secret
}
```

---

## Arquitectura de Enrutamiento de Modelos

PicoClaw usa un pipeline de 3 pasos:

### Campos de Configuración
- **`model_name`**: Alias interno — el nombre amigable que usas (ej. `antigravity-gemini-3-flash`)
- **`model`**: Instrucción de enrutamiento — debe contener `provider/model-id` (ej. `antigravity/gemini-3-flash`)

### El Pipeline

1. **Carga en Memoria**: Al iniciar, lee `model_list` de `~/.picoclaw/config.json` en RAM. Los cambios requieren reinicio.

2. **Enrutamiento (Factory)**: El alias se busca → el campo `model` se divide por `/` → el prefijo `antigravity` selecciona el proveedor Antigravity.

3. **Sanitización de Prefijos**: Antes de llamar a la API, el proveedor elimina todos los prefijos:
   - `antigravity/gemini-3-flash` → `gemini-3-flash` ✅
   - `antigravity-gemini-3-flash` → `gemini-3-flash` ✅ (prefijo con guión también se maneja)

> [!TIP]
> Tanto `antigravity/gemini-3-flash` (slash) como `antigravity-gemini-3-flash` (guión) son válidos.

---

## Uso en Producción (Coolify/Docker)

### Opción 1: Copiar Credenciales

```bash
# Autenticar localmente primero
./picoclaw-agents auth login --provider google-antigravity

# Copiar credenciales al servidor
scp ~/.picoclaw/auth.json usuario@tu-servidor:~/.picoclaw/
```

### Opción 2: Autenticar en el Servidor

```bash
# Ejecutar en el servidor (flujo headless)
./picoclaw-agents auth login --provider google-antigravity
# Copiar URL, abrir localmente, pegar URL de redireccionamiento de vuelta
```

---

## Solución de Problemas

| Error | Causa | Solución |
|-------|-------|----------|
| `403 PERMISSION_DENIED` | Token expirado/revocado | `./picoclaw-agents auth login --provider google-antigravity` |
| `ACCESS_TOKEN_SCOPE_INSUFFICIENT` | Token expirado/revocado | `./picoclaw-agents auth login --provider google-antigravity` |
| `404 NOT_FOUND` | Alias de modelo no resuelto | Verificar que `model` tenga prefijo `antigravity/` y `auth_method: "oauth"` |
| `401 invalid_api_key` | Proveedor incorrecto | Verificar que `model` tenga prefijo `antigravity/` |
| `429 Rate Limit` | Cuota agotada | Esperar tiempo de reset mostrado por PicoClaw, o cambiar de modelo |
| Respuesta vacía | Modelo restringido para el proyecto | Probar con `antigravity-gemini-3-flash` o `gemini-2.5-flash` |
| "Gemini for Google Cloud no está habilitado" | API no habilitada | Habilitar en [Google Cloud Console](https://console.cloud.google.com) |

---

## Requisitos

- **Cuenta de Google** con:
  - Plan Google One AI Premium, O
  - Complemento Workspace Gemini
- **Proyecto de Google Cloud** con API de Gemini habilitada
- **PicoClaw** v3.4.4 o posterior

---

## Referencia de Comandos

| Comando | Descripción |
|---------|-------------|
| `./picoclaw-agents auth login --provider google-antigravity` | Login con Antigravity |
| `./picoclaw-agents auth status` | Verificar estado de autenticación |
| `./picoclaw-agents auth models` | Listar modelos disponibles |
| `./picoclaw-agents auth logout --provider google-antigravity` | Cerrar sesión de Antigravity |
| `./picoclaw-agents agent -m "msg" --model <modelo>` | Usar modelo específico |

---

## Documentación Relacionada

- [ANTIGRAVITY.md](./ANTIGRAVITY.md) - Versión en inglés
- [IMAGE_GEN_util.md](./IMAGE_GEN_util.md) - Guías de generación de imágenes
- [Google Cloud Console](https://console.cloud.google.com) - Gestionar cuotas y facturación
- [Google AI Studio](https://aistudio.google.com) - Obtener claves API para generación de imágenes

---

**Inicio Rápido:** ¡Ejecuta `./picoclaw-agents auth login --provider google-antigravity` ahora! 🚀
