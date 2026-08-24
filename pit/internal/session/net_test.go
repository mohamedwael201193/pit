package session

import "testing"

func TestMatchNetwork(t *testing.T) {
	if err := MatchNetwork("mainnet", "mainnet"); err != nil {
		t.Fatal(err)
	}
	if err := MatchNetwork("mainnet", "testnet"); err == nil {
		t.Fatal("mix")
	}
}
