package sdk

import "testing"

func TestWatchNeverPlaces(t *testing.T) {
	c := Client{Network: "testnet"}
	if !c.WatchNeverPlaces() || c.WatchMayTrade() {
		t.Fatal("watch")
	}
}
