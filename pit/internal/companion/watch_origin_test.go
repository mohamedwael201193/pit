package companion

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWatchWebOriginOmitsAccountCapital(t *testing.T) {
	h := &Hub{Dir: t.TempDir()}
	r := httptest.NewRequest(http.MethodGet, "/watch?network=mainnet", nil)
	r.Header.Set("Origin", "https://pit0g.vercel.app")
	w := httptest.NewRecorder()
	h.watch(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("status %d", w.Code)
	}
	var body struct {
		BuyingPower float64 `json:"buyingPower"`
		Sign        bool    `json:"sign"`
		Trade       bool    `json:"trade"`
	}
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body.Sign || body.Trade {
		t.Fatal("sign/trade")
	}
	if body.BuyingPower != 0 {
		t.Fatalf("web origin leaked buying power %v", body.BuyingPower)
	}
}
