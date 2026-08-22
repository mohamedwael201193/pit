package exec

import "testing"

func TestCancelBoundRequiresPreview(t *testing.T) {
	e := NewExchange("https://api.hyperliquid.xyz/exchange")
	if _, err := CancelBound(e, 1, "0x11111111111111111111111111111111", "", "abc"); err == nil {
		t.Fatal("empty")
	}
	if _, err := CancelBound(e, 1, "0x11111111111111111111111111111111", "aaa", "bbb"); err == nil {
		t.Fatal("mismatch")
	}
}
