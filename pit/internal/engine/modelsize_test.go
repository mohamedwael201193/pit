package engine

import "testing"

func TestIgnoreModelSize(t *testing.T) {
	m := map[string]any{"sizeUsd": 1e12, "side": "buy"}
	if err := IgnoreModelSize(m); err != nil {
		t.Fatal(err)
	}
	if _, ok := m["sizeUsd"]; ok {
		t.Fatal("size still present")
	}
	if m["side"] != "buy" {
		t.Fatal(m)
	}
}

func TestHostStillSizesAfterHugeModelClip(t *testing.T) {
	m := map[string]any{"sizeUsd": 9e18}
	_ = IgnoreModelSize(m)
	sz, err := SizeOrder(SizerInput{
		MarkPx: 2500, SzDecimals: 4, MaxClipUSD: 50, RequestedUSD: 50,
		Side: "buy", Coin: "ETH", AllowedCoins: []string{"ETH"}, MaxLeverage: 1, RequestedLev: 1,
		Venue: "hyperliquid", AllowedVenue: "hyperliquid",
	})
	if err != nil {
		t.Fatal(err)
	}
	if sz.NotionalUSD > 50+1e-6 {
		t.Fatalf("%+v", sz)
	}
}
