package redteam

import (
	"testing"
	"time"

	"github.com/mohamedwael201193/pit/internal/engine"
	"github.com/mohamedwael201193/pit/internal/identity"
)

func TestMarketMutationDenied(t *testing.T) {
	p := engine.Preview{
		Market: "ETH", Side: "buy", Sz: 0.01, OrderType: "limit", LimitPx: "2500",
		SlippageBps: 20, PolicyVersion: "v1", SessionID: "s1",
		ExpiryUnixMs: time.Now().Add(time.Minute).UnixMilli(),
		Cloid: "0x11111111111111111111111111111111", ForecastID: "f1", Nonce: 1,
		WorkspaceID: identity.NewWorkspaceID(),
	}
	h, err := engine.CanonicalHash(p)
	if err != nil {
		t.Fatal(err)
	}
	if err := MarketMutation(p, h, time.Now().UnixMilli()); err != nil {
		t.Fatal(err)
	}
}
