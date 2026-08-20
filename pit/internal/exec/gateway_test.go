package exec

import (
	"testing"
	"time"

	"github.com/mohamedwael201193/pit/internal/engine"
)

func TestGatewayDeniesExtraActions(t *testing.T) {
	p := engine.Preview{
		Market: "hyperliquid:perp:ETH", Side: "buy", Sz: 0.004, OrderType: "limit",
		LimitPx: "1000", PolicyVersion: "v1", SessionID: "s", WorkspaceID: "w",
		ExpiryUnixMs: time.Now().Add(time.Minute).UnixMilli(),
		Cloid:        "0x11111111111111111111111111111111", ForecastID: "f1", Nonce: 1,
	}
	h, _ := engine.CanonicalHash(p)
	used := map[string]struct{}{}
	if err := Prepare(Intent{Action: "order", Preview: p, Hash: h, Workspace: "w"}, time.Now().UnixMilli(), used); err != nil {
		t.Fatal(err)
	}
	for _, a := range []string{"withdraw3", "updateLeverage", "sendAsset", "approveAgent"} {
		if err := Prepare(Intent{Action: a, Preview: p, Hash: h, Workspace: "w"}, time.Now().UnixMilli(), used); err == nil {
			t.Fatalf("%s", a)
		}
	}
	if _, err := RoundPx(2451.123456, 4); err != nil {
		t.Fatal(err)
	}
}
