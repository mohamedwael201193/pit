package redteam

import (
	"testing"

	"github.com/mohamedwael201193/pit/internal/policy"
	"github.com/mohamedwael201193/pit/internal/session"
)

func TestKillSwitchBlocksSessionAndPolicy(t *testing.T) {
	s := session.Session{ID: "s", Workspace: "w", AgentAddr: "0xa", Expires: 9e18, Kill: true, PolicyVer: "v1"}
	if err := session.CheckSession(s, 1, "v1", "w"); err == nil {
		t.Fatal("kill session")
	}
	p := policy.Default()
	p.KillSwitch = true
	err := policy.Check(p, policy.Context{RequestedUSD: 1, Coin: "ETH", MarketType: "perp", Venue: "hyperliquid", SessionAlive: true, NowUnix: 1})
	if err == nil {
		t.Fatal("kill policy")
	}
}
