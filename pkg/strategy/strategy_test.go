package strategy

import (
	"math"
	"testing"
)

func TestRSIExtremes(t *testing.T) {
	// Monotonic rising series → RSI pinned at 100.
	up := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	if got := RSI(up, 14); got != 100 {
		t.Fatalf("RSI(rising) = %.2f, want 100", got)
	}
	// Insufficient data → neutral 50.
	if got := RSI([]float64{1, 2, 3}, 14); got != 50 {
		t.Fatalf("RSI(short) = %.2f, want 50", got)
	}
}

func TestEvaluateFiresLongOnBullishCross(t *testing.T) {
	// A sustained decline followed by a sharp, sustained rally must eventually
	// produce a fresh bullish EMA cross with price above the fast EMA and RSI not
	// yet overbought — the strategy should go long at least once.
	closes := []float64{}
	p := 100.0
	for i := 0; i < 60; i++ { // decline
		p *= 0.992
		closes = append(closes, p)
	}
	for i := 0; i < 80; i++ { // gradual recovery (not a blow-off)
		p *= 1.006
		closes = append(closes, p)
	}
	highs := make([]float64, len(closes))
	lows := make([]float64, len(closes))
	for i, c := range closes {
		highs[i] = c * 1.01
		lows[i] = c * 0.99
	}

	fired := false
	params := DefaultParams()
	for i := params.EMASlowPeriod + 5; i <= len(closes); i++ {
		sig := Evaluate(closes[:i], highs[:i], lows[:i], params)
		if sig.Direction == "long" {
			fired = true
			if sig.StopLoss <= 0 || sig.TakeProfit <= sig.StopLoss {
				t.Fatalf("long fired with bad SL/TP: sl=%.4f tp=%.4f", sig.StopLoss, sig.TakeProfit)
			}
			if sig.Strength < 0 || sig.Strength > 1 {
				t.Fatalf("strength out of range: %.4f", sig.Strength)
			}
			break
		}
	}
	if !fired {
		t.Fatal("strategy never produced a long on a clear V-recovery (untradeable)")
	}
}

func TestBacktestProducesTradesOnRecoveries(t *testing.T) {
	// Repeated V-recoveries must yield a non-zero number of trades now that the
	// entry logic is tradeable.
	bars := []Bar{}
	p := 100.0
	for cycle := 0; cycle < 15; cycle++ {
		for k := 0; k < 25; k++ {
			p *= 0.99
			bars = append(bars, Bar{Close: p, High: p * 1.01, Low: p * 0.99})
		}
		for k := 0; k < 25; k++ {
			p *= 1.012
			bars = append(bars, Bar{Close: p, High: p * 1.01, Low: p * 0.99})
		}
	}
	res := Backtest(bars, DefaultParams(), 60)
	if res.Trades == 0 {
		t.Fatal("backtest produced zero trades on repeated recoveries")
	}
}

func TestAutoOptimizeSevereLossWidensStop(t *testing.T) {
	// The <0.35 branch must be reachable; the old ordering let <0.45 shadow it.
	p := DefaultParams()
	before := p.StopLossPct
	changed, reason := AutoOptimize(&p, TradeStats{WinRate: 0.30, AvgPnL: -0.1, TradeCount: 20})
	if !changed {
		t.Fatalf("expected change for winRate 0.30")
	}
	if p.StopLossPct <= before {
		t.Fatalf("StopLossPct = %.4f, want > %.4f (widened)", p.StopLossPct, before)
	}
	if reason == "" {
		t.Fatal("expected a non-empty reason")
	}
}

func TestAutoOptimizeThresholdsNeverInvert(t *testing.T) {
	p := DefaultParams()
	// Hammer the mild-loss branch many times; thresholds must stay ordered.
	for i := 0; i < 50; i++ {
		AutoOptimize(&p, TradeStats{WinRate: 0.40, TradeCount: 10})
	}
	if p.RSIOversold >= p.RSIOverbought {
		t.Fatalf("RSI thresholds inverted: oversold=%d overbought=%d", p.RSIOversold, p.RSIOverbought)
	}
	if p.RSIOversold > 45 || p.RSIOverbought < 55 {
		t.Fatalf("RSI thresholds out of clamp: oversold=%d overbought=%d", p.RSIOversold, p.RSIOverbought)
	}
}

func TestAutoOptimizeInsufficientTrades(t *testing.T) {
	p := DefaultParams()
	if changed, _ := AutoOptimize(&p, TradeStats{WinRate: 0.10, TradeCount: 3}); changed {
		t.Fatal("expected no change with fewer than 5 trades")
	}
}

func TestRiskAdjustedSizeRisksFixedFraction(t *testing.T) {
	// Equity 100 SOL, risk 1% (=1 SOL), 10% stop distance → notional 10 SOL,
	// whose 10% adverse move loses exactly the 1 SOL risk budget.
	in := SizingInput{
		EquitySOL:       100,
		RiskPerTradePct: 0.01,
		EntryPrice:      100,
		StopLossPrice:   90,
		Confidence:      1,
	}
	got := RiskAdjustedSize(in)
	if math.Abs(got-10) > 1e-9 {
		t.Fatalf("size = %.6f, want 10", got)
	}
	loss := got * (in.EntryPrice - in.StopLossPrice) / in.EntryPrice
	if math.Abs(loss-1) > 1e-9 {
		t.Fatalf("stop-out loss = %.6f SOL, want 1 (1%% of equity)", loss)
	}
}

func TestRiskAdjustedSizeScalesInverselyWithStop(t *testing.T) {
	base := SizingInput{EquitySOL: 100, RiskPerTradePct: 0.01, EntryPrice: 100, Confidence: 1}
	tight := base
	tight.StopLossPrice = 95 // 5% stop
	wide := base
	wide.StopLossPrice = 80 // 20% stop
	// Wider stop must yield a smaller position for the same risk budget.
	if RiskAdjustedSize(wide) >= RiskAdjustedSize(tight) {
		t.Fatalf("wide-stop size %.4f should be < tight-stop size %.4f",
			RiskAdjustedSize(wide), RiskAdjustedSize(tight))
	}
}

func TestRiskAdjustedSizeConfidenceAndCaps(t *testing.T) {
	in := SizingInput{
		EquitySOL:       100,
		RiskPerTradePct: 0.01,
		EntryPrice:      100,
		StopLossPrice:   90,
		Confidence:      0.5,
	}
	if got := RiskAdjustedSize(in); math.Abs(got-5) > 1e-9 {
		t.Fatalf("confidence-scaled size = %.6f, want 5", got)
	}

	capped := in
	capped.Confidence = 1
	capped.MaxPositionSOL = 3
	if got := RiskAdjustedSize(capped); got != 3 {
		t.Fatalf("MaxPositionSOL cap: size = %.6f, want 3", got)
	}

	pctCap := in
	pctCap.Confidence = 1
	pctCap.MaxPositionPct = 0.04 // 4% of 100 = 4
	if got := RiskAdjustedSize(pctCap); got != 4 {
		t.Fatalf("MaxPositionPct cap: size = %.6f, want 4", got)
	}
}

func TestRiskAdjustedSizeRejectsBadInputs(t *testing.T) {
	cases := []SizingInput{
		{EquitySOL: 0, RiskPerTradePct: 0.01, EntryPrice: 100, StopLossPrice: 90},
		{EquitySOL: 100, RiskPerTradePct: 0, EntryPrice: 100, StopLossPrice: 90},
		{EquitySOL: 100, RiskPerTradePct: 0.01, EntryPrice: 0, StopLossPrice: 90},
		{EquitySOL: 100, RiskPerTradePct: 0.01, EntryPrice: 100, StopLossPrice: 100}, // zero stop distance
	}
	for i, in := range cases {
		if got := RiskAdjustedSize(in); got != 0 {
			t.Fatalf("case %d: expected 0 for bad input, got %.6f", i, got)
		}
	}
}
