package redteam

import (
	"testing"
	"time"

	"github.com/mohamedwael201193/pit/internal/engine"
)

func TestPreviewMutationsFailClosed(t *testing.T) {
	host := engine.Preview{
		Market: "hyperliquid:perp:ETH", Side: "buy", Sz: 0.004, OrderType: "limit",
		LimitPx: "1000", PolicyVersion: "v1", SessionID: "s", WorkspaceID: "w",
		ExpiryUnixMs: time.Now().Add(time.Minute).UnixMilli(),
		Cloid:        "0x11111111111111111111111111111111", ForecastID: "f1", Nonce: 1,
	}
	if errs := PreviewMutations(host); len(errs) != 0 {
		t.Fatalf("%v", errs)
	}
}
