package watch

import (
	"testing"

	"github.com/mohamedwael201193/pit/internal/feasibility"
	"github.com/mohamedwael201193/pit/internal/hl"
	"github.com/mohamedwael201193/pit/internal/policy"
)

func TestApplyCapitalPrefersExecutableOverSignal(t *testing.T) {
	p := policy.Default()
	cands := []Candidate{
		{Coin: "BTC", Eligible: true, Book: hl.BookSnapshot{Coin: "BTC", MarkPx: 80000, OraclePx: 90000, OpenInterest: 1e9, SzDecimals: 5}, Reason: "mark_below_oracle"},
		{Coin: "ETH", Eligible: true, Book: hl.BookSnapshot{Coin: "ETH", MarkPx: 2500, OraclePx: 2501, OpenInterest: 1, SzDecimals: 4}, Reason: "in_universe"},
	}
	view := Public(cands, "mainnet")
	acct := feasibility.Account{BuyingPower: 0, SpotUSDC: 4.58, PowerSource: "spot_not_perp_margin", Note: "Spot is not perp margin."}
	view = ApplyCapital(view, acct, p, true, true)
	if view.Best == nil {
		t.Fatal("best")
	}
	if view.Best.ExecutionFeasible {
		t.Fatal("must not be executable")
	}
	if view.ExecGate != "insufficient_margin" {
		t.Fatalf("%s", view.ExecGate)
	}
}

func TestBestExecutableIgnoresUnsized(t *testing.T) {
	p := policy.Default()
	cands := []Candidate{
		{Coin: "ETH", Eligible: true, Book: hl.BookSnapshot{Coin: "ETH", MarkPx: 2500, OraclePx: 2510, OpenInterest: 1e9, SzDecimals: 4}},
	}
	acct := feasibility.Account{BuyingPower: 40, PowerSource: "perp_withdrawable"}
	got, ok := BestExecutable(cands, acct, p, true, true)
	if !ok || got.Coin != "ETH" {
		t.Fatalf("%v %+v", ok, got)
	}
}

func TestBestExecutableRejectsBelowMinimum(t *testing.T) {
	p := policy.Default()
	cands := []Candidate{
		{Coin: "ETH", Eligible: true, Book: hl.BookSnapshot{Coin: "ETH", MarkPx: 2500, OraclePx: 2510, OpenInterest: 1e9, SzDecimals: 4}},
	}
	acct := feasibility.Account{BuyingPower: 9.38, PowerSource: "unified_spot"}
	if _, ok := BestExecutable(cands, acct, p, true, true); ok {
		t.Fatal("must not invent size below $10")
	}
}
