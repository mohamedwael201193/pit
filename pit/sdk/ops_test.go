package sdk

import "testing"

func TestSDKCannotTrade(t *testing.T) {
	c := Client{Network: "mainnet"}
	if c.CanHoldSession() || c.WatchMayTrade() {
		t.Fatal("session")
	}
	if c.OpportunityCopy(0) != "No opportunities match your policy." {
		t.Fatal(c.OpportunityCopy(0))
	}
	if len(c.LawCards()) != 14 {
		t.Fatal(len(c.LawCards()))
	}
}
