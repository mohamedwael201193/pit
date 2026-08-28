package watch

import (
	"strings"
	"testing"

	"github.com/mohamedwael201193/pit/internal/hl"
)

func TestWhyHumanIsVenueFact(t *testing.T) {
	c := Candidate{Coin: "ETH", Reason: "mark_below_oracle", Book: hl.BookSnapshot{Coin: "ETH", MarkPx: 2400, OraclePx: 2420}}
	s := WhyHuman(c)
	if !strings.Contains(s, "ETH") || strings.Contains(s, "score") {
		t.Fatal(s)
	}
}

func TestRankUsesVenueGap(t *testing.T) {
	low := Rank(Candidate{Coin: "ETH", Reason: "in_universe", Book: hl.BookSnapshot{Coin: "ETH", MarkPx: 2500, OraclePx: 2500}})
	high := Rank(Candidate{Coin: "ETH", Reason: "mark_below_oracle", Book: hl.BookSnapshot{Coin: "ETH", MarkPx: 2400, OraclePx: 2500, Funding: 0.0001}})
	if high <= low {
		t.Fatalf("rank %d <= %d", high, low)
	}
}
