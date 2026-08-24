package cli

import "testing"

func TestWatchMayPlace(t *testing.T) {
	if err := WatchMayPlace(); err == nil {
		t.Fatal("watch")
	}
	if WatchCopy(0) != "No opportunities match your policy." {
		t.Fatal(WatchCopy(0))
	}
}
