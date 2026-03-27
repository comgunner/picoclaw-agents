# Native Skills - Lista Completa

**Versión:** 3.11.1  
**Última Actualización:** 26 de marzo de 2026  
**Total Skills Nativas:** 14

---

## 📋 Índice

1. [¿Qué son las Skills Nativas?](#qué-son-las-skills-nativas)
2. [Skills vs Tools](#skills-vs-tools)
3. [Lista Completa de Skills Nativas](#lista-completa-de-skills-nativas)
4. [Skills de Rol de Ingeniería](#skills-de-rol-de-ingeniería)
5. [Skills de Propósito General](#skills-de-propósito-general)
6. [Skills de Integración](#skills-de-integración)
7. [Configuración](#configuración)
8. [Ejemplos de Uso](#ejemplos-de-uso)
9. [Referencia](#referencia)

---

## ¿Qué son las Skills Nativas?

Las **Skills Nativas** son definiciones de roles especializados compiladas directamente en el binario de PicoClaw. A diferencia de las Tools (que el LLM llama para realizar acciones), las Skills definen **quién es el agente** — inyectando personalidad, experiencia y guías de comportamiento en el prompt del sistema.

### Características Clave

- ✅ **Cero Dependencias**: Compiladas en el binario, sin archivos externos
- ✅ **Seguridad Mejorada**: No pueden ser modificadas en runtime
- ✅ **Rendimiento**: Sin I/O de archivos, carga instantánea
- ✅ **Type-Safe**: Validación en tiempo de compilación

---

## Skills vs Tools

| Dimensión | Native Skill | Tool |
|-----------|-------------|------|
| **Propósito** | Define el rol/personalidad del agente | Función que el LLM invoca para acciones |
| **Interacción** | El agente "es" este rol | El agente "llama" esta función |
| **Efecto** | Modifica personalidad y comportamiento | Ejecuta acción real (I/O, API, shell) |
| **Interfaz** | `Name()`, `Description()`, `BuildSkillContext()` | `Name()`, `Description()`, `Parameters()`, `Execute()` |
| **Registro** | `listNativeSkills()` en `loader.go` | `ToolRegistry.Register()` en `instance.go` |
| **Ejemplo** | `backend_developer` → agente "es" dev backend | `read_file` → agente "llama" lector de archivos |
| **Output** | String inyectado en system prompt | `*ToolResult` con datos estructurados |

---

## Lista Completa de Skills Nativas

PicoClaw v3.11.1 incluye **14 skills nativas**:

### Skills de Rol de Ingeniería (7 skills)

| # | Skill Name | Descripción | Mejor Para |
|---|------------|-------------|------------|
| 1 | `backend_developer` | Experto en desarrollo backend | APIs REST, bases de datos, microservicios, seguridad |
| 2 | `frontend_developer` | Experto en desarrollo frontend | React, Vue, performance, accesibilidad, design systems |
| 3 | `devops_engineer` | Experto en DevOps | CI/CD, containers, IaC, monitoreo, SRE |
| 4 | `security_engineer` | Experto en seguridad | OWASP, penetration testing, hardening, threat modeling, compliance |
| 5 | `qa_engineer` | Experto en QA | Estrategias de testing, automatización, análisis de cobertura, quality gates |
| 6 | `data_engineer` | Experto en data engineering | ETL pipelines, data warehouses, streaming, data quality |
| 7 | `ml_engineer` | Experto en ML/AI | Entrenamiento de modelos, deployment, MLOps, feature engineering |

### Skills de Propósito General (4 skills)

| # | Skill Name | Descripción | Mejor Para |
|---|------------|-------------|------------|
| 8 | `fullstack_developer` | Desarrollador full-stack experto | Desarrollo general, arquitectura, best practices |
| 9 | `researcher` | Agente de investigación profunda | Búsqueda web, evaluación de fuentes, síntesis de información |
| 10 | `queue_batch` | Delegación de tareas en segundo plano | Tareas pesadas fire-and-forget |
| 11 | `agent_team_workflow` | Orquestador de equipos multi-agente | Coordinación de equipos, delegación de tareas |

### Skills de Integración (3 skills)

| # | Skill Name | Descripción | Mejor Para |
|---|------------|-------------|------------|
| 12 | `binance_mcp` | Integración con Binance MCP | Trading de cripto, datos de mercado |
| 13 | `n8n_workflow` | Experto en automatización n8n | Creación de workflows, validación JSON |
| 14 | `odoo_developer` | Arquitecto Odoo & QA Engineer | Ecosistemas Odoo, migración Pine Script, L10n-Mexico, CFDI 4.0 |

---

## Skills de Rol de Ingeniería

### 1. `backend_developer`

**Descripción:** Experto en desarrollo backend especializado en APIs REST, bases de datos, microservicios y seguridad.

**Responsabilidades Core:**
- Diseño e implementación de APIs REST y GraphQL
- Modelado de bases de datos relacionales y NoSQL
- Implementación de autenticación/autorización
- Optimización de queries y performance
- Arquitectura de microservicios

**Tecnologías:**
- **Lenguajes:** Go, Python, Node.js, Java
- **Bases de Datos:** PostgreSQL, MySQL, MongoDB, Redis
- **APIs:** REST, GraphQL, gRPC
- **Message Brokers:** Kafka, RabbitMQ, NATS

**Cuándo Usar:**
- ✅ Diseñando schemas de API
- ✅ Modelando bases de datos
- ✅ Implementando autenticación JWT/OAuth
- ✅ Optimizando queries lentos
- ✅ Construyendo microservicios

---

### 2. `frontend_developer`

**Descripción:** Experto en desarrollo frontend especializado en frameworks modernos, performance y accesibilidad.

**Responsabilidades Core:**
- Creación de componentes React/Vue/Svelte
- Implementación de state management
- Optimización de Core Web Vitals
- Implementación de WCAG 2.1 AA
- Diseño de design systems

**Tecnologías:**
- **Frameworks:** React, Vue, Svelte, Next.js, Nuxt
- **State:** Redux, Zustand, Pinia, Signals
- **Styling:** Tailwind CSS, CSS Modules, Styled Components
- **Testing:** Jest, Vitest, React Testing Library, Cypress

**Cuándo Usar:**
- ✅ Creando componentes UI
- ✅ Implementando routing
- ✅ Optimizando performance (LCP, FID, CLS)
- ✅ Asegurando accesibilidad
- ✅ Construyendo layouts responsive

---

### 3. `devops_engineer`

**Descripción:** Experto en DevOps especializado en CI/CD, containers, infraestructura como código y SRE.

**Responsabilidades Core:**
- Diseño de pipelines CI/CD
- Creación de Kubernetes manifests
- Writing Terraform modules
- Configuración de monitoreo/alerting
- Diseño de disaster recovery

**Tecnologías:**
- **CI/CD:** GitHub Actions, GitLab CI, Jenkins, ArgoCD
- **Containers:** Docker, Kubernetes, Helm
- **IaC:** Terraform, Pulumi, Ansible
- **Monitoreo:** Prometheus, Grafana, Datadog, New Relic

**Cuándo Usar:**
- ✅ Creando pipelines de deployment
- ✅ Escribiendo Kubernetes manifests
- ✅ Configurando Terraform
- ✅ Implementando monitoreo
- ✅ Diseñando estrategias de backup

---

### 4. `security_engineer`

**Descripción:** Experto en seguridad especializado en OWASP, penetration testing, hardening y compliance.

**Responsabilidades Core:**
- Threat modeling
- Security code reviews
- Implementación de controles OWASP
- Hardening de sistemas
- Compliance (SOC2, GDPR, HIPAA)

**Tecnologías:**
- **SAST/DAST:** SonarQube, Snyk, Dependabot
- **Scanning:** Trivy, Clair, Anchore
- **Secrets:** Vault, AWS Secrets Manager
- **Compliance:** SOC2, ISO 27001, GDPR

**Cuándo Usar:**
- ✅ Realizando threat modeling
- ✅ Auditando código por vulnerabilidades
- ✅ Implementando autenticación segura
- ✅ Revisando configuración de infraestructura
- ✅ Asegurando compliance

---

### 5. `qa_engineer`

**Descripción:** Experto en QA especializado en estrategias de testing, automatización y quality gates.

**Responsabilidades Core:**
- Diseño de estrategia de testing
- Writing unit/integration/E2E tests
- Configuración de test automation
- Análisis de code coverage
- Implementación de quality gates

**Tecnologías:**
- **Unit Testing:** Jest, Vitest, pytest, JUnit
- **E2E:** Cypress, Playwright, Selenium
- **API Testing:** Postman, REST Assured
- **Coverage:** Istanbul, coverage.py, JaCoCo

**Cuándo Usar:**
- ✅ Diseñando estrategia de tests
- ✅ Escribiendo tests automatizados
- ✅ Configurando CI con tests
- ✅ Analizando cobertura de código
- ✅ Implementando quality gates

---

### 6. `data_engineer`

**Descripción:** Experto en data engineering especializado en ETL pipelines, data warehouses y streaming.

**Responsabilidades Core:**
- Construcción de ETL/ELT pipelines
- Modelado de data warehouses
- Implementación de streaming pipelines
- Aseguramiento de data quality
- Configuración de data governance

**Tecnologías:**
- **Processing:** Spark, Flink, dbt
- **Warehouses:** Snowflake, BigQuery, Redshift
- **Streaming:** Kafka, Kinesis, Pulsar
- **Orchestration:** Airflow, Dagster, Prefect

**Cuándo Usar:**
- ✅ Construyendo pipelines de datos
- ✅ Modelando data warehouses
- ✅ Implementando streaming en tiempo real
- ✅ Asegurando calidad de datos
- ✅ Configurando gobernanza de datos

---

### 7. `ml_engineer`

**Descripción:** Experto en ML/AI especializado en entrenamiento, deployment y MLOps.

**Responsabilidades Core:**
- Entrenamiento de modelos ML
- Deployment de modelos a producción
- Configuración de pipelines MLOps
- Feature engineering
- Monitoreo de model drift

**Tecnologías:**
- **Frameworks:** PyTorch, TensorFlow, scikit-learn
- **Deployment:** SageMaker, Vertex AI, Azure ML
- **MLOps:** MLflow, Kubeflow, Weights & Biases
- **Monitoring:** Evidently AI, Arize, WhyLabs

**Cuándo Usar:**
- ✅ Entrenando modelos ML
- ✅ Deployando modelos a producción
- ✅ Configurando pipelines de retraining
- ✅ Implementando feature stores
- ✅ Monitoreando model drift

---

## Skills de Propósito General

### 8. `fullstack_developer`

**Descripción:** Desarrollador full-stack experto con conocimiento en frontend, backend y mejores prácticas.

**Cuándo Usar:**
- ✅ Desarrollo general de features
- ✅ Revisiones de arquitectura
- ✅ Implementación de best practices
- ✅ Refactorización de código
- ✅ Documentación técnica

---

### 9. `researcher`

**Descripción:** Agente de investigación profunda especializado en búsqueda web, evaluación de fuentes y síntesis.

**Capacidades:**
- Búsqueda web avanzada
- Evaluación crítica de fuentes
- Síntesis de información
- Reportes estructurados

**Cuándo Usar:**
- ✅ Investigando temas complejos
- ✅ Evaluando múltiples fuentes
- ✅ Sintetizando información
- ✅ Creando reportes de investigación

---

### 10. `queue_batch`

**Descripción:** Sistema de delegación de tareas en segundo plano usando patrón fire-and-forget.

**Características:**
- Procesamiento asíncrono
- Cola de tareas persistente
- Reintentos automáticos
- Monitoreo de estado

**Cuándo Usar:**
- ✅ Tareas pesadas de background
- ✅ Procesamiento por lotes
- ✅ Operaciones que no bloquean
- ✅ Reintentos automáticos

---

### 11. `agent_team_workflow`

**Descripción:** Orquestador de equipos multi-agente para coordinar tareas complejas.

**Capacidades:**
- Análisis de tareas
- Selección óptima de agentes
- Coordinación de ejecución
- Síntesis de resultados

**Cuándo Usar:**
- ✅ Tareas complejas multi-etapa
- ✅ Coordinación de especialistas
- ✅ Orquestación de flujos
- ✅ Gestión de dependencias

---

## Skills de Integración

### 12. `binance_mcp`

**Descripción:** Integración con servidor MCP de Binance para trading y datos de mercado.

**Capacidades:**
- Consultar balances spot/futures
- Obtener precios de ticker
- Ejecutar órdenes de trading
- Analizar order books

**Cuándo Usar:**
- ✅ Trading de criptomonedas
- ✅ Consulta de balances
- ✅ Análisis de mercado
- ✅ Ejecución de órdenes

---

### 13. `n8n_workflow`

**Descripción:** Experto en automatización n8n para creación de workflows production-ready.

**Capacidades:**
- Diseño de workflows n8n
- Validación de JSON
- Integración de nodos
- Best practices de automatización

**Cuándo Usar:**
- ✅ Creando workflows n8n
- ✅ Validando configuraciones
- ✅ Integrando APIs
- ✅ Automatizando procesos

---

### 14. `odoo_developer`

**Descripción:** Arquitecto Principal Odoo & QA Engineer especializado en ecosistemas Odoo.

**Capacidades:**
- Desarrollo de módulos Odoo
- Migración de Pine Script
- Localización L10n-Mexico
- CFDI 4.0 y facturación electrónica

**Cuándo Usar:**
- ✅ Desarrollo Odoo
- ✅ Migración de sistemas legacy
- ✅ Implementación L10n-Mexico
- ✅ Integración CFDI 4.0

---

## Configuración

### Skill Único Especializado

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

### Orquestador con Subagentes Especializados

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

## Ejemplos de Uso

### Ejemplo 1: Equipo de Desarrollo Full-Stack

**Configuración:**
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

**Uso:**
```bash
picoclaw-agents agent -m "Build a complete user authentication system with login, registration, and password recovery"
```

**Flujo:**
1. **Product Team** (orchestrator) analiza la tarea
2. Spawnea **Backend Dev** para APIs de autenticación
3. Spawnea **Frontend Dev** para formularios UI
4. Spawnea **QA Eng** para tests de seguridad
5. Sintetiza resultados en solución completa

---

### Ejemplo 2: Pipeline de Datos ML

**Configuración:**
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

**Uso:**
```bash
picoclaw-agents agent -m "Build an end-to-end ML pipeline for customer churn prediction"
```

**Flujo:**
1. **ML Engineer** diseña arquitectura del modelo
2. **Data Engineer** construye pipeline ETL
3. **Backend Dev** crea API de predicción
4. **ML Engineer** entrena y deploya modelo

---

### Ejemplo 3: Auditoría de Seguridad

**Configuración:**
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

**Uso:**
```bash
picoclaw-agents agent -m "Conduct a comprehensive security audit of our authentication system"
```

**Flujo:**
1. **Security Engineer** realiza threat modeling
2. **Backend Dev** revisa código de autenticación
3. **DevOps Eng** audita configuración de infraestructura
4. **Security Engineer** sintetiza hallazgos y recomendaciones

---

## Combinando Skills

### Multi-Skill Agents

Puedes asignar múltiples skills a un solo agente:

```json
{
  "id": "fullstack_tech_lead",
  "name": "Full-Stack Tech Lead",
  "model": "deepseek-chat",
  "skills": ["fullstack_developer", "backend_developer", "devops_engineer"],
  "subagents": {}
}
```

Esto crea un agente con experiencia combinada en full-stack, backend y DevOps.

### Patrón Skill + Subagentes

Para tareas complejas, combina skills con orquestación de subagentes:

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

El agente manager coordina subagentes especializados mientras mantiene oversight full-stack.

---

## Consideraciones de Rendimiento

### Uso de Tokens

Cada skill inyecta ~2,000-4,000 tokens en el system prompt. Considera:

- **Single skill:** ~3K tokens overhead
- **Múltiple skills:** Aditivo (3 skills = ~9K tokens)
- **Impacto:** Reduce ventana de contexto disponible para conversación

### Recomendaciones

1. **Usa 1-2 skills por agente** para experiencia enfocada
2. **Usa patrón orquestador** para necesidades multi-skill
3. **Monitorea uso de tokens** con configuración `context_management`
4. **Considera límites de contexto del modelo** (ej. 128K para Claude, 32K para GPT-4)

---

## Testing de Configuración

Después de agregar skills a tu `config.json`:

```bash
# Validar configuración
picoclaw-agents agents list

# Test con query
picoclaw-agents agent -m "Review this API design for security issues"

# Verificar skills cargadas
picoclaw-agents skills list
```

---

## Troubleshooting

### Skill No Carga

**Síntoma:** El agente no se comporta según la skill

**Verificar:**
1. Nombre de skill coincide exactamente (ej. `backend_developer` no `backend-dev`)
2. Skill está en array `skills` (no en `tools_override`)
3. No hay typos en config.json

### Skills No Aparecen en List

**Síntoma:** `picoclaw-agents skills list` no muestra skills nuevas

**Verificar:**
1. Binario está actualizado (rebuild con `make build`)
2. Skills están registradas en `loader.go` `listNativeSkills()`

### Agente Ignora Instrucciones de Skill

**Síntoma:** El agente no sigue guías de la skill

**Intentar:**
1. Mencionar explícitamente el rol en tu query: "As a backend developer, review this API..."
2. Verificar temperatura del modelo (más baja = más determinístico)
3. Verificar contenido de la skill es substancial (revisar output de `BuildSkillContext()`)

---

## Referencia

### Archivos Relacionados

- **[SKILLS.md](SKILLS.md)**: Guía completa de skills nativas
- **[ADDING_NATIVE_SKILLS.md](ADDING_NATIVE_SKILLS.md)**: Guía para desarrolladores creando nuevas skills
- **[config.example.json](config/config.example.json)**: Template de configuración completo
- **[CHANGELOG.md](CHANGELOG.md)**: Notas de release v3.11.1

### Enlaces Externos

- **Documentación Oficial MCP**: https://modelcontextprotocol.io
- **SDK de TypeScript**: https://github.com/modelcontextprotocol/typescript-sdk
- **Ejemplos de Servidores**: https://github.com/modelcontextprotocol/servers

---

**Última actualización:** 26 de marzo de 2026  
**Mantenido por:** @comgunner  
**Versión:** 3.11.1
