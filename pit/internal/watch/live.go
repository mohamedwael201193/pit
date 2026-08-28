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

type UniverseSource interface {
	PublicUniverse() ([]hl.BookSnapshot, error)
}

func loadBooks(src BookSource, p policy.Policy) []hl.BookSnapshot {
	if src == nil {
		return nil
	}
	if u, ok := src.(UniverseSource); ok {
		if got, err := u.PublicUniverse(); err == nil {
			return got
		}
	}
	if m, ok := src.(MultiBookSource); ok {
		if got, err := m.PublicBooks(p.AllowedAssets); err == nil {
			return got
		}
	}
	books := make([]hl.BookSnapshot, 0, len(p.AllowedAssets))
	for _, coin := range p.AllowedAssets {
		b, err := src.PublicBook(coin)
		if err != nil || b.MarkPx <= 0 {
			continue
		}
		books = append(books, b)
	}
	return books
}

// Live pulls venue books and returns policy-eligible opportunities only. It never places orders.
func Live(src BookSource, p policy.Policy) ([]Candidate, error) {
	cands, err := Scan(loadBooks(src, p), p)
	if err != nil {
		return nil, err
	}
	if err := MayPlaceOrder(true); err == nil {
		return nil, fmt.Errorf("watch_must_not_trade")
	}
	return cands, nil
}

// LiveUniverse returns every live book PIT can read, with an honest policy-fit flag.
func LiveUniverse(src BookSource, p policy.Policy) ([]Candidate, error) {
	all := Universe(loadBooks(src, p), p)
	if err := MayPlaceOrder(true); err == nil {
		return nil, fmt.Errorf("watch_must_not_trade")
	}
	return all, nil
}
