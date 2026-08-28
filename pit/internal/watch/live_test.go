package watch

import (
	"fmt"
	"testing"

	"github.com/mohamedwael201193/pit/internal/hl"
	"github.com/mohamedwael201193/pit/internal/policy"
)

type stubBooks struct{}

func (stubBooks) PublicBook(coin string) (hl.BookSnapshot, error) {
	if coin == "ETH" {
		return hl.BookSnapshot{Coin: "ETH", MarkPx: 2500, OraclePx: 2510, OpenInterest: 1e9}, nil
	}
	return hl.BookSnapshot{}, fmt.Errorf("no_book")
}

type batchStub struct{ n int }

func (s *batchStub) PublicBook(string) (hl.BookSnapshot, error) {
	s.n++
	return hl.BookSnapshot{}, fmt.Errorf("should_not_per_coin")
}

func (s *batchStub) PublicBooks([]string) ([]hl.BookSnapshot, error) {
	return []hl.BookSnapshot{{Coin: "ETH", MarkPx: 2500, OraclePx: 2510, OpenInterest: 1e9}}, nil
}

func TestLiveRespectsPolicyAndDoesNotTrade(t *testing.T) {
	cands, err := Live(stubBooks{}, policy.Default())
	if err != nil {
		t.Fatal(err)
	}
	if len(cands) != 1 || cands[0].Coin != "ETH" {
		t.Fatalf("%+v", cands)
	}
	if err := MayPlaceOrder(true); err == nil {
		t.Fatal("trade")
	}
}

func TestLiveUsesPublicBooksWhenAvailable(t *testing.T) {
	s := &batchStub{}
	cands, err := Live(s, policy.Default())
	if err != nil {
		t.Fatal(err)
	}
	if s.n != 0 {
		t.Fatalf("per-coin fetches %d", s.n)
	}
	if len(cands) != 1 || cands[0].Coin != "ETH" {
		t.Fatalf("%+v", cands)
	}
}
