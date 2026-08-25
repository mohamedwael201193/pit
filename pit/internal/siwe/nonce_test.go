package siwe

import "testing"

func TestRefuseNonceReplay(t *testing.T) {
	used := map[string]struct{}{"abc": {}}
	if err := RefuseNonceReplay(used, "abc"); err == nil {
		t.Fatal("replay")
	}
	if err := RefuseNonceReplay(used, "def"); err != nil {
		t.Fatal(err)
	}
}
