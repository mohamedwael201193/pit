package companion

import "testing"

func TestVenueTradeLink(t *testing.T) {
	if got := venueTradeLink("mainnet", "BTC"); got != "https://app.hyperliquid.xyz/trade/BTC" {
		t.Fatal(got)
	}
	if got := venueTradeLink("testnet", "hyperliquid:perp:ETH"); got != "https://app.hyperliquid-testnet.xyz/trade/ETH" {
		t.Fatal(got)
	}
	if got := venueTradeLink("mainnet", ""); got != "https://app.hyperliquid.xyz" {
		t.Fatal(got)
	}
}

func TestCommitteeReasonFromRoles(t *testing.T) {
	got := committeeReason([]map[string]any{
		{"role": "researcher", "verify_e2ee": "OK", "proposed_side": "buy"},
		{"role": "challenger", "verify_e2ee": "OK"},
		{"role": "risk", "verify_e2ee": "OK"},
	})
	if got != "researcher:OK:buy; challenger:OK; risk:OK" {
		t.Fatal(got)
	}
	if committeeReason(nil) != "" {
		t.Fatal("empty")
	}
}
