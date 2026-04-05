# Contribuir a PicoClaw

> **Última Actualización:** Marzo 2026 | **Versión:** v3.4.5+

¡Gracias por tu interés en contribuir a PicoClaw! Este proyecto es un esfuerzo comunitario para construir un asistente de IA personal ligero y versátil. Damos la bienvenida a contribuciones de todo tipo: corrección de errores, características, documentación, traducciones y pruebas.

PicoClaw mismo fue desarrollado sustancialmente con asistencia de IA — abrazamos este enfoque y hemos construido nuestro proceso de contribución alrededor de él.

## Tabla de Contenidos

- [Código de Conducta](#código-de-conducta)
- [Formas de Contribuir](#formas-de-contribuir)
- [Primeros Pasos](#primeros-pasos)
- [Configuración de Desarrollo](#configuración-de-desarrollo)
- [Realizar Cambios](#realizar-cambios)
- [Contribuciones Asistidas por IA](#contribuciones-asistidas-por-ia)
- [Proceso de Pull Request](#proceso-de-pull-request)
- [Estrategia de Ramas](#estrategia-de-ramas)
- [Revisión de Código](#revisión-de-código)
- [Comunicación](#comunicación)

---

## Código de Conducta

Estamos comprometidos a mantener una comunidad acogedora y respetuosa. Sé amable, constructivo y asume la buena fe. El acoso o discriminación de cualquier tipo no será tolerado.

**Nuestros Estándares:**
- Usar un lenguaje acogedor e inclusivo
- Ser respetuoso con los diferentes puntos de vista y experiencias
- Aceptar graciosamente las críticas constructivas
- Centrarse en lo que es mejor para la comunidad
- Mostrar empatía hacia otros miembros de la comunidad

**Comportamiento Inaceptable:**
- Uso de lenguaje o imágenes sexualizadas
- Trolling, comentarios insultantes/despectivos, y ataques personales o políticos
- Acoso público o privado
- Publicar información privada de otros sin permiso explícito
- Otra conducta que podría razonablemente considerarse inapropiada en un entorno profesional

**Reportar:**
Si experimentas o presencias un comportamiento inaceptable, por favor repórtalo abriendo un issue privado o contactando directamente a un mantenedor.

---

## Formas de Contribuir

### 1. Reportes de Errores

**Cuándo Reportar:**
- Fallos o errores inesperados
- Características que no funcionan como está documentado
- Problemas de rendimiento
- Vulnerabilidades de seguridad (ver [SECURITY.es.md](./SECURITY.es.md))

**Cómo Reportar:**
1. Busca issues existentes para evitar duplicados
2. Usa la plantilla de reporte de errores
3. Incluye:
   - Versión de PicoClaw
   - SO y hardware
   - Pasos para reproducir
   - Comportamiento esperado vs real
   - Logs relevantes (redacta secretos)

### 2. Solicitudes de Características

**Antes de Solicitar:**
- Busca solicitudes de características existentes
- Verifica si la característica se alinea con los objetivos del proyecto (ligero, multi-agente, portable)

**Cómo Solicitar:**
1. Usa la plantilla de solicitud de características
2. Describe el caso de uso
3. Explica por qué es necesario
4. Sugiere posibles implementaciones (opcional)

### 3. Contribuciones de Código

**Tipos de Contribuciones de Código:**
- Corrección de errores
- Nuevas características
- Mejoras de rendimiento
- Mejoras de seguridad
- Adición de pruebas
- Refactorización

**Antes de Codificar:**
- Abre un issue para discutir el cambio
- Verifica si hay trabajo similar en progreso
- Asegúrate de tener tiempo para completar el PR

### 4. Documentación

**Necesidades de Documentación:**
- Mejoras al README
- Documentación de API
- Creación de tutoriales
- Traducción a nuevos idiomas
- Mejoras de comentarios en código
- Guías de solución de problemas

### 5. Pruebas

**Oportunidades de Prueba:**
- Probar en nuevas plataformas de hardware
- Probar con diferentes proveedores LLM
- Probar nuevos canales de chat
- Reportar problemas de compatibilidad
- Escribir casos de prueba

---

## Primeros Pasos

### 1. Hacer Fork del Repositorio

```bash
# En GitHub, haz clic en el botón "Fork"
# Luego clona tu fork
git clone https://github.com/<tu-usuario>/picoclaw.git
cd picoclaw
```

### 2. Agregar Remote Upstream

```bash
git remote add upstream https://github.com/comgunner/picoclaw-agents.git
git remote -v
# Debería mostrar tanto origin (tu fork) como upstream
```

### 3. Crear una Rama

```bash
git checkout main
git pull upstream main
git checkout -b feature/tu-nombre-de-caracteristica
```

**Nombre de Ramas:**
- `feature/xyz` — Nuevas características
- `fix/xyz` — Corrección de errores
- `docs/xyz` — Documentación
- `test/xyz` — Adición de pruebas
- `refactor/xyz` — Refactorización de código

---

## Configuración de Desarrollo

### Prerrequisitos

| Herramienta | Versión | Instalación |
|------|---------|--------------|
| **Go** | 1.25.8+ | [go.dev](https://go.dev/dl/) |
| **make** | 4.0+ | `apt install make` / `brew install make` |
| **git** | 2.30+ | `apt install git` / `brew install git` |
| **docker** | 24.0+ (opcional) | [docker.com](https://docs.docker.com/get-docker/) |

### Compilación

```bash
# Descargar dependencias
make deps

# Compilar binario (ejecuta go generate primero)
make build

# Verificación pre-commit completa
make check  # deps + fmt + vet + test
```

### Ejecutar Pruebas

```bash
# Ejecutar todas las pruebas
make test

# Ejecutar una prueba específica
go test -run TestName -v ./pkg/session/

# Ejecutar benchmarks
go test -bench=. -benchmem -run='^$' ./...

# Ejecutar con cobertura
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Estilo de Código

```bash
# Formatear código
make fmt

# Análisis estático
make vet

# Ejecutar linter completo
make lint
```

**Todas las verificaciones de CI deben pasar antes de que un PR pueda ser fusionado.** Ejecuta `make check` localmente antes de hacer push para detectar problemas temprano.

---

## Realizar Cambios

### Ramas

Siempre haz branch desde `main` y apunta a `main` en tu PR. Nunca hagas push directamente a `main` o cualquier rama `release/*`:

```bash
git checkout main
git pull upstream main
git checkout -b tu-rama-de-caracteristica
```

Usa nombres descriptivos de ramas:
- ✅ `fix/telegram-timeout`
- ✅ `feat/ollama-provider`
- ✅ `docs/contributing-guide`
- ❌ `patch-1`
- ❌ `nueva-caracteristica`
- ❌ `fix`

### Commits

**Guías de Mensajes de Commit:**
- Escribe mensajes claros y concisos en inglés
- Usa modo imperativo: "Add retry logic" no "Added retry logic"
- Referencia issues relacionados: `Fix session leak (#123)`
- Mantén los commits enfocados: un cambio lógico por commit

**Formato Conventional Commits:**
```
<tipo>(<alcance>): <descripción>

[cuerpo opcional]

[pie opcional]
```

**Tipos:**
- `feat`: Nueva característica
- `fix`: Corrección de errores
- `docs`: Documentación
- `style`: Formateo (sin cambio de código)
- `refactor`: Refactorización de código
- `test`: Adición de pruebas
- `chore`: Mantenimiento

**Ejemplos:**
```bash
feat(agents): Add autonomous runtime for background processing

fix(telegram): Fix message timeout in long conversations

docs(security): Add security best practices guide

test(providers): Add integration tests for Antigravity provider
```

### Mantenerse Actualizado

Haz rebase de tu rama sobre upstream `main` antes de abrir un PR:

```bash
git fetch upstream
git rebase upstream/main

# Resolver conflictos si los hay
git add <archivos>
git rebase --continue

# Force push a tu fork
git push -f origin tu-rama-de-caracteristica
```

---

## Contribuciones Asistidas por IA

PicoClaw fue construido con asistencia sustancial de IA, y abrazamos completamente el desarrollo asistido por IA. Sin embargo, los contribuyentes deben entender sus responsabilidades al usar herramientas de IA.

### La Divulgación es Requerida

Cada PR debe divulgar la participación de IA usando la sección **🤖 Generación de Código con IA** de la plantilla del PR. Hay tres niveles:

| Nivel | Descripción |
|-------|-------------|
| 🤖 **Totalmente generado por IA** | La IA escribió el código; el contribuyente lo revisó y validó |
| 🛠️ **Mayormente generado por IA** | La IA produjo el borrador; el contribuyente hizo modificaciones significativas |
| 👨‍💻 **Mayormente escrito por humano** | El contribuyente lideró; la IA proporcionó sugerencias o nada |

**Se espera divulgación honesta.** No hay estigma asociado a ningún nivel — lo que importa es la calidad de la contribución.

### Eres Responsable de lo que Envías

Usar IA para generar código no reduce tu responsabilidad como contribuyente. Antes de abrir un PR con código generado por IA, debes:

1. **Leer y entender** cada línea del código generado
2. **Probarlo** en un entorno real (ver la sección Entorno de Prueba de la plantilla del PR)
3. **Verificar problemas de seguridad** — Los modelos de IA pueden generar código sutilmente inseguro (ej. path traversal, inyección, exposición de credenciales). Revisa cuidadosamente.
4. **Verificar corrección** — La lógica generada por IA puede sonar plausible pero estar equivocada. Valida el comportamiento, no solo la sintaxis.

**Los PRs donde es evidente que el contribuyente no ha leído o probado el código generado por IA serán cerrados sin revisión.**

### Estándares de Calidad de Código Generado por IA

Las contribuciones generadas por IA se mantienen al **mismo estándar de calidad** que el código escrito por humanos:

- Debe pasar todas las verificaciones de CI (`make check`)
- Debe ser Go idiomático y consistente con el estilo de código existente
- No debe introducir abstracciones innecesarias, código muerto, o sobre-ingeniería
- Debe incluir o actualizar pruebas donde sea apropiado

### Revisión de Seguridad

El código generado por IA requiere escrutinio de seguridad extra. Presta especial atención a:

- **Manejo de rutas de archivo y escapes de sandbox** (ver commit `244eb0b` para un ejemplo real)
- **Validación de entrada externa** en manejadores de canales e implementaciones de herramientas
- **Manejo de credenciales o secretos**
- **Ejecución de comandos** (`exec.Command`, invocaciones shell)

Si no estás seguro de si una pieza de código generado por IA es segura, dilo en el PR — los revisores ayudarán.

---

## Proceso de Pull Request

### Antes de Abrir un PR

**Lista de Verificación:**
- [ ] Ejecuta `make check` y asegúrate de que pase localmente
- [ ] Completa la plantilla del PR completamente, incluyendo la sección de divulgación de IA
- [ ] Vincula cualquier issue relacionado(s) en la descripción del PR
- [ ] Mantén el PR enfocado. Evita agrupar cambios no relacionados
- [ ] Actualiza documentación si es necesario
- [ ] Agrega o actualiza pruebas
- [ ] Actualiza CHANGELOG.md para cambios visibles al usuario

### Secciones de la Plantilla del PR

La plantilla del PR solicita:

1. **Descripción** — ¿Qué hace este cambio y por qué?
2. **Tipo de Cambio** — Corrección de error, característica, docs, o refactor
3. **Generación de Código con IA** — Divulgación de participación de IA (requerido)
4. **Issue Relacionado** — Vincula al issue que aborda
5. **Contexto Técnico** — URLs de referencia y razonamiento (omitir para PRs puros de docs)
6. **Entorno de Prueba** — Hardware, SO, modelo/proveedor, y canales usados para pruebas
7. **Evidencia** — Logs o capturas de pantalla opcionales demostrando que el cambio funciona
8. **Lista de Verificación** — Confirmación de auto-revisión

### Tamaño del PR

**Prefiere PRs pequeños y revisables:**
- Un PR que cambia 200 líneas en 5 archivos es mucho más fácil de revisar que uno que cambia 2000 líneas en 30 archivos
- Si tu característica es grande, considera dividirla en una serie de PRs más pequeños, lógicamente completos

**Ejemplo de Buena División de PR:**
```
PR 1: Agregar interfaz de proveedor
PR 2: Implementar proveedor OpenAI
PR 3: Implementar proveedor Anthropic
PR 4: Agregar factory de proveedores
```

### Ejemplo de Flujo de Trabajo de PR

```bash
# 1. Crear rama
git checkout -b feat/nuevo-proveedor

# 2. Realizar cambios
# Editar archivos...

# 3. Preparar y hacer commit
git add pkg/providers/nuevo_proveedor.go
git commit -m "feat(providers): Add new LLM provider

- Implement Provider interface
- Add authentication support
- Include comprehensive tests"

# 4. Push al fork
git push origin feat/nuevo-proveedor

# 5. Abrir PR en GitHub
# Navegar a https://github.com/comgunner/picoclaw-agents/pulls
# Hacer clic en "New Pull Request"
# Completar plantilla
```

---

## Estrategia de Ramas

### Ramas de Larga Duración

| Rama | Propósito | Protección |
|--------|---------|------------|
| **`main`** | Desarrollo activo | Requiere 1+ aprobación de mantenedor |
| **`release/x.y`** | Releases estables | Estrictamente protegida, sin pushes directos |

### Requisitos para Fusionar en `main`

Un PR solo puede ser fusionado cuando se satisfacen todos los siguientes:

1. **CI pasa** — Todos los workflows de GitHub Actions (lint, test, build) deben estar en verde
2. **Aprobación de revisor** — Al menos un mantenedor ha aprobado el PR
3. **Sin comentarios de revisión sin resolver** — Todos los hilos de revisión deben estar resueltos
4. **La plantilla del PR está completa** — Incluyendo divulgación de IA y entorno de prueba

### Quién Puede Fusionar

**Solo los mantenedores pueden fusionar PRs.** Los contribuyentes no pueden fusionar sus propios PRs, incluso si tienen acceso de escritura.

### Estrategia de Fusión

Usamos **squash merge** para la mayoría de los PRs para mantener el historial de `main` limpio y legible. Cada PR fusionado se convierte en un solo commit referenciando el número del PR:

```
feat: Add Ollama provider support (#491)
```

Si un PR consiste en múltiples commits independientes, bien separados, que cuentan una historia clara, se puede usar una fusión regular a discreción del mantenedor.

### Ramas de Release

Cuando una versión está lista, los mantenedores crean una rama `release/x.y` desde `main`. Después de ese punto:

- **Las nuevas características no se retroportan.** La rama de release no recibe nueva funcionalidad después de ser cortada.
- **Las correcciones de seguridad y errores críticos se cherry-pick.** Si una corrección en `main` califica (vulnerabilidad de seguridad, pérdida de datos, fallo), los mantenedores harán cherry-pick del commit(s) relevante(s) a la rama `release/x.y` afectada y emitirán una release de parche.

Si crees que una corrección en `main` debería ser retroportada a una rama de release, anótalo en la descripción del PR o abre un issue separado. La decisión recae en los mantenedores.

**Las ramas de release tienen protecciones más estrictas que `main` y nunca se les hace push directamente bajo ninguna circunstancia.**

---

## Revisión de Código

### Para Contribuyentes

**Responsabilidades:**
- Responde a los comentarios de revisión en un tiempo razonable (48 horas preferido)
- Cuando actualices un PR en respuesta a retroalimentación, nota brevemente qué cambió
- Si no estás de acuerdo con la retroalimentación, participa respetuosamente. Explica tu razonamiento; los revisores también pueden estar equivocados
- No hagas force-push después de que una revisión haya comenzado — hace más difícil para los revisores ver qué cambió. Usa commits adicionales en su lugar; el mantenedor hará squash al fusionar

**Ejemplo de Respuesta:**
```markdown
@reviewer ¡Gracias por la retroalimentación! He actualizado el código para:
- Usar `sync.RWMutex` en lugar de `sync.Mutex` para mejor rendimiento de lectura
- Agregar manejo de errores para el caso extremo X
- Actualizar pruebas para cubrir el nuevo comportamiento
```

### Para Revisores

**Revisar Para:**

1. **Corrección**
   - ¿El código hace lo que afirma?
   - ¿Hay casos extremos?
   - ¿Hay condiciones de carrera?

2. **Seguridad**
   - Especialmente para código generado por IA, implementaciones de herramientas y manejadores de canales
   - Verificar path traversal, inyección, exposición de credenciales

3. **Arquitectura**
   - ¿El enfoque es consistente con el diseño existente?
   - ¿Esto agrega complejidad innecesaria?

4. **Simplicidad**
   - ¿Hay una solución más simple?
   - ¿Esto introduce sobre-ingeniería?

5. **Pruebas**
   - ¿Los cambios están cubiertos por pruebas?
   - ¿Las pruebas existentes aún tienen sentido?

**Sé constructivo y específico:**
- ✅ "Esto podría tener una condición de carrera si dos goroutines llaman esto concurrentemente — considera usar un mutex aquí"
- ❌ "esto se ve mal"

### Lista de Revisores

Una vez que tu PR sea enviado, puedes contactar a los revisores asignados:

| Función | Revisor |
|----------|----------|
| Provider | @yinwm |
| Channel | @yinwm |
| Agent | @lxowalle |
| Tools | @lxowalle |
| Skill | — |
| MCP | — |
| Optimization | @lxowalle |
| Security | — |
| AI CI | @imguoguo |
| UX | — |
| Document | — |

---

## Comunicación

### Dónde Comunicarse

| Plataforma | Propósito |
|----------|---------|
| **GitHub Issues** | Reportes de errores, solicitudes de características, discusiones de diseño |
| **GitHub Discussions** | Preguntas generales, ideas, conversación comunitaria |
| **Comentarios de Pull Request** | Retroalimentación específica de código |
| **Discord** | [Próximamente] |

### En Caso de Duda

**Abre un issue antes de escribir código.** Cuesta poco y previene esfuerzo desperdiciado.

**Buenas Preguntas para Hacer:**
- "¿Esta característica está alineada con los objetivos del proyecto?"
- "¿Alguien ya ha trabajado en esto?"
- "¿Cuál es el mejor enfoque para X?"
- "¿Puedes revisar mi diseño antes de que lo implemente?"

### Expectativas de Respuesta

- **Mantenedores:** Objetivo responder dentro de 48 horas
- **Contribuyentes:** Responder a comentarios de revisión dentro de 48 horas
- **Comunidad:** Sé paciente y comprensivo — todos están donando su tiempo

---

## Una Nota sobre el Origen Impulsado por IA del Proyecto

La arquitectura de PicoClaw fue sustancialmente diseñada e implementada con asistencia de IA, guiada por supervisión humana. Si encuentras algo que se vea extraño o sobre-ingenierizado, puede ser un artefacto de ese proceso — abrir un issue para discutirlo siempre es bienvenido.

**Creemos que el desarrollo asistido por IA hecho responsablemente produce grandes resultados. También creemos que los humanos deben permanecer responsables de lo que envían. Estas dos creencias no están en conflicto.**

---

## Reconocimiento

Los contribuyentes son reconocidos de las siguientes maneras:

1. **CHANGELOG.md** — Contribuciones notables son mencionadas en el changelog
2. **GitHub Contributors Graph** — Visible en el repositorio
3. **Release Notes** — Los principales contribuyentes pueden ser mencionados en las notas de release
4. **Documentación** — Los contribuyentes significativos pueden ser listados en el README

---

## ¿Preguntas?

Si tienes preguntas sobre contribuir, por favor:

1. Revisa este documento primero
2. Busca issues y discusiones existentes
3. Abre una nueva discusión en GitHub
4. Contacta a un mantenedor directamente

**¡Gracias por contribuir a PicoClaw!** 🎉

---

*PicoClaw: IA Ultra-Eficiente en Go. Hardware de $10 · 10MB RAM · <1s de Inicio*
