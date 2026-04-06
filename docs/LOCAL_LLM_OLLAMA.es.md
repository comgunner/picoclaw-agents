# Ejecutar PicoClaw-Agents con LLMs Locales via Ollama

Ejecuta modelos de IA **100% offline** en tu propio hardware — sin API keys, sin nube, sin que tus datos salgan de tu máquina.

---

## ¿Qué es Ollama?

[Ollama](https://ollama.com) es la forma más sencilla de ejecutar modelos de lenguaje localmente. Provee una API compatible con el formato de OpenAI, lo que significa que picoclaw-agents se conecta a él sin configuración adicional.

---

## 1. Instalar Ollama

### macOS

```bash
# Opción A — Descarga directa (recomendada)
# https://ollama.com/download/mac
# Descarga y ejecuta el instalador .dmg

# Opción B — Homebrew
brew install ollama
```

### Linux

```bash
curl -fsSL https://ollama.com/install.sh | sh
```

### Windows

```powershell
irm https://ollama.com/install.ps1 | iex
```

### Termux (Android)

```bash
pkg update
pkg install ollama
```

---

## 2. Encontrar el Modelo Correcto para tu Hardware

¿No sabes qué modelo cabe en tu máquina? Usa `llm-checker` para obtener recomendaciones personalizadas:

```bash
# Instalar
npm install -g llm-checker

# Detectar tu hardware
llm-checker hw-detect

# Obtener recomendaciones por categoría
llm-checker recommend --category coding
llm-checker recommend --category general

# Ejecutar con selección automática
llm-checker ai-run --category coding --prompt "Write a hello world in Python"
```

**Ejemplo de salida en Mac M1 (8GB RAM):**
```
INFORMACIÓN DEL SISTEMA
│ CPU: M1 (8 núcleos, 2.4GHz)
│ Arquitectura: Apple Silicon
│ RAM: 8GB
│ GPU: Apple M1 / VRAM: 4GB (Integrado)
│ Nivel de Hardware: BAJO

RECOMENDACIONES INTELIGENTES
│ MEJOR GENERAL:  yi:6b          → ollama pull yi:6b
│ Programación:   deepseek-coder:6.7b
│ Razonamiento:   deepseek-coder:33b
│ Multimodal:     qwen2.5vl:7b
│ Creativo:       yi:6b
│ Chat:           yi:6b
│ Lectura:        qwen3:1.7b
```

> **Tip:** En Termux (Android), también instala llm-checker: `npm install -g llm-checker`

---

## 3. Modelos Recomendados (Selección de picoclaw-agents)

Estos modelos son ligeros, rápidos y funcionan bien como backend del agente:

| Modelo | Tamaño | Ideal Para | Comando |
|--------|--------|------------|---------|
| `llama3.2:1b` | ~800MB | Chat general, respuestas rápidas | `ollama pull llama3.2:1b` |
| `qwen2.5:0.5b` | ~400MB | Ultra-ligero, RAM mínima | `ollama pull qwen2.5:0.5b` |
| `qwen2.5-coder:0.5b` | ~400MB | Generación de código | `ollama pull qwen2.5-coder:0.5b` |

**Descargar y ejecutar:**
```bash
ollama pull llama3.2:1b
ollama run llama3.2:1b
```

## Modelos Populares Organizados por RAM

Visita **https://ollama.com/library** para el catálogo completo. Modelos agrupados por requisitos mínimos del sistema:

### Modelos para ≤ 4 GB de RAM

| Modelo | Tamaño | Ideal Para | Comando |
|--------|--------|------------|---------|
| `qwen2.5:0.5b` | ~400MB | Ultra-ligero, RAM mínima | `ollama pull qwen2.5:0.5b` |
| `qwen2.5-coder:0.5b` | ~400MB | Generación de código | `ollama pull qwen2.5-coder:0.5b` |
| `llama3.2:1b` | ~800MB | Chat general, respuestas rápidas | `ollama pull llama3.2:1b` |
| `qwen2.5:1.5b` | ~1GB | Multilingüe, uso general | `ollama pull qwen2.5:1.5b` |
| `deepseek-r1:1.5b` | ~1GB | Razonamiento chain-of-thought | `ollama pull deepseek-r1:1.5b` |
| `qwen3:0.6b` | ~400MB | Ultra-ligero, chat rápido | `ollama pull qwen3:0.6b` |
| `qwen3:1.7b` | ~1GB | Balance velocidad y capacidad | `ollama pull qwen3:1.7b` |
| `nomic-embed-text` | ~270MB | Embeddings de texto, RAG | `ollama pull nomic-embed-text` |

### Modelos para 8 GB de RAM

| Modelo | Tamaño | Ideal Para | Comando |
|--------|--------|------------|---------|
| `llama3.2:3b` | ~2GB | Balance velocidad y capacidad | `ollama pull llama3.2:3b` |
| `qwen2.5:3b` | ~2GB | Buen rendimiento general | `ollama pull qwen2.5:3b` |
| `qwen3:4b` | ~3GB | Último Qwen, excelente razonamiento | `ollama pull qwen3:4b` |
| `qwen3-coder:4b` | ~3GB | Último Qwen Coder, generación de código | `ollama pull qwen3-coder:4b` |
| `deepseek-coder:6.7b` | ~4GB | Generación y comprensión de código | `ollama pull deepseek-coder:6.7b` |
| `mistral:7b` | ~4GB | Buen razonamiento, uso general | `ollama pull mistral:7b` |
| `llava:7b` | ~4GB | Multimodal (texto + imágenes) | `ollama pull llava:7b` |
| `codellama:7b` | ~4GB | Enfocado en código, calidad Meta | `ollama pull codellama:7b` |
| `gemma4:e2b` | 7.2 GB | Último de Google, multilingüe | `ollama pull gemma4:e2b` |

### Modelos para 16 GB de RAM o Más

| Modelo | Tamaño | Ideal Para | Comando |
|--------|--------|------------|---------|
| `llama3.2:8b` | ~5GB | Uso general de alta calidad | `ollama pull llama3.2:8b` |
| `qwen2.5:7b` | ~5GB | Excelente en código y razonamiento | `ollama pull qwen2.5:7b` |
| `qwen3:8b` | ~5GB | Último Qwen, fuerte razonamiento | `ollama pull qwen3:8b` |
| `qwen3-coder:8b` | ~5GB | Último Qwen Coder, top programación | `ollama pull qwen3-coder:8b` |
| `phi4:14b` | ~9GB | El más inteligente de Microsoft | `ollama pull phi4:14b` |
| `qwen3:14b` | ~9GB | Qwen de alta calidad, multilingüe | `ollama pull qwen3:14b` |
| `qwen3-coder:14b` | ~9GB | Qwen Coder avanzado | `ollama pull qwen3-coder:14b` |
| `deepseek-r1:7b` | ~5GB | Razonamiento profundo, chain-of-thought | `ollama pull deepseek-r1:7b` |
| `deepseek-coder-v2:16b` | ~10GB | Modelo de programación top | `ollama pull deepseek-coder-v2:16b` |
| `qwen3:32b` | ~20GB | Máxima calidad Qwen | `ollama pull qwen3:32b` |
| `gemma4:e4b` | 9.6 GB | Último de Google, última variante | `ollama pull gemma4:e4b` |
| `gemma4:26b` | 18 GB | Mejor calidad Gemma 4 | `ollama pull gemma4:26b` |
| `gemma4:31b` | 20 GB | Máxima calidad Gemma 4 | `ollama pull gemma4:31b` |
| `llama3.1:70b` (Q4) | ~40GB | Máxima capacidad (cuantizado) | `ollama pull llama3.1:70b` |

### 🆕 Google Gemma 4 — Requisitos Detallados

Gemma 4 es la última familia de modelos abiertos de Google, con excelente capacidad multilingüe. Visita **[https://ollama.com/library/gemma4](https://ollama.com/library/gemma4)** para el catálogo completo.

| Variante | Peso | RAM/VRAM Mínima | RAM/VRAM Ideal | Comando |
|----------|------|-----------------|----------------|---------|
| `gemma4:e2b` | 7.2 GB | 8 GB | 12 GB | `ollama pull gemma4:e2b` |
| `gemma4:e4b` (latest) | 9.6 GB | 12 GB | 16 GB | `ollama pull gemma4:e4b` |
| `gemma4:26b` | 18 GB | 24 GB | 24 GB+ | `ollama pull gemma4:26b` |
| `gemma4:31b` | 20 GB | 24 GB | 32 GB | `ollama pull gemma4:31b` |

#### CPU vs GPU para Gemma 4

- **CPU:** Ollama funciona prácticamente con cualquier procesador moderno de los últimos años (requiere soporte para instrucciones AVX). Sin embargo, procesar estos modelos usando únicamente el CPU es un proceso lento.
- **GPU (Recomendado):** Ollama detecta automáticamente si tienes una tarjeta de video (Nvidia, AMD o un chip de Apple Silicon) y enviará el modelo allí para que funcione mucho más rápido. La clave es que tu tarjeta gráfica tenga suficiente VRAM para alojar el peso del modelo elegido.

#### Guía de Hardware para Gemma 4

| Tu Hardware | Variante Recomendada |
|-------------|---------------------|
| 8 GB RAM | `gemma4:e2b` (será lento, 12 GB ideal) |
| 12 GB RAM | `gemma4:e2b` o `gemma4:e4b` |
| 16 GB RAM | `gemma4:e4b` (latest) |
| 24 GB VRAM (GPU) | `gemma4:26b` |
| 32 GB VRAM (GPU) | `gemma4:31b` (mejor calidad) |

> **Nota para usuarios de PicoClaw-Agents:** Las variantes de Gemma 4 tienen ventanas de contexto moderadas. Para mejores resultados con el bucle del agente, usa `gemma4:e2b` o `gemma4:e4b` en sistemas con RAM limitada. Las variantes más grandes (26b, 31b) requieren VRAM significativa y funcionan mejor en configuraciones con GPU dedicada.

### 🆕 Qwen 3 y Qwen 3 Coder — Requisitos Detallados

Qwen 3 es la última familia de modelos abiertos de Alibaba con excelente capacidad multilingüe y de razonamiento. Qwen 3 Coder es la variante especializada para generación de código. Visita **[https://ollama.com/library/qwen3](https://ollama.com/library/qwen3)** y **[https://ollama.com/library/qwen3-coder](https://ollama.com/library/qwen3-coder)** para los catálogos completos.

#### Qwen 3 (Uso General)

| Variante | Peso | RAM/VRAM Mínima | RAM/VRAM Ideal | Comando |
|----------|------|-----------------|----------------|---------|
| `qwen3:0.6b` | ~400MB | 2 GB | 4 GB | `ollama pull qwen3:0.6b` |
| `qwen3:1.7b` | ~1GB | 4 GB | 4 GB | `ollama pull qwen3:1.7b` |
| `qwen3:4b` | ~3GB | 6 GB | 8 GB | `ollama pull qwen3:4b` |
| `qwen3:8b` | ~5GB | 8 GB | 16 GB | `ollama pull qwen3:8b` |
| `qwen3:14b` | ~9GB | 12 GB | 16 GB | `ollama pull qwen3:14b` |
| `qwen3:32b` | ~20GB | 24 GB | 32 GB | `ollama pull qwen3:32b` |

#### Qwen 3 Coder (Generación de Código)

| Variante | Peso | RAM/VRAM Mínima | RAM/VRAM Ideal | Comando |
|----------|------|-----------------|----------------|---------|
| `qwen3-coder:4b` | ~3GB | 6 GB | 8 GB | `ollama pull qwen3-coder:4b` |
| `qwen3-coder:8b` | ~5GB | 8 GB | 16 GB | `ollama pull qwen3-coder:8b` |
| `qwen3-coder:14b` | ~9GB | 12 GB | 16 GB | `ollama pull qwen3-coder:14b` |
| `qwen3-coder:32b` | ~20GB | 24 GB | 32 GB | `ollama pull qwen3-coder:32b` |

#### CPU vs GPU para Qwen 3

- **CPU:** Funciona con cualquier procesador moderno (requiere soporte AVX). Las variantes pequeñas (0.6b, 1.7b) corren razonablemente solo con CPU. Las variantes más grandes serán lentas.
- **GPU (Recomendado):** Ollama detecta automáticamente GPUs Nvidia, AMD o Apple Silicon y descarga la computación para una inferencia mucho más rápida. La clave es tener suficiente VRAM para tu variante elegida.

#### Guía de Hardware para Qwen 3

| Tu Hardware | Variante Recomendada |
|-------------|---------------------|
| 4 GB RAM | `qwen3:0.6b` o `qwen3:1.7b` |
| 8 GB RAM | `qwen3:4b` o `qwen3-coder:4b` |
| 16 GB RAM | `qwen3:8b`, `qwen3:14b`, o variantes coder |
| 32 GB RAM/VRAM | `qwen3:32b` o `qwen3-coder:32b` (mejor calidad) |

> **Nota para usuarios de PicoClaw-Agents:** Los modelos Qwen 3 tienen buenas ventanas de contexto y funcionan bien con el bucle del agente. Para tareas de código, prefiere las variantes `qwen3-coder`. Para uso general, `qwen3:4b` en sistemas de 8GB o `qwen3:8b` en sistemas de 16GB ofrecen el mejor balance.

---

## 3b. Limitar Uso de RAM, CPU y GPU

> **Verificado:** Todos los parámetros abajo son **configuraciones oficiales de Ollama** documentadas en la especificación [Ollama Modelfile](https://github.com/ollama/ollama/blob/main/docs/modelfile.md).

### Dónde Vive la Configuración de Ollama (Por Sistema Operativo)

Ollama **no usa un archivo `.env`**. La configuración se maneja mediante variables de entorno al iniciar el servidor:

| SO | Archivo / Ubicación de Config | Ejemplo |
|----|----------------------|---------|
| **macOS** | `~/.config/ollama/config.json` (raro)<br>O launchd plist: `~/Library/LaunchAgents/com.ollama.ollama.plist` | Ver sección macOS abajo |
| **Linux** | Override systemd: `systemctl edit ollama.service`<br>O `~/.bashrc` / `~/.zshrc` | Ver sección Linux abajo |
| **Windows** | Variables de Entorno del Sistema (Configuración → Sistema → Acerca de → Avanzado → Variables de Entorno)<br>O perfil PowerShell: `$PROFILE` | Ver sección Windows abajo |
| **Termux** | `~/.bashrc` o `~/.zshrc` | Ver sección Termux abajo |

**Variables de entorno más comunes:**

| Variable | Propósito | Default |
|----------|---------|---------|
| `OLLAMA_HOST` | Dirección de escucha (ej: `0.0.0.0:11434` para acceso en red) | `127.0.0.1:11434` |
| `OLLAMA_KEEP_ALIVE` | Cuánto tiempo el modelo permanece en RAM (`5m`, `1h`, `-1` siempre) | `5m` |
| `OLLAMA_NUM_PARALLEL` | Máx solicitudes concurrentes | `1` |
| `OLLAMA_MAX_LOADED_MODELS` | Máx modelos cargados simultáneamente | `1` |
| `OLLAMA_GPU_ENABLED` | Poner a `0` para desactivar GPU completamente | `1` |
| `OLLAMA_TMPDIR` | Directorio temporal para carga de modelos | Temp del sistema |

### Tres Formas de Aplicar Límites de Recursos

### Método A: Vía `/set` en el CLI (Interactivo, Solo Sesión)

Cuando ejecutas un modelo interactivamente (`ollama run llama3`), ajusta parámetros sobre la marcha:

```
/set parameter num_thread 4
/set parameter num_ctx 2048
/set parameter num_gpu 10
```

Los cambios tienen efecto inmediato pero **se pierden al salir de la sesión**.

### Método B: Vía Modelfile (Permanente — Recomendado para picoclaw-agents)

#### Dónde Guardar el Modelfile

El Modelfile es un archivo de texto plano que **creas y guardas donde quieras**. Ollama lo lee una sola vez para construir tu modelo personalizado, y luego el modelo se almacena permanentemente en el almacenamiento interno de Ollama. El Modelfile se puede borrar después de construir — pero es buena práctica conservarlo como referencia.

| SO | Ubicación Recomendada | Almacenamiento Interno del Modelo |
|----|----------------------|----------------------------------|
| **macOS** | `~/ollama-modelfiles/` | `~/.ollama/models/` |
| **Linux** | `~/ollama-modelfiles/` | `~/.ollama/models/` |
| **Windows** | `C:\Users\TuUsuario\ollama-modelfiles\` | `C:\Users\TuUsuario\.ollama\models\` |
| **Termux** | `~/ollama-modelfiles/` | `~/.ollama/models/` |

**Flujo de trabajo:**

```bash
# 1. Crear un directorio para tus Modelfiles (donde quieras)
mkdir -p ~/ollama-modelfiles
cd ~/ollama-modelfiles

# 2. Crear el Modelfile (usa cualquier editor de texto)
nano Modelfile    # o: vim Modelfile, code Modelfile

# 3. Construir el modelo personalizado (Ollama lee el archivo una vez)
ollama create mi-modelo-personalizado -f Modelfile

# 4. El modelo ahora está almacenado en el almacenamiento interno de Ollama
#    Puedes borrar el Modelfile si quieres, pero es útil conservarlo
ollama list

# 5. Usar el modelo
ollama run mi-modelo-personalizado
```

**Concepto clave:** El Modelfile es como una **receta**. Una vez que horneas el pastel (`ollama create`), ya no necesitas la receta — pero es útil si quieres hornear otro después.

#### Ejemplo 1: Gemma 4:e2b — Limitado a 8GB de RAM

Para una Mac o laptop con 8GB de RAM, manteniendo el uso de RAM bajo control:

```bash
mkdir -p ~/ollama-modelfiles
cd ~/ollama-modelfiles
```

Crea `Modelfile-gemma4-8gb`:

```Modelfile
FROM gemma4:e2b

# Hilos de CPU — mínimo que funciona (4 es seguro para la mayoría)
PARAMETER num_thread 4

# Ventana de contexto — reducida de 8192 a 2048 por defecto
# Esto ahorra ~600MB+ de RAM durante inferencia
PARAMETER num_ctx 2048

# Capas GPU — en 8GB de memoria unificada, deja espacio para SO + picoclaw
# gemma4:e2b tiene ~40 capas; 20 en GPU deja ~20 en CPU
PARAMETER num_gpu 20

# Tamaño de lote — lotes más pequeños = menos RAM a la vez
PARAMETER num_batch 128

# Mantener primeros 4 tokens (prompt del sistema) siempre en contexto
PARAMETER num_keep 4

# Opcional: prompt del sistema personalizado para picoclaw-agents
SYSTEM Eres PicoClaw, un asistente útil de IA. Sé conciso y orientado a la acción.
```

Construye y conecta:

```bash
# Construir el modelo (descarga gemma4:e2b si no está descargado)
ollama create picoclaw-gemma4-8gb -f Modelfile-gemma4-8gb

# Verificar
ollama list | grep gemma

# Probar
ollama run picoclaw-gemma4-8gb "Hola, cuánta RAM usas?"

# Configurar picoclaw-agents — agregar a ~/.picoclaw/config.json:
```

```json
{
  "model_list": [
    {
      "model_name": "picoclaw-gemma4-8gb",
      "model": "picoclaw-gemma4-8gb",
      "api_base": "http://localhost:11434/v1",
      "api_key": "ollama"
    }
  ],
  "agents": {
    "defaults": {
      "model_name": "picoclaw-gemma4-8gb",
      "max_tokens": 2048
    }
  }
}
```

**Uso esperado de RAM:** ~5-6GB total (modelo ~3GB + contexto + overhead)

#### Ejemplo 2: Qwen 3:8b — Configuración Mínima para Sistema de 16GB

Para un escritorio/laptop con 16GB de RAM, ejecutando Qwen 3:8b eficientemente:

```bash
cd ~/ollama-modelfiles
```

Crea `Modelfile-qwen3-16gb`:

```Modelfile
FROM qwen3:8b

# Hilos de CPU — iguala tu número de núcleos (6 es conservador)
PARAMETER num_thread 6

# Ventana de contexto — 4096 tokens, equilibrado para el bucle del agente
PARAMETER num_ctx 4096

# Capas GPU — qwen3:8b tiene ~35 capas; descarga 25 a GPU, 10 quedan en CPU
PARAMETER num_gpu 25

# Tamaño de lote — moderado para buen rendimiento sin picos de RAM
PARAMETER num_batch 256

# Mantener tokens del prompt del sistema
PARAMETER num_keep 4
```

Construye y conecta:

```bash
# Construir
ollama create picoclaw-qwen3-16gb -f Modelfile-qwen3-16gb

# Probar
ollama run picoclaw-qwen3-16gb "Escribe una función Python para invertir una cadena"

# Configurar picoclaw-agents:
```

```json
{
  "model_list": [
    {
      "model_name": "picoclaw-qwen3-16gb",
      "model": "picoclaw-qwen3-16gb",
      "api_base": "http://localhost:11434/v1",
      "api_key": "ollama"
    }
  ],
  "agents": {
    "defaults": {
      "model_name": "picoclaw-qwen3-16gb",
      "max_tokens": 4096
    }
  }
}
```

**Uso esperado de RAM/VRAM:** ~6-8GB total

#### Ejemplo 3: Qwen 2.5-Coder:0.5b — Mínimo Absoluto (Ultra-Baja RAM)

Para el menor consumo posible — Termux, Raspberry Pi, o cualquier sistema limitado. Esta es la **configuración mínima viable**:

```bash
cd ~/ollama-modelfiles
```

Crea `Modelfile-qwen-coder-minimal`:

```Modelfile
FROM qwen2.5-coder:0.5b

# Mínimo de hilos — 2 es el mínimo que funciona
PARAMETER num_thread 2

# Contexto mínimo — 512 tokens es el piso absoluto
# (por debajo de esto el modelo puede dar error o producir basura)
PARAMETER num_ctx 512

# Sin descarga a GPU — este modelo es tan pequeño que cabe en RAM de CPU
# num_gpu 0 asegura cero uso de VRAM
PARAMETER num_gpu 0

# Tamaño de lote mínimo — mínimo RAM durante procesamiento de prompts
PARAMETER num_batch 64

# Sin keep — ahorrar cada token de la pequeña ventana de contexto
PARAMETER num_keep 0
```

Construye y conecta:

```bash
# Construir (modelo ~400MB)
ollama create picoclaw-coder-minimal -f Modelfile-qwen-coder-minimal

# Probar — nota: respuestas cortas por la ventana de contexto tan pequeña
ollama run picoclaw-coder-minimal "def hola():"

# Configurar picoclaw-agents:
```

```json
{
  "model_list": [
    {
      "model_name": "picoclaw-coder-minimal",
      "model": "picoclaw-coder-minimal",
      "api_base": "http://localhost:11434/v1",
      "api_key": "ollama"
    }
  ],
  "agents": {
    "defaults": {
      "model_name": "picoclaw-coder-minimal",
      "max_tokens": 512,
      "max_tool_iterations": 5
    }
  }
}
```

**Uso esperado de RAM:** ~600MB total (modelo ~400MB + contexto + overhead)
**Uso de VRAM:** 0MB (solo CPU)

Esta configuración funciona en:
- Termux (Android, 2GB+ RAM)
- Raspberry Pi 4 (2GB RAM)
- Laptops antiguos (2GB+ RAM)
- Cualquier sistema donde necesites IA con mínimo consumo

#### Referencia Rápida: Rangos de Parámetros del Modelfile

| Parámetro | Mínimo Absoluto | Típico | Máximo |
|-----------|----------------|---------|--------|
| `num_thread` | 1 | 4-8 | Núcleos CPU |
| `num_ctx` | 512 | 2048-4096 | Máx del modelo (32K-128K) |
| `num_gpu` | 0 (solo CPU) | 20-35 | Todas las capas |
| `num_batch` | 32 | 128-512 | 2048 |
| `num_keep` | 0 | 4-8 | num_ctx / 2 |

#### Reconstruir Después de Cambios

Si editas el Modelfile, reconstruye el modelo:

```bash
# Editar el Modelfile
nano ~/ollama-modelfiles/Modelfile-gemma4-8gb

# Reconstruir (sobrescribe la versión anterior)
ollama create picoclaw-gemma4-8gb -f ~/ollama-modelfiles/Modelfile-gemma4-8gb

# El modelo anterior es reemplazado — no necesitas borrarlo primero
ollama list
```

### Método C: Vía Variables de Entorno (Nivel Servidor)

#### macOS — Usando launchd (App de Escritorio)

Si instalaste Ollama vía `.dmg`, se ejecuta como servicio launchd:

```bash
# Detener Ollama
launchctl stop com.ollama.ollama

# Editar el plist de launchd (si existe)
nano ~/Library/LaunchAgents/com.ollama.ollama.plist

# O usar variables de entorno al iniciar manualmente:
# Primero cierra la app de la barra de menú, luego:
OLLAMA_NUM_PARALLEL=2 OLLAMA_KEEP_ALIVE=1h ollama serve
```

**Usando la app de terminal en vez del servicio de fondo:**

```bash
# Primero cierra la app de la barra de menú
# Luego inicia con variables de entorno:
OLLAMA_KEEP_ALIVE=1h OLLAMA_NUM_PARALLEL=2 ollama serve
```

#### Linux — Usando systemd

```bash
# Editar el override del servicio systemd
systemctl edit ollama.service

# Agregar estas líneas en el editor:
[Service]
Environment="OLLAMA_NUM_PARALLEL=2"
Environment="OLLAMA_KEEP_ALIVE=1h"
Environment="OLLAMA_HOST=0.0.0.0:11434"

# Recargar y reiniciar
systemctl daemon-reload
systemctl restart ollama.service

# Verificar
systemctl status ollama.service
env | grep OLLAMA
```

#### Windows — Usando Variables de Entorno del Sistema

```powershell
# Método 1: Establecer permanentemente vía PowerShell (requiere reinicio)
[System.Environment]::SetEnvironmentVariable("OLLAMA_KEEP_ALIVE", "1h", "User")
[System.Environment]::SetEnvironmentVariable("OLLAMA_NUM_PARALLEL", "2", "User")

# Método 2: Solo para la sesión actual
$env:OLLAMA_KEEP_ALIVE = "1h"
$env:OLLAMA_NUM_PARALLEL = "2"

# Reiniciar servicio Ollama
Restart-Service Ollama

# O si se ejecuta manualmente, iniciar con variables:
$env:OLLAMA_KEEP_ALIVE = "1h"
ollama serve
```

#### Termux — Usando Config de Shell

```bash
# Agregar a ~/.bashrc o ~/.zshrc (persistente)
echo 'export OLLAMA_KEEP_ALIVE=30m' >> ~/.bashrc
echo 'export OLLAMA_NUM_PARALLEL=1' >> ~/.bashrc
source ~/.bashrc

# O establecer solo para la sesión actual
export OLLAMA_KEEP_ALIVE=30m
export OLLAMA_NUM_PARALLEL=1

# Iniciar Ollama
ollama serve
```

### Todos los Parámetros Oficiales de Recursos

| Parámetro | Tipo | Default | Descripción |
|-----------|------|---------|-------------|
| `num_thread` | int | Auto (núcleos CPU) | Hilos de CPU para inferencia |
| `num_ctx` | int | 2048 | Tamaño de ventana de contexto (tokens) |
| `num_gpu` | int | Todas las capas | Número de capas del modelo a descargar en GPU |
| `num_batch` | int | 512 | Tamaño de lote para procesamiento de prompts |
| `num_keep` | int | 0 | Tokens iniciales a mantener en contexto |
| `main_gpu` | int | 0 | GPU principal (configuraciones multi-GPU) |
| `use_mmap` | bool | true | Usar mapeo de memoria para cargar modelos |
| `numa` | bool | false | Habilitar optimización de memoria NUMA |

### Ejemplos por Plataforma para picoclaw-agents

#### Windows — Limitar VRAM de GPU

```powershell
# Modelfile: limitar a 20 capas en GPU (resto en CPU/RAM)
FROM gemma4:26b
PARAMETER num_gpu 20
PARAMETER num_ctx 4096
PARAMETER num_thread 8

# Construir y ejecutar
ollama create gemma4-limited -f Modelfile
ollama run gemma4-limited

# O vía /set durante sesión interactiva
ollama run gemma4:26b
/set parameter num_gpu 20
/set parameter num_ctx 4096
```

**Verificar uso de GPU en Windows:**
```powershell
# Administrador de tareas → Rendimiento → GPU → Memoria GPU Dedicada
# O vía PowerShell:
Get-Counter '\GPU Process Memory(*)\Local Usage'
```

#### macOS (Apple Silicon) — Limitar Memoria Unificada

```bash
# Modelfile: limitar hilos y contexto para Mac de 8GB
FROM qwen3:8b
PARAMETER num_thread 6
PARAMETER num_ctx 2048
PARAMETER num_gpu 30

# Construir y ejecutar
ollama create qwen3-lite -f Modelfile
ollama run qwen3-lite
```

**Verificar uso de memoria en macOS:**
```bash
# Monitorear memoria del proceso Ollama
ps aux | grep ollama | awk '{print $6/1024 " MB", $11}'

# O usar Monitor de Actividad → pestaña Memoria → filtrar "ollama"
```

#### Linux (GPU NVIDIA) — Limitar Capas GPU + RAM

```bash
# Modelfile: descarga parcial de GPU para VRAM limitada
FROM llama3.1:70b
PARAMETER num_gpu 35
PARAMETER num_ctx 4096
PARAMETER num_batch 256
PARAMETER num_thread 8

ollama create llama70b-limited -f Modelfile
ollama run llama70b-limited
```

**Verificar memoria GPU en Linux:**
```bash
# Memoria GPU NVIDIA
nvidia-smi --query-gpu=memory.used,memory.total --format=csv

# O ver en tiempo real
watch -n 1 nvidia-smi
```

#### Termux (Android) — Solo CPU, Mínimo Consumo

```bash
# Modelfile: ultra-ligero para móvil
FROM qwen3:0.6b
PARAMETER num_thread 4
PARAMETER num_ctx 1024
PARAMETER num_batch 128
PARAMETER num_gpu 0

ollama create qwen-mobile -f Modelfile
ollama run qwen-mobile
```

**Verificar memoria en Termux:**
```bash
# Memoria de procesos
ps -o pid,rss,comm | grep ollama | awk '{print $1, $2/1024 " MB", $3}'

# O usar htop
pkg install htop && htop
```

### Referencia Rápida: Valores por Hardware

| Hardware | `num_thread` | `num_ctx` | `num_gpu` | `num_batch` |
|----------|-------------|-----------|-----------|-------------|
| Termux (Android, 4GB) | 4 | 1024 | 0 | 128 |
| Laptop (8GB RAM, sin GPU) | 6 | 2048 | 0 | 256 |
| Mac M1 (8GB) | 6 | 2048 | 30 | 256 |
| Desktop (16GB + 8GB VRAM) | 8 | 4096 | 35 | 512 |
| Workstation (32GB + 24GB VRAM) | 12 | 8192 | auto | 512 |

> **Tip:** `num_gpu = 0` fuerza modo solo CPU. `num_gpu = auto` (u omitido) deja que Ollama decida según la VRAM disponible.

---

## 4. Conectar picoclaw-agents a Ollama

### Opción A — Editar `~/.picoclaw/config.json` directamente

Agrega entradas a tu `model_list`:

```json
{
  "model_list": [
    {
      "model_name": "llama3.2:1b",
      "model": "llama3.2:1b",
      "api_base": "http://localhost:11434/v1",
      "api_key": "ollama"  # pragma: allowlist secret
    },
    {
      "model_name": "qwen2.5:0.5b",
      "model": "qwen2.5:0.5b",
      "api_base": "http://localhost:11434/v1",
      "api_key": "ollama"  # pragma: allowlist secret
    },
    {
      "model_name": "qwen2.5-coder:0.5b",
      "model": "qwen2.5-coder:0.5b",
      "api_base": "http://localhost:11434/v1",
      "api_key": "ollama"  # pragma: allowlist secret
    }
  ]
}
```

> **Nota:** `api_key: "ollama"` `# pragma: allowlist secret` es requerido por el cliente OpenAI-compat aunque Ollama no use autenticación. Cualquier string no vacío funciona.

### Opción B — Web UI (picoclaw-agents-launcher)

1. Inicia el launcher: `picoclaw-agents-launcher`
2. Abre **http://localhost:18800/models**
3. Haz clic en **+ Add Model**
4. Completa el formulario:
   - **Name:** `llama3.2:1b`
   - **Model:** `llama3.2:1b`
   - **API Base:** `http://localhost:11434/v1`
   - **API Key:** `ollama`
5. Marca **Default Model** si quieres usarlo como predeterminado
6. Haz clic en **Save**

---

## 5. Ejecutar el Agente

```bash
# Mensaje único
picoclaw-agents agent --model llama3.2:1b -m "Hola, ¿estás corriendo localmente?"

# Modo interactivo
picoclaw-agents agent --model qwen2.5-coder:0.5b

# Con modelo predeterminado en config
picoclaw-agents agent -m "Escribe un script Python para listar archivos"
```

---

## 6. Verificar que Ollama está Corriendo

```bash
# Ver modelos instalados
ollama list

# Probar la API directamente
curl http://localhost:11434/v1/models

# Ollama corre en el puerto 11434 por defecto
# Iniciar manualmente si es necesario:
ollama serve
```

---

## Guía de Hardware

| RAM | Modelos Recomendados |
|-----|----------------------|
| 4GB | `qwen2.5:0.5b`, `llama3.2:1b` |
| 8GB | `llama3.2:3b`, `qwen2.5:3b`, `deepseek-coder:6.7b` |
| 16GB | `llama3.2:8b`, `qwen2.5:7b`, `mistral:7b` |
| 32GB+ | `llama3.1:70b` (cuantizado), `deepseek-r1:32b` |

- **Apple Silicon (M1/M2/M3):** Ollama usa aceleración GPU Metal automáticamente — los modelos corren más rápido que en hardware Intel equivalente
- **GPU NVIDIA:** CUDA se usa automáticamente si está disponible
- **Solo CPU:** Funciona bien para modelos de 1B–3B, más lento para modelos grandes

---

## Solución de Problemas

**Error `connection refused`:**
```bash
# Ollama no está corriendo — inícialo:
ollama serve
```

**Modelo no encontrado:**
```bash
# Descarga el modelo primero:
ollama pull llama3.2:1b
```

**Sin memoria (out of memory):**
```bash
# Usa un modelo más pequeño:
ollama pull qwen2.5:0.5b   # solo ~400MB
```

**Respuestas lentas:**
```bash
# Verifica qué backend GPU está activo:
ollama ps
```
