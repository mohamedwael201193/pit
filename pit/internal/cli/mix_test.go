package cli

import "testing"

func TestRefuseNetworkSwitch(t *testing.T) {
	if err := RefuseNetworkSwitch("mainnet", "mainnet"); err != nil {
		t.Fatal(err)
	}
	if err := RefuseNetworkSwitch("mainnet", "testnet"); err == nil {
		t.Fatal("mix")
	}
}
