// PicoClaw - Ultra-lightweight personal AI agent
// Inspired by and based on nanobot: https://github.com/HKUDS/nanobot
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Modified by comgunner (https://github.com/comgunner)
// Custom Fork: https://github.com/comgunner/picoclaw-agents

package utils

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const (
	EnvBinanceAPIKey    = "BINANCE_API_KEY"
	EnvBinanceSecretKey = "BINANCE_SECRET_KEY"
)

// NewBinanceMCPServer builds an MCP server with Binance tools.
// Public price tool is always available; private balance tool is added only
// when both API key and secret key are provided.
func NewBinanceMCPServer(apiKey, secretKey string) *server.MCPServer {
	s := server.NewMCPServer(
		"binance-mcp",
		"1.0.0",
		server.WithLogging(),
	)

	s.AddTool(
		mcp.NewTool(
			"get_ticker_price",
			mcp.WithDescription("Get current ticker price for a Binance symbol (e.g., BTCUSDT)"),
			mcp.WithString(
				"symbol",
				mcp.Description("Binance symbol to query (e.g., BTCUSDT)"),
				mcp.Required(),
			),
		),
		getTickerPriceHandler,
	)
	s.AddTool(
		mcp.NewTool(
			"get_order_book",
			mcp.WithDescription("Get Binance order book depth (bids/asks) for a symbol (e.g., BTCUSDT)."),
			mcp.WithString(
				"symbol",
				mcp.Description("Binance symbol to query (e.g., BTCUSDT)"),
				mcp.Required(),
			),
			mcp.WithNumber(
				"limit",
				mcp.Description("Optional depth levels (default 10, max 100)"),
			),
		),
		getOrderBookHandler,
	)
	s.AddTool(
		mcp.NewTool(
			"get_futures_order_book",
			mcp.WithDescription(
				"Get Binance USDT-M futures order book depth (bids/asks) for a symbol (e.g., BTCUSDT).",
			),
			mcp.WithString(
				"symbol",
				mcp.Description("Binance futures symbol to query (e.g., BTCUSDT)"),
				mcp.Required(),
			),
			mcp.WithNumber(
				"limit",
				mcp.Description("Optional depth levels (default 10, max 100)"),
			),
		),
		getFuturesOrderBookHandler,
	)
	s.AddTool(
		mcp.NewTool(
			"list_futures_volume",
			mcp.WithDescription("List USDT-M futures symbols sorted by 24h quote volume (descending)."),
			mcp.WithNumber("top", mcp.Description("How many symbols to return (default 10, max 50)")),
		),
		listFuturesVolumeHandler,
	)

	apiKey = strings.TrimSpace(apiKey)
	secretKey = strings.TrimSpace(secretKey)
	if apiKey != "" && secretKey != "" {
		s.AddTool(
			mcp.NewTool(
				"get_spot_balance",
				mcp.WithDescription("Get spot account balances (requires Binance API credentials)"),
			),
			getSpotBalanceHandler(apiKey, secretKey),
		)
		s.AddTool(
			mcp.NewTool(
				"get_futures_balance",
				mcp.WithDescription("Get USDT-M futures account balances (requires Binance API credentials)"),
			),
			getFuturesBalanceHandler(apiKey, secretKey),
		)
		s.AddTool(
			mcp.NewTool(
				"open_futures_position",
				mcp.WithDescription("Open Binance USDT-M futures MARKET position (LONG/SHORT). Requires confirm=true."),
				mcp.WithString("symbol", mcp.Description("Binance futures symbol, e.g. BTCUSDT"), mcp.Required()),
				mcp.WithString("side", mcp.Description("LONG or SHORT"), mcp.Enum("LONG", "SHORT"), mcp.Required()),
				mcp.WithString("quantity", mcp.Description("Quantity as decimal string, e.g. 0.001"), mcp.Required()),
				mcp.WithNumber("leverage", mcp.Description("Optional leverage (1-125)")),
				mcp.WithBoolean("confirm", mcp.Description("Must be true to execute real order"), mcp.Required()),
			),
			openFuturesPositionHandler(apiKey, secretKey),
		)
		s.AddTool(
			mcp.NewTool(
				"close_futures_position",
				mcp.WithDescription(
					"Close Binance USDT-M futures MARKET position (reduce-only). If quantity omitted, closes full net position. Requires confirm=true.",
				),
				mcp.WithString("symbol", mcp.Description("Binance futures symbol, e.g. BTCUSDT"), mcp.Required()),
				mcp.WithString("quantity", mcp.Description("Optional quantity as decimal string")),
				mcp.WithBoolean("confirm", mcp.Description("Must be true to execute real order"), mcp.Required()),
			),
			closeFuturesPositionHandler(apiKey, secretKey),
		)
	}

	return s
}

// ServeBinanceMCPStdio starts the Binance MCP server over stdio.
func ServeBinanceMCPStdio(apiKey, secretKey string) error {
	return server.ServeStdio(NewBinanceMCPServer(apiKey, secretKey))
}

// BinanceTickerPrice returns a human-readable price line for the given symbol.
func BinanceTickerPrice(ctx context.Context, symbol string) (string, error) {
	symbol = strings.ToUpper(strings.TrimSpace(symbol))
	if symbol == "" {
		return "", errors.New("symbol is required")
	}

	client := binance.NewClient("", "")
	res, err := client.NewListPricesService().Symbol(symbol).Do(ctx)
	if err != nil {
		return "", fmt.Errorf("binance api error: %w", err)
	}
	if len(res) == 0 {
		return "", fmt.Errorf("symbol %s not found", symbol)
	}

	return fmt.Sprintf("Current price for %s: %s", res[0].Symbol, res[0].Price), nil
}

// BinanceOrderBook returns a human-readable order book summary for the given symbol.
func BinanceOrderBook(ctx context.Context, symbol string, limit int) (string, error) {
	symbol = strings.ToUpper(strings.TrimSpace(symbol))
	if symbol == "" {
		return "", errors.New("symbol is required")
	}

	normalizedLimit, err := normalizeOrderBookLimit(limit)
	if err != nil {
		return "", err
	}

	client := binance.NewClient("", "")
	res, err := client.NewDepthService().Symbol(symbol).Limit(normalizedLimit).Do(ctx)
	if err != nil {
		return "", fmt.Errorf("binance api error: %w", err)
	}

	return formatOrderBook(symbol, normalizedLimit, "Order book", res.Bids, res.Asks), nil
}

type BookEntry struct {
	Price    string
	Quantity string
}

func BinanceFuturesOrderBook(ctx context.Context, symbol string, limit int) (string, error) {
	symbol = strings.ToUpper(strings.TrimSpace(symbol))
	if symbol == "" {
		return "", errors.New("symbol is required")
	}

	normalizedLimit, err := normalizeOrderBookLimit(limit)
	if err != nil {
		return "", err
	}

	client := futures.NewClient("", "")
	res, err := client.NewDepthService().Symbol(symbol).Limit(normalizedLimit).Do(ctx)
	if err != nil {
		return "", fmt.Errorf("binance futures api error: %w", err)
	}

	bids := make([]BookEntry, len(res.Bids))
	for i, b := range res.Bids {
		bids[i] = BookEntry{Price: b.Price, Quantity: b.Quantity}
	}
	asks := make([]BookEntry, len(res.Asks))
	for i, a := range res.Asks {
		asks[i] = BookEntry{Price: a.Price, Quantity: a.Quantity}
	}

	return formatOrderBookEntries(symbol, normalizedLimit, "Futures order book", bids, asks), nil
}

func formatOrderBook(symbol string, limit int, title string, bids []binance.Bid, asks []binance.Ask) string {
	eb := make([]BookEntry, len(bids))
	for i, b := range bids {
		eb[i] = BookEntry{Price: b.Price, Quantity: b.Quantity}
	}
	ea := make([]BookEntry, len(asks))
	for i, a := range asks {
		ea[i] = BookEntry{Price: a.Price, Quantity: a.Quantity}
	}
	return formatOrderBookEntries(symbol, limit, title, eb, ea)
}

func formatOrderBookEntries(symbol string, limit int, title string, bids []BookEntry, asks []BookEntry) string {

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s for %s (top %d)\n", title, symbol, limit))
	sb.WriteString("Bids:\n")
	if len(bids) == 0 {
		sb.WriteString("- (no bids)\n")
	} else {
		for _, bid := range bids {
			sb.WriteString(fmt.Sprintf("- %s | %s\n", bid.Price, bid.Quantity))
		}
	}

	sb.WriteString("Asks:\n")
	if len(asks) == 0 {
		sb.WriteString("- (no asks)")
	} else {
		for i, ask := range asks {
			sb.WriteString(fmt.Sprintf("- %s | %s", ask.Price, ask.Quantity))
			if i < len(asks)-1 {
				sb.WriteString("\n")
			}
		}
	}

	return sb.String()
}

// BinanceListFuturesByVolume returns USDT-M futures symbols sorted by 24h quote volume.
func BinanceListFuturesByVolume(ctx context.Context, top int) (string, error) {
	if top == 0 {
		top = 10
	}
	if top < 1 {
		return "", errors.New("top must be >= 1")
	}
	if top > 50 {
		return "", errors.New("top must be <= 50")
	}

	client := futures.NewClient("", "")
	stats, err := client.NewListPriceChangeStatsService().Do(ctx)
	if err != nil {
		return "", fmt.Errorf("binance futures api error: %w", err)
	}
	if len(stats) == 0 {
		return "No futures volume data available.", nil
	}

	type row struct {
		symbol      string
		lastPrice   string
		volume      string
		quoteVolume float64
		changePct   string
	}
	rows := make([]row, 0, len(stats))
	for _, s := range stats {
		qv, err := strconv.ParseFloat(strings.TrimSpace(s.QuoteVolume), 64)
		if err != nil {
			continue
		}
		rows = append(rows, row{
			symbol:      s.Symbol,
			lastPrice:   s.LastPrice,
			volume:      s.Volume,
			quoteVolume: qv,
			changePct:   s.PriceChangePercent,
		})
	}
	if len(rows) == 0 {
		return "No parseable futures volume data available.", nil
	}

	sort.Slice(rows, func(i, j int) bool { return rows[i].quoteVolume > rows[j].quoteVolume })
	if top > len(rows) {
		top = len(rows)
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Top %d futures by 24h quote volume:\n", top))
	for i := 0; i < top; i++ {
		r := rows[i]
		sb.WriteString(fmt.Sprintf("%d. %s | quoteVolume=%.2f | volume=%s | last=%s | change=%s%%",
			i+1, r.symbol, r.quoteVolume, r.volume, r.lastPrice, r.changePct))
		if i < top-1 {
			sb.WriteString("\n")
		}
	}
	return sb.String(), nil
}

// BinanceSpotBalances returns a human-readable summary of non-zero spot balances.
func BinanceSpotBalances(ctx context.Context, apiKey, secretKey string) (string, error) {
	apiKey = strings.TrimSpace(apiKey)
	secretKey = strings.TrimSpace(secretKey)
	if apiKey == "" || secretKey == "" {
		return "", errors.New("missing BINANCE_API_KEY or BINANCE_SECRET_KEY")
	}

	client := binance.NewClient(apiKey, secretKey)
	res, err := client.NewGetAccountService().Do(ctx)
	if err != nil {
		return "", fmt.Errorf("binance api error: %w", err)
	}

	var lines []string
	for _, b := range res.Balances {
		if isZeroAmount(b.Free) && isZeroAmount(b.Locked) {
			continue
		}
		lines = append(lines, fmt.Sprintf("- %s: Free=%s, Locked=%s", b.Asset, b.Free, b.Locked))
	}

	if len(lines) == 0 {
		return "Spot Balances:\nNo assets with balance found.", nil
	}
	return "Spot Balances:\n" + strings.Join(lines, "\n"), nil
}

// BinanceFuturesBalances returns a human-readable summary of non-zero USDT-M futures balances.
func BinanceFuturesBalances(ctx context.Context, apiKey, secretKey string) (string, error) {
	apiKey = strings.TrimSpace(apiKey)
	secretKey = strings.TrimSpace(secretKey)
	if apiKey == "" || secretKey == "" {
		return "", errors.New("missing BINANCE_API_KEY or BINANCE_SECRET_KEY")
	}

	client := futures.NewClient(apiKey, secretKey)
	res, err := client.NewGetBalanceService().Do(ctx)
	if err != nil {
		return "", fmt.Errorf("binance futures api error: %w", err)
	}

	var lines []string
	for _, b := range res {
		if isZeroAmount(b.Balance) && isZeroAmount(b.AvailableBalance) && isZeroAmount(b.CrossWalletBalance) {
			continue
		}
		lines = append(lines, fmt.Sprintf("- %s: Total=%s, Available=%s, CrossWallet=%s, UnrealizedPnL=%s",
			b.Asset, b.Balance, b.AvailableBalance, b.CrossWalletBalance, b.CrossUnPnl))
	}

	if len(lines) == 0 {
		return "Futures Balances:\nNo assets with balance found.", nil
	}
	return "Futures Balances:\n" + strings.Join(lines, "\n"), nil
}

func getTickerPriceHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	symbol, err := request.RequireString("symbol")
	if err != nil {
		return mcp.NewToolResultError("symbol is required and must be a string"), nil
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	priceLine, err := BinanceTickerPrice(ctx, symbol)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(priceLine), nil
}

func getOrderBookHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	symbol, err := request.RequireString("symbol")
	if err != nil {
		return mcp.NewToolResultError("symbol is required and must be a string"), nil
	}

	limit := request.GetInt("limit", 10)

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	book, err := BinanceOrderBook(ctx, symbol, limit)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(book), nil
}

func getFuturesOrderBookHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	symbol, err := request.RequireString("symbol")
	if err != nil {
		return mcp.NewToolResultError("symbol is required and must be a string"), nil
	}

	limit := request.GetInt("limit", 10)

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	book, err := BinanceFuturesOrderBook(ctx, symbol, limit)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(book), nil
}

func listFuturesVolumeHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	top := request.GetInt("top", 10)

	ctx, cancel := context.WithTimeout(ctx, 12*time.Second)
	defer cancel()

	listing, err := BinanceListFuturesByVolume(ctx, top)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(listing), nil
}

func getSpotBalanceHandler(apiKey, secretKey string) server.ToolHandlerFunc {
	return func(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
		defer cancel()

		balances, err := BinanceSpotBalances(ctx, apiKey, secretKey)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(balances), nil
	}
}

func getFuturesBalanceHandler(apiKey, secretKey string) server.ToolHandlerFunc {
	return func(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
		defer cancel()

		balances, err := BinanceFuturesBalances(ctx, apiKey, secretKey)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(balances), nil
	}
}

// BinanceFuturesOpenMarketPosition opens a USDT-M futures MARKET order.
func BinanceFuturesOpenMarketPosition(
	ctx context.Context,
	apiKey, secretKey, symbol, side, quantity string,
	leverage int,
	setLeverage bool,
) (string, error) {
	apiKey = strings.TrimSpace(apiKey)
	secretKey = strings.TrimSpace(secretKey)
	symbol = strings.ToUpper(strings.TrimSpace(symbol))
	side = strings.ToUpper(strings.TrimSpace(side))
	quantity = strings.TrimSpace(quantity)

	if apiKey == "" || secretKey == "" {
		return "", errors.New("missing BINANCE_API_KEY or BINANCE_SECRET_KEY")
	}
	if symbol == "" {
		return "", errors.New("symbol is required")
	}
	if quantity == "" {
		return "", errors.New("quantity is required")
	}

	sideType, err := futuresSideFromString(side)
	if err != nil {
		return "", err
	}

	client := futures.NewClient(apiKey, secretKey)
	if setLeverage {
		if leverage < 1 || leverage > 125 {
			return "", errors.New("leverage must be between 1 and 125")
		}
		if _, err := client.NewChangeLeverageService().Symbol(symbol).Leverage(leverage).Do(ctx); err != nil {
			return "", fmt.Errorf("binance change leverage api error: %w", err)
		}
	}

	order, err := client.NewCreateOrderService().
		Symbol(symbol).
		Side(sideType).
		Type(futures.OrderTypeMarket).
		Quantity(quantity).
		Do(ctx)
	if err != nil {
		return "", fmt.Errorf("binance create futures order api error: %w", err)
	}

	var leveragePart string
	if setLeverage {
		leveragePart = fmt.Sprintf(" leverage=%dx", leverage)
	}
	return fmt.Sprintf(
		"Futures position opened: symbol=%s side=%s quantity=%s order_id=%d status=%s%s",
		order.Symbol,
		order.Side,
		order.OrigQuantity,
		order.OrderID,
		order.Status,
		leveragePart,
	), nil
}

// BinanceFuturesCloseMarketPosition closes a USDT-M futures position with reduce-only MARKET order.
// If quantity is empty and hasQuantity is false, it closes the full net position size.
func BinanceFuturesCloseMarketPosition(
	ctx context.Context,
	apiKey, secretKey, symbol, quantity string,
	hasQuantity bool,
) (string, error) {
	apiKey = strings.TrimSpace(apiKey)
	secretKey = strings.TrimSpace(secretKey)
	symbol = strings.ToUpper(strings.TrimSpace(symbol))
	quantity = strings.TrimSpace(quantity)

	if apiKey == "" || secretKey == "" {
		return "", errors.New("missing BINANCE_API_KEY or BINANCE_SECRET_KEY")
	}
	if symbol == "" {
		return "", errors.New("symbol is required")
	}

	client := futures.NewClient(apiKey, secretKey)
	sideType, closeQty, err := determineCloseOrder(ctx, client, symbol, quantity, hasQuantity)
	if err != nil {
		return "", err
	}
	if closeQty == "" {
		return fmt.Sprintf("No open net futures position found for %s.", symbol), nil
	}

	order, err := client.NewCreateOrderService().
		Symbol(symbol).
		Side(sideType).
		Type(futures.OrderTypeMarket).
		ReduceOnly(true).
		Quantity(closeQty).
		Do(ctx)
	if err != nil {
		return "", fmt.Errorf("binance close futures order api error: %w", err)
	}

	return fmt.Sprintf(
		"Futures position close submitted: symbol=%s side=%s quantity=%s order_id=%d status=%s reduce_only=true",
		order.Symbol,
		order.Side,
		order.OrigQuantity,
		order.OrderID,
		order.Status,
	), nil
}

func openFuturesPositionHandler(apiKey, secretKey string) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		confirm, err := request.RequireBool("confirm")
		if err != nil || !confirm {
			return mcp.NewToolResultError("open_futures_position requires confirm=true"), nil
		}

		symbol, err := request.RequireString("symbol")
		if err != nil {
			return mcp.NewToolResultError("symbol is required"), nil
		}
		side, err := request.RequireString("side")
		if err != nil {
			return mcp.NewToolResultError("side is required"), nil
		}
		quantity, err := request.RequireString("quantity")
		if err != nil {
			return mcp.NewToolResultError("quantity is required"), nil
		}
		leverage := request.GetInt("leverage", 0)

		ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
		defer cancel()

		result, err := BinanceFuturesOpenMarketPosition(
			ctx,
			apiKey,
			secretKey,
			symbol,
			side,
			quantity,
			leverage,
			leverage > 0,
		)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(result), nil
	}
}

func closeFuturesPositionHandler(apiKey, secretKey string) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		confirm, err := request.RequireBool("confirm")
		if err != nil || !confirm {
			return mcp.NewToolResultError("close_futures_position requires confirm=true"), nil
		}

		symbol, err := request.RequireString("symbol")
		if err != nil {
			return mcp.NewToolResultError("symbol is required"), nil
		}
		quantity := strings.TrimSpace(request.GetString("quantity", ""))

		ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
		defer cancel()

		result, err := BinanceFuturesCloseMarketPosition(ctx, apiKey, secretKey, symbol, quantity, quantity != "")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(result), nil
	}
}

func determineCloseOrder(
	ctx context.Context,
	client *futures.Client,
	symbol, quantity string,
	hasQuantity bool,
) (futures.SideType, string, error) {
	positions, err := client.NewGetPositionRiskService().Symbol(symbol).Do(ctx)
	if err != nil {
		return "", "", fmt.Errorf("binance get position risk api error: %w", err)
	}
	if len(positions) == 0 {
		return "", "", fmt.Errorf("symbol %s not found in futures positions", symbol)
	}

	netQty := 0.0
	for _, pos := range positions {
		if !strings.EqualFold(pos.Symbol, symbol) {
			continue
		}
		amt, err := strconv.ParseFloat(strings.TrimSpace(pos.PositionAmt), 64)
		if err != nil {
			return "", "", fmt.Errorf("invalid position amount for %s: %q", symbol, pos.PositionAmt)
		}
		netQty += amt
	}

	if math.Abs(netQty) < 1e-12 {
		return "", "", nil
	}

	side := futures.SideTypeSell
	if netQty < 0 {
		side = futures.SideTypeBuy
	}

	if hasQuantity {
		parsedQty, err := strconv.ParseFloat(quantity, 64)
		if err != nil || parsedQty <= 0 {
			return "", "", errors.New("quantity must be a positive decimal value")
		}
		if parsedQty > math.Abs(netQty) {
			return "", "", fmt.Errorf(
				"quantity %.12g exceeds open position %.12g for %s",
				parsedQty,
				math.Abs(netQty),
				symbol,
			)
		}
		return side, quantity, nil
	}

	return side, strconv.FormatFloat(math.Abs(netQty), 'f', -1, 64), nil
}

func futuresSideFromString(side string) (futures.SideType, error) {
	switch strings.ToUpper(strings.TrimSpace(side)) {
	case "BUY", "LONG":
		return futures.SideTypeBuy, nil
	case "SELL", "SHORT":
		return futures.SideTypeSell, nil
	default:
		return "", errors.New("side must be LONG/SHORT or BUY/SELL")
	}
}

func normalizeOrderBookLimit(limit int) (int, error) {
	if limit == 0 {
		return 10, nil
	}
	if limit < 0 {
		return 0, errors.New("limit must be positive")
	}
	if limit > 100 {
		return 0, errors.New("limit must be <= 100")
	}
	return limit, nil
}

func isZeroAmount(amount string) bool {
	trimmed := strings.TrimSpace(amount)
	if trimmed == "" {
		return true
	}
	trimmed = strings.TrimLeft(trimmed, "+")
	trimmed = strings.TrimRight(trimmed, "0")
	trimmed = strings.TrimRight(trimmed, ".")
	return trimmed == "" || trimmed == "0"
}
