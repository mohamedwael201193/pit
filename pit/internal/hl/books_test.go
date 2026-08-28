package hl

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPublicBooksOneRoundTrip(t *testing.T) {
	hits := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		hits++
		_, _ = w.Write([]byte(`[{"universe":[{"name":"ETH","szDecimals":4},{"name":"BTC","szDecimals":5}]},[{"markPx":"2500","oraclePx":"2501","funding":"0.0001","openInterest":"1"},{"markPx":"70000","oraclePx":"70010","funding":"0.0002","openInterest":"2"}]]`))
	}))
	t.Cleanup(srv.Close)
	c := &Client{InfoURL: srv.URL, HTTP: srv.Client()}
	books, err := c.PublicBooks([]string{"ETH", "BTC", "NOPE"})
	if err != nil {
		t.Fatal(err)
	}
	if len(books) != 2 || books[0].Coin != "ETH" || books[1].Coin != "BTC" {
		t.Fatalf("%+v", books)
	}
	b, err := c.PublicBook("ETH")
	if err != nil {
		t.Fatal(err)
	}
	if b.MarkPx != 2500 {
		t.Fatalf("%+v", b)
	}
	if hits != 1 {
		t.Fatalf("hits %d want 1", hits)
	}
	uni, err := c.PublicUniverse()
	if err != nil || len(uni) != 2 {
		t.Fatalf("universe %+v %v", uni, err)
	}
}
