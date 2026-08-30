package feasibility

import (
	"strings"
	"testing"

	"github.com/mohamedwael201193/pit/internal/hl"
	"github.com/mohamedwael201193/pit/internal/policy"
)

func ethBook() hl.BookSnapshot {
	return hl.BookSnapshot{Coin: "ETH", MarkPx: 2500, OraclePx: 2510, Funding: 0.0001, OpenInterest: 1e9, DayNtlVlm: 1e8, SzDecimals: 4}
}

func btcBook() hl.BookSnapshot {
	return hl.BookSnapshot{Coin: "BTC", MarkPx: 80000, OraclePx: 80100, Funding: 0.0002, OpenInterest: 2e9, DayNtlVlm: 2e8, SzDecimals: 5}
}

func TestSpotOnlyIsNotExecutable(t *testing.T) {
	p := policy.Default()
	acct := Account{BuyingPower: 0, SpotUSDC: 4.58, SpotFree: 4.58, PowerSource: "spot_not_perp_margin", Note: "Spot is not perp margin."}
	eth := FitBook(ethBook(), p, acct, true, true)
	btc := FitBook(btcBook(), p, acct, true, true)
	if !eth.PolicyEligible || eth.ExecutionFeasible || eth.PreviewReady {
		t.Fatalf("eth %+v", eth)
	}
	if eth.Gate != "insufficient_margin" {
		t.Fatalf("gate %s", eth.Gate)
	}
	if RankGroup(btc) > RankGroup(eth) && btc.ExecutionFeasible {
		t.Fatal("btc must not become executable on signal alone")
	}
	if !strings.Contains(eth.WhyExecutable, "Spot") && !strings.Contains(eth.Why, "minimum") {
		t.Fatalf("%s %s", eth.Why, eth.WhyExecutable)
	}
}

func TestFundedAccountCanPreview(t *testing.T) {
	p := policy.Default()
	acct := Account{BuyingPower: 25, PowerSource: "perp_withdrawable", OpenPositions: 0}
	got := FitBook(ethBook(), p, acct, true, true)
	if !got.PreviewReady || !got.ExecutionFeasible || got.HostNotional < 10 {
		t.Fatalf("%+v", got)
	}
}

func TestKillSwitchBlocksExec(t *testing.T) {
	p := policy.Default()
	p.KillSwitch = true
	acct := Account{BuyingPower: 50, PowerSource: "perp_withdrawable"}
	got := FitBook(ethBook(), p, acct, true, true)
	if got.PolicyEligible || got.ExecutionFeasible {
		t.Fatalf("%+v", got)
	}
	if got.Gate != "kill_switch" {
		t.Fatalf("%s", got.Gate)
	}
}

func TestFitWhyNamesThisMarket(t *testing.T) {
	p := policy.Default()
	acct := Account{BuyingPower: 9.38, PowerSource: "unified_spot"}
	got := FitBook(btcBook(), p, acct, true, true)
	if got.ExecutionFeasible {
		t.Fatal("must not execute")
	}
	if !strings.Contains(got.Why, "BTC") || !strings.Contains(got.Why, "9.38") {
		t.Fatalf("%s", got.Why)
	}
}

func TestMinNotionalBlocksTinyPower(t *testing.T) {
	p := policy.Default()
	acct := Account{BuyingPower: 9.99, PowerSource: "perp_withdrawable"}
	got := FitBook(ethBook(), p, acct, true, true)
	if got.ExecutionFeasible || got.Gate != "insufficient_margin" {
		t.Fatalf("%+v", got)
	}
}

func TestUnpinnedIsExecNotPreview(t *testing.T) {
	p := policy.Default()
	acct := Account{BuyingPower: 40, PowerSource: "perp_withdrawable"}
	got := FitBook(ethBook(), p, acct, true, false)
	if !got.ExecutionFeasible || got.PreviewReady {
		t.Fatalf("%+v", got)
	}
	if got.Gate != "policy_unpinned" {
		t.Fatalf("%s", got.Gate)
	}
}

func TestSessionRequiredForPreview(t *testing.T) {
	p := policy.Default()
	acct := Account{BuyingPower: 40, PowerSource: "perp_withdrawable"}
	got := FitBook(ethBook(), p, acct, false, true)
	if got.PreviewReady || !got.ExecutionFeasible {
		t.Fatalf("%+v", got)
	}
}

func TestAssetOutsidePolicyIsResearchOnly(t *testing.T) {
	p := policy.Default()
	p.AllowedAssets = []string{"ETH"}
	acct := Account{BuyingPower: 50, PowerSource: "perp_withdrawable"}
	got := FitBook(btcBook(), p, acct, true, true)
	if got.PolicyEligible || !got.ResearchEligible {
		t.Fatalf("%+v", got)
	}
}

func TestFitPerAssetMinBTCExceedsETH(t *testing.T) {
	p := policy.Default()
	p.MaxClipUSD = 12
	acct := Account{BuyingPower: 10.20, PowerSource: "unified_spot"}
	btc := FitBook(hl.BookSnapshot{Coin: "BTC", MarkPx: 110000, OraclePx: 110100, SzDecimals: 5}, p, acct, true, true)
	eth := FitBook(ethBook(), p, acct, true, true)
	if btc.ExecutionFeasible {
		t.Fatalf("BTC min ~$11 must not be executable at $10.20: %+v", btc)
	}
	if btc.MinNotionalUSD < 10.99 {
		t.Fatalf("btc min %v", btc.MinNotionalUSD)
	}
	if !eth.ExecutionFeasible || eth.MinNotionalUSD > 10.01 {
		t.Fatalf("eth %+v", eth)
	}
}

func TestLiveEthMarkMakesDefaultClipTight(t *testing.T) {
	p := policy.Default()
	acct := Account{BuyingPower: 16.18, PowerSource: "unified_spot"}
	eth := FitBook(hl.BookSnapshot{Coin: "ETH", MarkPx: 2459, OraclePx: 2460, SzDecimals: 4}, p, acct, true, true)
	if eth.ExecutionFeasible {
		t.Fatalf("16.18 meets venue; $10 clip must not size ETH @ 2459: %+v", eth)
	}
	if eth.Gate != "policy_clip_tight" {
		t.Fatalf("gate %s why %s", eth.Gate, eth.Why)
	}
	if !strings.Contains(eth.Why, "too tight") || strings.Contains(strings.ToLower(eth.Why), "fund") {
		t.Fatal(eth.Why)
	}
	p.MaxClipUSD = 12
	ok := FitBook(hl.BookSnapshot{Coin: "ETH", MarkPx: 2459, OraclePx: 2460, SzDecimals: 4}, p, acct, true, true)
	if !ok.ExecutionFeasible {
		t.Fatalf("clip 12 should size ETH: %+v", ok)
	}
}

func TestFitDefaultClipCannotSizeBTC(t *testing.T) {
	p := policy.Default()
	acct := Account{BuyingPower: 40, PowerSource: "perp_withdrawable"}
	btc := FitBook(btcBook(), p, acct, true, true)
	eth := FitBook(ethBook(), p, acct, true, true)
	if btc.ExecutionFeasible {
		t.Fatalf("$10 policy clip cannot meet BTC rounded min: %+v", btc)
	}
	if !eth.ExecutionFeasible {
		t.Fatalf("eth %+v", eth)
	}
}

func TestRankGroupPrefersExecutable(t *testing.T) {
	if RankGroup(Fit{PreviewReady: true}) <= RankGroup(Fit{PolicyEligible: true}) {
		t.Fatal("preview must outrank signal-only")
	}
}

func TestLiquidityFloorIsResearchOnly(t *testing.T) {
	p := policy.Default()
	p.MinLiquidityUSD = 1e12
	acct := Account{BuyingPower: 50, PowerSource: "perp_withdrawable"}
	got := FitBook(ethBook(), p, acct, true, true)
	if got.PolicyEligible || !got.ResearchEligible {
		t.Fatalf("%+v", got)
	}
	if got.Gate != "liquidity" {
		t.Fatalf("gate %s", got.Gate)
	}
}

func TestOpenCeilingBlocksExec(t *testing.T) {
	p := policy.Default()
	acct := Account{BuyingPower: 50, PowerSource: "perp_withdrawable", OpenPositions: 1}
	got := FitBook(ethBook(), p, acct, true, true)
	if got.ExecutionFeasible || got.Gate != "max_open_positions" {
		t.Fatalf("%+v", got)
	}
}
