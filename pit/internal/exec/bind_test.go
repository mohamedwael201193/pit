package exec

import "testing"

func TestPostBoundRequiresPreview(t *testing.T) {
	e := NewExchange("https://api.hyperliquid.xyz/exchange")
	raw, err := jsonOrder()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := PostBound(e, raw, "", "abc"); err == nil {
		t.Fatal("empty")
	}
	if _, err := PostBound(e, raw, "aaa", "bbb"); err == nil {
		t.Fatal("mismatch")
	}
}
