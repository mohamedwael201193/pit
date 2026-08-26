package watch

import "testing"

func TestPublicNeverTrades(t *testing.T) {
	v := Public(nil, "mainnet")
	if v.Trade || v.Sign {
		t.Fatal("trade")
	}
	if v.Count != 0 || v.Copy != "No opportunities match your policy." {
		t.Fatalf("%+v", v)
	}
	if v.Coins == nil {
		t.Fatal("coins")
	}
}

func TestEmptyPublic(t *testing.T) {
	v := EmptyPublic("testnet")
	if v.Network != "testnet" || v.Count != 0 || v.Trade {
		t.Fatalf("%+v", v)
	}
}
