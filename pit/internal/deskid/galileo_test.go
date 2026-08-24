package deskid

import "testing"

func TestTestnetDeskRequired(t *testing.T) {
	if err := TestnetDeskRequired("", false); err != nil {
		t.Fatal(err)
	}
	if err := TestnetDeskRequired("", true); err == nil {
		t.Fatal("mint")
	}
}
