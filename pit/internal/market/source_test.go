package market

import "testing"

func TestRequireSource(t *testing.T) {
	if err := RequireSource("hyperliquid", "mainnet", "ETH"); err != nil {
		t.Fatal(err)
	}
	if err := RequireSource("", "mainnet", "ETH"); err == nil {
		t.Fatal("source")
	}
}
