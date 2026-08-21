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
