# Native Skills Guide

**Version:** 3.11.1
**Last Updated:** March 26, 2026
**Total Native Skills:** 14

## Overview

Native Skills are specialized role definitions compiled directly into the PicoClaw binary. Unlike Tools (which the LLM calls to perform actions), Skills define **who the agent is** — injecting personality, expertise, and behavioral guidelines into the system prompt.

---

## Skills vs Tools: Key Differences

| Dimension | Native Skill | Tool |
|-----------|-------------|------|
| **What it is** | Role instructions injected into system prompt | Function the LLM invokes via `tool_use` |
| **How LLM interacts** | "Is" this role / "has" this expertise | "Calls" this function to perform action |
| **Effect** | Modifies agent personality and behavior | Executes real action (I/O, API, shell, etc.) |
| **Interface** | `Name()`, `Description()`, `BuildSkillContext()` | `Name()`, `Description()`, `Parameters()`, `Execute()` |
| **Registration** | `listNativeSkills()` in `loader.go` | `ToolRegistry.Register()` in `instance.go` |
| **Example** | `backend_developer` → agent "is" a backend dev | `read_file` → agent "calls" file reader |
| **Output** | String injected into system prompt | `*ToolResult` with structured data |
| **Runtime dependencies** | None — compiled into binary | May require external binaries (git, docker, etc.) |

---

## Flow Comparison

### Skills Flow
```
config.json → agent.skills: ["backend_developer"]
  → SkillsLoader.LoadSkill("backend_developer")
    → listNativeSkills() — finds compiled skill
      → LoadNativeBackendDeveloperSkill()
        → BackendDeveloperSkill.BuildSkillContext()
          → string injected into LLM system prompt
```

### Tools Flow
```
config.json → agent.tools_override: ["read_file", "exec"]
  → ToolRegistry.Register(FilesystemTool{})
    → ToolRegistry.ToProviderDefs()
      → []providers.ToolDefinition sent to LLM
        → LLM generates: tool_use { name: "read_file", args: {...} }
          → ToolRegistry.Execute("read_file", ctx, args)
            → *ToolResult returned to LLM
```

**Rule of Thumb:** Engineering role skills (backend_developer, devops_engineer, etc.) are **Native Skills**. The LLM doesn't "call" them — it "has" them as part of its identity.

---

## Available Native Skills (v3.11.1)

PicoClaw v3.11.1 includes **14 native skills**:

### Engineering Role Skills (7 skills)

| Skill Name | Purpose | Best For |
|------------|---------|----------|
| `backend_developer` | Backend development expert | REST APIs, databases, microservices, security |
| `frontend_developer` | Frontend development expert | React, Vue, performance, accessibility |
| `devops_engineer` | DevOps expert | CI/CD, Kubernetes, Terraform, monitoring |
| `security_engineer` | Security expert | OWASP, penetration testing, threat modeling |
| `qa_engineer` | QA expert | Test automation, coverage analysis, E2E testing |
| `data_engineer` | Data engineering expert | ETL pipelines, data warehouses, streaming |
| `ml_engineer` | ML/AI expert | Model training, deployment, MLOps |

### General Purpose Skills (4 skills)

| Skill Name | Purpose | Best For |
|------------|---------|----------|
| `fullstack_developer` | Full-stack development assistant | General coding, architecture, best practices |
| `researcher` | Deep research agent | Web search, source evaluation, synthesis |
| `queue_batch` | Background task delegation | Fire-and-forget heavy tasks |
| `agent_team_workflow` | Multi-agent orchestrator | Team coordination, task delegation |

### Integration Skills (3 skills)

| Skill Name | Purpose | Best For |
|------------|---------|----------|
| `binance_mcp` | Binance trading integration | Crypto trading, market data |
| `n8n_workflow` | n8n automation expert | Workflow creation, JSON validation |
| `odoo_developer` | Odoo architect & QA engineer | Odoo ecosystems, L10n-Mexico, CFDI 4.0 |

---

## Configuration

### Single Specialized Agent

```json
{
  "agents": {
    "list": [
      {
        "id": "backend_dev",
        "name": "Backend Developer",
        "model": "deepseek-chat",
        "skills": ["backend_developer"],
        "tools_override": ["read_file", "write_file", "edit_file", "exec", "web_search"],
        "subagents": {}
      }
    ]
  }
}
```

### Orchestrator with Specialized Subagents

```json
{
  "agents": {
    "list": [
      {
        "id": "tech_lead",
        "name": "Technical Lead",
        "model": "deepseek-chat",
        "skills": ["fullstack_developer", "agent_team_workflow"],
        "subagents": {
          "allow_agents": ["backend_dev", "frontend_dev", "devops_eng", "qa_eng"],
          "max_spawn_depth": 2,
          "max_children_per_agent": 3
        }
      },
      {
        "id": "backend_dev",
        "name": "Backend Developer",
        "model": "deepseek-chat",
        "skills": ["backend_developer"],
        "tools_override": ["read_file", "write_file", "edit_file", "exec"]
      },
      {
        "id": "frontend_dev",
        "name": "Frontend Developer",
        "model": "deepseek-chat",
        "skills": ["frontend_developer"],
        "tools_override": ["read_file", "write_file", "edit_file"]
      },
      {
        "id": "devops_eng",
        "name": "DevOps Engineer",
        "model": "deepseek-chat",
        "skills": ["devops_engineer"],
        "tools_override": ["read_file", "write_file", "exec"]
      },
      {
        "id": "qa_eng",
        "name": "QA Engineer",
        "model": "deepseek-chat",
        "skills": ["qa_engineer"],
        "tools_override": ["read_file", "write_file", "exec"]
      }
    ]
  }
}
```

---

## Skill Content Structure

Each native skill includes:

### 1. Role Definition
```
**ROLE:** Expert Backend Developer specializing in...
```

### 2. Core Responsibilities
```markdown
## CORE RESPONSIBILITIES

### 1. API Design & Development
- Design RESTful APIs with clear contracts
- Implement GraphQL schemas when appropriate
...
```

### 3. Technology Stack
```markdown
## TECHNOLOGY STACK

### Languages
- **Primary**: Go, Python, Node.js, Java
...
```

### 4. Best Practices
```markdown
## BEST PRACTICES

### Code Organization
src/
├── controllers/    # Request handlers
├── services/       # Business logic
...
```

### 5. Quality Checklist
```markdown
## QUALITY CHECKLIST

Before considering a feature complete:

- [ ] API endpoints documented
- [ ] Input validation implemented
...
```

### 6. Anti-Patterns (with code examples)
```markdown
## ANTI-PATTERNS

### ❌ N+1 Query Problem
BAD:  for user in users:
        posts = db.query("SELECT * FROM posts WHERE user_id = ?", user.id)

GOOD: posts = db.query("""
        SELECT * FROM posts 
        WHERE user_id IN (?)
        """, [u.id for u in users])
```

### 7. Concrete Examples (with code)
```markdown
## EXAMPLE 1: CREATE REST API ENDPOINT

**Request:** "Create Express endpoint to create users with validation"

**Expert Response:**

```javascript
// routes/users.js
const express = require('express')
...
```

---

## When to Use Each Skill

### Backend Developer
Use when:
- Designing REST or GraphQL APIs
- Modeling database schemas
- Implementing authentication/authorization
- Optimizing database queries
- Building microservices

### Frontend Developer
Use when:
- Creating React/Vue components
- Implementing state management
- Optimizing performance (Core Web Vitals)
- Ensuring accessibility (WCAG)
- Building responsive layouts

### DevOps Engineer
Use when:
- Creating Kubernetes manifests
- Writing Terraform modules
- Setting up CI/CD pipelines
- Configuring monitoring/alerting
- Designing disaster recovery

### Security Engineer
Use when:
- Conducting threat modeling
- Performing security audits
- Implementing OWASP controls
- Reviewing code for vulnerabilities
- Ensuring compliance (SOC2, GDPR)

### QA Engineer
Use when:
- Designing test strategy
- Writing unit/integration/E2E tests
- Setting up test automation
- Analyzing code coverage
- Implementing quality gates

### Data Engineer
Use when:
- Building ETL/ELT pipelines
- Modeling data warehouses
- Implementing streaming pipelines
- Ensuring data quality
- Setting up data governance

### ML Engineer
Use when:
- Training ML models
- Deploying models to production
- Setting up MLOps pipelines
- Engineering features
- Monitoring model drift

---

## Combining Skills

### Multi-Skill Agents

You can assign multiple skills to a single agent:

```json
{
  "id": "fullstack_tech_lead",
  "name": "Full-Stack Tech Lead",
  "model": "deepseek-chat",
  "skills": ["fullstack_developer", "backend_developer", "devops_engineer"],
  "subagents": {}
}
```

This creates an agent with combined expertise in full-stack, backend, and DevOps.

### Skill + Subagents Pattern

For complex tasks, combine skills with subagent orchestration:

```json
{
  "id": "engineering_manager",
  "name": "Engineering Manager",
  "skills": ["fullstack_developer", "agent_team_workflow"],
  "subagents": {
    "allow_agents": ["backend_dev", "frontend_dev", "qa_eng"],
    "max_spawn_depth": 2
  }
}
```

The manager agent coordinates specialized subagents while maintaining full-stack oversight.

---

## Testing Your Configuration

After adding skills to your `config.json`:

```bash
# Validate configuration
picoclaw-agents agents list

# Test with a query
picoclaw-agents agent -m "Review this API design for security issues"

# Check which skills are loaded
picoclaw-agents skills list
```

---

## Performance Considerations

### Token Usage

Each skill injects ~2,000-4,000 tokens into the system prompt. Consider:

- **Single skill**: ~3K tokens overhead
- **Multiple skills**: Additive (3 skills = ~9K tokens)
- **Impact**: Reduces context window available for conversation

### Recommendations

1. **Use 1-2 skills per agent** for focused expertise
2. **Use orchestrator pattern** for multi-skill needs
3. **Monitor token usage** with `context_management` settings
4. **Consider model context limits** (e.g., 128K for Claude, 32K for GPT-4)

---

## Troubleshooting

### Skill Not Loading

**Symptom:** Agent doesn't behave according to skill

**Check:**
1. Skill name matches exactly (e.g., `backend_developer` not `backend-dev`)
2. Skill is in `skills` array (not `tools_override`)
3. No typos in config.json

### Skills Not Showing in List

**Symptom:** `picoclaw-agents skills list` doesn't show new skills

**Check:**
1. Binary is up to date (rebuild with `make build`)
2. Skills are registered in `loader.go` `listNativeSkills()`

### Agent Ignoring Skill Instructions

**Symptom:** Agent doesn't follow skill guidelines

**Try:**
1. Explicitly mention the role in your query: "As a backend developer, review this API..."
2. Check model temperature (lower = more deterministic)
3. Verify skill content is substantial (check `BuildSkillContext()` output)

---

## See Also

- **[ADDING_NATIVE_SKILLS.md](ADDING_NATIVE_SKILLS.md)**: Guide for developers creating new native skills
- **[config.example.json](config/config.example.json)**: Complete configuration template with examples
- **[CHANGELOG.md](CHANGELOG.md)**: v3.6.0 release notes
