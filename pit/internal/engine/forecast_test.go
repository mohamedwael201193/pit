package engine

import "testing"

func TestMarketAndHostP(t *testing.T) {
	if _, err := ParseMarket("binance:ETH"); err == nil {
		t.Fatal("market")
	}
	f, err := BuildForecast("hyperliquid:perp:ETH", "buy", "mark below 2000", 0.2, 0.4)
	if err != nil || f.P != 0.4 {
		t.Fatal(err)
	}
	if IgnoreModelP(map[string]any{"p": 0.99, "probability": 0.99}, 0.4) != 0.4 {
		t.Fatal("model p leaked")
	}
	if _, err := BuildForecast("none", "none", "", 0.1, 0.1); err == nil {
		t.Fatal("invalidation")
	}
}
