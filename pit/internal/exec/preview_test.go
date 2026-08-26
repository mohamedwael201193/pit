package exec

import (
	"testing"
	"time"

	"github.com/mohamedwael201193/pit/internal/engine"
)

func TestRequirePreview(t *testing.T) {
	p := engine.Preview{
		Market: "hyperliquid:perp:ETH", Side: "buy", Sz: 0.004, OrderType: "limit",
		LimitPx: "2500", PolicyVersion: "v1", SessionID: "s", WorkspaceID: "w",
		ExpiryUnixMs: time.Now().Add(time.Minute).UnixMilli(),
		Cloid:        "0x11111111111111111111111111111111", ForecastID: "f1", Nonce: 1,
	}
	h, err := engine.CanonicalHash(p)
	if err != nil {
		t.Fatal(err)
	}
	if err := RequirePreview(p, "0x"+h, time.Now().UnixMilli()); err != nil {
		t.Fatal(err)
	}
	p.Sz = 9
	if err := RequirePreview(p, "0x"+h, time.Now().UnixMilli()); err == nil {
		t.Fatal("mutated")
	}
}
