package cli

import "testing"

func TestParsePreviewFlags(t *testing.T) {
	m, s, f, err := ParsePreviewFlags([]string{"--market", "ETH", "--side", "buy", "--forecast", "f1"})
	if err != nil || m != "ETH" || s != "buy" || f != "f1" {
		t.Fatalf("%v %s %s %s", err, m, s, f)
	}
	if _, _, _, err := ParsePreviewFlags([]string{"--market", "ETH"}); err == nil {
		t.Fatal("flags")
	}
}
