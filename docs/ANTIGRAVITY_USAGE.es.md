# Uso del Proveedor Antigravity en PicoClaw

Esta guía explica cómo configurar y utilizar el proveedor **Antigravity** (Google Cloud Code Assist) en PicoClaw.

## Todos los Proveedores Soportados

| Proveedor            | Comando                                                            | Método de Auth                                              | Modelo Default                                      |
| -------------------- | ------------------------------------------------------------------ | ----------------------------------------------------------- | --------------------------------------------------- |
| `google-antigravity` | `./picoclaw auth antigravity`<br>o `--provider google-antigravity` | OAuth 2.0 + PKCE (browser)                                  | `gemini-flash` (`antigravity/gemini-3-flash`)       |
| `openai`             | `./picoclaw auth login --provider openai`                          | OAuth 2.0 + PKCE (browser)<br>`--device-code` para headless | `gpt-5.2` (`openai/gpt-5.2`)                        |
| `anthropic`          | `./picoclaw auth login --provider anthropic`                       | Pegar API key (sin OAuth)                                   | `claude-sonnet-4.6` (`anthropic/claude-sonnet-4.6`) |

> [!NOTE]
> `anthropic` usa una **API key estática** de [console.anthropic.com](https://console.anthropic.com) — sin login con browser, sin expiración de token. `openai` y `google-antigravity` usan OAuth con auto-refresh.

## Requisitos Previos de Antigravity

1.  Una cuenta de Google.
2.  Google Cloud Code Assist habilitado. Generalmente disponible a través del proceso de **Gemini for Google Cloud**.
3.  **Compatibilidad de Planes**: Antigravity está diseñado específicamente para funcionar con los planes **Google One AI Premium** y complementos de Gemini para **Google Workspace**.

> [!NOTE]
> Esto es **diferente** de la API estándar de Google AI Studio. En lugar de una clave de API, utiliza autenticación OAuth vinculada a los créditos y cuotas de tu plan de Google.

## 1. Autenticación

Para autenticarse con Antigravity, ejecute:

```bash
./picoclaw auth antigravity
```

Esto abre una ventana del navegador para el flujo OAuth de Google. Al completar el login, las credenciales se guardan en `~/.picoclaw/auth.json`.

### Autenticación Manual (Headless/VPS)
Si está ejecutando en un servidor (Coolify/Docker) y no puede acceder a `localhost`:
1.  Ejecute el comando anterior.
2.  Copie la URL proporcionada y ábrala en su navegador local.
3.  Complete el inicio de sesión.
4.  Su navegador se redireccionará a `localhost:51121` (que no va a cargar).
5.  **Copie esa URL final** de la barra de direcciones.
6.  **Péguela de nuevo en la terminal** donde PicoClaw está esperando.

PicoClaw extraerá el código de autorización y completará el proceso automáticamente.

### Expiración del Token y Auto-Refresh
El `access_token` expira en **1 hora** (estándar OAuth de Google). PicoClaw maneja el refresh en 3 capas:
1. **Daemon en background**: renueva proactivamente cada 20 min si quedan <30 min
2. **En cada request**: intenta refresh con `refresh_token` incluso si el token ya expiró (no solo pre-expiración)
3. **Comando `auth models`**: también se recupera de tokens expirados automáticamente

Solo necesita re-autenticar (`./picoclaw auth antigravity`) si:
- Revocó el acceso desde `myaccount.google.com > Seguridad > Apps con acceso`
- Cambió su contraseña de Google
- El `refresh_token` lleva 6+ meses sin uso

## 2. Gestión de Modelos

### Listar Modelos Disponibles
```bash
./picoclaw auth models
```

### Modelos Validados y Funcionando (al 2026-03-04)

Use estos valores exactos de `model_name` en su config o con `--model`:

| model_name                   | Descripción                                       |
| ---------------------------- | ------------------------------------------------- |
| `antigravity-gemini-3-flash` | Rápido, confiable — **recomendado por defecto**   |
| `gemini-3-flash`             | Igual al anterior (se resuelve automáticamente)   |
| `gemini-3-pro-high`          | Alto razonamiento Gemini 3                        |
| `gemini-3.1-pro-high`        | Alto razonamiento Gemini 3.1                      |
| `gemini-3.1-flash-image`     | Rápido, soporte multimodal de imágenes            |
| `gemini-2.5-pro`             | Gemini 2.5 Pro                                    |
| `gemini-2.5-flash`           | Gemini 2.5 Flash                                  |
| `gemini-2.5-flash-thinking`  | Flash con razonamiento                            |
| `gemini-2.5-flash-lite`      | Modelo ligero                                     |
| `claude-sonnet-4-6`          | Claude rápido con razonamiento                    |
| `claude-opus-4-6-thinking`   | Claude de primer nivel con bloques de pensamiento |
| `gpt-oss-120b-medium`        | Alternativa open source                           |

### Modelos de Gemini Estándar (vía API Key)
Si no utilizas el proveedor de Antigravity (OAuth) y prefieres usar una clave de API comercial o gratuita de Google AI Studio, puedes configurar el proveedor añadiendo el prefijo explícito `gemini/` en la directiva `model` de tu `config.json` para enrutar directamente hacia la API de Google, y usar cualquiera de estos modelos públicos:

| Nombre del Modelo                         | Nombre Público (Display Name) | Límite Input | Límite Output |
| :---------------------------------------- | :---------------------------- | :----------- | :------------ |
| `gemini-2.5-flash`                        | Gemini 2.5 Flash              | 1,048,576    | 65,536        |
| `gemini-2.5-pro`                          | Gemini 2.5 Pro                | 1,048,576    | 65,536        |
| `gemini-2.0-flash`                        | Gemini 2.0 Flash              | 1,048,576    | 8,192         |
| `gemini-2.5-flash-lite`                   | Gemini 2.5 Flash-Lite         | 1,048,576    | 65,536        |
| `gemini-2.5-flash-image`                  | Nano Banana (Flash Image)     | 32,768       | 32,768        |
| `gemini-3-pro-preview`                    | Gemini 3 Pro Preview          | 1,048,576    | 65,536        |
| `gemini-3-flash-preview`                  | Gemini 3 Flash Preview        | 1,048,576    | 65,536        |
| `gemini-3.1-pro-preview`                  | Gemini 3.1 Pro Preview        | 1,048,576    | 65,536        |
| `gemini-3.1-flash-lite-preview`           | Gemini 3.1 Flash Lite Preview | 1,048,576    | 65,536        |
| `gemini-3-pro-image-preview`              | Nano Banana Pro               | 131,072      | 32,768        |
| `gemini-3.1-flash-image-preview`          | Nano Banana 2                 | 65,536       | 65,536        |
| `deep-research-pro-preview-12-2025`       | Deep Research Pro Preview     | 131,072      | 65,536        |
| `gemini-2.5-computer-use-preview-10-2025` | Gemini Computer Use Preview   | 131,072      | 65,536        |
| `gemini-robotics-er-1.5-preview`          | Gemini Robotics-ER 1.5        | 1,048,576    | 65,536        |
| `imagen-4.0-generate-001`                 | Imagen 4                      | 480          | 8,192         |
| `imagen-4.0-ultra-generate-001`           | Imagen 4 Ultra                | 480          | 8,192         |
| `veo-3.1-generate-preview`                | Veo 3.1                       | 480          | 8,192         |
| `gemini-flash-latest`                     | Gemini Flash Latest           | 1,048,576    | 65,536        |
| `gemini-pro-latest`                       | Gemini Pro Latest             | 1,048,576    | 65,536        |

### Cambiar de Modelo
```bash
# Usando modelos verificados por Antigravity Auth (OAuth)
./picoclaw agent -m "Hola" --model claude-opus-4-6-thinking
./picoclaw agent -m "Hola" --model antigravity-gemini-3-flash

# Usando modelos estándar vía API pública (Google AI Studio Key)
./picoclaw agent -m "Hola" --model gemini-2.5-flash
./picoclaw agent -m "Hola" --model gemini-3-pro-preview
```

## 3. Configuración

### Configuración Predeterminada (deepseek-chat)

El archivo principal `config.example.json` usa **deepseek-chat** como modelo predeterminado para todos los agentes. Esto proporciona un punto de partida consistente y económico para nuevos usuarios:

```bash
cp config/config.example.json ~/.picoclaw/config.json
# Agrega tu clave de API de DeepSeek en config.json
./picoclaw agent -m "Hola"
```

### Configuración de Antigravity

El archivo `config/config.example_antigravity.json` es un config listo para usar donde todos los agentes usan `antigravity-gemini-3-flash` por defecto. Usa esto si prefieres el proveedor Antigravity de Google:

```bash
cp config/config.example_antigravity.json ~/.picoclaw/config.json
./picoclaw auth antigravity
./picoclaw agent -m "Hola"
```

### Entradas de model_list para Antigravity

Cada modelo de Antigravity requiere `"auth_method": "oauth"` y no necesita `api_base`:

```json
{
  "model_name": "antigravity-gemini-3-flash",
  "model": "antigravity/gemini-3-flash",
  "api_key": "",
  "auth_method": "oauth"
},
{
  "model_name": "claude-sonnet-4-6",
  "model": "antigravity/claude-sonnet-4-6",
  "api_key": "",
  "auth_method": "oauth"
},
{
  "model_name": "claude-opus-4-6-thinking",
  "model": "antigravity/claude-opus-4-6-thinking",
  "api_key": "",
  "auth_method": "oauth"
}
```

### Ejemplos Comparativos: Antigravity (OAuth) vs Gemini (API Key)

Dependiendo de si utilizas la cuota de tu cuenta de Google Cloud (Antigravity vía OAuth) o tu propia clave generada en Google AI Studio, las definiciones de `model_list` varían específicamente en el prefijo y el campo `auth_method`. A continuación tienes 3 ejemplos claros de cómo configurar el mismo modelo usando distintos métodos:

#### 1. Gemini 2.5 Flash
```json
// Vía Antigravity (OAuth, sin API Key)
{
  "model_name": "ag-gemini-2.5-flash",
  "model": "antigravity/gemini-2.5-flash",
  "api_key": "",
  "auth_method": "oauth"
}

// Vía API Pública (Requiere API Key)
{
  "model_name": "gemini-2.5-flash",
  "model": "gemini/gemini-2.5-flash",
  "api_key": "TU_API_KEY_DE_GEMINI"
}
```

#### 2. Gemini 3 Flash
```json
// Vía Antigravity (OAuth, sin API Key)
{
  "model_name": "ag-gemini-3-flash",
  "model": "antigravity/gemini-3-flash",
  "api_key": "",
  "auth_method": "oauth"
}

// Vía API Pública (Requiere API Key)
{
  "model_name": "gemini-3-flash-preview",
  "model": "gemini/gemini-3-flash-preview",
  "api_key": "TU_API_KEY_DE_GEMINI"
}
```

#### 3. Gemini 2.5 Pro
```json
// Vía Antigravity (OAuth, sin API Key)
{
  "model_name": "ag-gemini-2.5-pro",
  "model": "antigravity/gemini-2.5-pro",
  "api_key": "",
  "auth_method": "oauth"
}

// Vía API Pública (Requiere API Key)
{
  "model_name": "gemini-2.5-pro",
  "model": "gemini/gemini-2.5-pro",
  "api_key": "TU_API_KEY_DE_GEMINI"
}
```

## 4. Uso en Producción (Coolify/Docker)

1.  Autentíquese localmente primero, luego copie las credenciales al servidor:
    ```bash
    scp ~/.picoclaw/auth.json usuario@su-servidor:~/.picoclaw/
    ```
2.  *Alternativamente*, ejecute `./picoclaw auth antigravity` directamente en el servidor con el flujo headless.

## 5. Solución de Problemas

| Error                                                     | Causa                                     | Solución                                                                                          |
| --------------------------------------------------------- | ----------------------------------------- | ------------------------------------------------------------------------------------------------- |
| `403 PERMISSION_DENIED / ACCESS_TOKEN_SCOPE_INSUFFICIENT` | Token expirado o revocado                 | Ejecutar `./picoclaw auth antigravity` de nuevo                                                   |
| `404 NOT_FOUND`                                           | Alias de modelo no resuelto correctamente | Verificar que la entrada en `model_list` tenga el campo `model` correcto y `auth_method: "oauth"` |
| `401 invalid_api_key`                                     | Proveedor incorrecto para el modelo       | Verificar que el campo `model` tenga prefijo `antigravity/`                                       |
| `429 Rate Limit`                                          | Cuota alcanzada                           | PicoClaw muestra el tiempo de restablecimiento; esperar o cambiar de modelo                       |
| Respuesta vacía                                           | Modelo restringido para el proyecto       | Probar con `antigravity-gemini-3-flash` o `gemini-2.5-flash`                                      |

## 6. Arquitectura de Enrutamiento de Modelos

PicoClaw usa un pipeline de 3 pasos para resolver la selección de modelos:

### Definiciones en la Configuración
- **`model_name`**: ALIAS interno — el nombre amigable que usa el agente o usuario (ej. `antigravity-gemini-3-flash`)
- **`model`**: Instrucción de enrutamiento — debe contener `proveedor/id-modelo` (ej. `antigravity/gemini-3-flash`)

### El Pipeline de 3 Pasos
1. **Carga en Memoria**: Al iniciar, PicoClaw lee el `model_list` de `~/.picoclaw/config.json` en RAM. Los cambios requieren reinicio completo.
2. **Enrutamiento (Factory Router)**: El alias se busca → el campo `model` se divide por `/` → el prefijo `antigravity` selecciona el proveedor correspondiente.
3. **Sanitización Final**: Antes de llamar a la API, el proveedor elimina todos los prefijos:
   - `antigravity/gemini-3-flash` → `gemini-3-flash` ✅
   - `antigravity-gemini-3-flash` → `gemini-3-flash` ✅ (prefijo con guión también se maneja)

> [!TIP]
> Tanto `antigravity/gemini-3-flash` (con slash) como `antigravity-gemini-3-flash` (con guión) son valores válidos en `model_list`. El proveedor elimina ambos correctamente.
