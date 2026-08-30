package watch

import "testing"

func TestDecideRoutesPolicyClipSelectsWaitNotSwap(t *testing.T) {
	view := PublicView{
		ExecFeasibleN: 0,
		ExecGate:      "policy_clip_tight",
		ExecWhy:       "Policy cap is $0.08 too tight for ETH min $10.08. This account has $16.18. Raise max trade on Security, preview, then pin. PIT will not invent size.",
		Best:          &PublicCoin{Coin: "ETH"},
	}
	routes := DecideRoutes(view)
	if len(routes) != 5 {
		t.Fatalf("want 5 routes, got %d", len(routes))
	}
	got := map[string]CapitalRoute{}
	for _, r := range routes {
		got[r.Action] = r
	}
	if got["trade"].Execution != "blocked" {
		t.Fatalf("trade: %+v", got["trade"])
	}
	if got["wait"].Execution != "selected" {
		t.Fatalf("wait should be selected when clip is tight: %+v", got["wait"])
	}
	if got["swap"].Execution != "unavailable" || got["lp"].Execution != "unavailable" {
		t.Fatalf("swap/lp must stay unavailable: %+v %+v", got["swap"], got["lp"])
	}
	if got["hold"].Execution != "ready" {
		t.Fatalf("hold: %+v", got["hold"])
	}
}

func TestDecideRoutesExecutableSelectsTrade(t *testing.T) {
	view := PublicView{
		ExecFeasibleN: 1,
		Best:          &PublicCoin{Coin: "SOL"},
	}
	routes := DecideRoutes(view)
	for _, r := range routes {
		if r.Action == "trade" && r.Execution != "ready" {
			t.Fatalf("trade: %+v", r)
		}
		if r.Action == "wait" && r.Execution == "selected" {
			t.Fatalf("wait must not be selected when a clip exists: %+v", r)
		}
	}
}
