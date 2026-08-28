package policy

import "fmt"

type Card struct {
	Title string
	Value string
	Law   string
}

func Cards(p Policy) []Card {
	kill := "off"
	if p.KillSwitch {
		kill = "on"
	}
	return []Card{
		{Title: "Max trade", Value: fmt.Sprintf("%.2f USD", p.MaxClipUSD), Law: "The model cannot raise clip."},
		{Title: "Max daily loss", Value: fmt.Sprintf("%.2f USD", p.DailyLossUSD), Law: "A halt is a halt."},
		{Title: "Max leverage", Value: fmt.Sprintf("%dx", p.MaxLeverage), Law: "Session cannot change leverage."},
		{Title: "Allowed assets", Value: join(p.AllowedAssets), Law: "Unknown coins are blocked."},
		{Title: "Allowed venues", Value: join(p.AllowedVenues), Law: "One venue per workspace network."},
		{Title: "Minimum calibration", Value: fmt.Sprintf("%.2f", p.MinSkillCalibration), Law: "Empty health is not a pass."},
		{Title: "Max slippage", Value: fmt.Sprintf("%d bps", p.MaxSlippageBps), Law: "Impact above this fails closed."},
		{Title: "Cooldown", Value: fmt.Sprintf("%d s", p.CooldownSeconds), Law: "Replay waits."},
		{Title: "Session TTL", Value: "24h local mint; Hyperliquid approval until venue date", Law: "Policy still records sessionTtlSeconds 3600. Venue approval can last longer. Expired local sessions cannot sign unless Hyperliquid still lists the PIT agent."},
		{Title: "Kill switch", Value: kill, Law: "YOU flip this. The model cannot."},
		{Title: "Liquidity minimum", Value: fmt.Sprintf("%.2f USD", p.MinLiquidityUSD), Law: "Thin books are skipped."},
		{Title: "Uncertainty maximum", Value: fmt.Sprintf("%.2f", p.MaxUncertainty), Law: "Vague theses do not trade."},
	}
}

func join(xs []string) string {
	if len(xs) == 0 {
		return "(none)"
	}
	out := xs[0]
	for i := 1; i < len(xs); i++ {
		out += ", " + xs[i]
	}
	return out
}
