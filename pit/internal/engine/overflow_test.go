package engine

import "testing"

func TestHugeRequestDoesNotExceedClip(t *testing.T) {
	in := SizerInput{
		MarkPx: 2500, SzDecimals: 4, MaxClipUSD: 15, RequestedUSD: 1e18,
		Side: "buy", Coin: "ETH", AllowedCoins: []string{"ETH"},
		Venue: "hyperliquid", AllowedVenue: "hyperliquid",
	}
	o, err := SizeOrderGuarded(in)
	if err != nil {
		return
	}
	if err := ClipNotExceed(o.NotionalUSD, in.MaxClipUSD); err != nil {
		t.Fatal(err)
	}
}
