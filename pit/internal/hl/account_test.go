package hl

import "testing"

func TestSpotFundedWhenPerpZero(t *testing.T) {
	v := ParseAccount(0, 999)
	if v.State != FundedSpot {
		t.Fatalf("got %s", v.State)
	}
	v2 := ParseAccount(0, 0)
	if v2.State != Unfunded {
		t.Fatalf("got %s", v2.State)
	}
	v3 := ParseAccount(12, 0)
	if v3.State != FundedPerp {
		t.Fatalf("got %s", v3.State)
	}
}

func TestSpotUSDCParse(t *testing.T) {
	raw := []byte(`{"balances":[{"coin":"USDC","total":"4.8"},{"coin":"HYPE","total":"1"}]}`)
	if spotUSDCFromClearinghouse(raw) != 4.8 {
		t.Fatal("usdc")
	}
}
