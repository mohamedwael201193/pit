package compute

import (
	"strings"
	"testing"
)

func TestBuildPrivateBookDoesNotInventFills(t *testing.T) {
	b, err := BuildPrivateBook("0x1111111111111111111111111111111111111111", "11111111-1111-4111-8111-111111111111", "mainnet", "abc")
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)
	if !strings.Contains(s, `"positions":[]`) {
		t.Fatal(s)
	}
	if strings.Contains(strings.ToLower(s), "fill") && strings.Contains(s, "oid") {
		t.Fatal("invented fill")
	}
}
