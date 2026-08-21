package hl

import "testing"

func TestMid(t *testing.T) {
	b := L2Book{
		Coin: "ETH",
		Bids: []Level{{Px: 10, Sz: 1}},
		Asks: []Level{{Px: 12, Sz: 1}},
	}
	m, err := Mid(b)
	if err != nil || m != 11 {
		t.Fatalf("%v %v", m, err)
	}
}
