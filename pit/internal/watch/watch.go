package watch

import (
	"fmt"
	"sort"
	"strings"

	"github.com/mohamedwael201193/pit/internal/hl"
	"github.com/mohamedwael201193/pit/internal/policy"
)

type Candidate struct {
	Coin     string
	Reason   string
	Book     hl.BookSnapshot
	Eligible bool
	Block    string
	Risk     []string
}

func scanContext(p policy.Policy, b hl.BookSnapshot) policy.Context {
	return policy.Context{
		RequestedUSD: p.MaxClipUSD,
		RequestedLev: 1,
		Coin:         b.Coin,
		MarketType:   "perp",
		Venue:        "hyperliquid",
		SessionAlive: true,
		NowUnix:      1,
		ImpactUSD:    b.OpenInterest,
	}
}

func reasonFor(b hl.BookSnapshot) string {
	if b.OraclePx > 0 && b.MarkPx < b.OraclePx {
		return "mark_below_oracle"
	}
	if b.Funding != 0 {
		return "funding"
	}
	return "in_universe"
}

func Universe(books []hl.BookSnapshot, p policy.Policy) []Candidate {
	out := []Candidate{}
	for _, b := range books {
		if b.MarkPx <= 0 || strings.TrimSpace(b.Coin) == "" {
			continue
		}
		c := Candidate{
			Coin:     strings.ToUpper(b.Coin),
			Book:     b,
			Eligible: true,
			Reason:   reasonFor(b),
			Risk:     RiskFlags(b),
		}
		if err := policy.Check(p, scanContext(p, b)); err != nil {
			c.Eligible = false
			c.Block = err.Error()
			c.Reason = "blocked"
		}
		out = append(out, c)
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Eligible != out[j].Eligible {
			return out[i].Eligible
		}
		return Rank(out[i]) > Rank(out[j])
	})
	return out
}

func Scan(books []hl.BookSnapshot, p policy.Policy) ([]Candidate, error) {
	all := Universe(books, p)
	out := make([]Candidate, 0, len(all))
	for _, c := range all {
		if c.Eligible {
			out = append(out, c)
		}
	}
	return out, nil
}

func Best(cands []Candidate) (Candidate, bool) {
	for _, c := range cands {
		if c.Eligible {
			return c, true
		}
	}
	return Candidate{}, false
}

func EmptyCopy(n int) string {
	if n == 0 {
		return "No opportunities match your policy."
	}
	return fmt.Sprintf("%d opportunities match your policy.", n)
}
