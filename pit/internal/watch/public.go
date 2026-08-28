package watch

import (
	"time"

	"github.com/mohamedwael201193/pit/internal/policy"
)

type PublicCoin struct {
	Coin         string  `json:"coin"`
	Reason       string  `json:"reason"`
	Why          string  `json:"why"`
	Trend        string  `json:"trend,omitempty"`
	Rank         int     `json:"rank"`
	Freshness    string  `json:"freshness"`
	Mark         float64 `json:"mark"`
	Oracle       float64 `json:"oracle,omitempty"`
	Funding      float64 `json:"funding,omitempty"`
	OpenInterest float64 `json:"openInterest,omitempty"`
	Timestamp    string  `json:"timestamp"`
	Provenance   string  `json:"provenance"`
	Source       string  `json:"source"`
	Network      string  `json:"network"`
	Eligible     bool    `json:"eligible"`
}

type PublicView struct {
	OK      bool         `json:"ok"`
	Sign    bool         `json:"sign"`
	Trade   bool         `json:"trade"`
	Count   int          `json:"count"`
	Copy    string       `json:"copy"`
	Coins   []PublicCoin `json:"coins"`
	Source  string       `json:"source"`
	Network string       `json:"network"`
}

func Public(cands []Candidate, net string) PublicView {
	if cands == nil {
		cands = []Candidate{}
	}
	coins := make([]PublicCoin, 0, len(cands))
	now := time.Now().UTC().Format(time.RFC3339)
	for _, c := range cands {
		coins = append(coins, PublicCoin{
			Coin:         c.Coin,
			Reason:       c.Reason,
			Why:          WhyHuman(c),
			Trend:        Trend(c.Book),
			Rank:         Rank(c),
			Freshness:    "live",
			Mark:         c.Book.MarkPx,
			Oracle:       c.Book.OraclePx,
			Funding:      c.Book.Funding,
			OpenInterest: c.Book.OpenInterest,
			Timestamp:    now,
			Provenance:   "hyperliquid.info",
			Source:       "hyperliquid",
			Network:      net,
			Eligible:     true,
		})
	}
	return PublicView{
		OK:      true,
		Sign:    false,
		Trade:   false,
		Count:   len(coins),
		Copy:    Attention(len(coins)),
		Coins:   coins,
		Source:  "hyperliquid",
		Network: net,
	}
}

func EmptyPublic(net string) PublicView {
	v := Public(nil, net)
	v.Copy = Attention(0)
	return v
}

func PolicyForWatch() policy.Policy {
	return policy.Default()
}
