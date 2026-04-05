# Native Skills - Complete List

**Version:** 3.11.1  
**Last Updated:** March 26, 2026  
**Total Native Skills:** 14

---

## 📋 Table of Contents

1. [What are Native Skills?](#what-are-native-skills)
2. [Skills vs Tools](#skills-vs-tools)
3. [Complete List of Native Skills](#complete-list-of-native-skills)
4. [Engineering Role Skills](#engineering-role-skills)
5. [General Purpose Skills](#general-purpose-skills)
6. [Integration Skills](#integration-skills)
7. [Configuration](#configuration)
8. [Usage Examples](#usage-examples)
9. [Reference](#reference)

---

## What are Native Skills?

**Native Skills** are specialized role definitions compiled directly into the PicoClaw binary. Unlike Tools (which the LLM calls to perform actions), Skills define **who the agent is** — injecting personality, expertise, and behavioral guidelines into the system prompt.

### Key Features

- ✅ **Zero Dependencies**: Compiled into binary, no external files
- ✅ **Enhanced Security**: Cannot be modified at runtime
- ✅ **Performance**: No file I/O, instant loading
- ✅ **Type-Safe**: Compile-time validation

---

## Skills vs Tools

| Dimension | Native Skill | Tool |
|-----------|-------------|------|
| **Purpose** | Defines agent role/personality | Function LLM invokes for actions |
| **Interaction** | Agent "is" this role | Agent "calls" this function |
| **Effect** | Modifies personality and behavior | Executes real action (I/O, API, shell) |
| **Interface** | `Name()`, `Description()`, `BuildSkillContext()` | `Name()`, `Description()`, `Parameters()`, `Execute()` |
| **Registration** | `listNativeSkills()` in `loader.go` | `ToolRegistry.Register()` in `instance.go` |
| **Example** | `backend_developer` → agent "is" backend dev | `read_file` → agent "calls" file reader |
| **Output** | String injected into system prompt | `*ToolResult` with structured data |

---

## Complete List of Native Skills

PicoClaw v3.11.1 includes **14 native skills**:

### Engineering Role Skills (7 skills)

| # | Skill Name | Description | Best For |
|---|------------|-------------|----------|
| 1 | `backend_developer` | Backend development expert | REST APIs, databases, microservices, security |
| 2 | `frontend_developer` | Frontend development expert | React, Vue, performance, accessibility, design systems |
| 3 | `devops_engineer` | DevOps expert | CI/CD, containers, IaC, monitoring, SRE |
| 4 | `security_engineer` | Security expert | OWASP, penetration testing, hardening, threat modeling, compliance |
| 5 | `qa_engineer` | QA expert | Testing strategies, automation, coverage analysis, quality gates |
| 6 | `data_engineer` | Data engineering expert | ETL pipelines, data warehouses, streaming, data quality |
| 7 | `ml_engineer` | ML/AI expert | Model training, deployment, MLOps, feature engineering |

### General Purpose Skills (4 skills)

| # | Skill Name | Description | Best For |
|---|------------|-------------|----------|
| 8 | `fullstack_developer` | Full-stack development expert | General development, architecture, best practices |
| 9 | `researcher` | Deep research agent | Web search, source evaluation, information synthesis |
| 10 | `queue_batch` | Background task delegation | Heavy fire-and-forget tasks |
| 11 | `agent_team_workflow` | Multi-agent team orchestrator | Team coordination, task delegation |

### Integration Skills (3 skills)

| # | Skill Name | Description | Best For |
|---|------------|-------------|----------|
| 12 | `binance_mcp` | Binance MCP integration | Crypto trading, market data |
| 13 | `n8n_workflow` | n8n automation expert | Workflow creation, JSON validation |
| 14 | `odoo_developer` | Odoo architect & QA engineer | Odoo ecosystems, Pine Script migration, L10n-Mexico, CFDI 4.0 |

---

## Engineering Role Skills

### 1. `backend_developer`

**Description:** Backend development expert specialized in REST APIs, databases, microservices, and security.

**Core Responsibilities:**
- REST and GraphQL API design and implementation
- Relational and NoSQL database modeling
- Authentication/authorization implementation
- Query and performance optimization
- Microservices architecture

**Technologies:**
- **Languages:** Go, Python, Node.js, Java
- **Databases:** PostgreSQL, MySQL, MongoDB, Redis
- **APIs:** REST, GraphQL, gRPC
- **Message Brokers:** Kafka, RabbitMQ, NATS

**When to Use:**
- ✅ Designing API schemas
- ✅ Modeling databases
- ✅ Implementing JWT/OAuth authentication
- ✅ Optimizing slow queries
- ✅ Building microservices

---

### 2. `frontend_developer`

**Description:** Frontend development expert specialized in modern frameworks, performance, and accessibility.

**Core Responsibilities:**
- Creating React/Vue/Svelte components
- Implementing state management
- Core Web Vitals optimization
- WCAG 2.1 AA implementation
- Design system development

**Technologies:**
- **Frameworks:** React, Vue, Svelte, Next.js, Nuxt
- **State:** Redux, Zustand, Pinia, Signals
- **Styling:** Tailwind CSS, CSS Modules, Styled Components
- **Testing:** Jest, Vitest, React Testing Library, Cypress

**When to Use:**
- ✅ Creating UI components
- ✅ Implementing routing
- ✅ Optimizing performance (LCP, FID, CLS)
- ✅ Ensuring accessibility
- ✅ Building responsive layouts

---

### 3. `devops_engineer`

**Description:** DevOps expert specialized in CI/CD, containers, infrastructure as code, and SRE.

**Core Responsibilities:**
- CI/CD pipeline design
- Kubernetes manifests creation
- Terraform modules writing
- Monitoring/alerting configuration
- Disaster recovery design

**Technologies:**
- **CI/CD:** GitHub Actions, GitLab CI, Jenkins, ArgoCD
- **Containers:** Docker, Kubernetes, Helm
- **IaC:** Terraform, Pulumi, Ansible
- **Monitoring:** Prometheus, Grafana, Datadog, New Relic

**When to Use:**
- ✅ Creating deployment pipelines
- ✅ Writing Kubernetes manifests
- ✅ Configuring Terraform
- ✅ Implementing monitoring
- ✅ Designing backup strategies

---

### 4. `security_engineer`

**Description:** Security expert specialized in OWASP, penetration testing, hardening, and compliance.

**Core Responsibilities:**
- Threat modeling
- Security code reviews
- OWASP controls implementation
- Systems hardening
- Compliance (SOC2, GDPR, HIPAA)

**Technologies:**
- **SAST/DAST:** SonarQube, Snyk, Dependabot
- **Scanning:** Trivy, Clair, Anchore
- **Secrets:** Vault, AWS Secrets Manager
- **Compliance:** SOC2, ISO 27001, GDPR

**When to Use:**
- ✅ Performing threat modeling
- ✅ Auditing code for vulnerabilities
- ✅ Implementing secure authentication
- ✅ Reviewing infrastructure configuration
- ✅ Ensuring compliance

---

### 5. `qa_engineer`

**Description:** QA expert specialized in testing strategies, automation, and quality gates.

**Core Responsibilities:**
- Testing strategy design
- Writing unit/integration/E2E tests
- Test automation configuration
- Code coverage analysis
- Quality gates implementation

**Technologies:**
- **Unit Testing:** Jest, Vitest, pytest, JUnit
- **E2E:** Cypress, Playwright, Selenium
- **API Testing:** Postman, REST Assured
- **Coverage:** Istanbul, coverage.py, JaCoCo

**When to Use:**
- ✅ Designing testing strategy
- ✅ Writing automated tests
- ✅ Configuring CI with tests
- ✅ Analyzing code coverage
- ✅ Implementing quality gates

---

### 6. `data_engineer`

**Description:** Data engineering expert specialized in ETL pipelines, data warehouses, and streaming.

**Core Responsibilities:**
- ETL/ELT pipeline construction
- Data warehouse modeling
- Streaming pipeline implementation
- Data quality assurance
- Data governance configuration

**Technologies:**
- **Processing:** Spark, Flink, dbt
- **Warehouses:** Snowflake, BigQuery, Redshift
- **Streaming:** Kafka, Kinesis, Pulsar
- **Orchestration:** Airflow, Dagster, Prefect

**When to Use:**
- ✅ Building data pipelines
- ✅ Modeling data warehouses
- ✅ Implementing real-time streaming
- ✅ Ensuring data quality
- ✅ Configuring data governance

---

### 7. `ml_engineer`

**Description:** ML/AI expert specialized in training, deployment, and MLOps.

**Core Responsibilities:**
- ML model training
- Model deployment to production
- MLOps pipelines configuration
- Feature engineering
- Model drift monitoring

**Technologies:**
- **Frameworks:** PyTorch, TensorFlow, scikit-learn
- **Deployment:** SageMaker, Vertex AI, Azure ML
- **MLOps:** MLflow, Kubeflow, Weights & Biases
- **Monitoring:** Evidently AI, Arize, WhyLabs

**When to Use:**
- ✅ Training ML models
- ✅ Deploying models to production
- ✅ Configuring retraining pipelines
- ✅ Implementing feature stores
- ✅ Monitoring model drift

---

## General Purpose Skills

### 8. `fullstack_developer`

**Description:** Full-stack development expert with knowledge in frontend, backend, and best practices.

**When to Use:**
- ✅ General feature development
- ✅ Architecture reviews
- ✅ Best practices implementation
- ✅ Code refactoring
- ✅ Technical documentation

---

### 9. `researcher`

**Description:** Deep research agent specialized in web search, source evaluation, and synthesis.

**Capabilities:**
- Advanced web search
- Critical source evaluation
- Information synthesis
- Structured reporting

**When to Use:**
- ✅ Researching complex topics
- ✅ Evaluating multiple sources
- ✅ Synthesizing information
- ✅ Creating research reports

---

### 10. `queue_batch`

**Description:** Background task delegation system using fire-and-forget pattern.

**Capabilities:**
- Asynchronous processing
- Persistent task queue
- Automatic retries
- Status monitoring

**When to Use:**
- ✅ Heavy background tasks
- ✅ Batch processing
- ✅ Non-blocking operations
- ✅ Automatic retries

---

### 11. `agent_team_workflow`

**Description:** Multi-agent team orchestrator for coordinating complex tasks.

**Capabilities:**
- Task analysis
- Optimal agent selection
- Execution coordination
- Results synthesis

**When to Use:**
- ✅ Complex multi-stage tasks
- ✅ Specialist coordination
- ✅ Workflow orchestration
- ✅ Dependency management

---

## Integration Skills

### 12. `binance_mcp`

**Description:** Integration with Binance MCP server for trading and market data.

**Capabilities:**
- Query spot/futures balances
- Get ticker prices
- Execute trading orders
- Analyze order books

**When to Use:**
- ✅ Cryptocurrency trading
- ✅ Balance queries
- ✅ Market analysis
- ✅ Order execution

---

### 13. `n8n_workflow`

**Description:** n8n automation expert for creating production-ready workflows.

**Capabilities:**
- n8n workflow design
- JSON validation
- Node integration
- Automation best practices

**When to Use:**
- ✅ Creating n8n workflows
- ✅ Validating configurations
- ✅ Integrating APIs
- ✅ Automating processes

---

### 14. `odoo_developer`

**Description:** Principal Odoo Architect & QA Engineer specialized in Odoo ecosystems.

**Capabilities:**
- Odoo module development
- Pine Script migration
- L10n-Mexico localization
- CFDI 4.0 and electronic invoicing

**When to Use:**
- ✅ Odoo development
- ✅ Legacy system migration
- ✅ L10n-Mexico implementation
- ✅ CFDI 4.0 integration

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

## Usage Examples

### Example 1: Full-Stack Development Team

**Configuration:**
```json
{
  "agents": {
    "list": [
      {
        "id": "product_team",
        "name": "Product Development Team",
        "skills": ["fullstack_developer", "agent_team_workflow"],
        "subagents": {
          "allow_agents": ["backend_dev", "frontend_dev", "qa_eng"],
          "max_spawn_depth": 2
        }
      }
    ]
  }
}
```

**Usage:**
```bash
picoclaw-agents agent -m "Build a complete user authentication system with login, registration, and password recovery"
```

**Flow:**
1. **Product Team** (orchestrator) analyzes the task
2. Spawns **Backend Dev** for authentication APIs
3. Spawns **Frontend Dev** for forms UI
4. Spawns **QA Eng** for security tests
5. Synthesizes results into complete solution

---

### Example 2: ML Data Pipeline

**Configuration:**
```json
{
  "agents": {
    "list": [
      {
        "id": "ml_pipeline",
        "name": "ML Pipeline Team",
        "skills": ["ml_engineer", "agent_team_workflow"],
        "subagents": {
          "allow_agents": ["data_eng", "backend_dev"],
          "max_spawn_depth": 2
        }
      }
    ]
  }
}
```

**Usage:**
```bash
picoclaw-agents agent -m "Build an end-to-end ML pipeline for customer churn prediction"
```

**Flow:**
1. **ML Engineer** designs model architecture
2. **Data Engineer** builds ETL pipeline
3. **Backend Dev** creates prediction API
4. **ML Engineer** trains and deploys model

---

### Example 3: Security Audit

**Configuration:**
```json
{
  "agents": {
    "list": [
      {
        "id": "security_audit",
        "name": "Security Audit Team",
        "skills": ["security_engineer", "agent_team_workflow"],
        "subagents": {
          "allow_agents": ["backend_dev", "devops_eng"],
          "max_spawn_depth": 2
        }
      }
    ]
  }
}
```

**Usage:**
```bash
picoclaw-agents agent -m "Conduct a comprehensive security audit of our authentication system"
```

**Flow:**
1. **Security Engineer** performs threat modeling
2. **Backend Dev** reviews authentication code
3. **DevOps Eng** audits infrastructure configuration
4. **Security Engineer** synthesizes findings and recommendations

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

## Performance Considerations

### Token Usage

Each skill injects ~2,000-4,000 tokens into the system prompt. Consider:

- **Single skill:** ~3K tokens overhead
- **Multiple skills:** Additive (3 skills = ~9K tokens)
- **Impact:** Reduces context window available for conversation

### Recommendations

1. **Use 1-2 skills per agent** for focused expertise
2. **Use orchestrator pattern** for multi-skill needs
3. **Monitor token usage** with `context_management` configuration
4. **Consider model context limits** (e.g., 128K for Claude, 32K for GPT-4)

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

## Reference

### Related Files

- **[SKILLS.md](SKILLS.md)**: Complete native skills guide
- **[ADDING_NATIVE_SKILLS.md](ADDING_NATIVE_SKILLS.md)**: Developer guide for creating new skills
- **[config.example.json](../config/config.example.json)**: Complete configuration template
- **[CHANGELOG.md](../CHANGELOG.md)**: v3.11.1 release notes

### External Links

- **Official MCP Documentation:** https://modelcontextprotocol.io
- **TypeScript SDK:** https://github.com/modelcontextprotocol/typescript-sdk
- **Server Examples:** https://github.com/modelcontextprotocol/servers

---

**Last updated:** March 26, 2026  
**Maintained by:** @comgunner  
**Version:** 3.11.1
