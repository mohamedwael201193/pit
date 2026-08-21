package watch

import (
	"fmt"

	"github.com/mohamedwael201193/pit/internal/hl"
	"github.com/mohamedwael201193/pit/internal/policy"
)

type BookSource interface {
	PublicBook(coin string) (hl.BookSnapshot, error)
}

// Live pulls venue books for allowlisted coins only. It never places orders.
func Live(src BookSource, p policy.Policy) ([]Candidate, error) {
	if src == nil {
		return Scan(nil, p)
	}
	books := make([]hl.BookSnapshot, 0, len(p.AllowedAssets))
	for _, coin := range p.AllowedAssets {
		b, err := src.PublicBook(coin)
		if err != nil || b.MarkPx <= 0 {
			continue
		}
		books = append(books, b)
	}
	cands, err := Scan(books, p)
	if err != nil {
		return nil, err
	}
	if err := MayPlaceOrder(true); err == nil {
		return nil, fmt.Errorf("watch_must_not_trade")
	}
	return cands, nil
}
