package hl

import "testing"

func TestDecidePowerSpotNotPerpWhenUnknown(t *testing.T) {
	power, src, note := DecidePower("unknown", 4.58, 0, 0)
	if power != 0 || src != "spot_not_perp_margin" {
		t.Fatalf("%v %s %s", power, src, note)
	}
	if note == "" {
		t.Fatal("need explanation")
	}
}

func TestDecidePowerUnifiedUsesSpot(t *testing.T) {
	power, src, _ := DecidePower("unifiedAccount", 25, 0, 0)
	if power != 25 || src != "unified_spot" {
		t.Fatalf("%v %s", power, src)
	}
}

func TestDecidePowerStandardIgnoresSpot(t *testing.T) {
	power, src, note := DecidePower("disabled", 40, 0, 0)
	if power != 0 || src != "spot_not_perp_margin" {
		t.Fatalf("%v %s %s", power, src, note)
	}
}

func TestSpotHoldParse(t *testing.T) {
	raw := []byte(`{"balances":[{"coin":"USDC","total":"14.6","hold":"2.1"}]}`)
	tot, hold := spotUSDCHold(raw)
	if tot != 14.6 || hold != 2.1 {
		t.Fatalf("%v %v", tot, hold)
	}
}
