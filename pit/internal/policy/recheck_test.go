package policy

import "testing"

func TestRecheckDetectsMutation(t *testing.T) {
	a := Default()
	b := Default()
	b.MaxClipUSD = 999
	ctx := Context{RequestedUSD: 10, Coin: "ETH", MarketType: "perp", Venue: "hyperliquid", SessionAlive: true, NowUnix: 1}
	if err := Recheck(a, b, ctx); err == nil {
		t.Fatal("clip")
	}
	if err := Recheck(a, a, ctx); err != nil {
		t.Fatal(err)
	}
}
