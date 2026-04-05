# Binance Util

Guía rápida para usar herramientas de Binance en PicoClaw desde terminal y Telegram.

> **PicoClaw v3.4.1**: Ahora soporta **Comandos Slash Fast-path** para operaciones de trading instantáneas y **Global Tracker** para consistencia multi-agente.

## Requisitos

Configura tus credenciales en `~/.picoclaw/config.json`:

```json
{
  "tools": {
    "binance": {
      "api_key": "TU_BINANCE_API_KEY",
      "secret_key": "TU_BINANCE_SECRET_KEY"
    }
  }
}
```

También puedes usar variables de entorno:

```bash
export BINANCE_API_KEY="TU_BINANCE_API_KEY"
export BINANCE_SECRET_KEY="TU_BINANCE_SECRET_KEY"
```

## Herramientas disponibles

- `binance_get_ticker_price` (público, no requiere API)
- `binance_get_order_book` (público, libro de órdenes con bids/asks)
- `binance_get_futures_order_book` (público, libro de órdenes de futuros USDT-M)
- `binance_list_futures_volume` (público, ranking por volumen 24h en futuros)
- `binance_get_spot_balance` (requiere API/secret)
- `binance_get_futures_balance` (requiere API/secret)
- `binance_open_futures_position` (requiere API/secret + `confirm: true`)
- `binance_close_futures_position` (requiere API/secret + `confirm: true`)

## Interacción natural con el agente

Si el servidor MCP de Binance está conectado o si usas las tools nativas de PicoClaw, puedes pedirlo en lenguaje natural, por ejemplo:

```text
Obtén el libro de órdenes para BTCUSDT.
Revisa la profundidad del mercado para ver los bids y asks de BTCUSDT.
```

El agente debe resolverlo usando `binance_get_order_book` / `get_order_book`.

También puedes usar frases cortas:

```text
order book BTCUSDT
futures order book BTCUSDT
list future volume
spot balance
futures balance
```

## Uso desde terminal

### Atajos simples (recomendado)

```bash
./picoclaw-agents agent -m "open futures BTCUSDT long 0.001 leverage 5"
./picoclaw-agents agent -m "close futures BTCUSDT"
./picoclaw-agents agent -m "close futures partial BTCUSDT 0.0005"
./picoclaw-agents agent -m "order book BTCUSDT top 10"
./picoclaw-agents agent -m "futures order book BTCUSDT top 10"
./picoclaw-agents agent -m "list future volume"
./picoclaw-agents agent -m "spot balance"
./picoclaw-agents agent -m "futures balance"
```

### Abrir LONG

```bash
./picoclaw-agents agent -m "Usa binance_open_futures_position con symbol BTCUSDT, side LONG, quantity 0.001, leverage 5 y confirm true."
```

### Cerrar posición completa

```bash
./picoclaw-agents agent -m "Usa binance_close_futures_position con symbol BTCUSDT y confirm true."
```

### Cerrar parcial

```bash
./picoclaw-agents agent -m "Usa binance_close_futures_position con symbol BTCUSDT, quantity 0.0005 y confirm true."
```

## Uso desde Telegram

Con `picoclaw-agents gateway` activo, envía estos mensajes al bot:

### Atajos simples (recomendado)

```text
open futures BTCUSDT long 0.001 leverage 5
close futures BTCUSDT
close futures partial BTCUSDT 0.0005
order book BTCUSDT top 10
futures order book BTCUSDT top 10
list future volume
spot balance
futures balance
```

### Abrir LONG

```text
Usa binance_open_futures_position con symbol BTCUSDT, side LONG, quantity 0.001, leverage 5 y confirm true.
```

### Cerrar posición completa

```text
Usa binance_close_futures_position con symbol BTCUSDT y confirm true.
```

### Cerrar parcial

```text
Usa binance_close_futures_position con symbol BTCUSDT, quantity 0.0005 y confirm true.
```

## Consultas útiles

### Precio público (sin API)

```bash
./picoclaw-agents agent -m "Usa binance_get_ticker_price con symbol ETHUSDT y responde solo con el precio."
```

### Libro de órdenes / profundidad (sin API)

```bash
./picoclaw-agents agent -m "Usa binance_get_order_book con symbol BTCUSDT y limit 10."
```

### Libro de órdenes de futuros / profundidad (sin API)

```bash
./picoclaw-agents agent -m "Usa binance_get_futures_order_book con symbol BTCUSDT y limit 10."
```

### Ranking de futuros por volumen (sin API)

```bash
./picoclaw-agents agent -m "Usa binance_list_futures_volume con top 10."
./picoclaw-agents agent -m "list future volume"
./picoclaw-agents agent -m "list future volume top 20"
```

### Balance spot

```bash
./picoclaw-agents agent -m "Usa binance_get_spot_balance y muestra mis balances no cero."
```

### Balance de futuros

```bash
./picoclaw-agents agent -m "Usa binance_get_futures_balance y muestra mis balances de futuros no cero."
```

### Top 10 símbolos de futuros (libro de órdenes)

Ejemplo con símbolos líquidos frecuentes de futuros USDT-M:

`BTCUSDT`, `ETHUSDT`, `BNBUSDT`, `SOLUSDT`, `XRPUSDT`, `DOGEUSDT`, `ADAUSDT`, `AVAXUSDT`, `LINKUSDT`, `LTCUSDT`

Consulta rápida en terminal:

```bash
for s in BTCUSDT ETHUSDT BNBUSDT SOLUSDT XRPUSDT DOGEUSDT ADAUSDT AVAXUSDT LINKUSDT LTCUSDT; do
  ./picoclaw-agents agent -m "futures order book ${s} top 10"
done
```

También lo puedes pedir en un solo prompt:

```bash
./picoclaw-agents agent -m "Muestra el order book top 10 de futuros para BTCUSDT, ETHUSDT, BNBUSDT, SOLUSDT, XRPUSDT, DOGEUSDT, ADAUSDT, AVAXUSDT, LINKUSDT y LTCUSDT."
```

## Notas de seguridad

- Las órdenes reales de trading requieren `confirm true`.
- Verifica `symbol`, `quantity` y `leverage` antes de ejecutar.
- Si faltan API keys, las operaciones de trading quedan bloqueadas.

---

## ⚡ Estado de Comandos Slash Fast-path

**Importante:** A partir de v3.4.3, los comandos fast-path específicos de Binance (`/binance_open`, `/binance_price`, etc.) están **documentados para futura implementación** pero aún no están disponibles.

### Fast-paths Disponibles Actualmente

Estos comandos fast-path **sí** están disponibles:

```text
/status              - Mostrar estado del sistema
/help                - Mostrar ayuda interactiva
/show model          - Mostrar modelo activo
/show channel        - Mostrar canal de comunicación
/list models         - Listar modelos configurados
/list channels       - Listar canales configurados
#BATCH_ID            - Consultar estado de tarea (ej: #IMA_GEN_02_03_26_1500)
/bundle_approve      - Aprobar y publicar bundle
/bundle_regen        - Regenerar bundle
/bundle_edit         - Editar texto del bundle
```

### Trading Binance (Basado en Herramientas)

Para operaciones de Binance, usa funciones de herramienta directamente con lenguaje natural:

**Telegram:**
```text
Obtén el precio de Bitcoin
Muéstrame el order book de ETHUSDT
Abre una posición long en BTCUSDT con 0.001 y 5x leverage
Consulta mi balance de futures
```

**Discord:**
```text
Cuál es el precio actual de BTC?
Muestra order book de Ethereum con 10 niveles
Cierra mi posición de BTCUSDT
```

**CLI:**
```bash
./picoclaw-agents agent -m "Get BTCUSDT price"
./picoclaw-agents agent -m "Show order book for ETHUSDT limit 10"
./picoclaw-agents agent -m "Open long position BTCUSDT 0.001 leverage 5 confirm true"
```

**Beneficios:**
- ✅ **Funcionalidad completa**: 8 herramientas de Binance disponibles
- ✅ **Lenguaje natural**: Sin necesidad de memorizar sintaxis de comandos
- ✅ **Seguro por defecto**: Requiere `confirm=true` para órdenes reales
- ✅ **Consciente del contexto**: El LLM puede aclarar solicitudes ambiguas

### Futura Implementación de Fast-paths

Comandos fast-path de Binance planeados (aún no disponibles):
```text
/binance_price BTCUSDT        - Consulta rápida de precio
/binance_orderbook BTCUSDT 10 - Profundidad de order book
/binance_open BTCUSDT long 0.001 leverage 5 - Abrir posición
/binance_close BTCUSDT        - Cerrar posición
/binance_balance futures      - Consultar balance
```

Estarán disponibles en futuras versiones.

---

## 🌐 Global Tracker (v3.4.1+)

El **Global ImageGenTracker** ahora es compartido entre todos los agentes (PM, Subagentes), asegurando consistencia perfecta en flujos de trabajo multi-agente:

- **Subagente genera imagen** → **Agente Principal puede publicar inmediatamente**
- **Sin errores de "ID no encontrado"** entre límites de agentes
- **Estado compartido** para todas las operaciones de Binance y posts de redes sociales

Ver [queue_batch.es.md](queue_batch.es.md) para documentación completa del tracker.

---

## 📢 Publicar Datos de Binance en Redes Sociales

Puedes combinar las herramientas de Binance con las herramientas de Social Media para publicar automáticamente datos de mercado en tus redes sociales.

### Flujo: Consultar Binance → Publicar en Redes

```bash
# 1. Consultar datos de Binance
# 2. Publicar resultado en Facebook, Twitter, Discord
```

### Ejemplos en Terminal

```bash
# Order book de futuros + publicar en Twitter
./picoclaw-agents agent -m "futures order book BTCUSDT top 10 y publica el resultado en Twitter"

# Volumen de futuros + publicar en Discord
./picoclaw-agents agent -m "list future volume y publica el top 5 en Discord"

# Order book + publicar en Facebook
./picoclaw-agents agent -m "futures order book ETHUSDT top 10 y publica en Facebook"

# Múltiples símbolos + publicar resumen
./picoclaw-agents agent -m "Muestra order book de BTCUSDT, ETHUSDT, BNBUSDT y publica el más líquido en Twitter"
```

### Ejemplos en Telegram

```text
# Order book y publicar
futures order book BTCUSDT top 10 y publica en Twitter

# Volumen y publicar
list future volume y publica el top 5 en Discord

# Combinado
futures order book BTCUSDT top 10 y publica en Facebook y Twitter
```

### Ejemplos en Discord

```text
# Comandos directos al bot
futures order book BTCUSDT top 10 and post to Twitter
list future volume and post top 5 to Discord
order book ETHUSDT and share on Facebook
```

### Workflow Automatizado con Community Manager

```bash
# Generar post atractivo desde datos de Binance
./picoclaw-agents agent -m "
  Usa binance_get_futures_order_book con symbol BTCUSDT limit 10,
  luego usa community_manager_create_draft con raw_data del resultado, platform='twitter',
  luego publica el borrador generado
"

# Flujo completo: datos → post atractivo → publicar
./picoclaw-agents agent -m "
  Obtené order book de BTCUSDT,
  creá un post atractivo con community_manager para Twitter,
  publicalo con hashtags #BTC #Binance #Trading
"
```

### Ejemplos de Posts Automáticos

**Twitter:**
```text
📊 BTCUSDT Order Book - Top 10
💰 Best Bid: $95,432.50
💰 Best Ask: $95,435.00
📈 Spread: $2.50
#BTC #Binance #Trading
```

**Discord:**
```text
🚀 Top 5 Futuros por Volumen (24h)
1️⃣ BTCUSDT: $1.2B
2️⃣ ETHUSDT: $890M
3️⃣ BNBUSDT: $456M
4️⃣ SOLUSDT: $234M
5️⃣ XRPUSDT: $198M
```

**Facebook:**
```text
📈 Actualización de Mercado - Binance Futures

El order book de BTCUSDT muestra:
- Mejor oferta: $95,432.50
- Mejor demanda: $95,435.00
- Spread: $2.50 (0.003%)

Mercado líquido y estable. ¡Buen momento para operar!

#Trading #Binance #Bitcoin
```

### Loop de Monitoreo y Publicación

```bash
# Monitorear cada 5 minutos y publicar en Discord
while true; do
  ./picoclaw-agents agent -m "futures order book BTCUSDT top 10 y publica en Discord si el spread es < $5"
  sleep 300
done

# Publicar top volumen cada hora
while true; do
  ./picoclaw-agents agent -m "list future volume top 10 y publica en Twitter"
  sleep 3600
done
```

### Combinar con Generación de Imágenes

```bash
# Generar imagen con datos de Binance y publicar
./picoclaw-agents agent -m "
  Obtené order book de BTCUSDT,
  generá una imagen con gráfico usando image_gen_create,
  creá post con community_manager,
  publicá en Facebook y Twitter
"
```

---

## 🔗 Documentación Relacionada

- **Social Media:** Ver `SOCIAL_MEDIA.es.md` para configuración de Facebook, Twitter, Discord
- **Image Generation:** Ver `docs/IMAGE_GEN_util.es.md` para generar imágenes desde datos
- **Community Manager:** Ver herramientas `community_manager_create_draft` y `community_from_image`
