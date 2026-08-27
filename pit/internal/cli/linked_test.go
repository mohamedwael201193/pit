package cli

import "testing"

func TestLiveLinkedNeedsNetwork(t *testing.T) {
	if _, err := LiveLinked("not-a-net", "0x1", "abcdefgh-1234-5678-1234-567812345678", "0xabc", 1); err == nil {
		t.Fatal("network")
	}
}
