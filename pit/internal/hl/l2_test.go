package hl

import (
	"encoding/json"
	"testing"
)

func TestParseL2(t *testing.T) {
	raw := json.RawMessage(`{"coin":"ETH","time":1710000000000,"levels":[[["3500.1","1.2"],["3499.0","2"]],[["3500.5","0.8"],["3501.0","3"]]]}`)
	b, err := ParseL2(raw)
	if err != nil {
		t.Fatal(err)
	}
	if b.Coin != "ETH" || b.Bids[0].Px != 3500.1 || b.Asks[0].Sz != 0.8 {
		t.Fatalf("%+v", b)
	}
	bid, ask, err := BestBidAsk(b)
	if err != nil || bid >= ask {
		t.Fatalf("%v %v %v", bid, ask, err)
	}
}

func TestParseL2RejectsEmpty(t *testing.T) {
	_, err := ParseL2(json.RawMessage(`{"coin":"ETH","levels":[[],[]]}`))
	if err == nil {
		t.Fatal("empty")
	}
}

func TestL2CoinRequired(t *testing.T) {
	c := &Client{}
	if _, err := c.L2(""); err == nil {
		t.Fatal("coin")
	}
}
