package cli

import "testing"

func TestLiveOnVenueNeedsNetwork(t *testing.T) {
	if _, err := LiveOnVenue("not-a-net", "0x1", "0x11111111111111111111111111111111"); err == nil {
		t.Fatal("network")
	}
}
