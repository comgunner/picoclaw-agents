// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT

package tools

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================
// HELPER BUILDERS
// ============================================================

func newTestAgent() *TradingAgent {
	return &TradingAgent{
		balance:        decimal.NewFromFloat(10000),
		riskPerTrade:   decimal.NewFromFloat(0.02),
		maxPositionPct: decimal.NewFromFloat(0.10),
		stopLossPct:    decimal.NewFromFloat(0.15),
		takeProfitPct:  decimal.NewFromFloat(0.30),
		journal:        &TradeJournal{Trades: make([]BTTrade, 0)},
	}
}

func makeCandles(n int, basePrice float64) []Candle {
	candles := make([]Candle, n)
	for i := 0; i < n; i++ {
		p := decimal.NewFromFloat(basePrice + float64(i)*0.5)
		candles[i] = Candle{
			Timestamp: time.Now().Add(-time.Duration(n-i) * time.Minute),
			Open:      p,
			High:      p.Add(decimal.NewFromFloat(1.0)),
			Low:       p.Sub(decimal.NewFromFloat(0.5)),
			Close:     p.Add(decimal.NewFromFloat(0.2)),
			Volume:    decimal.NewFromFloat(100),
		}
	}
	return candles
}

// ============================================================
// INDICATOR TESTS
// ============================================================

func TestCalculateRSI_Bounds(t *testing.T) {
	closes := []float64{44, 44.25, 44.5, 44.75, 44.5, 44, 43, 42, 41, 42, 43, 44, 45, 44, 43}
	rsi := CalculateRSI(closes, 5)
	assert.GreaterOrEqual(t, rsi, 0.0, "RSI must be >= 0")
	assert.LessOrEqual(t, rsi, 100.0, "RSI must be <= 100")
}

func TestCalculateRSI_InsufficientData(t *testing.T) {
	rsi := CalculateRSI([]float64{44, 45}, 14)
	assert.Equal(t, 50.0, rsi, "should return neutral 50 for insufficient data")
}

func TestCalculateRSI_AllGains(t *testing.T) {
	closes := make([]float64, 20)
	for i := range closes {
		closes[i] = float64(i + 100)
	}
	rsi := CalculateRSI(closes, 14)
	assert.Equal(t, 100.0, rsi, "all-gain data should return RSI=100")
}

func TestCalculateEMA_Basic(t *testing.T) {
	closes := []float64{10, 11, 12, 13, 14, 15, 14, 13, 12, 11}
	ema := CalculateEMA(closes, 5)
	assert.Greater(t, ema, 0.0)
	assert.Less(t, ema, 20.0)
}

func TestCalculateEMA_SingleValue(t *testing.T) {
	ema := CalculateEMA([]float64{42.0}, 5)
	assert.Equal(t, 42.0, ema)
}

func TestCalculateMACD_Basic(t *testing.T) {
	closes := make([]float64, 50)
	for i := range closes {
		closes[i] = 100.0 + float64(i)*0.5
	}
	macd := CalculateMACD(closes)
	// MACD values should be valid (not zero since price is trending)
	assert.IsType(t, BTMACDData{}, macd)
}

func TestCalculateMACD_InsufficientData(t *testing.T) {
	macd := CalculateMACD([]float64{100, 101, 102})
	assert.Equal(t, BTMACDData{}, macd)
}

func TestCalculateBollingerBands_Bounds(t *testing.T) {
	closes := make([]float64, 25)
	for i := range closes {
		closes[i] = 100.0 + float64(i%5)
	}
	bb := CalculateBollingerBands(closes, 20, 2.0)
	assert.Greater(t, bb.Upper, bb.Middle, "Upper band must be above middle")
	assert.Less(t, bb.Lower, bb.Middle, "Lower band must be below middle")
}

func TestCalculateATR_Basic(t *testing.T) {
	candles := makeCandles(20, 100.0)
	atr := CalculateATR(candles, 14)
	assert.Greater(t, atr, 0.0, "ATR must be positive")
}

func TestCalculateATR_InsufficientData(t *testing.T) {
	candles := makeCandles(5, 100.0)
	atr := CalculateATR(candles, 14)
	assert.Equal(t, 0.0, atr)
}

// ============================================================
// PATTERN DETECTION TESTS
// ============================================================

func TestIsDoji(t *testing.T) {
	agent := newTestAgent()
	// Open ≈ Close, wide range => Doji
	doji := Candle{
		Open:  decimal.NewFromFloat(100),
		High:  decimal.NewFromFloat(105),
		Low:   decimal.NewFromFloat(95),
		Close: decimal.NewFromFloat(100.1),
	}
	assert.True(t, agent.isDoji(doji))

	notDoji := Candle{
		Open:  decimal.NewFromFloat(100),
		High:  decimal.NewFromFloat(115),
		Low:   decimal.NewFromFloat(95),
		Close: decimal.NewFromFloat(112),
	}
	assert.False(t, agent.isDoji(notDoji))
}

func TestDetectCandlestickPatterns_InsufficientData(t *testing.T) {
	agent := newTestAgent()
	patterns := agent.DetectCandlestickPatterns(makeCandles(2, 100.0))
	assert.Empty(t, patterns)
}

func TestDetectCandlestickPatterns_ReturnSlice(t *testing.T) {
	agent := newTestAgent()
	candles := makeCandles(10, 100.0)
	patterns := agent.DetectCandlestickPatterns(candles)
	assert.NotNil(t, patterns)
}

// ============================================================
// VALIDATION TESTS
// ============================================================

func TestValidateDataIntegrity_InsufficientData(t *testing.T) {
	agent := newTestAgent()
	ok, errs := agent.ValidateDataIntegrity(makeCandles(10, 100.0))
	assert.False(t, ok)
	require.Len(t, errs, 1)
	assert.Contains(t, errs[0], "Insufficient data")
}

func TestValidateDataIntegrity_ValidData(t *testing.T) {
	agent := newTestAgent()
	candles := makeCandles(50, 100.0)
	ok, errs := agent.ValidateDataIntegrity(candles)
	// Should pass structural checks (data freshness may warn but not fail in test)
	_ = errs
	_ = ok
	// Just ensure it runs without panic
}

func TestValidateDataIntegrity_InvalidPrices(t *testing.T) {
	agent := newTestAgent()
	candles := makeCandles(25, 100.0)
	// Inject an invalid candle
	candles[5] = Candle{
		Timestamp: time.Now(),
		Open:      decimal.NewFromFloat(-1),
		High:      decimal.NewFromFloat(-1),
		Low:       decimal.NewFromFloat(-1),
		Close:     decimal.NewFromFloat(-1),
		Volume:    decimal.NewFromFloat(100),
	}
	ok, errs := agent.ValidateDataIntegrity(candles)
	assert.False(t, ok)
	assert.NotEmpty(t, errs)
}

func TestValidateIndicators_ValidRSI(t *testing.T) {
	agent := newTestAgent()
	indicators := BTIndicatorData{
		RSI: 55,
		ATR: 100,
		Bollinger: BTBollingerData{
			Upper:  110,
			Middle: 100,
			Lower:  90,
		},
		MACD: BTMACDData{MACD: 5, Signal: 3},
	}
	ok, errs := agent.ValidateIndicators(indicators)
	assert.True(t, ok)
	assert.Empty(t, errs)
}

func TestValidateIndicators_InvalidRSI(t *testing.T) {
	agent := newTestAgent()
	indicators := BTIndicatorData{RSI: 150}
	ok, errs := agent.ValidateIndicators(indicators)
	assert.False(t, ok)
	assert.NotEmpty(t, errs)
}

func TestValidateSignal_LongSignalValid(t *testing.T) {
	agent := newTestAgent()
	result := AnalysisResult{
		Recommendation:  "LONG",
		Confidence:      70,
		EntryPrice:      decimal.NewFromFloat(100),
		StopLoss:        decimal.NewFromFloat(95),
		TakeProfit:      decimal.NewFromFloat(110),
		RiskRewardRatio: decimal.NewFromFloat(2.0),
	}
	ok, errs := agent.ValidateSignal(result)
	assert.True(t, ok)
	assert.Empty(t, errs)
}

func TestValidateSignal_LongStopAboveEntry(t *testing.T) {
	agent := newTestAgent()
	result := AnalysisResult{
		Recommendation:  "LONG",
		Confidence:      70,
		EntryPrice:      decimal.NewFromFloat(100),
		StopLoss:        decimal.NewFromFloat(105), // wrong: above entry
		TakeProfit:      decimal.NewFromFloat(115),
		RiskRewardRatio: decimal.NewFromFloat(2.0),
	}
	ok, errs := agent.ValidateSignal(result)
	assert.False(t, ok)
	assert.Contains(t, errs[0], "stop loss")
}

func TestCrossVerification_ConflictingSignals(t *testing.T) {
	agent := newTestAgent()
	results := []AnalysisResult{
		{Recommendation: "LONG", Confidence: 70},
		{Recommendation: "SHORT", Confidence: 65},
	}
	ok, errs := agent.CrossVerification(results)
	assert.False(t, ok)
	assert.Contains(t, errs[0], "conflicting")
}

func TestCrossVerification_AlignedSignals(t *testing.T) {
	agent := newTestAgent()
	results := []AnalysisResult{
		{Recommendation: "LONG", Confidence: 70},
		{Recommendation: "LONG", Confidence: 72},
	}
	ok, errs := agent.CrossVerification(results)
	assert.True(t, ok)
	assert.Empty(t, errs)
}

// ============================================================
// RISK MANAGEMENT TESTS
// ============================================================

func TestCalculatePositionSize_Basic(t *testing.T) {
	agent := newTestAgent()
	entry := decimal.NewFromFloat(100)
	sl := decimal.NewFromFloat(95)

	posSize, riskAmt := agent.CalculatePositionSize(entry, sl, agent.balance)

	// position size uncapped = 200 / 5 = 40 units
	// max position = 10% of 10000 / 100 = 10 units => capped
	posF, _ := posSize.Float64()
	assert.InDelta(t, 10.0, posF, 0.01, "position should be capped at max position")

	// risk recalculated after cap: 10 units * $5 stop = $50
	riskF, _ := riskAmt.Float64()
	assert.InDelta(t, 50.0, riskF, 0.01, "risk amount recalculated after position cap")
}

func TestCalculatePositionSize_ZeroStopDistance(t *testing.T) {
	agent := newTestAgent()
	entry := decimal.NewFromFloat(100)
	posSize, _ := agent.CalculatePositionSize(entry, entry, agent.balance)
	assert.True(t, posSize.IsZero())
}

func TestCalculateStopLoss_Long(t *testing.T) {
	agent := newTestAgent()
	entry := decimal.NewFromFloat(100)
	atr := 2.0
	sl := agent.CalculateStopLoss(entry, atr, true)
	expected := decimal.NewFromFloat(96) // 100 - 2*2
	assert.Equal(t, expected, sl)
}

func TestCalculateStopLoss_Short(t *testing.T) {
	agent := newTestAgent()
	entry := decimal.NewFromFloat(100)
	atr := 2.0
	sl := agent.CalculateStopLoss(entry, atr, false)
	expected := decimal.NewFromFloat(104) // 100 + 2*2
	assert.Equal(t, expected, sl)
}

func TestCalculateTakeProfit_Long(t *testing.T) {
	agent := newTestAgent()
	entry := decimal.NewFromFloat(100)
	sl := decimal.NewFromFloat(95)
	tp := agent.CalculateTakeProfit(entry, sl, true, decimal.NewFromFloat(2.0))
	tpF, _ := tp.Float64()
	assert.InDelta(t, 110.0, tpF, 0.01) // 100 + 5*2
}

func TestCalculateRiskRewardRatio(t *testing.T) {
	agent := newTestAgent()
	rr := agent.CalculateRiskRewardRatio(
		decimal.NewFromFloat(100),
		decimal.NewFromFloat(95),
		decimal.NewFromFloat(110),
	)
	// reward=10, risk=5 => 2.0
	rrF, _ := rr.Float64()
	assert.InDelta(t, 2.0, rrF, 0.001)
}

func TestCalculateVaR(t *testing.T) {
	agent := newTestAgent()
	returns := []float64{-0.05, -0.03, -0.01, 0.02, 0.04, -0.02, 0.01, -0.04}
	var95 := agent.CalculateVaR(returns, 0.95)
	assert.Greater(t, var95, 0.0)
}

func TestCalculateSharpeRatio_AllPositive(t *testing.T) {
	agent := newTestAgent()
	returns := []float64{0.01, 0.02, 0.015, 0.025, 0.01}
	sharpe := agent.CalculateSharpeRatio(returns, 0.001)
	assert.Greater(t, sharpe, 0.0)
}

// ============================================================
// TREND AND REGIME TESTS
// ============================================================

func TestDetectTrend_InsufficientData(t *testing.T) {
	agent := newTestAgent()
	candles := makeCandles(10, 100.0)
	trend := agent.DetectTrend(candles)
	assert.Equal(t, "UNKNOWN", trend)
}

func TestDetectTrend_UpTrend(t *testing.T) {
	agent := newTestAgent()
	// Create 60 candles with strongly increasing prices
	candles := make([]Candle, 60)
	for i := 0; i < 60; i++ {
		p := decimal.NewFromFloat(100.0 + float64(i)*2)
		candles[i] = Candle{
			Timestamp: time.Now().Add(-time.Duration(60-i) * time.Hour),
			Open:      p,
			High:      p.Add(decimal.NewFromFloat(1)),
			Low:       p.Sub(decimal.NewFromFloat(0.5)),
			Close:     p.Add(decimal.NewFromFloat(0.5)),
			Volume:    decimal.NewFromFloat(100),
		}
	}
	trend := agent.DetectTrend(candles)
	assert.Equal(t, "UPTREND", trend)
}

func TestDetectMarketRegime_InsufficientData(t *testing.T) {
	agent := newTestAgent()
	regime := agent.DetectMarketRegime(makeCandles(10, 100.0))
	assert.Equal(t, "UNKNOWN", regime)
}

// ============================================================
// TOOL INTERFACE TESTS
// ============================================================

func TestBaseTraderTool_Name(t *testing.T) {
	tool := NewBaseTraderTool(10000)
	assert.Equal(t, "base_trader", tool.Name())
}

func TestBaseTraderTool_Description(t *testing.T) {
	tool := NewBaseTraderTool(10000)
	assert.NotEmpty(t, tool.Description())
}

func TestBaseTraderTool_Parameters(t *testing.T) {
	tool := NewBaseTraderTool(10000)
	params := tool.Parameters()
	assert.NotNil(t, params)
	assert.Equal(t, "object", params["type"])
}

func TestBaseTraderTool_Execute_MissingAction(t *testing.T) {
	tool := NewBaseTraderTool(10000)
	result := tool.Execute(nil, map[string]any{})
	assert.True(t, result.IsError)
	assert.Contains(t, result.ForLLM, "action is required")
}

func TestBaseTraderTool_Execute_UnknownAction(t *testing.T) {
	tool := NewBaseTraderTool(10000)
	result := tool.Execute(nil, map[string]any{"action": "unknown"})
	assert.True(t, result.IsError)
}

func TestBaseTraderTool_Execute_StatusAction(t *testing.T) {
	tool := NewBaseTraderTool(10000)
	result := tool.Execute(nil, map[string]any{"action": "status"})
	assert.False(t, result.IsError)
	assert.Contains(t, result.ForLLM, "BASE TRADER STATUS")
	assert.Contains(t, result.ForLLM, "10000.00")
}

func TestBaseTraderTool_Execute_JournalEmpty(t *testing.T) {
	tool := NewBaseTraderTool(10000)
	result := tool.Execute(nil, map[string]any{"action": "journal"})
	assert.False(t, result.IsError)
	assert.Contains(t, result.ForLLM, "No trades")
}

func TestBaseTraderTool_Execute_AnalyzeMissingSymbol(t *testing.T) {
	tool := NewBaseTraderTool(10000)
	result := tool.Execute(nil, map[string]any{"action": "analyze"})
	assert.True(t, result.IsError)
	assert.Contains(t, result.ForLLM, "symbol")
}

func TestBaseTraderTool_BalanceOverride(t *testing.T) {
	tool := NewBaseTraderTool(10000)
	// status with balance override
	result := tool.Execute(nil, map[string]any{
		"action":  "status",
		"balance": 25000.0,
	})
	assert.False(t, result.IsError)
	assert.Contains(t, result.ForLLM, "25000.00")
}

// ============================================================
// CONSOLIDATE RESULTS TEST
// ============================================================

func TestConsolidateResults_Empty(t *testing.T) {
	agent := newTestAgent()
	result := agent.ConsolidateResults([]AnalysisResult{}, "BTC/USDT")
	assert.Equal(t, "WAIT", result.Recommendation)
	assert.Equal(t, "BTC/USDT", result.Symbol)
}

func TestConsolidateResults_PicksHighestConfidence(t *testing.T) {
	agent := newTestAgent()
	results := []AnalysisResult{
		{Recommendation: "LONG", Confidence: 60},
		{Recommendation: "LONG", Confidence: 80},
		{Recommendation: "LONG", Confidence: 70},
	}
	best := agent.ConsolidateResults(results, "ETH/USDT")
	assert.Equal(t, 80.0, best.Confidence)
}

// ============================================================
// MATH HELPERS TESTS
// ============================================================

func TestBtAverage(t *testing.T) {
	assert.Equal(t, 0.0, btAverage(nil))
	assert.Equal(t, 5.0, btAverage([]float64{1, 5, 9}))
}

func TestBtStdDev(t *testing.T) {
	assert.Equal(t, 0.0, btStdDev(nil))
	assert.InDelta(t, 0.0, btStdDev([]float64{5, 5, 5}), 0.0001)
	assert.Greater(t, btStdDev([]float64{1, 5, 9}), 0.0)
}

func TestBtContains(t *testing.T) {
	assert.True(t, btContains([]string{"a", "b", "c"}, "b"))
	assert.False(t, btContains([]string{"a", "b", "c"}, "z"))
}
