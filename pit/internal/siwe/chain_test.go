package siwe

import "testing"

func TestChainMatch(t *testing.T) {
	if err := ChainMatch(16661, 16661); err != nil {
		t.Fatal(err)
	}
	if err := ChainMatch(16602, 16661); err == nil {
		t.Fatal("mix")
	}
}
