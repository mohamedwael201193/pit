package verify

import "testing"

func TestSameNetwork(t *testing.T) {
	if err := SameNetwork("mainnet", "testnet"); err == nil {
		t.Fatal("mix")
	}
	if err := SameNetwork("mainnet", "mainnet"); err != nil {
		t.Fatal(err)
	}
}
