package mcp

import (
	"testing"
	"time"

	"github.com/mohamedwael201193/pit/internal/market"
)

func TestOpportunityQuoteRequiresTimestamp(t *testing.T) {
	r := OpportunityQuote(market.Quote{Source: "hyperliquid"})
	if r.OK {
		t.Fatal("timestamp")
	}
	q := market.Hyperliquid("mainnet", "ETH", 1, 1, 0, 0, time.Now().UTC())
	r = OpportunityQuote(q)
	if !r.OK {
		t.Fatal(r)
	}
	body, _ := r.Body.(map[string]any)
	if body["trade"] != false {
		t.Fatal(body)
	}
}
