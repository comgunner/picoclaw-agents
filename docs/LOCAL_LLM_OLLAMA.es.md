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
