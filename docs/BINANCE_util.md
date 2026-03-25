# Binance Util

Quick guide to use Binance tools in PicoClaw from terminal and Telegram.

> **PicoClaw v3.4.1**: Now supports **Fast-path Slash Commands** for instant trading operations and **Global Tracker** for multi-agent consistency.

## Requirements

Set your credentials in `~/.picoclaw/config.json`:

```json
{
  "tools": {
    "binance": {
      "api_key": "YOUR_BINANCE_API_KEY",
      "secret_key": "YOUR_BINANCE_SECRET_KEY"
    }
  }
}
```

You can also use environment variables:

```bash
export BINANCE_API_KEY="YOUR_BINANCE_API_KEY"
export BINANCE_SECRET_KEY="YOUR_BINANCE_SECRET_KEY"
```

## Available Tools

- `binance_get_ticker_price` (public, no API required)
- `binance_get_order_book` (public, order book depth with bids/asks)
- `binance_get_futures_order_book` (public, USDT-M futures order book depth)
- `binance_list_futures_volume` (public, 24h futures volume ranking)
- `binance_get_spot_balance` (requires API key/secret)
- `binance_get_futures_balance` (requires API key/secret)
- `binance_open_futures_position` (requires API key/secret + `confirm: true`)
- `binance_close_futures_position` (requires API key/secret + `confirm: true`)

## Natural Agent Interaction

If Binance MCP is connected (or you are using PicoClaw native Binance tools), you can ask naturally without crafting HTTP/WebSocket calls:

```text
Get the order book for BTCUSDT.
Check market depth and show bids and asks for BTCUSDT.
```

The agent should resolve this using `binance_get_order_book` / `get_order_book`.

You can also use short phrases:

```text
order book BTCUSDT
futures order book BTCUSDT
list future volume
spot balance
futures balance
```

## Terminal Usage

### Simple Shortcuts (recommended)

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

### Open LONG

```bash
./picoclaw-agents agent -m "Use binance_open_futures_position with symbol BTCUSDT, side LONG, quantity 0.001, leverage 5 and confirm true."
```

### Close Full Position

```bash
./picoclaw-agents agent -m "Use binance_close_futures_position with symbol BTCUSDT and confirm true."
```

### Partial Close

```bash
./picoclaw-agents agent -m "Use binance_close_futures_position with symbol BTCUSDT, quantity 0.0005 and confirm true."
```

## Telegram Usage

With `picoclaw-agents gateway` running, send these messages to your bot:

### Simple Shortcuts (recommended)

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

### Open LONG

```text
Use binance_open_futures_position with symbol BTCUSDT, side LONG, quantity 0.001, leverage 5 and confirm true.
```

### Close Full Position

```text
Use binance_close_futures_position with symbol BTCUSDT and confirm true.
```

### Partial Close

```text
Use binance_close_futures_position with symbol BTCUSDT, quantity 0.0005 and confirm true.
```

## Useful Queries

### Public Price (No API)

```bash
./picoclaw-agents agent -m "Use binance_get_ticker_price with symbol ETHUSDT and return only the numeric price."
```

### Order Book / Market Depth (No API)

```bash
./picoclaw-agents agent -m "Use binance_get_order_book with symbol BTCUSDT and limit 10."
```

### Futures Order Book / Market Depth (No API)

```bash
./picoclaw-agents agent -m "Use binance_get_futures_order_book with symbol BTCUSDT and limit 10."
```

### Futures Volume Ranking (No API)

```bash
./picoclaw-agents agent -m "Use binance_list_futures_volume with top 10."
./picoclaw-agents agent -m "list future volume"
./picoclaw-agents agent -m "list future volume top 20"
```

### Spot Balance

```bash
./picoclaw-agents agent -m "Use binance_get_spot_balance and show my non-zero balances."
```

### Futures Balance

```bash
./picoclaw-agents agent -m "Use binance_get_futures_balance and show my non-zero futures balances."
```

### Top 10 Futures Symbols (Order Book)

Example with commonly liquid USDT-M futures symbols:

`BTCUSDT`, `ETHUSDT`, `BNBUSDT`, `SOLUSDT`, `XRPUSDT`, `DOGEUSDT`, `ADAUSDT`, `AVAXUSDT`, `LINKUSDT`, `LTCUSDT`

Quick terminal loop:

```bash
for s in BTCUSDT ETHUSDT BNBUSDT SOLUSDT XRPUSDT DOGEUSDT ADAUSDT AVAXUSDT LINKUSDT LTCUSDT; do
  ./picoclaw-agents agent -m "futures order book ${s} top 10"
done
```

Single prompt alternative:

```bash
./picoclaw-agents agent -m "Show futures top 10 order book for BTCUSDT, ETHUSDT, BNBUSDT, SOLUSDT, XRPUSDT, DOGEUSDT, ADAUSDT, AVAXUSDT, LINKUSDT, and LTCUSDT."
```

## Safety Notes

- Real trading orders require `confirm true`.
- Verify `symbol`, `quantity`, and `leverage` before execution.
- If API keys are missing, trading operations are blocked.

---

## ⚡ Fast-path Slash Commands Status

**Important:** As of v3.4.3, Binance-specific fast-path commands (`/binance_open`, `/binance_price`, etc.) are **documented for future implementation** but not yet available.

### Current Available Fast-paths

These fast-path commands **are** available:

```text
/status              - Show system status
/help                - Show interactive help
/show model          - Show active model
/show channel        - Show communication channel
/list models         - List configured models
/list channels       - List configured channels
#BATCH_ID            - Query batch task status (e.g., #IMA_GEN_02_03_26_1500)
/bundle_approve      - Approve and publish bundle
/bundle_regen        - Regenerate bundle
/bundle_edit         - Edit bundle text
```

### Binance Trading (Tool-Based)

For Binance operations, use tool functions directly through natural language:

**Telegram:**
```text
Get the price of Bitcoin
Show me the order book for ETHUSDT
Open a long position on BTCUSDT with 0.001 and 5x leverage
Check my futures balance
```

**Discord:**
```text
What's the current price of BTC?
Show order book for Ethereum with 10 levels
Close my BTCUSDT position
```

**CLI:**
```bash
./picoclaw-agents agent -m "Get BTCUSDT price"
./picoclaw-agents agent -m "Show order book for ETHUSDT limit 10"
./picoclaw-agents agent -m "Open long position BTCUSDT 0.001 leverage 5 confirm true"
```

**Benefits:**
- ✅ **Full functionality**: All 8 Binance tools available
- ✅ **Natural language**: No need to memorize command syntax
- ✅ **Safe by default**: Requires `confirm=true` for real orders
- ✅ **Context-aware**: LLM can clarify ambiguous requests

### Future Fast-path Implementation

Planned Binance fast-path commands (not yet available):
```text
/binance_price BTCUSDT        - Quick price check
/binance_orderbook BTCUSDT 10 - Order book depth
/binance_open BTCUSDT long 0.001 leverage 5 - Open position
/binance_close BTCUSDT        - Close position
/binance_balance futures      - Check balance
```

Watch for updates in future releases.

---

## 🌐 Global Tracker (v3.4.1+)

The **Global ImageGenTracker** is now shared across all agents (PM, Subagents), ensuring perfect consistency in multi-agent workflows:

- **Subagent generates image** → **Main Agent can immediately publish**
- **No "ID not found" errors** across agent boundaries
- **Shared state** for all Binance operations and social media posts

See [queue_batch.md](queue_batch.md) for complete tracker documentation.

---

## 📢 Post Binance Data to Social Media

You can combine Binance tools with Social Media tools to automatically publish market data to your social networks.

### Flow: Query Binance → Post to Social Media

```bash
# 1. Query Binance data
# 2. Post result to Facebook, Twitter, Discord
```

### Terminal Examples

```bash
# Futures order book + post to Twitter
./picoclaw-agents agent -m "futures order book BTCUSDT top 10 and post result to Twitter"

# Futures volume + post to Discord
./picoclaw-agents agent -m "list future volume and post top 5 to Discord"

# Order book + post to Facebook
./picoclaw-agents agent -m "futures order book ETHUSDT top 10 and post to Facebook"

# Multiple symbols + post summary
./picoclaw-agents agent -m "Show order book for BTCUSDT, ETHUSDT, BNBUSDT and post the most liquid to Twitter"
```

### Telegram Examples

```text
# Order book and post
futures order book BTCUSDT top 10 and post to Twitter

# Volume and post
list future volume and post top 5 to Discord

# Combined
futures order book BTCUSDT top 10 and post to Facebook and Twitter
```

### Discord Examples

```text
# Direct commands to bot
futures order book BTCUSDT top 10 and post to Twitter
list future volume and post top 5 to Discord
order book ETHUSDT and share on Facebook
```

### Automated Workflow with Community Manager

```bash
# Generate engaging post from Binance data
./picoclaw-agents agent -m "
  Use binance_get_futures_order_book with symbol BTCUSDT limit 10,
  then use community_manager_create_draft with raw_data from result, platform='twitter',
  then publish the generated draft
"

# Complete flow: data → engaging post → publish
./picoclaw-agents agent -m "
  Get BTCUSDT order book,
  create engaging post with community_manager for Twitter,
  publish with hashtags #BTC #Binance #Trading
"
```

### Automatic Post Examples

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
🚀 Top 5 Futures by Volume (24h)
1️⃣ BTCUSDT: $1.2B
2️⃣ ETHUSDT: $890M
3️⃣ BNBUSDT: $456M
4️⃣ SOLUSDT: $234M
5️⃣ XRPUSDT: $198M
```

**Facebook:**
```text
📈 Market Update - Binance Futures

BTCUSDT order book shows:
- Best Bid: $95,432.50
- Best Ask: $95,435.00
- Spread: $2.50 (0.003%)

Liquid and stable market. Good time to trade!

#Trading #Binance #Bitcoin
```

### Monitoring and Posting Loop

```bash
# Monitor every 5 minutes and post to Discord
while true; do
  ./picoclaw-agents agent -m "futures order book BTCUSDT top 10 and post to Discord if spread < $5"
  sleep 300
done

# Post top volume every hour
while true; do
  ./picoclaw-agents agent -m "list future volume top 10 and post to Twitter"
  sleep 3600
done
```

### Combine with Image Generation

```bash
# Generate image with Binance data and post
./picoclaw-agents agent -m "
  Get BTCUSDT order book,
  generate chart image using image_gen_create,
  create post with community_manager,
  post to Facebook and Twitter
"
```

---

## 🔗 Related Documentation

- **Social Media:** See `SOCIAL_MEDIA.md` for Facebook, Twitter, Discord setup
- **Image Generation:** See `docs/IMAGE_GEN_util.md` for generating images from data
- **Community Manager:** See tools `community_manager_create_draft` and `community_from_image`
