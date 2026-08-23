package policy

import "testing"

func TestMatchPin(t *testing.T) {
	if err := MatchPin("abc", "abc"); err != nil {
		t.Fatal(err)
	}
	if err := MatchPin("abc", "def"); err == nil {
		t.Fatal("changed")
	}
	if err := MatchPin("", "abc"); err == nil {
		t.Fatal("empty")
	}
}
