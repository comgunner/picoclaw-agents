// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package tools

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/comgunner/picoclaw/pkg/utils"
)

type BinanceTickerPriceTool struct {
	hasPrivateCreds bool
}

func NewBinanceTickerPriceTool() *BinanceTickerPriceTool {
	return NewBinanceTickerPriceToolFromConfig("", "")
}

func NewBinanceTickerPriceToolFromConfig(configAPIKey, configSecretKey string) *BinanceTickerPriceTool {
	apiKey := strings.TrimSpace(os.Getenv(utils.EnvBinanceAPIKey))
	secretKey := strings.TrimSpace(os.Getenv(utils.EnvBinanceSecretKey))
	if apiKey == "" {
		apiKey = strings.TrimSpace(configAPIKey)
	}
	if secretKey == "" {
		secretKey = strings.TrimSpace(configSecretKey)
	}
	return &BinanceTickerPriceTool{
		hasPrivateCreds: apiKey != "" && secretKey != "",
	}
}

func (t *BinanceTickerPriceTool) Name() string {
	return "binance_get_ticker_price"
}

func (t *BinanceTickerPriceTool) Description() string {
	return "Get current ticker price from Binance for a symbol (e.g., BTCUSDT)."
}

func (t *BinanceTickerPriceTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"symbol": map[string]any{
				"type":        "string",
				"description": "Binance symbol, e.g. BTCUSDT",
			},
		},
		"required": []string{"symbol"},
	}
}

func (t *BinanceTickerPriceTool) Execute(ctx context.Context, args map[string]any) *ToolResult {
	symbol, _ := args["symbol"].(string)
	symbol = strings.TrimSpace(symbol)
	if symbol == "" {
		return ErrorResult("symbol is required and must be a string")
	}

	callCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	result, err := utils.BinanceTickerPrice(callCtx, symbol)
	if err != nil {
		return ErrorResult(fmt.Sprintf("binance ticker query failed: %v", err)).WithError(err)
	}
	if !t.hasPrivateCreds {
		result += " (source: Binance public endpoint; API credentials not configured)"
	}
	return UserResult(result)
}

type BinanceOrderBookTool struct{}

func NewBinanceOrderBookTool() *BinanceOrderBookTool {
	return &BinanceOrderBookTool{}
}

func (t *BinanceOrderBookTool) Name() string {
	return "binance_get_order_book"
}

func (t *BinanceOrderBookTool) Description() string {
	return "Get Binance order book depth (bids/asks) for a symbol (e.g., BTCUSDT). Useful for market depth analysis."
}

func (t *BinanceOrderBookTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"symbol": map[string]any{
				"type":        "string",
				"description": "Binance symbol, e.g. BTCUSDT",
			},
			"limit": map[string]any{
				"type":        "integer",
				"description": "Optional depth levels to return (default 10, max 100).",
			},
		},
		"required": []string{"symbol"},
	}
}

func (t *BinanceOrderBookTool) Execute(ctx context.Context, args map[string]any) *ToolResult {
	symbol, _ := args["symbol"].(string)
	symbol = strings.TrimSpace(symbol)
	if symbol == "" {
		return ErrorResult("symbol is required and must be a string")
	}

	limit, hasLimit, err := parseOptionalIntegerArg(args, "limit")
	if err != nil {
		return ErrorResult(err.Error())
	}
	if !hasLimit {
		limit = 10
	}

	callCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	result, err := utils.BinanceOrderBook(callCtx, symbol, limit)
	if err != nil {
		return ErrorResult(fmt.Sprintf("binance order book query failed: %v", err)).WithError(err)
	}
	return UserResult(result)
}

type BinanceFuturesOrderBookTool struct{}

func NewBinanceFuturesOrderBookTool() *BinanceFuturesOrderBookTool {
	return &BinanceFuturesOrderBookTool{}
}

func (t *BinanceFuturesOrderBookTool) Name() string {
	return "binance_get_futures_order_book"
}

func (t *BinanceFuturesOrderBookTool) Description() string {
	return "Get Binance USDT-M futures order book depth (bids/asks) for a symbol (e.g., BTCUSDT). Useful for futures market depth analysis."
}

func (t *BinanceFuturesOrderBookTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"symbol": map[string]any{
				"type":        "string",
				"description": "Binance futures symbol, e.g. BTCUSDT",
			},
			"limit": map[string]any{
				"type":        "integer",
				"description": "Optional depth levels to return (default 10, max 100).",
			},
		},
		"required": []string{"symbol"},
	}
}

func (t *BinanceFuturesOrderBookTool) Execute(ctx context.Context, args map[string]any) *ToolResult {
	symbol, _ := args["symbol"].(string)
	symbol = strings.TrimSpace(symbol)
	if symbol == "" {
		return ErrorResult("symbol is required and must be a string")
	}

	limit, hasLimit, err := parseOptionalIntegerArg(args, "limit")
	if err != nil {
		return ErrorResult(err.Error())
	}
	if !hasLimit {
		limit = 10
	}

	callCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	result, err := utils.BinanceFuturesOrderBook(callCtx, symbol, limit)
	if err != nil {
		return ErrorResult(fmt.Sprintf("binance futures order book query failed: %v", err)).WithError(err)
	}
	return UserResult(result)
}

type BinanceListFuturesVolumeTool struct{}

func NewBinanceListFuturesVolumeTool() *BinanceListFuturesVolumeTool {
	return &BinanceListFuturesVolumeTool{}
}

func (t *BinanceListFuturesVolumeTool) Name() string {
	return "binance_list_futures_volume"
}

func (t *BinanceListFuturesVolumeTool) Description() string {
	return "List Binance USDT-M futures symbols sorted by 24h quote volume (descending)."
}

func (t *BinanceListFuturesVolumeTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"top": map[string]any{
				"type":        "integer",
				"description": "How many symbols to return (default 10, max 50).",
			},
		},
	}
}

func (t *BinanceListFuturesVolumeTool) Execute(ctx context.Context, args map[string]any) *ToolResult {
	top, hasTop, err := parseOptionalIntegerArg(args, "top")
	if err != nil {
		return ErrorResult(err.Error())
	}
	if !hasTop {
		top = 10
	}

	callCtx, cancel := context.WithTimeout(ctx, 12*time.Second)
	defer cancel()

	result, err := utils.BinanceListFuturesByVolume(callCtx, top)
	if err != nil {
		return ErrorResult(fmt.Sprintf("binance futures volume listing failed: %v", err)).WithError(err)
	}
	return UserResult(result)
}

type BinanceSpotBalanceTool struct {
	apiKey    string
	secretKey string
}

func NewBinanceSpotBalanceTool(apiKey, secretKey string) *BinanceSpotBalanceTool {
	return &BinanceSpotBalanceTool{
		apiKey:    strings.TrimSpace(apiKey),
		secretKey: strings.TrimSpace(secretKey),
	}
}

func NewBinanceSpotBalanceToolFromEnv() *BinanceSpotBalanceTool {
	return NewBinanceSpotBalanceTool(
		os.Getenv(utils.EnvBinanceAPIKey),
		os.Getenv(utils.EnvBinanceSecretKey),
	)
}

func NewBinanceSpotBalanceToolFromConfig(configAPIKey, configSecretKey string) *BinanceSpotBalanceTool {
	apiKey := strings.TrimSpace(os.Getenv(utils.EnvBinanceAPIKey))
	secretKey := strings.TrimSpace(os.Getenv(utils.EnvBinanceSecretKey))
	if apiKey == "" {
		apiKey = strings.TrimSpace(configAPIKey)
	}
	if secretKey == "" {
		secretKey = strings.TrimSpace(configSecretKey)
	}
	return NewBinanceSpotBalanceTool(apiKey, secretKey)
}

func (t *BinanceSpotBalanceTool) Name() string {
	return "binance_get_spot_balance"
}

func (t *BinanceSpotBalanceTool) Description() string {
	return "Get spot balances from Binance account (requires BINANCE_API_KEY and BINANCE_SECRET_KEY)."
}

func (t *BinanceSpotBalanceTool) Parameters() map[string]any {
	return map[string]any{
		"type":       "object",
		"properties": map[string]any{},
	}
}

func (t *BinanceSpotBalanceTool) Execute(ctx context.Context, _ map[string]any) *ToolResult {
	if t.apiKey == "" || t.secretKey == "" {
		return UserResult(
			"BINANCE_API_KEY/BINANCE_SECRET_KEY were not found in .env or config.json (tools.binance). " +
				"You can query public price data without API keys using:\n" +
				"curl -s 'https://api.binance.com/api/v3/ticker/price?symbol=BTCUSDT'\n" +
				"Replace BTCUSDT with any valid Binance symbol (for example ETHUSDT, SOLUSDT). " +
				"If you want only the numeric price, ask: " +
				"\"Use binance_get_ticker_price with symbol BTCUSDT and return only the numeric price\".",
		)
	}

	callCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	result, err := utils.BinanceSpotBalances(callCtx, t.apiKey, t.secretKey)
	if err != nil {
		return ErrorResult(fmt.Sprintf("binance balance query failed: %v", err)).WithError(err)
	}
	return UserResult(result)
}

type BinanceFuturesBalanceTool struct {
	apiKey    string
	secretKey string
}

func NewBinanceFuturesBalanceToolFromConfig(configAPIKey, configSecretKey string) *BinanceFuturesBalanceTool {
	apiKey := strings.TrimSpace(os.Getenv(utils.EnvBinanceAPIKey))
	secretKey := strings.TrimSpace(os.Getenv(utils.EnvBinanceSecretKey))
	if apiKey == "" {
		apiKey = strings.TrimSpace(configAPIKey)
	}
	if secretKey == "" {
		secretKey = strings.TrimSpace(configSecretKey)
	}
	return &BinanceFuturesBalanceTool{
		apiKey:    apiKey,
		secretKey: secretKey,
	}
}

func (t *BinanceFuturesBalanceTool) Name() string {
	return "binance_get_futures_balance"
}

func (t *BinanceFuturesBalanceTool) Description() string {
	return "Get USDT-M futures balances from Binance account (requires BINANCE_API_KEY and BINANCE_SECRET_KEY)."
}

func (t *BinanceFuturesBalanceTool) Parameters() map[string]any {
	return map[string]any{
		"type":       "object",
		"properties": map[string]any{},
	}
}

func (t *BinanceFuturesBalanceTool) Execute(ctx context.Context, _ map[string]any) *ToolResult {
	if t.apiKey == "" || t.secretKey == "" {
		return UserResult(
			"BINANCE_API_KEY/BINANCE_SECRET_KEY were not found in .env or config.json (tools.binance). " +
				"Futures balance queries require API credentials.",
		)
	}

	callCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	result, err := utils.BinanceFuturesBalances(callCtx, t.apiKey, t.secretKey)
	if err != nil {
		return ErrorResult(fmt.Sprintf("binance futures balance query failed: %v", err)).WithError(err)
	}
	return UserResult(result)
}

type BinanceFuturesOpenPositionTool struct {
	apiKey    string
	secretKey string
}

func NewBinanceFuturesOpenPositionToolFromConfig(configAPIKey, configSecretKey string) *BinanceFuturesOpenPositionTool {
	apiKey := strings.TrimSpace(os.Getenv(utils.EnvBinanceAPIKey))
	secretKey := strings.TrimSpace(os.Getenv(utils.EnvBinanceSecretKey))
	if apiKey == "" {
		apiKey = strings.TrimSpace(configAPIKey)
	}
	if secretKey == "" {
		secretKey = strings.TrimSpace(configSecretKey)
	}
	return &BinanceFuturesOpenPositionTool{
		apiKey:    apiKey,
		secretKey: secretKey,
	}
}

func (t *BinanceFuturesOpenPositionTool) Name() string {
	return "binance_open_futures_position"
}

func (t *BinanceFuturesOpenPositionTool) Description() string {
	return "Open Binance USDT-M futures position with MARKET order (LONG/SHORT). Requires BINANCE_API_KEY/BINANCE_SECRET_KEY and confirm=true."
}

func (t *BinanceFuturesOpenPositionTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"symbol": map[string]any{
				"type":        "string",
				"description": "Binance futures symbol, e.g. BTCUSDT",
			},
			"side": map[string]any{
				"type":        "string",
				"description": "Direction to open: LONG or SHORT.",
				"enum":        []string{"LONG", "SHORT"},
			},
			"quantity": map[string]any{
				"type":        "string",
				"description": "Order quantity as decimal string, e.g. 0.001",
			},
			"leverage": map[string]any{
				"type":        "integer",
				"description": "Optional leverage (1-125). If provided, leverage is set before opening.",
			},
			"confirm": map[string]any{
				"type":        "boolean",
				"description": "Must be true to execute a real trading order.",
			},
		},
		"required": []string{"symbol", "side", "quantity", "confirm"},
	}
}

func (t *BinanceFuturesOpenPositionTool) Execute(ctx context.Context, args map[string]any) *ToolResult {
	if t.apiKey == "" || t.secretKey == "" {
		return UserResult(
			"BINANCE_API_KEY/BINANCE_SECRET_KEY are not configured in .env or config.json (tools.binance). " +
				"Trading operations are disabled without API credentials.",
		)
	}

	confirm, _ := args["confirm"].(bool)
	if !confirm {
		return ErrorResult("futures open operation requires confirm: true before executing a real order")
	}

	symbol, _ := args["symbol"].(string)
	symbol = strings.ToUpper(strings.TrimSpace(symbol))
	if symbol == "" {
		return ErrorResult("symbol is required")
	}

	sideRaw, _ := args["side"].(string)
	side, err := normalizeFuturesDirection(sideRaw)
	if err != nil {
		return ErrorResult(err.Error())
	}

	quantity, err := parsePositiveDecimalArg(args, "quantity")
	if err != nil {
		return ErrorResult(err.Error())
	}

	leverage, hasLeverage, err := parseOptionalIntegerArg(args, "leverage")
	if err != nil {
		return ErrorResult(err.Error())
	}
	if hasLeverage && (leverage < 1 || leverage > 125) {
		return ErrorResult("leverage must be between 1 and 125")
	}

	callCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	result, err := utils.BinanceFuturesOpenMarketPosition(
		callCtx,
		t.apiKey,
		t.secretKey,
		symbol,
		side,
		quantity,
		leverage,
		hasLeverage,
	)
	if err != nil {
		return ErrorResult(fmt.Sprintf("binance futures open failed: %v", err)).WithError(err)
	}
	return UserResult(result)
}

type BinanceFuturesClosePositionTool struct {
	apiKey    string
	secretKey string
}

func NewBinanceFuturesClosePositionToolFromConfig(
	configAPIKey, configSecretKey string,
) *BinanceFuturesClosePositionTool {
	apiKey := strings.TrimSpace(os.Getenv(utils.EnvBinanceAPIKey))
	secretKey := strings.TrimSpace(os.Getenv(utils.EnvBinanceSecretKey))
	if apiKey == "" {
		apiKey = strings.TrimSpace(configAPIKey)
	}
	if secretKey == "" {
		secretKey = strings.TrimSpace(configSecretKey)
	}
	return &BinanceFuturesClosePositionTool{
		apiKey:    apiKey,
		secretKey: secretKey,
	}
}

func (t *BinanceFuturesClosePositionTool) Name() string {
	return "binance_close_futures_position"
}

func (t *BinanceFuturesClosePositionTool) Description() string {
	return "Close Binance USDT-M futures position with MARKET reduce-only order. If quantity is omitted, closes full net position. Requires BINANCE_API_KEY/BINANCE_SECRET_KEY and confirm=true."
}

func (t *BinanceFuturesClosePositionTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"symbol": map[string]any{
				"type":        "string",
				"description": "Binance futures symbol, e.g. BTCUSDT",
			},
			"quantity": map[string]any{
				"type":        "string",
				"description": "Optional quantity as decimal string. If omitted, closes full net position.",
			},
			"confirm": map[string]any{
				"type":        "boolean",
				"description": "Must be true to execute a real trading order.",
			},
		},
		"required": []string{"symbol", "confirm"},
	}
}

func (t *BinanceFuturesClosePositionTool) Execute(ctx context.Context, args map[string]any) *ToolResult {
	if t.apiKey == "" || t.secretKey == "" {
		return UserResult(
			"BINANCE_API_KEY/BINANCE_SECRET_KEY are not configured in .env or config.json (tools.binance). " +
				"Trading operations are disabled without API credentials.",
		)
	}

	confirm, _ := args["confirm"].(bool)
	if !confirm {
		return ErrorResult("futures close operation requires confirm: true before executing a real order")
	}

	symbol, _ := args["symbol"].(string)
	symbol = strings.ToUpper(strings.TrimSpace(symbol))
	if symbol == "" {
		return ErrorResult("symbol is required")
	}

	quantity, hasQuantity, err := parseOptionalPositiveDecimalArg(args, "quantity")
	if err != nil {
		return ErrorResult(err.Error())
	}

	callCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	result, err := utils.BinanceFuturesCloseMarketPosition(
		callCtx,
		t.apiKey,
		t.secretKey,
		symbol,
		quantity,
		hasQuantity,
	)
	if err != nil {
		return ErrorResult(fmt.Sprintf("binance futures close failed: %v", err)).WithError(err)
	}
	return UserResult(result)
}

func normalizeFuturesDirection(sideRaw string) (string, error) {
	side := strings.ToUpper(strings.TrimSpace(sideRaw))
	switch side {
	case "LONG":
		return "BUY", nil
	case "SHORT":
		return "SELL", nil
	default:
		return "", fmt.Errorf("side must be LONG or SHORT")
	}
}

func parsePositiveDecimalArg(args map[string]any, key string) (string, error) {
	raw, ok := args[key]
	if !ok {
		return "", fmt.Errorf("%s is required", key)
	}

	switch v := raw.(type) {
	case string:
		q := strings.TrimSpace(v)
		if q == "" {
			return "", fmt.Errorf("%s is required", key)
		}
		f, err := strconv.ParseFloat(q, 64)
		if err != nil || f <= 0 {
			return "", fmt.Errorf("%s must be a positive decimal value", key)
		}
		return q, nil
	case float64:
		if v <= 0 {
			return "", fmt.Errorf("%s must be a positive decimal value", key)
		}
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	default:
		return "", fmt.Errorf("%s must be a decimal string or number", key)
	}
}

func parseOptionalPositiveDecimalArg(args map[string]any, key string) (string, bool, error) {
	raw, ok := args[key]
	if !ok || raw == nil {
		return "", false, nil
	}
	value, err := parsePositiveDecimalArg(args, key)
	if err != nil {
		return "", false, err
	}
	return value, true, nil
}

func parseOptionalIntegerArg(args map[string]any, key string) (int, bool, error) {
	raw, ok := args[key]
	if !ok || raw == nil {
		return 0, false, nil
	}

	switch v := raw.(type) {
	case float64:
		i := int(v)
		if float64(i) != v {
			return 0, false, fmt.Errorf("%s must be an integer", key)
		}
		return i, true, nil
	case int:
		return v, true, nil
	case string:
		text := strings.TrimSpace(v)
		if text == "" {
			return 0, false, nil
		}
		i, err := strconv.Atoi(text)
		if err != nil {
			return 0, false, fmt.Errorf("%s must be an integer", key)
		}
		return i, true, nil
	default:
		return 0, false, fmt.Errorf("%s must be an integer", key)
	}
}
