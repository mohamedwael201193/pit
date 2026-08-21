package sdk

import "testing"

func TestWatchCannotTrade(t *testing.T) {
	c := Client{}
	if err := c.WatchCannotTrade(); err == nil {
		t.Fatal("trade")
	}
	if c.WatchCopy(0) == "" {
		t.Fatal("copy")
	}
}
