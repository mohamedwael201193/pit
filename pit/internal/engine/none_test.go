package engine

import "testing"

func TestRefuseNoneMarket(t *testing.T) {
	if err := RefuseNoneMarket("none"); err == nil {
		t.Fatal("none")
	}
	if err := RefuseNoneMarket("hyperliquid:perp:ETH"); err != nil {
		t.Fatal(err)
	}
}
