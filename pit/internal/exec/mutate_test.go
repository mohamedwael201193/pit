package exec

import "testing"

func TestRejectMutation(t *testing.T) {
	if err := RejectMutation("ETH buy 10", "ETH buy 10"); err != nil {
		t.Fatal(err)
	}
	if err := RejectMutation("ETH buy 10", "ETH buy 99"); err == nil {
		t.Fatal("size")
	}
}
