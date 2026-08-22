package redteam

import (
	"testing"

	"github.com/mohamedwael201193/pit/internal/policy"
)

func TestSlippageAndLeverageFailClosed(t *testing.T) {
	p := policy.Default()
	ctx := policy.Context{
		RequestedUSD: 10, RequestedLev: 20, Coin: "ETH", MarketType: "perp",
		Venue: "hyperliquid", SessionAlive: true, NowUnix: 1, SlippageBps: 500,
	}
	if err := policy.Check(p, ctx); err == nil {
		t.Fatal("lev or slip")
	}
	ctx.RequestedLev = 1
	ctx.SlippageBps = 500
	if err := policy.Check(p, ctx); err == nil {
		t.Fatal("slip")
	}
	ctx.SlippageBps = 10
	if err := policy.Check(p, ctx); err != nil {
		t.Fatal(err)
	}
}
