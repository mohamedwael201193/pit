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
	if !strings.Contains(s, `"hypothesis":"none"`) {
		t.Fatal(s)
	}
	if !strings.Contains(s, "live market facts") {
		t.Fatal("none book must ask researcher to read live facts")
	}
}

func TestBuildPrivateBookHypothesisLong(t *testing.T) {
	b, err := BuildPrivateBookHypothesis("0x1111111111111111111111111111111111111111", "11111111-1111-4111-8111-111111111111", "mainnet", "abc", "long")
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)
	if !strings.Contains(s, `"hypothesis":"long"`) {
		t.Fatal(s)
	}
	if strings.Contains(s, `"oid"`) {
		t.Fatal("invented oid")
	}
}

func TestParseHypothesis(t *testing.T) {
	got, err := ParseHypothesis("")
	if err != nil || got != "none" {
		t.Fatal(got, err)
	}
	got, err = ParseHypothesis("LONG")
	if err != nil || got != "long" {
		t.Fatal(got, err)
	}
	if _, err := ParseHypothesis("withdraw"); err == nil {
		t.Fatal("withdraw")
	}
}
