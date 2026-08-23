package mcp

import "testing"

func TestVerifyHintNeverSigns(t *testing.T) {
	h := VerifyHint()
	if h["sign"] != false || h["authorize"] != false {
		t.Fatalf("%+v", h)
	}
}
