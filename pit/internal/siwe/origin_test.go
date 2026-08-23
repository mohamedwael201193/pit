package siwe

import "testing"

func TestOriginMatch(t *testing.T) {
	if err := OriginMatch("localhost:5173", "localhost:5173"); err != nil {
		t.Fatal(err)
	}
	if err := OriginMatch("localhost:5173", "evil.example"); err == nil {
		t.Fatal("mismatch")
	}
}
