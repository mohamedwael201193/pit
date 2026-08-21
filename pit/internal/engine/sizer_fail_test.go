package engine

import (
	"math"
	"testing"
)

func TestSizerRejectsNaNInfMissing(t *testing.T) {
	base := SizerInput{
		MarkPx: 2500, SzDecimals: 4, MaxClipUSD: 15, RequestedUSD: 15,
		Side: "buy", Coin: "ETH", AllowedCoins: []string{"ETH"},
		Venue: "hyperliquid", AllowedVenue: "hyperliquid",
	}
	bad := base
	bad.MarkPx = math.NaN()
	if _, err := SizeOrder(bad); err == nil {
		t.Fatal("nan")
	}
	bad = base
	bad.MarkPx = math.Inf(1)
	if _, err := SizeOrder(bad); err == nil {
		t.Fatal("inf")
	}
	bad = base
	bad.RequestedUSD = -1
	if _, err := SizeOrder(bad); err == nil {
		t.Fatal("neg")
	}
	bad = base
	bad.Coin = "SOL"
	if _, err := SizeOrder(bad); err == nil {
		t.Fatal("sol")
	}
}
