package watch

import (
	"fmt"

	"github.com/mohamedwael201193/pit/internal/hl"
	"github.com/mohamedwael201193/pit/internal/policy"
)

type BookSource interface {
	PublicBook(coin string) (hl.BookSnapshot, error)
}

type MultiBookSource interface {
	PublicBooks(coins []string) ([]hl.BookSnapshot, error)
}

// Live pulls venue books for allowlisted coins only. It never places orders.
func Live(src BookSource, p policy.Policy) ([]Candidate, error) {
	if src == nil {
		return Scan(nil, p)
	}
	var books []hl.BookSnapshot
	batched := false
	if m, ok := src.(MultiBookSource); ok {
		got, err := m.PublicBooks(p.AllowedAssets)
		if err == nil {
			books = got
			batched = true
		}
	}
	if !batched {
		books = make([]hl.BookSnapshot, 0, len(p.AllowedAssets))
		for _, coin := range p.AllowedAssets {
			b, err := src.PublicBook(coin)
			if err != nil || b.MarkPx <= 0 {
				continue
			}
			books = append(books, b)
		}
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
