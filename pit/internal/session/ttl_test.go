package session

import "testing"

func TestCheckTTL(t *testing.T) {
	if err := CheckTTL(3600, 3600); err != nil {
		t.Fatal(err)
	}
	if err := CheckTTL(7200, 3600); err == nil {
		t.Fatal("long")
	}
	if err := CheckTTL(0, 3600); err == nil {
		t.Fatal("zero")
	}
}
