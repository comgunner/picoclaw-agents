# Tutorial: Skill Creator

**ID del Skill:** `skill_creator`  
**Tipo:** Nativo (compilado en el binario)  
**Versión:** 1.0.0  
**Última actualización:** 12 de abril de 2026

---

## Resumen

El **Skill Creator** (`skill_creator`) es un skill nativo integrado en PicoClaw-Agents que guía al agente en la creación de nuevos **skills basados en archivos** (SKILL.md) en tu espacio de trabajo. Sigue un flujo de trabajo estructurado de 6 pasos para garantizar que cada skill que crees esté bien diseñado, sea seguro y cumpla con las convenciones de PicoClaw.

### Qué Hace

El Skill Creator enseña al agente cómo:

1. Entender qué skill necesita el usuario
2. Planificar la estructura y el contenido del skill
3. Crear el directorio y los archivos en `~/.picoclaw/workspace/skills/`
4. Escribir un SKILL.md válido con frontmatter YAML apropiado
5. Validar el skill contra reglas de nomenclatura y seguridad
6. Iterar basándose en el uso real

### Qué NO Hace

- ❌ **NO** crea skills nativos en Go (esas son tareas exclusivas de desarrolladores que requieren cambios de código y recompilación)
- ❌ **NO** instala skills desde ClawHub (usa la herramienta `install_skill` para eso)
- ❌ **NO** modifica la configuración del agente

---

## Inicio Rápido

Simplemente pide al agente que cree un skill:

```
Crea un skill para rotar páginas de PDF usando Python.
```

El agente invocará automáticamente el flujo de trabajo del Skill Creator y te guiará a través del proceso.

---

## El Flujo de Trabajo de 6 Pasos

### Paso 1: Entender

El agente hace preguntas concretas para entender el propósito del skill:

- "¿Qué problema resuelve este skill?"
- "¿Puedes dar un ejemplo de cómo se usaría?"
- "¿Necesita scripts (Python/Bash) o solo documentación?"

**Ejemplo de interacción:**

```
Usuario: "Necesito un skill para hacer copias de seguridad de mi base de datos."

Agente: "¡Claro! Déjame entender:
  1. ¿Qué sistema de base de datos? (PostgreSQL, MySQL, etc.)
  2. ¿Debe ser una copia completa o incremental?
  3. ¿Alguna política de retención para copias antiguas?"
```

### Paso 2: Planificar

El agente analiza el caso de uso y determina la estructura del skill:

| Necesidad | Solución |
|------|----------|
| Llamadas a API o procesamiento de datos | Scripts Python en `scripts/` |
| Automatización o comandos de shell | Scripts Bash en `scripts/` |
| Material de referencia, esquemas | Documentación en `references/` |
| Plantillas, archivos base | Recursos en `assets/` |

**Árbol de decisión:**

```
¿El skill necesita código ejecutable?
├─ SÍ → Usa scripts/ (Python preferido, Bash para tareas simples)
└─ NO  → Skill solo de documentación (solo SKILL.md)
```

### Paso 3: Crear

El agente crea la estructura de directorios:

```
~/.picoclaw/workspace/skills/{nombre-skill}/
├── SKILL.md          ← Obligatorio: instrucciones + frontmatter YAML
├── scripts/          ← Opcional: scripts Python, Bash
├── references/       ← Opcional: docs detallados, esquemas
└── assets/           ← Opcional: plantillas, archivos base
```

### Paso 4: Escribir

#### Plantilla SKILL.md

Cada skill comienza con un archivo SKILL.md que contiene frontmatter YAML e instrucciones en markdown:

```yaml
---
name: nombre-mi-skill
description: Descripción clara con desencadenantes de uso. Usar cuando X, Y o Z.
---
```

```markdown
# Título del Skill

Breve descripción de lo que hace.

## Cuándo Usar
- Caso 1: cuando el usuario pide...
- Caso 2: cuando necesites...

## Uso

```bash
python scripts/miscript.py arg1 arg2
```

## Scripts (si es necesario)

### miscript.py
Descripción de lo que hace el script.

```python
#!/usr/bin/env python3
"""Descripción del script."""
import sys

def main():
    pass

if __name__ == "__main__":
    main()
```
```

**Reglas de escritura:**

| Regla | Descripción |
|------|-------------|
| Mantener bajo 500 líneas | Mover contenido detallado a `references/` si es necesario |
| Usar forma imperativa | "Ejecuta el script" no "Deberías ejecutar el script" |
| Solo contexto no obvio | El agente ya conoce el conocimiento común |
| Sin secretos hardcodeados | Usa variables de entorno o archivos de configuración |

### Paso 5: Validar

El skill creator aplica estas reglas:

| Verificación | Regla |
|-------|------|
| Formato del nombre | `^[a-z0-9]+(-[a-z0-9]+)*$` (minúsculas, guiones) |
| Longitud del nombre | Menos de 64 caracteres |
| Separadores de ruta | Sin `/`, `\`, o `..` en el nombre |
| Node.js | ❌ Sin `package.json`, `node_modules`, ni dependencias npm |
| Secretos | ❌ Sin API keys, tokens o contraseñas hardcodeados |
| Tamaño de SKILL.md | Menos de 500 líneas (mover extras a `references/`) |
| Profundidad de referencias | Máximo 1 nivel desde SKILL.md |

### Paso 6: Iterar

Después de crear el skill, pruébalo en una tarea real:

1. Pide al agente que use el skill
2. Nota cualquier confusión o error
3. Actualiza el SKILL.md o los scripts para aclarar
4. Vuelve a probar

---

## Ejemplo Completo: Skill Rotador de PDF

### Solicitud del Usuario

```
Crea un skill para rotar páginas de PDF 90, 180 o 270 grados.
```

### Flujo del Agente

**Paso 1: Entender**
- El usuario quiere rotar páginas de PDF
- Necesita un script simple de Python usando `pypdf`
- Sin API compleja ni servicio externo

**Paso 2: Planificar**
- SKILL.md con instrucciones de uso
- `scripts/rotate_pdf.py` para la lógica de rotación

**Paso 3: Crear**
```bash
mkdir -p ~/.picoclaw/workspace/skills/pdf-rotator/scripts/
```

**Paso 4: Escribir**

SKILL.md:
```yaml
---
name: pdf-rotator
description: Rota páginas de PDF 90, 180 o 270 grados. Usar cuando el usuario necesita corregir orientación de PDF o rotar documentos escaneados.
---

# Rotador de PDF

Rota páginas en archivos PDF usando Python.

## Cuándo Usar
- El usuario pide rotar un PDF
- El usuario necesita corregir orientación de páginas
- El usuario tiene documentos escaneados al revés

## Requisitos

```bash
pip install pypdf
```

## Uso

```bash
python scripts/rotate_pdf.py entrada.pdf 90
```

Argumentos:
- `entrada.pdf` — Ruta al archivo PDF
- `90` — Ángulo de rotación (90, 180 o 270)

Salida: Crea `entrada_rotated.pdf` en el mismo directorio.

## Script

### scripts/rotate_pdf.py

```python
#!/usr/bin/env python3
"""Rota todas las páginas de un PDF por el ángulo especificado."""
import sys
from pypdf import PdfReader, PdfWriter

def rotate_pdf(input_path, angle):
    reader = PdfReader(input_path)
    writer = PdfWriter()

    for page in reader.pages:
        page.rotate(int(angle))
        writer.add_page(page)

    output_path = input_path.replace(".pdf", "_rotated.pdf")
    with open(output_path, "wb") as f:
        writer.write(f)

    print(f"Se rotaron {len(reader.pages)} páginas por {angle}°")
    print(f"Guardado en: {output_path}")
    return output_path

if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("Uso: rotate_pdf.py <entrada.pdf> <ángulo>")
        sys.exit(1)
    rotate_pdf(sys.argv[1], sys.argv[2])
```
```

**Paso 5: Validar**
- ✅ Nombre: `pdf-rotator` (kebab-case, < 64 chars)
- ✅ La descripción incluye desencadenantes de uso
- ✅ Sin secretos
- ✅ Solo Python, sin Node.js
- ✅ Menos de 500 líneas

**Paso 6: Probar**
```
Usuario: "Rota mi contrato.pdf 90 grados"
Agente: [Usa el skill pdf-rotator para rotar el PDF]
```

---

## Anti-Patrones a Evitar

### ❌ Nunca Usar Node.js

PicoClaw está diseñado para ser ligero (< 10MB RAM). Node.js + dependencias puede alcanzar fácilmente 50-150MB.

**Mal:**
```
skills/mi-skill/package.json
skills/mi-skill/node_modules/
```

**Bien:**
```
skills/mi-skill/scripts/miscript.py     # Python
skills/mi-skill/scripts/miscript.sh     # Bash
```

### ❌ Sin Documentación Extra

No crees `README.md`, `CHANGELOG.md`, `INSTALLATION_GUIDE.md`, etc. Estos archivos son solo para desarrolladores humanos — el agente solo necesita SKILL.md.

**Solo crea:**
- `SKILL.md` (obligatorio)
- `scripts/` (opcional)
- `references/` (opcional)
- `assets/` (opcional)

### ❌ No Inflar SKILL.md

Si tu SKILL.md excede 500 líneas, mueve el contenido detallado a `references/`:

```
skills/mi-skill/
├── SKILL.md              ← Guía de referencia rápida (< 500 líneas)
└── references/
    ├── api_docs.md       ← Documentación completa de API
    ├── schemas.md        ← Esquemas de base de datos
    └── workflows.md      ← Guías de flujo de trabajo complejas
```

Referencia desde SKILL.md:
```markdown
## Configuración Avanzada
Ver [references/api_docs.md](references/api_docs.md) para referencia completa de API.
```

### ❌ Sin Anidamiento Profundo de Referencias

Máximo 1 nivel de profundidad desde SKILL.md:

```
✅ SKILL.md → references/api_docs.md          (OK)
❌ SKILL.md → references/api_docs.md → sub/   (MAL)
```

### ❌ Sin Secretos Hardcodeados

**Mal:**
```python
API_KEY = "sk-abc123..."  # NO HAGAS ESTO
```

**Bien:**
```python
import os
API_KEY = os.environ.get("MY_API_KEY")  # pragma: allowlist secret
```

---

## Convenciones de Nomenclatura

| Convención | Regla | Ejemplos |
|-----------|------|----------|
| Mayúsculas | Solo minúsculas | `pdf-rotator` ✅, `PDF-Rotator` ❌ |
| Separador | Guiones o guiones bajos | `pdf-rotator` ✅, `pdf_rotator` ✅ |
| Longitud | < 64 caracteres | `backup-postgres` ✅, `backup-postgres-diario-con-politica-de-retencion` ❌ |
| Caracteres | Solo `a-z`, `0-9`, `-`, `_` | `mi-skill-1` ✅, `mi skill!` ❌ |
| Sin rutas | Sin `/`, `\`, `..` | `pdf-rotator` ✅, `skills/pdf-rotator` ❌ |

---

## Skills Basados en Archivos vs Skills Nativos en Go

| Aspecto | Skills Basados en Archivos (skill_creator) | Skills Nativos en Go |
|--------|-----------------------------------|-----------------|
| **Quién crea** | Agente (vía skill_creator) | Desarrollador (escribe código Go) |
| **Ubicación** | `~/.picoclaw/workspace/skills/` | `pkg/skills/` en el código fuente |
| **Compilación** | No — se cargan en tiempo de ejecución | Sí — compilados en el binario |
| **Formato** | SKILL.md + scripts opcionales | Structs + métodos Go |
| **Caso de uso** | Conocimiento personalizado, flujos de trabajo | Roles integrados, integraciones de herramientas |
| **Requiere recompilar** | No | Sí |

**Regla general:** Usa el Skill Creator para skills personalizados de usuario. Los skills nativos en Go son para características de nivel framework.

---

## Solución de Problemas

### Skill no reconocido

1. **Verifica el nombre:**
   ```bash
   ls ~/.picoclaw/workspace/skills/
   # Debería mostrar: nombre-tu-skill/
   ```

2. **Verifica que SKILL.md existe:**
   ```bash
   cat ~/.picoclaw/workspace/skills/nombre-tu-skill/SKILL.md
   ```

3. **Verifica el frontmatter YAML:**
   ```yaml
   ---
   name: nombre-tu-skill
   description: Una descripción con desencadenantes de uso.
   ---
   ```

### El script no se ejecuta

1. **Hazlo ejecutable:**
   ```bash
   chmod +x ~/.picoclaw/workspace/skills/nombre-tu-skill/scripts/miscript.py
   ```

2. **Verifica las dependencias:**
   ```bash
   pip install -r requirements.txt  # si aplica
   ```

3. **Prueba manualmente:**
   ```bash
   python ~/.picoclaw/workspace/skills/nombre-tu-skill/scripts/miscript.py
   ```

---

## Documentación Relacionada

- [ADDING_NATIVE_SKILLS.md](ADDING_NATIVE_SKILLS.md) — Guía para desarrolladores sobre cómo crear skills nativos en Go
- [NATIVE_SKILLS_LIST.md](NATIVE_SKILLS_LIST.md) — Lista completa de skills nativos integrados
- [SKILLS.md](SKILLS.md) — Documentación general de skills
- [SERVICE.md](SERVICE.md) — Gestión de servicios del sistema operativo

---

*¡El Skill Creator está integrado en el binario — no requiere instalación! Solo pide al agente que cree un skill.*
