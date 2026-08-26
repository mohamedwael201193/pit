package sdk

import "testing"

func TestClientCannotWatchTrade(t *testing.T) {
	if (Client{}).CanWatchTrade() {
		t.Fatal("trade")
	}
}
