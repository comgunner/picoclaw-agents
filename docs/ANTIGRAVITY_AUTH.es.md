# Guía de Autenticación e Integración de Antigravity

## Descripción General

**Antigravity** (Google Cloud Code Assist) es un proveedor de modelos de IA respaldado por Google que ofrece acceso a modelos como Claude Opus 4.6 y Gemini a través de la infraestructura de Google Cloud. A diferencia de la API estándar de Google AI Studio, Antigravity aprovecha los créditos y cuotas incluidos en los planes **Google One AI Premium** y **Workspace Gemini**.

---

## Proveedores de Autenticación Soportados

PicoClaw soporta **3 proveedores** con autenticación integrada:

| Proveedor            | Comando                                                                            | Método de Auth                           | Modelo default que activa                           |
| -------------------- | ---------------------------------------------------------------------------------- | ---------------------------------------- | --------------------------------------------------- |
| `google-antigravity` | `./picoclaw auth antigravity`<br>o `--provider google-antigravity`                 | OAuth 2.0 + PKCE (browser)               | `gemini-flash` → `antigravity/gemini-3-flash`       |
| `openai`             | `./picoclaw auth login --provider openai`<br>Agregar `--device-code` para headless | OAuth 2.0 + PKCE (browser o device code) | `gpt-5.2` → `openai/gpt-5.2`                        |
| `anthropic`          | `./picoclaw auth login --provider anthropic`                                       | Pegar API key manualmente (sin OAuth)    | `claude-sonnet-4.6` → `anthropic/claude-sonnet-4.6` |

> [!NOTE]
> `anthropic` usa una **API key estática** (se obtiene en [console.anthropic.com](https://console.anthropic.com)). No tiene `refresh_token` ni expira automáticamente. Los otros dos proveedores usan OAuth con auto-refresh.

### Inicio Rápido por Proveedor

```bash
# Google Antigravity (recomendado — gratuito con Google One AI Premium)
./picoclaw auth antigravity

# OpenAI (requiere cuenta OpenAI con suscripción ChatGPT)
./picoclaw auth login --provider openai

# OpenAI headless (VPS/servidor sin navegador)
./picoclaw auth login --provider openai --device-code

# Anthropic (requiere API key de console.anthropic.com)
./picoclaw auth login --provider anthropic

# Ver estado de todos los proveedores autenticados
./picoclaw auth status

# Cerrar sesión de un proveedor específico
./picoclaw auth logout --provider google-antigravity

# Cerrar sesión de todos los proveedores
./picoclaw auth logout
```

---

## Índice

1. [Flujo de Autenticación](#flujo-de-autenticación)
2. [Detalles de Implementación de OAuth](#detalles-de-implementación-de-oauth)
3. [Gestión de Tokens](#gestión-de-tokens)
4. [Obtención de la Lista de Modelos](#obtención-de-la-lista-de-modelos)
5. [Seguimiento de Uso](#seguimiento-de-uso)
6. [Estructura del Plugin del Proveedor](#estructura-del-plugin-del-proveedor)
7. [Requisitos de Integración](#requisitos-de-integración)
8. [Endpoints de la API](#endpoints-de-la-api)
9. [Configuración](#configuración)
10. [Creación de un Nuevo Proveedor en PicoClaw](#creación-de-un-nuevo-proveedor-en-picoclaw)

---

## Flujo de Autenticación

### 1. OAuth 2.0 con PKCE

Antigravity utiliza **OAuth 2.0 con PKCE (Proof Key for Code Exchange)** para una autenticación segura:

```
┌─────────────┐                                    ┌─────────────────┐
│   Cliente   │ ───(1) Generar Par PKCE──────────> │                 │
│             │ ───(2) Abrir URL de Autenticación> │  Servidor OAuth │
│             │                                    │    de Google    │
│             │ <──(3) Redirigir con Código─────── │                 │
│             │                                    └─────────────────┘
│             │ ───(4) Intercambiar Código───────> │   URL de Token  │
│             │        por Tokens                  │                 │
│             │ <──(5) Tokens de Acceso + Refresco │                 │
└─────────────┘                                    └─────────────────┘
```

### 2. Pasos Detallados

#### Paso 1: Generar Parámetros PKCE
```typescript
function generatePkce(): { verifier: string; challenge: string } {
  const verifier = randomBytes(32).toString("hex");
  const challenge = createHash("sha256").update(verifier).digest("base64url");
  return { verifier, challenge };
}
```

#### Paso 2: Construir URL de Autorización
```typescript
const AUTH_URL = "https://accounts.google.com/o/oauth2/v2/auth";
const REDIRECT_URI = "http://localhost:51121/oauth-callback";

function buildAuthUrl(params: { challenge: string; state: string }): string {
  const url = new URL(AUTH_URL);
  url.searchParams.set("client_id", CLIENT_ID);
  url.searchParams.set("response_type", "code");
  url.searchParams.set("redirect_uri", REDIRECT_URI);
  url.searchParams.set("scope", SCOPES.join(" "));
  url.searchParams.set("code_challenge", params.challenge);
  url.searchParams.set("code_challenge_method", "S256");
  url.searchParams.set("state", params.state);
  url.searchParams.set("access_type", "offline");
  url.searchParams.set("prompt", "consent");
  return url.toString();
}
```

**Permisos (Scopes) Requeridos:**
```typescript
const SCOPES = [
  "https://www.googleapis.com/auth/cloud-platform",
  "https://www.googleapis.com/auth/userinfo.email",
  "https://www.googleapis.com/auth/userinfo.profile",
  "https://www.googleapis.com/auth/cclog",
  "https://www.googleapis.com/auth/experimentsandconfigs",
];
```

#### Paso 3: Manejar el Callback de OAuth

**Modo Automático (Desarrollo Local):**
- Inicia un servidor HTTP local en el puerto 51121
- Espera la redirección de Google
- Extrae el código de autorización de los parámetros de la consulta

**Modo Manual (Remoto/Headless):**
- Muestra la URL de autorización al usuario
- El usuario completa la autenticación en su navegador
- El usuario pega la URL de redirección completa en la terminal
- Se analiza el código de la URL pegada

#### Paso 4: Intercambiar Código por Tokens
```typescript
const TOKEN_URL = "https://oauth2.googleapis.com/token";

async function exchangeCode(params: {
  code: string;
  verifier: string;
}): Promise<{ access: string; refresh: string; expires: number }> {
  const response = await fetch(TOKEN_URL, {
    method: "POST",
    headers: { "Content-Type": "application/x-www-form-urlencoded" },
    body: new URLSearchParams({
      client_id: CLIENT_ID,
      client_secret: CLIENT_SECRET,
      code: params.code,
      grant_type: "authorization_code",
      redirect_uri: REDIRECT_URI,
      code_verifier: params.verifier,
    }),
  });

  const data = await response.json();
  
  return {
    access: data.access_token,
    refresh: data.refresh_token,
    expires: Date.now() + data.expires_in * 1000 - 5 * 60 * 1000, // Margen de 5 min
  };
}
```

#### Paso 5: Obtener Datos Adicionales del Usuario

**Correo del Usuario:**
```typescript
async function fetchUserEmail(accessToken: string): Promise<string | undefined> {
  const response = await fetch(
    "https://www.googleapis.com/oauth2/v1/userinfo?alt=json",
    { headers: { Authorization: `Bearer ${accessToken}` } }
  );
  const data = await response.json();
  return data.email;
}
```

**ID del Proyecto (Requerido para llamadas a la API):**
```typescript
async function fetchProjectId(accessToken: string): Promise<string> {
  const headers = {
    Authorization: `Bearer ${accessToken}`,
    "Content-Type": "application/json",
    "User-Agent": "google-api-nodejs-client/9.15.1",
    "X-Goog-Api-Client": "google-cloud-sdk vscode_cloudshelleditor/0.1",
    "Client-Metadata": JSON.stringify({
      ideType: "IDE_UNSPECIFIED",
      platform: "PLATFORM_UNSPECIFIED",
      pluginType: "GEMINI",
    }),
  };

  const response = await fetch(
    "https://cloudcode-pa.googleapis.com/v1internal:loadCodeAssist",
    {
      method: "POST",
      headers,
      body: JSON.stringify({
        metadata: {
          ideType: "IDE_UNSPECIFIED",
          platform: "PLATFORM_UNSPECIFIED",
          pluginType: "GEMINI",
        },
      }),
    }
  );

  const data = await response.json();
  return data.cloudaicompanionProject || "rising-fact-p41fc"; // Valor predeterminado de respaldo
}
```

---

## Detalles de Implementación de OAuth

### Credenciales del Cliente

**Importante:** Estas están codificadas en base64 en el código fuente para la sincronización con pi-ai:

```typescript
const decode = (s: string) => Buffer.from(s, "base64").toString();

const CLIENT_ID = decode(
  "MTA3MTAwNjA2MDU5MS10bWhzc2luMmgyMWxjcmUyMzV2dG9sb2poNGc0MDNlcC5hcHBzLmdvb2dsZXVzZXJjb250ZW50LmNvbQ=="
);
const CLIENT_SECRET = decode("R09DU1BYLUs1OEZXUjQ4NkxkTEoxbUxCOHNYQzR6NnFEQWY=");
```

### Modos del Flujo OAuth

1. **Flujo Automático** (Máquinas locales con navegador):
   - Abre el navegador automáticamente
   - El servidor de callback local captura la redirección
   - No se requiere interacción del usuario después de la autenticación inicial

2. **Flujo Manual** (Remoto/headless/WSL2):
   - Se muestra la URL para copiar y pegar manualmente
   - El usuario completa la autenticación en un navegador externo
   - El usuario pega la URL de redirección completa de vuelta

```typescript
function shouldUseManualOAuthFlow(isRemote: boolean): boolean {
  return isRemote || isWSL2Sync();
}
```

---

## Gestión de Tokens

### Estructura del Perfil de Autenticación

```typescript
type OAuthCredential = {
  type: "oauth";
  provider: "google-antigravity";
  access: string;           // Token de acceso
  refresh: string;          // Token de refresco
  expires: number;          // Marca de tiempo de expiración (ms desde la época)
  email?: string;           // Correo del usuario
  projectId?: string;       // ID del proyecto de Google Cloud
};
```

### Refresco de Tokens

La credencial incluye un `refresh_token` que se usa para obtener nuevos `access_token` cuando el actual expire. El sistema maneja el refresco en **3 capas**:

1. **Daemon proactivo** (`pkg/auth/refresh_daemon.go`): Un goroutine en background revisa cada 20 minutos si el token expira en menos de 30 minutos, y lo renueva anticipadamente.
2. **Refresco reactivo** (`antigravity_provider.go`): Cada llamada al agente verifica si el token está por vencer (`NeedsRefresh`) **o ya expiró** (`IsExpired`). Si hay `refresh_token`, intenta renovar antes de fallar.
3. **Scope completo** (`oauth.go`): El refresh usa los mismos scopes que el login inicial (`cloud-platform`, `cclog`, etc.), evitando errores 403 post-refresh.

### Expiración del Token

- El `access_token` de Google dura **~1 hora** (3599 segundos). Esto es fijo en el servidor de Google y no es configurable.
- El `refresh_token` dura meses/indefinidamente, a menos que sea revocado.
- La re-autenticación manual solo es necesaria si el `refresh_token` fue revocado (cambio de contraseña, revocación manual, inactividad de 6+ meses).

---

## Obtención de la Lista de Modelos

### Obtener Modelos Disponibles

```typescript
const BASE_URL = "https://cloudcode-pa.googleapis.com";

async function fetchAvailableModels(
  accessToken: string,
  projectId: string
): Promise<Model[]> {
  const headers = {
    Authorization: `Bearer ${accessToken}`,
    "Content-Type": "application/json",
    "User-Agent": "antigravity",
    "X-Goog-Api-Client": "google-cloud-sdk vscode_cloudshelleditor/0.1",
  };

  const response = await fetch(
    `${BASE_URL}/v1internal:fetchAvailableModels`,
    {
      method: "POST",
      headers,
      body: JSON.stringify({ project: projectId }),
    }
  );

  const data = await response.json();
  
  // Devuelve modelos con información de cuota
  return Object.entries(data.models).map(([modelId, modelInfo]) => ({
    id: modelId,
    displayName: modelInfo.displayName,
    quotaInfo: {
      remainingFraction: modelInfo.quotaInfo?.remainingFraction,
      resetTime: modelInfo.quotaInfo?.resetTime,
      isExhausted: modelInfo.quotaInfo?.isExhausted,
    },
  }));
}
```

### Formato de Respuesta

```typescript
type FetchAvailableModelsResponse = {
  models?: Record<string, {
    displayName?: string;
    quotaInfo?: {
      remainingFraction?: number | string;
      resetTime?: string;      // Marca de tiempo ISO 8601
      isExhausted?: boolean;
    };
  }>;
};
```

---

## Seguimiento de Uso

### Obtener Datos de Uso

```typescript
export async function fetchAntigravityUsage(
  token: string,
  timeoutMs: number
): Promise<ProviderUsageSnapshot> {
  // 1. Obtener créditos e información del plan
  const loadCodeAssistRes = await fetch(
    `${BASE_URL}/v1internal:loadCodeAssist`,
    {
      method: "POST",
      headers: {
        Authorization: `Bearer ${token}`,
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        metadata: {
          ideType: "ANTIGRAVITY",
          platform: "PLATFORM_UNSPECIFIED",
          pluginType: "GEMINI",
        },
      }),
    }
  );

  // Extraer información de créditos
  const { availablePromptCredits, planInfo, currentTier } = data;
  
  // 2. Obtener cuotas de modelos
  const modelsRes = await fetch(
    `${BASE_URL}/v1internal:fetchAvailableModels`,
    {
      method: "POST",
      headers: { Authorization: `Bearer ${token}` },
      body: JSON.stringify({ project: projectId }),
    }
  );

  // Construir ventanas de uso
  return {
    provider: "google-antigravity",
    displayName: "Google Antigravity",
    windows: [
      { label: "Credits", usedPercent: calculateUsedPercent(available, monthly) },
      // Cuotas individuales de modelos...
    ],
    plan: currentTier?.name || planType,
  };
}
```

### Estructura de la Respuesta de Uso

```typescript
type ProviderUsageSnapshot = {
  provider: "google-antigravity";
  displayName: string;
  windows: UsageWindow[];
  plan?: string;
  error?: string;
};

type UsageWindow = {
  label: string;           // "Credits" o ID del modelo
  usedPercent: number;     // 0-100
  resetAt?: number;        // Marca de tiempo cuando la cuota se restablece
};
```

---

## Estructura del Plugin del Proveedor

### Definición del Plugin

```typescript
const antigravityPlugin = {
  id: "google-antigravity-auth",
  name: "Google Antigravity Auth",
  description: "OAuth flow para Google Antigravity (Cloud Code Assist)",
  configSchema: emptyPluginConfigSchema(),
  
  register(api: PicoClawPluginApi) {
    api.registerProvider({
      id: "google-antigravity",
      label: "Google Antigravity",
      docsPath: "/providers/models",
      aliases: ["antigravity"],
      
      auth: [
        {
          id: "oauth",
          label: "Google OAuth",
          hint: "PKCE + callback en localhost",
          kind: "oauth",
          run: async (ctx: ProviderAuthContext) => {
            // Implementación de OAuth aquí
          },
        },
      ],
    });
  },
};
```

---

## Requisitos de Integración

### 1. Entorno/Dependencias Requeridas

- Go ≥ 1.21
- Código base de PicoClaw (`pkg/providers/` y `pkg/auth/`)
- Paquetes de la biblioteca estándar `crypto` y `net/http`

### 2. Encabezados Requeridos para Llamadas API

```typescript
const REQUIRED_HEADERS = {
  "Authorization": `Bearer ${accessToken}`,
  "Content-Type": "application/json",
  "User-Agent": "antigravity",  // o "google-api-nodejs-client/9.15.1"
  "X-Goog-Api-Client": "google-cloud-sdk vscode_cloudshelleditor/0.1",
};

// Para llamadas loadCodeAssist, también incluya:
const CLIENT_METADATA = {
  ideType: "ANTIGRAVITY",  // o "IDE_UNSPECIFIED"
  platform: "PLATFORM_UNSPECIFIED",
  pluginType: "GEMINI",
};
```

### 3. Sanitización del Esquema del Modelo

Antigravity utiliza modelos compatibles con Gemini, por lo que los esquemas de herramientas deben ser sanitizados:

```typescript
const GOOGLE_SCHEMA_UNSUPPORTED_KEYWORDS = new Set([
  "patternProperties",
  "additionalProperties",
  "$schema",
  "$id",
  "$ref",
  "$defs",
  "definitions",
  "examples",
  "minLength",
  "maxLength",
  "minimum",
  "maximum",
  "multipleOf",
  "pattern",
  "format",
  "minItems",
  "maxItems",
  "uniqueItems",
  "minProperties",
  "maxProperties",
]);

// Limpiar esquema antes de enviar
function cleanToolSchemaForGemini(schema: Record<string, unknown>): unknown {
  // Eliminar palabras clave no soportadas
  // Asegurar que el nivel superior tenga tipo: "object"
  // Aplanar uniones anyOf/oneOf
}
```

---

## Endpoints de la API

### Endpoints de Autenticación

| Endpoint                                        | Método | Propósito                       |
| ----------------------------------------------- | ------ | ------------------------------- |
| `https://accounts.google.com/o/oauth2/v2/auth`  | GET    | Autorización OAuth              |
| `https://oauth2.googleapis.com/token`           | POST   | Intercambio de tokens           |
| `https://www.googleapis.com/oauth2/v1/userinfo` | GET    | Información del usuario (email) |

---

## Configuración

### Configuración en `config.json`

```json
{
  "model_list": [
    {
      "model_name": "gemini-flash",
      "model": "antigravity/gemini-3-flash",
      "auth_method": "oauth"
    }
  ],
  "agents": {
    "defaults": {
      "model": "gemini-flash"
    }
  }
}
```

### Almacenamiento del Perfil de Autenticación

Los perfiles de autenticación se almacenan en `~/.picoclaw/auth.json`:

```json
{
  "credentials": {
    "google-antigravity": {
      "access_token": "ya29...",
      "refresh_token": "1//...",
      "expires_at": "2026-01-01T00:00:00Z",
      "provider": "google-antigravity",
      "auth_method": "oauth",
      "email": "user@example.com",
      "project_id": "my-project-id"
    }
  }
}
```

---

## Creación de un Nuevo Proveedor en PicoClaw

Los proveedores de PicoClaw se implementan como paquetes de Go bajo `pkg/providers/`. Para agregar un nuevo proveedor:

### Paso a Paso de la Implementación

#### 1. Crear Archivo del Proveedor

Cree un nuevo archivo Go en `pkg/providers/`:

```
pkg/providers/
└── su_proveedor.go
```

#### 2. Implementar la Interfaz Provider

Su proveedor debe implementar la interfaz `Provider` definida en `pkg/providers/types.go`:

```go
package providers

type SuProveedor struct {
    apiKey  string
    apiBase string
}

func NewSuProveedor(apiKey, apiBase, proxy string) *SuProveedor {
    if apiBase == "" {
        apiBase = "https://api.punto-de-entrada.com/v1"
    }
    return &SuProveedor{apiKey: apiKey, apiBase: apiBase}
}

func (p *SuProveedor) Chat(ctx context.Context, messages []Message, tools []Tool, cb StreamCallback) error {
    // Implementar chat completion con streaming
}
```

#### 3. Registrar en la Fábrica

Agregue su proveedor al interruptor (switch) de protocolo en `pkg/providers/factory.go`:

```go
case "su-proveedor":
    return NewSuProveedor(sel.apiKey, sel.apiBase, sel.proxy), nil
```

#### 4. Agregar Configuración Predeterminada (Opcional)

Agregue una entrada predeterminada en `pkg/config/defaults.go`:

```go
{
    ModelName: "su-modelo",
    Model:     "su-proveedor/nombre-del-modelo",
    APIKey:    "",
},
```

#### 5. Agregar Soporte de Autenticación (Opcional)

Si su proveedor requiere OAuth o una autenticación especial, agregue un caso en `cmd/picoclaw/cmd_auth.go`:

```go
case "su-proveedor":
    authLoginSuProveedor()
```

#### 6. Configurar vía `config.json`

```json
{
  "model_list": [
    {
      "model_name": "su-modelo",
      "model": "su-proveedor/nombre-del-modelo",
      "api_key": "su-api-key",
      "api_base": "https://api.punto-de-entrada.com/v1"
    }
  ]
}
```

---

## Referencias

- **Archivos Fuente:**
  - `pkg/providers/antigravity_provider.go` - Implementación del proveedor Antigravity
  - `pkg/auth/oauth.go` - Implementación del flujo OAuth
  - `pkg/auth/store.go` - Almacenamiento de credenciales de autenticación (`~/.picoclaw/auth.json`)
  - `pkg/providers/factory.go` - Fábrica de proveedores y enrutamiento de protocolos
  - `pkg/providers/types.go` - Definiciones de la interfaz del proveedor
  - `cmd/picoclaw/cmd_auth.go` - Comandos de la CLI para autenticación

- **Documentación:**
  - `docs/ANTIGRAVITY_USAGE.md` - Guía de uso de Antigravity
  - `docs/migration/model-list-migration.md` - Guía de migración

---

## Notas

1. **Proyecto de Google Cloud:** Antigravity requiere que Gemini for Google Cloud esté habilitado en su proyecto de Google Cloud.
2. **Cuotas:** Utiliza las cuotas del proyecto de Google Cloud (no facturación separada).
3. **Acceso a Modelos:** Los modelos disponibles dependen de la configuración de su proyecto de Google Cloud.

---

## Solución de Problemas Comunes

### 1. Límite de Tasa (HTTP 429)

Antigravity devuelve un error 429 cuando se agotan las cuotas del proyecto/modelo. La respuesta de error a menudo contiene un `quotaResetDelay` en el campo `details`.

### 2. Respuestas Vacías (Modelos Restringidos)

Algunos modelos pueden aparecer en la lista de modelos disponibles pero devolver una respuesta vacía. Esto suele suceder con modelos en vista previa o restringidos para los que el proyecto actual no tiene permiso.

---

## "Token expirado" (403 PERMISSION_DENIED)
- El sistema ahora intenta auto-refresh incluso para tokens ya expirados.
- Si el auto-refresh falla, re-autentíquese: `./picoclaw auth antigravity`

## "Gemini for Google Cloud no está habilitado"
- Habilite la API en su Google Cloud Console

---
