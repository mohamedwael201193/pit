package market

import (
	"time"

	"github.com/mohamedwael201193/pit/internal/hl"
)

func FromBook(network string, b hl.BookSnapshot, asOf time.Time) (Quote, error) {
	q := Hyperliquid(network, b.Coin, b.MarkPx, b.OraclePx, b.Funding, b.OpenInterest, asOf)
	q.Ref = "info/metaAndAssetCtxs"
	return q, Validate(q)
}
