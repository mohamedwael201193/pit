package watch

import (
	"fmt"
	"strings"

	"github.com/mohamedwael201193/pit/internal/hl"
	"github.com/mohamedwael201193/pit/internal/policy"
)

type Candidate struct {
	Coin   string
	Reason string
	Book   hl.BookSnapshot
}

func Scan(books []hl.BookSnapshot, p policy.Policy) ([]Candidate, error) {
	out := []Candidate{}
	for _, b := range books {
		if b.MarkPx <= 0 {
			continue
		}
		ctx := policy.Context{
			RequestedUSD: p.MaxClipUSD,
			RequestedLev: 1,
			Coin:         b.Coin,
			MarketType:   "perp",
			Venue:        "hyperliquid",
			SessionAlive: true,
			NowUnix:      1,
			ImpactUSD:    b.OpenInterest,
		}
		if err := policy.Check(p, ctx); err != nil {
			continue
		}
		reason := ""
		if b.OraclePx > 0 && b.MarkPx < b.OraclePx {
			reason = "mark_below_oracle"
		} else if b.Funding != 0 {
			reason = "funding"
		} else {
			reason = "in_universe"
		}
		out = append(out, Candidate{Coin: strings.ToUpper(b.Coin), Reason: reason, Book: b})
	}
	return out, nil
}

func EmptyCopy(n int) string {
	if n == 0 {
		return "No opportunities match your policy."
	}
	return fmt.Sprintf("%d opportunities match your policy.", n)
}
