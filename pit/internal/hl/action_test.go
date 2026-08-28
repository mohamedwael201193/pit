package hl

import (
	"strings"
	"testing"
)

func TestOrderCancelOnly(t *testing.T) {
	b, err := BuildOrder(1, true, "2500", "0.004", "0x11111111111111111111111111111111")
	if err != nil {
		t.Fatal(err)
	}
	if err := AssertActionType(b); err != nil {
		t.Fatal(err)
	}
	c, err := BuildCancel(1, "0x11111111111111111111111111111111")
	if err != nil {
		t.Fatal(err)
	}
	if err := AssertActionType(c); err != nil {
		t.Fatal(err)
	}
	if err := AssertActionType([]byte(`{"type":"withdraw3"}`)); err == nil {
		t.Fatal("withdraw")
	}
	r, err := BuildOrderFlags(1, false, "2500", "0.0041", "0x11111111111111111111111111111111", true)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(r), `"r":true`) {
		t.Fatal(string(r))
	}
}
