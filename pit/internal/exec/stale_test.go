package exec

import (
	"testing"

	"github.com/mohamedwael201193/pit/internal/engine"
)

func TestRefuseStalePreview(t *testing.T) {
	p := engine.Preview{
		Market: "hyperliquid:perp:ETH", Side: "buy", Sz: 0.001,
		ForecastID: "f1", Cloid: "c1", SessionID: "s1", WorkspaceID: "w1",
		ExpiryUnixMs: 100, PolicyVersion: "v1",
	}
	h, err := engine.CanonicalHash(p)
	if err != nil {
		t.Fatal(err)
	}
	if err := RefuseStalePreview(p, h, 200); err == nil {
		t.Fatal("stale")
	}
}
