package sdk

import "testing"

func TestSDKCannotMoveFunds(t *testing.T) {
	c := Client{Network: "mainnet"}
	if c.CanTransfer() || c.CanApproveAgent() {
		t.Fatal("funds")
	}
}
