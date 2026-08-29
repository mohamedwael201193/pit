package auto

import "testing"

func TestHumanWhyNeverEmptyForKnownCodes(t *testing.T) {
	for _, code := range []string{
		"insufficient_margin", "below_min_notional", "policy_denied", "asset_not_allowed",
		"max_open_positions", "cooldown", "daily_loss", "consecutive_loss_limit",
		"liquidity_insufficient", "slippage_too_high", "committee_disagreement",
		"research_stood_down", "TEE_VERIFY_FAIL", "session_expired", "kill_switch",
		"opportunity_expired", "need_pin", "user_stop", "no_opportunity",
	} {
		if HumanWhy(code) == "" || HumanWhy(code) == code {
			t.Fatalf("%s", code)
		}
	}
}

func TestAwayJournalDedupAndCounts(t *testing.T) {
	dir := t.TempDir()
	ResetAway(dir)
	AppendAway(dir, AwayEvent{Kind: "detected", Coin: "BTC", Why: "mark_below_oracle"})
	AppendAway(dir, AwayEvent{Kind: "detected", Coin: "BTC", Why: "mark_below_oracle"})
	AppendAway(dir, AwayEvent{Kind: "rejected", Coin: "BTC", Why: "insufficient_margin"})
	got := LoadAway(dir)
	if got.Detected != 1 || got.Rejected != 1 {
		t.Fatalf("%+v", got)
	}
	if len(got.Events) != 2 {
		t.Fatalf("events %d", len(got.Events))
	}
	if got.Events[1].Human == "" {
		t.Fatal("human")
	}
}
