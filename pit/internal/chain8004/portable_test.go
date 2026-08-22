package chain8004

import "testing"

func TestPortableAcrossNetworks(t *testing.T) {
	if err := PortableAcrossNetworks("mainnet", "testnet"); err == nil {
		t.Fatal("mix")
	}
	if err := PortableAcrossNetworks("mainnet", "mainnet"); err != nil {
		t.Fatal(err)
	}
}
