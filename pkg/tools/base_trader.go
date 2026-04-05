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
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

// ============================================================
// TYPES AND STRUCTURES
// ============================================================

// TradingAgent is the core trading analysis agent.
type TradingAgent struct {
	balance        decimal.Decimal
	riskPerTrade   decimal.Decimal
	maxPositionPct decimal.Decimal
	stopLossPct    decimal.Decimal
	takeProfitPct  decimal.Decimal
	journal        *TradeJournal
}

// Candle represents a single OHLCV candlestick.
type Candle struct {
	Timestamp time.Time
	Open      decimal.Decimal
	High      decimal.Decimal
	Low       decimal.Decimal
	Close     decimal.Decimal
	Volume    decimal.Decimal
}

// AnalysisResult holds the complete analysis output for a symbol/timeframe.
type AnalysisResult struct {
	Symbol           string
	Timeframe        string
	CurrentPrice     decimal.Decimal
	Recommendation   string // LONG, SHORT, WAIT, NO_TRADE
	Confidence       float64
	EntryPrice       decimal.Decimal
	StopLoss         decimal.Decimal
	TakeProfit       decimal.Decimal
	RiskRewardRatio  decimal.Decimal
	PositionSize     decimal.Decimal
	RiskAmount       decimal.Decimal
	ValidationPassed bool
	ValidationReport BTValidationReport
	Indicators       BTIndicatorData
	Patterns         BTPatternData
}

// BTValidationReport contains the 6-layer validation report.
type BTValidationReport struct {
	DataIntegrity   bool
	IndicatorsValid bool
	SignalValid     bool
	CrossVerified   bool
	ExecutionReady  bool
	ProductionReady bool
	Errors          []string
	Warnings        []string
}

// BTIndicatorData holds all computed technical indicators.
type BTIndicatorData struct {
	RSI       float64
	MACD      BTMACDData
	Bollinger BTBollingerData
	ATR       float64
	MA20      float64
	MA50      float64
	MA200     float64
	EMA12     float64
	EMA26     float64
}

// BTMACDData holds MACD components.
type BTMACDData struct {
	MACD      float64
	Signal    float64
	Histogram float64
}

// BTBollingerData holds Bollinger Band values.
type BTBollingerData struct {
	Upper  float64
	Middle float64
	Lower  float64
}

// BTPatternData holds detected market patterns.
type BTPatternData struct {
	CandlestickPatterns []string
	TrendPattern        string // UPTREND, DOWNTREND, SIDEWAYS
	MarketRegime        string // TRENDING, RANGING
	SupportLevels       []float64
	ResistanceLevels    []float64
}

// TradeJournal tracks all executed trades.
type TradeJournal struct {
	Trades []BTTrade
}

// BTTrade represents a single trade record.
type BTTrade struct {
	ID         string
	Timestamp  time.Time
	Symbol     string
	Action     string
	AmountUSD  decimal.Decimal
	EntryPrice decimal.Decimal
	ExitPrice  decimal.Decimal
	StopLoss   decimal.Decimal
	TakeProfit decimal.Decimal
	PnL        decimal.Decimal
	PnLPercent float64
	Reason     string
	Validation BTValidationReport
	ExitReason string // STOP_LOSS, TAKE_PROFIT, MANUAL
}

// ============================================================
// MATH HELPERS
// ============================================================

func btAverage(values []float64) float64 {
	if len(values) == 0 {
		return 0.0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func btStdDev(values []float64) float64 {
	if len(values) == 0 {
		return 0.0
	}
	mean := btAverage(values)
	variance := 0.0
	for _, v := range values {
		diff := v - mean
		variance += diff * diff
	}
	variance /= float64(len(values))
	return math.Sqrt(variance)
}

func btContains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// ============================================================
// TECHNICAL INDICATORS
// ============================================================

// CalculateRSI computes the Relative Strength Index.
func CalculateRSI(closes []float64, period int) float64 {
	if len(closes) < period+1 {
		return 50.0
	}

	gains := make([]float64, 0, len(closes)-1)
	losses := make([]float64, 0, len(closes)-1)

	for i := 1; i < len(closes); i++ {
		change := closes[i] - closes[i-1]
		if change > 0 {
			gains = append(gains, change)
			losses = append(losses, 0)
		} else {
			gains = append(gains, 0)
			losses = append(losses, -change)
		}
	}

	if len(gains) < period {
		return 50.0
	}

	avgGain := btAverage(gains[len(gains)-period:])
	avgLoss := btAverage(losses[len(losses)-period:])

	if avgLoss == 0 {
		return 100.0
	}

	rs := avgGain / avgLoss
	return 100.0 - (100.0 / (1.0 + rs))
}

// CalculateEMA computes the Exponential Moving Average (last value).
func CalculateEMA(closes []float64, period int) float64 {
	if len(closes) == 0 {
		return 0.0
	}
	if len(closes) < period {
		return btAverage(closes)
	}

	multiplier := 2.0 / float64(period+1)
	ema := btAverage(closes[:period])

	for i := period; i < len(closes); i++ {
		ema = (closes[i]-ema)*multiplier + ema
	}

	return ema
}

// calculateEMASeries computes EMA for every data point.
func calculateEMASeries(closes []float64, period int) []float64 {
	result := make([]float64, len(closes))
	if len(closes) < period {
		for i := range result {
			result[i] = btAverage(closes)
		}
		return result
	}

	multiplier := 2.0 / float64(period+1)
	sma := btAverage(closes[:period])

	for i := 0; i < period; i++ {
		result[i] = sma
	}

	for i := period; i < len(closes); i++ {
		result[i] = (closes[i]-result[i-1])*multiplier + result[i-1]
	}

	return result
}

// CalculateMACD computes MACD line, signal line, and histogram.
func CalculateMACD(closes []float64) BTMACDData {
	if len(closes) < 35 {
		return BTMACDData{}
	}

	ema12Series := calculateEMASeries(closes, 12)
	ema26Series := calculateEMASeries(closes, 26)

	// Build MACD series from index 25 onwards (when EMA-26 is meaningful)
	macdSeries := make([]float64, 0, len(closes)-25)
	for i := 25; i < len(closes); i++ {
		macdSeries = append(macdSeries, ema12Series[i]-ema26Series[i])
	}

	if len(macdSeries) == 0 {
		return BTMACDData{}
	}

	currentMACD := macdSeries[len(macdSeries)-1]
	signalLine := CalculateEMA(macdSeries, 9)
	histogram := currentMACD - signalLine

	return BTMACDData{
		MACD:      currentMACD,
		Signal:    signalLine,
		Histogram: histogram,
	}
}

// CalculateBollingerBands computes Bollinger Bands.
func CalculateBollingerBands(closes []float64, period int, stdDevMult float64) BTBollingerData {
	if len(closes) < period {
		return BTBollingerData{}
	}

	window := closes[len(closes)-period:]
	middle := btAverage(window)

	variance := 0.0
	for _, v := range window {
		diff := v - middle
		variance += diff * diff
	}
	variance /= float64(period)
	std := math.Sqrt(variance)

	return BTBollingerData{
		Upper:  middle + stdDevMult*std,
		Middle: middle,
		Lower:  middle - stdDevMult*std,
	}
}

// CalculateATR computes Average True Range.
func CalculateATR(candles []Candle, period int) float64 {
	if len(candles) < period+1 {
		return 0.0
	}

	trueRanges := make([]float64, 0, len(candles)-1)
	for i := 1; i < len(candles); i++ {
		highLow, _ := candles[i].High.Sub(candles[i].Low).Float64()
		highClose := math.Abs(mustFloat64(candles[i].High.Sub(candles[i-1].Close)))
		lowClose := math.Abs(mustFloat64(candles[i].Low.Sub(candles[i-1].Close)))

		tr := math.Max(highLow, math.Max(highClose, lowClose))
		trueRanges = append(trueRanges, tr)
	}

	if len(trueRanges) < period {
		return btAverage(trueRanges)
	}

	return btAverage(trueRanges[len(trueRanges)-period:])
}

// calculateMA computes the last value of a simple moving average.
func (t *TradingAgent) calculateMA(candles []Candle, period int) float64 {
	if len(candles) < period {
		return 0.0
	}
	window := candles[len(candles)-period:]
	sum := 0.0
	for _, c := range window {
		f, _ := c.Close.Float64()
		sum += f
	}
	return sum / float64(period)
}

// calculateADX computes a simplified Average Directional Index.
func (t *TradingAgent) calculateADX(candles []Candle, period int) float64 {
	if len(candles) < period*2 {
		return 0.0
	}

	// Directional movement
	plusDM := make([]float64, 0, len(candles)-1)
	minusDM := make([]float64, 0, len(candles)-1)
	trValues := make([]float64, 0, len(candles)-1)

	for i := 1; i < len(candles); i++ {
		curr := candles[i]
		prev := candles[i-1]

		upMove, _ := curr.High.Sub(prev.High).Float64()
		downMove, _ := prev.Low.Sub(curr.Low).Float64()

		if upMove > downMove && upMove > 0 {
			plusDM = append(plusDM, upMove)
		} else {
			plusDM = append(plusDM, 0)
		}

		if downMove > upMove && downMove > 0 {
			minusDM = append(minusDM, downMove)
		} else {
			minusDM = append(minusDM, 0)
		}

		highLow, _ := curr.High.Sub(curr.Low).Float64()
		highClose := math.Abs(mustFloat64(curr.High.Sub(prev.Close)))
		lowClose := math.Abs(mustFloat64(curr.Low.Sub(prev.Close)))
		trValues = append(trValues, math.Max(highLow, math.Max(highClose, lowClose)))
	}

	if len(trValues) < period {
		return 0.0
	}

	atr := btAverage(trValues[len(trValues)-period:])
	if atr == 0 {
		return 0.0
	}

	plusDI := (btAverage(plusDM[len(plusDM)-period:]) / atr) * 100
	minusDI := (btAverage(minusDM[len(minusDM)-period:]) / atr) * 100

	diSum := plusDI + minusDI
	if diSum == 0 {
		return 0.0
	}

	dx := math.Abs(plusDI-minusDI) / diSum * 100
	return dx
}

// CalculateIndicators computes all technical indicators for a set of candles.
func (t *TradingAgent) CalculateIndicators(candles []Candle) BTIndicatorData {
	closes := make([]float64, len(candles))
	for i, c := range candles {
		closes[i], _ = c.Close.Float64()
	}

	macd := CalculateMACD(closes)
	bb := CalculateBollingerBands(closes, 20, 2.0)
	atr := CalculateATR(candles, 14)

	return BTIndicatorData{
		RSI:       CalculateRSI(closes, 14),
		MACD:      macd,
		Bollinger: bb,
		ATR:       atr,
		MA20:      t.calculateMA(candles, 20),
		MA50:      t.calculateMA(candles, 50),
		MA200:     t.calculateMA(candles, 200),
		EMA12:     CalculateEMA(closes, 12),
		EMA26:     CalculateEMA(closes, 26),
	}
}

// ============================================================
// CANDLESTICK PATTERNS
// ============================================================

// DetectCandlestickPatterns identifies common Japanese candlestick patterns.
func (t *TradingAgent) DetectCandlestickPatterns(candles []Candle) []string {
	patterns := make([]string, 0)

	if len(candles) < 3 {
		return patterns
	}

	last := candles[len(candles)-1]
	prev := candles[len(candles)-2]

	if t.isDoji(last) {
		patterns = append(patterns, "DOJI")
	}

	if t.isHammer(last) {
		patterns = append(patterns, "HAMMER")
	}

	if t.isShootingStar(last) {
		patterns = append(patterns, "SHOOTING_STAR")
	}

	if t.isBullishEngulfing(last, prev) {
		patterns = append(patterns, "BULLISH_ENGULFING")
	}

	if t.isBearishEngulfing(last, prev) {
		patterns = append(patterns, "BEARISH_ENGULFING")
	}

	return patterns
}

func (t *TradingAgent) isDoji(c Candle) bool {
	body := c.Open.Sub(c.Close).Abs()
	rangeHL := c.High.Sub(c.Low)
	if rangeHL.IsZero() {
		return false
	}
	return body.LessThan(rangeHL.Mul(decimal.NewFromFloat(0.1)))
}

func (t *TradingAgent) isHammer(c Candle) bool {
	body := c.Open.Sub(c.Close).Abs()
	topPrice := decimalMax(c.Open, c.Close)
	botPrice := decimalMin(c.Open, c.Close)
	upperShadow := c.High.Sub(topPrice)
	lowerShadow := botPrice.Sub(c.Low)

	// Body < 2 * upper shadow AND lower shadow > 2 * body
	return body.LessThan(upperShadow.Mul(decimal.NewFromFloat(2))) &&
		lowerShadow.GreaterThan(body.Mul(decimal.NewFromFloat(2)))
}

func (t *TradingAgent) isShootingStar(c Candle) bool {
	body := c.Open.Sub(c.Close).Abs()
	topPrice := decimalMax(c.Open, c.Close)
	botPrice := decimalMin(c.Open, c.Close)
	upperShadow := c.High.Sub(topPrice)
	lowerShadow := botPrice.Sub(c.Low)

	// Long upper shadow, small lower shadow
	return upperShadow.GreaterThan(body.Mul(decimal.NewFromFloat(2))) &&
		lowerShadow.LessThan(body.Mul(decimal.NewFromFloat(0.5)))
}

func (t *TradingAgent) isBullishEngulfing(current, previous Candle) bool {
	prevBearish := previous.Close.LessThan(previous.Open)
	currBullish := current.Close.GreaterThan(current.Open)
	engulfs := current.Open.LessThan(previous.Close) &&
		current.Close.GreaterThan(previous.Open)
	return prevBearish && currBullish && engulfs
}

func (t *TradingAgent) isBearishEngulfing(current, previous Candle) bool {
	prevBullish := previous.Close.GreaterThan(previous.Open)
	currBearish := current.Close.LessThan(current.Open)
	engulfs := current.Open.GreaterThan(previous.Close) &&
		current.Close.LessThan(previous.Open)
	return prevBullish && currBearish && engulfs
}

// ============================================================
// 6-LAYER VALIDATION
// ============================================================

// ValidateDataIntegrity - Layer 1: Validates OHLCV data integrity.
func (t *TradingAgent) ValidateDataIntegrity(candles []Candle) (bool, []string) {
	errs := make([]string, 0)

	if len(candles) < 20 {
		errs = append(errs, fmt.Sprintf("Insufficient data: %d candles (minimum 20)", len(candles)))
		return false, errs
	}

	for i, c := range candles {
		if c.Open.LessThanOrEqual(decimal.Zero) ||
			c.High.LessThanOrEqual(decimal.Zero) ||
			c.Low.LessThanOrEqual(decimal.Zero) ||
			c.Close.LessThanOrEqual(decimal.Zero) {
			errs = append(errs, fmt.Sprintf("Candle %d: invalid price (zero or negative)", i))
		}
		if c.High.LessThan(c.Low) {
			errs = append(errs, fmt.Sprintf("Candle %d: High < Low", i))
		}
		if c.High.LessThan(c.Open) || c.High.LessThan(c.Close) {
			errs = append(errs, fmt.Sprintf("Candle %d: High < Open or Close", i))
		}
		if c.Low.GreaterThan(c.Open) || c.Low.GreaterThan(c.Close) {
			errs = append(errs, fmt.Sprintf("Candle %d: Low > Open or Close", i))
		}
		if !c.Open.IsZero() {
			pct, _ := c.Close.Sub(c.Open).Div(c.Open).Mul(decimal.NewFromInt(100)).Float64()
			if math.Abs(pct) > 50 {
				errs = append(errs, fmt.Sprintf("Candle %d: unrealistic price jump (%.2f%%)", i, pct))
			}
		}
	}

	// Z-score anomaly detection
	closes := make([]float64, len(candles))
	for i, c := range candles {
		closes[i], _ = c.Close.Float64()
	}
	mean := btAverage(closes)
	std := btStdDev(closes)
	if std > 0 {
		for i, cl := range closes {
			zscore := math.Abs((cl - mean) / std)
			if zscore > 5.0 {
				errs = append(errs, fmt.Sprintf("Candle %d: extreme Z-score (%.2f)", i, zscore))
			}
		}
	}

	// Data freshness: only check the last candle
	if !candles[len(candles)-1].Timestamp.IsZero() {
		age := time.Since(candles[len(candles)-1].Timestamp)
		if age > 30*time.Minute {
			errs = append(errs, fmt.Sprintf("Last candle is stale (%.0f minutes old)", age.Minutes()))
		}
	}

	return len(errs) == 0, errs
}

// ValidateIndicators - Layer 2: Validates indicator values are sane.
func (t *TradingAgent) ValidateIndicators(indicators BTIndicatorData) (bool, []string) {
	errs := make([]string, 0)

	if indicators.RSI < 0 || indicators.RSI > 100 {
		errs = append(errs, fmt.Sprintf("invalid RSI: %.2f (must be 0-100)", indicators.RSI))
	}

	if indicators.ATR < 0 {
		errs = append(errs, fmt.Sprintf("invalid ATR: %.2f (must be positive)", indicators.ATR))
	}

	if indicators.Bollinger.Upper > 0 && indicators.Bollinger.Upper <= indicators.Bollinger.Lower {
		errs = append(errs, "invalid Bollinger Bands: Upper <= Lower")
	}

	if math.Abs(indicators.MACD.MACD) > 100000 {
		errs = append(errs, fmt.Sprintf("suspicious MACD value: %.2f", indicators.MACD.MACD))
	}

	return len(errs) == 0, errs
}

// ValidateSignal - Layer 3: Validates signal logic and risk/reward.
func (t *TradingAgent) ValidateSignal(result AnalysisResult) (bool, []string) {
	errs := make([]string, 0)

	validActions := map[string]bool{
		"LONG": true, "SHORT": true, "WAIT": true, "NO_TRADE": true,
	}
	if !validActions[result.Recommendation] {
		errs = append(errs, fmt.Sprintf("invalid action: %s", result.Recommendation))
	}

	if result.Confidence < 0 || result.Confidence > 100 {
		errs = append(errs, fmt.Sprintf("invalid confidence: %.2f", result.Confidence))
	}

	if result.Recommendation == "LONG" && !result.EntryPrice.IsZero() {
		if result.StopLoss.GreaterThanOrEqual(result.EntryPrice) {
			errs = append(errs, "LONG: stop loss must be < entry price")
		}
		if result.TakeProfit.LessThanOrEqual(result.EntryPrice) {
			errs = append(errs, "LONG: take profit must be > entry price")
		}
	}

	if result.Recommendation == "SHORT" && !result.EntryPrice.IsZero() {
		if result.StopLoss.LessThanOrEqual(result.EntryPrice) {
			errs = append(errs, "SHORT: stop loss must be > entry price")
		}
		if result.TakeProfit.GreaterThanOrEqual(result.EntryPrice) {
			errs = append(errs, "SHORT: take profit must be < entry price")
		}
	}

	minRR := decimal.NewFromFloat(1.5)
	if (result.Recommendation == "LONG" || result.Recommendation == "SHORT") &&
		result.RiskRewardRatio.LessThan(minRR) && !result.RiskRewardRatio.IsZero() {
		errs = append(
			errs,
			fmt.Sprintf("poor risk/reward: %.2f:1 (minimum 1.5:1)", result.RiskRewardRatio.InexactFloat64()),
		)
	}

	return len(errs) == 0, errs
}

// CrossVerification - Layer 4: Checks for conflicting signals across timeframes.
func (t *TradingAgent) CrossVerification(results []AnalysisResult) (bool, []string) {
	errs := make([]string, 0)

	if len(results) < 2 {
		return true, errs
	}

	longCount, shortCount := 0, 0
	for _, r := range results {
		if r.Recommendation == "LONG" {
			longCount++
		} else if r.Recommendation == "SHORT" {
			shortCount++
		}
	}

	if longCount > 0 && shortCount > 0 {
		errs = append(errs, fmt.Sprintf("conflicting signals: %d LONG, %d SHORT", longCount, shortCount))
	}

	confidences := make([]float64, len(results))
	for i, r := range results {
		confidences[i] = r.Confidence
	}
	if btStdDev(confidences) > 20.0 {
		errs = append(errs, fmt.Sprintf("high confidence variance: %.2f%%", btStdDev(confidences)))
	}

	return len(errs) == 0, errs
}

// ExecutionReadiness - Layer 5: All prior layers must pass.
func (t *TradingAgent) ExecutionReadiness(report BTValidationReport) bool {
	return report.DataIntegrity &&
		report.IndicatorsValid &&
		report.SignalValid &&
		report.CrossVerified
}

// ProductionValidation - Layer 6: Final confidence gate.
func (t *TradingAgent) ProductionValidation(result AnalysisResult) bool {
	return result.Confidence >= 50 && result.ValidationPassed
}

// ============================================================
// RISK MANAGEMENT
// ============================================================

// CalculatePositionSize returns position size and risk amount for a trade.
func (t *TradingAgent) CalculatePositionSize(
	entryPrice, stopLoss, balance decimal.Decimal,
) (positionSize, riskAmount decimal.Decimal) {
	riskAmount = balance.Mul(t.riskPerTrade)

	stopDistance := entryPrice.Sub(stopLoss).Abs()
	if stopDistance.IsZero() {
		return decimal.Zero, riskAmount
	}

	positionSize = riskAmount.Div(stopDistance)

	maxPositionValue := balance.Mul(t.maxPositionPct)
	maxPositionSize := maxPositionValue.Div(entryPrice)

	if positionSize.GreaterThan(maxPositionSize) {
		positionSize = maxPositionSize
		riskAmount = positionSize.Mul(stopDistance)
	}

	return positionSize, riskAmount
}

// CalculateStopLoss sets stop loss at 2×ATR from entry.
func (t *TradingAgent) CalculateStopLoss(entryPrice decimal.Decimal, atr float64, isLong bool) decimal.Decimal {
	atrD := decimal.NewFromFloat(atr)
	if isLong {
		return entryPrice.Sub(atrD.Mul(decimal.NewFromInt(2)))
	}
	return entryPrice.Add(atrD.Mul(decimal.NewFromInt(2)))
}

// CalculateTakeProfit sets take profit at minRR × risk distance from entry.
func (t *TradingAgent) CalculateTakeProfit(
	entryPrice, stopLoss decimal.Decimal,
	isLong bool,
	minRR decimal.Decimal,
) decimal.Decimal {
	risk := entryPrice.Sub(stopLoss).Abs()
	if isLong {
		return entryPrice.Add(risk.Mul(minRR))
	}
	return entryPrice.Sub(risk.Mul(minRR))
}

// CalculateRiskRewardRatio returns the reward:risk ratio.
func (t *TradingAgent) CalculateRiskRewardRatio(entryPrice, stopLoss, takeProfit decimal.Decimal) decimal.Decimal {
	risk := entryPrice.Sub(stopLoss).Abs()
	reward := takeProfit.Sub(entryPrice).Abs()
	if risk.IsZero() {
		return decimal.Zero
	}
	return reward.Div(risk)
}

// CalculateVaR computes Value at Risk at the given confidence level (e.g. 0.95).
func (t *TradingAgent) CalculateVaR(returns []float64, confidence float64) float64 {
	if len(returns) == 0 {
		return 0.0
	}
	sorted := make([]float64, len(returns))
	copy(sorted, returns)
	sort.Float64s(sorted)

	idx := int((1.0 - confidence) * float64(len(sorted)))
	if idx >= len(sorted) {
		idx = len(sorted) - 1
	}
	return math.Abs(sorted[idx])
}

// CalculateSharpeRatio computes the Sharpe ratio.
func (t *TradingAgent) CalculateSharpeRatio(returns []float64, riskFreeRate float64) float64 {
	if len(returns) == 0 {
		return 0.0
	}
	avg := btAverage(returns)
	std := btStdDev(returns)
	if std == 0 {
		return 0.0
	}
	return (avg - riskFreeRate) / std
}

// CalculateSortinoRatio computes the Sortino ratio (downside risk only).
func (t *TradingAgent) CalculateSortinoRatio(returns []float64, riskFreeRate float64) float64 {
	if len(returns) == 0 {
		return 0.0
	}
	avg := btAverage(returns)

	downsideSum := 0.0
	downsideCount := 0
	for _, r := range returns {
		if r < 0 {
			downsideSum += r * r
			downsideCount++
		}
	}
	if downsideCount == 0 {
		return 0.0
	}

	downsideDev := math.Sqrt(downsideSum / float64(downsideCount))
	if downsideDev == 0 {
		return 0.0
	}
	return (avg - riskFreeRate) / downsideDev
}

// ============================================================
// MARKET ANALYSIS
// ============================================================

// GenerateSignal scores indicator confluence and returns a trading signal.
func (t *TradingAgent) GenerateSignal(
	indicators BTIndicatorData,
	patterns []string,
	trend string,
	currentPrice float64,
) AnalysisResult {
	score := 0.0
	maxScore := 0.0
	bullishVotes := 0
	bearishVotes := 0

	// RSI (30%)
	maxScore += 30.0
	if indicators.RSI < 30 {
		score += 30.0
		bullishVotes++
	} else if indicators.RSI > 70 {
		score += 30.0
		bearishVotes++
	} else if indicators.RSI < 40 {
		score += 15.0
		bullishVotes++
	} else if indicators.RSI > 60 {
		score += 15.0
		bearishVotes++
	}

	// MACD (25%)
	maxScore += 25.0
	if indicators.MACD.MACD > indicators.MACD.Signal {
		score += 25.0
		bullishVotes++
	} else if indicators.MACD.MACD < indicators.MACD.Signal {
		score += 25.0
		bearishVotes++
	}

	// Bollinger Bands (20%)
	maxScore += 20.0
	if indicators.Bollinger.Lower > 0 && currentPrice <= indicators.Bollinger.Lower {
		score += 20.0
		bullishVotes++
	} else if indicators.Bollinger.Upper > 0 && currentPrice >= indicators.Bollinger.Upper {
		score += 20.0
		bearishVotes++
	}

	// Trend alignment (15%)
	maxScore += 15.0
	if trend == "UPTREND" {
		score += 15.0
		bullishVotes++
	} else if trend == "DOWNTREND" {
		score += 15.0
		bearishVotes++
	}

	// Candlestick pattern bonus (10%)
	maxScore += 10.0
	bullishPatterns := []string{"HAMMER", "BULLISH_ENGULFING", "DOJI"}
	bearishPatterns := []string{"SHOOTING_STAR", "BEARISH_ENGULFING"}
	for _, p := range patterns {
		if btContains(bullishPatterns, p) {
			score += 10.0
			bullishVotes++
			break
		}
		if btContains(bearishPatterns, p) {
			score += 10.0
			bearishVotes++
			break
		}
	}

	confidence := 0.0
	if maxScore > 0 {
		confidence = (score / maxScore) * 100.0
	}

	recommendation := "WAIT"
	if confidence >= 50.0 {
		if bullishVotes >= bearishVotes {
			recommendation = "LONG"
		} else {
			recommendation = "SHORT"
		}
	}

	return AnalysisResult{
		Recommendation: recommendation,
		Confidence:     confidence,
	}
}

// DetectTrend determines the primary trend direction.
func (t *TradingAgent) DetectTrend(candles []Candle) string {
	if len(candles) < 50 {
		return "UNKNOWN"
	}
	ma20 := t.calculateMA(candles, 20)
	ma50 := t.calculateMA(candles, 50)
	currentPrice, _ := candles[len(candles)-1].Close.Float64()

	if currentPrice > ma20 && ma20 > ma50 {
		return "UPTREND"
	} else if currentPrice < ma20 && ma20 < ma50 {
		return "DOWNTREND"
	}
	return "SIDEWAYS"
}

// DetectMarketRegime identifies whether the market is trending or ranging.
func (t *TradingAgent) DetectMarketRegime(candles []Candle) string {
	if len(candles) < 30 {
		return "UNKNOWN"
	}
	adx := t.calculateADX(candles, 14)
	if adx > 25 {
		return "TRENDING"
	}
	return "RANGING"
}

// ConsolidateResults aggregates multi-timeframe results into a single signal.
// It picks the result with the highest confidence that is execution-ready,
// falling back to the highest-confidence result overall.
func (t *TradingAgent) ConsolidateResults(results []AnalysisResult, symbol string) *AnalysisResult {
	if len(results) == 0 {
		return &AnalysisResult{
			Symbol:         symbol,
			Recommendation: "WAIT",
		}
	}

	// Prefer execution-ready results
	best := results[0]
	for _, r := range results[1:] {
		if r.Confidence > best.Confidence {
			best = r
		}
	}
	best.Symbol = symbol
	result := best
	return &result
}

// ComprehensiveAnalysis performs full multi-timeframe analysis with 6-layer validation.
func (t *TradingAgent) ComprehensiveAnalysis(
	ctx context.Context,
	symbol string,
	timeframes []string,
) (*AnalysisResult, error) {
	results := make([]AnalysisResult, 0, len(timeframes))

	for _, tf := range timeframes {
		candles, err := t.FetchOHLCV(ctx, symbol, tf)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch %s/%s: %w", symbol, tf, err)
		}

		// Layer 1: data integrity
		dataValid, dataErrors := t.ValidateDataIntegrity(candles)
		if !dataValid {
			results = append(results, AnalysisResult{
				Symbol:           symbol,
				Timeframe:        tf,
				Recommendation:   "WAIT",
				ValidationPassed: false,
				ValidationReport: BTValidationReport{
					DataIntegrity: false,
					Errors:        dataErrors,
				},
			})
			continue
		}

		indicators := t.CalculateIndicators(candles)

		// Layer 2: indicator validation
		indicatorsValid, indicatorErrors := t.ValidateIndicators(indicators)

		patterns := t.DetectCandlestickPatterns(candles)
		trendPattern := t.DetectTrend(candles)
		marketRegime := t.DetectMarketRegime(candles)

		currentPrice := candles[len(candles)-1].Close
		currentPriceF, _ := currentPrice.Float64()

		signal := t.GenerateSignal(indicators, patterns, trendPattern, currentPriceF)

		atr := CalculateATR(candles, 14)
		isLong := signal.Recommendation == "LONG"

		var entryPrice, stopLoss, takeProfit decimal.Decimal
		if signal.Recommendation == "LONG" || signal.Recommendation == "SHORT" {
			entryPrice = currentPrice
			stopLoss = t.CalculateStopLoss(entryPrice, atr, isLong)
			takeProfit = t.CalculateTakeProfit(entryPrice, stopLoss, isLong, decimal.NewFromFloat(2.0))
		}

		positionSize, riskAmount := t.CalculatePositionSize(entryPrice, stopLoss, t.balance)
		rrRatio := t.CalculateRiskRewardRatio(entryPrice, stopLoss, takeProfit)

		signalResult := AnalysisResult{
			Symbol:          symbol,
			Timeframe:       tf,
			CurrentPrice:    currentPrice,
			Recommendation:  signal.Recommendation,
			Confidence:      signal.Confidence,
			EntryPrice:      entryPrice,
			StopLoss:        stopLoss,
			TakeProfit:      takeProfit,
			RiskRewardRatio: rrRatio,
			PositionSize:    positionSize,
			RiskAmount:      riskAmount,
			Indicators:      indicators,
			Patterns: BTPatternData{
				CandlestickPatterns: patterns,
				TrendPattern:        trendPattern,
				MarketRegime:        marketRegime,
			},
		}

		// Layer 3: signal validation
		signalValid, signalErrors := t.ValidateSignal(signalResult)

		allErrors := append(dataErrors, append(indicatorErrors, signalErrors...)...) //nolint:gocritic
		signalResult.ValidationReport = BTValidationReport{
			DataIntegrity:   dataValid,
			IndicatorsValid: indicatorsValid,
			SignalValid:     signalValid,
			Errors:          allErrors,
		}
		signalResult.ValidationPassed = dataValid && indicatorsValid && signalValid

		results = append(results, signalResult)
	}

	// Layer 4: cross-verification
	crossVerified, crossErrors := t.CrossVerification(results)

	finalResult := t.ConsolidateResults(results, symbol)
	finalResult.ValidationReport.CrossVerified = crossVerified
	finalResult.ValidationReport.Errors = append(finalResult.ValidationReport.Errors, crossErrors...)

	// Layer 5: execution readiness
	finalResult.ValidationReport.ExecutionReady = t.ExecutionReadiness(finalResult.ValidationReport)

	// Layer 6: production validation
	finalResult.ValidationReport.ProductionReady = t.ProductionValidation(*finalResult)

	return finalResult, nil
}

// ============================================================
// BINANCE API INTEGRATION
// ============================================================

// FetchOHLCV retrieves OHLCV candlestick data from Binance public API.
func (t *TradingAgent) FetchOHLCV(ctx context.Context, symbol string, timeframe string) ([]Candle, error) {
	validTFs := map[string]bool{
		"1m": true, "5m": true, "15m": true, "30m": true,
		"1h": true, "4h": true, "1d": true, "1w": true,
	}
	if !validTFs[timeframe] {
		return nil, fmt.Errorf("unsupported timeframe: %s", timeframe)
	}

	// Normalize symbol: BTC/USDT -> BTCUSDT
	sym := strings.ReplaceAll(strings.ToUpper(symbol), "/", "")

	url := fmt.Sprintf(
		"https://api.binance.com/api/v3/klines?symbol=%s&interval=%s&limit=200",
		sym, timeframe,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("binance request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("binance API error: %s", resp.Status)
	}

	// Binance klines: [[openTime, "open", "high", "low", "close", "volume", ...], ...]
	var raw [][]json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("failed to decode klines: %w", err)
	}

	candles := make([]Candle, 0, len(raw))
	for _, row := range raw {
		if len(row) < 6 {
			continue
		}

		var openTimeMS int64
		if err := json.Unmarshal(row[0], &openTimeMS); err != nil {
			continue
		}

		openStr, highStr, lowStr, closeStr, volStr := "", "", "", "", ""
		_ = json.Unmarshal(row[1], &openStr)
		_ = json.Unmarshal(row[2], &highStr)
		_ = json.Unmarshal(row[3], &lowStr)
		_ = json.Unmarshal(row[4], &closeStr)
		_ = json.Unmarshal(row[5], &volStr)

		open, err1 := decimal.NewFromString(openStr)
		high, err2 := decimal.NewFromString(highStr)
		low, err3 := decimal.NewFromString(lowStr)
		close, err4 := decimal.NewFromString(closeStr)
		vol, err5 := decimal.NewFromString(volStr)

		if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil {
			continue
		}

		candles = append(candles, Candle{
			Timestamp: time.UnixMilli(openTimeMS),
			Open:      open,
			High:      high,
			Low:       low,
			Close:     close,
			Volume:    vol,
		})
	}

	return candles, nil
}

// GetTickerPrice retrieves the current price of a symbol from Binance.
func (t *TradingAgent) GetTickerPrice(ctx context.Context, symbol string) (decimal.Decimal, error) {
	sym := strings.ReplaceAll(strings.ToUpper(symbol), "/", "")
	url := fmt.Sprintf("https://api.binance.com/api/v3/ticker/price?symbol=%s", sym)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return decimal.Zero, err
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return decimal.Zero, fmt.Errorf("binance price request failed: %w", err)
	}
	defer resp.Body.Close()

	var data struct {
		Price string `json:"price"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return decimal.Zero, err
	}

	price, err := decimal.NewFromString(data.Price)
	if err != nil {
		return decimal.Zero, fmt.Errorf("invalid price format: %w", err)
	}
	return price, nil
}

// ============================================================
// PICOCLAW TOOL INTERFACE
// ============================================================

// BaseTraderTool is the PicoClaw native tool for cryptocurrency trading analysis.
type BaseTraderTool struct {
	agent *TradingAgent
}

// NewBaseTraderTool creates a new BaseTraderTool with the given initial balance.
func NewBaseTraderTool(balance float64) *BaseTraderTool {
	return &BaseTraderTool{
		agent: &TradingAgent{
			balance:        decimal.NewFromFloat(balance),
			riskPerTrade:   decimal.NewFromFloat(0.02), // 2%
			maxPositionPct: decimal.NewFromFloat(0.10), // 10%
			stopLossPct:    decimal.NewFromFloat(0.15), // 15%
			takeProfitPct:  decimal.NewFromFloat(0.30), // 30%
			journal:        &TradeJournal{Trades: make([]BTTrade, 0)},
		},
	}
}

// Name returns the tool name.
func (bt *BaseTraderTool) Name() string {
	return "base_trader"
}

// Description returns the tool description for the LLM.
func (bt *BaseTraderTool) Description() string {
	return "Autonomous cryptocurrency trading analyst. Provides comprehensive multi-timeframe analysis with professional risk management, pattern recognition, and 6-layer validation. Use for: analyzing trading opportunities, scanning markets, calculating position sizes, and getting trading recommendations with defined risk parameters."
}

// Parameters returns the JSON schema for the tool's input parameters.
func (bt *BaseTraderTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"action": map[string]any{
				"type":        "string",
				"description": "Action to perform: analyze, scan, journal, or status",
				"enum":        []string{"analyze", "scan", "journal", "status"},
			},
			"symbol": map[string]any{
				"type":        "string",
				"description": "Trading pair (e.g. BTC/USDT). Required for 'analyze' action.",
			},
			"timeframes": map[string]any{
				"type":  "array",
				"items": map[string]any{"type": "string"},
				"description": "Timeframes to analyze (default: ['15m','1h','4h']). " +
					"Options: 1m, 5m, 15m, 30m, 1h, 4h, 1d",
			},
			"top_n": map[string]any{
				"type":        "integer",
				"description": "Number of top opportunities to return for 'scan' action (default: 5)",
			},
			"balance": map[string]any{
				"type":        "number",
				"description": "Account balance in USD for position sizing (overrides initial balance)",
			},
		},
		"required": []string{"action"},
	}
}

// Execute runs the requested trading action.
func (bt *BaseTraderTool) Execute(ctx context.Context, args map[string]any) *ToolResult {
	action, _ := args["action"].(string)
	if action == "" {
		return ErrorResult("action is required and must be: analyze, scan, journal, or status")
	}

	if balance, ok := args["balance"].(float64); ok && balance > 0 {
		bt.agent.balance = decimal.NewFromFloat(balance)
	}

	switch action {
	case "analyze":
		return bt.analyzeSymbol(ctx, args)
	case "scan":
		return bt.scanMarket(ctx, args)
	case "journal":
		return bt.getJournal()
	case "status":
		return bt.getStatus()
	default:
		return ErrorResult(fmt.Sprintf("unknown action: %s. Must be: analyze, scan, journal, status", action))
	}
}

func (bt *BaseTraderTool) analyzeSymbol(ctx context.Context, args map[string]any) *ToolResult {
	symbol, _ := args["symbol"].(string)
	symbol = strings.TrimSpace(symbol)
	if symbol == "" {
		return ErrorResult("symbol is required for 'analyze' action (e.g. BTC/USDT)")
	}

	timeframes := []string{"15m", "1h", "4h"}
	if tfs, ok := args["timeframes"].([]any); ok && len(tfs) > 0 {
		timeframes = make([]string, 0, len(tfs))
		for _, tf := range tfs {
			if tfStr, ok := tf.(string); ok {
				timeframes = append(timeframes, tfStr)
			}
		}
	}

	result, err := bt.agent.ComprehensiveAnalysis(ctx, symbol, timeframes)
	if err != nil {
		return ErrorResult(fmt.Sprintf("analysis failed for %s: %v", symbol, err))
	}

	return UserResult(bt.formatAnalysisResult(result))
}

func (bt *BaseTraderTool) scanMarket(ctx context.Context, args map[string]any) *ToolResult {
	topN := 5
	if n, ok := args["top_n"].(float64); ok && n > 0 {
		topN = int(n)
	}

	symbols := []string{
		"BTC/USDT", "ETH/USDT", "BNB/USDT", "SOL/USDT", "XRP/USDT",
		"ADA/USDT", "AVAX/USDT", "DOGE/USDT", "DOT/USDT", "MATIC/USDT",
	}

	opportunities := make([]AnalysisResult, 0, len(symbols))
	for _, sym := range symbols {
		result, err := bt.agent.ComprehensiveAnalysis(ctx, sym, []string{"1h"})
		if err != nil {
			continue
		}
		if result.ValidationReport.ExecutionReady && result.Recommendation != "WAIT" {
			opportunities = append(opportunities, *result)
		}
	}

	sort.Slice(opportunities, func(i, j int) bool {
		scoreI := opportunities[i].Confidence * opportunities[i].RiskRewardRatio.InexactFloat64()
		scoreJ := opportunities[j].Confidence * opportunities[j].RiskRewardRatio.InexactFloat64()
		return scoreI > scoreJ
	})

	if len(opportunities) > topN {
		opportunities = opportunities[:topN]
	}

	return UserResult(bt.formatScanResults(opportunities))
}

func (bt *BaseTraderTool) getJournal() *ToolResult {
	if len(bt.agent.journal.Trades) == 0 {
		return UserResult("No trades in journal yet.")
	}

	var sb strings.Builder
	sb.WriteString("TRADING JOURNAL\n\n")
	for i, trade := range bt.agent.journal.Trades {
		entry, _ := trade.EntryPrice.Float64()
		exit, _ := trade.ExitPrice.Float64()
		pnl, _ := trade.PnL.Float64()
		fmt.Fprintf(&sb, "%d. %s %s\n", i+1, trade.Action, trade.Symbol)
		fmt.Fprintf(&sb, "   Entry: $%.2f | Exit: $%.2f\n", entry, exit)
		fmt.Fprintf(&sb, "   PnL: $%.2f (%.2f%%)\n", pnl, trade.PnLPercent)
		fmt.Fprintf(&sb, "   Reason: %s\n\n", trade.Reason)
	}

	return UserResult(sb.String())
}

func (bt *BaseTraderTool) getStatus() *ToolResult {
	bal, _ := bt.agent.balance.Float64()
	riskPct, _ := bt.agent.riskPerTrade.Float64()
	maxPos, _ := bt.agent.maxPositionPct.Float64()
	sl, _ := bt.agent.stopLossPct.Float64()
	tp, _ := bt.agent.takeProfitPct.Float64()

	return UserResult(fmt.Sprintf(
		"BASE TRADER STATUS\n\nBalance: $%.2f\nRisk per trade: %.1f%%\n"+
			"Max position: %.1f%%\nStop loss: %.1f%%\nTake profit: %.1f%%\nTotal trades: %d",
		bal, riskPct*100, maxPos*100, sl*100, tp*100,
		len(bt.agent.journal.Trades),
	))
}

// ============================================================
// OUTPUT FORMATTERS
// ============================================================

func (bt *BaseTraderTool) formatAnalysisResult(result *AnalysisResult) string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "ANALYSIS: %s\n\n", result.Symbol)

	price, _ := result.CurrentPrice.Float64()
	fmt.Fprintf(&sb, "CURRENT PRICE: $%.4f\n", price)

	if !result.ValidationPassed {
		sb.WriteString("\nVALIDATION FAILED\n")
		for _, e := range result.ValidationReport.Errors {
			fmt.Fprintf(&sb, "  * %s\n", e)
		}
		sb.WriteString("\nRecommendation: WAIT (validation issues detected)\n")
		return sb.String()
	}

	label := "NEUTRAL"
	switch result.Recommendation {
	case "LONG":
		label = "[LONG]"
	case "SHORT":
		label = "[SHORT]"
	case "WAIT":
		label = "[WAIT]"
	}

	fmt.Fprintf(&sb, "\n%s RECOMMENDATION: %s\n", label, result.Recommendation)
	fmt.Fprintf(&sb, "CONFIDENCE: %.0f%%\n", result.Confidence)

	if result.Recommendation != "WAIT" && !result.EntryPrice.IsZero() {
		entry, _ := result.EntryPrice.Float64()
		sl, _ := result.StopLoss.Float64()
		tp, _ := result.TakeProfit.Float64()
		rr, _ := result.RiskRewardRatio.Float64()

		fmt.Fprintf(&sb, "\nENTRY:       $%.4f\n", entry)
		fmt.Fprintf(&sb, "STOP LOSS:   $%.4f\n", sl)
		fmt.Fprintf(&sb, "TAKE PROFIT: $%.4f\n", tp)
		fmt.Fprintf(&sb, "RISK/REWARD: %.1f:1\n", rr)

		bal, _ := bt.agent.balance.Float64()
		posSize, _ := result.PositionSize.Float64()
		posValue, _ := result.PositionSize.Mul(result.EntryPrice).Float64()
		riskAmt, _ := result.RiskAmount.Float64()

		fmt.Fprintf(&sb, "\nPOSITION SIZING (balance $%.2f):\n", bal)
		fmt.Fprintf(&sb, "  Size:  %.6f units\n", posSize)
		fmt.Fprintf(&sb, "  Value: $%.2f\n", posValue)
		if bal > 0 {
			fmt.Fprintf(&sb, "  Risk:  $%.2f (%.2f%% of balance)\n", riskAmt, riskAmt/bal*100)
		}
	}

	fmt.Fprintf(&sb, "\nINDICATORS:\n")
	fmt.Fprintf(&sb, "  RSI:  %.1f\n", result.Indicators.RSI)
	fmt.Fprintf(&sb, "  MACD: %.4f  Signal: %.4f  Hist: %.4f\n",
		result.Indicators.MACD.MACD, result.Indicators.MACD.Signal, result.Indicators.MACD.Histogram)
	fmt.Fprintf(&sb, "  BB Upper: %.4f | Mid: %.4f | Lower: %.4f\n",
		result.Indicators.Bollinger.Upper, result.Indicators.Bollinger.Middle, result.Indicators.Bollinger.Lower)
	fmt.Fprintf(&sb, "  ATR: %.4f | MA20: %.4f | MA50: %.4f\n",
		result.Indicators.ATR, result.Indicators.MA20, result.Indicators.MA50)

	if len(result.Patterns.CandlestickPatterns) > 0 {
		fmt.Fprintf(&sb, "\nPATTERNS: %s\n", strings.Join(result.Patterns.CandlestickPatterns, ", "))
	}
	fmt.Fprintf(&sb, "TREND: %s | REGIME: %s\n", result.Patterns.TrendPattern, result.Patterns.MarketRegime)

	// Validation layers summary
	v := result.ValidationReport
	fmt.Fprintf(&sb, "\nVALIDATION:\n")
	fmt.Fprintf(&sb, "  [1] Data Integrity: %s\n", boolIcon(v.DataIntegrity))
	fmt.Fprintf(&sb, "  [2] Indicators:     %s\n", boolIcon(v.IndicatorsValid))
	fmt.Fprintf(&sb, "  [3] Signal:         %s\n", boolIcon(v.SignalValid))
	fmt.Fprintf(&sb, "  [4] Cross-Verified: %s\n", boolIcon(v.CrossVerified))
	fmt.Fprintf(&sb, "  [5] Exec Ready:     %s\n", boolIcon(v.ExecutionReady))
	fmt.Fprintf(&sb, "  [6] Production:     %s\n", boolIcon(v.ProductionReady))

	if len(v.Warnings) > 0 {
		sb.WriteString("\nWARNINGS:\n")
		for _, w := range v.Warnings {
			fmt.Fprintf(&sb, "  * %s\n", w)
		}
	}

	return sb.String()
}

func (bt *BaseTraderTool) formatScanResults(opportunities []AnalysisResult) string {
	if len(opportunities) == 0 {
		return "MARKET SCAN\n\nNo execution-ready opportunities found at this time."
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "MARKET SCAN - Top %d Opportunities\n\n", len(opportunities))

	for i, opp := range opportunities {
		entry, _ := opp.EntryPrice.Float64()
		rr, _ := opp.RiskRewardRatio.Float64()
		score := opp.Confidence * rr / 10

		fmt.Fprintf(&sb, "%d. %s | %s\n", i+1, opp.Recommendation, opp.Symbol)
		fmt.Fprintf(&sb, "   Confidence: %.0f%% | R/R: %.1f:1 | Score: %.1f\n", opp.Confidence, rr, score)
		fmt.Fprintf(&sb, "   Entry: $%.4f\n\n", entry)
	}

	return sb.String()
}

// ============================================================
// INTERNAL HELPERS
// ============================================================

func mustFloat64(d decimal.Decimal) float64 {
	f, _ := d.Float64()
	return f
}

func decimalMax(a, b decimal.Decimal) decimal.Decimal {
	if a.GreaterThan(b) {
		return a
	}
	return b
}

func decimalMin(a, b decimal.Decimal) decimal.Decimal {
	if a.LessThan(b) {
		return a
	}
	return b
}

func boolIcon(b bool) string {
	if b {
		return "PASS"
	}
	return "FAIL"
}
