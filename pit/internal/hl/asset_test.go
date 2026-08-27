package hl

import "testing"

func TestIndexInUniverse(t *testing.T) {
	u := []map[string]any{{"name": "BTC"}, {"name": "ETH"}}
	i, err := IndexInUniverse(u, "ETH")
	if err != nil || i != 1 {
		t.Fatalf("%d %v", i, err)
	}
	if _, err := IndexInUniverse(u, "SOL"); err == nil {
		t.Fatal("sol")
	}
}
