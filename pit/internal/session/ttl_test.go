package session

import "testing"

func TestCapTTLHours(t *testing.T) {
	if err := CapTTLHours(0); err == nil {
		t.Fatal("zero")
	}
	if err := CapTTLHours(MaxTTLHours + 1); err == nil {
		t.Fatal("over")
	}
	if err := CapTTLHours(1); err != nil {
		t.Fatal(err)
	}
	if err := CapTTLHours(DefaultTTLHours); err != nil {
		t.Fatal(err)
	}
}
