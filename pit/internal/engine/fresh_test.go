package engine

import (
	"testing"
	"time"

	"github.com/mohamedwael201193/pit/internal/market"
)

func TestRequireFreshQuote(t *testing.T) {
	now := time.Now().UTC()
	q := market.Hyperliquid("mainnet", "ETH", 3500, 3499, 0, 1, now)
	if err := RequireFreshQuote(q, now); err != nil {
		t.Fatal(err)
	}
	if err := RequireFreshQuote(q, now.Add(time.Minute)); err == nil {
		t.Fatal("stale")
	}
}

func TestRejectWrongMarket(t *testing.T) {
	if err := RejectWrongMarket("hyperliquid:perp:ETH", "hyperliquid:perp:BTC"); err == nil {
		t.Fatal("mkt")
	}
}
