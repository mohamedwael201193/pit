package mcp

import "testing"

func TestWatchNeverTrades(t *testing.T) {
	if !WatchNeverTrades() {
		t.Fatal("watch")
	}
}
