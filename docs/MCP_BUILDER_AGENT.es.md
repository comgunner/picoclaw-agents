# MCP Builder Agent - Guía Completa

**Versión:** 1.0.0  
**Categoría:** Specialized  
**Skill ID:** `specialized-mcp-builder`

---

## 📋 Índice

1. [¿Qué es MCP Builder Agent?](#qué-es-mcp-builder-agent)
2. [Casos de Uso](#casos-de-uso)
3. [Cómo Activar](#cómo-activar)
4. [Ejemplos Prácticos](#ejemplos-prácticos)
5. [Estructura de un Servidor MCP](#estructura-de-un-servidor-mcp)
6. [Mejores Prácticas](#mejores-prácticas)
7. [Referencia de API](#referencia-de-api)

---

## ¿Qué es MCP Builder Agent?

**MCP Builder Agent** es un skill especializado en construir servidores del **Model Context Protocol (MCP)**. Los servidores MCP extienden las capacidades de los agentes de IA al exponer herramientas, recursos y prompts personalizados.

### Características Principales

- 🛠️ **Diseño de Herramientas**: Nombres claros, parámetros tipados, descripciones útiles
- 📚 **Exposición de Recursos**: Expone fuentes de datos que los agentes pueden leer
- 🔄 **Manejo de Errores**: Fallos gracefully con mensajes de error accionables
- 🔐 **Seguridad**: Validación de inputs, autenticación, rate limiting
- ✅ **Testing**: Unit tests para herramientas, integration tests para el servidor

---

## Casos de Uso

### 1. Integración con APIs Externas

Crea herramientas que permiten a tu agente interactuar con APIs de terceros:

```typescript
// Ejemplo: Servidor MCP para API de clima
server.tool("get_weather", { 
  city: z.string(), 
  units: z.enum(["celsius", "fahrenheit"]).default("celsius") 
}, async ({ city, units }) => {
  const weather = await fetchWeatherAPI(city, units);
  return { content: [{ type: "text", text: JSON.stringify(weather, null, 2) }] };
});
```

### 2. Acceso a Bases de Datos

Expone datos de tu base de datos de forma segura:

```typescript
server.tool("search_users", { 
  query: z.string(), 
  limit: z.number().default(10) 
}, async ({ query, limit }) => {
  const users = await db.user.findMany({
    where: { name: { contains: query } },
    take: limit
  });
  return { content: [{ type: "text", text: JSON.stringify(users, null, 2) }] };
});
```

### 3. Automatización de Workflows

Automatiza tareas repetitivas de tu negocio:

```typescript
server.tool("create_invoice", { 
  client_id: z.string(), 
  amount: z.number(), 
  description: z.string() 
}, async ({ client_id, amount, description }) => {
  const invoice = await createInvoiceInERP(client_id, amount, description);
  return { content: [{ type: "text", text: `Invoice created: ${invoice.id}` }] };
});
```

### 4. Integración con Sistemas de Archivos

Permite lectura/escritura controlada de archivos:

```typescript
server.tool("read_config", { 
  path: z.string().refine(p => p.startsWith('/safe/')) 
}, async ({ path }) => {
  const content = await fs.readFile(path, 'utf-8');
  return { content: [{ type: "text", text: content }] };
});
```

---

## Cómo Activar

### Método 1: Vía CLI

```bash
# Invocar directamente al MCP Builder Agent
picoclaw-agents agent -m "Necesito crear un servidor MCP para integrar con la API de Stripe"
```

### Método 2: Vía Chat (Telegram/Discord)

Envía un mensaje a tu agente:

```
@agent Necesito un servidor MCP que permita:
1. Buscar productos en mi base de datos
2. Crear órdenes de compra
3. Consultar el estado de envíos

¿Puedes ayudarme a construirlo?
```

### Método 3: Configurar como Skill Predeterminada

En `~/.picoclaw/config.json`:

```json
{
  "agents": {
    "list": [
      {
        "id": "integration_specialist",
        "name": "Integration Specialist",
        "skills": ["specialized-mcp-builder"]
      }
    ]
  }
}
```

---

## Ejemplos Prácticos

### Ejemplo 1: Servidor MCP para API de GitHub

**Objetivo**: Crear herramientas para interactuar con GitHub API.

**Prompt para MCP Builder:**

```
Necesito un servidor MCP que permita:
1. Buscar repositorios por nombre
2. Obtener los últimos commits de un repo
3. Crear issues en un repositorio
4. Listar pull requests abiertos

Requisitos:
- Autenticación vía token de GitHub
- Rate limiting para evitar bloqueos
- Manejo de errores 404 y 403
```

**Resultado esperado:**

```typescript
// github-server.ts
import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";
import { z } from "zod";
import { Octokit } from "octokit";

const server = new McpServer({ 
  name: "github-server", 
  version: "1.0.0" 
});

// Inicializar Octokit con token
const octokit = new Octokit({ 
  auth: process.env.GITHUB_TOKEN 
});

// Herramienta: Buscar repositorios
server.tool(
  "search_repositories",
  { 
    query: z.string().describe("Search query (e.g., 'machine learning language:python')"),
    per_page: z.number().default(10).describe("Number of results (max 100)")
  },
  async ({ query, per_page }) => {
    try {
      const { data } = await octokit.request('GET /search/repositories', {
        q: query,
        per_page: Math.min(per_page, 100)
      });
      
      return { 
        content: [{ 
          type: "text", 
          text: JSON.stringify({
            total_count: data.total_count,
            items: data.items.map(repo => ({
              name: repo.full_name,
              description: repo.description,
              stars: repo.stargazers_count,
              url: repo.html_url
            }))
          }, null, 2) 
        }] 
      };
    } catch (error: any) {
      return { 
        content: [{ 
          type: "text", 
          text: `Error searching repositories: ${error.message}` 
        }],
        isError: true
      };
    }
  }
);

// Herramienta: Obtener commits
server.tool(
  "get_commits",
  { 
    owner: z.string().describe("Repository owner (e.g., 'facebook')"),
    repo: z.string().describe("Repository name (e.g., 'react')"),
    per_page: z.number().default(5)
  },
  async ({ owner, repo, per_page }) => {
    const { data } = await octokit.repos.listCommits({
      owner,
      repo,
      per_page
    });
    
    return { 
      content: [{ 
        type: "text", 
        text: JSON.stringify(data.map(commit => ({
          sha: commit.sha,
          message: commit.commit.message,
          author: commit.commit.author.name,
          date: commit.commit.author.date
        })), null, 2) 
      }] 
    };
  }
);

// Herramienta: Crear issue
server.tool(
  "create_issue",
  { 
    owner: z.string(),
    repo: z.string(),
    title: z.string(),
    body: z.string().optional(),
    labels: z.array(z.string()).optional()
  },
  async ({ owner, repo, title, body, labels }) => {
    const { data } = await octokit.issues.create({
      owner,
      repo,
      title,
      body,
      labels
    });
    
    return { 
      content: [{ 
        type: "text", 
        text: `Issue created: ${data.html_url}` 
      }] 
    };
  }
);

// Iniciar servidor
const transport = new StdioServerTransport();
await server.connect(transport);
```

**Configuración en PicoClaw:**

```json
{
  "tools": {
    "mcp": {
      "github": {
        "command": "node",
        "args": ["/path/to/github-server.ts"],
        "env": {
          "GITHUB_TOKEN": "ghp_..."
        }
      }
    }
  }
}
```

---

### Ejemplo 2: Servidor MCP para Base de Datos PostgreSQL

**Objetivo**: Exponer datos de usuarios y productos.

**Prompt para MCP Builder:**

```
Crea un servidor MCP para conectar a PostgreSQL con:
1. Buscar usuarios por email o nombre
2. Listar productos por categoría
3. Obtener órdenes de un usuario
4. Crear nueva orden de compra

Usa validación Zod para todos los inputs.
```

**Código generado:**

```typescript
// database-server.ts
import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";
import { z } from "zod";
import { Pool } from "pg";

const server = new McpServer({ 
  name: "database-server", 
  version: "1.0.0" 
});

// Configurar pool de conexiones
const pool = new Pool({
  connectionString: process.env.DATABASE_URL,
  max: 20,
  idleTimeoutMillis: 30000,
  connectionTimeoutMillis: 2000,
});

// Herramienta: Buscar usuarios
server.tool(
  "search_users",
  { 
    query: z.string().describe("Email or name to search"),
    limit: z.number().default(10).max(100)
  },
  async ({ query, limit }) => {
    const client = await pool.connect();
    try {
      const result = await client.query(
        `SELECT id, name, email, created_at 
         FROM users 
         WHERE name ILIKE $1 OR email ILIKE $1 
         LIMIT $2`,
        [`%${query}%`, limit]
      );
      
      return { 
        content: [{ 
          type: "text", 
          text: JSON.stringify(result.rows, null, 2) 
        }] 
      };
    } catch (error: any) {
      return { 
        content: [{ 
          type: "text", 
          text: `Database error: ${error.message}` 
        }],
        isError: true
      };
    } finally {
      client.release();
    }
  }
);

// Herramienta: Listar productos
server.tool(
  "list_products",
  { 
    category: z.string().optional().describe("Filter by category"),
    min_price: z.number().optional(),
    max_price: z.number().optional()
  },
  async ({ category, min_price, max_price }) => {
    const client = await pool.connect();
    try {
      let query = "SELECT * FROM products WHERE 1=1";
      const params: any[] = [];
      
      if (category) {
        params.push(category);
        query += ` AND category = $${params.length}`;
      }
      
      if (min_price !== undefined) {
        params.push(min_price);
        query += ` AND price >= $${params.length}`;
      }
      
      if (max_price !== undefined) {
        params.push(max_price);
        query += ` AND price <= $${params.length}`;
      }
      
      const result = await client.query(query, params);
      
      return { 
        content: [{ 
          type: "text", 
          text: JSON.stringify(result.rows, null, 2) 
        }] 
      };
    } catch (error: any) {
      return { 
        content: [{ 
          type: "text", 
          text: `Database error: ${error.message}` 
        }],
        isError: true
      };
    } finally {
      client.release();
    }
  }
);

// Iniciar servidor
const transport = new StdioServerTransport();
await server.connect(transport);
```

**Package.json:**

```json
{
  "name": "database-mcp-server",
  "version": "1.0.0",
  "type": "module",
  "dependencies": {
    "@modelcontextprotocol/sdk": "^0.5.0",
    "zod": "^3.22.4",
    "pg": "^8.11.3"
  },
  "scripts": {
    "start": "node database-server.ts"
  }
}
```

---

## Estructura de un Servidor MCP

### Componentes Principales

```typescript
// 1. Importar dependencias
import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";
import { z } from "zod";

// 2. Crear instancia del servidor
const server = new McpServer({ 
  name: "my-server", 
  version: "1.0.0",
  description: "My custom MCP server" // Opcional
});

// 3. Definir herramientas (tools)
server.tool(
  "tool_name",                    // Nombre único y descriptivo
  {                               // Schema de parámetros con Zod
    param1: z.string().describe("Description"),
    param2: z.number().optional()
  },
  async (args) => {               // Función asíncrona de implementación
    // Lógica de la herramienta
    const result = await doSomething(args.param1);
    
    // Retornar resultado estructurado
    return {
      content: [
        { type: "text", text: JSON.stringify(result, null, 2) }
      ]
    };
  }
);

// 4. Definir recursos (opcional)
server.resource(
  "resource_name",
  "uri://pattern/{id}",
  async (uri, params) => {
    const data = await fetchData(params.id);
    return {
      contents: [
        { 
          uri: uri.href, 
          text: JSON.stringify(data) 
        }
      ]
    };
  }
);

// 5. Conectar y escuchar
const transport = new StdioServerTransport();
await server.connect(transport);

console.error("MCP Server running on stdio");
```

### Anatomía de una Herramienta

```typescript
server.tool(
  // 1. Nombre: debe ser descriptivo y único
  "search_database",
  
  // 2. Schema: validación de tipos con Zod
  {
    query: z
      .string()
      .min(1, "Query cannot be empty")
      .max(100, "Query too long")
      .describe("Search term to find in database"),
    
    limit: z
      .number()
      .min(1)
      .max(100)
      .default(10)
      .describe("Maximum number of results to return"),
    
    filters: z
      .array(z.string())
      .optional()
      .describe("Additional filters to apply")
  },
  
  // 3. Implementación: función asíncrona
  async ({ query, limit, filters }) => {
    try {
      // Validación adicional si es necesaria
      if (!query.trim()) {
        throw new Error("Query must not be empty");
      }
      
      // Ejecutar lógica principal
      const results = await database.search(query, limit, filters);
      
      // Retornar resultado exitoso
      return {
        content: [
          { 
            type: "text", 
            text: JSON.stringify({
              success: true,
              count: results.length,
              data: results
            }, null, 2) 
          }
        ]
      };
      
    } catch (error: any) {
      // Retornar error gracefully
      return {
        content: [
          { 
            type: "text", 
            text: `Error searching database: ${error.message}` 
          }
        ],
        isError: true
      };
    }
  }
);
```

---

## Mejores Prácticas

### ✅ DO: Nombres Descriptivos

```typescript
// ✅ BIEN
server.tool("search_users_by_email", ...)
server.tool("create_invoice_for_client", ...)
server.tool("get_weather_forecast", ...)

// ❌ MAL
server.tool("tool1", ...)
server.tool("do_stuff", ...)
server.tool("query", ...)
```

### ✅ DO: Validación con Zod

```typescript
// ✅ BIEN - Todos los inputs validados
{
  email: z.string().email("Invalid email format"),
  age: z.number().min(18).max(120),
  status: z.enum(["active", "inactive", "pending"])
}

// ❌ MAL - Sin validación
{
  email: z.string(),  // ¿Qué formato?
  age: z.number(),    // ¿Qué rango?
  status: z.string()  // ¿Qué valores?
}
```

### ✅ DO: Descripciones Detalladas

```typescript
// ✅ BIEN
query: z
  .string()
  .min(1)
  .describe("Search query for finding users by name or email")

// ❌ MAL
query: z.string()  // Sin descripción
```

### ✅ DO: Manejo de Errores

```typescript
// ✅ BIEN - Error messages accionables
try {
  await database.connect();
} catch (error: any) {
  return {
    content: [{
      type: "text",
      text: `Database connection failed: ${error.message}. Check DATABASE_URL environment variable.`
    }],
    isError: true
  };
}

// ❌ MAL - Error messages crípticos
catch (error) {
  return { content: [{ type: "text", text: "Error occurred" }] };
}
```

### ✅ DO: Stateless Design

```typescript
// ✅ BIEN - Cada llamada es independiente
server.tool("get_user", { id: z.string() }, async ({ id }) => {
  return await db.user.findUnique({ where: { id } });
});

// ❌ MAL - Depende de estado previo
let lastUserId = null;
server.tool("get_next_user", async () => {
  // ¿Qué pasa si se llama en desorden?
  return await db.user.findNext(lastUserId);
});
```

### ✅ DO: Testing

```typescript
// ✅ BIEN - Unit tests para herramientas
import { test } from "vitest";

test("search_users returns valid results", async () => {
  const result = await searchUsers({ query: "john", limit: 10 });
  expect(result.content).toBeDefined();
  expect(result.content[0].text).toContain("john");
});

test("search_users handles empty query", async () => {
  const result = await searchUsers({ query: "", limit: 10 });
  expect(result.isError).toBe(true);
});
```

---

## Referencia de API

### Métodos del Servidor MCP

#### `server.tool(name, schema, handler)`

Registra una nueva herramienta.

**Parámetros:**
- `name` (string): Nombre único de la herramienta
- `schema` (object): Schema de Zod para validación de inputs
- `handler` (function): Función asíncrona que ejecuta la herramienta

**Retorna:** void

**Ejemplo:**
```typescript
server.tool(
  "greet",
  { name: z.string() },
  async ({ name }) => ({
    content: [{ type: "text", text: `Hello, ${name}!` }]
  })
);
```

#### `server.resource(name, uriTemplate, handler)`

Registra un recurso legible.

**Parámetros:**
- `name` (string): Nombre del recurso
- `uriTemplate` (string): Patrón URI (e.g., `users/{id}`)
- `handler` (function): Función que retorna el contenido

**Ejemplo:**
```typescript
server.resource(
  "user_profile",
  "users/{id}",
  async (uri, params) => ({
    contents: [{
      uri: uri.href,
      text: JSON.stringify(await getUser(params.id))
    }]
  })
);
```

#### `server.prompt(name, schema, handler)`

Registra un prompt predefinido.

**Ejemplo:**
```typescript
server.prompt(
  "code_review",
  { code: z.string() },
  async ({ code }) => ({
    messages: [{
      role: "user",
      content: { type: "text", text: `Review this code:\n\n${code}` }
    }]
  })
);
```

### Tipos de Retorno

#### ToolResult Exitoso

```typescript
{
  content: [
    { type: "text", text: "Resultados de la búsqueda..." },
    { type: "image", data: "base64...", mimeType: "image/png" } // Opcional
  ],
  isError: false
}
```

#### ToolResult con Error

```typescript
{
  content: [
    { type: "text", text: "Error: Conexión fallida" }
  ],
  isError: true
}
```

### Variables de Entorno Comunes

```bash
# Configuración del servidor
MCP_SERVER_NAME="my-server"
MCP_SERVER_VERSION="1.0.0"

# Autenticación
DATABASE_URL="postgresql://user:pass@localhost:5432/db"
API_KEY="sk-..."
GITHUB_TOKEN="ghp_..."

# Rate Limiting
RATE_LIMIT_PER_MINUTE="60"
MAX_CONCURRENT_REQUESTS="10"
```

---

## Recursos Adicionales

- **Documentación Oficial MCP**: https://modelcontextprotocol.io
- **SDK de TypeScript**: https://github.com/modelcontextprotocol/typescript-sdk
- **Ejemplos de Servidores**: https://github.com/modelcontextprotocol/servers
- **Zod Documentation**: https://zod.dev

---

## Soporte

Para reportar bugs o solicitar features del MCP Builder Agent:

1. **GitHub Issues**: https://github.com/comgunner/picoclaw-agents/issues
2. **Discord**: Únete al servidor de la comunidad
3. **Documentación**: Consulta `docs/` para más guías

---

**Última actualización:** 26 de marzo de 2026  
**Mantenido por:** @comgunner
