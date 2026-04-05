# Guía de Skills Nativas

**Versión:** 3.11.1  
**Última Actualización:** 26 de marzo de 2026  
**Total Skills Nativas:** 14

---

## 📋 Índice

1. [Descripción General](#descripción-general)
2. [Skills vs Tools](#skills-vs-tools)
3. [Skills Disponibles](#skills-disponibles)
4. [Configuración](#configuración)
5. [Combinando Skills](#combinando-skills)
6. [Consideraciones de Rendimiento](#consideraciones-de-rendimiento)
7. [Testing](#testing)
8. [Troubleshooting](#troubleshooting)

---

## Descripción General

Las **Skills Nativas** son definiciones de roles especializados compiladas directamente en el binario de PicoClaw. A diferencia de las Tools (que el LLM llama para realizar acciones), las Skills definen **quién es el agente** — inyectando personalidad, experiencia y guías de comportamiento en el prompt del sistema.

---

## Skills vs Tools

### Diferencias Clave

| Dimensión | Native Skill | Tool |
|-----------|-------------|------|
| **Qué es** | Instrucciones de rol inyectadas en system prompt | Función que el LLM invoca vía `tool_use` |
| **Cómo interactúa LLM** | "Es" este rol / "tiene" esta experiencia | "Llama" esta función para realizar acción |
| **Efecto** | Modifica personalidad y comportamiento del agente | Ejecuta acción real (I/O, API, shell, etc.) |
| **Interfaz** | `Name()`, `Description()`, `BuildSkillContext()` | `Name()`, `Description()`, `Parameters()`, `Execute()` |
| **Registro** | `listNativeSkills()` en `loader.go` | `ToolRegistry.Register()` en `instance.go` |
| **Ejemplo** | `backend_developer` → agente "es" dev backend | `read_file` → agente "llama" lector de archivos |
| **Output** | String inyectado en system prompt | `*ToolResult` con datos estructurados |
| **Dependencias runtime** | Ninguna — compilado en binario | Puede requerir binarios externos (git, docker, etc.) |

---

## Comparación de Flujo

### Flujo de Skills
```
config.json → agent.skills: ["backend_developer"]
  → SkillsLoader.LoadSkill("backend_developer")
    → listNativeSkills() — encuentra skill compilada
      → LoadNativeBackendDeveloperSkill()
        → BackendDeveloperSkill.BuildSkillContext()
          → string inyectado en system prompt del LLM
```

### Flujo de Tools
```
config.json → agent.tools_override: ["read_file", "exec"]
  → ToolRegistry.Register(FilesystemTool{})
    → ToolRegistry.ToProviderDefs()
      → []providers.ToolDefinition enviados al LLM
        → LLM genera: tool_use { name: "read_file", args: {...} }
          → ToolRegistry.Execute("read_file", ctx, args)
            → *ToolResult retornado al LLM
```

**Regla General:** Las skills de rol de ingeniería (backend_developer, devops_engineer, etc.) son **Native Skills**. El LLM no las "llama" — las "tiene" como parte de su identidad.

---

## Skills Disponibles (v3.11.1)

PicoClaw v3.11.1 incluye **14 skills nativas**:

### Skills de Rol de Ingeniería (7 skills)

| Skill Name | Propósito | Mejor Para |
|------------|-----------|------------|
| `backend_developer` | Experto en desarrollo backend | APIs REST, bases de datos, microservicios, seguridad |
| `frontend_developer` | Experto en desarrollo frontend | React, Vue, performance, accesibilidad |
| `devops_engineer` | Experto en DevOps | CI/CD, Kubernetes, Terraform, monitoreo |
| `security_engineer` | Experto en seguridad | OWASP, penetration testing, threat modeling |
| `qa_engineer` | Experto en QA | Test automation, coverage analysis, E2E testing |
| `data_engineer` | Experto en data engineering | ETL pipelines, data warehouses, streaming |
| `ml_engineer` | Experto en ML/AI | Entrenamiento de modelos, deployment, MLOps |

### Skills de Propósito General (4 skills)

| Skill Name | Propósito | Mejor Para |
|------------|-----------|------------|
| `fullstack_developer` | Asistente de desarrollo full-stack | Coding general, arquitectura, best practices |
| `researcher` | Agente de investigación profunda | Búsqueda web, evaluación de fuentes, síntesis |
| `queue_batch` | Delegación de tareas en segundo plano | Tareas pesadas fire-and-forget |
| `agent_team_workflow` | Orquestador multi-agente | Coordinación de equipos, delegación de tareas |

### Skills de Integración (3 skills)

| Skill Name | Propósito | Mejor Para |
|------------|-----------|------------|
| `binance_mcp` | Integración con Binance | Trading de cripto, datos de mercado |
| `n8n_workflow` | Experto en automatización n8n | Creación de workflows, validación JSON |
| `odoo_developer` | Arquitecto Odoo & QA engineer | Ecosistemas Odoo, L10n-Mexico, CFDI 4.0 |

---

## Configuración

### Agente Especializado Único

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

## Combinando Skills

### Agentes Multi-Skill

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

El agente manager coordina subagentes especializados mientras mantiene supervisión full-stack.

---

## Consideraciones de Rendimiento

### Uso de Tokens

Cada skill inyecta ~2,000-4,000 tokens en el system prompt. Considera:

- **Single skill:** ~3K tokens de overhead
- **Múltiple skills:** Aditivo (3 skills = ~9K tokens)
- **Impacto:** Reduce la ventana de contexto disponible para conversación

### Recomendaciones

1. **Usa 1-2 skills por agente** para experiencia enfocada
2. **Usa patrón orquestador** para necesidades multi-skill
3. **Monitorea uso de tokens** con configuración `context_management`
4. **Considera límites de contexto del modelo** (ej. 128K para Claude, 32K para GPT-4)

---

## Testing

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

## Ver También

- **[NATIVE_SKILLS_LIST.es.md](NATIVE_SKILLS_LIST.es.md)**: Lista completa y detallada de todas las skills
- **[ADDING_NATIVE_SKILLS.es.md](ADDING_NATIVE_SKILLS.es.md)**: Guía para desarrolladores creando nuevas skills
- **[config.example.json](../config/config.example.json)**: Template de configuración completo con ejemplos
- **[CHANGELOG.md](../CHANGELOG.md)**: Notas de release v3.11.1

---

**Última actualización:** 26 de marzo de 2026  
**Mantenido por:** @comgunner  
**Versión:** 3.11.1
