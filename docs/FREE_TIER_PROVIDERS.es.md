# Proveedores Free Tier — IA sin Costo con PicoClaw-Agents

**Última actualización:** 5 de abril de 2026
**Estado:** ✅ Todos los proveedores probados y operativos — sin tarjeta de crédito

---

## 🆓 Descripción General

PicoClaw-Agents soporta múltiples proveedores de LLM **100% gratuitos**. Sin tarjeta de crédito, sin plan de pago, sin expiración de prueba. Solo regístrate y comienza a usar modelos de IA potentes.

### Comparación Rápida

| Proveedor | Modelo Gratis | Costo | Registro | Ideal Para |
|-----------|--------------|------|----------|------------|
| **OpenRouter** | `openrouter/auto` | Gratis para siempre | [openrouter.ai](https://openrouter.ai/) | Auto-ruteo, simplicidad |
| **Zhipu AI** | `glm-4.5-flash` | Gratis para siempre | [z.ai](https://z.ai/) | Velocidad, tareas de código |
| **OpenAI** | `gpt-4.1-mini` | Free tier | [chatgpt.com](https://chatgpt.com/#settings/Security) | Calidad, razonamiento |
| **Qwen** | `qwen-plus` | Free tier | [dashscope.aliyun.com](https://dashscope.aliyun.com/) | Mejor calidad general ⭐ |

---

## 1. 🌐 OpenRouter Free (`openrouter-free`)

### Por Qué Usarlo
- **Cero configuración**: Auto-rutea al mejor modelo free disponible
- **Sin seleccionar modelo**: OpenRouter elige el óptimo automáticamente
- **Fallback integrado**: Si un modelo falla, cambia a otro

### Configuración

```bash
# Login interactivo
./picoclaw-agents auth login --provider openrouter-free

# O via wizard de onboard
./picoclaw-agents onboard --free
```

### Configuración Manual
```json
{
  "model_list": [
    {
      "model_name": "openrouter-free",
      "model": "openrouter/auto",
      "api_key": "sk-or-v1-..."
    }
  ]
}
```

### Modelos Gratuitos Disponibles
OpenRouter auto-rutea al mejor disponible, que puede incluir:
- `stepfun/step-3.5-flash` (256K contexto)
- `deepseek/deepseek-v3.2` (64K contexto)
- Variantes de Llama, Gemma y MiniMax

### Enlaces
- **Modelos Free:** https://openrouter.ai/collections/free-models
- **API Keys:** https://openrouter.ai/keys

---

## 2. 🧠 Zhipu AI (`zhipu`) — 100% Gratis Para Siempre

### Por Qué Usarlo
- **Completamente gratis**: Sin tarjeta de crédito, sin plan de pago
- **Inferencia rápida**: `glm-4.5-flash` optimizado para velocidad
- **Excelente para código**: Fuerte generación y comprensión de código
- **Sin límites de uso**: Free tier generoso para uso personal

### Configuración

```bash
# Login interactivo
./picoclaw-agents auth login --provider zhipu
```

El wizard:
1. Solicitará tu API key de Zhipu
2. Auto-configurará `glm-4.5-flash` como modelo default
3. Agregará todos los modelos GLM disponibles a tu config

### Modelos Disponibles

| Modelo | Contexto | Caso de Uso |
|--------|---------|-------------|
| `glm-4.5-flash` 🆓 | 128K | Default — rápido, gratis, capaz |
| `glm-4.7-flash` | 128K | Último modelo flash |
| `glm-5` | 128K | Modelo premium (puede requerir créditos) |
| `glm-5-turbo` | 128K | Optimizado para velocidad |
| `glm-4.5-air` | 128K | Variante ligera |

### Enlaces
- **Registro:** https://z.ai/
- **API Keys:** https://open.bigmodel.cn/usercenter/proj-mgmt/apikeys
- **Precios:** https://z.ai/pricing

---

## 3. 🔑 OpenAI (`openai`) — Free Tier via Device Code

### Por Qué Usarlo
- **Alta calidad**: Modelos GPT líderes en la industria
- **Free tier disponible**: Funciona con el plan gratuito de ChatGPT
- **Sin API key necesaria**: Usa flujo OAuth Device Code

### Configuración

```bash
# Device Code OAuth (método Browser NO soportado)
./picoclaw-agents auth login --provider openai --device-code
```

**Importante:** OpenAI solo soporta autenticación **Device Code**. OAuth por navegador NO está disponible.

### Cómo Funciona
1. Ejecuta el comando de arriba
2. Abre la URL mostrada en la terminal en cualquier dispositivo
3. Ingresa el código de 8 caracteres
4. Autoriza con tu cuenta de OpenAI/ChatGPT
5. Los modelos se agregan automáticamente a tu config

### Modelos Disponibles

| Modelo | ¿Gratis? | Contexto |
|--------|----------|----------|
| `gpt-4.1-mini` | ✅ Free tier | 32K |
| `gpt-4.1` | ✅ Free tier | 32K |
| `gpt-5` | ✅ Free tier | 32K |
| `o3-mini` | ✅ Free tier | 32K |
| `o3` | ✅ Free tier | 32K |
| `o1` | ✅ Free tier | 32K |

### Requisitos
- Habilitar autorización Device Code en [chatgpt.com/#settings/Security](https://chatgpt.com/#settings/Security)
- Cuenta gratuita de OpenAI/ChatGPT

---

## 4. ⭐ Qwen (`qwen`) — Mejor Free Tier en General

### Por Qué Usarlo
- **Máxima calidad**: Comparable a modelos de pago en benchmarks
- **Free tier generoso**: Límites de rate altos para uso personal
- **Multilingüe**: Excelente en inglés, chino y muchos otros idiomas
- **Contexto largo**: Hasta 128K tokens en algunos modelos

### ⚠️ Nota Importante
> **¿Cuánto durará el free tier de Qwen?** No lo sabemos. A abril de 2026, es la **mejor opción gratuita disponible** — muy recomendado usarlo mientras dure. No se ha anunciado fecha de finalización oficial.

### Configuración

```bash
# OAuth via navegador (recomendado)
./picoclaw-agents auth login --provider qwen

# O pega la API key directamente
./picoclaw-agents auth login --provider qwen --token
```

### Modelos Disponibles

| Modelo | ¿Gratis? | Contexto | Caso de Uso |
|--------|----------|----------|-------------|
| `qwen-plus` 🆓 | ✅ Free tier | 128K | **Default** — mejor balance |
| `qwen-max` | ✅ Free tier | 32K | Máxima calidad |
| `qwen-turbo` | ✅ Free tier | 128K | Inferencia más rápida |
| `qwen-long` | ✅ Free tier | 1M | Contexto ultra-largo |
| `qwen-vl-max` | ✅ Free tier | 32K | Visión/lenguaje |
| `qwen-vl-plus` | ✅ Free tier | 32K | Visión (más ligero) |

### Endpoints Regionales

| Región | URL |
|--------|-----|
| **US (Virginia)** | `https://dashscope-us.aliyuncs.com/compatible-mode/v1` |
| **Singapur** | `https://dashscope-sg.aliyuncs.com/compatible-mode/v1` |
| **China (Beijing)** | `https://dashscope.aliyuncs.com/compatible-mode/v1` |

### Enlaces
- **Registro:** https://dashscope.aliyun.com/
- **API Keys:** https://dashscope.console.aliyun.com/api-key
- **Modelos:** https://help.aliyun.com/zh/model-studio/getting-started/models

---

## 📊 Matriz de Recomendación

| Necesidad | Proveedor Recomendado |
|-----------|----------------------|
| **Cero setup, funciona ya** | `openrouter-free` |
| **Mejor calidad (úsalo mientras dure)** | `qwen` ⭐ |
| **Tareas de código, inferencia rápida** | `zhipu` |
| **Ecosistema OpenAI, razonamiento** | `openai --device-code` |
| **Documentos ultra-largos** | `qwen` (qwen-long: 1M contexto) |
| **Visión/comprensión de imágenes** | `qwen` (qwen-vl-max) |

---

## 🚀 Configuración Multi-Proveedor (Recomendado)

Para máxima confiabilidad, configura múltiples proveedores como respaldo:

```bash
# 1. Configura tu proveedor free primaria
./picoclaw-agents auth login --provider qwen

# 2. Agrega OpenRouter como respaldo
./picoclaw-agents auth login --provider openrouter-free

# 3. Agrega Zhipu como otro respaldo
./picoclaw-agents auth login --provider zhipu
```

Luego cambia entre ellos en la WebUI:
```
http://localhost:18800/credentials
```

O via CLI:
```bash
# Listar modelos disponibles
./picoclaw-agents models list

# Cambiar modelo/proveedor
./picoclaw-agents agent --model qwen-plus -m "Hola"
./picoclaw-agents agent --model glm-4.5-flash -m "Hola"
./picoclaw-agents agent --model openrouter/auto -m "Hola"
```

---

## ⚠️ Limitaciones Conocidas

| Proveedor | Limitación |
|-----------|------------|
| **OpenRouter Free** | El modelo puede cambiar sin aviso; límites de rate varían |
| **Zhipu** | `glm-5` puede requerir créditos de pago; usa `glm-4.5-flash` para gratis |
| **OpenAI** | Free tier tiene límites de uso; solo Device Code (sin Browser OAuth) |
| **Qwen** | Duración del free tier desconocida; puede cambiar la política en cualquier momento |

---

## 🔗 Documentación Relacionada

- **Guía OpenRouter:** [OPENROUTER_FREE.md](OPENROUTER_FREE.md) / [OPENROUTER_FREE.es.md](OPENROUTER_FREE.es.md)
- **Changelog:** [CHANGELOG.md](../CHANGELOG.md)
- **Fix Token Overflow:** [local_work/openrouter_free_token_fix.md](../local_work/openrouter_free_token_fix.md)
- **Referencia de Config:** [local_work/CONFIG_FIELD_REFERENCE.md](../local_work/CONFIG_FIELD_REFERENCE.md)

---

*FREE_TIER_PROVIDERS.es.md — Actualizado 5 de abril de 2026*
