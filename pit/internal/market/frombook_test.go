package market

import (
	"testing"
	"time"

	"github.com/mohamedwael201193/pit/internal/hl"
)

func TestFromBook(t *testing.T) {
	q, err := FromBook("mainnet", hl.BookSnapshot{Coin: "ETH", MarkPx: 3501.25, OraclePx: 3500, Funding: 0.0001, OpenInterest: 1.2e8}, time.Now().UTC())
	if err != nil {
		t.Fatal(err)
	}
	if q.Source != "hyperliquid" || q.Symbol != "ETH" {
		t.Fatalf("%+v", q)
	}
}

func TestFromBookRejectsZeroTime(t *testing.T) {
	_, err := FromBook("mainnet", hl.BookSnapshot{Coin: "ETH", MarkPx: 1}, time.Time{})
	if err == nil {
		t.Fatal("ts")
	}
}
