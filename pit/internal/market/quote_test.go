package market

import (
	"testing"
	"time"
)

func TestQuoteRequiresProvenance(t *testing.T) {
	if err := Validate(Quote{}); err == nil {
		t.Fatal("empty")
	}
	q := Hyperliquid("mainnet", "ETH", 2500, 2501, 0.0001, 1e9, time.Now().UTC())
	if err := Validate(q); err != nil {
		t.Fatal(err)
	}
	if q.Source != "hyperliquid" || q.Ref == "" {
		t.Fatal(q)
	}
}
