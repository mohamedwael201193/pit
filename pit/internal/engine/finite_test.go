package engine

import (
	"math"
	"testing"
)

func TestRejectNonFiniteRequestedSize(t *testing.T) {
	in := SizerInput{
		MarkPx: 2500, SzDecimals: 4, MaxClipUSD: 50, RequestedUSD: math.NaN(),
		Side: "buy", Coin: "ETH", AllowedCoins: []string{"ETH"},
		Venue: "hyperliquid", AllowedVenue: "hyperliquid",
	}
	if _, err := SizeOrderGuarded(in); err == nil {
		t.Fatal("nan size")
	}
	in.RequestedUSD = math.Inf(1)
	if _, err := SizeOrderGuarded(in); err == nil {
		t.Fatal("inf size")
	}
	in.RequestedUSD = 50
	if _, err := SizeOrderGuarded(in); err != nil {
		t.Fatal(err)
	}
}
