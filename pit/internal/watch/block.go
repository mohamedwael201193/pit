package watch

import "github.com/mohamedwael201193/pit/internal/policy"

func Blocked(coin string, p policy.Policy) HomeCard {
	ctx := policy.Context{
		RequestedUSD: p.MaxClipUSD,
		RequestedLev: 1,
		Coin:         coin,
		MarketType:   "perp",
		Venue:        "hyperliquid",
		SessionAlive: true,
		NowUnix:      1,
	}
	if err := policy.Check(p, ctx); err != nil {
		return BlockedCard(coin, err.Error())
	}
	return HomeCard{Coin: coin, Eligible: true, Reason: "in_universe"}
}
