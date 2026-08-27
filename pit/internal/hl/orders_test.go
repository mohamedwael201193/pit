package hl

import (
	"encoding/json"
	"testing"
)

func TestCloidOnVenue(t *testing.T) {
	raw := json.RawMessage(`[{"coin":"ETH","cloid":"0x11111111111111111111111111111111"}]`)
	if !CloidOnVenue(raw, "0x11111111111111111111111111111111") {
		t.Fatal("miss")
	}
	if CloidOnVenue(raw, "0x22222222222222222222222222222222") {
		t.Fatal("ghost")
	}
	if CloidOnVenue(json.RawMessage(`not-json`), "0x1") {
		t.Fatal("bad")
	}
}
