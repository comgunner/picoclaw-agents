# MCP Builder Agent - Complete Guide

**Version:** 1.0.0  
**Category:** Specialized  
**Skill ID:** `specialized-mcp-builder`

---

## 📋 Table of Contents

1. [What is MCP Builder Agent?](#what-is-mcp-builder-agent)
2. [Use Cases](#use-cases)
3. [How to Activate](#how-to-activate)
4. [Practical Examples](#practical-examples)
5. [MCP Server Structure](#mcp-server-structure)
6. [Best Practices](#best-practices)
7. [API Reference](#api-reference)

---

## What is MCP Builder Agent?

**MCP Builder Agent** is a specialized skill for building **Model Context Protocol (MCP)** servers. MCP servers extend AI agent capabilities by exposing custom tools, resources, and prompts.

### Key Features

- 🛠️ **Tool Design**: Clear names, typed parameters, helpful descriptions
- 📚 **Resource Exposure**: Expose data sources agents can read
- 🔄 **Error Handling**: Graceful failures with actionable error messages
- 🔐 **Security**: Input validation, auth handling, rate limiting
- ✅ **Testing**: Unit tests for tools, integration tests for the server

---

## Use Cases

### 1. External API Integration

Create tools that allow your agent to interact with third-party APIs:

```typescript
// Example: Weather API MCP Server
server.tool("get_weather", { 
  city: z.string(), 
  units: z.enum(["celsius", "fahrenheit"]).default("celsius") 
}, async ({ city, units }) => {
  const weather = await fetchWeatherAPI(city, units);
  return { content: [{ type: "text", text: JSON.stringify(weather, null, 2) }] };
});
```

### 2. Database Access

Safely expose data from your database:

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

### 3. Workflow Automation

Automate repetitive business tasks:

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

### 4. File System Integration

Enable controlled file read/write operations:

```typescript
server.tool("read_config", { 
  path: z.string().refine(p => p.startsWith('/safe/')) 
}, async ({ path }) => {
  const content = await fs.readFile(path, 'utf-8');
  return { content: [{ type: "text", text: content }] };
});
```

---

## How to Activate

### Method 1: Via CLI

```bash
# Invoke MCP Builder Agent directly
picoclaw-agents agent -m "I need to create an MCP server to integrate with Stripe API"
```

### Method 2: Via Chat (Telegram/Discord)

Send a message to your agent:

```
@agent I need an MCP server that allows:
1. Search products in my database
2. Create purchase orders
3. Check shipment status

Can you help me build it?
```

### Method 3: Configure as Default Skill

In `~/.picoclaw/config.json`:

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

## Practical Examples

### Example 1: MCP Server for GitHub API

**Goal**: Create tools to interact with GitHub API.

**Prompt for MCP Builder:**

```
I need an MCP server that allows:
1. Search repositories by name
2. Get latest commits from a repo
3. Create issues in a repository
4. List open pull requests

Requirements:
- Authentication via GitHub token
- Rate limiting to avoid blocks
- Handle 404 and 403 errors
```

**Expected Result:**

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

// Initialize Octokit with token
const octokit = new Octokit({ 
  auth: process.env.GITHUB_TOKEN 
});

// Tool: Search repositories
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

// Tool: Get commits
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

// Tool: Create issue
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

// Start server
const transport = new StdioServerTransport();
await server.connect(transport);
```

**Configuration in PicoClaw:**

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

### Example 2: MCP Server for PostgreSQL Database

**Goal**: Expose user and product data.

**Prompt for MCP Builder:**

```
Create an MCP server to connect to PostgreSQL with:
1. Search users by email or name
2. List products by category
3. Get orders for a user
4. Create new purchase order

Use Zod validation for all inputs.
```

**Generated Code:**

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

// Configure connection pool
const pool = new Pool({
  connectionString: process.env.DATABASE_URL,
  max: 20,
  idleTimeoutMillis: 30000,
  connectionTimeoutMillis: 2000,
});

// Tool: Search users
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

// Tool: List products
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

// Start server
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

## MCP Server Structure

### Main Components

```typescript
// 1. Import dependencies
import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";
import { z } from "zod";

// 2. Create server instance
const server = new McpServer({ 
  name: "my-server", 
  version: "1.0.0",
  description: "My custom MCP server" // Optional
});

// 3. Define tools
server.tool(
  "tool_name",                    // Unique, descriptive name
  {                               // Parameter schema with Zod
    param1: z.string().describe("Description"),
    param2: z.number().optional()
  },
  async (args) => {               // Async implementation function
    // Tool logic
    const result = await doSomething(args.param1);
    
    // Return structured result
    return {
      content: [
        { type: "text", text: JSON.stringify(result, null, 2) }
      ]
    };
  }
);

// 4. Define resources (optional)
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

// 5. Connect and listen
const transport = new StdioServerTransport();
await server.connect(transport);

console.error("MCP Server running on stdio");
```

### Tool Anatomy

```typescript
server.tool(
  // 1. Name: must be descriptive and unique
  "search_database",
  
  // 2. Schema: type validation with Zod
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
  
  // 3. Implementation: async function
  async ({ query, limit, filters }) => {
    try {
      // Additional validation if needed
      if (!query.trim()) {
        throw new Error("Query must not be empty");
      }
      
      // Execute main logic
      const results = await database.search(query, limit, filters);
      
      // Return successful result
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
      // Return error gracefully
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

## Best Practices

### ✅ DO: Descriptive Names

```typescript
// ✅ GOOD
server.tool("search_users_by_email", ...)
server.tool("create_invoice_for_client", ...)
server.tool("get_weather_forecast", ...)

// ❌ BAD
server.tool("tool1", ...)
server.tool("do_stuff", ...)
server.tool("query", ...)
```

### ✅ DO: Zod Validation

```typescript
// ✅ GOOD - All inputs validated
{
  email: z.string().email("Invalid email format"),
  age: z.number().min(18).max(120),
  status: z.enum(["active", "inactive", "pending"])
}

// ❌ BAD - No validation
{
  email: z.string(),  // What format?
  age: z.number(),    // What range?
  status: z.string()  // What values?
}
```

### ✅ DO: Detailed Descriptions

```typescript
// ✅ GOOD
query: z
  .string()
  .min(1)
  .describe("Search query for finding users by name or email")

// ❌ BAD
query: z.string()  // No description
```

### ✅ DO: Error Handling

```typescript
// ✅ GOOD - Actionable error messages
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

// ❌ BAD - Cryptic error messages
catch (error) {
  return { content: [{ type: "text", text: "Error occurred" }] };
}
```

### ✅ DO: Stateless Design

```typescript
// ✅ GOOD - Each call is independent
server.tool("get_user", { id: z.string() }, async ({ id }) => {
  return await db.user.findUnique({ where: { id } });
});

// ❌ BAD - Depends on previous state
let lastUserId = null;
server.tool("get_next_user", async () => {
  // What if called out of order?
  return await db.user.findNext(lastUserId);
});
```

### ✅ DO: Testing

```typescript
// ✅ GOOD - Unit tests for tools
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

## API Reference

### MCP Server Methods

#### `server.tool(name, schema, handler)`

Registers a new tool.

**Parameters:**
- `name` (string): Unique tool name
- `schema` (object): Zod schema for input validation
- `handler` (function): Async function that executes the tool

**Returns:** void

**Example:**
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

Registers a readable resource.

**Parameters:**
- `name` (string): Resource name
- `uriTemplate` (string): URI pattern (e.g., `users/{id}`)
- `handler` (function): Function that returns content

**Example:**
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

Registers a predefined prompt.

**Example:**
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

### Return Types

#### Successful ToolResult

```typescript
{
  content: [
    { type: "text", text: "Search results..." },
    { type: "image", data: "base64...", mimeType: "image/png" } // Optional
  ],
  isError: false
}
```

#### Error ToolResult

```typescript
{
  content: [
    { type: "text", text: "Error: Connection failed" }
  ],
  isError: true
}
```

### Common Environment Variables

```bash
# Server configuration
MCP_SERVER_NAME="my-server"
MCP_SERVER_VERSION="1.0.0"

# Authentication
DATABASE_URL="postgresql://user:pass@localhost:5432/db"
API_KEY="sk-..."
GITHUB_TOKEN="ghp_..."

# Rate Limiting
RATE_LIMIT_PER_MINUTE="60"
MAX_CONCURRENT_REQUESTS="10"
```

---

## Additional Resources

- **Official MCP Documentation**: https://modelcontextprotocol.io
- **TypeScript SDK**: https://github.com/modelcontextprotocol/typescript-sdk
- **Server Examples**: https://github.com/modelcontextprotocol/servers
- **Zod Documentation**: https://zod.dev

---

## Support

To report bugs or request features for MCP Builder Agent:

1. **GitHub Issues**: https://github.com/comgunner/picoclaw-agents/issues
2. **Discord**: Join the community server
3. **Documentation**: Check `docs/` for more guides

---

**Last updated:** March 26, 2026  
**Maintained by:** @comgunner
