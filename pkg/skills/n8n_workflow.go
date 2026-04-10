// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)

package skills

import (
	"strings"
)

// bt is a backtick constant for code block formatting (imported from common.go)

// N8NWorkflowSkill implements native skill for n8n workflow automation.
// Based on official n8n documentation - Workflow Automation Expert
type N8NWorkflowSkill struct {
	workspace string
}

// NewN8NWorkflowSkill creates a new N8NWorkflowSkill instance.
func NewN8NWorkflowSkill(workspace string) *N8NWorkflowSkill {
	return &N8NWorkflowSkill{
		workspace: workspace,
	}
}

// Name returns the skill identifier name.
func (n *N8NWorkflowSkill) Name() string {
	return "n8n_workflow"
}

// Description returns a brief description of the skill.
func (n *N8NWorkflowSkill) Description() string {
	return "n8n Workflow Automation Expert - Create production-ready workflows with valid JSON structure ready for import/export."
}

// GetInstructions returns the complete workflow creation guidelines.
func (n *N8NWorkflowSkill) GetInstructions() string {
	return n8nWorkflowInstructions
}

// GetAntiPatterns returns common n8n workflow anti-patterns.
func (n *N8NWorkflowSkill) GetAntiPatterns() string {
	return n8nWorkflowAntiPatterns
}

// GetExamples returns concrete n8n workflow JSON examples.
func (n *N8NWorkflowSkill) GetExamples() string {
	return n8nWorkflowExamples
}

// BuildSkillContext returns the complete skill context for prompt injection.
func (n *N8NWorkflowSkill) BuildSkillContext() string {
	var parts []string

	parts = append(parts, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	parts = append(parts, "🚀 NATIVE SKILL: n8n Workflow Automation Expert")
	parts = append(parts, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	parts = append(parts, "")
	parts = append(parts, "**ROLE:** n8n Workflow Automation Expert")
	parts = append(parts, "")
	parts = append(
		parts,
		"**OBJECTIVE:** Create production-ready n8n workflows with valid JSON structure ready for import/export via File or URL.",
	)
	parts = append(parts, "")
	parts = append(parts, n.GetInstructions())
	parts = append(parts, "")
	parts = append(parts, n.GetAntiPatterns())
	parts = append(parts, "")
	parts = append(parts, n.GetExamples())

	return strings.Join(parts, "\n")
}

// BuildSummary returns an XML summary for compact context injection.
func (n *N8NWorkflowSkill) BuildSummary() string {
	return `<skill name="n8n_workflow" type="native">
  <purpose>n8n Workflow Automation Expert</purpose>
  <pattern>Use for creating n8n workflows, JSON generation, automation design</pattern>
  <nodes>Webhook, HTTP Request, Function, Gmail, Slack, Sheets, PostgreSQL, etc.</nodes>
  <features>Import/Export, Expressions, Credentials, Version Control</features>
</skill>`
}

// ============================================================================
// DOCUMENTATION CONSTANTS
// ============================================================================

const n8nWorkflowInstructions = `## ROLE & OBJECTIVE

**Role:** n8n Workflow Automation Expert

**Objective:** Create production-ready n8n workflows with valid JSON structure ready for import/export via File or URL.

## WORKFLOW JSON STRUCTURE

### Root Object
` + bt + bt + bt + `json
{
  "name": "My Workflow",
  "nodes": [...],
  "connections": {...},
  "settings": {...},
  "active": false,
  "tags": [...]
}
` + bt + bt + bt + `

### Required Fields
| Field | Type | Description |
|-------|------|-------------|
| ` + bt + `nodes` + bt + ` | array | Array of node objects |
| ` + bt + `connections` + bt + ` | object | Node connection mappings |

### Optional Fields
| Field | Type | Default | Description |
|-------|------|---------|-------------|
| ` + bt + `name` + bt + ` | string | | Workflow name |
| ` + bt + `settings` + bt + ` | object | | Workflow settings |
| ` + bt + `active` + bt + ` | boolean | false | Whether workflow is active |
| ` + bt + `tags` + bt + ` | array | [] | Workflow tags |
| ` + bt + `id` + bt + ` | string | | Workflow ID (auto-generated) |

## WEBHOOK INTEGRATION (CRITICAL)

### Setup & Connection
1. **Create Webhook Node**: Use ` + bt + `n8n-nodes-base.webhook` + bt + ` as the trigger.
2. **HTTP Method**: Usually ` + bt + `POST` + bt + ` to receive JSON payloads from PicoClaw.
3. **Authentication**:
   - Recommended: ` + bt + `Header Auth` + bt + `.
   - Header Name: ` + bt + `X-Webhook-Secret` + bt + `.
   - This ensures only PicoClaw (or authorized users) can trigger the flow.
4. **Endpoint URLs**:
   - **Internal**: ` + bt + `http://n8n:5678/webhook/[PATH]` + bt + ` (Default for Docker setups).
   - **Production**: Public URL provided by n8n cloud or your tunnel/domain.

### Mini-Tutorial: Triggering n8n from PicoClaw
To trigger a workflow using the ` + bt + `exec` + bt + ` tool:
` + bt + bt + bt + `bash
curl -X POST "http://n8n:5678/webhook/your-path" \
  -H "Content-Type: application/json" \
  -H "X-Webhook-Secret: YOUR_SECRET_HERE" \
  -d '{"task": "Run data sync", "user": "pico-user"}'
` + bt + bt + bt + `

## NODE STRUCTURE

### Required Node Fields
| Field | Type | Description |
|-------|------|-------------|
| ` + bt + `name` + bt + ` | string | Unique node identifier |
| ` + bt + `type` + bt + ` | string | Node type (e.g., ` + bt + `n8n-nodes-base.webhook` + bt + `) |
| ` + bt + `typeVersion` + bt + ` | number | Node version (usually 1 or 2.1 depending on node) |
| ` + bt + `position` + bt + ` | array | [x, y] canvas coordinates |
| ` + bt + `parameters` + bt + ` | object | Node-specific configuration |

### Optional Node Fields
| Field | Type | Description |
|-------|------|-------------|
| ` + bt + `credentials` + bt + ` | object | Credential references |
| ` + bt + `disabled` + bt + ` | boolean | Whether node is disabled |
| ` + bt + `notes` + bt + ` | string | Node documentation |
| ` + bt + `id` + bt + ` | string | Node unique ID (auto-generated) |

## CONNECTION STRUCTURE

` + bt + bt + bt + `javascript
connections: {
  "SourceNodeName": {
    "main": [[           // Output type (main, audio, text)
      {                  // Output index 0
        "node": "TargetNodeName",
        "type": "main",  // Input type
        "index": 0       // Input index
      }
    ]],
    "false": []          // For IF nodes (false branch)
  }
}
` + bt + bt + bt + `

## NODE LIBRARY (Common Nodes)

### Trigger Nodes

#### 1. Webhook
` + bt + bt + bt + `json
{
  "name": "Webhook",
  "type": "n8n-nodes-base.webhook",
  "typeVersion": 1,
  "position": [240, 160],
  "parameters": {
    "path": "pico-trigger",
    "httpMethod": "POST",
    "authentication": "headerAuth",
    "options": {}
  }
}
` + bt + bt + bt + `

#### 2. Telegram Trigger
` + bt + bt + bt + `json
{
  "name": "Telegram Trigger",
  "type": "n8n-nodes-base.telegramTrigger",
  "typeVersion": 1.1,
  "position": [240, 160],
  "parameters": {
    "updates": ["message"]
  }
}
` + bt + bt + bt + `

### Action Nodes

#### 3. AI Agent (LangChain)
` + bt + bt + bt + `json
{
  "name": "AI Agent",
  "type": "@n8n/n8n-nodes-langchain.agent",
  "typeVersion": 1.1,
  "position": [480, 160],
  "parameters": {
    "options": {
      "systemMessage": "You are a helpful assistant."
    }
  }
}
` + bt + bt + bt + `

#### 4. PostgreSQL (Supabase/Local)
` + bt + bt + bt + `json
{
  "name": "PostgreSQL",
  "type": "n8n-nodes-base.postgres",
  "typeVersion": 2.5,
  "position": [480, 160],
  "parameters": {
    "operation": "executeQuery",
    "query": "SELECT * FROM my_table LIMIT 10"
  }
}
` + bt + bt + bt + `

... (Refer to official n8n docs for other nodes) ...

## EXPRESSION SYNTAX

### Accessing Data
- ` + bt + `{{$json.field}}` + bt + ` - Access field from previous node
- ` + bt + `{{$node["NodeName"].json.field}}` + bt + ` - Access specific node
- ` + bt + `{{$input.first().json.field}}` + bt + ` - First input item
- ` + bt + `{{$input.all()}}` + bt + ` - All input items

### Date Helpers
- ` + bt + `{{$now}}` + bt + ` - Current datetime
- ` + bt + `{{$today}}` + bt + ` - Current date
- ` + bt + `{{$formatDate(date, format)}}` + bt + ` - Format date

## IMPORT/EXPORT METHODS

### Import Methods

#### 1. Import from File
1. Open n8n editor -> Click three dots (top-right) -> "Import from File".

#### 2. Copy/Paste (Nodes Only)
1. Select nodes on canvas -> Ctrl+C -> Paste into target workflow (Ctrl+V).

## BEST PRACTICES

### Workflow Design
- Start with trigger node (Webhook, Schedule, Manual).
- Use descriptive node names.
- Add notes for complex logic.
- Use **Sticky Notes** for documentation.

### Error Handling
- Add ` + bt + `n8n-nodes-base.errorTrigger` + bt + ` nodes for critical workflows.
- Use ` + bt + `IF` + bt + ` nodes to validate data before processing.

## SECURITY CONSIDERATIONS

### Sanitization Checklist
` + bt + bt + bt + `json
// BEFORE (contains sensitive data)
{
  "credentials": {
    "httpHeaderAuth": {
      "id": "abc123",
      "name": "Production API Key"  // ❌ Remove
    }
  },
  "parameters": {
    "headers": {
      "Authorization": "Bearer sk-123456"  // ❌ Remove
    }
  }
}

// AFTER (sanitized)
{
  "credentials": {
    "httpHeaderAuth": {
      "id": "abc123"  // ✅ ID is safe
    }
  }
}
` + bt + bt + bt + `
`

const n8nWorkflowAntiPatterns = `## WORKFLOW DESIGN ANTI-PATTERNS

### ❌ No Trigger Node
` + bt + bt + bt + `json
// BAD: Workflow with no trigger
{
  "nodes": [
    {"name": "HTTP Request", "type": "n8n-nodes-base.httpRequest"}
  ]
}

// GOOD: Start with trigger
{
  "nodes": [
    {"name": "Webhook", "type": "n8n-nodes-base.webhook"},
    {"name": "HTTP Request", "type": "n8n-nodes-base.httpRequest"}
  ],
  "connections": {
    "Webhook": {"main": [[{"node": "HTTP Request"}]]}
  }
}
` + bt + bt + bt + `

### ❌ Hardcoded Credentials
` + bt + bt + bt + `json
// BAD: Credentials in parameters
{
  "parameters": {
    "headers": {
      "Authorization": "Bearer sk-1234567890"
    }
  }
}

// GOOD: Use credentials reference
{
  "credentials": {
    "httpHeaderAuth": "myAuthCredentials"
  }
}
` + bt + bt + bt + `

### ❌ No Error Handling
` + bt + bt + bt + `json
// BAD: No error handling for critical API call
{
  "nodes": [
    {"name": "HTTP Request", "type": "n8n-nodes-base.httpRequest"}
  ]
}

// GOOD: Add error trigger
{
  "nodes": [
    {"name": "Error Trigger", "type": "n8n-nodes-base.errorTrigger"},
    {"name": "Slack", "type": "n8n-nodes-base.slack"}
  ]
}
` + bt + bt + bt + `

### ❌ Complex Expressions Without Comments
` + bt + bt + bt + `json
// BAD: Complex expression without documentation
{
  "parameters": {
    "jsCode": "return $json.data.map(x => x.value * 2).filter(y => y > 100)"
  }
}

// GOOD: Add notes field
{
  "parameters": {
    "jsCode": "// Transform and filter data\nreturn $json.data.map(...)"
  },
  "notes": "Doubles values and filters results > 100"
}
` + bt + bt + bt + `

## SECURITY ANTI-PATTERNS

### ❌ Exposed API Keys
` + bt + bt + bt + `json
// BAD: API key visible in workflow
{
  "parameters": {
    "url": "https://api.example.com/data?key=sk-123456"
  }
}

// GOOD: Use credentials
{
  "credentials": {
    "httpQueryAuth": "myApiCredentials"
  }
}
` + bt + bt + bt + `

### ❌ Webhook Path with Secrets
` + bt + bt + bt + `json
// BAD: Secret in webhook path
{
  "parameters": {
    "path": "/webhook/secret-token-12345"
  }
}

// GOOD: Use webhook authentication
{
  "parameters": {
    "path": "/webhook",
    "options": {
      "httpHeaderAuth": "myAuthCredentials"
    }
  }
}
` + bt + bt + bt + `

### ❌ Sharing Workflows with Credential Names
` + bt + bt + bt + `json
// BAD: Credential name reveals environment
{
  "credentials": {
    "gmailOAuth2Api": {
      "id": "1",
      "name": "Production Gmail - admin@company.com"  // ❌
    }
  }
}

// GOOD: Generic name or remove
{
  "credentials": {
    "gmailOAuth2Api": {
      "id": "1"  // ✅ ID only
    }
  }
}
` + bt + bt + bt + `

## PERFORMANCE ANTI-PATTERNS

### ❌ No Batch Processing
` + bt + bt + bt + `json
// BAD: Process 1000 items one by one
{
  "nodes": [
    {"name": "HTTP Request", "type": "n8n-nodes-base.httpRequest"}
  ]
}

// GOOD: Split in batches
{
  "nodes": [
    {"name": "Split In Batches", "type": "n8n-nodes-base.splitInBatches",
     "parameters": {"batchSize": 10}}
  ]
}
` + bt + bt + bt + `

### ❌ Unnecessary API Calls
` + bt + bt + bt + `json
// BAD: Call API for each item
{
  "connections": {
    "Loop": {"main": [[{"node": "HTTP Request"}]]}
  }
}

// GOOD: Batch API call
{
  "nodes": [
    {"name": "Function", "parameters": {
      "jsCode": "// Combine all IDs into one request\nconst ids = $input.all().map(i => i.json.id);\nreturn [{json: {result: await api.get(ids)}}];"
    }}
  ]
}
` + bt + bt + bt + `

## CONNECTION ANTI-PATTERNS

### ❌ Missing Connections
` + bt + bt + bt + `json
// BAD: Nodes not connected
{
  "nodes": [
    {"name": "Webhook"},
    {"name": "HTTP Request"}
  ]
  // No connections object!
}

// GOOD: Proper connections
{
  "connections": {
    "Webhook": {
      "main": [[{"node": "HTTP Request", "type": "main", "index": 0}]]
    }
  }
}
` + bt + bt + bt + `

### ❌ Wrong Connection Syntax
` + bt + bt + bt + `json
// BAD: Incorrect nested array structure
{
  "connections": {
    "Webhook": {
      "main": [{"node": "HTTP Request"}]  // Missing inner array
    }
  }
}

// GOOD: Correct structure
{
  "connections": {
    "Webhook": {
      "main": [[{"node": "HTTP Request", "type": "main", "index": 0}]]
    }
  }
}
` + bt + bt + bt + `

## VERSION CONTROL ANTI-PATTERNS

### ❌ No Workflow Versioning
` + bt + bt + bt + `
// BAD: Manual changes without backup
- Edit workflow in UI
- No export to file
- No Git commit
- Changes lost if server crashes

// GOOD: Version controlled
- Export workflow to JSON
- Commit to Git with message
- Tag production versions
- Can rollback if needed
` + bt + bt + bt + `

### ❌ Large Workflow Files
` + bt + bt + bt + `json
// BAD: 5000-line workflow (hard to review)
{
  "nodes": [/* 200 nodes */]
}

// GOOD: Modular workflows
- Main workflow (orchestrator)
- Sub-workflows (specific tasks)
- Call sub-workflows via Execute Workflow node
` + bt + bt + bt + `
`

const n8nWorkflowExamples = `## EXAMPLE 1: BASIC WORKFLOW - Google Sheets → Gmail

**Request:** "Create n8n workflow to send welcome emails from Google Sheets"

**Expert Response:**

### Workflow JSON (Ready to Import)

` + bt + bt + bt + `json
{
  "name": "Google Sheets to Gmail Welcome",
  "nodes": [
    {
      "parameters": {},
      "name": "Start",
      "type": "n8n-nodes-base.start",
      "typeVersion": 1,
      "position": [240, 160]
    },
    {
      "parameters": {
        "spreadsheetId": "YOUR_SPREADSHEET_ID",
        "range": "Sheet1!A:B",
        "options": {}
      },
      "name": "Google Sheet",
      "type": "n8n-nodes-base.googleSheets",
      "typeVersion": 1,
      "position": [480, 160],
      "credentials": {
        "googleSheetsOAuth2Api": "YOUR_GOOGLE_SHEETS_CREDENTIALS"
      }
    },
    {
      "parameters": {
        "fromEmail": "YOUR_GMAIL_ADDRESS",
        "toEmail": "={{$node[\\"Google Sheet\\"].json[\\"email\\"]}}",
        "subject": "Welcome to Our Platform!",
        "text": "Hey {{$node[\\"Google Sheet\\"].json[\\"name\\"]}}!\\n\\nWelcome to our platform. We're excited to have you!"
      },
      "name": "Gmail",
      "type": "n8n-nodes-base.gmail",
      "typeVersion": 1,
      "position": [720, 160],
      "credentials": {
        "gmailOAuth2Api": "YOUR_GMAIL_CREDENTIALS"
      }
    }
  ],
  "connections": {
    "Start": {
      "main": [[{"node": "Google Sheet", "type": "main", "index": 0}]]
    },
    "Google Sheet": {
      "main": [[{"node": "Gmail", "type": "main", "index": 0}]]
    }
  },
  "settings": {
    "executionOrder": "v1"
  }
}
` + bt + bt + bt + `

### Import Instructions

1. **Save JSON:** Copy above JSON to file ` + bt + `sheets-to-gmail.json` + bt + `
2. **Open n8n:** Navigate to n8n editor
3. **Import:** Three dots → Import from File → Select JSON file
4. **Configure Credentials:**
   - Google Sheets OAuth2 API
   - Gmail OAuth2 API
5. **Update Parameters:**
   - Replace ` + bt + `YOUR_SPREADSHEET_ID` + bt + `
   - Replace ` + bt + `YOUR_GMAIL_ADDRESS` + bt + `
6. **Test:** Click "Execute Workflow"

---

## EXAMPLE 2: WEBHOOK TRIGGER WITH SECURITY

**Request:** "Configure a secure webhook to receive data from PicoClaw"

**Expert Response:**

### Workflow JSON

` + bt + bt + bt + `json
{
  "name": "Secure PicoClaw Inbox",
  "nodes": [
    {
      "parameters": {
        "path": "pico-inbox",
        "httpMethod": "POST",
        "authentication": "headerAuth",
        "options": {}
      },
      "name": "Webhook",
      "type": "n8n-nodes-base.webhook",
      "typeVersion": 1,
      "position": [240, 160],
      "credentials": {
        "httpHeaderAuth": "YOUR_AUTH_CREDENTIALID"
      }
    }
  ],
  "connections": {}
}
` + bt + bt + bt + `

### Setup in n8n
1. Create a ` + bt + `Header Auth` + bt + ` credential.
2. Name: ` + bt + `X-Webhook-Secret` + bt + `.
3. Value: ` + bt + `YOUR_SECURE_TOKEN` + bt + `.

---

## EXAMPLE 3: ADVANCED AI AGENT (n8n-claw Pattern)

**Request:** "Create an AI Agent in n8n that can search Postgres and reply to Telegram"

**Expert Response:**

### Workflow JSON

` + bt + bt + bt + `json
{
  "name": "n8n-claw AI Assistant",
  "nodes": [
    {
      "parameters": {
        "updates": ["message"]
      },
      "name": "Telegram Trigger",
      "type": "n8n-nodes-base.telegramTrigger",
      "typeVersion": 1.1,
      "position": [0, 0]
    },
    {
      "parameters": {
        "operation": "executeQuery",
        "query": "SELECT content FROM memory WHERE user_id = '{{$json.message.from.id}}'"
      },
      "name": "Load Memory",
      "type": "n8n-nodes-base.postgres",
      "typeVersion": 2.5,
      "position": [250, 0]
    },
    {
      "parameters": {
        "options": {
          "systemMessage": "You are n8n-claw. Use the loaded memory to assist the user."
        }
      },
      "name": "AI Agent",
      "type": "@n8n/n8n-nodes-langchain.agent",
      "typeVersion": 1.1,
      "position": [500, 0]
    },
    {
      "parameters": {
        "chatId": "={{$node[\\"Telegram Trigger\\"].json.message.chat.id}}",
        "text": "={{$json.output}}"
      },
      "name": "Telegram Reply",
      "type": "n8n-nodes-base.telegram",
      "typeVersion": 1.2,
      "position": [750, 0]
    }
  ],
  "connections": {
    "Telegram Trigger": { "main": [[{ "node": "Load Memory", "index": 0 }]] },
    "Load Memory": { "main": [[{ "node": "AI Agent", "index": 0 }]] },
    "AI Agent": { "main": [[{ "node": "Telegram Reply", "index": 0 }]] }
  }
}
` + bt + bt + bt + `
`
