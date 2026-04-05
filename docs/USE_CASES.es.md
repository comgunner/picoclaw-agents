# Caso de Uso Dinámico: Del Mensaje de Telegram al Fix en Local

Este documento describe el flujo **dinámico y en tiempo real** de cómo PicoClaw procesa una instrucción recibida por Telegram para ejecutar una reparación completa en el repositorio local usando la configuración de `config_dev_multiple_models.example.json`.

> **PicoClaw v3.4.1**: Ahora incluye **Comandos Slash Fast-path** para operaciones instantáneas y **Global Tracker** para consistencia multi-agente perfecta.

## Flujo Dinámico de Interacción

### 1. Entrada: El Disparador (Telegram)
El usuario, desde su móvil, detecta un problema o decide una mejora y envía un mensaje directo al bot:

**Usuario (Telegram):**
> 📲 *"Oye PM, he notado que las sesiones de usuario no expiran correctamente en la base de datos local. Por favor, revisa el código en `pkg/session/manager.go` y asegúrate de que el TTL se aplique. Una vez corregido, corre los tests y avísame cuando esté en GitHub."*

### 2. Procesamiento Central (Project Manager - GPT-4o)
PicoClaw recibe el webhook de Telegram. El **Project Manager** (PM) entra en acción:
- **Análisis de Intención**: El PM identifica tres tareas: Análisis/Fix, Testing y Despliegue.
- **Delegación Dinámica**: El PM no lo hace todo solo. Envía instrucciones técnicas al Senior Dev.

**PM ➔ Senior Dev (DeepSeek Reasoner):**
> 🤖 *"Tarea Crítica: El usuario reporta fallos en el TTL de sesiones. Analiza `pkg/session/manager.go`, implementa el fix y delega al QA la verificación de los tests unitarios."*

### 3. Ejecución Técnica (Senior Developer - DeepSeek Reasoner)
El **Senior Dev** recibe el contexto técnico. Gracias a su modelo de razonamiento profundo (`deepseek-reasoner`), analiza la lógica de tiempo en Go:
- Lee el archivo: `read_file("pkg/session/manager.go")`.
- Encuentra que el `time.After` no se está reiniciando.
- Aplica el fix dinámicamente: `edit_file(...)`.

### 4. Ciclo de Calidad (QA Engineer - Claude 3.5 Sonnet)
El Senior Dev invoca al **QA Engineer** de forma proactiva:
- **QA** ejecuta: `exec("go test ./pkg/session/...")`.
- Si fallan, el QA le reporta al Senior Dev para corregir. Si pasan, el QA procede.
- **QA ➔ GitHub**: Ejecuta los comandos de git para pushear la solución.

### 5. Cierre del Ciclo: Respuesta al Usuario (Telegram)
Una vez que la jerarquía termina, el PM consolida el reporte final y lo envía de vuelta al chat de Telegram del usuario.

**PM (Telegram):**
> ✅ *"¡Listo! El Senior Dev ha corregido el error en el manejo del TTL en `manager.go`. El QA ha verificado los tests (100% pass) y los cambios ya están en la rama principal. ¿Necesitas algo más?"*

## Por qué es "Dinámico"

1. **Inyección de Contexto**: El sistema lee el estado real de tu repositorio local *en el momento* que envías el mensaje.
2. **Jerarquía Autónoma**: El usuario solo habla con el PM; la "conversación" interna entre Senior Dev, QA y Junior Dev ocurre sin intervención humana.
3. **Multimodelado Inteligente**: Se activan diferentes LLMs según la complejidad de la subtarea que surgió de tu mensaje original.

## Beneficios de esta Configuración

1. **Eficiencia de Costos**: Usas modelos más económicos (`DeepSeek Chat`) para tareas simples y dejas los modelos potentes (`DeepSeek Reasoner`, `GPT-4o`) para el análisis crítico.
2. **Especialización**: Cada agente tiene un rol claro, evitando que un solo modelo se sature con demasiada información irrelevante.
3. **Paralelismo**: El QA puede estar verificando una parte del código mientras el Senior Dev analiza el siguiente paso, acelerando el ciclo de vida del fix.
---

## Caso de Uso 2: Generación Automatizada y Aprobación de Posts (Social Bundle)

Este flujo demuestra el uso del sistema de colas (**Batch System**) y los comandos rápidos (**Slash Commands**) para una gestión eficiente de contenido.

### 1. Instrucción Inicial (Telegram/Discord/CLI)
El usuario pide al agente que genere contenido para redes sociales:

**Usuario:**
> 📲 *"Genera un post para Facebook e Instagram sobre el lanzamiento de nuestra nueva versión v2.5. Quiero que incluya una imagen profesional y un texto persuasivo."*

### 2. Procesamiento en Segundo Plano (Social Manager)
El **Project Manager** delega la tarea al **Social Manager**:
- **Social Manager** inicia la herramienta `social_post_bundle`.
- El sistema responde inmediatamente con un ID: `⏳ Tarea iniciada: #IMA_GEN_02_03_26_1600. Te avisaré cuando esté listo.`
- El LLM se libera, ahorrando tokens mientras se genera la imagen (vía DALL-E/Gemini) y el texto.

### 3. Notificación de Entrega
Una vez completado el lote, PicoClaw envía el resultado al usuario con el texto final y la imagen adjunta.

**PicoClaw (Notificación):**
> 🎨 **Lote Generado: #IMA_GEN_02_03_26_1600**
> [Imagen Adjunta]
> *Texto sugerido: "¡Gran noticia! PicoClaw v2.5 ya está aquí..."*
> 💡 **Opciones (Copia y pega):**
> 1) `/bundle_approve id=20260302_1600`
> 2) `/bundle_regen id=20260302_1600`
> 3) `/bundle_edit id=20260302_1600`

### 4. Aprobación Instantánea (Fast-path)
El usuario revisa y decide aprobar desde su móvil o terminal:

**Usuario (Telegram/Discord):**
> `/bundle_approve id=20260302_1600`

**Resultado:**
- El sistema intercepta el comando `/` (Fast-path).
- **Procesamiento instantáneo**: Sin consultar a la IA, el sistema marca el lote como aprobado y lo publica directamente en Facebook/Instagram.
- El usuario recibe confirmación inmediata: `✅ Lote aprobado y publicado con éxito.`

## Ventajas del Sistema de Comandos
- **Cero Latencia**: La respuesta al comando `/` es inmediata.
- **Eficiencia**: No consume tokens de razonamiento para una simple aprobación.
- **Omnicanalidad**: Funciona igual en Discord (Slash Commands), Telegram y Terminal.

---

## Caso de Uso 3: Generación de Imágenes Multi-Agente con Global Tracker (v3.4.1+)

Con v3.4.1, el **Global Tracker** asegura consistencia perfecta en flujos de trabajo multi-agente:

### Escenario: Subagente genera, Agente Principal publica

1. **Usuario solicita contenido**:
   ```
   @picoclaw-agents spawn task='Genera imagen sobre IA y crea post para Twitter'
   ```

2. **Subagente trabaja**:
   - Genera imagen con `image_gen_create`
   - Crea post con `community_manager_create_draft`
   - Guarda en **Global Shared Workspace**

3. **Agente Principal puede acceder inmediatamente**:
   - Sin errores de "ID no encontrado"
   - Acceso instantáneo al trabajo del subagente
   - Puede aprobar y publicar sin demoras

### Beneficios

- ✅ **Estado Compartido**: Todos los agentes acceden al mismo workspace
- ✅ **Sin Problemas de Sync**: Los cambios son inmediatamente visibles
- ✅ **Escalable**: Funciona con cualquier número de subagentes
