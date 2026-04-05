# PicoClaw: Skills Sentinel (Centinela de Habilidades)

El **Sentinel** (`SkillsSentinelTool`) es un mecanismo de seguridad interno integrado en PicoClaw diseñado para defender al agente ante la inyección de prompts, extracción del sistema (system prompt extraction) y evitar la ejecución de rutinas maliciosas de código.

## ¿Cómo Funciona?

La herramienta intercepta e inspecciona texto que puede provenir de las interacciones con el usuario o escanear archivos internos y extensions/skills, operando con base en dos acciones principales (`actions`):

1. **Validación de Texto (`validate` - Acción por defecto):** Compara el texto de entrada (`input`) del agente contra una lista negra de expresiones regulares (regex) para identificar cadenas catalogadas de alto riesgo.
2. **Escaneo de Skills (`scan`):** Realiza una auditoría sobre el sistema de archivos local para buscar patrones maliciosos escritos en las *skills* instaladas.

---

## Categorías de Amenazas Detectadas (Lista Negra Mantenida)

El diseño del Sentinel bloquea por patrón diversos vectores de ataque muy comunes en el uso de LLMs y Agentes:

- **Inyección de Prompts y Extracción de Sistema:** 
  Impide que comandos evasivos obliguen al agente a olvidar sus lineamientos (`ignore previous instructions`, `bypass`, `override system`) o divulgar sus instrucciones base y configuración (`reveal system instructions`, `dump configuration`, o forzar modos como `DAN`).
- **Scripts de Ingeniería Social / Descargas (ClickFix):**
  Deshabilita vectores comunes para descargar y ejecutar malware automáticamente del tipo `curl ... | bash`, `wget ... | sh`, o usando recursos homólogos de PowerShell (`iex`).
- **RATs (Troyanos de Acceso Remoto) y Reverse Shells:**
  Veta la ejecución de llamadas a conexiones inversas utilizadas por atacantes (ej. `bash -i >& /dev/tcp/...`, utilidades como `netcat -e`, y la generación de *sockets binding* mediante Python).
- **Robo y Exfiltración de Información:**
  Detecta la extracción no autorizada de credenciales, variables del sistema e historiales. Algunos ejemplos bloqueados incluyen `cat .ssh/id_rsa`, `history | grep`, filtrado mediante `env | curl` o interacción con los `keychains` (`security find-internet-password`).

---

## Excepciones (Modo "Self-Aware")

Aumentar la seguridad utilizando filtros estrictos suele desencadenar falsos positivos, en especial si un usuario está preguntando cómo funciona PicoClaw. El Sentinel aborda esta situación incorporando un mecanismo *Self-Aware* o consciente del propio agente.

Si una cadena de entrada contiene términos relacionados al propio entorno (`picoclaw`, `herramienta`, `tool`, `sentinel`, `skill`), y se detecta una clara estructura de pregunta (como poseer signos de interrogación `?`, `¿` o palabras del tipo `qué`, `cómo`, `saber`, `how`, `what`), el Sentinel **no levantará la alerta** para esta consulta específica, considerándola una pregunta legítima acerca del sistema.

---

## Modos de Suspensión Temporal (Mantenimiento)

Bajo ciertas circunstancias controladas (ej. la configuración manual de herramientas seguras), el Sentinel debe ser deshabilitado. Se encuentra protegido por _mutex_ y posee controles para ser suspendido temporalmente.

- **`Disable(duration)`:** Apaga el centinela estrictamente durante el tiempo solicitado. En esta ventana, devolverá un estado señalando que está "suspendido para tareas de configuración".
- El sistema cuenta con reactivación automática al caducar el bloque de tiempo (`disabledUntil`), emitiendo un evento (`callback` notificador) que le indica a la capa de control superior que ha retornado a su función protectora (`onAutoReactivate`).
- **`Enable()`:** Permite el encendido manual inmediato, interrumpiendo un posible tiempo de Disable en curso.

---

## Escáner Profundo de Archivos (`scan`)

Cuando se llama al *Sentinel* con la acción `scan`, este audita los archivos ejecutables y de configuración locales de las skills vinculadas a PicoClaw. 

El centinela revisa:
1. El directorio `skills/` del *Workspace* actual de PicoClaw si está definido.
2. Los directorios base compartidos de *PicoClaw*: `.picoclaw/skills/` y `.picoclaw/extensions/`.
3. Hasta 2 niveles de profundidad en el directorio local de módulos (`.picoclaw/node_modules/`) para revisar dependencias inyectadas bajo packages genéricos.

Solamente revisa extensiones con posibilidad de ejecución agresiva o alteración paramétrica (`.js`, `.ts`, `.sh`, `.py`, `skill.md` y `package.json`). Si se encuentra una firma comprometida en cualquier archivo en esta lista, la ejecución de la acción de escaneo resulta en un error de seguridad que emite y enumera cada archivo viciado, aconsejando que la skill sea desinstalada del ordenador (`picoclaw skills remove <name>`).
