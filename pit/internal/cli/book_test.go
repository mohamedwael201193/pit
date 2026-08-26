package cli

import "testing"

func TestParseAskFlags(t *testing.T) {
	m, b, err := ParseAskFlags([]string{"--market", "m.json", "--book", "b.json"})
	if err != nil || m != "m.json" || b != "b.json" {
		t.Fatalf("%s %s %v", m, b, err)
	}
	if _, _, err := ParseAskFlags(nil); err == nil {
		t.Fatal("empty")
	}
}
