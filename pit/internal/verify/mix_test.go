package verify

import "testing"

func TestRefuseTestnetReceiptOnMainnet(t *testing.T) {
	if err := RefuseTestnetReceiptOnMainnet("testnet", "mainnet"); err == nil {
		t.Fatal("mix")
	}
	if err := RefuseTestnetReceiptOnMainnet("mainnet", "mainnet"); err != nil {
		t.Fatal(err)
	}
}
