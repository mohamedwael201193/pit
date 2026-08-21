package hl

import "testing"

func TestParseImpactPxs(t *testing.T) {
	im, err := ParseImpactPxs([]byte(`[2501.5, 2498.2]`))
	if err != nil || im.Buy <= 0 || im.Sell <= 0 {
		t.Fatal(im, err)
	}
	if _, err := ParseImpactPxs([]byte(`[]`)); err == nil {
		t.Fatal("empty")
	}
}
