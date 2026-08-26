package redteam

import (
	"testing"

	"github.com/mohamedwael201193/pit/internal/engine"
)

func TestMutatedPreview(t *testing.T) {
	a := engine.Preview{Market: "hyperliquid:perp:ETH", Side: "buy", Sz: 0.004, OrderType: "limit", LimitPx: "1", PolicyVersion: "v1", SessionID: "s", WorkspaceID: "w", ExpiryUnixMs: 9, Cloid: "0x1", ForecastID: "f", Nonce: 1}
	b := a
	b.Sz = 9
	if err := MutatedPreview(a, b); err != nil {
		t.Fatal(err)
	}
	if err := MutatedPreview(a, a); err == nil {
		t.Fatal("same")
	}
}
