package skills

import (
	"fmt"
	"math"
	"strings"

	"github.com/mohamedwael201193/pit/internal/hl"
)

type Skill struct {
	ID        string `json:"id"`
	Version   string `json:"version"`
	Kind      string `json:"kind"`
	Authority string `json:"authority"`
	Title     string `json:"title"`
}

type Finding struct {
	Skill   Skill  `json:"skill"`
	Fact    string `json:"fact"`
	Present bool   `json:"present"`
}

func Registry() []Skill {
	return []Skill{
		{ID: "mark_oracle_dislocation", Version: "1.0.0", Kind: "fact", Authority: "host", Title: "Mark vs oracle"},
		{ID: "funding", Version: "1.0.0", Kind: "fact", Authority: "host", Title: "Funding"},
		{ID: "open_interest", Version: "1.0.0", Kind: "fact", Authority: "host", Title: "Open interest"},
		{ID: "volume_liquidity", Version: "1.0.0", Kind: "fact", Authority: "host", Title: "Day notional"},
		{ID: "trend_mark_oracle", Version: "1.0.0", Kind: "fact", Authority: "host", Title: "Trend vs oracle"},
		{ID: "volatility_gap", Version: "1.0.0", Kind: "fact", Authority: "host", Title: "Mark/oracle gap"},
		{ID: "candle_structure", Version: "1.0.0", Kind: "fact", Authority: "host", Title: "Multi-timeframe candles"},
		{ID: "risk_reward_host", Version: "1.0.0", Kind: "fact", Authority: "host", Title: "Host clip vs min notional"},
		{ID: "execution_quality", Version: "1.0.0", Kind: "fact", Authority: "host", Title: "Tick / size decimals"},
		{ID: "regime_funding_oi", Version: "1.0.0", Kind: "fact", Authority: "host", Title: "Funding/OI regime"},
		{ID: "market_structure", Version: "1.0.0", Kind: "fact", Authority: "host", Title: "Market structure"},
		{ID: "momentum", Version: "1.0.0", Kind: "fact", Authority: "host", Title: "Momentum"},
		{ID: "mean_reversion", Version: "1.0.0", Kind: "fact", Authority: "host", Title: "Mean reversion"},
		{ID: "liquidation_context", Version: "1.0.0", Kind: "fact", Authority: "host", Title: "Liquidation context"},
		{ID: "support_resistance", Version: "1.0.0", Kind: "fact", Authority: "host", Title: "Support / resistance"},
		{ID: "breakout_rejection", Version: "1.0.0", Kind: "fact", Authority: "host", Title: "Breakout / rejection"},
		{ID: "order_book_context", Version: "1.0.0", Kind: "fact", Authority: "host", Title: "Order book"},
		{ID: "invalidation", Version: "1.0.0", Kind: "fact", Authority: "host", Title: "Invalidation"},
		{ID: "position_sizing", Version: "1.0.0", Kind: "fact", Authority: "host", Title: "Position sizing"},
		{ID: "post_trade_review", Version: "1.0.0", Kind: "fact", Authority: "host", Title: "Post-trade review"},
		{ID: "strategy_performance", Version: "1.0.0", Kind: "fact", Authority: "host", Title: "Strategy performance"},
		{ID: "calibration", Version: "1.0.0", Kind: "fact", Authority: "host", Title: "Calibration"},
	}
}

func ByID(id string) (Skill, bool) {
	for _, s := range Registry() {
		if s.ID == id {
			return s, true
		}
	}
	return Skill{}, false
}

func Apply(b hl.BookSnapshot, candles []Candle) []Finding {
	out := make([]Finding, 0, 22)
	out = append(out,
		markOracle(b), funding(b), openInterest(b), volume(b), trend(b), volGap(b),
		candleStructure(candles), tickQuality(b), regime(b), hostSize(b), positionSizing(b),
		marketStructure(b), momentum(b, candles), meanReversion(b),
		liquidationContext(b), supportResistance(candles), breakout(candles),
		orderBook(), invalidationFact(b), postTrade(), strategyPerf(), calibrationFact(),
	)
	return out
}

func markOracle(b hl.BookSnapshot) Finding {
	s, _ := ByID("mark_oracle_dislocation")
	if b.OraclePx <= 0 {
		return Finding{Skill: s, Present: false, Fact: "Oracle is not published on this book."}
	}
	gap := (b.MarkPx - b.OraclePx) / b.OraclePx
	return Finding{Skill: s, Present: true, Fact: fmt.Sprintf("Mark %.6g vs oracle %.6g (gap %.3f%%).", b.MarkPx, b.OraclePx, gap*100)}
}

func funding(b hl.BookSnapshot) Finding {
	s, _ := ByID("funding")
	return Finding{Skill: s, Present: b.Funding != 0, Fact: fmt.Sprintf("Funding %s.", pct(b.Funding))}
}

func openInterest(b hl.BookSnapshot) Finding {
	s, _ := ByID("open_interest")
	return Finding{Skill: s, Present: b.OpenInterest > 0, Fact: fmt.Sprintf("Open interest %.4g.", b.OpenInterest)}
}

func volume(b hl.BookSnapshot) Finding {
	s, _ := ByID("volume_liquidity")
	return Finding{Skill: s, Present: b.DayNtlVlm > 0, Fact: fmt.Sprintf("Day notional %.4g.", b.DayNtlVlm)}
}

func trend(b hl.BookSnapshot) Finding {
	s, _ := ByID("trend_mark_oracle")
	if b.OraclePx <= 0 {
		return Finding{Skill: s, Present: false, Fact: "Trend vs oracle unavailable."}
	}
	label := "in line with oracle"
	if b.MarkPx < b.OraclePx {
		label = "softer than oracle"
	} else if b.MarkPx > b.OraclePx {
		label = "firmer than oracle"
	}
	return Finding{Skill: s, Present: true, Fact: "Host trend: " + label + "."}
}

func volGap(b hl.BookSnapshot) Finding {
	s, _ := ByID("volatility_gap")
	if b.OraclePx <= 0 {
		return Finding{Skill: s, Present: false, Fact: "Gap unavailable."}
	}
	gap := math.Abs(b.MarkPx-b.OraclePx) / b.OraclePx
	return Finding{Skill: s, Present: gap > 0, Fact: fmt.Sprintf("Absolute mark/oracle gap %.3f%%.", gap*100)}
}

func candleStructure(candles []Candle) Finding {
	s, _ := ByID("candle_structure")
	if len(candles) < 3 {
		return Finding{Skill: s, Present: false, Fact: "Venue candles were not sampled. PIT will not invent support or resistance."}
	}
	last := candles[len(candles)-1]
	prev := candles[len(candles)-2]
	dir := "unchanged"
	if last.Close > prev.Close {
		dir = "last close above prior close"
	} else if last.Close < prev.Close {
		dir = "last close below prior close"
	}
	return Finding{Skill: s, Present: true, Fact: fmt.Sprintf("%d venue candles. %s. Host does not invent levels.", len(candles), dir)}
}

func tickQuality(b hl.BookSnapshot) Finding {
	s, _ := ByID("execution_quality")
	return Finding{Skill: s, Present: true, Fact: fmt.Sprintf("szDecimals %d. Host sizes to this tick. The model cannot.", b.SzDecimals)}
}

func regime(b hl.BookSnapshot) Finding {
	s, _ := ByID("regime_funding_oi")
	bits := []string{}
	if math.Abs(b.Funding) >= 0.0001 {
		bits = append(bits, "elevated funding")
	}
	if b.OpenInterest > 0 && b.OpenInterest < 1000 {
		bits = append(bits, "thin open interest")
	}
	if len(bits) == 0 {
		return Finding{Skill: s, Present: false, Fact: "No extra funding/OI regime flag."}
	}
	return Finding{Skill: s, Present: true, Fact: strings.Join(bits, "; ") + "."}
}

func hostSize(b hl.BookSnapshot) Finding {
	s, _ := ByID("risk_reward_host")
	if b.MarkPx <= 0 {
		return Finding{Skill: s, Present: false, Fact: "No mark. Host will not size."}
	}
	return Finding{Skill: s, Present: true, Fact: "Host sizes the clip. The model cannot set notional, leverage, or side as authority."}
}

func marketStructure(b hl.BookSnapshot) Finding {
	s, _ := ByID("market_structure")
	if b.OraclePx <= 0 {
		return Finding{Skill: s, Present: false, Fact: "Structure vs oracle is unavailable without an oracle print."}
	}
	return Finding{Skill: s, Present: true, Fact: "Live mark/oracle/funding/OI are the structure. PIT does not invent higher-timeframe levels from a single snapshot."}
}

func momentum(b hl.BookSnapshot, candles []Candle) Finding {
	s, _ := ByID("momentum")
	if len(candles) < 3 {
		return Finding{Skill: s, Present: false, Fact: "Venue candles were not sampled. PIT will not invent momentum."}
	}
	return Finding{Skill: s, Present: true, Fact: candleStructure(candles).Fact}
}

func meanReversion(b hl.BookSnapshot) Finding {
	s, _ := ByID("mean_reversion")
	if b.OraclePx <= 0 {
		return Finding{Skill: s, Present: false, Fact: "No oracle. Mean-reversion vs oracle is unavailable."}
	}
	gap := (b.MarkPx - b.OraclePx) / b.OraclePx
	return Finding{Skill: s, Present: math.Abs(gap) > 0.0005, Fact: fmt.Sprintf("Mark/oracle gap %.3f%% is a venue fact, not a mean-reversion signal the model may size.", gap*100)}
}

func liquidationContext(b hl.BookSnapshot) Finding {
	s, _ := ByID("liquidation_context")
	return Finding{Skill: s, Present: false, Fact: "Hyperliquid liquidation heat is not in this snapshot. PIT will not invent liquidations."}
}

func supportResistance(candles []Candle) Finding {
	s, _ := ByID("support_resistance")
	if len(candles) < 5 {
		return Finding{Skill: s, Present: false, Fact: "Not enough venue candles to name support or resistance. PIT will not invent levels."}
	}
	return Finding{Skill: s, Present: true, Fact: fmt.Sprintf("%d venue candles sampled. Host does not name fictional S/R.", len(candles))}
}

func breakout(candles []Candle) Finding {
	s, _ := ByID("breakout_rejection")
	if len(candles) < 5 {
		return Finding{Skill: s, Present: false, Fact: "Breakout/rejection needs sampled candles. None here."}
	}
	return Finding{Skill: s, Present: true, Fact: "Last sampled close vs prior close is recorded. Host does not call a breakout without venue candles."}
}

func orderBook() Finding {
	s, _ := ByID("order_book_context")
	return Finding{Skill: s, Present: false, Fact: "Full-universe L2 is not sampled. Estimated slippage is the host ceiling, not a live book walk."}
}

func invalidationFact(b hl.BookSnapshot) Finding {
	s, _ := ByID("invalidation")
	if b.OraclePx <= 0 {
		return Finding{Skill: s, Present: false, Fact: "Invalidation is a missing live mark or a policy fail."}
	}
	return Finding{Skill: s, Present: true, Fact: "Invalidation: mark/oracle gap closing, funding flipping, or open interest collapsing."}
}

func postTrade() Finding {
	s, _ := ByID("post_trade_review")
	return Finding{Skill: s, Present: false, Fact: "Post-trade review waits for a real OID/fill on this workspace. PIT will not invent an outcome."}
}

func strategyPerf() Finding {
	s, _ := ByID("strategy_performance")
	return Finding{Skill: s, Present: false, Fact: "Skill performance is computed from recorded outcomes later. This pass does not claim the model learned."}
}

func calibrationFact() Finding {
	s, _ := ByID("calibration")
	return Finding{Skill: s, Present: false, Fact: "Calibration stays NOT ENOUGH DATA until this workspace records verified outcomes."}
}

func positionSizing(b hl.BookSnapshot) Finding {
	s, _ := ByID("position_sizing")
	if b.SzDecimals < 0 {
		return Finding{Skill: s, Present: false, Fact: "szDecimals missing. Host will not size."}
	}
	return Finding{Skill: s, Present: true, Fact: fmt.Sprintf("Host sizes to szDecimals %d and the $10 venue floor. The model cannot.", b.SzDecimals)}
}

func pct(v float64) string {
	return fmt.Sprintf("%.4f%%", v*100)
}

func IDs(findings []Finding) []string {
	out := make([]string, 0, len(findings))
	for _, f := range findings {
		if f.Skill.ID != "" {
			out = append(out, f.Skill.ID+"@"+f.Skill.Version)
		}
	}
	return out
}
