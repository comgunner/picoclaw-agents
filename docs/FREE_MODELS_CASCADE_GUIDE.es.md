# Guia de Cascada de Modelos Gratuitos

> Guia de Cascada de Modelos Gratuitos documentation for PicoClaw-Agents.


---

## Que es Cascada de Fallback?

La cascada de fallback es un **mecanismo de supervivencia** para tu bot. Cuando el modelo gratuito primario falla (limite de velocidad, timeout, no disponible), el sistema automaticamente intenta el siguiente modelo gratuito en la cascada. Esto asegura que tu bot **se mantenga en linea 24/7** sin intervencion manual.

## Por que OpenRouter (Sin Tarjeta de Credito)

OpenRouter es el **unico proveedor** que ofrece:
- **Modelos gratuitos** sin requerir tarjeta de credito
- **Fallback automatico** entre multiples modelos gratuitos
- **API compatible con OpenAI** (reemplazo directo)
- **Sin vendor lock-in** (cambia modelos cuando quieras)

## Modelos Gratuitos Disponibles (Agosto 2026)

| Prioridad | Modelo | Proveedor | Caso de Uso |
|-----------|--------|-----------|-------------|
| 1 | `openrouter/free` | OpenRouter | Modelo gratuito primario (selecciona automaticamente el mejor disponible) |
| 2 | `nvidia/nemotron-3-ultra-550b-a55b:free` | NVIDIA | Contexto grande, razonamiento complejo |
| 3 | `openai/gpt-oss-20b:free` | OpenAI | Rendimiento equilibrado |
| 4 | `inclusionai/ling-3.0-tiny:free` | InclusionAI | Respuestas rapidas, baja latencia |
| 5 | `meta-llama/llama-3.1-8b-instruct:free` | Meta proposito general |

## Modelos de Fallback de Pago

Si todos los modelos gratuitos fallan, la cascada puede recurrir a modelos de pago:

| Modelo | Proveedor | Costo | Caso de Uso |
|--------|-----------|-------|-------------|
| `deepseek/deepseek-v4-flash-0731` | DeepSeek | $0.27/M tokens | Alta calidad, costo-efectivo |
| `xiaomi/mimo-v2.5` | Xiaomi | $0.10/M tokens | Opcion economica |

## Configuracion

### Configuracion Basica (Solo Gratuito)

```json
{
  "agents": {
    "defaults": {
      "provider": "openrouter",
      "model_name": "openrouter-free",
      "model": "openrouter/free",
      "cascade": {
        "enabled": true,
        "text_models": [
          "openrouter/free",
          "nvidia/nemotron-3-ultra-550b-a55b:free",
          "openai/gpt-oss-20b:free",
          "inclusionai/ling-3.0-tiny:free",
          "meta-llama/llama-3.1-8b-instruct:free"
        ]
      }
    }
  }
}
```

### Con Fallback de Pago

```json
{
  "agents": {
    "defaults": {
      "provider": "openrouter",
      "model_name": "openrouter-free",
      "model": "openrouter/free",
      "cascade": {
        "enabled": true,
        "text_models": [
          "openrouter/free",
          "nvidia/nemotron-3-ultra-550b-a55b:free",
          "openai/gpt-oss-20b:free",
          "inclusionai/ling-3.0-tiny:free",
          "meta-llama/llama-3.1-8b-instruct:free"
        ],
        "paid_fallbacks": [
          "deepseek/deepseek-v4-flash-0731",
          "xiaomi/mimo-v2.5"
        ]
      }
    }
  }
}
```

## Como Funciona la Cascada

```
Solicitud del Usuario
    ↓
Intentar Modelo 1: openrouter/free
    ↓ (si falla)
Intentar Modelo 2: nvidia/nemotron-3-ultra-550b-a55b:free
    ↓ (si falla)
Intentar Modelo 3: openai/gpt-oss-20b:free
    ↓ (si falla)
Intentar Modelo 4: inclusionai/ling-3.0-tiny:free
    ↓ (si falla)
Intentar Modelo 5: meta-llama/llama-3.1-8b-instruct:free
    ↓ (si falla)
Intentar Fallback de Pago: deepseek/deepseek-v4-flash-0731
    ↓ (si falla)
Intentar Fallback de Pago: xiaomi/mimo-v2.5
    ↓
Retornar error si todos fallan
```

## Importancia para el Uptime del Bot

### Sin Cascada
- Modelo unico falla → bot se desconecta
- Se requiere intervencion manual
- Experiencia de usuario degradada

### Con Cascada
- Failover automatico entre modelos
- Bot se mantiene en linea 24/7
- Sin intervencion manual necesaria
- Experiencia de usuario fluida

## Cascada de Generacion de Imagenes

El mismo patron de cascada aplica para generacion de imagenes:

```json
{
  "image_gen": {
    "provider": "openrouter",
    "openrouter_image_model": "krea/krea-2-medium-turbo",
    "openrouter_text_model": "openrouter/free"
  }
}
```

### Modelos de Imagen
| Prioridad | Modelo | Caso de Uso |
|-----------|--------|-------------|
| 1 | `krea/krea-2-medium-turbo` | Primario (mas barato) |
| 2 | `bytedance-seed/seedream-4.5` | Alternativo |
| 3 | `x-ai/grok-imagine-image-quality` | Alta calidad |

## Monitoreo

### Verificar Estado de Cascada
```bash
# Listar modelos disponibles
picoclaw models

# Probar modelo especifico
picoclaw chat --model openrouter/free "Hola"

# Verificar estado de OpenRouter
curl http://localhost:18792/health
```

### Logs
```bash
# Ver logs de cascada en tiempo real
journalctl -u picoclaw.service -f | grep cascade

# Verificar fallos de modelos
journalctl -u picoclaw.service | grep "model.*failed"
```

## Solucion de Problemas

### Todos los Modelos Gratuitos Fallan
1. Verificar estado de OpenRouter: https://status.openrouter.ai
2. Verificar API key: `picoclaw auth login`
3. Verificar limites de velocidad: Los modelos gratuitos tienen limites por minuto
4. Considerar agregar fallbacks de pago

### Alta Latencia
1. La cascada intenta multiples modelos secuencialmente
2. La primera respuesta exitosa se retorna
3. Considerar reducir la lista de cascada para respuesta mas rapida

### Limite de Velocidad
1. Los modelos gratuitos tienen limites de velocidad
2. La cascada automaticamente intenta el siguiente modelo
3. Esperar reinicio del limite o usar fallback de pago

## Mejores Practicas

1. **Siempre habilitar cascada** para bots en produccion
2. **Comenzar con modelos gratuitos** para minimizar costo
3. **Agregar fallbacks de pago** para aplicaciones criticas
4. **Monitorear logs** para verificar fallos de modelos
5. **Probar cascada** regularmente con `picoclaw chat`
6. **Mantener modelos actualizados** a medida que los proveedores agregan/remueven niveles gratuitos

## Documentacion Relacionada

- [Modelos Gratuitos OpenRouter](OPENROUTER_FREE.es.md)
- [Guia de Configuracion](CONFIGURATION.md)
- [Configuracion de Proveedores](PROVIDERS.md)
- [Solucion de Problemas](TROUBLESHOOTING.md)
