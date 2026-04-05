// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)

package skills

import (
	"strings"
)

// BinanceMCPSkill implements native skill for Binance MCP connection.
// It teaches the LLM how to connect to Binance MCP server and use trading tools.
type BinanceMCPSkill struct {
	workspace string
}

// NewBinanceMCPSkill creates a new BinanceMCPSkill instance.
func NewBinanceMCPSkill(workspace string) *BinanceMCPSkill {
	return &BinanceMCPSkill{
		workspace: workspace,
	}
}

// Name returns the skill identifier name.
func (b *BinanceMCPSkill) Name() string {
	return "binance_mcp"
}

// Description returns a brief description of the skill.
func (b *BinanceMCPSkill) Description() string {
	return "Connect to Binance MCP server for trading operations. Supports public data (no API) and private trading (with API credentials)."
}

// GetInstructions returns the complete usage instructions for the LLM.
func (b *BinanceMCPSkill) GetInstructions() string {
	return binanceMCPInstructions
}

// GetAntiPatterns returns common anti-patterns to avoid.
func (b *BinanceMCPSkill) GetAntiPatterns() string {
	return binanceMCPAntiPatterns
}

// GetExamples returns concrete usage examples.
func (b *BinanceMCPSkill) GetExamples() string {
	return binanceMCPExamples
}

// BuildSkillContext returns the complete skill context for prompt injection.
func (b *BinanceMCPSkill) BuildSkillContext() string {
	var parts []string

	parts = append(parts, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	parts = append(parts, "🚀 NATIVE SKILL: Binance MCP Connection")
	parts = append(parts, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	parts = append(parts, "")
	parts = append(parts, "**PURPOSE:** Connect to Binance MCP server for trading operations and market data.")
	parts = append(parts, "")
	parts = append(parts, b.GetInstructions())
	parts = append(parts, "")
	parts = append(parts, b.GetAntiPatterns())
	parts = append(parts, "")
	parts = append(parts, b.GetExamples())

	return strings.Join(parts, "\n")
}

// BuildSummary returns an XML summary for compact context injection.
func (b *BinanceMCPSkill) BuildSummary() string {
	return `<skill name="binance_mcp" type="native">
  <purpose>Connect to Binance MCP server for trading</purpose>
  <pattern>Use for market data and trading operations</pattern>
  <tools>get_ticker_price, get_order_book, open_futures_position, etc.</tools>
  <modes>Public (no API) + Private (with API credentials)</modes>
</skill>`
}

// ============================================================================
// DOCUMENTATION CONSTANTS
// ============================================================================

const binanceMCPInstructions = `## WHEN TO USE

Use Binance MCP connection when you need to:

### ✅ Public Data (No API Required)

1. **Market Data Queries**
   - Get ticker prices: ` + bt + `get_ticker_price(symbol="BTCUSDT")` + bt + `
   - Get order book depth: ` + bt + `get_order_book(symbol="BTCUSDT", limit=10)` + bt + `
   - Get futures order book: ` + bt + `get_futures_order_book(symbol="BTCUSDT")` + bt + `
   - List futures by volume: ` + bt + `list_futures_volume(top=10)` + bt + `

2. **Market Analysis**
   - Analyze bid/ask spreads
   - Check market depth
   - Compare volumes across symbols
   - Track price movements

### ✅ Private Trading (API Credentials Required)

**Prerequisites:**
- API keys configured in ` + bt + `config.json` + bt + ` or environment variables
- ` + bt + `BINANCE_API_KEY` + bt + ` and ` + bt + `BINANCE_SECRET_KEY` + bt + ` set

**Available Operations:**

1. **Check Balances**
   - Spot: ` + bt + `get_spot_balance()` + bt + `
   - Futures: ` + bt + `get_futures_balance()` + bt + `

2. **Open Positions**
   - ` + bt + `open_futures_position(symbol="BTCUSDT", side="LONG", quantity="0.001", leverage=5, confirm=true)` + bt + `

3. **Close Positions**
   - Full close: ` + bt + `close_futures_position(symbol="BTCUSDT", confirm=true)` + bt + `
   - Partial close: ` + bt + `close_futures_position(symbol="BTCUSDT", quantity="0.0005", confirm=true)` + bt + `

## CONNECTION MODES

### Mode 1: Public Mode (No API Keys)

**Automatically activated when:**
- No ` + bt + `api_key` + bt + ` in config.json
- No ` + bt + `BINANCE_API_KEY` + bt + ` environment variable

**Available tools:**
- ` + bt + `get_ticker_price` + bt + ` ✅
- ` + bt + `get_order_book` + bt + ` ✅
- ` + bt + `get_futures_order_book` + bt + ` ✅
- ` + bt + `list_futures_volume` + bt + ` ✅

**Limitations:**
- ❌ Cannot access account balances
- ❌ Cannot open/close positions
- ❌ Cannot execute trades

### Mode 2: Private Mode (With API Keys)

**Automatically activated when:**
- Both ` + bt + `api_key` + bt + ` and ` + bt + `secret_key` + bt + ` configured
- Or both ` + bt + `BINANCE_API_KEY` + bt + ` and ` + bt + `BINANCE_SECRET_KEY` + bt + ` set

**Available tools:**
- All public tools ✅
- ` + bt + `get_spot_balance` + bt + ` ✅
- ` + bt + `get_futures_balance` + bt + ` ✅
- ` + bt + `open_futures_position` + bt + ` ✅
- ` + bt + `close_futures_position` + bt + ` ✅

**Security Features:**
- ✅ All trading operations require ` + bt + `confirm=true` + bt + `
- ✅ API keys never exposed to LLM context
- ✅ Credential validation before enabling private tools

## CONFIGURATION

### Option 1: config.json

` + bt + bt + bt + `json
{
  "tools": {
    "binance": {
      "api_key": "YOUR_BINANCE_API_KEY",
      "secret_key": "YOUR_BINANCE_SECRET_KEY"
    }
  }
}
` + bt + bt + bt + `

### Option 2: Environment Variables

` + bt + bt + bt + `bash
export BINANCE_API_KEY="YOUR_BINANCE_API_KEY"
export BINANCE_SECRET_KEY="YOUR_BINANCE_SECRET_KEY"
` + bt + bt + bt + `

### Verification

After configuration, verify mode with:
` + bt + bt + bt + `
/binance_status
` + bt + bt + bt + `

Expected response:
- **Public Mode:** "Binance MCP connected (public data only)"
- **Private Mode:** "Binance MCP connected (trading enabled)"
`

const binanceMCPAntiPatterns = `## ANTI-PATTERNS TO AVOID

### ❌ Anti-Pattern 1: Trading Without Confirmation

` + bt + bt + bt + `
# BAD - Missing confirm parameter
open_futures_position(symbol="BTCUSDT", side="LONG", quantity="0.001")

# GOOD - Always confirm real orders
open_futures_position(symbol="BTCUSDT", side="LONG", quantity="0.001", confirm=true)
` + bt + bt + bt + `

### ❌ Anti-Pattern 2: Assuming API Keys Exist

` + bt + bt + bt + `
# BAD - Assume private tools available
get_spot_balance()  # May fail if no API configured

# GOOD - Check mode first
# 1. Use /binance_status to verify mode
# 2. If public mode, inform user API keys needed
# 3. If private mode, proceed with balance check
` + bt + bt + bt + `

### ❌ Anti-Pattern 3: Hardcoding Symbols

` + bt + bt + bt + `
# BAD - Assume symbol exists
get_ticker_price(symbol="BTCUSDT")  # What if user wants ETHUSDT?

# GOOD - Use user-provided symbol or ask
# If user says "check Bitcoin price", infer BTCUSDT
# If ambiguous, ask: "Which symbol? BTCUSDT or ETHUSDT?"
` + bt + bt + bt + `

### ❌ Anti-Pattern 4: Ignoring Rate Limits

` + bt + bt + bt + `
# BAD - Rapid repeated calls
get_ticker_price(symbol="BTCUSDT")
get_ticker_price(symbol="BTCUSDT")
get_ticker_price(symbol="BTCUSDT")

# GOOD - Cache results, respect limits
# 1. Call once
# 2. Cache result for reasonable time (e.g., 5s)
# 3. Reuse cached data for subsequent queries
` + bt + bt + bt + `

### ❌ Anti-Pattern 5: Unsafe Position Sizing

` + bt + bt + bt + `
# BAD - No balance check before opening
open_futures_position(symbol="BTCUSDT", side="LONG", quantity="1.0", confirm=true)
# What if user only has $100?

# GOOD - Check balance first
# 1. get_futures_balance()
# 2. Verify sufficient margin
# 3. Warn about risk before opening
` + bt + bt + bt + `
`

const binanceMCPExamples = `## CONCRETE EXAMPLES

### Example 1: Public Price Query (No API)

**User:** "What's the current Bitcoin price?"

**LLM (WITH this skill):**
` + bt + bt + bt + `
1. Infer symbol: BTCUSDT
2. Call: get_ticker_price(symbol="BTCUSDT")
3. Result: "Current price for BTCUSDT: $95,432.50"
4. Response: "Bitcoin (BTC/USDT) is currently trading at $95,432.50"
` + bt + bt + bt + `

**Tokens used:** ~50 (vs. ~200 for web search)

### Example 2: Market Depth Analysis

**User:** "Show me the order book for Ethereum"

**LLM (WITH this skill):**
` + bt + bt + bt + `
1. Infer symbol: ETHUSDT
2. Call: get_order_book(symbol="ETHUSDT", limit=10)
3. Parse result:
   - Best bid: $3,520.50 (10 ETH)
   - Best ask: $3,521.00 (15 ETH)
   - Spread: $0.50 (0.014%)
4. Response: "ETHUSDT order book (top 10):
   - Best Bid: $3,520.50 (10 ETH)
   - Best Ask: $3,521.00 (15 ETH)
   - Spread: $0.50 (0.014%)
   Market is liquid with tight spread."
` + bt + bt + bt + `

### Example 3: Check Balances (Private Mode)

**User:** "How much do I have in my spot account?"

**LLM (WITH this skill):**
` + bt + bt + bt + `
1. Verify mode: Check if API keys configured
2. If private mode: get_spot_balance()
3. Parse result, filter non-zero balances
4. Response: "Your spot account balances:
   - USDT: Free=1,250.50, Locked=0
   - BTC: Free=0.0523, Locked=0
   - ETH: Free=2.15, Locked=0.5"
` + bt + bt + bt + `

**If public mode:**
` + bt + bt + bt + `
1. Detect no API keys configured
2. Response: "I can't access your account balances because API credentials are not configured.
   
   To enable trading features:
   1. Add your Binance API keys to config.json
   2. Or set BINANCE_API_KEY and BINANCE_SECRET_KEY environment variables
   
   I can still show you public market data like prices and order books."
` + bt + bt + bt + `

### Example 4: Open Futures Position

**User:** "Open a small long position on Bitcoin with 5x leverage"

**LLM (WITH this skill):**
` + bt + bt + bt + `
1. Verify private mode (API keys configured)
2. Check balance: get_futures_balance()
3. Calculate quantity: "small" → 0.001 BTC (~$95)
4. Prepare order:
   - symbol: BTCUSDT
   - side: LONG
   - quantity: 0.001
   - leverage: 5
5. **CRITICAL:** Show draft to user for approval
6. User approves → add confirm=true
7. Execute: open_futures_position(symbol="BTCUSDT", side="LONG", quantity="0.001", leverage=5, confirm=true)
8. Response: "✅ Long position opened:
   - Symbol: BTCUSDT
   - Side: LONG
   - Quantity: 0.001 BTC (~$95)
   - Leverage: 5x
   - Order ID: 12345678
   - Status: FILLED"
` + bt + bt + bt + `

### Example 5: Close Position

**User:** "Close my BTCUSDT position"

**LLM (WITH this skill):**
` + bt + bt + bt + `
1. Verify position exists: Check open positions
2. Determine close type: Full close (no quantity specified)
3. Prepare close order:
   - symbol: BTCUSDT
   - confirm: true (user already approved by asking to close)
4. Execute: close_futures_position(symbol="BTCUSDT", confirm=true)
5. Response: "✅ Position closed:
   - Symbol: BTCUSDT
   - Side: SELL (closing LONG)
   - Quantity: 0.001 BTC
   - Order ID: 12345679
   - PnL: +$2.50 (+2.6%)"
` + bt + bt + bt + `

### Example 6: Futures Volume Analysis

**User:** "Which futures are most active today?"

**LLM (WITH this skill):**
` + bt + bt + bt + `
1. Call: list_futures_volume(top=10)
2. Parse result, extract top 10 by quote volume
3. Response: "🚀 Top 10 Futures by Volume (24h):
   1️⃣ BTCUSDT: $1.2B (+2.5%)
   2️⃣ ETHUSDT: $890M (+1.8%)
   3️⃣ BNBUSDT: $456M (+3.2%)
   4️⃣ SOLUSDT: $234M (+5.1%)
   5️⃣ XRPUSDT: $198M (+0.9%)
   ..."
` + bt + bt + bt + `

## QUICK REFERENCE

### Public Commands (No API Required)

| Command | Description | Example |
|---------|-------------|---------|
| ` + bt + `get_ticker_price` + bt + ` | Get current price | ` + bt + `get_ticker_price(symbol="BTCUSDT")` + bt + ` |
| ` + bt + `get_order_book` + bt + ` | Spot order book | ` + bt + `get_order_book(symbol="ETHUSDT", limit=10)` + bt + ` |
| ` + bt + `get_futures_order_book` + bt + ` | Futures order book | ` + bt + `get_futures_order_book(symbol="BTCUSDT")` + bt + ` |
| ` + bt + `list_futures_volume` + bt + ` | Volume ranking | ` + bt + `list_futures_volume(top=10)` + bt + ` |

### Private Commands (API Required)

| Command | Description | Example |
|---------|-------------|---------|
| ` + bt + `get_spot_balance` + bt + ` | Spot balances | ` + bt + `get_spot_balance()` + bt + ` |
| ` + bt + `get_futures_balance` + bt + ` | Futures balances | ` + bt + `get_futures_balance()` + bt + ` |
| ` + bt + `open_futures_position` + bt + ` | Open position | ` + bt + `open_futures_position(symbol="BTCUSDT", side="LONG", quantity="0.001", leverage=5, confirm=true)` + bt + ` |
| ` + bt + `close_futures_position` + bt + ` | Close position | ` + bt + `close_futures_position(symbol="BTCUSDT", confirm=true)` + bt + ` |

## Fast-path Commands (Limited Availability)

**Note:** These commands are documented for future implementation. Currently, use the tool functions directly:

| Tool Function | Example |
|---------------|---------|
| ` + bt + `get_ticker_price` + bt + ` | ` + bt + `get_ticker_price(symbol="BTCUSDT")` + bt + ` |
| ` + bt + `get_order_book` + bt + ` | ` + bt + `get_order_book(symbol="ETHUSDT", limit=10)` + bt + ` |
| ` + bt + `open_futures_position` + bt + ` | ` + bt + `open_futures_position(symbol="BTCUSDT", side="LONG", quantity="0.001", leverage=5, confirm=true)` + bt + ` |
| ` + bt + `close_futures_position` + bt + ` | ` + bt + `close_futures_position(symbol="BTCUSDT", confirm=true)` + bt + ` |

**Currently Available Fast-paths:**
- ` + bt + `/status` + bt + ` - Show system status
- ` + bt + `/help` + bt + ` - Show available commands
- ` + bt + `/show model` + bt + ` - Show active model
- ` + bt + `/show channel` + bt + ` - Show communication channel
- ` + bt + `/list models` + bt + ` - List configured models
- ` + bt + `/list channels` + bt + ` - List configured channels
- ` + bt + `#BATCH_ID` + bt + ` - Query batch task status (e.g., ` + bt + `#IMA_GEN_02_03_26_1500` + bt + `)
- ` + bt + `/bundle_approve` + bt + `, ` + bt + `/bundle_regen` + bt + `, ` + bt + `/bundle_edit` + bt + ` - Bundle management
- ` + bt + `/disable_sentinel [5m|15m|1h]` + bt + ` - Temporarily disable security sentinel
- ` + bt + `/activate_sentinel` + bt + ` - Immediately activate sentinel
- ` + bt + `/sentinel_status` + bt + ` - Check sentinel status

## SAFETY CHECKLIST

Before executing ANY trading operation:

- [ ] **Verify API Mode:** Confirm private mode is active
- [ ] **Check Balance:** Ensure sufficient margin
- [ ] **Validate Symbol:** Confirm symbol exists and is liquid
- [ ] **Review Quantity:** Verify quantity is appropriate
- [ ] **Confirm Leverage:** Ensure leverage is within safe range (1-125)
- [ ] **User Approval:** ALWAYS get explicit ` + bt + `confirm=true` + bt + ` for real orders
- [ ] **Risk Warning:** Inform user about potential losses

**Golden Rule:** "If in doubt, ask. Never execute a trade without explicit confirmation."
`
