package watch

import (
	"testing"

	"github.com/mohamedwael201193/pit/internal/hl"
)

func TestPublicNeverTrades(t *testing.T) {
	v := Public(nil, "mainnet")
	if v.Trade || v.Sign {
		t.Fatal("trade")
	}
	if v.Count != 0 || v.Copy != "No opportunities match your policy." {
		t.Fatalf("%+v", v)
	}
	if v.Coins == nil {
		t.Fatal("coins")
	}
}

func TestEmptyPublic(t *testing.T) {
	v := EmptyPublic("testnet")
	if v.Network != "testnet" || v.Count != 0 || v.Trade {
		t.Fatalf("%+v", v)
	}
}

func TestPublicCarriesVenueFields(t *testing.T) {
	v := Public([]Candidate{{
		Coin:   "ETH",
		Reason: "funding",
		Book:   hl.BookSnapshot{Coin: "ETH", MarkPx: 2500, OraclePx: 2510, Funding: 0.0001, OpenInterest: 1e9},
	}}, "mainnet")
	if v.Trade || len(v.Coins) != 1 {
		t.Fatalf("%+v", v)
	}
	if v.Coins[0].Oracle != 2510 || v.Coins[0].Funding == 0 || v.Coins[0].OpenInterest == 0 {
		t.Fatalf("%+v", v.Coins[0])
	}
	if v.Coins[0].Why == "" || v.Coins[0].Rank < 1 || v.Coins[0].Freshness != "live" {
		t.Fatalf("%+v", v.Coins[0])
	}
}
