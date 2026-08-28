package watch

import (
	"fmt"
	"math"

	"github.com/mohamedwael201193/pit/internal/hl"
)

// WhyHuman explains a host-detected public-book fact. It is not a model score.
func WhyHuman(c Candidate) string {
	if !c.Eligible {
		if c.Block != "" {
			return fmt.Sprintf("%s is live on Hyperliquid but policy blocked it (%s). Side is not decided here.", c.Coin, c.Block)
		}
		return fmt.Sprintf("%s is live on Hyperliquid but outside your policy.", c.Coin)
	}
	switch c.Reason {
	case "mark_below_oracle":
		return fmt.Sprintf("%s mark is below the oracle on Hyperliquid. Policy still allows it. Side is not decided here.", c.Coin)
	case "funding":
		return fmt.Sprintf("%s has nonzero funding on the live book. Policy still allows it. Side is not decided here.", c.Coin)
	default:
		return fmt.Sprintf("%s is in your policy universe with a live Hyperliquid mark.", c.Coin)
	}
}

func RiskFlags(b hl.BookSnapshot) []string {
	var out []string
	if b.OpenInterest > 0 && b.OpenInterest < 1000 {
		out = append(out, "thin_open_interest")
	}
	if math.Abs(b.Funding) >= 0.0001 {
		out = append(out, "elevated_funding")
	}
	if b.OraclePx > 0 {
		gap := math.Abs(b.MarkPx-b.OraclePx) / b.OraclePx
		if gap >= 0.01 {
			out = append(out, "mark_oracle_gap")
		}
	}
	return out
}

func Trend(b hl.BookSnapshot) string {
	if b.OraclePx <= 0 {
		return "mark only"
	}
	if b.MarkPx < b.OraclePx {
		return "softer than oracle"
	}
	if b.MarkPx > b.OraclePx {
		return "firmer than oracle"
	}
	return "in line with oracle"
}

// Rank is a host ranking of venue facts already on the book. It is not AI.
func Rank(c Candidate) int {
	score := 10
	if c.Book.OraclePx > 0 {
		gap := (c.Book.OraclePx - c.Book.MarkPx) / c.Book.OraclePx
		if gap > 0 {
			score += int(math.Min(gap*400, 50))
		}
	}
	if c.Book.Funding != 0 {
		score += 15
	}
	if c.Book.OpenInterest > 0 {
		score += 5
	}
	if score > 100 {
		score = 100
	}
	if score < 1 {
		score = 1
	}
	return score
}
