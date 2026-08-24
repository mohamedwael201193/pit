package session

import "testing"

func TestCapTTLHours(t *testing.T) {
	if err := CapTTLHours(2); err == nil {
		t.Fatal("long")
	}
	if err := CapTTLHours(1); err != nil {
		t.Fatal(err)
	}
}
