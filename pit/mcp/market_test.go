package mcp

import (
	"testing"
	"time"

	"github.com/mohamedwael201193/pit/internal/market"
)

func TestMarketQuoteRequiresTimestamp(t *testing.T) {
	r := MarketQuote(market.Quote{Source: "hyperliquid"})
	if r.OK {
		t.Fatal("ts")
	}
	r = MarketQuote(market.Hyperliquid("mainnet", "ETH", 2500, 2501, 0, 1, time.Now().UTC()))
	if !r.OK {
		t.Fatal(r.Error)
	}
}
