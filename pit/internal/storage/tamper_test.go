package storage

import "testing"

func TestTamperAndWrongRoot(t *testing.T) {
	h := PayloadHash([]byte("desk"))
	if err := RootsMatch(h, h); err != nil {
		t.Fatal(err)
	}
	if err := RootsMatch(h, PayloadHash([]byte("other"))); err == nil {
		t.Fatal("tamper")
	}
	if err := ComparePlain([]byte("a"), []byte("b")); err == nil {
		t.Fatal("plain")
	}
	if err := ComparePlain([]byte("a"), []byte("a")); err != nil {
		t.Fatal(err)
	}
}
