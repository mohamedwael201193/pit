package market

import (
	"testing"
	"time"
)

func TestPublicSourcesNeedTimestamp(t *testing.T) {
	now := time.Now().UTC()
	q, err := CoinGecko("ETH", 3512.4, now, "https://api.coingecko.com/api/v3/simple/price")
	if err != nil {
		t.Fatal(err)
	}
	if err := RejectStale(q, now, time.Minute); err != nil {
		t.Fatal(err)
	}
	if err := RejectStale(q, now.Add(10*time.Minute), time.Minute); err == nil {
		t.Fatal("stale")
	}
	if _, err := CoinGecko("ETH", 1, time.Time{}, "x"); err == nil {
		t.Fatal("ts")
	}
}

func TestRejectMockSource(t *testing.T) {
	q := Quote{Source: "mock", AsOf: time.Now().UTC(), MarkPx: 1}
	if err := RejectStale(q, time.Now().UTC(), time.Hour); err == nil {
		t.Fatal("mock")
	}
}
