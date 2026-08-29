package watch

import (
	"time"

	"github.com/mohamedwael201193/pit/internal/policy"
)

type PublicCoin struct {
	Coin         string   `json:"coin"`
	Venue        string   `json:"venue"`
	Reason       string   `json:"reason"`
	Why          string   `json:"why"`
	Trend        string   `json:"trend,omitempty"`
	Rank         int      `json:"rank"`
	Freshness    string   `json:"freshness"`
	Mark         float64  `json:"mark"`
	Oracle       float64  `json:"oracle,omitempty"`
	Funding      float64  `json:"funding,omitempty"`
	OpenInterest float64  `json:"openInterest,omitempty"`
	Volume       float64  `json:"volume,omitempty"`
	Timestamp    string   `json:"timestamp"`
	Provenance   string   `json:"provenance"`
	Source       string   `json:"source"`
	Network      string   `json:"network"`
	Eligible     bool     `json:"eligible"`
	PolicyFit    string   `json:"policyFit"`
	RiskFlags    []string `json:"riskFlags,omitempty"`
	Block        string   `json:"block,omitempty"`
	ExecGate     string   `json:"execGate,omitempty"`
	ExecWhy      string   `json:"execWhy,omitempty"`
}

type PublicView struct {
	OK      bool        `json:"ok"`
	Sign    bool        `json:"sign"`
	Trade   bool        `json:"trade"`
	Count   int         `json:"count"`
	Scanned int         `json:"scanned"`
	Copy    string       `json:"copy"`
	Coins   []PublicCoin `json:"coins"`
	Best    *PublicCoin  `json:"best,omitempty"`
	BestWhy string       `json:"bestWhy,omitempty"`
	ExecGate string      `json:"execGate,omitempty"`
	ExecWhy  string      `json:"execWhy,omitempty"`
	Source  string       `json:"source"`
	Network string      `json:"network"`
}

func toPublic(c Candidate, net, now string) PublicCoin {
	fit := "BLOCKED"
	if c.Eligible {
		fit = "PASS"
	}
	return PublicCoin{
		Coin:         c.Coin,
		Venue:        "hyperliquid",
		Reason:       c.Reason,
		Why:          WhyHuman(c),
		Trend:        Trend(c.Book),
		Rank:         Rank(c),
		Freshness:    "live",
		Mark:         c.Book.MarkPx,
		Oracle:       c.Book.OraclePx,
		Funding:      c.Book.Funding,
		OpenInterest: c.Book.OpenInterest,
		Volume:       c.Book.DayNtlVlm,
		Timestamp:    now,
		Provenance:   "hyperliquid.info metaAndAssetCtxs",
		Source:       "hyperliquid",
		Network:      net,
		Eligible:     c.Eligible,
		PolicyFit:    fit,
		RiskFlags:    c.Risk,
		Block:        c.Block,
	}
}

func Public(cands []Candidate, net string) PublicView {
	if cands == nil {
		cands = []Candidate{}
	}
	now := time.Now().UTC().Format(time.RFC3339)
	coins := make([]PublicCoin, 0, len(cands))
	eligible := 0
	var best *PublicCoin
	for _, c := range cands {
		row := toPublic(c, net, now)
		coins = append(coins, row)
		if c.Eligible {
			eligible++
			if best == nil {
				copyRow := row
				best = &copyRow
			}
		}
	}
	bestWhy := ""
	if best != nil {
		bestWhy = "Highest host rank among policy-eligible live Hyperliquid books. Rank uses mark/oracle gap, funding, and open interest already on the venue. It is not a model score."
	}
	return PublicView{
		OK:      true,
		Sign:    false,
		Trade:   false,
		Count:   eligible,
		Scanned: len(coins),
		Copy:    Attention(eligible),
		Coins:   coins,
		Best:    best,
		BestWhy: bestWhy,
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
