package ui

import "testing"

func TestYouCopyAndNoSeed(t *testing.T) {
	labs := Labels()
	if len(labs) < 6 {
		t.Fatal(labs)
	}
	if HasSeedField("") {
		t.Fatal("seed")
	}
}
