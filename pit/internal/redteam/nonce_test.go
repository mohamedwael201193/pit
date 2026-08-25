package redteam

import "testing"

func TestNonceReplayDenied(t *testing.T) {
	if !NonceReplayDenied(map[string]struct{}{"n": {}}, "n") {
		t.Fatal("replay")
	}
}
