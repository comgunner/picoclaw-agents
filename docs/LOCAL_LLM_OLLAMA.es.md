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

### Top 10 Modelos de la Librería de Ollama

Visita **https://ollama.com/library** para el catálogo completo. Los modelos más populares:

| # | Modelo | Destacado |
|---|--------|-----------|
| 1 | `gemma3` | Google Gemma 3, multilingüe |
| 2 | `llama3.2` | Meta Llama 3.2, rápido y capaz |
| 3 | `qwen2.5` | Alibaba Qwen 2.5, excelente en código |
| 4 | `phi4` | Microsoft Phi-4, pequeño pero inteligente |
| 5 | `mistral` | Mistral 7B, buen razonamiento |
| 6 | `deepseek-r1` | DeepSeek R1, chain-of-thought |
| 7 | `llava` | Multimodal (texto + imágenes) |
| 8 | `codellama` | Meta CodeLlama, enfocado en código |
| 9 | `deepseek-coder-v2` | DeepSeek Coder V2, top para programación |
| 10 | `nomic-embed-text` | Embeddings de texto, casos de uso RAG |

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
