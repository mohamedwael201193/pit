package watch

import "strings"

// CapitalRoute is a host-deterministic next-capital action. Not a smart-contract router.
// SWAP and LP stay labeled unavailable unless a verified execution rail exists.
type CapitalRoute struct {
	Action    string `json:"action"`
	Coin      string `json:"coin,omitempty"`
	Reason    string `json:"reason"`
	Execution string `json:"execution"`
}

func DecideRoutes(view PublicView) []CapitalRoute {
	why := strings.TrimSpace(view.ExecWhy)
	if why == "" {
		why = "No executable Hyperliquid clip under pinned policy, venue minimum, and this account."
	}
	best := ""
	if view.Best != nil {
		best = view.Best.Coin
	}
	trade := CapitalRoute{
		Action:    "trade",
		Coin:      best,
		Reason:    why,
		Execution: "blocked",
	}
	if view.ExecFeasibleN > 0 {
		trade = CapitalRoute{
			Action:    "trade",
			Coin:      best,
			Reason:    "A host-sized Hyperliquid perp clip fits this account, venue minimum, and pinned policy. Chat cannot AUTHORIZE.",
			Execution: "ready",
		}
	}
	wait := CapitalRoute{
		Action:    "wait",
		Reason:    why,
		Execution: "blocked",
	}
	if view.ExecGate == "policy_clip_tight" {
		wait.Execution = "selected"
		wait.Reason = why
	} else if view.ExecFeasibleN == 0 {
		wait.Execution = "selected"
	}
	if view.ExecFeasibleN > 0 {
		wait = CapitalRoute{
			Action:    "wait",
			Reason:    "A sizeable Hyperliquid clip exists. WAIT is not the selected next action.",
			Execution: "blocked",
		}
	}
	hold := CapitalRoute{
		Action:    "hold",
		Reason:    "No trade is also a valid desk result. PIT will not invent size or a side.",
		Execution: "ready",
	}
	swap := CapitalRoute{
		Action:    "swap",
		Reason:    "PIT does not execute DEX swaps. Hyperliquid perps only. No swap transaction is invented.",
		Execution: "unavailable",
	}
	lp := CapitalRoute{
		Action:    "lp",
		Reason:    "Zia/0G LP is not a PIT execution venue. Incentive APR is live market advertising, not a guaranteed return. Source: app.zia.finance/liquidity. Execution unavailable.",
		Execution: "unavailable",
	}
	return []CapitalRoute{trade, wait, hold, swap, lp}
}
