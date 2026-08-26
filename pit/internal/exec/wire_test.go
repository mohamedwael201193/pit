package exec

import (
	"testing"
	"time"

	"github.com/mohamedwael201193/pit/internal/engine"
	"github.com/mohamedwael201193/pit/internal/hl"
)

func TestWireFromPreviewIsOrderOnly(t *testing.T) {
	p := engine.Preview{
		Market: "hyperliquid:perp:ETH", Side: "buy", Sz: 0.004, OrderType: "limit",
		LimitPx: "2500", PolicyVersion: "v1", SessionID: "s", WorkspaceID: "w",
		ExpiryUnixMs: time.Now().Add(time.Minute).UnixMilli(),
		Cloid:        "0x11111111111111111111111111111111", ForecastID: "f1", Nonce: 1,
	}
	raw, err := WireFromPreview(p, 1)
	if err != nil {
		t.Fatal(err)
	}
	if err := hl.AssertActionType(raw); err != nil {
		t.Fatal(err)
	}
	if err := RefuseUnsigned(false); err == nil {
		t.Fatal("unsigned")
	}
	if _, err := CoinFromMarket("spot:ETH"); err == nil {
		t.Fatal("market")
	}
	if err := NeedAssetIndex(false, -1); err == nil {
		t.Fatal("asset")
	}
}
